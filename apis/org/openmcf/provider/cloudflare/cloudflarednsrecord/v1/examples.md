# Cloudflare DNS Record Examples

Concrete, copy-and-paste examples for common Cloudflare DNS record deployment scenarios.

## Zone ID Configuration

The `zoneId` field supports two formats:

### 1. Direct Value (String Literal)

Provide the zone ID directly:

```yaml
spec:
  zoneId: "abc123def456"
```

### 2. Reference to CloudflareDnsZone Resource

Reference an existing CloudflareDnsZone resource by name:

```yaml
spec:
  zoneId:
    valueFrom:
      name: my-dns-zone
```

This automatically resolves to `status.outputs.zone_id` of the referenced CloudflareDnsZone resource.

---

## Table of Contents

- [A Record (IPv4)](#a-record-ipv4)
- [AAAA Record (IPv6)](#aaaa-record-ipv6)
- [CNAME Record](#cname-record)
- [MX Records (Email)](#mx-records-email)
- [TXT Records](#txt-records)
- [CAA Record](#caa-record)
- [Root Domain (Apex)](#root-domain-apex)
- [Wildcard Subdomain](#wildcard-subdomain)
- [Development Environment](#development-environment)
- [Production Web Stack](#production-web-stack)
- [Complete Email Setup](#complete-email-setup)
- [Pulumi Go Example](#pulumi-go-example)
- [Terraform HCL Example](#terraform-hcl-example)

---

## A Record (IPv4)

Map a subdomain to an IPv4 address:

```yaml
apiVersion: cloudflare.openmcf.org/v1
kind: CloudflareDnsRecord
metadata:
  name: www-a-record
spec:
  zoneId: "abc123def456"
  name: "www"
  type: A
  value: "192.0.2.1"
  proxied: true
  ttl: 1
  comment: "Primary web server"
```

**Deploy:**
```bash
planton apply -f www-a-record.yaml
```

**Use Case:** Web servers, load balancers, any HTTP/HTTPS services.

### A Record with Zone Reference

Reference an existing CloudflareDnsZone instead of hardcoding the zone ID:

```yaml
apiVersion: cloudflare.openmcf.org/v1
kind: CloudflareDnsRecord
metadata:
  name: www-a-record
spec:
  zoneId:
    valueFrom:
      name: example-com-zone
  name: "www"
  type: A
  value: "192.0.2.1"
  proxied: true
  ttl: 1
  comment: "Primary web server"
```

This approach is recommended when managing both the zone and records together, as it creates an implicit dependency and automatically resolves the zone ID.

---

## AAAA Record (IPv6)

Map a subdomain to an IPv6 address:

```yaml
apiVersion: cloudflare.openmcf.org/v1
kind: CloudflareDnsRecord
metadata:
  name: www-aaaa-record
spec:
  zoneId: "abc123def456"
  name: "www"
  type: AAAA
  value: "2001:db8::1"
  proxied: true
  comment: "IPv6 web server"
```

**Use Case:** Dual-stack deployments supporting both IPv4 and IPv6.

---

## CNAME Record

Create an alias pointing to another hostname:

```yaml
apiVersion: cloudflare.openmcf.org/v1
kind: CloudflareDnsRecord
metadata:
  name: app-cname
spec:
  zoneId: "abc123def456"
  name: "app"
  type: CNAME
  value: "www.example.com"
  proxied: true
```

**Use Case:** Aliases, CDN endpoints, SaaS integrations.

**Note:** CNAME at the root domain (apex) is possible with Cloudflare's CNAME flattening feature.

---

## MX Records (Email)

### Primary Mail Server

```yaml
apiVersion: cloudflare.openmcf.org/v1
kind: CloudflareDnsRecord
metadata:
  name: mx-primary
spec:
  zoneId: "abc123def456"
  name: "@"
  type: MX
  value: "mail.example.com"
  priority: 10
  comment: "Primary mail server"
```

### Backup Mail Server

```yaml
apiVersion: cloudflare.openmcf.org/v1
kind: CloudflareDnsRecord
metadata:
  name: mx-backup
spec:
  zoneId: "abc123def456"
  name: "@"
  type: MX
  value: "backup-mail.example.com"
  priority: 20
  comment: "Backup mail server"
```

### Google Workspace MX Records

```yaml
apiVersion: cloudflare.openmcf.org/v1
kind: CloudflareDnsRecord
metadata:
  name: mx-google-1
spec:
  zoneId: "abc123def456"
  name: "@"
  type: MX
  value: "aspmx.l.google.com"
  priority: 1

---
apiVersion: cloudflare.openmcf.org/v1
kind: CloudflareDnsRecord
metadata:
  name: mx-google-2
spec:
  zoneId: "abc123def456"
  name: "@"
  type: MX
  value: "alt1.aspmx.l.google.com"
  priority: 5

---
apiVersion: cloudflare.openmcf.org/v1
kind: CloudflareDnsRecord
metadata:
  name: mx-google-3
spec:
  zoneId: "abc123def456"
  name: "@"
  type: MX
  value: "alt2.aspmx.l.google.com"
  priority: 5
```

**Important:** MX records cannot be proxied (always grey cloud).

---

## TXT Records

### SPF Record

```yaml
apiVersion: cloudflare.openmcf.org/v1
kind: CloudflareDnsRecord
metadata:
  name: spf-record
spec:
  zoneId: "abc123def456"
  name: "@"
  type: TXT
  value: "v=spf1 include:_spf.google.com ~all"
  comment: "SPF for email authentication"
```

### DKIM Record

```yaml
apiVersion: cloudflare.openmcf.org/v1
kind: CloudflareDnsRecord
metadata:
  name: dkim-google
spec:
  zoneId: "abc123def456"
  name: "google._domainkey"
  type: TXT
  value: "v=DKIM1; k=rsa; p=MIGfMA0GCSqGSIb3DQEB..."
  comment: "Google Workspace DKIM"
```

### DMARC Record

```yaml
apiVersion: cloudflare.openmcf.org/v1
kind: CloudflareDnsRecord
metadata:
  name: dmarc-record
spec:
  zoneId: "abc123def456"
  name: "_dmarc"
  type: TXT
  value: "v=DMARC1; p=quarantine; rua=mailto:dmarc@example.com"
  comment: "DMARC policy"
```

### Domain Verification

```yaml
apiVersion: cloudflare.openmcf.org/v1
kind: CloudflareDnsRecord
metadata:
  name: google-verification
spec:
  zoneId: "abc123def456"
  name: "@"
  type: TXT
  value: "google-site-verification=abc123xyz..."
  comment: "Google Search Console verification"
```

---

## CAA Record

Control which Certificate Authorities can issue certificates for your domain:

```yaml
apiVersion: cloudflare.openmcf.org/v1
kind: CloudflareDnsRecord
metadata:
  name: caa-letsencrypt
spec:
  zoneId: "abc123def456"
  name: "@"
  type: CAA
  value: "0 issue \"letsencrypt.org\""
  comment: "Allow only Let's Encrypt"
```

### Allow Multiple CAs

```yaml
apiVersion: cloudflare.openmcf.org/v1
kind: CloudflareDnsRecord
metadata:
  name: caa-digicert
spec:
  zoneId: "abc123def456"
  name: "@"
  type: CAA
  value: "0 issue \"digicert.com\""
  comment: "Also allow DigiCert"
```

---

## Root Domain (Apex)

### A Record at Root

```yaml
apiVersion: cloudflare.openmcf.org/v1
kind: CloudflareDnsRecord
metadata:
  name: apex-a-record
spec:
  zoneId: "abc123def456"
  name: "@"
  type: A
  value: "192.0.2.1"
  proxied: true
  comment: "Root domain"
```

### CNAME at Root (Cloudflare Flattening)

```yaml
apiVersion: cloudflare.openmcf.org/v1
kind: CloudflareDnsRecord
metadata:
  name: apex-cname
spec:
  zoneId: "abc123def456"
  name: "@"
  type: CNAME
  value: "cdn.example.net"
  proxied: true
  comment: "Root domain via CNAME flattening"
```

**Note:** Cloudflare automatically flattens CNAME at root to return A/AAAA records.

---

## Wildcard Subdomain

Route all unmatched subdomains:

```yaml
apiVersion: cloudflare.openmcf.org/v1
kind: CloudflareDnsRecord
metadata:
  name: wildcard-record
spec:
  zoneId: "abc123def456"
  name: "*"
  type: A
  value: "192.0.2.1"
  proxied: true
  comment: "Wildcard catch-all"
```

**Use Case:** Multi-tenant applications, catch-all routing.

---

## Development Environment

Non-proxied records for development:

```yaml
apiVersion: cloudflare.openmcf.org/v1
kind: CloudflareDnsRecord
metadata:
  name: dev-api
  labels:
    environment: development
spec:
  zoneId: "abc123def456"
  name: "dev-api"
  type: A
  value: "10.0.1.100"
  proxied: false  # No proxy for internal dev
  ttl: 60  # Low TTL for frequent changes
  comment: "Development API server"
```

---

## Production Web Stack

### Directory Structure

```
my-app/
├── dns/
│   ├── www.yaml
│   ├── api.yaml
│   ├── cdn.yaml
│   └── mail.yaml
└── deploy.sh
```

### WWW Record

```yaml
apiVersion: cloudflare.openmcf.org/v1
kind: CloudflareDnsRecord
metadata:
  name: prod-www
  labels:
    environment: production
spec:
  zoneId: "abc123def456"
  name: "www"
  type: A
  value: "192.0.2.10"
  proxied: true
  comment: "Production web"
```

### API Record

```yaml
apiVersion: cloudflare.openmcf.org/v1
kind: CloudflareDnsRecord
metadata:
  name: prod-api
  labels:
    environment: production
spec:
  zoneId: "abc123def456"
  name: "api"
  type: CNAME
  value: "api-lb.example.com"
  proxied: true
  comment: "Production API"
```

### CDN/Static Assets

```yaml
apiVersion: cloudflare.openmcf.org/v1
kind: CloudflareDnsRecord
metadata:
  name: prod-cdn
  labels:
    environment: production
spec:
  zoneId: "abc123def456"
  name: "cdn"
  type: CNAME
  value: "d123.cloudfront.net"
  proxied: false  # AWS CloudFront handles its own CDN
  comment: "Static assets CDN"
```

---

## Complete Email Setup

All records needed for professional email:

### MX Records

```yaml
apiVersion: cloudflare.openmcf.org/v1
kind: CloudflareDnsRecord
metadata:
  name: email-mx-1
spec:
  zoneId: "abc123def456"
  name: "@"
  type: MX
  value: "aspmx.l.google.com"
  priority: 1
---
apiVersion: cloudflare.openmcf.org/v1
kind: CloudflareDnsRecord
metadata:
  name: email-mx-2
spec:
  zoneId: "abc123def456"
  name: "@"
  type: MX
  value: "alt1.aspmx.l.google.com"
  priority: 5
```

### SPF Record

```yaml
apiVersion: cloudflare.openmcf.org/v1
kind: CloudflareDnsRecord
metadata:
  name: email-spf
spec:
  zoneId: "abc123def456"
  name: "@"
  type: TXT
  value: "v=spf1 include:_spf.google.com ~all"
```

### DKIM Record

```yaml
apiVersion: cloudflare.openmcf.org/v1
kind: CloudflareDnsRecord
metadata:
  name: email-dkim
spec:
  zoneId: "abc123def456"
  name: "google._domainkey"
  type: TXT
  value: "v=DKIM1; k=rsa; p=YOUR_DKIM_PUBLIC_KEY"
```

### DMARC Record

```yaml
apiVersion: cloudflare.openmcf.org/v1
kind: CloudflareDnsRecord
metadata:
  name: email-dmarc
spec:
  zoneId: "abc123def456"
  name: "_dmarc"
  type: TXT
  value: "v=DMARC1; p=quarantine; pct=100; rua=mailto:dmarc-reports@example.com"
```

---

## Pulumi Go Example

Direct Pulumi Go code for creating a DNS record (without OpenMCF CLI):

```go
package main

import (
	"github.com/pulumi/pulumi-cloudflare/sdk/v6/go/cloudflare"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Create A record
		record, err := cloudflare.NewRecord(ctx, "www-record", &cloudflare.RecordArgs{
			ZoneId:  pulumi.String("abc123def456"),
			Name:    pulumi.String("www"),
			Type:    pulumi.String("A"),
			Value:   pulumi.String("192.0.2.1"),
			Proxied: pulumi.Bool(true),
			Ttl:     pulumi.Int(1),
			Comment: pulumi.String("Web server"),
		})
		if err != nil {
			return err
		}

		// Export outputs
		ctx.Export("recordId", record.ID())
		ctx.Export("hostname", record.Hostname)

		return nil
	})
}
```

**Deploy:**
```bash
pulumi up
```

---

## Terraform HCL Example

Direct Terraform HCL code for creating a DNS record (without OpenMCF CLI):

```hcl
terraform {
  required_providers {
    cloudflare = {
      source  = "cloudflare/cloudflare"
      version = "~> 4.0"
    }
  }
}

provider "cloudflare" {
  api_token = var.cloudflare_api_token
}

variable "cloudflare_api_token" {
  description = "Cloudflare API token"
  type        = string
  sensitive   = true
}

variable "zone_id" {
  description = "Cloudflare zone ID"
  type        = string
}

# A Record
resource "cloudflare_record" "www" {
  zone_id = var.zone_id
  name    = "www"
  type    = "A"
  value   = "192.0.2.1"
  proxied = true
  ttl     = 1
  comment = "Web server"
}

# CNAME Record
resource "cloudflare_record" "api" {
  zone_id = var.zone_id
  name    = "api"
  type    = "CNAME"
  value   = "api-lb.example.com"
  proxied = true
}

# MX Record
resource "cloudflare_record" "mx" {
  zone_id  = var.zone_id
  name     = "@"
  type     = "MX"
  value    = "mail.example.com"
  priority = 10
}

output "www_record_id" {
  description = "The ID of the www record"
  value       = cloudflare_record.www.id
}

output "www_hostname" {
  description = "The hostname of the www record"
  value       = cloudflare_record.www.hostname
}
```

**Deploy:**
```bash
terraform init
terraform apply -var="cloudflare_api_token=$CLOUDFLARE_API_TOKEN" -var="zone_id=abc123..."
```

---

## Support

For questions or issues:
- **OpenMCF**: [openmcf.org](https://openmcf.org)
- **Cloudflare DNS Docs**: [developers.cloudflare.com/dns](https://developers.cloudflare.com/dns)
