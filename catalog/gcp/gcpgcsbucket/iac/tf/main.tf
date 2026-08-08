# Enable the Cloud Storage API — the control plane that owns buckets.
# disable_on_destroy is false: tearing down one bucket must never disable
# the API for everything else in the project (other buckets keep serving).
resource "google_project_service" "storage_api" {
  project = local.project_id
  service = "storage.googleapis.com"

  disable_dependent_services = true
  disable_on_destroy         = false
}

# The GCS bucket. Sharp edges, all taught by the API rather than invented
# here:
#
#   - name, location, project, custom placement, hierarchical namespace,
#     and enable_object_retention are ForceNew — changing any of them
#     replaces the bucket and everything in it.
#
#   - force_destroy defaults to false: destroying a non-empty bucket fails
#     instead of silently erasing data. When true, the provider deletes
#     every object version first (which can take a long time on large
#     buckets and refuses objects still under a locked retention policy).
#
#   - A locked retention policy is irreversible; the provider forces
#     bucket re-creation on any attempt to unlock.
#
#   - Soft delete is a server-side default (7 days) even when the block is
#     omitted; the module sends the block only when the spec sets it, so
#     unset specs follow GCP's default without a perpetual diff.
resource "google_storage_bucket" "this" {
  name     = local.bucket_name
  project  = local.project_id
  location = local.location
  labels   = local.final_labels

  storage_class               = local.storage_class
  force_destroy               = var.spec.force_destroy
  uniform_bucket_level_access = var.spec.uniform_bucket_level_access_enabled
  public_access_prevention    = local.public_access_prevention
  requester_pays              = var.spec.requester_pays
  default_event_based_hold    = var.spec.default_event_based_hold
  rpo                         = local.rpo

  # Destroy-time guard: PREVENT fails the destroy; ABANDON unmanages the
  # bucket without deleting it. Null falls back to the provider default
  # (DELETE). Orthogonal to force_destroy, which governs whether a
  # permitted deletion may erase contained objects.
  deletion_policy = var.spec.deletion_policy != "" ? var.spec.deletion_policy : null

  # Create-time-only surfaces (ForceNew).
  enable_object_retention = var.spec.enable_object_retention ? true : null

  dynamic "hierarchical_namespace" {
    for_each = var.spec.hierarchical_namespace_enabled ? [1] : []
    content {
      enabled = true
    }
  }

  dynamic "versioning" {
    for_each = var.spec.versioning_enabled ? [1] : []
    content {
      enabled = true
    }
  }

  # Autoclass: GCS moves each object between classes on observed access.
  # The spec's CEL guard already rejects combining it with SetStorageClass
  # lifecycle rules.
  dynamic "autoclass" {
    for_each = var.spec.autoclass != null ? [var.spec.autoclass] : []
    content {
      enabled                = autoclass.value.enabled
      terminal_storage_class = autoclass.value.terminal_storage_class != "" ? autoclass.value.terminal_storage_class : null
    }
  }

  # Nullable numeric conditions pass through as-is: null means "criterion
  # unset" while an explicit 0 is a meaningful value the provider sends
  # (its *_if_zero handling rides on attribute presence).
  dynamic "lifecycle_rule" {
    for_each = var.spec.lifecycle_rules
    content {
      action {
        type          = lifecycle_rule.value.action.type
        storage_class = lifecycle_rule.value.action.storage_class != "" ? lifecycle_rule.value.action.storage_class : null
      }
      condition {
        age                                     = lifecycle_rule.value.condition.age_days
        send_age_if_zero                        = lifecycle_rule.value.condition.age_days == 0 ? true : null
        created_before                          = lifecycle_rule.value.condition.created_before != "" ? lifecycle_rule.value.condition.created_before : null
        with_state                              = lifecycle_rule.value.condition.with_state != "" ? lifecycle_rule.value.condition.with_state : null
        matches_storage_class                   = length(lifecycle_rule.value.condition.matches_storage_class) > 0 ? lifecycle_rule.value.condition.matches_storage_class : null
        matches_prefix                          = length(lifecycle_rule.value.condition.matches_prefix) > 0 ? lifecycle_rule.value.condition.matches_prefix : null
        matches_suffix                          = length(lifecycle_rule.value.condition.matches_suffix) > 0 ? lifecycle_rule.value.condition.matches_suffix : null
        num_newer_versions                      = lifecycle_rule.value.condition.num_newer_versions
        send_num_newer_versions_if_zero         = lifecycle_rule.value.condition.num_newer_versions == 0 ? true : null
        days_since_noncurrent_time              = lifecycle_rule.value.condition.days_since_noncurrent_time
        send_days_since_noncurrent_time_if_zero = lifecycle_rule.value.condition.days_since_noncurrent_time == 0 ? true : null
        noncurrent_time_before                  = lifecycle_rule.value.condition.noncurrent_time_before != "" ? lifecycle_rule.value.condition.noncurrent_time_before : null
        days_since_custom_time                  = lifecycle_rule.value.condition.days_since_custom_time
        send_days_since_custom_time_if_zero     = lifecycle_rule.value.condition.days_since_custom_time == 0 ? true : null
        custom_time_before                      = lifecycle_rule.value.condition.custom_time_before != "" ? lifecycle_rule.value.condition.custom_time_before : null
        size_above_bytes                        = lifecycle_rule.value.condition.size_above_bytes
        size_below_bytes                        = lifecycle_rule.value.condition.size_below_bytes
      }
    }
  }

  # WORM retention. Locking cannot happen at create — GCP locks in a
  # follow-up call after the policy exists.
  dynamic "retention_policy" {
    for_each = var.spec.retention_policy != null ? [var.spec.retention_policy] : []
    content {
      retention_period = retention_policy.value.retention_period_seconds
      is_locked        = retention_policy.value.is_locked
    }
  }

  # Sent only when the spec sets it; an omitted block follows GCP's 7-day
  # default. A set 0 disables soft delete.
  dynamic "soft_delete_policy" {
    for_each = var.spec.soft_delete_policy != null ? [var.spec.soft_delete_policy] : []
    content {
      retention_duration_seconds = soft_delete_policy.value.retention_duration_seconds
    }
  }

  # One provider block carries both the default CMEK key and the
  # per-encryption-type enforcement for new objects, so the module emits
  # it when either half is configured.
  #
  # Default CMEK: the GCS service agent must hold
  # roles/cloudkms.cryptoKeyEncrypterDecrypter on this key before create.
  # Enforcement changes apply to NEW objects only.
  dynamic "encryption" {
    for_each = (local.kms_key_name != null || var.spec.encryption_enforcement != null) ? [1] : []
    content {
      default_kms_key_name = local.kms_key_name

      dynamic "google_managed_encryption_enforcement_config" {
        for_each = try(var.spec.encryption_enforcement.google_managed_restriction_mode, "") != "" ? [var.spec.encryption_enforcement.google_managed_restriction_mode] : []
        content {
          restriction_mode = google_managed_encryption_enforcement_config.value
        }
      }
      dynamic "customer_managed_encryption_enforcement_config" {
        for_each = try(var.spec.encryption_enforcement.customer_managed_restriction_mode, "") != "" ? [var.spec.encryption_enforcement.customer_managed_restriction_mode] : []
        content {
          restriction_mode = customer_managed_encryption_enforcement_config.value
        }
      }
      dynamic "customer_supplied_encryption_enforcement_config" {
        for_each = try(var.spec.encryption_enforcement.customer_supplied_restriction_mode, "") != "" ? [var.spec.encryption_enforcement.customer_supplied_restriction_mode] : []
        content {
          restriction_mode = customer_supplied_encryption_enforcement_config.value
        }
      }
    }
  }

  dynamic "website" {
    for_each = var.spec.website != null ? [var.spec.website] : []
    content {
      main_page_suffix = website.value.main_page_suffix != "" ? website.value.main_page_suffix : null
      not_found_page   = website.value.not_found_page != "" ? website.value.not_found_page : null
    }
  }

  dynamic "cors" {
    for_each = var.spec.cors_rules
    content {
      origin          = cors.value.origins
      method          = cors.value.methods
      response_header = length(cors.value.response_headers) > 0 ? cors.value.response_headers : null
      max_age_seconds = cors.value.max_age_seconds
    }
  }

  dynamic "logging" {
    for_each = var.spec.logging != null ? [var.spec.logging] : []
    content {
      log_bucket        = logging.value.log_bucket
      log_object_prefix = logging.value.log_object_prefix != "" ? logging.value.log_object_prefix : null
    }
  }

  # Custom dual-region placement (exactly two regions, enforced pre-deploy).
  dynamic "custom_placement_config" {
    for_each = var.spec.custom_placement_config != null ? [var.spec.custom_placement_config] : []
    content {
      data_locations = custom_placement_config.value.data_locations
    }
  }

  # Network-layer IP filtering: which CIDR ranges / VPC networks may reach
  # the bucket at all, evaluated before IAM. The spec's CEL guard rejects
  # an Enabled filter with no sources pre-deploy.
  dynamic "ip_filter" {
    for_each = var.spec.ip_filter != null ? [var.spec.ip_filter] : []
    content {
      mode                           = ip_filter.value.mode
      allow_cross_org_vpcs           = ip_filter.value.allow_cross_org_vpcs ? true : null
      allow_all_service_agent_access = ip_filter.value.allow_all_service_agent_access ? true : null

      dynamic "public_network_source" {
        for_each = ip_filter.value.public_network_source != null ? [ip_filter.value.public_network_source] : []
        content {
          allowed_ip_cidr_ranges = public_network_source.value.allowed_ip_cidr_ranges
        }
      }

      dynamic "vpc_network_sources" {
        for_each = ip_filter.value.vpc_network_sources
        content {
          network                = vpc_network_sources.value.network
          allowed_ip_cidr_ranges = vpc_network_sources.value.allowed_ip_cidr_ranges
        }
      }
    }
  }

  depends_on = [
    google_project_service.storage_api,
  ]
}

# Additive IAM grants: one (role, member) pair per resource, merging into
# the bucket's policy without touching grants made elsewhere — authoritative
# bindings/policies are deliberately not used.
resource "google_storage_bucket_iam_member" "members" {
  for_each = local.iam_members

  bucket = google_storage_bucket.this.name
  role   = each.value.role
  member = each.value.member

  dynamic "condition" {
    for_each = each.value.condition != null ? [each.value.condition] : []
    content {
      title       = condition.value.title
      expression  = condition.value.expression
      description = condition.value.description != "" ? condition.value.description : null
    }
  }
}
