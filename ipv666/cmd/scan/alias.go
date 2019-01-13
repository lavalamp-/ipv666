package scan

import (
	"github.com/lavalamp-/ipv666/common/app"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strings"
)

var aliasLongDesc = strings.TrimSpace(`
A utility for testing whether or not a network range exhibits traits of an aliased network range. 
Aliased network ranges are ranges in which every host responds to a ping request, thereby making 
it look like the range is full of IPv6 hosts. Pointing this utility at a network range will let 
tell you whether or not that network range is aliased and, if it is, the boundary of the network 
range that is aliased.
`)

var aliasCmd = &cobra.Command{
	Use:				"alias",
	Short:				"Test a network range for aliased characteristics",
	Long:				aliasLongDesc,
	Run: func(cmd *cobra.Command, args []string) {
		app.RunAlias(viper.GetString("ScanTargetNetwork"))
	},
}