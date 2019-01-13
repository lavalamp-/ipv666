package modeling

import (
	"github.com/lavalamp-/ipv666/internal/addressing"
	"github.com/lavalamp-/ipv666/internal/config"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"net"
	"testing"
	"time"
)

func init() {
	config.InitConfig()
	rand.Seed(time.Now().UTC().UnixNano())
}

func getNewAddressModel() *ProbabilisticAddressModel {
	return NewAddressModel("Testing Model")
}

func getFilterFunc() addrProcessFunc {
	checker := 0
	return func(toCheck *net.IP) (bool, error) {
		checker++
		return checker % 2 == 0, nil
	}
}

func TestProbabilisticAddressModel_GenerateMultiIPFromNetworkCount(t *testing.T) {
	model := getNewAddressModel()
	filterFunc := getFilterFunc()
	_, fromNetwork, _ := net.ParseCIDR("2000::/4")
	results, _ := model.GenerateMultiIPFromNetwork(fromNetwork, 20, filterFunc)
	assert.EqualValues(t, 20, len(results))
}

func TestProbabilisticAddressModel_GenerateMultiIPFromNetworkFilters(t *testing.T) {
	isChecked := false
	checkFunc := func(toCheck *net.IP) (bool, error) {
		isChecked = true
		return false, nil
	}
	model := getNewAddressModel()
	_, fromNetwork, _ := net.ParseCIDR("2000::/4")
	model.GenerateMultiIPFromNetwork(fromNetwork, 20, checkFunc)
	assert.True(t, isChecked)
}

func TestProbabilisticAddressModel_GenerateMultiIPFromNetworkContentNoMatch(t *testing.T) {
	model := getNewAddressModel()
	filterFunc := getFilterFunc()
	_, fromNetwork, _ := net.ParseCIDR("ffff:ffff:ffff:ffff::/16")
	newIPs, _ := model.GenerateMultiIPFromNetwork(fromNetwork, 20, filterFunc)
	expected := []byte{0x0f, 0x0f, 0x0f, 0x0f}
	for _, newIP := range newIPs {
		newIPBytes := addressing.GetNybblesFromIP(newIP, 4)
		assert.ElementsMatch(t, expected, newIPBytes)
	}
}

func TestProbabilisticAddressModel_GenerateMultiIPFromNetworkContentMustMatchLowNybble(t *testing.T) {
	model := getNewAddressModel()
	filterFunc := getFilterFunc()
	_, fromNetwork, _ := net.ParseCIDR("ffff:ffff:ffff:ffff::/14")
	newIPs, _ := model.GenerateMultiIPFromNetwork(fromNetwork, 20, filterFunc)
	expected := []byte{0x0f, 0x0f, 0x0f, 0x0c}
	for _, newIP := range newIPs {
		newIPBytes := addressing.GetNybblesFromIP(newIP, 4)
		newIPBytes[3] &= 0x0c
		assert.ElementsMatch(t, expected, newIPBytes)
	}
}

func TestProbabilisticAddressModel_GenerateMultiIPFromNetworkContentMustMatchHighNybble(t *testing.T) {
	model := getNewAddressModel()
	filterFunc := getFilterFunc()
	_, fromNetwork, _ := net.ParseCIDR("ffff:ffff:ffff:ffff::/18")
	newIPs, _ := model.GenerateMultiIPFromNetwork(fromNetwork, 20, filterFunc)
	expected := []byte{0x0f, 0x0f, 0x0f, 0x0f, 0x0c}
	for _, newIP := range newIPs {
		newIPBytes := addressing.GetNybblesFromIP(newIP, 5)
		newIPBytes[4] &= 0x0c
		assert.ElementsMatch(t, expected, newIPBytes)
	}
}

func TestProbabilisticAddressModel_GenerateSingleIPFromNybblesMinOffsetNoMatch(t *testing.T) {
	model := getNewAddressModel()
	fromNybbles := []byte{0x02}
	newIP := model.GenerateSingleIPFromNybbles(fromNybbles, 4)
	ipNybbles := addressing.GetNybblesFromIP(newIP, 1)
	assert.ElementsMatch(t, []byte{0x02}, ipNybbles[0:1])
}

func TestProbabilisticAddressModel_GenerateSingleIPFromNybblesMinOffsetMustMatch(t *testing.T) {
	model := getNewAddressModel()
	fromNybbles := []byte{0x02, 0x0c}
	newIP := model.GenerateSingleIPFromNybbles(fromNybbles, 6)
	ipNybbles := addressing.GetNybblesFromIP(newIP, 2)
	ipNybbles[1] &= 0xfc
	assert.ElementsMatch(t, fromNybbles, ipNybbles)
}

func TestProbabilisticAddressModel_GenerateSingleIPFromNybblesNoAlterInput(t *testing.T) {
	model := getNewAddressModel()
	fromNybbles := []byte{0x02, 0x0c}
	for i := 0; i < 100; i++ {
		model.GenerateSingleIPFromNybbles(fromNybbles, 6)
		assert.EqualValues(t, 0x0c, fromNybbles[1])
	}
}

func TestProbabilisticAddressModel_GenerateSingleIPFromNybblesMaxOffsetNoMatch(t *testing.T) {
	model := getNewAddressModel()
	fromNybbles := []byte{0x0f, 0x0f, 0x0f, 0x0f, 0x0f, 0x0f, 0x0f, 0x0f, 0x0f, 0x0f, 0x0f, 0x0f, 0x0f, 0x0f, 0x0f, 0x0f, 0x0f, 0x0f, 0x0f, 0x0f, 0x0f, 0x0f, 0x0f, 0x0f, 0x0f, 0x0f, 0x0f, 0x0f, 0x0f, 0x0f, 0x0f}
	newIP := model.GenerateSingleIPFromNybbles(fromNybbles, 124)
	ipNybbles := addressing.GetNybblesFromIP(newIP, 31)
	assert.ElementsMatch(t, fromNybbles, ipNybbles)
}

func TestProbabilisticAddressModel_GenerateSingleIPFromNybblesMaxOffsetMustMatch(t *testing.T) {
	model := getNewAddressModel()
	fromNybbles := []byte{0x0f, 0x0f, 0x0f, 0x0f, 0x0f, 0x0f, 0x0f, 0x0f, 0x0f, 0x0f, 0x0f, 0x0f, 0x0f, 0x0f, 0x0f, 0x0f, 0x0f, 0x0f, 0x0f, 0x0f, 0x0f, 0x0f, 0x0f, 0x0f, 0x0f, 0x0f, 0x0f, 0x0f, 0x0f, 0x0f, 0x0f, 0x0c}
	newIP := model.GenerateSingleIPFromNybbles(fromNybbles, 126)
	ipNybbles := addressing.GetNybblesFromIP(newIP, 32)
	ipNybbles[31] &= 0xfc
	assert.ElementsMatch(t, fromNybbles, ipNybbles)
}

func TestProbabilisticAddressModel_GenerateSingleIPFromNybblesNoMatch(t *testing.T) {
	model := getNewAddressModel()
	fromNybbles := []byte{0x0f, 0x0f, 0x0f, 0x0f, 0x0f, 0x0f}
	newIP := model.GenerateSingleIPFromNybbles(fromNybbles, 24)
	ipNybbles := addressing.GetNybblesFromIP(newIP, 6)
	assert.ElementsMatch(t, fromNybbles, ipNybbles)
}

func TestProbabilisticAddressModel_GenerateSingleIPFromNybblesMustMatch(t *testing.T) {
	model := getNewAddressModel()
	fromNybbles := []byte{0x0f, 0x0f, 0x0f, 0x0f, 0x0f, 0x0f, 0x08}
	newIP := model.GenerateSingleIPFromNybbles(fromNybbles, 25)
	ipNybbles := addressing.GetNybblesFromIP(newIP, 7)
	ipNybbles[6] &= 0xf8
	assert.ElementsMatch(t, fromNybbles, ipNybbles)
}
