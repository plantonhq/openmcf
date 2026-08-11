# Enable the Cloud Monitoring API so a fresh project can host the policy.
# disable_on_destroy is false: tearing down one policy must never disable
# monitoring for everything else in the project.
resource "google_project_service" "monitoring_api" {
  project = local.project_id
  service = "monitoring.googleapis.com"

  disable_dependent_services = true
  disable_on_destroy         = false
}

# A Cloud Monitoring alert policy — the rule that watches metrics or logs
# and notifies the referenced channels when incidents open.
#
# Each spec condition carries exactly one condition-type arm (the spec's
# validations enforce the API's oneof, which the provider leaves unchecked
# client-side), so each dynamic block below renders at most once per
# condition.
#
# `enabled` is sent EXPLICITLY on every apply: it is Optional in the
# provider with a server default of true, and a spec transition
# true -> false must reach the API rather than being omitted (the
# send-true-or-omit class — a silently still-enabled alert policy pages
# people at 3am).
resource "google_monitoring_alert_policy" "this" {
  display_name = local.display_name
  combiner     = var.spec.combiner
  project      = local.project_id

  severity = var.spec.severity != "" ? var.spec.severity : null
  enabled  = var.spec.enabled == null ? true : var.spec.enabled

  notification_channels = length(var.spec.notification_channels) > 0 ? var.spec.notification_channels : null

  user_labels = local.final_labels

  dynamic "conditions" {
    for_each = var.spec.conditions
    content {
      display_name = conditions.value.display_name

      dynamic "condition_threshold" {
        for_each = conditions.value.condition_threshold != null ? [conditions.value.condition_threshold] : []
        content {
          filter     = condition_threshold.value.filter
          comparison = condition_threshold.value.comparison
          # threshold_value 0 is a legal threshold (e.g. COMPARISON_GT 0 for
          # "any errors at all"), so it is always sent — zero and unset are
          # deliberately NOT distinguished for this field.
          threshold_value = condition_threshold.value.threshold_value
          duration        = condition_threshold.value.duration

          dynamic "aggregations" {
            for_each = condition_threshold.value.aggregations
            content {
              alignment_period     = aggregations.value.alignment_period != "" ? aggregations.value.alignment_period : null
              per_series_aligner   = aggregations.value.per_series_aligner != "" ? aggregations.value.per_series_aligner : null
              cross_series_reducer = aggregations.value.cross_series_reducer != "" ? aggregations.value.cross_series_reducer : null
              group_by_fields      = length(aggregations.value.group_by_fields) > 0 ? aggregations.value.group_by_fields : null
            }
          }

          denominator_filter = condition_threshold.value.denominator_filter != "" ? condition_threshold.value.denominator_filter : null

          dynamic "denominator_aggregations" {
            for_each = condition_threshold.value.denominator_aggregations
            content {
              alignment_period     = denominator_aggregations.value.alignment_period != "" ? denominator_aggregations.value.alignment_period : null
              per_series_aligner   = denominator_aggregations.value.per_series_aligner != "" ? denominator_aggregations.value.per_series_aligner : null
              cross_series_reducer = denominator_aggregations.value.cross_series_reducer != "" ? denominator_aggregations.value.cross_series_reducer : null
              group_by_fields      = length(denominator_aggregations.value.group_by_fields) > 0 ? denominator_aggregations.value.group_by_fields : null
            }
          }

          dynamic "forecast_options" {
            for_each = condition_threshold.value.forecast_options != null ? [condition_threshold.value.forecast_options] : []
            content {
              forecast_horizon = forecast_options.value.forecast_horizon
            }
          }

          dynamic "trigger" {
            for_each = condition_threshold.value.trigger != null ? [condition_threshold.value.trigger] : []
            content {
              count   = trigger.value.count != 0 ? trigger.value.count : null
              percent = trigger.value.percent != 0 ? trigger.value.percent : null
            }
          }

          evaluation_missing_data = condition_threshold.value.evaluation_missing_data != "" ? condition_threshold.value.evaluation_missing_data : null
        }
      }

      dynamic "condition_absent" {
        for_each = conditions.value.condition_absent != null ? [conditions.value.condition_absent] : []
        content {
          filter   = condition_absent.value.filter
          duration = condition_absent.value.duration

          dynamic "aggregations" {
            for_each = condition_absent.value.aggregations
            content {
              alignment_period     = aggregations.value.alignment_period != "" ? aggregations.value.alignment_period : null
              per_series_aligner   = aggregations.value.per_series_aligner != "" ? aggregations.value.per_series_aligner : null
              cross_series_reducer = aggregations.value.cross_series_reducer != "" ? aggregations.value.cross_series_reducer : null
              group_by_fields      = length(aggregations.value.group_by_fields) > 0 ? aggregations.value.group_by_fields : null
            }
          }

          dynamic "trigger" {
            for_each = condition_absent.value.trigger != null ? [condition_absent.value.trigger] : []
            content {
              count   = trigger.value.count != 0 ? trigger.value.count : null
              percent = trigger.value.percent != 0 ? trigger.value.percent : null
            }
          }
        }
      }

      dynamic "condition_matched_log" {
        for_each = conditions.value.condition_matched_log != null ? [conditions.value.condition_matched_log] : []
        content {
          filter           = condition_matched_log.value.filter
          label_extractors = length(condition_matched_log.value.label_extractors) > 0 ? condition_matched_log.value.label_extractors : null
        }
      }

      dynamic "condition_monitoring_query_language" {
        for_each = conditions.value.condition_monitoring_query_language != null ? [conditions.value.condition_monitoring_query_language] : []
        content {
          query    = condition_monitoring_query_language.value.query
          duration = condition_monitoring_query_language.value.duration

          dynamic "trigger" {
            for_each = condition_monitoring_query_language.value.trigger != null ? [condition_monitoring_query_language.value.trigger] : []
            content {
              count   = trigger.value.count != 0 ? trigger.value.count : null
              percent = trigger.value.percent != 0 ? trigger.value.percent : null
            }
          }

          evaluation_missing_data = condition_monitoring_query_language.value.evaluation_missing_data != "" ? condition_monitoring_query_language.value.evaluation_missing_data : null
        }
      }

      dynamic "condition_prometheus_query_language" {
        for_each = conditions.value.condition_prometheus_query_language != null ? [conditions.value.condition_prometheus_query_language] : []
        content {
          query                     = condition_prometheus_query_language.value.query
          duration                  = condition_prometheus_query_language.value.duration != "" ? condition_prometheus_query_language.value.duration : null
          evaluation_interval       = condition_prometheus_query_language.value.evaluation_interval != "" ? condition_prometheus_query_language.value.evaluation_interval : null
          labels                    = length(condition_prometheus_query_language.value.labels) > 0 ? condition_prometheus_query_language.value.labels : null
          rule_group                = condition_prometheus_query_language.value.rule_group != "" ? condition_prometheus_query_language.value.rule_group : null
          alert_rule                = condition_prometheus_query_language.value.alert_rule != "" ? condition_prometheus_query_language.value.alert_rule : null
          disable_metric_validation = condition_prometheus_query_language.value.disable_metric_validation ? true : null
        }
      }

      dynamic "condition_sql" {
        for_each = conditions.value.condition_sql != null ? [conditions.value.condition_sql] : []
        content {
          query = condition_sql.value.query

          dynamic "minutes" {
            for_each = condition_sql.value.minutes != null ? [condition_sql.value.minutes] : []
            content {
              periodicity = minutes.value.periodicity
            }
          }

          dynamic "hourly" {
            for_each = condition_sql.value.hourly != null ? [condition_sql.value.hourly] : []
            content {
              periodicity   = hourly.value.periodicity
              minute_offset = hourly.value.minute_offset
            }
          }

          dynamic "daily" {
            for_each = condition_sql.value.daily != null ? [condition_sql.value.daily] : []
            content {
              periodicity = daily.value.periodicity

              dynamic "execution_time" {
                for_each = daily.value.execution_time != null ? [daily.value.execution_time] : []
                content {
                  hours   = execution_time.value.hours
                  minutes = execution_time.value.minutes
                  seconds = execution_time.value.seconds
                  nanos   = execution_time.value.nanos
                }
              }
            }
          }

          dynamic "row_count_test" {
            for_each = condition_sql.value.row_count_test != null ? [condition_sql.value.row_count_test] : []
            content {
              comparison = row_count_test.value.comparison
              threshold  = row_count_test.value.threshold
            }
          }

          dynamic "boolean_test" {
            for_each = condition_sql.value.boolean_test != null ? [condition_sql.value.boolean_test] : []
            content {
              column = boolean_test.value.column
            }
          }
        }
      }
    }
  }

  dynamic "alert_strategy" {
    for_each = var.spec.alert_strategy != null ? [var.spec.alert_strategy] : []
    content {
      auto_close = alert_strategy.value.auto_close != "" ? alert_strategy.value.auto_close : null

      # Required for (and only legal on) log-based policies — the API's own
      # pairing, taught on the spec field.
      dynamic "notification_rate_limit" {
        for_each = alert_strategy.value.notification_rate_limit != null ? [alert_strategy.value.notification_rate_limit] : []
        content {
          period = notification_rate_limit.value.period != "" ? notification_rate_limit.value.period : null
        }
      }

      dynamic "notification_channel_strategy" {
        for_each = alert_strategy.value.notification_channel_strategy
        content {
          notification_channel_names = length(notification_channel_strategy.value.notification_channel_names) > 0 ? notification_channel_strategy.value.notification_channel_names : null
          renotify_interval          = notification_channel_strategy.value.renotify_interval != "" ? notification_channel_strategy.value.renotify_interval : null
        }
      }

      notification_prompts = length(alert_strategy.value.notification_prompts) > 0 ? alert_strategy.value.notification_prompts : null
    }
  }

  dynamic "documentation" {
    for_each = var.spec.documentation != null ? [var.spec.documentation] : []
    content {
      content   = documentation.value.content != "" ? documentation.value.content : null
      mime_type = documentation.value.mime_type != "" ? documentation.value.mime_type : null
      subject   = documentation.value.subject != "" ? documentation.value.subject : null

      dynamic "links" {
        for_each = documentation.value.links
        content {
          display_name = links.value.display_name != "" ? links.value.display_name : null
          url          = links.value.url != "" ? links.value.url : null
        }
      }
    }
  }

  # Empty defers to the provider default (DELETE).
  deletion_policy = var.spec.deletion_policy != "" ? var.spec.deletion_policy : null

  depends_on = [google_project_service.monitoring_api]
}
