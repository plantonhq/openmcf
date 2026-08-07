# GcpKmsKeyRing — Deep Dive

## The problem this resource solves

Every customer-managed encryption story on GCP starts with a container decision: where do the keys live, and who can touch them? The key ring is that container — a permanent, location-anchored grouping whose IAM flows down to every key inside it. Getting the ring layout right (one per environment or data domain, in the location where the protected data lives) is what keeps CMEK manageable at fleet scale; getting it wrong cannot be quietly fixed later, because rings can never be deleted or renamed.

## Where it sits in the composition

- **GcpKmsKeyRing** — this resource: the permanent container, owning project + location.
- **GcpKmsKey** — the keys inside the ring; each references the ring's `status.outputs.key_ring_id` (the fully qualified `projects/{p}/locations/{l}/keyRings/{name}` path) and inherits its project and location.
- **CMEK consumers** — BigQuery datasets/tables, Spanner, Cloud SQL, AlloyDB, GKE, Cloud Run/Functions, Pub/Sub, Vertex AI, Filestore, Bigtable, Firestore, Memorystore, Dataproc, Composer — reference the *keys*, not the ring.
- **Bare-name consumers** — components that take a ring name plus a separately supplied project/location (for example OpenBao's GCP KMS seal) compose from `key_ring_name` + `location`.

## Lifecycle contract

| Property | Behavior |
|---|---|
| `keyRingName`, `location`, `projectId` | Immutable (ForceNew) — any change abandons the old ring and creates a new one |
| Deletion | **No delete API exists.** Destroy removes the ring from IaC state only; the ring remains in GCP permanently, at no cost |
| Name reuse | Never possible within a project+location — the original ring still occupies the name |

The permanence is by design: a deleted ring would strand every key inside it, and destroyed key material is unrecoverable. GCP chose to make the container immortal instead.

## Design guidance

- **One ring per environment or data domain.** IAM granted on the ring reaches every key inside it, so the ring is the unit of access review. A `prod-data` ring and a `staging-data` ring keep those worlds separable; per-key rings multiply permanent objects for no isolation gain (key-level IAM exists for the exceptions).
- **Location follows data.** Regional CMEK integrations require the key (hence the ring) in the data's region; multi-region BigQuery/Spanner require the matching multi-region (`us`, `europe`, `asia`); `global` suits signing keys consumed from everywhere.
- **Name for forever.** `prod-encryption`, `us-compliance-keys` — names that stay correct as the fleet grows, because they cannot be recycled.

## 90/10 coverage vs the provider resource

| Provider field (`google_kms_key_ring`) | Modeled | Notes |
|---|---|---|
| `name` | ✅ `keyRingName` | ForceNew |
| `location` | ✅ `location` | ForceNew |
| `project` | ✅ `projectId` (ref → GcpProject, ambient fallback) | ForceNew |

The released provider resource has exactly these three fields — this kind models 100% of its surface.

## Deliberately not modeled (recorded reasons)

- **Per-ring IAM (`google_kms_key_ring_iam_*`)** — resource-scoped IAM stays out of the catalog pending concrete pull (the additive project-level grant, `GcpProjectIamMember`, covers the real cases).
- **`google_kms_key_ring_import_job`** — the BYOK import ceremony handles raw key material interactively (wrapping keys expire in 3 days); it is an operational act, not durable infrastructure. The import-only *key* container is modeled on `GcpKmsKey`.

## Provider mapping

Maps to `google_kms_key_ring` (`google/services/kms/resource_kms_key_ring.go`): `keyRingName` → `name` (ForceNew), `location` → `location` (ForceNew), `projectId` → `project` (ForceNew, ambient fallback when empty). The provider's Delete is a state-only removal with a warning — mirrored in both engines' behavior and taught in every doc surface of this kind.
