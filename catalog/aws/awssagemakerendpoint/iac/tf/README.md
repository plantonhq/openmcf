# AwsSagemakerEndpoint — Terraform/OpenTofu module

Deploys an Amazon SageMaker AI endpoint (`aws_sagemaker_endpoint`)
with its folded endpoint configuration
(`aws_sagemaker_endpoint_configuration`).

Module facts worth knowing before editing:

- **The configuration is immutable upstream** (every argument
  ForceNew; the provider's update is tags-only) while the endpoint's
  pointer to it updates in place (UpdateEndpoint, optionally shaped by
  `deployment_config`). The fold therefore rolls configurations:
  `name_prefix` + `create_before_destroy` mint a NEW suffixed
  configuration on any capacity change, UpdateEndpoint repoints, and
  only then is the old configuration destroyed — the endpoint never
  references a deleted configuration (AWS's own documented pattern).
- **Variant names default deterministically per position** (locals:
  `variant-0`, `shadow-variant-0`, …) so a re-plan never regenerates
  them — the provider would mint a random name per plan otherwise,
  forcing a config roll every apply. The Pulumi module derives the
  identical defaults.
- **The capacity-reservation preference has ONE legal value**
  (`capacity-reservations-only`) — the module owns the constant and
  sends it exactly when an ML reservation ARN is configured.
- **Exactly-one rules live in the spec** (serverless XOR instance
  settings, one deployment strategy, shadow one-each-side) — the
  module renders what validation already admitted.

Outputs mirror the Pulumi module key-for-key: `endpoint_name`,
`endpoint_arn`, `endpoint_config_name`, `endpoint_config_arn`.
