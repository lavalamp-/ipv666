package blacklist

import (
	"testing"
	"net"
	"github.com/stretchr/testify/assert"
)

//TODO fill this out

func TestNetworkBlacklist_AddNetworksAddedNoDupes(t *testing.T) {
	_, net1, _ := net.ParseCIDR("ffff:ffff:ffff:ffff::/64")
	_, net2, _ := net.ParseCIDR("ffff:ffff:ffff:fffe::/64")
	_, net3, _ := net.ParseCIDR("ffff:ffff:ffff:fffd::/64")
	_, net4, _ := net.ParseCIDR("ffff:ffff:ffff:fffc::/64")
	nets := []*net.IPNet{net1, net2, net3, net4}
	blacklist := NewNetworkBlacklist([]*net.IPNet{})
	added, _ := blacklist.AddNetworks(nets)
	assert.EqualValues(t, 4, added)
}

func TestNetworkBlacklist_AddNetworksSkippedNoDupes(t *testing.T) {
	_, net1, _ := net.ParseCIDR("ffff:ffff:ffff:ffff::/64")
	_, net2, _ := net.ParseCIDR("ffff:ffff:ffff:fffe::/64")
	_, net3, _ := net.ParseCIDR("ffff:ffff:ffff:fffd::/64")
	_, net4, _ := net.ParseCIDR("ffff:ffff:ffff:fffc::/64")
	nets := []*net.IPNet{net1, net2, net3, net4}
	blacklist := NewNetworkBlacklist([]*net.IPNet{})
	_, skipped := blacklist.AddNetworks(nets)
	assert.EqualValues(t, 0, skipped)
}

func TestNetworkBlacklist_AddNetworksAddedSomeDupes(t *testing.T) {
	_, net1, _ := net.ParseCIDR("ffff:ffff:ffff:ffff::/64")
	_, net2, _ := net.ParseCIDR("ffff:ffff:ffff:fffe::/64")
	_, net3, _ := net.ParseCIDR("ffff:ffff:ffff:ffff::/64")
	_, net4, _ := net.ParseCIDR("ffff:ffff:ffff:fffe::/64")
	nets := []*net.IPNet{net1, net2, net3, net4}
	blacklist := NewNetworkBlacklist([]*net.IPNet{})
	added, _ := blacklist.AddNetworks(nets)
	assert.EqualValues(t, 2, added)
}

func TestNetworkBlacklist_AddNetworksSkippedSomeDupes(t *testing.T) {
	_, net1, _ := net.ParseCIDR("ffff:ffff:ffff:ffff::/64")
	_, net2, _ := net.ParseCIDR("ffff:ffff:ffff:fffe::/64")
	_, net3, _ := net.ParseCIDR("ffff:ffff:ffff:ffff::/64")
	_, net4, _ := net.ParseCIDR("ffff:ffff:ffff:fffe::/64")
	nets := []*net.IPNet{net1, net2, net3, net4}
	blacklist := NewNetworkBlacklist([]*net.IPNet{})
	_, skipped := blacklist.AddNetworks(nets)
	assert.EqualValues(t, 2, skipped)
}

func TestNetworkBlacklist_AddNetworksAddedAllDupes(t *testing.T) {
	_, net1, _ := net.ParseCIDR("ffff:ffff:ffff:ffff::/64")
	_, net2, _ := net.ParseCIDR("ffff:ffff:ffff:fffe::/64")
	_, net3, _ := net.ParseCIDR("ffff:ffff:ffff:ffff::/64")
	_, net4, _ := net.ParseCIDR("ffff:ffff:ffff:fffe::/64")
	nets := []*net.IPNet{net1, net2, net3, net4}
	_, net_5, _ := net.ParseCIDR("::/0")
	blacklist := NewNetworkBlacklist([]*net.IPNet{net_5})
	added, _ := blacklist.AddNetworks(nets)
	assert.EqualValues(t, 0, added)
}

func TestNetworkBlacklist_AddNetworksSkippedAllDupes(t *testing.T) {
	_, net1, _ := net.ParseCIDR("ffff:ffff:ffff:ffff::/64")
	_, net2, _ := net.ParseCIDR("ffff:ffff:ffff:fffe::/64")
	_, net3, _ := net.ParseCIDR("ffff:ffff:ffff:ffff::/64")
	_, net4, _ := net.ParseCIDR("ffff:ffff:ffff:fffe::/64")
	nets := []*net.IPNet{net1, net2, net3, net4}
	_, net_5, _ := net.ParseCIDR("::/0")
	blacklist := NewNetworkBlacklist([]*net.IPNet{net_5})
	_, skipped := blacklist.AddNetworks(nets)
	assert.EqualValues(t, 4, skipped)
}

func TestNetworkBlacklist_AddNetworkReturnsTrue(t *testing.T) {
	_, net1, _ := net.ParseCIDR("ffff:ffff:ffff:ffff::/64")
	blacklist := NewNetworkBlacklist([]*net.IPNet{})
	added := blacklist.AddNetwork(net1)
	assert.True(t, added)
}

func TestNetworkBlacklist_AddNetworkReturnsFalse(t *testing.T) {
	_, net1, _ := net.ParseCIDR("ffff:ffff:ffff:ffff::/64")
	_, net2, _ := net.ParseCIDR("::/0")
	blacklist := NewNetworkBlacklist([]*net.IPNet{net2})
	added := blacklist.AddNetwork(net1)
	assert.False(t, added)
}

func TestNetworkBlacklist_AddNetworkAddsNetwork(t *testing.T) {
	_, net1, _ := net.ParseCIDR("ffff:ffff:ffff:ffff::/64")
	blacklist := NewNetworkBlacklist([]*net.IPNet{})
	startVal := len(blacklist.Networks)
	blacklist.AddNetwork(net1)
	assert.EqualValues(t, startVal + 1, len(blacklist.Networks))
}

func TestNetworkBlacklist_CleanIPListAllBlacklisted(t *testing.T) {
	ip1 := net.ParseIP("ffff:ffff:ffff:ffff::1")
	ip2 := net.ParseIP("ffff:ffff:ffff:ffff::2")
	ip3 := net.ParseIP("ffff:ffff:ffff:ffff::3")
	ip4 := net.ParseIP("ffff:ffff:ffff:ffff::4")
	ips := []*net.IP{&ip1, &ip2, &ip3, &ip4}
	_, net1, _ := net.ParseCIDR("::/0")
	blacklist := NewNetworkBlacklist([]*net.IPNet{net1})
	cleaned := blacklist.CleanIPList(ips, 9999)
	assert.Empty(t, cleaned)
}

func TestNetworkBlacklist_CleanIPListNoneBlacklisted(t *testing.T) {
	ip1 := net.ParseIP("ffff:ffff:ffff:ffff::1")
	ip2 := net.ParseIP("ffff:ffff:ffff:ffff::2")
	ip3 := net.ParseIP("ffff:ffff:ffff:ffff::3")
	ip4 := net.ParseIP("ffff:ffff:ffff:ffff::4")
	ips := []*net.IP{&ip1, &ip2, &ip3, &ip4}
	blacklist := NewNetworkBlacklist([]*net.IPNet{})
	cleaned := blacklist.CleanIPList(ips, 9999)
	assert.Len(t, cleaned, 4)
}

func TestNetworkBlacklist_CleanIPListSomeBlacklisted(t *testing.T) {
	ip1 := net.ParseIP("ffff:ffff:ffff:ffff::1")
	ip2 := net.ParseIP("ffff:ffff:ffff:ffff::2")
	ip3 := net.ParseIP("ffff:ffff:ffff:fffe::1")
	ip4 := net.ParseIP("ffff:ffff:ffff:fffe::2")
	ips := []*net.IP{&ip1, &ip2, &ip3, &ip4}
	_, net1, _ := net.ParseCIDR("ffff:ffff:ffff:fffe::/64")
	blacklist := NewNetworkBlacklist([]*net.IPNet{net1})
	cleaned := blacklist.CleanIPList(ips, 9999)
	assert.Len(t, cleaned, 2)
}

func TestNetworkBlacklist_IsNetworkBlacklistedTrue(t *testing.T) {
	_, net1, _ := net.ParseCIDR("ffff:ffff:ffff:ffff::/64")
	_, net2, _ := net.ParseCIDR("ffff:ffff:ffff:ffff::/65")
	blacklist := NewNetworkBlacklist([]*net.IPNet{net1})
	blacklisted := blacklist.IsNetworkBlacklisted(net2)
	assert.True(t, blacklisted)
}

func TestNetworkBlacklist_IsNetworkBlacklistedFalse(t *testing.T) {
	_, net1, _ := net.ParseCIDR("ffff:ffff:ffff:ffff::/64")
	_, net2, _ := net.ParseCIDR("ffff:ffff:ffff:ffff::/65")
	blacklist := NewNetworkBlacklist([]*net.IPNet{net2})
	blacklisted := blacklist.IsNetworkBlacklisted(net1)
	assert.False(t, blacklisted)
}

func TestNetworkBlacklist_IsNetworkBlacklistedMinMask(t *testing.T) {
	_, net1, _ := net.ParseCIDR("ffff:ffff:ffff:ffff::/0")
	_, net2, _ := net.ParseCIDR("ffff:ffff:ffff:ffff::/0")
	blacklist := NewNetworkBlacklist([]*net.IPNet{net1})
	blacklisted := blacklist.IsNetworkBlacklisted(net2)
	assert.True(t, blacklisted)
}

func TestNetworkBlacklist_IsNetworkBlacklistedMaxMask(t *testing.T) {
	_, net1, _ := net.ParseCIDR("ffff:ffff:ffff:ffff::/128")
	_, net2, _ := net.ParseCIDR("ffff:ffff:ffff:ffff::/128")
	blacklist := NewNetworkBlacklist([]*net.IPNet{net1})
	blacklisted := blacklist.IsNetworkBlacklisted(net2)
	assert.True(t, blacklisted)
}

func TestNetworkBlacklist_IsNetworkBlacklistedMidMask(t *testing.T) {
	_, net1, _ := net.ParseCIDR("ffff:ffff:ffff:ffff::/64")
	_, net2, _ := net.ParseCIDR("ffff:ffff:ffff:ffff::/64")
	blacklist := NewNetworkBlacklist([]*net.IPNet{net1})
	blacklisted := blacklist.IsNetworkBlacklisted(net2)
	assert.True(t, blacklisted)
}

func TestNetworkBlacklist_IsNetworkBlacklisted32Mask(t *testing.T) {
	_, net1, _ := net.ParseCIDR("ffff:ffff:ffff:ffff::/32")
	_, net2, _ := net.ParseCIDR("ffff:ffff:ffff:ffff::/32")
	blacklist := NewNetworkBlacklist([]*net.IPNet{net1})
	blacklisted := blacklist.IsNetworkBlacklisted(net2)
	assert.True(t, blacklisted)
}

func TestNetworkBlacklist_IsNetworkBlacklisted96Mask(t *testing.T) {
	_, net1, _ := net.ParseCIDR("ffff:ffff:ffff:ffff::/96")
	_, net2, _ := net.ParseCIDR("ffff:ffff:ffff:ffff::/96")
	blacklist := NewNetworkBlacklist([]*net.IPNet{net1})
	blacklisted := blacklist.IsNetworkBlacklisted(net2)
	assert.True(t, blacklisted)
}

func TestNetworkBlacklist_IsIPBlacklistedTrue(t *testing.T) {
	ip1 := net.ParseIP("ffff:ffff:ffff:ffff::1")
	_, net1, _ := net.ParseCIDR("::/0")
	blacklist := NewNetworkBlacklist([]*net.IPNet{net1})
	blacklisted := blacklist.IsIPBlacklisted(&ip1)
	assert.True(t, blacklisted)
}

func TestNetworkBlacklist_IsIPBlacklistedFalse(t *testing.T) {
	ip1 := net.ParseIP("ffff:ffff:ffff:ffff::1")
	_, net1, _ := net.ParseCIDR("::/128")
	blacklist := NewNetworkBlacklist([]*net.IPNet{net1})
	blacklisted := blacklist.IsIPBlacklisted(&ip1)
	assert.False(t, blacklisted)
}

func TestNetworkBlacklist_IsIPBlacklistedMinRange(t *testing.T) {
	ip1 := net.ParseIP("ffff:ffff:ffff:ffff::1")
	_, net1, _ := net.ParseCIDR("ffff:ffff:ffff:ffff::1/0")
	blacklist := NewNetworkBlacklist([]*net.IPNet{net1})
	blacklisted := blacklist.IsIPBlacklisted(&ip1)
	assert.True(t, blacklisted)
}

func TestNetworkBlacklist_IsIPBlacklistedMaxRange(t *testing.T) {
	ip1 := net.ParseIP("ffff:ffff:ffff:ffff::1")
	_, net1, _ := net.ParseCIDR("ffff:ffff:ffff:ffff::1/128")
	blacklist := NewNetworkBlacklist([]*net.IPNet{net1})
	blacklisted := blacklist.IsIPBlacklisted(&ip1)
	assert.True(t, blacklisted)
}

func TestNetworkBlacklist_GetBlacklistingNetworkNil(t *testing.T) {
	ip1 := net.ParseIP("ffff:ffff:ffff:fffb::1")
	_, net1, _ := net.ParseCIDR("ffff:ffff:ffff:ffff::/64")
	_, net2, _ := net.ParseCIDR("ffff:ffff:ffff:fffe::/64")
	_, net3, _ := net.ParseCIDR("ffff:ffff:ffff:fffd::/64")
	_, net4, _ := net.ParseCIDR("ffff:ffff:ffff:fffc::/64")
	nets := []*net.IPNet{net1, net2, net3, net4}
	blacklist := NewNetworkBlacklist(nets)
	blacklistNetwork := blacklist.GetBlacklistingNetwork(&ip1)
	assert.Nil(t, blacklistNetwork)
}

func TestNetworkBlacklist_GetBlacklistingNetworkPrecision(t *testing.T) {
	ip1 := net.ParseIP("ffff:ffff:ffff:ffff::1")
	_, net1, _ := net.ParseCIDR("ffff:ffff:ffff:ffff::/64")
	_, net2, _ := net.ParseCIDR("ffff:ffff:ffff:ffff::/65")
	_, net3, _ := net.ParseCIDR("ffff:ffff:ffff:ffff::/66")
	_, net4, _ := net.ParseCIDR("ffff:ffff:ffff:ffff::/67")
	nets := []*net.IPNet{net1, net2, net3, net4}
	blacklist := NewNetworkBlacklist(nets)
	blacklistNetwork := blacklist.GetBlacklistingNetwork(&ip1)
	assert.Equal(t, net1, blacklistNetwork)
}

func TestNetworkBlacklist_GetBlacklistingNetwork(t *testing.T) {
	ip1 := net.ParseIP("ffff:ffff:ffff:ffff::1")
	_, net1, _ := net.ParseCIDR("ffff:ffff:ffff:ffff::/64")
	_, net2, _ := net.ParseCIDR("ffff:ffff:ffff:fffe::/64")
	_, net3, _ := net.ParseCIDR("ffff:ffff:ffff:fffd::/64")
	_, net4, _ := net.ParseCIDR("ffff:ffff:ffff:fffc::/64")
	nets := []*net.IPNet{net1, net2, net3, net4}
	blacklist := NewNetworkBlacklist(nets)
	blacklistNetwork := blacklist.GetBlacklistingNetwork(&ip1)
	assert.Equal(t, net1, blacklistNetwork)
}
