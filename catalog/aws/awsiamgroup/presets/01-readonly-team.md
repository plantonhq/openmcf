# Read-Only Team

This preset creates the simplest useful group: declared members with
AWS's managed ReadOnlyAccess policy — auditors, analysts, or anyone
who should see everything and touch nothing.

## When to Use

- Auditors, finance, support — read access across the account
- The first group in any account moving off per-user policies

## What You Get

- A group under `/teams/` whose membership is exactly the declared
  list (out-of-band additions are removed on the next apply)
- AWS's maintained ReadOnlyAccess policy — new services gain read
  permissions without you editing anything

## Customize

- Reference AwsIamUser outputs instead of literal names to order
  creation in charts:
  `valueFrom: {kind: AwsIamUser, name: <user>, fieldPath: status.outputs.user_name}`
- Swap in `arn:aws:iam::aws:policy/SecurityAudit` for a
  security-review shape
- Add an inline policy for the one extra grant unique to this group
  (see the developers-group preset)
