package main

import (
	"github.com/miners-online/proxy/plugins/ping"
	"github.com/miners-online/proxy/plugins/tablist"
	"go.minekube.com/gate/cmd/gate"
	"go.minekube.com/gate/pkg/edition/java/proxy"
)

func main() {
	proxy.Plugins = append(proxy.Plugins,
		ping.Plugin,
		tablist.Plugin,
	)
	gate.Execute()
}
