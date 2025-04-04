<script setup lang="ts">
import {changeMenu} from "@/util/menu";
import {useMenuStore} from "@/store/menuStore";
import {useRouter} from "vue-router";

const menuStore = useMenuStore()
const router = useRouter()

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
        <span class="nav-info">{{ $t('sec-nav.rule') }} · 386</span>
      </el-text>
    </div>

    <div
        :class="menuStore.menu=='Connection'? 'nav-btn nav-btn-select':'nav-btn'"
        @click="changeMenu('Connection',router)">
      <el-text class="nav-text">
        <el-icon>
          <icon-mdi-lan-connect/>
        </el-icon>
        <span class="nav-info">{{ $t('sec-nav.conn') }} · 257</span>
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