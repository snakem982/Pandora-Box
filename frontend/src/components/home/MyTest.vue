<script setup lang="ts">

import {WebTest} from "@/types/webtest";
import createApi from "@/api";
import {useI18n} from "vue-i18n";
import {useMenuStore} from "@/store/menuStore";
import {useWebStore} from "@/store/webStore";
import {WS} from "@/util/ws";
import {onBeforeRouteLeave} from "vue-router";
import {pError, pSuccess} from "@/util/pLoad";


// i18n
const {t} = useI18n();

// 获取当前 Vue 实例的 proxy 对象
const {proxy} = getCurrentInstance()!;
const api = createApi(proxy);

// 当前页面使用store
const menuStore = useMenuStore();
const webStore = useWebStore();

// 编辑相关
const editShow = ref(false)

function handleDelete(data: any, index: number) {
  api.deleteWebTest(data)
  webTestList.splice(index, 1);
}


// 修改配置
const editFormVisible = ref(false)
let editForm = reactive<any>({})

function handleEdit(data: any) {
  editForm = data
  editFormVisible.value = true
}

async function saveUpdateProfile() {


  if (!editForm.title) {
    pError(t('profiles.edit.title-tip'))
    return
  }


  // await api.updateProfile(editForm)
  editFormVisible.value = false
  pSuccess(t('profiles.edit.success'))
}

// 列表显示
let webTestList = reactive<WebTest[]>([])

async function getWebTestList() {
  if (webTestList.length != 0) {
    webTestList.splice(0, webTestList.length)
  }
  const list = await api.getWebTest()
  list.forEach(item => {
    webTestList.push(item)
  })
}

// 保存排序
// webSocket相关操作
let wsOrder: WS

function sendOrder(data: any) {
  if (wsOrder) {
    wsOrder.send(JSON.stringify(data))
  }
}

// 路由切换前关闭 WebSocket
onBeforeRouteLeave(() => {
  wsOrder.close();
});
onBeforeUnmount(() => {
  wsOrder.close();
})

// vue 周期相关
onMounted(async () => {
  const urlTraffic = webStore.wsUrl + "/webtest/order?token=" + webStore.secret;
  wsOrder = new WS(urlTraffic);
  await getWebTestList()
})

</script>

<template>
  <el-row class="t-card" :gutter="20" style="margin-left: 12px">
    <el-col :span="24">
      <el-row>
        {{ $t('home.web.title') }}
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
          :data="webTestList"
          @getData="sendOrder"
          :gap="10"
          :top="8"
          draggable
      >
        <template v-slot:VDC="{data,index}">
          <div class="icon-item">
            <div class="icon">
              <img
                  draggable="false"
                  :src="data.src"
                  style="height: 48px;width: 48px;"
                  alt="C">
              <template v-if="editShow">
                <div class="delete-btn" @click="handleDelete(data,index)">
                  <icon-mdi-close/>
                </div>
                <div class="edit-btn" @click="handleEdit(data)">
                  <icon-mdi-pencil/>
                </div>
              </template>
            </div>
            <div
                class="icon-title"
                :title="data.title"
            >
              {{ data.title }}
            </div>
            <el-tag
                v-if="data.delay > 0"
                type="success"
                class="icon-delay">
              {{ data.delay }}
            </el-tag>
            <el-tag
                v-if="data.delay == -1"
                type="danger"
                class="icon-delay">
              {{ t('home.timeout') }}
            </el-tag>
          </div>
        </template>
      </VDContainer>
    </el-col>
  </el-row>


  <el-dialog v-model="editFormVisible"
             :title="t('edit')"
             width="520"
             draggable
             center
  >
    <el-form
        :model="editForm"
        label-position="top"
    >
      <el-form-item
          :label="t('home.web.edit')"
          label-width="120">
        <el-input
            v-model="editForm.title"
            clearable
            autocomplete="off"/>
      </el-form-item>
      <el-form-item
          :label="t('home.web.src')"
          label-width="120">
        <el-input
            v-model="editForm.src"
            clearable
            autocomplete="off"/>
      </el-form-item>
      <el-form-item
          :label="t('home.web.test')"
          label-width="120">
        <el-input
            v-model="editForm.testUrl"
            clearable
            autocomplete="off">
        </el-input>
      </el-form-item>

    </el-form>
    <template #footer>
      <div class="dialog-footer">
        <el-button @click="editFormVisible = false">
          {{ t('cancel') }}
        </el-button>
        <el-button type="primary"
                   @click="saveUpdateProfile">
          {{ t('confirm') }}
        </el-button>
      </div>
    </template>
  </el-dialog>


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