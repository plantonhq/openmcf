# SES Email Identities: The Trust Anchor of Sending

## Introduction

Nothing leaves SES without a verified identity. An identity is the domain or email address an application is allowed to send FROM, and it carries the sending domain's cryptographic reputation: DKIM signatures, the MAIL FROM envelope domain SPF evaluates, and the DMARC alignment both feed. A DOMAIN identity is the production shape — one verification covers every address at the domain. An EMAIL_ADDRESS identity verifies through a mailbox confirmation link — quick for testing, but it cannot carry its own DKIM configuration (it inherits the domain's when the domain is also verified).

## The Verification Composition Story

Domain identities verify through DNS, and this is where the kind earns first-class status: creating the identity returns three Easy DKIM tokens, exported as the `dkim_tokens` stack output. Each becomes a CNAME —

```
<token>._domainkey.<domain>  CNAME  <token>.dkim.amazonses.com
```

— composable directly into `AwsRoute53DnsRecord` nodes. SES flips `verification_status` from PENDING to SUCCESS once it observes them. The identity resource itself creates instantly; verification is asynchronous DNS convergence and never blocks the deploy.

## DKIM: Easy vs BYODKIM

- **Easy DKIM** (recommended): AWS generates and rotates the RSA key pair; you publish the three CNAMEs. `nextSigningKeyLength` selects RSA_2048_BIT (the default; prefer it) or RSA_1024_BIT (exists for DNS providers with 255-character TXT limits). Setting it on a live identity rotates the key.
- **BYODKIM**: bring your own key pair — the base64 PKCS #8 private key (a managed secret, never logged or exported) plus the selector you publish the public key under. The key and selector only work as a pair, and BYODKIM is mutually exclusive with Easy DKIM — both constraints are spec-level CEL.

## MAIL FROM and DMARC Alignment

By default SES uses its own bounce domain (amazonses.com) as the envelope MAIL FROM, which fails strict DMARC alignment on SPF. A custom MAIL FROM subdomain (`mail.example.com` with MX + SPF records, composable via `AwsRoute53DnsRecord`) aligns the envelope with the sending domain. `behaviorOnMxFailure` decides what happens when the MX record is missing: keep sending unaligned (`USE_DEFAULT_VALUE`, the AWS default) or fail the send (`REJECT_MESSAGE`).

## Design Notes

- **The identity string is the AWS identifier** — deliberately from `spec.emailIdentity`, not `metadata.name`, because it must be the exact DNS name mail is sent from. Changing it replaces the identity.
- **Folded satellites.** The custom MAIL FROM configuration, the bounce/complaint email-forwarding switch, and named authorization policies are each their own AWS sub-resource keyed by the identity — 1:1 (or named-many) with the identity's lifecycle and never referenced by anything else. Both engines materialize each independently.
- **Identity policies are the cross-account grant.** A named IAM-syntax policy on the identity lets another account or role send AS it (`ses:SendEmail`/`ses:SendRawEmail` with this identity as the source). AWS validates the document strictly at create time: `Resource` must be the identity's OWN ARN (`arn:aws:ses:<region>:<account>:identity/<identity>`) — a wildcard `Resource` is rejected with "Invalid ARN". Author policies with the concrete ARN of the deployed identity.
- **Feedback forwarding is a tri-state.** Unset accepts AWS's default (forwarding ON) without materializing a resource; an explicit value manages it.

## 90/10 Coverage Notes

The full `aws_sesv2_email_identity` surface plus its three satellite resources (`_mail_from_attributes`, `_feedback_attributes`, `_policy`) is modeled, including both DKIM arms with their pairing/exclusivity rules as CEL and the reject-requires-MX honesty of `behaviorOnMxFailure`.

## Deferred Surface (recorded reasons)

- **Dedicated IP pools, contact lists, tenants, account-level attributes** — see the configuration-set research doc; all are set- or account-plane surfaces, not identity shape.
- **Classic SES receiving** (receipt rules/filters) and `aws_ses_template` — receiving and templating are separate product surfaces; templates additionally have no SESv2 provider resource today.
