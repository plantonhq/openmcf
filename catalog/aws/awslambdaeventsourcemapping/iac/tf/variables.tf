variable "metadata" {
  description = "Cloud resource metadata"
  type = object({
    name        = string
    id          = optional(string, "")
    org         = optional(string, "")
    env         = optional(string, "")
    labels      = optional(map(string), {})
    annotations = optional(map(string), {})
    tags        = optional(list(string), [])
  })
}

variable "spec" {
  description = "AwsLambdaEventSourceMapping specification"
  type = object({
    region           = string
    function_arn     = string
    event_source_arn = optional(string, "")
    self_managed_kafka = optional(object({
      bootstrap_servers = list(string)
    }))
    disabled                        = optional(bool, false)
    batch_size                      = optional(number, 0)
    maximum_batching_window_seconds = optional(number, 0)
    filters = optional(list(object({
      pattern = optional(string, "")
    })), [])
    kms_key_arn                    = optional(string, "")
    function_response_types        = optional(list(string), [])
    scaling_max_concurrency        = optional(number, 0)
    metrics                        = optional(list(string), [])
    starting_position              = optional(string, "")
    starting_position_timestamp    = optional(string, "")
    parallelization_factor         = optional(number, 0)
    maximum_record_age_seconds     = optional(number, 0)
    maximum_retry_attempts         = optional(number)
    bisect_batch_on_function_error = optional(bool, false)
    tumbling_window_seconds        = optional(number, 0)
    on_failure_destination_arn     = optional(string, "")
    topics                         = optional(list(string), [])
    kafka_consumer_group_id        = optional(string, "")
    source_access_configurations = optional(list(object({
      type = string
      uri  = string
    })), [])
    schema_registry = optional(object({
      uri                   = string
      event_record_format   = string
      validation_attributes = optional(list(string), [])
      access_configurations = optional(list(object({
        type = string
        uri  = string
      })), [])
    }))
    provisioned_pollers = optional(object({
      minimum_pollers   = optional(number, 0)
      maximum_pollers   = optional(number, 0)
      poller_group_name = optional(string, "")
    }))
    mq_queue = optional(string, "")
    document_db = optional(object({
      database_name   = string
      collection_name = optional(string, "")
      full_document   = optional(string, "")
    }))
  })
}
