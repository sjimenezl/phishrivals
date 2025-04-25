package enrich

import (
	"errors"
	"time"

	"github.com/openrdap/rdap"
	"github.com/sjimenezl/phishrivals/internal/models"
)

var ErrNotFound = errors.New("rdap: not found")

type RDAPClient struct {
	client *rdap.Client
}

func NewRDAPClient() *RDAPClient {
	return &RDAPClient{
		client: &rdap.Client{},
	}
}

func (r *RDAPClient) LookupRDAP(domain string) (*models.DomainInfo, error) {
	d, err := r.client.QueryDomain(domain)
	if err != nil {
		if ce, ok := err.(*rdap.ClientError); ok && ce.Type == rdap.ObjectDoesNotExist {
			return nil, ErrNotFound
		}
		return nil, err
	}

	info := &models.DomainInfo{
		Domain: domain,
		Source: "rdap",
	}

	// registration date
	for _, ev := range d.Events {
		if ev.Action == "registration" {
			if t, err := time.Parse(time.RFC3339, ev.Date); err == nil {
				info.Created = &t
			}
			break
		}
	}

	// registrar name
	for _, ent := range d.Entities {
		for _, role := range ent.Roles {
			if role == "registrar" && ent.VCard != nil {
				info.Registrar = ent.VCard.Name()

				for _, rent := range ent.Entities {
					for _, r := range rent.Roles {
						if r == "abuse" && rent.VCard != nil {
							info.AbuseEmail = rent.VCard.Email()
							info.AbusePhone = rent.VCard.Tel()
						}
					}
				}
			}
		}
	}

	// namesservers
	for _, ns := range d.Nameservers {
		info.Nameservers = append(info.Nameservers, models.Nameserver{
			Domain:     domain,
			Nameserver: ns.LDHName,
		})
	}

	return info, nil
}
