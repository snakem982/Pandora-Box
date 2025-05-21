<template>
  <div v-if="isWindows">
    <el-tooltip
        :content="$t('minus')"
        placement="bottom">
      <span class="bar" @click="minus2tray">
          <el-icon>
              <icon-mdi-card-minus-outline/>
          </el-icon>
      </span>
    </el-tooltip>
    <el-tooltip
        :content="$t('mini')"
        placement="bottom">
      <span class="bar" @click="minus">
          <el-icon>
              <icon-mdi-minus/>
          </el-icon>
      </span>
    </el-tooltip>
    <el-tooltip
        v-if="isMaximized"
        :content="$t('restore')"
        placement="bottom">
      <span class="bar" @click="max">
          <el-icon>
              <icon-mdi-window-restore/>
          </el-icon>
      </span>
    </el-tooltip>
    <el-tooltip
        v-else
        :content="$t('max')"
        placement="bottom">
      <span class="bar" @click="max">
          <el-icon>
              <icon-mdi-window-maximize/>
          </el-icon>
      </span>
    </el-tooltip>
    <el-tooltip
        :content="$t('close')"
        placement="bottom">
      <span class="" @click="close">
          <el-icon>
              <icon-mdi-window-close/>
          </el-icon>
      </span>
    </el-tooltip>
  </div>
  <div v-else>
    <el-tooltip
        :content="$t('minus')"
        placement="left">
      <span class="" @click="minus2tray">
          <el-icon>
              <icon-mdi-card-minus-outline/>
          </el-icon>
      </span>
    </el-tooltip>
  </div>
</template>

<script setup lang="ts">
import {Events} from "@/runtime";

const isMaximized = ref(false)
const isWindows = ref(false)

onMounted(() => {
  // @ts-ignore
  if (window["pxShowBar"]) {
    isWindows.value = true;
  }
})

function close() {
  Events.Emit({name: "close", data: true});
}

function minus() {
  Events.Emit({name: "min", data: true});
}

function max() {
  isMaximized.value = !isMaximized.value;
  Events.Emit({name: "max", data: true});
}

// 最小化到托盘
function minus2tray() {
  Events.Emit({name: "hide", data: true});
}
</script>

<style scoped>
.bar {
  margin-right: 15px;
}
</style>
