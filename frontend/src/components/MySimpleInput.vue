<template>
  <div class="input-wrapper">
    <input
        autocapitalize="off"
        autocomplete="off"
        spellcheck="false"
        @input="handleInput"
        :placeholder="placeholder"
        v-model="inputValue"
        class="custom-input"
    />
    <button
      v-if="inputValue"
      @click="clearInput"
      class="clear-button"
    >
      ✕
    </button>
  </div>
</template>

<script setup>
const props = defineProps({
  placeholder: String, // 占位符文本
  onInputChange: Function // 值变化时触发的回调方法
});

const inputValue = ref(''); // 双向绑定的输入值
const inputTimeout = ref(null); // 用于存储防抖的定时器

function handleInput(event) {
  const newValue = event.target.value;

  // 清除之前的定时器
  if (inputTimeout.value) {
    clearTimeout(inputTimeout.value);
  }

  // 设置新的定时器，延迟触发回调函数
  inputTimeout.value = setTimeout(() => {
    if (props.onInputChange) {
      props.onInputChange(newValue); // 调用回调函数
    }
  }, 500); // 设置防抖延迟时间，单位为毫秒
}

function clearInput() {
  inputValue.value = ''; // 清空输入框
  if (props.onInputChange) {
    props.onInputChange(''); // 通知父组件输入框已清空
  }
}
</script>

<style scoped>
.input-wrapper {
  position: relative;
  width: 100%;
}

.custom-input {
  width: 100%; /* 撑满宽度 */
  padding: 8px;
  padding-right: 32px; /* 为清除按钮预留空间 */
  border: 2px solid var(--text-color); /* 边框 */
  border-radius: 8px; /* 圆角 */
  background-color: transparent; /* 背景透明 */
  color: var(--text-color);
  font-size: 14px; /* 字体大小 */
  box-sizing: border-box; /* 包含 padding 和边框 */
  outline: none; /* 移除默认 outline */
  transition: border-color 0.3s ease-in-out, box-shadow 0.3s ease-in-out; /* 动态效果 */
}

.custom-input:focus {
  background-color: rgba(255, 255, 255, 0.06);
}

.custom-input::placeholder {
  color: rgba(255, 255, 255, 0.6); /* 占位符颜色稍微透明 */
}

.clear-button {
  position: absolute;
  top: 50%;
  right: 10px;
  transform: translateY(-50%);
  background: none;
  border: none;
  color: var(--text-color);
  font-size: 14px;
  cursor: pointer;
  padding: 0;
}

.clear-button:hover {
  color: rgba(255, 255, 255, 0.8);
}
</style>
