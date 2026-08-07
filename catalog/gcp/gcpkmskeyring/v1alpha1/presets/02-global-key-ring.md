# Global Key Ring

This preset creates a KMS key ring in the `global` location, making encryption keys accessible from any GCP region without latency penalties associated with cross-region access.

## When to Use

- Encryption keys shared across workloads in multiple GCP regions
- Global services that don't have data residency restrictions
- Shared signing keys for multi-region CI/CD pipelines
- Application-level encryption where key location is not regulated

## Key Configuration Choices

- **Global location** — keys are accessible from any GCP region. No data residency guarantees.
- **Simplicity** — no need to match key ring location with workload location.

## Values to Adjust

- `projectId.value` — the sample `my-gcp-project` must be replaced with
  your project ID (or switch to a `valueFrom` reference to a `GcpProject`
  resource).
- `keyRingName` — the permanent GCP name for this ring (1-63 chars,
  letters/digits/hyphens/underscores).

## Important

Key rings **cannot be deleted** from GCP. The name you choose is permanent within the project and location.

## Related Presets

- **01-regional-key-ring** — Key ring in a specific region (data residency compliance)
- **03-multi-region-key-ring** — Key ring replicated across a continent (availability + residency)
