locals {
  resource_name = coalesce(try(var.metadata.name, null), "awselasticip")

  tags = merge({
    "Name" = local.resource_name
  }, try(var.metadata.labels, {}))

  # Allocation-shaping settings — null when not configured.
  public_ipv4_pool     = try(var.spec.public_ipv4_pool, null) != "" ? try(var.spec.public_ipv4_pool, null) : null
  address              = try(var.spec.address, null) != "" ? try(var.spec.address, null) : null
  network_border_group = try(var.spec.network_border_group, null) != "" ? try(var.spec.network_border_group, null) : null
  ipam_pool_id         = try(var.spec.ipam_pool_id, null) != "" ? try(var.spec.ipam_pool_id, null) : null

  # Association settings (references arrive resolved to literal ids).
  instance                  = try(var.spec.instance, null) != "" ? try(var.spec.instance, null) : null
  network_interface         = try(var.spec.network_interface, null) != "" ? try(var.spec.network_interface, null) : null
  associate_with_private_ip = try(var.spec.associate_with_private_ip, null) != "" ? try(var.spec.associate_with_private_ip, null) : null

  # Reverse DNS (PTR) — drives the count-gated aws_eip_domain_name resource.
  reverse_dns_domain_name = try(var.spec.reverse_dns_domain_name, null) != "" ? try(var.spec.reverse_dns_domain_name, null) : null
}
