package job

import (
	"fmt"
	"github.com/metacubex/mihomo/hub/executor"
	"github.com/metacubex/mihomo/log"
	"github.com/snakem982/pandora-box/pkg/cache"
	"github.com/snakem982/pandora-box/pkg/cron"
	sys "github.com/snakem982/pandora-box/pkg/sys/proxy"
	"github.com/snakem982/pandora-box/pkg/utils"
	"io"
	"net/http"
	"os"
	"sync"
	"time"
)

var aliveLock sync.Mutex

func AliveJob(name string, server string) {
	// 创建 HTTP 客户端
	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:      5,                // 最大空闲连接数
			IdleConnTimeout:   30 * time.Second, // 空闲连接超时时间
			DisableKeepAlives: false,            // 启用 Keep-Alive
		},
	}

	cron.AddTask(name, 3*time.Second, func() {
		if aliveLock.TryLock() {
			defer aliveLock.Unlock()
		} else {
			return
		}

		// 请求地址
		url := fmt.Sprintf("http://%s/pxAlive", server)

		// 创建请求
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return
		}

		// 设置 Keep-Alive 头
		req.Header.Set("Connection", "keep-alive")

		// 发送请求
		resp, err := client.Do(req)
		if err != nil {
			log.Infoln("[1]检测到父进程退出，准备退出...")
			Exit(true)
			return
		}
		defer func(Body io.ReadCloser) {
			_ = Body.Close()
		}(resp.Body)

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Infoln("[2]检测到父进程退出，准备退出...")
			Exit(true)
			return
		}

		if string(body) != "alive" {
			log.Infoln("[3]检测到父进程退出，准备退出...")
			Exit(true)
		}
	})
}

func Exit(needExit bool) {
	cache.Close()
	utils.UnlockSingleton()
	executor.Shutdown()
	sys.DisableProxy()
	if needExit {
		os.Exit(0)
	}
}
