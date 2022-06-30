package tcprouter

import (
	"fmt"
	"strings"
	"log"
	"net"
)


type ProxyProtoHeader struct {
	InetProtocol string
	InetProtocolFamily int
	SourceAddress string
	DestinationAddress string
	SourcePort int
	DestinationPort int
}


func (p *ProxyProtoHeader) SerializeV1() []byte {
	header := fmt.Sprintf(
		"PROXY %s%d %s %s %d %d\r\n",
		strings.ToUpper(p.InetProtocol), p.InetProtocolFamily, p.SourceAddress,
		p.DestinationAddress, p.SourcePort, p.DestinationPort)
	log.Printf("Proxy Header: %s", header)
	byteHeader := []byte(header)
	return byteHeader
}

func BuildHeader(version int, srcAddr, dstAddr *net.TCPAddr) []byte {
	if (version == 0) {
		log.Printf("Proxy Protocol Disabled")
		return []byte{}
	} else if (version != 1) {
		log.Printf("Proxy Protocol Version: %d is not supported", version)
		return []byte{}
	}
	
	family := 4
	if (srcAddr.IP.To4() == nil) {
		family = 6
	}
	srcIP := srcAddr.IP.String()
	dstIP := dstAddr.IP.String()

	proxyHeader := ProxyProtoHeader{
		InetProtocol: "tcp",
		InetProtocolFamily: family,
		SourceAddress: srcIP,
		DestinationAddress: dstIP,
		SourcePort: srcAddr.Port,
		DestinationPort: dstAddr.Port,
	}

	return proxyHeader.SerializeV1()
}