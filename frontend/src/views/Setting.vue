<script setup lang="ts">
import MyConfig from "@/components/setting/MyConfig.vue";

const distanceFromTop = ref(140)

</script>

<template>
  <MyLayout :top-height="distanceFromTop-15" :bottom-height="distanceFromTop+25">
    <template #top>
      <MySearch></MySearch>
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