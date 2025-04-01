<script setup lang="ts">
import {useRouter} from 'vue-router'

const menu = reactive({
  homeStatus: 1,
  setStatus: 1,
  proxyStatus: 1,
  subStatus: 1,
})

const router = useRouter();

function active(btn: string) {
  let temp = menu[btn]
  if (temp == 3 || temp == 23) {
    return
  }

  menu["homeStatus"] = 1
  menu["setStatus"] = 1
  menu["proxyStatus"] = 1
  menu["subStatus"] = 1

  if (temp == 22) {
    menu[btn] = 23
  } else {
    menu[btn] = 3
  }

  switch (btn) {
    case "homeStatus":
      router.push('/Home')
      break;
    case "setStatus":
      router.push('/Setting')
      break;
    case "proxyStatus":
      router.push('/Proxies')
      break;
    case "subStatus":
      router.push('/Profiles')
      break;
    default:
  }


}

function enter(btn: string) {
  if (menu[btn] === 1) {
    menu[btn] = 22
  } else {
    menu[btn] = 23
  }
}

function leave(btn: string) {
  if (menu[btn] === 22) {
    menu[btn] = 1
  } else if (menu[btn] === 23) {
    menu[btn] = 3
  }
}

</script>

<template>
  <div class="button-container">
    <button
        @click="active('homeStatus')"
        @mouseenter="enter('homeStatus')"
        @mouseleave="leave('homeStatus')"
        :class="{ active: menu.homeStatus==3 || menu.homeStatus==23 }"
    >
      <template v-if="menu.homeStatus==3">
        <icon-mdi-home/>
      </template>
      <template v-else-if="menu.homeStatus==22 || menu.homeStatus==23">
        主页
      </template>
      <template v-else>
        <icon-mdi-home-outline/>
      </template>
    </button>
    <button
        @click="active('setStatus')"
        @mouseenter="enter('setStatus')"
        @mouseleave="leave('setStatus')"
        :class="{ active: menu.setStatus==3 || menu.setStatus==23 }"
    >
      <template v-if="menu.setStatus==3">
        <icon-mdi-cog/>
      </template>
      <template v-else-if="menu.setStatus==22 || menu.setStatus==23">
        设置
      </template>
      <template v-else>
        <icon-mdi-cog-outline/>
      </template>
    </button>
    <button
        @click="active('proxyStatus')"
        @mouseenter="enter('proxyStatus')"
        @mouseleave="leave('proxyStatus')"
        :class="{ active: menu.proxyStatus==3 || menu.proxyStatus==23 }"
    >
      <template v-if="menu.proxyStatus==3">
        <icon-mdi-rocket-launch/>
      </template>
      <template v-else-if="menu.proxyStatus==22 || menu.proxyStatus==23">
        代理
      </template>
      <template v-else>
        <icon-mdi-rocket-launch-outline/>
      </template>
    </button>
    <button
        @click="active('subStatus')"
        @mouseenter="enter('subStatus')"
        @mouseleave="leave('subStatus')"
        :class="{ active: menu.subStatus==3 || menu.subStatus==23 }"
    >
      <template v-if="menu.subStatus==3">
        <icon-mdi-arrange-bring-forward/>
      </template>
      <template v-else-if="menu.subStatus==22 || menu.subStatus==23">
        订阅
      </template>
      <template v-else>
        <icon-mdi-arrange-send-backward/>
      </template>
    </button>
  </div>
</template>

<style scoped>
.button-container {
  margin-top: 20px;
  margin-left: 22px;
  display: grid;
  grid-template-columns: 1fr 1fr;
  column-gap: 8px;
  row-gap: 8px;
  width: 193px;
}

.button-container button {
  padding: 13px 20px;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-size: 16px;
  width: 90px;
  background-color: var(--left-nav-btn-bg);
  color: var(--text-color);
  box-shadow: var(--left-nav-shadow);
}

.button-container button:hover {
  background-color: var(--left-nav-btn-hover-bg);
  color: var(--text-color);
  font-size: 14px;
  box-shadow: var(--left-nav-hover-shadow);
}

.button-container button.active {
  background-color:var(--left-nav-btn-active-bg);
  color: var(--text-color);
}
</style>