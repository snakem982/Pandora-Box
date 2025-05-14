<script setup lang="ts">
// 获取当前 Vue 实例的 proxy 对象
import {useProxiesStore} from "@/store/proxiesStore";
import {useMenuStore} from "@/store/menuStore";
import createApi from "@/api";
import {Events} from "@/runtime";
import {useI18n} from "vue-i18n";
import {pError, pLoad, pSuccess} from "@/util/pLoad";
import {useSettingStore} from "@/store/settingStore";
import {useWebStore} from "@/store/webStore";

// i18n
const {t} = useI18n();

// 获取当前 Vue 实例的 proxy 对象
const {proxy} = getCurrentInstance()!;
const api = createApi(proxy);

// 当前页面使用store
const menuStore = useMenuStore();
const proxiesStore = useProxiesStore();
const settingStore = useSettingStore();
const webStore = useWebStore();

// 模式切换
Events.On("switchMode", (ev: any) => {
  console.log("switchMode=====", ev);
  menuStore.rule = ev;
});
watch(
    () => menuStore.rule,
    (newVal) => {
      Events.Emit({name: "mode", data: newVal});
    }
);


// 配置切换
Events.On("switchProfiles", async (ev: any) => {
  const data = ev;

  await pLoad(t('profiles.switch.ing'), async () => {
    try {
      await api.switchProfile(data)
      proxiesStore.active = ""

      await api.waitRunning()

      api.getRuleNum().then((res) => {
        menuStore.setRuleNum(res);
      });

      webStore.fProfile = data

      api.getProfileList().then((list) => {
        if (list && list.length != 0) {
          Events.Emit({
            name: "profiles",
            data: list
          })
        }
      })

      pSuccess(t('profiles.switch.success'))
    } catch (e) {
      if (e['message']) {
        pError(e['message'])
      }
    }
  })
});

onMounted(async () => {
  // 获取初始数据
  const res = await api.getMihomo()
  menuStore.setRule(res.mode)
  menuStore.setProxy(res.proxy)

  settingStore.setPort(res.port)
  settingStore.setBindAddress(res.bindAddress)
  settingStore.setStack(res.stack)
  settingStore.setDns(res.dns)
  settingStore.setIpv6(res.ipv6)


  // 发送代理模式数据
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
})

</script>

<template></template>
