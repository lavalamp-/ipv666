package main

import (
	"flag"
	"log"
	"github.com/lavalamp-/ipv666/common/config"
	"github.com/natefinch/lumberjack"
)

func setupLogging() {
  log.SetFlags(log.Flags() & (log.Ldate | log.Ltime))

  log.SetOutput(&lumberjack.Logger{
      Filename:   "/var/log/ipv666.log",
      MaxSize:    10,   // megabytes
      MaxBackups: 10,
      MaxAge:     120,  // days
      Compress:   false,
  })
}

func main() {

	setupLogging()

	var configPath string

	flag.StringVar(&configPath, "config", "config.json", "Local file path to the configuration file to use.")

	conf, err := config.LoadFromFile(configPath)

	if err != nil {
		log.Fatal("Can't proceed without loading valid configuration file.")
	}

	conf.Print()

	//fmt.Printf("Hello world\n")
	//addresses, err := common.GetAddressListFromBitStringsFile("/Users/lavalamp/Documents/Projects/IPv6/modeling/files/filtered_ipv6_addrs.dat")
	//addresses, err := common.GetAddressListFromBinaryFile("/Users/lavalamp/Documents/Projects/IPv6/modeling/files/ipv6_addresses_2.bin")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//addrModel := modeling.GenerateAddressModel(addresses, "My First Model")
	//addrModel.Save("firstmodel.model")
	//addrModel, err := modeling.GetProbablisticModelFromFile("firstmodel.model")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//addresses := addrModel.GenerateMulti(2, 10000000)
	//addresses, err := common.GetAddressListFromBinaryFile("generated_addys.bin")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//addresses.ToAddressesFile("generated_addys.txt")
	//addresses.ToBinaryFile("generated_addys.bin")
	//log.Printf("Woop woop!")
	//addresses.ToBinaryFile("/Users/lavalamp/Documents/Projects/IPv6/modeling/files/ipv6_addresses_2.bin")
	//log.Printf("Woop woop!")
}
