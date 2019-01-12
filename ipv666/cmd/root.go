package cmd

import (
	"github.com/spf13/cobra"
	"log"
)

func init() {
	rootCmd.AddCommand(scanCmd)
	rootCmd.AddCommand(aliasCmd)
	rootCmd.AddCommand(blgenCmd)
	rootCmd.AddCommand(cleanCmd)
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
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
