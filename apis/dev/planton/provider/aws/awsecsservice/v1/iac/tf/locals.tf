locals {
  service_name = var.metadata.name

  # Resource-identity tags, matching the Pulumi module key-for-key.
  aws_tags = {
    "Name"                     = local.service_name
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsEcsService"
    "planton.ai/resource-id"   = var.metadata.id
  }

  # Foreign-key fields are already flattened to primitive strings by the tofu
  # generator (the orchestrator resolves any value_from before the module
  # runs), so refs are consumed directly.

  # desired_count is a proto-optional tri-state: null seeds 1, explicit 0
  # deploys the wiring with nothing running.
  desired_count = coalesce(var.spec.desired_count, 1)

  # launch_type XOR capacity_provider_strategy is CEL-enforced in the spec;
  # with neither set, the module defaults to FARGATE explicitly rather than
  # inheriting the cluster's default -- the deployed result should never
  # depend on cluster-side mutable state.
  use_capacity_providers = length(var.spec.capacity_provider_strategy) > 0
  launch_type            = local.use_capacity_providers ? null : (var.spec.launch_type != "" ? var.spec.launch_type : "FARGATE")

  # The scalable target's resource id wants the cluster NAME; the spec
  # carries the cluster ARN (arn:aws:ecs:region:account:cluster/<name>).
  cluster_name = length(split("/", var.spec.cluster_arn)) > 1 ? element(split("/", var.spec.cluster_arn), length(split("/", var.spec.cluster_arn)) - 1) : var.spec.cluster_arn

  # The grace period only applies with load balancers (CEL-enforced); AWS
  # rejects it on a service without them.
  health_check_grace_period_seconds = length(var.spec.load_balancers) > 0 ? var.spec.health_check_grace_period_seconds : null

  # AWS defaults for target-tracking cooldowns, applied when the spec's
  # optional cooldown fields are unset.
  default_scale_in_cooldown  = 300
  default_scale_out_cooldown = 60

  autoscaling_enabled = var.spec.autoscaling != null
}
