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
  description = "AwsBatchComputeEnvironment specification"
  type = object({
    region = string
    state = optional(string)
    service_role = optional(string, "")
    compute_resources = object({
      type = string
      max_vcpus = number
      min_vcpus = optional(number)
      desired_vcpus = optional(number, 0)
      subnet_ids = list(string)
      security_group_ids = optional(list(string), [])
      instance_types = optional(list(string), [])
      allocation_strategy = optional(string, "")
      instance_role = optional(string, "")
      ec2_key_pair = optional(string, "")
      bid_percentage = optional(number)
      spot_iam_fleet_role = optional(string, "")
      launch_template = optional(object({
        launch_template_id = string
        version = optional(string, "")
      }))
      ec2_configurations = optional(list(object({
        image_type = optional(string, "")
        image_id_override = optional(string, "")
        image_kubernetes_version = optional(string, "")
      })), [])
      placement_group = optional(string, "")
      resource_tags = optional(map(string), {})
    })
    eks_configuration = optional(object({
      eks_cluster_arn = string
      kubernetes_namespace = string
    }))
    update_policy = optional(object({
      terminate_jobs_on_update = optional(bool, false)
      job_execution_timeout_minutes = optional(number)
    }))
  })
}
