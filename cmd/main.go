package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

const (
	ICMPProtocol   = 1
	ICMPv6Protocol = 58
)

func main() {
	if uid := os.Getuid(); uid != 0 {
		fmt.Println("use SUDO: ping/icmp uses raw socket, which requires super user permission")
		os.Exit(1)
	}
	config := &Config{}

	flag.IntVar(&config.count, "c", 0, "stop after <count> replies")
	flag.BoolVar(&config.ipv6, "6", false, "use Ipv6")
	flag.BoolVar(&config.flood, "f", false, "flood ping")
	flag.BoolVar(&config.quiet, "q", false, "Quiet — only print final summary")
	flag.BoolVar(&config.histogram, "H", false, "Show ASCII latency histogram after summary")
	//ttl, size, deadline, interval
	flag.IntVar(&config.ttl, "t", 64, "define time to live")
	flag.IntVar(&config.size, "s", 56, "use <size> as number of data bytes to send")
	flag.DurationVar(&config.interval, "i", time.Second, "seconds between sending each packet")
	flag.DurationVar(&config.timeout, "W", 3*time.Second, "Per-packet response timeout")
	flag.DurationVar(&config.deadLine, "w", time.Minute, "Exit after this duration")
	//multi ping
	var hostsString string
	flag.StringVar(&hostsString, "m", " ", "Write hosts inside double quot but seperate them by (,)")

	// flag.IntVar(&config.interval, "i", timetime.Second, "seconds between sending each packet")
	var usages string = `Usage
  goping [options] <destination>

Options:
  <destination>      DNS name or IP address
  -c <count>         Stop after <count> replies
  -f                 Flood ping
  -W <timeout>       Time to wait for response
  -i <interval>      Seconds between sending each packet
  -s <size>          Use <size> as number of data bytes to be sent
  -t <ttl>           Define time to live
  -w <deadline>   	 Exit after this duration (e.g. 10s, 1m)  
  -6                 Use IPv6
  -H                 Show ASCII latency histogram after summary
  -m				 Write hosts inside double quot("")but seperate them by (,)
`
	flag.Usage = func() { fmt.Print(usages) }

	flag.Parse()
	hosts := strings.Split(hostsString, ",")
	for i, v := range hosts {
		hosts[i] = strings.TrimSpace(v)
	}
	fmt.Println(hosts)
	//
	if len(hosts) > 1 {
		pingMultiple(*config, hosts)
		time.Sleep(2 * time.Minute)
	} else {

		if flag.NArg() < 1 {
			fmt.Print(usages)
			os.Exit(1)
		}

		config.destination = flag.Arg(0)

		ip := resolveHostIP(config.destination, config.ipv6)
		// if err != nil {
		// 	log.Fatal("couldnot resolve ip address")
		// }
		ping(ip, config)
	}

}
