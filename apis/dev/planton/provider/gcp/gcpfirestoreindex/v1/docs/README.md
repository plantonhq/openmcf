# GCP Firestore Index — Deep Dive

## The problem this resource solves

Firestore rejects queries that filter or order on multiple fields unless a matching composite index exists. Console error-message links create indexes ad hoc — invisible to review, unreproducible across environments. Declaring indexes as infrastructure makes a deployment's query capabilities explicit and version-controlled.

## Field roles and ordering

Each field plays exactly one role:

- `order` — scalar sort or range comparisons (`ASCENDING` / `DESCENDING`).
- `arrayConfig: CONTAINS` — array-membership queries.
- `vectorConfig` — nearest-neighbor (embedding) search; must be last.

List fields in query order: equality filters first, then inequality/sort fields, then vector. Firestore appends `__name__` automatically when needed.

## Immutability and replacement

Every index property is immutable in the provider (ForceNew). Changing anything replaces the index. Firestore rebuilds the new index in the background while the old one continues serving queries until the new one is ready — the operational reason create-before-destroy matters when applied manually, though Planton manages the replacement lifecycle.

## Scope and density

- `queryScope` — `COLLECTION` (default), `COLLECTION_GROUP`, or `COLLECTION_RECURSIVE` (Datastore Mode).
- `apiScope` — `ANY_API` (Native, default) or `DATASTORE_MODE_API`.
- `density` — leave empty for GCP's default (`SPARSE_ALL`); set `DENSE` for Datastore Mode query shapes that require it.

## Vector indexes

Vector fields require a `dimension` matching the embedding model output and a flat index marker in the provider (`flat {}` block in Terraform, `Flat: &IndexFieldVectorConfigFlatArgs{}` in Pulumi). Vector fields must come after any filter/sort fields.

## No labels surface

Firestore indexes do not support GCP labels — both engines skip label merge identically.

## IAM and API prerequisites

- `firestore.googleapis.com` enabled (both modules enable it with `disable_on_destroy=false`).
- `roles/datastore.indexAdmin` or broader Firestore admin on the project.

## Deliberately not modeled

- **`search_config` (text/geo search indexes)** — newer provider surface; Tier-2 candidate on concrete pull.
- **`multikey` / `unique`** — MongoDB-compatible API scope only; excluded from v1 spec.
- **`deletion_policy`** — client-side Terraform lever conflicting with Planton-managed destroy (catalog-wide decision).
