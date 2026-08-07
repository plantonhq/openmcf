# Input variables for Kubernetes Secret Terraform module

variable "metadata" {
  description = "Metadata for the secret resource"
  type = object({
    name = string
    org  = optional(string)
    env  = optional(string)
  })
}

variable "spec" {
  description = "Specification for the Kubernetes Secret"
  type = object({
    name        = string
    namespace   = optional(string, "default")
    labels      = optional(map(string), {})
    annotations = optional(map(string), {})
    immutable   = optional(bool, false)

    # Exactly one of the following should be provided.
    # The secret type is determined by which block is set.

    opaque = optional(object({
      # data carries plain-string values; binary_data carries base64-encoded values.
      # At least one entry across the two maps is required by the spec.
      data        = optional(map(string), {})
      binary_data = optional(map(string), {})
    }))

    tls = optional(object({
      tls_crt = string
      tls_key = string
    }))

    docker_config_json = optional(object({
      registry_server = string
      username        = string
      password        = string
      email           = optional(string, "")
    }))

    basic_auth = optional(object({
      username = string
      password = string
    }))

    ssh_auth = optional(object({
      ssh_private_key = string
    }))

    service_account_token = optional(object({
      service_account_name = string
    }))
  })
}
