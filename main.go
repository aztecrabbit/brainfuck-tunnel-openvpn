package main

import (
	"os"
	"fmt"
	"time"
	"strings"

	"github.com/aztecrabbit/liblog"
	"github.com/aztecrabbit/libutils"
	"github.com/aztecrabbit/libinject"
	"github.com/aztecrabbit/brainfuck-tunnel-openvpn/src/libopenvpn"
)

const (
	appName = "Brainfuck Tunnel"
	appVersionName = "Openvpn"
	appVersionCode = "200122"

	copyrightYear = "2020"
	copyrightAuthor = "Aztec Rabbit"
)

type Config struct {
	Inject *libinject.Config
	Openvpn *libopenvpn.Config
}

func init() {
	InterruptHandler := &libutils.InterruptHandler{
		Handle: func() {
			liblog.LogKeyboardInterrupt()
		},
	}
	InterruptHandler.Start()
}

func main() {
	liblog.Header(
		[]string{
			fmt.Sprintf("%s [%s Version. %s]", appName, appVersionName, appVersionCode),
			fmt.Sprintf("(c) %s %s.", copyrightYear, copyrightAuthor),
		},
		liblog.Colors["G1"],
	)

	config := new(Config)
	configDefault := new(Config)
	configDefault.Inject = libinject.DefaultConfig
	configDefault.Openvpn = libopenvpn.DefaultConfig

	libutils.JsonReadWrite(libutils.RealPath("config.json"), config, configDefault)

	Inject := new(libinject.Inject)
	Inject.Config = config.Inject

	if len(os.Args) > 1 {
		Inject.Config.Port = os.Args[1]
	}

	go Inject.Start()

	time.Sleep(200 * time.Millisecond)

	liblog.LogInfo("Inject running on port " + Inject.Config.Port, "INFO", liblog.Colors["G1"])
	liblog.LogInfo("Openvpn started", "INFO", liblog.Colors["G1"])

	Openvpn := new(libopenvpn.Openvpn)
	Openvpn.Config = config.Openvpn
	for proxyHostPort, _ := range Inject.Config.Proxies {
		Openvpn.ProxyHost = strings.Split(proxyHostPort, ":")[0]
		break
	}
	Openvpn.InjectPort = Inject.Config.Port
	Openvpn.Start()
}
