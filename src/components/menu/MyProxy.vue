<script setup lang="ts">
import {useMenuStore} from "@/store/menuStore";
import {useI18n} from "vue-i18n";
import createApi from "@/api";
import {Events} from "@/runtime";
import {pError, pLoad, pSuccess, pWarning} from "@/util/pLoad";
import {useSettingStore} from "@/store/settingStore";
import {pUpdateMihomo} from "@/util/mihomo";
import {useHomeStore} from "@/store/homeStore";

// 使用store
const menuStore = useMenuStore();
const settingStore = useSettingStore();
const homeStore = useHomeStore();

// 获取当前 Vue 实例的 proxy 对象
const {proxy} = getCurrentInstance()!;
const api = createApi(proxy);

// 国际化
const {t} = useI18n();

// 页面使用参数
const tunOn = ref(false)


async function selected() {
  const list = await api.getProfileList()
  if (!list || list.length == 0) {
    pWarning(t("no-profile-warning"));
    return false
  }

  for (let profile of list) {
    if (profile['selected']) {
      return true
    }
  }

  pWarning(t("select-profile-warning"));
  return false
}


// 代理开关
async function doSwitch() {
  let ok = false

  // 检测通过执行后续操作
  if (!menuStore.proxy) {
    try {
      // 添加配置后执行
      const select = await selected()
      if (!select) {
        return
      }
      // 检测端口是否被占用
      await api.checkAddressPort({
        "bindAddress": settingStore.bindAddress,
        "port": settingStore.port,
      })
      // 未被占用开启代理
      await api.enableProxy({
        "bindAddress": settingStore.bindAddress,
        "port": settingStore.port,
      })
      ok = true
      pSuccess(t("proxy-switch-on"));
    } catch (e) {
      if (e['message']) {
        pError(e['message'])
      }
    }
  } else {
    await api.disableProxy()
    ok = true
    pWarning(t("proxy-switch-off"));
  }

  // 同步配置
  if (ok) {
    menuStore.setProxy(!menuStore.proxy);
    pUpdateMihomo(menuStore, settingStore, api)
  }

  // 发送事件通知
  Events.Emit({name: "proxy", data: menuStore.proxy});
}

const proxySwitch = async () => {
  await pLoad(t('switch.ing'), doSwitch)
}

Events.On("switchProxy", async () => {
  await proxySwitch()
});


// 虚拟网卡开关
async function tunSwitch() {
  // 检测是否运行在管理员模式下
  const admin = await api.getAdmin();
  if (!admin.data) {
    pWarning(t("tun-warning"));
    return
  }

  // 添加配置后执行
  const select = await selected()
  if (!select) {
    return
  }

  // 检测通过执行后续操作
  menuStore.setTun(!tunOn.value);
  if (menuStore.tun) {
    api.updateConfigs({
      tun: {
        enable: true,
        stack: settingStore.stack,
      },
    }).then(() => {
      tunOn.value = true;
      pSuccess(t("tun-switch-on"));

      // 同步 mihomo 配置
      pUpdateMihomo(menuStore, settingStore, api)

      // 发送事件通知
      Events.Emit({name: "tun", data: menuStore.tun});
    });
  } else {
    api.updateConfigs({
      tun: {
        enable: false,
      },
    }).then(() => {
      tunOn.value = false;
      pWarning(t("tun-switch-off"));

      // 同步 mihomo 配置
      pUpdateMihomo(menuStore, settingStore, api)

      // 发送事件通知
      Events.Emit({name: "tun", data: menuStore.tun});
    });
  }
}

Events.On("switchTun", async () => {
  await tunSwitch()
});


onMounted(async () => {
  if (homeStore.os != "Windows" && menuStore.tun) {
    await api.waitRunning()
    await tunSwitch()
  }
})


</script>

<template>
  <div class="sub">
    <div class="switch-container">
      <span class="switch-label">
        {{ $t("proxy-switch") }}
      </span>
      <div
          :class="['switch', { 'switch-on': menuStore.proxy }]"
          @click="proxySwitch"
      >
        <div class="switch-circle"></div>
      </div>
    </div>
    <div class="switch-container">
      <span class="switch-label">
        {{ $t("tun-switch") }}
      </span>
      <div :class="['switch', { 'switch-on': tunOn }]" @click="tunSwitch">
        <div class="switch-circle"></div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.sub {
  width: 184px; /* 卡片宽度 */
  background-color: rgba(255, 255, 255, 0.1); /* 半透明白色背景 */
  border: 1px solid #ccc;
  border-radius: 8px;
  box-shadow: var(--left-nav-shadow);
  line-height: 1.5;
  margin-left: 22px;
  margin-top: 25px;
  padding-bottom: 12px;
}

.sub:hover {
  box-shadow: var(--left-nav-hover-shadow);
}

.switch-container {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 16px;
  padding-left: 10px;
  margin-top: 10px;
}

.switch-label {
  color: var(--text-color);
}

.switch {
  width: 54px;
  height: 26px;
  border: 2px solid var(--text-color);
  border-radius: 15px;
  position: relative;
  background-color: transparent;
  cursor: pointer;
  transition: background-color 0.3s ease;
}

.switch.switch-on {
  background-color: var(--left-item-selected-bg);
}

.switch-circle {
  width: 20px;
  height: 20px;
  background-color: var(--text-color);
  border-radius: 50%;
  position: absolute;
  top: 3px;
  left: 3px;
  transition: left 0.3s ease;
}

.switch-on .switch-circle {
  left: 31px;
}
</style>
