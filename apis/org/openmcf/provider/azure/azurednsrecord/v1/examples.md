# Azure DNS Record Examples

This document provides copy-paste ready examples for deploying Azure DNS Records using OpenMCF.

## Basic A Record

Creates a simple A record pointing to an IPv4 address.

```yaml
apiVersion: azure.openmcf.org/v1
kind: AzureDnsRecord
metadata:
  name: www-a-record
spec:
  resource_group: production-dns-rg
  zone_name:
    value: example.com
  record_type: A
  name: www
  values:
    - "192.0.2.1"
  ttl_seconds: 300
```

## Zone Apex A Record

Creates an A record at the zone apex (root domain).

```yaml
apiVersion: azure.openmcf.org/v1
kind: AzureDnsRecord
metadata:
  name: apex-a-record
spec:
  resource_group: production-dns-rg
  zone_name:
    value: example.com
  record_type: A
  name: "@"
  values:
    - "192.0.2.1"
    - "192.0.2.2"
  ttl_seconds: 300
```

## AAAA Record (IPv6)

Creates an AAAA record for IPv6 addresses.

```yaml
apiVersion: azure.openmcf.org/v1
kind: AzureDnsRecord
metadata:
  name: www-aaaa-record
spec:
  resource_group: production-dns-rg
  zone_name:
    value: example.com
  record_type: AAAA
  name: www
  values:
    - "2001:db8::1"
  ttl_seconds: 300
```

## CNAME Record

Creates a CNAME (alias) record.

```yaml
apiVersion: azure.openmcf.org/v1
kind: AzureDnsRecord
metadata:
  name: blog-cname
spec:
  resource_group: production-dns-rg
  zone_name:
    value: example.com
  record_type: CNAME
  name: blog
  values:
    - "example.ghost.io"
  ttl_seconds: 3600
```

## MX Record (Mail Exchange)

Creates MX records for email routing.

```yaml
apiVersion: azure.openmcf.org/v1
kind: AzureDnsRecord
metadata:
  name: mail-mx-record
spec:
  resource_group: production-dns-rg
  zone_name:
    value: example.com
  record_type: MX
  name: "@"
  values:
    - "aspmx.l.google.com"
    - "alt1.aspmx.l.google.com"
  mx_priority: 10
  ttl_seconds: 3600
```

## TXT Record (SPF)

Creates a TXT record for SPF email authentication.

```yaml
apiVersion: azure.openmcf.org/v1
kind: AzureDnsRecord
metadata:
  name: spf-txt-record
spec:
  resource_group: production-dns-rg
  zone_name:
    value: example.com
  record_type: TXT
  name: "@"
  values:
    - "v=spf1 include:_spf.google.com ~all"
  ttl_seconds: 3600
```

## TXT Record (Domain Verification)

Creates a TXT record for domain verification (e.g., Google, Microsoft).

```yaml
apiVersion: azure.openmcf.org/v1
kind: AzureDnsRecord
metadata:
  name: verification-txt-record
spec:
  resource_group: production-dns-rg
  zone_name:
    value: example.com
  record_type: TXT
  name: "@"
  values:
    - "google-site-verification=abc123xyz"
  ttl_seconds: 300
```

## CAA Record (Certificate Authority Authorization)

Creates a CAA record to control SSL certificate issuance.

```yaml
apiVersion: azure.openmcf.org/v1
kind: AzureDnsRecord
metadata:
  name: caa-record
spec:
  resource_group: production-dns-rg
  zone_name:
    value: example.com
  record_type: CAA
  name: "@"
  values:
    - "letsencrypt.org"
    - "digicert.com"
  ttl_seconds: 86400
```

## Wildcard A Record

Creates a wildcard A record for catch-all subdomain routing.

```yaml
apiVersion: azure.openmcf.org/v1
kind: AzureDnsRecord
metadata:
  name: wildcard-a-record
spec:
  resource_group: production-dns-rg
  zone_name:
    value: example.com
  record_type: A
  name: "*"
  values:
    - "192.0.2.100"
  ttl_seconds: 300
```

## Record with Zone Reference

Creates a record that references an AzureDnsZone resource for the zone name.

```yaml
apiVersion: azure.openmcf.org/v1
kind: AzureDnsRecord
metadata:
  name: api-record-with-ref
  env: production
spec:
  resource_group: production-dns-rg
  zone_name:
    value_from:
      name: production-dns-zone
  record_type: A
  name: api
  values:
    - "192.0.2.50"
  ttl_seconds: 60
```

## NS Record (Subdomain Delegation)

Creates NS records to delegate a subdomain to different nameservers.

```yaml
apiVersion: azure.openmcf.org/v1
kind: AzureDnsRecord
metadata:
  name: subdomain-ns-record
spec:
  resource_group: production-dns-rg
  zone_name:
    value: example.com
  record_type: NS
  name: staging
  values:
    - "ns1.staging-provider.com"
    - "ns2.staging-provider.com"
  ttl_seconds: 86400
```

## Deployment Commands

Deploy any of these examples using:

```bash
# Using Pulumi backend
openmcf pulumi up --manifest <example-file>.yaml

# Using Terraform/OpenTofu backend
openmcf tofu apply --manifest <example-file>.yaml
```
