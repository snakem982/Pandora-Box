<script setup lang="ts">
import createApi from '@/api';
import MyHr from "@/components/proxies/MyHr.vue";
import {useProxiesStore} from "@/store/proxiesStore";
import {useMenuStore} from "@/store/menuStore";

// 计算顶部遮挡
const distanceFromTop = ref(195)
const upFromTop = function (distance: number) {
  distanceFromTop.value = distance
}

// 获取当前 Vue 实例的 proxy 对象
const {proxy} = getCurrentInstance()!;
const api = createApi(proxy);

// 当前页面双向绑定对象
const groupList = ref<string[]>([]);
const nodeList = ref<any[]>([]);

// 当前页面使用store
const proxiesStore = useProxiesStore();
const menuStore = useMenuStore();

// 获取分组
async function groups() {
  const temp = await api.getGroups();
  switch (menuStore.rule) {
    case "rule":
      groupList.value = temp;
      break
    case "global":
      groupList.value = temp.concat(['GLOBAL'])
      break
    case "direct":
      groupList.value = []
      break
  }

  // 设置活跃分组
  const active = proxiesStore.active;
  if (!active) {
    proxiesStore.setActive(temp[0]);
  }
}

// 获取节点列表
async function nodes() {
  nodeList.value = await api.getProxies(proxiesStore.active, proxiesStore.isHide, proxiesStore.isSort); // 更新响应式数据
}

// 设置活跃分组
async function setActive(value: any) {
  proxiesStore.setActive(value)
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

// 设置代理
async function setProxy(now: any, name: string) {
  if (now) {
    return
  }
  try {
    await api.setProxy(proxiesStore.active, {name});
    await nodes();
  } catch (error) {
    console.error(error);
  }
}

onMounted(async () => {
  await groups();
  await nodes();
});


</script>

<template>
  <MyLayout :top-height="distanceFromTop-15" :bottom-height="distanceFromTop+25">
    <template #top>
      <MySearch></MySearch>
      <el-space class="space">
        <div class="title">
          {{ $t('proxies.title') }}
        </div>
        <div class="proxy-option">
          <el-tooltip
              :content="$t('proxies.test')"
              placement="top">
            <el-icon class="proxy-option-btn">
              <icon-mdi-speedometer/>
            </el-icon>
          </el-tooltip>

          <el-tooltip
              :content="proxiesStore.isHide?$t('proxies.hide-on'):$t('proxies.hide-off')"
              placement="top">
            <el-icon
                @click="setHide"
                class="proxy-option-btn">
              <icon-mdi-eye-off v-if="proxiesStore.isHide"/>
              <icon-mdi-eye v-else/>
            </el-icon>
          </el-tooltip>

          <el-tooltip
              :content="proxiesStore.isSort?$t('proxies.sort-on'):$t('proxies.sort-off')"
              placement="top">
            <el-icon
                @click="setSort"
                class="proxy-option-btn">
              <icon-mdi-sort-ascending v-if="proxiesStore.isSort"/>
              <icon-mdi-sort v-else/>
            </el-icon>
          </el-tooltip>
        </div>
      </el-space>

      <div class="proxy-group">
        <button
            :class="proxiesStore.active==item?'proxy-group-title proxy-group-title-select' : 'proxy-group-title'"
            @click="setActive(item)"
            v-for="item in groupList"
            :key="item + '-g'"
        >
          {{ item }}
        </button>
      </div>

      <MyHr :update="upFromTop"></MyHr>
    </template>
    <template #bottom>
      <div class="proxy-nodes">
        <div
            :class="node['now']?'proxy-nodes-card proxy-node-select':'proxy-nodes-card'"
            v-for="node in nodeList"
            @click="setProxy(node['now'],node['name'])"
            :key="node['name']"
        >
          <div class="proxy-nodes-title">
            {{ node['name'] }}
          </div>
          <div class="proxy-nodes-tags">
            <span class="proxy-nodes-tags-left">
              {{ node['type'] }}
            </span>
            <span :class="'proxy-nodes-tags-right ' + node['toClass']">
              {{ node['delay'] }} ms
            </span>
          </div>
        </div>
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

.proxy-group {
  display: flex;
  margin-left: 10px;
  margin-top: 10px;
  margin-bottom: 10px;
  gap: 10px;
  width: 95%;
  overflow-x: hidden;
}

.proxy-group-title {
  background-color: transparent;
  color: var(--text-color);
  border: 2px solid var(--hr-color);
  border-radius: 8px;
  padding: 6px 10px;
  font-size: 16px;
  text-align: center;
  cursor: pointer;
  box-shadow: var(--left-nav-shadow);
  white-space: nowrap;
}

.proxy-group-title:hover, .proxy-group-title-select {
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
  margin-top: 5px;
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

.proxy-nodes-card:hover, .proxy-node-select {
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

.toLow {
  color: #39FF14;
  text-shadow: 1px 1px 2px rgba(0, 0, 0, 0.5),
  0 0 2px rgba(50, 255, 50, 0.8);
}


.toMiddle {
  color: #FFD700;
  text-shadow: 1px 1px 2px rgba(0, 0, 0, 0.5),
  0 0 2px rgba(255, 215, 0, 0.8);
}

.toHigh {
  color: #FF4500;
  text-shadow: 1px 1px 2px rgba(0, 0, 0, 0.5),
  0 0 2px rgba(255, 69, 0, 0.8);
}
</style>