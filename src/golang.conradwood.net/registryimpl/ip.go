package registryimpl

import (
	"context"
	"flag"
	"fmt"
	"golang.conradwood.net/go-easyops/errors"
	"golang.conradwood.net/go-easyops/linux"
	"golang.conradwood.net/go-easyops/utils"
	"google.golang.org/grpc/peer"
	"net"
	"os"
	"time"
)

var (
	localip         = flag.String("local_ip", "", "if set this will be the local ipaddress returned for service local to registry.")
	currentip       string
	privateIPBlocks []*net.IPNet
	//xlocalip        string
)

func init() {
	for _, cidr := range []string{
		"127.0.0.0/8",    // IPv4 loopback
		"10.0.0.0/8",     // RFC1918
		"172.16.0.0/12",  // RFC1918
		"192.168.0.0/16", // RFC1918
		"169.254.0.0/16", // RFC3927 link-local
		"::1/128",        // IPv6 loopback
		"fe80::/10",      // IPv6 link-local
		"fc00::/7",       // IPv6 unique local addr
	} {
		_, block, err := net.ParseCIDR(cidr)
		if err != nil {
			panic(fmt.Errorf("parse error on %q: %v", cidr, err))
		}
		privateIPBlocks = append(privateIPBlocks, block)
	}

	currentip = myip()
	if currentip == "" {
		fmt.Printf("Warning: Currentip is \"\"\n")
	}
	go ipchecker()
}

type IP struct {
	literal    string
	determined string
	loopback   bool
}

func IPLocal() IP {
	return IP{loopback: true}
}

// no check whatsoever!
func IPFromLiteralString(ip string) IP {
	return IP{literal: ip}
}

// return an ip matching this string. no loopback addresses allowed - use IPLocal() instead
func IPFromString(ip string) IP {
	peerhost, _, err := net.SplitHostPort(ip)
	utils.Bail(fmt.Sprintf("invalid ip \"%s\"", ip), err)
	netip := net.ParseIP(peerhost)
	if netip == nil {
		fmt.Printf("invalid ip \"%s\"\n", ip)
		os.Exit(10)
	}
	if netip.IsLoopback() {
		fmt.Printf(" ip \"%s\" is loopback. must use IPLocal() instead", ip)
		os.Exit(10)
	}
	return IP{literal: ip}
}
func IPFromContext(ctx context.Context) (IP, error) {
	res := IP{}
	peer, ok := peer.FromContext(ctx)
	if !ok {
		fmt.Println("Error getting peer ")
		return res, errors.InvalidArgs(ctx, "Error getting peer from context", "error getting peer from context")
	}
	// transport.strAddr
	ad, ok := peer.Addr.(*net.IPNet)
	if ok {
		if ad.IP.IsLoopback() {
			res.loopback = true
			return res, nil
		}
	}
	peerhost, _, err := net.SplitHostPort(peer.Addr.String())
	if err != nil {
		return res, errors.InvalidArgs(ctx, "Invalid peer", "invalid peer: %s", err)
	}
	netip := net.ParseIP(peerhost)
	if netip == nil {
		return res, errors.InvalidArgs(ctx, "Invalid ip", "invalid ip: %s", peerhost)
	}
	if netip.IsLoopback() {
		res.loopback = true
		return res, nil
	}
	res.determined = peerhost
	return res, nil
}

func (i IP) Equals(ip IP) bool {
	if i.literal != ip.literal {
		return false
	}
	if i.determined != ip.determined {
		return false
	}
	if i.loopback != ip.loopback {
		return false
	}
	return true
}

// what do we tell clients this ip is at the moment?
// e.g. we do not want to tell them "localhost" style ips
func (i IP) ExposeAs() string {
	if i.literal != "" {
		return i.literal
	}
	if i.determined != "" {
		return i.determined
	}
	if i.loopback {
		return currentip
	}
	return "127.0.0.10"
}

// update ip addresses for localhost
func ipchecker() {
	for {
		time.Sleep(10 * time.Second)
		currentip = myip()
	}
}

func myip() string {
	if *localip != "" {
		return *localip
	}
	use_new_code := true
	if use_new_code {
		return linux.New().MyIP()
	}
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		os.Stderr.WriteString("Oops: " + err.Error() + "\n")
		os.Exit(1)
	}

	var ips []*net.IPNet
	for _, a := range addrs {
		ipnet, ok := a.(*net.IPNet)
		if !ok {
			continue
		}
		if ipnet.IP.IsLoopback() {
			continue
		}
		if ipnet.IP.IsLinkLocalMulticast() {
			continue
		}
		ips = append(ips, ipnet)
	}
	for _, ip := range ips {
		if isPrivateIP(ip.IP) {
			sip := ip.IP.String()
			if sip == "" {
				fmt.Printf("Warning returning empty ip\n")
			}
			return sip
		}
	}
	fmt.Printf("WARNING - no suitable ip found!\n")
	for _, ip := range ips {
		fmt.Printf("  %s\n", ip)
	}
	return ""
}

func isPrivateIP(ip net.IP) bool {
	if ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() {
		return true
	}

	for _, block := range privateIPBlocks {
		if block.Contains(ip) {
			return true
		}
	}
	return false
}
