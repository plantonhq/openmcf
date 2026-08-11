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
  description = "Specification for the GCP AlloyDB user"
  type = object({
    project_id = optional(string, "")

    cluster = string
    user_id = string

    user_type = optional(string, "ALLOYDB_BUILT_IN")
    password  = optional(string, "")

    database_roles = optional(list(string), [])

    # What happens to the user in GCP on destroy:
    # DELETE (provider default), PREVENT, or ABANDON.
    deletion_policy = optional(string, "")
  })
}
