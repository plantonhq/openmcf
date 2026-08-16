# AWS Backup Vault

The encrypted container your backups live in — a standard vault for
everyday recovery points, or a logically air-gapped vault whose
contents nobody (not even root) can delete before retention expires.

## What Gets Managed

- Exactly one vault type: standard (optional KMS key, drain-at-destroy
  switch) or logically air-gapped (locked-in retention bounds,
  immutable).
- On standard vaults: Vault Lock (governance or compliance mode), the
  vault access policy, and event notifications to SNS.

## Before You Deploy

### Planton Setup

- An organization and environment.
- An AWS provider connection with AWS Backup permissions.

### AWS Account

- Nothing else for the vault itself. For custom encryption, a KMS key
  ([AWS KMS Key](/cloud-catalog/aws-kms-key)); for notifications, an
  SNS topic whose policy lets `backup.amazonaws.com` publish
  ([AWS SNS Topic](/cloud-catalog/aws-sns-topic)).

## Deploy

### Console

Create the resource from the AWS catalog, pick the vault type, and
deploy. Compliance-mode locks are IRREVERSIBLE after their cooling-off
window — read the guide before setting `changeable_for_days`.

### CLI

```bash
planton apply -f backup-vault.yaml
```

## After Deploy

- Point backup plan rules at the vault by name
  ([AWS Backup Plan](/cloud-catalog/aws-backup-plan)).
- To destroy a standard vault holding recovery points, set
  `force_destroy: true` and apply first — AWS refuses to delete
  non-empty vaults.
- An air-gapped vault destroys only after its recovery points age out
  by retention.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
