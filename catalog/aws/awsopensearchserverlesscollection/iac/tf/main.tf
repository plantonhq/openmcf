# Amazon OpenSearch Serverless: a fully managed, auto-scaling OpenSearch
# workspace. The collection is inseparable from three account-level POLICY
# objects that attach by name-pattern matching -- this module renders each
# policy scoped to exactly this collection ("collection/<name>",
# "index/<name>/..."), so one manifest owns one collection and everything
# that makes it usable.

# The encryption policy MUST exist before the collection (AWS rejects
# CreateCollection without a matching encryption policy). Policy names share
# the collection's name -- types are separate namespaces, so the encryption
# and network policies may share it.
resource "aws_opensearchserverless_security_policy" "encryption" {
  name        = local.collection_name
  type        = "encryption"
  description = "Encryption for collection ${local.collection_name}"
  policy      = local.encryption_policy_json
}

resource "aws_opensearchserverless_collection" "this" {
  # Create-time immutable; doubles as the Name tag. metadata.name on both
  # engines.
  name        = local.collection_name
  description = var.spec.description != "" ? var.spec.description : null

  # Workload type and standby replicas (both ForceNew; defaults TIMESERIES /
  # ENABLED are materialized by the manifest loader -- sent explicitly so the
  # manifest's intent is visible in state on both engines).
  type             = var.spec.type
  standby_replicas = var.spec.standby_replicas

  # Collection-group membership (ForceNew). The group's standby-replicas
  # setting must match the collection's -- AWS rejects the mismatch at
  # create.
  collection_group_name = var.spec.collection_group_name != "" ? var.spec.collection_group_name : null

  # GPU-accelerated vector capacity -- VECTORSEARCH collections only (CEL
  # enforces the coupling at manifest time).
  dynamic "vector_options" {
    for_each = var.spec.serverless_vector_acceleration != "" ? [1] : []
    content {
      serverless_vector_acceleration = var.spec.serverless_vector_acceleration
    }
  }

  # The encryption KEY CHOICE lives on the module-rendered encryption
  # security policy above. The collection's own inline encryption_config
  # argument is the same setting's collection-group-era twin and is
  # deliberately not sent (recorded in the parity manifest).

  tags = local.aws_tags

  # Create after the matching encryption policy; destroy before it.
  depends_on = [aws_opensearchserverless_security_policy.encryption]
}

# Network reachability of the collection and Dashboards endpoints. Attaches
# by name pattern -- no create-order requirement against the collection.
resource "aws_opensearchserverless_security_policy" "network" {
  name        = local.collection_name
  type        = "network"
  description = "Network access for collection ${local.collection_name}"
  policy      = local.network_policy_json
}

# Data-plane permissions. Without at least one rule nothing can read or
# write data (IAM permissions alone grant nothing in OpenSearch
# Serverless) -- the spec comment says so loudly, so the skip is honest.
resource "aws_opensearchserverless_access_policy" "data" {
  count = length(var.spec.data_access) > 0 ? 1 : 0

  name        = local.collection_name
  type        = "data"
  description = "Data access for collection ${local.collection_name}"
  policy      = jsonencode(local.data_access_policy_document)
}

# Index retention. Skipped when no rules are declared (indexes are then
# retained indefinitely, AWS's default).
resource "aws_opensearchserverless_lifecycle_policy" "retention" {
  count = length(var.spec.retention_rules) > 0 ? 1 : 0

  name        = local.collection_name
  type        = "retention"
  description = "Index retention for collection ${local.collection_name}"
  policy      = local.lifecycle_policy_json
}
