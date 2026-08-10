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
  description = "GcpEventarcTrigger specification"
  type = object({
    project_id   = optional(string, "")
    location     = string
    trigger_name = optional(string, "")
    matching_criteria = list(object({
      attribute = string
      value     = string
      operator  = optional(string, "")
    }))
    destination = object({
      cloud_run_service = optional(object({
        service = string
        region  = optional(string, "")
        path    = optional(string, "")
      }))
      gke = optional(object({
        cluster   = string
        location  = string
        namespace = string
        service   = string
        path      = optional(string, "")
      }))
      workflow = optional(string, "")
      http_endpoint = optional(object({
        uri                = string
        network_attachment = string
      }))
    })
    service_account         = optional(string, "")
    transport_pubsub_topic  = optional(string, "")
    event_data_content_type = optional(string, "")
    retry_max_attempts      = optional(number, 0)
    labels                  = optional(map(string), {})
    partner_channel = optional(object({
      channel_name         = optional(string, "")
      third_party_provider = string
      crypto_key           = optional(string, "")
    }))
    google_channel_crypto_key = optional(string, "")
    deletion_policy           = optional(string, "")
  })
}