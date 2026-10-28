<script setup lang="ts">
defineProps<{
  isHide: boolean;
  isSort: boolean;
  isVertical: boolean;
}>();

const emit = defineEmits<{
  (event: "test"): void;
  (event: "locate"): void;
  (event: "toggle-hide"): void;
  (event: "toggle-sort"): void;
  (event: "toggle-vertical"): void;
}>();
</script>

<template>
  <el-space class="space">
    <div class="title">
      {{ $t("proxies.title") }}
    </div>
    <div class="proxy-option">
      <el-tooltip :content="$t('proxies.test')" placement="top">
        <el-icon class="proxy-option-btn" @click="emit('test')">
          <icon-mdi-speedometer />
        </el-icon>
      </el-tooltip>

      <el-tooltip :content="$t('proxies.locate')" placement="top">
        <el-icon class="proxy-option-btn" @click="emit('locate')">
          <icon-mdi-crosshairs-gps />
        </el-icon>
      </el-tooltip>

      <el-tooltip :content="isHide ? $t('proxies.hide-on') : $t('proxies.hide-off')" placement="top">
        <el-icon class="proxy-option-btn" @click="emit('toggle-hide')">
          <icon-mdi-eye-off v-if="isHide" />
          <icon-mdi-eye v-else />
        </el-icon>
      </el-tooltip>

      <el-tooltip :content="isSort ? $t('proxies.sort-on') : $t('proxies.sort-off')" placement="top">
        <el-icon class="proxy-option-btn" @click="emit('toggle-sort')">
          <icon-mdi-sort-ascending v-if="isSort" />
          <icon-mdi-sort v-else />
        </el-icon>
      </el-tooltip>

      <el-tooltip :content="isVertical ? $t('proxies.vertical-on') : $t('proxies.vertical-off')" placement="top">
        <el-icon class="proxy-option-btn" @click="emit('toggle-vertical')">
          <icon-mdi-arrow-expand-vertical v-if="isVertical" />
          <icon-mdi-arrow-expand-horizontal v-else />
        </el-icon>
      </el-tooltip>
    </div>
  </el-space>
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
</style>