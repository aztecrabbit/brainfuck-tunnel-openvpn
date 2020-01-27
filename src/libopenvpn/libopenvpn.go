package libopenvpn

import (
	"os/exec"
	"fmt"
	"bufio"
	"strings"
	"syscall"

	"github.com/aztecrabbit/liblog"
)

var (
	Loop = true
	DefaultConfig = &Config{
		FileName: "~/account.ovpn",
		AuthFileName: "~/account.ovpn.auth",
	}
)

func Stop() {
	Loop = false
}

type Config struct {
	FileName string
	AuthFileName string
}

type Openvpn struct {
	Config *Config
	ProxyHost string
	InjectPort string
}

func (o *Openvpn) LogInfoSplit(message string, slice int, color string) {
	if Loop {
		liblog.LogInfoSplit(message, slice, "INFO", color)
	}
}

func (o *Openvpn) LogInfo(message string, color string) {
	o.LogInfoSplit(message, 0, color)
}

func (o *Openvpn) Start() {
	command := exec.Command(
		"sh", "-c", fmt.Sprintf(
			"openvpn --config %s --auth-user-pass %s " +
				"--route %s 255.255.255.255 net_gateway " +
				"--http-proxy 127.0.0.1 %s",
			o.Config.FileName,
			o.Config.AuthFileName,
			o.ProxyHost,
			o.InjectPort,
		),
	)

	stdout, err := command.StdoutPipe()
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(stdout)
	go func() {
		var line string
		for Loop && scanner.Scan() {
			line = scanner.Text()

			if strings.Contains(line, "Initialization Sequence Completed") {
				o.LogInfo("Connected", liblog.Colors["Y1"])

			} else if strings.Contains(line, "Connection reset") {
				o.LogInfo("Reconnecting", liblog.Colors["G1"])

			} else if strings.Contains(line, "Exiting due to fatal error") {
				o.LogInfo(
					"Fatal Error:\n\n" +
						"|   Please run as root or something like that!\n" +
						"|   I don't know why exacly :D\n" +
						"|\n",
					liblog.Colors["R1"],
				)

			} else {
				o.LogInfoSplit(line[25:], 22, liblog.Colors["G2"])

			}
		}

		command.Process.Signal(syscall.SIGTERM)
	}()

	command.Start()
	command.Wait()
}
