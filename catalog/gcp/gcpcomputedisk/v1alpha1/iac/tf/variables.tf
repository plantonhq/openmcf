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
  description = "Specification for the Compute Engine persistent disk"
  type = object({
    # StringValueOrRef fields arrive from the proto→tfvars converter as
    # plain strings (already resolved), never as object({value}).
    # Empty falls back to the provider's default project.
    project_id = optional(string, "")

    # Disk name; empty falls back to metadata.name. Immutable.
    disk_name = optional(string, "")

    # Zone, e.g. "us-central1-a". Immutable.
    zone = string

    description = optional(string, "")

    # pd-standard / pd-balanced / pd-ssd / pd-extreme / hyperdisk-*.
    # Immutable.
    type = optional(string, "")

    # Size in GB. Required for empty disks; grows in place, never
    # shrinks.
    size_gb = optional(number)

    # At most one source (image / source_snapshot / source_disk) —
    # enforced pre-deploy by the spec's CEL. All create-time only.
    image           = optional(string, "")
    source_snapshot = optional(string, "")
    source_disk     = optional(string, "")

    # Resolved CMEK key id. Immutable.
    kms_key = optional(string, "")

    # Hyperdisk tuning; null for pd-* types (pd-extreme takes iops).
    provisioned_iops       = optional(number)
    provisioned_throughput = optional(number)

    access_mode  = optional(string, "")
    architecture = optional(string, "")

    # Confidential-compute disk (hyperdisk SKUs; requires kms_key).
    enable_confidential_compute = optional(bool, false)

    physical_block_size_bytes = optional(number)

    # Last-resort recovery net for precious volumes.
    create_snapshot_before_destroy = optional(bool, false)
    snapshot_before_destroy_prefix = optional(string, "")

    storage_pool = optional(string, "")

    # User labels merged beneath the platform attribution labels.
    labels = optional(map(string), {})

    # Resource Manager tags (tagKeys/{id} => tagValues/{id}). Create-time.
    resource_manager_tags = optional(map(string), {})
  })
}
