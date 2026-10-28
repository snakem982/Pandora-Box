<!-- src/views/profiles/components/ProfileHeader.vue -->
<script setup lang="ts">
import {useI18n} from 'vue-i18n'

defineProps<{
  headerShow: {
    available: string
    used: string
    expire: string
    update: string
  }
}>()

const emit = defineEmits(['add', 'paste', 'openFile', 'locate'])
const {t} = useI18n()
</script>

<template>
  <el-space class="space">
    <div class="title">{{ t('profiles.title') }}</div>
    <div class="profile-option">
      <el-tooltip :content="t('profiles.add')" placement="top">
        <el-icon @click="emit('add')" class="profile-option-btn">
          <icon-mdi-plus-thick/>
        </el-icon>
      </el-tooltip>
      <el-tooltip :content="t('profiles.paste')" placement="top">
        <el-icon @click="emit('paste')" class="profile-option-btn">
          <icon-mdi-content-paste/>
        </el-icon>
      </el-tooltip>
      <el-tooltip :content="t('profiles.open')" placement="top">
        <el-icon @click="emit('openFile')" class="profile-option-btn">
          <icon-mdi-folder-open/>
        </el-icon>
      </el-tooltip>
      <el-tooltip :content="t('profiles.locate')" placement="top">
        <el-icon @click="emit('locate')" class="profile-option-btn">
          <icon-mdi-crosshairs-gps/>
        </el-icon>
      </el-tooltip>
    </div>
  </el-space>

  <div class="sub-title">
    <template v-if="headerShow.used">
      <span>{{ t('profiles.use') }} {{ headerShow.used }}</span>
      <el-divider direction="vertical" border-style="dashed"/>
    </template>
    <template v-if="headerShow.available">
      <span>{{ t('profiles.available') }} {{ headerShow.available }}</span>
      <el-divider direction="vertical" border-style="dashed"/>
    </template>
    <template v-if="headerShow.expire">
      <span>{{ t('profiles.expire') }} {{ headerShow.expire }}</span>
      <el-divider direction="vertical" border-style="dashed"/>
    </template>
    <template v-if="headerShow.update">
      <span>{{ t('profiles.update') }} {{ headerShow.update }}</span>
    </template>
  </div>
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
  color: var(--top-hr-color);
  font-size: 14px;
  margin-top: 15px;
}

.profile-option {
  margin-left: 10px;
  font-size: 30px;
  padding-top: 10px;
}

.profile-option-btn {
  margin-right: 15px;
  cursor: pointer;
}

.profile-option-btn:hover {
  color: var(--hr-color);
}
</style>