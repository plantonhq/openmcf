# Enable the Artifact Registry API — the control plane that owns
# repositories. disable_on_destroy is false: tearing down one repository
# must never disable the API for everything else in the project (other
# repositories keep serving pulls for running workloads).
resource "google_project_service" "artifactregistry_api" {
  project = local.project_id
  service = "artifactregistry.googleapis.com"

  disable_dependent_services = true
  disable_on_destroy         = false
}

# The Artifact Registry repository — one repository, one format, one serving
# mode. Sharp edges, all taught by the API rather than invented here:
#
#   - format, mode, location, project, and kms_key_name are ForceNew:
#     changing any of them replaces the repository AND everything stored in
#     it. There is no in-place migration between formats or modes.
#
#   - The remote_repository_config block is ForceNew as a whole, but the
#     upstream credentials and disable_upstream_validation inside it are in
#     the API's update mask — credential rotation updates in place.
#
#   - Deleting a repository deletes every artifact version in it. Unlike
#     GCS buckets there is no force_destroy gate; protect precious
#     repositories with KEEP cleanup policies and IAM, not with the module.
resource "google_artifact_registry_repository" "this" {
  project       = local.project_id
  repository_id = local.repository_id
  location      = local.location
  format        = var.spec.format
  mode          = local.mode
  description   = local.description
  labels        = local.final_labels

  # DELETE (provider default) removes the repository and every artifact
  # in it on destroy; PREVENT fails the destroy; ABANDON leaves it
  # serving artifacts.
  deletion_policy = var.spec.deletion_policy != "" ? var.spec.deletion_policy : null

  # CMEK: the Artifact Registry service agent must hold
  # roles/cloudkms.cryptoKeyEncrypterDecrypter on this key before create.
  kms_key_name = local.kms_key_name

  dynamic "docker_config" {
    for_each = var.spec.docker_config != null ? [var.spec.docker_config] : []
    content {
      immutable_tags = docker_config.value.immutable_tags
    }
  }

  dynamic "maven_config" {
    for_each = var.spec.maven_config != null ? [var.spec.maven_config] : []
    content {
      version_policy            = maven_config.value.version_policy != "" ? maven_config.value.version_policy : null
      allow_snapshot_overwrites = maven_config.value.allow_snapshot_overwrites
    }
  }

  # Cleanup: DELETE policies remove matching versions; KEEP policies protect
  # them (KEEP wins on overlap). Dry-run logs matches without deleting.
  cleanup_policy_dry_run = var.spec.cleanup_policy_dry_run

  dynamic "cleanup_policies" {
    for_each = var.spec.cleanup_policies
    content {
      id     = cleanup_policies.value.id
      action = cleanup_policies.value.action

      dynamic "condition" {
        for_each = cleanup_policies.value.condition != null ? [cleanup_policies.value.condition] : []
        content {
          newer_than            = condition.value.newer_than != "" ? condition.value.newer_than : null
          older_than            = condition.value.older_than != "" ? condition.value.older_than : null
          package_name_prefixes = length(condition.value.package_name_prefixes) > 0 ? condition.value.package_name_prefixes : null
          tag_prefixes          = length(condition.value.tag_prefixes) > 0 ? condition.value.tag_prefixes : null
          tag_state             = condition.value.tag_state != "" ? condition.value.tag_state : null
          version_name_prefixes = length(condition.value.version_name_prefixes) > 0 ? condition.value.version_name_prefixes : null
        }
      }

      dynamic "most_recent_versions" {
        for_each = cleanup_policies.value.most_recent_versions != null ? [cleanup_policies.value.most_recent_versions] : []
        content {
          keep_count            = most_recent_versions.value.keep_count > 0 ? most_recent_versions.value.keep_count : null
          package_name_prefixes = length(most_recent_versions.value.package_name_prefixes) > 0 ? most_recent_versions.value.package_name_prefixes : null
        }
      }
    }
  }

  # REMOTE_REPOSITORY: a pull-through cache of exactly one upstream. The
  # spec enforces mode↔config coherence and exactly-one-upstream pre-deploy.
  dynamic "remote_repository_config" {
    for_each = var.spec.remote_repository_config != null ? [var.spec.remote_repository_config] : []
    content {
      description                 = remote_repository_config.value.description != "" ? remote_repository_config.value.description : null
      disable_upstream_validation = remote_repository_config.value.disable_upstream_validation

      dynamic "docker_repository" {
        for_each = remote_repository_config.value.docker_public_repository != "" ? [remote_repository_config.value.docker_public_repository] : []
        content {
          public_repository = docker_repository.value
        }
      }

      dynamic "maven_repository" {
        for_each = remote_repository_config.value.maven_public_repository != "" ? [remote_repository_config.value.maven_public_repository] : []
        content {
          public_repository = maven_repository.value
        }
      }

      dynamic "npm_repository" {
        for_each = remote_repository_config.value.npm_public_repository != "" ? [remote_repository_config.value.npm_public_repository] : []
        content {
          public_repository = npm_repository.value
        }
      }

      dynamic "python_repository" {
        for_each = remote_repository_config.value.python_public_repository != "" ? [remote_repository_config.value.python_public_repository] : []
        content {
          public_repository = python_repository.value
        }
      }

      dynamic "apt_repository" {
        for_each = remote_repository_config.value.apt_repository != null ? [remote_repository_config.value.apt_repository] : []
        content {
          public_repository {
            repository_base = apt_repository.value.repository_base
            repository_path = apt_repository.value.repository_path
          }
        }
      }

      dynamic "yum_repository" {
        for_each = remote_repository_config.value.yum_repository != null ? [remote_repository_config.value.yum_repository] : []
        content {
          public_repository {
            repository_base = yum_repository.value.repository_base
            repository_path = yum_repository.value.repository_path
          }
        }
      }

      # Custom upstream: another AR repository or any registry URI.
      dynamic "common_repository" {
        for_each = remote_repository_config.value.common_repository != null ? [remote_repository_config.value.common_repository] : []
        content {
          uri = common_repository.value.uri
        }
      }

      # Credential rotation updates in place — the password itself lives in
      # Secret Manager; only the secret-version PATH passes through here.
      dynamic "upstream_credentials" {
        for_each = remote_repository_config.value.upstream_credentials != null ? [remote_repository_config.value.upstream_credentials] : []
        content {
          username_password_credentials {
            username                = upstream_credentials.value.username
            password_secret_version = upstream_credentials.value.password_secret_version
          }
        }
      }
    }
  }

  # VIRTUAL_REPOSITORY: priority-ordered aggregation of other AR
  # repositories (highest priority wins on conflicts).
  dynamic "virtual_repository_config" {
    for_each = var.spec.virtual_repository_config != null ? [var.spec.virtual_repository_config] : []
    content {
      dynamic "upstream_policies" {
        for_each = virtual_repository_config.value.upstream_policies
        content {
          id         = upstream_policies.value.id
          repository = upstream_policies.value.repository
          priority   = upstream_policies.value.priority
        }
      }
    }
  }

  dynamic "vulnerability_scanning_config" {
    for_each = local.vulnerability_scanning_enablement != null ? [local.vulnerability_scanning_enablement] : []
    content {
      enablement_config = vulnerability_scanning_config.value
    }
  }

  depends_on = [
    google_project_service.artifactregistry_api,
  ]
}

# Additive IAM grants: one (role, member) pair per resource, merging into
# the repository's policy without touching grants made elsewhere —
# authoritative bindings/policies are deliberately not used.
resource "google_artifact_registry_repository_iam_member" "members" {
  for_each = local.iam_members

  project    = local.project_id
  location   = local.location
  repository = google_artifact_registry_repository.this.repository_id
  role       = each.value.role
  member     = each.value.member

  dynamic "condition" {
    for_each = each.value.condition != null ? [each.value.condition] : []
    content {
      title       = condition.value.title
      expression  = condition.value.expression
      description = condition.value.description != "" ? condition.value.description : null
    }
  }
}
