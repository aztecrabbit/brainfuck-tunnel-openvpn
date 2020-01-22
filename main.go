package main

import (
	"os"
	"fmt"
	"time"

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
	libutils.PathFile = os.Args[0]

	InterruptHandler := &libutils.InterruptHandler{
		Handle: func() {
			libopenvpn.Stop()
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
	configDefault.Inject = &libinject.Config{
		Port: "8989",
		ProxyHost: "202.152.240.50",
		ProxyPort: "80",
		ProxyPayload: "[raw][crlf]Host: t.co[crlf]Host: [crlf][crlf]",
		ProxyTimeout: 10,
		ShowLog: false,
	}
	configDefault.Openvpn = &libopenvpn.Config{
		FileName: "~/account.ovpn",
		AuthFileName: "~/account.ovpn.auth",
	}

	libutils.JsonReadWrite(libutils.RealPath("config.json"), config, configDefault)

	if len(os.Args) > 1 {
		config.Inject.Port = os.Args[1]
	}

	Inject := new(libinject.Inject)
	Inject.Config = config.Inject

	go Inject.Start()

	time.Sleep(200 * time.Millisecond)

	liblog.LogInfo("Inject running on port " + Inject.Config.Port, "INFO", liblog.Colors["G1"])
	liblog.LogInfo("Openvpn started", "INFO", liblog.Colors["G1"])

	Openvpn := new(libopenvpn.Openvpn)
	Openvpn.Config = config.Openvpn
	Openvpn.ProxyHost = Inject.Config.ProxyHost
	Openvpn.InjectPort = Inject.Config.Port
	Openvpn.Start()
}
