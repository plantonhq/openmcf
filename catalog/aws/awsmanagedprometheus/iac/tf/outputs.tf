output "workspace_id" {
  description = "The workspace's ID (the provider's import ID for the workspace and most satellites)"
  value       = aws_prometheus_workspace.this.id
}

output "workspace_arn" {
  description = "The workspace's ARN - what scrapers and remote-write clients reference"
  value       = aws_prometheus_workspace.this.arn
}

output "prometheus_endpoint" {
  description = "The workspace's Prometheus-compatible query/remote-write endpoint URL"
  value       = aws_prometheus_workspace.this.prometheus_endpoint
}

output "rule_group_namespace_arns" {
  description = "Rule-groups namespace ARNs keyed by namespace name (each namespace imports by its ARN)"
  value       = { for name, namespace in aws_prometheus_rule_group_namespace.this : name => namespace.arn }
}

output "anomaly_detector_ids" {
  description = "AWS-generated anomaly detector IDs keyed by detector alias (each imports as \"detector_id,workspace_id\")"
  value       = { for alias, detector in aws_prometheus_anomaly_detector.this : alias => detector.id }
}

output "anomaly_detector_arns" {
  description = "Anomaly detector ARNs keyed by detector alias"
  value       = { for alias, detector in aws_prometheus_anomaly_detector.this : alias => detector.arn }
}
