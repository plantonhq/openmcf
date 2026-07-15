locals {
  # Resource-identity tags follow the catalog convention. Tags land on
  # the identity itself; the mail-from/feedback/policy satellites are
  # sub-resources keyed by the identity and carry no tags of their own.
  aws_tags = {
    "Name"                     = var.metadata.name
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsSesEmailIdentity"
    "planton.ai/resource-id"   = var.metadata.id
  }

  # The DKIM block is emitted only when the manifest configures signing --
  # an absent block accepts AWS's Easy DKIM default (2048-bit key) without
  # materializing signing attributes into the resource.
  dkim = var.spec.dkim_signing

  # BYODKIM is selected by the key/selector pair (spec-level CEL keeps the
  # pair consistent and exclusive with next_signing_key_length). The key
  # carries the contract default "" when absent -- a plain != ""
  # comparison, never coalesce(x, ""), which errors when empty.
  byodkim = local.dkim != null ? local.dkim.domain_signing_private_key != "" : false

  # Identity policies are keyed by name so each maps to its own AWS
  # sub-resource and entries add/remove independently.
  policies = { for p in coalesce(var.spec.policies, []) : p.name => p }
}
