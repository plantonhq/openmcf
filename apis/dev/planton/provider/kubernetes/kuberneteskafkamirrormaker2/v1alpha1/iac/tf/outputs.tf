# Stack outputs — flattened onto KubernetesKafkaMirrorMaker2StackOutputs
# by the platform. Keep in lockstep with the Pulumi module's exports.

output "namespace" {
  description = "Namespace the MirrorMaker 2 deployment runs in"
  value       = local.namespace
}

output "mirrormaker_name" {
  description = "The deployment's name (metadata.name)"
  value       = local.mirrormaker_name
}

output "rest_api_endpoint" {
  description = "In-cluster Connect REST API endpoint of the underlying engine (http://<name>-mirrormaker2-api.<namespace>.svc.cluster.local:8083) — read-only inspection of mirror connector status"
  value       = local.rest_api_endpoint
}
