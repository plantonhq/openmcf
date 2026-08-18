# DigitalOcean DNS Zone

A DigitalOcean-hosted DNS zone described once in a Planton manifest: the domain itself, an inline list of managed records covering every type the DigitalOcean API accepts (A, AAAA, CNAME, MX, TXT, SRV, NS, CAA, SOA), and the create-only apex-A convenience. Adding a domain does not require owning it — the zone serves on DigitalOcean's name servers immediately and resolves publicly once the registrar delegates.

## What this component models

The spec maps onto DigitalOcean's `digitalocean_domain` plus one `digitalocean_record` per managed record value:

| Spec field | What it controls |
|---|---|
| `domainName` | The zone's domain (unique across ALL DigitalOcean accounts); changing it recreates the zone |
| `records` | The zone's managed records — each entry's `values` list fans out to one DigitalOcean record per value |
| `ipAddress` | Create-only convenience: seeds an apex A record DigitalOcean never tracks afterwards; prefer `records` |

Each record entry carries `name` (`@` for the apex), `type` (the shared record-type enum, restricted to what DigitalOcean accepts — ALIAS and PTR are rejected at validation time), one or more `values` (literals or references to other resources' outputs), an optional `ttlSeconds` (DigitalOcean defaults to 1800), and the per-type fields with the same validation floor as the standalone `DigitalOceanDnsRecord` kind: MX needs `priority`, SRV needs `priority`+`weight`+`port`, CAA needs `flags`+`tag`.

## Quick start

The smallest real zone:

```yaml
apiVersion: digital-ocean.planton.dev/v1alpha1
kind: DigitalOceanDnsZone
metadata:
  name: my-zone
spec:
  domainName: example.com
```

A website zone with records:

```yaml
spec:
  domainName: example.com
  records:
    - name: "@"
      type: A
      values:
        - value: 203.0.113.10
      ttlSeconds: 3600
    - name: www
      type: CNAME
      values:
        - value: example.com.
    - name: "@"
      type: MX
      values:
        - value: aspmx.l.google.com.
      priority: 1
```

## Behavior worth knowing

- **Delegation is the registrar's half** — the zone works inside DigitalOcean immediately, but public resolution starts when the registrar points at `ns1`/`ns2`/`ns3.digitalocean.com` (the `name_servers` output).
- **Domain names are globally unique across DigitalOcean** — adding a domain another account already holds fails at create.
- **Multi-value fan-out** — an entry with two `values` creates two records of the same name and type (e.g. round-robin A records).
- **Hostname values read back with a trailing dot** — author CNAME/MX/NS/SRV/CAA targets fully qualified (`mail.example.com.`).
- **`ipAddress` is a footgun by design** — the A record it seeds is invisible to later `records` edits; it exists for migrating configurations that already use it.
- **Zones import by name** — the domain name IS the resource ID.

## Outputs

| Output | Meaning |
|---|---|
| `zone_name` | The domain name — what DNS records reference (`status.outputs.zone_name`) |
| `zone_id` | The zone's resource identifier — the domain name itself, not a UUID |
| `name_servers` | DigitalOcean's fixed authoritative set — what the registrar must delegate to |
| `urn` | The uniform resource name (`do:domain:example.com`) |

## See also

- `GUIDE.md` — operational judgment calls (inline vs. standalone records, delegation, apex rules)
- `presets/` — simple-website and production-with-email starting points
- `v1alpha1/reference.md` — the generated field-by-field contract
