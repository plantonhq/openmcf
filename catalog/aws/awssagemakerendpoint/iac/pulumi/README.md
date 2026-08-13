# AwsSagemakerEndpoint — Pulumi module (Go)

Deploys an Amazon SageMaker AI endpoint (`sagemaker.Endpoint`) with
its folded endpoint configuration
(`sagemaker.EndpointConfiguration`).

Module facts worth knowing before editing:

- **The configuration is immutable upstream** (every argument
  ForceNew) while the endpoint's pointer to it updates in place
  (UpdateEndpoint, optionally shaped by `DeploymentConfig`). The fold
  therefore rolls configurations: `NamePrefix` + Pulumi's default
  create-before-delete replacement mint a NEW suffixed configuration
  on any capacity change, UpdateEndpoint repoints, and only then is
  the old configuration deleted — the endpoint never references a
  deleted configuration (AWS's own documented pattern).
- **Variant names default deterministically per position**
  (`variant-0`, `shadow-variant-0`, …) so re-previews never regenerate
  them — the Terraform module derives the identical defaults.
- **The capacity-reservation preference has ONE legal value**
  (`capacity-reservations-only`) — the module owns the constant and
  sends it exactly when an ML reservation ARN is configured.
- **Production and shadow variants share one spec message but two
  bridged types here** — parallel builders map the shared message onto
  each; keep them in lockstep.
- **Exactly-one rules live in the spec** (serverless XOR instance
  settings, one deployment strategy, shadow one-each-side) — the
  module renders what validation already admitted.

Outputs mirror the Terraform module key-for-key: `endpoint_name`,
`endpoint_arn`, `endpoint_config_name`, `endpoint_config_arn`.
