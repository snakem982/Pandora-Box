<template>
  <div id="titleBar" v-if="isWindows">
    <div id="window-controls">
      <div class="button close" @click="pxClose">
        <span>×</span>
      </div>
      <div class="button minimize" @click="pxMinimize">
        <span>–</span>
      </div>
      <div class="button maximize" @click="toggleMaximize">
        <span v-if="!isMaximized">+</span>
        <span v-else>↘︎</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import {Events} from "@/runtime";

const isMaximized = ref(false)
const isWindows = ref(false)

onMounted(() => {
  // @ts-ignore
  if (window["pxShowBar"]) {
    isWindows.value = true;
  }
})

function pxClose() {
  Events.Emit({name: "close", data: true});
}

function pxMinimize() {
  Events.Emit({name: "min", data: true});
}

function toggleMaximize() {
  Events.Emit({name: "max", data: true});
}
</script>

<style scoped>
#titleBar {
  position: fixed;
  top: 0;
  left: 0;
  width: 100px;
  height: 32px;
  background: transparent;
  z-index: 9999;
}

#window-controls {
  display: flex;
  align-items: center;
  height: 100%;
  padding-left: 12px;
  gap: 8px;
}

.button {
  width: 12px;
  height: 12px;
  border-radius: 50%;
  position: relative;
  flex-shrink: 0;
  overflow: hidden;
  transition: background 0.3s ease, filter 0.2s ease;
}

.button span {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  font-size: 10px;
  font-weight: bold;
  color: #000;
  opacity: 0;
  pointer-events: none;
  transition: opacity 0.2s ease;
}

.button.close {
  background-color: #ff5f56;
}

.button.minimize {
  background-color: #ffbd2e;
}

.button.maximize {
  background-color: #27c93f;
}

#window-controls .button {
  cursor: pointer;
}

.button:hover {
  filter: brightness(1.2);
}

.button:hover span {
  opacity: 1;
  color: #0e0e0e;
}

</style>
