package app

import (
	"github.com/lavalamp-/ipv666/internal/addressing"
	"github.com/lavalamp-/ipv666/internal/blacklist"
	"github.com/lavalamp-/ipv666/internal/config"
	"github.com/lavalamp-/ipv666/internal/data"
	"github.com/lavalamp-/ipv666/internal/fs"
	"github.com/lavalamp-/ipv666/internal/logging"
	"github.com/lavalamp-/ipv666/internal/shell"
	"github.com/spf13/viper"
	"net"
)

func RunBlgen(inputPath string) {

	var newBlacklist *blacklist.NetworkBlacklist

	approved, err := shell.AskForApproval("Would you like to add to the existing blacklist (if not, a new one will be created)? [y/N]")

	if err != nil {
		logging.ErrorF(err)
	}

	if approved {
		logging.Debugf("Loading existing blacklist...")
		newBlacklist, err = data.GetBlacklist()
		if err != nil {
			logging.ErrorF(err)
		}
	} else {
		newBlacklist = blacklist.NewNetworkBlacklist([]*net.IPNet{})
	}

	networks, err := addressing.ReadIPv6NetworksFromHexFile(inputPath)

	if err != nil {
		logging.ErrorStringFf("Error thrown when reading IPv6 networks from file '%s': %e", inputPath, err)
	}

	uniqueNetworks := addressing.GetUniqueNetworks(networks, viper.GetInt("LogLoopEmitFreq"))
	logging.Debugf("%d networks trimmed down to %d unique networks.", len(networks), len(uniqueNetworks))

	logging.Debugf("Adding %d network ranges to blacklist (starting with %d addresses).", len(uniqueNetworks), newBlacklist.GetCount())

	added, skipped := newBlacklist.AddNetworks(uniqueNetworks)
	logging.Debugf("%d addresses were added and %d were skipped.", added, skipped)

	startCount := newBlacklist.GetCount()
	newBlacklist.Clean(viper.GetInt("LogLoopEmitFreq"))
	logging.Infof("Cleaned up duplicated networks from blacklist. Down to %d networks (from %d).", newBlacklist.GetCount(), startCount)

	outputPath := fs.GetTimedFilePath(config.GetNetworkBlacklistDirPath())

	logging.Debugf("Writing network blacklist with %d network ranges to file at path '%s'.", newBlacklist.GetCount(), outputPath)

	err = blacklist.WriteNetworkBlacklistToFile(outputPath, newBlacklist)

	if err != nil {
		logging.Warnf("Error thrown when writing blacklist to file '%s': %e", outputPath, err)
	}

	logging.Successf("Successfully generated blacklist file at path '%s' using input addresses from file '%s' (list was %d long).", outputPath, inputPath, newBlacklist.GetCount())

}
