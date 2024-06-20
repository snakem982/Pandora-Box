<script setup lang="ts">
import {onMounted, ref} from "vue";

const nima = defineProps({
  proxy: Object
})

const delay = ref(0)

onMounted(() => {
  let history: any = (nima.proxy as any).history
  if (history.length > 0) {
    delay.value = history[history.length - 1]['delay']
  }
})

</script>

<template>
  <el-text size="large" type="danger" v-if="delay>1000">{{ delay }} ms</el-text>
  <el-text size="large" type="warning" v-else-if="delay>500">{{ delay }} ms</el-text>&emsp;
  <el-text size="large" type="success" v-else-if="delay>0">{{ delay }} ms</el-text>&emsp;
</template>

<style scoped>

</style>