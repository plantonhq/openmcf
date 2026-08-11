variable "metadata" {
  description = "Cloud resource metadata"
  type = object({
    name = string
    id = optional(string, "")
    org = optional(string, "")
    env = optional(string, "")
    labels = optional(map(string), {})
    annotations = optional(map(string), {})
    tags = optional(list(string), [])
  })
}

variable "spec" {
  description = "AwsRoute53DnsRecord specification"
  type = object({
    region = string
    zone_id = string
    name = string
    type = string
    ttl = optional(number)
    values = optional(list(string), [])
    alias_target = optional(object({
      dns_name = string
      zone_id = string
      evaluate_target_health = optional(bool, false)
    }))
    routing_policy = optional(object({
      weighted = optional(object({
        weight = optional(number, 0)
      }))
      latency = optional(object({
        region = string
      }))
      failover = optional(object({
        failover_type = string
      }))
      geolocation = optional(object({
        continent = optional(string, "")
        country = optional(string, "")
        subdivision = optional(string, "")
      }))
      geoproximity = optional(object({
        aws_region = optional(string, "")
        coordinates = optional(object({
          latitude = string
          longitude = string
        }))
        local_zone_group = optional(string, "")
        bias = optional(number, 0)
      }))
      cidr = optional(object({
        collection_id = string
        location_name = string
      }))
      multivalue_answer = optional(object({}))
    }))
    health_check_id = optional(string, "")
    set_identifier = optional(string, "")
    allow_overwrite = optional(bool, false)
  })
}
