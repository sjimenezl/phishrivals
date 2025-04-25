package enrich

import (
	"fmt"

	"github.com/sjimenezl/phishrivals/internal/models"
)

type Enricher struct {
	rdap  *RDAPClient
	whois *WhoisClient
}

func NewEnricher() *Enricher {
	return &Enricher{
		rdap:  NewRDAPClient(),
		whois: NewWhoisClient(),
	}
}

func (e *Enricher) Lookup(domain string) (*models.DomainInfo, error) {
	info, err := e.rdap.LookupRDAP(domain)
	if err == ErrNotFound {
		info, err = e.whois.LookupWhois(domain)
		if err == ErrWhoisNotFound {
			return nil, fmt.Errorf("domain %q not registered", domain)
		}
	}

	return info, nil
}
