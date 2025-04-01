<template>
  <div class="custom-style">
    <span class="liable">Tun:</span>
    <el-segmented v-model="value" :options="options" />
  </div>
</template>

<script lang="ts" setup>
import { ref } from 'vue'

const value = ref('Mixed')

const options = ['Mixed', 'Gvisor', 'System']
</script>

<style scoped>
.custom-style .el-segmented{
  height: 30px;
  margin-left: 10px;
  border: 1px solid var(--text-color);
  background: rgba(255, 255, 255, 0.1);
  --el-segmented-item-selected-color: var(--text-color);
  --el-segmented-item-selected-bg-color: var(--left-item-selected-bg);
  --el-border-radius-base: 5px;
  color: var(--text-color);
  font-size: 15px;
}

.liable {
  font-size: 18px;
  font-weight: bold;
}
</style>