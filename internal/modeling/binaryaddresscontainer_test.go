package modeling

import (
	"github.com/lavalamp-/ipv666/internal/addressing"
	"github.com/lavalamp-/ipv666/internal/config"
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
)

func init() {
	config.InitConfig()
}

func bytesToUint(bytes []byte) uint64 {  //TODO move helper functions to common test lib
	toReturn := uint64(0)
	for i := 0; i < 8; i++ {
		toReturn ^= uint64(bytes[i]) << uint((7 - i) * 8)
	}
	return toReturn
}

func byteListToUints(byteLists [][]byte) []uint64 {
	var toReturn []uint64
	for _, curList := range byteLists {
		toReturn = append(toReturn, bytesToUint(curList))
	}
	return toReturn
}

func getTestUints(odd bool) []uint64 {
	toReturn := []uint64{
		1,
		2,
		9,
		10,
		11,
		19,
		20,
		21,
		29,
		30,
		31,
		39,
		40,
		41,
		49,
		50,
	}
	if odd {
		toReturn = append(toReturn, 51)
	}
	return toReturn
}

func getSeekRangeTestUints() []uint64 {
	return []uint64{
		9,
		10,
		11,
		19,
		20,
		21,
		29,
		30,
		31,
		39,
		40,
		41,
		49,
		50,
	}
}

func getBinaryContainer() *BinaryAddressContainer {
	return ContainerFromAddrs(getDefaultIPs())
}

func getEmptyBinaryContainer() *BinaryAddressContainer {
	return ContainerFromAddrs([]*net.IP{})
}

func getExtendedBinaryContainer() *BinaryAddressContainer {
	return ContainerFromAddrs(getExtendedDefaultIPs())
}

func TestContainerFromAddrsReturns(t *testing.T) {
	container := getBinaryContainer()
	assert.NotNil(t, container)
}

func TestContainerFromAddrsEmpty(t *testing.T) {
	container := getEmptyBinaryContainer()
	assert.NotNil(t, container)
}

func TestContainerFromAddrsCount(t *testing.T) {
	container := getBinaryContainer()
	assert.EqualValues(t, len(getDefaultIPs()), container.Size())
}

func TestBinaryAddressContainer_AddIPCount(t *testing.T) {
	container := getBinaryContainer()
	newIP := net.ParseIP("2600:0:1:0001:0000:0000:0000:0001")
	firstCount := container.Size()
	container.AddIP(&newIP)
	assert.Equal(t, firstCount + 1, container.Size())
}

func TestBinaryAddressContainer_AddIPAdds(t *testing.T) {
	container := getBinaryContainer()
	newIP := net.ParseIP("2600:0:1:0001:0000:0000:0000:0001")
	container.AddIP(&newIP)
	assert.True(t, container.ContainsIP(&newIP))
}

func TestBinaryAddressContainer_AddIPsCount(t *testing.T) {
	container := getBinaryContainer()
	newIPs := addressing.GetIPsFromStrings([]string {
		"2600:0:1:0001:0000:0000:0000:0001",
		"2600:0:1:0001:0000:0000:0000:0002",
		"2600:0:1:0001:0000:0000:0000:0003",
		"2600:0:1:0001:0000:0000:0000:0004",
	})
	firstCount := container.Size()
	container.AddIPs(newIPs, 100)
	assert.Equal(t, firstCount + len(newIPs), container.Size())
}

func TestBinaryAddressContainer_AddIPsAdds(t *testing.T) {
	container := getBinaryContainer()
	newIPs := addressing.GetIPsFromStrings([]string {
		"2600:0:1:0001:0000:0000:0000:0001",
		"2600:0:1:0001:0000:0000:0000:0002",
		"2600:0:1:0001:0000:0000:0000:0003",
		"2600:0:1:0001:0000:0000:0000:0004",
	})
	container.AddIPs(newIPs, 100)
	for _, newIP := range newIPs {
		assert.True(t, container.ContainsIP(newIP))
	}
}

func TestBinaryAddressContainer_AddIPsAdded(t *testing.T) {
	container := getBinaryContainer()
	newIPs := addressing.GetIPsFromStrings([]string {
		"2600:0:1:0001:0000:0000:0000:0001",
		"2600:0:1:0001:0000:0000:0000:0002",
		"2600:0:1:0001:0000:0000:0000:0003",
		"2600:0:1:0001:0000:0000:0000:0004",
	})
	added, _ := container.AddIPs(newIPs, 100)
	assert.Equal(t, 4, added)
}

func TestBinaryAddressContainer_AddIPsSkipped(t *testing.T) {
	container := getBinaryContainer()
	newIPs := addressing.GetIPsFromStrings([]string {
		"2600:0:1:0001:0000:0000:0000:0001",
		"2600:0:1:0001:0000:0000:0000:0002",
		"2600:0:1:0001:0000:0000:0000:0003",
		"2600:0:1:0001:0000:0000:0000:0004",
	})
	_, skipped := container.AddIPs(newIPs, 100)
	assert.Equal(t, 0, skipped)
}

func TestBinaryAddressContainer_GetAllIPsEmptyValue(t *testing.T) {
	container := getEmptyBinaryContainer()
	assert.Empty(t, container.GetAllIPs())
}

func TestBinaryAddressContainer_GetAllIPsNotEmpty(t *testing.T) {
	container := getBinaryContainer()
	assert.NotEmpty(t, container.GetAllIPs())
}

func TestBinaryAddressContainer_GetAllIPsMatchesCount(t *testing.T) {
	container := getBinaryContainer()
	assert.EqualValues(t, container.Size(), len(container.GetAllIPs()))
}

func TestBinaryAddressContainer_GetAllIPsMatchesInput(t *testing.T) {
	container := getBinaryContainer()
	inputAddrs := getDefaultIPs()
	assert.ElementsMatch(t, inputAddrs, container.GetAllIPs())
}

func TestBinaryAddressContainer_GetIPsInRangeNoError(t *testing.T) {
	container := getBinaryContainer()
	_, testNet, _ := net.ParseCIDR("2600:0:1:0000:0000:0000:0000:0000/0")
	_, err := container.GetIPsInRange(testNet)
	assert.Nil(t, err)
}

func TestBinaryAddressContainer_GetIPsInRange0(t *testing.T) {
	container := getBinaryContainer()
	_, testNet, _ := net.ParseCIDR("2600:0:1:0000:0000:0000:0000:0000/0")
	addrs, _ := container.GetIPsInRange(testNet)
	assert.Equal(t, 64, len(addrs))
}

func TestBinaryAddressContainer_GetIPsInRange32(t *testing.T) {
	container := getBinaryContainer()
	_, testNet, _ := net.ParseCIDR("2600:0:1:0000:0000:0000:0000:0000/32")
	addrs, _ := container.GetIPsInRange(testNet)
	assert.Equal(t, 64, len(addrs))
}

func TestBinaryAddressContainer_GetIPsInRange64(t *testing.T) {
	container := getBinaryContainer()
	_, testNet, _ := net.ParseCIDR("2600:0:1:0000:0000:0000:0000:0000/64")
	addrs, _ := container.GetIPsInRange(testNet)
	assert.Equal(t, 64, len(addrs))
}

func TestBinaryAddressContainer_GetIPsInRange96(t *testing.T) {
	container := getBinaryContainer()
	_, testNet, _ := net.ParseCIDR("2600:0:1:0000:0000:0000:0000:0000/96")
	addrs, _ := container.GetIPsInRange(testNet)
	assert.Equal(t, 16, len(addrs))
}

func TestBinaryAddressContainer_GetIPsInRange128(t *testing.T) {
	container := getBinaryContainer()
	_, testNet, _ := net.ParseCIDR("2600:0:1:0000:0000:0000:0000:0000/128")
	addrs, _ := container.GetIPsInRange(testNet)
	assert.Equal(t, 1, len(addrs))
}

func TestBinaryAddressContainer_CountIPsInRangeNoError(t *testing.T) {
	container := getBinaryContainer()
	_, testNet, _ := net.ParseCIDR("2600:0:1:0000:0000:0000:0000:0000/0")
	_, err := container.CountIPsInRange(testNet)
	assert.Nil(t, err)
}

func TestBinaryAddressContainer_CountIPsInRange0(t *testing.T) {
	container := getBinaryContainer()
	_, testNet, _ := net.ParseCIDR("2600:0:1:0000:0000:0000:0000:0000/0")
	count, _ := container.CountIPsInRange(testNet)
	assert.EqualValues(t, 64, count)
}

func TestBinaryAddressContainer_CountIPsInRange32(t *testing.T) {
	container := getBinaryContainer()
	_, testNet, _ := net.ParseCIDR("2600:0:1:0000:0000:0000:0000:0000/32")
	count, _ := container.CountIPsInRange(testNet)
	assert.EqualValues(t, 64, count)
}

func TestBinaryAddressContainer_CountIPsInRange64(t *testing.T) {
	container := getBinaryContainer()
	_, testNet, _ := net.ParseCIDR("2600:0:1:0000:0000:0000:0000:0000/64")
	count, _ := container.CountIPsInRange(testNet)
	assert.EqualValues(t, 64, count)
}

func TestBinaryAddressContainer_CountIPsInRange96(t *testing.T) {
	container := getBinaryContainer()
	_, testNet, _ := net.ParseCIDR("2600:0:1:0000:0000:0000:0000:0000/96")
	count, _ := container.CountIPsInRange(testNet)
	assert.EqualValues(t, 16, count)
}

func TestBinaryAddressContainer_CountIPsInRange128(t *testing.T) {
	container := getBinaryContainer()
	_, testNet, _ := net.ParseCIDR("2600:0:1:0000:0000:0000:0000:0000/128")
	count, _ := container.CountIPsInRange(testNet)
	assert.EqualValues(t, 1, count)
}

func TestBinaryAddressContainer_ExtendedGetIPsInRangeNoError(t *testing.T) {
	container := getExtendedBinaryContainer()
	_, testNet, _ := net.ParseCIDR("2600:0:1:0000:0000:0000:0000:0000/0")
	_, err := container.GetIPsInRange(testNet)
	assert.Nil(t, err)
}

func TestBinaryAddressContainer_ExtendedGetIPsInRange0(t *testing.T) {
	container := getExtendedBinaryContainer()
	_, testNet, _ := net.ParseCIDR("2600:0:1:0000:0000:0000:0000:0000/0")
	addrs, _ := container.GetIPsInRange(testNet)
	assert.Equal(t, 128, len(addrs))
}

func TestBinaryAddressContainer_ExtendedGetIPsInRange32(t *testing.T) {
	container := getExtendedBinaryContainer()
	_, testNet, _ := net.ParseCIDR("2600:0:1:0000:0000:0000:0000:0000/32")
	addrs, _ := container.GetIPsInRange(testNet)
	assert.Equal(t, 128, len(addrs))
}

func TestBinaryAddressContainer_ExtendedGetIPsInRange64(t *testing.T) {
	container := getExtendedBinaryContainer()
	_, testNet, _ := net.ParseCIDR("2600:0:1:0000:0000:0000:0000:0000/64")
	addrs, _ := container.GetIPsInRange(testNet)
	assert.Equal(t, 64, len(addrs))
}

func TestBinaryAddressContainer_ExtendedGetIPsInRange96(t *testing.T) {
	container := getExtendedBinaryContainer()
	_, testNet, _ := net.ParseCIDR("2600:0:1:0000:0000:0000:0000:0000/96")
	addrs, _ := container.GetIPsInRange(testNet)
	assert.Equal(t, 16, len(addrs))
}

func TestBinaryAddressContainer_ExtendedGetIPsInRange128(t *testing.T) {
	container := getExtendedBinaryContainer()
	_, testNet, _ := net.ParseCIDR("2600:0:1:0000:0000:0000:0000:0000/128")
	addrs, _ := container.GetIPsInRange(testNet)
	assert.Equal(t, 1, len(addrs))
}



//func TestBinaryAddressContainer_SaveCreatesFile(t *testing.T) {
//	container := getBinaryContainer()
//	addrTree.Save("/tmp/tester")
//	_, err := os.Stat("/tmp/tester")
//	exists := !os.IsNotExist(err)
//	assert.True(t, exists)
//	os.Remove("/tmp/tester")
//}
//
//func TestBinaryAddressContainer_SaveNoError(t *testing.T) {
//	container := getBinaryContainer()
//	err := addrTree.Save("/tmp/tester")
//	os.Remove("/tmp/tester")
//	assert.Nil(t, err)
//}
//
//func TestLoadAddressTreeFromFileNoError(t *testing.T) {
//	container := getBinaryContainer()
//	addrTree.Save("/tmp/tester")
//	_, err := LoadAddressTreeFromFile("/tmp/tester")
//	os.Remove("/tmp/tester")
//	assert.Nil(t, err)
//}
//
//func TestLoadAddressTreeFromFileLoads(t *testing.T) {
//	container := getBinaryContainer()
//	addrTree.Save("/tmp/tester")
//	newTree, _ := LoadAddressTreeFromFile("/tmp/tester")
//	os.Remove("/tmp/tester")
//	assert.NotNil(t, newTree)
//}
//
//func TestLoadAddressTreeFromFileContent(t *testing.T) {
//	container := getBinaryContainer()
//	addrTree.Save("/tmp/tester")
//	newTree, _ := LoadAddressTreeFromFile("/tmp/tester")
//	os.Remove("/tmp/tester")
//	assert.ElementsMatch(t, getDefaultIPs(), newTree.GetAllIPs())
//}

func TestBinaryAddressContainer_ContainsIPEmpty(t *testing.T) {
	container := getEmptyBinaryContainer()
	ip := net.ParseIP("2600:0:1:0000:0000:0000:0000:0000")
	assert.False(t, container.ContainsIP(&ip))
}

func TestBinaryAddressContainer_ContainsIPFalse(t *testing.T) {
	container := getBinaryContainer()
	ip := net.ParseIP("2700:0:1:0000:0000:0000:0000:0000")
	assert.False(t, container.ContainsIP(&ip))
}

func TestBinaryAddressContainer_ContainsIPTrue(t *testing.T) {
	container := getBinaryContainer()
	ip := net.ParseIP("2600:0:1:0000:0000:0000:0000:0000")
	assert.True(t, container.ContainsIP(&ip))
}

func TestBinaryAddressContainer_GetIPsInGenRangeEmpty(t *testing.T) {
	container := getEmptyBinaryContainer()
	newIPs := addressing.GetIPsFromStrings([]string {
		"2600:0:1:0000:0000:0000:0000:0001",
		"2600:0:1:0000:0000:0000:0000:0002",
		"2600:0:1:0000:0000:0000:0000:0003",
		"2600:0:1:0000:0000:0000:0000:0004",
	})
	genRange := GetGenRangeFromIPs(newIPs)
	results := container.GetIPsInGenRange(genRange)
	assert.Empty(t, results)
}

func TestBinaryAddressContainer_GetIPsInGenRangeNoWild(t *testing.T) {
	container := getBinaryContainer()
	newIPs := addressing.GetIPsFromStrings([]string {
		"2600:0:1:0000:0000:0000:0000:0001",
	})
	genRange := GetGenRangeFromIPs(newIPs)
	results := container.GetIPsInGenRange(genRange)
	assert.Equal(t, 1, len(results))
}

func TestBinaryAddressContainer_GetIPsInGenRangeOneWild(t *testing.T) {
	container := getBinaryContainer()
	newIPs := addressing.GetIPsFromStrings([]string {
		"2600:0:1:0000:0000:0000:0000:0001",
		"2600:0:1:0000:0000:0000:0000:0002",
	})
	genRange := GetGenRangeFromIPs(newIPs)
	results := container.GetIPsInGenRange(genRange)
	assert.Equal(t, 8, len(results))
}

func TestBinaryAddressContainer_GetIPsInGenRangeTwoWild(t *testing.T) {
	container := getBinaryContainer()
	newIPs := addressing.GetIPsFromStrings([]string {
		"2600:0:1:0000:0000:0000:0000:0001",
		"2600:0:1:0000:0000:0000:0000:0002",
		"2600:0:1:0000:0001:0000:0000:0000",
	})
	genRange := GetGenRangeFromIPs(newIPs)
	results := container.GetIPsInGenRange(genRange)
	assert.Equal(t, 16, len(results))
}

func TestBinaryAddressContainer_CountIPsInGenRangeEmpty(t *testing.T) {
	container := getEmptyBinaryContainer()
	newIPs := addressing.GetIPsFromStrings([]string {
		"2600:0:1:0000:0000:0000:0000:0001",
		"2600:0:1:0000:0000:0000:0000:0002",
		"2600:0:1:0000:0000:0000:0000:0003",
		"2600:0:1:0000:0000:0000:0000:0004",
	})
	genRange := GetGenRangeFromIPs(newIPs)
	results := container.CountIPsInGenRange(genRange)
	assert.Equal(t, 0, results)
}

func TestBinaryAddressContainer_CountIPsInGenRangeNoWild(t *testing.T) {
	container := getBinaryContainer()
	newIPs := addressing.GetIPsFromStrings([]string {
		"2600:0:1:0000:0000:0000:0000:0001",
	})
	genRange := GetGenRangeFromIPs(newIPs)
	results := container.CountIPsInGenRange(genRange)
	assert.Equal(t, 1, results)
}

func TestBinaryAddressContainer_CountIPsInGenRangeOneWild(t *testing.T) {
	container := getBinaryContainer()
	newIPs := addressing.GetIPsFromStrings([]string {
		"2600:0:1:0000:0000:0000:0000:0001",
		"2600:0:1:0000:0000:0000:0000:0002",
	})
	genRange := GetGenRangeFromIPs(newIPs)
	results := container.CountIPsInGenRange(genRange)
	assert.Equal(t, 8, results)
}

func TestBinaryAddressContainer_CountIPsInGenRangeTwoWild(t *testing.T) {
	container := getBinaryContainer()
	newIPs := addressing.GetIPsFromStrings([]string {
		"2600:0:1:0000:0000:0000:0000:0001",
		"2600:0:1:0000:0000:0000:0000:0002",
		"2600:0:1:0000:0001:0000:0000:0000",
	})
	genRange := GetGenRangeFromIPs(newIPs)
	results := container.CountIPsInGenRange(genRange)
	assert.Equal(t, 16, results)
}

func TestFilterByMaskEmpty(t *testing.T) {
	result := filterByMask([]uint64{}, 1, 1)
	assert.Empty(t, result)
}

func TestFilterByMaskNoHits(t *testing.T) {
	candidates := byteListToUints([][]byte{
		{ 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff },
		{ 0xff, 0xff, 0xff, 0xff, 0x00, 0x00, 0x00, 0x00 },
		{ 0x00, 0x00, 0x00, 0x00, 0xff, 0xff, 0xff, 0xff },
		{ 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00 },
	})
	mask := bytesToUint([]byte{
		0xff, 0xff, 0xff, 0xff, 0x00, 0x00, 0x00, 0x00,
	})
	expected := bytesToUint([]byte{
		0x01, 0x02, 0x03, 0x04, 0x00, 0x00, 0x00, 0x00,
	})
	result := filterByMask(candidates, mask, expected)
	assert.Empty(t, result)
}

func TestFilterByMaskFullHits(t *testing.T) {
	candidates := byteListToUints([][]byte{
		{ 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff },
		{ 0xff, 0xff, 0xff, 0xff, 0x00, 0x00, 0x00, 0x00 },
		{ 0x00, 0x00, 0x00, 0x00, 0xff, 0xff, 0xff, 0xff },
		{ 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00 },
	})
	mask := bytesToUint([]byte{
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	})
	expected := bytesToUint([]byte{
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	})
	result := filterByMask(candidates, mask, expected)
	assert.ElementsMatch(t, candidates, result)
}

func TestFilterByMaskHalfHits(t *testing.T) {
	candidates := byteListToUints([][]byte{
		{ 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff },
		{ 0xff, 0xff, 0xff, 0xff, 0x00, 0x00, 0x00, 0x00 },
		{ 0x00, 0x00, 0x00, 0x00, 0xff, 0xff, 0xff, 0xff },
		{ 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00 },
	})
	mask := bytesToUint([]byte{
		0xff, 0xff, 0xff, 0xff, 0x00, 0x00, 0x00, 0x00,
	})
	expected := bytesToUint([]byte{
		0xff, 0xff, 0xff, 0xff, 0x00, 0x00, 0x00, 0x00,
	})
	success := byteListToUints([][]byte{
		{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
		{0xff, 0xff, 0xff, 0xff, 0x00, 0x00, 0x00, 0x00},
	})
	result := filterByMask(candidates, mask, expected)
	assert.ElementsMatch(t, success, result)
}

func TestInsertEmptyResult(t *testing.T) {
	var addrs []uint64
	result, _ := insert(addrs, 1)
	assert.ElementsMatch(t, []uint64{ 1 }, result)
}

func TestInsertEmptyAdded(t *testing.T) {
	var addrs []uint64
	_, added := insert(addrs, 1)
	assert.True(t, added)
}

func TestInsert1Below(t *testing.T) {
	addrs := getSeekRangeTestUints()
	result, _ := insert(addrs, 1)
	assert.ElementsMatch(t, append([]uint64{ 1 }, addrs...), result)
}

func TestInsert1Above(t *testing.T) {
	addrs := getSeekRangeTestUints()
	result, _ := insert(addrs, 51)
	assert.ElementsMatch(t, append(addrs, 51), result)
}

func TestInsertBottomHit(t *testing.T) {
	addrs := getSeekRangeTestUints()
	_, added := insert(addrs, 9)
	assert.False(t, added)
}

func TestInsertTopHit(t *testing.T) {
	addrs := getSeekRangeTestUints()
	_, added := insert(addrs, 50)
	assert.False(t, added)
}

func TestInsertMiddleMiss(t *testing.T) {
	addrs := getSeekRangeTestUints()
	result, _ := insert(addrs, 32)
	expected := []uint64{9, 10, 11, 19, 20, 21, 29, 30, 31, 32, 39, 40, 41, 49, 50 }
	assert.ElementsMatch(t, expected, result)
}

func TestInsertMiddleHit(t *testing.T) {
	addrs := getSeekRangeTestUints()
	_, added := insert(addrs, 31)
	assert.False(t, added)
}

func TestSeekRangeLowerAbove(t *testing.T) {
	uints := getSeekRangeTestUints()
	results := seekRange(uints, 51, 60)
	assert.Empty(t, results)
}

func TestSeekRangeUpperBelow(t *testing.T) {
	uints := getSeekRangeTestUints()
	results := seekRange(uints, 0, 8)
	assert.Empty(t, results)
}

func TestSeekRangeEquivalentHit(t *testing.T) {
	uints := getSeekRangeTestUints()
	results := seekRange(uints, 9, 9)
	assert.ElementsMatch(t, []uint64{ 9 }, results)
}

func TestSeekRangeEquivalentMiss(t *testing.T) {
	uints := getSeekRangeTestUints()
	results := seekRange(uints, 15, 15)
	assert.Empty(t, results)
}

func TestSeekRangeCaptureLowHit(t *testing.T) {
	uints := getSeekRangeTestUints()
	results := seekRange(uints, 9, 15)
	assert.ElementsMatch(t, []uint64{ 9, 10, 11 }, results)
}

func TestSeekRangeCaptureLowMiss(t *testing.T) {
	uints := getSeekRangeTestUints()
	results := seekRange(uints, 8, 15)
	assert.ElementsMatch(t, []uint64{ 9, 10, 11 }, results)
}

func TestSeekRangeCaptureHighHit(t *testing.T) {
	uints := getSeekRangeTestUints()
	results := seekRange(uints, 9, 19)
	assert.ElementsMatch(t, []uint64{ 9, 10, 11, 19 }, results)
}

func TestSeekRangeCaptureHighMiss(t *testing.T) {
	uints := getSeekRangeTestUints()
	results := seekRange(uints, 9, 22)
	assert.ElementsMatch(t, []uint64{ 9, 10, 11, 19, 20, 21 }, results)
}

func TestSeekRangeCaptureBothMiss(t *testing.T) {
	uints := getSeekRangeTestUints()
	results := seekRange(uints, 8, 22)
	assert.ElementsMatch(t, []uint64{ 9, 10, 11, 19, 20, 21 }, results)
}

func TestSeekRangeCaptureBothHit(t *testing.T) {
	uints := getSeekRangeTestUints()
	results := seekRange(uints, 9, 21)
	assert.ElementsMatch(t, []uint64{ 9, 10, 11, 19, 20, 21 }, results)
}

func TestSeekRangeCaptureLowMissHighHit(t *testing.T) {
	uints := getSeekRangeTestUints()
	results := seekRange(uints, 8, 21)
	assert.ElementsMatch(t, []uint64{ 9, 10, 11, 19, 20, 21 }, results)
}

func TestSeekRangeCaptureLowHitHighMiss(t *testing.T) {
	uints := getSeekRangeTestUints()
	results := seekRange(uints, 9, 22)
	assert.ElementsMatch(t, []uint64{ 9, 10, 11, 19, 20, 21 }, results)
}

func TestSeekRangeCaptureMid(t *testing.T) {
	uints := getSeekRangeTestUints()
	results := seekRange(uints, 29, 39)
	assert.ElementsMatch(t, []uint64{ 29, 30, 31, 39}, results)
}

func TestSeekRangeCaptureHigh(t *testing.T) {
	uints := getSeekRangeTestUints()
	results := seekRange(uints, 40, 51)
	assert.ElementsMatch(t, []uint64{ 40, 41, 49, 50 }, results)
}

func TestSeekRangeCaptureFull(t *testing.T) {
	uints := getSeekRangeTestUints()
	results := seekRange(uints, 0, 51)
	assert.ElementsMatch(t, uints, results)
}

func TestSeekRangeDoubleMissOneLength(t *testing.T) {
	uints := []uint64{ 5 }
	results := seekRange(uints, 0, 10)
	assert.ElementsMatch(t, uints, results)
}

func TestSeekEmptyIndex(t *testing.T) {
	index, _ := seek([]uint64{}, 123)
	assert.EqualValues(t, 0, index)
}

func TestSeekEmptyBool(t *testing.T) {
	_, found := seek([]uint64{}, uint64(123))
	assert.False(t, found)
}

func TestSeekOddMissingBool(t *testing.T) {
	uints := getTestUints(true)
	_, found := seek(uints, 35)
	assert.False(t, found)
}

func TestSeekOddMissingIndex(t *testing.T) {
	uints := getTestUints(true)
	index, _ := seek(uints, 35)
	assert.EqualValues(t, 11, index)
}

func TestSeekOddFoundBool(t *testing.T) {
	uints := getTestUints(true)
	_, found := seek(uints, 31)
	assert.True(t, found)
}

func TestSeekOddFoundIndex(t *testing.T) {
	uints := getTestUints(true)
	index, _ := seek(uints, 31)
	assert.EqualValues(t, 10, index)
}

func TestSeekEvenMissingBool(t *testing.T) {
	uints := getTestUints(false)
	_, found := seek(uints, 35)
	assert.False(t, found)
}

func TestSeekEvenMissingIndex(t *testing.T) {
	uints := getTestUints(false)
	index, _ := seek(uints, 35)
	assert.EqualValues(t, 11, index)
}

func TestSeekEvenFoundBool(t *testing.T) {
	uints := getTestUints(false)
	_, found := seek(uints, 31)
	assert.True(t, found)
}

func TestSeekEvenFoundIndex(t *testing.T) {
	uints := getTestUints(false)
	index, _ := seek(uints, 31)
	assert.EqualValues(t, 10, index)
}

func TestSeekOddFoundBottom(t *testing.T) {
	uints := getTestUints(true)
	index, _ := seek(uints, 1)
	assert.EqualValues(t, 0, index)
}

func TestSeekOddFoundLow(t *testing.T) {
	uints := getTestUints(true)
	index, _ := seek(uints, 11)
	assert.EqualValues(t, 4, index)
}

func TestSeekOddFoundMid(t *testing.T) {
	uints := getTestUints(true)
	index, _ := seek(uints, 21)
	assert.EqualValues(t, 7, index)
}

func TestSeekOddFoundHigh(t *testing.T) {
	uints := getTestUints(true)
	index, _ := seek(uints, 39)
	assert.EqualValues(t, 11, index)
}

func TestSeekOddFoundTop(t *testing.T) {
	uints := getTestUints(true)
	index, _ := seek(uints, 51)
	assert.EqualValues(t, 16, index)
}

func TestSeekOddMissingBottom(t *testing.T) {
	uints := getTestUints(true)
	index, _ := seek(uints, 0)
	assert.EqualValues(t, 0, index)
}

func TestSeekOddMissingLow(t *testing.T) {
	uints := getTestUints(true)
	index, _ := seek(uints, 5)
	assert.EqualValues(t, 2, index)
}

func TestSeekOddMissingMid(t *testing.T) {
	uints := getTestUints(true)
	index, _ := seek(uints, 22)
	assert.EqualValues(t, 8, index)
}

func TestSeekOddMissingHigh(t *testing.T) {
	uints := getTestUints(true)
	index, _ := seek(uints, 45)
	assert.EqualValues(t, 14, index)
}

func TestSeekOddMissingTop(t *testing.T) {
	uints := getTestUints(true)
	index, _ := seek(uints, 52)
	assert.EqualValues(t, 17, index)
}

func TestSeekEvenFoundBottom(t *testing.T) {
	uints := getTestUints(false)
	index, _ := seek(uints, 1)
	assert.EqualValues(t, 0, index)
}

func TestSeekEvenFoundLow(t *testing.T) {
	uints := getTestUints(false)
	index, _ := seek(uints, 11)
	assert.EqualValues(t, 4, index)
}

func TestSeekEvenFoundMid(t *testing.T) {
	uints := getTestUints(false)
	index, _ := seek(uints, 21)
	assert.EqualValues(t, 7, index)
}

func TestSeekEvenFoundHigh(t *testing.T) {
	uints := getTestUints(false)
	index, _ := seek(uints, 39)
	assert.EqualValues(t, 11, index)
}

func TestSeekEvenFoundTop(t *testing.T) {
	uints := getTestUints(false)
	index, _ := seek(uints, 50)
	assert.EqualValues(t, 15, index)
}

func TestSeekEvenMissingBottom(t *testing.T) {
	uints := getTestUints(false)
	index, _ := seek(uints, 0)
	assert.EqualValues(t, 0, index)
}

func TestSeekEvenMissingLow(t *testing.T) {
	uints := getTestUints(false)
	index, _ := seek(uints, 5)
	assert.EqualValues(t, 2, index)
}

func TestSeekEvenMissingMid(t *testing.T) {
	uints := getTestUints(false)
	index, _ := seek(uints, 22)
	assert.EqualValues(t, 8, index)
}

func TestSeekEvenMissingHigh(t *testing.T) {
	uints := getTestUints(false)
	index, _ := seek(uints, 45)
	assert.EqualValues(t, 14, index)
}

func TestSeekEvenMissingTop(t *testing.T) {
	uints := getTestUints(false)
	index, _ := seek(uints, 51)
	assert.EqualValues(t, 16, index)
}