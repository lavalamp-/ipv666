package scan

import (
	"github.com/lavalamp-/ipv666/internal/config"
	"github.com/lavalamp-/ipv666/internal/logging"
	"github.com/lavalamp-/ipv666/internal/validation"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"runtime"
	"strings"
)

func init() {
	var bandwidth string
	var targetNetwork string
	Cmd.PersistentFlags().StringVarP(&bandwidth, "bandwidth", "b", viper.GetString("PingScanBandwidth"), "The maximum bandwidth to use for ping scanning")
	Cmd.PersistentFlags().StringVarP(&targetNetwork, "network", "n", viper.GetString("ScanTargetNetwork"), "The IPv6 CIDR range to scan.")
	viper.BindPFlag("PingScanBandwidth", Cmd.PersistentFlags().Lookup("bandwidth"))
	viper.BindPFlag("ScanTargetNetwork", Cmd.PersistentFlags().Lookup("network"))
	Cmd.AddCommand(discoverCmd)
	Cmd.AddCommand(aliasCmd)
}

var scanLongDesc = strings.TrimSpace(`
The scanning utilities of IPv666 include (1) scanning a target network range (or 
the global IPv6 address space) for live hosts over IPv6 and (2) determining whether 
or not a target network range is an aliased network range.
`)

var Cmd = &cobra.Command{
	Use:			"scan",
	Short:			"Perform scanning functions",
	Long:			scanLongDesc,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {

		bandwidth := viper.GetString("PingScanBandwidth")

		if err := validation.ValidateScanBandwidth(bandwidth); err != nil {
			logging.ErrorF(err)
		}

		targetNetwork := viper.GetString("ScanTargetNetwork")

		if err := validation.ValidateIPv6NetworkString(targetNetwork); err != nil {
			logging.ErrorF(err)
		}

		_, err := config.GetTargetNetwork()
		if err != nil {
			logging.ErrorF(err)
		}

		if runtime.GOOS != "linux" {
			logging.ErrorStringFf("%s is not a supported platform - ipv666's scanning tools only work on Linux systems", runtime.GOOS)
		}

	},
}