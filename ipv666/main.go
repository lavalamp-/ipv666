package main

import (
	"github.com/lavalamp-/ipv666/internal/config"
	"github.com/lavalamp-/ipv666/internal/logging"
	"github.com/lavalamp-/ipv666/internal/metrics"
	"github.com/lavalamp-/ipv666/internal/setup"
	"github.com/lavalamp-/ipv666/ipv666/cmd"
	"math/rand"
	"time"
)

func main() {
	config.InitConfig()
	logging.SetupLogging()
	err := setup.InitFilesystem()
	if err != nil {
		logging.ErrorF(err)
	}
	err = metrics.InitMetrics()
	if err != nil {
		logging.ErrorF(err)
	}
	rand.Seed(time.Now().UTC().UnixNano())
	cmd.Execute()
}
