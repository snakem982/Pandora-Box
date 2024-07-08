<script setup lang="ts">
import {onMounted, ref} from "vue";
import {IsUnifiedDelay} from "../../wailsjs/go/main/App";

const nima = defineProps({
  proxy: Object
})

const delay = ref(0)
const warning = ref(500)
const danger = ref(1000)

onMounted(async () => {
  let history: any = (nima.proxy as any).history
  if (history.length > 0) {
    delay.value = history[history.length - 1]['delay']
  }

  const isUnifiedDelay = await IsUnifiedDelay()
  if (isUnifiedDelay == "true") {
    warning.value = 200
    danger.value = 500
  }
})

</script>

<template>
  <el-text size="large" type="danger" v-if="delay>danger">{{ delay }} ms</el-text>
  <el-text size="large" type="warning" v-else-if="delay>warning">{{ delay }} ms</el-text>&emsp;
  <el-text size="large" type="success" v-else-if="delay>0">{{ delay }} ms</el-text>&emsp;
</template>

<style scoped>

</style>