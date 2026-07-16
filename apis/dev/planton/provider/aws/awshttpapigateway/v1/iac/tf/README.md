# AwsHttpApiGateway — Terraform Module

## Overview

This Terraform module provisions an API Gateway HTTP API (v2) as one coherent unit: the API, its stage, deduplicated integrations (Lambda proxy, HTTP proxy, private VPC-link proxy, or AWS service actions), optional authorizers, and the routes that bind them.

## Module Structure

```
main.tf       — api, stage, integrations, authorizers, routes
locals.tf     — identity tags, stage presence semantics, whole-object integration dedup
outputs.tf    — api id/endpoint/arn, execution arn, stage invoke url/name
variables.tf  — Generator-owned typed contract (metadata, spec)
provider.tf   — AWS provider configuration (>= 6.29.0)
```

## Usage

```hcl
module "http_api" {
  source = "./path/to/module"

  metadata = {
    name = "orders-api"
    org  = "my-org"
    env  = "prod"
    id   = "orders-api-prod"
  }

  spec = {
    region = "us-east-1"
    routes = [
      {
        route_key = "$default"
        integration = {
          integration_type = "AWS_PROXY"
          integration_uri  = "arn:aws:lambda:us-east-1:123456789012:function:orders"
        }
      }
    ]
  }
}
```

Note: `variables.tf` is generated from the proto spec by `planton tofu
generate-variables AwsHttpApiGateway` and guarded against drift -- reference
fields (`integration_uri`, `connection_id`, `credentials_arn`,
`access_log.destination_arn`, authorizer refs) arrive pre-resolved as plain
strings.

## Outputs

| Output | Description |
|--------|-------------|
| `api_id` | The API identifier (the join key domain mappings reference) |
| `api_endpoint` | The default endpoint URL |
| `api_arn` | ARN of the API |
| `execution_arn` | Execution ARN prefix for Lambda resource policies |
| `stage_invoke_url` | Invoke URL of the deployed stage |
| `stage_name` | Name of the deployed stage |

## Implementation Notes

- **Integration dedup**: routes whose integration blocks are identical across every field share one integration resource, keyed on a hash of the whole object -- a partial key would silently merge integrations that differ in newer fields.
- **auto_deploy presence rule**: an omitted `auto_deploy` means true (a declarative spec should be self-applying); only an explicit `false` turns it off. Identical rule in the Pulumi module.
- **Stage after routes**: per-route `route_settings` reference routes by key, so the stage carries `depends_on` for the routes.
- **Service integrations**: `integration_subtype` integrations send no `integration_uri` (the action's parameters ride `request_parameters`) and are pinned to payload format 1.0 by AWS.
- **WebSocket-only knobs** (`logging_level`, `data_trace_enabled`) are deliberately not modeled -- AWS ignores them for HTTP APIs.
