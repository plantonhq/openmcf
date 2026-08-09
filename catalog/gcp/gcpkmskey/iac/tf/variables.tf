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
  description = "Specification for the GCP KMS crypto key"
  type = object({
    # StringValueOrRef fields arrive from the proto→tfvars converter as
    # plain strings (already resolved), never as object({value}).
    # The fully qualified key ring path
    # (projects/{project}/locations/{location}/keyRings/{name}) — the key
    # inherits its project and location from the ring, so this module
    # needs no project variable of its own.
    key_ring_id = string

    # Key name (the GCP resource name). Immutable (ForceNew) — and since
    # crypto keys can never be deleted from GCP, a name is permanently
    # consumed within its ring once used.
    key_name = string

    # Key purpose (e.g. ENCRYPT_DECRYPT). Immutable (ForceNew). Empty
    # defers to GCP's default, ENCRYPT_DECRYPT.
    purpose = optional(string, "")

    # Automatic rotation cadence (e.g. "7776000s"). Empty disables
    # automatic rotation. Mutable in place.
    rotation_period = optional(string, "")

    # DESTROY_SCHEDULED recovery window for destroyed versions
    # (e.g. "2592000s"). Immutable (ForceNew). Empty defers to GCP's
    # default of 30 days.
    destroy_scheduled_duration = optional(string, "")

    version_template = optional(object({
      # Algorithm for new versions. Mutable — affects future versions only.
      algorithm = string
      # SOFTWARE / HSM / EXTERNAL / EXTERNAL_VPC. Immutable (ForceNew).
      protection_level = optional(string, "")
    }), null)

    # Create the key without an initial version (required for import_only
    # keys). Consumed only at create time.
    skip_initial_version_creation = optional(bool, false)

    # BYOK container: the key may only ever hold imported versions.
    # Immutable (ForceNew).
    import_only = optional(bool, false)

    # EKM connection path backing EXTERNAL_VPC keys (resolved from a
    # reference). Immutable (ForceNew). Empty for all other protection
    # levels.
    crypto_key_backend = optional(string, "")

    # User labels, merged beneath the platform attribution labels
    # (see locals.tf).
    labels = optional(map(string), {})

    # DELETE (default) destroys every key version on destroy — data
    # encrypted under them becomes unrecoverable; PREVENT fails the
    # destroy; ABANDON leaves the key and all versions intact.
    deletion_policy = optional(string, "")
  })

  validation {
    condition     = var.spec.key_ring_id != ""
    error_message = "key_ring_id is required."
  }

  validation {
    condition     = var.spec.key_name != ""
    error_message = "key_name is required."
  }
}
