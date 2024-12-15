//go:build windows

package update

import (
	"fmt"
	"pandora-box/backend/constant"
	"pandora-box/backend/tools"
)

func IsNeedUpdate() (bool, string) {
	url := fmt.Sprintf("%s?v=%s", constant.PandoraVersionUrl, tools.Dec(8))
	all, _ := tools.ConcurrentHttpGet(url, nil)
	if all != nil && len(all) > 0 {
		ver := string(all)
		return ver != constant.PandoraVersion, ver
	}

	return false, constant.PandoraVersion
}

func Replace() bool {
	return false
}
