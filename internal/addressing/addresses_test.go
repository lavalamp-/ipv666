package addressing

import (
	"github.com/lavalamp-/ipv666/internal/config"
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
)

func init() {
	config.InitConfig()
}

func getTestingIP() *net.IP {
	toReturn := net.ParseIP("2600::1")
	return &toReturn
}

func getTestingIPs() []*net.IP {
	ips := []net.IP{
		net.ParseIP("2600::0"),
		net.ParseIP("2600::1"),
		net.ParseIP("2601::1"),
	}
	var toReturn []*net.IP
	for i := 0; i < len(ips); i++ {
		toReturn = append(toReturn, &ips[i])
	}
	return toReturn
}

func getExpectedAdjacentIPs() []*net.IP {
	ips := []net.IP{
		net.ParseIP("2600::0"),
		net.ParseIP("2600::1"),
		net.ParseIP("2600::2"),
		net.ParseIP("2600::3"),
		net.ParseIP("2600::4"),
		net.ParseIP("2600::5"),
		net.ParseIP("2600::6"),
		net.ParseIP("2600::7"),
		net.ParseIP("2600::8"),
		net.ParseIP("2600::9"),
		net.ParseIP("2600::a"),
		net.ParseIP("2600::b"),
		net.ParseIP("2600::c"),
		net.ParseIP("2600::d"),
		net.ParseIP("2600::e"),
		net.ParseIP("2600::f"),
	}
	var toReturn []*net.IP
	for i := 0; i < len(ips); i++ {
		toReturn = append(toReturn, &ips[i])
	}
	return toReturn
}

func TestGetAdjacentNetworkAddressesFromIPsBadFromNybble(t *testing.T) {
	_, err := GetAdjacentNetworkAddressesFromIPs(getTestingIPs(), -1, 32)
	assert.NotNil(t, err)
}

func TestGetAdjacentNetworkAddressesFromIPsBadToNybble(t *testing.T) {
	_, err := GetAdjacentNetworkAddressesFromIPs(getTestingIPs(), 0, 33)
	assert.NotNil(t, err)
}

func TestGetAdjacentNetworkAddressesFromIPsBadFromAndToNybbles(t *testing.T) {
	_, err := GetAdjacentNetworkAddressesFromIPs(getTestingIPs(), 10, 10)
	assert.NotNil(t, err)
}

func TestGetAdjacentNetworkAddressesFromIPsCount(t *testing.T) {
	results, _ := GetAdjacentNetworkAddressesFromIPs(getTestingIPs(), 0, 32)
	assert.EqualValues(t, 1410, len(results))
}

func TestGetAdjacentNetworkAddressesFromIPsNoDuplicates(t *testing.T) {
	results, _ := GetAdjacentNetworkAddressesFromIPs(getTestingIPs(), 0, 32)
	firstCount := len(results)
	results = GetUniqueIPs(results, 99999)
	secondCount := len(results)
	assert.Equal(t, firstCount, secondCount)
}

func TestGetAdjacentNetworkAddressesFromIPBadFromNybble(t *testing.T) {
	_, err := GetAdjacentNetworkAddressesFromIP(getTestingIP(), -1, 32)
	assert.NotNil(t, err)
}

func TestGetAdjacentNetworkAddressesFromIPBadToNybble(t *testing.T) {
	_, err := GetAdjacentNetworkAddressesFromIP(getTestingIP(), 0, 33)
	assert.NotNil(t, err)
}

func TestGetAdjacentNetworkAddressesFromIPBadFromAndToNybbles(t *testing.T) {
	_, err := GetAdjacentNetworkAddressesFromIP(getTestingIP(), 10, 10)
	assert.NotNil(t, err)
}

func TestGetAdjacentNetworkAddressesFromIPSmallCount(t *testing.T) {
	results, _ := GetAdjacentNetworkAddressesFromIP(getTestingIP(), 31, 32)
	assert.EqualValues(t, 16, len(results))
}

func TestGetAdjacentNetworkAddressesFromIPLargeCount(t *testing.T) {
	results, _ := GetAdjacentNetworkAddressesFromIP(getTestingIP(), 0, 32)
	assert.EqualValues(t, 16 * 32 - 32 + 1, len(results))
}

func TestGetAdjacentNetworkAddressesFromIPSmallContent(t *testing.T) {
	results, _ := GetAdjacentNetworkAddressesFromIP(getTestingIP(), 31, 32)
	assert.ElementsMatch(t, getExpectedAdjacentIPs(), results)
}

func TestFlipBitsInAddressOnBoundaries(t *testing.T) {
	testAddr := net.ParseIP("aaaa:aaaa:aaaa:aaaa:aaaa:aaaa:aaaa:aaaa")
	flipAddr := FlipBitsInAddress(&testAddr, 8, 15)
	expectedAddr := net.ParseIP("aa55:aaaa:aaaa:aaaa:aaaa:aaaa:aaaa:aaaa")
	assert.Equal(t, expectedAddr.String(), flipAddr.String())
}

func TestFlipBitsInAddressOffBoundaryStart(t *testing.T) {
	testAddr := net.ParseIP("aaaa:aaaa:aaaa:aaaa:aaaa:aaaa:aaaa:aaaa")
	flipAddr := FlipBitsInAddress(&testAddr, 12, 15)
	expectedAddr := net.ParseIP("aaa5:aaaa:aaaa:aaaa:aaaa:aaaa:aaaa:aaaa")
	assert.Equal(t, expectedAddr.String(), flipAddr.String())
}

func TestFlipBitsInAddressOffBoundaryEnd(t *testing.T) {
	testAddr := net.ParseIP("aaaa:aaaa:aaaa:aaaa:aaaa:aaaa:aaaa:aaaa")
	flipAddr := FlipBitsInAddress(&testAddr, 8, 19)
	expectedAddr := net.ParseIP("aa55:5aaa:aaaa:aaaa:aaaa:aaaa:aaaa:aaaa")
	assert.Equal(t, expectedAddr.String(), flipAddr.String())
}

func TestFlipBitsInAddressSameByte(t *testing.T) {
	testAddr := net.ParseIP("aaaa:aaaa:aaaa:aaaa:aaaa:aaaa:aaaa:aaaa")
	flipAddr := FlipBitsInAddress(&testAddr, 8, 11)
	expectedAddr := net.ParseIP("aa5a:aaaa:aaaa:aaaa:aaaa:aaaa:aaaa:aaaa")
	assert.Equal(t, expectedAddr.String(), flipAddr.String())
}

func TestFlipBitsInAddressMultiByteAway(t *testing.T) {
	testAddr := net.ParseIP("aaaa:aaaa:aaaa:aaaa:aaaa:aaaa:aaaa:aaaa")
	flipAddr := FlipBitsInAddress(&testAddr, 8, 31)
	expectedAddr := net.ParseIP("aa55:5555:aaaa:aaaa:aaaa:aaaa:aaaa:aaaa")
	assert.Equal(t, expectedAddr.String(), flipAddr.String())
}

func TestFlipBitsInAddress(t *testing.T) {
	testAddr := net.ParseIP("aaaa:aaaa:aaaa:aaaa:aaaa:aaaa:aaaa:aaaa")
	flipAddr := FlipBitsInAddress(&testAddr, 64, 127)
	expectedAddr := net.ParseIP("aaaa:aaaa:aaaa:aaaa:5555:5555:5555:5555")
	assert.Equal(t, expectedAddr.String(), flipAddr.String())
}

func TestAddressToUintsZeroesFirst(t *testing.T) {
	testAddr := net.ParseIP("0000:0000:0000:0000:0000:0000:0000:0000")
	first, _ := AddressToUints(testAddr)
	assert.EqualValues(t, 0, first)
}

func TestAddressToUintsZeroesSecond(t *testing.T) {
	testAddr := net.ParseIP("0000:0000:0000:0000:0000:0000:0000:0000")
	_, second := AddressToUints(testAddr)
	assert.EqualValues(t, 0, second)
}

func TestAddressToUintsOnesFirst(t *testing.T) {
	testAddr := net.ParseIP("ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff")
	first, _ := AddressToUints(testAddr)
	assert.EqualValues(t, ^uint64(0), first)
}

func TestAddressToUintsOnesSecond(t *testing.T) {
	testAddr := net.ParseIP("ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff")
	_, second := AddressToUints(testAddr)
	assert.EqualValues(t, ^uint64(0), second)
}

func TestAddressToUintsOneFirst(t *testing.T) {
	testAddr := net.ParseIP("0000:0000:0000:0001:0000:0000:0000:0001")
	first, _ := AddressToUints(testAddr)
	assert.EqualValues(t, 1, first)
}

func TestAddressToUintsOneSecond(t *testing.T) {
	testAddr := net.ParseIP("0000:0000:0000:0001:0000:0000:0000:0001")
	_, second := AddressToUints(testAddr)
	assert.EqualValues(t, 1, second)
}

func TestUintsToAddressZeroes(t *testing.T) {
	testAddr := net.ParseIP("0000:0000:0000:0000:0000:0000:0000:0000")
	resultAddr := UintsToAddress(uint64(0), uint64(0))
	assert.EqualValues(t, testAddr, *resultAddr)
}

func TestUintsToAddressOnes(t *testing.T) {
	testAddr := net.ParseIP("ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff")
	resultAddr := UintsToAddress(^uint64(0), ^uint64(0))
	assert.EqualValues(t, testAddr, *resultAddr)
}

func TestUintsToAddressLow(t *testing.T) {
	testAddr := net.ParseIP("0000:0000:0000:0001:0000:0000:0000:0001")
	resultAddr := UintsToAddress(uint64(1), uint64(1))
	assert.EqualValues(t, testAddr, *resultAddr)
}

func TestUintsToAddressHigh(t *testing.T) {
	testAddr := net.ParseIP("8000:0000:0000:0000:8000:0000:0000:0000")
	baseUint := uint64(1) << 63
	resultAddr := UintsToAddress(baseUint, baseUint)
	assert.EqualValues(t, testAddr, *resultAddr)
}
