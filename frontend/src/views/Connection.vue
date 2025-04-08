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
        <div class="title">
          {{ $t('connections.title') }}
        </div>
      </el-space>
      <MyHr :update="upFromTop" v-show="false"></MyHr>
    </template>
    <template #bottom>
      <div class="conn">
        <el-space class="op">
          <div class="search">
            <MySimpleInput
                :onInputChange="handleInputChange"
                :placeholder="$t('connections.search')"
                class="search"
            ></MySimpleInput>
          </div>
          <el-button>
            {{ $t('connections.close') }}
          </el-button>
        </el-space>
      </div>

      <div class="content">
        <div class="info-list" @scroll="handleScroll">
          <el-row
              class="info"
              v-for="(item, i) in paginatedData"
              :key="i"
          >
            <el-col :span="24">
              <el-tag type="success" size="small">HTTPS</el-tag>
              &emsp;
              <el-tag type="primary" size="small">Google Chrome Helper</el-tag>
              &emsp;
              <el-tag type="danger" size="small">less than a minute</el-tag>
              <div class="od">
                <span class="ot">{{ $t('connections.host') }} : </span>otheve.beacon.qq.com:443
              </div>
              <div class="od">
                <span class="ot">{{ $t('connections.download') }} : </span>118 KB
                &emsp;
                &#8595;
                20 MB/s
                &emsp;
                <span class="ot">{{ $t('connections.upload') }} : </span>26.7 KB
                &emsp;
                &#8593;
                120 KB/s
              </div>
              <div class="od">
                <span class="ot">{{ $t('connections.rule') }} : </span>DomainKeyword &#8594; google
              </div>
              <div class="od">
                <span class="ot">{{ $t('connections.chains') }} : </span>🐟 漏网之鱼 / 🚀 节点选择 / 🇯🇵 日本IEPL 专线 02
              </div>
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

:deep(.el-button) {
  padding: 2px 10px;
  --el-button-bg-color: transparent;
  --el-button-text-color: var(--text-color);
  --el-button-hover-text-color: var(--left-item-selected-bg);
  --el-button-hover-bg-color: var(--text-color)
}

.content {
  border: 2px solid var(--text-color);
  margin-top: 20px;
  width: calc(95% - 10px);
  margin-left: 10px;
  border-radius: 10px;
}

.info-list {
  max-height: calc(100vh - 250px);
  overflow-y: auto;
  border-radius: 10px;
}

.info {
  border-bottom: 1px solid #ccc;
  padding: 5px 10px;
  font-size: 15px;
  line-height: 1.6;
  background-color: rgba(0, 0, 0, 0.1);
  border-radius: 10px;
}

.od {
  -webkit-user-select: text;
  user-select: text;
}

.ot {
  font-weight: bold;
  font-size: 15px;
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