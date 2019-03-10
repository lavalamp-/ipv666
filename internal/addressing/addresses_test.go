package addressing

import (
	"testing"
	"net"
	"github.com/stretchr/testify/assert"
)

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
