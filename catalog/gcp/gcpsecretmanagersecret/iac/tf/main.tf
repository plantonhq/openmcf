# Enable the Secret Manager API so a fresh project can host the secret.
# disable_on_destroy is false: tearing down one secret must never disable
# the API for everything else in the project.
resource "google_project_service" "secretmanager_api" {
  project = local.project_id
  service = "secretmanager.googleapis.com"

  disable_dependent_services = true
  disable_on_destroy         = false
}

# A Secret Manager secret — the container for versioned secret payloads.
#
# One kind, two GCP API surfaces: an empty spec.region creates the GLOBAL
# trio (secret / version / iam_member) and a set region creates the REGIONAL
# trio. Exactly one branch of each pair is created — never both. The
# surfaces differ exactly where the spec does: replication is global-only,
# and the regional secret takes CMEK directly.
#
# An OMITTED spec.replication renders the API's `auto {}` mode — the
# provider REQUIRES a replication block on the global secret, and automatic
# placement is the right default when no residency regime applies.
resource "google_secret_manager_secret" "this" {
  count = local.is_regional ? 0 : 1

  secret_id = local.secret_id
  project   = local.project_id

  labels      = local.final_labels
  annotations = length(var.spec.annotations) > 0 ? var.spec.annotations : null
  tags        = length(var.spec.tags) > 0 ? var.spec.tags : null

  # Exactly one replication mode renders: user_managed when the spec pins
  # replica regions, auto otherwise (with its optional CMEK when the spec
  # set the arm explicitly).
  replication {
    dynamic "auto" {
      for_each = (var.spec.replication == null || var.spec.replication.auto != null) && (var.spec.replication == null || var.spec.replication.user_managed == null) ? [1] : []
      content {
        dynamic "customer_managed_encryption" {
          for_each = var.spec.replication != null && var.spec.replication.auto != null && var.spec.replication.auto.customer_managed_encryption != null ? [var.spec.replication.auto.customer_managed_encryption] : []
          content {
            kms_key_name = customer_managed_encryption.value.kms_key
          }
        }
      }
    }

    dynamic "user_managed" {
      for_each = var.spec.replication != null && var.spec.replication.user_managed != null ? [var.spec.replication.user_managed] : []
      content {
        dynamic "replicas" {
          for_each = user_managed.value.replicas
          content {
            location = replicas.value.location

            dynamic "customer_managed_encryption" {
              for_each = replicas.value.customer_managed_encryption != null ? [replicas.value.customer_managed_encryption] : []
              content {
                kms_key_name = customer_managed_encryption.value.kms_key
              }
            }
          }
        }
      }
    }
  }

  expire_time = var.spec.expire_time != "" ? var.spec.expire_time : null
  ttl         = var.spec.ttl != "" ? var.spec.ttl : null
  # GCP validates aliases against EXISTING versions at create/update, so an
  # alias cannot land in the same apply that seeds its version — add aliases
  # on a subsequent apply (live API: "Aliases cannot be assigned to versions
  # that don't exist").
  version_aliases     = length(var.spec.version_aliases) > 0 ? var.spec.version_aliases : null
  version_destroy_ttl = var.spec.version_destroy_ttl != "" ? var.spec.version_destroy_ttl : null

  dynamic "rotation" {
    for_each = var.spec.rotation != null ? [var.spec.rotation] : []
    content {
      rotation_period    = rotation.value.rotation_period != "" ? rotation.value.rotation_period : null
      next_rotation_time = rotation.value.next_rotation_time != "" ? rotation.value.next_rotation_time : null
    }
  }

  dynamic "topics" {
    for_each = var.spec.topics
    content {
      name = topics.value
    }
  }

  deletion_protection = var.spec.deletion_protection

  # Empty defers to the provider default (DELETE).
  deletion_policy = var.spec.deletion_policy != "" ? var.spec.deletion_policy : null

  depends_on = [google_project_service.secretmanager_api]
}

# The regional variant — payloads never leave spec.region; CMEK attaches
# directly instead of through replication.
resource "google_secret_manager_regional_secret" "this" {
  count = local.is_regional ? 1 : 0

  secret_id = local.secret_id
  location  = var.spec.region
  project   = local.project_id

  labels      = local.final_labels
  annotations = length(var.spec.annotations) > 0 ? var.spec.annotations : null
  tags        = length(var.spec.tags) > 0 ? var.spec.tags : null

  dynamic "customer_managed_encryption" {
    for_each = var.spec.customer_managed_encryption != null ? [var.spec.customer_managed_encryption] : []
    content {
      kms_key_name = customer_managed_encryption.value.kms_key
    }
  }

  expire_time = var.spec.expire_time != "" ? var.spec.expire_time : null
  ttl         = var.spec.ttl != "" ? var.spec.ttl : null
  # Same alias temporal constraint as the global variant above.
  version_aliases     = length(var.spec.version_aliases) > 0 ? var.spec.version_aliases : null
  version_destroy_ttl = var.spec.version_destroy_ttl != "" ? var.spec.version_destroy_ttl : null

  dynamic "rotation" {
    for_each = var.spec.rotation != null ? [var.spec.rotation] : []
    content {
      rotation_period    = rotation.value.rotation_period != "" ? rotation.value.rotation_period : null
      next_rotation_time = rotation.value.next_rotation_time != "" ? rotation.value.next_rotation_time : null
    }
  }

  dynamic "topics" {
    for_each = var.spec.topics
    content {
      name = topics.value
    }
  }

  deletion_protection = var.spec.deletion_protection

  deletion_policy = var.spec.deletion_policy != "" ? var.spec.deletion_policy : null

  depends_on = [google_project_service.secretmanager_api]
}

# The first payload, stored as version 1 — so one manifest yields a
# READABLE secret. enabled is sent explicitly (spec default true) so a
# true -> false transition reaches the API.
resource "google_secret_manager_secret_version" "initial" {
  count = !local.is_regional && var.spec.initial_version != null ? 1 : 0

  secret      = google_secret_manager_secret.this[0].id
  secret_data = var.spec.initial_version.data

  enabled               = var.spec.initial_version.enabled == null ? true : var.spec.initial_version.enabled
  is_secret_data_base64 = var.spec.initial_version.is_base64 ? true : null

  # Empty defers to the provider default (DELETE). DISABLE keeps the
  # payload recoverable; ABANDON leaves it untouched.
  deletion_policy = var.spec.initial_version.deletion_policy != "" ? var.spec.initial_version.deletion_policy : null
}

resource "google_secret_manager_regional_secret_version" "initial" {
  count = local.is_regional && var.spec.initial_version != null ? 1 : 0

  secret      = google_secret_manager_regional_secret.this[0].id
  secret_data = var.spec.initial_version.data

  enabled               = var.spec.initial_version.enabled == null ? true : var.spec.initial_version.enabled
  is_secret_data_base64 = var.spec.initial_version.is_base64 ? true : null

  deletion_policy = var.spec.initial_version.deletion_policy != "" ? var.spec.initial_version.deletion_policy : null
}

# Additive IAM grants: one (role, member) pair per resource, merging into
# the secret's policy without touching grants made elsewhere —
# authoritative bindings/policies are deliberately not used.
resource "google_secret_manager_secret_iam_member" "members" {
  for_each = local.is_regional ? {} : local.iam_members

  project   = local.project_id
  secret_id = google_secret_manager_secret.this[0].secret_id
  role      = each.value.role
  member    = each.value.member

  dynamic "condition" {
    for_each = each.value.condition != null ? [each.value.condition] : []
    content {
      title       = condition.value.title
      expression  = condition.value.expression
      description = condition.value.description != "" ? condition.value.description : null
    }
  }
}

resource "google_secret_manager_regional_secret_iam_member" "members" {
  for_each = local.is_regional ? local.iam_members : {}

  project   = local.project_id
  location  = var.spec.region
  secret_id = google_secret_manager_regional_secret.this[0].secret_id
  role      = each.value.role
  member    = each.value.member

  dynamic "condition" {
    for_each = each.value.condition != null ? [each.value.condition] : []
    content {
      title       = condition.value.title
      expression  = condition.value.expression
      description = condition.value.description != "" ? condition.value.description : null
    }
  }
}
