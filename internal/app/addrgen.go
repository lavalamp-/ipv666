package app

import (
	"github.com/lavalamp-/ipv666/internal/addressing"
	"github.com/lavalamp-/ipv666/internal/data"
	"github.com/lavalamp-/ipv666/internal/logging"
	"github.com/lavalamp-/ipv666/internal/modeling"
	"github.com/spf13/viper"
	"net"
)

func RunAddrGen(modelPath string, outputPath string, fromNetwork string, genCount int) {

	var clusterModel *modeling.ClusterModel
	var err error

	if modelPath == "" {
		logging.Info("No model path specified. Using default model packaged with IPv666.")
		clusterModel, err = data.GetProbabilisticClusterModel()
	} else {
		logging.Infof("Using cluster model found at path '%s'.", modelPath)
		clusterModel, err = modeling.LoadModelFromFile(modelPath)
	}

	if err != nil {
		logging.ErrorF(err)
	}

	var generatedAddrs []*net.IP

	if fromNetwork == "" {
		logging.Info("No network specified. Generating addresses in the global address space.")
		generatedAddrs = clusterModel.GenerateAddresses(genCount, viper.GetFloat64("ModelGenerationJitter"))
	} else {
		_, ipnet, _ := net.ParseCIDR(fromNetwork)
		logging.Infof("Generating addresses in specified network range of '%s'.", ipnet)
		generatedAddrs, err = clusterModel.GenerateAddressesFromNetwork(genCount, viper.GetFloat64("ModelGenerationJitter"), ipnet)
		if err != nil {
			logging.ErrorF(err)
		}
	}

	logging.Infof("Successfully generated %d IP addresses. Writing results to file at path '%s'.", genCount, outputPath)

	err = addressing.WriteIPsToHexFile(outputPath, generatedAddrs)  //TODO allow users to specify what type of file to write

	if err != nil {
		logging.Error(err)
	}

	logging.Infof("Successfully wrote addresses to file '%s'.", outputPath)

}
