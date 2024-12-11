<script lang="ts" setup>
import {onMounted, ref, watch} from 'vue'

import {CheckboxValueType, ElMessage} from 'element-plus'
import {mdiShare} from "@mdi/js";
import SvgIcon from "@jamescoyle/vue-icon";
import {get, post} from "../api/http";
import {ExportCrawl} from "../../wailsjs/go/main/App";

const protocolCheckAll = ref(false)
const protocolIndeterminate = ref(false)
const protocol = ref<CheckboxValueType[]>([])
const protocols = ref<any[]>([])

const countryCheckAll = ref(false)
const countryIndeterminate = ref(false)
const country = ref<CheckboxValueType[]>([])
const countries = ref<any[]>([])

const nodeCount = ref(0)
const nodeFilterValue = ref(128)

onMounted(async () => {
  const msg = await get<any>("/nodeCache");
  protocols.value = msg.protocol
  countries.value = msg.country
  nodeCount.value = msg.count
  if (msg.count > 128) {
    nodeFilterValue.value = 128
  } else if (msg.count > 64) {
    nodeFilterValue.value = 64
  } else if (msg.count > 32) {
    nodeFilterValue.value = 32
  }
})


watch(protocol, (val) => {
  if (val.length === 0) {
    protocolCheckAll.value = false
    protocolIndeterminate.value = false
  } else if (val.length === protocols.value.length) {
    protocolCheckAll.value = true
    protocolIndeterminate.value = false
  } else {
    protocolIndeterminate.value = true
  }
})

const handleProtocolCheckAll = (val: CheckboxValueType) => {
  protocolIndeterminate.value = false
  if (val) {
    protocol.value = protocols.value.map((_) => _.value)
  } else {
    protocol.value = []
  }
}


watch(country, (val) => {
  if (val.length === 0) {
    countryCheckAll.value = false
    countryIndeterminate.value = false
  } else if (val.length === countries.value.length) {
    countryCheckAll.value = true
    countryIndeterminate.value = false
  } else {
    countryIndeterminate.value = true
  }
})

const handleCountryCheckAll = (val: CheckboxValueType) => {
  countryIndeterminate.value = false
  if (val) {
    country.value = countries.value.map((_) => _.value)
  } else {
    country.value = []
  }
}

const options = [
  {
    value: 32,
    label: 32,
  },
  {
    value: 64,
    label: 64,
  },
  {
    value: 128,
    label: 128,
  },
  {
    value: 256,
    label: 256,
  },
  {
    value: 512,
    label: 512,
  },
]

const nodeFilterCount = ref(0)
const nodeResultShow = ref(true)
const nodeFilterProtocols = ref<any[]>([])
const OverrideFlag = ref(false)
const GenerateFlag = ref(false)

function getReq(option: number) {
  const req = {
    protocol: [],
    country: [],
    count: nodeFilterValue.value,
    option: option
  }

  req["protocol"] = protocol.value.map((k, _) => k) as any
  req["country"] = country.value.map((k, _) => k) as any

  return req
}

async function getFilterNodes() {
  const msg = await post<any>("/nodeFilter", getReq(1));
  nodeFilterProtocols.value = msg.protocol
  nodeFilterCount.value = msg.count

  if (nodeFilterCount.value > 0) {
    nodeResultShow.value = false
    OverrideFlag.value = false
    GenerateFlag.value = false
    ElMessage.success("筛选成功Filter successfully")
  }
}


async function Override() {
  const msg = await post<any>("/nodeFilter", getReq(2));
  if (msg) {
    OverrideFlag.value = true
    ElMessage.success("覆盖默认配置成功 Override Default Config Successfully")
  }
}


async function Generate() {
  const msg = await post<any>("/nodeFilter", getReq(3));
  if (msg) {
    GenerateFlag.value = true
    ElMessage.success("生成新配置成功 Generate New Config Successfully")
  }
}

async function Export() {
  const msg = await post<any>("/nodeFilter", getReq(4));
  if (!msg) {
    ElMessage.error("生成新配置失败 Generate New Config Failed")
    return
  }
  const ok = await ExportCrawl()
  switch (ok) {
    case "false":
      break;
    case "true":
      ElMessage({
        showClose: true,
        message: "导出成功 Export Success",
        type: 'success',
      })
      break;
    default:
      ElMessage({
        showClose: true,
        message: "导出失败 Export failed : " + ok,
        type: 'error',
      })
      break;
  }
}


</script>

<template>
  <div style="width: 650px;">
    <el-text>节点库已缓存节点
      <el-text type="danger">{{ nodeCount }}</el-text>
      个，其中
      <div v-for="item in protocols" :key="item.value + item.count">
        <el-text type="warning">{{ item.value }}</el-text>&nbsp;
        <el-text type="danger"> {{ item.count }}</el-text>&nbsp;个
      </div>
    </el-text>
  </div>
  <br>
  <div>
    <el-text>协议类型 Protocol Type</el-text>&emsp;
    <el-select
        v-model="protocol"
        multiple
        clearable
        placeholder="选择 Select"
        popper-class="custom-header"
        style="width: 650px"
    >
      <template #header>
        <el-checkbox
            v-model="protocolCheckAll"
            :indeterminate="protocolIndeterminate"
            @change="handleProtocolCheckAll"
        >
          全部 All
        </el-checkbox>
      </template>
      <el-option
          v-for="item in protocols"
          :key="item.value"
          :label="item.label"
          :value="item.value"
      />
    </el-select>
  </div>
  <br>
  <div>
    <el-text>国家地区 Country</el-text>&emsp;
    <el-select
        v-model="country"
        multiple
        clearable
        filterable
        placeholder="选择 Select"
        popper-class="custom-header"
        style="width: 689px"
    >
      <template #header>
        <el-checkbox
            v-model="countryCheckAll"
            :indeterminate="countryIndeterminate"
            @change="handleCountryCheckAll"
        >
          全部 All
        </el-checkbox>
      </template>
      <el-option
          v-for="item in countries"
          :key="item.value"
          :label="item.label"
          :value="item.value"
      />
    </el-select>
  </div>
  <br>
  <div>
    <el-text>节点个数 Node Number</el-text>&emsp;
    <el-select v-model="nodeFilterValue" placeholder="Select" style="width: 100px">
      <el-option
          v-for="item in options"
          :key="item.value"
          :label="item.label"
          :value="item.value"
      />
    </el-select>&emsp;
    <el-text>总节点数不会超过 {{ nodeFilterValue }} 个</el-text>&emsp;
    <el-text>The total number of nodes will not exceed {{ nodeFilterValue }}</el-text>
  </div>
  <br>
  <div>
    <el-button type="primary" @click="getFilterNodes">筛选 Filter</el-button>
  </div>
  <br>
  <div :hidden="nodeResultShow">
    <div style="width: 650px;">
      <el-text type="primary">本次总共筛选出节点
        <el-text type="danger">{{ nodeFilterCount }}</el-text>
        个，其中
        <div v-for="item in nodeFilterProtocols" :key="item.value + item.count">
          <el-text type="warning">{{ item.value }}</el-text>&nbsp;
          <el-text type="danger"> {{ item.count }}</el-text>&nbsp;个
        </div>
      </el-text>
    </div>
    <br>
    <div>
      <el-button type="danger" @click="Override" :disabled="OverrideFlag">覆盖默认配置 Override Default Config
      </el-button>
      <el-button type="success" @click="Generate" :disabled="GenerateFlag">生成新配置 Generate New Config</el-button>
      <el-button type="primary" @click="Export">
        <el-icon>
          <svg-icon type="mdi" :path="mdiShare" :size="24"></svg-icon>
        </el-icon>&nbsp;
        导出配置 Export Config
      </el-button>
    </div>
  </div>


</template>


<style scoped>

</style>