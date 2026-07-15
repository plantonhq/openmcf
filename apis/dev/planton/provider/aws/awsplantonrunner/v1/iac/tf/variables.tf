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
  description = "AwsPlantonRunner specification"
  type = object({
    region = string

    # Foreign-key fields arrive pre-resolved as plain strings (the
    # orchestrator resolves any value_from before the module runs).
    subnets = list(string)
    security_groups = optional(list(string), [])
    assign_public_ip = optional(bool, false)

    # Module-level defaults mirror the spec's platform defaults so the
    # tfvars path (which prunes unset fields) lands on the same values
    # the Pulumi module receives.
    cpu = optional(number, 512)
    memory = optional(number, 1024)
    runner_version = optional(string, "latest")
    image_repository = optional(string, "ghcr.io/plantonhq/planton/runner")
    execution_mode = optional(string, "temporal")

    # The runner's credentials document -- resolved from its managed-secret
    # reference before the module runs; it only ever leaves this module
    # into Secrets Manager, never into the task definition.
    credentials = string

    task_role = optional(string, "")
    log_retention_days = optional(number, 30)
  })
}
