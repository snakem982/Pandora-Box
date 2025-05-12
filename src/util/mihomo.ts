export function pUpdateMihomo(menuStore: any, settingStore: any, api: any): void {
    api.updateMihomo({
        mode: menuStore.rule,
        proxy: menuStore.proxy,
        tun: menuStore.tun,

        port: Number(settingStore.port),
        bindAddress: settingStore.bindAddress,
        stack: settingStore.stack,
        dns: settingStore.dns,
        ipv6: settingStore.ipv6,
    })
}