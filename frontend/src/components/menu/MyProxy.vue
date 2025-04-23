<script setup lang="ts">
import {useMenuStore} from "@/store/menuStore";
import {useI18n} from "vue-i18n";
import createApi from "@/api";
import {Events} from "@wailsio/runtime";
import {pError, pSuccess, pWarning} from "@/util/pLoad";
import {useSettingStore} from "@/store/settingStore";
import {pUpdateMihomo} from "@/util/mihomo";

// 使用store
const menuStore = useMenuStore();
const settingStore = useSettingStore();

// 获取当前 Vue 实例的 proxy 对象
const {proxy} = getCurrentInstance()!;
const api = createApi(proxy);

// 国际化
const {t} = useI18n();

// 代理开关
async function proxySwitch() {
  let ok = false

  // 检测通过执行后续操作
  if (!menuStore.proxy) {
    try {
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
}

Events.On("switchProxy", async (ev: any) => {
  await proxySwitch()
  // 发送事件通知
  Events.Emit({name: "proxy", data: menuStore.proxy});
});


// 虚拟网卡开关
function tunSwitch() {
  // 检测是否运行在管理员模式下


  // 检测通过执行后续操作
  menuStore.setTun(!menuStore.tun);
  if (menuStore.tun) {
    api.updateConfigs({
      tun: {
        enable: true,
      },
    }).then((res: any) => {
      pSuccess(t("tun-switch-on"));

      // 同步 mihomo 配置
      pUpdateMihomo(menuStore, settingStore, api)
    });
  } else {
    api.updateConfigs({
      tun: {
        enable: false,
      },
    }).then((res: any) => {
      pWarning(t("tun-switch-off"));

      // 同步 mihomo 配置
      pUpdateMihomo(menuStore, settingStore, api)
    });
  }
}


onMounted(() => {
  api.getMihomo().then((res) => {
    menuStore.setRule(res.mode)
    menuStore.setProxy(res.proxy)
    menuStore.setTun(res.tun)

    settingStore.setPort(res.port)
    settingStore.setBindAddress(res.bindAddress)
    settingStore.setStack(res.stack)
    settingStore.setDns(res.dns)
    settingStore.setIpv6(res.ipv6)
  })
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
      <div :class="['switch', { 'switch-on': menuStore.tun }]" @click="tunSwitch">
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
