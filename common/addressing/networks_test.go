package addressing

import (
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
