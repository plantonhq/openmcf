# DigitalOcean Database User
#
# Provisions an additional user on a DigitalOcean managed database cluster,
# modeling the complete digitalocean_database_user resource surface: the
# MySQL authentication plugin choice and the Kafka / OpenSearch access
# control lists. DigitalOcean generates the password (and, on Kafka, the
# mTLS certificate pair) server-side; they surface only as outputs.
#
# The DigitalOcean API serializes user creation/deletion per cluster, so
# composing many users on one cluster deploys sequentially by design.

resource "digitalocean_database_user" "user" {
  cluster_id = var.spec.cluster
  name       = var.spec.user_name

  # MySQL clusters only (API-enforced). Unset defers to DigitalOcean's
  # caching_sha2_password default; updates apply through a
  # password-preserving auth reset.
  mysql_auth_plugin = local.mysql_auth_plugin

  # Engine-specific ACLs. DigitalOcean returns these only in the CREATE
  # response -- reads never include them -- so this configuration is the
  # source of truth and imports can never recover it (recorded as a
  # config-only import tolerance). Each ACL row also carries a computed
  # server-side id, which is provisioning noise and not modeled.
  dynamic "settings" {
    for_each = var.spec.settings != null ? [var.spec.settings] : []
    content {
      dynamic "acl" {
        for_each = settings.value.kafka_acls
        content {
          topic      = acl.value.topic
          permission = acl.value.permission
        }
      }
      dynamic "opensearch_acl" {
        for_each = settings.value.opensearch_acls
        content {
          index      = opensearch_acl.value.index
          permission = opensearch_acl.value.permission
        }
      }
    }
  }
}
