# Zero Trust infrastructure target: the hostname/IP inventory row that
# Access infrastructure applications select and broker SSH access to. A
# plain CRUD resource: real create, in-place updates (hostname, addresses,
# virtual networks), real delete. Only the account forces replacement.
#
# An omitted virtual_network_id is not sent -- Cloudflare then assigns the
# account's default virtual network (the attribute is computed on the
# provider side, so the assigned value never drifts).
resource "cloudflare_zero_trust_access_infrastructure_target" "main" {
  account_id = var.spec.account_id
  hostname   = var.spec.hostname

  ip = {
    ipv4 = try(var.spec.ip.ipv4, null) != null ? {
      ip_addr            = var.spec.ip.ipv4.ip_addr
      virtual_network_id = try(var.spec.ip.ipv4.virtual_network_id, "") != "" ? var.spec.ip.ipv4.virtual_network_id : null
    } : null
    ipv6 = try(var.spec.ip.ipv6, null) != null ? {
      ip_addr            = var.spec.ip.ipv6.ip_addr
      virtual_network_id = try(var.spec.ip.ipv6.virtual_network_id, "") != "" ? var.spec.ip.ipv6.virtual_network_id : null
    } : null
  }
}
