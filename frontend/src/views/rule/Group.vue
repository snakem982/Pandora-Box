<script setup>
import {VAceEditor} from 'vue3-ace-editor';
import 'ace-builds/src-noconflict/ace';
import 'ace-builds/src-noconflict/mode-yaml'; // YAML 支持
import 'ace-builds/src-noconflict/theme-cloud9_night'; // 主题支持

const yamlContent = ref(`# 示例 YAML
key: 🚀 节点选择
list:
  - 🇺🇸 美国
  - item2 🇭🇰`);

const editorOptions = {
  showPrintMargin: false,
};


</script>

<template>
  <div class="group">
    <VAceEditor
        v-model:value="yamlContent"
        lang="yaml"
        theme="cloud9_night"
        :options="editorOptions"
        style="width: 100%;height: 300px;"
    />
  </div>
</template>

<style scoped>
.group {
  width: 95%;
  margin-left: 10px;
  margin-top: 5px;
}

:deep(.ace_editor) {
  border: 2px solid var(--text-color);
  border-radius: 8px;
  font: 15px 'Twemoji', 'Monaco', 'Menlo', 'Ubuntu Mono', 'Consolas', 'Source Code Pro', 'source-code-pro', monospace;
}

:deep(.ace_gutter) {
  border-top-left-radius: 8px;
  border-bottom-left-radius: 8px;
}

</style>