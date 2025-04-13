import { ElLoading, ElMessage } from "element-plus";

export async function load(tip: any, callback: any) {
  const loading = ElLoading.service({
    lock: true,
    text: tip,
    background: "rgba(0, 0, 0, 0.2)",
  });
  await callback();
  loading.close();
}


export async function copy(textToCopy: any, t: any) {
  try {
    await navigator.clipboard.writeText(textToCopy);
    success(t("copy.success"));
  } catch (error) {
    error(t("copy.fail"));
  }
}

export function success(msg: any) {
  ElMessage({
    message: msg,
    type: "success",
  });
}

export function error(msg: any) {
  ElMessage({
    message: msg,
    type: "error",
  });
}

export function warning(msg: any) {
  ElMessage({
    message: msg,
    type: "warning",
  });
}
