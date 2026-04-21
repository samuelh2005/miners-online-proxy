package tablist

import (
	"context"

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

		event.Subscribe(p.Event(), 0, onPing())

		return nil
	},
}

func onPing() func(*proxy.PostLoginEvent) {
	return func(e *proxy.PostLoginEvent) {

		header := &c.Text{
			Content: "Header",
		}
		footer := &c.Text{
			Content: "Footer",
		}

		e.Player().TabList().SetHeaderFooter(header, footer)
	}
}
