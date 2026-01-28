##############################################
# variables.tf
#
# Input variables for KubernetesGhaRunnerScaleSetController
# deployment using Terraform.
#
# These variables map directly to the spec.proto fields.
##############################################

variable "metadata" {
  description = "Resource metadata including name, org, and environment"
  type = object({
    name = string
    id   = optional(string, "")
    org  = optional(string, "")
    env  = optional(string, "")
  })
}

variable "namespace" {
  description = "Kubernetes namespace where the controller will be installed"
  type        = string
}

variable "create_namespace" {
  description = "Whether to create the namespace"
  type        = bool
  default     = true
}

variable "helm_chart_version" {
  description = "Version of the Helm chart to deploy"
  type        = string
  default     = "0.13.1"
}

variable "replica_count" {
  description = "Number of controller replicas"
  type        = number
  default     = 1
}

variable "container" {
  description = "Container configuration"
  type = object({
    resources = optional(object({
      requests = optional(object({
        cpu    = optional(string, "100m")
        memory = optional(string, "128Mi")
      }), {})
      limits = optional(object({
        cpu    = optional(string, "500m")
        memory = optional(string, "512Mi")
      }), {})
    }), {})
    image = optional(object({
      repository  = optional(string, "")
      tag         = optional(string, "")
      pull_policy = optional(string, "IfNotPresent")
    }), {})
  })
  default = {}
}

variable "flags" {
  description = "Controller behavior flags"
  type = object({
    log_level                           = optional(string, "debug")
    log_format                          = optional(string, "text")
    watch_single_namespace              = optional(string, "")
    runner_max_concurrent_reconciles    = optional(number, 2)
    update_strategy                     = optional(string, "immediate")
    exclude_label_propagation_prefixes  = optional(list(string), [])
    k8s_client_rate_limiter_qps         = optional(number, 0)
    k8s_client_rate_limiter_burst       = optional(number, 0)
  })
  default = {}
}

variable "metrics" {
  description = "Metrics configuration for monitoring"
  type = object({
    controller_manager_addr = optional(string, "")
    listener_addr           = optional(string, "")
    listener_endpoint       = optional(string, "")
  })
  default = null
}

variable "image_pull_secrets" {
  description = "List of image pull secret names"
  type        = list(string)
  default     = []
}

variable "priority_class_name" {
  description = "Priority class name for the controller pods"
  type        = string
  default     = ""
}

