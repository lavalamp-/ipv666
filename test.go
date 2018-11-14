package main

import (
	"net"
	"github.com/lavalamp-/ipv666/common/addressing"
)

func main() {
	ip_1 := net.ParseIP("aaaa:aaaa:aaaa:ffff:ffff:ffff:ffff:ffff")

	addressing.FlipBitsInAddress(&ip_1, 12, 35)

	//var i uint8
	//for i = 8; i < 32; i++ {
	//	addressing.FlipBitsInAddress(ip_1, 8, i)
	//}
	//statemachine.SetStateFile("ipv6results/state.txt", statemachine.NETWORK_GROUP)
	//log.Fatal("LULZ!")
	//rand.Seed(time.Now().UTC().UnixNano())
	//log.Printf("Hello world")
	//log.Printf("%d", 65 % 8)
	//var i uint
	//for i = 0; i < 9; i++ {
	//	curByte := addressing.GetByteWithBitsMasked(i)
	//	log.Printf("Here (%d): %x", i, curByte)
	//}
	//byteMask := addressing.GetByteMask(66)
	//log.Printf("Here: %x", byteMask)
	//for i := 0; i < 129; i++ {
	//	hostBits := zrandom.GenerateHostBits(i)
	//	log.Printf("Here (%d): %08b", i, hostBits)
	//	//byteMask := addressing.GetByteMask(i)
	//	//log.Printf("Here (%d): %x", i, byteMask)
	//}
	//var files []string
	//err := filepath.Walk("ipv6results", func(path string, info os.FileInfo, err error) (error) {
	//	files = append(files, path)
	//	return nil
	//})
	//if err != nil {
	//	panic(err)
	//}
	//for _, file := range files {
	//	fmt.Println(file)
	//}
	//addrs, err := addressing.ReadIPsFromBinaryFile("ipv6_addresses_2.bin")
	//if err != nil {
	//	panic(err)
	//}
	//for i, addr := range addrs {
	//	addrBytes := ([]byte)(*addr)
	//	entropy := zrandom.GetEntropyOfBitsFromRight(addrBytes, 64)
	//	fmt.Printf("(%d) %s %f\n", i, addr, entropy)
	//}
	//ip_1 := net.ParseIP("2606:6000:6008:af00:10c5:7e1b:a7ad:c990")
	//ip_2 := net.ParseIP("2606:6000:6008:af00:10c5:7e1b:a7ad:c990")
	//if ip_1.Equal(ip_2) {
	//	log.Printf("Yup")
	//} else {
	//	log.Printf("Nope")
	//}
	//_, ipnet1, _ := net.ParseCIDR("2001:db8::/32")
	//randAddrs := addressing.GenerateRandomAddressesInNetwork(ipnet1, 20)
	//log.Printf("Here: %s", randAddrs)
	//_, ipnet2, _ := net.ParseCIDR("2002:db8::/32")
	//_, ipnet3, _ := net.ParseCIDR("2003:db8::/32")
	//_, ipnet4, _ := net.ParseCIDR("2001:db8::/32")
	//_, ipnet5, _ := net.ParseCIDR("2002:db8::/32")
	//_, ipnet6, _ := net.ParseCIDR("2003:db8::/32")
	//ipnets := []*net.IPNet{ipnet1, ipnet2, ipnet3, ipnet4, ipnet5, ipnet6}
	//othernets := addressing.GetUniqueNetworks(ipnets, 100)
	//log.Printf("Othernets: %s", othernets)
	//filter := bloom.New(10000000000, 7)
	//model, err := data.GetProbabilisticAddressModel("ipv6results/models")
	//if err != nil {
	//	panic(err)
	//}
	//log.Printf("Generating addresses")
	//addrs := model.GenerateMultiIP(2,10000000, 100000)
	//log.Printf("Updating filter")
	//for i, addr := range addrs {
	//	if i % 100000 == 0 {
	//		log.Printf("Here: %d out of %d", i, len(addrs))
	//	}
	//	addrBytes := ([]byte)(*addr)
	//	filter.Add(addrBytes)
	//}
	//log.Printf("Updated filter with 10m addresses.")
	//file, err := os.OpenFile("tester.bin", os.O_WRONLY | os.O_CREATE, 0600)
	//if err != nil {
	//	panic(err)
	//}
	//bytesWritten, err := filter.WriteTo(file)
	//if err != nil {
	//	panic(err)
	//}
	//if bytesWritten <= 0 {
	//	panic(fmt.Sprintf("Dog wtf: %d", bytesWritten))
	//}
	//log.Printf("WE DID IT WOOOO")
	//falsePositiveRate := filter.EstimateFalsePositiveRate(100000000)
	//log.Printf("False positive rate: %f", falsePositiveRate)
	//err := addressing.WriteIPv6NetworksToFile("test_networks", ipnets)
	//ipnets, err := addressing.ReadIPv6NetworksFromFile("test_networks")
	//log.Printf("Error: %e", err)
	//log.Printf("IPNets: %s", ipnets)
	//othernets := addressing.GetUniqueNetworks(ipnets)
	//log.Printf("Othernets: %s", othernets)
}