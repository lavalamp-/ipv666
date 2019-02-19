package generate

import (
	"github.com/spf13/cobra"
	"strings"
)

func init() {
	Cmd.AddCommand(blgenCmd)
	Cmd.AddCommand(modelgenCmd)
}

var generateLongDesc = strings.TrimSpace(`
The generation utilities of IPv666 include (1) generating a network range blacklist 
and (2) generating a predictive clustering model.
`)

var Cmd = &cobra.Command{
	Use:		"generate",
	Short:		"Perform generation functions",
	Long:		generateLongDesc,
}