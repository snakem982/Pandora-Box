<script setup lang="ts">
import MyHr from "@/components/proxies/MyHr.vue";

const distanceFromTop = ref(195)
const upFromTop = function (distance: number) {
  distanceFromTop.value = distance
}
</script>

<template>
  <MyLayout
      :top-height="distanceFromTop - 15"
      :bottom-height="distanceFromTop + 25"
  >
    <template #top>
      <MySearch></MySearch>
      <MyHr :update="upFromTop" v-show="false"></MyHr>
    </template>
    <template #bottom>
      <div class="content">
        <MyChart></MyChart>
        <MyTest></MyTest>
        <MyIp></MyIp>
      </div>
    </template>
  </MyLayout>
</template>

<style scoped>
.content {
  margin-top: 10px;
}
</style>