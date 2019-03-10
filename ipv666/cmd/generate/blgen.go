package generate

import (
	"github.com/lavalamp-/ipv666/internal/app"
	"github.com/lavalamp-/ipv666/internal/logging"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

//TODO for some reason the PersisteentPreRun is happening before the flag being missing throws an error. fix

func init() {
	var inputPath string
	blgenCmd.PersistentFlags().StringVarP(&inputPath, "input", "i", "", "An input file containing IPv6 network ranges to build a blacklist from.")
	blgenCmd.MarkPersistentFlagRequired("input")
}

var blgenLongDesc = strings.TrimSpace(`
This utility takes a list of IPv6 CIDR ranges from a text file (new-line delimited),
adds them to the current network blacklist, and sets the new blacklist as the one to use
for the 'scan' command.
`)

var blgenCmd = &cobra.Command{
	Use:			"blacklist",
	Short:			"Generate a scanning blacklist",
	Long:			blgenLongDesc,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {

		inputPath, err := cmd.PersistentFlags().GetString("input")

		if err != nil {
			logging.ErrorF(err)
		}

		if _, err := os.Stat(inputPath); os.IsNotExist(err) {
			logging.ErrorStringFf("No file found at path '%s'. Please supply a valid file path.", inputPath)
		}

	},
	Run: func(cmd *cobra.Command, args []string) {
		inputPath, _ := cmd.PersistentFlags().GetString("input")
		app.RunBlgen(inputPath)
	},
}
