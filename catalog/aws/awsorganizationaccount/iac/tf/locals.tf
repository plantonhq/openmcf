locals {
  # Resource-identity tags match the Pulumi module key-for-key. The
  # account resource is the kind's one taggable surface (the contact
  # and region satellites are untaggable).
  aws_tags = {
    "Name"                     = var.metadata.name
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsOrganizationAccount"
    "planton.ai/resource-id"   = var.metadata.id
  }

  # Alternate contacts key by the same type token the provider imports
  # them under ("{account_id}/{BILLING|OPERATIONS|SECURITY}").
  alternate_contacts = var.spec.alternate_contacts == null ? {} : {
    for type, contact in {
      BILLING    = var.spec.alternate_contacts.billing
      OPERATIONS = var.spec.alternate_contacts.operations
      SECURITY   = var.spec.alternate_contacts.security
    } : type => contact if contact != null
  }

  # Region enablements key by region name (each imports as
  # "{account_id},{region_name}").
  regions = {
    for r in var.spec.regions : r.region_name => r
  }
}
