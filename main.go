package main

import (
	"log"
	"github.com/lavalamp-/ipv666/common/config"
	"os"
	"github.com/natefinch/lumberjack"
	"github.com/lavalamp-/ipv666/common/fs"
	"flag"
	"github.com/lavalamp-/ipv666/common/shell"
	"github.com/lavalamp-/ipv666/common/statemachine"
)

func setupLogging(conf *config.Configuration) {
	log.Print("Now setting up logging.")
	log.SetFlags(log.Flags() & (log.Ldate | log.Ltime))
  	log.SetOutput(&lumberjack.Logger{
  		Filename:   conf.LogFilePath,
  		MaxSize:    conf.LogFileMBSize,		// megabytes
  		MaxBackups: conf.LogFileMaxBackups,
  		MaxAge:     conf.LogFileMaxAge,		// days
  		Compress:   conf.CompressLogFiles,
  	})
	log.Print("Logging set up successfully.")
}
  
func initializeFilesystem(conf *config.Configuration) (error) {
	log.Print("Now initializing filesystem for IPv6 address discovery process...")
	for _, dirPath := range conf.GetAllDirectories() {
		err := fs.CreateDirectoryIfNotExist(dirPath)
		if err != nil {
			return err
		}
	}
	log.Printf("Initializing state file at '%s'.", conf.GetStateFilePath())
	if _, err := os.Stat(conf.GetStateFilePath()); os.IsNotExist(err) {
		log.Printf("State file does not exist at path '%s'. Creating now.", conf.GetStateFilePath())
		err = statemachine.InitStateFile(conf.GetStateFilePath())
		if err != nil {
			return err
		}
	} else {
		log.Printf("State file already exists at path '%s'.", conf.GetStateFilePath())
	}
	log.Print("Local filesystem initialized for IPv6 address discovery process.")
	return nil
}
  
func main() {

	var configPath string

	flag.StringVar(&configPath, "config", "config.json", "Local file path to the configuration file to use.")

	conf, err := config.LoadFromFile(configPath)

	if err != nil {
		log.Fatal("Can't proceed without loading valid configuration file.")
	}

	if !conf.LogToFile {
		log.Printf("Not configured to log to file. Logging to stdout instead.")
	} else {
		setupLogging(&conf)
	}

	err = initializeFilesystem(&conf)

	if err != nil {
		log.Fatal("Error thrown during initialization: ", err)
	}

	zmapAvailable, err := shell.IsZmapAvailable(&conf)

	if err != nil {
		log.Fatal("Error thrown when checking for Zmap: ", err)
	} else if !zmapAvailable {
		log.Fatal("Zmap not found. Please install Zmap.")
	}

	log.Printf("Zmap found and working at path '%s'.", conf.ZmapExecPath)

	log.Print("All systems are green. Entering state machine.")

	err = statemachine.RunStateMachine(&conf)

	if err != nil {
		log.Fatal(err)
	}

}
