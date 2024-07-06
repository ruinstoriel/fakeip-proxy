//go:build linux

package tproxy

import (
	"fmt"
	"net"
	"syscall"
)

const (
	IPV6_TRANSPARENT     = 0x4b
	IPV6_RECVORIGDSTADDR = 0x4a
)

func ListenTCP(network string, laddr *net.TCPAddr) (net.Listener, error) {
	listener, err := net.ListenTCP(network, laddr)
	if err != nil {
		return nil, err
	}

	fileDescriptorSource, err := listener.File()
	if err != nil {
		return nil, &net.OpError{Op: "listen", Net: network, Source: nil, Addr: laddr, Err: fmt.Errorf("get file descriptor: %s", err)}
	}
	defer fileDescriptorSource.Close()

	if err = syscall.SetsockoptInt(int(fileDescriptorSource.Fd()), syscall.SOL_IP, syscall.IP_TRANSPARENT, 1); err != nil {
		return nil, &net.OpError{Op: "listen", Net: network, Source: nil, Addr: laddr, Err: fmt.Errorf("set socket option: IP_TRANSPARENT: %s", err)}
	}

	return listener, nil
}
