package main

import (
	"github.com/miners-online/proxy/plugins/ping"
	"go.minekube.com/gate/cmd/gate"
	"go.minekube.com/gate/pkg/edition/java/proxy"
)

func main() {
	proxy.Plugins = append(proxy.Plugins,
		ping.Plugin,
	)
	gate.Execute()
}
