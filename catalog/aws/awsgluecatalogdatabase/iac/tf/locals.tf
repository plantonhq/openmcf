locals {
  # The database's cloud name is the resource's metadata.name. Glue database
  # names are create-time-immutable (changing the name replaces the database)
  # and AWS rejects uppercase characters -- the name is passed verbatim, never
  # silently transformed. Same basis as the Pulumi module.
  database_name = var.metadata.name

  # Resource-identity tags match the Pulumi module key-for-key. Identity
  # tagging is the only tagging surface this module manages; user-defined
  # custom tags are a platform-wide concern, not per-kind spec surface.
  # (Tags land on the Glue resource itself; the spec's `parameters` map is
  # catalog metadata inside the database, a different surface.)
  tags = {
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsGlueCatalogDatabase"
    "planton.ai/resource-id"   = var.metadata.id
  }

  # Shape switches: a database is exactly one of regular / resource link /
  # federated (spec CEL enforces the exclusivity); the grants list rides on
  # the regular shape.
  has_target_database    = var.spec.target_database != null
  has_federated_database = var.spec.federated_database != null
  has_create_table_perms = length(var.spec.create_table_default_permissions) > 0
}
