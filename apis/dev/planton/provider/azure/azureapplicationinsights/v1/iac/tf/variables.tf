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
  description = "Azure Application Insights specification"
  type = object({
    # The Azure region -- should match the monitored application's region
    region = string

    # The Azure Resource Group name (references resolved by the platform
    # before the module runs)
    resource_group = string

    # The Application Insights resource name (ForceNew)
    application_insights_name = string

    # The application type as the spec enum's value name (WEB / JAVA /
    # NODE_JS / OTHER / IOS / PHONE / STORE / MOBILE_CENTER). Absent
    # means WEB. ForceNew.
    application_type = optional(string)

    # The Log Analytics Workspace ARM id storing the telemetry
    # (references resolved by the platform). Repointable, never removable.
    workspace_id = string

    # Telemetry retention in days (one of Azure's fixed values)
    retention_in_days = optional(number, 90)

    # Daily telemetry cap in GB
    daily_data_cap_in_gb = optional(number, 100)

    # Notify by email when the daily cap is reached
    daily_data_cap_notifications_enabled = optional(bool, true)

    # Percentage of telemetry sampled (0-100)
    sampling_percentage = optional(number, 100)

    # Mask client IPs to 0.0.0.0 in stored telemetry (Azure's default)
    ip_masking_enabled = optional(bool, true)

    # Whether instrumentation-key-only ingestion works in addition to
    # Entra ID
    local_authentication_enabled = optional(bool, true)

    # Whether ingestion is accepted over the public internet
    internet_ingestion_enabled = optional(bool, true)

    # Whether queries are served over the public internet
    internet_query_enabled = optional(bool, true)

    # Force customer-owned storage for the .NET Profiler / Snapshot
    # Debugger artifacts
    force_customer_storage_for_profiler = optional(bool, false)

    # User tags, merged over the metadata-derived tags (user wins)
    tags = optional(map(string), {})
  })
}
