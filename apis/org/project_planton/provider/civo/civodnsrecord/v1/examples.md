# Civo DNS Record Examples

Concrete, copy-and-paste examples for common Civo DNS record deployment scenarios.

## Table of Contents

- [A Record (IPv4)](#a-record-ipv4)
- [AAAA Record (IPv6)](#aaaa-record-ipv6)
- [CNAME Record](#cname-record)
- [MX Records (Email)](#mx-records-email)
- [TXT Records](#txt-records)
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
apiVersion: civo.project-planton.org/v1
kind: CivoDnsRecord
metadata:
  name: www-a-record
spec:
  zoneId: "zone-abc123"
  name: "www"
  type: A
  value: "192.0.2.1"
  ttl: 3600
```

**Deploy:**
```bash
planton apply -f www-a-record.yaml
```

**Use Case:** Web servers, load balancers, any HTTP/HTTPS services.

---

## AAAA Record (IPv6)

Map a subdomain to an IPv6 address:

```yaml
apiVersion: civo.project-planton.org/v1
kind: CivoDnsRecord
metadata:
  name: www-aaaa-record
spec:
  zoneId: "zone-abc123"
  name: "www"
  type: AAAA
  value: "2001:db8::1"
  ttl: 3600
```

**Use Case:** Dual-stack deployments supporting both IPv4 and IPv6.

---

## CNAME Record

Create an alias pointing to another hostname:

```yaml
apiVersion: civo.project-planton.org/v1
kind: CivoDnsRecord
metadata:
  name: app-cname
spec:
  zoneId: "zone-abc123"
  name: "app"
  type: CNAME
  value: "www.example.com"
```

**Use Case:** Aliases, CDN endpoints, SaaS integrations.

---

## MX Records (Email)

### Primary Mail Server

```yaml
apiVersion: civo.project-planton.org/v1
kind: CivoDnsRecord
metadata:
  name: mx-primary
spec:
  zoneId: "zone-abc123"
  name: "@"
  type: MX
  value: "mail.example.com"
  priority: 10
```

### Backup Mail Server

```yaml
apiVersion: civo.project-planton.org/v1
kind: CivoDnsRecord
metadata:
  name: mx-backup
spec:
  zoneId: "zone-abc123"
  name: "@"
  type: MX
  value: "backup-mail.example.com"
  priority: 20
```

### Google Workspace MX Records

```yaml
apiVersion: civo.project-planton.org/v1
kind: CivoDnsRecord
metadata:
  name: mx-google-1
spec:
  zoneId: "zone-abc123"
  name: "@"
  type: MX
  value: "aspmx.l.google.com"
  priority: 1

---
apiVersion: civo.project-planton.org/v1
kind: CivoDnsRecord
metadata:
  name: mx-google-2
spec:
  zoneId: "zone-abc123"
  name: "@"
  type: MX
  value: "alt1.aspmx.l.google.com"
  priority: 5

---
apiVersion: civo.project-planton.org/v1
kind: CivoDnsRecord
metadata:
  name: mx-google-3
spec:
  zoneId: "zone-abc123"
  name: "@"
  type: MX
  value: "alt2.aspmx.l.google.com"
  priority: 5
```

---

## TXT Records

### SPF Record

```yaml
apiVersion: civo.project-planton.org/v1
kind: CivoDnsRecord
metadata:
  name: spf-record
spec:
  zoneId: "zone-abc123"
  name: "@"
  type: TXT
  value: "v=spf1 include:_spf.google.com ~all"
```

### DKIM Record

```yaml
apiVersion: civo.project-planton.org/v1
kind: CivoDnsRecord
metadata:
  name: dkim-google
spec:
  zoneId: "zone-abc123"
  name: "google._domainkey"
  type: TXT
  value: "v=DKIM1; k=rsa; p=MIGfMA0GCSqGSIb3DQEB..."
```

### DMARC Record

```yaml
apiVersion: civo.project-planton.org/v1
kind: CivoDnsRecord
metadata:
  name: dmarc-record
spec:
  zoneId: "zone-abc123"
  name: "_dmarc"
  type: TXT
  value: "v=DMARC1; p=quarantine; rua=mailto:dmarc@example.com"
```

### Domain Verification

```yaml
apiVersion: civo.project-planton.org/v1
kind: CivoDnsRecord
metadata:
  name: google-verification
spec:
  zoneId: "zone-abc123"
  name: "@"
  type: TXT
  value: "google-site-verification=abc123xyz..."
```

---

## Root Domain (Apex)

### A Record at Root

```yaml
apiVersion: civo.project-planton.org/v1
kind: CivoDnsRecord
metadata:
  name: apex-a-record
spec:
  zoneId: "zone-abc123"
  name: "@"
  type: A
  value: "192.0.2.1"
  ttl: 3600
```

---

## Wildcard Subdomain

Route all unmatched subdomains:

```yaml
apiVersion: civo.project-planton.org/v1
kind: CivoDnsRecord
metadata:
  name: wildcard-record
spec:
  zoneId: "zone-abc123"
  name: "*"
  type: A
  value: "192.0.2.1"
  ttl: 3600
```

**Use Case:** Multi-tenant applications, catch-all routing.

---

## Development Environment

Records for development environment:

```yaml
apiVersion: civo.project-planton.org/v1
kind: CivoDnsRecord
metadata:
  name: dev-api
  labels:
    environment: development
spec:
  zoneId: "zone-abc123"
  name: "dev-api"
  type: A
  value: "10.0.1.100"
  ttl: 60  # Low TTL for frequent changes
```

---

## Production Web Stack

### Directory Structure

```
my-app/
├── dns/
│   ├── www.yaml
│   ├── api.yaml
│   └── mail.yaml
└── deploy.sh
```

### WWW Record

```yaml
apiVersion: civo.project-planton.org/v1
kind: CivoDnsRecord
metadata:
  name: prod-www
  labels:
    environment: production
spec:
  zoneId: "zone-abc123"
  name: "www"
  type: A
  value: "192.0.2.10"
  ttl: 3600
```

### API Record

```yaml
apiVersion: civo.project-planton.org/v1
kind: CivoDnsRecord
metadata:
  name: prod-api
  labels:
    environment: production
spec:
  zoneId: "zone-abc123"
  name: "api"
  type: CNAME
  value: "api-lb.example.com"
```

---

## Complete Email Setup

All records needed for professional email:

### MX Records

```yaml
apiVersion: civo.project-planton.org/v1
kind: CivoDnsRecord
metadata:
  name: email-mx-1
spec:
  zoneId: "zone-abc123"
  name: "@"
  type: MX
  value: "aspmx.l.google.com"
  priority: 1
---
apiVersion: civo.project-planton.org/v1
kind: CivoDnsRecord
metadata:
  name: email-mx-2
spec:
  zoneId: "zone-abc123"
  name: "@"
  type: MX
  value: "alt1.aspmx.l.google.com"
  priority: 5
```

### SPF Record

```yaml
apiVersion: civo.project-planton.org/v1
kind: CivoDnsRecord
metadata:
  name: email-spf
spec:
  zoneId: "zone-abc123"
  name: "@"
  type: TXT
  value: "v=spf1 include:_spf.google.com ~all"
```

### DKIM Record

```yaml
apiVersion: civo.project-planton.org/v1
kind: CivoDnsRecord
metadata:
  name: email-dkim
spec:
  zoneId: "zone-abc123"
  name: "google._domainkey"
  type: TXT
  value: "v=DKIM1; k=rsa; p=YOUR_DKIM_PUBLIC_KEY"
```

### DMARC Record

```yaml
apiVersion: civo.project-planton.org/v1
kind: CivoDnsRecord
metadata:
  name: email-dmarc
spec:
  zoneId: "zone-abc123"
  name: "_dmarc"
  type: TXT
  value: "v=DMARC1; p=quarantine; pct=100; rua=mailto:dmarc-reports@example.com"
```

---

## Pulumi Go Example

Direct Pulumi Go code for creating a DNS record (without Project Planton CLI):

```go
package main

import (
	"github.com/pulumi/pulumi-civo/sdk/v2/go/civo"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Create A record
		record, err := civo.NewDnsRecord(ctx, "www-record", &civo.DnsRecordArgs{
			DomainId: pulumi.String("zone-abc123"),
			Name:     pulumi.String("www"),
			Type:     pulumi.String("A"),
			Value:    pulumi.String("192.0.2.1"),
			Ttl:      pulumi.Int(3600),
		})
		if err != nil {
			return err
		}

		// Export outputs
		ctx.Export("recordId", record.ID())
		ctx.Export("accountId", record.AccountId)

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

Direct Terraform HCL code for creating a DNS record (without Project Planton CLI):

```hcl
terraform {
  required_providers {
    civo = {
      source  = "civo/civo"
      version = "~> 1.0"
    }
  }
}

provider "civo" {
  token = var.civo_api_key
}

variable "civo_api_key" {
  description = "Civo API key"
  type        = string
  sensitive   = true
}

variable "zone_id" {
  description = "Civo DNS zone ID"
  type        = string
}

# A Record
resource "civo_dns_domain_record" "www" {
  domain_id = var.zone_id
  name      = "www"
  type      = "A"
  value     = "192.0.2.1"
  ttl       = 3600
}

# CNAME Record
resource "civo_dns_domain_record" "api" {
  domain_id = var.zone_id
  name      = "api"
  type      = "CNAME"
  value     = "api-lb.example.com"
  ttl       = 3600
}

# MX Record
resource "civo_dns_domain_record" "mx" {
  domain_id = var.zone_id
  name      = "@"
  type      = "MX"
  value     = "mail.example.com"
  priority  = 10
  ttl       = 3600
}

output "www_record_id" {
  description = "The ID of the www record"
  value       = civo_dns_domain_record.www.id
}
```

**Deploy:**
```bash
terraform init
terraform apply -var="civo_api_key=$CIVO_API_KEY" -var="zone_id=zone-abc123"
```

---

## Support

For questions or issues:
- **Project Planton**: [project-planton.org](https://project-planton.org)
- **Civo DNS Docs**: [civo.com/docs/dns](https://www.civo.com/docs/dns)
