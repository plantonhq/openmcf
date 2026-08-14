---
title: "Preset: Threat Protected"
description: "**Rank**: 4"
type: "preset"
rank: "04"
presetSlug: "04-threat-protected"
componentSlug: "cognito-user-pool"
componentTitle: "Cognito User Pool"
provider: "aws"
icon: "package"
order: 4
---

# Preset: Threat Protected

**Rank**: 4

## When to Use

- User directories under active credential-stuffing or account-takeover pressure
- Applications whose compliance posture requires automated responses to risky sign-ins, not just detection
- Teams that want users NOTIFIED when the platform blocks or challenges an attempt on their account

## What It Provides

- PLUS feature tier with threat protection ENFORCED, plus the full automated-response policy the mode acts on:
  - High-risk sign-ins are BLOCKED and the user is emailed
  - Medium-risk sign-ins require MFA when the user has a factor enrolled, with a notification
  - Compromised credentials (username/password pairs found in known breach data) are blocked outright
  - A trusted office CIDR skips risk evaluation entirely
- Notification emails sent through your verified SES identity with block and MFA templates
- User groups (`admins`, `members`) whose membership lands in the `cognito:groups` token claim -- the admin group carries an IAM role for identity-pool federation
- The production-hardened baseline: strict password policy with history, optional TOTP MFA, verified-before-update email, deletion protection, hosted-UI domain

## What You Might Add

- Per-client overrides: an `AwsCognitoUserPoolClient` with its own `riskConfiguration` when one app needs a stricter or looser posture than the pool-wide policy
- `logConfigurations` routing `userAuthEvents` to a CloudWatch log group -- threat protection's risk verdicts become queryable telemetry
- `blockedIpRanges` under `riskException` for known-hostile networks
- `noActionEmail` template if you set a `NO_ACTION` response with `notify: true` (watch-only posture)
