# AwsSesEmailIdentity Pulumi Module

Provisions an Amazon SES (SESv2) email identity with DKIM signing, optional
default configuration set, MAIL FROM domain, feedback forwarding, and
per-name authorization policies.

## Resources Created

- `aws:sesv2:EmailIdentity` — The identity (`spec.email_identity` is the AWS id)
- `aws:sesv2:EmailIdentityMailFromAttributes` — When `mail_from` is set
- `aws:sesv2:EmailIdentityFeedbackAttributes` — When `email_forwarding_enabled` is explicit
- `aws:sesv2:EmailIdentityPolicy` — One per named policy

## Outputs

| Key | Description |
|-----|-------------|
| `identity_arn` | ARN of the email identity |
| `email_identity` | The identity string |
| `identity_type` | `DOMAIN` or `EMAIL_ADDRESS` |
| `verification_status` | Verification status at deploy time |
| `dkim_tokens` | Easy DKIM CNAME tokens |

## Local Development

```bash
./debug.sh preview
./debug.sh up
./debug.sh destroy
```
