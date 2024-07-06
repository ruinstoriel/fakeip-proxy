package listener

import (
	"context"
	"errors"
	"fakeip-proxy/dns"
	"net"
	"strconv"
	"strings"
	"time"

	"fakeip-proxy/tproxy"
)

const (
	udpBufferSize  = 4096
	defaultTimeout = 60 * time.Second
)

type UDPTProxy struct {
	Timeout     time.Duration
	EventLogger UDPEventLogger
}

type UDPEventLogger interface {
	Connect(addr, reqAddr net.Addr)
	Error(addr, reqAddr net.Addr, err error)
}

func (r *UDPTProxy) ListenAndServe(laddr *net.UDPAddr) error {
	conn, err := tproxy.ListenUDP("udp", laddr)
	if err != nil {
		return err
	}
	defer conn.Close()
	buf := make([]byte, udpBufferSize)
	for {
		// We will only get the first packet of each src/dst pair here,
		// because newPair will create a TProxy connection and take over
		// the src/dst pair. Later packets will be sent there instead of here.
		n, srcAddr, dstAddr, err := tproxy.ReadFromUDP(conn, buf)
		if err != nil {
			return err
		}
		r.newPair(srcAddr, dstAddr, buf[:n])
	}
}

func (r *UDPTProxy) newPair(srcAddr, dstAddr *net.UDPAddr, initPkt []byte) {
	if r.EventLogger != nil {
		r.EventLogger.Connect(srcAddr, dstAddr)
	}
	var closeErr error
	defer func() {
		// If closeErr is nil, it means we at least successfully sent the first packet
		// and started forwarding, in which case we don't call the error logger.
		if r.EventLogger != nil && closeErr != nil {
			r.EventLogger.Error(srcAddr, dstAddr, closeErr)
		}
	}()
	conn, err := tproxy.DialUDP("udp", dstAddr, srcAddr)
	if err != nil {
		closeErr = err
		return
	}
	ip := strings.Split(conn.LocalAddr().String(), ":")[0]
	port, _ := strconv.Atoi(strings.Split(conn.LocalAddr().String(), ":")[1])
	domain, b := dns.Cache.Get(int(dns.Ip2int(net.ParseIP(ip))))
	if b {
		realIp, err := net.DefaultResolver.LookupIP(context.Background(), "udp", domain)
		if err != nil {
			return
		}
		hyConn, err := net.DialUDP("udp", nil, &net.UDPAddr{
			IP:   realIp[0],
			Port: port,
		})
		if err != nil {
			_ = conn.Close()
			closeErr = err
			return
		}
		// Send the first packet
		if err != nil {
			_ = conn.Close()
			_ = hyConn.Close()
			closeErr = err
			return
		}
		// Start forwarding
		go func() {
			err := r.forwarding(conn, hyConn)
			_ = conn.Close()
			_ = hyConn.Close()
			if r.EventLogger != nil {
				var netErr net.Error
				if errors.As(err, &netErr) && netErr.Timeout() {
					// We don't consider deadline exceeded (timeout) an error
					err = nil
				}
				r.EventLogger.Error(srcAddr, dstAddr, err)
			}
		}()
	}

}

func (r *UDPTProxy) forwarding(conn *net.UDPConn, hyConn *net.UDPConn) error {
	errChan := make(chan error, 2)
	// Local <- Remote
	go func() {
		for {
			oob := make([]byte, udpBufferSize)
			n, _, err := hyConn.ReadFromUDP(oob)
			if err != nil {
				errChan <- err
				return
			}
			_, err = conn.Write(oob[:n])
			if err != nil {
				errChan <- err
				return
			}
			_ = r.updateConnDeadline(conn)
		}
	}()
	// Local -> Remote
	go func() {
		buf := make([]byte, udpBufferSize)
		for {
			_ = r.updateConnDeadline(conn)
			n, err := conn.Read(buf)
			if n > 0 {
				_, err := hyConn.Write(buf[:n])
				if err != nil {
					errChan <- err
					return
				}
			}
			if err != nil {
				errChan <- err
				return
			}
		}
	}()
	return <-errChan
}

func (r *UDPTProxy) updateConnDeadline(conn *net.UDPConn) error {
	if r.Timeout == 0 {
		return conn.SetReadDeadline(time.Now().Add(defaultTimeout))
	} else {
		return conn.SetReadDeadline(time.Now().Add(r.Timeout))
	}
}
