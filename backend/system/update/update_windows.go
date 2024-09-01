//go:build windows

package update

import (
	"fmt"
	"os"
	"os/exec"
	"pandora-box/backend/constant"
	"pandora-box/backend/tools"
	"time"
)

func IsNeedUpdate() bool {
	all, _ := tools.ConcurrentHttpGet(constant.PandoraVersionUrl, nil)
	if all != nil && len(all) > 0 {
		return string(all) != constant.PandoraVersion
	}

	return false
}

func Replace() bool {
	err := os.Rename("./cms-test.exe", "./cms-test-temp.exe")
	if err != nil {
		fmt.Print(err)
		return
	}
	time.Sleep(5 * time.Second)
	os.Rename("./cms.exe", "./cms-test.exe")
	time.Sleep(2 * time.Second)
	os.Remove("./cms-test-temp.exe")
	time.Sleep(5 * time.Second)

	// 到程序安装路径下去执行启动命令(预防相对路径方式启动)
	daemon := "timeout /T 3 & E:\\cms\\cms-test.exe 2>&1 &"
	_ = exec.Command("cmd.exe", "/C", daemon).Start()

	return false
}
