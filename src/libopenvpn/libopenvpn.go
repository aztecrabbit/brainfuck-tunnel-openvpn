package libopenvpn

import (
	"os"
	"os/exec"
	"fmt"
	"bufio"
	"strings"

	"github.com/aztecrabbit/liblog"
	"github.com/aztecrabbit/libutils"
)

func init() {
	libutils.PathFile = os.Args[0]
}

func Stop() {
	command := exec.Command("killall openvpn > /dev/null 2>&1 &")
	command.Start()
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

func (o *Openvpn) Start() {
	command := exec.Command(
		"dash", "-c", fmt.Sprintf(
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
		for scanner.Scan() {
			line = scanner.Text()

			if strings.Contains(line, "Initialization Sequence Completed") {
				liblog.LogInfo("Connected", "INFO", liblog.Colors["Y1"])

			} else if strings.Contains(line, "Exiting due to fatal error") {
				liblog.LogInfo(
					"Fatal Error:\n\n" +
						"|   Please run as root or something like that!\n" +
						"|   I don't know why exacly :D\n" +
						"|\n",
					"INFO", liblog.Colors["R1"])

			} else {
				liblog.LogInfoSplit(line[25:], 22, "INFO", liblog.Colors["G2"])

			}
		}
	}()

	command.Start()
	command.Wait()
}
