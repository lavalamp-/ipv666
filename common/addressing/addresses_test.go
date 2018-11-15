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
