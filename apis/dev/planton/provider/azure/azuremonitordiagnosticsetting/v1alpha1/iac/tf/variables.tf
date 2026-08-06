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
  description = "Azure Monitor diagnostic setting specification"
  type = object({
    # The setting name, unique among the target's settings (ForceNew)
    setting_name = string

    # The ARM id of the resource whose telemetry is routed (references
    # resolved by the platform before the module runs). ForceNew.
    target_resource_id = string

    # Log selections: each entry names exactly one of category or
    # category_group (spec-enforced XOR).
    enabled_logs = optional(list(object({
      category       = optional(string)
      category_group = optional(string)
    })), [])

    # Metric selections (most resource types expose "AllMetrics")
    enabled_metrics = optional(list(object({
      category = string
    })), [])

    # Destination: Log Analytics Workspace ARM id (references resolved
    # by the platform)
    log_analytics_workspace_id = optional(string)

    # Workspace table layout as the spec enum's value name
    # (DEDICATED / AZURE_DIAGNOSTICS); absent lets Azure pick.
    log_analytics_destination_type = optional(string)

    # Destination: Storage Account ARM id (archival)
    storage_account_id = optional(string)

    # Destination: Event Hub namespace authorization rule ARM id
    # (streaming); eventhub_name optionally picks a specific hub.
    eventhub_authorization_rule_id = optional(string)
    eventhub_name                  = optional(string)

    # Destination: Azure Native ISV partner solution ARM id
    partner_solution_id = optional(string)
  })
}
