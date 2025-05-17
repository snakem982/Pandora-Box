<template>
  <div class="dropdown-container"
       @mouseenter="showDropdown"
       @mouseleave="hideDropdown">
    <el-icon class="dropdown-button">
      <icon-mdi-tshirt-crew-outline/>
    </el-icon>
    <div class="dropdown-content"
         v-show="isDropdownVisible"
         @mouseenter="cancelHide">
      <div class="dropdown-item"
           @click="changeBackground(item)"
           v-for="(item,index) in theme"
           :key="index">
        {{ t("bg." + item.id) }}
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import {useI18n} from 'vue-i18n';
import {useMenuStore} from "@/store/menuStore";

// 存储背景主题
const menuStore = useMenuStore()

// 国际化
const {t} = useI18n();

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

// 获取随机元素
function getRandom(arr: any[]) {
  if (arr.length === 1) return arr[0];
  return arr[Math.floor(Math.random() * arr.length)];
}

// 切换背景
const changeBackground = (item: any) => {
  let url = item.bg;
  if (Array.isArray(item.bg)) {
    url = getRandom(item.bg);
    if (item["rand"]) {
      url = "url('" + url + "&date=" + Date.now() + "')"
    }
  }
  menuStore.setBackground(url)
}

const theme = ref(null);
onMounted(async () => {
  try {
    const response = await fetch("/json/theme.json");
    theme.value = await response.json();
  } catch (error) {
    console.error("获取 JSON 失败", error);
  }
});

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
