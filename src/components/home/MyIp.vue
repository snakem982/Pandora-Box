<script setup lang="ts">


import {useHomeStore} from "@/store/homeStore";
import {useI18n} from "vue-i18n";
import createApi from "@/api";
import {pError} from "@/util/pLoad";
import {useMenuStore} from "@/store/menuStore";

// 获取当前 Vue 实例的 proxy 对象
const {proxy} = getCurrentInstance()!;
const api = createApi(proxy);

const {t} = useI18n()
const homeStore = useHomeStore()
const menuStore = useMenuStore()

// 预计算常量，减少重复运算
const dayInMs = 1000 * 60 * 60 * 24;
const hourInMs = 1000 * 60 * 60;
const minuteInMs = 1000 * 60;

// 优化计时器更新函数
function updateTimer() {
  const elapsed = Date.now() - homeStore.startTime; // 使用 `Date.now()` 获取当前时间戳

  // 将时间差转换为天、时、分、秒
  const days = Math.floor(elapsed / dayInMs);
  const hours = Math.floor((elapsed % dayInMs) / hourInMs);
  const minutes = Math.floor((elapsed % hourInMs) / minuteInMs);
  const seconds = Math.floor((elapsed % minuteInMs) / 1000);

  let show = `${hours.toString().padStart(2, '0')}:${minutes.toString().padStart(2, '0')}:${seconds.toString().padStart(2, '0')}`

  if (days) {
    show = `${days} ${t('home.system.day')} ` + show
  }

  // 更新计时器显示
  time.value = show;
}

// 获取 ip 信息
async function getIpInfo(hide: boolean = true) {
  ipInfo.value = homeStore.ip;
  try {
    // 切换节点后才进行 ip 请求
    let md6 = await api.getGroupMd5()
    md6 += menuStore.language
    if (homeStore.md6 === md6) {
      return
    } else {
      homeStore.setMd6(md6)
    }

    // 进行ip探测
    const url = "http://ip-api.com/json/?lang=" + t('home.ip.lang');
    const data = await api.getWebTestIp({url});
    data['as'] = data['as'].split(" ")[0];

    // 绑定数据
    ipInfo.value = data;
    homeStore.setIp(data)

  } catch (e) {
    if (hide) {
      // 隐藏错误提示
      return
    }
    // 显示错误提示
    if (e['message']) {
      pError(e['message'])
    }
  }
}

// 页面变量
const time = ref("");
const admin = ref("off");
const version = ref("");
const port = ref("");
const ipInfo = ref({
  query: '',
  regionName: '',
  country: '',
  city: '',
  isp: '',
  timezone: '',
  as: '',
})

onMounted(async () => {
  // 每秒更新
  setInterval(updateTimer, 1000);
  // 获取版本
  version.value = await api.getVersion()
  // 获取端口
  const configs = await api.getConfigs();
  port.value = configs['mixed-port'];
  // 获取ip
  await getIpInfo(true)

  // 检测是否运行在管理员模式下
  const res = await api.getAdmin();
  if (res.data) {
    admin.value = "on"
  } else {
    admin.value = "off"
  }
})


</script>

<template>
  <el-row :gutter="20" class="spark"
          style="margin-left: 2px;">
    <el-col :span="12">
      <div class="box box1">
        <div class="title">
          {{ $t('home.ip.title') }}
          <el-icon size="22"
                   @click="getIpInfo(false)"
                   class="refreshIp">
            <icon-mdi-refresh/>
          </el-icon>
        </div>
        <hr/>
        <ul class="info-list">
          <li><strong>{{ $t('home.ip.real') }} : </strong>
            {{ ipInfo['query'] }}
          </li>
          <li><strong>{{ $t('home.ip.city') }} : </strong>
            {{ ipInfo['city'] }}
          </li>
          <li><strong>{{ $t('home.ip.country') }} : </strong>
            {{ ipInfo['country'] }}
          </li>
          <li><strong>{{ $t('home.ip.isp') }} : </strong>
            {{ ipInfo['isp'] }}
          </li>
          <li><strong>{{ $t('home.ip.asn') }} : </strong>
            {{ ipInfo['as'] }}
          </li>
          <li><strong>{{ $t('home.ip.time-zone') }} : </strong>
            {{ ipInfo['timezone'] }}
          </li>
        </ul>
      </div>
    </el-col>

    <el-col :span="12">
      <div class="box box2">
        <div class="title">
          {{ $t('home.system.title') }}
        </div>
        <hr/>
        <ul class="info-list">
          <li><strong>{{ $t('home.system.os') }} : </strong> {{ homeStore.os }}</li>
          <li><strong>{{ $t('home.system.runtime') }} : </strong>
            {{ time }}
          </li>
          <li><strong>{{ $t('home.system.startup') }} : </strong> {{ $t('off') }}</li>
          <li><strong>{{ $t('home.system.admin') }} : </strong> {{ $t(admin) }}</li>
          <li><strong>{{ $t('home.system.port') }} : </strong>
            {{ port }}
          </li>
          <li><strong>{{ $t('home.system.version') }} : </strong>
            {{ version }}
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
  line-height: 20px;
}

.box1 {
  box-shadow: var(--right-box-shadow);
}

.box2 {
  box-shadow: var(--right-box-shadow);
}

.refreshIp {
  position: absolute;
  margin-left: 8px;
  margin-top: -4px
}

.refreshIp:hover {
  cursor: pointer;
}

</style>