locals {
  # Resource-identity tags match the Pulumi module key-for-key. The
  # detector, filters, and IP/threat lists are taggable; the feature
  # patches, org surface, and members are not, and the publishing
  # destination is deliberately untagged (ForceNew tags - see main.tf).
  aws_tags = {
    "Name"                     = var.metadata.name
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsGuardDuty"
    "planton.ai/resource-id"   = var.metadata.id
  }

  # Stable entry-name keys - the for_each keys both engines share and
  # the output map keys.
  features          = { for f in var.spec.features : f.name => f }
  filters           = { for f in var.spec.filters : f.name => f }
  ip_sets           = { for s in var.spec.ip_sets : s.name => s }
  threat_intel_sets = { for s in var.spec.threat_intel_sets : s.name => s }
  members           = { for m in var.spec.members : m.account_id => m }

  org_features = var.spec.organization != null ? { for f in var.spec.organization.features : f.name => f } : {}

  # Member feature patches, keyed "account/feature" (two stable keys).
  member_features = merge([
    for m in var.spec.members : {
      for f in m.features : "${m.account_id}/${f.name}" => {
        account_id = m.account_id
        feature    = f
      }
    }
  ]...)
}
