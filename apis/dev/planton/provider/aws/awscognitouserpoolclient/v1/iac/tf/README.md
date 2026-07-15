# AwsCognitoUserPoolClient — Terraform Module

## Overview

This Terraform module provisions a Cognito User Pool app client — the OAuth 2.0 / OIDC contract between one application and a user pool: grant types, scopes, redirect URLs, explicit auth flows, token lifetimes with units, refresh-token rotation, user-existence-error posture, attribute access, and Pinpoint analytics.

## Module Structure

```
main.tf       — aws_cognito_user_pool_client
locals.tf     — naming basis (the client resource is not taggable)
outputs.tf    — client_id + client_secret (sensitive)
variables.tf  — Generator-owned typed contract (metadata, spec)
provider.tf   — AWS provider configuration (>= 6.0.0)
```

## Usage

```hcl
module "web_client" {
  source = "./path/to/module"

  metadata = {
    name = "web-app"
    org  = "my-org"
    env  = "prod"
    id   = "web-app-prod"
  }

  spec = {
    region       = "us-west-2"
    user_pool_id = "us-west-2_Ab1Cd2EfG"

    allowed_oauth_flows_user_pool_client = true
    allowed_oauth_flows                  = ["code"]
    allowed_oauth_scopes                 = ["openid", "email", "profile"]
    callback_urls                        = ["https://app.example.com/callback"]
    logout_urls                          = ["https://app.example.com/"]
    explicit_auth_flows                  = ["ALLOW_USER_SRP_AUTH", "ALLOW_REFRESH_TOKEN_AUTH"]
    prevent_user_existence_errors        = "ENABLED"
  }
}
```

## Notes

- `generate_secret` is ForceNew: confidential (server-side, secret-holding) vs public (SPA/mobile, PKCE) is decided at creation.
- Token lifetimes pair with `token_validity_units` (defaults: hours/hours/days); AWS bounds access/ID tokens to 5 minutes – 24 hours and refresh tokens to 60 minutes – 10 years.
- `enable_token_revocation` is presence-aware: leaving it unset keeps AWS's default (enabled) instead of silently flipping it off.
