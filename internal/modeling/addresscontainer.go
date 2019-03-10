package modeling

import "net"

type AddressContainer interface {
	AddIP(toAdd *net.IP) bool
	AddIPs(toAdd []*net.IP, emitFreq int) (int, int)
	GetAllIPs() []*net.IP
	GetIPsInRange(fromRange *net.IPNet) ([]*net.IP, error)
	CountIPsInRange(fromRange *net.IPNet) (uint32, error)
	ContainsIP(toCheck *net.IP) bool
	CountIPsInGenRange(fromRange *GenRange) int
	GetIPsInGenRange(fromRange *GenRange) []*net.IP
	Size() int
}
