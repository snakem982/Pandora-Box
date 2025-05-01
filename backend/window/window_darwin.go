package window

/*
#cgo darwin CFLAGS: -x objective-c -I./window
#cgo darwin LDFLAGS: -framework Cocoa
// 引入必要的头文件
#import <Cocoa/Cocoa.h>
#import <objc/runtime.h>

// 用二进制数据设置 Dock 图标
void setDockIconFromBytes(void* data, int length) {
    NSData *nsdata = [NSData dataWithBytes:data length:length];
    NSImage *image = [[NSImage alloc] initWithData:nsdata];
    if (image) {
        [NSApp setApplicationIconImage:image];
    }
}

NSWindow *globalWindow = nil;

// --- 自定义 WindowDelegate 处理关闭按钮 ---
@interface WindowDelegate : NSObject <NSWindowDelegate>
@end

@implementation WindowDelegate
- (BOOL)windowShouldClose:(id)sender {
    [globalWindow orderOut:nil]; // 隐藏窗口
    [NSApp hide:nil];            // 隐藏应用
    return NO;                   // 阻止真正退出
}
@end

// --- 设置窗口样式 ---
void setupWindow(void* window) {
    NSWindow *nswindow = (__bridge NSWindow *)window;
    globalWindow = nswindow;

    // 不改 class，保留原生 NSWindow

    // 隐藏标题栏，只保留按钮
    [globalWindow setTitleVisibility:NSWindowTitleHidden];
    [globalWindow setTitlebarAppearsTransparent:YES];
    [globalWindow setStyleMask:[globalWindow styleMask] | NSWindowStyleMaskFullSizeContentView];

    // 设置代理，拦截关闭按钮
    [globalWindow setDelegate:(id<NSWindowDelegate>)[[WindowDelegate alloc] init]];
}

// --- 设置 App 被激活时恢复窗口 ---
void setupAppDelegate() {
    [NSApp setActivationPolicy:NSApplicationActivationPolicyRegular];
    [[NSNotificationCenter defaultCenter] addObserverForName:NSApplicationDidBecomeActiveNotification
                                                      object:nil
                                                       queue:nil
                                                  usingBlock:^(NSNotification *note) {
        if (globalWindow) {
            [globalWindow makeKeyAndOrderFront:nil];
        }
    }];
}

// --- 设置默认菜单，包括退出项 ---
// --- 转语言字符串 ---
NSString *localizedString(const char *clang, NSString *en, NSString *zh) {
    if (clang == NULL) {
        return en;
    }
    NSString *language = [NSString stringWithUTF8String:clang];
    if ([language.lowercaseString hasPrefix:@"zh"]) {
        return zh;
    }
    return en;
}

// --- 刷新菜单，并指定语言 ---
void refreshMenu(const char *lang) {
    NSMenu *mainMenu = [[NSMenu alloc] init];
    [NSApp setMainMenu:mainMenu];

    NSMenuItem *appMenuItem = [[NSMenuItem alloc] init];
    [mainMenu addItem:appMenuItem];

    NSMenu *appMenu = [[NSMenu alloc] init];

    NSString *appName = [[NSProcessInfo processInfo] processName];

    // 关于
    NSString *aboutTitle = localizedString(lang, @"About", @"关于");
    NSMenuItem *aboutItem = [[NSMenuItem alloc] initWithTitle:aboutTitle
                                                       action:@selector(orderFrontStandardAboutPanel:)
                                                keyEquivalent:@""];
    [appMenu addItem:aboutItem];

    [appMenu addItem:[NSMenuItem separatorItem]];

    // 退出
    NSString *quitTitle = localizedString(lang, @"Quit", @"退出");
    NSMenuItem *quitItem = [[NSMenuItem alloc] initWithTitle:quitTitle
                                                      action:@selector(terminate:)
                                               keyEquivalent:@"q"];
    [appMenu addItem:quitItem];

    [appMenuItem setSubmenu:appMenu];
}


// --- 支持 Cmd+C/V/X/A 快捷键 ---
void setupKeyMonitor() {
    [NSEvent addLocalMonitorForEventsMatchingMask:NSEventMaskKeyDown handler:^NSEvent*(NSEvent *event) {
        if ((event.modifierFlags & NSEventModifierFlagCommand) == NSEventModifierFlagCommand) {
            NSString *characters = [event charactersIgnoringModifiers];
            id responder = [globalWindow firstResponder];
            if ([characters isEqualToString:@"c"] && [responder respondsToSelector:@selector(copy:)]) {
                [responder performSelector:@selector(copy:) withObject:nil];
                return nil;
            } else if ([characters isEqualToString:@"v"] && [responder respondsToSelector:@selector(paste:)]) {
                [responder performSelector:@selector(paste:) withObject:nil];
                return nil;
            } else if ([characters isEqualToString:@"x"] && [responder respondsToSelector:@selector(cut:)]) {
                [responder performSelector:@selector(cut:) withObject:nil];
                return nil;
            } else if ([characters isEqualToString:@"a"] && [responder respondsToSelector:@selector(selectAll:)]) {
                [responder performSelector:@selector(selectAll:) withObject:nil];
                return nil;
            }
        }
        return event; // 正常传递其他按键
    }];
}

// --- 支持区域拖拽窗口 ---
void beginDrag(void* window) {
    NSWindow *nswindow = (__bridge NSWindow *)window;
    [nswindow makeKeyAndOrderFront:nil]; // 确保窗口已激活

    NSEvent *event = [NSApp currentEvent];
    if (event.type == NSEventTypeLeftMouseDown) {
        [nswindow performWindowDragWithEvent:event];
    }
}

// --- 设置窗口最小尺寸 ---
void setMinWindowSize(void* window, int minWidth, int minHeight) {
NSWindow *nswindow = (__bridge NSWindow *)window;
[nswindow setContentMinSize:NSMakeSize(minWidth, minHeight)];
}

void activateSelf() {
    NSRunningApplication *app = [NSRunningApplication currentApplication];
    if (![app isActive]) {
        [app activateWithOptions:(NSApplicationActivateIgnoringOtherApps | NSApplicationActivateAllWindows)];
    }

    dispatch_async(dispatch_get_main_queue(), ^{
        NSWindow *keyWindow = [NSApp keyWindow];
        if (!keyWindow) {
            keyWindow = [NSApp mainWindow];
        }

        if (!keyWindow) {
            NSArray<NSWindow *> *windows = [NSApp windows];
            for (NSWindow *win in windows) {
                if (win.isVisible && !win.isMiniaturized) {
                    keyWindow = win;
                    break;
                }
            }
        }

        if (keyWindow) {
            [keyWindow makeKeyAndOrderFront:nil];
            [keyWindow orderFrontRegardless];
        }
    });
}

// --- 隐藏窗口 ---
void hideWindow() {
	[globalWindow orderOut:nil];
	[NSApp hide:nil];
	[NSApp setActivationPolicy:NSApplicationActivationPolicyAccessory];
}

// --- 显示窗口 ---
void showWindow() {
	[NSApp setActivationPolicy:NSApplicationActivationPolicyRegular];
    NSRunningApplication *app = [NSRunningApplication currentApplication];

    if ([NSApp isHidden]) {
        [NSApp unhide:nil];
    }

    if (![app isActive]) {
        if (globalWindow) {
            [globalWindow orderOut:nil]; // 隐藏一下避免残留
            [globalWindow makeKeyAndOrderFront:nil];
        }
        [app activateWithOptions:(NSApplicationActivateIgnoringOtherApps | NSApplicationActivateAllWindows)];
    } else {
        if (globalWindow && ![globalWindow isVisible]) {
            [globalWindow makeKeyAndOrderFront:nil];
        }
    }
}

const char* getClipboardText() {
    @autoreleasepool {
        NSPasteboard *pasteboard = [NSPasteboard generalPasteboard];
        NSString *contents = [pasteboard stringForType:NSPasteboardTypeString];
        if (contents == nil) {
            return NULL;
        }
        const char* cstr = [contents UTF8String];
        return strdup(cstr); // 复制一份，Go负责释放
    }
}

*/
import "C"
import (
	_ "embed"
	"errors"
	"fmt"
	sys "github.com/snakem982/pandora-box/pkg/sys/admin"
	webview "github.com/webview/webview_go"
	"os"
	"os/exec"
	"pandora-box/static"
	"strings"
	"time"
	"unsafe"
)

//go:embed icon-128.png
var Icon []byte

func Init(w webview.WebView) {
	w.SetTitle("Pandora-Box")
	w.SetSize(1100, 760, webview.HintNone)
	C.setupWindow(w.Window())
	C.setupAppDelegate()
	C.setupKeyMonitor() // 加载快捷键监听
	RefreshMenu("zh")
	SetMinWindowSize(w, 960, 660)
	SetDockIconBytes(w, Icon)
	// 绑定前端调用拖拽
	_ = w.Bind("pxDrag", func() {
		C.beginDrag(w.Window())
	})
	// 剪贴板
	_ = w.Bind("pxClipboard", GetClipboard)
	// 打开地址
	_ = w.Bind("pxOpen", Open)
	// 隐藏窗口
	_ = w.Bind("pxHide", HideWindow)

	go func() {
		time.Sleep(1 * time.Second)
		C.activateSelf()
	}()
}

// SetDockIconBytes 用 []byte 设置 Dock 图标
func SetDockIconBytes(_ webview.WebView, icon []byte) {
	if len(icon) == 0 {
		return
	}
	ptr := unsafe.Pointer(&icon[0])
	C.setDockIconFromBytes(ptr, C.int(len(icon)))
}

func SetMinWindowSize(w webview.WebView, width, height int) {
	if w.Window() != nil {
		C.setMinWindowSize(w.Window(), C.int(width), C.int(height))
	}
}

func RefreshMenu(lang string) {
	clang := C.CString(lang)
	defer C.free(unsafe.Pointer(clang))

	C.refreshMenu(clang)
}

// ShowWindow 显示窗口
func ShowWindow() {
	C.showWindow()
}

// HideWindow 隐藏窗口
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
	return exec.Command("open", url).Start()
}

// RunAsAdmin 在后台运行管理员命令
func RunAsAdmin(prompt string, cmd string, args ...string) error {
	// 已经是管理员直接返回
	if sys.IsAdmin() {
		return errors.New("already admin")
	}

	// 构建完整命令
	fullCmd := fmt.Sprintf("%s %s", cmd, strings.Join(args, " "))
	escapedCmd := strings.ReplaceAll(fullCmd, `"`, `\"`)
	escapedPrompt := strings.ReplaceAll(prompt, `"`, `\"`)

	// 构建 AppleScript
	script := fmt.Sprintf(`do shell script "%s  &> /dev/null &" with prompt "%s" with administrator privileges`, escapedCmd, escapedPrompt)
	return exec.Command("osascript", "-e", script).Run()
}

// RunAsNoAdmin 在后台运行命令
func RunAsNoAdmin(cmd string, args ...string) error {

	// 构建完整命令
	fullCmd := fmt.Sprintf("%s %s", cmd, strings.Join(args, " "))
	escapedCmd := strings.ReplaceAll(fullCmd, `"`, `\"`)

	// 构建命令
	script := fmt.Sprintf(`%s  &> /dev/null &`, escapedCmd)
	return exec.Command("sh", "-c", script).Run()
}

func TryAdmin() string {
	server := static.StartFileServer()

	// 尝试提权
	exePath, _ := os.Executable()
	tip := "Pandora-Box 需要授权才能使用 TUN 模式。[Pandora-Box requires authorization to enable TUN.]"
	err := RunAsAdmin(tip, exePath,
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
