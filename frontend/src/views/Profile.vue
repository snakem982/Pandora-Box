<script setup lang="ts">
import {onMounted, reactive, ref} from 'vue'
import {ElLoading, ElMessage, UploadFile, UploadFiles} from "element-plus";
import {del, get, patch, post, put} from "../api/http";
import {Profile} from "../api/pandora";
import {ClipboardGetText, ClipboardSetText} from "../../wailsjs/runtime";
import {mdiDownload, mdiFileReplace, mdiFolderOpen} from "@mdi/js";
import SvgIcon from "@jamescoyle/vue-icon";
import {useRouter} from "vue-router";
import {GetFreePort, GetSecret} from "../../wailsjs/go/main/App";

const subOrShare = ref('')
const drawer = ref(false)
const zh = navigator.language.toLowerCase().indexOf("zh") > -1

function uploadError(error: Error, uploadFile: UploadFile, uploadFiles: UploadFiles) {
  const err: any = JSON.parse(error.message);
  ElMessage({
    showClose: true,
    dangerouslyUseHTMLString: true,
    message: uploadFile.name + '<br/><br/>导入失败错误信息<br/>Import failed error message:<br/><br/>' + err.message,
    type: 'error',
    duration: 10000,
  })
}

async function uploadSuccess(error: Error, uploadFile: UploadFile, uploadFiles: UploadFiles) {
  await getProfile()
  ElMessage({
    showClose: true,
    message: uploadFile.name + ' 导入成功! Import success!',
    type: 'success',
  })
}

const drawerUpload: any = ref(null);

function beforeDrawClose(done: (cancel?: boolean) => void) {
  drawerUpload.value.clearFiles()
  done(false)
}

async function downSubOrShare() {
  const trim = subOrShare.value.trim();
  if (trim == "") {
    ElMessage.error("订阅地址或分享链接不能为空、Subscription address or sharing link cannot be empty")
    return
  }

  const loading = ElLoading.service({
    lock: true,
    text: '解析中Parsing...',
    background: 'rgba(0, 0, 0, 0.7)',
  })
  try {
    await post("/profile", {data: subOrShare.value})
  } catch (error) {
    console.error(error);
    loading.close()
    ElMessage.error("解析失败Parsing failed")
    return
  }
  const length = data.length
  await getProfile()
  loading.close()
  if (data.length > length) {
    ElMessage.success("解析成功Parsed successfully")
    subOrShare.value = ""
  } else {
    ElMessage.error("解析失败Parsing failed")
  }
}


const data: Profile[] = reactive([])

async function getProfile() {
  try {
    const ok = await get<Profile[]>("/profile")
    data.length = 0
    data.push(...ok)
  } catch (error) {
    console.error(error);
  }
}

async function delProfile(id: string) {
  try {
    await del<Profile[]>("/profile/" + id)
    await getProfile()
  } catch (error) {
    console.error(error);
  }
}

const uploadUrl = ref('')
const uploadHeader = reactive({
  Authorization: ''
})

onMounted(async () => {
  await getProfile()
  const host = await GetFreePort()
  uploadUrl.value = "http://" + host + "/profile/file"
  const secret = await GetSecret()
  uploadHeader.Authorization = 'Bearer ' + secret
})

const toolFormVisible = ref(false)
const formLabelWidth = '140px'
const form: Profile = reactive({
  id: '',
  type: 0,
  title: '',
  path: '',
  url: '',
  order: -1,
  selected: false
})

function toolSet(profile: Profile) {
  form.id = profile.id
  form.type = profile.type
  form.title = profile.title
  form.path = profile.path
  form.url = profile.url
  form.order = profile.order
  form.selected = profile.selected
  toolFormVisible.value = true
}

async function setProfile() {
  try {
    await put<any>("/profile/" + form.id, form)
    await getProfile()
  } catch (error) {
    console.error(error);
  }
  toolFormVisible.value = false
}

function cardSelectClass(isSelect: boolean): string {
  if (isSelect) {
    return "box-card box-card-select"
  } else {
    return "box-card"
  }
}

async function selectProfile(profile: Profile) {
  if (profile.selected) {
    return
  }
  for (let datum of data) {
    if (datum.selected) {
      datum.selected = false
      await put("/profile/" + datum.id, datum)
      break
    }
  }
  profile.selected = true
  await put("/profile/" + profile.id, profile)
  await patch<any>("/profile/" + profile.id, profile)

  setTimeout(() => del("/connections"), 3000)
}

async function pastTxt() {
  const txt = await ClipboardGetText()
  if (txt) {
    subOrShare.value = txt
  }
}

interface Msg {
  message: string
}

async function refreshProfile(profile: Profile) {
  const loading = ElLoading.service({
    lock: true,
    text: '更新中Refreshing...',
    background: 'rgba(0, 0, 0, 0.7)',
  })
  let msg: any
  try {
    msg = await put<Msg>("/profile/refresh", profile);
  } catch (error) {
    msg = error
  }
  loading.close()
  if (msg && msg.message) {
    ElMessage({
      showClose: true,
      dangerouslyUseHTMLString: true,
      message: '更新失败<br/>Refresh failed' + '<br/><br/>错误信息<br/>Error message:<br/><br/>' + msg.message,
      type: 'error',
      duration: 10000,
    })
  } else {
    if (profile.selected) {
      await patch<any>("/profile/" + profile.id, profile)
    }
    ElMessage.success("更新成功Refresh successfully")
  }
}

const router = useRouter()

function toCrawl() {
  router.push("/crawl")
}


</script>

<template>
  <div class="header">
    <el-row :gutter="24">
      <el-col :span="18">
        <el-input
            v-model="subOrShare"
            placeholder="订阅地址、分享链接、Base64、Yaml、Subscription、ShareLink"
            size="large"
            type="textarea"
            autocomplete="off"
        ></el-input>
      </el-col>
      <el-col :span="6" style="padding-top: 5px">
        <el-tooltip
            content="下载 Download"
            placement="bottom"
        >
          <el-button
              @click.stop="downSubOrShare"
              style="color: #817df7;background-color: #eaeafd"
              size="large"
              circle>
            <svg-icon type="mdi" :path="mdiDownload" :size="22"></svg-icon>
          </el-button>
        </el-tooltip>

        <el-tooltip
            content="粘贴 Past"
            placement="bottom"
        >
          <el-button
              @click.stop="pastTxt"
              type="primary"
              size="large"
              circle>
            <svg-icon type="mdi" :path="mdiFileReplace" :size="20"></svg-icon>
          </el-button>
        </el-tooltip>

        <el-tooltip
            content="打开文件 Open File"
            placement="bottom"
        >
          <el-button
              @click.stop="drawer = true"
              type="danger"
              size="large"
              circle>
            <svg-icon type="mdi" :path="mdiFolderOpen" :size="23"></svg-icon>
          </el-button>
        </el-tooltip>
      </el-col>
    </el-row>
  </div>
  <div class="content">

    <el-space wrap :size="20">

      <template v-for="(profile,index) in data"
                :key="index">
        <el-card
            :class="cardSelectClass(profile.selected)"
            @click="selectProfile(profile)"
            v-if="profile.type==1">
          <el-space alignment="flex-start">
            <el-button color="#eaeafd" circle style="width: 50px;height: 50px;font-size: 28px">
              <el-icon color="#817df7">
                <Platform/>
              </el-icon>
            </el-button>
            <el-space direction="vertical" alignment="flex-start">
              <el-text type="primary" truncated style="font-size: 16px;margin-left: 5px">
                默认配置
              </el-text>
              <el-text truncated style="font-size: 16px;margin-left: 5px">
                Default Config
              </el-text>
            </el-space>
          </el-space>
          <el-divider/>
          <el-icon class="ca-tool" @click.stop="toCrawl">
            <MagicStick/>
          </el-icon>
        </el-card>

        <el-card
            :class="cardSelectClass(profile.selected)"
            @click="selectProfile(profile)"
            v-if="profile.type==2">
          <el-space alignment="flex-start">
            <el-button color="#eaeafd" circle style="width: 50px;height: 50px;font-size: 28px">
              <el-icon color="#817df7">
                <Share/>
              </el-icon>
            </el-button>
            <el-space direction="vertical" alignment="flex-start">
              <el-text type="primary" truncated style="font-size: 16px;margin-left: 5px">
                {{ profile.title }}
              </el-text>
              <el-text truncated style="font-size: 16px;margin-left: 5px">
                Share Link
              </el-text>
            </el-space>
          </el-space>
          <el-divider/>
          <el-icon
              @click.stop="toolSet(profile)"
              class="ca-tool">
            <Tools/>
          </el-icon>
          <el-icon
              @click.stop="delProfile(profile.id)"
              class="ca-del">
            <DeleteFilled/>
          </el-icon>
        </el-card>

        <el-card
            :class="cardSelectClass(profile.selected)"
            @click="selectProfile(profile)"
            v-if="profile.type==32 || profile.type==31">
          <el-space alignment="flex-start">
            <el-button color="#eaeafd" circle style="width: 50px;height: 50px;font-size: 28px">
              <el-icon color="#817df7">
                <Link/>
              </el-icon>
            </el-button>
            <el-space direction="vertical" alignment="flex-start">
              <el-text type="primary" truncated style="width: 150px; font-size: 16px;margin-left: 5px">
                {{ profile.title }}
              </el-text>
              <el-text truncated style="width: 150px;font-size: 16px;margin-left: 5px">
                {{ profile.url }}
              </el-text>
            </el-space>
            <el-icon
                @click.stop="refreshProfile(profile)"
                class="ca-refresh">
              <RefreshRight/>
            </el-icon>
          </el-space>
          <el-divider/>
          <el-icon
              @click.stop="toolSet(profile)"
              class="ca-tool">
            <Tools/>
          </el-icon>
          <el-icon
              @click.stop="ClipboardSetText(profile.url)"
              class="ca-copy">
            <CopyDocument/>
          </el-icon>
          <el-icon
              @click.stop="delProfile(profile.id)"
              class="ca-del">
            <DeleteFilled/>
          </el-icon>
        </el-card>

        <el-card
            :class="cardSelectClass(profile.selected)"
            @click="selectProfile(profile)"
            v-if="profile.type==42 || profile.type==41">
          <el-space alignment="flex-start">
            <el-button color="#eaeafd" circle style="width: 50px;height: 50px;font-size: 28px">
              <el-icon color="#817df7">
                <HomeFilled/>
              </el-icon>
            </el-button>
            <el-space direction="vertical" alignment="flex-start">
              <el-text type="primary" truncated style="width: 150px; font-size: 16px;margin-left: 5px">
                {{ profile.title }}
              </el-text>
              <el-text truncated style="width: 150px;font-size: 16px;margin-left: 5px">
                Local File
              </el-text>
            </el-space>
          </el-space>
          <el-divider/>
          <el-icon
              @click.stop="toolSet(profile)"
              class="ca-tool">
            <Tools/>
          </el-icon>
          <el-icon
              @click.stop="delProfile(profile.id)"
              class="ca-del">
            <DeleteFilled/>
          </el-icon>
        </el-card>


      </template>

    </el-space>
  </div>

  <el-drawer
      v-model="drawer"
      title="打开文件 Open File"
      direction="btt"
      size="100vh"
      :before-close="beforeDrawClose"
  >
    <el-upload
        drag
        :action="uploadUrl"
        multiple
        :on-error="uploadError"
        :on-success="uploadSuccess"
        ref="drawerUpload"
        :headers = uploadHeader
    >
      <el-icon size="50" class="el-icon--upload">
        <UploadFilled/>
      </el-icon>
      <div class="el-upload__text">
        将文件拖放到此处&nbsp;或&nbsp;<em>单击导入 </em><br>
        Drop file here or <em>click to import</em>
      </div>
      <template #tip>
        <div class="el-upload__tip">
          文件小于2MB &emsp;File less than 2MB
        </div>
      </template>
    </el-upload>
  </el-drawer>

  <el-dialog v-model="toolFormVisible" title="设置 Config">
    <el-form :model="form">
      <el-form-item label="标题 Title" :label-width="formLabelWidth">
        <el-input v-model="form.title" clearable autocomplete="off"/>
      </el-form-item>
      <el-form-item v-if="form.type==31 || form.type==32" label="订阅地址 Url" :label-width="formLabelWidth">
        <el-input v-model="form.url" clearable autocomplete="off"/>
      </el-form-item>
    </el-form>
    <template #footer>
      <span class="dialog-footer">
        <el-button @click.stop="toolFormVisible = false">取消 Cancel</el-button>
        <el-button type="primary" @click.stop="setProfile">
          确认 Confirm
        </el-button>
      </span>
    </template>
  </el-dialog>

</template>

<style scoped>

.box-card {
  padding: 0;
  width: 275px;
}

.box-card:hover {
  cursor: pointer;
}

.box-card-select {
  border-left-style: solid;
  border-left-color: #817df7;
  border-left-width: 5px;
  width: 270px;
}

.header {
  margin-top: 20px;
  margin-left: 10px;
}

.content {
  margin-top: 20px;
  margin-left: 10px;
}

:deep(.el-collapse-item__header) {
  font-size: 18px;
}

:deep(.el-card__body) {
  padding: 10px;
}

:deep(.el-divider--horizontal) {
  margin: 5px 0;
}

.ca-refresh {
  margin-left: 5px
}

.ca-refresh:hover {
  color: #67C23A;
}

.ca-tool {
  margin-top: 5px
}

.ca-tool:hover {
  color: #409EFF;
}

.ca-copy {
  margin-top: 5px;
  margin-left: 10px
}

.ca-copy:hover {
  color: #E6A23C;
}

.ca-del {
  margin-top: 5px;
  margin-right: 10px;
  float: right;
}

.ca-del:hover {
  color: #F56C6C;
}

</style>