locals {
  # Resource-identity tags match the Pulumi module key-for-key.
  aws_tags = {
    "Name"                     = var.metadata.name
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsBedrockAgentCoreEvaluation"
    "planton.ai/resource-id"   = var.metadata.id
  }

  # Arms keyed by their stable entry names (the for_each keys both
  # engines share and the output map keys).
  evaluators     = { for e in var.spec.evaluators : e.name => e }
  harnesses      = { for h in var.spec.harnesses : h.name => h }
  online_configs = { for c in var.spec.online_evaluation_configs : c.name => c }

  # The set of in-bundle evaluator names: online-config evaluator
  # entries naming one resolve to the CREATED evaluator's
  # AWS-generated ID (and gain the dependency edge); builtins and full
  # custom IDs pass through as literals.
  bundle_evaluator_names = toset([for e in var.spec.evaluators : e.name])
}
