# AwsIamSamlProvider — Terraform/OpenTofu module

Manages one IAM SAML identity provider (`aws_iam_saml_provider`).

Module facts worth knowing before editing:

- **The name comes from `metadata.name` and is WRITE-ONCE at AWS** —
  a rename replaces the provider and invalidates every role trust
  policy naming its ARN.
- **The metadata document updates IN PLACE** — certificate rotations
  are ordinary updates; the `valid_until` output is the
  rotate-by date.
- **IAM is global**: the provider exists account-wide regardless of
  the endpoint region the stack ran against.
- **The provider is taggable** and carries the catalog's identity
  tags.

Outputs mirror the Pulumi module key-for-key: `provider_arn`,
`saml_provider_uuid`, `valid_until`.
