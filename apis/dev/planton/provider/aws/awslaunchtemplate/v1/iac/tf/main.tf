# A launch template is the reusable blueprint EC2 fleets launch from. Only
# the template NAME is create-only in AWS: every other change creates a new
# immutable template VERSION, and update_default_version promotes it to the
# default -- so consumers following "$Default" (the common ASG and
# node-group setup) pick up the change on their next launch or instance
# refresh, while consumers pinned to a numeric version keep exactly what
# they tested.
resource "aws_launch_template" "this" {
  name = local.launch_template_name

  # Promote every new version to the template default. This is the
  # declarative-model contract: the spec describes ONE desired
  # configuration, so the newest version is always the intended one.
  update_default_version = true

  description   = var.spec.description != "" ? var.spec.description : null
  image_id      = var.spec.image_id != "" ? var.spec.image_id : null
  instance_type = var.spec.instance_type != "" ? var.spec.instance_type : null
  key_name      = var.spec.key_name != "" ? var.spec.key_name : null

  # The spec carries plain text so manifests stay readable; the EC2 API
  # requires base64 -- encode here, identically to the Pulumi module.
  user_data = var.spec.user_data != "" ? base64encode(var.spec.user_data) : null

  # ebs_optimized is a nullable tri-state at AWS (support varies by
  # instance type), so the provider takes a string; only an explicit true
  # is sent and unset keeps the type's own default.
  ebs_optimized = var.spec.ebs_optimized ? "true" : null

  vpc_security_group_ids = length(var.spec.security_group_ids) > 0 ? var.spec.security_group_ids : null

  disable_api_stop                     = var.spec.disable_api_stop ? true : null
  disable_api_termination              = var.spec.disable_api_termination ? true : null
  instance_initiated_shutdown_behavior = var.spec.instance_initiated_shutdown_behavior != "" ? var.spec.instance_initiated_shutdown_behavior : null

  dynamic "iam_instance_profile" {
    for_each = var.spec.instance_profile != "" ? [var.spec.instance_profile] : []
    content {
      arn = iam_instance_profile.value
    }
  }

  # Attribute-based instance selection: memory_mib and vcpu_count are the
  # two AWS-required dimensions (the spec enforces their presence); every
  # other field narrows the candidate set and is sent only when set so
  # AWS's own defaults keep applying.
  dynamic "instance_requirements" {
    for_each = var.spec.instance_requirements != null ? [var.spec.instance_requirements] : []
    content {
      memory_mib {
        min = instance_requirements.value.memory_mib.min
        max = instance_requirements.value.memory_mib.max > 0 ? instance_requirements.value.memory_mib.max : null
      }
      vcpu_count {
        min = instance_requirements.value.vcpu_count.min
        max = instance_requirements.value.vcpu_count.max > 0 ? instance_requirements.value.vcpu_count.max : null
      }

      allowed_instance_types  = length(instance_requirements.value.allowed_instance_types) > 0 ? instance_requirements.value.allowed_instance_types : null
      excluded_instance_types = length(instance_requirements.value.excluded_instance_types) > 0 ? instance_requirements.value.excluded_instance_types : null
      instance_generations    = length(instance_requirements.value.instance_generations) > 0 ? instance_requirements.value.instance_generations : null
      cpu_manufacturers       = length(instance_requirements.value.cpu_manufacturers) > 0 ? instance_requirements.value.cpu_manufacturers : null

      bare_metal                = instance_requirements.value.bare_metal != "" ? instance_requirements.value.bare_metal : null
      burstable_performance     = instance_requirements.value.burstable_performance != "" ? instance_requirements.value.burstable_performance : null
      require_hibernate_support = instance_requirements.value.require_hibernate_support ? true : null

      spot_max_price_percentage_over_lowest_price                    = instance_requirements.value.spot_max_price_percentage_over_lowest_price > 0 ? instance_requirements.value.spot_max_price_percentage_over_lowest_price : null
      max_spot_price_as_percentage_of_optimal_on_demand_price        = instance_requirements.value.max_spot_price_as_percentage_of_optimal_on_demand_price > 0 ? instance_requirements.value.max_spot_price_as_percentage_of_optimal_on_demand_price : null
      on_demand_max_price_percentage_over_lowest_price               = instance_requirements.value.on_demand_max_price_percentage_over_lowest_price > 0 ? instance_requirements.value.on_demand_max_price_percentage_over_lowest_price : null

      local_storage       = instance_requirements.value.local_storage != "" ? instance_requirements.value.local_storage : null
      local_storage_types = length(instance_requirements.value.local_storage_types) > 0 ? instance_requirements.value.local_storage_types : null

      dynamic "total_local_storage_gb" {
        for_each = instance_requirements.value.total_local_storage_gb != null ? [instance_requirements.value.total_local_storage_gb] : []
        content {
          min = total_local_storage_gb.value.min > 0 ? total_local_storage_gb.value.min : null
          max = total_local_storage_gb.value.max > 0 ? total_local_storage_gb.value.max : null
        }
      }
      dynamic "memory_gib_per_vcpu" {
        for_each = instance_requirements.value.memory_gib_per_vcpu != null ? [instance_requirements.value.memory_gib_per_vcpu] : []
        content {
          min = memory_gib_per_vcpu.value.min > 0 ? memory_gib_per_vcpu.value.min : null
          max = memory_gib_per_vcpu.value.max > 0 ? memory_gib_per_vcpu.value.max : null
        }
      }
      dynamic "network_interface_count" {
        for_each = instance_requirements.value.network_interface_count != null ? [instance_requirements.value.network_interface_count] : []
        content {
          min = network_interface_count.value.min > 0 ? network_interface_count.value.min : null
          max = network_interface_count.value.max > 0 ? network_interface_count.value.max : null
        }
      }
      dynamic "network_bandwidth_gbps" {
        for_each = instance_requirements.value.network_bandwidth_gbps != null ? [instance_requirements.value.network_bandwidth_gbps] : []
        content {
          min = network_bandwidth_gbps.value.min > 0 ? network_bandwidth_gbps.value.min : null
          max = network_bandwidth_gbps.value.max > 0 ? network_bandwidth_gbps.value.max : null
        }
      }
      dynamic "baseline_ebs_bandwidth_mbps" {
        for_each = instance_requirements.value.baseline_ebs_bandwidth_mbps != null ? [instance_requirements.value.baseline_ebs_bandwidth_mbps] : []
        content {
          min = baseline_ebs_bandwidth_mbps.value.min > 0 ? baseline_ebs_bandwidth_mbps.value.min : null
          max = baseline_ebs_bandwidth_mbps.value.max > 0 ? baseline_ebs_bandwidth_mbps.value.max : null
        }
      }
      dynamic "accelerator_count" {
        for_each = instance_requirements.value.accelerator_count != null ? [instance_requirements.value.accelerator_count] : []
        content {
          min = accelerator_count.value.min > 0 ? accelerator_count.value.min : null
          max = accelerator_count.value.max > 0 ? accelerator_count.value.max : null
        }
      }
      accelerator_manufacturers = length(instance_requirements.value.accelerator_manufacturers) > 0 ? instance_requirements.value.accelerator_manufacturers : null
      accelerator_names         = length(instance_requirements.value.accelerator_names) > 0 ? instance_requirements.value.accelerator_names : null
      accelerator_types         = length(instance_requirements.value.accelerator_types) > 0 ? instance_requirements.value.accelerator_types : null
      dynamic "accelerator_total_memory_mib" {
        for_each = instance_requirements.value.accelerator_total_memory_mib != null ? [instance_requirements.value.accelerator_total_memory_mib] : []
        content {
          min = accelerator_total_memory_mib.value.min > 0 ? accelerator_total_memory_mib.value.min : null
          max = accelerator_total_memory_mib.value.max > 0 ? accelerator_total_memory_mib.value.max : null
        }
      }
    }
  }

  # Unset EBS fields are omitted so the AMI's own mapping (size, type) and
  # the account default (encryption) keep deciding -- what makes a minimal
  # root-volume override safe.
  dynamic "block_device_mappings" {
    for_each = var.spec.block_device_mappings
    content {
      device_name  = block_device_mappings.value.device_name
      virtual_name = block_device_mappings.value.virtual_name != "" ? block_device_mappings.value.virtual_name : null
      # no_device is a presence-signal at AWS: an empty string suppresses
      # the AMI's device. The spec's bool maps onto that convention.
      no_device = block_device_mappings.value.no_device ? "" : null

      dynamic "ebs" {
        for_each = block_device_mappings.value.ebs != null ? [block_device_mappings.value.ebs] : []
        content {
          volume_size = ebs.value.volume_size_gb > 0 ? ebs.value.volume_size_gb : null
          volume_type = ebs.value.volume_type != "" ? ebs.value.volume_type : null
          iops        = ebs.value.iops > 0 ? ebs.value.iops : null
          throughput  = ebs.value.throughput_mibps > 0 ? ebs.value.throughput_mibps : null
          encrypted   = ebs.value.encrypted ? "true" : null
          kms_key_id  = ebs.value.kms_key_id != "" ? ebs.value.kms_key_id : null
          snapshot_id = ebs.value.snapshot_id != "" ? ebs.value.snapshot_id : null
          # delete_on_termination is a nullable tri-state at AWS (the AMI
          # mapping decides the default): null keeps the AMI default, an
          # explicit value overrides it.
          delete_on_termination = ebs.value.delete_on_termination == null ? null : tostring(ebs.value.delete_on_termination)
        }
      }
    }
  }

  dynamic "network_interfaces" {
    for_each = var.spec.network_interfaces
    content {
      device_index         = network_interfaces.value.device_index
      network_card_index   = network_interfaces.value.network_card_index > 0 ? network_interfaces.value.network_card_index : null
      description          = network_interfaces.value.description != "" ? network_interfaces.value.description : null
      interface_type       = network_interfaces.value.interface_type != "" ? network_interfaces.value.interface_type : null
      network_interface_id = network_interfaces.value.network_interface_id != "" ? network_interfaces.value.network_interface_id : null

      # associate_public_ip_address and delete_on_termination are nullable
      # tri-states at AWS (the subnet / interface origin decides the
      # default): null inherits, an explicit value overrides.
      associate_public_ip_address = network_interfaces.value.associate_public_ip_address == null ? null : tostring(network_interfaces.value.associate_public_ip_address)
      delete_on_termination       = network_interfaces.value.delete_on_termination == null ? null : tostring(network_interfaces.value.delete_on_termination)

      subnet_id       = network_interfaces.value.subnet_id != "" ? network_interfaces.value.subnet_id : null
      security_groups = length(network_interfaces.value.security_group_ids) > 0 ? network_interfaces.value.security_group_ids : null

      private_ip_address = network_interfaces.value.private_ip_address != "" ? network_interfaces.value.private_ip_address : null
      ipv4_address_count = network_interfaces.value.ipv4_address_count > 0 ? network_interfaces.value.ipv4_address_count : null
      ipv4_addresses     = length(network_interfaces.value.ipv4_addresses) > 0 ? network_interfaces.value.ipv4_addresses : null
      ipv6_address_count = network_interfaces.value.ipv6_address_count > 0 ? network_interfaces.value.ipv6_address_count : null
      ipv6_addresses     = length(network_interfaces.value.ipv6_addresses) > 0 ? network_interfaces.value.ipv6_addresses : null
      ipv4_prefix_count  = network_interfaces.value.ipv4_prefix_count > 0 ? network_interfaces.value.ipv4_prefix_count : null
      ipv4_prefixes      = length(network_interfaces.value.ipv4_prefixes) > 0 ? network_interfaces.value.ipv4_prefixes : null
      ipv6_prefix_count  = network_interfaces.value.ipv6_prefix_count > 0 ? network_interfaces.value.ipv6_prefix_count : null
      ipv6_prefixes      = length(network_interfaces.value.ipv6_prefixes) > 0 ? network_interfaces.value.ipv6_prefixes : null
    }
  }

  # Only explicitly set IMDS fields are sent, so AWS keeps its own defaults
  # for the rest.
  dynamic "metadata_options" {
    for_each = var.spec.metadata_options != null ? [var.spec.metadata_options] : []
    content {
      http_endpoint               = metadata_options.value.http_endpoint != "" ? metadata_options.value.http_endpoint : null
      http_tokens                 = metadata_options.value.http_tokens != "" ? metadata_options.value.http_tokens : null
      http_put_response_hop_limit = metadata_options.value.http_put_response_hop_limit > 0 ? metadata_options.value.http_put_response_hop_limit : null
      http_protocol_ipv6          = metadata_options.value.http_protocol_ipv6 != "" ? metadata_options.value.http_protocol_ipv6 : null
      instance_metadata_tags      = metadata_options.value.instance_metadata_tags != "" ? metadata_options.value.instance_metadata_tags : null
    }
  }

  dynamic "monitoring" {
    for_each = var.spec.detailed_monitoring ? [true] : []
    content {
      enabled = true
    }
  }

  dynamic "placement" {
    for_each = var.spec.placement != null ? [var.spec.placement] : []
    content {
      availability_zone = placement.value.availability_zone != "" ? placement.value.availability_zone : null
      group_name        = placement.value.group_name != "" ? placement.value.group_name : null
      partition_number  = placement.value.partition_number > 0 ? placement.value.partition_number : null
      tenancy           = placement.value.tenancy != "" ? placement.value.tenancy : null
    }
  }

  dynamic "cpu_options" {
    for_each = var.spec.cpu_options != null ? [var.spec.cpu_options] : []
    content {
      core_count       = cpu_options.value.core_count > 0 ? cpu_options.value.core_count : null
      threads_per_core = cpu_options.value.threads_per_core > 0 ? cpu_options.value.threads_per_core : null
      amd_sev_snp      = cpu_options.value.amd_sev_snp != "" ? cpu_options.value.amd_sev_snp : null
    }
  }

  dynamic "credit_specification" {
    for_each = var.spec.cpu_credits != "" ? [var.spec.cpu_credits] : []
    content {
      cpu_credits = credit_specification.value
    }
  }

  # Configuring spot_options makes every launch a Spot request; the market
  # type is implied by the block's presence.
  dynamic "instance_market_options" {
    for_each = var.spec.spot_options != null ? [var.spec.spot_options] : []
    content {
      market_type = "spot"
      spot_options {
        max_price                      = instance_market_options.value.max_price != "" ? instance_market_options.value.max_price : null
        spot_instance_type             = instance_market_options.value.spot_instance_type != "" ? instance_market_options.value.spot_instance_type : null
        instance_interruption_behavior = instance_market_options.value.instance_interruption_behavior != "" ? instance_market_options.value.instance_interruption_behavior : null
        valid_until                    = instance_market_options.value.valid_until != "" ? instance_market_options.value.valid_until : null
      }
    }
  }

  dynamic "enclave_options" {
    for_each = var.spec.enclave_enabled ? [true] : []
    content {
      enabled = true
    }
  }

  dynamic "hibernation_options" {
    for_each = var.spec.hibernation_enabled ? [true] : []
    content {
      configured = true
    }
  }

  dynamic "maintenance_options" {
    for_each = var.spec.auto_recovery != "" ? [var.spec.auto_recovery] : []
    content {
      auto_recovery = maintenance_options.value
    }
  }

  dynamic "private_dns_name_options" {
    for_each = var.spec.private_dns_name_options != null ? [var.spec.private_dns_name_options] : []
    content {
      hostname_type                        = private_dns_name_options.value.hostname_type != "" ? private_dns_name_options.value.hostname_type : null
      enable_resource_name_dns_a_record    = private_dns_name_options.value.enable_resource_name_dns_a_record ? true : null
      enable_resource_name_dns_aaaa_record = private_dns_name_options.value.enable_resource_name_dns_aaaa_record ? true : null
    }
  }

  # Identity tags land in three places on purpose: on the template itself,
  # and via tag_specifications on the instances and volumes each launch
  # creates (template tags do not propagate to launched resources).
  tag_specifications {
    resource_type = "instance"
    tags          = local.aws_tags
  }
  tag_specifications {
    resource_type = "volume"
    tags          = local.aws_tags
  }

  tags = local.aws_tags
}
