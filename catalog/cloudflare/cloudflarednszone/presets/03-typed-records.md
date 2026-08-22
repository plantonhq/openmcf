# Typed Records

Bootstraps a zone with its full record portfolio inline: a proxied web A
record, mail MX with an SPF TXT, a structured SRV service record, and a CAA
record controlling which certificate authority may issue for the domain.
Demonstrates both value paths -- simple records through `content` and
structured records through their typed blocks (`srv`, `caa`).

## When to Use

- Standing up a complete zone in one apply: web, mail, service discovery, and issuance policy together
- Migrating a domain whose record set is small and owned by one team (records share the zone's lifecycle)
- A template for the structured-record syntax (SRV/CAA here; CERT, DNSKEY, DS, HTTPS, LOC, NAPTR, SMIMEA, SSHFP, SVCB, TLSA, URI follow the same pattern)

## Key Configuration Choices

- **Proxied web record** (`records[0].proxied: true`) -- traffic flows through Cloudflare's CDN/WAF; TTL stays automatic.
- **MX priority** (`records[1].priority`) -- top-level priority is an MX concern; SRV carries its own inside `srv`.
- **Structured SRV** (`records[3].srv`) -- typed block instead of `content`; exactly one of the two is set per record.
- **CAA issuance policy** (`records[4].caa`) -- `tag: issue` restricts certificate issuance to the named CA.

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|-------------|-------------|---------------|
| `<your-domain.com>` | Fully qualified domain for the zone | Your registered domain |
| `<cloudflare-account-id>` | Cloudflare account ID | Cloudflare Dashboard → Overview → Account ID (right sidebar) |
| `<web-server-ipv4>` | Origin IPv4 for the proxied web record | Your origin infrastructure |
| `<mail-host.your-domain.com>` | Mail exchanger hostname | Your mail provider's setup docs |
| `<mail-provider.com>` | SPF include domain | Your mail provider's SPF docs |
| `<sip-host.your-domain.com>` | SRV target host | Your service's deployment |
| `<ca-domain.org>` | Certificate authority domain (e.g. letsencrypt.org) | Your certificate issuance policy |

## Related Presets

- **01-basic-zone** -- A clean zone with no inline records
- **02-dnssec-signed** -- A zone with DNSSEC enabled and the DS material exported
