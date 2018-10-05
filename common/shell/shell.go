package shell

import (
	"os/exec"
	"github.com/lavalamp-/ipv666/common/config"
	"fmt"
)

func IsCommandAvailable(command string, args ...string) (bool, error) {
	cmd := exec.Command(command, args...)
	if err := cmd.Run(); err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func IsZmapAvailable(conf *config.Configuration) (bool, error) {
	return IsCommandAvailable(conf.ZmapExecPath, "-h")
}

func ZmapScan(conf *config.Configuration, inputFile string, outputFile string, bandwidth string, sourceAddress string) (string, error) {
	var args []string
	args = append(args, fmt.Sprintf("--bandwidth=%s", bandwidth))
	args = append(args, fmt.Sprintf("--output-file=%s", outputFile))
	args = append(args, fmt.Sprintf("--ipv6-target-file=%s", inputFile))
	args = append(args, fmt.Sprintf("--ipv6-source-ip=%s", sourceAddress))
	args = append(args, "--probe-module=icmp6_echoscan")
	cmd := exec.Command(conf.ZmapExecPath, args...)
	fmt.Print(cmd)
	if err := cmd.Run(); err != nil {
		return "", err
	} else {
		return "", nil
	}
}

func ZmapScanFromConfig(conf *config.Configuration, inputFile string, outputFile string) (string, error) {
	// return "", nil
	return ZmapScan(conf, inputFile, outputFile, conf.ZmapBandwidth, conf.ZmapSourceAddress)
}
