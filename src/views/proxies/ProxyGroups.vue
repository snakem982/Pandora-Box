<script setup lang="ts">
import {nextTick, onBeforeUnmount, onMounted, ref, watch} from "vue";

const props = defineProps<{
  groups: string[];
  active: string;
  isVertical: boolean;
  rule: string;
}>();

const emit = defineEmits<{
  (event: "select", value: string): void;
}>();

const proxyGroup = ref<HTMLElement | null>(null);
const atStart = ref(true);
const atEnd = ref(true);
const isDropdownOpen = ref(false);
let scrollTimer: ReturnType<typeof setTimeout> | undefined;
let dropdownTimer: ReturnType<typeof setTimeout> | undefined;

const updateButtonVisibility = () => {
  if (!proxyGroup.value) return;
  const {scrollLeft, scrollWidth, clientWidth} = proxyGroup.value;
  atStart.value = scrollLeft === 0;
  atEnd.value = scrollLeft + clientWidth >= scrollWidth;
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

const handleScroll = () => {
  clearTimeout(scrollTimer);
  scrollTimer = setTimeout(updateButtonVisibility, 200);
};

const hideDropdown = () => {
  dropdownTimer = setTimeout(() => {
    isDropdownOpen.value = false;
  }, 200);
};

const enterDropDown = () => {
  clearTimeout(dropdownTimer);
  isDropdownOpen.value = true;
};

watch(() => [props.groups, props.isVertical], async () => {
  await nextTick();
  atStart.value = true;
  atEnd.value = true;
  updateButtonVisibility();
});

onMounted(() => {
  updateButtonVisibility();
  window.addEventListener("resize", updateButtonVisibility);
});

onBeforeUnmount(() => {
  window.removeEventListener("resize", updateButtonVisibility);
  clearTimeout(scrollTimer);
  clearTimeout(dropdownTimer);
});
</script>

<template>
  <div
    v-if="isVertical && rule !== 'direct' && groups.length > 0"
    class="dropdown"
  >
    <button class="dropdown-btn" @mouseenter="enterDropDown" @mouseleave="hideDropdown">
      {{ active }}
    </button>
    <ul
      v-if="isDropdownOpen"
      class="dropdown-list"
      @mouseenter="enterDropDown"
      @mouseleave="hideDropdown"
    >
      <li
        v-for="item in groups"
        :key="item + '-gv'"
        class="dropdown-item"
        @click="emit('select', item)"
      >
        {{ item }}
      </li>
    </ul>
  </div>

  <div
    v-if="!isVertical && rule !== 'direct' && groups.length > 0"
    class="button-container"
  >
    <el-icon v-if="!atStart" class="scroll-left" @click="scrollLeft">
      <icon-mdi-arrow-expand-left />
    </el-icon>
    <div ref="proxyGroup" class="proxy-group" @scroll="handleScroll">
      <button
        v-for="item in groups"
        :key="item + '-g'"
        :class="active === item ? 'proxy-group-title proxy-group-title-select' : 'proxy-group-title'"
        @click="emit('select', item)"
      >
        {{ item }}
      </button>
    </div>
    <el-icon v-if="!atEnd" class="scroll-right" @click="scrollRight">
      <icon-mdi-arrow-expand-right />
    </el-icon>
  </div>
</template>

<style scoped>
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
  overflow-x: auto;
  -webkit-overflow-scrolling: touch;
  scroll-behavior: smooth;
  scrollbar-width: none;
  -ms-overflow-style: none;
  overscroll-behavior-x: contain;
  scroll-snap-type: x proximity;
}

.proxy-group::-webkit-scrollbar {
  display: none;
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
  font-family: inherit;
}

.proxy-group-title:hover,
.proxy-group-title-select {
  background-color: var(--left-item-selected-bg);
  box-shadow: var(--left-nav-hover-shadow);
  border-color: var(--text-color);
}

.proxy-group-title-select:hover {
  cursor: default;
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
  background: var(--skin-bg-color);
  border: 2px solid var(--text-color);
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
  color: var(--text-color);
  padding: 8px;
  cursor: pointer;
}

.dropdown-item:hover {
  background: var(--skin-hover-color);
}
</style>