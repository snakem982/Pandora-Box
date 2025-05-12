package job

import (
	"github.com/metacubex/mihomo/log"
	"github.com/snakem982/pandora-box/pkg/cron"
	"github.com/snakem982/pandora-box/pkg/utils"
	"os"
	"sync"
	"time"
)

var logLock sync.Mutex

func LogJob(name string) {
	cron.AddTask(name, 15*time.Minute, func() {
		if logLock.TryLock() {
			defer logLock.Unlock()
		} else {
			return
		}

		// 日志路径
		filePath := utils.GetUserHomeDir("logs", name)
		// 获取文件信息
		fileInfo, err := os.Stat(filePath)
		if err != nil {
			log.Infoln("无法获取文件[%s]信息:%v", name, err)
			return
		}

		// 判断文件大小是否超过 5MB
		if fileInfo.Size() > 5*1024*1024 {
			// 清空文件
			err := os.Truncate(filePath, 0)
			if err != nil {
				log.Errorln("清空文件[%s]失败:%v", name, err)
			} else {
				log.Infoln("文件[%s]已清空", name)
			}
		} else {
			log.Infoln("[%s]文件大小未超过 5MB", name)
		}

	})
}
