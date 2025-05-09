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
const int BORDER_WIDTH = 5; // Width of the draggable resize border

// 自定义窗口过程
LRESULT CALLBACK WindowProc(HWND hwnd, UINT uMsg, WPARAM wParam, LPARAM lParam) {
    switch (uMsg) {
        case WM_CLOSE:
            ShowWindow(hwnd, SW_MINIMIZE); // 关闭按钮改成最小化
            return 0;

        case WM_NCCALCSIZE:
            // Prevents the standard frame from being drawn when wParam is TRUE
            if (wParam == TRUE) {
                 // If you want to preserve some client area, adjust the NCCALCSIZE_PARAMS structure pointed to by lParam.
                 // For a completely custom frame, returning 0 tells Windows that the client area occupies the entire window.
                return 0;
            }
            break;

        case WM_NCHITTEST: {
            // Get the point coordinates for the hit test.
            POINT ptMouse = { GET_X_LPARAM(lParam), GET_Y_LPARAM(lParam) };

            // Get the window rectangle.
            RECT rcWindow;
            GetWindowRect(hwnd, &rcWindow);

            // Get the client rectangle.
            RECT rcClient;
            GetClientRect(hwnd, &rcClient);
            // Convert client rectangle to screen coordinates
            ClientToScreen(hwnd, (POINT*)&rcClient.left);
            ClientToScreen(hwnd, (POINT*)&rcClient.right);


            // Determine if the cursor is on the resize border.
            // The BORDER_WIDTH defines how thick the resize handles are.
            BOOL onLeft = ptMouse.x >= rcWindow.left && ptMouse.x < rcWindow.left + BORDER_WIDTH;
            BOOL onRight = ptMouse.x < rcWindow.right && ptMouse.x >= rcWindow.right - BORDER_WIDTH;
            BOOL onTop = ptMouse.y >= rcWindow.top && ptMouse.y < rcWindow.top + BORDER_WIDTH;
            BOOL onBottom = ptMouse.y < rcWindow.bottom && ptMouse.y >= rcWindow.bottom - BORDER_WIDTH;

            if (onTop && onLeft) return HTTOPLEFT;
            if (onTop && onRight) return HTTOPRIGHT;
            if (onBottom && onLeft) return HTBOTTOMLEFT;
            if (onBottom && onRight) return HTBOTTOMRIGHT;
            if (onLeft) return HTLEFT;
            if (onRight) return HTRIGHT;
            if (onTop) return HTTOP;
            if (onBottom) return HTBOTTOM;

            // If the cursor is within the client area (but not on the borders we defined),
            // and you want the main content area to be draggable like a caption,
            // you can return HTCAPTION.
            // This example assumes the drag is initiated by a specific UI element via `beginDrag`.
            // If you want the whole window (or a part of it) to be draggable, return HTCAPTION here.
            // For instance, if ptMouse is within rcClient after adjusting for borders:
            // if (ptMouse.x > rcWindow.left + BORDER_WIDTH && ptMouse.x < rcWindow.right - BORDER_WIDTH &&
            //     ptMouse.y > rcWindow.top + BORDER_WIDTH && ptMouse.y < rcClient.bottom - BORDER_WIDTH) // Assuming client area for caption
            // {
            //      // If you have a specific draggable region (e.g., a custom title bar)
            //      // you'd check if ptMouse is within that region's coordinates.
            //      // For now, let's assume beginDrag handles caption dragging.
            // }


            // If not on our custom borders or specific drag regions,
            // let the default procedure handle it, or return HTCLIENT for client area.
            // For a simple borderless window where dragging is handled by beginDrag,
            // and resizing by the borders above, areas not hit by borders should be client.
            // However, to make the window draggable from anywhere *not* a resize border:
            // This might be too broad if you have interactive elements.
            // A common pattern is to have a specific area (like a custom title bar) return HTCAPTION.
            // For now, we will rely on your existing `pxDrag` bind for dragging.
            // If no resize handle is hit, let the default processing occur, which might return HTCLIENT.
            LRESULT hit = CallWindowProc(originalWndProc, hwnd, uMsg, wParam, lParam);
            if (hit == HTCLIENT) {
                 // If you want the general client area (not resize borders) to be draggable, return HTCAPTION.
                 // This is a common approach for custom title bars.
                 // For example, if your "pxDrag" is bound to the entire window body:
                 // return HTCAPTION;
            }
            // Otherwise, return what the original proc decided or HTNOWHERE
            return hit == HTCLIENT ? HTNOWHERE : hit; // Or HTCLIENT if you don't want to make the whole client area HTCAPTION
        }

        case WM_SIZE: {
            // 可以在这里处理窗口大小变化后的操作
            // int width = LOWORD(lParam);
            // int height = HIWORD(lParam);
            // InvalidateRect(hwnd, NULL, TRUE); // May be useful to force repaint
            break;
        }
    }
    return CallWindowProc(originalWndProc, hwnd, uMsg, wParam, lParam);
}

void setupWindow(HWND hwnd) {
    globalHWND = hwnd;
    originalWndProc = (WNDPROC)SetWindowLongPtr(hwnd, GWLP_WNDPROC, (LONG_PTR)WindowProc);

    LONG style = GetWindowLong(hwnd, GWL_STYLE);
    // Remove caption, sysmenu, standard minimize/maximize buttons, and standard border.
    // We are removing WS_THICKFRAME as well, as we'll handle resizing via WM_NCHITTEST.
    style &= ~(WS_CAPTION | WS_SYSMENU | WS_MINIMIZEBOX | WS_MAXIMIZEBOX | WS_BORDER | WS_THICKFRAME);
    // Ensure it's a popup window if it's top-level, or child if it's meant to be embedded.
    // WS_POPUP is common for this kind of custom window.
    // style |= WS_POPUP; // If it's a top-level window. Often set during CreateWindowEx.

    SetWindowLong(hwnd, GWL_STYLE, style);

    // Adjust window size. For a borderless window (no WS_THICKFRAME, no WS_CAPTION),
    // the client area is the window area.
    // The initial size is set by webview.SetSize().
    // We might need to force a redraw or recalculation of styles.
    SetWindowPos(hwnd, NULL, 0, 0, 0, 0, SWP_NOMOVE | SWP_NOSIZE | SWP_NOZORDER | SWP_FRAMECHANGED | SWP_SHOWWINDOW);
    // The call to SetWindowPos with SWP_FRAMECHANGED is important to make the system recognize style changes.
    // SWP_NOMOVE and SWP_NOSIZE ensure we don't change position/size here,
    // relying on initial SetSize or subsequent resize operations.
}

// 设置任务栏图标
void setTaskbarIcon(HWND hwnd, void* data, int length) {
    HICON hIcon = CreateIconFromResource((PBYTE)data, length, TRUE, 0x00030000);
    if (hIcon != NULL) {
        SendMessage(hwnd, WM_SETICON, ICON_BIG, (LPARAM)hIcon);
        SendMessage(hwnd, WM_SETICON, ICON_SMALL, (LPARAM)hIcon);
        DestroyIcon(hIcon); // Important to destroy the icon handle after setting it
    }
}

// --- 支持区域拖拽窗口 ---
// This function is called from Go via w.Bind("pxDrag", ...)
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
    sei.nShow = SW_HIDE; // Consider SW_SHOW if you want to see the UAC prompt clearly
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
