package setup

import (
	"github.com/lavalamp-/ipv666/common/config"
	"github.com/lavalamp-/ipv666/common/fs"
	"github.com/lavalamp-/ipv666/common/logging"
	"github.com/lavalamp-/ipv666/common/statemachine"
	"os"
)

func InitFilesystem() error {
	logging.Debug("Now initializing filesystem...")
	for _, dirPath := range config.GetAllDirectories() {
		err := fs.CreateDirectoryIfNotExist(dirPath)
		if err != nil {
			return err
		}
	}
	logging.Debugf("Initializing state file at '%s'.", config.GetStateFilePath())
	if _, err := os.Stat(config.GetStateFilePath()); os.IsNotExist(err) {
		logging.Debugf("State file does not exist at path '%s'. Creating now.", config.GetStateFilePath())
		err = statemachine.InitStateFile(config.GetStateFilePath())
		if err != nil {
			return err
		}
	} else {
		logging.Debugf("State file already exists at path '%s'.", config.GetStateFilePath())
	}
	logging.Debug("Local filesystem initialized.")
	return nil
}
