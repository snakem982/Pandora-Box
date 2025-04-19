<template>
  <div class="mask"
       @click="webStore.dnd=false"
       v-show="webStore.dnd">
    <h3>{{ t('drag.hear') }}</h3>
  </div>
</template>

<script setup lang="ts">
import {useWebStore} from "@/store/webStore.js";
import {pError, pLoad, pSuccess, pWarning} from "@/util/pLoad";
import {useI18n} from "vue-i18n";
import {Profile} from "@/types/profile.js";
import createApi from "@/api/index.js";

const {t} = useI18n();
const webStore = useWebStore();
// 获取当前 Vue 实例的 proxy 对象
const {proxy} = getCurrentInstance()!;
const api = createApi(proxy);

onMounted(() => manageDragEvents('add'));
onUnmounted(() => manageDragEvents('remove'));

function manageDragEvents(action) {
  const method = action === 'add' ? 'addEventListener' : 'removeEventListener';
  document.body[method]('dragenter', handleDragEnter);
  document.body[method]('dragover', preventDefault);
  document.body[method]('drop', handleDrop);
}

function handleDragEnter(e) {
  if (e.dataTransfer && e.dataTransfer.types.includes('Files')) {
    webStore.dnd = true;
  }
}

function preventDefault(e) {
  e.preventDefault();
}

function handleDrop(e) {
  e.preventDefault();
  webStore.dnd = false;

  const files = Array.from(e.dataTransfer.files);
  if (files.length > 1) {
    pWarning(t('drag.size'))
    return
  }

  files.forEach((file: any) => {
    const reader = new FileReader();

    reader.onload = async (event) => {
      console.log(`Content of ${file.name}:`);
      // console.log(event.target.result);

      await pLoad(t('drag.add'), async () => {
        const p = new Profile()
        p.content = event.target.result
        p.title = file.name
        try {
          const pList = await api.addProfileFromInput(p)
          if (pList && pList.length > 0) {
            webStore.dProfile = pList;
            pSuccess(t('drag.success'))
          }
        } catch (e) {
          if (e['message']) {
            pError(e['message'])
          }
        }
      })

    };

    reader.onerror = (error) => {
      console.error(`Error reading ${file.name}:`, error);
    };

    // 使用 readAsText 方法读取文件内容
    reader.readAsText(file);
  });
}


</script>

<style scoped>
.mask {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 9999;
  background: rgba(0, 0, 0, 0.8);
  display: flex;
  justify-content: center;
  align-items: center;
  color: var(--text-color);
  font-size: 1.5rem;
}

h3 {
  margin: 0;
  width: 450px;
  height: 100px;
  border: 2px dashed var(--text-color);
  text-align: center;
  padding-top: 70px;
  border-radius: 10px;
}
</style>
