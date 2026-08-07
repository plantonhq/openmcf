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
  description = "Azure Monitor action group specification"
  type = object({
    # The Azure Resource Group name (references resolved by the platform
    # before the module runs). Action groups are global -- no region.
    resource_group = string

    # The 1-12 character short name shown in SMS and push notifications
    short_name = string

    # Whether the group is active (a disabled group swallows alerts)
    enabled = optional(bool, true)

    # Receiver lists -- all optional; a receiver-less group is a legal
    # "null" routing target.
    email_receivers = optional(list(object({
      name                    = string
      email_address           = string
      use_common_alert_schema = optional(bool, false)
    })), [])

    sms_receivers = optional(list(object({
      name         = string
      country_code = string
      phone_number = string
    })), [])

    voice_receivers = optional(list(object({
      name         = string
      country_code = string
      phone_number = string
    })), [])

    webhook_receivers = optional(list(object({
      name                    = string
      service_uri             = string
      use_common_alert_schema = optional(bool, false)
      aad_auth = optional(object({
        object_id      = string
        identifier_uri = optional(string)
        tenant_id      = optional(string)
      }))
    })), [])

    azure_app_push_receivers = optional(list(object({
      name          = string
      email_address = string
    })), [])

    automation_runbook_receivers = optional(list(object({
      name                    = string
      automation_account_id   = string
      runbook_name            = string
      webhook_resource_id     = string
      is_global_runbook       = optional(bool, false)
      service_uri             = string
      use_common_alert_schema = optional(bool, false)
    })), [])

    logic_app_receivers = optional(list(object({
      name                    = string
      resource_id             = string
      callback_url            = string
      use_common_alert_schema = optional(bool, false)
    })), [])

    azure_function_receivers = optional(list(object({
      name                    = string
      function_app_resource_id = string
      function_name           = string
      http_trigger_url        = string
      use_common_alert_schema = optional(bool, false)
    })), [])

    arm_role_receivers = optional(list(object({
      name                    = string
      role_id                 = string
      use_common_alert_schema = optional(bool, false)
    })), [])

    event_hub_receivers = optional(list(object({
      name                    = string
      event_hub_name          = string
      event_hub_namespace     = string
      tenant_id               = optional(string)
      subscription_id         = optional(string)
      use_common_alert_schema = optional(bool, false)
    })), [])

    itsm_receivers = optional(list(object({
      name                 = string
      workspace_id         = string
      connection_id        = string
      ticket_configuration = string
      region               = string
    })), [])

    # User tags, merged over the metadata-derived tags (user wins)
    tags = optional(map(string), {})
  })
}
