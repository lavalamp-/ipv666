package blacklist

import (
	"testing"
	"net"
	"github.com/stretchr/testify/assert"
	"github.com/lavalamp-/ipv666/common/addressing"
)

func getTestAddress() (*net.IP) {
	toReturn := net.ParseIP("2001:0:4137:9e76:38c1:3e16:6ac9:506a")
	return &toReturn
}

func getFoundAddrsMap(toInclude []*net.IP, addrCount int) (map[string]*Empty) {
	toReturn := make(map[string]*Empty)
	placeholder := &Empty{}
	for _, curInclude := range toInclude {
		toReturn[curInclude.String()] = placeholder
	}
	for len(toReturn) < addrCount {
		newAddr := addressing.GenerateRandomAddress()
		toReturn[newAddr.String()] = placeholder
	}
	return toReturn
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

func TestAliasCheckState_GetLeftTestIndexZero(t *testing.T) {
	acs, _ := NewAliasCheckState(getTestAddress(), 16, 17)
	assert.EqualValues(t, acs.GetRight(), acs.GetLeftTestIndex())
}

func TestAliasCheckState_GetLeftTestIndexOne(t *testing.T) {
	acs, _ := NewAliasCheckState(getTestAddress(), 16, 18)
	assert.EqualValues(t, 17, acs.GetLeftTestIndex())
}

func TestAliasCheckState_GetLeftTestIndexOdd(t *testing.T) {
	acs, _ := NewAliasCheckState(getTestAddress(), 10, 19)
	assert.EqualValues(t, 15, acs.GetLeftTestIndex())
}

func TestAliasCheckState_GetLeftTestIndexEven(t *testing.T) {
	acs, _ := NewAliasCheckState(getTestAddress(), 10, 20)
	assert.EqualValues(t, 15, acs.GetLeftTestIndex())
}

func TestAliasCheckState_GetLeftTestIndex(t *testing.T) {
	acs, _ := NewAliasCheckState(getTestAddress(), 0, 127)
	assert.EqualValues(t, 64, acs.GetLeftTestIndex())
}

func TestAliasCheckState_GetRightTestIndexZero(t *testing.T) {
	acs, _ := NewAliasCheckState(getTestAddress(), 16, 17)
	assert.EqualValues(t, acs.GetRight(), acs.GetRightTestIndex())
}

func TestAliasCheckState_GetRightTestIndexOne(t *testing.T) {
	acs, _ := NewAliasCheckState(getTestAddress(), 16, 18)
	assert.EqualValues(t, 17, acs.GetRightTestIndex())
}

func TestAliasCheckState_GetRightTestIndexOdd(t *testing.T) {
	acs, _ := NewAliasCheckState(getTestAddress(), 10, 19)
	assert.EqualValues(t, 18, acs.GetRightTestIndex())
}

func TestAliasCheckState_GetRightTestIndexEven(t *testing.T) {
	acs, _ := NewAliasCheckState(getTestAddress(), 10, 20)
	assert.EqualValues(t, 19, acs.GetRightTestIndex())
}

func TestAliasCheckState_GetRightTestIndex(t *testing.T) {
	acs, _ := NewAliasCheckState(getTestAddress(), 0, 127)
	assert.EqualValues(t, 126, acs.GetRightTestIndex())
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
	assert.EqualValues(t, 63, acs.GetTestBitCount())
}

func TestAliasCheckState_GetTestBitCount(t *testing.T) {
	acs, _ := NewAliasCheckState(getTestAddress(), 0, 7)
	assert.EqualValues(t, 3, acs.GetTestBitCount())
}

//func TestAliasCheckState_UpdateFoundPosition(t *testing.T) {
//	addr, addrMap := getDefaultTestAddrAndMap(true)
//	acs, _ := NewAliasCheckState(addr, 0, 127)
//	middle := acs.GetMiddle()
//	acs.Update(addrMap)
//	assert.Equal(t, middle, acs.GetRight())
//}
//
//func TestAliasCheckState_UpdateNotFoundPosition(t *testing.T) {
//	addr, addrMap := getDefaultTestAddrAndMap(false)
//	acs, _ := NewAliasCheckState(addr, 0, 127)
//	middle := acs.GetMiddle()
//	acs.Update(addrMap)
//	assert.Equal(t, middle, acs.GetLeft())
//}
//
//func TestAliasCheckState_UpdateFoundNotFinished(t *testing.T) {
//	addr, addrMap := getDefaultTestAddrAndMap(true)
//	acs, _ := NewAliasCheckState(addr, 0, 127)
//	acs.Update(addrMap)
//	assert.False(t, acs.GetFound())
//}
//
//func TestAliasCheckState_UpdateFoundFinished(t *testing.T) {
//	addr, addrMap := getDefaultTestAddrAndMap(true)
//	acs, _ := NewAliasCheckState(addr, 0, 127)
//	acs.Update(addrMap)
//	assert.False(t, acs.GetFound())
//}
//
//func TestAliasCheckState_UpdateEmptyAddr(t *testing.T) {
//
//}