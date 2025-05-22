<script setup lang="ts">
import MyConfig from "@/components/setting/MyConfig.vue";
</script>

<template>
  <MyLayout>
    <template #top>
      <el-space class="space">
        <div class="title">
          {{ $t('setting.title') }}
        </div>
      </el-space>
    </template>
    <template #bottom>
      <MyConfig></MyConfig>
    </template>
  </MyLayout>
</template>

<style scoped>

.space {
  margin-top: 20px;
}

.title {
  font-size: 32px;
  font-weight: bold;
  margin-left: 10px;
}
</style>