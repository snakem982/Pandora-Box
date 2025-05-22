<script setup lang="ts">
import MyEditor from "@/components/MyEditor.vue";
import createApi from "@/api";
import {pError, pSuccess} from "@/util/pLoad";
import {useI18n} from "vue-i18n";

// 获取当前 Vue 实例的 proxy 对象
const {proxy} = getCurrentInstance()!;
const api = createApi(proxy);

// i18n
const {t} = useI18n();

const load = function (yamlContent: any) {
  api.getDNS().then((dns) => {
    yamlContent.value = dns.content;
  })
}

const save = async function (yamlContent: any) {
  try {
    await api.updateDNS({
      data: yamlContent
    })
    pSuccess(t('dns.success'))
  } catch (e) {
    if (e['message']) {
      pError(e['message'])
    }
  }
}
</script>

<template>
  <MyLayout>
    <template #top>
      <el-space class="space">
        <div class="title">
          {{ $t("dns.title") }}
        </div>
      </el-space>
    </template>
    <template #bottom>
      <MyEditor :load="load" :save="save"></MyEditor>
    </template>
  </MyLayout>
</template>

<style scoped>
.space {
  margin-top: 20px;
}

.title {
  font-size: 32px;
  font-weight: bold;
  margin-left: 10px;
}
</style>
