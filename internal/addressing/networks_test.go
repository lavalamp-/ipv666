package addressing

import (
	"math"
	"testing"
	"github.com/stretchr/testify/assert"
	"net"
)

func TestGetNetworkFromUintsZeroLength(t *testing.T) {
	addrBytes := [2]uint64{65535, 65535}
	network := GetNetworkFromUints(addrBytes, 0)
	ones, _ := network.Mask.Size()
	assert.EqualValues(t, 0, ones)
}

func TestGetNetworkFromUintsMaxLength(t *testing.T) {
	addrBytes := [2]uint64{65535, 65535}
	network := GetNetworkFromUints(addrBytes, 128)
	ones, _ := network.Mask.Size()
	assert.EqualValues(t, 128, ones)
}

func TestGetNetworkFromUintsLeftBytesMaxLength(t *testing.T) {
	addrBytes := [2]uint64{65535, 65535}
	network := GetNetworkFromUints(addrBytes, 128)
	expected := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xff, 0xff}
	assert.ElementsMatch(t, expected, network.IP[:8])
}

func TestGetNetworkFromUintsRightBytesMaxLength(t *testing.T) {
	addrBytes := [2]uint64{65535, 65535}
	network := GetNetworkFromUints(addrBytes, 128)
	expected := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xff, 0xff}
	assert.ElementsMatch(t, expected, network.IP[8:])
}

func TestGetNetworkFromUintsLeftBytesMidLength(t *testing.T) {
	addrBytes := [2]uint64{65535, 65535}
	network := GetNetworkFromUints(addrBytes, 64)
	expected := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xff, 0xff}
	assert.ElementsMatch(t, expected, network.IP[:8])
}

func TestGetNetworkFromUintsRightBytesMidLength(t *testing.T) {
	addrBytes := [2]uint64{65535, 65535}
	network := GetNetworkFromUints(addrBytes, 64)
	expected := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	assert.ElementsMatch(t, expected, network.IP[8:])
}

func TestGetNetworkFromUintsLeftBytesNoLength(t *testing.T) {
	addrBytes := [2]uint64{65535, 65535}
	network := GetNetworkFromUints(addrBytes, 0)
	expected := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	assert.ElementsMatch(t, expected, network.IP[:8])
}

func TestGetNetworkFromUintsRightBytesNoLength(t *testing.T) {
	addrBytes := [2]uint64{65535, 65535}
	network := GetNetworkFromUints(addrBytes, 0)
	expected := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	assert.ElementsMatch(t, expected, network.IP[8:])
}

func TestGetBorderAddressesFromNetworkMinMaskBase(t *testing.T) {
	_, network, _ := net.ParseCIDR("ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff/0")
	base, _ := GetBorderAddressesFromNetwork(network)
	expected := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	assert.ElementsMatch(t, expected, *base)
}

func TestGetBorderAddressesFromNetworkMinMaskTop(t *testing.T) {
	_, network, _ := net.ParseCIDR("ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff/0")
	_, top := GetBorderAddressesFromNetwork(network)
	expected := []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
	assert.ElementsMatch(t, expected, *top)
}

func TestGetBorderAddressesFromNetworkMaxMaskBase(t *testing.T) {
	_, network, _ := net.ParseCIDR("ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff/128")
	base, _ := GetBorderAddressesFromNetwork(network)
	expected := []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
	assert.ElementsMatch(t, expected, *base)
}

func TestGetBorderAddressesFromNetworkMaxMaskTop(t *testing.T) {
	_, network, _ := net.ParseCIDR("ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff/128")
	_, top := GetBorderAddressesFromNetwork(network)
	expected := []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
	assert.ElementsMatch(t, expected, *top)
}

func TestGetBorderAddressesFromNetworkMidMaskBase(t *testing.T) {
	_, network, _ := net.ParseCIDR("ffff:0000:ffff:ffff:ffff:ffff:ffff:ffff/64")
	base, _ := GetBorderAddressesFromNetwork(network)
	expected := []byte{0xff, 0xff, 0x00, 0x00, 0xff, 0xff, 0xff, 0xff, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	assert.ElementsMatch(t, expected, *base)
}

func TestGetBorderAddressesFromNetworkMidMaskTop(t *testing.T) {
	_, network, _ := net.ParseCIDR("ffff:0000:ffff:ffff:ffff:ffff:ffff:ffff/64")
	_, top := GetBorderAddressesFromNetwork(network)
	expected := []byte{0xff, 0xff, 0x00, 0x00, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
	assert.ElementsMatch(t, expected, *top)
}

func TestNetworkToUintsZeroMaskLowerFirst(t *testing.T) {
	_, network, _ := net.ParseCIDR("8000:0000:0000:0001:8000:0000:0000:0001/0")
	lowerFirst, _, _, _ := NetworkToUints(network)
	assert.EqualValues(t, 0, lowerFirst)
}

func TestNetworkToUintsZeroMaskLowerSecond(t *testing.T) {
	_, network, _ := net.ParseCIDR("8000:0000:0000:0001:8000:0000:0000:0001/0")
	_, lowerSecond, _, _ := NetworkToUints(network)
	assert.EqualValues(t, 0, lowerSecond)
}

func TestNetworkToUintsZeroMaskUpperFirst(t *testing.T) {
	_, network, _ := net.ParseCIDR("8000:0000:0000:0001:8000:0000:0000:0001/0")
	_, _, upperFirst, _ := NetworkToUints(network)
	assert.EqualValues(t, 0, upperFirst)
}

func TestNetworkToUintsZeroMaskUpperSecond(t *testing.T) {
	_, network, _ := net.ParseCIDR("8000:0000:0000:0001:8000:0000:0000:0001/0")
	_, _, _, upperSecond := NetworkToUints(network)
	assert.EqualValues(t, 0, upperSecond)
}

func TestNetworkToUintsFullMaskLowerFirst(t *testing.T) {
	_, network, _ := net.ParseCIDR("8000:0000:0000:0001:8000:0000:0000:0001/128")
	lowerFirst, _, _, _ := NetworkToUints(network)
	assert.EqualValues(t, (uint(1) << 63) + 1, lowerFirst)
}

func TestNetworkToUintsFullMaskLowerSecond(t *testing.T) {
	_, network, _ := net.ParseCIDR("8000:0000:0000:0001:8000:0000:0000:0001/128")
	_, lowerSecond, _, _ := NetworkToUints(network)
	assert.EqualValues(t, (uint(1) << 63) + 1, lowerSecond)
}

func TestNetworkToUintsFullMaskUpperFirst(t *testing.T) {
	_, network, _ := net.ParseCIDR("8000:0000:0000:0001:8000:0000:0000:0001/128")
	_, _, upperFirst, _ := NetworkToUints(network)
	assert.EqualValues(t, (uint(1) << 63) + 1, upperFirst)
}

func TestNetworkToUintsFullMaskUpperSecond(t *testing.T) {
	_, network, _ := net.ParseCIDR("8000:0000:0000:0001:8000:0000:0000:0001/128")
	_, _, _, upperSecond := NetworkToUints(network)
	assert.EqualValues(t, (uint(1) << 63) + 1, upperSecond)
}

func TestNetworkToUints32MaskLowerFirst(t *testing.T) {
	_, network, _ := net.ParseCIDR("8000:0000:0000:0001:8000:0000:0000:0001/32")
	lowerFirst, _, _, _ := NetworkToUints(network)
	assert.EqualValues(t, uint(1) << 63, lowerFirst)
}

func TestNetworkToUints32MaskLowerSecond(t *testing.T) {
	_, network, _ := net.ParseCIDR("8000:0000:0000:0001:8000:0000:0000:0001/32")
	_, lowerSecond, _, _ := NetworkToUints(network)
	assert.EqualValues(t, 0, lowerSecond)
}

func TestNetworkToUints32MaskUpperFirst(t *testing.T) {
	_, network, _ := net.ParseCIDR("8000:0000:0000:0001:8000:0000:0000:0001/32")
	_, _, upperFirst, _ := NetworkToUints(network)
	assert.EqualValues(t, (uint(1) << 63) ^ uint(math.Pow(2, 32) - 1), upperFirst)
}

func TestNetworkToUints32MaskUpperSecond(t *testing.T) {
	_, network, _ := net.ParseCIDR("8000:0000:0000:0001:8000:0000:0000:0001/32")
	_, _, _, upperSecond := NetworkToUints(network)
	assert.EqualValues(t, ^uint(0), upperSecond)
}

func TestNetworkToUints64MaskLowerFirst(t *testing.T) {
	_, network, _ := net.ParseCIDR("8000:0000:0000:0001:8000:0000:0000:0001/64")
	lowerFirst, _, _, _ := NetworkToUints(network)
	assert.EqualValues(t, (uint(1) << 63) + 1, lowerFirst)
}

func TestNetworkToUints64MaskLowerSecond(t *testing.T) {
	_, network, _ := net.ParseCIDR("8000:0000:0000:0001:8000:0000:0000:0001/64")
	_, lowerSecond, _, _ := NetworkToUints(network)
	assert.EqualValues(t, 0, lowerSecond)
}

func TestNetworkToUints64MaskUpperFirst(t *testing.T) {
	_, network, _ := net.ParseCIDR("8000:0000:0000:0001:8000:0000:0000:0001/64")
	_, _, upperFirst, _ := NetworkToUints(network)
	assert.EqualValues(t, (uint(1) << 63) + 1, upperFirst)
}

func TestNetworkToUints64MaskUpperSecond(t *testing.T) {
	_, network, _ := net.ParseCIDR("8000:0000:0000:0001:8000:0000:0000:0001/64")
	_, _, _, upperSecond := NetworkToUints(network)
	assert.EqualValues(t, ^uint(0), upperSecond)
}

func TestNetworkToUints96MaskLowerFirst(t *testing.T) {
	_, network, _ := net.ParseCIDR("8000:0000:0000:0001:8000:0000:0000:0001/96")
	lowerFirst, _, _, _ := NetworkToUints(network)
	assert.EqualValues(t, (uint(1) << 63) + 1, lowerFirst)
}

func TestNetworkToUints96MaskLowerSecond(t *testing.T) {
	_, network, _ := net.ParseCIDR("8000:0000:0000:0001:8000:0000:0000:0001/96")
	_, lowerSecond, _, _ := NetworkToUints(network)
	assert.EqualValues(t, uint(1) << 63, lowerSecond)
}

func TestNetworkToUints96MaskUpperFirst(t *testing.T) {
	_, network, _ := net.ParseCIDR("8000:0000:0000:0001:8000:0000:0000:0001/96")
	_, _, upperFirst, _ := NetworkToUints(network)
	assert.EqualValues(t, (uint(1) << 63) + 1, upperFirst)
}

func TestNetworkToUints96MaskUpperSecond(t *testing.T) {
	_, network, _ := net.ParseCIDR("8000:0000:0000:0001:8000:0000:0000:0001/96")
	_, _, _, upperSecond := NetworkToUints(network)
	assert.EqualValues(t, (uint(1) << 63) ^ uint(math.Pow(2, 32) - 1), upperSecond)
}

func TestNetworkToUints63MaskLowerFirst(t *testing.T) {
	_, network, _ := net.ParseCIDR("8000:0000:0000:0001:8000:0000:0000:0001/63")
	lowerFirst, _, _, _ := NetworkToUints(network)
	assert.EqualValues(t, uint(1) << 63, lowerFirst)
}

func TestNetworkToUints63MaskLowerSecond(t *testing.T) {
	_, network, _ := net.ParseCIDR("8000:0000:0000:0001:8000:0000:0000:0001/63")
	_, lowerSecond, _, _ := NetworkToUints(network)
	assert.EqualValues(t, 0, lowerSecond)
}

func TestNetworkToUints63MaskUpperFirst(t *testing.T) {
	_, network, _ := net.ParseCIDR("8000:0000:0000:0001:8000:0000:0000:0001/63")
	_, _, upperFirst, _ := NetworkToUints(network)
	assert.EqualValues(t, (uint(1) << 63) + 1, upperFirst)
}

func TestNetworkToUints63MaskUpperSecond(t *testing.T) {
	_, network, _ := net.ParseCIDR("8000:0000:0000:0001:8000:0000:0000:0001/63")
	_, _, _, upperSecond := NetworkToUints(network)
	assert.EqualValues(t, ^uint(0), upperSecond)
}

func TestNetworkToUints65MaskLowerFirst(t *testing.T) {
	_, network, _ := net.ParseCIDR("8000:0000:0000:0001:8000:0000:0000:0001/65")
	lowerFirst, _, _, _ := NetworkToUints(network)
	assert.EqualValues(t, (uint(1) << 63) + 1, lowerFirst)
}

func TestNetworkToUints65MaskLowerSecond(t *testing.T) {
	_, network, _ := net.ParseCIDR("8000:0000:0000:0001:8000:0000:0000:0001/65")
	_, lowerSecond, _, _ := NetworkToUints(network)
	assert.EqualValues(t, uint(1) << 63, lowerSecond)
}

func TestNetworkToUints65MaskUpperFirst(t *testing.T) {
	_, network, _ := net.ParseCIDR("8000:0000:0000:0001:8000:0000:0000:0001/65")
	_, _, upperFirst, _ := NetworkToUints(network)
	assert.EqualValues(t, (uint(1) << 63) + 1, upperFirst)
}

func TestNetworkToUints65MaskUpperSecond(t *testing.T) {
	_, network, _ := net.ParseCIDR("8000:0000:0000:0001:8000:0000:0000:0001/65")
	_, _, _, upperSecond := NetworkToUints(network)
	assert.EqualValues(t, ^uint(0), upperSecond)
}

func TestNetworkToUints1MaskLowerFirst(t *testing.T) {
	_, network, _ := net.ParseCIDR("8000:0000:0000:0001:8000:0000:0000:0001/1")
	lowerFirst, _, _, _ := NetworkToUints(network)
	assert.EqualValues(t, uint(1) << 63, lowerFirst)
}

func TestNetworkToUints1MaskLowerSecond(t *testing.T) {
	_, network, _ := net.ParseCIDR("8000:0000:0000:0001:8000:0000:0000:0001/1")
	_, lowerSecond, _, _ := NetworkToUints(network)
	assert.EqualValues(t, uint(0), lowerSecond)
}

func TestNetworkToUints1MaskUpperFirst(t *testing.T) {
	_, network, _ := net.ParseCIDR("8000:0000:0000:0001:8000:0000:0000:0001/1")
	_, _, upperFirst, _ := NetworkToUints(network)
	assert.EqualValues(t, ^uint(0), upperFirst)
}

func TestNetworkToUints1MaskUpperSecond(t *testing.T) {
	_, network, _ := net.ParseCIDR("8000:0000:0000:0001:8000:0000:0000:0001/1")
	_, _, _, upperSecond := NetworkToUints(network)
	assert.EqualValues(t, ^uint(0), upperSecond)
}

func TestNetworkToUints127MaskLowerFirst(t *testing.T) {
	_, network, _ := net.ParseCIDR("8000:0000:0000:0001:8000:0000:0000:0001/127")
	lowerFirst, _, _, _ := NetworkToUints(network)
	assert.EqualValues(t, (uint(1) << 63) + 1, lowerFirst)
}

func TestNetworkToUints127MaskLowerSecond(t *testing.T) {
	_, network, _ := net.ParseCIDR("8000:0000:0000:0001:8000:0000:0000:0001/127")
	_, lowerSecond, _, _ := NetworkToUints(network)
	assert.EqualValues(t, uint(1) << 63, lowerSecond)
}

func TestNetworkToUints127MaskUpperFirst(t *testing.T) {
	_, network, _ := net.ParseCIDR("8000:0000:0000:0001:8000:0000:0000:0001/127")
	_, _, upperFirst, _ := NetworkToUints(network)
	assert.EqualValues(t, (uint(1) << 63) + 1, upperFirst)
}

func TestNetworkToUints127MaskUpperSecond(t *testing.T) {
	_, network, _ := net.ParseCIDR("8000:0000:0000:0001:8000:0000:0000:0001/127")
	_, _, _, upperSecond := NetworkToUints(network)
	assert.EqualValues(t, (uint(1) << 63) + 1, upperSecond)
}
