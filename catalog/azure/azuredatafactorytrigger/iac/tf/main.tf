# One kind, four provider resources: Azure stores all four trigger
# types in the SAME factory-scoped trigger namespace
# ({factory_id}/triggers/{name}), so the spec's variant block selects
# which resource is created. All four share the started/stopped
# lifecycle: the provider stops a started trigger before any update
# or delete, then starts it again when `activated` is true (the
# platform default, sent explicitly).

# Fire on a wall-clock recurrence.
resource "azurerm_data_factory_trigger_schedule" "main" {
  count = var.spec.schedule != null ? 1 : 0

  name            = var.spec.name
  data_factory_id = var.spec.data_factory_id

  description = var.spec.description != "" ? var.spec.description : null
  annotations = length(var.spec.annotations) > 0 ? var.spec.annotations : null
  activated   = coalesce(var.spec.activated, true)

  frequency = coalesce(var.spec.schedule.frequency, "Minute")
  interval  = coalesce(var.spec.schedule.interval, 1)

  # Omitted start_time lets Azure fill in the moment of deployment
  # (the provider injects now-UTC itself when unset).
  start_time = var.spec.schedule.start_time != "" ? var.spec.schedule.start_time : null
  end_time   = var.spec.schedule.end_time != "" ? var.spec.schedule.end_time : null
  time_zone  = var.spec.schedule.time_zone != "" ? var.spec.schedule.time_zone : null

  # The provider's `schedule` block -- the recurrence narrowing rules.
  dynamic "schedule" {
    for_each = var.spec.schedule.recurrence_schedule != null ? [var.spec.schedule.recurrence_schedule] : []
    content {
      days_of_month = length(schedule.value.days_of_month) > 0 ? schedule.value.days_of_month : null
      days_of_week  = length(schedule.value.days_of_week) > 0 ? schedule.value.days_of_week : null
      hours         = length(schedule.value.hours) > 0 ? schedule.value.hours : null
      minutes       = length(schedule.value.minutes) > 0 ? schedule.value.minutes : null

      dynamic "monthly" {
        for_each = schedule.value.monthly
        content {
          weekday = monthly.value.weekday
          week    = monthly.value.week
        }
      }
    }
  }

  # The modern plural pipeline blocks (the legacy singular
  # pipeline_name/pipeline_parameters pair covers the same wire
  # surface and is never sent).
  dynamic "pipeline" {
    for_each = var.spec.schedule.pipelines
    content {
      name       = pipeline.value.name
      parameters = length(pipeline.value.parameters) > 0 ? pipeline.value.parameters : null
    }
  }
}

# Fire once per contiguous, non-overlapping time window.
resource "azurerm_data_factory_trigger_tumbling_window" "main" {
  count = var.spec.tumbling_window != null ? 1 : 0

  name            = var.spec.name
  data_factory_id = var.spec.data_factory_id

  description = var.spec.description != "" ? var.spec.description : null
  annotations = length(var.spec.annotations) > 0 ? var.spec.annotations : null
  activated   = coalesce(var.spec.activated, true)

  frequency  = var.spec.tumbling_window.frequency
  interval   = var.spec.tumbling_window.interval
  start_time = var.spec.tumbling_window.start_time
  end_time   = var.spec.tumbling_window.end_time != "" ? var.spec.tumbling_window.end_time : null
  delay      = var.spec.tumbling_window.delay != "" ? var.spec.tumbling_window.delay : null

  max_concurrency = coalesce(var.spec.tumbling_window.max_concurrency, 50)

  dynamic "retry" {
    for_each = var.spec.tumbling_window.retry != null ? [var.spec.tumbling_window.retry] : []
    content {
      count    = retry.value.count
      interval = coalesce(retry.value.interval, 30)
    }
  }

  # An entry without trigger_name is a SELF-dependency (this trigger's
  # own earlier windows) -- the provider's own convention.
  dynamic "trigger_dependency" {
    for_each = var.spec.tumbling_window.dependencies
    content {
      trigger_name = trigger_dependency.value.trigger_name != "" ? trigger_dependency.value.trigger_name : null
      offset       = trigger_dependency.value.offset != "" ? trigger_dependency.value.offset : null
      size         = trigger_dependency.value.size != "" ? trigger_dependency.value.size : null
    }
  }

  additional_properties = length(var.spec.tumbling_window.additional_properties) > 0 ? var.spec.tumbling_window.additional_properties : null

  # Tumbling window triggers drive exactly ONE pipeline (Azure's own
  # model).
  pipeline {
    name       = var.spec.tumbling_window.pipeline.name
    parameters = length(var.spec.tumbling_window.pipeline.parameters) > 0 ? var.spec.tumbling_window.pipeline.parameters : null
  }
}

# Fire on blob creation/deletion in a storage account (Azure wires
# this through Event Grid on the storage account behind the scenes).
resource "azurerm_data_factory_trigger_blob_event" "main" {
  count = var.spec.blob_event != null ? 1 : 0

  name            = var.spec.name
  data_factory_id = var.spec.data_factory_id

  description = var.spec.description != "" ? var.spec.description : null
  annotations = length(var.spec.annotations) > 0 ? var.spec.annotations : null
  activated   = coalesce(var.spec.activated, true)

  storage_account_id = var.spec.blob_event.storage_account_id
  events             = var.spec.blob_event.events

  blob_path_begins_with = var.spec.blob_event.blob_path_begins_with != "" ? var.spec.blob_event.blob_path_begins_with : null
  blob_path_ends_with   = var.spec.blob_event.blob_path_ends_with != "" ? var.spec.blob_event.blob_path_ends_with : null

  ignore_empty_blobs = coalesce(var.spec.blob_event.ignore_empty_blobs, false)

  additional_properties = length(var.spec.blob_event.additional_properties) > 0 ? var.spec.blob_event.additional_properties : null

  dynamic "pipeline" {
    for_each = var.spec.blob_event.pipelines
    content {
      name       = pipeline.value.name
      parameters = length(pipeline.value.parameters) > 0 ? pipeline.value.parameters : null
    }
  }
}

# Fire on events published to an Event Grid custom topic.
resource "azurerm_data_factory_trigger_custom_event" "main" {
  count = var.spec.custom_event != null ? 1 : 0

  name            = var.spec.name
  data_factory_id = var.spec.data_factory_id

  description = var.spec.description != "" ? var.spec.description : null
  annotations = length(var.spec.annotations) > 0 ? var.spec.annotations : null
  activated   = coalesce(var.spec.activated, true)

  eventgrid_topic_id = var.spec.custom_event.eventgrid_topic_id
  events             = var.spec.custom_event.events

  subject_begins_with = var.spec.custom_event.subject_begins_with != "" ? var.spec.custom_event.subject_begins_with : null
  subject_ends_with   = var.spec.custom_event.subject_ends_with != "" ? var.spec.custom_event.subject_ends_with : null

  additional_properties = length(var.spec.custom_event.additional_properties) > 0 ? var.spec.custom_event.additional_properties : null

  dynamic "pipeline" {
    for_each = var.spec.custom_event.pipelines
    content {
      name       = pipeline.value.name
      parameters = length(pipeline.value.parameters) > 0 ? pipeline.value.parameters : null
    }
  }
}
