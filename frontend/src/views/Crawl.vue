<script lang="ts" setup>
import {computed, onMounted, reactive, ref} from 'vue'
import {mdiBookInformationVariant, mdiPlusBoxMultiple, mdiRadar, mdiSquareEditOutline, mdiTrashCan} from "@mdi/js";
import SvgIcon from "@jamescoyle/vue-icon";
import {del, get, patch, post, put} from "../api/http";
import {ElLoading, ElMessage} from "element-plus";
import {IsAdmin} from "../../wailsjs/go/main/App";

interface Getter {
  id: string
  type: string
  url: string
}

const drawer = ref(false)
const search = ref('')
const filterTableData = computed(() =>
    tableData.filter(
        (data) =>
            !search.value ||
            data.url.toLowerCase().includes(search.value.toLowerCase())
    )
)

const tableData: Getter[] = reactive([])

const form = reactive({
  id: '',
  type: '',
  url: '',
})
const dialogFormVisible = ref(false)
const addFlag = ref(true)
const formLabelWidth = '150px'

async function getGetter() {
  try {
    const ok = await get<Getter[]>("/getter")
    tableData.length = 0
    tableData.push(...ok)
  } catch (error) {
    console.error(error);
  }
}

async function postGetter(getter: Getter) {
  try {
    await post<any>("/getter", getter)
    await getGetter()
  } catch (error) {
    console.error(error);
  }
}

async function putGetter(getter: Getter) {
  try {
    await put<any>("/getter/" + getter.id, getter)
    await getGetter()
  } catch (error) {
    console.error(error);
  }
}

async function delGetter(getter: Getter) {
  try {
    await del<any>("/getter/" + getter.id)
    await getGetter()
  } catch (error) {
    console.error(error);
  }
}

function addShow() {
  dialogFormVisible.value = true;
  addFlag.value = true;
  form.id = ""
  form.type = ""
  form.url = search.value || ""
}

function editShow(getter: Getter) {
  dialogFormVisible.value = true;
  addFlag.value = false;
  form.id = getter.id
  form.type = getter.type
  form.url = getter.url
}

function isValidHttpUrl(url: string): boolean {
  try {
    const newUrl = new URL(url);
    return newUrl.protocol === 'http:' || newUrl.protocol === 'https:';
  } catch (err) {
    return false;
  }
}

async function addOrEdit() {
  if (form.url == "") {
    ElMessage.error("地址不能为空 Url cannot be empty")
    return
  }
  if (form.type == "") {
    ElMessage.error("类型不能为空 Type cannot be empty")
    return
  }
  const filter = tableData.filter(data => data.type == form.type && data.url == form.url);
  if (filter.length != 0) {
    ElMessage.error("地址已存在 Url already exists")
    return
  }
  if (addFlag.value) {
    if (!isValidHttpUrl(form.url)) {
      ElMessage.error("地址格式不正确 Url format is incorrect")
      return
    }
    await postGetter(form)
  } else {
    await putGetter(form)
  }
  dialogFormVisible.value = false
}

onMounted(getGetter)

const ToDelay = (delay: number) => new Promise((resolve) => setTimeout(resolve, delay))

async function crawling() {
  const loading = ElLoading.service({
    lock: true,
    text: '抓取中Crawling...',
    background: 'rgba(0, 0, 0, 0.7)',
  })

  try {
    let needTun = false
    const isAdmin = await IsAdmin()
    const tun = localStorage.getItem("tun")
    if (isAdmin == "true") {
      if (tun != "off") {
        await patch("/configs", {tun: {enable: false}})
        needTun = true
        await ToDelay(500)
      }
    }

    await get<any>("/crawl",{timeout:1800000})
    ElMessage.success("抓取成功Crawling Success")

    if (needTun) {
      await patch("/configs", {tun: {enable: true, "stack": tun}})
    }
  } catch (error) {
    console.error(error);
    ElMessage.error("抓取失败Crawling Failed")
  }

  loading.close()
}


</script>

<template>
  <div class="header">
    <el-row :gutter="24">
      <el-col :span="18">
        <el-input
            v-model="search"
            placeholder="搜索地址 Search Url"
            size="large"
            autocomplete="off"
        ></el-input>
      </el-col>
      <el-col :span="6">

        <el-tooltip
            content="添加 Add"
            placement="bottom"
        >
          <el-button
              @click="addShow"
              type="success"
              size="large"
              circle>
            <svg-icon type="mdi" :path="mdiPlusBoxMultiple" :size="20"></svg-icon>
          </el-button>
        </el-tooltip>

        <el-tooltip
            content="抓取 Crawl"
            placement="bottom"
        >
          <el-button
              @click="crawling"
              type="warning"
              size="large"
              circle>
            <svg-icon type="mdi" :path="mdiRadar" :size="23"></svg-icon>
          </el-button>
        </el-tooltip>

        <el-tooltip
            content="使用说明 Manual"
            placement="bottom"
        >
          <el-button
              @click.stop="drawer = true"
              type="primary"
              size="large"
              circle>
            <svg-icon type="mdi" :path="mdiBookInformationVariant" :size="23"></svg-icon>
          </el-button>
        </el-tooltip>
      </el-col>
    </el-row>
  </div>
  <el-table :data="filterTableData"
            table-layout="fixed"
            max-height="85vh"
            empty-text="暂无数据 No Data"
            stripe>
    <el-table-column fixed prop="type" label="抓取类型 Type" width="150em"/>
    <el-table-column prop="url" label="抓取地址 Url"/>
    <el-table-column label="操作 Option" width="150em" align="center">
      <template #default="scope">
        <svg-icon type="mdi"
                  @click="editShow(scope.row)"
                  class="s_edit"
                  :path="mdiSquareEditOutline"
                  :size="30"></svg-icon>
        <svg-icon type="mdi"
                  @click="delGetter(scope.row)"
                  class="s_del"
                  :path="mdiTrashCan"
                  :size="30"></svg-icon>
      </template>
    </el-table-column>
  </el-table>

  <el-dialog v-model="dialogFormVisible" :title="addFlag?'添加 Add':'编辑 Edit'">
    <el-form :model="form">
      <el-form-item label="抓取类型 Type" :label-width="formLabelWidth">
        <el-select v-model="form.type" placeholder="选择类型 Select type" style="width: 100%">
          <el-option label="clash订阅(clash subscription)" value="clash"/>
          <el-option label="v2ray订阅(v2ray subscription)" value="v2ray"/>
          <el-option label="分享链接(share link)" value="share"/>
          <el-option label="模糊抓取(fuzzy crawling)" value="fuzzy"/>
        </el-select>
      </el-form-item>
      <el-form-item label="抓取地址 Url" :label-width="formLabelWidth">
        <el-input v-model="form.url" autocomplete="off" clearable/>
      </el-form-item>
    </el-form>
    <template #footer>
      <span class="dialog-footer">
        <el-button @click="dialogFormVisible = false">取消 Cancel</el-button>
        <el-button type="primary" @click="addOrEdit">
          确认 Confirm
        </el-button>
      </span>
    </template>
  </el-dialog>

  <el-drawer
      v-model="drawer"
      title="使用说明 Manual"
      direction="btt"
      size="100vh"
  >
    <div>
      <el-text class="el-text--primary title">1、抓取逻辑 Crawl Logic</el-text>
      <div class="content">
        <el-text>- 节点可用</el-text>
        <br>
        <el-text>- 节点延迟在4s内</el-text>
        <br>
        <el-text>- 节点开启Tls</el-text>
        <br>
        <el-text>- 检测节点server字段是否为Cloudflare Ip,
          如果为Cloudflare Ip则根据本机网络环境自动优选Cloudflare Ip后，新增一个节点到列表中
        </el-text>
        <br>
        <el-text class="el-text--danger">- 抓取成功会更新默认配置</el-text>
        <br><br>

        <el-text>- Node available</el-text>
        <br>
        <el-text>- Node delay is within 4s</el-text>
        <br>
        <el-text>- Node enables Tls</el-text>
        <br>
        <el-text>- Detect whether the node server field is Cloudflare IP. If it is Cloudflare IP, Cloudflare IP will be
          automatically selected according to the local network environment and a new node will be added to the list.
        </el-text>
        <br>
        <el-text class="el-text--danger">- If the crawl is successful, the default profile will be updated.</el-text>
        <br>
        <br>
      </div>
    </div>

    <div>
      <el-text class="el-text--primary title">2、抓取类型 Crawl Type</el-text>
      <div class="content">
        <el-text>- clash订阅(clash subscription)</el-text>
        <br>
        <el-text>&emsp;一般用yaml编码</el-text>
        <br>
        <el-text>&emsp;Generally encoded in yaml</el-text>
        <br>
        <br>
        <el-text>- v2ray订阅(v2ray subscription)</el-text>
        <br>
        <el-text>&emsp;一般用Base64编码</el-text>
        <br>
        <el-text>&emsp;Generally encoded in Base64</el-text>
        <br>
        <br>
        <el-text>- 分享链接(share link)</el-text>
        <br>
        <el-text>&emsp;以下开头的文字视为分享链接</el-text>
        <br>
        <el-text>&emsp;Text starting with the following text is considered a share link.</el-text>
        <br>
        <el-text>&emsp;ss://...</el-text>
        <br>
        <el-text>&emsp;ssr://...</el-text>
        <br>
        <el-text>&emsp;vmess://...</el-text>
        <br>
        <el-text>&emsp;vless://...</el-text>
        <br>
        <el-text>&emsp;trojan://...</el-text>
        <br>
        <el-text>&emsp;tuic://...</el-text>
        <br>
        <el-text>&emsp;hysteria://...</el-text>
        <br>
        <el-text>&emsp;hysteria2://...</el-text>
        <br><br>
        <el-text>- 模糊抓取(fuzzy crawling)</el-text>
        <br>
        <el-text>&emsp;当内容既包含订阅地址又包含分享链接时使用。</el-text>
        <br>
        <el-text>&emsp;不知道填写什么类型时也可以使用此选项。</el-text>
        <br>
        <el-text>&emsp;Used when content contains subscription addresses and sharing links.</el-text>
        <br>
        <el-text>&emsp;You can also use this option when you don't know what type to fill in.</el-text>
        <br>
        <br>
      </div>
    </div>

  </el-drawer>

</template>

<style scoped>
:deep(.cell) {
  font-size: 16px;
}

.s_edit {
  margin-right: 5px;
}

.s_edit:hover {
  cursor: pointer;
  color: #409EFF;
}

.s_del:hover {
  cursor: pointer;
  color: #F56C6C;
}

.header {
  margin-bottom: 10px;
}

:deep(.el-drawer__header) {
  margin-bottom: 0;
}

.title {
  font-size: 17px;
}

.content {
  margin-top: 10px;
  padding-left: 1.5rem;
  font-size: 16px;
  line-height: 1.6rem;
}


</style>
