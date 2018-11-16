package blacklist

import (
	"testing"
	"net"
	"github.com/stretchr/testify/assert"
	"github.com/lavalamp-/ipv666/common/addressing"
)

//func TestNetworkBlacklist_AddNetworksAddedNoDupes(t *testing.T) {
//	_, net1, _ := net.ParseCIDR("ffff:ffff:ffff:ffff::/64")
//	_, net2, _ := net.ParseCIDR("ffff:ffff:ffff:fffe::/64")
//	_, net3, _ := net.ParseCIDR("ffff:ffff:ffff:fffd::/64")
//	_, net4, _ := net.ParseCIDR("ffff:ffff:ffff:fffc::/64")
//	nets := []*net.IPNet{net1, net2, net3, net4}
//	blacklist := NewNetworkBlacklist([]*net.IPNet{})
//	added, _ := blacklist.AddNetworks(nets)
//	assert.EqualValues(t, 4, added)
//}
//
//func TestNetworkBlacklist_AddNetworksSkippedNoDupes(t *testing.T) {
//	_, net1, _ := net.ParseCIDR("ffff:ffff:ffff:ffff::/64")
//	_, net2, _ := net.ParseCIDR("ffff:ffff:ffff:fffe::/64")
//	_, net3, _ := net.ParseCIDR("ffff:ffff:ffff:fffd::/64")
//	_, net4, _ := net.ParseCIDR("ffff:ffff:ffff:fffc::/64")
//	nets := []*net.IPNet{net1, net2, net3, net4}
//	blacklist := NewNetworkBlacklist([]*net.IPNet{})
//	_, skipped := blacklist.AddNetworks(nets)
//	assert.EqualValues(t, 0, skipped)
//}
//
//func TestNetworkBlacklist_AddNetworksAddedSomeDupes(t *testing.T) {
//	_, net1, _ := net.ParseCIDR("ffff:ffff:ffff:ffff::/64")
//	_, net2, _ := net.ParseCIDR("ffff:ffff:ffff:fffe::/64")
//	_, net3, _ := net.ParseCIDR("ffff:ffff:ffff:ffff::/64")
//	_, net4, _ := net.ParseCIDR("ffff:ffff:ffff:fffe::/64")
//	nets := []*net.IPNet{net1, net2, net3, net4}
//	blacklist := NewNetworkBlacklist([]*net.IPNet{})
//	added, _ := blacklist.AddNetworks(nets)
//	assert.EqualValues(t, 2, added)
//}
//
//func TestNetworkBlacklist_AddNetworksSkippedSomeDupes(t *testing.T) {
//	_, net1, _ := net.ParseCIDR("ffff:ffff:ffff:ffff::/64")
//	_, net2, _ := net.ParseCIDR("ffff:ffff:ffff:fffe::/64")
//	_, net3, _ := net.ParseCIDR("ffff:ffff:ffff:ffff::/64")
//	_, net4, _ := net.ParseCIDR("ffff:ffff:ffff:fffe::/64")
//	nets := []*net.IPNet{net1, net2, net3, net4}
//	blacklist := NewNetworkBlacklist([]*net.IPNet{})
//	_, skipped := blacklist.AddNetworks(nets)
//	assert.EqualValues(t, 2, skipped)
//}
//
//func TestNetworkBlacklist_AddNetworksAddedAllDupes(t *testing.T) {
//	_, net1, _ := net.ParseCIDR("ffff:ffff:ffff:ffff::/64")
//	_, net2, _ := net.ParseCIDR("ffff:ffff:ffff:fffe::/64")
//	_, net3, _ := net.ParseCIDR("ffff:ffff:ffff:ffff::/64")
//	_, net4, _ := net.ParseCIDR("ffff:ffff:ffff:fffe::/64")
//	nets := []*net.IPNet{net1, net2, net3, net4}
//	_, net_5, _ := net.ParseCIDR("::/0")
//	blacklist := NewNetworkBlacklist([]*net.IPNet{net_5})
//	added, _ := blacklist.AddNetworks(nets)
//	assert.EqualValues(t, 0, added)
//}
//
//func TestNetworkBlacklist_AddNetworksSkippedAllDupes(t *testing.T) {
//	_, net1, _ := net.ParseCIDR("ffff:ffff:ffff:ffff::/64")
//	_, net2, _ := net.ParseCIDR("ffff:ffff:ffff:fffe::/64")
//	_, net3, _ := net.ParseCIDR("ffff:ffff:ffff:ffff::/64")
//	_, net4, _ := net.ParseCIDR("ffff:ffff:ffff:fffe::/64")
//	nets := []*net.IPNet{net1, net2, net3, net4}
//	_, net_5, _ := net.ParseCIDR("::/0")
//	blacklist := NewNetworkBlacklist([]*net.IPNet{net_5})
//	_, skipped := blacklist.AddNetworks(nets)
//	assert.EqualValues(t, 4, skipped)
//}
//
//func TestNetworkBlacklist_AddNetworkReturnsTrue(t *testing.T) {
//	_, net1, _ := net.ParseCIDR("ffff:ffff:ffff:ffff::/64")
//	blacklist := NewNetworkBlacklist([]*net.IPNet{})
//	added := blacklist.AddNetwork(net1)
//	assert.True(t, added)
//}
//
//func TestNetworkBlacklist_AddNetworkReturnsFalse(t *testing.T) {
//	_, net1, _ := net.ParseCIDR("ffff:ffff:ffff:ffff::/64")
//	_, net2, _ := net.ParseCIDR("::/0")
//	blacklist := NewNetworkBlacklist([]*net.IPNet{net2})
//	added := blacklist.AddNetwork(net1)
//	assert.False(t, added)
//}
//
//func TestNetworkBlacklist_AddNetworkAddsNetwork(t *testing.T) {
//	_, net1, _ := net.ParseCIDR("ffff:ffff:ffff:ffff::/64")
//	blacklist := NewNetworkBlacklist([]*net.IPNet{})
//	startVal := blacklist.GetCount()
//	blacklist.AddNetwork(net1)
//	assert.EqualValues(t, startVal + 1, blacklist.GetCount())
//}
//
//func TestNetworkBlacklist_CleanIPListAllBlacklisted(t *testing.T) {
//	ip1 := net.ParseIP("ffff:ffff:ffff:ffff::1")
//	ip2 := net.ParseIP("ffff:ffff:ffff:ffff::2")
//	ip3 := net.ParseIP("ffff:ffff:ffff:ffff::3")
//	ip4 := net.ParseIP("ffff:ffff:ffff:ffff::4")
//	ips := []*net.IP{&ip1, &ip2, &ip3, &ip4}
//	_, net1, _ := net.ParseCIDR("::/0")
//	blacklist := NewNetworkBlacklist([]*net.IPNet{net1})
//	cleaned := blacklist.CleanIPList(ips, 9999)
//	assert.Empty(t, cleaned)
//}
//
//func TestNetworkBlacklist_CleanIPListNoneBlacklisted(t *testing.T) {
//	ip1 := net.ParseIP("ffff:ffff:ffff:ffff::1")
//	ip2 := net.ParseIP("ffff:ffff:ffff:ffff::2")
//	ip3 := net.ParseIP("ffff:ffff:ffff:ffff::3")
//	ip4 := net.ParseIP("ffff:ffff:ffff:ffff::4")
//	ips := []*net.IP{&ip1, &ip2, &ip3, &ip4}
//	blacklist := NewNetworkBlacklist([]*net.IPNet{})
//	cleaned := blacklist.CleanIPList(ips, 9999)
//	assert.Len(t, cleaned, 4)
//}
//
//func TestNetworkBlacklist_CleanIPListSomeBlacklisted(t *testing.T) {
//	ip1 := net.ParseIP("ffff:ffff:ffff:ffff::1")
//	ip2 := net.ParseIP("ffff:ffff:ffff:ffff::2")
//	ip3 := net.ParseIP("ffff:ffff:ffff:fffe::1")
//	ip4 := net.ParseIP("ffff:ffff:ffff:fffe::2")
//	ips := []*net.IP{&ip1, &ip2, &ip3, &ip4}
//	_, net1, _ := net.ParseCIDR("ffff:ffff:ffff:fffe::/64")
//	blacklist := NewNetworkBlacklist([]*net.IPNet{net1})
//	cleaned := blacklist.CleanIPList(ips, 9999)
//	assert.Len(t, cleaned, 2)
//}
//
//func TestNetworkBlacklist_IsNetworkBlacklistedTrue(t *testing.T) {
//	_, net1, _ := net.ParseCIDR("ffff:ffff:ffff:ffff::/64")
//	_, net2, _ := net.ParseCIDR("ffff:ffff:ffff:ffff::/65")
//	blacklist := NewNetworkBlacklist([]*net.IPNet{net1})
//	blacklisted := blacklist.IsNetworkBlacklisted(net2)
//	assert.True(t, blacklisted)
//}
//
//func TestNetworkBlacklist_IsNetworkBlacklistedFalse(t *testing.T) {
//	_, net1, _ := net.ParseCIDR("ffff:ffff:ffff:ffff::/64")
//	_, net2, _ := net.ParseCIDR("ffff:ffff:ffff:ffff::/65")
//	blacklist := NewNetworkBlacklist([]*net.IPNet{net2})
//	blacklisted := blacklist.IsNetworkBlacklisted(net1)
//	assert.False(t, blacklisted)
//}
//
//func TestNetworkBlacklist_IsNetworkBlacklistedMinMask(t *testing.T) {
//	_, net1, _ := net.ParseCIDR("ffff:ffff:ffff:ffff::/0")
//	_, net2, _ := net.ParseCIDR("ffff:ffff:ffff:ffff::/0")
//	blacklist := NewNetworkBlacklist([]*net.IPNet{net1})
//	blacklisted := blacklist.IsNetworkBlacklisted(net2)
//	assert.True(t, blacklisted)
//}
//
//func TestNetworkBlacklist_IsNetworkBlacklistedMaxMask(t *testing.T) {
//	_, net1, _ := net.ParseCIDR("ffff:ffff:ffff:ffff::/128")
//	_, net2, _ := net.ParseCIDR("ffff:ffff:ffff:ffff::/128")
//	blacklist := NewNetworkBlacklist([]*net.IPNet{net1})
//	blacklisted := blacklist.IsNetworkBlacklisted(net2)
//	assert.True(t, blacklisted)
//}
//
//func TestNetworkBlacklist_IsNetworkBlacklistedMidMask(t *testing.T) {
//	_, net1, _ := net.ParseCIDR("ffff:ffff:ffff:ffff::/64")
//	_, net2, _ := net.ParseCIDR("ffff:ffff:ffff:ffff::/64")
//	blacklist := NewNetworkBlacklist([]*net.IPNet{net1})
//	blacklisted := blacklist.IsNetworkBlacklisted(net2)
//	assert.True(t, blacklisted)
//}
//
//func TestNetworkBlacklist_IsNetworkBlacklisted32Mask(t *testing.T) {
//	_, net1, _ := net.ParseCIDR("ffff:ffff:ffff:ffff::/32")
//	_, net2, _ := net.ParseCIDR("ffff:ffff:ffff:ffff::/32")
//	blacklist := NewNetworkBlacklist([]*net.IPNet{net1})
//	blacklisted := blacklist.IsNetworkBlacklisted(net2)
//	assert.True(t, blacklisted)
//}
//
//func TestNetworkBlacklist_IsNetworkBlacklisted96Mask(t *testing.T) {
//	_, net1, _ := net.ParseCIDR("ffff:ffff:ffff:ffff::/96")
//	_, net2, _ := net.ParseCIDR("ffff:ffff:ffff:ffff::/96")
//	blacklist := NewNetworkBlacklist([]*net.IPNet{net1})
//	blacklisted := blacklist.IsNetworkBlacklisted(net2)
//	assert.True(t, blacklisted)
//}
//
//func TestNetworkBlacklist_IsIPBlacklistedTrue(t *testing.T) {
//	ip1 := net.ParseIP("ffff:ffff:ffff:ffff::1")
//	_, net1, _ := net.ParseCIDR("::/0")
//	blacklist := NewNetworkBlacklist([]*net.IPNet{net1})
//	blacklisted := blacklist.IsIPBlacklisted(&ip1)
//	assert.True(t, blacklisted)
//}
//
//func TestNetworkBlacklist_IsIPBlacklistedFalse(t *testing.T) {
//	ip1 := net.ParseIP("ffff:ffff:ffff:ffff::1")
//	_, net1, _ := net.ParseCIDR("::/128")
//	blacklist := NewNetworkBlacklist([]*net.IPNet{net1})
//	blacklisted := blacklist.IsIPBlacklisted(&ip1)
//	assert.False(t, blacklisted)
//}
//
//func TestNetworkBlacklist_IsIPBlacklistedMinRange(t *testing.T) {
//	ip1 := net.ParseIP("ffff:ffff:ffff:ffff::1")
//	_, net1, _ := net.ParseCIDR("ffff:ffff:ffff:ffff::1/0")
//	blacklist := NewNetworkBlacklist([]*net.IPNet{net1})
//	blacklisted := blacklist.IsIPBlacklisted(&ip1)
//	assert.True(t, blacklisted)
//}
//
//func TestNetworkBlacklist_IsIPBlacklistedMaxRange(t *testing.T) {
//	ip1 := net.ParseIP("ffff:ffff:ffff:ffff::1")
//	_, net1, _ := net.ParseCIDR("ffff:ffff:ffff:ffff::1/128")
//	blacklist := NewNetworkBlacklist([]*net.IPNet{net1})
//	blacklisted := blacklist.IsIPBlacklisted(&ip1)
//	assert.True(t, blacklisted)
//}
//
//func TestNetworkBlacklist_GetBlacklistingNetworkFromIPNil(t *testing.T) {
//	ip1 := net.ParseIP("ffff:ffff:ffff:fffb::1")
//	_, net1, _ := net.ParseCIDR("ffff:ffff:ffff:ffff::/64")
//	_, net2, _ := net.ParseCIDR("ffff:ffff:ffff:fffe::/64")
//	_, net3, _ := net.ParseCIDR("ffff:ffff:ffff:fffd::/64")
//	_, net4, _ := net.ParseCIDR("ffff:ffff:ffff:fffc::/64")
//	nets := []*net.IPNet{net1, net2, net3, net4}
//	blacklist := NewNetworkBlacklist(nets)
//	blacklistNetwork := blacklist.GetBlacklistingNetworkFromIP(&ip1)
//	assert.Nil(t, blacklistNetwork)
//}
//
//func TestNetworkBlacklist_GetBlacklistingNetworkFromIPPrecision(t *testing.T) {
//	ip1 := net.ParseIP("ffff:ffff:ffff:ffff::1")
//	_, net1, _ := net.ParseCIDR("ffff:ffff:ffff:ffff::/64")
//	_, net2, _ := net.ParseCIDR("ffff:ffff:ffff:ffff::/65")
//	_, net3, _ := net.ParseCIDR("ffff:ffff:ffff:ffff::/66")
//	_, net4, _ := net.ParseCIDR("ffff:ffff:ffff:ffff::/67")
//	nets := []*net.IPNet{net1, net2, net3, net4}
//	blacklist := NewNetworkBlacklist(nets)
//	blacklistNetwork := blacklist.GetBlacklistingNetworkFromIP(&ip1)
//	assert.Equal(t, net1, blacklistNetwork)
//}
//
//func TestNetworkBlacklist_GetBlacklistingNetworkFromIP(t *testing.T) {
//	ip1 := net.ParseIP("ffff:ffff:ffff:ffff::1")
//	_, net1, _ := net.ParseCIDR("ffff:ffff:ffff:ffff::/64")
//	_, net2, _ := net.ParseCIDR("ffff:ffff:ffff:fffe::/64")
//	_, net3, _ := net.ParseCIDR("ffff:ffff:ffff:fffd::/64")
//	_, net4, _ := net.ParseCIDR("ffff:ffff:ffff:fffc::/64")
//	nets := []*net.IPNet{net1, net2, net3, net4}
//	blacklist := NewNetworkBlacklist(nets)
//	blacklistNetwork := blacklist.GetBlacklistingNetworkFromIP(&ip1)
//	assert.Equal(t, net1, blacklistNetwork)
//}
//
//func TestNetworkBlacklist_GetBlacklistingNetworkFromNetworkNil(t *testing.T) {
//	_, net1, _ := net.ParseCIDR("ffff:ffff:ffff:ffff::/64")
//	_, net2, _ := net.ParseCIDR("ffff:ffff:ffff:fffe::/64")
//	_, net3, _ := net.ParseCIDR("ffff:ffff:ffff:fffd::/64")
//	_, net4, _ := net.ParseCIDR("ffff:ffff:ffff:fffc::/64")
//	_, net5, _ := net.ParseCIDR("ffff:ffff:ffff:fffc::/60")
//	nets := []*net.IPNet{net1, net2, net3, net4}
//	blacklist := NewNetworkBlacklist(nets)
//	blacklistNetwork := blacklist.GetBlacklistingNetworkFromNetwork(net5)
//	assert.Nil(t, blacklistNetwork)
//}
//
//func TestNetworkBlacklist_GetBlacklistingNetworkFromNetworkPrecision(t *testing.T) {
//	_, net1, _ := net.ParseCIDR("ffff:ffff:ffff:ffff::/64")
//	_, net2, _ := net.ParseCIDR("ffff:ffff:ffff:ffff::/65")
//	_, net3, _ := net.ParseCIDR("ffff:ffff:ffff:ffff::/66")
//	_, net4, _ := net.ParseCIDR("ffff:ffff:ffff:ffff::/67")
//	_, net5, _ := net.ParseCIDR("ffff:ffff:ffff:ffff::/68")
//	nets := []*net.IPNet{net1, net2, net3, net4}
//	blacklist := NewNetworkBlacklist(nets)
//	blacklistNetwork := blacklist.GetBlacklistingNetworkFromNetwork(net5)
//	assert.Equal(t, net1, blacklistNetwork)
//}
//
//func TestNetworkBlacklist_GetBlacklistingNetworkFromNetwork(t *testing.T) {
//	_, net1, _ := net.ParseCIDR("ffff:ffff:ffff:ffff::/64")
//	_, net2, _ := net.ParseCIDR("ffff:ffff:ffff:fffe::/64")
//	_, net3, _ := net.ParseCIDR("ffff:ffff:ffff:fffd::/64")
//	_, net4, _ := net.ParseCIDR("ffff:ffff:ffff:fffc::/64")
//	_, net5, _ := net.ParseCIDR("ffff:ffff:ffff:ffff::/64")
//	nets := []*net.IPNet{net1, net2, net3, net4}
//	blacklist := NewNetworkBlacklist(nets)
//	blacklistNetwork := blacklist.GetBlacklistingNetworkFromNetwork(net5)
//	assert.Equal(t, net1, blacklistNetwork)
//}
//
//func TestNetworkBlacklist_GetCount(t *testing.T) {
//	_, net1, _ := net.ParseCIDR("ffff:ffff:ffff:ffff::/64")
//	_, net2, _ := net.ParseCIDR("ffff:ffff:ffff:fffe::/64")
//	_, net3, _ := net.ParseCIDR("ffff:ffff:ffff:fffd::/64")
//	_, net4, _ := net.ParseCIDR("ffff:ffff:ffff:fffc::/64")
//	nets := []*net.IPNet{net1, net2, net3, net4}
//	blacklist := NewNetworkBlacklist(nets)
//	assert.EqualValues(t, 4, blacklist.GetCount())
//}
//
//func TestNetworkBlacklist_GetMaskLengthsSingle(t *testing.T) {
//	netStrings := []string{
//		"2001:ee0:4501:e46a:189:42e5:76ea:fdd9/96",
//		"2406:e001:389d:0:ca0e:14ff:fe42:32d2/96",
//		"2600:9000:201a:a416:c566:363c:e75e:ee76/96",
//		"2404:e800:e700:1501::628/96",
//		"2405:4800:2063:1b92:4884:439e:24fc:829e/96",
//		"2a02:e980:65:541c:adad:b851:565a:1813/96",
//		"2a02:e980:53:9209:6f13:f4b8:7a80:ed6c/96",
//		"2600:9000:203c:3a00:1d:dcd1:c580:93a1/96",
//		"2a01:488:42:1000:50ed:84e6:ff7d:999d/96",
//		"2800:4f0:62:a1f8:e0b9:794c:2246:b1b4/96",
//		"2800:370:55:9734:e940:9470:6f11:efd8/96",
//		"2800:4f0:2:1a8b:19d:c983:e42:431f/96",
//	}
//	nets := addressing.GetNetworksFromStrings(netStrings)
//	blacklist := NewNetworkBlacklist(nets)
//	maskLengths := blacklist.GetMaskLengths()
//	assert.Equal(t, 1, len(maskLengths))
//}
//
//func TestNetworkBlacklist_GetMaskLengthsMulti(t *testing.T) {
//	netStrings := []string{
//		"2001:ee0:4501:e46a:189:42e5:76ea:fdd9/96",
//		"2406:e001:389d:0:ca0e:14ff:fe42:32d2/95",
//		"2600:9000:201a:a416:c566:363c:e75e:ee76/94",
//		"2404:e800:e700:1501::628/93",
//		"2405:4800:2063:1b92:4884:439e:24fc:829e/92",
//		"2a02:e980:65:541c:adad:b851:565a:1813/91",
//		"2a02:e980:53:9209:6f13:f4b8:7a80:ed6c/90",
//	}
//	nets := addressing.GetNetworksFromStrings(netStrings)
//	blacklist := NewNetworkBlacklist(nets)
//	maskLengths := blacklist.GetMaskLengths()
//	assert.Equal(t, 7, len(maskLengths))
//}
//
//func TestNetworkBlacklist_GetMaskLengthsSorted(t *testing.T) {
//	netStrings := []string{
//		"2405:4800:2063:1b92:4884:439e:24fc:829e/92",
//		"2001:ee0:4501:e46a:189:42e5:76ea:fdd9/96",
//		"2600:9000:201a:a416:c566:363c:e75e:ee76/94",
//		"2a02:e980:53:9209:6f13:f4b8:7a80:ed6c/90",
//		"2404:e800:e700:1501::628/93",
//		"2a02:e980:65:541c:adad:b851:565a:1813/91",
//		"2406:e001:389d:0:ca0e:14ff:fe42:32d2/95",
//	}
//	nets := addressing.GetNetworksFromStrings(netStrings)
//	blacklist := NewNetworkBlacklist(nets)
//	maskLengths := blacklist.GetMaskLengths()
//	expected := []int{90, 91, 92, 93, 94, 95, 96}
//	assert.ElementsMatch(t, expected, maskLengths)
//}
//
//func TestNetworkBlacklist_CleanNoChange(t *testing.T) {
//	_, net1, _ := net.ParseCIDR("ffff:ffff:ffff:ffff::/64")
//	_, net2, _ := net.ParseCIDR("ffff:ffff:ffff:fffe::/64")
//	_, net3, _ := net.ParseCIDR("ffff:ffff:ffff:fffd::/64")
//	_, net4, _ := net.ParseCIDR("ffff:ffff:ffff:fffc::/64")
//	nets := []*net.IPNet{net1, net2, net3, net4}
//	blacklist := NewNetworkBlacklist(nets)
//	blacklist.Clean(9999)
//	assert.EqualValues(t, 4, blacklist.GetCount())
//}
//
//func TestNetworkBlacklist_CleanAllChange(t *testing.T) {
//	_, net1, _ := net.ParseCIDR("ffff:ffff:ffff:ffff::/64")
//	_, net2, _ := net.ParseCIDR("ffff:ffff:ffff:ffff::/63")
//	_, net3, _ := net.ParseCIDR("ffff:ffff:ffff:ffff::/62")
//	_, net4, _ := net.ParseCIDR("ffff:ffff:ffff:ffff::/61")
//	nets := []*net.IPNet{net1, net2, net3, net4}
//	blacklist := NewNetworkBlacklist(nets)
//	blacklist.Clean(9999)
//	assert.EqualValues(t, 1, blacklist.GetCount())
//}
//
//func TestWriteNetworkBlacklistToFileNoError(t *testing.T) {
//	_, net1, _ := net.ParseCIDR("ffff:ffff:ffff:ffff::/64")
//	_, net2, _ := net.ParseCIDR("ffff:ffff:ffff:fffe::/64")
//	_, net3, _ := net.ParseCIDR("ffff:ffff:ffff:fffd::/64")
//	_, net4, _ := net.ParseCIDR("ffff:ffff:ffff:fffc::/64")
//	nets := []*net.IPNet{net1, net2, net3, net4}
//	blacklist := NewNetworkBlacklist(nets)
//	err := WriteNetworkBlacklistToFile(fs.GetTemporaryFilePath(), blacklist)
//	assert.Nil(t, err)
//}
//
//func TestWriteNetworkBlacklistToFileWrites(t *testing.T) {
//	_, net1, _ := net.ParseCIDR("ffff:ffff:ffff:ffff::/64")
//	_, net2, _ := net.ParseCIDR("ffff:ffff:ffff:fffe::/64")
//	_, net3, _ := net.ParseCIDR("ffff:ffff:ffff:fffd::/64")
//	_, net4, _ := net.ParseCIDR("ffff:ffff:ffff:fffc::/64")
//	nets := []*net.IPNet{net1, net2, net3, net4}
//	blacklist := NewNetworkBlacklist(nets)
//	filePath := fs.GetTemporaryFilePath()
//	WriteNetworkBlacklistToFile(filePath, blacklist)
//	exists := fs.CheckIfFileExists(filePath)
//	assert.True(t, exists)
//}
//
//func TestWriteNetworkBlacklistToFileLength(t *testing.T) {
//	_, net1, _ := net.ParseCIDR("ffff:ffff:ffff:ffff::/64")
//	_, net2, _ := net.ParseCIDR("ffff:ffff:ffff:fffe::/64")
//	_, net3, _ := net.ParseCIDR("ffff:ffff:ffff:fffd::/64")
//	_, net4, _ := net.ParseCIDR("ffff:ffff:ffff:fffc::/64")
//	nets := []*net.IPNet{net1, net2, net3, net4}
//	blacklist := NewNetworkBlacklist(nets)
//	filePath := fs.GetTemporaryFilePath()
//	WriteNetworkBlacklistToFile(filePath, blacklist)
//	size, _ := fs.CountFileSize(filePath)
//	assert.Zero(t, size % 17)
//}
//
//func TestReadNetworkBlacklistFromFileNoError(t *testing.T) {
//	_, net1, _ := net.ParseCIDR("ffff:ffff:ffff:ffff::/64")
//	_, net2, _ := net.ParseCIDR("ffff:ffff:ffff:fffe::/64")
//	_, net3, _ := net.ParseCIDR("ffff:ffff:ffff:fffd::/64")
//	_, net4, _ := net.ParseCIDR("ffff:ffff:ffff:fffc::/64")
//	nets := []*net.IPNet{net1, net2, net3, net4}
//	blacklist := NewNetworkBlacklist(nets)
//	filePath := fs.GetTemporaryFilePath()
//	WriteNetworkBlacklistToFile(filePath, blacklist)
//	_, err := ReadNetworkBlacklistFromFile(filePath)
//	assert.Nil(t, err)
//}
//
//func TestReadNetworkBlacklistFromFileContent(t *testing.T) {
//	_, net1, _ := net.ParseCIDR("ffff:ffff:ffff:ffff::/64")
//	_, net2, _ := net.ParseCIDR("ffff:ffff:ffff:fffe::/64")
//	_, net3, _ := net.ParseCIDR("ffff:ffff:ffff:fffd::/64")
//	_, net4, _ := net.ParseCIDR("ffff:ffff:ffff:fffc::/64")
//	nets := []*net.IPNet{net1, net2, net3, net4}
//	blacklist := NewNetworkBlacklist(nets)
//	filePath := fs.GetTemporaryFilePath()
//	WriteNetworkBlacklistToFile(filePath, blacklist)
//	newBlacklist, _ := ReadNetworkBlacklistFromFile(filePath)
//	net1Check := newBlacklist.IsNetworkBlacklisted(net1)
//	net2Check := newBlacklist.IsNetworkBlacklisted(net2)
//	net3Check := newBlacklist.IsNetworkBlacklisted(net3)
//	net4Check := newBlacklist.IsNetworkBlacklisted(net4)
//	assert.True(t, net1Check && net2Check && net3Check && net4Check)
//}

//func TestNetworkBlacklist_IsIPBlacklistedSmallBlacklistSize(t *testing.T) {
//	ipStrings := []string{
//		"2600:9000:203c:3a00:1d:dcd1:57cc:3270",
//		"2600:9000:203c:3a00:1d:dcd1:1863:5325",
//		"2600:9000:203c:3a00:1d:dcd1:a3b2:685e",
//		"2600:9000:203c:3a00:1d:dcd1:bafb:8348",
//		"2600:9000:203c:3a00:1d:dcd1:ab0a:2d70",
//		"2600:9000:203c:3a00:1d:dcd1:c07a:893a",
//		"2600:9000:203c:3a00:1d:dcd1:9087:3ec2",
//		"2600:9000:203c:3a00:1d:dcd1:3441:7b8",
//		"2a02:e980:53:9209:6f13:f4b8:4dbf:f58e",
//		"2a02:e980:53:9209:6f13:f4b8:8ec0:7233",
//		"2a02:e980:53:9209:6f13:f4b8:11d5:c243",
//		"2a02:e980:53:9209:6f13:f4b8:b529:f05b",
//		"2a02:e980:53:9209:6f13:f4b8:3a86:2ea2",
//		"2a02:e980:53:9209:6f13:f4b8:5f89:51a7",
//		"2a02:e980:53:9209:6f13:f4b8:31d2:834e",
//		"2a02:e980:53:9209:6f13:f4b8:9f96:be24",
//	}
//	ips := addressing.GetIPsFromStrings(ipStrings)
//	netStrings := []string{
//		"2001:ee0:4501:e46a:189:42e5:76ea:fdd9/96",
//		"2406:e001:389d:0:ca0e:14ff:fe42:32d2/96",
//		"2600:9000:201a:a416:c566:363c:e75e:ee76/96",
//		"2404:e800:e700:1501::628/96",
//		"2405:4800:2063:1b92:4884:439e:24fc:829e/96",
//		"2a02:e980:65:541c:adad:b851:565a:1813/96",
//		"2a02:e980:53:9209:6f13:f4b8:7a80:ed6c/96",
//		"2600:9000:203c:3a00:1d:dcd1:c580:93a1/96",
//		"2a01:488:42:1000:50ed:84e6:ff7d:999d/96",
//		"2800:4f0:62:a1f8:e0b9:794c:2246:b1b4/96",
//		"2800:370:55:9734:e940:9470:6f11:efd8/96",
//		"2800:4f0:2:1a8b:19d:c983:e42:431f/96",
//	}
//	nets := addressing.GetNetworksFromStrings(netStrings)
//	blacklist := NewNetworkBlacklist(nets)
//	cleanedIPs := blacklist.CleanIPList(ips, 9999)
//	assert.Empty(t, cleanedIPs)
//}
//
//func TestNetworkBlacklist_GetNetworksEmpty(t *testing.T) {
//	blacklist := NewNetworkBlacklist([]*net.IPNet{})
//	networks := blacklist.GetNetworks()
//	assert.Empty(t, networks)
//}
//
//func TestNetworkBlacklist_GetNetworksSingleLength(t *testing.T) {
//	_, net1, _ := net.ParseCIDR("ffff:ffff:ffff:ffff::/64")
//	blacklist := NewNetworkBlacklist([]*net.IPNet{net1})
//	networks := blacklist.GetNetworks()
//	assert.Equal(t, 1, len(networks))
//}
//
//func TestNetworkBlacklist_GetNetworksSingleContent(t *testing.T) {
//	_, net1, _ := net.ParseCIDR("ffff:ffff:ffff:ffff::/64")
//	blacklist := NewNetworkBlacklist([]*net.IPNet{net1})
//	networks := blacklist.GetNetworks()
//	assert.EqualValues(t, "ffff:ffff:ffff:ffff::/64", networks[0].String())
//}
//
//func TestNetworkBlacklist_GetNetworksManyLength(t *testing.T) {
//	netStrings := []string{
//		"2600:9000:203d:982b:4faf:485:410a:af42/96",
//		"2a02:e980:a4:901:e0a3:7148:22ac:f6bd/96",
//		"2600:9000:202d:2c00:1d:476d:9c40:93a1/96",
//		"2800:370:2:d873:1994:df93:82bb:5fbf/96",
//		"2a02:e980:64:426b:bfec:afb7:3f76:ea84/96",
//		"2604:2d80:400e:0:9105:998f:b472:ad96/96",
//		"2406:e000:9ac1:0:3a10:d5ff:fe3b:5590/96",
//		"2806:108e:3:f0ad:e63e:d7ff:fe0c:a7ce/96",
//		"2600:9000:20ab:d4ba:8e98:b1b6:27de:7d48/96",
//		"2604:2d80:c804:0:2514:203:80d4:cbed/96",
//		"2800:370:61:488:89ea:8ed9:a332:ee96/96",
//		"2604:2d80:882b:0:2d51:53ab:cb56:3525/96",
//		"2600:3400:2:417::3/96",
//		"2001:16a2:5:6ca3:e48f:78a9:2340:521d/96",
//		"2600:9000:2037:3136:1581:bc2a:ad72:8c5c/96",
//		"2600:9000:202d:ca00:14:5d3a:200:21/96",
//		"2806:10a6:a:298b:26bc:f8ff:feba:89ba/96",
//		"2a04:4e42:2:d051:b8cb:bb42:ccd6:a959/96",
//		"2600:9000:202d:de00:d:85cc:4780:93a1/96",
//		"2405:4800:4248:e2e1:6135:dbcd:a5a4:fa27/96",
//		"2a04:4e42:c:5945:e264:7c98:79bc:d73e/96",
//		"2800:68:10:23::2/96",
//		"2600:9000:202d:a400:15:6b13:f9c0:93a1/96",
//		"2600:9000:214c:a479:a91e:2174:9885:c7c8/96",
//		"2600:9000:203b:738a:f791:3f1f:5e02:66ca/96",
//		"e553:1ed0::/96",
//		"2600:9000:20d3:da00:f:db28:d680:93a1/96",
//		"2800:370:55:e6a1:2d49:2051:d201:cfe4/96",
//		"2001:610:600:35f::1/96",
//		"2600:9000:2043:b4fb:db16:6395:c3fc:5370/96",
//		"2003:0:5a00:40f::1/96",
//		"2a0a:3407:100:b423:dd30:6dea:cf36:9f86/96",
//		"2a02:e980:7:a05:1988:f373:46c6:b3d7/96",
//		"2001:bc8:2800:6cec:62ba:3bbd:9cb5:d975/96",
//		"2a02:e980:88:42bc:576a:e83b:9bfa:14af/96",
//		"2600:9000:202f:f8e1:114e:58ed:8e93:978/96",
//		"2600:9000:20d3:b000:14:f09b:1900:93a1/96",
//		"2001:978:2:40::12:1/96",
//	}
//	nets := addressing.GetNetworksFromStrings(netStrings)
//	blacklist := NewNetworkBlacklist(nets)
//	networks := blacklist.GetNetworks()
//	assert.EqualValues(t, len(netStrings), len(networks))
//}
//
//func TestNetworkBlacklist_GetNetworksManyContent(t *testing.T) {
//	netStrings := []string{
//		"2600:9000:203d:982b:4faf:485:410a:af42/96",
//		"2a02:e980:a4:901:e0a3:7148:22ac:f6bd/96",
//		"2600:9000:202d:2c00:1d:476d:9c40:93a1/96",
//		"2800:370:2:d873:1994:df93:82bb:5fbf/96",
//		"2a02:e980:64:426b:bfec:afb7:3f76:ea84/96",
//		"2604:2d80:400e:0:9105:998f:b472:ad96/96",
//		"2406:e000:9ac1:0:3a10:d5ff:fe3b:5590/96",
//		"2806:108e:3:f0ad:e63e:d7ff:fe0c:a7ce/96",
//		"2600:9000:20ab:d4ba:8e98:b1b6:27de:7d48/96",
//		"2604:2d80:c804:0:2514:203:80d4:cbed/96",
//		"2800:370:61:488:89ea:8ed9:a332:ee96/96",
//		"2604:2d80:882b:0:2d51:53ab:cb56:3525/96",
//		"2600:3400:2:417::3/96",
//		"2001:16a2:5:6ca3:e48f:78a9:2340:521d/96",
//		"2600:9000:2037:3136:1581:bc2a:ad72:8c5c/96",
//		"2600:9000:202d:ca00:14:5d3a:200:21/96",
//		"2806:10a6:a:298b:26bc:f8ff:feba:89ba/96",
//		"2a04:4e42:2:d051:b8cb:bb42:ccd6:a959/96",
//		"2600:9000:202d:de00:d:85cc:4780:93a1/96",
//		"2405:4800:4248:e2e1:6135:dbcd:a5a4:fa27/96",
//		"2a04:4e42:c:5945:e264:7c98:79bc:d73e/96",
//		"2800:68:10:23::2/96",
//		"2600:9000:202d:a400:15:6b13:f9c0:93a1/96",
//		"2600:9000:214c:a479:a91e:2174:9885:c7c8/96",
//		"2600:9000:203b:738a:f791:3f1f:5e02:66ca/96",
//		"e553:1ed0::/96",
//		"2600:9000:20d3:da00:f:db28:d680:93a1/96",
//		"2800:370:55:e6a1:2d49:2051:d201:cfe4/96",
//		"2001:610:600:35f::1/96",
//		"2600:9000:2043:b4fb:db16:6395:c3fc:5370/96",
//		"2003:0:5a00:40f::1/96",
//		"2a0a:3407:100:b423:dd30:6dea:cf36:9f86/96",
//		"2a02:e980:7:a05:1988:f373:46c6:b3d7/96",
//		"2001:bc8:2800:6cec:62ba:3bbd:9cb5:d975/96",
//		"2a02:e980:88:42bc:576a:e83b:9bfa:14af/96",
//		"2600:9000:202f:f8e1:114e:58ed:8e93:978/96",
//		"2600:9000:20d3:b000:14:f09b:1900:93a1/96",
//		"2001:978:2:40::12:1/96",
//	}
//	nets := addressing.GetNetworksFromStrings(netStrings)
//	blacklist := NewNetworkBlacklist(nets)
//	networks := blacklist.GetNetworks()
//	var resultStrings []string
//	for _, network := range networks {
//		resultStrings = append(resultStrings, network.String())
//	}
//	sort.Strings(resultStrings)
//	var expectedStrings []string
//	for _, network := range nets {
//		expectedStrings = append(expectedStrings, addressing.GetBaseAddressString(network))
//	}
//	sort.Strings(expectedStrings)
//	assert.ElementsMatch(t, expectedStrings, resultStrings)
//}

func TestNetworkBlacklist_AddNetworkShortShortFails(t *testing.T) {
	_, net1, _ := net.ParseCIDR("2001:0:4137:9e76:101c:b89::/96")
	_, net2, _ := net.ParseCIDR("2001:0:4137:9e76:101c:b89::/96")
	blacklist := NewNetworkBlacklist([]*net.IPNet{net1})
	added := blacklist.AddNetwork(net2)
	assert.False(t, added)
}

func TestNetworkBlacklist_AddNetworkShortLongFails(t *testing.T) {
	_, net1, _ := net.ParseCIDR("2001:0:4137:9e76:101c:b89::/96")
	_, net2, _ := net.ParseCIDR("2001:0:4137:9e76:101c:b89:ce63:392a/96")
	blacklist := NewNetworkBlacklist([]*net.IPNet{net1})
	added := blacklist.AddNetwork(net2)
	assert.False(t, added)
}

func TestNetworkBlacklist_AddNetworkLongShortFails(t *testing.T) {
	_, net1, _ := net.ParseCIDR("2001:0:4137:9e76:101c:b89:ce63:392a/96")
	_, net2, _ := net.ParseCIDR("2001:0:4137:9e76:101c:b89::/96")
	blacklist := NewNetworkBlacklist([]*net.IPNet{net1})
	added := blacklist.AddNetwork(net2)
	assert.False(t, added)
}

func TestNetworkBlacklist_AddNetworkLongLongFails(t *testing.T) {
	_, net1, _ := net.ParseCIDR("2001:0:4137:9e76:101c:b89:ce63:392a/96")
	_, net2, _ := net.ParseCIDR("2001:0:4137:9e76:101c:b89:ce63:392a/96")
	blacklist := NewNetworkBlacklist([]*net.IPNet{net1})
	added := blacklist.AddNetwork(net2)
	assert.False(t, added)
}

func TestNetworkBlacklist_AddNetworkWith100Fails(t *testing.T) {
	_, net1, _ := net.ParseCIDR("2001:0:4137:9e76:101c:b89:ce63:392a/96")
	_, net2, _ := net.ParseCIDR("2001:0:4137:9e76:101c:b89:ce63:392a/96")
	nets := addressing.GenerateRandomNetworks(99, 96)
	nets = append(nets, net1)
	blacklist := NewNetworkBlacklist(nets)
	added := blacklist.AddNetwork(net2)
	assert.False(t, added)
}

func TestNetworkBlacklist_AddNetworkWith1000Fails(t *testing.T) {
	_, net1, _ := net.ParseCIDR("2001:0:4137:9e76:101c:b89:ce63:392a/96")
	_, net2, _ := net.ParseCIDR("2001:0:4137:9e76:101c:b89:ce63:392a/96")
	nets := addressing.GenerateRandomNetworks(999, 96)
	nets = append(nets, net1)
	blacklist := NewNetworkBlacklist(nets)
	added := blacklist.AddNetwork(net2)
	assert.False(t, added)
}

func TestNetworkBlacklist_AddNetworkWith10000Fails(t *testing.T) {
	_, net1, _ := net.ParseCIDR("2001:0:4137:9e76:101c:b89:ce63:392a/96")
	_, net2, _ := net.ParseCIDR("2001:0:4137:9e76:101c:b89:ce63:392a/96")
	nets := addressing.GenerateRandomNetworks(9999, 96)
	nets = append(nets, net1)
	blacklist := NewNetworkBlacklist(nets)
	added := blacklist.AddNetwork(net2)
	assert.False(t, added)
}

func TestNetworkBlacklist_IsIPBlacklistedLongNetwork(t *testing.T) {
	_, net1, _ := net.ParseCIDR("2001:0:4137:9e76:101c:b89:ce63:392a/96")
	blacklist := NewNetworkBlacklist([]*net.IPNet{net1})
	ip := net.ParseIP("2001:0:4137:9e76:101c:b89:ffff:392a")
	isBlacklisted := blacklist.IsIPBlacklisted(&ip)
	assert.True(t, isBlacklisted)
}

func TestNetworkBlacklist_IsIPBlacklistedShortNetwork(t *testing.T) {
	_, net1, _ := net.ParseCIDR("2001:0:4137:9e76:101c:b89::/96")
	blacklist := NewNetworkBlacklist([]*net.IPNet{net1})
	ip := net.ParseIP("2001:0:4137:9e76:101c:b89:ffff:392a")
	isBlacklisted := blacklist.IsIPBlacklisted(&ip)
	assert.True(t, isBlacklisted)
}
