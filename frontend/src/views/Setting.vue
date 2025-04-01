<script setup lang="ts">

import MyConfig from "@/components/setting/MyConfig.vue";
</script>

<template>
  <MyLayout>
    <template #top>
      <MySearch></MySearch>
    </template>
    <template #bottom>
      <div class="title">设置</div>
      <MyConfig></MyConfig>
    </template>
  </MyLayout>
</template>

<style scoped>
.title {
  font-size: 32px;
  font-weight: bold;
  margin-left: 38px;
}
</style>