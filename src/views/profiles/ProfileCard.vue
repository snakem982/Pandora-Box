<!-- src/views/profiles/components/ProfileCard.vue -->
<script setup lang="ts">
import {prettyBytes} from '@/util/format'

const props = defineProps<{
  data: any
}>()

const emit = defineEmits(['switch', 'contextmenu', 'dragStart', 'dragEnd'])

function getRandomChar(flag: number) {
  let chars = ['-', '^', '>'];
  if (flag === 2) {
    chars = ['^', '-', '<']
  }
  const randomIndex = Math.floor(Math.random() * chars.length);
  return chars[randomIndex];
}
</script>

<template>
  <div
      :class="['sub-card', { 'sub-card-select': data.selected }]"
      @click="emit('switch', data)"
      @contextmenu.prevent="(e) => emit('contextmenu', e, data)"
  >
    <div class="row">
      <span class="s-left">
        <el-icon
            @mouseenter.stop="emit('dragStart')"
            @mouseleave.stop="emit('dragEnd')"
            @contextmenu.stop
            size="25"
            class="drag"
        >
          <icon-mdi-drag/>
        </el-icon>
      </span>
      <span class="s-title" :title="data.title">{{ data.title }}</span>
      <span class="s-right">
        <el-icon size="25" class="ops" @click.stop="(e: any) => emit('contextmenu', e, data)">
          <icon-mdi-dots-horizontal/>
        </el-icon>
      </span>
    </div>

    <div class="system-info">
      {{ data.used ? prettyBytes(data.used) : getRandomChar(1) }} {{ data.total || data.used ? '/' : '_' }}
      {{ data.total ? prettyBytes(data.total) : getRandomChar(2) }}
    </div>

    <div class="system-info">
      {{ data.update }}
    </div>
  </div>
</template>

<style scoped>
.sub-card {
  padding: 5px 8px 5px 5px;
  border: 2px solid var(--sub-card-border);
  border-radius: 8px;
  background: var(--sub-card-bg);
  color: var(--text-color);
  box-shadow: var(--left-nav-shadow);
  margin-top: 5px;
}

.sub-card-select {
  background-color: var(--left-item-selected-bg);
  border: 2px solid var(--text-color);
  cursor: pointer;
}

.sub-card-select:hover {
  cursor: default;
}

.profile-card-flash {
  animation: profile-card-flash 0.35s ease-in-out 0.2s 2;
}

.profile-card-flash * {
  animation: profile-card-text-flash 0.35s ease-in-out 0.2s 2;
}

@keyframes profile-card-flash {
  50% {
    background-color: var(--text-color);
  }
}

@keyframes profile-card-text-flash {
  50% {
    color: var(--left-item-selected-bg);
  }
}

.sub-card .row {
  display: flex;
  align-items: center;
  width: 100%;
}

.sub-card .row .drag:hover {
  cursor: grab;
}

.sub-card .row .s-left {
  flex-shrink: 0;
}

.sub-card .row .s-title {
  flex: 1;
  min-width: 0;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  font-size: 16px;
  color: var(--text-color);
  margin-left: 3px;
  margin-right: 2px;
}

.sub-card .row .s-right {
  flex-shrink: 0;
}

.ops:hover {
  cursor: pointer;
}

.system-info {
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
  text-align: left;
  font-size: 14px;
  padding: 7px 0 5px 26px;
  color: var(--text-color);
}
</style>