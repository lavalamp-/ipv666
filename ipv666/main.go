package main

import (
	"github.com/lavalamp-/ipv666/common/config"
	"github.com/lavalamp-/ipv666/common/logging"
	"github.com/lavalamp-/ipv666/common/metrics"
	"github.com/lavalamp-/ipv666/common/setup"
	"github.com/lavalamp-/ipv666/ipv666/cmd"
	"log"
	"math/rand"
	"time"
)

func main() {
	config.InitConfig()
	logging.SetupLogging()
	err := setup.InitFilesystem()
	if err != nil {
		log.Fatal(err)
	}
	err = metrics.InitMetrics()
	if err != nil {
		log.Fatal(err)
	}
	rand.Seed(time.Now().UTC().UnixNano())
	cmd.Execute()
}
