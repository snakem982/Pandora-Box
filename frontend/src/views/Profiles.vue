<script setup lang="ts">
let sArray = [
  {title: "亚马逊机房", selected: false},
  {title: "CloudFlare 转发", selected: false},
  {title: "大威天龙", selected: false},
  {title: "降龙十八掌", selected: false},
  {title: "打狗棒法", selected: false},
  {title: "一阳指", selected: false},
  {title: "庐山升龙霸", selected: false},
  {title: "天马流星拳", selected: false},
  {title: "巴拉拉小魔仙", selected: false},
  {title: "榴莲", selected: true},
  {
    title: "百分百空百分百空手接白刃百分百空手接白刃百分百空手接白刃百分百空手接白刃百分百空手接白刃手接白刃",
    selected: false
  },
  {title: "士力架", selected: false},
]
let configs = reactive([])


onBeforeMount(function (): void {
  configs = sArray
})


const canDrag = ref(false)

function mouseEnter() {
  canDrag.value = true
}

function mouseLeave() {
  canDrag.value = false
}


function handleEmit(value: any) {
  console.log(configs)
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
            <el-icon class="profile-option-btn">
              <icon-mdi-plus-thick/>
            </el-icon>
          </el-tooltip>

          <el-tooltip
              :content="$t('profiles.paste')"
              placement="top">
            <el-icon class="profile-option-btn">
              <icon-mdi-content-paste/>
            </el-icon>
          </el-tooltip>

          <el-tooltip
              :content="$t('profiles.open')"
              placement="top">
            <el-icon class="profile-option-btn">
              <icon-mdi-folder-open/>
            </el-icon>
          </el-tooltip>
        </div>
      </el-space>

      <div class="sub-title">
        <span>{{ $t('profiles.available') }} 50G</span>
        <el-divider direction="vertical" border-style="dashed"/>
        <span>{{ $t('profiles.use') }} 100G</span>
        <el-divider direction="vertical" border-style="dashed"/>
        <span>2025-05-06 23:59 {{ $t('profiles.expire') }}</span>
        <el-divider direction="vertical" border-style="dashed"/>
        <span>2025-04-06 23:59 {{ $t('profiles.update') }}</span>
      </div>
    </template>
    <template #bottom>

      <VDContainer
          :data="configs"
          @getData="handleEmit"
          :gap="15"
          :draggable="canDrag"
          style="margin-left: 10px;width: 95%;"
      >
        <template v-slot:VDC="{data,index}">
          <div :class="data.selected?'sub-card sub-card-select':'sub-card'">
            <div class="row">
              <el-icon
                  @mouseenter="mouseEnter"
                  @mouseleave="mouseLeave"
                  size="22"
                  class="drag">
                <icon-mdi-drag/>
              </el-icon>
              <el-icon size="22">
                <icon-mdi-refresh/>
              </el-icon>
            </div>
            <div class="system-info">
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
  width: calc(33% - 30px);
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