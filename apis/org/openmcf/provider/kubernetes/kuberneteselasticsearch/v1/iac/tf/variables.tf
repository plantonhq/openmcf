variable "metadata" {
  description = "Metadata for the resource, including name and labels"
  type = object({
    name = string,
    id = optional(string),
    org = optional(string),
    env = optional(string),
    labels = optional(map(string)),
    tags = optional(list(string)),
    version = optional(object({ id = string, message = string }))
  })
}


variable "spec" {
  description = "spec"
  type = object({
    # Kubernetes namespace
    namespace = string

    # Flag to indicate if the namespace should be created
    create_namespace = bool

    # Elasticsearch configuration
    elasticsearch = object({
      # Container configuration
      container = object({
        replicas             = number
        persistence_enabled  = bool
        disk_size            = string
        resources = object({
          limits = object({
            cpu    = string
            memory = string
          })
          requests = object({
            cpu    = string
            memory = string
          })
        })
      })

      # Ingress configuration
      ingress = object({
        enabled  = bool
        hostname = string
      })
    })

    # Kibana configuration
    kibana = object({
      enabled = bool
      container = object({
        replicas = number
        resources = object({
          limits = object({
            cpu    = string
            memory = string
          })
          requests = object({
            cpu    = string
            memory = string
          })
        })
      })
      ingress = object({
        enabled  = bool
        hostname = string
      })
    })
  })
}