<script setup lang="ts">

</script>

<template>
  <div class="ignore">
    <el-space class="op">
      <el-button>保存</el-button>
      <el-divider direction="vertical" border-style="dashed"/>
      <el-text class="st">按行分割</el-text>
    </el-space>
    <div class="content">
      <textarea
          class="custom-textarea"
          placeholder="在这里输入"
      ></textarea>
    </div>
  </div>
</template>

<style scoped>
.ignore {
  width: 95%;
  margin-left: 10px;
  margin-top: 5px;
}

.op {
  margin-top: 12px;
}

:deep(.el-button) {
  padding: 2px 10px;
  --el-button-bg-color: transparent;
  --el-button-text-color: var(--text-color);
  --el-button-hover-text-color: var(--left-item-selected-bg);
  --el-button-hover-bg-color: var(--text-color)
}

.st {
  color: var(--text-color);
}

.content {
  margin-top: 25px;
}

.custom-textarea {
  background-color: transparent; /* 背景透明 */
  border: 2px solid var(--text-color);       /* 边界为 2px 的白色 */
  color: white;                  /* 文字颜色为白色 */
  padding: 8px;                  /* 内间距，确保内容不贴边 */
  border-radius: 8px;            /* 圆角样式（可选） */
  font-size: 16px;               /* 字体大小 */
  resize: none;                  /* 禁止调整大小（可选） */
  outline: none;                 /* 去掉点击时的默认高亮框 */
  width: 96%;
  height: calc(100vh - 340px);
}

.custom-textarea::placeholder {
  color: rgba(255, 255, 255, 0.6); /* 占位符文字颜色，设置为白色半透明 */
}

.custom-textarea:focus {
  box-shadow: var(--right-box-shadow); /* 焦点时添加阴影效果（可选） */
}
</style>