locals {
  # Basic naming and tags
  resource_name = coalesce(try(var.metadata.name, null), "aws-ecs-service")
  tags          = merge({
    "Name" = local.resource_name
  }, try(var.metadata.labels, {}))

  # Safe cluster ARN handling
  safe_cluster_arn = coalesce(
    try(var.spec.cluster_arn.value, null),
    try(var.spec.cluster_arn.value_from.name, null)
  )

  # Extract cluster name from ARN
  cluster_name = length(split("/", local.safe_cluster_arn)) > 1 ? element(split("/", local.safe_cluster_arn), 1) : ""

  # Container settings
  container_port   = try(var.spec.container.port, null)
  desired_count    = coalesce(try(var.spec.container.replicas, null), 1)
  cpu              = var.spec.container.cpu
  memory           = var.spec.container.memory
  
  # Container image construction
  image_repo       = try(var.spec.container.image.repo, null)
  image_tag        = try(var.spec.container.image.tag, null)
  container_image  = local.image_repo != null && local.image_repo != "" ? (
    local.image_tag != null && local.image_tag != "" ? "${local.image_repo}:${local.image_tag}" : local.image_repo
  ) : null
  
  logging_enabled  = try(var.spec.container.logging.enabled, true)

  # Safe networking configuration
  safe_subnet_ids = [
    for s in try(var.spec.network.subnets, []) : 
    coalesce(try(s.value, null), try(s.value_from.name, null))
  ]
  
  safe_security_group_ids = [
    for sg in try(var.spec.network.security_groups, []) : 
    coalesce(try(sg.value, null), try(sg.value_from.name, null))
  ]

  # Safe IAM role handling
  safe_task_execution_role_arn = coalesce(
    try(var.spec.iam.task_execution_role_arn.value, null),
    try(var.spec.iam.task_execution_role_arn.value_from.name, null)
  )
  
  safe_task_role_arn = coalesce(
    try(var.spec.iam.task_role_arn.value, null),
    try(var.spec.iam.task_role_arn.value_from.name, null)
  )

  # ALB configuration
  alb_enabled          = try(var.spec.alb.enabled, false)
  safe_alb_arn         = coalesce(
    try(var.spec.alb.arn.value, null),
    try(var.spec.alb.arn.value_from.name, null)
  )
  alb_listener_port    = try(var.spec.alb.listener_port, 80)
  alb_listener_priority = try(var.spec.alb.listener_priority, 100)
  alb_routing_type     = lower(coalesce(try(var.spec.alb.routing_type, null), ""))
  alb_path             = try(var.spec.alb.path, null)
  alb_hostname         = try(var.spec.alb.hostname, null)

  # Boolean flags for conditional logic
  has_container_port = local.container_port != null
  has_alb_config     = local.alb_enabled && local.safe_alb_arn != null
  has_iam_roles      = local.safe_task_execution_role_arn != null || local.safe_task_role_arn != null

  # Health check grace period
  health_check_grace_period_seconds = coalesce(try(var.spec.health_check_grace_period_seconds, null), 60)

  # Auto scaling configuration
  autoscaling_enabled                 = try(var.spec.autoscaling.enabled, false)
  autoscaling_min_tasks               = try(var.spec.autoscaling.min_tasks, 1)
  autoscaling_max_tasks               = try(var.spec.autoscaling.max_tasks, 10)
  autoscaling_target_cpu_percent      = try(var.spec.autoscaling.target_cpu_percent, null)
  autoscaling_target_memory_percent   = try(var.spec.autoscaling.target_memory_percent, null)
}


