//go:build !linux

package tproxy

import (
	"net"
)

// ListenUDP will construct a new UDP listener
// socket with the Linux IP_TRANSPARENT option
// set on the underlying socket
func ListenUDP(network string, laddr *net.UDPAddr) (*net.UDPConn, error) {
	return nil, nil
}

// DialUDP connects to the remote address raddr on the network net,
// which must be "udp", "udp4", or "udp6".  If laddr is not nil, it is
// used as the local address for the connection.
func DialUDP(network string, laddr *net.UDPAddr, raddr *net.UDPAddr) (*net.UDPConn, error) {
	return nil, nil
}

// ReadFromUDP reads a UDP packet from c, copying the payload into b.
// It returns the number of bytes copied into b and the return address
// that was on the packet.
//
// Out-of-band data is also read in so that the original destination
// address can be identified and parsed.
func ReadFromUDP(conn *net.UDPConn, b []byte) (int, *net.UDPAddr, *net.UDPAddr, error) {
	return 0, nil, nil, nil
}
