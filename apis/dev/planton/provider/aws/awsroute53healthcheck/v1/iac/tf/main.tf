resource "aws_route53_health_check" "this" {
  # The monitoring model (ForceNew). The spec's CEL rules guarantee the
  # arguments below already match the chosen type, so this resource maps
  # fields 1:1 without re-validating.
  type = var.spec.check_type

  # --- Endpoint checks (HTTP / HTTPS / *_STR_MATCH / TCP) -------------------
  # fqdn alone: Route 53 resolves and probes it (and sends it as the Host
  # header). ip_address alone or with fqdn: the probe goes to the IP.
  fqdn          = local.fqdn
  ip_address    = local.ip_address
  port          = local.port
  resource_path = local.resource_path
  search_string = local.search_string

  # request_interval and measure_latency are create-time (ForceNew).
  request_interval  = local.request_interval
  failure_threshold = local.failure_threshold
  measure_latency   = local.measure_latency
  enable_sni        = local.enable_sni
  regions           = local.regions

  # --- State shaping (any type) ----------------------------------------------
  invert_healthcheck = var.spec.invert_healthcheck
  disabled           = var.spec.disabled

  # --- CALCULATED: aggregate child checks ------------------------------------
  child_healthchecks     = local.child_healthchecks
  child_health_threshold = local.child_health_threshold

  # --- CLOUDWATCH_METRIC: mirror an alarm (the private-resource pattern) -----
  cloudwatch_alarm_name           = local.cloudwatch_alarm_name
  cloudwatch_alarm_region         = local.cloudwatch_alarm_region
  insufficient_data_health_status = local.insufficient_data_health_status

  # --- RECOVERY_CONTROL: mirror an ARC routing control (ForceNew) ------------
  routing_control_arn = local.routing_control_arn

  tags = local.aws_tags
}
