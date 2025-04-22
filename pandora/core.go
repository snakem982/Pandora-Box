package pandora

import (
	"github.com/metacubex/mihomo/hub/route"
	"github.com/metacubex/mihomo/log"
	"github.com/snakem982/pandora-box/pandora/api/handlers"
	"github.com/snakem982/pandora-box/pandora/internal"
)

func StartCore() (port int, secret string) {
	internal.Init()
	log.Errorln("Initialization completed")

	route.Register(handlers.Profile)
	route.Register(handlers.WebTest)
	route.Register(handlers.Rule)
	route.Register(handlers.DNS)
	route.Register(handlers.Mihomo)
	route.Register(handlers.Pandora)

	port = 9686
	secret = "Y8IUaPeFLTRvsrdf2mUJkLMBuphVZRE5"
	cors := route.Cors{AllowOrigins: []string{"*"}, AllowPrivateNetwork: true}
	addr := route.StartByPandoraBox("127.0.0.1", port, secret, cors)
	log.Errorln("Routing startup completed, Address: %s", addr)

	// 开启mihomo
	internal.SwitchProfile()

	return port, secret
}
