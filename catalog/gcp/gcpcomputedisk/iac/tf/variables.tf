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

    # At most one source (image / source_snapshot / source_instant_snapshot /
    # source_storage_object / source_disk) — enforced pre-deploy by the
    # spec's CEL. All create-time only.
    image                   = optional(string, "")
    source_snapshot         = optional(string, "")
    source_instant_snapshot = optional(string, "")
    source_storage_object   = optional(string, "")
    source_disk             = optional(string, "")

    # Resolved CMEK key id. Immutable.
    kms_key = optional(string, "")

    # Service account for the CMEK encryption request; empty uses the
    # Compute Engine default service agent. Immutable.
    kms_key_service_account = optional(string, "")

    # CMEK decryption parameters for encrypted sources (raw CSEK keys are
    # deliberately not part of the contract — key material never flows
    # through manifests or state).
    source_image_encryption = optional(object({
      kms_key                 = string
      kms_key_service_account = optional(string, "")
    }), null)
    source_snapshot_encryption = optional(object({
      kms_key                 = string
      kms_key_service_account = optional(string, "")
    }), null)

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

    # Guest OS features for bootable disks (e.g. UEFI_COMPATIBLE,
    # SECURE_BOOT, GVNIC). Create-time only.
    guest_os_features = optional(list(string), [])

    # License URIs (bring-your-own-license imports). Create-time only.
    licenses = optional(list(string), [])

    # Destroy-time guard: "" / "DELETE" deletes, "PREVENT" fails the
    # destroy, "ABANDON" unmanages without deleting.
    deletion_policy = optional(string, "")

    # Resolved self link of the async replication PRIMARY disk; makes
    # this disk the replication secondary. Create-time only.
    async_primary_disk = optional(string, "")
  })
}
