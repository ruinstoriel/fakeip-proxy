package dns

import (
	"context"
	"encoding/binary"
	"fakeip-proxy/v2geo"
	"fmt"
	"github.com/go-faster/city"
	"github.com/hashicorp/golang-lru/v2/expirable"
	"github.com/spf13/viper"
	"log"
	"net"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/miekg/dns"
)

var Cache = expirable.NewLRU[int, string](1000, nil, time.Minute*1)

var Site, _ = v2geo.LoadGeoSite("/root/dns/geosite.dat")

func parseQuery(m *dns.Msg) {

	for _, q := range m.Question {
		switch q.Qtype {
		case dns.TypeA:
			domain := q.Name[0 : len(q.Name)-1]
			log.Printf("Query for %s\n", domain)

			domains := Site["google"].Domain
			b := Match(domains, domain)

			if !b {
				ip, err := net.DefaultResolver.LookupIP(context.Background(), "ip4", q.Name)
				if err != nil {
					fmt.Println(err)
					return
				}

				rr, err := dns.NewRR(fmt.Sprintf("%s A %s", q.Name, ip[0]))
				if err == nil {
					m.Answer = append(m.Answer, rr)
				}
				return

			}
			prefix := strconv.Itoa(viper.Get("prefix").(int))
			ip := int2ip(Ip2int(net.ParseIP(prefix+".19.0.0")) + city.Hash32([]byte(q.Name))>>16).String()
			if ip != "" {
				rr, err := dns.NewRR(fmt.Sprintf("%s A %s", q.Name, ip))
				if err == nil {
					m.Answer = append(m.Answer, rr)
					Cache.Add(int(Ip2int(net.ParseIP(ip))), q.Name)
				}
			}
		}
	}
}

func Match(m []*v2geo.Domain, host string) bool {
	for _, domain := range m {
		if matchDomain(domain, host) {
			return true
		}
	}
	return false
}
func matchDomain(domain *v2geo.Domain, host string) bool {

	switch domain.Type {
	case v2geo.Domain_Plain:
		return strings.Contains(host, domain.Value)
	case v2geo.Domain_Regex:
		regex, err := regexp.Compile(domain.Value)
		if err != nil {
			return false
		}
		return regex.MatchString(host)
	case v2geo.Domain_Full:
		return host == domain.Value
	case v2geo.Domain_RootDomain:
		if host == domain.Value {
			return true
		}
		return strings.HasSuffix(host, "."+domain.Value)
	default:
		return false
	}
	return false
}
func handleDnsRequest(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(r)
	m.Compress = false

	switch r.Opcode {
	case dns.OpcodeQuery:
		parseQuery(m)
	}

	w.WriteMsg(m)
}

func Ip2int(ip net.IP) uint32 {
	if len(ip) == 16 {
		return binary.BigEndian.Uint32(ip[12:16])
	}
	return binary.BigEndian.Uint32(ip)
}

func int2ip(nn uint32) net.IP {
	ip := make(net.IP, 4)
	binary.BigEndian.PutUint32(ip, nn)
	return ip
}

func Start() {

	// attach request handler func
	dns.HandleFunc(".", handleDnsRequest)

	ipPrefix := viper.Get("prefix").(int)
	if ipPrefix > 255 || ipPrefix < 0 {
		panic("Error: the prefix range must be between 0 and 255")
	}
	// start server

	server := &dns.Server{Addr: ":" + strconv.Itoa(53), Net: "udp"}
	log.Printf("Starting at %d\n", 53)
	err := server.ListenAndServe()
	defer server.Shutdown()
	if err != nil {
		log.Fatalf("Failed to start server: %s\n ", err.Error())
	}
}
