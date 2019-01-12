package statemachine

import (
	"bufio"
	"fmt"
	"github.com/lavalamp-/ipv666/common/config"
	"github.com/lavalamp-/ipv666/common/data"
	"github.com/rcrowley/go-metrics"
	"github.com/spf13/viper"
	"log"
	"os"
	"time"
)

var addressUpdateTimer = metrics.NewTimer()

func init() {
	metrics.Register("addrupdate.file_write.time", addressUpdateTimer)
}

func updateAddressFile() error {
	cleanPings, err := data.GetCleanPingResults(config.GetCleanPingDirPath())
	if err != nil {
		return err
	}
	//TODO don't write addresses in input file in output file
	outputPath := config.GetOutputFilePath()
	log.Printf("Updating file at path '%s' with %d newly-found IP addresses.", outputPath, len(cleanPings))
	file, err := os.OpenFile(outputPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	writer := bufio.NewWriter(file)
	if err != nil {
		return err
	}
	defer file.Close()
	start := time.Now()
	if viper.GetString("OutputFileType") != "bin" {
		if !(viper.GetString("OutputFileType") == "text") { //TODO figure out why the != check fails but this works
			log.Printf("Unexpected file format for output (%s). Defaulting to text.", viper.GetString("OutputFileType"))
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
	log.Printf("Finished writing %d addresses to '%s'.", len(cleanPings), outputPath)
	return nil
}
