# Cloudflare DNS Record

Provision and manage individual DNS records in Cloudflare zones using OpenMCF's unified API.

## Overview

Cloudflare DNS provides authoritative DNS served from 330+ global locations via native anycast, with built-in DDoS protection, zero per-query charges, and optional integrated CDN/WAF/proxy capabilities. This component enables you to create individual DNS records within a Cloudflare-managed zone.

DNS records are the fundamental building blocks of your domain's DNS configuration—mapping hostnames to IP addresses (A/AAAA records), creating aliases (CNAME), routing email (MX), and storing verification data (TXT).

This component provides a clean, protobuf-defined API for provisioning DNS records, following the **80/20 principle**: exposing only the essential configuration fields that cover the most common use cases.

## Key Features

- **All Major Record Types**: Support for A, AAAA, CNAME, MX, TXT, SRV, NS, and CAA records
- **Orange Cloud Integration**: Optional Cloudflare proxy (CDN/WAF) for A, AAAA, and CNAME records
- **TTL Control**: Automatic or custom TTL settings
- **Priority Support**: Required for MX records, optional for SRV
- **Comment Support**: Document record purpose (up to 100 characters)
- **Validation**: Built-in validation for record types, TTL ranges, and cross-field rules

## Prerequisites

1. **Cloudflare DNS Zone**: An existing zone where records will be created (use CloudflareDnsZone component)
2. **Zone ID**: The Cloudflare Zone ID (from CloudflareDnsZone outputs or dashboard)
3. **API Token**: Cloudflare API token with DNS:Edit permissions
4. **OpenMCF CLI**: Install from [openmcf.org](https://openmcf.org)

## Quick Start

### A Record (IPv4)

Point a subdomain to an IPv4 address:

```yaml
apiVersion: cloudflare.openmcf.org/v1
kind: CloudflareDnsRecord
metadata:
  name: www-a-record
spec:
  zone_id: "your-zone-id-here"
  name: "www"
  type: A
  value: "192.0.2.1"
  proxied: true
```

Deploy:

```bash
planton apply -f record.yaml
```

### CNAME Record

Create an alias to another hostname:

```yaml
apiVersion: cloudflare.openmcf.org/v1
kind: CloudflareDnsRecord
metadata:
  name: app-cname
spec:
  zone_id: "your-zone-id-here"
  name: "app"
  type: CNAME
  value: "www.example.com"
  proxied: true
```

### MX Record (Email)

Route email to your mail server:

```yaml
apiVersion: cloudflare.openmcf.org/v1
kind: CloudflareDnsRecord
metadata:
  name: mx-primary
spec:
  zone_id: "your-zone-id-here"
  name: "@"
  type: MX
  value: "mail.example.com"
  priority: 10
```

### TXT Record (SPF, DKIM, Verification)

Add SPF for email authentication:

```yaml
apiVersion: cloudflare.openmcf.org/v1
kind: CloudflareDnsRecord
metadata:
  name: spf-record
spec:
  zone_id: "your-zone-id-here"
  name: "@"
  type: TXT
  value: "v=spf1 include:_spf.google.com ~all"
```

## Configuration Reference

### Required Fields

| Field | Type | Description |
|-------|------|-------------|
| `zone_id` | string | Cloudflare Zone ID where the record will be created |
| `name` | string | Record name (e.g., "www", "@" for root, "api") |
| `type` | enum | DNS record type: A, AAAA, CNAME, MX, TXT, SRV, NS, CAA |
| `value` | string | Record value (IP address, hostname, or text) |

### Optional Fields

| Field | Type | Description | Default |
|-------|------|-------------|---------|
| `proxied` | bool | Route through Cloudflare CDN/WAF (A, AAAA, CNAME only) | false |
| `ttl` | int32 | Time to live: 1 (auto) or 60-86400 seconds | 1 (auto) |
| `priority` | int32 | Priority for MX/SRV records (0-65535) | 0 |
| `comment` | string | Optional note (max 100 characters) | "" |

### Record Types

| Type | Description | Example Value |
|------|-------------|---------------|
| **A** | IPv4 address | `192.0.2.1` |
| **AAAA** | IPv6 address | `2001:db8::1` |
| **CNAME** | Canonical name (alias) | `www.example.com` |
| **MX** | Mail exchange | `mail.example.com` |
| **TXT** | Text record | `v=spf1 include:...` |
| **SRV** | Service locator | `0 5 5269 xmpp.example.com` |
| **NS** | Nameserver | `ns1.example.com` |
| **CAA** | Certificate Authority Authorization | `0 issue "letsencrypt.org"` |

## Outputs

After deployment, the following outputs are available:

- `record_id`: Unique identifier of the created DNS record
- `hostname`: Fully qualified hostname (e.g., "www.example.com")
- `record_type`: The DNS record type that was created
- `proxied`: Whether the record is proxied through Cloudflare

Access outputs:

```bash
planton output record_id
planton output hostname
```

## Orange Cloud vs Grey Cloud

Cloudflare's unique feature is the **proxy toggle**:

- **Orange Cloud (proxied: true)**: Traffic flows through Cloudflare
  - Enables: CDN caching, WAF, DDoS protection, SSL termination
  - Hides origin IP address
  - Use for: Web services (www, app, api)

- **Grey Cloud (proxied: false)**: DNS-only resolution
  - Direct connection to origin server
  - Origin IP is visible
  - Use for: Email (MX), SSH, VPN, non-HTTP services

**Important**: Only A, AAAA, and CNAME records can be proxied. Other types (MX, TXT, NS, etc.) are always grey-cloud.

## Common Use Cases

### 1. Web Server

```yaml
spec:
  zone_id: "zone-123"
  name: "www"
  type: A
  value: "192.0.2.1"
  proxied: true  # CDN + protection
  comment: "Primary web server"
```

### 2. API Endpoint

```yaml
spec:
  zone_id: "zone-123"
  name: "api"
  type: CNAME
  value: "api-lb.example.com"
  proxied: true  # WAF protection
```

### 3. Root Domain (Apex)

```yaml
spec:
  zone_id: "zone-123"
  name: "@"
  type: A
  value: "192.0.2.1"
  proxied: true
```

### 4. Email Configuration

Primary MX record:
```yaml
spec:
  zone_id: "zone-123"
  name: "@"
  type: MX
  value: "mail.example.com"
  priority: 10
```

Backup MX record:
```yaml
spec:
  zone_id: "zone-123"
  name: "@"
  type: MX
  value: "backup.mail.example.com"
  priority: 20
```

### 5. Domain Verification

```yaml
spec:
  zone_id: "zone-123"
  name: "@"
  type: TXT
  value: "google-site-verification=abc123..."
```

### 6. DKIM for Email

```yaml
spec:
  zone_id: "zone-123"
  name: "google._domainkey"
  type: TXT
  value: "v=DKIM1; k=rsa; p=MIGfMA0GCS..."
```

### 7. CAA Record

Restrict certificate issuance:
```yaml
spec:
  zone_id: "zone-123"
  name: "@"
  type: CAA
  value: "0 issue \"letsencrypt.org\""
```

## Best Practices

1. **Use Proxied Records for Web Traffic**: Enable orange cloud for A/AAAA/CNAME web records
2. **Document Your Records**: Use the comment field to explain each record's purpose
3. **Set Appropriate TTLs**: Use automatic (1) for proxied records, lower TTLs (60-300) for records that change frequently
4. **MX Priority Ordering**: Use 10, 20, 30 for primary, secondary, tertiary mail servers
5. **SPF Records**: Limit to one TXT record per domain with SPF
6. **Test Before Production**: Verify records resolve correctly with `dig` before DNS propagation

## Testing Records

Before relying on DNS propagation:

```bash
# Test directly against Cloudflare nameservers
dig @gina.ns.cloudflare.com www.example.com A

# Check MX records
dig @gina.ns.cloudflare.com example.com MX

# Verify TXT records
dig @gina.ns.cloudflare.com example.com TXT
```

## Troubleshooting

### "Record Already Exists" Error
Cloudflare may have a conflicting record. Check the dashboard or import existing records.

### Proxied Record Not Working
Ensure the record type supports proxying (only A, AAAA, CNAME). MX and TXT cannot be proxied.

### TTL Appears Different Than Set
Proxied records always show TTL as "Auto" because Cloudflare manages caching.

### Email Not Working After Adding MX
Ensure you're not proxying the MX record's target. Also verify SPF/DKIM TXT records.

## Examples

For detailed usage examples, see [examples.md](examples.md).

## Architecture Details

For in-depth architectural guidance and production best practices, see [docs/README.md](docs/README.md).

## Terraform and Pulumi

This component supports both Pulumi (default) and Terraform:

- **Pulumi**: `iac/pulumi/` - Go-based implementation
- **Terraform**: `iac/tf/` - HCL-based implementation

Both produce identical infrastructure. Choose based on your team's preference.

## Support

- **Documentation**: [docs/README.md](docs/README.md)
- **Cloudflare DNS Docs**: [developers.cloudflare.com/dns](https://developers.cloudflare.com/dns)
- **OpenMCF**: [openmcf.org](https://openmcf.org)

## License

This component is part of OpenMCF and follows the same license.
