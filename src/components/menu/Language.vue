<template>
  <div class="dropdown-container"
       @mouseenter="showDropdown"
       @mouseleave="hideDropdown">
    <el-icon class="dropdown-button">
      <icon-mdi-translate/>
    </el-icon>
    <div class="dropdown-content"
         v-show="isDropdownVisible"
         @mouseenter="cancelHide">
      <div class="dropdown-item" @click="changeLang('zh')">简体中文</div>
      <div class="dropdown-item" @click="changeLang('en')">English</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import {useI18n} from 'vue-i18n';
import {useMenuStore} from "@/store/menuStore";
import {Events} from "@/runtime"

// 存储语言
const menuStore = useMenuStore()

// 国际化
const {locale, t} = useI18n();

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

// tray 翻译id
const trayMenuId = [
  'tray.show',
  'tray.rule',
  'tray.global',
  'tray.direct',
  'tray.profiles',
  'tray.proxy',
  'tray.tun',
  'tray.quit'
]

// 发送 tray 翻译
function sendTranslation() {
  const translate: any = {}
  trayMenuId.forEach(item => {
    translate[item] = t(item)
  })
  Events.Emit({
    name: "translate",
    data: translate
  })
}


// 切换语言
const changeLang = (value: any) => {
  locale.value = value
  menuStore.setLanguage(value)
  sendTranslation()
}

onMounted(() => {
  // 设置语言
  if (menuStore.language) {
    locale.value = menuStore.language
    sendTranslation()
  }
})
</script>

<style scoped>
.dropdown-container {
  position: relative;
  display: inline-block;
}

.dropdown-button {
  margin-left: 20px;
  font-size: 20px;
  color: var(--text-color);
  border: none;
  border-radius: 5px;
  cursor: pointer;
}

.dropdown-content {
  font-size: 14px;
  width: 80px;
  position: absolute;
  bottom: 32px;
  margin-left: 30px;
  transform: translateX(-50%);
  background-color: rgba(0, 0, 0, 0.8);
  color: var(--text-color);;
  padding: 10px;
  border-radius: 5px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  text-align: center;
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
  background-color: #555;
}
</style>
