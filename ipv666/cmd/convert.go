package cmd

import (
	"github.com/lavalamp-/ipv666/internal/app"
	"github.com/lavalamp-/ipv666/internal/logging"
	"github.com/lavalamp-/ipv666/internal/validation"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strings"
)

func init() {
	var inputPath string
	var outputPath string
	var outputType string
	convertCmd.PersistentFlags().StringVarP(&inputPath, "input", "i", "", "The file to process IPv6 addresses out of.")
	convertCmd.PersistentFlags().StringVarP(&outputPath, "out", "o", "", "The file path to write the converted file to.")
	convertCmd.PersistentFlags().StringVarP(&outputType, "type", "t", viper.GetString("OutputFileType"), "The format to write the IPv6 addresses in (one of 'txt', 'bin', 'hex', 'tree').")
	convertCmd.MarkPersistentFlagRequired("input")
	convertCmd.MarkPersistentFlagRequired("out")
}

var convertLongDesc = strings.TrimSpace(`
This utility will process the contents of a file as containing IPv6 addresses, convert
those addresses to another format, and then write a new file with the same addresses in
the new format. This functionality is (hopefully) intelligent enough to determine how the
addresses are stored in the file without having to specify an input type.
`)

var convertCmd = &cobra.Command{
	Use:			"convert",
	Short:			"Convert a file of IPv6 addresses to another format",
	Long:			convertLongDesc,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {

		inputPath, err := cmd.PersistentFlags().GetString("input")

		if err != nil {
			logging.ErrorF(err)
		}

		if err := validation.ValidateFileExists(inputPath); err != nil {
			logging.ErrorF(err)
		}

		outputPath, err := cmd.PersistentFlags().GetString("out")

		if err != nil {
			logging.ErrorF(err)
		}

		if err := validation.ValidateFileNotExist(outputPath); err != nil {
			logging.ErrorF(err)
		}

		outputType, err := cmd.PersistentFlags().GetString("type")

		if err != nil {
			logging.ErrorF(err)
		}

		if err := validation.ValidateOutputFileType(outputType); err != nil {
			logging.ErrorF(err)
		}

	},
	Run: func(cmd *cobra.Command, args []string) {
		inputPath, _ := cmd.PersistentFlags().GetString("input")
		outputPath, _ := cmd.PersistentFlags().GetString("out")
		outputType, _ := cmd.PersistentFlags().GetString("type")
		app.RunConvert(inputPath, outputPath, outputType)
	},
}
