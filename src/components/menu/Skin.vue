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
           v-for="(item,index) in theme"
           :key="index">
        <span class="dropdown-label"
              @click="changeBackground(item)">
          {{ t("bg." + item.id) }}
        </span>
        <span v-if="item.allowUpload"
              class="upload-button"
              :title="t('bg.upload')"
              :aria-label="t('bg.upload')"
              @click.stop="triggerUpload(item)">
          <icon-mdi-upload/>
        </span>
      </div>
    </div>
    <input ref="fileInput"
           class="file-input"
           type="file"
           accept="image/*"
           @change="handleFileChange"/>
  </div>
</template>

<script setup lang="ts">
import {useI18n} from 'vue-i18n';
import {useMenuStore} from "@/store/menuStore";

interface ThemeOption {
  id: string;
  bg?: string | string[];
  rand?: boolean;
  allowUpload?: boolean;
}

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
const changeBackground = (item: ThemeOption) => {
  if (item.allowUpload) {
    const custom = localStorage.getItem(getCustomBackgroundKey(item.id));
    if (custom) {
      menuStore.setBackground(custom);
      return;
    }
    triggerUpload(item);
    return;
  }
  let url: string;
  if (Array.isArray(item.bg)) {
    url = getRandom(item.bg);
    if (item["rand"]) {
      url = "url('" + url + "&date=" + Date.now() + "')";
    }
  } else if (typeof item.bg === "string") {
    url = item.bg;
  } else {
    console.warn(`Theme "${item.id}" is missing a background definition.`);
    return;
  }
  menuStore.setBackground(url);
};

const theme = ref<ThemeOption[]>([]);
const fileInput = ref<HTMLInputElement | null>(null);
const pendingThemeId = ref<string | null>(null);

const getCustomBackgroundKey = (id: string) => `custom-bg-${id}`;

const triggerUpload = (item: ThemeOption) => {
  pendingThemeId.value = item.id;
  fileInput.value?.click();
};

const handleFileChange = (event: Event) => {
  const target = event.target as HTMLInputElement;
  const file = target.files?.[0];
  if (!file) {
    target.value = '';
    return;
  }
  const reader = new FileReader();
  reader.onload = () => {
    const result = reader.result;
    if (typeof result === 'string' && pendingThemeId.value) {
      const cssValue = `url('${result}')`;
      menuStore.setBackground(cssValue);
      localStorage.setItem(getCustomBackgroundKey(pendingThemeId.value), cssValue);
    }
  };
  reader.onloadend = () => {
    target.value = '';
    pendingThemeId.value = null;
  };
  reader.readAsDataURL(file);
};
onMounted(async () => {
  try {
    const response = await fetch("/json/theme.json");
    theme.value = await response.json() as ThemeOption[];
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
  min-width: 80px;
  position: absolute;
  bottom: 32px;
  margin-left: 30px;
  transform: translateX(-50%);
  background-color: var(--skin-bg-color);
  color: var(--text-color);
  padding: 10px;
  border-radius: 5px;
  box-shadow: var(--skin-box-shadow);
  text-align: center;
  z-index: 1;
  transition: all 0.3s ease;
}

.dropdown-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
}

.dropdown-label {
  flex: 1;
  padding: 5px 10px;
  border-radius: 3px;
  cursor: pointer;
  transition: background-color 0.3s ease;
}

.dropdown-label:hover {
  background-color: var(--skin-hover-color);
}

.upload-button {
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 16px;
  padding: 4px 6px;
  border-radius: 3px;
  cursor: pointer;
  transition: background-color 0.3s ease;
}

.upload-button:hover {
  background-color: var(--skin-hover-color);
}

.file-input {
  display: none;
}
</style>
