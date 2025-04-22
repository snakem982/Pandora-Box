<script setup lang="ts">
import {VAceEditor} from "vue3-ace-editor";
import "ace-builds/src-noconflict/ace";
import "ace-builds/src-noconflict/ext-searchbox"; // 查找替换
import "ace-builds/src-noconflict/mode-yaml"; // YAML 支持
import "ace-builds/src-noconflict/ext-beautify";
import "ace-builds/src-noconflict/ext-language_tools"; // YAML 支持
import "ace-builds/src-noconflict/theme-monokai"; // 主题支持
import createApi from "@/api";
import {useI18n} from "vue-i18n";
import {pError, pSuccess} from "@/util/pLoad";
import {useMenuStore} from "@/store/menuStore";
import {useProxiesStore} from "@/store/proxiesStore";

// 编辑器使用
const editorOptions = {
  showPrintMargin: false,
};
// 编辑器显示内容
const yamlContent = ref("");

// 当前页面使用store
const menuStore = useMenuStore();
const proxiesStore = useProxiesStore();

// i18n
const {t} = useI18n();

// 获取当前 Vue 实例的 proxy 对象 和 api
const {proxy} = getCurrentInstance()!;
const api = createApi(proxy);


// Template列表
let tList = reactive([]);
// Template
let now = reactive({
  id: "",
  title: "",
  selected: false
})

onMounted(async () => {
  // 初始化
  tList = await api.getTemplateList();
  Object.assign(now, tList[0]);

  // 处理选中项
  for (const item of tList) {
    if (item.selected) {
      Object.assign(now, item);
      break;
    }
  }

  // 处理编辑器内容
  yamlContent.value = await api.getTemplateById(now.id);
});

// Template 下拉列表逻辑
const isDropdownOpen = ref(false);
const selectOption = async (item: any) => {
  Object.assign(now, item);
  // 处理编辑器内容
  yamlContent.value = await api.getTemplateById(item.id);
  isDropdownOpen.value = false;
};

// 测试逻辑
const testTemplate = async () => {
  try {
    await api.testTemplate({
      data: yamlContent.value,
    });
    pSuccess(t('rule.group.test-success'))
  } catch (e) {
    if (e['message']) {
      pError(e['message'])
    }
  }
}

// 保存逻辑
const saveTemplate = async () => {
  try {
    // 先测试
    await api.testTemplate({
      data: yamlContent.value,
    });
    // 再保存
    await api.updateTemplate({
      data: yamlContent.value,
      template: now,
    });
    // 如果是启用中的 进行切换
    if (now.selected) {
      await api.switchTemplate(now);
      proxiesStore.active = ""
      api.getRules().then((res) => {
        menuStore.setRuleNum(res.length);
      });
    }
    pSuccess(t('rule.success'))
  } catch (e) {
    if (e['message']) {
      pError(e['message'])
    }
  }
}

// 切换逻辑
const switchTemplate = async () => {
  try {
    await api.switchTemplate(now);
    tList = await api.getTemplateList();
    pSuccess(t('rule.group.switch-success'))

    proxiesStore.active = ""
    api.getRules().then((res) => {
      menuStore.setRuleNum(res.length);
    });

  } catch (e) {
    if (e['message']) {
      pError(e['message'])
    }
  }
}


</script>

<template>
  <div class="group">
    <el-space class="op">
      <div class="dropdown">
        <button class="dropdown-btn" @click="isDropdownOpen = !isDropdownOpen">
          {{ t("rule.group." + now.title) }}
        </button>
        <ul v-if="isDropdownOpen" class="dropdown-list">
          <li
              :key="item.id + index"
              @click="selectOption(item)"
              class="dropdown-item"
              v-for="(item, index) in tList"
          >
            {{ t("rule.group." + item.title) }}
          </li>
        </ul>
      </div>
      <el-divider direction="vertical" border-style="dashed"/>
      <el-button @click="testTemplate">
        {{ t("rule.group.test") }}
      </el-button>
      <el-button @click="saveTemplate">
        {{ t("save") }}
      </el-button>
      <el-divider direction="vertical" border-style="dashed"/>
      <el-text :class="now.selected ? '' : 'st'">{{ t("off") }}</el-text>
      <el-switch @click="switchTemplate" v-model="now.selected" class="set-switch"/>
      <el-text :class="now.selected ? 'st' : ''">{{ t("on") }}</el-text>
    </el-space>

    <VAceEditor
        v-model:value="yamlContent"
        lang="yaml"
        theme="monokai"
        :options="editorOptions"
        style="width: 100%; height: calc(100vh - 325px)"
        class="editor"
    />
  </div>
</template>

<style scoped>
.group {
  width: 95%;
  margin-left: 10px;
  margin-top: 5px;
}

.op {
  margin-top: 8px;
}

.dropdown {
  position: relative;
  display: inline-block;
}

.dropdown-btn {
  background: transparent;
  color: white;
  border: 2px solid white;
  padding: 5px 10px;
  cursor: pointer;
  font-size: 15px;
  outline: none;
  border-radius: 8px;
  min-width: 150px;
}

.dropdown-btn:hover {
  opacity: 0.8;
}

.dropdown-list {
  position: absolute;
  background: rgba(0, 0, 0, 0.8);
  border: 2px solid white;
  margin-top: 4px;
  padding: 0;
  list-style: none;
  min-width: 146px;
  z-index: 20;
  border-radius: 8px;
  font-size: 15px;
  text-align: center;
}

.dropdown-item {
  color: white;
  padding: 8px;
  cursor: pointer;
}

.dropdown-item:hover {
  background: rgba(255, 255, 255, 0.2);
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
  --el-button-hover-bg-color: var(--text-color);
}

.st {
  color: var(--text-color);
}

.editor {
  margin-top: 25px;
}

:deep(.ace_editor) {
  border: 2px solid var(--text-color);
  border-radius: 8px;
  font: 15px "Twemoji", "Monaco", "Menlo", "Ubuntu Mono", "Consolas",
  "Source Code Pro", "source-code-pro", monospace;
}

:deep(.ace_gutter) {
  border-top-left-radius: 8px;
  border-bottom-left-radius: 8px;
}

:deep(.ace_search.right) {
  width: 420px;
  margin-left: 10px;
  margin-right: -4px;
  padding-left: 8px;
  margin-top: 0;
  border: none;
  float: right;
  color: var(--text-color);
}

:deep(.ace_search_form, .ace_replace_form) {
  margin: 0;
}

:deep(.ace_search_form.ace_nomatch) {
  width: 374px;
}

:deep(.ace_button, .ace_searchbtn_close) {
  color: #cccccc;
}

:deep(.ace_button:hover) {
  color: black;
}
</style>
