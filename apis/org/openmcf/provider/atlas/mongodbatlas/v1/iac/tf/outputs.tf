# Terraform Outputs
# These outputs map to the MongodbAtlasStackOutputs proto definition
# Documentation: https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/advanced_cluster#attributes-reference

output "id" {
  description = "The provider-assigned unique ID for the MongoDB Atlas cluster resource (cluster_id)"
  value       = mongodbatlas_advanced_cluster.main.cluster_id
}

# Note: The stack_outputs.proto currently has fields for bootstrap_endpoint, crn, and rest_endpoint
# which appear to be copied from a Kafka cluster resource. These fields don't directly map to 
# MongoDB Atlas cluster outputs. We provide connection string and state information instead.

output "bootstrap_endpoint" {
  description = "MongoDB Atlas standard connection string (SRV format)"
  value       = mongodbatlas_advanced_cluster.main.connection_strings[0].standard_srv
}

output "crn" {
  description = "MongoDB Atlas cluster resource name (cluster_id for identification)"
  value       = mongodbatlas_advanced_cluster.main.cluster_id
}

output "rest_endpoint" {
  description = "MongoDB Atlas standard connection string (for backwards compatibility)"
  value       = mongodbatlas_advanced_cluster.main.connection_strings[0].standard
}

# Additional useful outputs not in the proto but helpful for users
# These provide comprehensive connection information for various use cases

output "connection_strings" {
  description = "Complete connection strings object with all available formats"
  value       = mongodbatlas_advanced_cluster.main.connection_strings
  sensitive   = true
}

output "connection_string_standard" {
  description = "Standard format connection string"
  value       = mongodbatlas_advanced_cluster.main.connection_strings[0].standard
  sensitive   = true
}

output "connection_string_standard_srv" {
  description = "Standard SRV format connection string (recommended for drivers)"
  value       = mongodbatlas_advanced_cluster.main.connection_strings[0].standard_srv
  sensitive   = true
}

output "connection_string_private" {
  description = "Private endpoint connection string (if private endpoints are configured)"
  value       = try(mongodbatlas_advanced_cluster.main.connection_strings[0].private, "")
  sensitive   = true
}

output "connection_string_private_srv" {
  description = "Private endpoint SRV connection string (if private endpoints are configured)"
  value       = try(mongodbatlas_advanced_cluster.main.connection_strings[0].private_srv, "")
  sensitive   = true
}

output "cluster_name" {
  description = "The name of the MongoDB Atlas cluster"
  value       = mongodbatlas_advanced_cluster.main.name
}

output "cluster_type" {
  description = "The type of the cluster (REPLICASET, SHARDED, or GEOSHARDED)"
  value       = mongodbatlas_advanced_cluster.main.cluster_type
}

output "state_name" {
  description = "Current state of the cluster (IDLE, CREATING, UPDATING, DELETING, etc.)"
  value       = mongodbatlas_advanced_cluster.main.state_name
}

output "mongo_db_version" {
  description = "Version of MongoDB the cluster is running"
  value       = mongodbatlas_advanced_cluster.main.mongo_db_version
}

output "project_id" {
  description = "The unique ID for the Atlas project containing this cluster"
  value       = mongodbatlas_advanced_cluster.main.project_id
}

