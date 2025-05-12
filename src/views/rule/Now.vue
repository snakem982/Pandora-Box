<script setup lang="ts">

import MySimpleInput from "@/components/MySimpleInput.vue";
import createApi from "@/api";
import {useWebStore} from "@/store/webStore";


// 获取当前 Vue 实例的 proxy 对象 和 api
const {proxy} = getCurrentInstance()!;
const api = createApi(proxy);


// 原始大数据集合
const allData = ref([]);
const filterData = ref([]);

// 分页数据状态
const itemsPerPage = 50; // 每页加载50条数据
const currentPage = ref(1); // 当前页数
const paginatedData = ref([]);

// 加载下一页数据
function loadMore() {
  if (currentPage.value * itemsPerPage >= filterData.value.length) return; // 没有更多数据时停止加载
  currentPage.value++;
  const nextPageData = filterData.value.slice(
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


onMounted(() => {
  api.getRules().then((res) => {
    allData.value = res;
    filterData.value = res;
    // 初始化分页数据
    paginatedData.value = allData.value.slice(0, itemsPerPage);
  });
});

// 过滤数据
function handleInputChange(value: any) {
  if (value) {
    filterData.value = allData.value.filter((item: any) => {
      return (
          item.type.toLowerCase().includes(value.toLowerCase()) ||
          item.payload.toLowerCase().includes(value.toLowerCase()) ||
          item.proxy.toLowerCase().includes(value.toLowerCase())
      );
    });
  } else {
    filterData.value = allData.value;
  }
  // 重置分页数据
  currentPage.value = 1;
  paginatedData.value = filterData.value.slice(0, itemsPerPage);
}

// 监控配置切换
const webStore = useWebStore();
watch(() => webStore.fProfile, async () => {
  await api.waitRunning()
  api.getRules().then((res) => {
    allData.value = res;
    filterData.value = res;
    // 初始化分页数据
    paginatedData.value = allData.value.slice(0, itemsPerPage);
  });
})

</script>

<template>
  <div class="now">
    <MySimpleInput
        :onInputChange="handleInputChange"
        :placeholder="$t('rule.now.search')"
        class="search"
    ></MySimpleInput>

    <div class="content">
      <el-row class="title">
        <el-col :span="5">
          {{ $t('rule.now.type') }}
        </el-col>
        <el-col :span="14">
          {{ $t('rule.now.payload') }}
        </el-col>
        <el-col :span="5">
          {{ $t('rule.now.proxy') }}
        </el-col>
      </el-row>
      <div class="info-list" @scroll="handleScroll">
        <el-row
            :class="i%2 == 1? 'info info-s' : 'info'"
            v-for="(item, i) in paginatedData"
            :key="i"
        >
          <el-col :span="5">{{ item.type }}</el-col>
          <el-col :span="14">{{ item.payload }}</el-col>
          <el-col :span="5">{{ item.proxy }}</el-col>
        </el-row>
      </div>
    </div>

  </div>
</template>

<style scoped>
.now {
  width: 95%;
  margin-left: 10px;
  margin-top: 5px;
}

.search {
  margin-top: 12px;
}

.content {
  border: 2px solid var(--text-color);
  border-radius: 10px;
  margin-top: 25px;
}

.title {
  border-bottom: 2px solid #f4f4f4;
  padding: 8px 10px;
  font-weight: bold;
}

.info-list {
  max-height: calc(100vh - 365px);
  overflow-y: auto;
}

.info {
  border-bottom: 1px solid #ccc;
  padding: 8px 10px;
}

.info-s {
  border-bottom: 1px solid #ccc;
  padding: 8px 10px;
  background-color: rgba(128, 128, 128, 0.2); /* 深灰色，透明度为50% */
}

.info:hover {
  background-color: rgba(0, 0, 0, 0.2);
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
