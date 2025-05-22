<script setup lang="ts">
import {Profile} from "@/types/profile";
import createApi from "@/api";
import {pError, pLoad, pSuccess, pWarning} from "@/util/pLoad";
import {useProxiesStore} from "@/store/proxiesStore";
import {useMenuStore} from "@/store/menuStore";
import {isHttpOrHttps, prettyBytes} from "@/util/format";
import {useI18n} from "vue-i18n";
import {Browser, Clipboard, Events} from "@/runtime"
import {useWebStore} from "@/store/webStore";
import {WS} from "@/util/ws";
import {onBeforeRouteLeave} from "vue-router";

const distanceFromTop = ref(195)
const upFromTop = function (distance: number) {
  distanceFromTop.value = distance
}

// i18n
const {t} = useI18n();

// 获取当前 Vue 实例的 proxy 对象
const {proxy} = getCurrentInstance()!;
const api = createApi(proxy);

// 当前页面使用store
const menuStore = useMenuStore();
const proxiesStore = useProxiesStore();
const webStore = useWebStore();

// 头部几个按钮操作
const addFormVisible = ref(false)
const isNowAdd = ref(false)
const addForm = reactive({
  content: '',
})

async function add() {
  if (!addForm.content) {
    return
  }

  isNowAdd.value = true
  const p = new Profile()
  p.content = addForm.content
  try {
    const pList = await api.addProfileFromInput(p)
    if (pList && pList.length > 0) {
      pList.forEach(item => profiles.push(item))
    }
    sendOrder(profiles)
    addForm.content = ""
    addFormVisible.value = false
  } catch (e) {
    if (e['message']) {
      pError(e['message'])
    }
  }
  isNowAdd.value = false
}

function handleAdd() {
  addForm.content = ""
  addFormVisible.value = true
}

function handlePaste() {
  addForm.content = Clipboard.Text()
  addFormVisible.value = true
}

function openFile() {
  webStore.dnd = true
}

// 头部显示
const headerShow = reactive({
  available: '',
  used: '',
  expire: '',
  update: '',
})

function setHeaderShow(item: any) {
  if (item['available']) {
    headerShow.available = prettyBytes(item['available'])
  } else {
    headerShow.available = ''
  }
  if (item['used']) {
    headerShow.used = prettyBytes(item['used'])
  } else {
    headerShow.used = ''
  }
  if (item['expire']) {
    headerShow.expire = item['expire']
  } else {
    headerShow.expire = ''
  }
  if (item['update']) {
    headerShow.update = item['update']
  } else {
    headerShow.update = ''
  }
}

// 列表显示
let profiles = reactive<any[]>([])

async function getProfileList() {
  if (profiles.length != 0) {
    profiles.splice(0, profiles.length)
  }
  const list = await api.getProfileList()
  if (list && list.length != 0) {
    list.forEach(item => {
      profiles.push(item)
      if (item['selected']) {
        setHeaderShow(item)
      }
    })

    Events.Emit({
      name: "profiles",
      data: list
    })

  }
}

// 拖动相关
const canDrag = ref(false)

function mouseEnter() {
  canDrag.value = true
}

function mouseLeave() {
  canDrag.value = false
}

// 切换订阅配置
async function switchProfile(data: any) {
  if (data['selected']) {
    return
  }

  await pLoad(t('profiles.switch.ing'), async () => {
    try {
      await api.switchProfile(data)
      proxiesStore.active = ""

      await api.waitRunning()

      for (let profile of profiles) {
        if (profile['selected']) {
          profile['selected'] = false
        }
      }
      data['selected'] = true
      setHeaderShow(data)

      api.getRuleNum().then((res) => {
        menuStore.setRuleNum(res);
      });

      Events.Emit({
        name: "profiles",
        data: toRaw(profiles)
      })

      // 关闭之前的连接
      api.closeAllConnection()

      pSuccess(t('profiles.switch.success'))
    } catch (e) {
      if (e['message']) {
        pError(e['message'])
      }
    }
  })

}


watch(() => webStore.fProfile, async (data: any) => {
  for (let profile of profiles) {
    if (profile['selected']) {
      profile['selected'] = false
    }
    if (profile['id'] == data['id']) {
      data = profile
    }
  }

  data['selected'] = true
  setHeaderShow(data)
})


// 更新订阅
async function refresh(data: any) {
  await pLoad(t('profiles.refresh.ing'), async () => {
    try {
      const re = await api.refreshProfile(data)
      if (data['selected']) {
        setHeaderShow(re)
      }
      Object.assign(data, re);
      pSuccess(t('profiles.refresh.success'))
    } catch (e) {
      if (e['message']) {
        pError(e['message'])
      }
    }
  })
}

// 几个按钮操作
// 到主页
function goHome(data: any) {
  Browser.OpenURL(data.home)
}

// 修改配置
const editFormVisible = ref(false)
let editForm = reactive<any>({})
let editFormD = {}

function updateProfile(data: any) {
  editFormD = data
  editForm = reactive<any>({})
  Object.assign(editForm, data)
  editFormVisible.value = true
}

function validateField(value: any) {
  // 如果为空，则通过校验
  if (value === "" || value === null || value === undefined) {
    return true;
  }

  // 如果不为空，验证是否是大于0且小于等于128的整数
  const regex = /^[1-9][0-9]?$|^1[0-2][0-8]$/;
  return regex.test(value.toString());
}

async function saveUpdateProfile() {

  switch (editForm.type) {
    case 2:
      if (!editForm.title) {
        pError(t('profiles.edit.title-tip'))
        return
      }
      break
    case 1:
      if (!editForm.title) {
        pError(t('profiles.edit.title-tip'))
        return
      }

      if (!editForm.content) {
        pError(t('profiles.edit.url-tip'))
        return
      }

      if (!isHttpOrHttps(editForm.content)) {
        pError(t('profiles.edit.url-error'))
        return
      }

      if (!validateField(editForm.interval)) {
        pError(t('profiles.edit.update-tip'))
        return
      }
  }


  await api.updateProfile(editForm)
  // 更新当前页面的值
  Object.assign(editFormD, editForm)
  editFormVisible.value = false
  pSuccess(t('profiles.edit.success'))

  Events.Emit({
    name: "profiles",
    data: toRaw(profiles)
  })
}

// 删除配置
async function deleteProfile(data: any, index: any) {
  if (data['selected']) {
    pWarning(t('profiles.del-tip'))
    return
  }

  try {
    await api.deleteProfile(data)
    profiles.splice(index, 1)
    Events.Emit({
      name: "profiles",
      data: toRaw(profiles)
    })
  } catch (e) {
    if (e['message']) {
      pError(e['message'])
    }
  }
}

// webSocket相关操作
let wsOrder: WS

function sendOrder(data: any) {
  if (wsOrder) {
    Events.Emit({
      name: "profiles",
      data: toRaw(data)
    })
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
  const urlTraffic = webStore.wsUrl + "/profile/order?token=" + webStore.secret;
  wsOrder = new WS(urlTraffic);

  await getProfileList()
})

watch(() => webStore.dProfile, async (pList) => {
  if (pList && pList.length > 0) {
    pList.forEach(item => profiles.push(item))
  }
})
</script>

<template>
  <MyLayout>
    <template #top>
      <el-space class="space">
        <div class="title">
          {{ $t('profiles.title') }}
        </div>
        <div class="profile-option">
          <el-tooltip
              :content="$t('profiles.add')"
              placement="top">
            <el-icon
                @click="handleAdd"
                class="profile-option-btn">
              <icon-mdi-plus-thick/>
            </el-icon>
          </el-tooltip>

          <el-tooltip
              :content="$t('profiles.paste')"
              placement="top">
            <el-icon
                @click="handlePaste"
                class="profile-option-btn">
              <icon-mdi-content-paste/>
            </el-icon>
          </el-tooltip>

          <el-tooltip
              :content="$t('profiles.open')"
              placement="top">
            <el-icon
                @click="openFile"
                class="profile-option-btn">
              <icon-mdi-folder-open/>
            </el-icon>
          </el-tooltip>
        </div>
      </el-space>

      <div class="sub-title">
        <template v-if="headerShow.available">
          <span>{{ $t('profiles.available') }} {{ headerShow.available }}</span>
          <el-divider direction="vertical" border-style="dashed"/>
        </template>
        <template v-if="headerShow.used">
          <span>{{ $t('profiles.use') }} {{ headerShow.used }}</span>
          <el-divider direction="vertical" border-style="dashed"/>
        </template>
        <template v-if="headerShow.expire">
          <span>{{ $t('profiles.expire') }} {{ headerShow.expire }}</span>
          <el-divider direction="vertical" border-style="dashed"/>
        </template>
        <template v-if="headerShow.update">
          <span>{{ $t('profiles.update') }} {{ headerShow.update }}</span>
        </template>
      </div>
    </template>

    <template #bottom>
      <VDContainer
          :data="profiles"
          @getData="sendOrder"
          :gap="15"
          :draggable="canDrag"
          style="margin-left: 10px;width: 95%;"
      >
        <template v-slot:VDC="{data,index}">
          <div
              :class="data.selected?'sub-card sub-card-select':'sub-card'"
              @click="switchProfile(data)"
          >
            <div class="row">
              <el-icon
                  @mouseenter.stop="mouseEnter"
                  @mouseleave.stop="mouseLeave"
                  size="22"
                  class="drag">
                <icon-mdi-drag/>
              </el-icon>

              <el-icon size="22"
                       v-if="data.type == 1"
                       class="ops"
                       @click.stop="refresh(data)">
                <icon-mdi-refresh/>
              </el-icon>

            </div>
            <div
                class="system-info"
            >
              <span :title="data.title">
                {{ data.title }}
              </span>
            </div>
            <div class="bottom-row">
              <el-tooltip
                  v-if="data.home"
                  :content="$t('profiles.home')"
                  placement="top">
                <el-icon
                    class="ops"
                    @click.stop="goHome(data)"
                    size="20">
                  <icon-mdi-home-import-outline/>
                </el-icon>
              </el-tooltip>
              <el-tooltip
                  :content="$t('edit')"
                  placement="top">
                <el-icon
                    class="ops"
                    @click.stop="updateProfile(data)"
                    size="20">
                  <icon-mdi-square-edit-outline/>
                </el-icon>
              </el-tooltip>
              <el-tooltip
                  :content="$t('delete')"
                  placement="top">
                <el-icon
                    class="ops"
                    @click.stop="deleteProfile(data,index)"
                    size="20">
                  <icon-mdi-trash-can/>
                </el-icon>
              </el-tooltip>
            </div>
          </div>
        </template>
      </VDContainer>

    </template>
  </MyLayout>

  <el-dialog v-model="addFormVisible"
             :title="t('profiles.add')"
             width="520"
             draggable
             center
  >
    <el-form :model="addForm">
      <el-form-item>
        <el-input
            :rows="3"
            type="textarea"
            autocapitalize="off"
            autocomplete="off"
            spellcheck="false"
            :placeholder="t('profiles.placeholder')"
            v-model="addForm.content"
        />
      </el-form-item>
    </el-form>
    <template #footer>
      <div class="dialog-footer">
        <el-button @click="addFormVisible = false">
          {{ t('cancel') }}
        </el-button>
        <el-button
            :loading="isNowAdd"
            type="primary"
            @click="add">
          {{ t('confirm') }}
        </el-button>
      </div>
    </template>
  </el-dialog>

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
          :label="t('profiles.edit.title')"
          label-width="120">
        <el-input
            v-model="editForm.title"
            clearable
            autocapitalize="off"
            autocomplete="off"
            spellcheck="false"/>
      </el-form-item>
      <el-form-item
          v-if="editForm.type == 1"
          :label="t('profiles.edit.url')"
          label-width="120">
        <el-input
            v-model="editForm.content"
            clearable
            autocapitalize="off"
            autocomplete="off"
            spellcheck="false"/>
      </el-form-item>
      <el-form-item
          v-if="editForm.type == 1"
          :label="t('profiles.edit.update')"
          label-width="120">
        <el-input
            v-model="editForm.interval"
            clearable
            autocapitalize="off"
            autocomplete="off"
            spellcheck="false">
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
.space {
  margin-top: 15px;
}

.title {
  font-size: 32px;
  font-weight: bold;
  margin-left: 10px;
}

.sub-title {
  margin-left: 10px;
  color: #FFD700;
  font-size: 14px;
  margin-top: 5px;
}

.profile-option {
  margin-left: 10px;
  font-size: 30px;
  padding-top: 10px;
}

.profile-option-btn {
  margin-right: 15px;
}

.profile-option-btn:hover {
  cursor: pointer;
  color: var(--hr-color);
}

:deep(.vdc-item-container) {
  width: calc(33% - 10px);
  max-width: 245px;
}

.sub-card {
  padding: 5px 8px 5px 5px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.1);
  color: white;
  box-shadow: var(--left-nav-shadow);
}

.sub-card:hover {
  cursor: pointer;
}

.sub-card-select {
  background-color: var(--left-item-selected-bg);
  box-shadow: var(--left-nav-hover-shadow);
  border: 2px solid var(--text-color);
  cursor: default;
}

.sub-card-select:hover {
  cursor: default;
}

.sub-card .row {
  display: flex;
  justify-content: space-between;
}

.sub-card .row .drag:hover {
  cursor: grab;
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
  padding: 5px 10px 2px 15px;
}

.bottom-row {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  margin-top: 5px;
  margin-bottom: 2px;
}


</style>