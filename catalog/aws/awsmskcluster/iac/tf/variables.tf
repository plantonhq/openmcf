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
  description = "AwsMskCluster specification"
  type = object({
    region = string
    kafka_version = string
    number_of_broker_nodes = number
    instance_type = string
    subnet_ids = list(string)
    security_group_ids = list(string)
    public_access_type = optional(string, "")
    vpc_connectivity = optional(object({
      sasl_iam_enabled = optional(bool, false)
      sasl_scram_enabled = optional(bool, false)
      tls_enabled = optional(bool, false)
    }))
    network_type = optional(string, "")
    ebs_volume_size_gib = optional(number)
    provisioned_throughput_enabled = optional(bool, false)
    provisioned_throughput_mbs = optional(number, 0)
    storage_mode = optional(string, "")
    kms_key_arn = optional(string, "")
    client_broker_encryption = optional(string)
    in_cluster_encryption = optional(bool)
    authentication = optional(object({
      sasl_iam_enabled = optional(bool, false)
      sasl_scram_enabled = optional(bool, false)
      tls_enabled = optional(bool, false)
      tls_certificate_authority_arns = optional(list(string), [])
      unauthenticated = optional(bool, false)
    }))
    scram_secret_arns = optional(list(string), [])
    cluster_policy = optional(any)
    configuration_arn = optional(string, "")
    configuration_revision = optional(number, 0)
    server_properties = optional(map(string), {})
    logging = optional(object({
      cloudwatch_logs = optional(object({
        enabled = optional(bool, false)
        log_group = optional(string, "")
      }))
      firehose = optional(object({
        enabled = optional(bool, false)
        delivery_stream = optional(string, "")
      }))
      s3 = optional(object({
        enabled = optional(bool, false)
        bucket = optional(string, "")
        prefix = optional(string, "")
      }))
    }))
    enhanced_monitoring = optional(string, "")
    jmx_exporter_enabled = optional(bool, false)
    node_exporter_enabled = optional(bool, false)
    rebalancing_status = optional(string, "")
    topics = optional(list(object({
      name = string
      partition_count = number
      replication_factor = number
      configs = optional(map(string), {})
    })), [])
  })
}
