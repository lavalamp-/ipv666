package blacklist

import (
	"testing"
	"net"
	"github.com/stretchr/testify/assert"
)

func getTestAddress() (*net.IP) {
	toReturn := net.ParseIP("2001:0:4137:9e76:38c1:3e16:6ac9:506a")
	return &toReturn
}

func TestNewAliasCheckStateRightError(t *testing.T) {
	_, err := NewAliasCheckState(getTestAddress(), 0, 128)
	assert.NotNil(t, err, "No error thrown when creating new AliasCheckState with right value of 128.")
}

func TestNewAliasCheckStateRightLessThanLeftError(t *testing.T) {
	_, err := NewAliasCheckState(getTestAddress(), 65, 64)
	assert.NotNil(t, err, "No error thrown when creating new AliasCheckState with left value higher than right value (65 and 64).")
}

func TestAliasCheckState_GetLeft(t *testing.T) {
	acs, _ := NewAliasCheckState(getTestAddress(), 16, 98)
	assert.EqualValues(t, acs.GetLeft(), 16)
}

func TestAliasCheckState_GetRight(t *testing.T) {
	acs, _ := NewAliasCheckState(getTestAddress(), 16, 98)
	assert.EqualValues(t, acs.GetRight(), 98)
}

func TestAliasCheckState_GetFound(t *testing.T) {
	acs, _ := NewAliasCheckState(getTestAddress(), 16, 98)
	assert.False(t, acs.GetFound())
}

func TestAliasCheckState_GetTestAddrs(t *testing.T) {
	acs, _ := NewAliasCheckState(getTestAddress(), 16, 98)
	assert.Nil(t, acs.GetTestAddr())
}

func TestAliasCheckState_GetBaseAddress(t *testing.T) {
	testAddress := getTestAddress()
	acs, _ := NewAliasCheckState(testAddress, 16, 98)
	assert.Equal(t, testAddress, acs.GetBaseAddress())
}

func TestAliasCheckState_GetMiddleOdd(t *testing.T) {
	acs, _ := NewAliasCheckState(getTestAddress(), 10, 21)
	assert.EqualValues(t, 16, acs.GetMiddle())
}

func TestAliasCheckState_GetMiddleEven(t *testing.T) {
	acs, _ := NewAliasCheckState(getTestAddress(), 10, 20)
	assert.EqualValues(t, 16, acs.GetMiddle())
}

func TestAliasCheckState_GetMiddleOneDistance(t *testing.T) {
	acs, _ := NewAliasCheckState(getTestAddress(), 10, 11)
	assert.EqualValues(t, 11, acs.GetMiddle())
}

func TestAliasCheckState_GetMiddleMaxDistance(t *testing.T) {
	acs, _ := NewAliasCheckState(getTestAddress(), 0, 127)
	assert.EqualValues(t, 64, acs.GetMiddle())
}

func TestAliasCheckState_GetTestDistanceOne(t *testing.T) {
	acs, _ := NewAliasCheckState(getTestAddress(), 10, 11)
	assert.EqualValues(t, 1, acs.GetTestDistance())
}

func TestAliasCheckState_GetTestDistanceMax(t *testing.T) {
	acs, _ := NewAliasCheckState(getTestAddress(), 0, 127)
	assert.EqualValues(t, 127, acs.GetTestDistance())
}

func TestAliasCheckState_GetTestDistance(t *testing.T) {
	acs, _ := NewAliasCheckState(getTestAddress(), 16, 96)
	assert.EqualValues(t,80, acs.GetTestDistance())
}

func TestAliasCheckState_GetTestBitCountOne(t *testing.T) {
	acs, _ := NewAliasCheckState(getTestAddress(), 10, 11)
	assert.EqualValues(t,1, acs.GetTestBitCount())
}

func TestAliasCheckState_GetTestBitCountMax(t *testing.T) {
	acs, _ := NewAliasCheckState(getTestAddress(), 0, 127)
	assert.EqualValues(t, 64, acs.GetTestBitCount())
}

func TestAliasCheckState_GetTestBitCount(t *testing.T) {
	acs, _ := NewAliasCheckState(getTestAddress(), 0, 7)
	assert.EqualValues(t, 4, acs.GetTestBitCount())
}