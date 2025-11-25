// Package bootstrap calls and load modules
package bootstrap

import (
	"log"

	dns "github.com/FhmiSddq/ProyekDNS/internal/app/dns/handler"
	"github.com/FhmiSddq/ProyekDNS/internal/app/dns/interface/udp"
	"github.com/FhmiSddq/ProyekDNS/internal/env"
)

func Start() {
	env := env.New()

	dns := dns.New()

	log.Println("app started")

	dns.Register("google.com", "8.8.8.8")
	dns.Register("router", "192.168.100.153")
	dns.Register("peer1", "10.0.0.2")
	dns.Register("peer2", "10.0.0.3")
	dns.Register("peer3", "10.0.0.4")

	udp.New(env, dns)
}
