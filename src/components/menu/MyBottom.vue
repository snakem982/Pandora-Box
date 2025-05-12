<script setup lang="ts">
import {useMenuStore} from "@/store/menuStore";
import {useRouter} from "vue-router";
import Language from "@/components/menu/Language.vue";
import Off from "@/components/menu/Off.vue";
import Skin from "@/components/menu/Skin.vue";

const menuStore = useMenuStore()
const router = useRouter()


onBeforeMount(() => {
  // 上次打开页面
  if (menuStore.path) {
    router.push(menuStore.path)
  }
})
</script>

<template>
  <div class="bottom-text">

    <Off></Off>
    <Language></Language>
    <Skin></Skin>

  </div>
</template>

<style scoped>
.bottom-text {
  position: absolute;
  bottom: 18px;
  left: 22px;
  color: var(--text-color);
  font-size: 20px;
}
</style>