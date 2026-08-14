locals {
  # The gateway name is metadata.name -- the naming basis both engines
  # share so a manifest deploys identically on either. AWS allows letters
  # and digits with single hyphens, max 100 characters (no underscores,
  # no consecutive hyphens).
  gateway_name = var.metadata.name

  # Resource-identity tags match the Pulumi module key-for-key.
  aws_tags = {
    "Name"                     = var.metadata.name
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsBedrockAgentCoreGateway"
    "planton.ai/resource-id"   = var.metadata.id
  }

  # Optional single-entry surfaces render only when declared.
  has_jwt           = var.spec.custom_jwt_authorizer != null
  has_mcp           = var.spec.mcp != null
  has_policy_engine = var.spec.policy_engine != null

  # Targets keyed by their stable entry names (the for_each keys both
  # engines share and the output map keys). The collection is `any`-typed
  # in variables.tf (heterogeneous JSON-document members defeat HCL's
  # object-type unification), so entry access below is try()-based sparse
  # access against the tfvars converter's always-complete keys.
  targets = { for t in var.spec.targets : t.name => t }
}
