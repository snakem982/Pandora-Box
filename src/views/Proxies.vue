<script setup lang="ts">
import createApi from "@/api";
import {useProxiesStore} from "@/store/proxiesStore";
import {useMenuStore} from "@/store/menuStore";
import {useSettingStore} from "@/store/settingStore";
import {useI18n} from "vue-i18n";
import {pError, pLoad} from "@/util/pLoad";
import {useWebStore} from "@/store/webStore";
import ProxyOptions from "./proxies/ProxyOptions.vue";
import ProxyGroups from "./proxies/ProxyGroups.vue";
import ProxyNodeList from "./proxies/ProxyNodeList.vue";

const {t} = useI18n();

// 获取当前 Vue 实例的 proxy 对象
const {proxy} = getCurrentInstance()!;
const api = createApi(proxy);

// 当前页面双向绑定对象
const groupList = ref<string[]>([]);
const nodeList = ref<any[]>([]);
const proxyNodeList = ref<InstanceType<typeof ProxyNodeList> | null>(null);

// 当前页面使用store
const proxiesStore = useProxiesStore();
const menuStore = useMenuStore();
const settingStore = useSettingStore();
const webStore = useWebStore();

// 获取分组
async function groups() {
  // 活跃分组
  const active = proxiesStore.active;

  const temp = await api.getGroups();
  switch (menuStore.rule) {
    case "rule":
      groupList.value = temp;
      let find = false
      for (let i = 0; i < temp.length; i++) {
        if (temp[i] === active) {
          find = true
          break
        }
      }
      if (find) {
        if (active == "GLOBAL") {
          proxiesStore.setActive(temp[0]);
        }
      } else {
        proxiesStore.setActive(temp[0]);
      }
      break;
    case "global":
      groupList.value = temp.concat(["GLOBAL"]);
      if (!active) {
        proxiesStore.setActive(temp[0]);
      }
      break;
    case "direct":
      groupList.value = [];
      break;
  }
}

// 获取节点列表
async function nodes() {
  if (menuStore.rule == "direct") {
    nodeList.value = [];
    return;
  }
  nodeList.value = await api.getProxies(
      proxiesStore.active,
      proxiesStore.isHide,
      proxiesStore.isSort
  ); // 更新响应式数据
}

// 设置活跃分组
async function setActive(value: any) {
  if (proxiesStore.active == value) {
    return;
  }
  proxiesStore.setActive(value);
  await nodes();
}

// 设置隐藏
async function setHide() {
  proxiesStore.setHide(!proxiesStore.isHide);
  await nodes();
}

// 设置排序
async function setSort() {
  proxiesStore.setSort(!proxiesStore.isSort);
  await nodes();
}

// 设置代理
async function setProxy(now: any, name: string) {
  if (now) {
    return;
  }
  try {
    await api.setProxy(proxiesStore.active, {name});
    await nodes();
  } catch (error) {
    console.error(error);
  }
}

// 测试延迟
function testDelay() {
  pLoad(t("proxies.loading"), async () => {
    try {
      await api.getDelay(proxiesStore.active, settingStore.testUrl, 3000);
      await nodes();
    } catch (e) {
      if (e['message']) {
        pError(e['message'])
      }
    }
  });
}

function locateSelectedNode() {
  proxyNodeList.value?.locateSelected();
}

let fresh: any = null;
onMounted(async () => {
  await groups();
  await nodes();
  // 创建刷新定时器
  fresh = setInterval(async () => {
    await nodes();
  }, 10000);
});

onBeforeUnmount(() => {
  // 清除定时器
  clearInterval(fresh);
});

// 监听具体状态
watch(() => menuStore.rule, // 监听 store 中的某个状态
    async () => {
      await groups();
      await nodes();
    }
);

watch(() => webStore.fProfile, async () => {
  await groups();
  await nodes();
})

watch(() => proxiesStore.now, async () => {
  await nodes();
})

</script>

<template>
  <MyLayout hr-show>
    <template #top>
      <ProxyOptions
          :is-hide="proxiesStore.isHide"
          :is-sort="proxiesStore.isSort"
          :is-vertical="proxiesStore.isVertical"
          @test="testDelay"
          @locate="locateSelectedNode"
          @toggle-hide="setHide"
          @toggle-sort="setSort"
          @toggle-vertical="proxiesStore.setVertical(!proxiesStore.isVertical)"
      />
      <ProxyGroups
          :groups="groupList"
          :active="proxiesStore.active"
          :is-vertical="proxiesStore.isVertical"
          :rule="menuStore.rule"
          @select="setActive"
      />
    </template>


    <template #bottom>
      <ProxyNodeList
          ref="proxyNodeList"
          :nodes="nodeList"
          :is-direct="menuStore.rule == 'direct'"
          :direct-label="$t('proxies.direct')"
          @select="setProxy"
      />
    </template>
  </MyLayout>
</template>

<style scoped>
.space {
  margin-top: 15px;
}
</style>
