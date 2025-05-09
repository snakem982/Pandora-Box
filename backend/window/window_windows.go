package window

/*
#cgo windows CFLAGS: -DUNICODE
#cgo windows LDFLAGS: -lole32 -lgdi32 -luser32
#include <windows.h>
#include <windowsx.h> // For GET_X_LPARAM, GET_Y_LPARAM
#include <shellapi.h>
#include <stdlib.h>

static WNDPROC originalWndProc;
static HWND globalHWND = NULL;
const int BORDER_WIDTH = 8;     // 定义可拖拽调整大小的边框宽度 (像素)
// const int TITLE_BAR_HEIGHT = 30; // 如果您想定义一个顶部拖拽区域的高度（可选）

// 自定义窗口过程
LRESULT CALLBACK WindowProc(HWND hwnd, UINT uMsg, WPARAM wParam, LPARAM lParam) {
    switch (uMsg) {
        case WM_CLOSE:
            ShowWindow(hwnd, SW_MINIMIZE); // 关闭按钮改成最小化
            return 0;

        case WM_NCCALCSIZE:
            // 当 wParam 为 TRUE 时，返回 0 会使客户区占据整个窗口。
            // 这会移除标准的窗口边框和标题栏，但保留 WS_THICKFRAME 样式的可调整大小特性。
            if (wParam == TRUE) {
                return 0;
            }
            // 对于 wParam == FALSE 的情况，让默认过程处理。
            break; // 注意：这里需要 break，然后由函数末尾的 CallWindowProc 处理

        case WM_NCHITTEST: {
            // 如果窗口已最大化，则通常不允许通过拖动边框来调整大小。
            if (IsZoomed(hwnd)) { // IsZoomed 检查窗口是否最大化
                // 可以调用原始窗口过程处理，或者根据需要返回 HTCAPTION (如果允许拖动最大化窗口)
                return CallWindowProc(originalWndProc, hwnd, uMsg, wParam, lParam);
            }

            POINT ptMouse = { GET_X_LPARAM(lParam), GET_Y_LPARAM(lParam) }; // 屏幕坐标
            RECT rcWindow;
            GetWindowRect(hwnd, &rcWindow); // 获取窗口的屏幕坐标

            // 检查鼠标是否在可调整大小的边框上 (仅当窗口未最大化时)
            BOOL onLeft = (ptMouse.x >= rcWindow.left) && (ptMouse.x < rcWindow.left + BORDER_WIDTH);
            BOOL onRight = (ptMouse.x < rcWindow.right) && (ptMouse.x >= rcWindow.right - BORDER_WIDTH);
            BOOL onTop = (ptMouse.y >= rcWindow.top) && (ptMouse.y < rcWindow.top + BORDER_WIDTH);
            BOOL onBottom = (ptMouse.y < rcWindow.bottom) && (ptMouse.y >= rcWindow.bottom - BORDER_WIDTH);

            if (onTop && onLeft) return HTTOPLEFT;
            if (onTop && onRight) return HTTOPRIGHT;
            if (onBottom && onLeft) return HTBOTTOMLEFT;
            if (onBottom && onRight) return HTBOTTOMRIGHT;
            if (onLeft) return HTLEFT;
            if (onRight) return HTRIGHT;
            if (onTop) return HTTOP;
            if (onBottom) return HTBOTTOM;

            // （可选）定义一个可拖拽的区域，例如顶部的自定义标题栏。
            // 如果您的 `pxDrag` 是从 JS 控制的，则此部分可能不是必需的。
            // 例如，使窗口顶部 BORDER_WIDTH 以下，TITLE_BAR_HEIGHT 高度的区域可拖拽：
            // if (ptMouse.y > rcWindow.top + BORDER_WIDTH && ptMouse.y < rcWindow.top + BORDER_WIDTH + TITLE_BAR_HEIGHT) {
            //     return HTCAPTION;
            // }


            // 如果鼠标不在我们定义的边框上，则认为是客户区。
            // 您的 pxDrag 绑定仍然可以从客户区内的HTML元素发起拖拽。
            return HTCLIENT;
        }

        case WM_SIZE: {
            // 窗口大小已更改。
            // int width = LOWORD(lParam);
            // int height = HIWORD(lParam);
            // 如果 webview 内容在调整大小后未正确刷新，可能需要强制重绘。
            // InvalidateRect(hwnd, NULL, TRUE);
            break;
        }
    }
    // 对于未处理的消息，调用原始窗口过程。
    return CallWindowProc(originalWndProc, hwnd, uMsg, wParam, lParam);
}

void setupWindow(HWND hwnd) {
    globalHWND = hwnd;
    originalWndProc = (WNDPROC)SetWindowLongPtr(hwnd, GWLP_WNDPROC, (LONG_PTR)WindowProc);

    LONG style = GetWindowLong(hwnd, GWL_STYLE);

    // 移除标准的视觉元素，但务必保留 WS_THICKFRAME 以便系统处理调整大小。
    // WS_THICKFRAME (也称为 WS_SIZEBOX) 是必需的，这样系统才能
    // 根据 WM_NCHITTEST 返回的 HTLEFT, HTRIGHT 等值来执行窗口大小调整操作。
    style &= ~(WS_CAPTION | WS_SYSMENU | WS_MINIMIZEBOX | WS_MAXIMIZEBOX | WS_BORDER);
    style |= WS_THICKFRAME; // 确保此样式存在，以实现可调整大小

    SetWindowLong(hwnd, GWL_STYLE, style);

    // 强制窗口重新计算其框架和客户区。
    // 上面处理的 WM_NCCALCSIZE 将使客户区填充整个窗口的大小。
    // 因此，这里不再需要之前的 GetClientRect/SetWindowPos 来手动调整客户区。
    SetWindowPos(hwnd, NULL, 0, 0, 0, 0,
                 SWP_FRAMECHANGED | SWP_NOMOVE | SWP_NOSIZE | SWP_NOZORDER | SWP_SHOWWINDOW);
}

// 设置任务栏图标
void setTaskbarIcon(HWND hwnd, void* data, int length) {
    HICON hIcon = CreateIconFromResource((PBYTE)data, length, TRUE, 0x00030000);
    if (hIcon != NULL) {
        SendMessage(hwnd, WM_SETICON, ICON_BIG, (LPARAM)hIcon);
        SendMessage(hwnd, WM_SETICON, ICON_SMALL, (LPARAM)hIcon);
        DestroyIcon(hIcon); // 创建的图标句柄使用后需要销毁
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
        ShowWindow(globalHWND, SW_RESTORE);
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
    int len = lstrlenW(wText);
    int size = WideCharToMultiByte(CP_UTF8, 0, wText, len, NULL, 0, NULL, NULL);
    char* buffer = (char*)malloc(size + 1);
    if (buffer != NULL) {
        WideCharToMultiByte(CP_UTF8, 0, wText, len, buffer, size, NULL, NULL);
        buffer[size] = '\0';
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
    sei.lpDirectory = NULL;
    sei.nShow = SW_HIDE;
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
	"github.com/snakem982/pandora-box/static"
	"os"
	"os/exec"
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
	// 获取系统
	_ = w.Bind("pxOs", geOs)
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

func geOs() string {
	return "Windows"
}
