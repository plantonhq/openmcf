# Preset: Production Hardened

**Rank**: 3

## When to Use

- Production user directories holding real customer accounts
- Applications with security or compliance requirements
- Pools that must survive accidental deletion and credential-stuffing attacks

## What It Provides

- Strict password policy (12+ characters, all character classes, 5-password history, 3-day temporary passwords)
- Optional TOTP MFA users can enroll in
- Production email through a verified SES identity (`DEVELOPER` mode)
- Email updates verified before they take effect (no recovery lockout from a typo)
- PLUS feature tier with threat protection ENFORCED -- Cognito blocks or challenges risky sign-ins
- Deletion protection and a hosted-UI prefix domain

## What You Might Add

- `AwsCognitoUserPoolClient` resources for each application (web SPA, mobile, machine-to-machine)
- `logConfigurations` routing `userAuthEvents` to a CloudWatch log group for auth-event telemetry
- `mfaConfiguration: ON` to require MFA for every user
- `allowedFirstAuthFactors` with `WEB_AUTHN` (plus `webAuthn.relyingPartyId`) for passkey sign-in
- A custom domain with an ACM certificate in us-east-1
