<script lang="ts" setup>
import {onBeforeUnmount, onMounted, reactive, ref} from 'vue'
import {WS} from '../api/ws'
import {GetFreePort, GetSecret} from "../../wailsjs/go/main/App";

interface Getter {
  time: string
  type: string
  payload: string
}

const search = ref('')
const tableData: any = reactive({
  array: []
})

function pad(num: any) {
  return num.toString().padStart(2, '0');
}

function formatDateTime(): string {
  const date = new Date()
  const year = date.getFullYear();
  const month = date.getMonth() + 1;
  const day = date.getDate();
  const hour = date.getHours();
  const minute = date.getMinutes();
  const second = date.getSeconds();
  return `${year}/${pad(month)}/${pad(day)} ${pad(hour)}:${pad(minute)}:${pad(second)}`;
}

function onmessage(ws: WS, ev: MessageEvent) {
  tableData.array.unshift({
    time: formatDateTime(),
    ...JSON.parse(ev.data)
  })
  if (tableData.array.length > 100) {
    tableData.array.pop();
  }
}

function filterData(cacheData: Getter[]): Getter[] {
  return cacheData.filter(
      (data) =>
          !search.value ||
          data.payload.toLowerCase().includes(search.value.toLowerCase())
  )
}

let ws: WS
onMounted(async () => {
  const addr = await GetFreePort()
  const split = addr.split(":")
  const secret = await GetSecret()
  ws = new WS(split[0], split[1], '/logs?token=' + secret + '&level=info', null, onmessage);
})

onBeforeUnmount(() => {
  if (ws) {
    ws.ws.close()
  }
})

</script>

<template>
  <el-input
      v-model="search"
      placeholder="内容搜索 Content Search"
      class="search"
      size="large"
      autocomplete="off"
  >
  </el-input>
  <el-table :data="filterData(tableData.array)" table-layout="fixed" max-height="85vh" empty-text="暂无数据 No Data">
    <el-table-column fixed align="center" prop="time" label="时间 Time" width="200px"/>
    <el-table-column align="center" label="类型 Type" width="100px">
      <template #default="scope">
        <el-text size="large" type="primary" v-if="scope.row.type=='info'">INFO</el-text>
        <el-text size="large" type="warning" v-if="scope.row.type=='warning'">WARN</el-text>
        <el-text size="large" type="danger" v-if="scope.row.type=='error'">ERROR</el-text>
      </template>
    </el-table-column>
    <el-table-column prop="payload" label="内容 Content" header-align="center"/>
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
