locals {
  # The bucket arrives as a plain name: the StringValueOrRef foreign key is
  # flattened to its resolved string before the module runs (a referenced
  # AwsS3Bucket resolves to status.outputs.bucket_id).
  bucket_name = var.spec.bucket

  # Resource-identity tags follow the catalog convention; user labels merge in
  # without being able to override the identity keys. Every object in the set
  # carries them (S3 object tags), attributing each object back to the owning
  # set for auditing and orphan cleanup.
  base_tags = merge(try(var.metadata.labels, {}), {
    "Name"                     = var.metadata.name
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsS3ObjectSet"
    "planton.ai/resource-id"   = var.metadata.id
  })

  # Set-level spec tags layer on top of the identity tags; per-object tags are
  # merged last in main.tf so they win on key collisions.
  set_tags = merge(local.base_tags, try(var.spec.tags, {}))

  # Key the provider resources by the S3 object key, so adding, removing, or
  # REORDERING entries in the manifest never churns unrelated objects — each
  # object's Terraform identity is its key, mirroring its S3 identity.
  objects_map = {
    for obj in var.spec.objects : obj.key => obj
  }
}
