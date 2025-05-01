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
func AddTask(interval uint64, task func()) {
	cron := GetInstance()
	job, err := cron.scheduler.Every(interval).Do(task)
	if err != nil {
		log.Infoln("添加任务失败: %v", err)
		return
	}
	log.Infoln("任务已成功添加: %v", job)
}

// Start 启动调度器
func Start() {
	GetInstance().scheduler.StartAsync()
}

// Stop 停止调度器
func Stop() {
	GetInstance().scheduler.Stop()
}
