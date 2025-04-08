<script setup lang="ts">
import MyHr from "@/components/proxies/MyHr.vue";
import MySimpleInput from "@/components/MySimpleInput.vue";

const distanceFromTop = ref(195)
const upFromTop = function (distance: number) {
  distanceFromTop.value = distance
}

function handleInputChange(value: any) {
  console.log("输入框的值发生了变化：", value);
}

// 原始大数据集合
const allData = Array.from({length: 1000}, (_, i) => ({
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
  <MyLayout :top-height="distanceFromTop-15" :bottom-height="distanceFromTop+25">
    <template #top>
      <MySearch></MySearch>
      <el-space class="space">
        <div class="title">日志</div>
      </el-space>
      <MyHr :update="upFromTop" v-show="false"></MyHr>
    </template>
    <template #bottom>
      <div class="conn">
          <div class="search">
            <MySimpleInput
                :onInputChange="handleInputChange"
                placeholder="搜索内容"
                class="search"
            ></MySimpleInput>
          </div>
      </div>

      <div class="content">
        <div class="info-list" @scroll="handleScroll">
          <el-row
              class="info"
              v-for="(item, i) in paginatedData"
              :key="i"
          >
            <el-col :span="24">
              <div>2025/04/08 16:57:25 INFO</div>
              [TCP] 127.0.0.1:50611(夸克网盘 Helper) --> track.lc.quark.cn:80 match DomainSuffix(cn) using 🎯 全球直连[DIRECT]
            </el-col>
          </el-row>
        </div>
      </div>

    </template>
  </MyLayout>
</template>

<style scoped>
.space {
  margin-top: 20px;
}

.conn {
  width: 95%;
  margin-left: 10px;
  margin-top: 5px;
}

.title {
  font-size: 32px;
  font-weight: bold;
  margin-left: 10px;
}

.search {
  width: 400px;
}

.content {
  border: 2px solid var(--text-color);
  border-radius: 10px;
  margin-top: 20px;
  width: 95%;
  margin-left: 10px;
}

.info-list {
  max-height: calc(100vh - 250px);
  overflow-y: auto;
}

.info {
  border-bottom: 1px solid #ccc;
  padding: 5px 10px;
  font-size: 14px;
  line-height: 1.5;
  -webkit-user-select: text;
  user-select: text;
  background-color: rgba(0, 0, 0, 0.1);
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