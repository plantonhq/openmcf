# Create the Azure Monitor autoscale setting -- the rule book that
# automatically adds and removes instances of ONE scalable target (VM
# Scale Set, App Service plan, ...) based on metric rules and schedules.
#
# Azure allows a single autoscale setting per target resource, and the
# setting must live in the target's region -- both enforced at deploy
# time. Exactly one profile is in effect at any moment: a matching
# fixed_date profile wins over a matching recurrence profile, which wins
# over the default (schedule-less) profile.
resource "azurerm_monitor_autoscale_setting" "main" {
  name                = var.spec.name
  resource_group_name = var.spec.resource_group
  location            = var.spec.region

  # The scalable resource this setting controls.
  target_resource_id = var.spec.target_resource_id

  # The provider requires an explicit value; the platform default is
  # true (a setting exists to scale).
  enabled = coalesce(var.spec.enabled, true)

  # Predictive autoscale (VM Scale Set targets only). Omitting the block
  # IS the "disabled" state -- the API exposes no Disabled mode.
  dynamic "predictive" {
    for_each = var.spec.predictive != null ? [var.spec.predictive] : []
    content {
      scale_mode = predictive.value.scale_mode
      # The provider validates a non-empty ISO 8601 duration -- send the
      # argument only when set.
      look_ahead_time = predictive.value.look_ahead_time != "" ? predictive.value.look_ahead_time : null
    }
  }

  dynamic "profile" {
    for_each = var.spec.profiles
    content {
      name = profile.value.name

      capacity {
        minimum = profile.value.capacity.minimum
        maximum = profile.value.capacity.maximum
        default = profile.value.capacity.default
      }

      dynamic "rule" {
        for_each = profile.value.rules
        content {
          metric_trigger {
            metric_name        = rule.value.metric_trigger.metric_name
            metric_resource_id = rule.value.metric_trigger.metric_resource_id
            time_grain         = rule.value.metric_trigger.time_grain
            statistic          = rule.value.metric_trigger.statistic
            time_window        = rule.value.metric_trigger.time_window
            time_aggregation   = rule.value.metric_trigger.time_aggregation
            operator           = rule.value.metric_trigger.operator
            threshold          = rule.value.metric_trigger.threshold
            # The provider validates a non-empty namespace -- send the
            # argument only when set (platform metrics imply their own).
            metric_namespace         = rule.value.metric_trigger.metric_namespace != "" ? rule.value.metric_trigger.metric_namespace : null
            divide_by_instance_count = rule.value.metric_trigger.divide_by_instance_count

            dynamic "dimensions" {
              for_each = rule.value.metric_trigger.dimensions
              content {
                name     = dimensions.value.name
                operator = dimensions.value.operator
                values   = dimensions.value.values
              }
            }
          }

          scale_action {
            direction = rule.value.scale_action.direction
            type      = rule.value.scale_action.type
            value     = rule.value.scale_action.value
            cooldown  = rule.value.scale_action.cooldown
          }
        }
      }

      dynamic "fixed_date" {
        for_each = profile.value.fixed_date != null ? [profile.value.fixed_date] : []
        content {
          # The platform default is UTC, always sent explicitly.
          timezone = coalesce(fixed_date.value.timezone, "UTC")
          start    = fixed_date.value.start
          end      = fixed_date.value.end
        }
      }

      dynamic "recurrence" {
        for_each = profile.value.recurrence != null ? [profile.value.recurrence] : []
        content {
          timezone = coalesce(recurrence.value.timezone, "UTC")
          days     = recurrence.value.days
          # The provider takes one-item lists here; the spec models the
          # single hour/minute the schedule fires at.
          hours   = [recurrence.value.hour]
          minutes = [recurrence.value.minute]
        }
      }
    }
  }

  dynamic "notification" {
    for_each = var.spec.notification != null ? [var.spec.notification] : []
    content {
      dynamic "email" {
        for_each = notification.value.email != null ? [notification.value.email] : []
        content {
          # The retired admin/co-admin flags stay unset (provider default
          # false) -- ARM rejects true since the April 2024 classic-
          # administrator retirement.
          custom_emails = email.value.custom_emails
        }
      }

      dynamic "webhook" {
        for_each = notification.value.webhooks
        content {
          service_uri = webhook.value.service_uri
          properties  = webhook.value.properties
        }
      }
    }
  }

  tags = local.final_tags
}
