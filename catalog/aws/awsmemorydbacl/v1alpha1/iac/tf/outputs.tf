output "acl_name" {
  description = "The ACL's AWS name -- what clusters attach via their acl_name."
  value       = aws_memorydb_acl.this.name
}

output "acl_arn" {
  description = "The ACL's Amazon Resource Name -- used in IAM policies and cross-service permissions."
  value       = aws_memorydb_acl.this.arn
}

output "minimum_engine_version" {
  description = "The minimum engine version the ACL's combined user set requires -- a cluster attaching this ACL must run at least this version."
  value       = aws_memorydb_acl.this.minimum_engine_version
}
