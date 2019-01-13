package cmd

import (
	"github.com/lavalamp-/ipv666/common/validation"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	var logLevel string
	rootCmd.PersistentFlags().StringVarP(&logLevel, "log", "l", viper.GetString("LogLevel"), "The log level to emit logs at (one of debug, info, success, warn, error).")
	viper.BindPFlag("LogLevel", rootCmd.PersistentFlags().Lookup("log"))

	rootCmd.AddCommand(scanCmd)
	rootCmd.AddCommand(aliasCmd)
	rootCmd.AddCommand(blgenCmd)
	rootCmd.AddCommand(cleanCmd)
	rootCmd.AddCommand(convertCmd)
}

var rootLongDesc = `
An IPv6 host enumeration tool set intended for the discovery of hosts within the 
vast IPv6 address space. This tool set includes capabilities for scanning the global 
address space, specific network ranges, converting IPv6 address lists to various 
representations, testing IPv6 network ranges for aliased traits, and cleaning IPv6 
address lists based on known bad networks.
`

var rootCmd = &cobra.Command{
	Use:		"ipv666",
	Short:		"IPv6 address enumeration tool set",
	Long:		rootLongDesc,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		logLevel := viper.GetString("LogLevel")
		return validation.ValidateLogLevel(logLevel)
	},
}

func Execute() {
	rootCmd.Execute()
}
