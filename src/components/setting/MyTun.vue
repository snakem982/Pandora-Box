<template>
  <div class="custom-style">
    <span class="liable">Tun Stack:</span>
    <el-segmented v-model="settingStore.stack" :options="options"/>
  </div>
</template>

<script lang="ts" setup>
import {useMenuStore} from "@/store/menuStore";
import {useSettingStore} from "@/store/settingStore";
import createApi from "@/api";
import {pUpdateMihomo} from "@/util/mihomo";

// 使用 store
const menuStore = useMenuStore()
const settingStore = useSettingStore()

// 获取当前 Vue 实例的 proxy 对象 和 api
const {proxy} = getCurrentInstance()!;
const api = createApi(proxy);

// 更新值
watch(() => settingStore.stack, async () => {
  // 先关闭一下再开，会切换更快
  if (menuStore.tun) {
    await api.updateConfigs({tun: {enable: false}})
  }

  // 更新配置
  api.updateConfigs({
    tun: {
      enable: menuStore.tun,
      stack: settingStore.stack,
    },
  }).then(() => {
    // 同步 mihomo 配置
    pUpdateMihomo(menuStore, settingStore, api)
  });
});

const options = ['Mixed', 'gVisor', 'System']
</script>

<style scoped>
.custom-style .el-segmented {
  height: 30px;
  margin-left: 10px;
  border: 1px solid var(--text-color);
  background: rgba(255, 255, 255, 0.1);
  --el-segmented-item-selected-color: var(--text-color);
  --el-segmented-item-selected-bg-color: var(--left-item-selected-bg);
  --el-border-radius-base: 5px;
  color: var(--text-color);
  font-size: 15px;
}

.liable {
  font-size: 18px;
  font-weight: bold;
}
</style>