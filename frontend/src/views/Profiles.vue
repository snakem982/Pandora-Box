<script setup lang="ts">
import {Profile} from "@/types/profile";
import createApi from "@/api";
import {error, pLoad, success} from "@/util/pLoad";
import {useProxiesStore} from "@/store/proxiesStore";
import {useMenuStore} from "@/store/menuStore";
import {prettyBytes} from "@/util/format";
import {useI18n} from "vue-i18n";
import {Clipboard} from "@wailsio/runtime"
import {useWebStore} from "@/store/webStore";

const {t} = useI18n();

// 获取当前 Vue 实例的 proxy 对象
const {proxy} = getCurrentInstance()!;
const api = createApi(proxy);

// 当前页面使用store
const menuStore = useMenuStore();
const proxiesStore = useProxiesStore();
const webStore = useWebStore();

let profiles = reactive<any[]>([])
const profileInfo = reactive({
  available: '',
  used: '',
  expire: '',
  update: '',
})

function setProfile(item: any) {
  if (item['available']) {
    profileInfo.available = prettyBytes(item['available'])
  } else {
    profileInfo.available = ''
  }
  if (item['used']) {
    profileInfo.used = prettyBytes(item['used'])
  } else {
    profileInfo.used = ''
  }
  if (item['expire']) {
    profileInfo.expire = item['expire']
  } else {
    profileInfo.expire = ''
  }
  if (item['update']) {
    profileInfo.update = item['update']
  } else {
    profileInfo.update = ''
  }
}

async function getProfileList() {
  if (profiles.length != 0) {
    profiles.splice(0, profiles.length)
  }
  const list = await api.getProfileList()
  list.forEach(item => {
    profiles.push(item)
    if (item['selected']) {
      setProfile(item)
    }
  })
}

onMounted(async () => {
  await getProfileList()
})

const canDrag = ref(false)

function mouseEnter() {
  canDrag.value = true
}
function mouseLeave() {
  canDrag.value = false
}


function handleEmit(value: any) {
  console.log(profiles)
}

const dialogFormVisible = ref(false)
const profile = reactive({
  content: '',
})

const isNowAdd = ref(false)

async function add() {
  if (!profile.content) {
    return
  }

  isNowAdd.value = true
  const p = new Profile()
  p.content = profile.content
  try {
    const pList = await api.addProfileFromInput(p)
    if (pList && pList.length > 0) {
      pList.forEach(item => profiles.push(item))
    }
    profile.content = ""
    dialogFormVisible.value = false
  } catch (e) {
    isNowAdd.value = false
    if (e['message']) {
      error(e['message'])
    }
  }
}

watch(() => webStore.dProfile, async (pList, oldValue) => {
  if (pList && pList.length > 0) {
    pList.forEach(item => profiles.push(item))
  }
})

async function switchProfile(data: any) {
  if (data['selected']) {
    return
  }

  try {
    await api.switchProfile(data)
    proxiesStore.active = ""

    for (let profile of profiles) {
      if (profile['selected']) {
        profile['selected'] = false
        break
      }
    }
    data['selected'] = true
    setProfile(data)

    api.getRules().then((res) => {
      menuStore.setRuleNum(res.length);
    });

  } catch (e) {
    if (e['message']) {
      error(e['message'])
    }
  }

}

async function refresh(data: any) {
  await pLoad(t('proxies.refresh.ing'), async () => {
    try {
      await api.refreshProfile(data)
      setProfile(data)
      success(t('proxies.refresh.success'))
    } catch (e) {
      if (e['message']) {
        error(e['message'])
      }
    }
  })
}

function handlePaste() {
  Clipboard.Text().then(text => {
    profile.content = text
    dialogFormVisible.value = true
  })
}

function openFile() {
  webStore.dnd = true
}

function deleteProfile(data: any) {

}


</script>

<template>
  <MyLayout :top-height="150" :bottom-height="180">
    <template #top>
      <MySearch></MySearch>
      <el-space class="space">
        <div class="title">
          {{ $t('profiles.title') }}
        </div>
        <div class="profile-option">
          <el-tooltip
              :content="$t('profiles.add')"
              placement="top">
            <el-icon
                @click="dialogFormVisible = true"
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
        <template v-if="profileInfo.available">
          <span>{{ $t('profiles.available') }} {{ profileInfo.available }}</span>
          <el-divider direction="vertical" border-style="dashed"/>
        </template>
        <template v-if="profileInfo.used">
          <span>{{ $t('profiles.use') }} {{ profileInfo.used }}</span>
          <el-divider direction="vertical" border-style="dashed"/>
        </template>
        <template v-if="profileInfo.expire">
          <span>{{ $t('profiles.expire') }} {{ profileInfo.expire }}</span>
          <el-divider direction="vertical" border-style="dashed"/>
        </template>
        <template v-if="profileInfo.update">
          <span>{{ $t('profiles.update') }} {{ profileInfo.update }}</span>
        </template>
      </div>
    </template>
    <template #bottom>

      <VDContainer
          :data="profiles"
          @getData="handleEmit"
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
                       class="refresh"
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
              <el-icon size="20">
                <icon-mdi-cog/>
              </el-icon>
              <el-icon size="20">
                <icon-mdi-trash-can/>
              </el-icon>
            </div>
          </div>
        </template>
      </VDContainer>

    </template>
  </MyLayout>

  <el-dialog v-model="dialogFormVisible"
             :title="t('profiles.add')"
             width="600"
             draggable
  >
    <el-form :model="profile">
      <el-form-item>
        <el-input
            type="textarea"
            :placeholder="t('profiles.placeholder')"
            v-model="profile.content"
            autocomplete="off"/>
      </el-form-item>
    </el-form>
    <template #footer>
      <div class="dialog-footer">
        <el-button @click="dialogFormVisible = false">
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

.sub-card .row .refresh:hover {
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