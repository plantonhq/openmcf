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
  description = "Azure Event Hubs Dedicated Cluster specification"
  type = object({
    # The Azure region for the cluster (e.g. "eastus"). ForceNew.
    region = string

    # The resource group name. References are resolved to a literal by
    # the platform before the module runs. ForceNew.
    resource_group = string

    # The cluster's name -- unique within the resource group. ForceNew.
    cluster_name = string

    # The cluster's size in capacity units (CUs). Scales in place;
    # unset deploys 1 CU -- Azure's entry size.
    capacity_units = optional(number)

    # User tags, merged over the Planton-derived identity tags (user
    # values win on key conflicts).
    tags = optional(map(string), {})
  })
}
