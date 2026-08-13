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
  description = "Specification for the GCP Bigtable instance"
  type = object({
    # The GCP project for the instance. The CLI's tfvars converter
    # resolves StringValueOrRef fields to their literal string before the
    # module runs, so this arrives as a plain string.
    # If empty, the provider's default project is used (see locals.tf).
    project_id = optional(string, "")

    # Instance ID (GCP resource name), 6-33 chars. Immutable (ForceNew).
    instance_name = string

    # Human-readable display name; empty defaults to instance_name
    # server-side.
    display_name = optional(string, "")

    # Terraform-side deletion guard. The spec defaults this to true
    # (Planton middleware materializes the default); destroy fails until
    # it is explicitly set false.
    deletion_protection = optional(bool, true)

    # Delete all backups in the instance on destroy — Bigtable otherwise
    # blocks instance deletion while backups exist.
    force_destroy = optional(bool, false)

    # Physical replicas. Each cluster's zone, storage_type, kms_key_name,
    # and node_scaling_factor are immutable server-side.
    clusters = list(object({
      cluster_id          = string
      zone                = string
      num_nodes           = optional(number, 0)
      storage_type        = optional(string, "SSD")
      kms_key_name        = optional(string, "")
      node_scaling_factor = optional(string, "")
      autoscaling_config = optional(object({
        min_nodes      = number
        max_nodes      = number
        cpu_target     = number
        storage_target = optional(number, 0)
      }), null)
    }))

    # User labels merged beneath Planton platform labels (platform keys
    # win on conflict).
    labels = optional(map(string), {})

    # ENTERPRISE (GCP default) or ENTERPRISE_PLUS; upgrade in place only.
    edition = optional(string, "")

    # Resource Manager tags (tagKeys/{id} -> tagValues/{id}); ForceNew.
    resource_manager_tags = optional(map(string), {})

    # Client-side destroy behavior: DELETE (default), PREVENT, ABANDON.
    deletion_policy = optional(string, "")
  })

  validation {
    condition     = var.spec.instance_name != ""
    error_message = "instance_name is required."
  }

  validation {
    condition     = length(var.spec.clusters) > 0
    error_message = "at least one cluster is required."
  }
}
