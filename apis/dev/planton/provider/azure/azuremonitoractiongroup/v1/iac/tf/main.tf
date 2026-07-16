# The action group is the notification hub alert rules fire into. It is
# a GLOBAL resource (location "global"): notifications keep flowing
# during regional outages, which is when they matter most. The location
# is a constant because "global" is the only value that makes sense for
# alerting infrastructure -- Azure accepts regional action groups only
# for niche data-residency cases.
resource "azurerm_monitor_action_group" "main" {
  name                = var.metadata.name
  resource_group_name = var.spec.resource_group
  short_name          = var.spec.short_name
  enabled             = var.spec.enabled

  dynamic "email_receiver" {
    for_each = var.spec.email_receivers
    content {
      name                    = email_receiver.value.name
      email_address           = email_receiver.value.email_address
      use_common_alert_schema = email_receiver.value.use_common_alert_schema
    }
  }

  # SMS and voice carry no schema toggle -- they are not payload-aware.
  dynamic "sms_receiver" {
    for_each = var.spec.sms_receivers
    content {
      name         = sms_receiver.value.name
      country_code = sms_receiver.value.country_code
      phone_number = sms_receiver.value.phone_number
    }
  }

  dynamic "voice_receiver" {
    for_each = var.spec.voice_receivers
    content {
      name         = voice_receiver.value.name
      country_code = voice_receiver.value.country_code
      phone_number = voice_receiver.value.phone_number
    }
  }

  dynamic "webhook_receiver" {
    for_each = var.spec.webhook_receivers
    content {
      name                    = webhook_receiver.value.name
      service_uri             = webhook_receiver.value.service_uri
      use_common_alert_schema = webhook_receiver.value.use_common_alert_schema

      # Entra-authenticated webhooks: the keyless posture -- no secret
      # baked into the URL.
      dynamic "aad_auth" {
        for_each = webhook_receiver.value.aad_auth != null ? [webhook_receiver.value.aad_auth] : []
        content {
          object_id      = aad_auth.value.object_id
          identifier_uri = aad_auth.value.identifier_uri
          tenant_id      = aad_auth.value.tenant_id
        }
      }
    }
  }

  dynamic "azure_app_push_receiver" {
    for_each = var.spec.azure_app_push_receivers
    content {
      name          = azure_app_push_receiver.value.name
      email_address = azure_app_push_receiver.value.email_address
    }
  }

  dynamic "automation_runbook_receiver" {
    for_each = var.spec.automation_runbook_receivers
    content {
      name                    = automation_runbook_receiver.value.name
      automation_account_id   = automation_runbook_receiver.value.automation_account_id
      runbook_name            = automation_runbook_receiver.value.runbook_name
      webhook_resource_id     = automation_runbook_receiver.value.webhook_resource_id
      is_global_runbook       = automation_runbook_receiver.value.is_global_runbook
      service_uri             = automation_runbook_receiver.value.service_uri
      use_common_alert_schema = automation_runbook_receiver.value.use_common_alert_schema
    }
  }

  dynamic "logic_app_receiver" {
    for_each = var.spec.logic_app_receivers
    content {
      name                    = logic_app_receiver.value.name
      resource_id             = logic_app_receiver.value.resource_id
      callback_url            = logic_app_receiver.value.callback_url
      use_common_alert_schema = logic_app_receiver.value.use_common_alert_schema
    }
  }

  dynamic "azure_function_receiver" {
    for_each = var.spec.azure_function_receivers
    content {
      name                     = azure_function_receiver.value.name
      function_app_resource_id = azure_function_receiver.value.function_app_resource_id
      function_name            = azure_function_receiver.value.function_name
      http_trigger_url         = azure_function_receiver.value.http_trigger_url
      use_common_alert_schema  = azure_function_receiver.value.use_common_alert_schema
    }
  }

  # Role-based fan-out: every user holding the role on the subscription
  # is notified -- no address list to maintain.
  dynamic "arm_role_receiver" {
    for_each = var.spec.arm_role_receivers
    content {
      name                    = arm_role_receiver.value.name
      role_id                 = arm_role_receiver.value.role_id
      use_common_alert_schema = arm_role_receiver.value.use_common_alert_schema
    }
  }

  dynamic "event_hub_receiver" {
    for_each = var.spec.event_hub_receivers
    content {
      name                    = event_hub_receiver.value.name
      event_hub_name          = event_hub_receiver.value.event_hub_name
      event_hub_namespace     = event_hub_receiver.value.event_hub_namespace
      tenant_id               = event_hub_receiver.value.tenant_id
      subscription_id         = event_hub_receiver.value.subscription_id
      use_common_alert_schema = event_hub_receiver.value.use_common_alert_schema
    }
  }

  # ITSM work-item creation; the ticket_configuration JSON must carry
  # PayloadRevision and WorkItemType (spec-enforced).
  dynamic "itsm_receiver" {
    for_each = var.spec.itsm_receivers
    content {
      name                 = itsm_receiver.value.name
      workspace_id         = itsm_receiver.value.workspace_id
      connection_id        = itsm_receiver.value.connection_id
      ticket_configuration = itsm_receiver.value.ticket_configuration
      region               = itsm_receiver.value.region
    }
  }

  tags = local.final_tags
}
