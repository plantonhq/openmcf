locals {
  # Resource-identity tags match the Pulumi module key-for-key. The
  # organization resource itself is untaggable - the tags land on the
  # taggable folded satellite (the resource policy).
  aws_tags = {
    "Name"                     = var.metadata.name
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsOrganization"
    "planton.ai/resource-id"   = var.metadata.id
  }

  # Delegated administrator registrations key by account and principal.
  # The "//" separator is the import machinery's segment delimiter --
  # each half of the provider's "{account_id}/{service_principal}"
  # import composite derives from its own key segment.
  delegated_administrators = {
    for d in var.spec.delegated_administrators :
    "${d.account_id}//${d.service_principal}" => d
  }
}
