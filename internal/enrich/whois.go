package enrich

import (
	"errors"
	"strings"
	"time"

	"github.com/likexian/whois"
	whoisparser "github.com/likexian/whois-parser"
	"github.com/sjimenezl/phishrivals/internal/models"
)

var ErrWhoisNotFound = errors.New("whois: not registered")

type WhoisClient struct{}

func NewWhoisClient() *WhoisClient {
	return &WhoisClient{}
}

func (w *WhoisClient) LookupWhois(domain string) (*models.DomainInfo, error) {
	raw, err := whois.Whois(domain)
	if err != nil {
		return nil, err
	}
	low := strings.ToLower(raw)
	if strings.Contains(low, "available for registration") ||
		strings.Contains(low, "no match for") ||
		strings.Contains(low, "not found") {
		return nil, ErrWhoisNotFound
	}

	parsed, err := whoisparser.Parse(raw)
	if err != nil {
		return nil, err
	}
	di := &models.DomainInfo{
		Domain:    domain,
		Source:    "whois",
		Registrar: parsed.Registrar.Name,
	}
	if cd := parsed.Domain.CreatedDate; cd != "" {
		if t, err := time.Parse(time.RFC3339, cd); err == nil {
			di.Created = &t
		}
	}
	// abuse
	di.AbuseEmail = parsed.Registrar.Email
	di.AbusePhone = parsed.Registrar.Phone
	return di, nil
}
