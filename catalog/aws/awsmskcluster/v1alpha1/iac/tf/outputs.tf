output "cluster_arn" {
  description = "ARN of the MSK cluster"
  value       = aws_msk_cluster.this.arn
}

output "cluster_name" {
  description = "Name of the MSK cluster"
  value       = aws_msk_cluster.this.cluster_name
}

output "cluster_uuid" {
  description = "UUID of the MSK cluster, for use in IAM policies"
  value       = aws_msk_cluster.this.cluster_uuid
}

output "current_version" {
  description = "Current version of the MSK cluster, required for cluster updates"
  value       = aws_msk_cluster.this.current_version
}

output "bootstrap_brokers" {
  description = "Comma-separated plaintext broker endpoints (port 9092); empty when client-broker encryption is TLS-only"
  value       = aws_msk_cluster.this.bootstrap_brokers
}

output "bootstrap_brokers_tls" {
  description = "Comma-separated TLS broker endpoints (port 9094)"
  value       = aws_msk_cluster.this.bootstrap_brokers_tls
}

output "bootstrap_brokers_sasl_iam" {
  description = "Comma-separated SASL/IAM broker endpoints (port 9098)"
  value       = aws_msk_cluster.this.bootstrap_brokers_sasl_iam
}

output "bootstrap_brokers_sasl_scram" {
  description = "Comma-separated SASL/SCRAM broker endpoints (port 9096)"
  value       = aws_msk_cluster.this.bootstrap_brokers_sasl_scram
}

output "bootstrap_brokers_public_tls" {
  description = "Comma-separated public TLS broker endpoints"
  value       = aws_msk_cluster.this.bootstrap_brokers_public_tls
}

output "bootstrap_brokers_public_sasl_iam" {
  description = "Comma-separated public SASL/IAM broker endpoints"
  value       = aws_msk_cluster.this.bootstrap_brokers_public_sasl_iam
}

output "bootstrap_brokers_public_sasl_scram" {
  description = "Comma-separated public SASL/SCRAM broker endpoints"
  value       = aws_msk_cluster.this.bootstrap_brokers_public_sasl_scram
}

output "bootstrap_brokers_vpc_connectivity_tls" {
  description = "Comma-separated PrivateLink (multi-VPC) mutual-TLS broker endpoints"
  value       = aws_msk_cluster.this.bootstrap_brokers_vpc_connectivity_tls
}

output "bootstrap_brokers_vpc_connectivity_sasl_iam" {
  description = "Comma-separated PrivateLink (multi-VPC) SASL/IAM broker endpoints"
  value       = aws_msk_cluster.this.bootstrap_brokers_vpc_connectivity_sasl_iam
}

output "bootstrap_brokers_vpc_connectivity_sasl_scram" {
  description = "Comma-separated PrivateLink (multi-VPC) SASL/SCRAM broker endpoints"
  value       = aws_msk_cluster.this.bootstrap_brokers_vpc_connectivity_sasl_scram
}

output "zookeeper_connect_string" {
  description = "Comma-separated ZooKeeper endpoints (plaintext); empty on KRaft-mode clusters"
  value       = aws_msk_cluster.this.zookeeper_connect_string
}

output "zookeeper_connect_string_tls" {
  description = "Comma-separated ZooKeeper endpoints (TLS); empty on KRaft-mode clusters"
  value       = aws_msk_cluster.this.zookeeper_connect_string_tls
}

output "configuration_arn" {
  description = "ARN of the module-managed MSK Configuration created from server_properties, if any"
  value       = local.manage_configuration ? aws_msk_configuration.this[0].arn : ""
}
