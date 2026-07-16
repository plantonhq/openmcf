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
  description = "Azure Log Analytics Workspace specification"
  type = object({
    # The Azure region where the workspace will be deployed
    region = string

    # The Azure Resource Group name (references resolved by the platform
    # before the module runs)
    resource_group = string

    # The workspace name (4-63 letters/digits/hyphens; ForceNew)
    workspace_name = string

    # The pricing tier as the spec enum's value name (PER_GB_2018 /
    # CAPACITY_RESERVATION / PER_NODE / STANDALONE). Absent means
    # PER_GB_2018 -- Azure's recommended pay-as-you-go tier.
    sku = optional(string)

    # The commitment tier in GB/day; only with CAPACITY_RESERVATION
    # (spec-enforced pairing).
    reservation_capacity_in_gb_per_day = optional(number)

    # Workspace-level data retention in days (30-730)
    retention_in_days = optional(number, 30)

    # Daily ingestion cap in GB; -1 means unlimited
    daily_quota_gb = optional(number, -1)

    # Managed identity: type is the spec enum's value name
    # (SYSTEM_ASSIGNED / USER_ASSIGNED / SYSTEM_AND_USER_ASSIGNED);
    # user_assigned_identity_ids carry resolved ARM ids.
    identity = optional(object({
      type                       = string
      user_assigned_identity_ids = optional(list(string), [])
    }))

    # Whether workspace shared keys work in addition to Entra ID
    local_authentication_enabled = optional(bool, true)

    # Whether ingestion is accepted over the public internet
    internet_ingestion_enabled = optional(bool, true)

    # Whether queries are served over the public internet
    internet_query_enabled = optional(bool, true)

    # Resource-context query access (Azure's modern default) vs
    # workspace-context-only when false
    allow_resource_only_permissions = optional(bool, true)

    # Force customer-managed storage for query artifacts
    cmk_for_query_forced = optional(bool, false)

    # Purge data immediately after 30 days (no grace store)
    immediate_data_purge_on_30_days_enabled = optional(bool, false)

    # Default Data Collection Rule ARM id (literal; no DCR kind exists)
    data_collection_rule_id = optional(string)

    # User tags, merged over the metadata-derived tags (user wins)
    tags = optional(map(string), {})
  })
}
