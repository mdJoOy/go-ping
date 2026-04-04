package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"

	//	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
)

func ping(ip net.IP, cf *Config) {
	stat := &Stats{}
	var (
		network     string
		listener    string
		icmpMsgType icmp.Type
		protocol    int
	)

	//deciding using ipv4 or ipv6
	if ip.To4() != nil && !cf.ipv6 {
		network = "ip4:icmp"
		listener = "0.0.0.0"
		icmpMsgType = ipv4.ICMPTypeEcho
		protocol = ICMPProtocol
	} else if ip.To4() == nil && cf.ipv6 {
		network = "ip6:ipv6-icmp"
		listener = "::"
		icmpMsgType = ipv6.ICMPTypeEchoRequest
		protocol = ICMPv6Protocol
	}
	c, err := icmp.ListenPacket(network, listener)
	if err != nil {
		log.Fatalf("cannot open listener connection %w\n", err)
	}
	defer c.Close()

	//setting up ttl
	if !cf.ipv6 {
		if err := c.IPv4PacketConn().SetTTL(cf.ttl); err != nil {
			fmt.Errorf("couldn't set ttl: %w", err)
		}

	} else {
		if err := c.IPv6PacketConn().SetHopLimit(cf.ttl); err != nil {
			fmt.Errorf("couldn't set ttl: %w", err)
		}

	}

	//body for icmp echo request, which contains id, seq number or data in raw bytes
	echoRequestBody := &icmp.Echo{ID: os.Getpid() & 0xffff, Seq: 1, Data: makeBody(cf.size)}
	icmpEchoMsg := icmp.Message{Type: icmpMsgType, Code: 0, Body: echoRequestBody}

	//marshaling the body
	icmpEchoMsgInBytes, err := icmpEchoMsg.Marshal(nil)
	if err != nil {
		log.Fatalf("could not marshal the icmp echo req, %w\n", err)
	}
	//setting deadline

	//	c.SetWriteDeadline(time.Now().Add(3 * time.Second))
	//send the echo msg
	dst := &net.IPAddr{IP: ip}
	for seq := 0; cf.count == 0 || seq < cf.count; seq++ {

		timeNow := time.Now()
		if _, err := c.WriteTo(icmpEchoMsgInBytes, dst); err != nil {
			fmt.Println(err)
			log.Fatalf("write failed:%w\n", err)

		}
		//reading the request
		rb := make([]byte, 1500)
		n, peer, err := c.ReadFrom(rb)
		fmt.Println(n, peer, err)
		if err != nil {
			log.Fatalf("couldnot read the echo reply message: %w\n", err)
		}
		stat.rtt = int(time.Second * time.Since(timeNow))
		rm, err := icmp.ParseMessage(protocol, rb[:n])
		if err != nil {
			fmt.Println("error while parsing the icmp request msg")
			log.Fatal(err)
		}
		stat.sent++
		stat.received++
		stat.loss = (stat.sent - stat.received) / stat.sent * 100
		switch rm.Type {
		case ipv4.ICMPTypeEchoReply, ipv6.ICMPTypeEchoReply:
			fmt.Printf("refelection from %v icmp_seq=%d ttl=%d time=%v ms\n", cf.destination, seq, cf.ttl, stat.rtt*1000)
		default:
			fmt.Printf("expected %v but got %v\n", ipv6.ICMPTypeEchoReply, rm.Type)
		}

	}
	fmt.Printf("sent %d packets, received %d packets, packet loss=%d%%\n", stat.sent, stat.received, stat.loss)
}
