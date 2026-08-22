# DigitalOcean VPC Peering
#
# Provisions a private-network peering connection between exactly two VPCs
# -- the complete digitalocean_vpc_peering resource surface. The spec
# models the two members as named references (the API requires exactly
# two); the provider's unordered vpc_ids set is assembled here.
#
# Only the name updates in place; changing either VPC replaces the
# peering. The provider waits for ACTIVE on create and retries the delete
# while DigitalOcean returns 403 during peering settlement (both within
# 2-minute timeouts).

resource "digitalocean_vpc_peering" "peering" {
  name = var.spec.peering_name

  # References resolve to literal VPC UUIDs before the module runs. The
  # provider treats the pair as an unordered set, so member order never
  # produces a diff.
  vpc_ids = [
    var.spec.vpc_1,
    var.spec.vpc_2,
  ]
}
