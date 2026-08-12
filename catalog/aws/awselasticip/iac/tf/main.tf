resource "aws_eip" "this" {
  # Every managed EIP is a VPC EIP: AWS retired EC2-Classic ("standard"
  # domain) addresses and the provider itself refuses to read or delete
  # them, so the legacy domain is deliberately not configurable.
  domain = "vpc"

  # BYOIP: allocate from a specific IPv4 address pool.
  public_ipv4_pool = local.public_ipv4_pool

  # BYOIP / IPAM: request a specific IP address the pool holds.
  address = local.address

  # IPAM: allocate from an Amazon VPC IP Address Manager public pool.
  ipam_pool_id = local.ipam_pool_id

  # Location scope for Local Zones and Wavelength zones.
  network_border_group = local.network_border_group

  # Association: attach the address to an instance XOR a network interface
  # (spec CEL enforces at-most-one; AWS associates with exactly one target).
  instance                  = local.instance
  network_interface         = local.network_interface
  associate_with_private_ip = local.associate_with_private_ip

  tags = local.tags
}

# Reverse DNS (PTR) record for the address. AWS validates SERVER-SIDE that a
# forward A record for the domain already resolves to this EIP before
# granting the PTR — a fresh EIP therefore typically sets this on a
# follow-up apply, after DNS points at the address.
resource "aws_eip_domain_name" "this" {
  count = local.reverse_dns_domain_name != null ? 1 : 0

  allocation_id = aws_eip.this.allocation_id
  domain_name   = local.reverse_dns_domain_name
}
