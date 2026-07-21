# Input variables for KubernetesStorageClass Terraform module.
# These mirror the KubernetesStorageClassSpec protobuf schema; enum fields
# arrive as the proto enum value names (e.g. "delete", "retain",
# "immediate", "wait_for_first_consumer").

variable "metadata" {
  description = "Metadata for the storage class resource"
  type = object({
    name = string
    id   = optional(string)
    org  = optional(string)
    env  = optional(string)
  })
}

variable "spec" {
  description = "Specification for the Kubernetes StorageClass"
  type = object({
    name        = string
    labels      = optional(map(string), {})
    annotations = optional(map(string), {})

    # The CSI driver (or in-tree plugin) provisioning volumes. Immutable.
    provisioner = string

    # Provisioner-specific parameters (each driver defines its own keys).
    # Immutable.
    parameters = optional(map(string), {})

    # "delete" (default) or "retain".
    reclaim_policy = optional(string, "delete")

    # "immediate" (default) or "wait_for_first_consumer".
    volume_binding_mode = optional(string, "immediate")

    allow_volume_expansion = optional(bool, false)

    mount_options = optional(list(string), [])

    # Zone/topology restrictions; only meaningful with
    # wait_for_first_consumer binding.
    allowed_topologies = optional(list(object({
      match_label_expressions = list(object({
        key    = string
        values = list(string)
      }))
    })), [])

    # Renders the storageclass.kubernetes.io/is-default-class annotation.
    is_default_class = optional(bool, false)
  })
}
