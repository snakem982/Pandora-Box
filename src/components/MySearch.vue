<script lang="ts" setup>
import {useRouter} from "vue-router";
import {debounce} from "lodash";
import createApi from "@/api";
import {useProxiesStore} from "@/store/proxiesStore";

// 获取当前 Vue 实例的 proxy 对象
const {proxy} = getCurrentInstance()!;
const api = createApi(proxy);

// 当前组件使用store
const proxiesStore = useProxiesStore();

// 搜索值
const searchValue = ref('');
const searchList = ref<any[]>([]);

// 清空搜索
const clearSearch = () => {
  searchValue.value = ''
};

// 控制下拉菜单的显示状态
const isDropdownVisible = ref(false);

// 防抖
const debouncedSearch = debounce((keyword) => {
  const lowerKeyword = keyword.toLowerCase();

  // 获取前可用节点
  const promise = api.getProxies(proxiesStore.active, true, true)
  promise.then(arr => {
    if (arr && arr.length === 0) return;
    const result: any[] = [];
    arr.some(item => {
      if (item.name.toLowerCase().includes(lowerKeyword)) {
        result.push(item);
      }

      // 找到 10 个后，提前终止遍历
      return result.length >= 10;
    });
    if (result.length === 0) return;
    searchList.value = result;
    isDropdownVisible.value = true;
  })

}, 300);

watch(searchValue, (keyword) => {
  if (keyword.trim().length == 0) {
    isDropdownVisible.value = false;
    return
  }

  debouncedSearch(keyword.trim());
});

// 组件销毁时，取消防抖任务
onUnmounted(() => {
  debouncedSearch.cancel();
});

// 添加延时隐藏下拉菜单
const hideDropdown = () => {
  setTimeout(() => {
    isDropdownVisible.value = false;
  }, 200); // 延迟 200 毫秒
};

const searchInputRef = ref<HTMLInputElement | null>(null);

const router = useRouter()

const isWindows = ref(false)
onMounted(() => {
  if (searchInputRef.value) {
    searchInputRef.value.blur();
  }
  // @ts-ignore
  if (window["pxShowBar"]) {
    isWindows.value = true;
  }
});

// 设置代理
async function changeProxy(now: any, name: any) {
  if (now) {
    return;
  }
  try {
    await api.setProxy(proxiesStore.active, {name});
    proxiesStore.setNow(name)
    searchValue.value = '';
  } catch (error) {
    console.error(error);
  }
}

</script>

<template>
  <div :class="isWindows?'search-container win':'search-container'">
    <el-space @click.stop class="no-drag">
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
      <span class="search">
        <el-icon>
          <icon-ep-search/>
        </el-icon>
      </span>
      <input
          type="text"
          ref="searchInputRef"
          autocapitalize="off"
          autocomplete="off"
          spellcheck="false"
          :placeholder="$t('search')"
          v-model="searchValue"
          @blur="hideDropdown"
      />

      <span class="clear"
            v-show="searchValue.trim().length > 0"
            @click="clearSearch">
        <el-icon>
          <icon-mdi-close/>
        </el-icon>
      </span>
    </el-space>

    <MyTitleBar class="minus no-drag"></MyTitleBar>

    <div
        class="dropdown no-drag"
        id="dropdown"
        :style="{ display: isDropdownVisible ? 'block' : 'none' }">
      <ul>
        <!--        <li class="group">A</li>-->
        <!--        <li>Alice</li>-->
        <li v-for="item in searchList" @click="changeProxy(item['now'], item['name'])">
          <span class="sName"> {{ item.name }} </span>
          <span :class="'sDelay ' + item['toClass']">{{ item.delay }} ms</span>
        </li>
      </ul>
    </div>
  </div>
</template>

<style scoped>
.search-container {
  padding-top: 25px;
  position: relative;
  -webkit-app-region: drag;
}

.win {
  padding-top: 15px;
}

.no-drag {
  -webkit-app-region: no-drag;
}

.search-container input {
  width: 138px;
  padding: 10px 30px 10px 40px;
  border: none;
  border-radius: 10px;
  background-color: rgba(255, 255, 255, 0.06);
  color: var(--text-color);
  font-size: 12px;
  margin-left: -38px;
  margin-top: -3px;
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
  margin-left: 10px;
  font-size: 18px;
  color: var(--text-color);
}

.forward:hover {
  cursor: pointer;
}

.search-container .search {
  margin-left: 20px;
  font-size: 18px;
  color: var(--text-color);
}

.search-container .clear {
  margin-left: -28px;
  margin-top: -5px;
  font-size: 14px;
  color: var(--text-color);
  cursor: pointer;
  display: block;
}

.minus {
  margin-right: 25px;
  float: right;
  font-size: 18px;
  color: var(--text-color);
  cursor: pointer;
}

.dropdown {
  position: absolute; /* 确保定位基于父容器 */
  margin-top: 8px;
  transform: translateX(66px); /* 调整偏移量，与输入框左边对齐 */
  width: 235px; /* 与输入框宽度一致 */
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
  padding: 10px 15px;
  cursor: pointer;
}

.dropdown ul li:hover {
  background-color: rgba(255, 255, 255, 0.1); /* 鼠标悬停时背景微亮 */
}

.sName {
  display: inline-block;
  max-width: 135px;
  font-weight: bold;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.sDelay {
  float: right;
}

</style>
