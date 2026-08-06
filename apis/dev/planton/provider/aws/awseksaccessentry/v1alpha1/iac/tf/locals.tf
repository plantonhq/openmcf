locals {
  # Key each association by its policy's name (the last ARN segment,
  # e.g. "AmazonEKSViewPolicy") -- stable across scope edits, unique per
  # entry because AWS allows one association per policy per principal.
  policy_associations_by_name = {
    for association in var.spec.policy_associations :
    element(split("/", association.policy_arn), length(split("/", association.policy_arn)) - 1) => association
  }

  # Resource-identity tags, matching the Pulumi module key-for-key.
  aws_tags = {
    "Name"                     = var.metadata.name
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsEksAccessEntry"
    "planton.ai/resource-id"   = var.metadata.id
  }
}
