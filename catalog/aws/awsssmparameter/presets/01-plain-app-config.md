# Plain App Config

This preset creates a readable configuration parameter under a
hierarchical app path, with a validation pattern AWS enforces on every
future write.

## When to Use

- Ordinary application configuration that is not a secret — log
  levels, feature switches, endpoints
- Config organized by path so services read their whole subtree with
  one by-path call

## What You Get

- A free Standard-tier String parameter whose value stays readable in
  plans (the plain arm renders as the provider's `insecure_value` —
  by design)
- An `allowedPattern` guard rejecting invalid future writes at the AWS
  API

## Customize

- Switch `type: StringList` for comma-separated lists (one string —
  values with commas need String)
- Add `dataType: aws:ec2:image` to have AWS validate values as real
  AMI IDs on every write
- Secrets belong in the `secureValue` arm with `type: SecureString` —
  see the secure-secret preset
