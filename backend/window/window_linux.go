package window

/*
#cgo linux CFLAGS:
#cgo linux LDFLAGS: -lX11
#include <X11/Xlib.h>
#include <X11/Xatom.h>
#include <stdlib.h>
#include <stdio.h>
#include <string.h>

// 全局变量保存 X11 的 Display 与 Window
static Display* globalDisplay = NULL;
static Window globalWindow = 0;

// 利用 _MOTIF_WM_HINTS 移除窗口装饰（标题栏、边框）
// 参考：https://stackoverflow.com/a/25409273/316319
typedef struct {
    unsigned long flags;
    unsigned long functions;
    unsigned long decorations;
    long inputMode;
    unsigned long status;
} MotifWmHints;
#define MWM_HINTS_DECORATIONS   (1L << 1)

void removeWindowDecorations(Display* display, Window window) {
    if (display == NULL || window == 0) return;
    MotifWmHints hints;
    hints.flags = MWM_HINTS_DECORATIONS;
    hints.decorations = 0;
    Atom property = XInternAtom(display, "_MOTIF_WM_HINTS", False);
    XChangeProperty(display, window, property, property, 32, PropModeReplace, (unsigned char *)&hints, 5);
    XFlush(display);
}

// 初始化窗口：传入 native 窗口句柄（来自 webview）后，打开 X11 display 并移除装饰
void setupWindow(void* winPtr) {
    globalWindow = (Window) winPtr;
    globalDisplay = XOpenDisplay(NULL);
    if(!globalDisplay) {
        fprintf(stderr, "Cannot open display\n");
        return;
    }
    removeWindowDecorations(globalDisplay, globalWindow);
    // 如果需要，可在此处加入其它窗口调整操作，例如调整客户区尺寸
}

// 设置任务栏图标：Linux 下通常需要将图片解码为 ARGB 数据后通过 _NET_WM_ICON 设置，较为复杂
// 此处提供 stub 实现，确保接口与 Windows 版本兼容
void setTaskbarIcon(void* winPtr, void* data, int length) {
    // stub：真正的实现需要PNG解码及调用 XChangeProperty 设置 _NET_WM_ICON
    (void)winPtr;
    (void)data;
    (void)length;
}

// 开始拖拽：由于窗口移动通常由窗口管理器处理，在 Linux 下“手工”模拟拖拽较困难
// 这里提供一个空的 stub 实现，或根据需求使用 XGrabPointer/XSendEvent 实现自定义拖拽
void beginDrag(void* winPtr) {
    (void)winPtr;
    // stub: 没有实现
}

// 显示窗口：使用 XMapRaised 将窗口映射并置顶
void showWindow() {
    if (globalDisplay && globalWindow) {
        XMapRaised(globalDisplay, globalWindow);
        XFlush(globalDisplay);
    }
}

// 隐藏窗口：使用 XUnmapWindow 隐藏窗口
void hideWindow() {
    if (globalDisplay && globalWindow) {
        XUnmapWindow(globalDisplay, globalWindow);
        XFlush(globalDisplay);
    }
}

// 获取剪贴板文本：利用 xclip 命令作为简化实现
char* getClipboardText() {
    FILE *fp = popen("xclip -selection clipboard -o 2>/dev/null", "r");
    char buffer[4096];
    size_t len = 0;

    if (fp != NULL) {
        len = fread(buffer, 1, sizeof(buffer) - 1, fp);
        pclose(fp);
    }

    // 如果使用 xclip 未获得有效内容，则尝试使用 xsel
    if (len == 0) {
        fp = popen("xsel --clipboard --output 2>/dev/null", "r");
        if (fp != NULL) {
            len = fread(buffer, 1, sizeof(buffer) - 1, fp);
            pclose(fp);
        }
    }

    if (len > 0) {
        buffer[len] = '\0';
        char* result = (char*)malloc(len + 1);
        if (result) {
            strcpy(result, buffer);
            return result;
        }
    }

    return NULL;
}


// 弹出消息框：调用 zenity 显示询问对话框，用户点击“确定”返回 1，否则返回 0
int ShowPrompt(const char* message, const char* title) {
    char command[1024];
    snprintf(command, sizeof(command), "zenity --question --title=\"%s\" --text=\"%s\"", title, message);
    int ret = system(command);
    // zenity执行成功（用户点击确定）返回 0
    return (ret == 0) ? 1 : 0;
}

// 提权执行：使用 pkexec 提升权限执行命令
int RunAsAdmin(const char* appPath, const char* cmdArgs) {
    char command[1024];
    // 构造命令：pkexec <appPath> <cmdArgs>
    snprintf(command, sizeof(command), "pkexec %s %s", appPath, cmdArgs);
    int ret = system(command);
    return ret;
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

// getNativeWindow 把 webview.WebView 的 Window() 转成 native X11 窗口句柄
func getNativeWindow(w webview.WebView) unsafe.Pointer {
	// 在 Linux 下，webview 返回的句柄通常为 X11 Window（unsigned long）类型
	return unsafe.Pointer(w.Window())
}

// Init 初始化窗口：设置标题、尺寸、调用 C.setupWindow 注入自定义处理，以及绑定各功能
func Init(w webview.WebView) {
	w.SetTitle("Pandora-Box")
	w.SetSize(1100, 760, webview.HintNone)

	// 调用 C 函数配置窗口（如移除装饰）
	C.setupWindow(getNativeWindow(w))
	SetDockIconBytes(w, Icon)

	_ = w.Bind("pxDrag", func() {
		Drag(w)
	})

	_ = w.Bind("pxClose", func() {
		// Linux 下没有 WM_CLOSE 消息，使用 hideWindow 作为替代
		C.hideWindow()
	})

	_ = w.Bind("pxMinimize", func() {
		// 采用 xdotool 方式最小化窗口（依赖 xdotool 工具）
		go exec.Command("xdotool", "windowminimize", fmt.Sprintf("%d", w.Window())).Run()
	})

	_ = w.Bind("pxHide", func() {
		C.hideWindow()
	})

	_ = w.Bind("pxToggleMaximize", func() {
		// 切换最大化在 X11 中较为依赖窗口管理器。这里通过 xdotool 模拟激活窗口，并发送 F11 键（依赖窗口管理器配置）
		go exec.Command("xdotool", "windowactivate", "--sync", fmt.Sprintf("%d", w.Window()), "key", "F11").Run()
	})

	_ = w.Bind("pxShowBar", func() {})

	_ = w.Bind("pxClipboard", GetClipboard)
	_ = w.Bind("pxOpen", Open)
	_ = w.Bind("pxConfigDir", openConfigDir)
	_ = w.Bind("pxOs", geOs)
}

// SetDockIconBytes 设置任务栏图标（stub 实现，Linux 下需要专门处理 _NET_WM_ICON）
func SetDockIconBytes(w webview.WebView, icon []byte) {
	if w.Window() != nil && len(icon) > 0 {
		C.setTaskbarIcon(getNativeWindow(w), unsafe.Pointer(&icon[0]), C.int(len(icon)))
	}
}

// Drag 支持拖拽窗口（stub 实现，真正拖拽需要较为复杂的 X11 交互）
func Drag(w webview.WebView) {
	if w.Window() != nil {
		C.beginDrag(getNativeWindow(w))
	}
}

// ShowWindow 显示窗口
func ShowWindow() {
	C.showWindow()
}

// HideWindow 隐藏窗口
func HideWindow() {
	C.hideWindow()
}

// GetClipboard 获取剪贴板中的文本（调用 xclip，需要系统中已安装 xclip 工具）
func GetClipboard() string {
	text := C.getClipboardText()
	if text == nil {
		return ""
	}
	defer C.free(unsafe.Pointer(text))
	return C.GoString(text)
}

// Open 使用 xdg-open 打开 URL
func Open(url string) error {
	return exec.Command("xdg-open", url).Start()
}

// RunAsAdmin 使用提权执行：首先用 zenity 弹出提示，再调用 pkexec (RunAsAdmin 在 C 中实现)
func RunAsAdmin(promptTitle, promptMessage, exe string, args ...string) error {
	// 如果已经是管理员则直接返回
	if sys.IsAdmin() {
		return errors.New("already admin")
	}

	argStr := strings.Join(args, " ")
	// 弹出提示：用户点确定才继续
	cMessage := C.CString(promptMessage)
	cTitle := C.CString(promptTitle)
	defer C.free(unsafe.Pointer(cMessage))
	defer C.free(unsafe.Pointer(cTitle))

	result := C.ShowPrompt(cMessage, cTitle)
	if result != 1 {
		return fmt.Errorf("用户取消操作")
	}

	cExe := C.CString(exe)
	cArgs := C.CString(argStr)
	defer C.free(unsafe.Pointer(cExe))
	defer C.free(unsafe.Pointer(cArgs))

	ret := C.RunAsAdmin(cExe, cArgs)
	if ret != 0 {
		return fmt.Errorf("提权执行失败，错误码：%d", int(ret))
	}
	return nil
}

// RunAsNoAdmin 在后台以普通权限启动命令
func RunAsNoAdmin(cmd string, args ...string) error {
	command := exec.Command(cmd, args...)
	command.SysProcAttr = &syscall.SysProcAttr{}
	command.Stdout = nil
	command.Stderr = nil

	err := command.Start()
	if err != nil {
		return fmt.Errorf("启动失败: %v", err)
	}
	return command.Process.Release()
}

// TryAdmin 尝试以管理员权限启动并获取页面信息，类似原 Windows 代码逻辑
func TryAdmin() string {
	server := static.StartFileServer()
	exePath, _ := os.Executable()
	tip := "Pandora-Box 需要授权才能使用 TUN 模式。[Pandora-Box requires authorization to enable TUN.]"
	err := RunAsAdmin("Pandora-Box", tip, exePath,
		"-back=true", "-addr="+server)
	// 如果提权失败则用普通用户权限运行
	if err != nil {
		_ = RunAsNoAdmin(exePath, "-back=true", "-addr="+server)
	}
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

// openConfigDir 打开配置目录：调用 xdg-open 打开用户目录
func openConfigDir() {
	c := exec.Command("xdg-open", utils.GetUserHomeDir())
	// Linux 下 HideWindow 并不适用，但这里保持了接口一致性
	c.SysProcAttr = &syscall.SysProcAttr{}
	_ = c.Start()
}

// geOs 返回当前运行系统名称
func geOs() string {
	return "Linux"
}
