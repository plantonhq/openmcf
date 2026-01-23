# AWS Route53 DNS Record Examples

This document provides working, copy-paste ready examples for common Route53 DNS record configurations.

## Basic A Record

Simple A record pointing a subdomain to an IP address.

```yaml
apiVersion: aws.project-planton.org/v1
kind: AwsRoute53DnsRecord
metadata:
  name: www-a-record
spec:
  hosted_zone_id: Z1234567890ABC
  name: www.example.com
  type: A
  ttl: 300
  values:
    - 192.0.2.1
```

## A Record with Multiple IPs (Round Robin)

DNS-based load balancing across multiple servers.

```yaml
apiVersion: aws.project-planton.org/v1
kind: AwsRoute53DnsRecord
metadata:
  name: api-round-robin
spec:
  hosted_zone_id: Z1234567890ABC
  name: api.example.com
  type: A
  ttl: 60
  values:
    - 192.0.2.1
    - 192.0.2.2
    - 192.0.2.3
```

## CNAME Record

Alias a subdomain to another domain name.

```yaml
apiVersion: aws.project-planton.org/v1
kind: AwsRoute53DnsRecord
metadata:
  name: blog-cname
spec:
  hosted_zone_id: Z1234567890ABC
  name: blog.example.com
  type: CNAME
  ttl: 300
  values:
    - example.ghost.io
```

## MX Records for Email

Configure mail exchange servers with priority.

```yaml
apiVersion: aws.project-planton.org/v1
kind: AwsRoute53DnsRecord
metadata:
  name: mail-mx-records
spec:
  hosted_zone_id: Z1234567890ABC
  name: example.com
  type: MX
  ttl: 3600
  values:
    - "10 mail1.example.com"
    - "20 mail2.example.com"
    - "30 mail3.example.com"
```

## TXT Record for SPF

Email authentication SPF record.

```yaml
apiVersion: aws.project-planton.org/v1
kind: AwsRoute53DnsRecord
metadata:
  name: spf-record
spec:
  hosted_zone_id: Z1234567890ABC
  name: example.com
  type: TXT
  ttl: 300
  values:
    - "v=spf1 include:_spf.google.com include:servers.mcsv.net ~all"
```

## TXT Record for Domain Verification

Verify domain ownership for Google Workspace, etc.

```yaml
apiVersion: aws.project-planton.org/v1
kind: AwsRoute53DnsRecord
metadata:
  name: google-verification
spec:
  hosted_zone_id: Z1234567890ABC
  name: example.com
  type: TXT
  ttl: 300
  values:
    - "google-site-verification=abc123xyz789"
```

## Wildcard A Record

Catch-all record for any subdomain.

```yaml
apiVersion: aws.project-planton.org/v1
kind: AwsRoute53DnsRecord
metadata:
  name: wildcard-record
spec:
  hosted_zone_id: Z1234567890ABC
  name: "*.example.com"
  type: A
  ttl: 300
  values:
    - 192.0.2.1
```

## Alias Record to CloudFront Distribution

Point zone apex to CloudFront (free queries, no CNAME restriction).

```yaml
apiVersion: aws.project-planton.org/v1
kind: AwsRoute53DnsRecord
metadata:
  name: apex-cloudfront
spec:
  hosted_zone_id: Z1234567890ABC
  name: example.com
  type: A
  alias_target:
    dns_name: d1234abcd.cloudfront.net
    hosted_zone_id: Z2FDTNDATAQYW2
    evaluate_target_health: false
```

## Alias Record to Application Load Balancer

Point subdomain to ALB with health evaluation.

```yaml
apiVersion: aws.project-planton.org/v1
kind: AwsRoute53DnsRecord
metadata:
  name: api-alb-alias
spec:
  hosted_zone_id: Z1234567890ABC
  name: api.example.com
  type: A
  alias_target:
    dns_name: my-alb-1234567890.us-east-1.elb.amazonaws.com
    hosted_zone_id: Z35SXDOTRQ7X7K
    evaluate_target_health: true
```

## Alias Record to S3 Website

Point domain to S3 static website hosting.

```yaml
apiVersion: aws.project-planton.org/v1
kind: AwsRoute53DnsRecord
metadata:
  name: static-s3-alias
spec:
  hosted_zone_id: Z1234567890ABC
  name: static.example.com
  type: A
  alias_target:
    dns_name: my-bucket.s3-website-us-east-1.amazonaws.com
    hosted_zone_id: Z3AQBSTGFYJSTF
    evaluate_target_health: false
```

## Weighted Routing - Blue/Green Deployment

Split traffic between two versions (70% blue, 30% green).

**Blue environment (70% traffic):**

```yaml
apiVersion: aws.project-planton.org/v1
kind: AwsRoute53DnsRecord
metadata:
  name: api-weighted-blue
spec:
  hosted_zone_id: Z1234567890ABC
  name: api.example.com
  type: A
  ttl: 60
  values:
    - 192.0.2.1
  routing_policy:
    weighted:
      weight: 70
  set_identifier: blue
```

**Green environment (30% traffic):**

```yaml
apiVersion: aws.project-planton.org/v1
kind: AwsRoute53DnsRecord
metadata:
  name: api-weighted-green
spec:
  hosted_zone_id: Z1234567890ABC
  name: api.example.com
  type: A
  ttl: 60
  values:
    - 192.0.2.2
  routing_policy:
    weighted:
      weight: 30
  set_identifier: green
```

## Latency-Based Routing - Multi-Region

Route users to the lowest-latency endpoint.

**US East endpoint:**

```yaml
apiVersion: aws.project-planton.org/v1
kind: AwsRoute53DnsRecord
metadata:
  name: api-latency-us-east
spec:
  hosted_zone_id: Z1234567890ABC
  name: api.example.com
  type: A
  ttl: 60
  values:
    - 192.0.2.1
  routing_policy:
    latency:
      region: us-east-1
  set_identifier: us-east-1
```

**EU West endpoint:**

```yaml
apiVersion: aws.project-planton.org/v1
kind: AwsRoute53DnsRecord
metadata:
  name: api-latency-eu-west
spec:
  hosted_zone_id: Z1234567890ABC
  name: api.example.com
  type: A
  ttl: 60
  values:
    - 192.0.2.2
  routing_policy:
    latency:
      region: eu-west-1
  set_identifier: eu-west-1
```

## Failover Routing - Disaster Recovery

Automatic failover when primary fails health check.

**Primary record:**

```yaml
apiVersion: aws.project-planton.org/v1
kind: AwsRoute53DnsRecord
metadata:
  name: www-failover-primary
spec:
  hosted_zone_id: Z1234567890ABC
  name: www.example.com
  type: A
  ttl: 60
  values:
    - 192.0.2.1
  routing_policy:
    failover:
      failover_type: primary
  set_identifier: primary
  health_check_id: abcd1234-5678-90ab-cdef-example
```

**Secondary record:**

```yaml
apiVersion: aws.project-planton.org/v1
kind: AwsRoute53DnsRecord
metadata:
  name: www-failover-secondary
spec:
  hosted_zone_id: Z1234567890ABC
  name: www.example.com
  type: A
  ttl: 60
  values:
    - 192.0.2.2
  routing_policy:
    failover:
      failover_type: secondary
  set_identifier: secondary
```

## Geolocation Routing - GDPR Compliance

Route EU users to EU servers for data residency.

**EU users:**

```yaml
apiVersion: aws.project-planton.org/v1
kind: AwsRoute53DnsRecord
metadata:
  name: api-geo-eu
spec:
  hosted_zone_id: Z1234567890ABC
  name: api.example.com
  type: A
  ttl: 300
  values:
    - 192.0.2.1
  routing_policy:
    geolocation:
      continent: EU
  set_identifier: europe
```

**US users:**

```yaml
apiVersion: aws.project-planton.org/v1
kind: AwsRoute53DnsRecord
metadata:
  name: api-geo-us
spec:
  hosted_zone_id: Z1234567890ABC
  name: api.example.com
  type: A
  ttl: 300
  values:
    - 192.0.2.2
  routing_policy:
    geolocation:
      country: US
  set_identifier: us
```

**Default (catch-all for other locations):**

```yaml
apiVersion: aws.project-planton.org/v1
kind: AwsRoute53DnsRecord
metadata:
  name: api-geo-default
spec:
  hosted_zone_id: Z1234567890ABC
  name: api.example.com
  type: A
  ttl: 300
  values:
    - 192.0.2.3
  routing_policy:
    geolocation: {}
  set_identifier: default
```

## CAA Record - Certificate Authority Authorization

Restrict which CAs can issue certificates for your domain.

```yaml
apiVersion: aws.project-planton.org/v1
kind: AwsRoute53DnsRecord
metadata:
  name: caa-record
spec:
  hosted_zone_id: Z1234567890ABC
  name: example.com
  type: CAA
  ttl: 3600
  values:
    - '0 issue "letsencrypt.org"'
    - '0 issue "amazon.com"'
    - '0 issuewild ";"'
```

## Deployment

Deploy any of these examples using the Project Planton CLI:

```bash
# Deploy a single record
project-planton pulumi up --manifest dns-record.yaml

# Deploy with OpenTofu/Terraform
project-planton tofu apply --manifest dns-record.yaml
```
