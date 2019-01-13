package shell

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/lavalamp-/ipv666/common/logging"
	"io"
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

func RequireApproval(prompt string, errString string) error {
	resp, err := PromptForInput(prompt)
	if err != nil {
		return err
	}
	fmt.Println()
	resp = strings.TrimSpace(strings.ToLower(resp))
	if resp != "y" {
		logging.ErrorStringF(errString)
	}
	return nil
}
