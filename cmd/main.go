package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

const (
	ICMPProtocol   = 1
	ICMPv6Protocol = 58
)

// echo body creator functin
func makeBody(size int) []byte {
	b := make([]byte, size)
	for i := range b {
		b[i] = byte(i & 0xff)
	}
	return b
}
func ping() {
	c, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0.")
	if err != nil {
		fmt.Errorf("cannot open listener connection %w\n", err)
	}
	defer c.Close()
	//body for icmp echo request, which contains id, seq number or data in raw bytes
	echoRequestBody := &icmp.Echo{ID: os.Getpid() & 0xffff, Seq: 1, Data: makeBody(56)}
	icmpEchoMsg := icmp.Message{Type: ipv4.ICMPTypeEcho, Code: 0, Body: echoRequestBody}

	//marshaling the body
	icmpEchoMsgInBytes, err := icmpEchoMsg.Marshal(nil)
	if err != nil {
		fmt.Errorf("could not marshal the icmp echo req, %w\n", err)
	}

	//send the echo msg
	dst := &net.IPAddr{IP: net.ParseIP("8.8.8.8")}
	if _, err := c.WriteTo(icmpEchoMsgInBytes, dst); err != nil {
		fmt.Errorf("write fail to '8.8.8.8':%w\n", err)
	}
	//reading the request
	rb := make([]byte, 1500)
	n, peer, err := c.ReadFrom(rb)
	fmt.Println(n, peer, err)
	if err != nil {
		fmt.Errorf("couldnot read the echo reply message: %w\n", err)
	}
	rm, err := icmp.ParseMessage(ICMPProtocol, rb[:n])
	if err != nil {
		fmt.Println("error while parsing the icmp request msg")
		log.Fatal(err)
	}
	switch rm.Type {
	case ipv4.ICMPTypeEchoReply:
		fmt.Printf("refelection from %v\n", peer)
	default:
		fmt.Printf("expected %v but got %v\n", ipv4.ICMPTypeEchoReply, rm.Type)
	}
}

func main() {
	ping()
}
