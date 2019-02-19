package modeling

import (
	"github.com/magiconair/properties/assert"
	"net"
	"testing"
)

func getTestRange() *GenRange {
	ip1 := net.ParseIP("8000:0000:0001:0001:8000:0001:0000:0001")
	ip2 := net.ParseIP("8000:0001:0000:0001:8000:0000:0000:0001")
	ip3 := net.ParseIP("8000:0000:0000:0001:8000:0001:0000:0001")
	ip4 := net.ParseIP("8000:0001:0000:0001:8000:0000:0001:0001")
	toReturn := newGenRange(&ip1)
	toReturn.AddIP(&ip2)
	toReturn.AddIP(&ip3)
	toReturn.AddIP(&ip4)
	return toReturn
}

func TestGenRange_GetMaskFirstMask(t *testing.T) {
	testRange := getTestRange()
	expected := bytesToUint([]byte{ 0xff, 0xff, 0xff, 0xf0, 0xff, 0xf0, 0xff, 0xff })
	assert.Equal(t, expected, testRange.GetMask().FirstMask)
}

func TestGenRange_GetMaskFirstExpected(t *testing.T) {
	testRange := getTestRange()
	expected := bytesToUint([]byte{ 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01 })
	assert.Equal(t, expected, testRange.GetMask().FirstExpected)
}

func TestGenRange_GetMaskFirstMin(t *testing.T) {
	testRange := getTestRange()
	expected := bytesToUint([]byte{ 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01 })
	assert.Equal(t, expected, testRange.GetMask().FirstMin)
}

func TestGenRange_GetMaskFirstMax(t *testing.T) {
	testRange := getTestRange()
	expected := bytesToUint([]byte{ 0x80, 0x00, 0x00, 0x0f, 0x00, 0x0f, 0x00, 0x01 })
	assert.Equal(t, expected, testRange.GetMask().FirstMax)
}

func TestGenRange_GetMaskSecondMask(t *testing.T) {
	testRange := getTestRange()
	expected := bytesToUint([]byte{ 0xff, 0xff, 0xff, 0xf0, 0xff, 0xf0, 0xff, 0xff })
	assert.Equal(t, expected, testRange.GetMask().SecondMask)
}

func TestGenRange_GetMaskSecondExpected(t *testing.T) {
	testRange := getTestRange()
	expected := bytesToUint([]byte{ 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01 })
	assert.Equal(t, expected, testRange.GetMask().SecondExpected)
}

func TestGenRange_GetMaskSecondMin(t *testing.T) {
	testRange := getTestRange()
	expected := bytesToUint([]byte{ 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01 })
	assert.Equal(t, expected, testRange.GetMask().SecondMin)
}

func TestGenRange_GetMaskSecondMax(t *testing.T) {
	testRange := getTestRange()
	expected := bytesToUint([]byte{ 0x80, 0x00, 0x00, 0x0f, 0x00, 0x0f, 0x00, 0x01 })
	assert.Equal(t, expected, testRange.GetMask().SecondMax)
}