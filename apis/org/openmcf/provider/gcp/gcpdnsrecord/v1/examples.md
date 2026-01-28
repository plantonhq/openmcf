# GcpDnsRecord Examples

This document provides copy-paste ready examples for common DNS record configurations.

## Basic A Record

Create an A record pointing to a single IPv4 address.

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpDnsRecord
metadata:
  name: www-example-a-record
spec:
  projectId: my-gcp-project
  managedZone: example-zone
  recordType: A
  name: www.example.com.
  values:
    - 192.0.2.1
```

## Round-Robin A Record

Distribute traffic across multiple IP addresses.

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpDnsRecord
metadata:
  name: api-loadbalanced
spec:
  projectId: my-gcp-project
  managedZone: example-zone
  recordType: A
  name: api.example.com.
  values:
    - 192.0.2.1
    - 192.0.2.2
    - 192.0.2.3
  ttlSeconds: 60
```

## CNAME Record

Create an alias pointing to another domain.

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpDnsRecord
metadata:
  name: blog-cname
spec:
  projectId: my-gcp-project
  managedZone: example-zone
  recordType: CNAME
  name: blog.example.com.
  values:
    - example.github.io.
```

## Wildcard Record

Route all subdomains to a single IP.

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpDnsRecord
metadata:
  name: wildcard-record
spec:
  projectId: my-gcp-project
  managedZone: example-zone
  recordType: A
  name: "*.example.com."
  values:
    - 192.0.2.1
  ttlSeconds: 300
```

## MX Record (Mail Exchange)

Configure email routing for your domain.

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpDnsRecord
metadata:
  name: google-workspace-mx
spec:
  projectId: my-gcp-project
  managedZone: example-zone
  recordType: MX
  name: example.com.
  values:
    - "10 aspmx.l.google.com."
    - "20 alt1.aspmx.l.google.com."
    - "20 alt2.aspmx.l.google.com."
    - "30 alt3.aspmx.l.google.com."
    - "30 alt4.aspmx.l.google.com."
  ttlSeconds: 3600
```

## TXT Record (SPF)

Add SPF record for email authentication.

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpDnsRecord
metadata:
  name: spf-record
spec:
  projectId: my-gcp-project
  managedZone: example-zone
  recordType: TXT
  name: example.com.
  values:
    - "v=spf1 include:_spf.google.com ~all"
  ttlSeconds: 3600
```

## TXT Record (Domain Verification)

Add verification record for third-party services.

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpDnsRecord
metadata:
  name: google-site-verification
spec:
  projectId: my-gcp-project
  managedZone: example-zone
  recordType: TXT
  name: example.com.
  values:
    - "google-site-verification=abc123xyz"
```

## AAAA Record (IPv6)

Create a DNS record for IPv6 addresses.

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpDnsRecord
metadata:
  name: ipv6-record
spec:
  projectId: my-gcp-project
  managedZone: example-zone
  recordType: AAAA
  name: www.example.com.
  values:
    - "2001:db8::1"
```

## CAA Record

Specify which Certificate Authorities can issue certificates for your domain.

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpDnsRecord
metadata:
  name: caa-letsencrypt
spec:
  projectId: my-gcp-project
  managedZone: example-zone
  recordType: CAA
  name: example.com.
  values:
    - '0 issue "letsencrypt.org"'
    - '0 issuewild "letsencrypt.org"'
  ttlSeconds: 86400
```

## SRV Record

Service location record for specific services.

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpDnsRecord
metadata:
  name: xmpp-server-srv
spec:
  projectId: my-gcp-project
  managedZone: example-zone
  recordType: SRV
  name: _xmpp-server._tcp.example.com.
  values:
    - "10 5 5269 xmpp.example.com."
```

## Production Configuration with Low TTL

For records that may need quick updates during incidents.

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpDnsRecord
metadata:
  name: prod-api-failover-ready
spec:
  projectId: production-project
  managedZone: production-zone
  recordType: A
  name: api.example.com.
  values:
    - 203.0.113.10
  ttlSeconds: 60
```

## Using Project Reference

Reference a GcpProject resource instead of hardcoding the project ID.

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpDnsRecord
metadata:
  name: referenced-project-record
spec:
  projectId:
    ref:
      kind: GcpProject
      name: my-gcp-project-resource
      fieldPath: status.outputs.project_id
  managedZone: example-zone
  recordType: A
  name: app.example.com.
  values:
    - 192.0.2.100
```

## Deployment Commands

### Deploy with Pulumi

```bash
openmcf pulumi up --manifest dns-record.yaml
```

### Deploy with Terraform

```bash
openmcf terraform apply --manifest dns-record.yaml
```

### Validate Manifest

```bash
openmcf validate --manifest dns-record.yaml
```

### Destroy Record

```bash
openmcf pulumi destroy --manifest dns-record.yaml
```
