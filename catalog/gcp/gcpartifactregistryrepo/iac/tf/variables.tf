variable "metadata" {
  description = "Metadata for the resource, including name and labels"
  type = object({
    name    = string,
    id      = optional(string),
    org     = optional(string),
    env     = optional(string),
    labels  = optional(map(string)),
    tags    = optional(list(string)),
    version = optional(object({ id = string, message = string }))
  })
}

variable "spec" {
  description = "Specification for the GCP Artifact Registry repository"
  type = object({
    # StringValueOrRef fields arrive from the proto→tfvars converter as
    # plain strings (already resolved), never as object({value}).
    # Empty falls back to the provider's default project.
    project_id = optional(string, "")

    # The last segment of the repository resource name (appears in registry
    # URLs). Empty falls back to metadata.name. Immutable (ForceNew).
    repository_id = optional(string, "")

    # Region ("us-central1") or multi-region ("us", "europe", "asia").
    # Immutable (ForceNew).
    location = string

    # Package format: DOCKER, MAVEN, NPM, PYTHON, GO, APT, YUM, GENERIC,
    # KFP, ... (free string — GCP matches case-insensitively and adds
    # formats over time). Immutable (ForceNew).
    format = string

    # Serving mode: STANDARD_REPOSITORY (default when empty),
    # REMOTE_REPOSITORY, or VIRTUAL_REPOSITORY. Immutable (ForceNew).
    mode = optional(string, "")

    # Human-readable description. Mutable in place.
    description = optional(string, "")

    # User labels, merged beneath the platform attribution labels
    # (see locals.tf). Mutable in place.
    labels = optional(map(string), {})

    # CMEK crypto key path (resolved from a GcpKmsKey reference).
    # Immutable (ForceNew). Empty means Google-managed encryption.
    kms_key_name = optional(string, "")

    # Docker-format settings. Mutable in place.
    docker_config = optional(object({
      immutable_tags = optional(bool, false)
    }), null)

    # Maven-format settings. Effectively immutable: GCP rejects changing
    # them on an existing repository.
    maven_config = optional(object({
      version_policy            = optional(string, "")
      allow_snapshot_overwrites = optional(bool, false)
    }), null)

    # Automatic version cleanup. KEEP policies always win over DELETE
    # policies on overlap. Mutable in place.
    cleanup_policies = optional(list(object({
      id     = string
      action = string
      condition = optional(object({
        newer_than            = optional(string, "")
        older_than            = optional(string, "")
        package_name_prefixes = optional(list(string), [])
        tag_prefixes          = optional(list(string), [])
        tag_state             = optional(string, "")
        version_name_prefixes = optional(list(string), [])
      }), null)
      most_recent_versions = optional(object({
        keep_count            = optional(number, 0)
        package_name_prefixes = optional(list(string), [])
      }), null)
    })), [])

    # Log-only mode for cleanup policies (validate before deleting).
    cleanup_policy_dry_run = optional(bool, false)

    # Upstream source for REMOTE_REPOSITORY mode. The block is immutable
    # EXCEPT upstream_credentials and disable_upstream_validation, which
    # rotate in place.
    remote_repository_config = optional(object({
      description                 = optional(string, "")
      docker_public_repository    = optional(string, "")
      maven_public_repository     = optional(string, "")
      npm_public_repository       = optional(string, "")
      python_public_repository    = optional(string, "")
      disable_upstream_validation = optional(bool, false)
      apt_repository = optional(object({
        repository_base = string
        repository_path = string
      }), null)
      yum_repository = optional(object({
        repository_base = string
        repository_path = string
      }), null)
      common_repository = optional(object({
        # Resolved from a GcpArtifactRegistryRepo reference or a literal
        # registry URI.
        uri = string
      }), null)
      upstream_credentials = optional(object({
        username = string
        # Secret Manager secret version path — a reference to the secret,
        # never the secret material itself.
        password_secret_version = string
      }), null)
    }), null)

    # Upstream aggregation for VIRTUAL_REPOSITORY mode. Mutable in place.
    virtual_repository_config = optional(object({
      upstream_policies = list(object({
        id = string
        # Resolved from a GcpArtifactRegistryRepo reference — the full
        # projects/{p}/locations/{l}/repositories/{r} path.
        repository = string
        priority   = optional(number, 0)
      }))
    }), null)

    # Artifact Analysis scanning: "" or "INHERITED" follow the project
    # setting; "DISABLED" opts this repository out. Mutable in place.
    vulnerability_scanning_enablement = optional(string, "")

    # Additive IAM grants: one role to one member each, composing safely
    # with grants made elsewhere.
    iam_members = optional(list(object({
      role = string
      # Resolved from a GcpServiceAccount reference or a literal IAM
      # member string (serviceAccount:..., user:..., allUsers, ...).
      member = string
      condition = optional(object({
        title       = string
        expression  = string
        description = optional(string, "")
      }), null)
    })), [])

    # DELETE (default) removes the repository and every artifact in it on
    # destroy; PREVENT fails the destroy; ABANDON leaves the repository
    # serving artifacts in GCP.
    deletion_policy = optional(string, "")
  })

  validation {
    condition     = var.spec.location != ""
    error_message = "location is required."
  }

  validation {
    condition     = var.spec.format != ""
    error_message = "format is required."
  }
}
