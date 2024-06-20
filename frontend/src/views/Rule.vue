<script lang="ts" setup>
import {computed, onMounted, reactive, ref} from "vue";
import {get, post, put} from "../api/http";
import {ElMessage} from "element-plus";
import YamlEditor from "../weight/YamlEditor.vue";

const search = ref('')
const profileRulesTemp = reactive([])
const bypass = ref('')
const bypassBtn = ref(false)

const currentPage = ref(1)
const currentTotal = ref(0)

async function getRules() {
  try {
    const data = await get<any>("/rules")
    const rules: [] = data.rules
    currentTotal.value = rules.length
    profileRulesTemp.push(...rules)

    const ignore: [] = await get("/ignore")
    bypass.value = ignore.join("\n")
  } catch (error) {
    console.error(error);
  }
}

const filterTableData = computed(() => {
      search.value = search.value.toLowerCase()
      const filter = profileRulesTemp.filter(
          (data: any) =>
              !search.value ||
              data.payload.includes(search.value)
      )
      currentTotal.value = filter.length
      return filter.slice((currentPage.value - 1) * 10, currentPage.value * 10)
    }
)

function changePage(page: number) {
  currentPage.value = page
}

async function savaIgnore() {
  let value = bypass.value;
  if (!value || value.trim() === '') {
    return
  }
  value = value.trim()
  let msg: any
  try {
    bypassBtn.value = true
    msg = await put("/ignore", value.split("\n"))
  } catch (error) {
    msg = error
  }
  bypassBtn.value = false
  if (msg && msg.message) {
    ElMessage.error("保存失败 Save failed")
  } else {
    ElMessage.success("保存成功 Saved successfully")
  }
}

const isOn = ref(false)
let sNos = 0;
onMounted(async function () {
  await getRules()
  const data: any = await get("/myRules/on")
  if (data['on'] == "on") {
    isOn.value = true
  } else {
    sNos = 1
  }
})

async function load() {
  let data: any = await get("/myRules")
  if (data['status'] == "ok") {
    return data['buf']
  }

  data = await get("/myRules/default")
  if (data['status'] == "ok") {
    return data['buf']
  }

  return ""
}

async function reset() {
  const data: any = await get("/myRules/default")
  if (data['status'] == "ok") {
    return data['buf']
  }

  return ""
}

async function test(code: any, silent: Boolean) {
  const data: any = await post("/myRules/test", {data: code})
  if (data['status'] == "ok") {
    if (!silent) {
      ElMessage.success("测试成功 Test success")
    }
    return true
  } else {
    ElMessage({
      showClose: true,
      message: "测试失败 Test failed : " + data['status'],
      type: 'error',
      duration: 10000,
    })
    return false
  }
}

async function save(code: any, silent: Boolean) {
  const data: any = await post("/myRules/save", {data: code})
  if (data['status'] == "ok") {
    if (!silent) {
      ElMessage.success("保存成功 Save success")
    }
    return true
  } else {
    ElMessage({
      showClose: true,
      message: "保存失败 Save failed : " + data['status'],
      type: 'error',
      duration: 10000,
    })
    return false
  }
}

async function onSwitch(on: any) {
  if (sNos == 0) {
    sNos++
    return
  }
  const data: any = await put("/myRules/on", {data: on})
  if (data['status'] == "ok") {
    ElMessage.success("切换成功 Switch success")
  } else {
    ElMessage({
      showClose: true,
      message: "切换失败 Switch failed : " + data['status'],
      type: 'error',
      duration: 10000,
    })
  }
}

</script>

<template>
  <el-tabs stretch>
    <el-tab-pane label="查看规则 View Rules">
      <div class="header">
        <el-row :gutter="24">
          <el-col :span="24">
            <el-input
                v-model="search"
                placeholder="搜索内容 Search Payload"
                size="large"
                autocomplete="off"
            ></el-input>
          </el-col>
        </el-row>
      </div>
      <el-table
          :data="filterTableData"
          table-layout="fixed"
          stripe
          style="width: 100%"
          max-height="65vh"
          empty-text="暂无数据 No Data">
        <el-table-column prop="type" label="类型 Type" width="225em"/>
        <el-table-column prop="payload" label="内容 Payload"/>
        <el-table-column prop="proxy" label="代理 Proxy" width="250em"/>
      </el-table>
      <el-pagination
          :current-page="currentPage"
          page-size="10"
          layout="prev, pager, next, jumper"
          :total="currentTotal"
          @current-change="changePage"
          style="margin-top: 10px;"
      />
    </el-tab-pane>
    <el-tab-pane label="统一规则分组 Unified rule grouping" lazy>
      <yaml-editor
          :load="load"
          :reset="reset"
          :test="test"
          :save="save"
          :onSwitch="onSwitch"
          :is-on="isOn"
      ></yaml-editor>
    </el-tab-pane>
    <el-tab-pane label="绕过 Bypass">
      <el-text>按行分割 Split by row</el-text>
      <el-input
          v-model="bypass"
          type="textarea"
          autocomplete="off"
      ></el-input>
      <el-button
          type="primary"
          style="margin-top: 15px"
          :disabled="bypassBtn"
          @click="savaIgnore"
      >
        保存 Save
      </el-button>
    </el-tab-pane>
  </el-tabs>
</template>


<style scoped>
.header {
  margin-bottom: 10px;
}

:deep(.cell) {
  font-size: 16px;
}

:deep(.el-textarea__inner) {
  height: 50vh;
  margin-top: 10px;
}

</style>
