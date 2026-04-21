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
	"go.minekube.com/gate/pkg/edition/java/proto/packet"
	"go.minekube.com/gate/pkg/edition/java/proto/packet/chat"
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
		player := e.Player()

		server := player.CurrentServer().Server()
		serverName := server.ServerInfo().Name()
		playerCount := server.Players().Len()

		// -------------------------
		// Tablist header/footer
		// -------------------------
		headerRaw := config.GetString("tablist.header")
		headerRaw = replace(headerRaw, "{server}", serverName)
		headerRaw = replace(headerRaw, "{players}", strconv.Itoa(playerCount))

		header := &c.Text{Content: headerRaw}
		footer := &c.Text{Content: config.GetString("tablist.footer")}

		if err := player.TabList().SetHeaderFooter(header, footer); err != nil {
			log.Error(err, "failed to set tab list header/footer")
		}

		// -------------------------
		// Server links
		// -------------------------
		type ServerLinkConfig struct {
			Label any    `mapstructure:"label"`
			URL   string `mapstructure:"url"`
		}

		var linksConfig []ServerLinkConfig

		if err := config.UnmarshalKey("serverLinks", &linksConfig); err != nil {
			log.Error(err, "failed to read serverLinks config")
			linksConfig = nil
		}

		if len(linksConfig) > 0 {
			links := make([]*packet.ServerLink, 0, len(linksConfig))

			for _, l := range linksConfig {
				if l.URL == "" {
					continue
				}

				link := &packet.ServerLink{
					ID:  -1, // server links enum type, 0 to 9. -1 for custom links without an ID (-1 gets ignored and the component is encoded instead).
					URL: l.URL,
				}

				switch v := l.Label.(type) {
				case int:
					link.ID = v
				case int64:
					link.ID = int(v)
				case float64:
					link.ID = int(v)
				case string:
					if n, err := strconv.Atoi(v); err == nil {
						link.ID = n
					} else {
						link.DisplayName = *chat.FromComponent(&c.Text{Content: v})
					}
				}

				links = append(links, link)
			}

			player.WritePacket(&packet.ServerLinks{
				ServerLinks: links,
			})
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
