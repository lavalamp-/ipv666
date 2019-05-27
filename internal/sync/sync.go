package sync

import (
	"encoding/json"
	"fmt"
	"github.com/lavalamp-/ipv666/internal/addressing"
	"github.com/lavalamp-/ipv666/internal/logging"
	"github.com/rcrowley/go-metrics"
	"github.com/spf13/viper"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"
)

var syncIpCount = metrics.NewCounter()
var syncIpSuccessCount = metrics.NewCounter()
var syncIpFailureCount = metrics.NewCounter()
var syncFailures []time.Time

func init() {
	metrics.Register("sync.ip.count", syncIpCount)
	metrics.Register("sync.attempts.success.count", syncIpSuccessCount)
	metrics.Register("sync.attempts.failure.count", syncIpFailureCount)
}

type urlResponse struct {
	UploadUrl string `json:"upload_url"`
}

func SyncIpAddresses(toSync []*net.IP) {
	if !checkForSync() {
		goodTime := syncFailures[0].Add(time.Duration(viper.GetFloat64("SyncBackoffSeconds")) * time.Second)
		logging.Warnf("Not syncing IP addresses (%d failures seen in last %d seconds, waiting until %s).", len(syncFailures), viper.GetInt("SyncBackoffSeconds"), goodTime)
	}
	go func(addrs []*net.IP) {
		err := syncIpAddressesRoutine(addrs)
		if err != nil {
			recordFailure()
		} else {
			syncIpCount.Inc(int64(len(toSync)))
			syncIpSuccessCount.Inc(1)
		}
	}(toSync)
}

func checkForSync() bool {
	if len(syncFailures) != viper.GetInt("SyncFailureThreshold") {
		return true
	} else {
		return time.Since(syncFailures[0]).Seconds() > viper.GetFloat64("SyncBackoffSeconds")
	}
}

func recordFailure() {
	syncIpFailureCount.Inc(1)
	syncFailures = append(syncFailures, time.Now())
	if len(syncFailures) > viper.GetInt("SyncFailureThreshold") {
		syncFailures = syncFailures[1:]
	}
}

func syncIpAddressesRoutine(toSync []*net.IP) error {

	logging.Debugf("Attempting to sync %d addresses to remote server.", len(toSync))

	client := http.Client {
		Timeout: time.Duration(viper.GetInt("SyncTimeout")) * time.Second,
	}

	uploadUrl, err := fetchUploadUrl(&client)
	if err != nil {
		logging.Warnf("Error thrown when attempting to retrieve fetch URL: %e", err)
		return err
	}

	err = putAddressesToUrl(toSync, uploadUrl, &client)
	if err != nil {
		logging.Warnf("Error thrown when pushing addresses to remote server: %s", err)
		return err
	}

	logging.Successf("Successfully synced %d addresses to remote server.", len(toSync))
	return nil

}

// https://gist.github.com/slav123/cbb3309052de5a870667

func putAddressesToUrl(toPut []*net.IP, url string, client *http.Client) error {

	logging.Debugf("Putting %d addresses to URL '%s'.", len(toPut), url)

	stringContent := addressing.GetTextLinesFromIPs(toPut)

	request, err := http.NewRequest("PUT", url, strings.NewReader(stringContent))
	if err != nil {
		return err
	}

	request.ContentLength = int64(len(stringContent))

	logging.Debug("Sending request...")

	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return fmt.Errorf("did not get 200 response from upload URL push (got %d)", response.StatusCode)
	}

	logging.Debugf("Successfully pushed %d addresses to URL '%s'.", len(toPut), url)
	return nil

}

func fetchUploadUrl(client *http.Client) (string, error) {

	logging.Debugf("Fetching upload URL from ipv666 server.")

	req, err := http.NewRequest("GET", viper.GetString("SyncUrl"), nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", viper.GetString("SyncUserAgent"))
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("did not get 200 response from upload URL retrieval (got %d)", resp.StatusCode)
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	urlContent := urlResponse{}
	err = json.Unmarshal(content, &urlContent)
	if err != nil {
		return "", err
	} else {
		logging.Debugf("Successfully fetched upload URL of '%s'", urlContent.UploadUrl)
		return urlContent.UploadUrl, nil
	}

}
