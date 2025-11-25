package handler

import (
	"fmt"
	"net"

	"github.com/miekg/dns"
)

type DNSHandlerItf interface {
	ServeDNS(w dns.ResponseWriter, r *dns.Msg)
	Register(domain string, address string)
	Deregister(domain string)
}

type DNSHandler struct {
	address map[string]string
}

func New() DNSHandlerItf {
	dnsHandler := DNSHandler{
		address: make(map[string]string),
	}

	return &dnsHandler
}

func (d *DNSHandler) ServeDNS(writer dns.ResponseWriter, request *dns.Msg) {
	msg := dns.Msg{}
	msg.SetReply(request)

	switch request.Question[0].Qtype {
	case dns.TypeA:
		msg.Authoritative = true
		domain := msg.Question[0].Name
		address, ok := d.address[domain]
		if ok {
			msg.Answer = append(
				msg.Answer,
				&dns.A{
					Hdr: dns.RR_Header{
						Name:   domain,
						Rrtype: dns.TypeA,
						Class:  dns.ClassINET,
						Ttl:    60,
					},
					A: net.ParseIP(address),
				})
		}
	}

	writer.WriteMsg(&msg)
}

func (d *DNSHandler) Register(domain string, address string) {
	domainParsed := fmt.Sprintf("%s.", domain)

	d.address[domainParsed] = address
}

func (d *DNSHandler) Deregister(domain string) {
	domainParsed := fmt.Sprintf("%s.", domain)

	delete(d.address, domainParsed)
}
