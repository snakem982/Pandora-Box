<script setup lang="ts">
import {ref} from "vue";
import {EditPen} from "@element-plus/icons-vue";
import {useSettingStore} from "@/store/settingStore";
import {pError} from "@/util/pLoad";
import {useI18n} from "vue-i18n";
import {pUpdateMihomo} from "@/util/mihomo";
import {useMenuStore} from "@/store/menuStore";
import createApi from "@/api";

// 使用 store
const menuStore = useMenuStore()
const settingStore = useSettingStore()
const {t} = useI18n()

// 获取当前 Vue 实例的 proxy 对象 和 api
const {proxy} = getCurrentInstance()!;
const api = createApi(proxy);

// 定义数据
const isEditing = ref(false);
const port = ref(0);

// 切换编辑模式
const toggleEditing = () => {
  isEditing.value = !isEditing.value;
};

// 检查端口值是否在有效范围内
function isValidIntegerRegex(str: any) {
  return /^[1-9]\d{0,4}$/.test(str) && Number(str) <= 65535;
}

// 保存端口值
const savePort = async () => {
  // 检查端口值是否在有效范围内
  if (!isValidIntegerRegex(port.value)) {
    pError(t('setting.mihomo.port-error'))
    return;
  }

  // 端口号没有变化
  if (port.value === settingStore.port) {
    isEditing.value = false;
    return;
  }

  try {
    // 检测端口是否被占用
    await api.checkAddressPort({
      "bindAddress": settingStore.bindAddress,
      "port": Number(port.value),
    })

    // 更新配置
    api.updateConfigs({
      "mixed-port": Number(port.value),
    }).then((res: any) => {
      settingStore.setPort(port.value);
      isEditing.value = false; // 退出编辑模式
      // 同步 mihomo 配置
      pUpdateMihomo(menuStore, settingStore, api)

      if (menuStore.proxy) {
        // 未被占用开启代理
        api.enableProxy({
          "bindAddress": settingStore.bindAddress,
          "port": settingStore.port,
        })
      }
    });

  } catch (e) {
    if (e['message']) {
      pError(e['message'])
    }
  }

};

// 取消编辑
const cancelEdit = () => {
  isEditing.value = false;
  port.value = settingStore.port; // 恢复原始值
};


onMounted(() => {
  // 初始化端口值
  port.value = settingStore.port;
});
</script>

<template>
  <div class="input-container">
    <span>{{ $t('setting.mihomo.port') }} :</span>
    <template v-if="isEditing">
      <input
          type="text"
          v-model="port"
          placeholder="请输入端口号"
          autocapitalize="off"
          autocomplete="off"
          autocorrect="off"
          spellcheck="false"
      />
    </template>
    <template v-else>
      <span class="content">{{ settingStore.port }}</span>
    </template>
    <el-icon
        class="btn"
        @click="toggleEditing"
        v-if="!isEditing">
      <EditPen/>
    </el-icon>
    <el-icon
        class="btn"
        @click="savePort"
        v-if="isEditing">
      <icon-ep-select/>
    </el-icon>
    <el-icon
        class="btn"
        @click="cancelEdit"
        v-if="isEditing">
      <icon-ep-close-bold/>
    </el-icon>
  </div>
</template>

<style scoped>
.input-container {
  display: flex;
  align-items: center;
  gap: 10px;
  height: 30px;
}

span {
  color: var(--text-color);
  font-size: 18px;
  font-weight: bold;
}

.content {
  font-weight: normal;
}

input {
  width: 100px;
  padding: 5px 8px;
  border: 1px solid var(--text-color);
  border-radius: 5px;
  background-color: rgba(255, 255, 255, 0.1);
  color: var(--text-color);
  font-size: 16px;
}

input:focus {
  outline: none;
}

.btn {
  font-size: 18px;
}

.btn:hover {
  cursor: pointer;
  color: var(--hr-color);
}
</style>
