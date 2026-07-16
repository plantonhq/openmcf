# AWS Batch MANAGED compute environment.
#
# Only MANAGED environments are modeled: Batch owns the instance lifecycle.
# (UNMANAGED means bring-your-own ECS container instances -- a different
# operating model.) Job queues and scheduling policies are separate
# resources (AwsBatchJobQueue / AwsBatchSchedulingPolicy) that compose onto
# this environment through its exported ARN.
resource "aws_batch_compute_environment" "this" {
  # The cloud name comes from metadata.name (the catalog naming basis) --
  # set explicitly so both engines create the same environment name.
  name  = var.metadata.name
  type  = "MANAGED"
  state = var.spec.state

  # Leaving service_role unset lets AWS use (and auto-create) the Batch
  # service-linked role -- which is also what keeps the environment eligible
  # for in-place infrastructure updates.
  service_role = var.spec.service_role != "" ? var.spec.service_role : null

  compute_resources {
    type      = local.cr.type
    max_vcpus = local.cr.max_vcpus

    # Subnets are required for every type: even "serverless" Fargate tasks
    # get ENIs placed into these subnets.
    subnets            = local.cr.subnet_ids
    security_group_ids = length(local.cr.security_group_ids) > 0 ? local.cr.security_group_ids : null

    # EC2/SPOT-only knobs. Fargate environments must not send them at all
    # (AWS rejects the request), hence the null gating rather than defaults.
    # min_vcpus is platform-defaulted to 0 so coalesce never fires in
    # practice; it exists for direct tfvars runs outside the platform.
    min_vcpus           = local.is_ec2_family ? coalesce(local.cr.min_vcpus, 0) : null
    desired_vcpus       = local.is_ec2_family && local.cr.desired_vcpus > 0 ? local.cr.desired_vcpus : null
    instance_type       = local.is_ec2_family && length(local.cr.instance_types) > 0 ? local.cr.instance_types : null
    allocation_strategy = local.is_ec2_family && local.cr.allocation_strategy != "" ? local.cr.allocation_strategy : null
    instance_role       = local.is_ec2_family && local.cr.instance_role != "" ? local.cr.instance_role : null
    ec2_key_pair        = local.is_ec2_family && local.cr.ec2_key_pair != "" ? local.cr.ec2_key_pair : null
    placement_group     = local.is_ec2_family && local.cr.placement_group != "" ? local.cr.placement_group : null

    # These tags land on the EC2 instances / Spot requests Batch launches --
    # deliberately NOT merged with the environment's own identity tags
    # (local.aws_tags), which tag the CE resource itself.
    tags = local.is_ec2_family && length(local.cr.resource_tags) > 0 ? local.cr.resource_tags : null

    # SPOT-only knobs.
    bid_percentage      = local.is_spot ? local.cr.bid_percentage : null
    spot_iam_fleet_role = local.is_spot && local.cr.spot_iam_fleet_role != "" ? local.cr.spot_iam_fleet_role : null

    dynamic "launch_template" {
      for_each = local.cr.launch_template != null ? [local.cr.launch_template] : []
      content {
        launch_template_id = launch_template.value.launch_template_id
        version            = launch_template.value.version != "" ? launch_template.value.version : null
      }
    }

    dynamic "ec2_configuration" {
      for_each = local.cr.ec2_configurations
      content {
        image_type               = ec2_configuration.value.image_type != "" ? ec2_configuration.value.image_type : null
        image_id_override        = ec2_configuration.value.image_id_override != "" ? ec2_configuration.value.image_id_override : null
        image_kubernetes_version = ec2_configuration.value.image_kubernetes_version != "" ? ec2_configuration.value.image_kubernetes_version : null
      }
    }
  }

  # Batch-on-EKS attachment: create-time only (the provider replaces the
  # environment on any change here).
  dynamic "eks_configuration" {
    for_each = var.spec.eks_configuration != null ? [var.spec.eks_configuration] : []
    content {
      eks_cluster_arn      = eks_configuration.value.eks_cluster_arn
      kubernetes_namespace = eks_configuration.value.kubernetes_namespace
    }
  }

  dynamic "update_policy" {
    for_each = var.spec.update_policy != null ? [var.spec.update_policy] : []
    content {
      terminate_jobs_on_update = update_policy.value.terminate_jobs_on_update
      # Tri-state: unset lets AWS apply its own 30-minute default (matching
      # the Pulumi module, which also omits it when unset).
      job_execution_timeout_minutes = update_policy.value.job_execution_timeout_minutes
    }
  }

  tags = local.aws_tags
}
