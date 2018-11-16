package statemachine

import (
	"testing"
)

func TestGetSeekPairsFromScanResults(t *testing.T) {
	//nets, _ := addressing.ReadIPv6NetworksFromFile("../../data/networkgroups")
	//log.Printf("Loaded %d networks.", len(nets))
	//netList := blacklist.NewNetworkBlacklist(nets)
	//ip := net.ParseIP("2a04:4e42:20:7400:3386:9ce4:7cde:f0b9")
	//check := netList.IsIPBlacklisted(&ip)
	//log.Printf("Check for IP %s: %v", ip, check)
	//log.Printf("Added %d networks to blacklist. Blacklist count is %d", len(nets), netList.GetCount())
	//bestNets := netList.GetNetworks()
	//log.Printf("Got %d networks back out from blacklist.", len(bestNets))
	//var netStrings []string
	//for _, curNet := range bestNets {
	//	netStrings = append(netStrings, addressing.GetBaseAddressString(curNet))
	//}
	//sort.Strings(netStrings)
	//var initStrings []string
	//for _, curNet := range nets {
	//	initStrings = append(initStrings, addressing.GetBaseAddressString(curNet))
	//}
	//sort.Strings(initStrings)
	//log.Printf("Writing init strings to /tmp/initstrings and retrieved strings to /tmp/retstrings")
	//fs.WriteStringsToFile(initStrings, "/tmp/initstrings")
	//fs.WriteStringsToFile(netStrings, "/tmp/retstrings")
	//
	//ips, err := addressing.ReadIPsFromHexFile("../../data/networkscanresults")
	//log.Printf("Here2: %e", err)
	//conf, err := config.LoadFromFile("../../config.json")
	//log.Printf("Here3: %e", err)
	//getSeekPairsFromScanResults(nets, ips, &conf)
	//addressing.WriteIPv6NetworksToHexFile("../../data/networkgroups.txt", nets)
	//netList := blacklist.NewNetworkBlacklist(nets)
	//bestNets := netList.GetNetworks()
	//addressing.WriteIPv6NetworksToHexFile("/tmp/foobywooby", bestNets)
	//_, net1, _ := net.ParseCIDR("ffff:ffff:ffff:ffff::1/64")
	//_, net2, _ := net.ParseCIDR("ffff:ffff:ffff:fffe::1/64")
	//_, net3, _ := net.ParseCIDR("ffff:ffff:ffff:fffd::1/64")
	//_, net4, _ := net.ParseCIDR("ffff:ffff:ffff:fffc::1/64")
	//nets := []*net.IPNet{net1, net2, net3, net4}
	//ip1 := net.ParseIP("ffff:ffff:ffff:ffff::1")
	//ip2 := net.ParseIP("ffff:ffff:ffff:ffff::2")
	//ip3 := net.ParseIP("ffff:ffff:ffff:ffff::3")
	//ip4 := net.ParseIP("ffff:ffff:ffff:ffff::4")
	//ip5 := net.ParseIP("ffff:ffff:ffff:fffe::1")
	//ip6 := net.ParseIP("ffff:ffff:ffff:fffe::2")
	//ip7 := net.ParseIP("ffff:ffff:ffff:fffe::3")
	//ip8 := net.ParseIP("ffff:ffff:ffff:fffe::4")
	//ip9 := net.ParseIP("ffff:ffff:ffff:fffd::1")
	//ip10 := net.ParseIP("ffff:ffff:ffff:fffd::2")
	//ip11 := net.ParseIP("ffff:ffff:ffff:fffd::3")
	//ip12 := net.ParseIP("ffff:ffff:ffff:fffd::4")
	//ip13 := net.ParseIP("ffff:ffff:ffff:fffc::1")
	//ip14 := net.ParseIP("ffff:ffff:ffff:fffc::2")
	//ip15 := net.ParseIP("ffff:ffff:ffff:fffc::3")
	//ip16 := net.ParseIP("ffff:ffff:ffff:fffc::4")
	//ips := []*net.IP{&ip1, &ip2, &ip3, &ip4, &ip5, &ip6, &ip7, &ip8, &ip9, &ip10, &ip11, &ip12, &ip13, &ip14, &ip15, &ip16}
	//conf, _ := config.LoadFromFile("../../config.json")
	//getSeekPairsFromScanResults(nets, ips, &conf)
}
