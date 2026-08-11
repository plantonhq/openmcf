locals {
  # Resource-identity tags follow the catalog convention. The Name tag is
  # what the Route 53 console displays as the health check's name.
  aws_tags = {
    "Name"                     = var.metadata.name
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsRoute53HealthCheck"
    "planton.ai/resource-id"   = var.metadata.id
  }

  # The check type gates which argument families apply (CEL-enforced in the
  # spec): endpoint probing, calculated aggregation, CloudWatch mirroring, or
  # recovery-control mirroring.
  is_endpoint_check = contains(["HTTP", "HTTPS", "HTTP_STR_MATCH", "HTTPS_STR_MATCH", "TCP"], var.spec.check_type)

  # Null-when-unset scalars so the provider applies its own defaults (port 80
  # for HTTP / 443 for HTTPS, resource path "/", all checker regions) instead
  # of this module freezing them.
  fqdn          = var.spec.fqdn != "" ? var.spec.fqdn : null
  ip_address    = var.spec.ip_address != "" ? var.spec.ip_address : null
  port          = var.spec.port != 0 ? var.spec.port : null
  resource_path = var.spec.resource_path != "" ? var.spec.resource_path : null
  search_string = var.spec.search_string != "" ? var.spec.search_string : null
  regions       = length(var.spec.regions) > 0 ? var.spec.regions : null

  # Probe tuning only exists where probing happens.
  request_interval  = local.is_endpoint_check ? var.spec.request_interval : null
  failure_threshold = local.is_endpoint_check ? var.spec.failure_threshold : null
  measure_latency   = local.is_endpoint_check ? var.spec.measure_latency : null
  enable_sni        = local.is_endpoint_check ? var.spec.enable_sni : null

  # CALCULATED aggregation — the generator flattens the child refs to plain
  # strings. The threshold passes through by PRESENCE: an explicit 0 means
  # "always healthy" per AWS's contract (a different configuration from
  # omitting it, which lets AWS apply its server-side default).
  child_healthchecks     = length(var.spec.child_health_checks) > 0 ? var.spec.child_health_checks : null
  child_health_threshold = var.spec.child_health_threshold

  # CLOUDWATCH_METRIC mirroring.
  cloudwatch_alarm_name           = var.spec.cloudwatch_alarm_name != "" ? var.spec.cloudwatch_alarm_name : null
  cloudwatch_alarm_region         = var.spec.cloudwatch_alarm_region != "" ? var.spec.cloudwatch_alarm_region : null
  insufficient_data_health_status = var.spec.insufficient_data_health_status != "" ? var.spec.insufficient_data_health_status : null

  # RECOVERY_CONTROL mirroring.
  routing_control_arn = var.spec.routing_control_arn != "" ? var.spec.routing_control_arn : null
}
