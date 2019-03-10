package generate

import (
	"github.com/lavalamp-/ipv666/internal/app"
	"github.com/lavalamp-/ipv666/internal/logging"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

func init() {
	var inputPath string
	var outputPath string
	modelgenCmd.PersistentFlags().StringVarP(&inputPath, "input", "i", "", "An input file containing IPv6 addresses to use for the model.")
	modelgenCmd.PersistentFlags().StringVarP(&outputPath, "out", "o", "", "The file path to write the resulting model to.")
	modelgenCmd.MarkPersistentFlagRequired("input") //TODO figure out why persistentflagrequired ain't working
	modelgenCmd.MarkPersistentFlagRequired("out")
}

var modelgenLongDesc = strings.TrimSpace(`
This utility will generate a predictive clustering model based on the contents of  
an IPv6 address file. 
`)

var modelgenCmd = &cobra.Command{
	Use:		"model",
	Short:		"Generate a predictive model",
	Long:		modelgenLongDesc,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {

		inputPath, err := cmd.PersistentFlags().GetString("input")

		if err != nil {
			logging.ErrorF(err)
		}

		if inputPath == "" {
			logging.ErrorStringFf("No input path supplied (-input or -i)")
		}

		if _, err := os.Stat(inputPath); os.IsNotExist(err) {
			logging.ErrorStringFf("No file found at path '%s'. Please supply a valid file path.", inputPath)
		}

		outputPath, err := cmd.PersistentFlags().GetString("out")

		if err != nil {
			logging.ErrorF(err)
		}

		if _, err := os.Stat(outputPath); !os.IsNotExist(err) {
			logging.ErrorStringFf("A file already exists at the path '%s'. Please choose a different file path.", outputPath)
		}

	},
	Run: func(cmd *cobra.Command, args []string) {
		inputPath, _ := cmd.PersistentFlags().GetString("input")
		outputPath, _ := cmd.PersistentFlags().GetString("out")
		app.RunModelgen(inputPath, outputPath)
	},
}