<script lang="ts" setup>
import {onBeforeUnmount, onMounted, reactive, ref} from 'vue'
import {WS} from '../api/ws'
import {Connection, Connections, Metadata} from "../api/pandora";
import {formatDistance} from 'date-fns';
import {enUS, zhCN} from 'date-fns/locale'
import {del} from "../api/http";
import {GetFreePort} from "../../wailsjs/go/main/App";


const search = ref('')
const tableData = reactive({
  array: []
})
const local = (navigator.language.toLowerCase().indexOf("zh") > -1 ? zhCN : enUS)


function onmessage(ws: WS, ev: MessageEvent): any {
  const myObj: Connections = JSON.parse(ev.data);
  myObj.connections.sort((a, b) => {
    const dateA = new Date(a.start);
    const dateB = new Date(b.start);
    return dateB > dateA ? 1 : -1
  })
  tableData.array = myObj.connections as any
}

function filterData(cacheData: Connection[]): Connection[] {
  return cacheData.filter(
      (data) =>
          !search.value ||
          fHost(data.metadata).toLowerCase().includes(search.value.toLowerCase())
  )
}

let ws: WS
onMounted(async () => {
  const addr = await GetFreePort()
  const split = addr.split(":")
  ws = new WS(split[0], split[1], '/connections?token=', null, onmessage);
})

onBeforeUnmount(() => {
  if (ws) {
    ws.ws.close()
  }
})

function fHost(metadata: Metadata): string {
  return (metadata.host || metadata.destinationIP) + ':' + metadata.destinationPort
}

function fDate(ok: string): string {
  const date = new Date(ok);
  return formatDistance(new Date(), date, {locale: local})
}

const UNITS = ['B', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB'];

function prettyBytes(n: number) {
  if (n < 1000) {
    return n + ' B';
  }
  const exponent = Math.min(Math.floor(Math.log10(n) / 3), UNITS.length - 1);
  n = Number((n / Math.pow(1000, exponent)).toPrecision(3));
  const unit = UNITS[exponent];
  return n + ' ' + unit;
}

function close(id: string) {
  del("/connections/" + id)
}

function fProcess(row: any, column: any, cellValue: any): string {
  let ok = cellValue || "-"
  ok = ok.replace(' Helper', '')
  return ok
}

function closeAll() {
  const data = filterData(tableData.array)
  if (data.length > 0) {
    if (search.value) {
      for (let connection of data) {
        close(connection.id)
      }
    } else {
      del("/connections")
    }
  }
}

</script>

<template>
  <el-row :gutter="24">
    <el-col :span="18">
      <el-input
          v-model="search"
          placeholder="域名搜索 Host Search"
          class="search"
          size="large"
          autocomplete="off"
      >
      </el-input>
    </el-col>
    <el-col :span="6">
      <el-button
          type="danger"
          style="margin-top: 3px"
          @click="closeAll">一键关闭 Close All
      </el-button>
    </el-col>
  </el-row>
  <el-table :data="filterData(tableData.array)" max-height="85vh" empty-text="暂无数据 No Data" stripe>
    <el-table-column prop="metadata.type" label="类型 Type" width="100"/>
    <el-table-column prop="metadata.process"
                     label="程序 Process"
                     width="220"
                     align="center"
                     :formatter="fProcess"/>
    <el-table-column label="域名 Host" width="400">
      <template #default="scope">
        {{ fHost(scope.row.metadata) }}
      </template>
    </el-table-column>
    <el-table-column label="规则 Rule" width="400">
      <template #default="scope">
        {{ scope.row.rule + (scope.row.rulePayload ? " → " + scope.row.rulePayload : "") }}
      </template>
    </el-table-column>
    <el-table-column label="节点链 Chains" width="400">
      <template #default="scope">
        {{ scope.row.chains[scope.row.chains.length - 1] + " → " + scope.row.chains[0] }}
      </template>
    </el-table-column>
    <el-table-column label="连接时间 Time" width="150">
      <template #default="scope">
        {{ fDate(scope.row.start) }}
      </template>
    </el-table-column>
    <el-table-column label="上传 Upload" width="150">
      <template #default="scope">
        <el-text type="danger">
          <el-icon>
            <Top/>
          </el-icon>
          {{ prettyBytes(scope.row.upload) }}
        </el-text>
      </template>
    </el-table-column>
    <el-table-column label="下载 Download" width="150">
      <template #default="scope">
        <el-text type="success">
          <el-icon>
            <Bottom/>
          </el-icon>
          {{ prettyBytes(scope.row.download) }}
        </el-text>
      </template>
    </el-table-column>
    <el-table-column fixed="right" label="关闭 Close" width="120" align="center">
      <template #default="scope">
        <el-button type="danger" circle @click="close(scope.row.id)">
          <el-icon>
            <CloseBold/>
          </el-icon>
        </el-button>
      </template>
    </el-table-column>
  </el-table>
</template>

<style scoped>
:deep(.cell) {
  font-size: 16px;
}

.search {
  margin-bottom: 10px;
}
</style>
