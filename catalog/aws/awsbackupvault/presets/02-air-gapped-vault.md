# Air-Gapped Vault

This preset creates a logically air-gapped vault: recovery points that
nobody — compromised credentials included — can delete or re-encrypt
before retention expires. The ransomware-recovery tier.

## When to Use

- A last-resort copy of critical backups, isolated from the account's
  everyday credentials
- Compliance regimes requiring immutable, non-deletable backups

## What You Get

- An air-gapped vault enforcing 7–35 day retention on every recovery
  point it holds
- AWS-owned-key encryption and immutability by construction — every
  field forces replacement

## Customize

- Raise `maxRetentionDays` for longer immutable retention — but note
  the vault cannot be destroyed until its recovery points age out
- Fill the vault from a backup plan rule's
  `targetLogicallyAirGappedBackupVaultArn` (the rule still needs a
  standard target vault)
- Add `encryptionKeyArn` to substitute a customer-managed key (fixed
  at creation)
