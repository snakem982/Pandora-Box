<template>
  <div class="mask" v-show="isDragging">
    <h3>拖拽配置文件到这里导入</h3>
  </div>
</template>

<script setup>
const isDragging = ref(false);

onMounted(() => manageDragEvents('add'));
onUnmounted(() => manageDragEvents('remove'));

function manageDragEvents(action) {
  const method = action === 'add' ? 'addEventListener' : 'removeEventListener';
  document.body[method]('dragenter', handleDragEnter);
  document.body[method]('dragover', preventDefault);
  document.body[method]('drop', handleDrop);
}

function handleDragEnter(e) {
  if (e.dataTransfer && e.dataTransfer.types.includes('Files')) {
    isDragging.value = true;
  }
}

function preventDefault(e) {
  e.preventDefault();
}

function handleDrop(e) {
  e.preventDefault();
  isDragging.value = false;

  const files = Array.from(e.dataTransfer.files);

  files.forEach((file) => {
    const reader = new FileReader();

    reader.onload = (event) => {
      console.log(`Content of ${file.name}:`);
      console.log(event.target.result);
    };

    reader.onerror = (error) => {
      console.error(`Error reading ${file.name}:`, error);
    };

    // 使用 readAsText 方法读取文件内容
    reader.readAsText(file);
  });
}


</script>

<style scoped>
.mask {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 9999;
  background: rgba(0, 0, 0, 0.6);
  display: flex;
  justify-content: center;
  align-items: center;
  color: white;
  font-size: 1.5rem;
}

h3 {
  margin: 0;
}
</style>
