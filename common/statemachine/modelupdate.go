package statemachine

import (
	"github.com/lavalamp-/ipv666/common/config"
	"github.com/lavalamp-/ipv666/common/data"
	"github.com/lavalamp-/ipv666/common/addresses"
	"log"
)

func updateModelWithSuccessfulHosts(conf *config.Configuration) (error) {
	// TODO what happens if Zmap fails silently and we keep adding the same file to our model
	resultsPath, err := data.GetMostRecentFilePathFromDir(conf.GetCleanPingDirPath())
	if err != nil {
		return err
	}
	results, err := addresses.GetAddressListFromHexStringsFile(resultsPath)
	if err != nil {
		return err
	}
	model, err := data.GetProbabilisticAddressModel(conf.GetGeneratedModelDirPath())
	if err != nil {
		return err
	}
	model.UpdateMulti(results)
	outputPath := getTimedFilePath(conf.GetGeneratedModelDirPath())
	log.Printf("Now saving updated model '%s' (%d digest count) to file at path '%s'.", model.Name, model.DigestCount, outputPath)
	model.Save(outputPath)
	log.Printf("Model '%s' was saved to file at path '%s' successfully.", model.Name, outputPath)
	return nil
}
