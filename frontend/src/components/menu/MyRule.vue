<template>
  <div class="custom-style">
    <el-segmented v-model="menuStore.rule" :options="getOptions()">
      <template #default="scope">
        <div>{{ scope.item["label"] }}</div>
      </template>
    </el-segmented>
  </div>
</template>

<script lang="ts" setup>
import { useMenuStore } from "@/store/menuStore";
import { useI18n } from "vue-i18n";
import createApi from "@/api";
import { success } from "@/util/load";

// 存储规则模式
const menuStore = useMenuStore();

// 获取当前 Vue 实例的 proxy 对象
const { proxy } = getCurrentInstance()!;
const api = createApi(proxy);

// 国际化
const { t } = useI18n();
const getOptions = function (): any[] {
  return [
    {
      label: t("rules.rule"),
      value: "rule",
    },
    {
      label: t("rules.global"),
      value: "global",
    },
    {
      label: t("rules.direct"),
      value: "direct",
    },
  ];
};

// 监听 store.rule 的变化
watch(
  () => menuStore.rule,
  (newValue, oldValue) => {
    api.updateConfigs({
      mode: newValue,
    }).then((res: any) => {
      success(t("rules." + newValue + "-switch"));
    });
  }
);
</script>

<style scoped>
.custom-style {
  margin-left: 22px;
  margin-top: 23px;
}

.custom-style .el-segmented {
  min-width: 185px;
  border: 1px solid #ccc;
  background: rgba(255, 255, 255, 0.1);
  box-shadow: var(--left-nav-shadow);
  --el-segmented-item-selected-color: var(--text-color);
  --el-segmented-item-selected-bg-color: var(--left-item-selected-bg);
  --el-border-radius-base: 5px;
  color: var(--text-color);
  font-size: 15px;
}

.custom-style .el-segmented:hover {
  box-shadow: var(--left-nav-hover-shadow);
}
</style>
