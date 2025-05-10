# Phishrivals

> [!IMPORTANT]  
> This project is a work in progress so some functionalites might be slightly broken or missing (e.g. reporting).

A playful Go toolkit for hunting down phishing sites—and learning the ropes of security tooling.

Phishrivals is a side project I started when my iPhone got stolen and I kept getting sketchy “Apple ID” texts. I poked around and found Pastebin and CT logs full of scammy domains.

## Why I Built This

- **Personal itch**: After my phone was stolen, I got some convincing phishing SMSs and figured I could automate the boring bits of takedown.

- **Learning playground**: I wanted to deepen my Go skills, dabble in RDAP/WHOIS, and stitch together a full security pipeline.

- **Interviews**: Having a real-world demo always makes for a better story when interviewing.


## What It Does

1. Ingest suspect domains from:
    - Pastebin/Gist dumps

    - Batch CRT.sh JSON
    - Real-time CertStream feeds

2. Enrich each domain with:

    - RDAP + WHOIS fallback for creation date & registrar info

    - Nameserver lookups

3. Score the badness based on:

    - Age, registrar, cert validity, login-form hits, etc.

4. Report by auto-generating takedown emails to abuse contacts

## Getting started

### Requirements

- Go 1.18+

- Docker (optional, for local CertStream server)

## License

[MIT](https://choosealicense.com/licenses/mit/)