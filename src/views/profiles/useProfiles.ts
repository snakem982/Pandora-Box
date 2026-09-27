// src/views/profiles/useProfiles.ts
import {ref, reactive, toRaw, watch, onMounted, onBeforeUnmount, getCurrentInstance} from 'vue'
import {useI18n} from 'vue-i18n'
import {onBeforeRouteLeave} from 'vue-router'
import createApi from '@/api'
import {pError, pLoad, pSuccess, pWarning} from '@/util/pLoad'
import {useProxiesStore} from '@/store/proxiesStore'
import {useMenuStore} from '@/store/menuStore'
import {useWebStore} from '@/store/webStore'
import {prettyBytes} from '@/util/format'
import {WS} from '@/util/ws'
import {Events, Browser} from '@/runtime'
import {useProfileStore} from "@/store/profileStore";

export function useProfiles() {
    const {t} = useI18n()
    const {proxy} = getCurrentInstance()!
    const api = createApi(proxy)

    const menuStore = useMenuStore()
    const proxiesStore = useProxiesStore()
    const webStore = useWebStore()
    const profileStore = useProfileStore()

    const profiles = reactive<any[]>([])
    const tList = ref<any[]>([])
    let wsOrder: WS

    const headerShow = reactive({
        available: '',
        used: '',
        expire: '',
        update: '',
    })

    // 更新头部统计显示
    function setHeaderShow(item: any) {
        headerShow.available = item['available'] ? prettyBytes(item['available']) : ''
        headerShow.used = item['used'] ? prettyBytes(item['used']) : ''
        headerShow.expire = item['expire'] || ''
        headerShow.update = item['update'] || ''
    }

    // 安全数值转化
    function num2SafeNumber(data: any, key: string) {
        if (data[key] !== undefined && data[key] !== null) {
            let num = Number(data[key])
            if (!Number.isFinite(num)) return
            if (num > Number.MAX_SAFE_INTEGER) data[key] = Number.MAX_SAFE_INTEGER
            else if (num < Number.MIN_SAFE_INTEGER) data[key] = Number.MIN_SAFE_INTEGER
            else data[key] = num
        }
    }

    // WS 发送排序/更新
    function sendOrder(data: any) {
        if (wsOrder) {
            Events.Emit({name: 'profiles', data: toRaw(data)})
            for (let i = 0; i < data.length; i++) {
                num2SafeNumber(data[i], 'available')
                num2SafeNumber(data[i], 'used')
                num2SafeNumber(data[i], 'total')
            }
            wsOrder.send(JSON.stringify(data))
        }
    }

    // 获取 Profile 列表
    async function getProfileList() {
        if (profiles.length !== 0) profiles.splice(0, profiles.length)
        const list = await api.getProfileList()
        if (list && list.length !== 0) {
            list.forEach((item) => {
                profiles.push(item)
                if (item['selected']) setHeaderShow(item)
            })
            Events.Emit({name: 'profiles', data: list})
        }
    }

    // 刷新/更新订阅配置
    async function refreshProfile(data: any) {
        await pLoad(t('profiles.refresh.ing'), async () => {
            try {
                const re = await api.refreshProfile(data)
                if (data['selected']) {
                    setHeaderShow(re)
                }
                Object.assign(data, re)
                pSuccess(t('profiles.refresh.success'))
            } catch (e: any) {
                if (e['message']) pError(e['message'])
            }
        })
    }

    // 切换 Profile
    async function switchProfile(data: any) {
        if (data['selected']) return
        await pLoad(t('profiles.switch.ing'), async () => {
            try {
                await api.switchProfile(data)
                proxiesStore.active = ''
                await api.waitRunning()

                for (let profile of profiles) {
                    if (profile['selected']) profile['selected'] = false
                }
                data['selected'] = true
                setHeaderShow(data)

                api.getRuleNum().then((res) => menuStore.setRuleNum(res))
                Events.Emit({name: 'profiles', data: toRaw(profiles)})
                api.closeAllConnection()
                pSuccess(t('profiles.switch.success'))
            } catch (e: any) {
                if (e['message']) pError(e['message'])
            }
        })
    }

    // 删除 Profile
    async function deleteProfile(data: any, index: number) {
        if (data['selected']) {
            pWarning(t('profiles.del-tip'))
            return
        }
        try {
            await api.deleteProfile(data)
            profiles.splice(index, 1)
            Events.Emit({name: 'profiles', data: toRaw(profiles)})
        } catch (e: any) {
            if (e['message']) pError(e['message'])
        }
    }

    // 浏览器打开链接（主页/支持）
    function goHome(data: any) {
        if (data?.home) Browser.OpenURL(data.home)
    }

    function goSupport(data: any) {
        if (data?.support) Browser.OpenURL(data.support)
    }

    // 监听外部选中的 profile 变化
    watch(() => webStore.fProfile, async (data: any) => {
        if (!data) return
        for (let profile of profiles) {
            if (profile['selected']) {
                profile['selected'] = false
            }
            if (profile['id'] == data['id']) {
                data = profile
            }
        }
        data['selected'] = true
        setHeaderShow(data)
    })

    // 监听文件拖拽/导入的 profile 列表变化
    watch(() => webStore.dProfile, async (pList) => {
        if (pList && pList.length > 0) {
            pList.forEach((item) => profiles.push(item))
        }
    })

    // 生命周期 & 路由钩子
    onMounted(async () => {
        const urlTraffic = `${webStore.wsUrl}/profile/order?token=${webStore.secret}`
        wsOrder = new WS(urlTraffic)

        await getProfileList()
        const templates = await api.getTemplateList()
        tList.value = [{title: 'm0', id: 'm0'}, ...templates]
    })

    onBeforeRouteLeave(() => wsOrder?.close())
    onBeforeUnmount(() => wsOrder?.close())

    return {
        api,
        profiles,
        headerShow,
        tList,
        menuStore,
        webStore,
        profileStore,
        sendOrder,
        getProfileList,
        switchProfile,
        deleteProfile,
        refreshProfile,
        goHome,
        goSupport
    }
}