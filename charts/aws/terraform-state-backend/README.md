# AWS Terraform State Backend

A production-hardened remote state backend for Terraform and OpenTofu, ready in
one deploy: a versioned, encrypted, public-access-blocked S3 bucket for state
files plus a DynamoDB lock table that serializes concurrent runs. This is the
first piece of infrastructure most teams need — you cannot collaborate on IaC
without a remote backend — and it is deliberately the fastest chart in the
catalog to go from zero to done.

It is also the fastest way to use **Planton with your own state backend**: deploy
this chart, register the bucket and table as a `StateBackend` in Planton, and
every subsequent deployment stores its state in infrastructure you own.

## Architecture

```
                      ┌─────────────────────────────────┐
 terraform / tofu ───▶│  AwsS3Bucket (state files)      │
   plan / apply       │   • versioning: Enabled         │
        │             │   • encryption: SSE-S3 (AES256) │
        │             │   • public access: blocked ×4   │
        │             │   • lifecycle: version cleanup  │
        │             └─────────────────────────────────┘
        │             ┌─────────────────────────────────┐
        └────────────▶│  AwsDynamodb (lock table)       │  ← lock_table_enabled
         lock LockID  │   • key: LockID (S, HASH)       │    (default: true)
                      │   • billing: on-demand          │
                      │   • deletion protection + PITR  │
                      └─────────────────────────────────┘
```

The two resources are independent (no reference between them); the DAG deploys
them in parallel.

## Included Cloud Resources

| Resource | Kind | Purpose |
|----------|------|---------|
| State bucket | `AwsS3Bucket` | Stores `.tfstate` objects — versioned for recovery, encrypted at rest, all public access blocked, noncurrent versions cleaned up by lifecycle rules |
| Lock table | `AwsDynamodb` | Serializes concurrent plans/applies via the S3 backend's `LockID` contract — on-demand billing, deletion-protected, point-in-time recovery (conditional) |

## Parameters

| Parameter | Description | Default | Type |
|-----------|-------------|---------|------|
| `aws_region` | Region for both resources — pick the one closest to where runs execute | `us-east-1` | string |
| `bucket_name` | Globally unique S3 bucket name (`<org>-terraform-state` convention) | `my-org-terraform-state` | string |
| `lock_table_name` | DynamoDB lock table name (per account+region) | `terraform-state-lock` | string |
| `lock_table_enabled` | Create the DynamoDB lock table (see "Locking" below) | `true` | bool |

## After deploying

### Use it with Planton (your own state backend in minutes)

Register the bucket as a `StateBackend` so your organization's deployments
store state in it:

1. In the Planton console, create a **State Backend** with provisioner
   **OpenTofu** (or Terraform), type **S3**, your bucket name and region, and —
   when the lock table is enabled — the DynamoDB table name for locking.
2. Mark it as the organization **default** for that provisioner. Every cloud
   resource created from then on pins this backend at creation and keeps it
   for life (existing resources keep the backend they were created with —
   rebinding is a controlled migration, not an edit).

Credentials: with `inline` auth, provide an access key that can read/write the
bucket and the table; with `runner` auth, the runner's own AWS identity (IRSA,
instance profile, or environment) is used and no keys are stored.

### Use it directly from Terraform / OpenTofu

```hcl
terraform {
  backend "s3" {
    bucket         = "my-org-terraform-state"
    key            = "my-project/terraform.tfstate"   # one key per state file
    region         = "us-east-1"
    dynamodb_table = "terraform-state-lock"           # when lock_table_enabled
    encrypt        = true
  }
}
```

Each project or environment gets its own `key`; the bucket serves any number of
state files.

## Locking: DynamoDB table vs native S3 lockfile

Terraform 1.10+ and OpenTofu 1.10+ can lock state natively in S3
(`use_lockfile = true` in the backend block), which removes the need for a
DynamoDB table. The toggle default keeps the table because:

- it works with **every** Terraform/OpenTofu version a team might still run,
- it is the locking shape Planton's S3 state-backend configuration integrates
  with today.

If your team is standardized on ≥ 1.10 and uses the bucket outside Planton, set
`lock_table_enabled: false` and add `use_lockfile = true` to your backend block
— the bucket needs no changes.

## Day-2 guidance

- **Recovering state**: the bucket keeps the 30 most recent noncurrent versions
  of every state object (and expires older ones after 90 days beyond those).
  Restore by copying a prior version over the current object, or with
  `terraform state push` after downloading it.
- **KMS instead of SSE-S3**: deploy an `AwsKmsKey` and switch the bucket's
  `encryption.sseAlgorithm` to `aws:kms` with `kmsKeyId` referencing the key —
  useful when compliance requires customer-managed rotation and key policy.
- **Tearing down**: the lock table has deletion protection enabled — disable it
  on the resource first, deliberately. The bucket refuses deletion while it
  holds objects unless `forceDestroy` is set; that is intentional for a bucket
  whose objects are your infrastructure's memory.

---

© Planton. Licensed under [Apache-2.0](../../../LICENSE).
