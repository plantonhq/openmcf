locals {
  # The bucket arrives as a plain name: the StringValueOrRef foreign key is
  # flattened to its resolved string before the module runs (a referenced
  # AwsS3Bucket resolves to status.outputs.bucket_id).
  bucket_name = var.spec.bucket

  # Resource-identity tags match the Pulumi module key-for-key (the canonical
  # six-key identity map — user labels never merge into cloud tags). Every
  # object in the set carries them (S3 object tags), attributing each object
  # back to the owning set for auditing and orphan cleanup.
  base_tags = {
    "Name"                     = var.metadata.name
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsS3ObjectSet"
    "planton.ai/resource-id"   = var.metadata.id
  }

  # Set-level spec tags layer on top of the identity tags; per-object tags are
  # merged last in main.tf so they win on key collisions.
  set_tags = merge(local.base_tags, try(var.spec.tags, {}))

  # Key the provider resources by the S3 object key, so adding, removing, or
  # REORDERING entries in the manifest never churns unrelated objects — each
  # object's Terraform identity is its key, mirroring its S3 identity. The
  # source arm picks the resource type: inline content renders aws_s3_object,
  # a copy_from source renders aws_s3_object_copy.
  content_objects_map = {
    for obj in var.spec.objects : obj.key => obj if try(obj.copy_from, null) == null
  }
  copy_objects_map = {
    for obj in var.spec.objects : obj.key => obj if try(obj.copy_from, null) != null
  }

  # The provider's copy source string: "bucket/key" for bucket-named sources,
  # "arn.../object/key" when the source bucket arrives as an S3 access-point
  # ARN (the provider's two documented source formats).
  copy_sources = {
    for key, obj in local.copy_objects_map :
    key => startswith(obj.copy_from.source_bucket, "arn:") ? "${obj.copy_from.source_bucket}/object/${obj.copy_from.source_key}" : "${obj.copy_from.source_bucket}/${obj.copy_from.source_key}"
  }
}
