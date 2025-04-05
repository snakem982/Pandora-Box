<script setup lang="ts">

import MyHr from "@/components/proxies/MyHr.vue";

const eye = ref(false)
const sort = ref(true)

const nodes = reactive([
  {
    name: "🇯🇵 JP_01",
    kind: "Hysteria2",
    delay: 123
  },
  {
    name: "🇺🇸 US_02",
    kind: "ShadowSocks",
    delay: 423
  },
  {
    name: "🇰🇷韩国7-高速专线-hy2",
    kind: "Tuic",
    delay: 23
  },
  {
    name: "🇰🇷 韩国IEPL",
    kind: "Vless",
    delay: 123
  },
  {
    name: "🇸🇬 新加坡IEPL 专线 02",
    kind: "AnyTls",
    delay: 123
  },
  {
    name: " 🇯🇵 JP_11",
    kind: "Hysteria2",
    delay: 123
  },
  {
    name: " 🇯🇵 JP_11",
    kind: "Vless",
    delay: 123
  },

])

const distanceFromTop = ref(195)
const upFromTop = function (distance: number) {
  distanceFromTop.value = distance
}

</script>

<template>
  <MyLayout :top-height="distanceFromTop-15" :bottom-height="distanceFromTop+25">
    <template #top>
      <MySearch></MySearch>
      <el-space>
        <div class="title">代理</div>
        <div class="proxy-option">
          <el-tooltip
              content="测试延迟"
              placement="top">
            <el-icon class="proxy-option-btn">
              <icon-mdi-speedometer/>
            </el-icon>
          </el-tooltip>

          <el-tooltip
              :content="eye?'隐藏不可用节点':'显示不可用节点'"
              placement="top">
            <el-icon class="proxy-option-btn">
              <icon-mdi-eye v-if="eye"/>
              <icon-mdi-eye-off v-else/>
            </el-icon>
          </el-tooltip>

          <el-tooltip
              :content="sort?'按延迟排序':'默认排序'"
              placement="top">
            <el-icon class="proxy-option-btn">
              <icon-mdi-sort v-if="sort"/>
              <icon-mdi-sort-ascending v-else/>
            </el-icon>
          </el-tooltip>
        </div>
      </el-space>

      <div class="proxy-group">
        <button class="proxy-group-title">🚀 节点选择</button>
        <button class="proxy-group-title proxy-group-title-select">♻️ 自动选择</button>
        <button class="proxy-group-title">🎯 全球直连</button>
        <button class="proxy-group-title">🛑 全球拦截</button>
        <button class="proxy-group-title">🐟 漏网之鱼</button>
      </div>

      <MyHr :update="upFromTop"></MyHr>
    </template>
    <template #bottom>
      <div class="proxy-nodes">
        <div class="proxy-nodes-card" v-for="(node, index) in nodes" :key="index">
          <div class="proxy-nodes-title">
            {{ node.name }}
          </div>
          <div class="proxy-nodes-tags">
        <span class="proxy-nodes-tags-left">
          {{ node.kind }}
        </span>
            <span class="proxy-nodes-tags-right">
          {{ node.delay }} ms
        </span>
          </div>
        </div>
      </div>
    </template>
  </MyLayout>
</template>

<style scoped>
.title {
  margin-top: 25px;
  font-size: 32px;
  font-weight: bold;
  margin-left: 10px;
}

.proxy-option {
  margin-left: 10px;
  font-size: 30px;
  margin-top: 25px;
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
  margin-top: 14px;
  margin-left: 10px;
  flex-wrap: wrap;
  gap: 10px;
  width: 95%;
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
  margin-top: 15px;
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


</style>