import {useMenuStore} from "@/store/menuStore";

export function changeMenu(value: string, router: any): void {
    let path = ''
    if (!value.startsWith("/")) {
        path = "/" + value
    }
    const split = path.split("/");
    const menuStore = useMenuStore();
    menuStore.setMenu(split[1]);
    menuStore.setPath(path);
    router.push(path);
}
