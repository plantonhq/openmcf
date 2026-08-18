# AwsEventBridgeApiDestination

An EventBridge API destination with its connection: the authenticated HTTP(S) endpoint that rules, pipes, and schedules invoke. Two independently deployable arms — the connection (the shareable auth trust anchor) and the destination (endpoint + method + rate limit), bound own-XOR-existing.

## Highlights

- **The auth type can never lie**: the spec models api-key XOR basic XOR OAuth as arms and the modules DERIVE `authorization_type` from whichever is set — no separate field to contradict the credentials.
- **Secrets are secret-typed and never round-trip**: credential fields carry the platform's sensitive marking; AWS stores the values in a Secrets Manager secret it creates (the `connection_secret_arn` output) and no read API returns them — declared write-normalized in the import catalog.
- **One connection, many destinations**: a shared connection lives in one owning instance; other instances' destinations reference it by ARN (`connection_arn`, chart-wired to this kind's own output).
- **Private endpoints modeled honestly**: both VPC Lattice homes — the private OAuth authorization endpoint (inside auth) and the private invocation endpoint (top-level) — as intention-named fields.

## Both Engines

Both modules render the two arms identically (the connection's 20-minute auth state machine gates creates) and export the same outputs: `connection_arn`, `connection_secret_arn`, `api_destination_arn`.

## Chart Wiring

`destination.connection_arn` → another AwsEventBridgeApiDestination's `connection_arn` (the shared-connection topology); `kms_key_identifier` → AwsKmsKey `key_arn`. The `api_destination_arn` output is what AwsEventBridgeRule targets, pipes, and schedules invoke.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
