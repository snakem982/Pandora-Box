<script lang="ts" setup>
import {useRouter} from 'vue-router'
import {onMounted, ref, watch} from 'vue'
import {get, patch, put} from "./api/http";
import {WindowSetDarkTheme} from "../wailsjs/runtime";
import {IsAdmin} from "../wailsjs/go/main/App";

const active = ref("")
const version = ref("v0.2.0")

async function getVersion() {
  try {
    const temp = await get<any>("/version")
    version.value = temp.version
  } catch (error) {
    console.error(error);
  }
}

const router = useRouter()
watch(() => router.currentRoute.value.path, (toPath) => {
  if (toPath === "/") {
    return
  }
  //要执行的方法
  if (active.value != toPath) {
    active.value = toPath
    localStorage.setItem("ActiveHome", toPath)
  }
}, {immediate: true, deep: true})

onMounted(async () => {
  await getVersion()

  const isDark = localStorage.getItem("vueuse-color-scheme") == "dark"
  if (isDark) {
    WindowSetDarkTheme()
  }

  const item = localStorage.getItem("ActiveHome");
  if (item) {
    await router.push(item)
    if (item != "/general") {
      const system_proxy_port = localStorage.getItem("system_proxy_port");
      if (system_proxy_port) {
        await patch("/configs", {'mixed-port': Number(system_proxy_port)})
        const system_proxy = localStorage.getItem("system_proxy") == "0"
        if (system_proxy) {
          await put<any>("/system/" + system_proxy_port)
        }
      }

      const isAdmin = await IsAdmin()
      if (isAdmin == "true") {
        const tun = localStorage.getItem("tun")
        if (tun != "off") {
          await patch("/configs", {tun: {enable: false}})
          await patch("/configs", {tun: {enable: true, "stack": tun}})
        }
      }

    }
  } else {
    await router.push("/general")
  }
})

</script>

<template>
  <el-container>
    <el-aside class="ch aside" style="widows: 1">
      <h5 class="title">Pandora-Box</h5>
      <el-menu :default-active="active" router>
        <el-menu-item index="/general">
          <el-icon>
            <Setting/>
          </el-icon>
          <span>通用 GENERAL</span>
        </el-menu-item>
        <el-menu-item index="/proxy">
          <el-icon>
            <Promotion/>
          </el-icon>
          <span>节点 PROXIES</span>
        </el-menu-item>
        <el-menu-item index="/profile">
          <el-icon>
            <Files/>
          </el-icon>
          <span>配置 PROFILES</span>
        </el-menu-item>
        <el-menu-item index="/crawl">
          <el-icon>
            <MagicStick/>
          </el-icon>
          <span>抓取 CRAWL</span>
        </el-menu-item>
        <el-menu-item index="/connection">
          <el-icon>
            <Connection/>
          </el-icon>
          <span>连接 CONNECTION</span>
        </el-menu-item>
        <el-menu-item index="/rule">
          <el-icon>
            <Guide/>
          </el-icon>
          <span>规则 RULE</span>
        </el-menu-item>
        <el-menu-item index="/log">
          <el-icon>
            <Tickets/>
          </el-icon>
          <span>日志 LOG</span>
        </el-menu-item>
        <el-menu-item index="/about">
          <el-icon>
            <InfoFilled/>
          </el-icon>
          <span>关于 ABOUT</span>
        </el-menu-item>
      </el-menu>

      <el-space direction="vertical" style="position: absolute;
        bottom: 30px;
        left: 30px;
        width: 180px;
        text-align: center;
        ">
        <el-text>
          JUST FOR SIMPLE
        </el-text>
        <el-text>
          版本号：{{ version }}
        </el-text>
      </el-space>
    </el-aside>
    <el-main class="ch main">
      <router-view></router-view>
    </el-main>
  </el-container>
</template>

<style scoped>

.title {
  background-image: url("./assets/images/appicon.png");
  background-repeat: no-repeat;
  background-size: 45px;
  background-position-x: 20px;
  line-height: 45px;
  padding-left: 70px;
  font-size: 20px;
}

.title:hover {
  cursor: move;
}

.ch {
  height: 100vh;
}

.aside {
  width: 230px;
  border-right: solid 1px var(--el-menu-border-color);
}

.main {
  //background-color: #ecf2ec;
}

.is-active {
  border-left-style: solid;
  border-left-width: 5px;
}

.el-menu {
  border-right: none;
}
</style>
