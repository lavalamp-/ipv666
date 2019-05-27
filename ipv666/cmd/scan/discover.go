package scan

import (
	"fmt"
	"github.com/lavalamp-/ipv666/internal/app"
	"github.com/lavalamp-/ipv666/internal/config"
	"github.com/lavalamp-/ipv666/internal/logging"
	"github.com/lavalamp-/ipv666/internal/shell"
	"github.com/lavalamp-/ipv666/internal/validation"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"strings"
)

func init() {
	var outputFileName string
	var outputFileType string
	discoverCmd.PersistentFlags().StringVarP(&outputFileName, "output", "o", viper.GetString("OutputFileName"), "The path to the file where discovered addresses should be written.")
	discoverCmd.PersistentFlags().StringVarP(&outputFileType, "output-type", "t", viper.GetString("OutputFileType"), "The type of output to write to the output file (txt or bin).")
	viper.BindPFlag("OutputFileName", discoverCmd.PersistentFlags().Lookup("output"))
	viper.BindPFlag("OutputFileType", discoverCmd.PersistentFlags().Lookup("output-type"))
}

var discoverLongDesc = strings.TrimSpace(`
This utility scans for live hosts over IPv6 based on the network range you specify. If no 
range is specified, then this utility scans the global IPv6 address space (e.g. 2000::/4). 
The scanning process generates candidate addresses, scans for them, tests the network ranges 
where live addresses are found for aliased conditions, and adds legitimate discovered IPv6 
addresses to an output list.
`)

var discoverCmd = &cobra.Command{
	Use:			"discover",
	Short:			"Discover live hosts over IPv6",
	Long:			discoverLongDesc,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {

		err := validation.ValidateOutputFileType(viper.GetString("OutputFileType"))
		if err != nil {
			logging.ErrorF(err)
		}

		if _, err := os.Stat(config.GetOutputFilePath()); !os.IsNotExist(err) {
			if !viper.GetBool("ForceAcceptPrompts") {
				prompt := fmt.Sprintf("Output file already exists at path '%s,' continue (will append to existing file)? [y/N]", config.GetOutputFilePath())
				errMsg := fmt.Sprintf("Exiting. Please move the file at path '%s' and try again.", config.GetOutputFilePath())
				err := shell.RequireApproval(prompt, errMsg)
				if err != nil {
					logging.ErrorF(err)
				}
			} else {
				logging.Infof("Force accept configured. Not asking for permission to append to file '%s'.", config.GetOutputFilePath())
			}
		}

	},
	Run: func(cmd *cobra.Command, args []string) {
		app.RunDiscovery()
	},
}