package generate

import (
	"github.com/lavalamp-/ipv666/internal/app"
	"github.com/lavalamp-/ipv666/internal/logging"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"strings"
)

func init() {
	var inputPath string
	var outputPath string
	var densityThreshold float64 //TODO add functionality for providing model and pick sizes to command line
	modelgenCmd.PersistentFlags().StringVarP(&inputPath, "input", "i", "", "An input file containing IPv6 addresses to use for the model.")
	modelgenCmd.PersistentFlags().StringVarP(&outputPath, "out", "o", "", "The file path to write the resulting model to.")
	modelgenCmd.PersistentFlags().Float64VarP(&densityThreshold, "density", "d", 0.2, "The minimum density threshold to build a model against.")
	modelgenCmd.MarkPersistentFlagRequired("input") //TODO figure out why persistentflagrequired ain't working
	modelgenCmd.MarkPersistentFlagRequired("out")
}

var modelgenLongDesc = strings.TrimSpace(`
This utility will generate a predictive clustering model based on the contents of  
an IPv6 address file. 
`)

var modelgenCmd = &cobra.Command{
	Use:		"modelgen",
	Short:		"Generate a predictive model",
	Long:		modelgenLongDesc,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {

		inputPath, err := cmd.PersistentFlags().GetString("input")

		if inputPath == "" {
			logging.ErrorStringFf("No input path supplied (-input or -i)")
		}

		if err != nil {
			logging.ErrorF(err)
		}

		if _, err := os.Stat(inputPath); os.IsNotExist(err) {
			logging.ErrorStringFf("No file found at path '%s'. Please supply a valid file path.", inputPath)
		}

		outputPath, err := cmd.PersistentFlags().GetString("out")

		if _, err := os.Stat(outputPath); !os.IsNotExist(err) {
			logging.ErrorStringFf("A file already exists at the path '%s'. Please choose a different file path.", outputPath)
		}

		density, err := cmd.PersistentFlags().GetFloat64("density")

		if err != nil {
			logging.ErrorF(err)
		}

		if density < 0.0 || density > 1.0 {
			logging.ErrorStringFf("Density must be between 0.0 and 1.0 (with a recommended value of around %f).", viper.GetFloat64("ModelDefaultDensityThreshold"))
		}

	},
	Run: func(cmd *cobra.Command, args []string) {
		inputPath, _ := cmd.PersistentFlags().GetString("input")
		outputPath, _ := cmd.PersistentFlags().GetString("out")
		density, _ := cmd.PersistentFlags().GetFloat64("density")
		app.RunModelgen(inputPath, outputPath, density)
	},
}