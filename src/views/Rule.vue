<script setup lang="ts">
import MyHr from "@/components/MyHr.vue";
import {useMenuStore} from "@/store/menuStore";
import {useRouter} from "vue-router";

const menuStore = useMenuStore()
const router = useRouter()

const getActive = function (value: string): string {
  return menuStore.ruleMenu === value ? 'proxy-group-title proxy-group-title-select' : 'proxy-group-title';
}

const setActive = function (value: string) {
  router.push("/Rule/" + value);
}
</script>

<template>
  <MyLayout hr-show>
    <template #top>
      <el-space class="space">
        <div class="title">
          {{ $t('rule.title') }}
        </div>
      </el-space>

      <div class="proxy-group">
        <button
            :class="getActive('Now')"
            @click="setActive('Now')"
        >
          <icon-mdi-eye-arrow-right class="pre"/>
          <span class="suf">
            {{ $t('rule.now.title') }}
          </span>
        </button>
        <button
            :class="getActive('Group')"
            @click="setActive('Group')"
        >
          <icon-mdi-view-dashboard class="pre"/>
          <span class="suf">
            {{ $t('rule.group.title') }}
          </span>
        </button>
        <button
            :class="getActive('Ignore')"
            @click="setActive('Ignore')"
        >
          <icon-mdi-cancel class="pre"/>
          <span class="suf">
            {{ $t('rule.ignore.title') }}
          </span>
        </button>
      </div>
    </template>
    <template #bottom>
      <router-view/>
    </template>
  </MyLayout>
</template>

<style scoped>
.space {
  margin-top: 20px;
}

.title {
  font-size: 32px;
  font-weight: bold;
  margin-left: 10px;
}

.proxy-group {
  display: flex;
  margin-top: 20px;
  margin-left: 10px;
  flex-wrap: wrap;
  gap: 10px;
  width: 95%;
}

.pre {
  position: absolute;
}

.suf {
  margin-left: 25px;
}

.proxy-group-title {
  background-color: transparent;
  color: var(--text-color);
  border: 2px solid var(--hr-color);
  border-radius: 8px;
  padding: 6px 8px 6px 8px;
  font-size: 16px;
  text-align: center;
  cursor: pointer;
  box-shadow: var(--left-nav-shadow);
}

.proxy-group-title:hover, .proxy-group-title-select {
  background-color: var(--left-item-selected-bg);
  box-shadow: var(--left-nav-hover-shadow);
  border-color: var(--text-color);
}
</style>