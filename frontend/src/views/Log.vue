<script setup lang="ts">
import MyHr from "@/components/proxies/MyHr.vue";
import MySimpleInput from "@/components/MySimpleInput.vue";
import { useWebStore } from "@/store/webStore";

// 调整顶部高度
const distanceFromTop = ref(195);
const upFromTop = function (distance: number) {
  distanceFromTop.value = distance;
};

// 获取Store
const webStore = useWebStore();

// 搜索框
const search = ref("");
function handleInputChange(value: any) {
  search.value = value;
}

// 过滤数据
function filterData() {
  return webStore.logs.filter((data: any) => {
    const searchLower = search.value.toLowerCase();
    return (
      !search.value ||
      data.payload.toLowerCase().includes(searchLower) || // 内容过滤
      data.type.toLowerCase().includes(searchLower) // 类型过滤
    );
  });
}
</script>

<template>
  <MyLayout
    :top-height="distanceFromTop - 15"
    :bottom-height="distanceFromTop + 25"
  >
    <template #top>
      <MySearch></MySearch>
      <el-space class="space">
        <div class="title">
          {{ $t("log.title") }}
        </div>
      </el-space>
      <MyHr :update="upFromTop" v-show="false"></MyHr>
    </template>
    <template #bottom>
      <div class="conn">
        <div class="search">
          <MySimpleInput
            :onInputChange="handleInputChange"
            :placeholder="$t('log.search')"
            class="search"
          ></MySimpleInput>
        </div>
      </div>

      <div class="content">
        <div class="info-list">
          <el-row class="info" v-for="(item, i) in filterData()" :key="i">
            <el-col :span="24">
              <div>
                {{ item.time }}&emsp;[{{ item.type }}]
                <br>
                {{ item.payload }}
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
