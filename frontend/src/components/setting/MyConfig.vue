<script setup lang="ts">

import MyInput from "@/components/setting/MyInput.vue";
import MyTun from "@/components/setting/MyTun.vue";
import {EditPen} from "@element-plus/icons-vue";

const value1 = ref(true)
const value2 = ref(false)
const value3 = ref(true)

</script>

<template>
  <el-row :gutter="20" class="spark"
          style="margin-left: 0;
          margin-right: 0;">
    <el-col :span="12">
      <div class="box box1">
        <div class="title">
          Mihomo
        </div>
        <hr/>
        <ul class="info-list">
          <li>
            <MyInput></MyInput>
          </li>
          <li>
            <MyTun></MyTun>
          </li>
          <li>
            <strong>
              DNS覆写:
            </strong>
            <el-icon class="btn">
              <EditPen/>
            </el-icon>
            <el-switch
                v-model="value1"
                class="set-switch"
                style="margin-left: 28px"
            />
          </li>
          <li>
            <strong>局域网:</strong>
            <el-switch
                v-model="value2"
                class="set-switch"
            />
          </li>
          <li>
            <strong>Ipv6:</strong>
            <el-switch
                v-model="value3"
                class="set-switch"
            />
          </li>
          <li style="height: 30px">
            <strong>Api:</strong>
            http://127.0.0.1:7981
            <el-button>复制</el-button>
          </li>
          <li style="height: 30px">
            <strong>Secret:</strong>
            syeSsdeddu
            <el-button>复制</el-button>
          </li>
        </ul>
      </div>
    </el-col>

    <el-col :span="12">
      <div class="box box2">
        <div class="title">
          Pandora-Box
        </div>
        <hr/>
        <ul class="info-list">
          <li>
            <strong>开机启动:</strong>
            <el-switch
                v-model="value2"
                class="set-switch"
            />
          </li>
          <li style="height: 30px">
            <strong>配置目录:</strong>
            <el-button style="margin-left: 10px">打开</el-button>
            <el-button>导出</el-button>
            <el-button>导入</el-button>
          </li>
          <li style="height: 30px">
            <strong>软件更新:</strong>
            <el-button style="margin-left: 10px">检查</el-button>
          </li>
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

.set-switch {
  margin-left: 10px;
  --el-switch-border-color: var(--text-color);
  --el-switch-on-color: var(--left-item-selected-bg);
  --el-switch-off-color: transparent;
}

:deep(.el-switch__core) {
  width: 46px;
  height: 26px;
  border-radius: 12px;
  border: 2px solid var(--text-color);
}

:deep(.el-switch__core .el-switch__action) {
  margin-left: 2px;
}

:deep(.el-switch.is-checked .el-switch__core .el-switch__action) {
  left: calc(100% - 21px);
}

:deep(.el-button) {
  padding: 2px 10px;
  --el-button-bg-color: transparent;
  --el-button-text-color: var(--text-color);
  --el-button-hover-text-color: var(--left-item-selected-bg);
  --el-button-hover-bg-color: var(--text-color)
}

.btn {
  font-size: 18px;
  position: absolute;
  margin-top: 6px;
}

.btn:hover {
  cursor: pointer;
  color: var(--hr-color);
}


</style>