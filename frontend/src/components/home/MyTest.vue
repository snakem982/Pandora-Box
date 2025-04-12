<script setup lang="ts">

import {cJoin} from "@/util/format";

let configs = reactive([1, 2, 3, 4, 5, 6])



function handleEmit(value: any) {
  console.log(cJoin(value, ","))
}

function handleDelete(index: number) {
  configs.splice(index, 1);
}

</script>

<template>
  <el-row class="t-card" :gutter="20" style="margin-left: 12px">
    <el-col :span="24">
      <el-row>
        {{ $t('home.web') }}
        <el-icon size="22" style="margin-left: 8px;margin-top: -4px">
          <icon-mdi-refresh/>
        </el-icon>
        <el-icon size="22" style="margin-left: 8px;margin-top: -4px">
          <icon-mdi-link-edit/>
        </el-icon>
      </el-row>
      <hr>

      <FlatSortable>
        <FlatSortableContent
            class="icon-container"
            :gap="16"
            :vl="configs.length"
            @update:model-value="handleEmit"
            direction="row">
          <FlatSortableItem
              draggable
              class="icon-item"
              v-for="(item,index) in configs"
              :key="'ph-'+index">
            <div class="icon">
              <img
                  src="https://studiostaticassetsprod.azureedge.net/bundle-cmc/favicon.svg"
                  style="height: 48px;width: 48px;"
                  alt="C">
              <div class="overlay"></div>
              <!-- 添加删除按钮 -->
              <div class="delete-btn" @click="handleDelete(index)">×</div>
            </div>
            <div class="icon-title">
              {{ item }}
            </div>
            <el-tag type="success" class="icon-delay">
              200 ms
            </el-tag>
          </FlatSortableItem>
        </FlatSortableContent>
      </FlatSortable>


    </el-col>
  </el-row>
</template>

<style scoped>
/* 整体卡片样式 */
.t-card {
  width: calc(95% - 20px);
  margin-top: 30px;
  padding: 10px 0 10px 0;
  border-radius: 8px;
  text-align: left;
  box-shadow: var(--right-box-shadow);
}

/* 分割线样式 */
.t-card hr {
  border: none;
  height: 1px;
  background-color: var(--hr-color);
  margin: 10px 0;
}

/* 图标容器样式 */
.icon-container {
  display: flex;
  flex-wrap: wrap;
  justify-content: flex-start;
  gap: 10px;
  margin-top: 10px;
}

/* 单个图标和标题样式 */
.icon-item {
  text-align: center;
}

/* 图标样式 */
.icon {
  width: 60px;
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  font-weight: bold;
  background-color: var(--home-test-icon-bg);
  border-radius: 10px;
}

/* 图标标题样式 */
.icon-title {
  font-size: 13px;
  color: var(--text-color);
  margin-top: 5px;
}

.icon-delay {
  border-radius: 5px;
}

.overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: transparent;
  display: none;
  z-index: 100;
}

.icon:hover .overlay {
  display: block;
}

/* 删除按钮样式 */
.delete-btn {
  position: absolute;
  top: -5px;
  right: -5px;
  width: 20px;
  height: 20px;
  background-color: red;
  color: white;
  font-size: 14px;
  font-weight: bold;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  z-index: 200;
}

</style>