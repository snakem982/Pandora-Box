<script setup lang="ts">
const targetHr = ref(null);
const distanceFromTop = ref(0);
let observer = null;

const props = defineProps({
  update: Function,
});

const updateDistance = () => {
  if (targetHr.value) {
    const rect = targetHr.value.getBoundingClientRect();
    distanceFromTop.value = rect.top + window.scrollY;
    props.update(distanceFromTop.value)
  }
};

const handleResize = () => {
  updateDistance();
};

onMounted(() => {
  // 初次计算
  updateDistance();

  // 监听 resize 事件
  window.addEventListener('resize', handleResize);

  // 创建一个 MutationObserver 实例
  observer = new MutationObserver(() => {
    updateDistance(); // 监测到变化时重新计算距离
  });

  // 观察整个页面或者某些父级容器
  observer.observe(document.body, {
    childList: true, // 观察子节点的添加或移除
    subtree: true, // 观察整个子树
  });
});

onUnmounted(() => {
  // 清除 resize 监听
  window.removeEventListener('resize', handleResize);

  // 解除观察器绑定
  if (observer) {
    observer.disconnect();
  }
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