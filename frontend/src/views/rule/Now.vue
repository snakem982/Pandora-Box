<script setup lang="ts">

import MySimpleInput from "@/components/MySimpleInput.vue";

function handleInputChange(value: any) {
  console.log("输入框的值发生了变化：", value);
}

// 原始大数据集合
const allData = Array.from({ length: 20000 }, (_, i) => ({
  name: `User ${i + 1}`,
  age: Math.floor(Math.random() * 50) + 20,
  city: `City ${i % 10}`,
}));

// 分页数据状态
const itemsPerPage = 50; // 每页加载50条数据
const currentPage = ref(1); // 当前页数
const paginatedData = ref(allData.slice(0, itemsPerPage));

// 加载下一页数据
function loadMore() {
  if (currentPage.value * itemsPerPage >= allData.length) return; // 没有更多数据时停止加载
  currentPage.value++;
  const nextPageData = allData.slice(
      (currentPage.value - 1) * itemsPerPage,
      currentPage.value * itemsPerPage
  );
  paginatedData.value = [...paginatedData.value, ...nextPageData];
}

// 监听滚动事件
function handleScroll(event: Event) {
  const target = event.target as HTMLElement;
  if (
      target.scrollTop + target.clientHeight >= target.scrollHeight - 10
  ) {
    loadMore(); // 滚动到底部时加载更多
  }
}


</script>

<template>
  <div class="now">
    <MySimpleInput
        :onInputChange="handleInputChange"
        placeholder="搜索规则"
    ></MySimpleInput>
    <el-row class="title">
      <el-col :span="5">
        类型
      </el-col>
      <el-col :span="14">
        内容
      </el-col>
      <el-col :span="5">
        代理
      </el-col>
    </el-row>
    <div class="info-list" @scroll="handleScroll">
      <el-row
          class="info"
          v-for="(item, i) in paginatedData"
          :key="i"
      >
        <el-col :span="5">{{ item.name }}</el-col>
        <el-col :span="14">{{ item.age }}</el-col>
        <el-col :span="5">{{ item.city }}</el-col>
      </el-row>
    </div>
  </div>
</template>

<style scoped>
.now {
  width: 95%;
  margin-left: 10px;
  margin-top: 5px;
}

.title {
  border-bottom: 2px solid #f4f4f4;
  padding: 18px 5px 8px 5px;
  font-weight: bold;
}

.info-list {
  max-height: calc(100vh - 330px);
  overflow-y: auto;
}

.info {
  border-bottom: 1px solid #ccc;
  padding: 15px 5px 8px 5px;
}

.info:hover {
  background-color: rgba(0, 0, 0, 0.2);
  border-radius: 8px;
}

.info-list::-webkit-scrollbar {
  width: 5px;
  padding-bottom: 20px;
}

.info-list::-webkit-scrollbar-track {
  background: transparent;
}

.info-list::-webkit-scrollbar-thumb {
  background: var(--scrollbar-bg);
  border-radius: 2px;
  transition: background 0.3s ease, box-shadow 0.3s ease;
}

.info-list::-webkit-scrollbar-thumb:hover {
  background: var(--scrollbar-hover-bg);
  box-shadow: var(--scrollbar-hover-shadow);
}

</style>
