package generate

import (
	"github.com/lavalamp-/ipv666/internal/app"
	"github.com/lavalamp-/ipv666/internal/logging"
	"github.com/spf13/cobra"
	"net"
	"os"
	"strings"
)

func init() {
	var modelPath string
	var outputPath string
	var genNetwork string
	var genCount int
	addrgenCmd.PersistentFlags().StringVarP(&modelPath, "model", "m", "", "Local file path to the model to generate addresses from (if empty, uses the default model packaged with ipv666).")
	addrgenCmd.PersistentFlags().StringVarP(&outputPath, "out", "o", "", "File path to where the generated IP addresses should be written.")
	addrgenCmd.PersistentFlags().StringVarP(&genNetwork, "network", "n", "", "The address range to generate addresses within (if empty, generates addresses in the global address space of ::/0).")
	addrgenCmd.PersistentFlags().IntVarP(&genCount, "count", "c", 1000000, "The number of IP addresses to generate.")
	addrgenCmd.MarkPersistentFlagRequired("out")
}

var addrgenLongDesc = strings.TrimSpace(`
This utility will generate IPv6 addresses in target network range (or in the global address space) based on
the default included cluster model or a cluster model that you specify.
`)

var addrgenCmd = &cobra.Command{
	Use:			"addresses",
	Short:			"Generate IPv6 addresses",
	Long:			addrgenLongDesc,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {

		modelPath, err := cmd.PersistentFlags().GetString("model")

		if err != nil {
			logging.ErrorF(err)
		}

		if modelPath != "" {
			if _, err := os.Stat(modelPath); os.IsNotExist(err) {
				logging.ErrorStringFf("No file found at path '%s'. Please supply a valid model path.", modelPath)
			}
		}

		outputPath, err := cmd.PersistentFlags().GetString("out")

		if err != nil {
			logging.ErrorF(err)
		}

		if _, err := os.Stat(outputPath); !os.IsNotExist(err) {
			logging.ErrorStringFf("A file already exists at the path '%s'. Please choose a different file path.", outputPath)
		}

		networkString, err := cmd.PersistentFlags().GetString("network")

		if err != nil {
			logging.ErrorF(err)
		}

		if networkString != "" {
			_, _, err := net.ParseCIDR(networkString)
			if err != nil {
				logging.ErrorF(err)
			}
		}

		genCount, err := cmd.PersistentFlags().GetInt("count")

		if err != nil {
			logging.ErrorF(err)
		}

		if genCount <= 0 {
			logging.ErrorStringFf("You must supply a generate count of greater than zero (got %d).", genCount)
		}

	},
	Run: func(cmd *cobra.Command, args []string) {
		modelPath, _ := cmd.PersistentFlags().GetString("model")
		outputPath, _ := cmd.PersistentFlags().GetString("out")
		networkString, _ := cmd.PersistentFlags().GetString("network")
		genCount, _ := cmd.PersistentFlags().GetInt("count")
		app.RunAddrGen(modelPath, outputPath, networkString, genCount)
	},
}
