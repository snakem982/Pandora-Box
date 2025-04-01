<script setup lang="ts">

</script>

<template>
  <el-row :gutter="20" class="spark"
          style="margin-left: 30px;
          margin-right: 0;">
    <el-col :span="12">
      <div class="box box1">
        <div class="title">
          IP 信息
          <el-icon size="22" style="position: absolute;margin-left: 8px;">
            <icon-mdi-refresh/>
          </el-icon>
        </div>
        <hr/>
        <ul class="info-list">
          <li><strong>落地ip:</strong> 1.2.3.4</li>
          <li><strong>位置:</strong> 圣何塞</li>
          <li><strong>服务商:</strong> 甲骨文</li>
          <li><strong>组织:</strong> 甲骨文 BMC</li>
          <li><strong>自治域:</strong> AS31898</li>
          <li><strong>时区:</strong> 美国/落山鸡</li>
        </ul>
      </div>
    </el-col>

    <el-col :span="12">
      <div class="box box2">
        <div class="title">
          系统信息
        </div>
        <hr/>
        <ul class="info-list">
          <li><strong>操作系统:</strong> Windows</li>
          <li><strong>运行时间:</strong> 1:06:42</li>
          <li><strong>开机自启:</strong> 已启用</li>
          <li><strong>运行模式:</strong> 服务模式</li>
          <li><strong>代理端口:</strong> 12345</li>
          <li><strong>Pandora-Box 版本:</strong> v1.0.1</li>
        </ul>
      </div>
    </el-col>
  </el-row>
</template>

<style scoped>
.spark {
  max-width: 92%;
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