# AwsCloudwatchLogAccountPolicy

One CloudWatch Logs account-level policy: a per-(name, type) policy object AWS applies account-wide in the region — data-protection masking, account-wide subscription forwarding, field indexing, ingest-time transformation, or metric extraction — optionally narrowed by selection criteria.

## Highlights

- **Identity is the (policy_name, policy_type) pair** — one AWS object per pair; two instances sharing both would fight over it. The provider imports the pair as `policy_name:policy_type`.
- **All five 2026 policy types modeled**: DATA_PROTECTION_POLICY, SUBSCRIPTION_FILTER_POLICY, FIELD_INDEX_POLICY, TRANSFORMER_POLICY, METRIC_EXTRACTION_POLICY. The document is structured configuration in the type's own schema, validated server-side at Put time.
- **`selection_criteria` narrows the account-wide scope** (e.g. a log-group-name prefix) and replaces on change; the provider's `scope` argument is module-pinned to its only legal value (ALL).
- **Untaggable at AWS** — the deliberate absence from the catalog's tag convention.

## Both Engines

Both modules render the single resource identically and export the same outputs: `policy_name` and `policy_type` (together, the import ID).
