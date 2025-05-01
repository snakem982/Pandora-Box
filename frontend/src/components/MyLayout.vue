<script setup lang="ts">
defineProps({
  topHeight: {
    type: Number,
    default: 65
  },
  bottomHeight: {
    type: Number,
    default: 95
  },
})
</script>

<template>
  <div v-px-drag class="top" :style="{ '--layout-top-height': topHeight + 'px' }">
    <slot name="top"></slot>
  </div>
  <div v-px-drag class="bottom" :style="{ '--layout-bottom-height': bottomHeight + 'px' }">
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