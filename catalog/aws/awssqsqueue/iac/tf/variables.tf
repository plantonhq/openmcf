variable "metadata" {
  description = "Cloud resource metadata"
  type = object({
    name = string
    id = optional(string, "")
    org = optional(string, "")
    env = optional(string, "")
    labels = optional(map(string), {})
    annotations = optional(map(string), {})
    tags = optional(list(string), [])
  })
}

variable "spec" {
  description = "AwsSqsQueue specification"
  type = object({
    region = string
    fifo_queue = optional(bool, false)
    visibility_timeout_seconds = optional(number, 0)
    message_retention_seconds = optional(number, 0)
    max_message_size_bytes = optional(number, 0)
    delay_seconds = optional(number, 0)
    receive_wait_time_seconds = optional(number, 0)
    content_based_deduplication = optional(bool, false)
    deduplication_scope = optional(string, "")
    fifo_throughput_limit = optional(string, "")
    dead_letter_config = optional(object({
      target_arn = string
      max_receive_count = optional(number, 0)
    }))
    kms_key_id = optional(string, "")
    kms_data_key_reuse_period_seconds = optional(number, 0)
    sqs_managed_sse_enabled = optional(bool)
    policy = optional(any)
    redrive_allow_policy = optional(object({
      redrive_permission = string
      source_queue_arns = optional(list(string), [])
    }))
  })
}
