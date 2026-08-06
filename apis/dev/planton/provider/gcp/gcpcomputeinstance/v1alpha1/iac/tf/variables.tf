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
  description = "Specification for the Compute Engine instance"
  type = object({
    # StringValueOrRef fields arrive from the proto→tfvars converter as
    # plain strings (already resolved), never as object({value}).
    # Empty falls back to the provider's default project.
    project_id = optional(string, "")

    # Instance name; empty falls back to metadata.name. Immutable.
    instance_name = optional(string, "")

    # Zone, e.g. "us-central1-a". Immutable.
    zone = string

    # Machine type. Mutable via stop/update/start.
    machine_type = string

    description = optional(string, "")

    # Custom FQDN hostname. Immutable.
    hostname = optional(string, "")

    # Boot disk. Exactly one source (image / source_snapshot /
    # source_disk) — enforced pre-deploy by the spec's CEL.
    boot_disk = object({
      image           = optional(string, "")
      source_snapshot = optional(string, "")
      source_disk     = optional(string, "")
      size_gb         = optional(number)
      type            = optional(string, "")
      auto_delete     = optional(bool)
      device_name     = optional(string, "")
      kms_key         = optional(string, "")
      disk_labels     = optional(map(string), {})
      # Hyperdisk tuning; null for pd-* types.
      provisioned_iops       = optional(number)
      provisioned_throughput = optional(number)
      architecture           = optional(string, "")
      # Confidential-mode boot disk (hyperdisk SKUs; requires kms_key).
      enable_confidential_compute = optional(bool, false)
      resource_policies           = optional(list(string), [])
      storage_pool                = optional(string, "")
    })

    # Existing GcpComputeDisk attachments (source is the resolved disk
    # name/self link).
    attached_disks = optional(list(object({
      source      = string
      device_name = optional(string, "")
      mode        = optional(string, "")
      kms_key     = optional(string, "")
    })), [])

    # Ephemeral local SSDs. Create-time only.
    scratch_disks = optional(list(object({
      interface   = string
      size_gb     = optional(number)
      device_name = optional(string, "")
    })), [])

    # Network interfaces (at least one).
    network_interfaces = list(object({
      network            = optional(string, "")
      subnetwork         = optional(string, "")
      subnetwork_project = optional(string, "")
      # Static internal IP (resolved GcpAddress or literal IP).
      network_ip = optional(string, "")
      # At most one; presence grants an external IPv4.
      access_configs = optional(list(object({
        nat_ip                 = optional(string, "")
        network_tier           = optional(string, "")
        public_ptr_domain_name = optional(string, "")
      })), [])
      # At most one; presence grants external IPv6 (dual-stack subnets).
      ipv6_access_configs = optional(list(object({
        network_tier           = string
        public_ptr_domain_name = optional(string, "")
      })), [])
      stack_type  = optional(string, "")
      nic_type    = optional(string, "")
      queue_count = optional(number)
      alias_ip_ranges = optional(list(object({
        ip_cidr_range         = string
        subnetwork_range_name = optional(string, "")
      })), [])
    }))

    # VM identity. scopes required when the block is present.
    service_account = optional(object({
      email  = optional(string, "")
      scopes = list(string)
    }), null)

    # Provisioning model, maintenance, lifetime limits, sole-tenancy.
    scheduling = optional(object({
      provisioning_model          = optional(string, "")
      automatic_restart           = optional(bool)
      on_host_maintenance         = optional(string, "")
      instance_termination_action = optional(string, "")
      max_run_duration_seconds    = optional(number)
      termination_time            = optional(string, "")
      discard_local_ssds_on_stop  = optional(bool)
      availability_domain         = optional(number)
      min_node_cpus               = optional(number)
      node_affinities = optional(list(object({
        key      = string
        operator = string
        values   = list(string)
      })), [])
      local_ssd_recovery_timeout_seconds = optional(number)
    }), null)

    shielded_instance_config = optional(object({
      enable_secure_boot          = optional(bool)
      enable_vtpm                 = optional(bool)
      enable_integrity_monitoring = optional(bool)
    }), null)

    confidential_instance_config = optional(object({
      confidential_instance_type = string
    }), null)

    advanced_machine_features = optional(object({
      enable_nested_virtualization = optional(bool)
      threads_per_core             = optional(number)
      visible_core_count           = optional(number)
      enable_uefi_networking       = optional(bool)
      performance_monitoring_unit  = optional(string, "")
      turbo_mode                   = optional(string, "")
    }), null)

    guest_accelerators = optional(list(object({
      type  = string
      count = number
    })), [])

    reservation_affinity = optional(object({
      type = string
      specific_reservation = optional(object({
        key    = string
        values = list(string)
      }), null)
    }), null)

    total_egress_bandwidth_tier = optional(string, "")

    metadata       = optional(map(string), {})
    startup_script = optional(string, "")
    ssh_keys       = optional(list(string), [])

    labels                = optional(map(string), {})
    tags                  = optional(list(string), [])
    resource_manager_tags = optional(map(string), {})
    resource_policies     = optional(list(string), [])

    min_cpu_platform           = optional(string, "")
    can_ip_forward             = optional(bool, false)
    enable_display             = optional(bool, false)
    deletion_protection        = optional(bool, false)
    desired_status             = optional(string, "")
    allow_stopping_for_update  = optional(bool)
    key_revocation_action_type = optional(string, "")
  })
}
