package helper

import (
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/sjimenezl/phishrivals/internal/models"
)

func IsAlive(domain string) bool {
	// dns lookup
	ips, err := net.LookupHost(domain)
	if err != nil || len(ips) == 0 {
		return false
	}

	// https check
	client := &http.Client{Timeout: 3 * time.Second}
	res, err := client.Get("https://" + domain)
	if err != nil {
		return false
	}

	res.Body.Close()
	return true
}

func ThreatScore(info *models.DomainInfo) float64 {
	var score float64

	// domain age < 7 days (0.2)
	if info.Created != nil && time.Since(*info.Created) < 7*24*time.Hour {
		score += 0.2
	}

	// wildcard or punny domain (0.1)
	// eg. starts with * or contains --
	if strings.HasPrefix(info.Domain, "*.") ||
		strings.HasPrefix(info.Domain, "--") ||
		strings.HasPrefix(info.Domain, "xn--") {
		score += 0.1
	}

	if score > 1.0 {
		score = 1.0
	}

	return score
}
