<script setup lang="ts">
import MyHr from "@/components/MyHr.vue";

const props = defineProps({
  hrShow: {
    type: Boolean,
    default: false
  }
})

const topHeight = ref(65)
const bottomHeight = ref(95)

// 调整顶部高度
const upFromTop = function (distance: number) {
  if (props.hrShow) {
    topHeight.value = distance - 10;
    bottomHeight.value = distance + 10;
  } else {
    topHeight.value = distance - 15;
    bottomHeight.value = distance + 15;
  }
};

</script>

<template>
  <div class="top" :style="{ '--layout-top-height': topHeight + 'px' }">
    <MySearch></MySearch>
    <slot name="top"></slot>
    <MyHr :update="upFromTop" v-show="hrShow" style="margin-top: 10px"></MyHr>
  </div>
  <div class="bottom" :style="{ '--layout-bottom-height': bottomHeight + 'px' }">
    <slot name="bottom"></slot>
  </div>
</template>

<style scoped>
.top {
  height: var(--layout-top-height);
  flex-shrink: 0;
  padding-left: 18px;
}

.bottom {
  flex-grow: 1;
  overflow-y: auto;
  height: calc(100% - var(--layout-bottom-height));
  padding-top: 10px;
  padding-bottom: 20px;
  overscroll-behavior: none;
  padding-left: 18px;
}

.bottom::-webkit-scrollbar {
  width: 5px;
  padding-bottom: 20px;
}

.bottom::-webkit-scrollbar-track {
  background: transparent;
}

.bottom::-webkit-scrollbar-thumb {
  background: var(--scrollbar-bg);
  border-radius: 2px;
  transition: background 0.3s ease, box-shadow 0.3s ease;
}

.bottom::-webkit-scrollbar-thumb:hover {
  background: var(--scrollbar-hover-bg);
  box-shadow: var(--scrollbar-hover-shadow);
}
</style>