<script setup>
import {ref} from "vue";
import {EditPen} from "@element-plus/icons-vue";

// 定义数据
const isEditing = ref(false);
const port = ref(3456);

// 切换编辑模式
const toggleEditing = () => {
  if (isEditing.value) {
    // 保存为数字
    const numericPort = parseInt(port.value, 10);
    if (!isNaN(numericPort)) {
      port.value = numericPort;
    }
  }
  isEditing.value = !isEditing.value;
};
</script>

<template>
  <div class="input-container">
    <span>端口:</span>
    <template v-if="isEditing">
      <input
          type="text"
          v-model="port"
          placeholder="请输入端口号"
          autocapitalize="off"
          autocomplete="off"
          autocorrect="off"
          spellcheck="false"
      />
    </template>
    <template v-else>
      <span class="content">{{ port }}</span>
    </template>
    <el-icon
        class="btn"
        @click="toggleEditing"
        v-if="!isEditing">
      <EditPen/>
    </el-icon>
    <el-icon
        class="btn"
        @click="toggleEditing"
        v-if="isEditing">
      <icon-ep-select/>
    </el-icon>
    <el-icon
        class="btn"
        @click="toggleEditing"
        v-if="isEditing">
      <icon-ep-close-bold/>
    </el-icon>
  </div>
</template>

<style scoped>
.input-container {
  display: flex;
  align-items: center;
  gap: 10px;
  height: 30px;
}

span {
  color: var(--text-color);
  font-size: 18px;
  font-weight: bold;
}

.content {
  font-weight: normal;
}

input {
  width: 100px;
  padding: 5px 8px;
  border: 1px solid var(--text-color);
  border-radius: 5px;
  background-color: rgba(255, 255, 255, 0.1);
  color: var(--text-color);
  font-size: 16px;
}

input:focus {
  outline: none;
}

.btn {
  font-size: 18px;
}

.btn:hover {
  cursor: pointer;
  color: var(--hr-color);
}
</style>
