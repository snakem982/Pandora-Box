//go:build windows

package update

import (
	"pandora-box/backend/constant"
	"pandora-box/backend/tools"
)

func IsNeedUpdate() (bool, string) {
	all, _ := tools.ConcurrentHttpGet(constant.PandoraVersionUrl, nil)
	if all != nil && len(all) > 0 {
		ver := string(all)
		return ver != constant.PandoraVersion, ver
	}

	return false, constant.PandoraVersion
}

func Replace() bool {
	return false
}
