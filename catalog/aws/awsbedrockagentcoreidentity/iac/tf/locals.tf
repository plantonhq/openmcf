locals {
  # Resource-identity tags match the Pulumi module key-for-key.
  aws_tags = {
    "Name"                     = var.metadata.name
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsBedrockAgentCoreIdentity"
    "planton.ai/resource-id"   = var.metadata.id
  }

  # Arms keyed by their stable entry names (the for_each keys both
  # engines share and the output map keys).
  workload_identities = { for w in var.spec.workload_identities : w.name => w }
  api_key_providers   = { for p in var.spec.api_key_credential_providers : p.name => p }
  oauth2_providers    = { for p in var.spec.oauth2_credential_providers : p.name => p }

  has_policy_engine = var.spec.policy_engine != null
  policies          = local.has_policy_engine ? { for p in var.spec.policy_engine.policies : p.name => p } : {}

  # The spec's clean vendor vocabulary maps to the provider's enum values
  # AND selects which vendor block renders -- the vendor field IS the
  # discriminator over six structurally-identical provider blocks.
  oauth2_vendor_values = {
    CUSTOM     = "CustomOauth2"
    GITHUB     = "GithubOauth2"
    GOOGLE     = "GoogleOauth2"
    MICROSOFT  = "MicrosoftOauth2"
    SALESFORCE = "SalesforceOauth2"
    SLACK      = "SlackOauth2"
  }
}
