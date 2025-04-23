<script setup lang="ts">
// 获取当前 Vue 实例的 proxy 对象
import { useWebStore } from "@/store/webStore";
import { useProxiesStore } from "@/store/proxiesStore";
import { useMenuStore } from "@/store/menuStore";
import createApi from "@/api";
import { Events } from "@wailsio/runtime";

const { proxy } = getCurrentInstance()!;
const api = createApi(proxy);

// 当前页面使用store
const menuStore = useMenuStore();
const proxiesStore = useProxiesStore();
const webStore = useWebStore();

// 模式切换
Events.On("switchMode", (data: any) => {
  console.log("switchMode", data);
  menuStore.rule = data;
});
watch(
  () => menuStore.rule,
  (newVal) => {
    Events.Emit({ name: "mode", data: newVal });
  }
);


// 配置切换
Events.On("switchProfiles", (data: any) => {
  console.log("switchProfiles", data);
  api.switchProfile(data).then((res: any) => {
    proxiesStore.active = "";

    api.getRules().then((res) => {
      menuStore.setRuleNum(res.length);
    });

    Events.Emit({ name: "hasSwitchProfile", data: res });
  });
});

</script>

<template></template>
