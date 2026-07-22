# AWS Pulumi State Backend

A self-managed Pulumi backend on S3, ready in one deploy: a versioned,
encrypted, public-access-blocked bucket that holds your stack checkpoints,
deployment history, and locks — with no Pulumi Cloud account required and no
per-seat pricing. Point `pulumi login` at it and every stack you run is stored
in infrastructure you own.

It is also the fastest way to use **Planton with your own Pulumi state
backend**: deploy this chart, register the bucket as a `StateBackend` in
Planton, and every Pulumi-provisioned deployment stores its state here.

## Architecture

```
                      ┌───────────────────────────────────┐
  pulumi up ─────────▶│  AwsS3Bucket (Pulumi backend)     │
  (checkpoints,       │   • versioning: Enabled           │
   history, locks)    │   • encryption: SSE-S3 (AES256)   │
                      │   • public access: blocked ×4     │
                      │   • lifecycle: version cleanup    │
                      └───────────────────────────────────┘
```

One resource — Pulumi locks natively inside the bucket, so no companion lock
table is needed (unlike a Terraform S3 backend).

## Included Cloud Resources

| Resource | Kind | Purpose |
|----------|------|---------|
| State bucket | `AwsS3Bucket` | Stores Pulumi checkpoints, history, and locks — versioned for recovery, encrypted at rest, all public access blocked, noncurrent versions cleaned up by lifecycle rules |

## Parameters

| Parameter | Description | Default | Type |
|-----------|-------------|---------|------|
| `aws_region` | Region for the bucket — pick the one closest to where deployments run | `us-east-1` | string |
| `bucket_name` | Globally unique S3 bucket name (`<org>-pulumi-state` convention) | `my-org-pulumi-state` | string |

## The one thing Pulumi DIY backends require: a secrets passphrase

Pulumi encrypts stack secrets **client-side** before they reach the backend.
On Pulumi Cloud that key is managed for you; on a self-managed backend YOU
choose the secrets provider — for most teams that is a passphrase:

```bash
export PULUMI_CONFIG_PASSPHRASE="<a strong passphrase from your secret manager>"
```

Treat the passphrase as a production credential: store it in a secret manager,
never rotate the underlying value out of band (existing stack secrets become
undecryptable), and prefer `pulumi stack change-secrets-provider` for
controlled rotation.

## After deploying

### Use it with Planton (your own state backend in minutes)

1. In the Planton console, create a **State Backend** with provisioner
   **Pulumi**, type **S3**, your bucket name and region — and set the
   **secrets passphrase** (Planton resolves it just-in-time for every run; it
   never lands in logs or state).
2. Mark it as the organization **default** for Pulumi. Every cloud resource
   created from then on pins this backend at creation and keeps it for life
   (rebinding an existing resource is a controlled migration, not an edit).

Credentials: with `inline` auth, provide an access key that can read/write the
bucket; with `runner` auth, the runner's own AWS identity (IRSA, instance
profile, or environment) is used and no keys are stored.

### Use it directly from the Pulumi CLI

```bash
pulumi login "s3://my-org-pulumi-state?region=us-east-1"
export PULUMI_CONFIG_PASSPHRASE="<your passphrase>"
pulumi stack init my-project/dev
pulumi up
```

Every stack managed after `pulumi login` lives under `.pulumi/` in the bucket;
one bucket serves any number of projects and stacks.

## Day-2 guidance

- **Recovering a checkpoint**: the bucket keeps the 30 most recent noncurrent
  versions of every object (expiring older ones after 90 days beyond those).
  Restore by copying a prior version of the stack's checkpoint object over the
  current one — or use `pulumi stack export`/`import` for surgical repairs.
- **KMS secrets provider instead of a passphrase**: `pulumi stack init
  --secrets-provider="awskms://<key-arn>"` moves secret encryption to a KMS
  key — no passphrase to manage; combine with an `AwsKmsKey` you own.
- **KMS bucket encryption**: for customer-managed at-rest encryption of the
  whole bucket, switch `encryption.sseAlgorithm` to `aws:kms` with `kmsKeyId`
  referencing an `AwsKmsKey`.
- **Tearing down**: the bucket refuses deletion while it holds objects unless
  `forceDestroy` is set — intentional for a bucket whose objects are your
  infrastructure's memory. Export any stacks you want to keep first.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
