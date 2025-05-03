package window

/*
#cgo windows CFLAGS: -DUNICODE
#cgo windows LDFLAGS: -lole32 -lgdi32 -luser32
#include <windows.h>
#include <shellapi.h>
#include <stdlib.h>
#include <shlwapi.h>

static WNDPROC originalWndProc;
static HWND globalHWND = NULL;

// 自定义窗口过程，拦截关闭按钮和调整大小
LRESULT CALLBACK WindowProc(HWND hwnd, UINT uMsg, WPARAM wParam, LPARAM lParam) {
    if (uMsg == WM_CLOSE) {
        ShowWindow(hwnd, SW_MINIMIZE); // 关闭按钮改成最小化
        return 0;
    }

    if (uMsg == WM_SIZE) {
        // 可以在这里处理窗口大小变化后的操作
        int width = LOWORD(lParam);
        int height = HIWORD(lParam);
        // 可以在此处添加额外的逻辑，例如更新 UI 或其他操作
    }

    return CallWindowProc(originalWndProc, hwnd, uMsg, wParam, lParam);
}

void setupWindow(HWND hwnd) {
	// 保存全局 hwnd
	globalHWND = hwnd;
    // 保存原始的 WndProc
    originalWndProc = (WNDPROC)SetWindowLongPtr(hwnd, GWLP_WNDPROC, (LONG_PTR)WindowProc);

    // 取当前窗口样式
    LONG style = GetWindowLong(hwnd, GWL_STYLE);

    // 移除标题栏、系统菜单、最小化按钮、最大化按钮
    style &= ~(WS_CAPTION | WS_SYSMENU | WS_MINIMIZEBOX | WS_MAXIMIZEBOX | WS_BORDER);

    // 允许调整窗口大小（保持可调整大小的边框）
    style |= WS_THICKFRAME;

    // 设置新的窗口样式
    SetWindowLong(hwnd, GWL_STYLE, style);

    // 获取窗口的客户区域（客户端区域）
    RECT rect;
    GetClientRect(hwnd, &rect);
    int clientWidth = rect.right - rect.left;
    int clientHeight = rect.bottom - rect.top;

    // 重新调整窗口尺寸，确保不出现白边
    SetWindowPos(hwnd, NULL, 0, 0, clientWidth, clientHeight, SWP_NOMOVE | SWP_NOZORDER | SWP_FRAMECHANGED);
}


// 设置任务栏图标
void setTaskbarIcon(HWND hwnd, void* data, int length) {
    HICON hIcon = CreateIconFromResource((PBYTE)data, length, TRUE, 0x00030000);
    if (hIcon != NULL) {
        SendMessage(hwnd, WM_SETICON, ICON_BIG, (LPARAM)hIcon);
        SendMessage(hwnd, WM_SETICON, ICON_SMALL, (LPARAM)hIcon);
    }
}

// --- 支持区域拖拽窗口 ---
void beginDrag(HWND hwnd) {
    ReleaseCapture();
    SendMessage(hwnd, WM_NCLBUTTONDOWN, HTCAPTION, 0);
}

// 显示窗口
void showWindow() {
    if (globalHWND != NULL) {
        ShowWindow(globalHWND, SW_RESTORE); // 从最小化或最大化恢复
        SetForegroundWindow(globalHWND);
    }
}

// 隐藏窗口
void hideWindow() {
    if (globalHWND != NULL) {
        ShowWindow(globalHWND, SW_HIDE);
    }
}

char* getClipboardText() {
    if (!OpenClipboard(NULL)) {
        return NULL;
    }

    HANDLE hData = GetClipboardData(CF_UNICODETEXT);
    if (hData == NULL) {
        CloseClipboard();
        return NULL;
    }

    wchar_t* wText = (wchar_t*)GlobalLock(hData);
    if (wText == NULL) {
        CloseClipboard();
        return NULL;
    }

    // 计算需要的缓冲区大小
    int len = lstrlenW(wText);
    int size = WideCharToMultiByte(CP_UTF8, 0, wText, len, NULL, 0, NULL, NULL);

    char* buffer = (char*)malloc(size + 1);
    if (buffer != NULL) {
        WideCharToMultiByte(CP_UTF8, 0, wText, len, buffer, size, NULL, NULL);
        buffer[size] = '\0'; // Null-terminate
    }

    GlobalUnlock(hData);
    CloseClipboard();

    return buffer;
}

// 弹出消息框
int ShowPrompt(const wchar_t* message, const wchar_t* title) {
    return MessageBoxW(NULL, message, title, MB_OKCANCEL | MB_ICONEXCLAMATION | MB_TOPMOST);
}

// 提权执行
int RunAsAdmin(const wchar_t* appPath, const wchar_t* cmdArgs) {
    SHELLEXECUTEINFOW sei;
    ZeroMemory(&sei, sizeof(sei));
    sei.cbSize = sizeof(sei);
    sei.fMask = SEE_MASK_NOCLOSEPROCESS;
    sei.hwnd = NULL;
    sei.lpVerb = L"runas";
    sei.lpFile = appPath;
    sei.lpParameters = cmdArgs;

	// 设置程序所在目录作为工作目录
    wchar_t dir[MAX_PATH];
    wcscpy_s(dir, MAX_PATH, appPath)
    PathRemoveFileSpecW(dir);
    sei.lpDirectory = dir;

    sei.nShow = SW_SHOWNORMAL;

    if (!ShellExecuteExW(&sei)) {
        return GetLastError();
    }

	if (sei.hProcess) {
		CloseHandle(sei.hProcess);
	}

    return 0;
}

*/
import "C"

import (
	_ "embed"
	"errors"
	"fmt"
	sys "github.com/snakem982/pandora-box/pkg/sys/admin"
	"github.com/snakem982/pandora-box/pkg/utils"
	"os"
	"os/exec"
	"pandora-box/static"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"github.com/webview/webview_go"
)

//go:embed icon-128.png
var Icon []byte

// getHWND 把 webview.WebView 的 Window() 转成 HWND
func getHWND(w webview.WebView) C.HWND {
	return (C.HWND)(unsafe.Pointer(w.Window()))
}

// Init 初始化窗口，设置大小 + 注入窗口过程
func Init(w webview.WebView) {
	// 设置初始标题和大小
	w.SetTitle("Pandora-Box")
	w.SetSize(1100, 760, webview.HintNone)

	// 注入自定义窗口过程
	C.setupWindow(getHWND(w))
	SetDockIconBytes(w, Icon)

	_ = w.Bind("pxDrag", func() {
		Drag(w)
	})

	_ = w.Bind("pxClose", func() {
		C.PostMessage(getHWND(w), C.WM_CLOSE, 0, 0)
	})

	_ = w.Bind("pxMinimize", func() {
		C.ShowWindow(getHWND(w), C.SW_MINIMIZE)
	})

	_ = w.Bind("pxHide", func() {
		C.ShowWindow(getHWND(w), C.SW_HIDE)
	})

	_ = w.Bind("pxToggleMaximize", func() {
		hwnd := getHWND(w)
		style := C.GetWindowLongPtr(hwnd, C.GWL_STYLE)
		if style&C.WS_MAXIMIZE == 0 {
			C.ShowWindow(hwnd, C.SW_MAXIMIZE)
		} else {
			C.ShowWindow(hwnd, C.SW_RESTORE)
		}
	})

	// 是否显示 自定义标题
	_ = w.Bind("pxShowBar", func() {})

	// 剪贴板
	_ = w.Bind("pxClipboard", GetClipboard)
	// 打开地址
	_ = w.Bind("pxOpen", Open)
	// 打开配置目录
	_ = w.Bind("pxConfigDir", openConfigDir)
}

// SetDockIconBytes 设置任务栏图标
func SetDockIconBytes(w webview.WebView, icon []byte) {
	if w.Window() != nil && len(icon) > 0 {
		ptr := unsafe.Pointer(&icon[0])
		C.setTaskbarIcon(getHWND(w), ptr, C.int(len(icon)))
	}
}

// RefreshMenu 刷新菜单（空实现，兼容接口）
func RefreshMenu(lang string) {}

// Drag 用来支持拖拽（如果需要）
func Drag(w webview.WebView) {
	if w.Window() != nil {
		C.beginDrag(getHWND(w))
	}
}

func ShowWindow() {
	C.showWindow()
}

func HideWindow() {
	C.hideWindow()
}

// GetClipboard 返回剪贴板内容
func GetClipboard() string {
	text := C.getClipboardText()
	if text == nil {
		return ""
	}
	defer C.free(unsafe.Pointer(text))
	return C.GoString(text)
}

func Open(url string) error {
	u := C.CString(url)
	verb := C.CString("open")
	defer C.free(unsafe.Pointer(u))
	defer C.free(unsafe.Pointer(verb))

	C.ShellExecuteA(
		(*C.struct_HWND__)(unsafe.Pointer(uintptr(0))), // 把0安全地转成 HWND
		verb,
		u,
		nil,
		nil,
		C.SW_SHOWNORMAL,
	)
	return nil
}

// 转 wchar
func utf16Ptr(s string) (*C.wchar_t, error) {
	u16, err := syscall.UTF16FromString(s)
	if err != nil {
		return nil, err
	}
	return (*C.wchar_t)(unsafe.Pointer(&u16[0])), nil
}

// RunAsAdmin 带提示 + 提权执行
func RunAsAdmin(promptTitle, promptMessage, exe string, args ...string) error {
	// 已经是管理员直接返回
	if sys.IsAdmin() {
		return errors.New("already admin")
	}

	exePtr, err := utf16Ptr(exe)
	if err != nil {
		return err
	}

	argStr := strings.Join(args, " ")
	argPtr, err := utf16Ptr(argStr)
	if err != nil {
		return err
	}

	titlePtr, _ := utf16Ptr(promptTitle)
	msgPtr, _ := utf16Ptr(promptMessage)

	result := C.ShowPrompt(msgPtr, titlePtr)
	if result != C.IDOK {
		return fmt.Errorf("用户取消操作")
	}

	ret := C.RunAsAdmin(exePtr, argPtr)
	if ret != 0 {
		return fmt.Errorf("提权执行失败，错误码：%d", int(ret))
	}
	return nil
}

// RunAsNoAdmin 在后台运行普通用户权限的命令
func RunAsNoAdmin(cmd string, args ...string) error {
	command := exec.Command(cmd, args...)

	// 释放进程，不阻塞主线程
	command.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	command.Stdout = nil
	command.Stderr = nil

	err := command.Start()
	if err != nil {
		return fmt.Errorf("启动失败: %v", err)
	}

	// 释放进程句柄，确保不阻塞主线程
	err = command.Process.Release()
	if err != nil {
		return fmt.Errorf("释放进程失败: %v", err)
	}

	return nil
}

func TryAdmin() string {
	server := static.StartFileServer()

	// 尝试提权
	exePath, _ := os.Executable()
	tip := "Pandora-Box 需要授权才能使用 TUN 模式。[Pandora-Box requires authorization to enable TUN.]"
	err := RunAsAdmin("Pandora-Box", tip, exePath,
		"-back=true", "-addr="+server)

	// 提权失败，普通运行
	if err != nil {
		_ = RunAsNoAdmin(exePath, "-back=true", "-addr="+server)
	}

	// 获取端口加载页面
	var port interface{}
	var secret interface{}
	for {
		port = static.Get("port")
		secret = static.Get("secret")
		if port != nil && secret != nil {
			break
		}
		time.Sleep(1 * time.Second)
	}

	return fmt.Sprintf("http://%s/index.html?port=%v&secret=%v", server, port, secret)
}

func openConfigDir() {
	c := exec.Command(`cmd`, `/c`, `explorer`, utils.GetUserHomeDir())
	c.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	_ = c.Start()
}
