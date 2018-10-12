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
	"github.com/rcrowley/go-metrics"
	"time"
	"math/rand"
	"fmt"
	"github.com/lavalamp-/ipv666/common/input"
)

var mainLoopRunTimer = metrics.NewTimer()

//TODO switch all the various interval logging emissions to the single config value

func init() {
	metrics.Register("main_loop_timer", mainLoopRunTimer)
}

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
  
func initFilesystem(conf *config.Configuration) (error) {
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

func initMetrics(conf *config.Configuration) () {
	if conf.MetricsToStdout {
		log.Printf("Setting up metrics to print to stdout every %d seconds.", conf.MetricsStdoutFreq)
		go metrics.Log(metrics.DefaultRegistry, time.Duration(conf.MetricsStdoutFreq) * time.Second, log.New(os.Stdout, "metrics: ", log.Lmicroseconds))
	} else {
		log.Printf("Not printing metrics to stdout.")
	}
}

func isValidFileType(toCheck string) (bool) {
	return toCheck == "txt" || toCheck == "bin"
}
  
func main() {

	//TODO refactor into a input struct and its own function (input handling)

	var configPath string
	var inputFile string
	var inputType string
	var outputFile string
	var outputType string

	flag.StringVar(&configPath, "config", "config.json", "Local file path to the configuration file to use.")
	flag.StringVar(&inputFile, "input", "", "An input file containing IPv6 addresses to initiate scanning from.")
	flag.StringVar(&inputType, "input-type", "txt", "The type of file pointed to by the 'input' argument (bin or txt).")
	flag.StringVar(&outputFile, "output", "", "The path to the file where discovered addresses should be written.")
	flag.StringVar(&outputType, "output-type", "txt", "The type of output to write to the output file (txt or bin).")

	flag.Parse()

	if inputFile != "" && !isValidFileType(inputType) {
		log.Fatal(fmt.Sprintf("%s is not a valid input file type (requires txt or bin).", inputType))
	}

	if !isValidFileType(outputType) {
		log.Fatal(fmt.Sprintf("%s is not a valid output file type (requires txt or bin).", outputType))
	}

	if inputFile != "" {
		if _, err := os.Stat(inputFile); os.IsNotExist(err) {
			log.Fatal(fmt.Sprintf("No file found at input file path of '%s'.", inputFile))
		}
	}

	conf, err := config.LoadFromFile(configPath)

	if err != nil {
		log.Fatal("Can't proceed without loading valid configuration file.")
	}

	if !(outputFile == "") { //TODO figure out why the straight != check is failing
		conf.OutputFileName = outputFile
		conf.OutputFileType = outputType
	}

	if _, err := os.Stat(conf.GetOutputFilePath()); !os.IsNotExist(err) {
		prompt := fmt.Sprintf("Output file already exists at path '%s,' continue (will append to existing file)? [y/N]", conf.GetOutputFilePath())
		errMsg := fmt.Sprintf("Exiting. Please move the file at path '%s' and try again.", conf.GetOutputFilePath())
		err := shell.PromptForApproval(prompt, errMsg)
		if err != nil {
			log.Fatal(err)
		}
	}

	if !conf.LogToFile {
		log.Printf("Not configured to log to file. Logging to stdout instead.")
	} else {
		setupLogging(&conf)
	}

	err = initFilesystem(&conf)

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

	initMetrics(&conf)
	rand.Seed(time.Now().UTC().UnixNano())

	if inputFile != "" {
		err := input.PrepareFromInputFile(inputFile, inputType, &conf)
		if err != nil {
			log.Fatal(fmt.Sprintf("Error thrown when preparing to scan from input file '%s': %e", inputFile, err))
		}
	}

	log.Print("All systems are green. Entering state machine.")

	start := time.Now()
	err = statemachine.RunStateMachine(&conf)
	elapsed := time.Since(start)
	mainLoopRunTimer.Update(elapsed)

	//TODO push metrics

	if err != nil {
		log.Fatal(err)
	}

}
