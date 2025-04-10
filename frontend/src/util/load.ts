import {ElLoading} from "element-plus";

export async function load(tip: any, callback: any) {
    const loading = ElLoading.service({
        lock: true,
        text: tip,
        background: 'rgba(0, 0, 0, 0.2)',
    })
    await callback();
    loading.close();
}