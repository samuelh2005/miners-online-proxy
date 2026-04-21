package tablist

import (
	"context"
	"strconv"

	"github.com/go-logr/logr"
	"github.com/robinbraemer/event"
	c "go.minekube.com/common/minecraft/component"
	"go.minekube.com/gate/pkg/edition/java/proxy"
)

// Plugin is a ping plugin that handles ping events.
var Plugin = proxy.Plugin{
	Name: "TabList",
	Init: func(ctx context.Context, p *proxy.Proxy) error {
		log := logr.FromContextOrDiscard(ctx)
		log.Info("Hello from TabList plugin!")

		event.Subscribe(p.Event(), 0, onConnect(log))

		return nil
	},
}

func onConnect(log logr.Logger) func(*proxy.ServerPostConnectEvent) {
	return func(e *proxy.ServerPostConnectEvent) {
		server := e.Player().CurrentServer().Server()
		serverName := server.ServerInfo().Name()
		playerCount := server.Players().Len()

		header := &c.Text{
			Content: "§e⛏ §9lMiners Online §r§e⛏\n§7You are playing on §a" + serverName + " §7with §a" + strconv.Itoa(playerCount) + " §7players online!",
		}
		footer := &c.Text{
			Content: "§eWebsite: §bwww.minersonline.uk §7| §eDiscord: §bdiscord.gg/aeRReEaNnm",
		}

		setError := e.Player().TabList().SetHeaderFooter(header, footer)
		if setError != nil {
			log.Error(setError, "Failed to set tab list header/footer")
		}
	}
}
