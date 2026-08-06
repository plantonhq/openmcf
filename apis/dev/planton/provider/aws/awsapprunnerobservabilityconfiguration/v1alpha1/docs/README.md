# AwsAppRunnerObservabilityConfiguration -- Research Notes

## Provider Surface

Modeled 1:1 on `aws_apprunner_observability_configuration`:

- `observability_configuration_name` (Required, ForceNew) -- derived from `metadata.name`, never a spec field.
- `trace_configuration` (Optional block) -- spec `trace_configuration` message.
  - `vendor` (Optional, enum `TracingVendor`) -- `AWSXRAY` is the only value the App Runner API defines today; the spec validates the closed set so a future vendor is an additive spec change, not a silent pass-through.
- Computed: `arn`, `observability_configuration_revision`, `latest`, `status`.

The trace settings are ForceNew -- AWS versions the resource exactly like the auto scaling configuration (same revision/latest semantics).

## Design Decisions

- **The service reference is the enable switch** -- the provider's service-side block requires a separate `observability_enabled` boolean alongside the ARN; the spec deliberately models only the reference and both engines set the boolean from its presence, making an enabled-without-configuration state unrepresentable.
- **`trace_configuration` kept as a message** (not flattened to one string) -- mirrors the provider block so a future vendor or additional trace knobs land additively in place.

## Deferral Ledger

- Nothing deferred: the resource has exactly one block and App Runner exposes no further observability surface today.

## Verification

- Spec tests cover the vendor closed set and the inert (no trace block) shape.
- E2E: register -> DescribeObservabilityConfiguration (active) -> destroy -> verify the revision flipped to inactive.
