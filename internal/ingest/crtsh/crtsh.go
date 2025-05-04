package crtsh

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func Fetch(keywords []string) ([]string, error) {
	seen := make(map[string]struct{})

	for _, kw := range keywords {
		queryURL := fmt.Sprintf(
			"https://crt.sh/?q=%%25.%s&output=json",
			kw,
		)

		res, err := http.Get(queryURL)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()

		var entries []struct {
			NameValue string `json:"name_value"`
		}
		body, _ := io.ReadAll(res.Body)
		if err := json.Unmarshal(body, &entries); err != nil {
			return nil, err
		}

		// split on newlines & commas, dedupe
		for _, e := range entries {
			for _, part := range strings.FieldsFunc(e.NameValue, func(r rune) bool {
				return r == ',' || r == '\n'
			}) {
				part = strings.TrimSpace(part)
				for _, kw := range keywords {
					if strings.Contains(part, kw) {
						seen[part] = struct{}{}
						break
					}
				}
			}
		}
	}

	// unique domains
	out := make([]string, 0, len(seen))
	for d := range seen {
		out = append(out, d)
	}
	return out, nil
}
