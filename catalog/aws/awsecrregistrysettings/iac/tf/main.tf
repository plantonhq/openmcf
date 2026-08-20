# One region's ECR registry-level configuration: the registry policy,
# scanning, replication, pull-through cache rules, repository creation
# templates, account settings, and pull-time update exclusions.
#
# Lifecycle facts the render below depends on:
#   - the policy/scanning/replication arms are one-per-registry
#     singletons keyed by the account id; scanning and replication
#     have RESET-not-delete semantics at the provider (destroy puts
#     the empty default back);
#   - account settings are PutAccountSetting upserts with a NO-OP
#     delete - the last-applied values persist after destroy;
#   - pull-through cache rules and creation templates are prefix-keyed
#     (ForceNew prefixes - the for_each keys below); clearing a cache
#     rule's credential/custom-role ARN back to empty is NOT
#     propagated by the provider (replace the rule to drop
#     credentials);
#   - pull-time update exclusions are immutable per principal ARN.
#
# Every resource here is untaggable at AWS - hence no tags anywhere
# (a creation template's resource_tags are the STAMPED repositories'
# tags, user surface from the spec).

# The registry's identity is the account itself - resolved once for
# imports and the pull-URL join key.
data "aws_caller_identity" "this" {}

resource "aws_ecr_registry_policy" "this" {
  count = var.spec.registry_policy != "" ? 1 : 0

  policy = var.spec.registry_policy
}

resource "aws_ecr_registry_scanning_configuration" "this" {
  count = var.spec.scanning != null ? 1 : 0

  scan_type = var.spec.scanning.scan_type

  dynamic "rule" {
    for_each = var.spec.scanning.rules
    content {
      scan_frequency = rule.value.scan_frequency

      dynamic "repository_filter" {
        for_each = rule.value.filters
        content {
          filter = repository_filter.value
          # WILDCARD is the only filter type AWS supports - pinned
          # here, never spec surface.
          filter_type = "WILDCARD"
        }
      }
    }
  }
}

resource "aws_ecr_replication_configuration" "this" {
  count = length(var.spec.replication_rules) > 0 ? 1 : 0

  replication_configuration {
    dynamic "rule" {
      for_each = var.spec.replication_rules
      content {
        dynamic "destination" {
          for_each = rule.value.destinations
          content {
            region      = destination.value.region
            registry_id = destination.value.registry_id
          }
        }

        dynamic "repository_filter" {
          for_each = rule.value.repository_filters
          content {
            filter = repository_filter.value
            # PREFIX_MATCH is the only filter type AWS supports -
            # pinned here, never spec surface.
            filter_type = "PREFIX_MATCH"
          }
        }
      }
    }
  }
}

# Pull-through cache rules, keyed by prefix (the rules' own import
# IDs).
resource "aws_ecr_pull_through_cache_rule" "this" {
  for_each = { for cache_rule in var.spec.pull_through_cache_rules : cache_rule.ecr_repository_prefix => cache_rule }

  ecr_repository_prefix      = each.value.ecr_repository_prefix
  upstream_registry_url      = each.value.upstream_registry_url
  upstream_repository_prefix = each.value.upstream_repository_prefix != "" ? each.value.upstream_repository_prefix : null
  credential_arn             = each.value.credential_arn != "" ? each.value.credential_arn : null
  custom_role_arn            = each.value.custom_role_arn != "" ? each.value.custom_role_arn : null
}

# Repository creation templates, keyed by prefix (the templates' own
# import IDs).
resource "aws_ecr_repository_creation_template" "this" {
  for_each = { for template in var.spec.repository_creation_templates : template.prefix => template }

  prefix          = each.value.prefix
  description     = each.value.description != "" ? each.value.description : null
  applied_for     = each.value.applied_for
  custom_role_arn = each.value.custom_role_arn != "" ? each.value.custom_role_arn : null

  image_tag_mutability = each.value.image_tag_mutability != "" ? each.value.image_tag_mutability : null

  dynamic "image_tag_mutability_exclusion_filter" {
    for_each = each.value.image_tag_mutability_exclusion_filters
    content {
      filter = image_tag_mutability_exclusion_filter.value
      # WILDCARD is the only filter type AWS supports.
      filter_type = "WILDCARD"
    }
  }

  dynamic "encryption_configuration" {
    for_each = each.value.encryption != null ? [each.value.encryption] : []
    content {
      encryption_type = encryption_configuration.value.type != "" ? encryption_configuration.value.type : null
      kms_key         = encryption_configuration.value.kms_key != "" ? encryption_configuration.value.kms_key : null
    }
  }

  lifecycle_policy  = each.value.lifecycle_policy != "" ? each.value.lifecycle_policy : null
  repository_policy = each.value.repository_policy != "" ? each.value.repository_policy : null
  resource_tags     = length(each.value.resource_tags) > 0 ? each.value.resource_tags : null
}

locals {
  # Each configured toggle is its own PutAccountSetting upsert; unset
  # toggles keep the account's current values (and ALL of them persist
  # after destroy).
  account_settings = var.spec.account_settings == null ? {} : merge(
    var.spec.account_settings.basic_scan_type_version != "" ? { BASIC_SCAN_TYPE_VERSION = var.spec.account_settings.basic_scan_type_version } : {},
    var.spec.account_settings.blob_mounting != null ? { BLOB_MOUNTING = var.spec.account_settings.blob_mounting ? "ENABLED" : "DISABLED" } : {},
    var.spec.account_settings.registry_policy_scope != "" ? { REGISTRY_POLICY_SCOPE = var.spec.account_settings.registry_policy_scope } : {}
  )
}

resource "aws_ecr_account_setting" "this" {
  for_each = local.account_settings

  name  = each.key
  value = each.value
}

# Pull-time update exclusions, keyed by the resolved principal ARN
# (the exclusions' own import IDs).
resource "aws_ecr_pull_time_update_exclusion" "this" {
  for_each = toset(var.spec.pull_time_update_exclusions)

  principal_arn = each.value
}
