package statemachine

import (
	"bufio"
	"fmt"
	"github.com/lavalamp-/ipv666/internal/config"
	"github.com/lavalamp-/ipv666/internal/data"
	"github.com/lavalamp-/ipv666/internal/logging"
	"github.com/lavalamp-/ipv666/internal/sync"
	"github.com/rcrowley/go-metrics"
	"github.com/spf13/viper"
	"os"
	"time"
)

var addressUpdateTimer = metrics.NewTimer()

func init() {
	metrics.Register("addrupdate.file_write.time", addressUpdateTimer)
}

func updateAddressFile() error {
	cleanPings, err := data.GetCleanPingResults()
	if err != nil {
		return err
	}
	//TODO don't write addresses in input file in output file
	outputPath := config.GetOutputFilePath()
	logging.Infof("Updating file at path '%s' with %d newly-found IP addresses.", outputPath, len(cleanPings))
	file, err := os.OpenFile(outputPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	writer := bufio.NewWriter(file)
	if err != nil {
		return err
	}
	defer file.Close()
	start := time.Now()
	if viper.GetString("OutputFileType") != "bin" {
		if !(viper.GetString("OutputFileType") == "txt") { //TODO figure out why the != check fails but this works
			logging.Warnf("Unexpected file format for output (%s). Defaulting to text.", viper.GetString("OutputFileType"))
		}
		for _, addr := range cleanPings {
			writer.WriteString(fmt.Sprintf("%s\n", addr))
		}
	} else {
		for _, addr := range cleanPings {
			toWrite := ([]byte)(*addr)
			writer.Write(toWrite)
		}
	}
	writer.Flush()
	elapsed := time.Since(start)
	addressUpdateTimer.Update(elapsed)
	logging.Successf("%d new live IPv6 addresses were found.", len(cleanPings))
	logging.Debugf("Finished writing %d addresses to '%s'.", len(cleanPings), outputPath)
	if viper.GetBool("CloudSyncOptIn") {
		sync.SyncIpAddresses(cleanPings)
	}
	return nil
}
