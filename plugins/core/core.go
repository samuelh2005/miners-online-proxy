package core

import (
	"context"
	"strconv"
	"sync/atomic"

	"github.com/fsnotify/fsnotify"
	"github.com/go-logr/logr"
	"github.com/robinbraemer/event"
	"github.com/spf13/viper"
	c "go.minekube.com/common/minecraft/component"
	"go.minekube.com/gate/pkg/edition/java/proxy"
	"go.minekube.com/gate/pkg/util/favicon"
)

var index uint64
var config = viper.New()

var Plugin = proxy.Plugin{
	Name: "Core",
	Init: func(ctx context.Context, p *proxy.Proxy) error {
		log := logr.FromContextOrDiscard(ctx)

		config.SetConfigName("branding")
		config.SetConfigType("yaml")
		config.AddConfigPath(".")
		config.AddConfigPath("./config")
		config.AddConfigPath("./plugins/core")

		if err := config.ReadInConfig(); err != nil {
			return err
		}

		config.WatchConfig()
		config.OnConfigChange(func(e fsnotify.Event) {
			log.Info("branding.yml reloaded", "file", e.Name)
		})

		log.Info("branding.yml loaded", "file", config.ConfigFileUsed())

		event.Subscribe(p.Event(), 0, onConnect(log))
		event.Subscribe(p.Event(), 0, onPing())

		return nil
	},
}

func onConnect(log logr.Logger) func(*proxy.ServerPostConnectEvent) {
	return func(e *proxy.ServerPostConnectEvent) {
		server := e.Player().CurrentServer().Server()
		serverName := server.ServerInfo().Name()
		playerCount := server.Players().Len()

		out := config.GetString("tablist.header")
		out = replace(out, "{server}", serverName)
		out = replace(out, "{players}", strconv.Itoa(playerCount))

		header := &c.Text{
			Content: out,
		}
		footer := &c.Text{
			Content: config.GetString("tablist.footer"),
		}

		setError := e.Player().TabList().SetHeaderFooter(header, footer)
		if setError != nil {
			log.Error(setError, "Failed to set tab list header/footer")
		}
	}
}

func onPing() func(*proxy.PingEvent) {
	return func(e *proxy.PingEvent) {
		line1 := config.GetString("ping.header")
		lines := config.GetStringSlice("ping.lines")

		i := atomic.AddUint64(&index, 1)
		line2 := lines[i%uint64(len(lines))]

		p := e.Ping()
		p.Description = &c.Text{
			Content: line1 + "\n" + line2,
		}

		max := config.GetInt("ping.maxPlayers.value")
		inflateOnline := config.GetBool("ping.maxPlayers.inflateOnline")
		if inflateOnline && p.Players.Online >= max {
			max = p.Players.Online + 1
		}
		p.Players.Max = max

		icon := config.GetString("ping.favicon")
		if icon != "" {
			p.Favicon = favicon.Favicon(icon)
		}
	}
}

func replace(s, old, new string) string {
	for {
		i := len(s)
		s = stringReplaceOnce(s, old, new)
		if len(s) == i {
			break
		}
	}
	return s
}

func stringReplaceOnce(s, old, new string) string {
	for i := 0; i+len(old) <= len(s); i++ {
		if s[i:i+len(old)] == old {
			return s[:i] + new + s[i+len(old):]
		}
	}
	return s
}
