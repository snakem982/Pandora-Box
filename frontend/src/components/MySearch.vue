<script lang="ts" setup>
import {ref} from 'vue';
import {useRouter} from "vue-router";

// 搜索值
const searchValue = ref('');

// 清空搜索
const clearSearch = () => {
  searchValue.value = ''
};

// 控制下拉菜单的显示状态
const isDropdownVisible = ref(false);

// 搜索逻辑
watch(searchValue, (newValue) => {
  if (newValue) {
    isDropdownVisible.value = true;
  }
})

// 添加延时隐藏下拉菜单
const hideDropdown = () => {
  setTimeout(() => {
    isDropdownVisible.value = false;
  }, 200); // 延迟 200 毫秒
};

const searchInputRef = ref<HTMLInputElement | null>(null);

const router = useRouter()

onMounted(() => {
  if (searchInputRef.value) {
    searchInputRef.value.blur();
  }
});


</script>

<template>
  <div class="search-container" @click.stop>
    <span class="back"
          @click="router.back()">
      <el-icon>
        <icon-ep-arrow-left/>
      </el-icon>
    </span>
    <span class="forward"
          @click="router.forward()">
      <el-icon>
        <icon-ep-arrow-right/>
      </el-icon>
    </span>
    <input
        type="text"
        ref="searchInputRef"
        autocapitalize="off"
        autocomplete="off"
        spellcheck="false"
        placeholder="搜索节点"
        v-model="searchValue"
        @blur="hideDropdown"
    />

    <span class="search">
      <el-icon>
        <icon-ep-search/>
      </el-icon>
    </span>
    <span class="clear"
          @click="clearSearch">
      <el-icon>
        <icon-ep-close/>
      </el-icon>
    </span>

    <span class="minus">
      <el-icon>
        <icon-mdi-card-minus-outline/>
      </el-icon>
    </span>

    <div
        class="dropdown"
        id="dropdown"
        :style="{ display: isDropdownVisible ? 'block' : 'none' }">
      <ul>
        <li class="group">A</li>
        <li>Alice</li>
        <li>Adam</li>
        <li class="group">B</li>
        <li>Bob</li>
        <li>Betty</li>
        <li class="group">C</li>
        <li>Charlie</li>
        <li>Chloe</li>
      </ul>
    </div>
  </div>
</template>

<style scoped>
.search-container {
  padding-top: 25px;
  position: relative;
}

.search-container input {
  width: 138px;
  padding: 10px 30px 10px 40px;
  border: none;
  border-radius: 10px;
  background-color: rgba(255, 255, 255, 0.06);
  color: var(--text-color);
  font-size: 12px;
  margin-left: 26px;
  margin-top: -8px;
  position: absolute;
}

.search-container input:focus {
  outline: none;
}

.back {
  margin-left: 8px;
  font-size: 18px;
  color: var(--text-color);
}

.back:hover {
  cursor: pointer;
}

.forward {
  margin-left: 20px;
  font-size: 18px;
  color: var(--text-color);
}

.forward:hover {
  cursor: pointer;
}

.search-container .search {
  margin-left: 38px;
  margin-top: 28px;
  font-size: 18px;
  color: var(--text-color);
}

.search-container .clear {
  margin-left: 275px;
  margin-top: -18px;
  font-size: 12px;
  color: var(--text-color);
  cursor: pointer;
  display: block;
  position: absolute;
}

.search-container .minus {
  margin-right: 26px;
  float: right;
  font-size: 16px;
  color: var(--text-color);
}

.dropdown {
  position: absolute; /* 确保定位基于父容器 */
  margin-top: 8px;
  transform: translateX(66px); /* 调整偏移量，与输入框左边对齐 */
  width: 250px; /* 与输入框宽度一致 */
  border-radius: 5px;
  background-color: rgba(0, 0, 0, 0.8); /* 背景透明 */
  z-index: 9999; /* 确保下拉框显示在最上层 */
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.2);
  color: var(--text-color); /* 文字颜色为白色 */
  font-size: 12px;
}

.dropdown ul {
  margin: 0;
  padding: 0;
  list-style: none;
}

.dropdown ul .group {
  padding: 5px;
  font-weight: bold;
  color: #ccc; /* 分组标题颜色稍浅 */
  background-color: transparent; /* 背景透明 */
  border-bottom: 1px solid rgba(255, 255, 255, 0.2); /* 分割线透明 */
}

.dropdown ul li {
  padding: 10px;
  cursor: pointer;
}

.dropdown ul li:hover {
  background-color: rgba(255, 255, 255, 0.1); /* 鼠标悬停时背景微亮 */
}
</style>
