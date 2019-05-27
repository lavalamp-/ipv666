package fanout

import (
  "context"
  "fmt"
  "github.com/alecthomas/units"
  "github.com/lavalamp-/ipv666/internal/logging"
  "github.com/lavalamp-/ipv666/internal/data"
  "github.com/lavalamp-/ipv666/internal/addressing"
  "github.com/lavalamp-/ipv666/internal/fs"
  "github.com/lavalamp-/ipv666/internal/config"
  "github.com/spf13/viper"
  "golang.org/x/net/icmp"
  "golang.org/x/net/ipv6"
  "golang.org/x/time/rate"
  "net"
  "os"
  "sync/atomic"
  "time"
)

func Slash64s(bandwidth string) error {
  _, err := fanOut(bandwidth, true, false)
  return err
}

func NybbleAdjacent(bandwidth string) error {
  _, err := fanOut(bandwidth, false, true)
  return err
}

func fanOut(bandwidth string, slash64FanOut bool, nybbleFanOut bool) (string, error) {

  // Instantiate ICMPv6 packet listener
  listener, err := net.ListenPacket("ip6:58", "::")
  if err != nil {
    logging.Warnf("Error thrown when listening for IPv6 packets: %s", err.Error())
    return "", err
  }

  // Instantiate IPv6 packet connection
  conn := ipv6.NewPacketConn(listener)
  if err := conn.SetControlMessage(ipv6.FlagHopLimit|ipv6.FlagSrc|ipv6.FlagDst|ipv6.FlagInterface, true); err != nil {
    logging.Warnf("Error thrown when setting control message: %s", err.Error())
    return "", err
  }

  // Apply ICMP echo reply filter
  var filter ipv6.ICMPFilter
  filter.SetAll(true)
  filter.Accept(ipv6.ICMPTypeEchoReply)
  if err := conn.SetICMPFilter(&filter); err != nil {
    logging.Warnf("Error thrown when setting ICMP filter: %s", err.Error())
    return "", err
  }

  // Ping configuration
  // - 10-byte payload
  // - 255-hop limit
  echoData := []byte("0123456789")
  wcm := &ipv6.ControlMessage{HopLimit: 255}

  // Use the zmap kp/s rates to estimate our bandwidth-constrained ping rate
  maxBandwidthInt, err := units.ParseBase2Bytes(bandwidth)
  targetRate := float64(maxBandwidthInt) / 1e6 * 1300
  rateLimit := rate.Limit(targetRate)
  rateLimiter := rate.NewLimiter(rateLimit, 10)
  ctx := context.Background()

  // Kick off the receive processor
  outputPath := fs.GetTimedFilePath(config.GetPingResultDirPath())
  hitCount := uint64(0)
  ips :=  make(chan net.IPAddr)
  newIps := make(map[string]struct{})
  go processReplies(conn, outputPath, newIps, &hitCount)

  // Generate neighboring networks
  done := make(chan bool, 1)
  go func(newIps map[string]struct{}) {

    if slash64FanOut == true {

      // Generate neighboring /64s
      netIps := make(map[*net.IP]struct{})
      err := generateNeighboring64Networks(ips, netIps)
      if err != nil {
        return
      }

      // Wait until the networks are all processed
      for len(ips) > 0 {
        logging.Debugf("Fan-out has %d networks remaining", len(ips))
        time.Sleep(1*time.Second)
      }
      time.Sleep(2*time.Second)

      // Generate neighboring /64s
      err = generate64NetworkHosts(ips, netIps, newIps)
      if err != nil {
        return
      }

      // Wait until they are all processed
      for len(ips) > 0 {
        logging.Debugf("Fan-out has %d networks remaining", len(ips))
        time.Sleep(1*time.Second)
      }
      time.Sleep(1*time.Second)

    }

    if nybbleFanOut == true {

      // Generate addresses
      err := generateNybbleAdjacentAddrs(ips)
      if err != nil {
        return
      }

      // Wait until the networks are all processed
      for len(ips) > 0 {
        logging.Debugf("Fan-out has %d networks remaining", len(ips))
        time.Sleep(1*time.Second)
      }
      time.Sleep(2*time.Second) 
    }


    done <- true    
  }(newIps)

  bloom, err := data.GetBloomFilter()
  if err != nil {
    return "", err
  }
  blacklist, err := data.GetBlacklist()
  if err != nil {
    return "", err
  }

  // Ping each address
  seq := uint16(0)
  lastSecondCount := uint64(0)
  count := uint64(0)
  lastStatus := time.Now().Unix()
  for len(ips) > 0 || len(done) == 0 {

    // Attempt to read from ips with a 1 second timeout
    select {

      // Read
      case ip := <-ips:
       
        if blacklist.IsIPBlacklisted(&ip.IP) {
          continue
        } else if bloom.Test(ip.IP) {
          continue
        } else {
          bloom.Add(ip.IP)
        }

        // Rate limit outgoing connections
        rateLimiter.Wait(ctx)

        // Build the packet
        ping := icmp.Message{
          Type: ipv6.ICMPTypeEchoRequest,
          Code: 0,
          Body: &icmp.Echo{ID: int(seq), Seq: int(seq), Data: echoData},
        }
        req, err := ping.Marshal(nil)
        if err != nil {
          logging.Warnf("error encoding ICMP echo packet with destination %s (%s)", ip, err)
          return "", err
        }
        seq += 1

        // Send the packet
        _, werr := conn.WriteTo(req, wcm, &ip)
        if werr != nil {

          // Requeue the packet if it failed (i.e. due to network buffer backpressure)
          go func() { ips <- ip }()
          continue
        }

        // Increment the counter
        lastSecondCount += 1
        count += 1
        t := time.Now().Unix()
        if t != lastStatus {
          lastStatus = t
          logging.Infof("Ping-scanned %d addresses (%d hits, %d packets/second)", count, atomic.LoadUint64(&hitCount), lastSecondCount)
          lastSecondCount = 0
        }

      // Timeout 
      case <- time.After(1*time.Second):
        continue
    }
  }

  // Close handle to stop the packet processor
  conn.Close()

  // Close the listener
  listener.Close()

  // Exit
  return "", err
}


func generateNybbleAdjacentAddrs(ips chan net.IPAddr) error {

  // Load the discovered addresses
  cleanPings, err := data.GetCleanPingResults()
  if err != nil {
    return err
  }

  logging.Infof("Performing nybble-adjacent ping scan from %d discovered addresses", len(cleanPings))

  // Get the target network
  network, err := config.GetTargetNetwork()
  if err != nil {
    return err
  }

  nybbleCount := 32
  for x := 0; x < 16; x++ {
    if network.Mask[x] & 0xF0 == 0xF0 {
      nybbleCount -= 1
    } else {
      break
    }
    if network.Mask[x] & 0x0F == 0x0F {
      nybbleCount -= 1
    } else {
      break
    }
  }

  // Generate nybble-adjacent addresses
  addrs, err := addressing.GetAdjacentNetworkAddressesFromIPs(cleanPings, 32-nybbleCount, nybbleCount)
  if err != nil {
    return err
  }
  for _, v := range addrs {
    ips <- net.IPAddr{ IP: *v }
  }

  return nil
}


func generate64NetworkHosts(ips chan net.IPAddr, netIps map[*net.IP]struct{}, newIps map[string]struct{}) error {

  logging.Infof("Fanning out from %d discovered /64 networks (host disovery)", (len(netIps) + len(newIps)))

  // Host discovery
  toScan := make(map[string]struct{})
  for k, _ := range newIps {
    toScan[k] = struct{}{}
  }
  for k, _ := range netIps {
    toScan[k.String()] = struct{}{}
  }

  blockSize := viper.GetInt("FanOutHostBlockSize")
  maxHosts := viper.GetInt("FanOutMaxHosts")
  genIps := make(map[string]struct{})
  count := 0
  for k, _ := range toScan {

    seed := net.ParseIP(k)

    // Generate $blockSize addresses
    for x := 0; x < blockSize; x++ {
      for d := 15; d >= 8; d-- {
        seed[d] += 1
        if seed[d] != 0 {
          break
        }
      }
      ip := net.IP(seed)
      if _, ok := genIps[ip.String()]; !ok {
        ips <- net.IPAddr{ IP: ip }
        genIps[ip.String()] = struct{}{}
        count += 1
      }
    }

    if count >= maxHosts {
      break
    }
  }

  return nil
}


func generateNeighboring64Networks(ips chan net.IPAddr, netIps map[*net.IP]struct{}) error {

  // Load the discovered addresses
  cleanPings, err := data.GetCleanPingResults()
  if err != nil {
    return err
  }

  // Find the /64 networks
  for _, v := range cleanPings {
    _, v2 := addressing.AddressToUints(*v)
    if v2 == 1 {
      netIps[v] = struct{}{}
    }
  }

  logging.Infof("Fanning out from %d discovered /64 networks (network disovery)", len(netIps))

  // Generate neighboring /64 networks
  genIps := make(map[string]struct{})
  blockSize := viper.GetInt("FanOutNetworkBlockSize")
  maxNetworks := viper.GetInt("FanOutMaxHosts")
  count := 0
  for k, _ := range netIps {

    seedUp := net.IP(*k)
    seedDown := net.IP(*k)

    // Generate $blockSize addresses
    for x := 0; x < blockSize; x++ {
      for d := 7; d >= 1; d-- {
        seedUp[d] += 1
        if seedUp[d] != 0 {
          break
        }
      }

      ip := net.IP(seedUp)
      if _, ok := genIps[ip.String()]; !ok {
        ips <- net.IPAddr{ IP: ip }
        genIps[ip.String()] = struct{}{}
        count += 1
      }
    }

    // Generate $blockSize addresses
    for x := 0; x < blockSize; x++ {
      for d := 7; d >= 1; d-- {
        seedDown[d] -= 1
        if seedDown[d] != 0 {
          break
        }
      }

      ip := net.IP(seedDown)
      if _, ok := genIps[ip.String()]; !ok {
        ips <- net.IPAddr{ IP: ip }
        genIps[ip.String()] = struct{}{}
        count += 1
      }
    }

    if count >= maxNetworks {
      break
    }
  }

  return nil
}


func processReplies(conn *ipv6.PacketConn, outputPath string, newIps map[string]struct{}, hitCount *uint64) {

  // Output file
  file, err := os.OpenFile(outputPath, os.O_CREATE|os.O_WRONLY, 0644) 
  if err != nil {
    logging.ErrorF(err)
    return
  }
  defer file.Close()

  // Receive loop
  rxIps := make(map[string]struct{})
  buff := make([]byte, 1500)
  for {

    // Read the next ping response
    rlen, _, raddr, rerr := conn.ReadFrom(buff)
    if rerr != nil {

      // Read timeout
      if nerr, ok := rerr.(net.Error); ok && nerr.Timeout() {
        continue
      }

      // Temporary error
      nerr, ok := rerr.(*net.OpError)
      if ok && nerr.Temporary() {
        continue
      }

      // Permanent error
      break
    }

    newIps[raddr.String()] = struct{}{}

    // Deduplicate received packets
    if _, ok := rxIps[raddr.String()]; !ok {

      // Parse the response
      rm, err := icmp.ParseMessage(58, buff[:rlen])
      if err != nil {
        logging.ErrorF(err)
      }
      atomic.AddUint64(hitCount, 1)
      fmt.Fprintf(file, "%s\n", raddr.String())
      file.Sync()
      logging.Debugf("receiver got response from %s %v (%v)", raddr, buff[:rlen], rm)
    }

    continue
  }
}

