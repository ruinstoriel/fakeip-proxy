package listener

import (
	"net"
)

func Start() {
	tcp := TCPTProxy{}
	addr := net.TCPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: 9898,
	}
	err := tcp.ListenAndServe(&addr)
	if err != nil {
		return
	}
	udpAddr := net.UDPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: 9898,
	}
	udp := UDPTProxy{}
	err = udp.ListenAndServe(&udpAddr)
	if err != nil {
		return
	}
}
