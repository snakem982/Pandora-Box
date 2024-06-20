export class Group {
    constructor(proxies: Proxy[]) {
        this.proxies = proxies;
    }

    proxies: Proxy[];
}

export class Proxy {
    constructor(alive: boolean, all: string[], extra: any, history: History[], name: string, now: string, tfo: boolean, type: string, udp: boolean, xudp: boolean) {
        this.alive = alive;
        this.all = all;
        this.extra = extra;
        this.history = history;
        this.name = name;
        this.now = now;
        this.tfo = tfo;
        this.type = type;
        this.udp = udp;
        this.xudp = xudp;
    }

    alive: boolean;
    all: string[];
    extra: any;
    history: History[];
    name: string;
    now: string;
    tfo: boolean;
    type: string;
    udp: boolean;
    xudp: boolean;
}

export class History {
    constructor(time: string, delay: number) {
        this.time = time;
        this.delay = delay;
    }

    time: string;
    delay: number;
}

export interface Metadata {
    network: string;
    type: string;
    sourceIP: string;
    destinationIP: string;
    sourcePort: string;
    destinationPort: string;
    inboundIP: string;
    inboundPort: string;
    inboundName: string;
    inboundUser: string;
    host: string;
    dnsMode: string;
    uid: number;
    process: string;
    processPath: string;
    specialProxy: string;
    specialRules: string;
    remoteDestination: string;
    sniffHost: string;
}

export interface Connection {
    id: string;
    metadata: Metadata;
    upload: number;
    download: number;
    start: string;
    chains: string[];
    rule: string;
    rulePayload: string;
}

export interface Connections {
    downloadTotal: number;
    uploadTotal: number;
    connections: Connection[];
    memory: number;
}

export interface Profile {
    id: string;
    type: number;
    title: string;
    path: string;
    url: string;
    order: number;
    selected: boolean;
}