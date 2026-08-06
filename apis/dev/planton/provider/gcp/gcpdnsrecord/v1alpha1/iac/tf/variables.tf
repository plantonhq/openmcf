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
  description = "Specification for the GCP DNS record set"
  type = object({
    # StringValueOrRef fields arrive as PLAIN STRINGS: the tfvars converter
    # flattens refs before the module ever sees them.
    project_id   = optional(string, "")
    managed_zone = string

    type = string
    name = string

    # Static RRDATA — mutually exclusive with routing_policy (enforced by
    # the spec's pre-deploy validation, mirrored by the provider's
    # ExactlyOneOf on rrdatas/routing_policy).
    values = optional(list(string), [])

    ttl_seconds = optional(number, 300)

    routing_policy = optional(object({
      wrr = optional(list(object({
        weight = number
        values = optional(list(string), [])
        health_checked_targets = optional(object({
          internal_load_balancers = optional(list(object({
            ip_address         = string
            ip_protocol        = string
            load_balancer_type = optional(string, "")
            network_url        = string
            port               = string
            project            = string
            region             = optional(string, "")
          })), [])
          external_endpoints = optional(list(string), [])
        }), null)
      })), [])

      geo = optional(list(object({
        location = string
        values   = optional(list(string), [])
        health_checked_targets = optional(object({
          internal_load_balancers = optional(list(object({
            ip_address         = string
            ip_protocol        = string
            load_balancer_type = optional(string, "")
            network_url        = string
            port               = string
            project            = string
            region             = optional(string, "")
          })), [])
          external_endpoints = optional(list(string), [])
        }), null)
      })), [])

      enable_geo_fencing = optional(bool, false)

      primary_backup = optional(object({
        primary = object({
          internal_load_balancers = optional(list(object({
            ip_address         = string
            ip_protocol        = string
            load_balancer_type = optional(string, "")
            network_url        = string
            port               = string
            project            = string
            region             = optional(string, "")
          })), [])
          external_endpoints = optional(list(string), [])
        })
        backup_geo = list(object({
          location = string
          values   = optional(list(string), [])
          health_checked_targets = optional(object({
            internal_load_balancers = optional(list(object({
              ip_address         = string
              ip_protocol        = string
              load_balancer_type = optional(string, "")
              network_url        = string
              port               = string
              project            = string
              region             = optional(string, "")
            })), [])
            external_endpoints = optional(list(string), [])
          }), null)
        }))
        trickle_ratio                  = optional(number)
        enable_geo_fencing_for_backups = optional(bool, false)
      }), null)

      health_check = optional(string, "")
    }), null)
  })
}
