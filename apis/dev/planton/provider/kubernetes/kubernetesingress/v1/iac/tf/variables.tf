# Input variables for KubernetesIngress Terraform module.
# These mirror the KubernetesIngressSpec protobuf schema; StringValueOrRef
# fields (namespace, backend service names, TLS secret names) arrive flattened
# to plain strings, and enum fields arrive as the proto enum value names
# (e.g. "prefix", "exact", "implementation_specific").

variable "metadata" {
  description = "Metadata for the ingress resource"
  type = object({
    name = string
    id   = optional(string)
    org  = optional(string)
    env  = optional(string)
  })
}

variable "spec" {
  description = "Specification for the Kubernetes Ingress"
  type = object({
    namespace   = optional(string, "default")
    name        = string
    labels      = optional(map(string), {})
    annotations = optional(map(string), {})

    # The IngressClass selecting which controller serves this Ingress; ""
    # defers to the cluster's default class.
    ingress_class_name = optional(string, "")

    # Backend for requests no rule matches. Exactly one of port_number /
    # port_name is set (spec CEL contract).
    default_backend = optional(object({
      service_name = string
      port_number  = optional(number, 0)
      port_name    = optional(string, "")
    }))

    tls = optional(list(object({
      hosts       = optional(list(string), [])
      secret_name = optional(string, "")
    })), [])

    rules = optional(list(object({
      host = optional(string, "")
      paths = list(object({
        path = optional(string, "")
        # prefix / exact / implementation_specific; null defaults to prefix.
        path_type = optional(string)
        backend = object({
          service_name = string
          port_number  = optional(number, 0)
          port_name    = optional(string, "")
        })
      }))
    })), [])
  })
}
