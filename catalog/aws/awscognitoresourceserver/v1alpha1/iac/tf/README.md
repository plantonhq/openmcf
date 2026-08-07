# AwsCognitoResourceServer — Terraform Module

## Overview

This Terraform module provisions a Cognito resource server — the OAuth 2.0 resource (an API) a user pool mints custom access-token scopes for. Machine-to-machine clients using the `client_credentials` grant can only request custom scopes, which is exactly what this resource defines.

## Module Structure

```
main.tf       — aws_cognito_resource_server
locals.tf     — computed scope identifiers ({identifier}/{scope_name})
outputs.tf    — resource_server_identifier + scope_identifiers
variables.tf  — Generator-owned typed contract (metadata, spec)
provider.tf   — AWS provider configuration (>= 6.0.0)
```

## Usage

```hcl
module "orders_api" {
  source = "./path/to/module"

  metadata = {
    name = "orders-api"
    org  = "my-org"
    env  = "prod"
    id   = "orders-api-prod"
  }

  spec = {
    region       = "us-west-2"
    user_pool_id = "us-west-2_Ab1Cd2EfG"
    identifier   = "https://api.example.com"
    name         = "orders-api"

    scopes = [
      { scope_name = "read", scope_description = "Read access to orders" },
      { scope_name = "orders:write", scope_description = "Write access to orders" },
    ]
  }
}
```

## Notes

- `identifier` is ForceNew — it is the resource server's identity within the pool and the prefix of every scope it mints.
- Scopes update in place; removing one invalidates it for future tokens while already-issued tokens carry it until they expire.
- The exported `scope_identifiers` are the exact strings app clients list in `allowed_oauth_scopes`.
