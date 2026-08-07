# Object Storage Bucket on OCI

Deploys an Oracle Cloud Infrastructure Object Storage bucket -- a durable, scalable object store with configurable access controls, versioning, auto-tiering, retention rules, lifecycle management, and cross-region replication. Retention rules enforce minimum object lifetimes for compliance. Lifecycle rules automate tiering, archival, and deletion based on object age. Replication policies asynchronously copy objects to a destination bucket in another OCI region for disaster recovery. The component integrates with Planton's Provider Connections for OCI credential management and supports ValueFromRef wiring to compartments and encryption keys.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Object Storage Bucket** -- the bucket in the specified compartment and namespace with access type, storage tier, versioning, auto-tiering, event emission, optional KMS encryption, user metadata, and inline retention rules (up to 100)
- **Object Lifecycle Policy** -- created only when `lifecycleRules` is non-empty; a single policy resource containing all lifecycle rules that automate object archival, tiering transitions, deletion, and multipart upload cleanup
- **Replication Policies** -- one policy per entry in `replicationPolicies`; each asynchronously copies objects to a pre-existing destination bucket in another OCI region. All replication policy fields are immutable after creation.
- **Freeform Tags** -- resource metadata tags (organization, environment, resource kind, resource ID) applied to the bucket

## Before You Deploy

### Planton Setup

- **OCI Provider Connection** -- an active connection in the Connect module with credentials for the target OCI tenancy. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline credentials.

### OCI Tenancy

- A compartment to place the bucket in. Provide the compartment OCID directly or reference an OciCompartment Cloud Resource via ValueFromRef.
- The Object Storage namespace for the tenancy (a unique string like `axe1234abc`). Retrieve it via `oci os ns get` or from the OCI Console.
- For customer-managed encryption: an OCI KMS key. When omitted, Oracle-managed encryption is used.
- For cross-region replication: destination buckets must already exist in the target regions before creating replication policies.

## Deploy

### Console

Open the deployment store, find **Object Storage Bucket on OCI**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Private Versioned** preset in the [Presets](#presets) tab to pre-populate a production bucket with KMS encryption, versioning, auto-tiering, and lifecycle rules for old versions and incomplete uploads.

### CLI

```yaml
apiVersion: oci.planton.dev/v1
kind: OciObjectStorageBucket
metadata:
  name: app-data
  org: acme-corp
  env: prod
spec:
  compartmentId:
    value: "ocid1.compartment.oc1..example"
  namespace: "axe1234abc"
  name: "acme-prod-app-data"
  accessType: no_public_access
  storageTier: standard
  versioning: enabled
  autoTiering: infrequent_access
```

```shell
planton apply -f bucket.yaml
```

This creates a private bucket with Standard tier, versioning enabled, and auto-tiering to InfrequentAccess for cost optimization. Oracle-managed encryption is used. No retention rules, lifecycle policies, or replication are configured.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the bucket to a compartment and encryption key deployed in the same InfraPipeline:

```yaml
spec:
  compartmentId:
    valueFrom:
      kind: OciCompartment
      name: storage
      fieldPath: status.outputs.compartmentId
  kmsKeyId:
    valueFrom:
      kind: OciKmsKey
      name: storage-key
      fieldPath: status.outputs.keyId
```

The InfraPipeline resolves the dependency graph, deploys the compartment and KMS key first, then provisions the bucket with the resolved values.

## Key Configuration

These are the most important decisions when configuring an Object Storage bucket. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Storage tier** -- Set `storageTier` to `standard` for frequently accessed data or `archive` for cold storage with lower cost and higher retrieval latency. The storage tier is immutable after creation. Enable `autoTiering` with `infrequent_access` on Standard-tier buckets to automatically move infrequently accessed objects to a lower-cost tier.

**Versioning** -- Set `versioning` to `enabled` to protect objects from accidental overwrites and deletions by maintaining version history. Once enabled, versioning can be `suspended` (new versions stop being created) but cannot be reverted to `disabled`. Pair with lifecycle rules targeting `previous-object-versions` to archive or delete old versions after a retention period.

**Retention rules** -- Add entries to `retentionRules` to enforce minimum object lifetimes. Each rule has a `displayName`, `duration` (amount + unit), and optional `timeRuleLocked` datetime. Once a rule is locked (past the lock time), it can only be deleted by deleting the bucket, and only duration increases are permitted.

**Lifecycle rules** -- Add entries to `lifecycleRules` to automate object management. Actions include `lifecycle_archive` (move to Archive tier), `lifecycle_infrequent_access` (move to InfrequentAccess), `lifecycle_delete` (permanently remove), and `lifecycle_abort` (clean up incomplete multipart uploads). Each rule targets `objects`, `previous-object-versions`, or `multipart-uploads` with optional name filters.

**Cross-region replication** -- Add entries to `replicationPolicies` to asynchronously replicate objects to destination buckets in other OCI regions. Destination buckets must exist before policy creation. All replication policy fields are immutable after creation.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **OciCompartment** | `compartmentId` | `status.outputs.compartmentId` |
| **OciKmsKey** (optional) | `kmsKeyId` | `status.outputs.keyId` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `bucket_id` | OCID of the Object Storage bucket | IAM policy scoping, replication destination references, monitoring alarms, application configuration |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Private versioned** -- A production bucket with no public access, KMS encryption, versioning, auto-tiering, and lifecycle rules that archive old object versions after 90 days and abort incomplete multipart uploads after 7 days. Start from the **Private Versioned** preset.

**Archive storage** -- A compliance-oriented bucket using the Archive storage tier with a 7-year retention rule. Versioning is disabled (Archive tier stores immutable objects). A lifecycle rule cleans up incomplete uploads. Start from the **Archive Storage** preset.

**Public read** -- A bucket allowing individual object reads without directory listing (ObjectReadWithoutList). Event emission is enabled to track downloads via the OCI Events service. Suitable for public static assets. Start from the **Public Read** preset.

## Works With

- [**Compartment on OCI**](/cloud-catalog/oci-compartment) -- provides the compartment that scopes this bucket
- [**KMS Key on OCI**](/cloud-catalog/oci-kms-key) -- provides the customer-managed encryption key for server-side encryption