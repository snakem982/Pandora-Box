<!-- src/views/profiles/Profiles.vue -->
<script setup lang="ts">
import {onBeforeUnmount, onMounted, ref} from 'vue'
import {Clipboard} from '@/runtime'
import {DropdownInstance} from 'element-plus'
import {Delete, Edit, EditPen, HomeFilled, RefreshRight, Service} from '@element-plus/icons-vue'

import {useProfiles} from './profiles/useProfiles'
import ProfileHeader from './profiles/ProfileHeader.vue'
import ProfileCard from './profiles/ProfileCard.vue'
import ProfileDialogs from './profiles/ProfileDialogs.vue'

import {useRouter} from "vue-router";

const router = useRouter()

const {
  api, profiles, headerShow, tList, profileStore,
  webStore, sendOrder, getProfileList, switchProfile, deleteProfile,
  refreshProfile, goHome, goSupport
} = useProfiles()

const dialogsRef = ref<InstanceType<typeof ProfileDialogs>>()
const canDrag = ref(false)

// 右键菜单相关
const dropdownRef = ref<DropdownInstance>()
const cr_home = ref(false)
const cr_support = ref(false)
const position = ref({top: 0, left: 0, bottom: 0, right: 0} as DOMRect)
const triggerRef = ref({getBoundingClientRect: () => position.value})


const handleContextmenu = (event: MouseEvent, data: any) => {
  profileStore.profileConfig = data
  cr_home.value = !!data?.home;
  cr_support.value = !!data?.support;
  position.value = DOMRect.fromRect({x: event.clientX, y: event.clientY})
  event.preventDefault()
  dropdownRef.value?.handleOpen()
}

const handleCommand = (command: string) => {
  dropdownRef.value?.handleClose()
  const data = profileStore.profileConfig
  if (!data) return

  switch (command) {
    case 'refresh':
      refreshProfile(data)
      break
    case 'edit':
      dialogsRef.value?.openEdit(data)
      break
    case 'editConfig':
      router.push('/Profiles/Config');
      break
    case 'home':
      goHome(data)
      break
    case 'support':
      goSupport(data)
      break
    case 'delete':
      const idx = profiles.findIndex((p) => p.id === data.id)
      if (idx !== -1) deleteProfile(data, idx)
      break
    default:
      return;
  }
}

// 头部事件触发
const handleAdd = () => dialogsRef.value?.openAdd()
const handlePaste = () => dialogsRef.value?.openAdd(Clipboard.Text())
const openFile = () => {
  webStore.dnd = true
  webStore.dSelect = true
}

// 监听 Deeplink 导入事件
function handleProfilesImported(event: Event) {
  const detail = (event as CustomEvent).detail
  if (!detail?.profiles || !Array.isArray(detail.profiles)) return

  let added = false
  for (const item of detail.profiles) {
    if (item && !profiles.some((p) => p.id === item.id)) {
      profiles.push(item)
      added = true
    }
  }
  if (added) sendOrder(profiles)
}

onMounted(() => {
  window.addEventListener('deeplink-profile-imported', handleProfilesImported as EventListener)
})
onBeforeUnmount(() => {
  window.removeEventListener('deeplink-profile-imported', handleProfilesImported as EventListener)
})


</script>

<template>
  <MyLayout>
    <template #top>
      <ProfileHeader
          :header-show="headerShow"
          @add="handleAdd"
          @paste="handlePaste"
          @openFile="openFile"
      />
    </template>

    <template #bottom>
      <VDContainer
          :data="profiles"
          @getData="sendOrder"
          :gap="15"
          :draggable="canDrag"
          style="margin-left: 10px; width: 95%"
      >
        <template v-slot:VDC="{ data }">
          <ProfileCard
              :data="data"
              @switch="switchProfile"
              @contextmenu="handleContextmenu"
              @dragStart="canDrag = true"
              @dragEnd="canDrag = false"
          />
        </template>
      </VDContainer>
    </template>
  </MyLayout>

  <!-- 弹窗统一管理组件 -->
  <ProfileDialogs
      ref="dialogsRef"
      :api="api"
      :tList="tList"
      :profiles="profiles"
      @refreshList="getProfileList"
  />

  <!-- 上下文右键菜单 -->
  <el-dropdown
      ref="dropdownRef"
      :virtual-ref="triggerRef"
      :show-arrow="false"
      virtual-triggering
      trigger="contextmenu"
      placement="bottom-start"
      size="large"
      @command="handleCommand"
  >
    <template #dropdown>
      <el-dropdown-menu>
        <el-dropdown-item :icon="RefreshRight" command="refresh">{{ $t('refresh') }}</el-dropdown-item>
        <el-dropdown-item :icon="Edit" command="edit">{{ $t('edit') }}</el-dropdown-item>
        <el-dropdown-item :icon="EditPen" command="editConfig">{{ $t('modify') }}</el-dropdown-item>
        <el-dropdown-item :icon="HomeFilled" command="home" v-if="cr_home">
          {{ $t('profiles.home') }}
        </el-dropdown-item>
        <el-dropdown-item :icon="Service" command="support" v-if="cr_support">
          {{ $t('profiles.support') }}
        </el-dropdown-item>
        <el-dropdown-item :icon="Delete" divided command="delete">{{ $t('delete') }}</el-dropdown-item>
      </el-dropdown-menu>
    </template>
  </el-dropdown>

</template>

<style scoped>
:deep(.vdc-item-container) {
  width: calc(33% - 10px);
  max-width: 245px;
}
</style>