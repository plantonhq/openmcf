output "bucket_id" {
  description = "The bucket the objects were uploaded to"
  value       = local.bucket_name
}

output "object_arns" {
  description = "Map of object key to its ARN (arn:aws:s3:::bucket/key), for IAM policy Resource lists"
  value = {
    for key, obj in aws_s3_object.objects : key => obj.arn
  }
}

output "object_etags" {
  description = "Map of object key to its ETag (content hash)"
  value = {
    for key, obj in aws_s3_object.objects : key => obj.etag
  }
}

output "object_version_ids" {
  description = "Map of object key to its version ID (if bucket versioning is enabled)"
  value = {
    for key, obj in aws_s3_object.objects : key => obj.version_id
  }
}
