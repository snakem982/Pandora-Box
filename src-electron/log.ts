import {app} from 'electron';
import * as path from 'path';
import * as fs from 'fs';
import log from 'electron-log/main';

// 获取用户根目录路径
const userHomeDir = app.getPath('home');

// 定义日志目录和文件路径
const logDir = path.join(userHomeDir, 'Pandora-Box-V3', 'logs');
const logFilePath = path.join(logDir, 'px-client.log');

// 确保目录存在
if (!fs.existsSync(logDir)) {
    fs.mkdirSync(logDir, {recursive: true});
}

// 自定义日志文件路径
log.transports.file.resolvePathFn = () => logFilePath;

// 设置日志等级和最大文件大小等配置
log.transports.file.level = 'info';  // 可以设置日志等级，比如 'info', 'warn', 'error'
log.transports.file.maxSize = 5 * 1024 * 1024;  // 设置最大文件大小为 5MB

// 日志记录功能的封装
export default {
    info: (...args: any[]) => {
        log.info(...args);
    },
    warn: (...args: any[]) => {
        log.warn(...args);
    },
    error: (...args: any[]) => {
        log.error(...args);
    },
    debug: (...args: any[]) => {
        log.debug(...args);
    },
    getLogFilePath: () => {
        return logFilePath;
    },
    getHomeDir: () => {
        return path.join(userHomeDir, 'Pandora-Box-V3');
    }
};
