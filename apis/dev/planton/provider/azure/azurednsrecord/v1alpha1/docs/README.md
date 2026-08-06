# AzureDnsRecord -- Design Research

## The Resource

One record set in an Azure public DNS zone
(`Microsoft.Network/dnsZones/{A|AAAA|CNAME|MX|SRV|CAA|TXT|NS|PTR}`).
ARM models each record TYPE as its own child recordset of the zone;
azurerm splits that into nine per-type resources
(`internal/services/dns/dns_*_record_resource.go`, DNS API 2018-05-01)
as Terraform ergonomics. The component stays ONE polymorphic kind whose
typed payload selects the recordset type, dispatching to the matching
provider resource on both engines -- parity-verified against
pulumi-azure v6 (`dns.ARecord` ... `dns.PtrRecord`).

## Shape Decision: Typed Payloads, No Discriminator Enum

The record type IS whichever payload message/list is present (an
exactly-one-of message CEL), so a type field that could disagree with
the payload never exists. Each payload carries DNS's real value shape:

| Payload | Shape | azurerm resource |
|---|---|---|
| `a` / `aaaa` | address list XOR alias `target_resource_id` | `dns_a_record` / `dns_aaaa_record` |
| `cname` | single value XOR alias | `dns_cname_record` |
| `mx` | (preference 0-65535, exchange) entries | `dns_mx_record` (string-typed preference, converted) |
| `srv` | (priority, weight, port, target) entries | `dns_srv_record` |
| `caa` | (flags 0-255, tag enum, value) entries | `dns_caa_record` (lowercase wire tags) |
| `txt` | strings 1-4096 chars | `dns_txt_record` (provider chunks to 254-char DNS strings) |
| `ns` | name-server hostnames | `dns_ns_record` |
| `ptr` | pointer hostnames | `dns_ptr_record` |

The flat value-list shape this replaces could not express per-entry MX
preferences or SRV endpoints at all -- any such shape forces the module
to synthesize the missing fields, which is a misdeploy by construction.

## Field Mapping (azurerm -> spec, common surface)

| azurerm | spec | Notes |
|---|---|---|
| `zone_name` + `resource_group_name` | `zone_name` + `resource_group` | Azure's only record-set addressing mode on both engines (no ARM-id input exists). Both FK-defaulted; ForceNew |
| `name` | `name` | `@`/wildcard/underscore-label contract as CEL. ForceNew (MX's azurerm-side `@` default not adopted -- the name is required uniformly) |
| `ttl` | `ttl_seconds` | azurerm requires TTL with no default; the platform default (300) keeps the field ergonomic |
| `tags` | `tags` | Recordset metadata on every type |
| `target_resource_id` | `a`/`aaaa`/`cname` `target_resource_id` | Alias records -- A/AAAA/CNAME only, XOR the literal values (azurerm's contract, front-loaded as CELs). Bare polymorphic ref: no kind dominates alias targets |

## Front-Loaded Contracts

- Exactly one typed payload (message CEL) -- ARM stores each type as its
  own record set.
- A/AAAA: literal addresses XOR alias; addresses are format-validated
  (IPv4/IPv6 -- azurerm validates only AAAA's, the asymmetry is not
  copied).
- CNAME: single value XOR alias (the provider's ExactlyOneOf).
- CAA: the tag vocabulary azurerm accepts (issue/issuewild/iodef/
  contactemail) as a closed enum; flags 0-255.
- MX/SRV/CAA integers are `optional` + required so meaningful zeros
  (null-MX preference 0, SRV priority 0, CAA flags 0) validate.

## Deliberately Not Modeled (with reasons)

- **A standalone SOA record** -- azurerm exposes SOA only as a zone
  block (folded on `AzureDnsZone`) and a data source.
- **DS/TLSA/NAPTR/DNSSEC records** -- absent from azurerm v4 entirely.
- **azurerm's per-type resource split** -- TF ergonomics, not ARM's
  model (the orchestration-mode dispatch precedent).

## Operational Behavior Worth Knowing

- Alias A/AAAA records are the way to point a zone APEX at an Azure
  resource (Public IP, Traffic Manager, CDN/Front Door endpoints) --
  DNS itself forbids CNAME at the apex.
- The zone's own apex NS record set is Azure-managed; the `ns` payload
  is for delegating CHILD subdomains.
- `fqdn` output carries DNS's trailing dot (`www.example.com.`).
