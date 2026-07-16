# AWS SES Email Identity

Verify a domain or email address for sending with Amazon SES. A domain identity is the production shape: DKIM-signed mail for every address at the domain, DNS-based verification whose CNAME tokens compose directly into Route 53 records, and a custom MAIL FROM domain for DMARC-aligned SPF.

## What Gets Created

- An SESv2 email identity (domain or email address) with Easy DKIM or BYODKIM signing.
- Optionally, the identity's custom MAIL FROM configuration and bounce/complaint email-forwarding switch.
- One SESv2 identity policy per named `policies[]` entry — the cross-account sending grants.

## Prerequisites

- Control of the domain's DNS (to publish the DKIM CNAMEs) for domain identities, or access to the mailbox for email-address identities.
- Optionally an `AwsSesConfigurationSet` to attach as the identity's default sending rules.

## Quick Start

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsSesEmailIdentity
metadata:
  name: prod-sender
spec:
  region: us-west-2
  emailIdentity: example.com
  dkimSigning:
    nextSigningKeyLength: RSA_2048_BIT
  mailFrom:
    mailFromDomain: mail.example.com
```

## Configuration Reference

### Required Fields

| Field | Description |
|---|---|
| `region` | AWS region; verify the same domain in each region you send from. |
| `emailIdentity` | The domain (`example.com`) or address (`sender@example.com`) to verify. Immutable. |

### Common Optional Fields

| Field | Description |
|---|---|
| `configurationSet` | Default configuration set for every message sent from this identity. |
| `dkimSigning.nextSigningKeyLength` | Easy DKIM key rotation length: `RSA_2048_BIT` (recommended) or `RSA_1024_BIT`. |
| `dkimSigning.domainSigningPrivateKey` + `domainSigningSelector` | BYODKIM: bring your own key pair (the private key is a managed secret). |
| `mailFrom.mailFromDomain` | Custom MAIL FROM subdomain for aligned SPF (needs MX + SPF records). |
| `mailFrom.behaviorOnMxFailure` | `USE_DEFAULT_VALUE` (keep sending) or `REJECT_MESSAGE` (strict alignment). |
| `emailForwardingEnabled` | Forward bounces/complaints by email (default true); disable once event destinations carry feedback. |
| `policies[]` | Named authorization policies granting other accounts `ses:SendEmail` AS this identity. |

## Stack Outputs

| Output | Description |
|---|---|
| `identity_arn` | The identity's ARN for policy grants and IAM statements. |
| `email_identity` | The identity string — the DNS-name join key. |
| `identity_type` | `DOMAIN` or `EMAIL_ADDRESS`. |
| `verification_status` | `PENDING` until verification completes, then `SUCCESS`. |
| `dkim_tokens` | Easy DKIM's three CNAME tokens, composable into `AwsRoute53DnsRecord` nodes. |

## Related Components

- [AWS SES Configuration Set](/docs/catalog/aws/awssesconfigurationset) — the sending rules and event publishing this identity opts into.
- [AWS Route53 DNS Record](/docs/catalog/aws/awsroute53dnsrecord) — publishes the DKIM CNAMEs and MAIL FROM MX/SPF records.
- [AWS Route53 Zone](/docs/catalog/aws/awsroute53zone) — the hosted zone those records live in.
