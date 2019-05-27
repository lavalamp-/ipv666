package scan

import (
	"bufio"
	"fmt"
	"github.com/lavalamp-/ipv666/internal/app"
	"github.com/lavalamp-/ipv666/internal/config"
	"github.com/lavalamp-/ipv666/internal/logging"
	"github.com/lavalamp-/ipv666/internal/shell"
	"github.com/lavalamp-/ipv666/internal/validation"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

func init() {
	var outputFileName string
	var outputFileType string
	discoverCmd.PersistentFlags().StringVarP(&outputFileName, "output", "o", viper.GetString("OutputFileName"), "The path to the file where discovered addresses should be written.")
	discoverCmd.PersistentFlags().StringVarP(&outputFileType, "output-type", "t", viper.GetString("OutputFileType"), "The type of output to write to the output file (txt or bin).")
	viper.BindPFlag("OutputFileName", discoverCmd.PersistentFlags().Lookup("output"))
	viper.BindPFlag("OutputFileType", discoverCmd.PersistentFlags().Lookup("output-type"))
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
		ok, err := shell.AskForApproval("Would you like to give back to the community and contribute the cloud-sourced IPv6 dataset @ ipv6.exposed? [y/N]:")
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

var discoverLongDesc = strings.TrimSpace(`
This utility scans for live hosts over IPv6 based on the network range you specify. If no 
range is specified, then this utility scans the global IPv6 address space (e.g. 2000::/4). 
The scanning process generates candidate addresses, scans for them, tests the network ranges 
where live addresses are found for aliased conditions, and adds legitimate discovered IPv6 
addresses to an output list.
`)

var discoverCmd = &cobra.Command{
	Use:			"discover",
	Short:			"Discover live hosts over IPv6",
	Long:			discoverLongDesc,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {

		err := validation.ValidateOutputFileType(viper.GetString("OutputFileType"))
		if err != nil {
			logging.ErrorF(err)
		}

		if _, err := os.Stat(config.GetOutputFilePath()); !os.IsNotExist(err) {
			if !viper.GetBool("ForceAcceptPrompts") {
				prompt := fmt.Sprintf("Output file already exists at path '%s,' continue (will append to existing file)? [y/N]", config.GetOutputFilePath())
				errMsg := fmt.Sprintf("Exiting. Please move the file at path '%s' and try again.", config.GetOutputFilePath())
				err := shell.RequireApproval(prompt, errMsg)
				if err != nil {
					logging.ErrorF(err)
				}
			} else {
				logging.Infof("Force accept configured. Not asking for permission to append to file '%s'.", config.GetOutputFilePath())
			}
		}

		err = cloudSyncOptIn()
		if err != nil {
			logging.ErrorF(err)
		}

	},
	Run: func(cmd *cobra.Command, args []string) {
		app.RunDiscovery()
	},
}