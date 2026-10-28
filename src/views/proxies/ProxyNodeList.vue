<script setup lang="ts">
import {ref} from "vue";

defineProps<{
  nodes: any[];
  isDirect: boolean;
  directLabel: string;
}>();

const emit = defineEmits<{
  (event: "select", now: boolean, name: string): void;
}>();

const nodesContainer = ref<HTMLElement | null>(null);

function locateSelected() {
  nodesContainer.value
    ?.querySelector<HTMLElement>(".proxy-node-select")
    ?.scrollIntoView({behavior: "smooth", block: "center", inline: "nearest"});
}

defineExpose({locateSelected});
</script>

<template>
  <div v-if="!isDirect" ref="nodesContainer" class="proxy-nodes">
    <div
      v-for="node in nodes"
      :key="node['name']"
      :class="node['now'] ? 'proxy-nodes-card proxy-node-select' : 'proxy-nodes-card'"
      @click="emit('select', node['now'], node['name'])"
    >
      <div class="proxy-nodes-title">
        <span :title="node['name']">{{ node["name"] }}</span>
      </div>
      <div class="proxy-nodes-tags">
        <span class="proxy-nodes-tags-left">{{ node["type"] }}</span>
        <span :class="'proxy-nodes-tags-right ' + node['toClass']">
          {{ node["delay"] }} ms
        </span>
      </div>
    </div>
  </div>

  <div v-else class="proxy-nodes">
    {{ directLabel }}
  </div>
</template>

<style scoped>
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
  border: 2px solid var(--sub-card-border);
  border-radius: 8px;
  padding: 8px 12px;
  background: var(--sub-card-bg);
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  line-height: 1.3;
  box-shadow: var(--left-nav-shadow);
  margin-top: 3px;
}

.proxy-nodes-card:hover,
.proxy-node-select {
  background-color: var(--left-item-selected-bg);
  border: 2px solid var(--text-color);
  cursor: pointer;
}

.proxy-node-select:hover {
  cursor: default;
}

.proxy-nodes-title {
  font-size: 16px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.proxy-nodes-tags {
  font-size: 14px;
  display: flex;
  margin-top: 12px;
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
</style>