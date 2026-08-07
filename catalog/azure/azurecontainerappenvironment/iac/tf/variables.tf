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
  description = "Azure Container App Environment specification"
  type = object({
    # The Azure region where the environment will be created. ForceNew.
    region = string

    # The Azure Resource Group name. References are resolved to a literal
    # name by the platform before the module runs. ForceNew.
    resource_group = string

    # The name of the Container App Environment (2-60 alphanumerics/
    # hyphens; no leading/trailing hyphen). ForceNew.
    environment_name = string

    # Where application logs persist, as the spec enum's name string
    # (LOG_ANALYTICS / AZURE_MONITOR). Absent with a workspace deploys
    # log-analytics; absent without one means streaming-only.
    logs_destination = optional(string)

    # Log Analytics Workspace ARM ID. Required with LOG_ANALYTICS,
    # forbidden with AZURE_MONITOR (spec-enforced).
    log_analytics_workspace_id = optional(string)

    # Application Insights connection string Dapr uses for service-to-
    # service telemetry. Write-only in ARM; ForceNew.
    dapr_application_insights_connection_string = optional(string)

    # Existing subnet ARM ID for VNet injection (/21 or larger). ForceNew.
    infrastructure_subnet_id = optional(string)

    # Name for the platform-managed infrastructure resource group. Only
    # valid with workload profiles (spec-enforced). ForceNew.
    infrastructure_resource_group_name = optional(string)

    # Internal load balancing mode -- apps reachable only inside the VNet.
    # Requires infrastructure_subnet_id (spec-enforced). ForceNew.
    internal_load_balancer_enabled = optional(bool, false)

    # Spread infrastructure across availability zones. Requires
    # infrastructure_subnet_id (spec-enforced). ForceNew.
    zone_redundancy_enabled = optional(bool, false)

    # Platform-level public network access, as the spec enum's name
    # string (ENABLED / DISABLED). Absent lets Azure derive it from the
    # network configuration.
    public_network_access = optional(string)

    # Mutual TLS between apps in the environment.
    mutual_tls_enabled = optional(bool, false)

    # Dedicated / GPU workload profiles. workload_profile_type is the
    # spec enum's name string (e.g. D4, E8, NC24_A100,
    # CONSUMPTION_GPU_NC8AS_T4). Counts only apply to dedicated families
    # (spec-enforced).
    workload_profiles = optional(list(object({
      name                  = string
      workload_profile_type = string
      minimum_count         = optional(number)
      maximum_count         = optional(number)
    })), [])

    # Environment-wide custom DNS suffix backed by a wildcard PFX
    # certificate. ARM models it as a patch on the environment itself.
    custom_domain = optional(object({
      dns_suffix              = string
      certificate_blob_base64 = string
      certificate_password    = string
    }))

    # Managed identity: type is the spec enum's name string
    # (SYSTEM_ASSIGNED / USER_ASSIGNED / SYSTEM_AND_USER_ASSIGNED);
    # user_assigned_identity_ids carry AzureUserAssignedIdentity ARM ids.
    identity = optional(object({
      type                       = string
      user_assigned_identity_ids = optional(list(string), [])
    }))

    # Free-form user tags, merged over the metadata-derived tags (user
    # tags win on key collision).
    tags = optional(map(string), {})
  })
}
