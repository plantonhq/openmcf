locals {
  runner_name = var.metadata.name

  # Resource-identity tags, matching the Pulumi module key-for-key.
  aws_tags = {
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsPlantonRunner"
    "planton.ai/resource-id"   = var.metadata.id
  }

  # The group name is decided here (not read back from the resource) so the
  # container-definitions JSON can embed it as a plain string; the task
  # definition depends on the group explicitly so a task never starts
  # before its log destination exists.
  log_group_name = "/ecs/${local.runner_name}"

  # The tunnel exists for the real-time CloudOps channel: dual/grpc modes
  # join it (outbound-initiated, using the tunnel material in the
  # credentials document); a temporal-only worker polls its queue and
  # needs no tunnel at all.
  tunnel_enabled = var.spec.execution_mode == "temporal" ? "false" : "true"

  # The runtime identity the runner holds while executing work: the
  # referenced task_role when supplied (the module never mutates a
  # resource it merely references), else the permissionless role created
  # below. one() over the splat stays null-safe when count is 0.
  task_role_arn = var.spec.task_role != "" ? var.spec.task_role : one(aws_iam_role.runtime[*].arn)

  # The service's security groups: always the created outbound-only group,
  # plus any referenced extras (groups a private target trusts).
  security_group_ids = concat([aws_security_group.runner.id], var.spec.security_groups)
}
