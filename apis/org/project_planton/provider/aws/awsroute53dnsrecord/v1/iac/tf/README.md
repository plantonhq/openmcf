# AWS Route53 DNS Record - Terraform Module

This Terraform module creates and manages AWS Route53 DNS records.

## Overview

The module provisions a single DNS record in an existing Route53 hosted zone with support for:

- Standard DNS record types (A, AAAA, CNAME, MX, TXT, etc.)
- Alias records (Route53's killer feature for AWS resources)
- Advanced routing policies (weighted, latency, failover, geolocation)
- Health check integration

## Usage

### With Project Planton CLI

```bash
# Deploy a DNS record
project-planton tofu apply --manifest dns-record.yaml

# Preview changes
project-planton tofu plan --manifest dns-record.yaml

# Destroy the record
project-planton tofu destroy --manifest dns-record.yaml
```

### Standalone Usage

```hcl
module "dns_record" {
  source = "./path/to/module"

  metadata = {
    name = "www-example"
  }

  spec = {
    hosted_zone_id = "Z1234567890ABC"
    name           = "www.example.com"
    type           = "A"
    ttl            = 300
    values         = ["192.0.2.1"]
  }
}
```

## Variables

### Required Variables

| Variable | Type | Description |
|----------|------|-------------|
| `metadata.name` | string | Resource name |
| `spec.hosted_zone_id` | string | Route53 hosted zone ID |
| `spec.name` | string | DNS record name (FQDN) |
| `spec.type` | string | Record type (A, AAAA, CNAME, etc.) |

### Optional Variables

| Variable | Type | Default | Description |
|----------|------|---------|-------------|
| `spec.ttl` | number | 300 | TTL in seconds (ignored for alias) |
| `spec.values` | list(string) | [] | Record values (for standard records) |
| `spec.alias_target` | object | null | Alias target configuration |
| `spec.routing_policy` | object | null | Routing policy configuration |
| `spec.health_check_id` | string | null | Health check ID for failover |
| `spec.set_identifier` | string | null | Set ID for routing policies |

## Outputs

| Output | Description |
|--------|-------------|
| `fqdn` | Fully qualified domain name |
| `record_type` | DNS record type |
| `hosted_zone_id` | Hosted zone ID |
| `is_alias` | Whether this is an alias record |
| `set_identifier` | Set identifier if using routing policies |

## Examples

### Basic A Record

```hcl
spec = {
  hosted_zone_id = "Z1234567890ABC"
  name           = "www.example.com"
  type           = "A"
  ttl            = 300
  values         = ["192.0.2.1", "192.0.2.2"]
}
```

### Alias to CloudFront

```hcl
spec = {
  hosted_zone_id = "Z1234567890ABC"
  name           = "example.com"
  type           = "A"
  alias_target = {
    dns_name               = "d1234abcd.cloudfront.net"
    hosted_zone_id         = "Z2FDTNDATAQYW2"
    evaluate_target_health = false
  }
}
```

### Weighted Routing

```hcl
spec = {
  hosted_zone_id = "Z1234567890ABC"
  name           = "api.example.com"
  type           = "A"
  ttl            = 60
  values         = ["192.0.2.1"]
  routing_policy = {
    weighted = {
      weight = 70
    }
  }
  set_identifier = "primary"
}
```

### Failover Routing

```hcl
spec = {
  hosted_zone_id  = "Z1234567890ABC"
  name            = "www.example.com"
  type            = "A"
  ttl             = 60
  values          = ["192.0.2.1"]
  routing_policy  = {
    failover = {
      failover_type = "primary"
    }
  }
  set_identifier  = "primary"
  health_check_id = "abcd1234-5678-90ab-cdef-example"
}
```

## Notes

- TTL is ignored for alias records (uses target resource's TTL)
- `set_identifier` is required for weighted, latency, failover, and geolocation routing
- For alias targets, use the AWS service's hosted zone ID, not your Route53 zone ID
