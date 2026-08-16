# Cross-Account Backup

This preset enables cross-account backup copies for the organization —
the global switch behind copy actions that target vaults in other
accounts.

## When to Use

- Organizations copying recovery points to a dedicated backup account
  (the strongest blast-radius isolation short of air-gapped vaults)
- Before any backup plan copy action targets a cross-account vault ARN

## What You Get

- `isCrossAccountBackupEnabled` on, account-wide (the global arm's
  identity is the ACCOUNT — deploy this in exactly one instance across
  all regions)

## Customize

- Combine with the region arm in one instance when a single region
  owns both levers
- This is an organizations-level control: flipping it affects every
  account's backup posture — treat it as change-managed
- Destroy reverts NOTHING (a no-op at AWS) — to disable, apply
  `"false"` first
