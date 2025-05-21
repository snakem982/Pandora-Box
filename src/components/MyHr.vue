<script setup lang="ts">
import {defineProps, nextTick, onMounted, onUnmounted, ref} from 'vue';

const targetHr = ref<HTMLElement | null>(null);
const distanceFromTop = ref(0);
let observer: MutationObserver | null = null;

const props = defineProps({
  update: Function,
});

/**
 * 获取元素距离页面顶部的距离，即使它当前被 v-show 隐藏
 */
const getDistanceEvenIfHidden = (el: HTMLElement): number => {
  const originalDisplay = el.style.display;
  const computedStyle = getComputedStyle(el);

  let needRestore = false;

  if (computedStyle.display === 'none') {
    el.style.display = 'block';
    needRestore = true;
  }

  const rect = el.getBoundingClientRect();
  const top = rect.top + window.scrollY;

  if (needRestore) {
    el.style.display = originalDisplay;
  }

  return top;
};

/**
 * 原始的、带“值是否变化”判断的 updateDistance 函数
 */
const handleResize = () => {
  if (targetHr.value) {
    const newTop = getDistanceEvenIfHidden(targetHr.value);

    // ⚠️ 如果值没变就不更新
    if (newTop === distanceFromTop.value) return;

    distanceFromTop.value = newTop;
    props.update?.(newTop);
  }
};

onMounted(async () => {
  await nextTick(); // 确保 DOM 渲染完成
  handleResize();

  window.addEventListener('resize', handleResize);
  window.addEventListener('scroll', handleResize);

  observer = new MutationObserver(() => {
    handleResize();
  });

  observer.observe(document.body, {
    childList: true,
    subtree: true,
  });
});

onUnmounted(() => {
  window.removeEventListener('resize', handleResize);
  window.removeEventListener('scroll', handleResize);
  observer?.disconnect();
});
</script>

<template>
  <hr ref="targetHr" class="proxy-hr">
</template>

<style scoped>
.proxy-hr {
  border: none;
  border-top: 1px dashed #FFD700;
  width: 95%;
  margin: 0 10px;
}
</style>
