# Standard Vault

This preset creates the everyday backup vault: encrypted with the AWS
Backup service key, no lock, ready for backup plan rules to target by
name.

## When to Use

- The default container for a backup plan's recovery points
- One vault per application or environment (empty vaults are free)

## What You Get

- A standard backup vault named `metadata.name`
- Service-key encryption (no KMS key to manage)
- Deletable any time while empty — AWS refuses to delete vaults
  holding recovery points

## Customize

- Add `kmsKeyArn` for customer-managed encryption (fixed at creation)
- Add `lock` for write-once protection — omit `changeableForDays` for
  removable governance mode; setting it opts into IRREVERSIBLE
  compliance mode
- Add `notifications` (an SNS topic whose policy allows
  `backup.amazonaws.com`) to alert on failed jobs
- Set `forceDestroy: true` only on disposable vaults — it deletes ALL
  recovery points at destroy
