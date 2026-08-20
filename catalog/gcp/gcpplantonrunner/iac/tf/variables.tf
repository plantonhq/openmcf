variable "metadata" {
  description = "Cloud resource metadata"
  type = object({
    name        = string
    id          = optional(string, "")
    org         = optional(string, "")
    env         = optional(string, "")
    labels      = optional(map(string), {})
    annotations = optional(map(string), {})
    tags        = optional(list(string), [])
  })
}

variable "spec" {
  description = "GcpPlantonRunner specification"
  type = object({
    project_id             = optional(string, "")
    region                 = string
    token                  = string
    control_plane_endpoint = optional(string, "")
    runner_version         = optional(string)
    image_repository       = optional(string)
    service_account        = optional(string, "")
    vpc_access = optional(object({
      network    = string
      subnetwork = string
      tags       = optional(list(string), [])
    }))
    cpu    = optional(string)
    memory = optional(string)
  })
}