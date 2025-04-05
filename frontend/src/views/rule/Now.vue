<script setup lang="ts">

const searchValue = ref('')

const tableData = [
  {
    date: '2016-05-03',
    name: 'Tom',
    address: 'No. 189, Grove St, Los Angeles',
  },
  {
    date: '2016-05-02',
    name: 'Tom',
    address: 'No. 189, Grove St, Los Angeles',
  },
  {
    date: '2016-05-04',
    name: 'Tom',
    address: 'No. 189, Grove St, Los Angeles',
  },
  {
    date: '2016-05-01',
    name: 'Tom',
    address: 'No. 189, Grove St, Los Angeles',
  },
]

</script>

<template>
<div class="now">
    <el-input 
    v-model="searchValue"
    clearable
    placeholder="搜索内容" 
    autocomplete="off"/>
    <el-table :data="tableData" class="table">
      <el-table-column prop="date" label="Date" width="180" />
      <el-table-column prop="name" label="Name" width="180" />
      <el-table-column prop="address" label="Address" />
    </el-table>
</div>
</template>

<style scoped>
.now {
    width: 95%;
    margin-left: 10px;
    margin-top: 5px;
}

:deep(.el-input) {
  --el-input-border-radius: 5px;
  background-color: rgba(255, 255, 255, 0.2);
  --el-input-text-color: var(--el-text-color-regular);
  --el-input-border: var(--el-border);
  --el-input-hover-border: var(--el-border-color-hover);
  --el-input-focus-border: var(--el-color-primary);
  --el-input-transparent-border: 0 0 0 1px transparent inset;
  --el-input-border-color: var(--el-border-color);
  --el-input-bg-color: transparent;
  --el-input-icon-color: var(--el-text-color-placeholder);
  --el-input-placeholder-color: var(--el-text-color-placeholder);
  --el-input-hover-border-color: var(--el-border-color-hover);
  --el-input-clear-hover-color: var(--el-text-color-secondary);
  --el-input-focus-border-color: var(--el-color-primary);
  --el-input-width: 100%;
  --el-input-height: var(--el-component-size);
  box-sizing: border-box;
  display: inline-flex;
  font-size: var(--el-font-size-base);
  line-height: var(--el-input-height);
  position: relative;
  vertical-align: middle;
  width: var(--el-input-width);
}

:deep(.el-input__inner) {
  color: var(--text-color);
}



.table {
    width: 100%;
    margin-top: 15px;
    border-radius: 8px;
  }

</style>
