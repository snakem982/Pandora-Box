<script setup lang="ts">
import {changeMenu} from "@/util/menu";
import {useMenuStore} from "@/store/menuStore";
import {useRouter} from "vue-router";
import {WS} from "@/util/ws";
import {useWebStore} from "@/store/webStore";
import createApi from "@/api";


// 获取当前 Vue 实例的 proxy 对象 和 api
const {proxy} = getCurrentInstance()!;
const api = createApi(proxy);

// 获取Store
const menuStore = useMenuStore()
const webStore = useWebStore()

// 获取路由
const router = useRouter()

const conn = ref(0)

function onConn(ev: MessageEvent) {
  const parsedData = JSON.parse(ev.data);
  conn.value = parsedData['connections'].length
}

let wsConn: WS
onMounted(()=>{
  const urlTraffic = webStore.wsUrl + "/connections?token=" + webStore.secret;
  wsConn = new WS(urlTraffic, null, onConn);


  api.getRules().then((res) => {
    menuStore.setRuleNum(res.length) 
  })
})

onUnmounted(()=>{
  wsConn.close()
})


</script>

<template>
  <div class="nav">
    <div
        :class="menuStore.menu=='Rule'? 'nav-btn nav-btn-select':'nav-btn'"
        @click="changeMenu('Rule',router)">
      <el-text class="nav-text">
        <el-icon>
          <icon-mdi-source-branch/>
        </el-icon>
        <span class="nav-info">{{ $t('sec-nav.rule') }} · {{ menuStore.ruleNum }}</span>
      </el-text>
    </div>

    <div
        :class="menuStore.menu=='Connection'? 'nav-btn nav-btn-select':'nav-btn'"
        @click="changeMenu('Connection',router)">
      <el-text class="nav-text">
        <el-icon>
          <icon-mdi-lan-connect/>
        </el-icon>
        <span class="nav-info">{{ $t('sec-nav.conn') }} · {{ conn }}</span>
      </el-text>
    </div>

    <div
        :class="menuStore.menu=='Log'? 'nav-btn nav-btn-select':'nav-btn'"
        @click="changeMenu('Log',router)">
      <el-text class="nav-text">
        <el-icon>
          <icon-mdi-text-box-outline/>
        </el-icon>
        <span class="nav-info">{{ $t('sec-nav.log') }}</span>
      </el-text>
    </div>

    <div
        :class="menuStore.menu=='Crawl'? 'nav-btn nav-btn-select':'nav-btn'"
        @click="changeMenu('Crawl',router)">
      <el-text class="nav-text">
        <el-icon>
          <icon-mdi-spider-outline/>
        </el-icon>
        <span class="nav-info">{{ $t('sec-nav.crawl') }} · 530</span>
      </el-text>
    </div>

  </div>
</template>

<style scoped>
.nav {
  margin-top: 18px;
  margin-left: 22px;
}

.nav-btn {
  padding: 11px;
}

.nav-btn-select, .nav-btn:hover {
  background-color: var(--left-sec-nav-hover-bg);
  width: 164px;
  border-radius: 8px;
  cursor: pointer;
  box-shadow: var(--left-sec-nav-hover-shadow);
}

.nav-text {
  color: var(--text-color);
  font-size: 18px;
}

.nav-info {
  font-size: 14px;
  margin-left: 12px;
}

</style>