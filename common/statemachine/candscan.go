package statemachine

import (
	"github.com/lavalamp-/ipv666/common/config"
	"github.com/lavalamp-/ipv666/common/data"
	"log"
	"time"
	"github.com/lavalamp-/ipv666/common/shell"
)

func zmapScanCandidateAddresses(conf *config.Configuration) (error) {
	inputPath, err := data.GetMostRecentFilePathFromDir(conf.GetCandidateAddressDirPath())
	if err != nil {
		return err
	}
	outputPath := getTimedFilePath(conf.GetPingResultDirPath())
	log.Printf(
		"Now Zmap scanning IPv6 addresses found in file at path '%s'. Results will be written to '%s'.",
		inputPath,
		outputPath,
	)
	start := time.Now()
	_, err = shell.ZmapScanFromConfig(conf, inputPath, outputPath)
	elapsed := time.Since(start)
	//  TODO do something with result
	if err != nil {
		log.Printf("An error was thrown when trying to run zmap: %s", err)
		log.Printf("Zmap elapsed time was %s.", elapsed)
		return err
	}
	log.Printf("Zmap completed successfully in %s. Results written to file at '%s'.", elapsed, outputPath)
	return nil
}
