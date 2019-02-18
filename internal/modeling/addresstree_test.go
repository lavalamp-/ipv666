package modeling

import (
	"github.com/lavalamp-/ipv666/internal/addressing"
	"github.com/stretchr/testify/assert"
	"net"
	"os"
	"testing"
)

func getDefaultIPs() []*net.IP {
	var ipStrings = []string{
		"2600:0:1:0000:0000:0000:0000:0000",
		"2600:0:1:0000:0000:0000:0000:0001",
		"2600:0:1:0000:0000:0000:0000:0002",
		"2600:0:1:0000:0000:0000:0000:0003",
		"2600:0:1:0000:0000:0000:0000:0004",
		"2600:0:1:0000:0000:0000:0000:0005",
		"2600:0:1:0000:0000:0000:0000:0006",
		"2600:0:1:0000:0000:0000:0000:0007",
		"2600:0:1:0000:0000:0000:0001:0000",
		"2600:0:1:0000:0000:0000:0001:0001",
		"2600:0:1:0000:0000:0000:0001:0002",
		"2600:0:1:0000:0000:0000:0001:0003",
		"2600:0:1:0000:0000:0000:0001:0004",
		"2600:0:1:0000:0000:0000:0001:0005",
		"2600:0:1:0000:0000:0000:0001:0006",
		"2600:0:1:0000:0000:0000:0001:0007",
		"2600:0:1:0000:0000:0001:0000:0000",
		"2600:0:1:0000:0000:0001:0000:0001",
		"2600:0:1:0000:0000:0001:0000:0002",
		"2600:0:1:0000:0000:0001:0000:0003",
		"2600:0:1:0000:0000:0001:0000:0004",
		"2600:0:1:0000:0000:0001:0000:0005",
		"2600:0:1:0000:0000:0001:0000:0006",
		"2600:0:1:0000:0000:0001:0000:0007",
		"2600:0:1:0000:0000:0001:0001:0000",
		"2600:0:1:0000:0000:0001:0001:0001",
		"2600:0:1:0000:0000:0001:0001:0002",
		"2600:0:1:0000:0000:0001:0001:0003",
		"2600:0:1:0000:0000:0001:0001:0004",
		"2600:0:1:0000:0000:0001:0001:0005",
		"2600:0:1:0000:0000:0001:0001:0006",
		"2600:0:1:0000:0000:0001:0001:0007",
		"2600:0:1:0000:0001:0000:0000:0000",
		"2600:0:1:0000:0001:0000:0000:0001",
		"2600:0:1:0000:0001:0000:0000:0002",
		"2600:0:1:0000:0001:0000:0000:0003",
		"2600:0:1:0000:0001:0000:0000:0004",
		"2600:0:1:0000:0001:0000:0000:0005",
		"2600:0:1:0000:0001:0000:0000:0006",
		"2600:0:1:0000:0001:0000:0000:0007",
		"2600:0:1:0000:0001:0000:0001:0000",
		"2600:0:1:0000:0001:0000:0001:0001",
		"2600:0:1:0000:0001:0000:0001:0002",
		"2600:0:1:0000:0001:0000:0001:0003",
		"2600:0:1:0000:0001:0000:0001:0004",
		"2600:0:1:0000:0001:0000:0001:0005",
		"2600:0:1:0000:0001:0000:0001:0006",
		"2600:0:1:0000:0001:0000:0001:0007",
		"2600:0:1:0000:0001:0001:0000:0000",
		"2600:0:1:0000:0001:0001:0000:0001",
		"2600:0:1:0000:0001:0001:0000:0002",
		"2600:0:1:0000:0001:0001:0000:0003",
		"2600:0:1:0000:0001:0001:0000:0004",
		"2600:0:1:0000:0001:0001:0000:0005",
		"2600:0:1:0000:0001:0001:0000:0006",
		"2600:0:1:0000:0001:0001:0000:0007",
		"2600:0:1:0000:0001:0001:0001:0000",
		"2600:0:1:0000:0001:0001:0001:0001",
		"2600:0:1:0000:0001:0001:0001:0002",
		"2600:0:1:0000:0001:0001:0001:0003",
		"2600:0:1:0000:0001:0001:0001:0004",
		"2600:0:1:0000:0001:0001:0001:0005",
		"2600:0:1:0000:0001:0001:0001:0006",
		"2600:0:1:0000:0001:0001:0001:0007",
	}
	return addressing.GetIPsFromStrings(ipStrings)
}

func getAddressTree() *AddressTree {
	return CreateFromAddresses(getDefaultIPs(), 100)
}

func getEmptyAddressTree() *AddressTree {
	return CreateFromAddresses([]*net.IP{}, 100)
}

func TestCreateFromAddressesReturns(t *testing.T) {
	addrTree := getAddressTree()
	assert.NotNil(t, addrTree)
}

func TestCreateFromAddressesEmpty(t *testing.T) {
	addrTree := getEmptyAddressTree()
	assert.NotNil(t, addrTree)
}

func TestCreateFromAddressesCount(t *testing.T) {
	addrTree := getAddressTree()
	assert.EqualValues(t, len(getDefaultIPs()), addrTree.ChildrenCount)
}

func TestAddressTree_AddIPCount(t *testing.T) {
	addrTree := getAddressTree()
	newIP := net.ParseIP("2600:0:1:0001:0000:0000:0000:0001")
	firstCount := addrTree.ChildrenCount
	addrTree.AddIP(&newIP)
	assert.Equal(t, addrTree.ChildrenCount, firstCount + 1)
}

func TestAddressTree_AddIPAdds(t *testing.T) {
	addrTree := getAddressTree()
	newIP := net.ParseIP("2600:0:1:0001:0000:0000:0000:0001")
	addrTree.AddIP(&newIP)
	assert.True(t, addrTree.ContainsIP(&newIP))
}

func TestAddressTree_AddIPsCount(t *testing.T) {
	addrTree := getAddressTree()
	newIPs := addressing.GetIPsFromStrings([]string {
		"2600:0:1:0001:0000:0000:0000:0001",
		"2600:0:1:0001:0000:0000:0000:0002",
		"2600:0:1:0001:0000:0000:0000:0003",
		"2600:0:1:0001:0000:0000:0000:0004",
	})
	firstCount := addrTree.ChildrenCount
	addrTree.AddIPs(newIPs, 100)
	assert.Equal(t, firstCount + uint32(len(newIPs)), addrTree.ChildrenCount)
}

func TestAddressTree_AddIPsAdds(t *testing.T) {
	addrTree := getAddressTree()
	newIPs := addressing.GetIPsFromStrings([]string {
		"2600:0:1:0001:0000:0000:0000:0001",
		"2600:0:1:0001:0000:0000:0000:0002",
		"2600:0:1:0001:0000:0000:0000:0003",
		"2600:0:1:0001:0000:0000:0000:0004",
	})
	addrTree.AddIPs(newIPs, 100)
	for _, newIP := range newIPs {
		assert.True(t, addrTree.ContainsIP(newIP))
	}
}

func TestAddressTree_AddIPsAdded(t *testing.T) {
	addrTree := getAddressTree()
	newIPs := addressing.GetIPsFromStrings([]string {
		"2600:0:1:0001:0000:0000:0000:0001",
		"2600:0:1:0001:0000:0000:0000:0002",
		"2600:0:1:0001:0000:0000:0000:0003",
		"2600:0:1:0001:0000:0000:0000:0004",
	})
	added, _ := addrTree.AddIPs(newIPs, 100)
	assert.Equal(t, 4, added)
}

func TestAddressTree_AddIPsSkipped(t *testing.T) {
	addrTree := getAddressTree()
	newIPs := addressing.GetIPsFromStrings([]string {
		"2600:0:1:0001:0000:0000:0000:0001",
		"2600:0:1:0001:0000:0000:0000:0002",
		"2600:0:1:0001:0000:0000:0000:0003",
		"2600:0:1:0001:0000:0000:0000:0004",
	})
	_, skipped := addrTree.AddIPs(newIPs, 100)
	assert.Equal(t, 0, skipped)
}

func TestAddressTree_GetAllIPsEmptyNotNil(t *testing.T) {
	addrTree := getEmptyAddressTree()
	assert.NotNil(t, addrTree.GetAllIPs())
}

func TestAddressTree_GetAllIPsEmptyValue(t *testing.T) {
	addrTree := getEmptyAddressTree()
	assert.Empty(t, addrTree.GetAllIPs())
}

func TestAddressTree_GetAllIPsNotEmpty(t *testing.T) {
	addrTree := getAddressTree()
	assert.NotEmpty(t, addrTree.GetAllIPs())
}

func TestAddressTree_GetAllIPsMatchesCount(t *testing.T) {
	addrTree := getAddressTree()
	assert.EqualValues(t, addrTree.ChildrenCount, len(addrTree.GetAllIPs()))
}

func TestAddressTree_GetAllIPsMatchesInput(t *testing.T) {
	addrTree := getAddressTree()
	inputAddrs := getDefaultIPs()
	assert.ElementsMatch(t, inputAddrs, addrTree.GetAllIPs())
}

func TestAddressTree_GetIPsInRangeNoError(t *testing.T) {
	addrTree := getAddressTree()
	_, testNet, _ := net.ParseCIDR("2600:0:1:0000:0000:0000:0000:0000/0")
	_, err := addrTree.GetIPsInRange(testNet)
	assert.Nil(t, err)
}

func TestAddressTree_GetIPsInRange0(t *testing.T) {
	addrTree := getAddressTree()
	_, testNet, _ := net.ParseCIDR("2600:0:1:0000:0000:0000:0000:0000/0")
	addrs, _ := addrTree.GetIPsInRange(testNet)
	assert.Equal(t, 64, len(addrs))
}

func TestAddressTree_GetIPsInRange32(t *testing.T) {
	addrTree := getAddressTree()
	_, testNet, _ := net.ParseCIDR("2600:0:1:0000:0000:0000:0000:0000/32")
	addrs, _ := addrTree.GetIPsInRange(testNet)
	assert.Equal(t, 64, len(addrs))
}

func TestAddressTree_GetIPsInRange64(t *testing.T) {
	addrTree := getAddressTree()
	_, testNet, _ := net.ParseCIDR("2600:0:1:0000:0000:0000:0000:0000/64")
	addrs, _ := addrTree.GetIPsInRange(testNet)
	assert.Equal(t, 64, len(addrs))
}

func TestAddressTree_GetIPsInRange96(t *testing.T) {
	addrTree := getAddressTree()
	_, testNet, _ := net.ParseCIDR("2600:0:1:0000:0000:0000:0000:0000/96")
	addrs, _ := addrTree.GetIPsInRange(testNet)
	assert.Equal(t, 16, len(addrs))
}

func TestAddressTree_GetIPsInRange128(t *testing.T) {
	addrTree := getAddressTree()
	_, testNet, _ := net.ParseCIDR("2600:0:1:0000:0000:0000:0000:0000/128")
	addrs, _ := addrTree.GetIPsInRange(testNet)
	assert.Equal(t, 1, len(addrs))
}

func TestAddressTree_CountIPsInRangeNoError(t *testing.T) {
	addrTree := getAddressTree()
	_, testNet, _ := net.ParseCIDR("2600:0:1:0000:0000:0000:0000:0000/0")
	_, err := addrTree.CountIPsInRange(testNet)
	assert.Nil(t, err)
}

func TestAddressTree_CountIPsInRange0(t *testing.T) {
	addrTree := getAddressTree()
	_, testNet, _ := net.ParseCIDR("2600:0:1:0000:0000:0000:0000:0000/0")
	count, _ := addrTree.CountIPsInRange(testNet)
	assert.EqualValues(t, 64, count)
}

func TestAddressTree_CountIPsInRange32(t *testing.T) {
	addrTree := getAddressTree()
	_, testNet, _ := net.ParseCIDR("2600:0:1:0000:0000:0000:0000:0000/32")
	count, _ := addrTree.CountIPsInRange(testNet)
	assert.EqualValues(t, 64, count)
}

func TestAddressTree_CountIPsInRange64(t *testing.T) {
	addrTree := getAddressTree()
	_, testNet, _ := net.ParseCIDR("2600:0:1:0000:0000:0000:0000:0000/64")
	count, _ := addrTree.CountIPsInRange(testNet)
	assert.EqualValues(t, 64, count)
}

func TestAddressTree_CountIPsInRange96(t *testing.T) {
	addrTree := getAddressTree()
	_, testNet, _ := net.ParseCIDR("2600:0:1:0000:0000:0000:0000:0000/96")
	count, _ := addrTree.CountIPsInRange(testNet)
	assert.EqualValues(t, 16, count)
}

func TestAddressTree_CountIPsInRange128(t *testing.T) {
	addrTree := getAddressTree()
	_, testNet, _ := net.ParseCIDR("2600:0:1:0000:0000:0000:0000:0000/128")
	count, _ := addrTree.CountIPsInRange(testNet)
	assert.EqualValues(t, 1, count)
}

func TestAddressTree_SaveCreatesFile(t *testing.T) {
	addrTree := getAddressTree()
	addrTree.Save("/tmp/tester")
	_, err := os.Stat("/tmp/tester")
	exists := !os.IsNotExist(err)
	assert.True(t, exists)
	os.Remove("/tmp/tester")
}

func TestAddressTree_SaveNoError(t *testing.T) {
	addrTree := getAddressTree()
	err := addrTree.Save("/tmp/tester")
	os.Remove("/tmp/tester")
	assert.Nil(t, err)
}

func TestLoadAddressTreeFromFileNoError(t *testing.T) {
	addrTree := getAddressTree()
	addrTree.Save("/tmp/tester")
	_, err := LoadAddressTreeFromFile("/tmp/tester")
	os.Remove("/tmp/tester")
	assert.Nil(t, err)
}

func TestLoadAddressTreeFromFileLoads(t *testing.T) {
	addrTree := getAddressTree()
	addrTree.Save("/tmp/tester")
	newTree, _ := LoadAddressTreeFromFile("/tmp/tester")
	os.Remove("/tmp/tester")
	assert.NotNil(t, newTree)
}

func TestLoadAddressTreeFromFileContent(t *testing.T) {
	addrTree := getAddressTree()
	addrTree.Save("/tmp/tester")
	newTree, _ := LoadAddressTreeFromFile("/tmp/tester")
	os.Remove("/tmp/tester")
	assert.ElementsMatch(t, getDefaultIPs(), newTree.GetAllIPs())
}

func TestAddressTree_ContainsIPEmpty(t *testing.T) {
	addrTree := getEmptyAddressTree()
	ip := net.ParseIP("2600:0:1:0000:0000:0000:0000:0000")
	assert.False(t, addrTree.ContainsIP(&ip))
}

func TestAddressTree_ContainsIPFalse(t *testing.T) {
	addrTree := getAddressTree()
	ip := net.ParseIP("2700:0:1:0000:0000:0000:0000:0000")
	assert.False(t, addrTree.ContainsIP(&ip))
}

func TestAddressTree_ContainsIPTrue(t *testing.T) {
	addrTree := getAddressTree()
	ip := net.ParseIP("2600:0:1:0000:0000:0000:0000:0000")
	assert.True(t, addrTree.ContainsIP(&ip))
}
