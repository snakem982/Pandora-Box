<script setup lang="ts">

import MyPort from "@/components/setting/MyPort.vue";
import MyTun from "@/components/setting/MyTun.vue";
import {EditPen} from "@element-plus/icons-vue";
import {useWebStore} from "@/store/webStore";
import { copy } from "@/util/pLoad";
import {useI18n} from "vue-i18n";
import { useSettingStore } from "@/store/settingStore";
import createApi from "@/api";

// 获取当前 Vue 实例的 proxy 对象 和 api
const {proxy} = getCurrentInstance()!;
const api = createApi(proxy);

// 使用 store
const webStore = useWebStore()
const settingStore = useSettingStore()
const {t} = useI18n()


const dns = ref(false)
const lan = ref(false)
const ipv6 = ref(false)
const startup = ref(false)


onMounted(() => {
  dns.value = settingStore.dns
  lan.value = settingStore.lan
  ipv6.value = settingStore.ipv6
  startup.value = settingStore.startup
})


</script>

<template>
  <el-row :gutter="20" class="spark"
          style="margin-left: 0;
          margin-right: 0;">
    <el-col :span="24">
      <div class="box box1">
        <div class="title">
          Mihomo
        </div>
        <hr/>
        <ul class="info-list">
          <li>
            <MyPort></MyPort>
          </li>
          <li>
            <MyTun></MyTun>
          </li>
          <li>
            <strong>
              {{ $t('setting.mihomo.dns') }} :
            </strong>
            <el-icon class="btn">
              <EditPen/>
            </el-icon>
            <el-switch
                v-model="dns"
                class="set-switch"
                style="margin-left: 28px"
            />
          </li>
          <li>
            <strong>{{ $t('setting.mihomo.lan') }} :</strong>
            <el-switch
                v-model="lan"
                class="set-switch"
            />
          </li>
          <li>
            <strong>IPV6 :</strong>
            <el-switch
                v-model="ipv6"
                class="set-switch"
            />
          </li>
          <li style="height: 30px">
            <strong>API :</strong>
            {{ webStore.baseUrl }}
            <el-button 
            @click="copy(webStore.baseUrl,t)"
            >复制</el-button>
          </li>
          <li style="height: 30px">
            <strong>Secret:</strong>
            {{ webStore.secret }}
            <el-button
            @click="copy(webStore.secret,t)"
            >复制</el-button>
          </li>
        </ul>
      </div>
    </el-col>
  </el-row>

  <el-row :gutter="20" class="spark"
          style="margin-left: 0;
          margin-top: 30px;
          margin-right: 0;">
    <el-col :span="24">
      <div class="box box2">
        <div class="title">
          Pandora-Box
        </div>
        <hr/>
        <ul class="info-list">
          <li>
            <strong>{{ $t('setting.px.startup') }} :</strong>
            <el-switch
                disabled
                v-model="startup"
                class="set-switch"
            />
          </li>
          <li style="height: 30px">
            <strong>{{ $t('setting.px.dir') }} :</strong>
            <el-button style="margin-left: 10px">
              {{ $t('setting.px.open') }}
            </el-button>
            <el-button>{{ $t('setting.px.export') }}</el-button>
            <el-button>{{ $t('setting.px.import') }}</el-button>
          </li>
          <li style="height: 30px">
            <strong>{{ $t('setting.px.update') }} :</strong>
            <el-button style="margin-left: 10px">{{ $t('setting.px.check') }}</el-button>
          </li>
        </ul>
      </div>
    </el-col>
  </el-row>



</template>

<style scoped>
.spark {
  max-width: 95%;
}

.box {
  padding: 10px;
  border-radius: 8px;
  text-align: left;
}

.box hr {
  border: none;
  height: 1px;
  background-color: var(--hr-color);
  margin: 10px 0;
}

.info-list {
  list-style: none;
  padding: 0;
}

.info-list li {
  font-size: 18px;
  margin: 8px 0;
}

.box1 {
  box-shadow: var(--right-box-shadow);
}

.box2 {
  box-shadow: var(--right-box-shadow);
}

.set-switch {
  margin-left: 10px;
  --el-switch-border-color: var(--text-color);
  --el-switch-on-color: var(--left-item-selected-bg);
  --el-switch-off-color: transparent;
}

:deep(.el-switch__core) {
  width: 46px;
  height: 26px;
  border-radius: 12px;
  border: 2px solid var(--text-color);
}

:deep(.el-switch__core .el-switch__action) {
  margin-left: 2px;
}

:deep(.el-switch.is-checked .el-switch__core .el-switch__action) {
  left: calc(100% - 21px);
}

:deep(.el-button) {
  padding: 2px 10px;
  --el-button-bg-color: transparent;
  --el-button-text-color: var(--text-color);
  --el-button-hover-text-color: var(--left-item-selected-bg);
  --el-button-hover-bg-color: var(--text-color)
}

.btn {
  font-size: 18px;
  position: absolute;
  margin-top: 6px;
}

.btn:hover {
  cursor: pointer;
  color: var(--hr-color);
}


</style>