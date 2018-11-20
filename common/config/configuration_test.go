package config

import (
	"testing"
	"github.com/prometheus/common/log"
	"net"
	"github.com/stretchr/testify/assert"
)

func getTestingConfig() (*Configuration) {
	conf, err := LoadFromFile("../../config.json")
	if err != nil {
		log.Fatal(err)
	}
	return &conf
}

func TestConfiguration_SetTargetNetworkSets(t *testing.T) {
	conf := getTestingConfig()
	_, network, _ := net.ParseCIDR("2600::1/64")
	conf.SetTargetNetwork(network)
	assert.NotNil(t, conf.targetNetwork)
}

func TestConfiguration_SetTargetNetworkSetsValue(t *testing.T) {
	conf := getTestingConfig()
	_, network, _ := net.ParseCIDR("2600::1/64")
	conf.SetTargetNetwork(network)
	expected := network.String()
	assert.EqualValues(t, expected, conf.targetNetwork.String())
}

func TestConfiguration_GetTargetNetwork(t *testing.T) {
	conf := getTestingConfig()
	_, err := conf.GetTargetNetwork()
	assert.Nil(t, err)
}

func TestConfiguration_GetTargetNetworkNilCreates(t *testing.T) {
	conf := getTestingConfig()
	network, _ := conf.GetTargetNetwork()
	expected := "2000::/4"
	assert.EqualValues(t, expected, network.String())
}

func TestConfiguration_GetTargetNetworkNotNilReturns(t *testing.T) {
	conf := getTestingConfig()
	_, network, _ := net.ParseCIDR("2600::1/64")
	conf.SetTargetNetwork(network)
	expected := network.String()
	newNet, _ := conf.GetTargetNetwork()
	assert.EqualValues(t, expected, newNet.String())
}
