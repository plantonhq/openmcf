# Input variables for KubernetesPersistentVolumeClaim Terraform module.
# These mirror the KubernetesPersistentVolumeClaimSpec protobuf schema; the
# namespace and storage_class_name StringValueOrRef fields arrive flattened
# to plain strings, and enum fields arrive as the proto enum value names
# (e.g. "filesystem", "block").

variable "metadata" {
  description = "Metadata for the persistent volume claim resource"
  type = object({
    name = string
    id   = optional(string)
    org  = optional(string)
    env  = optional(string)
  })
}

variable "spec" {
  description = "Specification for the Kubernetes PersistentVolumeClaim"
  type = object({
    namespace   = optional(string, "default")
    name        = string
    labels      = optional(map(string), {})
    annotations = optional(map(string), {})

    # ReadWriteOnce / ReadOnlyMany / ReadWriteMany / ReadWriteOncePod.
    # Empty defaults to ["ReadWriteOnce"].
    access_modes = optional(list(string), [])

    storage_request = string
    storage_limit   = optional(string, "")

    # The StorageClass name; empty/null means the cluster default applies.
    storage_class_name = optional(string)

    # Pins storageClassName to "" (bind only pre-provisioned volumes) —
    # distinct from omitting the class (cluster default).
    disable_dynamic_provisioning = optional(bool, false)

    # "filesystem" (default) or "block".
    volume_mode = optional(string, "filesystem")

    # Static binding to one specific PersistentVolume.
    volume_name = optional(string, "")

    # Narrows which pre-provisioned volumes may bind.
    selector = optional(object({
      match_labels = optional(map(string), {})
      match_expressions = optional(list(object({
        key      = string
        operator = string
        values   = optional(list(string), [])
      })), [])
    }))

    # Populate from a clone (persistent_volume_claim) or a snapshot restore
    # (volume_snapshot).
    data_source = optional(object({
      kind = string
      name = string
    }))
  })
}
