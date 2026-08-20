# AWS Secrets Manager Secret

Deploys an AWS Secrets Manager secret -- a named, versioned, KMS-encrypted container for credential material with optional automatic rotation and cross-region replication. The secret value is supplied as a managed-secret reference and resolved just-in-time at deploy, so plaintext never lives in the control plane. It integrates with Planton's Provider Connections for AWS credential management and supports ValueFromRef wiring to KMS keys, rotation Lambda functions, and IAM roles.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Secret** -- the named, KMS-encrypted container; hierarchical names like `prod/payments/db` are legal (up to 512 characters of alphanumeric plus `/_+=.@-`)
- **Secret Version** -- created when a value arm (`stringValue` or `binaryValue`) is set, staged `AWSCURRENT` with optional custom staging labels riding alongside
- **Resource Policy** -- created when `policy` is declared, rendered through the standalone policy resource so `blockPublicPolicy` (default on) rejects policies that grant anonymous access
- **Rotation Configuration** -- created when `rotation` is declared: exactly one of a self-managed rotation Lambda or a partner-managed external rotation role
- **AWS Tags** -- resource metadata tags (organization, environment, resource kind, resource ID) applied automatically for tracking and governance

## Before You Deploy

### Planton Setup

- **AWS Provider Connection** -- an active connection in the Connect module with credentials for the target AWS account. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **A managed secret** holding the value -- create it in Planton's secrets store and reference it as `$secret/<slug>` in `stringValue`; the platform rejects plaintext on sensitive fields.

### AWS Account

- **A KMS key** (optional) for customer-managed encryption -- required in practice when other AWS accounts must read the secret (the AWS-managed `aws/secretsmanager` key cannot be granted cross-account). Provide the ARN directly or reference an AwsKmsKey Cloud Resource.
- **A rotation Lambda** (optional) for self-managed rotation -- the function must implement the four rotation steps and grant Secrets Manager invoke permission (`principal: secretsmanager.amazonaws.com`). Reference an AwsLambda Cloud Resource or provide the ARN.

## Deploy

### Console

Open the deployment store, find **AWS Secrets Manager Secret**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Application Credentials** preset in the [Presets](#presets) tab to pre-populate a working configuration.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsSecretsManagerSecret
metadata:
  name: prod/payments/db
spec:
  region: us-west-2
  description: Payments database credentials
  stringValue: $secret/payments-db-credentials
```

```bash
planton apply -f secret.yaml
```

## Operational Notes

- **Deletion is soft by default**: the secret is recoverable for `recoveryWindowInDays` (default 30), and its NAME stays reserved for the window. Use `0` for immediate permanent deletion of ephemeral secrets.
- **The secret ARN carries a random suffix** (`arn:...:secret:name-AbCdEf`) -- wire consumers through the `secret_arn` output, never by deriving the ARN from the name.
- **Replicas are read-only regional copies** kept in sync by AWS; removing a `replicaRegions` entry deletes that replica. `forceOverwriteReplicaSecret` controls collision behavior with pre-existing same-named secrets in the replica region.
- **Rotation invokes immediately by default** (`rotateImmediately: true`) -- the secret must hold a value the rotation function can read.
