package cron

import (
	"github.com/go-co-op/gocron"
	"github.com/metacubex/mihomo/log"
	"sync"
	"time"
)

type Cron struct {
	scheduler *gocron.Scheduler
}

var (
	instance *Cron
	once     sync.Once
)

// GetInstance 获取 Cron 单例
func GetInstance() *Cron {
	once.Do(func() {
		instance = &Cron{
			scheduler: gocron.NewScheduler(time.Local),
		}
	})
	return instance
}

// AddTask 添加任务
func AddTask(name string, interval interface{}, task func()) {
	cron := GetInstance()
	_, err := cron.scheduler.Every(interval).Do(task)
	if err != nil {
		log.Infoln("添加任务 %s 失败: %v", name, err)
		return
	}
	log.Infoln("已成功添加任务: %s", name)
}

// Start 启动调度器
func Start() {
	GetInstance().scheduler.StartAsync()
}

// Stop 停止调度器
func Stop() {
	GetInstance().scheduler.Stop()
}
