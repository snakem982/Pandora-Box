<script setup lang="ts">

import {cJoin} from "@/util/format";

let configs = reactive([1, 2, 3, 4, 5, 6, 7, 8, 9, 10]);


function getData() {
  console.log(cJoin(configs, ","))
}

function handleDelete(index: number) {
  configs.splice(index, 1);
  configs.push(index + 10);
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


      <VDContainer
          :data="configs"
          @getData="getData"
          :gap="10"
          :top="8"
          draggable
      >
        <template v-slot:VDC="{data,index}">
          <div class="icon-item">
            <div class="icon">
              <img
                  draggable="false"
                  src="https://studiostaticassetsprod.azureedge.net/bundle-cmc/favicon.svg"
                  style="height: 48px;width: 48px;"
                  alt="C">
              <div class="delete-btn" @click="handleDelete(index)">×</div>
            </div>
            <div class="icon-title">
              {{ data }}
            </div>
            <el-tag type="success" class="icon-delay">
              200 ms
            </el-tag>
          </div>
        </template>
      </VDContainer>
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

/* 删除按钮样式 */
.delete-btn {
  position: absolute;
  margin-top: -50px;
  margin-left: 50px;
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