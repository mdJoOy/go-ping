package main

import (
	"fmt"
	"log"
	"net"
	"os"
	//	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
)

func ping(ip net.IP, cf *Config) {
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

	//body for icmp echo request, which contains id, seq number or data in raw bytes
	echoRequestBody := &icmp.Echo{ID: os.Getpid() & 0xffff, Seq: 1, Data: makeBody(56)}
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
	rm, err := icmp.ParseMessage(protocol, rb[:n])
	if err != nil {
		fmt.Println("error while parsing the icmp request msg")
		log.Fatal(err)
	}
	switch rm.Type {
	case ipv4.ICMPTypeEchoReply, ipv6.ICMPTypeEchoReply:
		fmt.Printf("refelection from %v\n", peer)
	default:
		fmt.Printf("expected %v but got %v\n", ipv6.ICMPTypeEchoReply, rm.Type)
	}
}
