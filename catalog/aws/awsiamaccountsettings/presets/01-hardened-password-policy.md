# Hardened Password Policy

This preset applies the CIS-benchmark-shaped password posture (14+
characters, all four character classes, 90-day rotation, 24-password
reuse memory) plus the v2 STS token version that works in every
region.

## When to Use

- Any account with IAM user console access — the security baseline
  auditors look for first
- Accounts using opt-in regions (the v2 token requirement)

## What You Get

- A password policy IAM enforces on every user's next password change
- Self-service password changes kept ON (users fix their own expiry)
- The global STS endpoint issuing v2 tokens, valid in all regions

## Customize

- The policy is replaced WHOLE on every apply — capture any existing
  posture in the spec before adopting an account, or omitted fields
  drop to AWS defaults
- `hardExpiry: true` requires an admin reset at expiry — a lockout
  policy, pair deliberately
- Add `accountAlias` for the friendly sign-in URL — knowing it
  REPLACES whatever alias the account already had
