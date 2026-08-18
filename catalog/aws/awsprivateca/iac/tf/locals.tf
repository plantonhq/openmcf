locals {
  # Resource-identity tags match the Pulumi module key-for-key.
  aws_tags = {
    "Name"                     = var.metadata.name
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsPrivateCa"
    "planton.ai/resource-id"   = var.metadata.id
  }

  # Which composed-activation path applies (ROOT self-signs; a
  # SUBORDINATE activates from a parent AwsPrivateCa when configured;
  # otherwise the CA waits in PENDING_CERTIFICATE for out-of-band
  # activation).
  is_root                = var.spec.type == "ROOT"
  activates_subordinate  = var.spec.type == "SUBORDINATE" && var.spec.subordinate_activation != null
  has_composed_activation = var.spec.type == "ROOT" || (var.spec.type == "SUBORDINATE" && var.spec.subordinate_activation != null)
}
