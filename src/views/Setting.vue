<script setup lang="ts">
import MyConfig from "@/components/setting/MyConfig.vue";
import MyHr from "@/components/proxies/MyHr.vue";

const distanceFromTop = ref(140)
const upFromTop = function (distance: number) {
  distanceFromTop.value = distance
}

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
      <MyHr :update="upFromTop" v-show="false"></MyHr>
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