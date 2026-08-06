# AwsHttpApiDomain — Terraform Module

## Overview

This Terraform module provisions an API Gateway v2 custom domain name and its API mappings: TLS termination with an ACM certificate, optional mutual TLS, and one mapping resource per configured API.

## Module Structure

```
main.tf       — aws_apigatewayv2_domain_name + aws_apigatewayv2_api_mapping (per entry)
locals.tf     — identity tags, mapping addressing (root alias for the empty key)
outputs.tf    — domain name/arn, target domain, hosted zone id
variables.tf  — Generator-owned typed contract (metadata, spec)
provider.tf   — AWS provider configuration (>= 6.29.0)
```

## Usage

```hcl
module "api_domain" {
  source = "./path/to/module"

  metadata = {
    name = "api-example-com"
    org  = "my-org"
    env  = "prod"
    id   = "api-example-com-prod"
  }

  spec = {
    region          = "us-east-1"
    domain_name     = "api.example.com"
    certificate_arn = "arn:aws:acm:us-east-1:123456789012:certificate/abc-123"
    api_mappings = [
      { api_id = "a1b2c3d4", stage = "$default", api_mapping_key = "" }
    ]
  }
}
```

Note: `variables.tf` is generated from the proto spec by `planton tofu
generate-variables AwsHttpApiDomain` and guarded against drift -- reference
fields (`certificate_arn`, `api_mappings[].api_id`) arrive pre-resolved as
plain strings.

## Outputs

| Output | Description |
|--------|-------------|
| `domain_name` | The custom domain name (the domain's join key) |
| `domain_name_arn` | ARN of the domain name resource |
| `target_domain_name` | The regional domain to alias/CNAME to from DNS |
| `hosted_zone_id` | Route 53 alias target zone of the regional endpoint |

## Implementation Notes

- **Endpoint type / security policy hardcoded**: v2 domains accept only REGIONAL + TLS_1_2, so neither is a spec field.
- **Root mapping addressing**: the empty mapping key (the domain root) is addressed internally as `(root)` because `for_each` needs non-empty keys; `api_mapping_key` still sends the real (absent) value to AWS.
- **domain_name is immutable**: renaming replaces the domain and its mappings; the DNS alias must follow.
