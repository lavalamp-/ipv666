package cmd

import (
	"github.com/lavalamp-/ipv666/internal/logging"
	"github.com/lavalamp-/ipv666/internal/splash"
	"github.com/lavalamp-/ipv666/internal/validation"
	"github.com/lavalamp-/ipv666/ipv666/cmd/scan"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strings"
)

func init() {
	cobra.OnInitialize(splash.PrintSplash)

	var logLevel string
	var forceAccept bool
	rootCmd.PersistentFlags().StringVarP(&logLevel, "log", "l", viper.GetString("LogLevel"), "The log level to emit logs at (one of debug, info, success, warn, error).")
	rootCmd.PersistentFlags().BoolVarP(&forceAccept, "force", "f", viper.GetBool("ForceAcceptPrompts"), "Whether or not to force accept all prompts (useful for daemonized scanning).")
	viper.BindPFlag("LogLevel", rootCmd.PersistentFlags().Lookup("log"))
	viper.BindPFlag("ForceAcceptPrompts", rootCmd.PersistentFlags().Lookup("force"))

	rootCmd.AddCommand(blgenCmd)
	rootCmd.AddCommand(cleanCmd)
	rootCmd.AddCommand(convertCmd)
	rootCmd.AddCommand(scan.Cmd)
}

var rootLongDesc = strings.TrimSpace(`
An IPv6 host enumeration tool set intended for the discovery of hosts within the 
vast IPv6 address space. This tool set includes capabilities for scanning the global 
address space, specific network ranges, converting IPv6 address lists to various 
representations, testing IPv6 network ranges for aliased traits, and cleaning IPv6 
address lists based on known bad networks.
`)

var rootCmd = &cobra.Command{
	Use:		"ipv666",
	Short:		"IPv6 address enumeration tool set",
	Long:		rootLongDesc,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		logLevel := viper.GetString("LogLevel")
		if err := validation.ValidateLogLevel(logLevel); err != nil {
			logging.ErrorF(err)
		}
	},
}

func Execute() {
	rootCmd.Execute()
}
