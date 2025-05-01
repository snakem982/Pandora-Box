import {useMenuStore} from "@/store/menuStore";

export function changeMenu(value: string, router: any): void {
    let path = ''
    if (!value.startsWith("/")) {
        path = "/" + value
    }
    const menuStore = useMenuStore();
    // 对rule特殊处理
    if (path === "/Rule") {
        path += "/" + menuStore.ruleMenu;
    }
    menuStore.setPath(path);
    router.push(path);
}
