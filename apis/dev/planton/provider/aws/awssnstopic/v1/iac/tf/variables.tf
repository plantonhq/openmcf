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
  description = "AwsSnsTopic specification"
  type = object({
    region = string
    fifo_topic = optional(bool, false)
    content_based_deduplication = optional(bool, false)
    fifo_throughput_scope = optional(string, "")
    archive_policy = optional(any)
    display_name = optional(string, "")
    kms_key_id = optional(string, "")
    policy = optional(any)
    data_protection_policy = optional(any)
    delivery_policy = optional(string, "")
    delivery_feedback = optional(object({
      application = optional(object({
        success_feedback_role = optional(string, "")
        failure_feedback_role = optional(string, "")
        success_feedback_sample_rate = optional(number, 0)
      }))
      firehose = optional(object({
        success_feedback_role = optional(string, "")
        failure_feedback_role = optional(string, "")
        success_feedback_sample_rate = optional(number, 0)
      }))
      http = optional(object({
        success_feedback_role = optional(string, "")
        failure_feedback_role = optional(string, "")
        success_feedback_sample_rate = optional(number, 0)
      }))
      lambda = optional(object({
        success_feedback_role = optional(string, "")
        failure_feedback_role = optional(string, "")
        success_feedback_sample_rate = optional(number, 0)
      }))
      sqs = optional(object({
        success_feedback_role = optional(string, "")
        failure_feedback_role = optional(string, "")
        success_feedback_sample_rate = optional(number, 0)
      }))
    }))
    tracing_config = optional(string, "")
    signature_version = optional(number, 0)
  })
}
