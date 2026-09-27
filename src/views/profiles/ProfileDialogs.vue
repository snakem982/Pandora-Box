<!-- src/views/profiles/components/ProfileDialogs.vue -->
<script setup lang="ts">
import {reactive, ref} from 'vue'
import {useI18n} from 'vue-i18n'
import {Profile} from '@/types/profile'
import {pError, pSuccess} from '@/util/pLoad'
import {getTemplateTitle, isHttpOrHttps} from '@/util/format'
import {Events} from '@/runtime'
import {useMenuStore} from "@/store/menuStore";

const props = defineProps<{
  api: any
  tList: any[]
  profiles: any[]
}>()

const emit = defineEmits(['refreshList', 'updateSuccess'])
const {t} = useI18n()

const menuStore = useMenuStore()

// 新增相关
const addVisible = ref(false)
const isNowAdd = ref(false)
const addForm = reactive({content: ''})

// 编辑相关
const editVisible = ref(false)
const isNowEdit = ref(false)
let editForm = reactive<any>({})
let editFormOrigin: any = null

function openAdd(content = '') {
  addForm.content = content
  addVisible.value = true
}

function openEdit(data: any) {
  editFormOrigin = data
  editForm = reactive({...data})
  editVisible.value = true
}

// 暴露打开方法给父组件调用
defineExpose({openAdd, openEdit})

async function handleAdd() {
  if (!addForm.content) return
  isNowAdd.value = true
  const p = new Profile()
  p.content = addForm.content
  try {
    const pList = await props.api.addProfileFromInput(p)
    if (pList?.length) pList.forEach((item: any) => props.profiles.push(item))
    addForm.content = ''
    addVisible.value = false
    emit('refreshList')
  } catch (e: any) {
    if (e['message']) pError(e['message'])
    emit('refreshList')
  } finally {
    isNowAdd.value = false
  }
}

function validateInterval(value: any) {
  if (!value) return true
  return /^[1-9][0-9]?$|^1[0-2][0-8]$/.test(value.toString())
}

async function handleSaveEdit() {
  if (editForm.type === 2 && !editForm.title) {
    pError(t('profiles.edit.title-tip'))
    return
  }
  if (editForm.type === 1) {
    if (!editForm.title) return pError(t('profiles.edit.title-tip'))
    if (!editForm.content) return pError(t('profiles.edit.url-tip'))
    if (!isHttpOrHttps(editForm.content)) return pError(t('profiles.edit.url-error'))
    if (!validateInterval(editForm.interval)) return pError(t('profiles.edit.update-tip'))
  }

  isNowEdit.value = true
  await props.api.updateProfile(editForm)
  isNowEdit.value = false

  Object.assign(editFormOrigin, editForm)
  editVisible.value = false
  pSuccess(t('profiles.edit.success'))

  Events.Emit({name: 'profiles', data: props.profiles})
  props.api.getRuleNum().then((res: any) => menuStore.setRuleNum(res))
}
</script>

<template>
  <!-- 新增弹窗 -->
  <el-dialog v-model="addVisible" :title="t('profiles.add')" width="520px" draggable center>
    <el-form :model="addForm">
      <el-form-item>
        <el-input
            v-model="addForm.content"
            :rows="3"
            type="textarea"
            autocapitalize="off"
            autocomplete="off"
            spellcheck="false"
            :placeholder="t('profiles.placeholder')"
        />
      </el-form-item>
    </el-form>
    <template #footer>
      <div class="dialog-footer">
        <el-button @click="addVisible = false">{{ t('cancel') }}</el-button>
        <el-button :loading="isNowAdd" type="primary" @click="handleAdd">{{ t('confirm') }}</el-button>
      </div>
    </template>
  </el-dialog>

  <!-- 编辑弹窗 -->
  <el-dialog v-model="editVisible" :title="t('edit')" width="520px" draggable center>
    <el-form :model="editForm" label-position="top">
      <el-form-item :label="t('profiles.edit.title')" label-width="120">
        <el-input v-model="editForm.title" clearable autocapitalize="off" autocomplete="off" spellcheck="false"/>
      </el-form-item>
      <el-form-item v-if="editForm.type == 1" :label="t('profiles.edit.url')" label-width="120">
        <el-input v-model="editForm.content" clearable autocapitalize="off" autocomplete="off" spellcheck="false"/>
      </el-form-item>
      <el-form-item v-if="editForm.type == 1" :label="t('profiles.edit.update')" label-width="120">
        <el-input v-model="editForm.interval" clearable autocapitalize="off" autocomplete="off" spellcheck="false"/>
      </el-form-item>
      <el-form-item :label="t('profiles.edit.template')" label-width="120">
        <el-select v-model="editForm.template" placeholder="" clearable>
          <el-option
              v-for="item in tList"
              :key="item.id"
              :label="getTemplateTitle(t, item.title)"
              :value="item.id"
          />
        </el-select>
      </el-form-item>
    </el-form>
    <template #footer>
      <div class="dialog-footer">
        <el-button @click="editVisible = false">{{ t('cancel') }}</el-button>
        <el-button type="primary" :loading="isNowEdit" @click="handleSaveEdit">{{ t('confirm') }}</el-button>
      </div>
    </template>
  </el-dialog>
</template>