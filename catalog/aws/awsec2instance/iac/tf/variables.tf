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
  description = "AwsEc2Instance specification"
  type = object({
    region = string
    ami = optional(string, "")
    instance_type = optional(string, "")
    launch_template = optional(object({
      id = optional(string, "")
      name = optional(string, "")
      version = optional(string, "")
    }))
    instance_profile = optional(string, "")
    key_name = optional(string, "")
    subnet_id = optional(string, "")
    security_group_ids = optional(list(string), [])
    primary_network_interface_id = optional(string, "")
    private_ip = optional(string, "")
    secondary_private_ips = optional(list(string), [])
    associate_public_ip_address = optional(bool)
    source_dest_check = optional(bool)
    ipv6_address_count = optional(number, 0)
    ipv6_addresses = optional(list(string), [])
    enable_primary_ipv6 = optional(bool)
    private_dns_name_options = optional(object({
      hostname_type = optional(string, "")
      enable_resource_name_dns_a_record = optional(bool, false)
      enable_resource_name_dns_aaaa_record = optional(bool, false)
    }))
    secondary_network_interfaces = optional(list(object({
      network_card_index = optional(number, 0)
      device_index = optional(number, 0)
      subnet_id = string
      private_ip_address_count = optional(number, 0)
      delete_on_termination = optional(bool)
    })), [])
    root_block_device = optional(object({
      volume_size_gb = optional(number, 0)
      volume_type = optional(string, "")
      iops = optional(number, 0)
      throughput_mibps = optional(number, 0)
      encrypted = optional(bool, false)
      kms_key_id = optional(string, "")
      delete_on_termination = optional(bool)
    }))
    ebs_block_devices = optional(list(object({
      device_name = string
      volume_size_gb = optional(number, 0)
      volume_type = optional(string, "")
      iops = optional(number, 0)
      throughput_mibps = optional(number, 0)
      encrypted = optional(bool, false)
      kms_key_id = optional(string, "")
      snapshot_id = optional(string, "")
      delete_on_termination = optional(bool)
    })), [])
    ephemeral_block_devices = optional(list(object({
      device_name = string
      virtual_name = optional(string, "")
      no_device = optional(bool, false)
    })), [])
    ebs_optimized = optional(bool, false)
    metadata_options = optional(object({
      http_endpoint = optional(string, "")
      http_tokens = optional(string, "")
      http_put_response_hop_limit = optional(number, 0)
      http_protocol_ipv6 = optional(string, "")
      instance_metadata_tags = optional(string, "")
    }))
    detailed_monitoring = optional(bool, false)
    cpu_options = optional(object({
      core_count = optional(number, 0)
      threads_per_core = optional(number, 0)
      amd_sev_snp = optional(string, "")
      nested_virtualization = optional(string, "")
    }))
    cpu_credits = optional(string, "")
    spot_options = optional(object({
      max_price = optional(string, "")
      spot_instance_type = optional(string, "")
      instance_interruption_behavior = optional(string, "")
      valid_until = optional(string, "")
    }))
    capacity_reservation = optional(object({
      preference = optional(string, "")
      capacity_reservation_id = optional(string, "")
      capacity_reservation_resource_group_arn = optional(string, "")
    }))
    placement = optional(object({
      availability_zone = optional(string, "")
      group_name = optional(string, "")
      group_id = optional(string, "")
      partition_number = optional(number, 0)
      tenancy = optional(string, "")
      host_id = optional(string, "")
      host_resource_group_arn = optional(string, "")
    }))
    enclave_enabled = optional(bool, false)
    hibernation_enabled = optional(bool, false)
    auto_recovery = optional(string, "")
    instance_initiated_shutdown_behavior = optional(string, "")
    disable_api_stop = optional(bool)
    disable_api_termination = optional(bool)
    user_data = optional(string, "")
    user_data_base64 = optional(string, "")
    user_data_replace_on_change = optional(bool, false)
  })
}
