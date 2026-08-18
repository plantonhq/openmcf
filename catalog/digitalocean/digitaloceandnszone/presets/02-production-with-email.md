# Production Zone with Email

This preset creates a production-shaped DNS zone: website records (apex A + `www`), Google Workspace mail routing (primary and backup MX with priorities), an SPF policy authorizing Google to send mail for the domain, and CAA pinning so only Let's Encrypt may issue certificates.

## When to Use

- A production domain that serves a website AND receives email
- Any zone where certificate issuance should be restricted to a known CA

## Key Configuration Choices

- **Two MX records with priorities** — mail goes to `aspmx.l.google.com` (priority 1) and falls back to `alt1` (priority 5); lower priority wins. MX records require a priority (validation enforces it).
- **SPF as TXT at the apex** — `~all` soft-fails senders outside Google's ranges; tighten to `-all` once mail flow is verified.
- **CAA with `flags: 0` and `tag: issue`** — certificate authorities other than Let's Encrypt must refuse to issue for this domain. Add an `issuewild` record to control wildcard issuance separately.
- **Trailing dots on mail targets** — DigitalOcean stores MX targets fully qualified; authoring the dot avoids a permanent diff.

## Placeholders to Replace

- `metadata.name` — your zone resource's name.
- `domainName` (`example.com` is a documentation example) — your domain.
- The A record's `values` (`203.0.113.10` is a documentation example) — your server's public IPv4 address.
- The MX targets and SPF include — your mail provider's values if not Google Workspace.
- The CAA value (`letsencrypt.org`) — your certificate authority.

## Related Presets

- **01-simple-website** — the website-only starting point without mail or CAA.
