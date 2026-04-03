package main

import (
	"flag"
	"fmt"
	"log"
)

const (
	ICMPProtocol   = 1
	ICMPv6Protocol = 58
)

func main() {
	config := &Config{}

	flag.IntVar(&config.count, "c", 0, "stop after <count> replies")
	flag.BoolVar(&config.ipv6, "6", false, "use Ipv6")
	flag.Parse()
	config.destination = flag.Arg(0)
	fmt.Println(config.destination)
	fmt.Println(config.ipv6)
	ip, err := resolveHostIP(config.destination, config.ipv6)
	fmt.Println(ip)
	if err != nil {
		log.Fatal("couldnot resolve ip address")
	}
	ping(ip, config)
}
