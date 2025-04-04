<script setup lang="ts">
import {useMenuStore} from "@/store/menuStore";
import {changeMenu} from "@/util/menu";
import {useRouter} from "vue-router";

const menuStore = useMenuStore()
const router = useRouter()

const menu = reactive({
  homeStatus: false,
  setStatus: false,
  proxyStatus: false,
  subStatus: false,
})

function enter(btn: string) {
  menu[btn] = true
}

function leave(btn: string) {
  menu[btn] = false
}

</script>

<template>
  <div class="button-container">
    <button
        @mouseenter="enter('homeStatus')"
        @mouseleave="leave('homeStatus')"
        @click="changeMenu('Home',router)"
        :class="{ active: menuStore.menu=='Home' }"
    >
      <template v-if="menu.homeStatus">
        主页
      </template>
      <template v-else>
        <template v-if="menuStore.menu=='Home'">
          <icon-mdi-home/>
        </template>
        <template v-else>
          <icon-mdi-home-outline/>
        </template>
      </template>
    </button>

    <button
        @mouseenter="enter('setStatus')"
        @mouseleave="leave('setStatus')"
        @click="changeMenu('Setting',router)"
        :class="{ active: menuStore.menu=='Setting' }"
    >
      <template v-if="menu.setStatus">
        设置
      </template>
      <template v-else>
        <template v-if="menuStore.menu=='Setting'">
          <icon-mdi-cog/>
        </template>
        <template v-else>
          <icon-mdi-cog-outline/>
        </template>
      </template>
    </button>

    <button
        @mouseenter="enter('proxyStatus')"
        @mouseleave="leave('proxyStatus')"
        @click="changeMenu('Proxies',router)"
        :class="{ active: menuStore.menu=='Proxies' }"
    >
      <template v-if="menu.proxyStatus">
        设置
      </template>
      <template v-else>
        <template v-if="menuStore.menu=='Proxies'">
          <icon-mdi-rocket-launch/>
        </template>
        <template v-else>
          <icon-mdi-rocket-launch-outline/>
        </template>
      </template>
    </button>
    <button
        @mouseenter="enter('subStatus')"
        @mouseleave="leave('subStatus')"
        @click="changeMenu('Profiles',router)"
        :class="{ active: menuStore.menu=='Profiles' }"
    >
      <template v-if="menu.subStatus">
        订阅
      </template>
      <template v-else>
        <template v-if="menuStore.menu=='Profiles'">
          <icon-mdi-arrange-bring-forward/>
        </template>
        <template v-else>
          <icon-mdi-arrange-send-backward/>
        </template>
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
  background-color: var(--left-nav-btn-active-bg);
  color: var(--text-color);
}
</style>