<script setup lang="ts">

import {cJoin} from "@/util/format";

const editShow = ref(false)

let configs = reactive(["google", "youtube", "chatgptsdfsdfsfsfd"]);


function getData() {
  console.log(cJoin(configs, ","))
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
        <el-tooltip
            :content="$t('refresh')"
            placement="top">
          <el-icon size="22" class="tip">
            <icon-mdi-refresh/>
          </el-icon>
        </el-tooltip>
        <el-tooltip
            :content="$t('add')"
            placement="top">
          <el-icon size="22" class="tip">
            <icon-mdi-plus-thick/>
          </el-icon>
        </el-tooltip>
        <el-tooltip
            :content="$t('edit')"
            placement="top">
          <el-icon
              @click="editShow=!editShow"
              size="22"
              class="tip">
            <icon-mdi-link-edit/>
          </el-icon>
        </el-tooltip>
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
              <template v-if="editShow">
                <div class="delete-btn" @click="handleDelete(index)">
                  <icon-mdi-close/>
                </div>
                <div class="edit-btn" @click="handleDelete(index)">
                  <icon-mdi-pencil/>
                </div>
              </template>
            </div>
            <div
                class="icon-title"
                :title="data"
            >
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

.tip {
  margin-left: 8px;
  margin-top: -4px
}

.tip:hover {
  color: #cccccc;
  cursor: pointer;
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
  width: 60px;
  text-overflow: ellipsis;
  overflow: hidden;
}

.icon-delay {
  border-radius: 5px;
  margin-top: 5px;
}

/* 删除按钮样式 */
.delete-btn {
  position: absolute;
  margin-top: -50px;
  margin-left: 60px;
  width: 17px;
  height: 17px;
  background-color: red;
  color: var(--text-color);
  font-size: 15px;
  border-radius: 50%;
  display: flex;
  justify-content: center; /* 水平居中 */
  align-items: center; /* 垂直居中 */
  cursor: pointer;
  z-index: 200;
}

/* 编辑按钮样式 */
.edit-btn {
  position: absolute;
  margin-top: -10px;
  margin-left: 60px;
  width: 17px;
  height: 17px;
  background-color: blue;
  color: var(--text-color);
  font-size: 9px;
  border-radius: 50%;
  display: flex;
  justify-content: center; /* 水平居中 */
  align-items: center; /* 垂直居中 */
  cursor: pointer;
  z-index: 200;
}

</style>