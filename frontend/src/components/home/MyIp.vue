<script setup lang="ts">

</script>

<template>
  <el-row :gutter="20" class="spark"
          style="margin-left: 2px;">
    <el-col :span="12">
      <div class="box box1">
        <div class="title">
          {{ $t('home.ip.title') }}
          <el-icon size="22" style="position: absolute;margin-left: 8px;">
            <icon-mdi-refresh/>
          </el-icon>
        </div>
        <hr/>
        <ul class="info-list">
          <li><strong>{{ $t('home.ip.real') }} : </strong> 1.2.3.4</li>
          <li><strong>{{ $t('home.ip.city') }} : </strong> 圣何塞</li>
          <li><strong>{{ $t('home.ip.country') }} : </strong> 甲骨文</li>
          <li><strong>{{ $t('home.ip.isp') }} : </strong> 甲骨文 BMC</li>
          <li><strong>{{ $t('home.ip.asn') }} : </strong> AS31898</li>
          <li><strong>{{ $t('home.ip.time-zone') }} : </strong> 美国/落山鸡</li>
        </ul>
      </div>
    </el-col>

    <el-col :span="12">
      <div class="box box2">
        <div class="title">
          {{ $t('home.system.title') }}
        </div>
        <hr/>
        <ul class="info-list">
          <li><strong>{{ $t('home.system.os') }} : </strong> Windows</li>
          <li><strong>{{ $t('home.system.runtime') }} : </strong> 1:06:42</li>
          <li><strong>{{ $t('home.system.startup') }} : </strong> 已启用</li>
          <li><strong>{{ $t('home.system.service') }} : </strong> 服务模式</li>
          <li><strong>{{ $t('home.system.port') }} : </strong> 12345</li>
          <li><strong>{{ $t('home.system.version') }} : </strong> v1.0.1</li>
        </ul>
      </div>
    </el-col>
  </el-row>
</template>

<style scoped>
.spark {
  max-width: 95%;
  margin-top: 30px;
}

.box {
  padding: 10px;
  border-radius: 8px;
  text-align: left;
}

.box hr {
  border: none;
  height: 1px;
  background-color: var(--hr-color);
  margin: 10px 0;
}

.info-list {
  list-style: none;
  padding: 0;
}

.info-list li {
  font-size: 18px;
  margin: 8px 0;
}

.box1 {
  box-shadow: var(--right-box-shadow);
}

.box2 {
  box-shadow: var(--right-box-shadow);
}
</style>