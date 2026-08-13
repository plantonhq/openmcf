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
  description = "AwsBatchJobDefinition specification"
  type = object({
    region = string
    container = optional(object({
      image = string
      command = optional(list(string), [])
      vcpus = optional(number, 0)
      memory_mib = optional(number, 0)
      gpus = optional(number, 0)
      job_role = optional(string, "")
      execution_role = optional(string, "")
      environment = optional(map(string), {})
      secrets = optional(map(string), {})
      log_configuration = optional(object({
        log_driver = string
        options = optional(map(string), {})
        secret_options = optional(map(string), {})
      }))
      mount_points = optional(list(object({
        source_volume = string
        container_path = string
        read_only = optional(bool, false)
      })), [])
      volumes = optional(list(object({
        name = string
        efs = optional(object({
          file_system_id = string
          root_directory = optional(string, "")
          access_point_id = optional(string, "")
          iam_authorization = optional(bool, false)
        }))
        host_path = optional(string, "")
      })), [])
      ulimits = optional(list(object({
        name = string
        soft_limit = optional(number, 0)
        hard_limit = optional(number, 0)
      })), [])
      linux_parameters = optional(object({
        init_process_enabled = optional(bool, false)
        devices = optional(list(object({
          host_path = string
          container_path = optional(string, "")
          permissions = optional(list(string), [])
        })), [])
        shared_memory_size_mib = optional(number, 0)
        max_swap_mib = optional(number, 0)
        swappiness = optional(number, 0)
        tmpfs = optional(list(object({
          container_path = string
          size_mib = optional(number, 0)
          mount_options = optional(list(string), [])
        })), [])
      }))
      privileged = optional(bool, false)
      user = optional(string, "")
      readonly_root_filesystem = optional(bool, false)
      repository_credentials_secret_arn = optional(string, "")
      runtime_platform = optional(object({
        cpu_architecture = optional(string, "")
        operating_system_family = optional(string, "")
      }))
      fargate_platform_version = optional(string, "")
      assign_public_ip = optional(bool, false)
      ephemeral_storage_gib = optional(number, 0)
    }))
    platform_capabilities = optional(list(string), [])
    parameters = optional(map(string), {})
    retry_strategy = optional(object({
      attempts = optional(number, 0)
      evaluate_on_exit = optional(list(object({
        action = string
        on_exit_code = optional(string, "")
        on_reason = optional(string, "")
        on_status_reason = optional(string, "")
      })), [])
    }))
    timeout = optional(object({
      attempt_duration_seconds = optional(number, 0)
    }))
    scheduling_priority = optional(number, 0)
    propagate_tags = optional(bool, false)
    deregister_on_new_revision = optional(bool)
    eks = optional(object({
      containers = list(object({
        image = string
        name = optional(string, "")
        command = optional(list(string), [])
        args = optional(list(string), [])
        env = optional(map(string), {})
        image_pull_policy = optional(string, "")
        resources = optional(object({
          limits = optional(map(string), {})
          requests = optional(map(string), {})
        }))
        security_context = optional(object({
          run_as_user = optional(number)
          run_as_group = optional(number)
          run_as_non_root = optional(bool, false)
          allow_privilege_escalation = optional(bool)
          privileged = optional(bool, false)
          read_only_root_file_system = optional(bool, false)
        }))
        volume_mounts = optional(list(object({
          name = string
          mount_path = string
          read_only = optional(bool, false)
        })), [])
      }))
      init_containers = optional(list(object({
        image = string
        name = optional(string, "")
        command = optional(list(string), [])
        args = optional(list(string), [])
        env = optional(map(string), {})
        image_pull_policy = optional(string, "")
        resources = optional(object({
          limits = optional(map(string), {})
          requests = optional(map(string), {})
        }))
        security_context = optional(object({
          run_as_user = optional(number)
          run_as_group = optional(number)
          run_as_non_root = optional(bool, false)
          allow_privilege_escalation = optional(bool)
          privileged = optional(bool, false)
          read_only_root_file_system = optional(bool, false)
        }))
        volume_mounts = optional(list(object({
          name = string
          mount_path = string
          read_only = optional(bool, false)
        })), [])
      })), [])
      host_network = optional(bool)
      dns_policy = optional(string, "")
      service_account_name = optional(string, "")
      pod_labels = optional(map(string), {})
      image_pull_secret_names = optional(list(string), [])
      share_process_namespace = optional(bool, false)
      volumes = optional(list(object({
        name = string
        empty_dir = optional(object({
          medium = optional(string, "")
          size_limit = string
        }))
        host_path = optional(string, "")
        secret = optional(object({
          secret_name = string
          optional = optional(bool, false)
        }))
      })), [])
    }))
  })
}
