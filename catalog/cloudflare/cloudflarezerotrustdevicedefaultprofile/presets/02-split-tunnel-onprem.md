# Split tunnel with on-prem DNS

The hybrid-office shape: private LAN ranges and printer discovery stay off the tunnel (exclude mode, plus Microsoft 365 IPs via `exclude_office_ips`), short hostnames complete in the corp domain, and `corp.internal` resolves against the on-prem resolver through the fallback list. The fallback list is FULL-REPLACEMENT -- this preset re-declares `localhost` and `home.arpa` because declaring any row replaces Cloudflare's seeded defaults.
