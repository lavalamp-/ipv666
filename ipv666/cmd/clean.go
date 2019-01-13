package cmd

import (
	"github.com/lavalamp-/ipv666/internal/app"
	"github.com/lavalamp-/ipv666/internal/blacklist"
	"github.com/lavalamp-/ipv666/internal/data"
	"github.com/lavalamp-/ipv666/internal/logging"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

func init() {
	var inputPath string
	var outputPath string
	var blacklistPath string
	cleanCmd.PersistentFlags().StringVarP(&inputPath, "input", "i", "", "An input file containing IPv6 addresses to clean via a blacklist.")
	cleanCmd.PersistentFlags().StringVarP(&outputPath, "out", "o", "", "The file path where the cleaned results should be written to.")
	cleanCmd.PersistentFlags().StringVarP(&blacklistPath, "blacklist", "b", "", "The local file path to the blacklist to use. If not specified, defaults to the most recent blacklist in the configured blacklist directory.")
	cleanCmd.MarkPersistentFlagRequired("input")
	cleanCmd.MarkPersistentFlagRequired("out")
}

var cleanLongDesc = strings.TrimSpace(`
This utility will clean the contents of an IPv6 address file (new-line delimited, 
standard ASCII hex representation) based on the contents of an IPv6 network blacklist
file. If no blacklist path is supplied then the utility will use the default blacklist. 
The cleaned results will then be written to an output file.
`)

var cleanCmd = &cobra.Command{
	Use:			"clean",
	Short:			"Clean a list of IPv6 addresses",
	Long:			cleanLongDesc,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {

		inputPath, err := cmd.PersistentFlags().GetString("input")

		if err != nil {
			logging.ErrorF(err)
		}

		if _, err := os.Stat(inputPath); os.IsNotExist(err) {
			logging.ErrorStringFf("No file found at path '%s'. Please supply a valid file path.", inputPath)
		}

		outPath, err := cmd.PersistentFlags().GetString("out")

		if err != nil {
			logging.ErrorF(err)
		}

		if _, err := os.Stat(outPath); !os.IsNotExist(err) {
			logging.ErrorStringFf("File already exists at output path of '%s'. Please choose a different output path or delete the file in question.", outPath)
		}

		blacklistPath, err := cmd.PersistentFlags().GetString("blacklist")

		if err != nil {
			logging.ErrorF(err)
		}

		if blacklistPath != "" {
			if _, err := os.Stat(blacklistPath); os.IsNotExist(err) {
				logging.ErrorStringFf("No blacklist file found at path '%s'. Please either specify a valid blacklist file path or don't specify one.", blacklistPath)
			}
		}

	},
	Run: func(cmd *cobra.Command, args []string) {
		inputPath, _ := cmd.PersistentFlags().GetString("input")
		outputPath, _ := cmd.PersistentFlags().GetString("out")
		blacklistPath, _ := cmd.PersistentFlags().GetString("blacklist")
		var processBlacklist *blacklist.NetworkBlacklist
		var err error

		if blacklistPath != "" {
			processBlacklist, err = blacklist.ReadNetworkBlacklistFromFile(blacklistPath)
		} else {
			processBlacklist, err = data.GetBlacklist()
		}

		if err != nil {
			logging.ErrorF(err)
		}

		app.RunClean(inputPath, outputPath, processBlacklist)
	},
}
