<script setup lang="ts">
// 获取当前 Vue 实例的 proxy 对象
import {useWebStore} from "@/store/webStore";
import {useProxiesStore} from "@/store/proxiesStore";
import {useMenuStore} from "@/store/menuStore";
import createApi from "@/api";
import {Events} from "@wailsio/runtime";
import {useI18n} from "vue-i18n";
import {pError, pLoad, pSuccess} from "@/util/pLoad";

// i18n
const {t} = useI18n();

// 获取当前 Vue 实例的 proxy 对象
const {proxy} = getCurrentInstance()!;
const api = createApi(proxy);

// 当前页面使用store
const menuStore = useMenuStore();
const proxiesStore = useProxiesStore();
const webStore = useWebStore();

// 模式切换
Events.On("switchMode", (ev: any) => {
  menuStore.rule = ev.data[0];
});
watch(
    () => menuStore.rule,
    (newVal) => {
      Events.Emit({name: "mode", data: newVal});
    }
);


// 配置切换
Events.On("switchProfiles", async (ev: any) => {
  const data = ev.data[0];

  await pLoad(t('profiles.switch.ing'), async () => {
    try {
      await api.switchProfile(data)
      proxiesStore.active = ""

      await api.waitRunning()

      api.getRules().then((res) => {
        menuStore.setRuleNum(res.length);
      });

      await Events.Emit({name: "hasSwitchProfile", data});

      pSuccess(t('profiles.switch.success'))
    } catch (e) {
      if (e['message']) {
        pError(e['message'])
      }
      api.getProfileList().then((list) => {
        if (list && list.length != 0) {
          Events.Emit({
            name: "profiles",
            data: list
          })
        }
      })
    }
  })
});

onMounted(() => {
  // 发送模式数据
  Events.Emit({name: "mode", data: menuStore.rule});
  // 发送订阅配置数据
  api.getProfileList().then((list) => {
    if (list && list.length != 0) {
      Events.Emit({
        name: "profiles",
        data: list
      })
    }
  })
  // 发送系统代理数据
  Events.Emit({
    name: "proxy",
    data: menuStore.proxy
  })
  // 发送虚拟网卡数据
  Events.Emit({
    name: "tun",
    data: menuStore.tun
  })
})

</script>

<template></template>
