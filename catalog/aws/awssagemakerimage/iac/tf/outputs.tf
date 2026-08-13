output "image_name" {
  description = "The image name (the AWS identity Studio configurations reference)"
  value       = aws_sagemaker_image.this.image_name
}

output "image_arn" {
  description = "The Amazon Resource Name of the image"
  value       = aws_sagemaker_image.this.arn
}

output "version_numbers" {
  description = "AWS-assigned version numbers keyed by each versions entry's position"
  value       = { for k, v in aws_sagemaker_image_version.this : k => tostring(v.version) }
}
