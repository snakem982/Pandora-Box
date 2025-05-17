import express from "express";
import path from "path";
import log from './log';

const app = express();

let storedPort: string | undefined;
let storedSecret: string | undefined;
let listenAddr: string | undefined;
let goFlag: Function;

// 解析请求参数
app.use(express.urlencoded({extended: true}));
app.use(express.json());

// **提供静态文件服务**
app.use(express.static(path.join(__dirname, '../renderer/px_window')));

// **检测 PX 后端是否存活**
app.get("/pxAlive", (req, res) => {
    res.status(200).send("alive");
});

// **存储 PX 后端的端口和密钥**
// @ts-ignore
app.get("/pxStore", (req, res) => {
    const {port, secret} = req.query;

    if (!port || !secret) {
        return res.status(400).json({error: "缺少参数 port 或 secret"});
    }

    storedPort = port as string;
    storedSecret = secret as string;

    log.info("已获取 port:", storedPort);
    log.info("已获取 secret:", storedSecret);
    res.status(200).send("ok");

    if (goFlag) {
        goFlag(); // 通知主流程可以继续了
    }
});

// **处理所有未匹配的请求，返回 index.html**
app.use((req, res) => {
    res.redirect(302, '/index.html');
});


// **启动服务器**
export const startServer = (c1: Function, c2: Function) => {
    goFlag = c1
    const server = app.listen(0, "127.0.0.1", () => {
        // @ts-ignore
        listenAddr = `127.0.0.1:${server.address().port}`;
        c2(listenAddr)
    });
}

// 获取端口密钥
export const storeInfo = {
    port: () => storedPort,
    secret: () => storedSecret,
    listenAddr: () => listenAddr,
}
