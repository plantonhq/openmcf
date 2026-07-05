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
  description = "Specification for the GCP Serverless VPC Access connector"
  type = object({
    project_id     = optional(string, "")
    region         = string
    connector_name = optional(string, "")

    # Network placement: carve a new /28 out of the VPC. Exactly one of
    # network(+ip_cidr_range) or subnet is set — enforced pre-deploy.
    network       = optional(string, "")
    ip_cidr_range = optional(string, "")

    # Subnet placement: occupy an existing dedicated /28 subnetwork
    # (the Shared-VPC-capable mode).
    subnet = optional(object({
      name       = string
      project_id = optional(string, "")
    }), null)

    machine_type  = optional(string, "")
    min_instances = optional(number, null)
    max_instances = optional(number, null)
  })
}
