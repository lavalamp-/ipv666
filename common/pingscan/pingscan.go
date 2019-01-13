package pingscan

import (
	"bufio"
	"context"
	"fmt"
	"github.com/alecthomas/units"
	"github.com/lavalamp-/ipv666/common/logging"
	"github.com/spf13/viper"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv6"
	"golang.org/x/time/rate"
	"net"
	"os"
	"sync/atomic"
	"time"
)

func processReplies(conn *ipv6.PacketConn, outputFile string, done chan bool, hitCount *uint64) {

	// Output file
	file, err := os.OpenFile(outputFile, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logging.ErrorF(err)
		return
	}
	defer file.Close()

	// Receive loop
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

		// Parse the response
		_, err := icmp.ParseMessage(58, buff[:rlen])
		if err != nil {
			logging.ErrorF(err)
			continue
		}
		atomic.AddUint64(hitCount, 1)
		fmt.Fprintf(file, "%s\n", raddr.String())
		file.Sync()
	}
	done <- true
}

func Scan(inputFile string, outputFile string, bandwidth string) (string, error) {

	logging.Infof("Performing ping scan on addresses defined in %s", inputFile)

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

	// Read the addresses from disk and queue them in the channel
	ips := make(chan net.IPAddr)
	done := make(chan bool, 2)
	go func() {
		file, err := os.Open(inputFile)
		if err != nil {
			logging.Warnf("Error thrown when opening IP input file: %s", err.Error())
			return
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			ip := scanner.Text()
			parsedAddr := net.ParseIP(ip)
			dstAddr := net.IPAddr{IP: parsedAddr}
			ips <- dstAddr
		}
		done <- true
	}()

	// Kick off the receive processor
	hitCount := uint64(0)
	go processReplies(conn, outputFile, done, &hitCount)

	// Ping each address
	seq := uint16(0)
	finished := false
	count := uint64(0)
	lastSecondCount := uint64(0)
	lastStatus := time.Now().Unix()
	for finished == false {

		// Attempt to read from ips with a 1 second timeout
		select {

		// Read
		case ip := <-ips:

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
		case <-time.After(5 * time.Second):
			finished = true
			break
		}
	}

	// Wait for the file read goroutine to finish
	<-done

	// Close handle to stop the packet processor
	conn.Close()

	// Close the listener
	listener.Close()

	// Wait for the receiver goroutine to finish
	<-done

	return "", nil
}

func ScanFromConfig(inputFile string, outputFile string) (string, error) {
	return Scan(inputFile, outputFile, viper.GetString("PingScanBandwidth"))
}
