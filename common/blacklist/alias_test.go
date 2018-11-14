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

func getTestAddress2() (*net.IP) {
	toReturn := net.ParseIP("aaaa:aaaa:aaaa:aaaa:aaaa:aaaa:aaaa:aaaa")
	return &toReturn
}

func getTestAddresses() ([]*net.IP) {
	return []*net.IP{getTestAddress(), getTestAddress2()}
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

func TestAliasCheckState_GenerateTestAddressNotNil(t *testing.T) {
	acs, _ := NewAliasCheckState(getTestAddress(), 16, 127)
	acs.GenerateTestAddress()
	assert.NotNil(t, acs.GetTestAddr())
}

func TestAliasCheckState_GenerateTestAddress(t *testing.T) {
	testAddr := net.ParseIP("aaaa:aaaa:aaaa:aaaa:aaaa:aaaa:aaaa:aaaa")
	acs, _ := NewAliasCheckState(&testAddr, 16, 127)
	acs.GenerateTestAddress()
	expectedAddr := net.ParseIP("aaaa:aaaa:aaaa:aaaa:aa55:5555:5555:5554")
	assert.Equal(t, expectedAddr.String(), acs.GetTestAddr().String())
}

func TestAliasCheckState_UpdateFoundPosition(t *testing.T) {
	testAddr := net.ParseIP("aaaa:aaaa:aaaa:aaaa:aaaa:aaaa:aaaa:aaaa")
	acs, _ := NewAliasCheckState(&testAddr, 16, 127)
	acs.GenerateTestAddress()
	foundAddrs := getFoundAddrsMap([]*net.IP{acs.GetTestAddr()}, 10)
	expectedPosition := acs.GetLeftTestIndex()
	acs.Update(foundAddrs)
	assert.EqualValues(t, expectedPosition, acs.GetRight())
}

func TestAliasCheckState_UpdateNotFoundPosition(t *testing.T) {
	testAddr := net.ParseIP("aaaa:aaaa:aaaa:aaaa:aaaa:aaaa:aaaa:aaaa")
	acs, _ := NewAliasCheckState(&testAddr, 16, 127)
	acs.GenerateTestAddress()
	foundAddrs := getFoundAddrsMap([]*net.IP{}, 10)
	expectedPosition := acs.GetLeftTestIndex()
	acs.Update(foundAddrs)
	assert.EqualValues(t, expectedPosition, acs.GetLeft())
}

func TestAliasCheckState_UpdateFoundNotFinished(t *testing.T) {
	testAddr := net.ParseIP("aaaa:aaaa:aaaa:aaaa:aaaa:aaaa:aaaa:aaaa")
	acs, _ := NewAliasCheckState(&testAddr, 16, 127)
	acs.GenerateTestAddress()
	foundAddrs := getFoundAddrsMap([]*net.IP{acs.GetTestAddr()}, 10)
	acs.Update(foundAddrs)
	assert.False(t, acs.GetFound())
}

func TestAliasCheckState_UpdateFoundFinished(t *testing.T) {
	testAddr := net.ParseIP("aaaa:aaaa:aaaa:aaaa:aaaa:aaaa:aaaa:aaaa")
	acs, _ := NewAliasCheckState(&testAddr, 16, 18)
	acs.GenerateTestAddress()
	foundAddrs := getFoundAddrsMap([]*net.IP{acs.GetTestAddr()}, 10)
	acs.Update(foundAddrs)
	assert.True(t, acs.GetFound())
}

func TestAliasCheckState_UpdateEmptyAddr(t *testing.T) {
	testAddr := net.ParseIP("aaaa:aaaa:aaaa:aaaa:aaaa:aaaa:aaaa:aaaa")
	acs, _ := NewAliasCheckState(&testAddr, 16, 18)
	acs.GenerateTestAddress()
	foundAddrs := getFoundAddrsMap([]*net.IP{acs.GetTestAddr()}, 10)
	acs.Update(foundAddrs)
	assert.Nil(t, acs.GetTestAddr())
}

func TestAliasCheckState_GetAliasedNetworkError(t *testing.T) {
	acs, _ := NewAliasCheckState(getTestAddress(), 0, 127)
	_, err := acs.GetAliasedNetwork()
	assert.NotNil(t, err)
}

func TestAliasCheckState_GetAliasedNetwork(t *testing.T) {
	testAddr := net.ParseIP("aaaa:aaaa:aaaa:aaaa:aaaa:aaaa:aaaa:aaaa")
	acs, _ := NewAliasCheckState(&testAddr, 15, 17)
	acs.GenerateTestAddress()
	foundAddrs := getFoundAddrsMap([]*net.IP{acs.GetTestAddr()}, 10)
	acs.Update(foundAddrs)
	network, _ := acs.GetAliasedNetwork()
	assert.Equal(t, "aaaa::/16", network.String())
}

func TestNewAliasCheckStatesRightError(t *testing.T) {
	_, err := NewAliasCheckStates(getTestAddresses(), 16, 128)
	assert.NotNil(t, err)
}

func TestNewAliasCheckStatesRightLessThanLeftError(t *testing.T) {
	_, err := NewAliasCheckStates(getTestAddresses(), 16, 15)
	assert.NotNil(t, err)
}

func TestAliasCheckStates_GetChecksCount(t *testing.T) {
	acs, _ := NewAliasCheckStates(getTestAddresses(), 16, 127)
	assert.EqualValues(t, 2, acs.GetChecksCount())
}

func TestAliasCheckStates_GetFoundCountInit(t *testing.T) {
	acs, _ := NewAliasCheckStates(getTestAddresses(), 16, 127)
	assert.EqualValues(t, 0, acs.GetFoundCount())
}

func TestAliasCheckStates_GetFoundCount(t *testing.T) {
	addr1 := net.ParseIP("aaaa:aaaa:aaaa:aaaa:aaaa:aaaa:aaaa:aaaa")
	addr2 := net.ParseIP("aaaa:aaaa:aaaa:aaaa:ffff:ffff:ffff:ffff")
	addrs := []*net.IP{&addr1, &addr2}
	acs, _ := NewAliasCheckStates(addrs, 95, 97)
	foundAddr1 := net.ParseIP("aaaa:aaaa:aaaa:aaaa:aaaa:aaaa:2aaa:aaaa")
	foundAddr2 := net.ParseIP("aaaa:aaaa:aaaa:aaaa:ffff:ffff:7fff:ffff")
	foundAddrs := getFoundAddrsMap([]*net.IP{&foundAddr1, &foundAddr2}, 10)
	acs.GenerateTestAddresses()
	acs.Update(foundAddrs)
	assert.EqualValues(t, 2, acs.GetFoundCount())
}

func TestAliasCheckStates_GetAllFoundInit(t *testing.T) {
	acs, _ := NewAliasCheckStates(getTestAddresses(), 16, 127)
	assert.False(t, acs.GetAllFound())
}

func TestAliasCheckStates_GetAllFoundTrue(t *testing.T) {
	addr1 := net.ParseIP("aaaa:aaaa:aaaa:aaaa:aaaa:aaaa:aaaa:aaaa")
	addr2 := net.ParseIP("aaaa:aaaa:aaaa:aaaa:ffff:ffff:ffff:ffff")
	addrs := []*net.IP{&addr1, &addr2}
	acs, _ := NewAliasCheckStates(addrs, 95, 97)
	foundAddr1 := net.ParseIP("aaaa:aaaa:aaaa:aaaa:aaaa:aaaa:2aaa:aaaa")
	foundAddr2 := net.ParseIP("aaaa:aaaa:aaaa:aaaa:ffff:ffff:7fff:ffff")
	foundAddrs := getFoundAddrsMap([]*net.IP{&foundAddr1, &foundAddr2}, 10)
	acs.GenerateTestAddresses()
	acs.Update(foundAddrs)
	assert.True(t, acs.GetAllFound())
}

func TestAliasCheckStates_GetAllFoundFalse(t *testing.T) {
	addr1 := net.ParseIP("aaaa:aaaa:aaaa:aaaa:aaaa:aaaa:aaaa:aaaa")
	addr2 := net.ParseIP("aaaa:aaaa:aaaa:aaaa:ffff:ffff:ffff:ffff")
	addrs := []*net.IP{&addr1, &addr2}
	acs, _ := NewAliasCheckStates(addrs, 10, 97)
	foundAddr1 := net.ParseIP("aaaa:aaaa:aaaa:aaaa:aaaa:aaaa:2aaa:aaaa")
	foundAddr2 := net.ParseIP("aaaa:aaaa:aaaa:aaaa:ffff:ffff:7fff:ffff")
	foundAddrs := getFoundAddrsMap([]*net.IP{&foundAddr1, &foundAddr2}, 10)
	acs.GenerateTestAddresses()
	acs.Update(foundAddrs)
	assert.False(t, acs.GetAllFound())
}

func TestAliasCheckStates_GetAliasedNetworksError(t *testing.T) {
	addr1 := net.ParseIP("aaaa:aaaa:aaaa:aaaa:aaaa:aaaa:aaaa:aaaa")
	addr2 := net.ParseIP("aaaa:aaaa:aaaa:aaaa:ffff:ffff:ffff:ffff")
	addrs := []*net.IP{&addr1, &addr2}
	acs, _ := NewAliasCheckStates(addrs, 95, 97)
	_, err := acs.GetAliasedNetworks()
	assert.NotNil(t, err)
}

func TestAliasCheckStates_GetAliasedNetworks(t *testing.T) {
	addr1 := net.ParseIP("aaaa:aaaa:aaaa:aaaa:aaaa:aaaa:aaaa:aaaa")
	addr2 := net.ParseIP("aaaa:aaaa:aaaa:aaaa:ffff:ffff:ffff:ffff")
	addrs := []*net.IP{&addr1, &addr2}
	acs, _ := NewAliasCheckStates(addrs, 95, 97)
	foundAddr1 := net.ParseIP("aaaa:aaaa:aaaa:aaaa:aaaa:aaaa:2aaa:aaaa")
	foundAddr2 := net.ParseIP("aaaa:aaaa:aaaa:aaaa:ffff:ffff:7fff:ffff")
	foundAddrs := getFoundAddrsMap([]*net.IP{&foundAddr1, &foundAddr2}, 10)
	acs.GenerateTestAddresses()
	acs.Update(foundAddrs)
	nets, _ := acs.GetAliasedNetworks()
	assert.NotNil(t, nets)
}
