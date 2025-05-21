<script setup lang="ts">
import createApi from "@/api";
import {useProxiesStore} from "@/store/proxiesStore";
import {useMenuStore} from "@/store/menuStore";
import {useSettingStore} from "@/store/settingStore";
import {useI18n} from "vue-i18n";
import {pError, pLoad} from "@/util/pLoad";
import {useWebStore} from "@/store/webStore";

const {t} = useI18n();

// 获取当前 Vue 实例的 proxy 对象
const {proxy} = getCurrentInstance()!;
const api = createApi(proxy);

// 当前页面双向绑定对象
const groupList = ref<string[]>([]);
const nodeList = ref<any[]>([]);

// 当前页面使用store
const proxiesStore = useProxiesStore();
const menuStore = useMenuStore();
const settingStore = useSettingStore();
const webStore = useWebStore();

// 获取分组
async function groups() {
  // 活跃分组
  const active = proxiesStore.active;

  const temp = await api.getGroups();
  switch (menuStore.rule) {
    case "rule":
      groupList.value = temp;
      let find = false
      for (let i = 0; i < temp.length; i++) {
        if (temp[i] === active) {
          find = true
          break
        }
      }
      if (find) {
        if (active == "GLOBAL") {
          proxiesStore.setActive(temp[0]);
        }
      } else {
        proxiesStore.setActive(temp[0]);
      }
      break;
    case "global":
      groupList.value = temp.concat(["GLOBAL"]);
      if (!active) {
        proxiesStore.setActive(temp[0]);
      }
      break;
    case "direct":
      groupList.value = [];
      break;
  }
}

// 获取节点列表
async function nodes() {
  if (menuStore.rule == "direct") {
    nodeList.value = [];
    return;
  }
  nodeList.value = await api.getProxies(
      proxiesStore.active,
      proxiesStore.isHide,
      proxiesStore.isSort
  ); // 更新响应式数据
}

// 设置活跃分组
async function setActive(value: any) {
  if (proxiesStore.active == value) {
    return;
  }
  proxiesStore.setActive(value);
  await nodes();
}

// 设置隐藏
async function setHide() {
  proxiesStore.setHide(!proxiesStore.isHide);
  await nodes();
}

// 设置排序
async function setSort() {
  proxiesStore.setSort(!proxiesStore.isSort);
  await nodes();
}

// 设置分组
function setVertical() {
  proxiesStore.setVertical(!proxiesStore.isVertical);
  atStart.value = true;
  atEnd.value = true;
  setTimeout(() => {
    updateButtonVisibility();
  }, 200);
}

// 设置代理
async function setProxy(now: any, name: string) {
  if (now) {
    return;
  }
  try {
    await api.setProxy(proxiesStore.active, {name});
    await nodes();
  } catch (error) {
    console.error(error);
  }
}

// 测试延迟
function testDelay() {
  pLoad(t("proxies.loading"), async () => {
    try {
      await api.getDelay(proxiesStore.active, settingStore.testUrl, 3000);
      await nodes();
    } catch (e) {
      if (e['message']) {
        pError(e['message'])
      }
    }
  });
}

const proxyGroup = ref(null);
const atStart = ref(true); // 标记是否在最左边
const atEnd = ref(true); // 标记是否在最右边

const updateButtonVisibility = () => {
  if (proxyGroup.value) {
    const scrollLeft = proxyGroup.value.scrollLeft;
    const scrollWidth = proxyGroup.value.scrollWidth;
    const clientWidth = proxyGroup.value.clientWidth;

    atStart.value = scrollLeft === 0;
    atEnd.value = scrollLeft + clientWidth >= scrollWidth;
  }
};

const scrollLeft = () => {
  if (proxyGroup.value) {
    proxyGroup.value.scrollLeft -= proxyGroup.value.clientWidth + 15;
  }
};

const scrollRight = () => {
  if (proxyGroup.value) {
    proxyGroup.value.scrollLeft += proxyGroup.value.clientWidth - 15;
  }
};

let isScrolling: any;
const handleScroll = () => {
  clearTimeout(isScrolling);
  isScrolling = setTimeout(() => {
    updateButtonVisibility();
  }, 200); // 200ms 延迟
};

const isDropdownOpen = ref(false);

// 添加延时隐藏下拉菜单
let isOvering: any;
const hideDropdown = () => {
  isOvering = setTimeout(() => {
    isDropdownOpen.value = false;
  }, 200); // 延迟 200 毫秒
};

// 鼠标进入下拉菜单时，清除延时隐藏
const enterDropDown = () => {
  clearTimeout(isOvering);
  isDropdownOpen.value = true;
};

let fresh: any = null;
onMounted(async () => {
  await groups();
  await nodes();
  updateButtonVisibility();
  // 监听 resize 事件
  window.addEventListener("resize", updateButtonVisibility);
  // 创建刷新定时器
  fresh = setInterval(async () => {
    await nodes();
  }, 10000);
});

onBeforeUnmount(() => {
  // 清除定时器
  clearInterval(fresh);
  // 移除 resize 事件监听
  window.removeEventListener("resize", updateButtonVisibility);
});

// 监听具体状态
watch(() => menuStore.rule, // 监听 store 中的某个状态
    async () => {
      await groups();
      await nodes();
      updateButtonVisibility();
    }
);

watch(() => webStore.fProfile, async () => {
  await groups();
  await nodes();
  updateButtonVisibility();
})

watch(() => proxiesStore.now, async () => {
  await nodes();
})

</script>

<template>
  <MyLayout hr-show>
    <template #top>
      <MySearch></MySearch>
      <el-space class="space">
        <div class="title">
          {{ $t("proxies.title") }}
        </div>
        <div class="proxy-option">
          <el-tooltip :content="$t('proxies.test')" placement="top">
            <el-icon @click="testDelay" class="proxy-option-btn">
              <icon-mdi-speedometer/>
            </el-icon>
          </el-tooltip>

          <el-tooltip
              :content="
              proxiesStore.isHide
                ? $t('proxies.hide-on')
                : $t('proxies.hide-off')
            "
              placement="top"
          >
            <el-icon @click="setHide" class="proxy-option-btn">
              <icon-mdi-eye-off v-if="proxiesStore.isHide"/>
              <icon-mdi-eye v-else/>
            </el-icon>
          </el-tooltip>

          <el-tooltip
              :content="
              proxiesStore.isSort
                ? $t('proxies.sort-on')
                : $t('proxies.sort-off')
            "
              placement="top"
          >
            <el-icon @click="setSort" class="proxy-option-btn">
              <icon-mdi-sort-ascending v-if="proxiesStore.isSort"/>
              <icon-mdi-sort v-else/>
            </el-icon>
          </el-tooltip>

          <el-tooltip
              :content="
              proxiesStore.isVertical
                ? $t('proxies.vertical-on')
                : $t('proxies.vertical-off')
            "
              placement="top"
          >
            <el-icon @click="setVertical" class="proxy-option-btn">
              <icon-mdi-arrow-expand-vertical v-if="proxiesStore.isVertical"/>
              <icon-mdi-arrow-expand-horizontal v-else/>
            </el-icon>
          </el-tooltip>
        </div>
      </el-space>

      <div
          class="dropdown"
          v-if="proxiesStore.isVertical && menuStore.rule != 'direct' && groupList.length > 0"
      >
        <button
            class="dropdown-btn"
            @mouseenter="enterDropDown"
            @mouseleave="hideDropdown"
        >
          {{ proxiesStore.active }}
        </button>
        <ul
            v-if="isDropdownOpen"
            @mouseenter="enterDropDown"
            @mouseleave="hideDropdown"
            class="dropdown-list"
        >
          <li
              v-for="item in groupList"
              :key="item + '-gv'"
              @click="setActive(item)"
              class="dropdown-item"
          >
            {{ item }}
          </li>
        </ul>
      </div>

      <div
          class="button-container"
          v-if="!proxiesStore.isVertical && menuStore.rule != 'direct' && groupList.length > 0"
      >
        <el-icon v-if="!atStart" @click="scrollLeft" class="scroll-left">
          <icon-mdi-arrow-expand-left/>
        </el-icon>
        <div @scroll="handleScroll" ref="proxyGroup" class="proxy-group">
          <button
              :class="
              proxiesStore.active == item
                ? 'proxy-group-title proxy-group-title-select'
                : 'proxy-group-title'
            "
              @click="setActive(item)"
              v-for="item in groupList"
              :key="item + '-g'"
          >
            {{ item }}
          </button>
        </div>
        <el-icon v-if="!atEnd" class="scroll-right" @click="scrollRight">
          <icon-mdi-arrow-expand-right/>
        </el-icon>
      </div>
    </template>


    <template #bottom>
      <div class="proxy-nodes">
        <div
            :class="
            node['now']
              ? 'proxy-nodes-card proxy-node-select'
              : 'proxy-nodes-card'
          "
            v-for="node in nodeList"
            @click="setProxy(node['now'], node['name'])"
            :key="node['name']"
        >
          <div class="proxy-nodes-title">
            <span :title="node['name']">  
              {{ node["name"] }}
            </span>

          </div>
          <div class="proxy-nodes-tags">
            <span class="proxy-nodes-tags-left">
              {{ node["type"] }}
            </span>
            <span :class="'proxy-nodes-tags-right ' + node['toClass']">
              {{ node["delay"] }} ms
            </span>
          </div>
        </div>
      </div>

      <div class="proxy-nodes" v-if="menuStore.rule == 'direct'">
        {{ $t("proxies.direct") }}
      </div>
    </template>
  </MyLayout>
</template>

<style scoped>
.space {
  margin-top: 15px;
}

.title {
  font-size: 32px;
  font-weight: bold;
  margin-left: 10px;
}

.proxy-option {
  margin-left: 10px;
  font-size: 30px;
  padding-top: 10px;
}

.proxy-option-btn {
  margin-right: 15px;
}

.proxy-option-btn:hover {
  cursor: pointer;
  color: var(--hr-color);
}

.button-container {
  display: flex;
  align-items: center;
  width: 95%;
  margin-left: 10px;
  min-height: 50px;
}

.proxy-group {
  display: flex;
  gap: 10px;
  margin: 12px 0 3px 0;
  overflow-x: hidden;
  scroll-behavior: smooth;
}

.scroll-left {
  cursor: pointer;
  border: none;
  margin-right: 10px;
}

.scroll-right {
  cursor: pointer;
  border: none;
  margin-left: 10px;
}

.scroll-left[hidden],
.scroll-right[hidden] {
  display: none;
}

.proxy-group-title {
  background-color: transparent;
  color: var(--text-color);
  border: 2px solid var(--hr-color);
  border-radius: 8px;
  padding: 6px 10px;
  font-size: 15px;
  text-align: center;
  cursor: pointer;
  box-shadow: var(--left-nav-shadow);
  white-space: nowrap;
}

.proxy-group-title:hover,
.proxy-group-title-select {
  background-color: var(--left-item-selected-bg);
  box-shadow: var(--left-nav-hover-shadow);
  border-color: var(--text-color);
}

.proxy-nodes {
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
  padding: 0;
  color: var(--text-color);
  margin-left: 12px;
  width: 95%;
}

.proxy-nodes-card {
  width: calc(33% - 41px);
  max-width: 210px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-radius: 8px;
  padding: 8px 12px;
  background: rgba(255, 255, 255, 0.1);
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  box-shadow: var(--left-nav-shadow);
  line-height: 1.3;
}

.proxy-nodes-card:hover,
.proxy-node-select {
  background-color: var(--left-item-selected-bg);
  box-shadow: var(--left-nav-hover-shadow);
  border: 2px solid var(--text-color);
  cursor: pointer;
}

.proxy-nodes-title {
  font-size: 14px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.proxy-nodes-tags {
  font-size: 14px;
  display: flex;
  margin-top: 8px;
  justify-content: space-between;
}

.proxy-nodes-tags-left {
  flex: 1;
}

.proxy-nodes-tags-right {
  text-align: right;
}

.toHidden {
  display: none;
}

.dropdown {
  position: relative;
  display: inline-block;
  width: 95%;
  margin: 12px 10px 5px 10px;
}

.dropdown-btn {
  background-color: var(--left-item-selected-bg);
  box-shadow: var(--left-nav-hover-shadow);
  border: 2px solid var(--text-color);
  color: var(--text-color);
  padding: 5px 10px;
  cursor: pointer;
  font-size: 15px;
  outline: none;
  border-radius: 8px;
  min-width: 204px;
}

.dropdown-btn:hover {
  opacity: 0.8;
}

.dropdown-list {
  position: absolute;
  background: rgba(0, 0, 0, 0.8);
  border: 2px solid white;
  margin-top: 4px;
  padding: 0;
  list-style: none;
  min-width: 200px;
  z-index: 20;
  border-radius: 8px;
  font-size: 15px;
  text-align: center;
  max-height: calc(100vh - 230px);
  overflow-y: auto;
}

.dropdown-item {
  color: white;
  padding: 8px;
  cursor: pointer;
}

.dropdown-item:hover {
  background: rgba(255, 255, 255, 0.2);
}
</style>
