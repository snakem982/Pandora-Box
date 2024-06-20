<script lang="ts" setup>
import {computed, onMounted, reactive, ref, watch} from "vue"
import type {CmComponentRef} from "codemirror-editor-vue3"
import Codemirror from "codemirror-editor-vue3"
import type {Editor, EditorConfiguration} from "codemirror"
import "codemirror/mode/yaml/yaml.js"
import SvgIcon from "@jamescoyle/vue-icon";

// search
import "codemirror/addon/search/search.js"
import "codemirror/addon/search/searchcursor.js"
import "codemirror/addon/search/jump-to-line.js"
import "codemirror/addon/scroll/simplescrollbars.js"
import "codemirror/addon/scroll/simplescrollbars.css"
import "codemirror/addon/dialog/dialog.js"
import "codemirror/addon/dialog/dialog.css"

// theme
import "codemirror/theme/eclipse.css"
import "codemirror/theme/material-darker.css"
import {mdiArrowUpBold} from "@mdi/js";

const props = defineProps({
  load: Function,
  reset: Function,
  test: Function,
  save: Function,
  onSwitch: Function,
  isOn: Boolean,
  height: {
    type: String,
    default: "70vh"
  },
});

const isOn = ref(false)
const code = ref("")
const cmRef = ref<CmComponentRef>()
const cmOptions: EditorConfiguration = reactive({
  mode: "text/x-yaml",
  theme: "eclipse",
  readOnly: false,
  lineNumbers: true,
  lineWiseCopyCut: true,
  scrollbarStyle: 'overlay'
});

function goTop() {
  const ele: any = document.querySelector('.CodeMirror-scroll')
  ele.scrollTop = 0;
}

const onReady = (cm: Editor) => {
  cm.focus()
}

onMounted(async () => {
  const isDark = localStorage.getItem("vueuse-color-scheme") == "dark"
  if (isDark) {
    cmOptions.theme = "material-darker"
  }
  cmRef.value?.refresh()
  code.value = await props.load?.call(null)
  isOn.value = props.isOn
})


const zh = navigator.language.toLowerCase().indexOf("zh") > -1
const isMac = navigator.userAgent.indexOf("Mac") > -1
const getShortCut = computed(() => {
  let title = `快捷键:
    查找 Cmd+F
    查找下一个 Cmd+G
    替换 Cmd+Alt+F`
  if (zh) {
    if (!isMac) {
      title = `快捷键:
    查找 Ctrl+F
    查找下一个 Ctrl+G
    替换 Shift+Ctrl+F`
    }
  } else {
    title = `Shortcut Key:
    find Cmd+F
    findNext Cmd+G
    replace Cmd+Alt+F`
    if (!isMac) {
      title = `Shortcut Key:
    find Ctrl+F
    findNext Ctrl+G
    replace Shift+Ctrl+F`
    }
  }
  return title
})

async function RESET() {
  code.value = await props.reset?.call(null)
}

async function TEST() {
  await props.test?.call(null, code.value, false)
}

async function SAVE() {
  const ok = await props.test?.call(null, code.value, true)
  if (ok) {
    await props.save?.call(null, code.value, false)
  }
}

async function SWITCH() {
  let on = "off"
  if (isOn.value) {
    on = "on"
  }
  await props.onSwitch?.call(null, on)
}

watch(isOn, async (newValue, oldValue) => {
  await SWITCH()
});

async function beforeChange() {
  if (isOn.value == false) {
    const ok = await props.test?.call(null, code.value, true)
    if (ok) {
      return await props.save?.call(null, code.value, true)
    } else {
      return false
    }
  } else {
    return true
  }
}

</script>

<template>
  <el-alert type="warning" :title="getShortCut" style="margin-bottom: 10px"/>
  <Codemirror
      v-model:value="code"
      :options="cmOptions"
      border
      style="font-size: 16px;"
      ref="cmRef"
      :height="height"
      width="99%"
      @ready="onReady"
  >
  </Codemirror>
  <el-button
      type="danger"
      style="margin-top: 15px"
      @click="RESET">
    重置 Reset
  </el-button>

  <el-button
      type="warning"
      style="margin-top: 15px"
      @click="TEST">
    测试 Test
  </el-button>
  <el-button
      type="primary"
      style="margin-top: 15px"
      @click="SAVE">
    保存 Save
  </el-button>
  <el-switch
      v-model="isOn"
      :before-change="beforeChange"
      style="margin: 15px 0 0 20px;--el-switch-on-color: #13ce66; --el-switch-off-color: #ff4949"
      active-text="启用 On"
      inactive-text="关闭 Off"
  />
  <el-button
      style="margin: 15px 5px 0 0;float: right;"
      @click="goTop"
      type="primary"
      circle
      title="返回顶部Back To Top">
    <svg-icon type="mdi" :path="mdiArrowUpBold" :size="18"></svg-icon>
  </el-button>
</template>

<style scoped>
:deep(.CodeMirror), :deep(.CodeMirror-code) {
  font-family: 'Twemoji', 'Microsoft YaHei', '微软雅黑', 'Helvetica Neue', Helvetica, 'PingFang SC', 'Hiragino Sans GB',
  Arial, 'Nunito', -apple-system, BlinkMacSystemFont, 'Segoe UI', 'Roboto',
  'Oxygen', 'Ubuntu', 'Cantarell', 'Fira Sans', 'Droid Sans', 'Helvetica Neue',
  sans-serif;
}
</style>