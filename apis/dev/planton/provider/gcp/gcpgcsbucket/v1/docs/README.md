# GcpGcsBucket — Deep Dive

## The problem this resource solves

Object storage is the connective tissue of a GCP environment: nearly every other service either reads from or writes to a bucket — Dataproc stages jobs in one, Cloud Functions deploy from one, Pub/Sub sinks events into one, load balancers serve static content out of one, BigQuery loads external data from one. This kind models the bucket as a first-class node so all of those relationships are explicit references to one reviewable resource — and so the dangerous decisions (public access, force-destroy, retention locks, soft-delete windows) are deliberate spec choices rather than tool defaults.

## Where it sits in the composition

Eight catalog kinds reference a bucket's `bucket_id` output today:

- **GcpBackendBucket** — the L7 load-balancer origin for static/CDN serving.
- **GcpCloudFunction** — the source-archive bucket builds deploy from.
- **GcpCloudRun / GcpCloudRunJob** — GCS volumes mounted into services and jobs.
- **GcpDataprocCluster** — staging and temp buckets.
- **GcpCloudComposerEnvironment** — the environment's DAG bucket.
- **GcpPubSubTopic / GcpPubSubSubscription** — Cloud Storage ingestion and delivery.

Inbound references the bucket itself makes: `projectId` → GcpProject, `kmsKeyName` → GcpKmsKey, `logging.logBucket` → another GcpGcsBucket, `iamMembers[].member` → GcpServiceAccount.

## Lifecycle contract

| Property | Behavior |
|---|---|
| `bucketName`, `location`, `projectId`, `customPlacementConfig`, `hierarchicalNamespaceEnabled`, `enableObjectRetention` | Immutable (ForceNew) — replacement deletes the bucket and its objects |
| `retentionPolicy.isLocked` | One-way: locking is irreversible; an unlock attempt forces re-creation, and a locked bucket cannot be deleted until every object passes retention |
| Everything else | Mutable in place |
| Deletion | Fails on a non-empty bucket unless `forceDestroy: true`, which deletes every object version first (slow on large buckets; refuses objects under a locked retention policy) |

## The safety posture (what defaults protect)

- **`forceDestroy` defaults to false.** A teardown of a data-bearing bucket should fail loudly, not erase quietly. Enable it only for ephemeral/derived data.
- **Soft delete is on by default server-side** (7 days): deleted objects remain recoverable — and billed. The spec makes the window explicit: 0 disables it for high-churn scratch buckets; up to 90 days lengthens it for precious data. An omitted block follows GCP's default without a perpetual diff.
- **`publicAccessPrevention: enforced`** makes public exposure impossible regardless of IAM — the right posture for every bucket not deliberately public.
- **Versioning + a bounded-history lifecycle rule** is the standard recovery pattern: overwrites keep the old version, the rule caps how many.

## Lifecycle rules and the explicit-zero contract

Numeric conditions (`ageDays`, `numNewerVersions`, `daysSinceNoncurrentTime`, `daysSinceCustomTime`) use explicit presence: unset means "criterion not used", while a set `0` is a real value (age 0 matches every object — useful with a prefix for TTL-on-write semantics). Both modules translate a set zero into the provider's send-zero flags so the realized rule is identical on either engine. `AbortIncompleteMultipartUpload` deserves a place in almost every bucket: abandoned multipart uploads hold invisible billed storage until aborted.

## Autoclass vs. hand-written transitions

Autoclass moves each object between storage classes based on its observed access pattern (down to `terminalStorageClass`, back to STANDARD on read). SetStorageClass lifecycle rules do the same by fixed schedule. They fight over the same property, so the spec rejects enabling both — pick Autoclass when access patterns are uncertain, explicit rules when they are fully predictable.

## Access model

Additive IAM only: each `iamMembers` entry grants one role to one member and composes safely with grants made anywhere else. IAM conditions scope grants to object prefixes or expiry dates. Public serving is the same mechanism (`allUsers` + `roles/storage.objectViewer`) — never a flag. Authoritative bindings/policies (which clobber grants they do not list) are deliberately not modeled.

Beneath IAM sits the network layer: `ipFilter` restricts which public CIDR ranges and which VPC networks (by reference to a `GcpVpcNetwork`) may reach the bucket at all, before any IAM evaluation. IAM decides WHO; the IP filter decides FROM WHERE. `allowAllServiceAgentAccess` keeps Google-managed integrations (transfer agents, log delivery) working when the filter is enabled.

## Recorded scope decisions

- **Bucket objects (`google_storage_bucket_object`)** — data plane, not infrastructure; deployment content flows through CI/tools, not IaC nodes.
- **Notifications, HMAC keys, managed/HNS folders, Anywhere Cache** — real sub-resources with independent lifecycles, deferred as Tier-2 kinds until composition demand exists.
- **Legacy ACL resources** (`bucket_acl`, `object_acl`, `default_object_acl`, `*_access_control`) — superseded by uniform bucket-level access; not modeled.
- **Encryption enforcement-config blocks** (GMEK/CMEK/CSEK restriction modes) — absent from the released 6.x provider line (schema-probe verified); revisit when the floor moves.
- **`lifecycle_rule` size conditions** (`size_above_bytes`/`size_below_bytes`) — absent from the released 6.x line (schema-probe verified).
- **`deletion_policy`** — absent from the released 6.x line; destroy semantics ride on `forceDestroy` identically on both engines.
