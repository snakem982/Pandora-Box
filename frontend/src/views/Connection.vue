<script setup lang="ts">
import MyHr from "@/components/proxies/MyHr.vue";
import MySimpleInput from "@/components/MySimpleInput.vue";
import {WS} from "@/util/ws";
import {useWebStore} from "@/store/webStore";
import {prettyBytes, rJoin} from "@/util/format";
import {onBeforeRouteLeave} from "vue-router";
import {formatDistance} from 'date-fns';
import {enUS, zhCN} from 'date-fns/locale'
import {useI18n} from "vue-i18n";
import createApi from "@/api";

// 获取当前 Vue 实例的 proxy 对象 和 api
const {proxy} = getCurrentInstance()!;
const api = createApi(proxy);

// 获取 i18n
const {t} = useI18n()
const localeMap = {
  '简体中文': zhCN,
  'English': enUS,
};

function fDate(start: any): string {
  const startTime = new Date(start);
  return formatDistance(new Date(), startTime, {locale: localeMap[t('language')]})
}

const distanceFromTop = ref(195)
const upFromTop = function (distance: number) {
  distanceFromTop.value = distance
}

const search = ref('')

function handleInputChange(value: any) {
  search.value = value
}

function fHost(metadata: any): string {
  return (metadata.host || metadata.destinationIP) + ':' + metadata.destinationPort
}

function filterData(cacheData: any): any {

  if(!cacheData || cacheData.length === 0) {
    return
  }

  const cache = cacheData.filter((data: any) => {
    const searchLower = search.value.toLowerCase();
    return (
      (!search.value || fHost(data.metadata).toLowerCase().includes(searchLower)) || // 主机过滤
      data.rule.toLowerCase().includes(searchLower) || // 规则过滤
      (data.metadata.process && data.metadata.process.toLowerCase().includes(searchLower)) // 程序过滤
    );
  });

  cache.sort((obj1: any, obj2: any) => obj2.start.localeCompare(obj1.start));

  return cache;
}

// 分页数据状态
const paginatedData = ref([]);

function onConn(ev: MessageEvent) {
  const parsedData = JSON.parse(ev.data);
  paginatedData.value = parsedData['connections']
}

const webStore = useWebStore()
let wsConn: WS
onMounted(() => {
  const urlTraffic = webStore.wsUrl + "/connections?token=" + webStore.secret;
  wsConn = new WS(urlTraffic, null, onConn);
})

// 路由切换前关闭 WebSocket
onBeforeRouteLeave(() => {
  wsConn.close();
});

onBeforeUnmount(() => {
  wsConn.close();
})


function closeAll() {
  const data = filterData(paginatedData.value)
  if (data.length > 0) {
    if (search.value) {
      for (let connection of data) {
        api.closeConnection(connection.id)
      }
    } else {
      api.closeAllConnection()
    }
  }
}


</script>

<template>
  <MyLayout :top-height="distanceFromTop-15" :bottom-height="distanceFromTop+25">
    <template #top>
      <MySearch></MySearch>
      <el-space class="space">
        <div class="title">
          {{ $t('connections.title') }}
        </div>
      </el-space>
      <MyHr :update="upFromTop" v-show="false"></MyHr>
    </template>
    <template #bottom>
      <div class="conn">
        <el-space class="op">
          <div class="search">
            <MySimpleInput
                :onInputChange="handleInputChange"
                :placeholder="$t('connections.search')"
                class="search"
            ></MySimpleInput>
          </div>
          <el-button @click="closeAll">
            {{ $t('connections.close') }}
          </el-button>
        </el-space>
      </div>

      <div class="content">
        <div class="info-list">
          <el-row
              class="info"
              v-for="(item, i) in filterData(paginatedData)"
              :key="i"
          >
            <el-col :span="24">
              <el-tag type="success" size="small">{{ item.metadata.type }}</el-tag>
              &emsp;
              <el-tag type="danger" size="small">
                {{ fDate(item.start) }}
              </el-tag>
              <template v-if="item.metadata.process">
                &emsp;
                <el-tag type="primary" size="small">{{ item.metadata.process }}</el-tag>
              </template>
              <div class="od">
                <span class="ot">{{ $t('connections.host') }} : </span>
                {{ item.metadata.host }}:{{ item.metadata.destinationPort }}
              </div>
              <div class="od">
                <span class="ot">{{ $t('connections.download') }} : </span>
                {{ prettyBytes(item.download) }}
                &emsp;
                <span class="ot">{{ $t('connections.upload') }} : </span>
                {{ prettyBytes(item.upload) }}
              </div>
              <div class="od" v-if="item.rule">
                <span class="ot">{{ $t('connections.rule') }} : </span>
                {{ item.rule }}
                {{ item.rulePayload ? ' / ' + item.rulePayload : '' }}
              </div>
              <div class="od">
                <span class="ot">{{ $t('connections.chains') }} : </span>
                {{ rJoin(item.chains, '&nbsp;/&nbsp;') }}
              </div>
            </el-col>
          </el-row>
        </div>
      </div>

    </template>
  </MyLayout>
</template>

<style scoped>
.space {
  margin-top: 20px;
}

.conn {
  width: 95%;
  margin-left: 10px;
  margin-top: 5px;
}

.title {
  font-size: 32px;
  font-weight: bold;
  margin-left: 10px;
}

.search {
  width: 400px;
}

:deep(.el-button) {
  padding: 2px 10px;
  --el-button-bg-color: transparent;
  --el-button-text-color: var(--text-color);
  --el-button-hover-text-color: var(--left-item-selected-bg);
  --el-button-hover-bg-color: var(--text-color)
}

.content {
  border: 2px solid var(--text-color);
  margin-top: 20px;
  width: calc(95% - 10px);
  margin-left: 10px;
  border-radius: 10px;
}

.info-list {
  max-height: calc(100vh - 250px);
  overflow-y: auto;
  border-radius: 10px;
}

.info {
  border-bottom: 1px solid #ccc;
  padding: 5px 10px;
  font-size: 15px;
  line-height: 1.6;
  background-color: rgba(0, 0, 0, 0.1);
  border-radius: 10px;
}

.od {
  -webkit-user-select: text;
  user-select: text;
}

.ot {
  font-weight: bold;
  font-size: 15px;
}

.info-list::-webkit-scrollbar {
  width: 5px;
  padding-bottom: 20px;
}

.info-list::-webkit-scrollbar-track {
  background: transparent;
}

.info-list::-webkit-scrollbar-thumb {
  background: var(--scrollbar-bg);
  border-radius: 2px;
  transition: background 0.3s ease, box-shadow 0.3s ease;
}

.info-list::-webkit-scrollbar-thumb:hover {
  background: var(--scrollbar-hover-bg);
  box-shadow: var(--scrollbar-hover-shadow);
}


</style>