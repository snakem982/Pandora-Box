import { ElLoading , ElMessage } from "element-plus";

export async function load(tip: any, callback: any) {
  const loading = ElLoading.service({
    lock: true,
    text: tip,
    background: "rgba(0, 0, 0, 0.2)",
  });
  await callback();
  loading.close();
}

export async function copy(textToCopy: any,t: any) {
  try {
    await navigator.clipboard.writeText(textToCopy);
    ElMessage({
        message: t('copy.success'),
        type: 'success',
      })
  } catch (error) {
    ElMessage({
        message: t('copy.fail'),
        type: 'error',
      })
  }
}

