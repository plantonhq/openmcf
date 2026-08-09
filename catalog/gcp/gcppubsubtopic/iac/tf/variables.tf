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
  description = "Specification for the GCP Pub/Sub topic"
  type = object({
    # StringValueOrRef fields arrive from the proto→tfvars converter as
    # plain strings (already resolved), never as object({value}). If
    # project_id is empty, the provider's default project is used
    # (see locals.tf).
    project_id = optional(string, "")

    # Topic name (the GCP resource name). Immutable (ForceNew).
    topic_name = string

    # CMEK key path (resolved from a GcpKmsKey reference). Empty means
    # Google-managed encryption.
    kms_key_name = optional(string, "")

    # Topic-level retention window (e.g. "604800s"). Empty leaves
    # retention to individual subscriptions.
    message_retention_duration = optional(string, "")

    # User labels, merged beneath the platform attribution labels
    # (see locals.tf).
    labels = optional(map(string), {})

    message_storage_policy = optional(object({
      allowed_persistence_regions = list(string)
      enforce_in_transit          = optional(bool, false)
    }), null)

    schema_settings = optional(object({
      # Resolved from a GcpPubSubSchema reference to the fully qualified
      # projects/{project}/schemas/{name} path.
      schema   = string
      encoding = optional(string, "")

      # Revision-range pinning (resolved from GcpPubSubSchema revision_id
      # references or given as literal revision IDs). Empty means any
      # revision on that bound.
      first_revision_id = optional(string, "")
      last_revision_id  = optional(string, "")
    }), null)

    ingestion_data_source_settings = optional(object({
      aws_kinesis = optional(object({
        stream_arn          = string
        consumer_arn        = string
        aws_role_arn        = string
        gcp_service_account = string
      }), null)

      aws_msk = optional(object({
        cluster_arn         = string
        topic               = string
        aws_role_arn        = string
        gcp_service_account = string
      }), null)

      azure_event_hubs = optional(object({
        resource_group      = optional(string, "")
        namespace           = optional(string, "")
        event_hub           = optional(string, "")
        client_id           = optional(string, "")
        tenant_id           = optional(string, "")
        subscription_id     = optional(string, "")
        gcp_service_account = optional(string, "")
      }), null)

      cloud_storage = optional(object({
        bucket                     = string
        match_glob                 = optional(string, "")
        minimum_object_create_time = optional(string, "")
        avro_format                = optional(bool, false)
        pubsub_avro_format         = optional(bool, false)
        text_format = optional(object({
          delimiter = optional(string, "")
        }), null)
      }), null)

      confluent_cloud = optional(object({
        bootstrap_server    = string
        topic               = string
        identity_pool_id    = string
        gcp_service_account = string
        cluster_id          = optional(string, "")
      }), null)

      platform_logs_settings = optional(object({
        severity = optional(string, "")
      }), null)
    }), null)

    # Ordered transform pipeline applied to every published message.
    # Each step carries exactly one transform arm (spec-enforced):
    # javascript_udf or ai_inference.
    message_transforms = optional(list(object({
      javascript_udf = optional(object({
        function_name = string
        code          = string
      }), null)
      ai_inference = optional(object({
        # Vertex AI endpoint path (resolved from a GcpVertexAiEndpoint
        # reference or given as a literal dedicated-endpoint or
        # publisher-model path).
        endpoint              = string
        service_account_email = optional(string, "")
        unstructured_inference = optional(object({
          parameters = optional(map(string), {})
        }), null)
      }), null)
      disabled = optional(bool, false)
    })), [])

    # Resource Manager tags bound at create time (tagKeys/{id} =>
    # tagValues/{id}). Changing them later replaces the topic.
    resource_manager_tags = optional(map(string), {})

    # Deletion policy: "", "DELETE" (default), "PREVENT" (destroy fails),
    # or "ABANDON" (remove from management, leave serving in GCP).
    deletion_policy = optional(string, "")
  })

  validation {
    condition     = var.spec.topic_name != ""
    error_message = "topic_name is required."
  }

  validation {
    condition     = contains(["", "DELETE", "PREVENT", "ABANDON"], var.spec.deletion_policy)
    error_message = "deletion_policy must be one of: DELETE, PREVENT, ABANDON."
  }
}
