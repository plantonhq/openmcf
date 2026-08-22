locals {
  runner_name = var.metadata.name

  # Resource-identity tags, matching the Pulumi module key-for-key.
  # DELIBERATELY five keys, no Name: the appliance's resources (cluster,
  # service, roles, secret, log group) each carry their own explicit resource
  # name -- a shared Name tag would mislabel ten distinct resources with one
  # value (the same recorded convention as AwsEcrRepo/AwsGlobalAccelerator).
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

  # The name the runner registers itself under when it joins the control
  # plane: "<env>-<metadata.name>" (metadata.name outside an environment)
  # -- the SAME derivation the platform uses for records that reference
  # this runner (its minted token, its managed destroy); changing this
  # formula breaks arrival attribution and managed teardown.
  registration_name = try(var.metadata.env, "") != "" ? "${var.metadata.env}-${var.metadata.name}" : var.metadata.name

  # The runtime identity the runner holds while executing work: the
  # referenced task_role when supplied (the module never mutates a
  # resource it merely references), else the permissionless role created
  # below. one() over the splat stays null-safe when count is 0.
  task_role_arn = var.spec.task_role != "" ? var.spec.task_role : one(aws_iam_role.runtime[*].arn)

  # The service's security groups: always the created outbound-only group,
  # plus any referenced extras (groups a private target trusts).
  security_group_ids = concat([aws_security_group.runner.id], var.spec.security_groups)
}
