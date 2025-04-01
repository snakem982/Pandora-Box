<script setup lang="ts">

</script>

<template>
  <div class="nav">
    <div class="nav-btn">
      <el-text class="nav-text">
        <el-icon>
          <icon-mdi-source-branch/>
        </el-icon>
        <span class="nav-info">规则 · 386</span>
      </el-text>
    </div>

    <div class="nav-btn">
      <el-text class="nav-text">
        <el-icon>
          <icon-mdi-connection/>
        </el-icon>
        <span class="nav-info">连接 · 257</span>
      </el-text>
    </div>

    <div class="nav-btn">
      <el-text class="nav-text">
        <el-icon>
          <icon-mdi-text-box-outline/>
        </el-icon>
        <span class="nav-info">日志</span>
      </el-text>
    </div>

    <div class="nav-btn">
      <el-text class="nav-text">
        <el-icon>
          <icon-mdi-spider-outline/>
        </el-icon>
        <span class="nav-info">爬取 · 530</span>
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