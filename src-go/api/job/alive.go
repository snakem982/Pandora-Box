package job

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/metacubex/mihomo/hub/executor"
	"github.com/metacubex/mihomo/log"
	"github.com/snakem982/pandora-box/pkg/cache"
	"github.com/snakem982/pandora-box/pkg/cron"
	sys "github.com/snakem982/pandora-box/pkg/sys/proxy"
	"github.com/snakem982/pandora-box/pkg/utils"
)

var aliveLock sync.Mutex

func AliveJob(name string, server string) {
	transport := &http.Transport{
		DisableKeepAlives: true,
	}
	client := &http.Client{
		Transport: transport,
		Timeout:   2 * time.Second,
	}

	url := fmt.Sprintf("http://%s/pxAlive", server)
	lastCheck := time.Now().Round(0)
	failCount := 0
	// 平时连续失败 2 次（约 3~6 秒）即退出，避免拖延重启
	const maxFailCount = 2

	cron.AddTask(name, 3*time.Second, func() {
		if !aliveLock.TryLock() {
			return
		}
		defer aliveLock.Unlock()

		// 仅在系统休眠唤醒时给 10 秒恢复缓冲
		now := time.Now().Round(0)
		if now.Sub(lastCheck) > 10*time.Second {
			log.Warnln("检测到系统从休眠中唤醒，进入 10 秒心跳宽限期...")
			transport.CloseIdleConnections()
			time.Sleep(10 * time.Second)
			lastCheck = time.Now().Round(0)
			failCount = 0
			return
		}
		lastCheck = now

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return
		}

		resp, err := client.Do(req)
		if err != nil {
			failCount++
			log.Warnln("[1]心跳检测失败 (%d/%d): %v", failCount, maxFailCount, err)
			if failCount >= maxFailCount {
				log.Warnln("[1]检测到父进程退出，准备退出...")
				Exit(true)
			}
			return
		}
		defer func(Body io.ReadCloser) {
			_ = Body.Close()
		}(resp.Body)

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			failCount++
			log.Warnln("[2]读取心跳响应失败 (%d/%d): %v", failCount, maxFailCount, err)
			if failCount >= maxFailCount {
				log.Warnln("[2]检测到父进程退出，准备退出...")
				Exit(true)
			}
			return
		}

		if string(body) != "alive" {
			failCount++
			log.Warnln("[3]心跳响应异常 (%d/%d), status=%d, body=%s", failCount, maxFailCount, resp.StatusCode, string(body))
			if failCount >= maxFailCount {
				log.Warnln("[3]检测到父进程退出，准备退出...")
				Exit(true)
			}
			return
		}

		failCount = 0
	})
}

func Exit(needExit bool) {
	cache.Close()
	utils.UnlockSingleton()
	executor.Shutdown()
	sys.DisableProxy()
	if needExit {
		// 给异步日志刷盘留出缓冲时间
		time.Sleep(100 * time.Millisecond)
		os.Exit(0)
	}
}
