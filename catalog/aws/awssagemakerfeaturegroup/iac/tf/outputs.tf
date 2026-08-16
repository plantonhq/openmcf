output "feature_group_name" {
  description = "The feature group name (the AWS identity ingestion and serving calls use)"
  value       = aws_sagemaker_feature_group.this.feature_group_name
}

output "feature_group_arn" {
  description = "The Amazon Resource Name of the feature group"
  value       = aws_sagemaker_feature_group.this.arn
}
