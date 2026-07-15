---
title: "Route53 DNS Record"
description: "Route53 DNS Record deployment documentation"
icon: "package"
order: 100
componentName: "awsroute53dnsrecord"
---

# AWS Route53 DNS Record

Deploys an individual DNS record in an existing AWS Route53 hosted zone, with the full record-type set (A, AAAA, CAA, CNAME, DS, HTTPS, MX, NAPTR, NS, PTR, SOA, SPF, SRV, SSHFP, SVCB, TLSA, TXT), alias records pointing to AWS resources, health-check gating, and all seven routing policies (weighted, latency, failover, geolocation, geoproximity, CIDR, multivalue answer).

## What Gets Created

When you deploy an AwsRoute53DnsRecord resource, Planton provisions:

- **Route53 DNS Record** — a `route53.Record` (AWS Classic) resource in the specified hosted zone, configured with the given record type, values or alias target, and optional routing policy

Depending on configuration, the record may be:

- A **standard record** with explicit values and TTL (e.g., A record with IP addresses)
- An **alias record** pointing to an AWS resource (ALB, CloudFront, S3 website, API Gateway) with automatic TTL from the target
- A **policy-routed record** using weighted, latency, failover, geolocation, geoproximity, CIDR, or multivalue answer routing

## Prerequisites

- **AWS credentials** configured via environment variables or Planton provider config
- **An existing Route53 hosted zone** — provide the zone ID directly or reference an AwsRoute53Zone resource via `valueFrom`
- **A health check** (optional) for gated answers — reference an AwsRoute53HealthCheck resource via `valueFrom`
- **The target resource's DNS name and hosted zone ID** if creating alias records — or reference an AwsAlb resource via `valueFrom`

## Quick Start

Create a file `dns-record.yaml`:

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsRoute53DnsRecord
metadata:
  name: www-example
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AwsRoute53DnsRecord.www-example
spec:
  region: us-east-1
  zoneId:
    value: Z1234567890ABCDEF
  name: www.example.com
  type: A
  ttl: 300
  values:
    - "203.0.113.10"
```

Deploy:

```shell
planton apply -f dns-record.yaml
```

This creates an A record for `www.example.com` pointing to `203.0.113.10` with a 5-minute TTL.

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `region` | `string` | The AWS region where the resource will be created. | Required |
| `zoneId` | `StringValueOrRef` | The Route53 hosted zone ID where this record is created. Can reference an AwsRoute53Zone resource via `valueFrom` (default kind: `AwsRoute53Zone`, default field: `status.outputs.zone_id`). | Required |
| `name` | `string` | The fully qualified domain name or subdomain (e.g., `www.example.com`, `*.example.com`). | Required. Pattern: hostname, subdomain, or wildcard. |
| `type` | `string` | The DNS record type. Valid values: `A`, `AAAA`, `CAA`, `CNAME`, `DS`, `HTTPS`, `MX`, `NAPTR`, `NS`, `PTR`, `SOA`, `SPF`, `SRV`, `SSHFP`, `SVCB`, `TLSA`, `TXT`. | Required. |

At least one of `values` or `aliasTarget` must be specified. They are mutually exclusive.

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `ttl` | `int32` | `300` | Time to live in seconds. Valid range: 0-604800 (1 week). Ignored for alias records. Common values: 60 (fast changes), 300 (default), 86400 (static records). |
| `values` | `string[]` | `[]` | Record values. Format depends on type: A = IPv4 addresses, AAAA = IPv6, CNAME = target hostname, MX = priority + mail server, TXT = text values. Mutually exclusive with `aliasTarget`. |
| `aliasTarget` | `object` | — | Alias target for Route53 alias records. Use for zone apex records or free-query routing to AWS resources. Mutually exclusive with `values`. |
| `aliasTarget.dnsName` | `StringValueOrRef` | — | DNS name of the target resource. Can reference an AwsAlb resource via `valueFrom` (default kind: `AwsAlb`, default field: `status.outputs.load_balancer_dns_name`). Required for alias records. |
| `aliasTarget.zoneId` | `StringValueOrRef` | — | Hosted zone ID of the target AWS service (not your Route53 zone). Can reference an AwsAlb resource via `valueFrom` (default kind: `AwsAlb`, default field: `status.outputs.load_balancer_hosted_zone_id`). Required for alias records. |
| `aliasTarget.evaluateTargetHealth` | `bool` | `false` | When `true`, Route53 checks the health of the target before responding. Useful for failover scenarios. |
| `routingPolicy` | `object` | — | Routing policy for advanced traffic management. If not specified, simple routing is used. Only one policy type can be set. |
| `routingPolicy.weighted` | `object` | — | Weighted routing. Distributes traffic based on assigned weights. |
| `routingPolicy.weighted.weight` | `int32` | — | Weight value (0-255). Higher weight means more traffic. Weight of 0 stops traffic to this record. |
| `routingPolicy.latency` | `object` | — | Latency-based routing. Routes to the lowest-latency endpoint. |
| `routingPolicy.latency.region` | `string` | — | AWS region where this resource is located (e.g., `us-east-1`). Required for latency routing. |
| `routingPolicy.failover` | `object` | — | Failover routing. Automatic failover to secondary when primary fails. |
| `routingPolicy.failover.failoverType` | `string` | — | `PRIMARY` or `SECONDARY`. Required for failover routing. |
| `routingPolicy.geolocation` | `object` | — | Geolocation routing. Routes based on user location. |
| `routingPolicy.geolocation.continent` | `string` | — | Two-letter continent code (e.g., `NA`, `EU`, `AS`). Use continent or country, not both. |
| `routingPolicy.geolocation.country` | `string` | — | Two-letter ISO 3166-1 country code (e.g., `US`, `GB`, `DE`), or `*` for the default member. |
| `routingPolicy.geolocation.subdivision` | `string` | — | US state code (e.g., `CA`, `NY`). Only valid when country is `US`. |
| `routingPolicy.geoproximity` | `object` | — | Geoproximity routing: distance-based with a bias dial. Exactly one of `awsRegion`, `coordinates` (`latitude`/`longitude`), or `localZoneGroup`, plus optional `bias` (−99..99). |
| `routingPolicy.cidr` | `object` | — | CIDR routing by resolver IP block: `collectionId` + `locationName` (`*` for the default location). |
| `routingPolicy.multivalueAnswer` | `object` | — | Multivalue answer routing: up to eight healthy members per query. Not valid with alias records. The empty object is the configuration. |
| `healthCheckId` | `string \| ref` | — | Health check gating this record's answers. Can reference an AwsRoute53HealthCheck resource (default field: `status.outputs.health_check_id`). |
| `setIdentifier` | `string` | — | Unique identifier within the routing group. Required by every routing policy. Must be unique among records with the same name and type. |
| `allowOverwrite` | `bool` | `false` | Overwrite an existing record set with the same name and type instead of failing — the adoption path for records created outside the graph. |

## Examples

### Simple A Record

A basic A record pointing a subdomain to an IP address:

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsRoute53DnsRecord
metadata:
  name: api-example
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AwsRoute53DnsRecord.api-example
spec:
  zoneId:
    value: Z1234567890ABCDEF
  name: api.example.com
  type: A
  ttl: 300
  values:
    - "203.0.113.10"
    - "203.0.113.11"
```

### Alias Record Pointing to an ALB

An alias record at the zone apex pointing to an Application Load Balancer. Alias records at the apex are not possible with CNAME, and alias queries to AWS resources are free:

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsRoute53DnsRecord
metadata:
  name: apex-example
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AwsRoute53DnsRecord.apex-example
spec:
  zoneId:
    value: Z1234567890ABCDEF
  name: example.com
  type: A
  aliasTarget:
    dnsName:
      value: my-alb-123456.us-east-1.elb.amazonaws.com
    zoneId:
      value: Z35SXDOTRQ7X7K
    evaluateTargetHealth: true
```

### Using Foreign Key References

Reference other Planton-managed resources instead of hardcoding IDs. The `zoneId` defaults to kind `AwsRoute53Zone` and the alias target fields default to kind `AwsAlb`:

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsRoute53DnsRecord
metadata:
  name: app-alias
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AwsRoute53DnsRecord.app-alias
spec:
  zoneId:
    valueFrom:
      name: my-zone
  name: app.example.com
  type: A
  aliasTarget:
    dnsName:
      valueFrom:
        name: my-alb
    zoneId:
      valueFrom:
        name: my-alb
    evaluateTargetHealth: true
```

### Weighted Routing for Canary Deployment

Split traffic between two endpoints. Create two AwsRoute53DnsRecord resources with the same name and type but different `setIdentifier` and weights:

**Primary (90% traffic):**

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsRoute53DnsRecord
metadata:
  name: api-stable
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AwsRoute53DnsRecord.api-stable
spec:
  zoneId:
    value: Z1234567890ABCDEF
  name: api.example.com
  type: A
  ttl: 60
  values:
    - "203.0.113.10"
  routingPolicy:
    weighted:
      weight: 90
  setIdentifier: stable
```

**Canary (10% traffic):**

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsRoute53DnsRecord
metadata:
  name: api-canary
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AwsRoute53DnsRecord.api-canary
spec:
  region: us-east-1
  zoneId:
    value: Z1234567890ABCDEF
  name: api.example.com
  type: A
  ttl: 60
  values:
    - "203.0.113.20"
  routingPolicy:
    weighted:
      weight: 10
  setIdentifier: canary
```

### Failover with Health Check

Active-passive failover between a primary and secondary endpoint. The primary record includes a health check; when it fails, Route53 returns the secondary:

**Primary:**

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsRoute53DnsRecord
metadata:
  name: app-primary
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AwsRoute53DnsRecord.app-primary
spec:
  region: us-east-1
  zoneId:
    value: Z1234567890ABCDEF
  name: app.example.com
  type: A
  ttl: 60
  values:
    - "203.0.113.10"
  routingPolicy:
    failover:
      failoverType: PRIMARY
  healthCheckId:
    valueFrom:
      kind: AwsRoute53HealthCheck
      name: app-primary-check
      fieldPath: status.outputs.health_check_id
  setIdentifier: primary
```

**Secondary:**

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsRoute53DnsRecord
metadata:
  name: app-secondary
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AwsRoute53DnsRecord.app-secondary
spec:
  region: us-east-1
  zoneId:
    value: Z1234567890ABCDEF
  name: app.example.com
  type: A
  ttl: 60
  values:
    - "203.0.113.20"
  routingPolicy:
    failover:
      failoverType: SECONDARY
  setIdentifier: secondary
```

## Stack Outputs

After deployment, the following outputs are available in `status.outputs`:

| Output | Type | Description |
|--------|------|-------------|
| `fqdn` | `string` | The fully qualified domain name of the created record (e.g., `www.example.com`) |
| `record_type` | `string` | The DNS record type that was created (e.g., `A`, `CNAME`, `TXT`) |
| `zone_id` | `string` | The hosted zone ID where the record was created |
| `is_alias` | `bool` | Whether this is an alias record (pointing to an AWS resource) |
| `set_identifier` | `string` | The set identifier, if using routing policies. Empty for simple routing. |

## Related Components

- [AwsRoute53Zone](/docs/catalog/aws/route53-zone) — creates the hosted zone where DNS records are placed
- [AwsRoute53HealthCheck](/docs/catalog/aws/route53-health-check) — gates this record's answers on endpoint health
- [AwsAlb](/docs/catalog/aws/alb) — provides `load_balancer_dns_name` and `load_balancer_hosted_zone_id` outputs for alias record targets
- [AwsCloudfront](/docs/catalog/aws/cloudfront) — provides the distribution DNS name for alias record targets
- [AwsCertManagerCert](/docs/catalog/aws/certificate) — DNS-validated certificates require TXT or CNAME records in the zone
