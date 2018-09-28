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

func ZmapScan(conf *config.Configuration, inputFile string, outputFile string, bandwidth string, port int) (string, error) {
	var args []string
	args = append(args, fmt.Sprint("--bandwidth=%s", bandwidth))
	args = append(args, fmt.Sprintf("--output-file=%s", outputFile))
	args = append(args, fmt.Sprintf("--whitelist-file=%s", inputFile))
	args = append(args, fmt.Sprintf("--target-port=%d", port))
	cmd := exec.Command(conf.ZmapExecPath, args...)
	if err := cmd.Run(); err != nil {
		return "", err
	} else {
		return "", nil
	}
}

func ZmapScanFromConfig(conf *config.Configuration, inputFile string, outputFile string) (string, error) {
	return ZmapScan(conf, inputFile, outputFile, conf.ZmapBandwidth, 80)
}
