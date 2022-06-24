package tcprouter

import (
	"testing"
	"bytes"
	"net"
)

func TestProxyProtocolV1Serialize(t *testing.T) {
	h := ProxyProtoHeader{
		InetProtocol: "tcp",
		InetProtocolFamily: 4,
		SourceAddress: "192.168.1.100",
		DestinationAddress: "10.100.1.200",
		SourcePort:	53863,
		DestinationPort: 443,
	}
	bresult := h.SerializeV1()
	if (len(bresult) == 0) {
		t.Errorf("serialization failed")
	}
}


func TestProxyProtocolBuildHeader(t *testing.T) {
	srcTCPAddr := &net.TCPAddr{IP: net.ParseIP("192.168.1.100"), Port: 53863}
	dstTCPAddr := &net.TCPAddr{IP: net.ParseIP("10.100.1.200"), Port: 443}
	bresult1 := BuildHeader(1, srcTCPAddr, dstTCPAddr)
	h := ProxyProtoHeader{
		InetProtocol: "tcp",
		InetProtocolFamily: 4,
		SourceAddress: "192.168.1.100",
		DestinationAddress: "10.100.1.200",
		SourcePort:	53863,
		DestinationPort: 443,
	}
	bresult2 := h.SerializeV1()

	if !bytes.Equal(bresult1, bresult2) {
		t.Errorf("header build is not correct")
	}
}
