output "namespace" {
  description = "The namespace in which the Postgres resources are deployed."
  value       = local.namespace
}

output "service" {
  description = "Name of the Postgres service (master)."
  value       = local.kube_service_name
}

output "port_forward_command" {
  description = "Convenient command to port-forward to the Postgres service."
  value       = local.kube_port_forward_command
}

output "kube_endpoint" {
  description = "FQDN of the Postgres service within the cluster."
  value       = local.kube_service_fqdn
}

output "external_hostname" {
  description = "The external hostname for Postgres if ingress is enabled."
  value       = local.ingress_external_hostname
}

output "password_secret_name" {
  description = "Name of the secret holding the Postgres password."
  # Matches the pattern used by Zalando's operator for password credentials
  value       = "postgres.db-${local.resource_id}.credentials.postgresql.acid.zalan.do"
}

output "password_secret_key" {
  description = "Key within the secret that contains the Postgres password."
  value       = "password"
}

output "username_secret_name" {
  description = "Name of the secret holding the Postgres username."
  value       = "postgres.db-${local.resource_id}.credentials.postgresql.acid.zalan.do"
}

output "username_secret_key" {
  description = "Key within the secret that contains the Postgres username."
  value       = "username"
}
