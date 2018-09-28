package ping

import (
  "fmt"
  "time"
  "github.com/sparrc/go-ping"
)

// Ping a host and return the number of received packets
func Ping(host string, timeout time.Duration, interval time.Duration, count int, privileged bool, verbose bool) (int, error) {

  // Instantiate the pinger
  pinger, err := ping.NewPinger(host)
  if err != nil {
    fmt.Printf("ERROR: %s\n", err.Error())
    return 0, err
  }

  // Handle echo responses
  pinger.OnRecv = func(pkt *ping.Packet) {
    if verbose {
      fmt.Printf("%d bytes from %s: icmp_seq=%d time=%v\n",
        pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt)
    }
  }

  // Handle ping completed
  pinger.OnFinish = func(stats *ping.Statistics) {
    if verbose {
      fmt.Printf("\n--- %s ping statistics ---\n", stats.Addr)
      fmt.Printf("%d packets transmitted, %d packets received, %v%% packet loss\n",
        stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss)
      fmt.Printf("round-trip min/avg/max/stddev = %v/%v/%v/%v\n",
        stats.MinRtt, stats.AvgRtt, stats.MaxRtt, stats.StdDevRtt)
    }
  }

  // Configure the pinger instance
  pinger.Count = count
  pinger.Interval = interval
  pinger.Timeout = timeout
  pinger.SetPrivileged(privileged)

  // Ping the target
  if verbose {
    fmt.Printf("PING %s (%s):\n", pinger.Addr(), pinger.IPAddr())
  }
  pinger.Run()

  // Return the received packet count
  return pinger.PacketsRecv, nil
}