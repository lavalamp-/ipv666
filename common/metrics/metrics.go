package metrics

import (
	"fmt"
	"github.com/cyberdelia/go-metrics-graphite"
	"github.com/lavalamp-/ipv666/common/config"
	"github.com/rcrowley/go-metrics"
	"github.com/spf13/viper"
	"log"
	"net"
	"os"
	"time"
)

func InitMetrics() error {
	if viper.GetBool("MetricsToStdout") {
		log.Printf("Setting up metrics to print to stdout every %d seconds.", viper.GetInt64("MetricsStdoutFreq"))
		go metrics.Log(metrics.DefaultRegistry, time.Duration(viper.GetInt64("MetricsStdoutFreq")) * time.Second, log.New(os.Stdout, "metrics: ", log.Lmicroseconds))
	} else {
		log.Printf("Not printing metrics to stdout.")
	}
	if viper.GetBool("GraphiteExportEnabled") {
		graphiteEndpoint := fmt.Sprintf("%s:%d", viper.GetString("GraphiteHost"), viper.GetInt("GraphitePort"))
		log.Printf("Configured to export to Graphite at %s (%s frequency).", graphiteEndpoint, config.GetGraphiteEmitDuration())
		addr, err := net.ResolveTCPAddr("tcp", graphiteEndpoint)
		if err != nil {
			log.Printf("Error thrown when resolving TCP address %s: %e", graphiteEndpoint, err)
			return err
		}
		go graphite.Graphite(metrics.DefaultRegistry, config.GetGraphiteEmitDuration(), "metrics", addr)
		log.Printf("Export to Graphite at %s set up and running.", graphiteEndpoint)
	}
	return nil
}