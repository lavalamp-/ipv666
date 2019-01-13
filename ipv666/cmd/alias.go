package cmd

import (
	"github.com/lavalamp-/ipv666/common/app"
	"github.com/lavalamp-/ipv666/common/config"
	"github.com/lavalamp-/ipv666/common/logging"
	"github.com/lavalamp-/ipv666/common/shell"
	"github.com/lavalamp-/ipv666/common/validation"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	var targetNetwork string
	aliasCmd.PersistentFlags().StringVarP(&targetNetwork, "target", "t", viper.GetString("ScanTargetNetwork"), "An IPv6 CIDR range to test as an aliased network.")
	aliasCmd.MarkPersistentFlagRequired("target")
}

//TODO move all validation into PersistentPreRun

func validateAliasCommand(cmd *cobra.Command) {

	targetNetwork, err := cmd.PersistentFlags().GetString("target")

	if err != nil {
		logging.ErrorF(err)
	}

	_, err = validation.ValidateIPv6NetworkStringForScanning(targetNetwork)
	if err != nil {
		logging.ErrorF(err)
	}

	_, err = config.GetTargetNetwork()
	if err != nil {
		logging.ErrorF(err)
	}

	zmapAvailable, err := shell.IsZmapAvailable()

	if err != nil {
		logging.ErrorStringFf("Error thrown when checking for Zmap: %s", err)
	} else if !zmapAvailable {
		logging.ErrorStringF("Zmap not found. Please install Zmap.")
	}

}

var aliasLongDesc = `
A utility for testing whether or not a network range exhibits traits of an aliased network range. 
Aliased network ranges are ranges in which every host responds to a ping request, thereby making 
it look like the range is full of IPv6 hosts. Pointing this utility at a network range will let 
tell you whether or not that network range is aliased and, if it is, the boundary of the network 
range that is aliased.
`

var aliasCmd = &cobra.Command{
	Use:			"alias",
	Short:			"Test a network range for aliased characteristics",
	Long:			aliasLongDesc,
	Run: func(cmd *cobra.Command, args []string) {
		validateAliasCommand(cmd)
		targetNetwork, _ := cmd.PersistentFlags().GetString("target")
		app.RunAlias(targetNetwork)
	},
}