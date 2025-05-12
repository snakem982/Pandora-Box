// 导入组件
import {App, Component} from 'vue'
import VDContainer from './VDContainer/src/VDContainer.vue'

interface comp {
    [keys: string]: Component
}

// 存储组件列表
const components: comp = {
    VDContainer
}

function install(Vue: App): void {
    const keys = Object.keys(components)
    keys.forEach(name => {
        const component = components[name]
        Vue.component(component.name || name, component)
    })
}

export default {
    // 导出的对象必须具有 install，才能被 Vue.use() 方法安装
    install,
    // 以下是具体的组件列表
    ...components
}
