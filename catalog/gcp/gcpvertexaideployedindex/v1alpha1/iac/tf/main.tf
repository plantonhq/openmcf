# The deployed index — the placement joining a Vector Search index to an
# index endpoint, with its own serving compute. Deploying is a
# long-running operation (tens of minutes; provider timeout 45). Only the
# replica bounds inside the sizing arm update in place (the provider
# PATCHes them via mutateDeployedIndex); every other change undeploys and
# redeploys.
#
# No aiplatform API enablement here, deliberately: a deployment cannot
# exist without its index endpoint, and creating the endpoint (its own
# kind) already enabled the API. This resource also carries no project
# field — the project rides inside the index_endpoint resource path.
resource "google_vertex_ai_index_endpoint_deployed_index" "this" {
  deployed_index_id = local.deployed_index_id
  index             = var.spec.index
  index_endpoint    = var.spec.index_endpoint

  # The provider resolves the regional Vertex AI API host
  # (https://{region}-aiplatform.googleapis.com) from `region`; without
  # it the deploy fails unless the provider config happens to carry one.
  # Must match the endpoint's own region — deployments cannot cross
  # regions.
  region = local.location

  # Unusually for a display name, the API treats it as immutable on a
  # deployed index.
  display_name = var.spec.display_name != "" ? var.spec.display_name : null

  # Vertex-managed serving compute. 0 means "not set": GCP then applies
  # its defaults (min 2, max = min).
  dynamic "automatic_resources" {
    for_each = var.spec.automatic_resources != null ? [var.spec.automatic_resources] : []
    content {
      min_replica_count = automatic_resources.value.min_replica_count > 0 ? automatic_resources.value.min_replica_count : null
      max_replica_count = automatic_resources.value.max_replica_count > 0 ? automatic_resources.value.max_replica_count : null
    }
  }

  # Explicitly pinned serving compute. machine_spec is a required block
  # even when machine_type is left to the API's default; min_replica_count
  # is required by the API (>= 1).
  dynamic "dedicated_resources" {
    for_each = var.spec.dedicated_resources != null ? [var.spec.dedicated_resources] : []
    content {
      machine_spec {
        machine_type = dedicated_resources.value.machine_type != "" ? dedicated_resources.value.machine_type : null
      }
      min_replica_count = dedicated_resources.value.min_replica_count
      max_replica_count = dedicated_resources.value.max_replica_count > 0 ? dedicated_resources.value.max_replica_count : null
    }
  }

  # IP-space partitioning: empty lets GCP default ("default"). The API
  # HOLDS the group↔ranges pairing — a non-default group, once used with
  # a set of reserved ranges, can only ever be used with exactly that set
  # (taught in the spec).
  deployment_group = var.spec.deployment_group != "" ? var.spec.deployment_group : null

  # False is the provider default, so only true is sent.
  enable_access_logging = var.spec.enable_access_logging ? true : null

  # Names of reserved VPC_PEERING address ranges under the endpoint's
  # peered network; only meaningful on a peered endpoint.
  reserved_ip_ranges = length(var.spec.reserved_ip_ranges) > 0 ? var.spec.reserved_ip_ranges : null

  # JWT auth on the private query endpoint. The provider nests the
  # API's single-child deployedIndexAuthConfig.authProvider wrapper;
  # the spec flattens it to one honest auth_config block.
  dynamic "deployed_index_auth_config" {
    for_each = var.spec.auth_config != null ? [var.spec.auth_config] : []
    content {
      auth_provider {
        allowed_issuers = length(deployed_index_auth_config.value.allowed_issuers) > 0 ? deployed_index_auth_config.value.allowed_issuers : null
        audiences       = length(deployed_index_auth_config.value.audiences) > 0 ? deployed_index_auth_config.value.audiences : null
      }
    }
  }
}
