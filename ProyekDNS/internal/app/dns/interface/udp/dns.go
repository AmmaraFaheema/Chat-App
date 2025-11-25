package udp

import (
	"log"
	"strconv"

	"github.com/FhmiSddq/ProyekDNS/internal/app/dns/handler"
	"github.com/FhmiSddq/ProyekDNS/internal/env"
	"github.com/miekg/dns"
)

type DNS struct {
	Env        *env.Env
	DNSHandler handler.DNSHandlerItf
}

func New(env *env.Env, dnsHandler handler.DNSHandlerItf) {
	dns := DNS{
		Env:        env,
		DNSHandler: dnsHandler,
	}

	dns.Start()
}

func (d *DNS) Start() {
	srv := &dns.Server{
		Addr:    ":" + strconv.Itoa(d.Env.Port),
		Net:     "udp",
		Handler: d.DNSHandler,
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Fatalf("failed to set udp listener %s\n", err.Error())
	}
}
