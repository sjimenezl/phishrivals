package ingest

import (
	"bufio"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/publicsuffix"
)

type Ingestor struct {
	Client   *http.Client
	FeedURL  string
	Keywords []string
}

func NewIngestor(feedURL string, keywords []string) *Ingestor {
	return &Ingestor{
		Client:   &http.Client{Timeout: 10 * time.Second},
		FeedURL:  feedURL,
		Keywords: keywords,
	}
}

func (i *Ingestor) Fetch(ctx context.Context) ([]string, error) {
	// 1) Fetch the raw pastebin URL
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, i.FeedURL, nil)
	if err != nil {
		return nil, err
	}

	res, err := i.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	scanner := bufio.NewScanner(res.Body)
	seen := make(map[string]struct{})
	var domains []string

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		// split off anything after a space, take the first token
		parts := strings.SplitN(line, " ", 2)
		raw := parts[0]

		// parse & strip www.
		u, err := url.Parse(raw)
		if err != nil {
			// you might want to log and continue instead of panic
			fmt.Printf("bad URL %q: %v\n", raw, err)
			continue
		}
		host := strings.TrimPrefix(u.Host, "www.")

		// normalize to registered domain
		base, err := publicsuffix.EffectiveTLDPlusOne(host)
		if err != nil {
			fmt.Printf("cannot normalize %q: %v\n", host, err)
			continue
		}

		// skip duplicates
		if _, ok := seen[base]; ok {
			continue
		}
		seen[base] = struct{}{}
		domains = append(domains, base)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return domains, nil
}
