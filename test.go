package main

import (
	"log"
	"github.com/lavalamp-/ipv666/common/addressing"

)

func main() {
	log.Printf("Hello world")
	log.Printf("%d", 65 % 8)
	var i uint
	for i = 0; i < 9; i++ {
		curByte := addressing.GetByteWithBitsMasked(i)
		log.Printf("Here (%d): %x", i, curByte)
	}
	//byteMask := addressing.GetByteMask(66)
	//log.Printf("Here: %x", byteMask)
	for i := 0; i < 129; i++ {
		//byteMask := addressing.GetByteMask(i)
		//log.Printf("Here (%d): %x", i, byteMask)
	}
	//_, ipnet1, _ := net.ParseCIDR("2001:db8::/32")
	//_, ipnet2, _ := net.ParseCIDR("2002:db8::/32")
	//_, ipnet3, _ := net.ParseCIDR("2003:db8::/32")
	//_, ipnet4, _ := net.ParseCIDR("2001:db8::/32")
	//_, ipnet5, _ := net.ParseCIDR("2002:db8::/32")
	//_, ipnet6, _ := net.ParseCIDR("2003:db8::/32")
	//ipnets := []*net.IPNet{ipnet1, ipnet2, ipnet3, ipnet4, ipnet5, ipnet6}
	//othernets := addressing.GetUniqueNetworks(ipnets)
	//log.Printf("Othernets: %s", othernets)
	//err := addressing.WriteIPv6NetworksToFile("test_networks", ipnets)
	ipnets, err := addressing.ReadIPv6NetworksFromFile("test_networks")
	log.Printf("Error: %e", err)
	log.Printf("IPNets: %s", ipnets)
	othernets := addressing.GetUniqueNetworks(ipnets)
	log.Printf("Othernets: %s", othernets)
}