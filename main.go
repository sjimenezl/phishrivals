package main

import (
	"fmt"
	"log"

	"github.com/sjimenezl/phishrivals/internal/config"
	"github.com/sjimenezl/phishrivals/internal/db"
	"github.com/sjimenezl/phishrivals/internal/enrich"
	"github.com/sjimenezl/phishrivals/internal/helper"
	"github.com/sjimenezl/phishrivals/internal/ingest/crtsh"
	"github.com/sjimenezl/phishrivals/internal/models"
)

// const URL = "https://pastebin.com/raw/fHt0aScX"
const RISK_THRESHOLD = 0.5

func main() {
	db.InitDB()

	if err := models.Migrate(db.DB); err != nil {
		panic(err)
	}

	//TODO: check if i want the pastebin fetcher
	// ctx := context.Background()
	// urls, err := ingestor.Fetch(ctx)
	enricher := enrich.NewEnricher()

	cfg, err := config.Load("sus.yaml")
	if err != nil {
		log.Fatalf("couldn't load config: %v", err)
	}

	crtDomains, err := crtsh.Fetch(cfg.Keywords)
	if err != nil {
		log.Fatalf("failed to fetch crtsh URLs: %v", err)
	}

	// check if the domains are alive
	for _, domain := range crtDomains {
		if !helper.IsAlive(domain) {
			continue
		}

		// check for cert
		info, err := enricher.Lookup(domain)
		if err != nil {
			log.Printf("skip %s: %v", domain, err)
			continue
		}

		// check threat score
		score := helper.ThreatScore(info)
		if score < RISK_THRESHOLD {
			continue
		}

		// found high risk, feed into DB and takedown pipeline
		fmt.Println(domain)
		fmt.Println(score)

		// save to db
		// db.SaveEnrichment(info)
		// todo: send report email
	}

}
