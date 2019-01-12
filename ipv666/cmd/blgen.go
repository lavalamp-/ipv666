package cmd

import (
	"github.com/lavalamp-/ipv666/common/app"
	"github.com/spf13/cobra"
	"log"
	"os"
)

func init() {
	var inputPath string
	blgenCmd.PersistentFlags().StringVarP(&inputPath, "input", "i", "", "An input file containing IPv6 network ranges to build a blacklist from.")
	blgenCmd.MarkFlagRequired("input")
}

func validateBlgenCommand(cmd *cobra.Command) {

	inputPath, err := cmd.PersistentFlags().GetString("input")

	if err != nil {
		log.Fatal(err)
	}

	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		log.Fatalf("No file found at path '%s'. Please supply a valid file path.", inputPath)
	}

}

var blgenLongDesc = `
This utility takes a list of IPv6 CIDR ranges from a text file (new-line delimited),
adds them to the current network blacklist, and sets the new blacklist as the one to use
for the 'scan' command.
`

var blgenCmd = &cobra.Command{
	Use:			"blgen",
	Short:			"Generate a scanning blacklist",
	Long:			blgenLongDesc,
	Run: func(cmd *cobra.Command, args []string) {
		validateBlgenCommand(cmd)
		inputPath, _ := cmd.PersistentFlags().GetString("input")
		app.RunBlgen(inputPath)
	},
}
