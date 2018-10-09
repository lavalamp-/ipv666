package statemachine

import (
	"github.com/lavalamp-/ipv666/common/config"
	"github.com/lavalamp-/ipv666/common/data"
	"log"
	"time"
)

func generateCandidateAddresses(conf *config.Configuration) (error) {
	model, err := data.GetProbabilisticAddressModel(conf.GetGeneratedModelDirPath())
	if err != nil {
		return err
	}
	log.Printf(
		"Generating a total of %d addressing based on the content of model '%s' (%d digest count). Starting nybble is %d.",
		conf.GenerateAddressCount,
		model.Name,
		model.DigestCount,
		conf.GenerateFirstNybble,
	)
	start := time.Now()
	// TODO: filter out from blacklist
	// TODO: add metrics to how many are filtered out
	addresses := model.GenerateMulti(conf.GenerateFirstNybble, conf.GenerateAddressCount, conf.GenerateUpdateFreq)
	elapsed := time.Since(start)
	log.Printf("Took a total of %s to generate %d candidate addressing", elapsed, conf.GenerateAddressCount)
	outputPath := getTimedFilePath(conf.GetCandidateAddressDirPath())
	log.Printf("Writing results of candidate address generation to file at '%s'.", outputPath)
	addresses.ToAddressesFile(outputPath, conf.GenWriteUpdateFreq)
	log.Printf("Successfully wrote %d candidate addressing to file at '%s'.", conf.GenerateAddressCount, outputPath)
	return nil
}