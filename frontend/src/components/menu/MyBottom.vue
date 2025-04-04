<script setup lang="ts">
import {useMenuStore} from "@/store/menuStore";
import {useRouter} from "vue-router";
import Language from "@/components/menu/Language.vue";

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
    <el-icon>
      <icon-mdi-power/>
    </el-icon>

    <Language></Language>

    <el-icon style="margin-left:20px">
      <icon-mdi-tshirt-crew-outline/>
    </el-icon>
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