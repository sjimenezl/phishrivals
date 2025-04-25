package main

import (
	"context"
	"log"

	"github.com/sjimenezl/phishrivals/internal/db"
	"github.com/sjimenezl/phishrivals/internal/enrich"
	"github.com/sjimenezl/phishrivals/internal/ingest"
	"github.com/sjimenezl/phishrivals/internal/models"
)

const URL = "https://pastebin.com/raw/fHt0aScX"

func main() {
	db.InitDB()

	if err := models.Migrate(db.DB); err != nil {
		panic(err)
	}

	ctx := context.Background()

	ingestor := ingest.NewIngestor(URL)
	enricher := enrich.NewEnricher()

	urls, err := ingestor.Fetch(ctx)
	if err != nil {
		log.Fatalf("failed to fetch URLs: %v", err)
	}

	for _, u := range urls {
		info, err := enricher.Lookup(u)
		if err != nil {
			log.Printf("skip %s: %v", u, err)
			continue
		}
		// save to db
		db.SaveEnrichment(info)
		// todo: send report email
	}

}
