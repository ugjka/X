// MIT License | ugjka@proton.me
// This code may not be used to train artificial intelligence computer models or retrieved by artificial intelligence software or hardware.
// Broadcast your machine on mdns
// Hardcoded to 192.168.1.0/24 subnet
// Usage: bcast mymachine.local
package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/ugjka/mdns"
)

func main() {
	out, err := exec.Command("pidof", os.Args[0]).CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
	arr := bytes.Split(out, []byte(" "))
	if len(arr) > 1 {
		return
	}
	hostname := os.Args[1]
	ip, err := getIP()
	if err != nil {
		log.Println(err)
		return
	}
	const TTL = 60 * 15
	revip := reverse(ip)
	zone, err := mdns.New(true, false)
	if err != nil {
		log.Fatal(err)
	}
	publish(zone, fmt.Sprintf("%s.local. %d IN A %s", hostname, TTL, ip))
	publish(zone, fmt.Sprintf("%s.in-addr.arpa. %d IN PTR %s.local.", revip, TTL, hostname))
	go func() {
		for {
			time.Sleep(time.Second * 10)
			newIP, err := getIP()
			if err != nil {
				continue
			}
			if newIP != ip {
				zone.Shutdown()
				var err error
				zone, err = mdns.New(true, false)
				if err != nil {
					log.Fatal(err)
				}
				ip = newIP
				revip = reverse(ip)
				publish(zone, fmt.Sprintf("%s.local. %d IN A %s", hostname, TTL, ip))
				publish(zone, fmt.Sprintf("%s.in-addr.arpa. %d IN PTR %s.local.", revip, TTL, hostname))
			}
		}
	}()
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	<-sig
	log.Println("Shutting down.")
}

func publish(zone *mdns.Zone, rr string) {
	err := zone.Publish(rr)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func unpublish(zone *mdns.Zone, rr string) {
	err := zone.Unpublish(rr)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func reverse(ip string) string {
	split := strings.Split(ip, ".")
	for i, j := 0, len(split)-1; i < j; i, j = i+1, j-1 {
		split[i], split[j] = split[j], split[i]
	}
	return strings.Join(split, ".")
}

func getIP() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			ip, ok := addr.(*net.IPNet)
			if !ok {
				continue
			}
			if strings.HasPrefix(ip.IP.String(), "192.") {
				return ip.IP.String(), nil
			}
		}
	}
	return "", fmt.Errorf("no ip")
}
