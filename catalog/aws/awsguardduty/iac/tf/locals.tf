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

  # AWS materializes a feature's FULL sub-toggle family server-side:
  # RUNTIME_MONITORING's trio reads back with every undeclared member
  # DISABLED, so a partial send leaves the refreshed state carrying
  # blocks the config lacks and every later plan proposes stripping
  # sub-toggles AWS will re-materialize (live-caught by the post-apply
  # idempotency re-plan). The module therefore always sends the
  # complete family: declared members with their declared value,
  # undeclared members an explicit DISABLED. Family membership is the
  # provider enum at the pin (~> 6.58).
  runtime_monitoring_subtoggles = [
    "EC2_AGENT_MANAGEMENT",
    "ECS_FARGATE_AGENT_MANAGEMENT",
    "EKS_ADDON_MANAGEMENT",
  ]

  feature_additional_configurations = {
    for name, f in local.features : name => (
      name == "RUNTIME_MONITORING"
      ? [
        for sub in local.runtime_monitoring_subtoggles : {
          name = sub
          status = lookup(
            { for ac in f.additional_configuration : ac.name => (ac.enabled == null ? "ENABLED" : (ac.enabled ? "ENABLED" : "DISABLED")) },
            sub,
            "DISABLED",
          )
        }
      ]
      : [
        for ac in f.additional_configuration : {
          name   = ac.name
          status = ac.enabled == null ? "ENABLED" : (ac.enabled ? "ENABLED" : "DISABLED")
        }
      ]
    )
  }
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
