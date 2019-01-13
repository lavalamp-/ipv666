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
	"net"
	"github.com/cyberdelia/go-metrics-graphite"
	"github.com/lavalamp-/ipv666/common/validation"
	"github.com/lavalamp-/ipv666/common/data"
	"github.com/lavalamp-/ipv666/common/setup"
)

var mainLoopRunTimer = metrics.NewTimer()

//TODO switch all the various interval logging emissions to the single config value

func init() {
	metrics.Register("main.run.time", mainLoopRunTimer)
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

func initMetrics(conf *config.Configuration) (error) {
	if conf.MetricsToStdout {
		log.Printf("Setting up metrics to print to stdout every %d seconds.", conf.MetricsStdoutFreq)
		go metrics.Log(metrics.DefaultRegistry, time.Duration(conf.MetricsStdoutFreq) * time.Second, log.New(os.Stdout, "metrics: ", log.Lmicroseconds))
	} else {
		log.Printf("Not printing metrics to stdout.")
	}
	if conf.GraphiteExportEnabled {
		graphiteEndpoint := fmt.Sprintf("%s:%d", conf.GraphiteHost, conf.GraphitePort)
		log.Printf("Configured to export to Graphite at %s (%s frequency).", graphiteEndpoint, conf.GetGraphiteEmitDuration())
		addr, err := net.ResolveTCPAddr("tcp", graphiteEndpoint)
		if err != nil {
			log.Printf("Error thrown when resolving TCP address %s: %e", graphiteEndpoint, err)
			return err
		}
		go graphite.Graphite(metrics.DefaultRegistry, conf.GetGraphiteEmitDuration(), "metrics", addr)
		log.Printf("Export to Graphite at %s set up and running.", graphiteEndpoint)
	}
	return nil
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
	var forceAccept bool
	var targetNetworkString string

	flag.StringVar(&configPath, "config", "config.json", "Local file path to the configuration file to use.")
	flag.StringVar(&inputFile, "input", "", "An input file containing IPv6 addresses to initiate scanning from.")
	flag.StringVar(&inputType, "input-type", "txt", "The type of file pointed to by the 'input' argument (bin or txt).")
	flag.StringVar(&outputFile, "output", "", "The path to the file where discovered addresses should be written.")
	flag.StringVar(&outputType, "output-type", "txt", "The type of output to write to the output file (txt or bin).")
	flag.StringVar(&targetNetworkString, "network", "", "The target IPv6 network range to scan in. If empty, defaults to 2000::/4")
	flag.BoolVar(&forceAccept, "force", false, "Whether or not to force accept all prompts (useful for daemonized scanning).")

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
	conf.ForceAcceptPrompts = forceAccept

	if _, err := os.Stat(conf.GetOutputFilePath()); !os.IsNotExist(err) {
		if !conf.ForceAcceptPrompts {
			prompt := fmt.Sprintf("Output file already exists at path '%s,' continue (will append to existing file)? [y/N]", conf.GetOutputFilePath())
			errMsg := fmt.Sprintf("Exiting. Please move the file at path '%s' and try again.", conf.GetOutputFilePath())
			err := shell.RequireApproval(prompt, errMsg)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			log.Printf("Force accept configured. Not asking for permission to append to file '%s'.", conf.GetOutputFilePath())
		}

	}

	var targetNetwork *net.IPNet
	if targetNetworkString != "" {
		targetNetwork, err = validation.ValidateIPv6NetworkStringForScanning(targetNetworkString, &conf)
		if err != nil {
			log.Fatalf("The target network of '%s' was not valid: %e", targetNetworkString, err)
		}
		conf.SetTargetNetwork(targetNetwork)
	} else {
		targetNetwork, err = conf.GetTargetNetwork()
		if err != nil {
			log.Fatalf("Error thrown when creating default target network: %e", err)
		}
	}
	log.Printf("Target network to scan will be %s.", targetNetwork)

	mostRecentNetworkString, err := data.GetMostRecentTargetNetworkString(&conf)
	if err != nil {
		log.Fatalf("Error thrown when reading most recent network string: %e", err)
	}
	if mostRecentNetworkString != targetNetwork.String() {
		if mostRecentNetworkString == "" {
			log.Printf("No prior record of a scanned network exists. Resetting state machine to scan %s appropriately.", targetNetwork)
		} else {
			log.Printf("Target network (%s) is not the most recently scanned network (%s). Resetting state machine and Bloom filter accordingly.", targetNetwork, mostRecentNetworkString)
		}
		err := statemachine.ResetStateFile(conf.GetStateFilePath())
		if err != nil {
			log.Fatalf("Error thrown when resetting state file: %e", err)
		}
		_, _, err = fs.DeleteAllFilesInDirectory(conf.GetBloomDirPath(), []string{})
		if err != nil {
			log.Fatalf("Error thrown when deleting Bloom directory files (path '%s'): %e", conf.GetBloomDirPath(), err)
		}
		err = data.WriteMostRecentTargetNetwork(targetNetwork, &conf)
		if err != nil {
			log.Fatalf("Error thrown when writing most recent target network: %e", err)
		}
	} else {
		log.Printf("The network %s is the last network that was targeted. Picking up from where we left off.", targetNetwork)
	}

	if !conf.LogToFile {
		log.Printf("Not configured to log to file. Logging to stdout instead.")
	} else {
		setupLogging(&conf)
	}

	err = setup.InitFilesystem(&conf)

	if err != nil {
		log.Fatal("Error thrown during filesystem initialization: ", err)
	}

	err = initMetrics(&conf)
	if err != nil {
		log.Fatal(err)
	}

	if inputFile != "" {
		err := input.PrepareFromInputFile(inputFile, inputType, &conf)
		if err != nil {
			log.Fatal(fmt.Sprintf("Error thrown when preparing to scan from input file '%s': %e", inputFile, err))
		}
	}

	rand.Seed(time.Now().UTC().UnixNano())

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
