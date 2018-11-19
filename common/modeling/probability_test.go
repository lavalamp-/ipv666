package modeling

import (
	"testing"
	"github.com/lavalamp-/ipv666/common/config"
	"github.com/stretchr/testify/assert"
	"github.com/lavalamp-/ipv666/common/addressing"
	"log"
	"time"
	"math/rand"
)

var conf *config.Configuration

func init() {
	loadedConf, err := config.LoadFromFile("../../config.json")
	rand.Seed(time.Now().UTC().UnixNano())
	if err != nil {
		log.Fatal(err)
	}
	conf = &loadedConf
}

func getNewAddressModel() (*ProbabilisticAddressModel) {
	return NewAddressModel("Testing Model", conf)
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
