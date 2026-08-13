---
title: "Secrets Manager Secret"
description: "AWS Secrets Manager secret deployment documentation"
icon: "package"
order: 100
componentName: "awssecretsmanagersecret"
---

# AWS Secrets Manager Secret

Deploys an AWS Secrets Manager secret — a named, versioned, KMS-encrypted container for credential material (database passwords, API keys, tokens, key/value JSON documents) with optional automatic rotation and cross-region replication. Applications retrieve the value at runtime through the Secrets Manager API instead of carrying credentials in configuration.

## What Gets Created

When you deploy an AwsSecretsManagerSecret resource, Planton provisions:

- **Secret** — the named container, encrypted under the AWS-managed `aws/secretsmanager` key or a customer-managed KMS key you reference
- **Secret Version** — created when a value is supplied (`stringValue` or `binaryValue`); staged `AWSCURRENT`, with optional custom staging labels alongside
- **Resource Policy** — created when `policy` is declared, with `blockPublicPolicy` guarding against policies that grant anonymous access (default on)
- **Rotation Configuration** — created when `rotation` is declared: a rotation Lambda you own, or a partner-managed external rotation

## How the Secret Value Stays Secret

`stringValue` and `binaryValue` are sensitive fields: the platform stores them as managed-secret references and resolves them just-in-time at deploy, so plaintext never lives in the control plane. Omitting both creates the secret shell with no value — useful when an application or rotation function writes the first version itself.

## Prerequisites

- **AWS credentials** configured via environment variables or Planton provider config
- **A KMS key** (optional) — required in practice for cross-account access; the AWS-managed key cannot be granted to other accounts
- **A rotation Lambda** (optional) — required for self-managed rotation; the function must grant Secrets Manager invoke permission (`principal: secretsmanager.amazonaws.com`)

## Quick Start

Create a file `secret.yaml`:

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsSecretsManagerSecret
metadata:
  name: prod/payments/db
  annotations:
    planton.dev/provisioner: pulumi
spec:
  region: us-west-2
  description: Payments database credentials
  stringValue: $secret/payments-db-credentials
```

Deploy it:

```bash
planton apply -f secret.yaml
```

## Deletion Behavior

Deletion is soft by default: AWS schedules the secret for deletion after `recoveryWindowInDays` (default 30) during which it can be restored — and during which the NAME stays reserved. Set `recoveryWindowInDays: 0` for immediate, unrecoverable deletion (the right choice for ephemeral or test secrets that must be recreatable immediately).

## Cross-Region Replication

Add entries to `replicaRegions` to keep read-only copies (same name, regional ARN) in sync in other regions — consumers there read locally with no cross-region call. Each replica encrypts under its own region's key. Removing an entry deletes that replica.

## Rotation

Two mechanisms, exactly one per secret:

- **`rotation.rotationLambdaArn`** — the classic self-managed path: Secrets Manager invokes your Lambda through the createSecret/setSecret/testSecret/finishSecret steps on the configured cadence.
- **`rotation.externalRotationRoleArn`** — partner-managed external rotation for managed external secrets (pairs with the spec's `type` partner identifier).

Cadence is `automaticallyAfterDays` (1–1000) or a `scheduleExpression` (`rate(...)` / `cron(...)`, UTC). By default AWS rotates once immediately when rotation is configured; set `rotateImmediately: false` to only test the configuration.

## Spec Reference

See [reference](v1alpha1/reference.md) for the complete field reference, and the [presets](presets/) for ready-to-deploy configurations.
