import path from "node:path";
import {spawn} from "child_process";
import {app, dialog} from "electron";
import fs from "node:fs";
import log from './log';
import {storeGet} from "./store";

// 是否在开发模式
const isDev = !app.isPackaged;

// 获取px路径
function getBackendPath() {
    const execName = process.platform === 'win32' ? 'px.exe' : 'px';
    return isDev
        ? path.join(__dirname, '../../src-go', execName)
        : path.join(process.resourcesPath, execName);
}

// 检查是否有管理员权限
function checkAdminRights(callback: any) {
    const platform = process.platform;

    if (platform === 'win32') {
        // Windows 上检查管理员权限
        const command = spawn('net', ['session']);

        command.on('error', (error) => {
            log.info('net session :', error);
            callback(false);
        });

        command.on('exit', (code) => {
            if (code === 0) {
                callback(true);   // 如果退出码为 0，表示有管理员权限
            } else {
                callback(false);  // 否则没有管理员权限
            }
        });

    } else if (platform === 'darwin' || platform === 'linux') {
        // macOS 或 Linux 上检查是否为 root 用户
        if (process.getuid && process.getuid() === 0) {
            callback(true);  // 有管理员权限
        } else {
            callback(false); // 没有管理员权限
        }
    } else {
        // 其他平台默认认为没有管理员权限
        callback(false);
    }
}

// 开启后端
export function startBackend(addr: string) {
    const backendPath = getBackendPath();
    const args = ['-addr=' + addr, '-home=' + encodeURIComponent(log.getHomeDir())];

    // 检查管理权限
    checkAdminRights((isAdmin: boolean) => {
        if (isAdmin) {
            log.info('有管理员权限');

            startNormally(backendPath, args);
        } else {
            log.info('没有管理员权限');

            // 判断是否关闭提权
            const setting = storeGet("setting")
            if (!!setting) {
                const set = JSON.parse(setting as string);
                if (set.hasOwnProperty("auth") && !set["auth"]) {
                    log.info('关闭了提权提示');
                    startNormally(backendPath, args);
                    return;
                }
            }

            // 只在 Windows 和 Linux 平台上弹出提权提示，macOS 也需要显示提权提示
            if (process.platform !== 'darwin') {
                const tip = "Px 需要授权才能使用 TUN 模式。\n[Px requires authorization to enable TUN.]";
                const confirmed = dialog.showMessageBoxSync({
                    type: 'info',
                    buttons: ['继续', '取消'],
                    defaultId: 0,
                    cancelId: 1,
                    title: 'Pandora-Box',
                    message: tip,
                });

                if (confirmed === 1) {
                    // 用户取消提权 → 普通模式启动
                    log.info('用户取消了提权，使用普通权限启动');
                    startNormally(backendPath, args);
                    return;
                }
            }

            // 尝试以管理员权限运行，失败则降级
            tryRunAsAdmin(backendPath, args, (success) => {
                if (!success) {
                    log.info('管理员权限启动失败，使用普通模式启动');
                    startNormally(backendPath, args);
                }
            });
        }
    });
}

// 尝试管理员启动
function tryRunAsAdmin(executable: string, args: string[], callback: (success: boolean) => void) {
    switch (process.platform) {
        case 'darwin': {
            // macOS 使用 AppleScript 提权
            const tip = "Px 需要授权才能使用 TUN 模式。\n[Px requires authorization to enable TUN.]";
            const command = `${[executable, ...args].map(escapeShell).join(' ')}`;
            // 使用 `with prompt` 来直接在授权对话框中显示提示信息
            const script = `
                do shell script "${command}" with administrator privileges with prompt "${tip}"
            `;
            const osa = spawn('osascript', ['-e', script]);
            log.info("[Admin] 启动px命令行：", osa.spawnargs);
            osa.on('exit', (code) => callback(code === 0));
            osa.on('error', () => callback(false));
            break;
        }

        case 'win32': {
            // Windows 使用 PowerShell 提权并隐藏窗口
            const psArgs = [
                '-Command',
                `Start-Process -FilePath '${executable}' -ArgumentList '${args.join(' ')}' -Verb RunAs -WindowStyle Hidden`
            ];
            const ps = spawn('powershell.exe', psArgs);
            log.info("[Admin] 启动px命令行：", ps.spawnargs);
            ps.on('exit', (code) => callback(code === 0));
            ps.on('error', () => callback(false));
            break;
        }

        case 'linux': {
            // Linux: 提权方式依次尝试 pkexec → gksudo → kdesudo → sudo
            const env = {
                ...process.env,
                PATH: process.env.PATH || "/usr/bin:/bin:/usr/sbin:/sbin",
                DISPLAY: process.env.DISPLAY,
                XAUTHORITY: process.env.XAUTHORITY,
            };

            const methods = [
                '/usr/bin/pkexec',
                '/usr/bin/gksudo',
                '/usr/bin/kdesudo',
                '/usr/bin/sudo',
                'pkexec',
                'gksudo',
                'kdesudo',
                'sudo'
            ];

            let tried = false;

            (function tryNext(index = 0) {
                if (index >= methods.length) {
                    log.error("No available elevation method succeeded.");
                    callback(false);
                    return;
                }

                const method = methods[index];
                if (!fs.existsSync(method) && !method.includes('/')) {
                    // Skip fallback names like 'sudo' if not full path
                    return tryNext(index + 1);
                }

                log.info(`Trying to elevate with: ${method}`);
                tried = true;

                const elevated = spawn(method, [executable, ...args], {
                    env,
                    stdio: 'inherit',
                });
                log.info("[Admin] 启动px命令行：", elevated.spawnargs);

                elevated.on('error', (err) => {
                    log.error(`Error using ${method}:`, err);
                    tryNext(index + 1);
                });

                elevated.on('exit', (code, signal) => {
                    if (code === 0) {
                        log.info(`${method} succeeded`);
                        callback(true);
                    } else {
                        log.warn(`${method} exited with code ${code} or signal ${signal}`);
                        tryNext(index + 1);
                    }
                });
            })();

            break;
        }


        default:
            log.error('不支持的平台:', process.platform);
            callback(false);
    }
}

function startNormally(executable: string, args: string[]) {
    const backend = spawn(executable, args, {
        stdio: ['ignore', 'pipe', 'pipe']
    });

    log.info("[Normal] 启动px命令行：", backend.spawnargs);

    backend.stdout.on('data', () => {
        // log.info(`[backend stdout]: ${data}`);
    });

    backend.stderr.on('data', (data) => {
        log.error(`[backend stderr]: ${data}`);
    });

    backend.on('error', (err) => log.error('Backend error:', err));
    backend.on('exit', (code) => log.info('Backend exited with code:', code));
}

function escapeShell(cmd: string): string {
    return cmd.replace(/"/g, '\\"').replace(/\$/g, '\\$');
}
