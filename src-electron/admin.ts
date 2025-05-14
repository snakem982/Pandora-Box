import path from "node:path";
import {spawn} from "child_process";
import {app} from "electron";

// 是否在开发模式
const isDev = !app.isPackaged;

function getBackendPath() {
    const execName = 'px';

    return isDev
        ? path.join(__dirname, '../../src-go', execName)
        : path.join(process.resourcesPath, execName);
}

export function startBackend(addr: string) {
    const backendPath = getBackendPath();

    const backend = spawn(backendPath, ['-addr=' + addr], {
        stdio: ['ignore', 'pipe', 'pipe']
    });

    backend.stdout.on('data', () => {
        // console.log(`[backend stdout]: ${data}`);
    });

    backend.stderr.on('data', (data) => {
        console.error(`[backend stderr]: ${data}`);
    });

    backend.on('error', (err) => console.error('Backend error:', err));
    backend.on('exit', (code) => console.log('Backend exited with code:', code));
}