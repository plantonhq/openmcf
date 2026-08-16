# Secure Secret

This preset creates a KMS-encrypted SecureString under your own key,
sourced from a managed-secret reference so plaintext never lives in
the control plane or plan output.

## When to Use

- Database passwords, API keys, tokens — anything applications read
  from Parameter Store that must never appear in plans or state
- Accounts where secrets must be encrypted under a customer-managed
  key rather than the default aws/ssm key

## What You Get

- A SecureString whose value is resolved just-in-time at deploy (the
  spec REJECTS the plain arm for SecureString — this cannot regress
  silently)
- Encryption under the referenced KMS key; readers need both SSM and
  KMS access

## Customize

- Drop `keyId` to use the account's default aws/ssm key
- `tier: Advanced` unlocks values beyond 4KB (bills per parameter and
  cannot be downgraded in place)
- Rotation is the writer's job — every write publishes a new version
  and readers pick it up on the next fetch
