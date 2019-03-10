package generate

import (
	"github.com/spf13/cobra"
	"strings"
)

func init() {
	Cmd.AddCommand(blgenCmd)
	Cmd.AddCommand(modelgenCmd)
	Cmd.AddCommand(addrgenCmd)
}

var generateLongDesc = strings.TrimSpace(`
The generation utilities of IPv666 include (1) generating a network range blacklist, 
(2) generating a predictive clustering model, and (3) generating IPv6 addresses.
`)

var Cmd = &cobra.Command{
	Use:		"generate",
	Short:		"Perform generation functions",
	Long:		generateLongDesc,
}