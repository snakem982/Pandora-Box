// 排除的分组类型
const exclude = {
    DIRECT: true,
    REJECT: true,
    "REJECT-DROP": true,
    PASS: true,
    COMPATIBLE: true
}

// 不排除的节点类型
const include = {
    Direct: true,
    Reject: true,
    RejectDrop: true,
    URLTest: true,
    LoadBalance: true,
    Selector: true,
    Pass: true,
    Relay: true,
    Fallback: true,
}

// 计算类名
const getClass = (delay: any) => {
    if (delay === 99999) {
        return 'toHidden'
    }

    if (delay <= 300) {
        return 'toLow'
    } else if (delay <= 600) {
        return 'toMiddle'
    } else {
        return 'toHigh'
    }
}

// 获取节点延迟
const getDelay = (proxy: any) => {
    if (!proxy['alive']) {
        return 99999;
    }

    const history = proxy['history']
    if (!history || history.length === 0) {
        return 99999;
    }

    return history[history.length - 1]['delay']
}

export default function createProxiesApi(proxy: any) {
    return {
        // 获取分组延迟
        async getDelay(group: any, url: any, timeout: any) {
            await proxy.$http.get('/group/' + group + '/delay?timeout=' + timeout + "&url=" + url);
        },
        // 获取分组列表
        async getGroups() {
            // 获取所有节点分组列表
            const data = await proxy.$http.get('/proxies');
            const proxies = data['proxies']

            // 判空
            if (!proxies['GLOBAL']) {
                return []
            }

            // 获取分组
            const group = []
            for (const name of proxies['GLOBAL']['all']) {
                if (exclude[name]) {
                    continue
                }
                if (!include[proxies[name]['type']]) {
                    continue
                }
                group.push(name)
            }

            return group
        },
        // 获取相应的分组节点列表
        async getProxies(active: string, isHide: boolean, isSort: boolean) {
            // 获取所有节点分组列表
            const data = await proxy.$http.get('/proxies')
            const proxies = data['proxies']

            // 判空
            if (!proxies[active]) {
                return []
            }

            // 获取分组节点列表
            const proxiesNames = proxies[active]['all']
            const nowName = proxies[active]['now']

            // 获取节点延迟
            const activeProxies = []
            const inProxies = []
            for (const name of proxiesNames) {
                const proxy = proxies[name]
                const type = proxy['type'];
                const delay = getDelay(proxy)
                if (include[type]) {
                    inProxies.push({
                        name,
                        type,
                        delay: delay,
                        now: name === nowName,
                        toClass: getClass(delay)
                    })
                } else {
                    activeProxies.push({
                        name,
                        type,
                        delay,
                        now: name === nowName,
                        toClass: getClass(delay)
                    })
                }
            }

            // 获取显示的节点
            let showProxies = []
            if (isHide) {
                for (const proxy of activeProxies) {
                    if (proxy['delay'] != 99999) {
                        showProxies.push(proxy)
                    }
                }
            } else {
                showProxies = activeProxies
            }

            // 构建哈希表
            const GLOBAL = proxies['GLOBAL']['all'];
            const map = new Map();
            GLOBAL.forEach((value: any, index: any) => {
                map.set(value, index);
            });

            // 进行排序
            if (isSort) {
                inProxies.sort((obj1, obj2) => {
                    if (obj1.delay != obj2.delay) {
                        return obj1.delay - obj2.delay
                    }

                    return map.get(obj1.name) - map.get(obj2.name)
                });
                showProxies.sort((obj1, obj2) => obj1.delay - obj2.delay);
            } else {
                showProxies.sort((obj1, obj2) => map.get(obj1.name) - map.get(obj2.name));
                inProxies.sort((obj1, obj2) => map.get(obj1.name) - map.get(obj2.name));
            }

            return inProxies.concat(showProxies);
        },
        // 设置代理
        async setProxy(group: any, name: any) {
            await proxy.$http.put("/proxies/" + group, name);
        },
    };
}
