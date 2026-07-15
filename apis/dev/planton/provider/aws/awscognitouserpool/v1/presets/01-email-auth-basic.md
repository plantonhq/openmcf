# Preset: Email Auth Basic

**Rank**: 1 (most common starting point)

## When to Use

- Getting started with Cognito
- Development and testing environments
- Simple applications needing email-based sign-up and sign-in

## What It Provides

- Email as the sign-in identifier (ForceNew -- chosen deliberately up front)
- Auto-verified email addresses
- Password recovery via email
- No MFA, no domain, no custom attributes

## What You Might Add

- An `AwsCognitoUserPoolClient` resource so an application can authenticate (the pool alone has no app client)
- `passwordPolicy` for stronger password requirements
- `mfaConfiguration: OPTIONAL` with `softwareTokenMfaEnabled: true`
- `domain` for hosted UI login pages
