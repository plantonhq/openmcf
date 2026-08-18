# AwsEventBridgeApiDestination — Terraform/OpenTofu module

Manages an EventBridge API destination and/or its connection (`aws_cloudwatch_event_api_destination`, `aws_cloudwatch_event_connection`) as two independently deployable arms.

Module facts worth knowing before editing:

- **`authorization_type` is derived** from whichever auth block the spec sets (exactly one, by CEL) — the type and the credentials can never disagree.
- **Secrets never round-trip** — AWS stores credential values in a Secrets Manager secret it creates (`secret_arn` output) and no read API returns them; the provider reads them back from prior state, so imports cannot recover them (declared write-normalized in the import catalog).
- **The auth state machine gates creates/updates** — CREATING/AUTHORIZING → AUTHORIZED, up to 20 minutes; failures surface the connection's StateReason.
- **Neither resource is taggable at AWS** — the deliberate tag-convention absence (the AwsCloudwatchDashboard precedent).
- **The destination binds to exactly one connection** — the owned arm's ARN when present, else the spec's external `connection_arn` (CEL-enforced exactly-one).
- **Two private-endpoint homes exist deliberately** — `connectivity_parameters` inside auth reaches a private OAuth authorization endpoint; top-level `invocation_connectivity_parameters` invokes a private API. Both are VPC Lattice resource configurations.

Outputs mirror the Pulumi module key-for-key: `connection_arn`, `connection_secret_arn`, `api_destination_arn`.
