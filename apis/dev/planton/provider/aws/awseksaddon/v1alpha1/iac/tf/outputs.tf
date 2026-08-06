output "addon_arn" {
  description = "The Amazon Resource Name of the add-on."
  value       = aws_eks_addon.this.arn
}

output "addon_name" {
  description = "The EKS catalog name the add-on was installed under."
  value       = aws_eks_addon.this.addon_name
}

output "addon_version" {
  description = "The version actually running -- the resolved AWS default when the spec pinned nothing."
  value       = aws_eks_addon.this.addon_version
}
