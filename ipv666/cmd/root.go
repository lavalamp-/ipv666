package cmd

import (
	"bufio"
	"fmt"
	"github.com/lavalamp-/ipv666/internal/config"
	"github.com/lavalamp-/ipv666/internal/logging"
	"github.com/lavalamp-/ipv666/internal/shell"
	"github.com/lavalamp-/ipv666/internal/validation"
	"github.com/lavalamp-/ipv666/ipv666/cmd/generate"
	"github.com/lavalamp-/ipv666/ipv666/cmd/scan"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

func init() {
	var logLevel string
	var forceAccept bool
	rootCmd.PersistentFlags().StringVarP(&logLevel, "log", "l", viper.GetString("LogLevel"), "The log level to emit logs at (one of debug, info, success, warn, error).")
	rootCmd.PersistentFlags().BoolVarP(&forceAccept, "force", "f", viper.GetBool("ForceAcceptPrompts"), "Whether or not to force accept all prompts (useful for daemonized scanning).")
	viper.BindPFlag("LogLevel", rootCmd.PersistentFlags().Lookup("log"))
	viper.BindPFlag("ForceAcceptPrompts", rootCmd.PersistentFlags().Lookup("force"))

	rootCmd.AddCommand(cleanCmd)
	rootCmd.AddCommand(convertCmd)
	rootCmd.AddCommand(scan.Cmd)
	rootCmd.AddCommand(generate.Cmd)
}

func cloudSyncOptIn() error {

	// Get the cloud sync opt-in timestamp path
	path := config.GetCloudSyncOptInPath()

	// Populate the file if it doesn't exist
	if _, err := os.Stat(path); os.IsNotExist(err) {
		f, err := os.Create(path)
		if err != nil {
			return err
		}
		f.WriteString("0\n0\n")
		defer f.Close()
	}

	// Read the opt-in status
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	r := bufio.NewReader(f)
	optInLine, _, err := r.ReadLine()
	if err == io.EOF {
		optInLine = []byte("0")
	} else if err != nil {
		return err
	}
	lastAskLine, _, err := r.ReadLine()
	if err == io.EOF {
		lastAskLine = []byte("0")
	} else if err != nil {
		return err
	}

	// Convert the status to ints
	optIn, err := strconv.ParseInt(string(optInLine), 10, 64)
	if err != nil {
		return err
	}
	lastAsk, err := strconv.ParseInt(string(lastAskLine), 10, 64)
	if err != nil {
		return err
	}

	// Check if already opted in
	if optIn == 1 {
		config.SetCloudSyncOptIn(true)
		return nil
	}

	// Check if the user isn't opted in to cloud sync, and it's been 7+ days since the last nag
	now := time.Now().Unix()
	if optIn == 0 && now - lastAsk > 604800 /* 7 days in seconds */ {

		// Prompt to opt-in to cloud sync
		ok, err := shell.AskForApproval("Would you like to give back to the community and contribute to the cloud-sourced IPv6 dataset @ ipv6.exposed? [y/N]:")
		if err != nil {
			return err
		}

		// Persist the state to disk
		ff, err := os.Create(path)
		if err != nil {
			return err
		}
		defer ff.Close()
		if ok {
			config.SetCloudSyncOptIn(true)
			fmt.Fprintf(ff, "%d\n%d\n", 1, now)
		} else {
			config.SetCloudSyncOptIn(false)
			fmt.Fprintf(ff, "%d\n%d\n", 0, now)
		}
	}
	return nil
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
	err := cloudSyncOptIn()
	if err != nil {
		logging.ErrorF(err)
	}
	rootCmd.Execute()
}
