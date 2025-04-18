package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/metacubex/mihomo/log"
	"github.com/snakem982/pandora-box/pandora/api/models"
	"github.com/snakem982/pandora-box/pandora/pkg/cache"
	"github.com/snakem982/pandora-box/pandora/pkg/constant"
)

// 保存排序后的配置文件
func saveProfilesOrder(w http.ResponseWriter, r *http.Request) {
	// 必须是 websocket 请求
	if !(r.Header.Get("Upgrade") == "websocket") {
		ErrorResponse(w, r, fmt.Errorf("must be a websocket connection"))
		return
	}

	// 升级为 WebSocket 连接
	conn, _, _, err := ws.UpgradeHTTP(r, w)
	if err != nil {
		return
	}
	defer conn.Close()

	// 处理 WebSocket 消息
	for {
		msg, op, err := wsutil.ReadClientData(conn)
		if err != nil {
			log.Errorln("Read error:%v", err)
			break
		}

		// 解析消息
		var profiles []models.Profile
		if err := json.Unmarshal(msg, &profiles); err != nil {
			log.Errorln("Decode error:%v", err)
			break
		}

		// 保存配置文件顺序
		if len(profiles) > 0 {
			_ = cache.Put(constant.ProfilesOrder, profiles)
		}

		// 回显消息
		err = wsutil.WriteServerMessage(conn, op, []byte("success"))
		if err != nil {
			log.Errorln("Read error:%v", err)
			break
		}
	}
}

// 保存排序后的配置文件
func saveWebTestOrder(w http.ResponseWriter, r *http.Request) {
	// 必须是 websocket 请求
	if !(r.Header.Get("Upgrade") == "websocket") {
		ErrorResponse(w, r, fmt.Errorf("must be a websocket connection"))
		return
	}

	// 升级为 WebSocket 连接
	conn, _, _, err := ws.UpgradeHTTP(r, w)
	if err != nil {
		return
	}
	defer conn.Close()

	// 处理 WebSocket 消息
	for {
		msg, op, err := wsutil.ReadClientData(conn)
		if err != nil {
			log.Errorln("Read error:%v", err)
			break
		}

		// 解析消息
		var profiles []models.Profile
		if err := json.Unmarshal(msg, &profiles); err != nil {
			log.Errorln("Decode error:%v", err)
			break
		}

		// 保存配置文件顺序
		if len(profiles) > 0 {
			_ = cache.Put(constant.WebTestOrder, profiles)
		}

		// 回显消息
		err = wsutil.WriteServerMessage(conn, op, []byte("success"))
		if err != nil {
			log.Errorln("Read error:%v", err)
			break
		}
	}
}
