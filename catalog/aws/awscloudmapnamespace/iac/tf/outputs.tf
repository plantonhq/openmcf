output "namespace_id" {
  description = "The namespace's id (ns-...) - the provider's import ID (a PRIVATE_DNS namespace imports as namespace_id:vpc_id)"
  value       = local.namespace_id
}

output "namespace_arn" {
  description = "The namespace's ARN"
  value       = local.namespace_arn
}

output "hosted_zone_id" {
  description = "The Route 53 hosted zone Cloud Map created for a DNS namespace; empty for HTTP"
  value       = local.hosted_zone_id
}

output "http_name" {
  description = "The name an HTTP namespace's DiscoverInstances calls use; empty for DNS namespaces"
  value       = local.http_name
}

output "service_ids" {
  description = "AWS-generated service IDs (srv-...) keyed by service name"
  value       = { for name, service in aws_service_discovery_service.this : name => service.id }
}

output "service_arns" {
  description = "Service ARNs keyed by service name - what ECS service registries wire as the registry_arn"
  value       = { for name, service in aws_service_discovery_service.this : name => service.arn }
}

output "instance_service_ids" {
  description = "Each registration's owning service ID keyed by service_name//instance_id (the first half of the instance's composite import ID)"
  # try(): during an import round-trip the state holds a PARTIAL echo
  # (services import one at a time), and an eager index on a
  # not-yet-imported service name hard-errors the whole plan instead
  # of deferring - the documented import-time partial-echo class. The
  # tolerated null disappears once every service is in state.
  value = { for key, entry in local.service_instances_by_key : key => try(aws_service_discovery_service.this[entry.service_name].id, null) }
}
