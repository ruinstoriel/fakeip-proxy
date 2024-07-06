//go:build !linux

package tproxy

import "net"

func ListenTCP(network string, laddr *net.TCPAddr) (net.Listener, error) {

	return nil, nil
}
