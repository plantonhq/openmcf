# AwsSesEmailIdentity Terraform Module

Provisions an Amazon SES (SESv2) email identity with optional DKIM signing,
default configuration set, custom MAIL FROM domain, feedback forwarding, and
per-name authorization policies.

## Resources Created

- `aws_sesv2_email_identity.this` — The identity (`spec.email_identity` is the AWS id)
- `aws_sesv2_email_identity_mail_from_attributes.this` — When `mail_from` is set
- `aws_sesv2_email_identity_feedback_attributes.this` — When `email_forwarding_enabled` is explicit
- `aws_sesv2_email_identity_policy.this` — One per named policy

## Inputs

| Variable | Description |
|----------|-------------|
| `metadata` | Resource metadata (name, org, env, id, labels) |
| `spec` | AwsSesEmailIdentitySpec — desired configuration |

## Outputs

| Output | Description |
|--------|-------------|
| `identity_arn` | ARN of the email identity |
| `email_identity` | The identity string (domain or address) |
| `identity_type` | `DOMAIN` or `EMAIL_ADDRESS` |
| `verification_status` | Verification status at deploy time |
| `dkim_tokens` | Easy DKIM CNAME tokens (empty for BYODKIM / email-address identities) |

## Provider

Requires `hashicorp/aws` provider version `>= 6.0.0`.
