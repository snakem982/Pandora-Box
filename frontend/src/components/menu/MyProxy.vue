<script setup lang="ts">
import { useMenuStore } from "@/store/menuStore";
import { useI18n } from "vue-i18n";
import createApi from "@/api";
import { success, warning } from "@/util/load";

// 存储规则模式
const menuStore = useMenuStore();

// 获取当前 Vue 实例的 proxy 对象
const { proxy } = getCurrentInstance()!;
const api = createApi(proxy);

// 国际化
const { t } = useI18n();

// 代理开关
function toggle() {
  // 检测端口是否被占用

  // 检测通过执行后续操作
  menuStore.setProxy(!menuStore.proxy);
  if (menuStore.proxy) {
    success(t("proxy-switch-on"));
  } else {
    warning(t("proxy-switch-off"));
  }
}

// 虚拟网卡开关
function toggle2() {
  // 检测是否运行在管理员模式下


  // 检测通过执行后续操作
  menuStore.setTun(!menuStore.tun);
  if (menuStore.tun) {
    api.updateConfigs({
      tun: {
        enable: true,
      },
    }).then((res: any) => {
      success(t("tun-switch-on"));
    });
  } else {
    api.updateConfigs({
      tun: {
        enable: false,
      },
    }).then((res: any) => {
      warning(t("tun-switch-off"));
    });
  }
}
</script>

<template>
  <div class="sub">
    <div class="switch-container">
      <span class="switch-label">
        {{ $t("proxy-switch") }}
      </span>
      <div
        :class="['switch', { 'switch-on': menuStore.proxy }]"
        @click="toggle"
      >
        <div class="switch-circle"></div>
      </div>
    </div>
    <div class="switch-container">
      <span class="switch-label">
        {{ $t("tun-switch") }}
      </span>
      <div :class="['switch', { 'switch-on': menuStore.tun }]" @click="toggle2">
        <div class="switch-circle"></div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.sub {
  width: 184px; /* 卡片宽度 */
  background-color: rgba(255, 255, 255, 0.1); /* 半透明白色背景 */
  border: 1px solid #ccc;
  border-radius: 8px;
  box-shadow: var(--left-nav-shadow);
  line-height: 1.5;
  margin-left: 22px;
  margin-top: 25px;
  padding-bottom: 12px;
}

.sub:hover {
  box-shadow: var(--left-nav-hover-shadow);
}

.switch-container {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 16px;
  padding-left: 10px;
  margin-top: 10px;
}

.switch-label {
  color: var(--text-color);
}

.switch {
  width: 54px;
  height: 26px;
  border: 2px solid var(--text-color);
  border-radius: 15px;
  position: relative;
  background-color: transparent;
  cursor: pointer;
  transition: background-color 0.3s ease;
}

.switch.switch-on {
  background-color: var(--left-item-selected-bg);
}

.switch-circle {
  width: 20px;
  height: 20px;
  background-color: var(--text-color);
  border-radius: 50%;
  position: absolute;
  top: 3px;
  left: 3px;
  transition: left 0.3s ease;
}

.switch-on .switch-circle {
  left: 31px;
}
</style>
