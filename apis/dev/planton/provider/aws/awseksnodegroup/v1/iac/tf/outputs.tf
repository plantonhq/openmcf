output "nodegroup_name" {
  description = "The name of the managed node group."
  value       = aws_eks_node_group.this.node_group_name
}

output "nodegroup_arn" {
  description = "The Amazon Resource Name of the node group."
  value       = aws_eks_node_group.this.arn
}

output "asg_name" {
  description = "The name of the EC2 Auto Scaling group AWS manages behind the node group."
  value       = try(aws_eks_node_group.this.resources[0].autoscaling_groups[0].name, "")
}

output "remote_access_sg_id" {
  description = "The ID of the security group AWS creates for SSH access when remote access is enabled without explicit source groups; empty otherwise."
  value       = try(aws_eks_node_group.this.resources[0].remote_access_security_group_id, "")
}
