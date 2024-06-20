<script setup lang="ts">

import {BrowserOpenURL} from "../../wailsjs/runtime";
</script>

<template>
  <el-card class="box-card">
    <template #header>
      <div class="card-header">
        <span>关于 Pandora-Box</span>
      </div>
    </template>
    <div>
      <el-text class="el-text--primary title">1、开发初衷 Original development intention</el-text>
      <div class="content">
        &emsp;
        <el-text>一切为了简单 &emsp;JUST FOR SIMPLE</el-text>
      </div>
    </div>

    <div style="margin-top: 20px">
      <el-text class="el-text--primary title">2、依赖开源项目 Dependent open source projects</el-text>
      <div class="content">
        &emsp;
        <el-popover
            placement="bottom"
            :width="400"
            trigger="hover"
            content="Build beautiful cross-platform applications using Go"
        >
          <template #reference>
            <el-link @click="BrowserOpenURL('https://wails.io')" target="_blank">Wails</el-link>
          </template>
        </el-popover>
        &emsp;
        <el-popover
            placement="bottom"
            :width="180"
            trigger="hover"
            content="An awesome kernel"
        >
          <template #reference>
            <el-link @click="BrowserOpenURL('https://github.com/MetaCubeX/mihomo')" target="_blank">Mihomo</el-link>
          </template>
        </el-popover>
        &emsp;
        <el-popover
            placement="bottom"
            :width="460"
            trigger="hover"
            content="A Vue3 based component library for designers and developers"
        >
          <template #reference>
            <el-link @click="BrowserOpenURL('https://element-plus.org')" target="_blank">Element-Plus</el-link>
          </template>
        </el-popover>
      </div>
    </div>

    <div style="margin-top: 20px">
      <el-text class="el-text--primary title">3、免责声明 Disclaimer</el-text>
      <div class="content">
        &emsp;
        <el-text>该软件仅供个人和教育用途。 <br>&emsp;
          使用本软件时请遵守相关法律法规。 <br>&emsp;
          软件所有者对用户行为不承担任何责任。
        </el-text>
        <br><br>&emsp;
        <el-text>This software is meant for personal and educational uses only. <br>&emsp;
          Please follow relevant laws and regulations when using this software. <br>&emsp;
          The software owner assumes no responsibility for user behavior.
        </el-text>
      </div>
    </div>

  </el-card>
</template>

<style scoped>
.title {
  font-size: 17px;
}

.content {
  margin-top: 10px;
  font-size: 16px;
}

</style>