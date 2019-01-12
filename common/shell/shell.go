package shell

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/spf13/viper"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

func IsCommandAvailable(command string, args ...string) (bool, error) {
	cmd := exec.Command(command, args...)
	if err := cmd.Run(); err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func IsZmapAvailable() (bool, error) {
	return IsCommandAvailable(viper.GetString("ZmapExecPath"), "-h")
}

func ZmapScan(inputFile string, outputFile string, bandwidth string, sourceAddress string) (string, error) {
	var args []string
	args = append(args, fmt.Sprintf("--bandwidth=%s", bandwidth))
	args = append(args, fmt.Sprintf("--output-file=%s", outputFile))
	args = append(args, fmt.Sprintf("--ipv6-target-file=%s", inputFile))
	args = append(args, fmt.Sprintf("--ipv6-source-ip=%s", sourceAddress))
	args = append(args, "--probe-module=icmp6_echoscan")
	cmd := exec.Command(viper.GetString("ZmapExecPath"), args...)
	log.Printf("Zmap command is: %s %s", cmd.Path, cmd.Args)
	if err := RunCommandToStdout(cmd); err != nil {
		return "", err
	} else {
		return "", nil
	}
}

func ZmapScanFromConfig(inputFile string, outputFile string) (string, error) {
	return ZmapScan(inputFile, outputFile, viper.GetString("ZmapBandwidth"), viper.GetString("ZmapSourceAddress"))
}

func RunCommandToStdout(cmd *exec.Cmd) (error) {
	var stdBuffer bytes.Buffer
	mw := io.MultiWriter(os.Stdout, &stdBuffer)
	cmd.Stdout = mw
	cmd.Stderr = mw
	toReturn := cmd.Run()
	return toReturn
}

func PromptForInput(prompt string) (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println(fmt.Sprintf("\n%s\n", prompt))
	text, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	} else {
		return text, nil
	}
}

func AskForApproval(prompt string) (bool, error) {
	resp, err := PromptForInput(prompt)
	if err != nil {
		return false, err
	}
	fmt.Println()
	resp = strings.TrimSpace(strings.ToLower(resp))
	return resp == "y", nil
}

func RequireApproval(prompt string, error string) (error) {
	resp, err := PromptForInput(prompt)
	if err != nil {
		return err
	}
	fmt.Println()
	resp = strings.TrimSpace(strings.ToLower(resp))
	if resp != "y" {
		log.Fatal(error)
	}
	return nil
}
