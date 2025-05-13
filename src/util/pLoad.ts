import {ElLoading, ElMessage} from "element-plus";

export async function pLoad(tip: any, callback: any) {
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
        pSuccess(t("copy.success"));
    } catch (error) {
        error(t("copy.fail"));
    }
}

export function pSuccess(msg: any) {
    ElMessage({
        message: msg,
        type: "success",
        grouping: true
    });
}

export function pError(msg: any) {
    ElMessage({
        message: msg,
        type: "error",
        duration: 5000,
        grouping: true
    });
}

export function pWarning(msg: any) {
    ElMessage({
        message: msg,
        type: "warning",
        duration: 5000,
        grouping: true
    });
}
