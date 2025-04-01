<template>
  <div class="custom-style">
    <el-segmented v-model="value" :options="options" />
  </div>
</template>

<script lang="ts" setup>
import { ref } from 'vue'

const value = ref('规则')

const options = ['规则', '全局', '直连']
</script>

<style scoped>
.custom-style .el-segmented{
  width: 184px;
  margin-top: 23px;
  margin-left: 22px;
  border: 1px solid #ccc;
  background: rgba(255, 255, 255, 0.1);
  box-shadow: var(--left-nav-shadow);
  --el-segmented-item-selected-color: var(--text-color);
  --el-segmented-item-selected-bg-color: var(--left-item-selected-bg);
  --el-border-radius-base: 5px;
  color: var(--text-color);
  font-size: 15px;
}

.custom-style .el-segmented:hover{
  box-shadow: var(--left-nav-hover-shadow);
}

</style>