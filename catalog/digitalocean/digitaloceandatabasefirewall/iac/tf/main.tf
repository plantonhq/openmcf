# DigitalOcean Database Firewall
#
# Provisions the inbound trusted-sources rule set of a DigitalOcean managed
# database cluster -- the complete digitalocean_database_firewall resource
# surface. The spec's five typed lists (IPs/CIDRs, Droplets, Kubernetes
# clusters, apps, Droplet tags) fan out to the provider's polymorphic
# {type, value} rows in locals.tf.
#
# Semantics worth knowing before editing:
#   - The rule set is a PROPERTY of the cluster: at most one per cluster,
#     and every update replaces the FULL set (the API is a PUT).
#   - "Destroy" does not delete an object -- it PUTs an EMPTY rule list,
#     after which the cluster accepts connections from anywhere again.
#   - Each rule row carries computed uuid/created_at noise; the rule set
#     hashes as a set, so server-assigned fields never cause diffs.
#   - The Terraform state id is a random unique string minted at create
#     (not stable across imports); the cluster UUID is the real identity.

resource "digitalocean_database_firewall" "firewall" {
  cluster_id = var.spec.cluster

  dynamic "rule" {
    for_each = local.rules
    content {
      type  = rule.value.type
      value = rule.value.value
    }
  }
}
