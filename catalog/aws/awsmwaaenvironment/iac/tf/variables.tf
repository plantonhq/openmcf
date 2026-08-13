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
  description = "AwsMwaaEnvironment specification"
  type = object({
    region                           = string
    airflow_version                  = optional(string, "")
    airflow_configuration_options    = optional(map(string), {})
    source_bucket_arn                = string
    dag_s3_path                      = string
    plugins_s3_path                  = optional(string, "")
    plugins_s3_object_version        = optional(string, "")
    requirements_s3_path             = optional(string, "")
    requirements_s3_object_version   = optional(string, "")
    startup_script_s3_path           = optional(string, "")
    startup_script_s3_object_version = optional(string, "")
    execution_role_arn               = string
    subnet_ids                       = list(string)
    security_group_ids               = list(string)
    kms_key_arn                      = optional(string, "")
    environment_class                = optional(string, "")
    min_workers                      = optional(number, 0)
    max_workers                      = optional(number, 0)
    min_webservers                   = optional(number, 0)
    max_webservers                   = optional(number, 0)
    schedulers                       = optional(number, 0)
    webserver_access_mode            = optional(string)
    endpoint_management              = optional(string, "")
    logging_configuration = optional(object({
      dag_processing_logs = optional(object({
        enabled   = optional(bool, false)
        log_level = optional(string, "")
      }))
      scheduler_logs = optional(object({
        enabled   = optional(bool, false)
        log_level = optional(string, "")
      }))
      task_logs = optional(object({
        enabled   = optional(bool, false)
        log_level = optional(string, "")
      }))
      webserver_logs = optional(object({
        enabled   = optional(bool, false)
        log_level = optional(string, "")
      }))
      worker_logs = optional(object({
        enabled   = optional(bool, false)
        log_level = optional(string, "")
      }))
    }))
    weekly_maintenance_window_start = optional(string, "")
    worker_replacement_strategy     = optional(string, "")
  })
}
