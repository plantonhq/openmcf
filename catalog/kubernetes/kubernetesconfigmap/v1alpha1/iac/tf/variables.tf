# Input variables for Kubernetes ConfigMap Terraform module

variable "metadata" {
  description = "Metadata for the configmap resource"
  type = object({
    name = string
    org  = optional(string)
    env  = optional(string)
  })
}

variable "spec" {
  description = "Specification for the Kubernetes ConfigMap"
  type = object({
    name        = string
    namespace   = optional(string, "default")
    labels      = optional(map(string), {})
    annotations = optional(map(string), {})
    immutable   = optional(bool, false)

    # UTF-8 configuration entries. Keys become file names (volume mounts) or
    # environment variable names (envFrom).
    data = optional(map(string), {})

    # Binary configuration entries. Values are already base64-encoded strings —
    # the exact wire form Kubernetes uses for binaryData — and pass through unchanged.
    binary_data = optional(map(string), {})
  })
}
