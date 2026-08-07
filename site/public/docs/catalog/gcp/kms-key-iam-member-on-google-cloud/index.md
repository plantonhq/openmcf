---
title: "KMS Key IAM Member on Google Cloud"
description: "KMS Key IAM Member on Google Cloud deployment documentation"
icon: "package"
order: 100
componentName: "gcpkmskeyiammember"
---

# KMS Key IAM Member on Google Cloud

Grants one role, to one identity, on ONE Cloud KMS crypto key — the least-privilege unit of CMEK access control. Creating a key grants nobody anything: every consumer (a GCP service agent encrypting a bucket, a workload signing artifacts) needs an explicit role on the key, and this resource is exactly one such grant. Additive: it merges into the key's IAM policy without touching any other member's bindings, and removal subtracts only this exact (role, member) pair. Key-scoped beats ring- or project-scoped — a ring-level grant hands the member every key in the ring. Integrates with Planton's Provider Connections and composes through ValueFromRef with GcpKmsKey (`key_id` output), GcpServiceAccount (`member` output), and GcpIamCustomRole (`name` output).

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Crypto Key IAM Member Binding** -- a `kms.CryptoKeyIAMMember` merging the (role, member) pair into the target key's IAM policy, with an optional IAM Condition attached

## Before You Deploy

### Planton Setup

- **GCP Provider Connection** -- an active connection in the Connect module with credentials permitted to set IAM policy on the target key (e.g. `roles/cloudkms.admin` on the key or its ring). Map it as the default for your environment, or specify it explicitly.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline credentials or browser OAuth authentication modes.

### GCP Project

- **A KMS crypto key** whose IAM policy receives the grant. Provide the fully qualified key path directly or reference a GcpKmsKey Cloud Resource via ValueFromRef.
- **The identity** receiving the grant must already exist (deleted principals are not grantable).

## Deploy

### Console

Open the deployment store, find **KMS Key IAM Member on Google Cloud**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and the grant definition. Start from the **Storage CMEK Grant** preset in the [Presets](#presets) tab for the most common shape.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: gcp.planton.dev/v1
kind: GcpKmsKeyIamMember
metadata:
  name: gcs-agent-cmek-grant
  org: acme-corp
  env: prod
spec:
  cryptoKeyId:
    value: projects/acme-prod-12345/locations/us-central1/keyRings/prod-encryption/cryptoKeys/cmek-encrypt-key
  role:
    value: roles/cloudkms.cryptoKeyEncrypterDecrypter
  member:
    value: serviceAccount:service-123456789@gs-project-accounts.iam.gserviceaccount.com
```

```shell
planton apply -f gcp-kms-key-iam-member.yaml
```

This merges one binding into the key's policy. A Stack Job tracks the provisioning in real time.

### InfraChart

The composed form is where this kind shines — the full CMEK chain wired in one InfraPipeline:

```yaml
spec:
  cryptoKeyId:
    valueFrom:
      kind: GcpKmsKey
      name: cmek-encrypt-key
      fieldPath: status.outputs.key_id
  member:
    valueFrom:
      kind: GcpServiceAccount
      name: orders-api-worker
      fieldPath: status.outputs.member
```

The InfraPipeline resolves the dependency graph, deploys the key (and the account) first, then provisions the grant with the resolved values.

## Key Configuration

These are the most important decisions when configuring a grant. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**The CMEK pattern** -- to encrypt a GCP service's resources with your own key, grant that service's SERVICE AGENT (e.g. GCS's `service-<project-number>@gs-project-accounts.iam.gserviceaccount.com`) the `roles/cloudkms.cryptoKeyEncrypterDecrypter` role on the key, then point the resource's KMS field at the key. Both halves are required — a missing grant is the classic CMEK deploy failure.

**Role** -- a predefined KMS role (`roles/cloudkms.cryptoKeyEncrypterDecrypter` for CMEK, `signerVerifier` for signing keys, the one-directional `cryptoKeyEncrypter`/`cryptoKeyDecrypter` split for pipelines that must never read what they write, `publicKeyViewer` for verification-only consumers) or a custom role's fully-qualified name.

**Member format** -- the prefix declares the identity type: `serviceAccount:<email>`, `user:<email>`, `group:<email>`, `domain:<domain>`, or `principal://`/`principalSet://` federation principals. Format validation happens at deploy time because values usually arrive through references.

**IAM Condition** -- an optional CEL expression scoping WHEN the grant applies (a break-glass expiry is the everyday case on keys). The condition is part of the grant's identity: the same role with and without a condition are two independent grants.

**Everything replaces atomically** -- an IAM grant has no update. Changing key, role, member, or condition replaces the grant — and for the brief moment between delete and create, the member's encrypt/decrypt calls fail. Schedule changes on keys serving live CMEK traffic accordingly.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **GcpKmsKey** | `cryptoKeyId` | `status.outputs.key_id` |
| **GcpIamCustomRole** (optional) | `role` | `status.outputs.name` |
| **GcpServiceAccount** (optional) | `member` | `status.outputs.member` |

### What This Component Provides

After provisioning, `status.outputs` contains the grant's post-resolution facts:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `crypto_key_id` | The key after reference resolution | Audit tooling, key-access reviews |
| `role` | The role after reference resolution | Audit tooling |
| `member` | The member after reference resolution | Audit tooling |
| `etag` | The key IAM policy fingerprint when this grant last merged | Drift detection |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Storage CMEK grant** -- the GCS service agent gets encrypter/decrypter on the bucket's key; the grant every CMEK bucket needs. Start from the **Storage CMEK Grant** preset.

**Workload key user** -- an application's service account composed by reference for application-level encryption or signing. Start from the **Workload Key User** preset.

**Conditional key access** -- a human identity with a time-boxed condition; break-glass access that removes itself. Start from the **Conditional Key Access** preset.

## Works With

- [**GCP KMS Key**](/cloud-catalog/gcp-kms-key) -- its `key_id` output feeds the crypto key field; the key this grant controls access to
- [**GCP KMS Key Ring**](/cloud-catalog/gcp-kms-key-ring) -- the key's container; IAM granted at the ring flows down to every key (this kind exists to avoid that blast radius)
- [**GCP Service Account**](/cloud-catalog/gcp-service-account) -- its `member` output feeds the member field for workload key access
- [**GCP IAM Custom Role**](/cloud-catalog/gcp-iam-custom-role) -- its `name` output feeds the role field for curated permission bundles
