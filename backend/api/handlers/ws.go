package handlers

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/metacubex/mihomo/log"
	"github.com/snakem982/pandora-box/api/models"
	"github.com/snakem982/pandora-box/pkg/cache"
	"github.com/snakem982/pandora-box/pkg/constant"
)

// 保存排序后的 Profile 文件
func saveProfileOrder(w http.ResponseWriter, r *http.Request) {
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
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)

	// 处理 WebSocket 消息
	for {
		msg, op, err := wsutil.ReadClientData(conn)
		if err != nil {
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
			_ = cache.Put(constant.ProfileOrder, profiles)
		}

		// 回显消息
		err = wsutil.WriteServerMessage(conn, op, []byte("success"))
		if err != nil {
			log.Errorln("Read error:%v", err)
			break
		}
	}
}

// 保存排序后的 WebTest 文件
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
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)

	// 处理 WebSocket 消息
	for {
		msg, op, err := wsutil.ReadClientData(conn)
		if err != nil {
			break
		}

		// 解析消息
		var webTests []models.WebTest
		if err := json.Unmarshal(msg, &webTests); err != nil {
			log.Errorln("Decode error:%v", err)
			break
		}

		// 保存配置文件顺序
		if len(webTests) > 0 {
			_ = cache.Put(constant.WebTestOrder, webTests)
		}

		// 回显消息
		err = wsutil.WriteServerMessage(conn, op, []byte("success"))
		if err != nil {
			log.Errorln("Read error:%v", err)
			break
		}
	}
}
