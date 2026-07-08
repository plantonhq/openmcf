---
title: "Domain with Easy DKIM"
description: "Verifies a sending domain with AWS-managed DKIM keys. The stack output `dkim_tokens` carries three CNAME names to publish via `AwsRoute53DnsRecord`."
type: "preset"
rank: "01"
presetSlug: "01-domain-easy-dkim"
componentSlug: "ses-email-identity"
componentTitle: "SES Email Identity"
provider: "aws"
icon: "package"
order: 1
---

# Domain with Easy DKIM

Verifies a sending domain with AWS-managed DKIM keys. The stack output
`dkim_tokens` carries three CNAME names to publish via `AwsRoute53DnsRecord`.

## When to Use

- Production sending from any address at a domain you control
- The common path: Easy DKIM with a 2048-bit key (AWS default)

## What It Configures

- **`emailIdentity: example.com`** — the domain being verified (this is the AWS id)
- **Easy DKIM** via `nextSigningKeyLength: RSA_2048_BIT`

## What to Customize

- Replace `example.com` and `<aws-region>`
- Publish the three `dkim_tokens` CNAMEs to complete verification
