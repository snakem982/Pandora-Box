<template>
  <div class="dropdown-container"
       @mouseenter="showDropdown"
       @mouseleave="hideDropdown">
    <el-icon @click="quit" class="dropdown-button">
      <icon-mdi-power/>
    </el-icon>
    <div class="dropdown-content"
         v-show="isDropdownVisible"
         @mouseenter="cancelHide">
      <div class="dropdown-item" @click="quit">
        {{ t('quit') }}
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import {useI18n} from 'vue-i18n';
import {ElLoading} from "element-plus";
import {Events} from "@/runtime";
import createApi from "@/api";

// 国际化
const {t} = useI18n();

// 获取当前 Vue 实例的 proxy 对象
const {proxy} = getCurrentInstance()!;
const api = createApi(proxy);

// 下拉框
const isDropdownVisible = ref(false);
let hideTimeout: any;

// 显示下拉框
const showDropdown = () => {
  clearTimeout(hideTimeout);
  isDropdownVisible.value = true;
};

// 隐藏下拉框（带延迟）
const hideDropdown = () => {
  hideTimeout = setTimeout(() => {
    isDropdownVisible.value = false;
  }, 200); // 延迟200ms隐藏
};

// 鼠标进入下拉框内容时取消隐藏
const cancelHide = () => {
  clearTimeout(hideTimeout);
};

// 退出
const isQuitting = ref(false);
let quitInterval: ReturnType<typeof setInterval> | undefined;
let quitLoading: ReturnType<typeof ElLoading.service> | undefined;

const showQuitCountdown = (success: boolean) => {
  let countdown = 3;
  quitLoading = ElLoading.service({
    lock: true,
    fullscreen: true,
    text: `${t('quit-countdown')} ${countdown}s`,
    background: "color-mix(in srgb, var(--skin-bg-color) 72%, transparent)",
    customClass: "quit-loading",
  });

  quitInterval = setInterval(() => {
    countdown -= 1;
    quitLoading?.setText(`${t('quit-countdown')} ${countdown}s`);
    if (countdown <= 0) {
      clearInterval(quitInterval);
      quitInterval = undefined;
      quitLoading?.close();
      quitLoading = undefined;
      Events.Emit({name: "doQuit", data: success});
    }
  }, 1000);
};

const quit = () => {
  if (isQuitting.value) return;
  isQuitting.value = true;

  api.exit()
      .then(() => showQuitCountdown(true))
      .catch(() => showQuitCountdown(false));
}

// 监听准备退出
onMounted(() => Events.On("readyToQuit", quit))
onBeforeUnmount(() => {
  if (quitInterval) clearInterval(quitInterval);
  quitLoading?.close();
})

</script>

<style scoped>
.dropdown-container {
  position: relative;
  display: inline-block;
}

.dropdown-button {
  margin-left: 0;
  font-size: 20px;
  color: var(--text-color);
  border: none;
  border-radius: 5px;
  cursor: pointer;
}

.dropdown-content {
  font-size: 14px;
  min-width: 50px;
  position: absolute;
  bottom: 32px;
  margin-left: 30px;
  transform: translateX(-50%);
  background-color: var(--skin-bg-color);
  color: var(--text-color);
  padding: 10px;
  border-radius: 5px;
  text-align: left;
  z-index: 1;
  transition: all 0.3s ease;
}

.dropdown-item {
  padding: 5px 10px;
  border-radius: 3px;
  cursor: pointer;
  transition: background-color 0.3s ease;
}

.dropdown-item:hover {
  background-color: var(--skin-hover-color);
}

:global(.quit-loading) {
  --el-color-primary: var(--text-color);
}

:global(.quit-loading .el-loading-text) {
  color: var(--text-color) !important;
}

:global(.quit-loading .el-loading-spinner .path) {
  stroke: var(--text-color) !important;
}
</style>
