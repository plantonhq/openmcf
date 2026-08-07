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
  description = "Azure Front Door endpoint specification"
  type = object({
    # The Front Door profile the endpoint lives in, by ARM ID.
    # References are resolved to a literal ID by the platform before the
    # module runs. ForceNew.
    profile_id = string

    # The endpoint's name -- unique within the profile and the prefix of
    # the generated *.azurefd.net hostname. ForceNew.
    endpoint_name = string

    # Whether the endpoint accepts traffic. Azure defaults to true when
    # omitted (tfvars drops zero-valued proto fields; the platform
    # materializes the documented default centrally).
    enabled = optional(bool)

    # Free-form user tags merged over the metadata-derived tags (user
    # tags win on collision).
    tags = optional(map(string), {})
  })
}
