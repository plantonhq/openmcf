# The EC2 instance. Launch identity (AMI, type, subnet, key pair,
# placement, CPU topology, purchase option) is create-time in EC2 --
# changing those fields replaces or restarts the instance -- while the
# operational posture (security groups, IAM profile, metadata options,
# protections, monitoring) edits in place.
#
# Field-presence discipline: the generated variable contract keeps
# tri-state optionals null when unset, and every argument below passes
# null through so an omitted field keeps the AWS (or launch-template)
# default instead of forcing a zero value. Nested optional OBJECTS are
# guarded with ternaries, never `!= null &&` -- HCL's && does not
# short-circuit and errors on the null dereference.
resource "aws_instance" "this" {
  # AMI and type are each optional when the launch template supplies
  # them (CEL enforces that at validate time); empty strings lower to
  # null so the template's values win.
  ami           = var.spec.ami != "" ? var.spec.ami : null
  instance_type = var.spec.instance_type != "" ? var.spec.instance_type : null

  # Launch from a referenced template; inline fields set on this
  # resource override the template's values -- the template is the
  # baseline, this instance's spec is the deviation.
  dynamic "launch_template" {
    for_each = var.spec.launch_template != null ? [var.spec.launch_template] : []
    content {
      id      = launch_template.value.id != "" ? launch_template.value.id : null
      name    = launch_template.value.name != "" ? launch_template.value.name : null
      version = launch_template.value.version != "" ? launch_template.value.version : null
    }
  }

  # The EC2 instance API takes the profile by NAME (launch templates
  # accept an ARN; instances do not) -- the spec's ref resolves the
  # profile's name output for exactly this reason.
  iam_instance_profile = var.spec.instance_profile != "" ? var.spec.instance_profile : null
  key_name             = var.spec.key_name != "" ? var.spec.key_name : null

  # Either a pre-provisioned primary ENI carries the network identity,
  # or the inline fields shape a new primary interface (CEL guarantees
  # the two never mix).
  dynamic "primary_network_interface" {
    for_each = var.spec.primary_network_interface_id != "" ? [1] : []
    content {
      network_interface_id = var.spec.primary_network_interface_id
    }
  }

  subnet_id              = var.spec.subnet_id != "" ? var.spec.subnet_id : null
  vpc_security_group_ids = length(var.spec.security_group_ids) > 0 ? var.spec.security_group_ids : null
  private_ip             = var.spec.private_ip != "" ? var.spec.private_ip : null
  secondary_private_ips  = length(var.spec.secondary_private_ips) > 0 ? var.spec.secondary_private_ips : null

  # Tri-state: null inherits the subnet's map-public-IP setting.
  associate_public_ip_address = var.spec.associate_public_ip_address

  # Platform middleware defaults this to true (the AWS default); an
  # explicit false is the NAT/router posture.
  source_dest_check = var.spec.source_dest_check

  ipv6_address_count  = var.spec.ipv6_address_count != 0 ? var.spec.ipv6_address_count : null
  ipv6_addresses      = length(var.spec.ipv6_addresses) > 0 ? var.spec.ipv6_addresses : null
  enable_primary_ipv6 = var.spec.enable_primary_ipv6

  dynamic "private_dns_name_options" {
    for_each = var.spec.private_dns_name_options != null ? [var.spec.private_dns_name_options] : []
    content {
      hostname_type                        = private_dns_name_options.value.hostname_type != "" ? private_dns_name_options.value.hostname_type : null
      enable_resource_name_dns_a_record    = private_dns_name_options.value.enable_resource_name_dns_a_record
      enable_resource_name_dns_aaaa_record = private_dns_name_options.value.enable_resource_name_dns_aaaa_record
    }
  }

  # Additional interfaces on non-primary network cards, for
  # high-bandwidth instance types with multiple cards. Security groups
  # are not configurable at launch on these (the EC2 API applies the
  # VPC default group).
  dynamic "secondary_network_interface" {
    for_each = var.spec.secondary_network_interfaces
    content {
      network_card_index       = secondary_network_interface.value.network_card_index
      device_index             = secondary_network_interface.value.device_index != 0 ? secondary_network_interface.value.device_index : null
      secondary_subnet_id      = secondary_network_interface.value.subnet_id
      private_ip_address_count = secondary_network_interface.value.private_ip_address_count != 0 ? secondary_network_interface.value.private_ip_address_count : null
      delete_on_termination    = secondary_network_interface.value.delete_on_termination
    }
  }

  # Root-volume override: unset fields stay null so the AMI's block
  # device mapping keeps deciding them.
  dynamic "root_block_device" {
    for_each = var.spec.root_block_device != null ? [var.spec.root_block_device] : []
    content {
      volume_size           = root_block_device.value.volume_size_gb != 0 ? root_block_device.value.volume_size_gb : null
      volume_type           = root_block_device.value.volume_type != "" ? root_block_device.value.volume_type : null
      iops                  = root_block_device.value.iops != 0 ? root_block_device.value.iops : null
      throughput            = root_block_device.value.throughput_mibps != 0 ? root_block_device.value.throughput_mibps : null
      encrypted             = root_block_device.value.encrypted ? true : null
      kms_key_id            = root_block_device.value.kms_key_id != "" ? root_block_device.value.kms_key_id : null
      delete_on_termination = root_block_device.value.delete_on_termination
      # Per-volume tags apply post-creation (see volume_tags for the
      # at-creation alternative; the provider forbids mixing them).
      tags = length(root_block_device.value.tags) > 0 ? root_block_device.value.tags : null
    }
  }

  dynamic "ebs_block_device" {
    for_each = var.spec.ebs_block_devices
    content {
      device_name           = ebs_block_device.value.device_name
      volume_size           = ebs_block_device.value.volume_size_gb != 0 ? ebs_block_device.value.volume_size_gb : null
      volume_type           = ebs_block_device.value.volume_type != "" ? ebs_block_device.value.volume_type : null
      iops                  = ebs_block_device.value.iops != 0 ? ebs_block_device.value.iops : null
      throughput            = ebs_block_device.value.throughput_mibps != 0 ? ebs_block_device.value.throughput_mibps : null
      encrypted             = ebs_block_device.value.encrypted ? true : null
      kms_key_id            = ebs_block_device.value.kms_key_id != "" ? ebs_block_device.value.kms_key_id : null
      snapshot_id           = ebs_block_device.value.snapshot_id != "" ? ebs_block_device.value.snapshot_id : null
      delete_on_termination = ebs_block_device.value.delete_on_termination
      tags                  = length(ebs_block_device.value.tags) > 0 ? ebs_block_device.value.tags : null
    }
  }

  # Uniform at-creation tags for EVERY volume (incl. AMI-inherited
  # mappings) -- the ABAC-compliant arm; mutually exclusive with the
  # per-device tags above.
  volume_tags = length(var.spec.volume_tags) > 0 ? var.spec.volume_tags : null

  dynamic "ephemeral_block_device" {
    for_each = var.spec.ephemeral_block_devices
    content {
      device_name  = ephemeral_block_device.value.device_name
      virtual_name = ephemeral_block_device.value.virtual_name != "" ? ephemeral_block_device.value.virtual_name : null
      no_device    = ephemeral_block_device.value.no_device ? true : null
    }
  }

  # Only forced on: most current-generation types are EBS-optimized by
  # default at no charge, so an omitted field keeps the type's default
  # rather than pinning false.
  ebs_optimized = var.spec.ebs_optimized ? true : null

  # IMDSv2 posture: http_tokens = "required" is the recommended
  # hardening for every new instance.
  dynamic "metadata_options" {
    for_each = var.spec.metadata_options != null ? [var.spec.metadata_options] : []
    content {
      http_endpoint               = metadata_options.value.http_endpoint != "" ? metadata_options.value.http_endpoint : null
      http_tokens                 = metadata_options.value.http_tokens != "" ? metadata_options.value.http_tokens : null
      http_put_response_hop_limit = metadata_options.value.http_put_response_hop_limit != 0 ? metadata_options.value.http_put_response_hop_limit : null
      http_protocol_ipv6          = metadata_options.value.http_protocol_ipv6 != "" ? metadata_options.value.http_protocol_ipv6 : null
      instance_metadata_tags      = metadata_options.value.instance_metadata_tags != "" ? metadata_options.value.instance_metadata_tags : null
    }
  }

  monitoring = var.spec.detailed_monitoring ? true : null

  # CPU topology is fixed at launch; unset fields keep the type's
  # defaults (e.g. hyper-threading on x86).
  dynamic "cpu_options" {
    for_each = var.spec.cpu_options != null ? [var.spec.cpu_options] : []
    content {
      core_count            = cpu_options.value.core_count != 0 ? cpu_options.value.core_count : null
      threads_per_core      = cpu_options.value.threads_per_core != 0 ? cpu_options.value.threads_per_core : null
      amd_sev_snp           = cpu_options.value.amd_sev_snp != "" ? cpu_options.value.amd_sev_snp : null
      nested_virtualization = cpu_options.value.nested_virtualization != "" ? cpu_options.value.nested_virtualization : null
    }
  }

  dynamic "credit_specification" {
    for_each = var.spec.cpu_credits != "" ? [1] : []
    content {
      cpu_credits = var.spec.cpu_credits
    }
  }

  # The purchase market: an explicit market_type, or presence of
  # spot_options implying "spot" (the classic shape). On-Demand needs no
  # market options at all. Reservation-backed markets (capacity-block,
  # interruptible-capacity-reservation) carry no spot_options -- CEL
  # enforces that pairing and the required reservation target.
  dynamic "instance_market_options" {
    for_each = var.spec.market_type != "" || var.spec.spot_options != null ? [1] : []
    content {
      market_type = var.spec.market_type != "" ? var.spec.market_type : "spot"
      dynamic "spot_options" {
        for_each = var.spec.spot_options != null ? [var.spec.spot_options] : []
        content {
          max_price                      = spot_options.value.max_price != "" ? spot_options.value.max_price : null
          spot_instance_type             = spot_options.value.spot_instance_type != "" ? spot_options.value.spot_instance_type : null
          instance_interruption_behavior = spot_options.value.instance_interruption_behavior != "" ? spot_options.value.instance_interruption_behavior : null
          valid_until                    = spot_options.value.valid_until != "" ? spot_options.value.valid_until : null
        }
      }
    }
  }

  dynamic "capacity_reservation_specification" {
    for_each = var.spec.capacity_reservation != null ? [var.spec.capacity_reservation] : []
    content {
      capacity_reservation_preference = capacity_reservation_specification.value.preference != "" ? capacity_reservation_specification.value.preference : null

      dynamic "capacity_reservation_target" {
        for_each = capacity_reservation_specification.value.capacity_reservation_id != "" || capacity_reservation_specification.value.capacity_reservation_resource_group_arn != "" ? [1] : []
        content {
          capacity_reservation_id                 = capacity_reservation_specification.value.capacity_reservation_id != "" ? capacity_reservation_specification.value.capacity_reservation_id : null
          capacity_reservation_resource_group_arn = capacity_reservation_specification.value.capacity_reservation_resource_group_arn != "" ? capacity_reservation_specification.value.capacity_reservation_resource_group_arn : null
        }
      }
    }
  }

  # Placement is fixed at launch. Ternary guards throughout: the whole
  # placement object may be pruned from tfvars.
  availability_zone          = var.spec.placement != null ? (var.spec.placement.availability_zone != "" ? var.spec.placement.availability_zone : null) : null
  placement_group            = var.spec.placement != null ? (var.spec.placement.group_name != "" ? var.spec.placement.group_name : null) : null
  placement_group_id         = var.spec.placement != null ? (var.spec.placement.group_id != "" ? var.spec.placement.group_id : null) : null
  placement_partition_number = var.spec.placement != null ? (var.spec.placement.partition_number != 0 ? var.spec.placement.partition_number : null) : null
  tenancy                    = var.spec.placement != null ? (var.spec.placement.tenancy != "" ? var.spec.placement.tenancy : null) : null
  host_id                    = var.spec.placement != null ? (var.spec.placement.host_id != "" ? var.spec.placement.host_id : null) : null
  host_resource_group_arn    = var.spec.placement != null ? (var.spec.placement.host_resource_group_arn != "" ? var.spec.placement.host_resource_group_arn : null) : null

  # Nitro Enclaves and hibernation are mutually exclusive (CEL rejects
  # the combination before AWS would).
  dynamic "enclave_options" {
    for_each = var.spec.enclave_enabled ? [1] : []
    content {
      enabled = true
    }
  }
  hibernation = var.spec.hibernation_enabled ? true : null

  dynamic "maintenance_options" {
    for_each = var.spec.auto_recovery != "" ? [1] : []
    content {
      auto_recovery = var.spec.auto_recovery
    }
  }

  instance_initiated_shutdown_behavior = var.spec.instance_initiated_shutdown_behavior != "" ? var.spec.instance_initiated_shutdown_behavior : null
  disable_api_stop                     = var.spec.disable_api_stop
  disable_api_termination              = var.spec.disable_api_termination

  # Destroy-time escape hatch: lift stop/termination protection before
  # terminating instead of failing the destroy.
  force_destroy = var.spec.force_destroy ? true : null

  # Plain text arrives with HCL template introducers already escaped by
  # the tfvars writer, so ${...} content in shell scripts passes through
  # literally -- never add escaping here.
  user_data                   = var.spec.user_data != "" ? var.spec.user_data : null
  user_data_base64            = var.spec.user_data_base64 != "" ? var.spec.user_data_base64 : null
  user_data_replace_on_change = var.spec.user_data_replace_on_change

  tags = local.aws_tags
}
