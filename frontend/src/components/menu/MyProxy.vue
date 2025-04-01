<script setup lang="ts">
const isOn = ref(false)
const isOn2 = ref(false)

function toggle() {
  isOn.value = !isOn.value
}

function toggle2() {
  isOn2.value = !isOn2.value
}
</script>

<template>

  <div class="sub">
    <div class="switch-container" title="代理已关闭">
      <span class="switch-label">代理</span>
      <div :class="['switch', { 'switch-on': isOn }]" @click="toggle">
        <div class="switch-circle"></div>
      </div>
    </div>
    <div class="switch-container" title="虚拟网卡已关闭">
      <span class="switch-label">虚拟网卡</span>
      <div :class="['switch', { 'switch-on': isOn2 }]" @click="toggle2">
        <div class="switch-circle"></div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.sub {
  width: 184px; /* 卡片宽度 */
  background-color: rgba(255, 255, 255, 0.1); /* 半透明白色背景 */
  border: 1px solid #ccc;
  border-radius: 8px;
  box-shadow: var(--left-nav-shadow);
  line-height: 1.5;
  margin-left: 22px;
  margin-top: 25px;
  padding-bottom: 12px;
}

.sub:hover {
  box-shadow: var(--left-nav-hover-shadow);
}

.switch-container {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 16px;
  padding-left: 10px;
  margin-top: 10px;
}

.switch-label {
  color: var(--text-color);
}

.switch {
  width: 60px;
  height: 30px;
  border: 2px solid var(--text-color);
  border-radius: 15px;
  position: relative;
  background-color: transparent;
  cursor: pointer;
  transition: background-color 0.3s ease;
}

.switch.switch-on {
  background-color: var(--left-item-selected-bg);
}

.switch-circle {
  width: 26px;
  height: 26px;
  background-color: var(--text-color);
  border-radius: 50%;
  position: absolute;
  top: 2px;
  left: 2px;
  transition: left 0.3s ease;
}

.switch-on .switch-circle {
  left: 32px;
}
</style>