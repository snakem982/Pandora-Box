<script setup lang="ts">
// 获取当前 Vue 实例的 proxy 对象
import {useWebStore} from "@/store/webStore";
import {useProxiesStore} from "@/store/proxiesStore";
import {useMenuStore} from "@/store/menuStore";
import createApi from "@/api";
import {Events} from "@wailsio/runtime"

const {proxy} = getCurrentInstance()!;
const api = createApi(proxy);

// 当前页面使用store
const menuStore = useMenuStore();
const proxiesStore = useProxiesStore();
const webStore = useWebStore();

onMounted(() => {
  Events.OnMultiple()
})


</script>

<template>

</template>