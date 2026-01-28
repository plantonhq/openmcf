# Cloudflare DNS Zone Examples

Concrete, copy-and-paste examples for common Cloudflare DNS zone deployment scenarios.

## Table of Contents

- [Minimal Configuration](#minimal-configuration)
- [Free Plan Zone](#free-plan-zone)
- [Pro Plan Zone](#pro-plan-zone)
- [Business Plan Zone](#business-plan-zone)
- [Paused Zone (DNS-Only)](#paused-zone-dns-only)
- [Default Proxied Zone](#default-proxied-zone)
- [Zone with DNS Records](#zone-with-dns-records)
- [Complete Zone with Multiple Records](#complete-zone-with-multiple-records)
- [Multi-Environment Setup](#multi-environment-setup)
- [Pulumi Go Example](#pulumi-go-example)
- [Terraform HCL Example](#terraform-hcl-example)

---

## Minimal Configuration

The simplest possible DNS zone with only required fields:

\`\`\`yaml
apiVersion: cloudflare.openmcf.org/v1
kind: CloudflareDnsZone
metadata:
  name: minimal-zone
spec:
  zoneName: "example.com"
  accountId: "abc123def456..."
\`\`\`

**Deploy:**
\`\`\`bash
planton apply -f minimal-zone.yaml
\`\`\`

**Use Case:** Quick experimentation or proof-of-concept.

---

## Free Plan Zone

Default free plan with all standard features:

\`\`\`yaml
apiVersion: cloudflare.openmcf.org/v1
kind: CloudflareDnsZone
metadata:
  name: free-zone
  labels:
    environment: production
    plan: free
spec:
  zoneName: "myblog.com"
  accountId: "abc123..."
  plan: FREE
  paused: false
  defaultProxied: true  # Orange-cloud new records by default
\`\`\`

**Features:**
- Unlimited DNS queries
- Basic CDN
- Universal SSL
- DDoS protection
- Free forever

**Use Case:** Personal blogs, small projects, portfolio sites.

---

## Pro Plan Zone

Pro plan with advanced features:

\`\`\`yaml
apiVersion: cloudflare.openmcf.org/v1
kind: CloudflareDnsZone
metadata:
  name: pro-zone
  labels:
    environment: production
    plan: pro
spec:
  zoneName: "mybusiness.com"
  accountId: "abc123..."
  plan: PRO
  paused: false
  defaultProxied: true
\`\`\`

**Features:**
- All Free features
- WAF (Web Application Firewall)
- Polish image optimization
- Mobile optimization
- Advanced SSL

**Use Case:** Small to medium businesses, professional blogs, e-commerce sites.

---

## Business Plan Zone

Business plan for enterprise needs:

\`\`\`yaml
apiVersion: cloudflare.openmcf.org/v1
kind: CloudflareDnsZone
metadata:
  name: business-zone
  labels:
    environment: production
    plan: business
spec:
  zoneName: "enterprise.example.com"
  accountId: "xyz789..."
  plan: BUSINESS
  paused: false
  defaultProxied: false  # Explicit control per record
\`\`\`

**Features:**
- All Pro features
- Custom SSL certificates
- 99.95% SLA
- Advanced DDoS protection
- 24/7 priority support
- 100 page rules

**Use Case:** E-commerce, SaaS platforms, high-traffic sites.

---

## Paused Zone (DNS-Only)

Create a zone in paused state (DNS without CDN/WAF):

\`\`\`yaml
apiVersion: cloudflare.openmcf.org/v1
kind: CloudflareDnsZone
metadata:
  name: dns-only-zone
spec:
  zoneName: "dnsonly.example.com"
  accountId: "abc123..."
  paused: true  # No proxy/CDN/WAF features
\`\`\`

**Use Case:**
- Using Cloudflare purely for DNS (no proxy features)
- Development/staging environments where you don't want caching
- Gradual migration (DNS first, enable proxy later)

**Note:** Paused zones still get fast global DNS resolution, just no CDN/proxy features.

---

## Default Proxied Zone

Zone where all new DNS records default to orange-cloud (proxied):

\`\`\`yaml
apiVersion: cloudflare.openmcf.org/v1
kind: CloudflareDnsZone
metadata:
  name: proxied-zone
spec:
  zoneName: "secure.example.com"
  accountId: "abc123..."
  plan: PRO
  defaultProxied: true  # New records orange-cloud by default
\`\`\`

**Use Case:**
- Security-first approach (hides origin IPs)
- Ensures all web traffic goes through Cloudflare's edge
- Prevents accidental origin IP exposure

**Note:** You can still manually set individual records to grey-cloud when creating them.

---

## Zone with DNS Records

Zone with inline DNS records managed together:

\`\`\`yaml
apiVersion: cloudflare.openmcf.org/v1
kind: CloudflareDnsZone
metadata:
  name: zone-with-records
spec:
  zoneName: "myapp.com"
  accountId: "abc123..."
  plan: PRO
  records:
    # A record for root domain
    - name: "@"
      type: A
      value: "192.0.2.1"
      proxied: true
      ttl: 1  # Automatic TTL

    # A record for www subdomain
    - name: "www"
      type: A
      value: "192.0.2.1"
      proxied: true

    # CNAME for API subdomain
    - name: "api"
      type: CNAME
      value: "api-gateway.example.com"
      proxied: true

    # MX record for email
    - name: "@"
      type: MX
      value: "mail.example.com"
      priority: 10
      ttl: 3600
\`\`\`

**Use Case:** Manage zone and DNS records together as a single deployment unit.

---

## Complete Zone with Multiple Records

A production-ready zone with comprehensive DNS setup:

\`\`\`yaml
apiVersion: cloudflare.openmcf.org/v1
kind: CloudflareDnsZone
metadata:
  name: production-zone
  labels:
    environment: production
spec:
  zoneName: "production.example.com"
  accountId: "abc123..."
  plan: BUSINESS
  defaultProxied: true
  records:
    # Web traffic records (proxied)
    - name: "@"
      type: A
      value: "203.0.113.50"
      proxied: true
      comment: "Production web server"

    - name: "www"
      type: CNAME
      value: "production.example.com"
      proxied: true

    - name: "api"
      type: A
      value: "203.0.113.51"
      proxied: true
      comment: "API server"

    # Email records (not proxied)
    - name: "@"
      type: MX
      value: "mail1.example.com"
      priority: 10
      ttl: 3600
      comment: "Primary mail server"

    - name: "@"
      type: MX
      value: "mail2.example.com"
      priority: 20
      ttl: 3600
      comment: "Secondary mail server"

    # TXT records for verification and SPF
    - name: "@"
      type: TXT
      value: "v=spf1 include:_spf.google.com ~all"
      ttl: 3600
      comment: "SPF record"

    - name: "_dmarc"
      type: TXT
      value: "v=DMARC1; p=quarantine; rua=mailto:dmarc@example.com"
      ttl: 3600
      comment: "DMARC policy"

    # CAA record for certificate authority
    - name: "@"
      type: CAA
      value: "0 issue \"letsencrypt.org\""
      ttl: 3600
      comment: "Allow Let's Encrypt certificates"
\`\`\`

**Features:**
- Web traffic proxied through Cloudflare (CDN, WAF, DDoS protection)
- Email records with proper MX priorities
- SPF and DMARC for email security
- CAA for certificate issuance control

**Use Case:** Full production setup with web, email, and security configurations.

---

## Multi-Environment Setup

Complete multi-environment setup with dev, staging, and production zones.

### Directory Structure

\`\`\`
my-project/
├── zones/
│   ├── dev-zone.yaml
│   ├── staging-zone.yaml
│   └── prod-zone.yaml
└── deploy.sh
\`\`\`

### Development Zone

**File:** \`zones/dev-zone.yaml\`

\`\`\`yaml
apiVersion: cloudflare.openmcf.org/v1
kind: CloudflareDnsZone
metadata:
  name: dev-zone
  labels:
    environment: development
spec:
  zoneName: "dev.example.com"
  accountId: "abc123..."
  plan: FREE
  paused: false
  defaultProxied: true
\`\`\`

### Staging Zone

**File:** \`zones/staging-zone.yaml\`

\`\`\`yaml
apiVersion: cloudflare.openmcf.org/v1
kind: CloudflareDnsZone
metadata:
  name: staging-zone
  labels:
    environment: staging
spec:
  zoneName: "staging.example.com"
  accountId: "abc123..."
  plan: PRO
  paused: false
  defaultProxied: true
\`\`\`

### Production Zone

**File:** \`zones/prod-zone.yaml\`

\`\`\`yaml
apiVersion: cloudflare.openmcf.org/v1
kind: CloudflareDnsZone
metadata:
  name: prod-zone
  labels:
    environment: production
spec:
  zoneName: "example.com"
  accountId: "abc123..."
  plan: BUSINESS
  paused: false
  defaultProxied: true
\`\`\`

### Deployment Script

**File:** \`deploy.sh\`

\`\`\`bash
#!/bin/bash
set -e

ENVIRONMENT=\$1

if [ "\$ENVIRONMENT" = "dev" ]; then
  planton apply -f zones/dev-zone.yaml
elif [ "\$ENVIRONMENT" = "staging" ]; then
  planton apply -f zones/staging-zone.yaml
elif [ "\$ENVIRONMENT" = "prod" ]; then
  planton apply -f zones/prod-zone.yaml
else
  echo "Usage: ./deploy.sh [dev|staging|prod]"
  exit 1
fi
\`\`\`

**Usage:**
\`\`\`bash
chmod +x deploy.sh
./deploy.sh dev      # Deploy development zone
./deploy.sh staging  # Deploy staging zone
./deploy.sh prod     # Deploy production zone
\`\`\`

---

## Pulumi Go Example

Direct Pulumi Go code for provisioning a DNS zone (without OpenMCF CLI):

\`\`\`go
package main

import (
	"github.com/pulumi/pulumi-cloudflare/sdk/v6/go/cloudflare"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Create Cloudflare DNS zone
		zone, err := cloudflare.NewZone(ctx, "my-zone", &cloudflare.ZoneArgs{
			Account: cloudflare.ZoneAccountArgs{
				Id: pulumi.String("abc123..."),
			},
			Name:   pulumi.String("example.com"),
			Paused: pulumi.Bool(false),
			Plan:   pulumi.String("free"),
		})
		if err != nil {
			return err
		}

		// Export outputs
		ctx.Export("zoneId", zone.ID())
		ctx.Export("nameservers", zone.NameServers)

		return nil
	})
}
\`\`\`

**Deploy:**
\`\`\`bash
pulumi up
\`\`\`

---

## Terraform HCL Example

Direct Terraform HCL code for provisioning a DNS zone (without OpenMCF CLI):

\`\`\`hcl
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

variable "account_id" {
  description = "Cloudflare account ID"
  type        = string
}

resource "cloudflare_zone" "main" {
  account_id = var.account_id
  zone       = "example.com"
  plan       = "free"
  paused     = false
  type       = "full"
  jump_start = false
}

output "zone_id" {
  description = "The ID of the created zone"
  value       = cloudflare_zone.main.id
}

output "nameservers" {
  description = "The nameservers for the zone"
  value       = cloudflare_zone.main.name_servers
}
\`\`\`

**Deploy:**
\`\`\`bash
terraform init
terraform apply -var="cloudflare_api_token=\$CLOUDFLARE_API_TOKEN" -var="account_id=abc123..."
\`\`\`

---

## Migration Example

Migrating an existing domain from another DNS provider to Cloudflare:

\`\`\`yaml
apiVersion: cloudflare.openmcf.org/v1
kind: CloudflareDnsZone
metadata:
  name: migrated-zone
spec:
  zoneName: "oldprovider.com"
  accountId: "abc123..."
  plan: PRO
  paused: false  # Active immediately
  defaultProxied: false  # Start with grey-cloud for safety
\`\`\`

**Migration Steps:**

1. **Create zone** (but don't update nameservers yet):
   \`\`\`bash
   planton apply -f migrated-zone.yaml
   \`\`\`

2. **Get Cloudflare nameservers**:
   \`\`\`bash
   planton output nameservers
   \`\`\`

3. **Import DNS records** (via dashboard DNS scan or manual creation)

4. **Test before going live**:
   \`\`\`bash
   dig @gina.ns.cloudflare.com oldprovider.com A
   \`\`\`

5. **Update nameservers** at registrar once records are verified

---

## Support

For questions or issues:
- **OpenMCF**: [openmcf.org](https://openmcf.org)
- **Cloudflare DNS Docs**: [developers.cloudflare.com/dns](https://developers.cloudflare.com/dns)
