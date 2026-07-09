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
  description = "Azure Container App Environment Dapr component specification"
  type = object({
    # The Container App Environment ARM ID. ForceNew.
    container_app_environment_id = string

    # The Dapr component name application code passes to the Dapr API
    # (max 60 lowercase alphanumerics/hyphens, starts with a letter).
    # ForceNew.
    component_name = string

    # The Dapr component type in dotted notation (e.g.
    # state.azure.blobstorage, pubsub.azure.servicebus). ForceNew.
    component_type = string

    # The component version ("v1" for virtually all stable components).
    version = string

    # Sidecar initialisation timeout in whole seconds/minutes/hours
    # (documented default applied here because the platform never
    # materializes proto defaults).
    init_timeout = optional(string, "5s")

    # Whether the sidecar continues initialising when this component
    # fails to load.
    ignore_errors = optional(bool, false)

    # Component secrets referenced from metadata by secret_name.
    secrets = optional(list(object({
      name  = string
      value = string
    })), [])

    # Configuration entries: literal value XOR secret_name per entry
    # (spec-enforced).
    metadata = optional(list(object({
      name        = string
      value       = optional(string)
      secret_name = optional(string)
    })), [])

    # The dapr.app_id values allowed to use this component; empty exposes
    # it to every Dapr-enabled app in the environment.
    scopes = optional(list(string), [])
  })
}
