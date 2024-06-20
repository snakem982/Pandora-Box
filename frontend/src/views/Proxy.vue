<script lang="ts" setup>
import {onBeforeUnmount, onMounted, reactive, ref, watch} from 'vue'
import {del, get, put} from "../api/http";
import {Group, Proxy} from "../api/pandora";
import Delay from "../weight/Delay.vue";
import {ElLoading} from 'element-plus'
import {mdiArrowUpBold, mdiFlash, mdiSortAlphabeticalAscending} from "@mdi/js";
import SvgIcon from "@jamescoyle/vue-icon";


const isHide = ref(true)
const activeNames: any = ref("")
const labelPosition = ref('')
const groupCache = reactive(new Group([]))

const GROUP: any = reactive({
  proxies: []
})
let PROVIDER: any = reactive({})

function getProxies() {
  GROUP.proxies.length = 0
  let globalProxy: any
  for (let proxy of groupCache.proxies) {
    if (labelPosition.value == "rule" && proxy.name != "GLOBAL" && proxy.name != "❌  默认拦截") {
      GROUP.proxies.push(proxy)
    }

    if (labelPosition.value == "global" && proxy.name == "GLOBAL") {
      GROUP.proxies.push(proxy)
    }

    if (proxy.name == "GLOBAL") {
      globalProxy = proxy
    }
  }
  if (GROUP.proxies.length > 0) {
    if (globalProxy) {
      GROUP.proxies.sort(function (a: any, b: any): number {
        return (globalProxy as Proxy).all.indexOf(a.name) - (globalProxy as Proxy).all.indexOf(b.name)
      })
    }
  }
}

async function getGroup() {
  try {
    let temp = await get<any>("/group");
    groupCache.proxies = temp.proxies
  } catch (error) {
    console.error(error);
  }
}

async function getConfig() {
  try {
    let temp = await get<any>("/configs");
    labelPosition.value = temp.mode
  } catch (error) {
    console.error(error);
  }
}

async function getProvider() {
  try {
    let res = await get<any>("/proxies");
    PROVIDER = res.proxies
  } catch (error) {
    console.error(error);
  }
}

let interval = 0

const flushDelay = async () => {
  for (let i = 0; i < GROUP.proxies.length; i++) {
    if (GROUP.proxies[i].name == activeNames.value) {
      await getProvider()
      GROUP.proxies[i].all.length = 0
      setTimeout(() => {
        GROUP.proxies[i].all.push(...PROVIDER[activeNames.value].all)
      }, 0)
      break
    }
  }
}

onBeforeUnmount(() => {
  window.clearInterval(interval)
})

onMounted(async () => {
  isHide.value = localStorage.getItem('isHide') == "0";
  isSort.value = localStorage.getItem('isSort') == "0";
  isGrid.value = localStorage.getItem('isGrid') == "0" ? 6 : 8;
  const acName = localStorage.getItem("activeNames")

  await getConfig()
  await getProvider()
  await getGroup()
  getProxies()
  activeNames.value = acName

  interval = setInterval(flushDelay, 15000)
})

watch(labelPosition, getProxies)

function cardStyle(proxy: Proxy, name: string): string {
  return proxy.now === name ? "grid-content-select" : "grid-content"
}

function cardShadow(proxy: Proxy, name: string): any {
  return proxy.now === name ? "always" : "hover"
}

function colShow(proxy: any, hide: boolean): boolean {
  if (proxy['name'] == "❌  默认拦截") {
    return false
  }
  if (hide) {
    return true
  }
  const type = proxy['type'];
  return proxy['alive'] || type == 'Reject' || type == 'URLTest' || type == 'Fallback' || type == 'LoadBalance'
}

async function setProxy(proxy: Proxy, name: string) {
  if (proxy.now === name) {
    return
  }
  try {
    await del("/connections")
    await put("/proxies/" + proxy.name, {name});
    await getGroup()
    getProxies()
  } catch (error) {
    console.error(error);
  }
}

const isGrid = ref(8)

function setGrid() {
  if (isGrid.value == 8) {
    isGrid.value = 6
  } else {
    isGrid.value = 8
  }
  localStorage.setItem("isGrid", isGrid.value == 6 ? "0" : "1")
}

function setHide() {
  isHide.value = !isHide.value
  localStorage.setItem("isHide", isHide.value ? "0" : "1")
}

function goTop() {
  location.href = "#goTop"
}


async function getDelay() {
  if (GROUP.proxies.length == 0) {
    return
  }

  if (activeNames.value === "") {
    activeNames.value = GROUP.proxies[0].name
  }
  const tempName = activeNames.value
  window.clearInterval(interval)
  const loading = ElLoading.service({
    lock: true,
    text: '测试中Testing...',
    background: 'rgba(0, 0, 0, 0.7)',
  })
  activeNames.value = ""
  try {
    const param = "/delay?url=https%3A%2F%2Fwww.gstatic.com%2Fgenerate_204&timeout=3000"
    await get("/group/" + tempName + param);
    await getProvider()
  } catch (error) {
    console.error(error);
  }
  activeNames.value = tempName
  loading.close()
  interval = setInterval(flushDelay, 15000)
}

function getLocal() {
  if (GROUP.proxies.length == 0) {
    return
  }

  if (activeNames.value === "") {
    activeNames.value = GROUP.proxies[0].name
    setTimeout(function () {
      location.href = "#" + activeNames.value + PROVIDER[activeNames.value]['now']
    }, 500)
  } else {
    location.href = "#" + activeNames.value + PROVIDER[activeNames.value]['now']
  }
}

function allDelay(cao: any): number {
  let history: any = cao.history
  let type: any = cao.type
  const types = ["URLTest", "Direct", "Selector", "Reject", "Reject", "Fallback", "LoadBalance"]
  if (types.includes(type)) {
    return -9999
  }

  if (history.length > 0) {
    const delay = history[history.length - 1]['delay'];
    return delay > 0 ? delay : 99999
  }

  return 9999
}

const isSort = ref(false)

function changeSortFlag() {
  isSort.value = !isSort.value
  localStorage.setItem("isSort", isSort.value ? "0" : "1")
  const temp = activeNames.value
  if (activeNames.value != "") {
    activeNames.value = ""
  }
  setTimeout(() => activeNames.value = temp, 50)
}

function reSort(all: any): any {
  let ok: any = []
  if (all.length == 0) {
    return ok
  }
  ok = ok.concat(all);
  if (isSort.value) {
    ok.sort(function (a: any, b: any): number {
      return allDelay(PROVIDER[a]) - allDelay(PROVIDER[b])
    })
  }

  return ok
}

watch(activeNames, () => {
  localStorage.setItem("activeNames", activeNames.value)
})

function fType(type: any): any {
  if (type == "Shadowsocks") {
    return "SS"
  }

  if (type == "ShadowsocksR") {
    return "SSr"
  }

  return type
}

</script>

<template>
  <el-affix :offset="20" id="goTop">
    <el-row>
      <el-button @click="setGrid" type="info" icon="Grid" circle title="分栏Columns"/>
      <el-button @click="setHide" type="primary" icon="Hide" circle
                 title="隐藏不可用节点Hide Unavailable Proxies"/>
      <el-button @click="getDelay" type="success" circle
                 title="延迟测速Delay Test">
        <svg-icon type="mdi" :path="mdiFlash" :size="18"></svg-icon>
      </el-button>
      <el-button @click="getLocal" type="danger" icon="Aim" circle title="当前节点Current Proxy"/>
      <el-button @click="changeSortFlag" type="warning" circle title="排序Sort">
        <svg-icon type="mdi" :path="mdiSortAlphabeticalAscending" :size="18"></svg-icon>
      </el-button>
      <el-button @click="goTop" type="primary" circle title="返回顶部Back To Top">
        <svg-icon type="mdi" :path="mdiArrowUpBold" :size="18"></svg-icon>
      </el-button>
    </el-row>
  </el-affix>

  <div class="proxy-collapse">
    <el-collapse accordion
                 v-model="activeNames"
                 v-for="(item,index) in GROUP.proxies"
                 :key="index"
    >
      <el-collapse-item
          :title="item.name"
          :name="item.name"
      >
        <el-row v-if="activeNames==item.name">
          <el-col :span="isGrid" v-for="(all,index) in reSort(item.all)"
                  :key="'all' + index"
                  v-show="colShow(PROVIDER[all],isHide)"
                  @click="setProxy(item,all)"
                  :id="item.name+all"
          >
            <el-card :shadow="cardShadow(item,all)" :class="cardStyle(item,all)">
              <el-text size="large" truncated>
                {{ all }}
              </el-text>
              <br>
              <el-text size="large" type="info" style="margin-right: 20px">{{ fType(PROVIDER[all]['type']) }}</el-text>
              <delay :proxy="PROVIDER[all]"></delay>
            </el-card>
          </el-col>
        </el-row>
      </el-collapse-item>
    </el-collapse>
  </div>
</template>

<style scoped>

.proxy-mode {
  width: 300px;
  margin: auto;
}

.proxy-collapse {
  margin-top: 25px;
}

:deep(.el-collapse-item__header) {
  font-size: 18px;
}

:deep(.el-card__body) {
  padding: 10px;
}

.grid-content {
  border-radius: 4px;
  min-width: 150px;
  margin: 10px;
  font-size: 18px;
}

.grid-content:hover {
  cursor: pointer;
}

.grid-content-select {
  border-radius: 4px;
  min-width: 150px;
  margin: 10px;
  font-size: 18px;
  border-left-style: solid;
  border-left-color: #75e043;
  border-left-width: 5px;
}
</style>


