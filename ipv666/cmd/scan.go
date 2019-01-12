package cmd

import (
	"fmt"
	"github.com/lavalamp-/ipv666/common/app"
	"github.com/lavalamp-/ipv666/common/config"
	"github.com/lavalamp-/ipv666/common/shell"
	"github.com/lavalamp-/ipv666/common/validation"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
)

func init() {
	var outputFileName string
	var outputFileType string
	var targetNetwork string
	var forceAccept bool
	scanCmd.PersistentFlags().StringVarP(&outputFileName, "output", "o", viper.GetString("OutputFileName"), "The path to the file where discovered addresses should be written.")
	scanCmd.PersistentFlags().StringVarP(&outputFileType, "output-type", "t", viper.GetString("OutputFileType"), "The type of output to write to the output file (txt or bin).")
	scanCmd.PersistentFlags().StringVarP(&targetNetwork, "network", "n", viper.GetString("ScanTargetNetwork"), "The target IPv6 network range to scan.")
	scanCmd.PersistentFlags().BoolVarP(&forceAccept, "force", "f", viper.GetBool("ForceAcceptPrompts"), "Whether or not to force accept all prompts (useful for daemonized scanning).")
	viper.BindPFlag("OutputFileName", scanCmd.PersistentFlags().Lookup("output"))
	viper.BindPFlag("OutputFileType", scanCmd.PersistentFlags().Lookup("output-type"))
	viper.BindPFlag("ScanTargetNetwork", scanCmd.PersistentFlags().Lookup("network"))
	viper.BindPFlag("ForceAcceptPrompts", scanCmd.PersistentFlags().Lookup("force"))
}

func validateScanCommand() {

	err := validation.ValidateOutputFileType(viper.GetString("OutputFileType"))
	if err != nil {
		log.Fatal(err)
	}

	if _, err := os.Stat(config.GetOutputFilePath()); !os.IsNotExist(err) {
		if !viper.GetBool("ForceAcceptPrompts") {
			prompt := fmt.Sprintf("Output file already exists at path '%s,' continue (will append to existing file)? [y/N]", config.GetOutputFilePath())
			errMsg := fmt.Sprintf("Exiting. Please move the file at path '%s' and try again.", config.GetOutputFilePath())
			err := shell.RequireApproval(prompt, errMsg)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			log.Printf("Force accept configured. Not asking for permission to append to file '%s'.", config.GetOutputFilePath())
		}
	}

	_, err = validation.ValidateIPv6NetworkStringForScanning(viper.GetString("ScanTargetNetwork"))
	if err != nil {
		log.Fatal(err)
	}

	_, err = config.GetTargetNetwork()
	if err != nil {
		log.Fatal(err)
	}

	zmapAvailable, err := shell.IsZmapAvailable()

	if err != nil {
		log.Fatal("Error thrown when checking for Zmap: ", err)
	} else if !zmapAvailable {
		log.Fatal("Zmap not found. Please install Zmap.")
	}

}

var scanLongDesc = `
This utility scans for live hosts over IPv6 based on the network range you specify. If no 
range is specified, then this utility scans the global IPv6 address space (e.g. 2000::/4). 
The scanning process generates candidate addresses, scans for them, tests the network ranges 
where live addresses are found for aliased conditions, and adds legitimate discovered IPv6 
addresses to an output list.
`

var scanCmd = &cobra.Command{
	Use:			"scan",
	Short:			"Scan for live hosts on IPv6 addresses",
	Long:			scanLongDesc,
	Run: func(cmd *cobra.Command, args []string) {
		validateScanCommand()
		app.RunScanning()
	},
}