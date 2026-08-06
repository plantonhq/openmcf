# AwsCognitoUserPool — Terraform Module

## Overview

This Terraform module provisions an Amazon Cognito User Pool at full provider depth — identity model, password and sign-in (passwordless) policies, MFA (TOTP/email/WebAuthn), SMS delivery, verification/invitation messaging, custom schema attributes, Lambda triggers, threat protection — plus the two pool-scoped folded satellites: the hosted-UI domain (one per pool) and log delivery.

App clients, identity providers, and resource servers are separate kinds (`AwsCognitoUserPoolClient`, `AwsCognitoIdentityProvider`, `AwsCognitoResourceServer`) that compose onto the pool by reference.

## Module Structure

```
main.tf       — aws_cognito_user_pool + aws_cognito_user_pool_domain + aws_cognito_log_delivery_configuration
locals.tf     — identity tags, domain-shape detection (prefix vs custom)
outputs.tf    — pool id/arn/endpoint, issuer, domain join keys, CloudFront alias trio
variables.tf  — Generator-owned typed contract (metadata, spec)
provider.tf   — AWS provider configuration (>= 6.11.0)
```

## Usage

```hcl
module "user_pool" {
  source = "./path/to/module"

  metadata = {
    name = "my-app-users"
    org  = "my-org"
    env  = "prod"
    id   = "my-app-users-prod"
  }

  spec = {
    region              = "us-west-2"
    username_attributes = ["email"]

    password_policy = {
      minimum_length    = 12
      require_lowercase = true
      require_uppercase = true
      require_numbers   = true
      require_symbols   = true
    }

    mfa_configuration          = "OPTIONAL"
    software_token_mfa_enabled = true

    auto_verified_attributes = ["email"]

    domain = {
      domain = "my-app-auth"
    }
  }
}
```

## Notes

- The identity model (`username_attributes` / `alias_attributes` / `username_case_sensitive`) is ForceNew: changing it replaces the pool and destroys every user in it.
- The custom-attribute schema is append-only in AWS: additions apply in place; removing or modifying an existing attribute errors.
- A domain containing a dot is a custom domain and requires `certificate_arn` (an ACM certificate in us-east-1); AWS fronts it with a managed CloudFront distribution exported for DNS aliasing.
- `user_pool_endpoint` is exported exactly as AWS reports it (no scheme); `issuer` carries the `https://` scheme and is the value JWT authorizers consume.
