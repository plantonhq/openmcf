# Create the Entra access grant on the Mongo vCore cluster.
# "MicrosoftEntraID" is the identity provider's only legal value today
# -- deliberately not part of the spec; both engines send it
# explicitly. Every argument is create-only (the resource has no update
# path -- changing anything replaces the grant, a harmless drop-and-
# re-add of an access binding); the target cluster must allow
# MicrosoftEntraID authentication (deploy-time contract, documented on
# the spec).
resource "azurerm_mongo_cluster_user" "main" {
  object_id              = local.object_id
  mongo_cluster_id       = var.spec.mongo_cluster_id
  identity_provider_type = "MicrosoftEntraID"
  principal_type         = var.spec.principal_type

  dynamic "role" {
    for_each = var.spec.roles
    content {
      database = role.value.database
      name     = role.value.role
    }
  }
}
