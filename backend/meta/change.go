package meta

import (
	"encoding/json"
	"github.com/metacubex/mihomo/log"
	"pandora-box/backend/cache"
	"pandora-box/backend/constant"
	"pandora-box/backend/resolve"
)

func SwitchProfile(reload bool) {
	defer func() {
		if e := recover(); e != nil {
			log.Errorln("meta.SwitchProfile field:", e)
		}
	}()

	values := cache.GetList(constant.PrefixProfile)
	for _, value := range values {
		profile := resolve.Profile{}
		_ = json.Unmarshal(value, &profile)
		if profile.Selected {
			StartCore(profile, reload)
			return
		}
	}

	profile := resolve.Profile{
		Id:       constant.DefaultProfile,
		Type:     1,
		Order:    0,
		Path:     "uploads/" + constant.DefaultProfile + ".yaml",
		Selected: true,
	}
	bytes, _ := json.Marshal(profile)
	_ = cache.Put(profile.Id, bytes)
	StartCore(profile, reload)
}
