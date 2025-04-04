<script setup lang="ts">
import MyHr from "@/components/proxies/MyHr.vue";
import {useMenuStore} from "@/store/menuStore";
import {useRouter} from "vue-router";

const menuStore = useMenuStore()
const router = useRouter()

const distanceFromTop = ref(195)
const upFromTop = function (distance: number) {
  distanceFromTop.value = distance
}

function changeRuleMenu(active: string): void {
  menuStore.setRuleMenu(active)
  router.push(active);
}

const getActive = function (value: string): string {
  return menuStore.ruleMenu === value ? 'proxy-group-title proxy-group-title-select' : 'proxy-group-title';
}

const setActive = function (value: string) {
  changeRuleMenu(value)
}

onMounted(() => {
  changeRuleMenu(menuStore.ruleMenu)
})

</script>

<template>
  <MyLayout :top-height="distanceFromTop-15" :bottom-height="distanceFromTop+25">
    <template #top>
      <MySearch></MySearch>
      <el-space>
        <div class="title">规则</div>
      </el-space>

      <div class="proxy-group">
        <button
            :class="getActive('Now')"
            @click="setActive('Now')"
        >
          <icon-mdi-eye-arrow-right class="pre"/>
          <span class="suf">查看目前规则</span>
        </button>
        <button
            :class="getActive('Group')"
            @click="setActive('Group')"
        >
          <icon-mdi-view-dashboard class="pre"/>
          <span class="suf">统一规则分组</span>
        </button>
        <button
            :class="getActive('Ignore')"
            @click="setActive('Ignore')"
        >
          <icon-mdi-cancel class="pre"/>
          <span class="suf">忽略这些域名</span>
        </button>
      </div>

      <MyHr :update="upFromTop"></MyHr>
    </template>
    <template #bottom>
      <router-view/>
    </template>
  </MyLayout>
</template>

<style scoped>
.title {
  margin-top: 25px;
  font-size: 32px;
  font-weight: bold;
  margin-left: 38px;
}

.proxy-group {
  display: flex;
  margin-top: 14px;
  margin-left: 36px;
  flex-wrap: wrap;
  gap: 10px;
  width: 92%;
}
</style>