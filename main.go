package main

import (
	"fmt"
	"log"
	"github.com/lavalamp-/ipv666/common/modeling"
)

func main() {
	fmt.Printf("Hello world\n")
	//addresses, err := common.GetAddressListFromBitStringsFile("/Users/lavalamp/Documents/Projects/IPv6/modeling/files/filtered_ipv6_addrs.dat")
	//addresses, err := common.GetAddressListFromBinaryFile("/Users/lavalamp/Documents/Projects/IPv6/modeling/files/ipv6_addresses_2.bin")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//addrModel := modeling.GenerateAddressModel(addresses, "My First Model")
	//addrModel.Save("firstmodel.model")
	addrModel, err := modeling.GetProbablisticModelFromFile("firstmodel.model")
	if err != nil {
		log.Fatal(err)
	}
	addresses := addrModel.GenerateMulti(2, 10000000)
	//addresses, err := common.GetAddressListFromBinaryFile("generated_addys.bin")
	//if err != nil {
	//	log.Fatal(err)
	//}
	addresses.ToAddressesFile("generated_addys.txt")
	//addresses.ToBinaryFile("generated_addys.bin")
	//log.Printf("Woop woop!")
	//addresses.ToBinaryFile("/Users/lavalamp/Documents/Projects/IPv6/modeling/files/ipv6_addresses_2.bin")
	log.Printf("Woop woop!")
}
