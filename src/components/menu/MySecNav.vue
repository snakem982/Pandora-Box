<script setup lang="ts">
import {changeMenu} from "@/util/menu";
import {useMenuStore} from "@/store/menuStore";
import {useRouter} from "vue-router";
import {WS} from "@/util/ws";
import {useWebStore} from "@/store/webStore";
import createApi from "@/api";
import {formatDate} from "@/util/format";

// 获取当前 Vue 实例的 proxy 对象 和 api
const {proxy} = getCurrentInstance()!;
const api = createApi(proxy);

// 获取Store
const menuStore = useMenuStore();
const webStore = useWebStore();

// 获取路由
const router = useRouter();

const conn = ref(0);

// 连接数
function onConn(ev: MessageEvent) {
  const parsedData = JSON.parse(ev.data);
  if (parsedData["connections"]) {
    conn.value = parsedData["connections"].length;
  } else {
    conn.value = 0;
  }
}

// 日志
function onLog(ev: MessageEvent) {
  const parsedData = JSON.parse(ev.data);
  webStore.addLog({
    time: formatDate(new Date()),
    type: parsedData["type"].toUpperCase(),
    payload: parsedData["payload"]
  });
}

function aliveTest(conn: WS, cb: Function) {
  setInterval(() => {
    try {
      if (conn.ws.readyState === WebSocket.OPEN) {
        conn.ws.send("ping");
      } else {
        console.log("WebSocket 连接可能已断开");
        if (cb) {
          cb()
        }
      }
    } catch (error) {
      console.error("发送失败，WebSocket 可能已断开:", error);
      if (cb) {
        cb()
      }
    }
  }, 10000);
}

let wsConn: WS;
let logConn: WS;
onMounted(() => {
  const urlTraffic = webStore.wsUrl + "/connections?token=" + webStore.secret;
  wsConn = new WS(urlTraffic, null, onConn);
  aliveTest(wsConn, () => {
    wsConn.close()
    wsConn = new WS(urlTraffic, null, onConn);
  })

  const logTraffic = webStore.wsUrl + "/logs?token=" + webStore.secret;
  logConn = new WS(logTraffic, null, onLog);
  aliveTest(logConn, () => {
    logConn.close()
    logConn = new WS(logTraffic, null, onLog);
  })

  api.getRuleNum().then((res) => {
    menuStore.setRuleNum(res);
  });
});

</script>

<template>
  <div class="nav">
    <div
        :class="menuStore.menu == 'Rule' ? 'nav-btn nav-btn-select' : 'nav-btn'"
        @click="changeMenu('Rule', router)"
    >
      <el-text class="nav-text">
        <el-icon>
          <icon-mdi-source-branch/>
        </el-icon>
        <span class="nav-info"
        >{{ $t("sec-nav.rule") }} · {{ menuStore.ruleNum }}</span
        >
      </el-text>
    </div>

    <div
        :class="
        menuStore.menu == 'Connection' ? 'nav-btn nav-btn-select' : 'nav-btn'
      "
        @click="changeMenu('Connection', router)"
    >
      <el-text class="nav-text">
        <el-icon>
          <icon-mdi-lan-connect/>
        </el-icon>
        <span class="nav-info">{{ $t("sec-nav.conn") }} · {{ conn }}</span>
      </el-text>
    </div>

    <div
        :class="menuStore.menu == 'Log' ? 'nav-btn nav-btn-select' : 'nav-btn'"
        @click="changeMenu('Log', router)"
    >
      <el-text class="nav-text">
        <el-icon>
          <icon-mdi-text-box-outline/>
        </el-icon>
        <span class="nav-info">{{ $t("sec-nav.log") }}</span>
      </el-text>
    </div>

    <!--    <div-->
    <!--        :class="menuStore.menu == 'Crawl' ? 'nav-btn nav-btn-select' : 'nav-btn'"-->
    <!--        @click="changeMenu('Crawl', router)"-->
    <!--    >-->
    <!--      <el-text class="nav-text">-->
    <!--        <el-icon>-->
    <!--          <icon-mdi-spider-outline/>-->
    <!--        </el-icon>-->
    <!--        <span class="nav-info">{{ $t("sec-nav.crawl") }} · 530</span>-->
    <!--      </el-text>-->
    <!--    </div>-->
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

.nav-btn-select,
.nav-btn:hover {
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
