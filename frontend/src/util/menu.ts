import {useMenuStore} from "@/store/menuStore";

export function changeMenu(value: string,router:any): void {
    const menuStore = useMenuStore();
    const split = value.split("/");
    menuStore.setMenu(split[0]);
    menuStore.setPath(value);
    router.push(value);
}
