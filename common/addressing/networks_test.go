package addressing

import (
	"testing"
	"github.com/stretchr/testify/assert"
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
