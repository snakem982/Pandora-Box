export class Profile {
    id!: string;
    type!: number; // 1: 远程订阅, 2: 本地配置, 3: 爬取合并
    title?: string; // 可选
    order!: string;
    selected?: boolean; // 可选
    path!: string;
    content?: string | ArrayBuffer; // 可选
    used?: bigint; // 可选
    available?: bigint; // 可选
    total?: bigint; // 可选
    expire?: string; // 可选
    interval?: string; // 可选
    home?: string; // 可选
    update?: string; // 可选
    template?: string; // 可选
}
