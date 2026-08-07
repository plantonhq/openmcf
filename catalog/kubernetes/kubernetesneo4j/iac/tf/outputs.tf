# Stack outputs — identical names and derivations in the Pulumi module's
# outputs.go / main.go exports.
#
# The service name is the chart's always-created ClusterIP Service:
# templates/neo4j-svc.yaml names it after neo4j.fullname, which is the
# release name here (no name overrides are rendered).

output "namespace" {
  description = "Kubernetes namespace the Neo4j server runs in"
  value       = local.namespace
}

output "release_name" {
  description = "Helm release name (= metadata.name)"
  value       = local.release_name
}

output "service_name" {
  description = "Name of the main Neo4j Service (bolt/http ports; = the release name)"
  value       = local.service_name
}

output "bolt_endpoint" {
  description = "In-cluster bolt endpoint drivers connect to"
  value       = "neo4j://${local.service_name}.${local.namespace}.svc.cluster.local:7687"
}

output "http_endpoint" {
  description = "In-cluster HTTP API / Browser endpoint"
  value       = "http://${local.service_name}.${local.namespace}.svc.cluster.local:7474"
}

output "auth_secret_name" {
  description = "Secret holding the admin credentials (NEO4J_AUTH key): the module-materialized <name>-auth, the referenced existing Secret, or empty when the chart generated a random password"
  value       = local.auth_secret_name
}

output "port_forward_command" {
  description = "kubectl one-liner for reaching bolt from a workstation"
  value       = "kubectl port-forward svc/${local.service_name} -n ${local.namespace} 7687:7687"
}
