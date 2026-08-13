variable "metadata" {
  description = "Planton resource metadata"
  type = object({
    name    = string
    id      = optional(string, "")
    org     = optional(string, "")
    env     = optional(string, "")
    labels  = optional(map(string), {})
    tags    = optional(list(string), [])
    version = optional(string, "")
  })
}

variable "spec" {
  description = "GcpVertexAiIndex spec"
  type = object({
    # StringValueOrRef fields arrive from the proto→tfvars converter as
    # plain strings (already resolved), never as object({value}). If
    # project_id is empty, the provider's default project is used
    # (see locals.tf).
    project_id = optional(string, "")

    # Region. Immutable (ForceNew).
    location = string

    display_name = string

    description = optional(string, "")

    # BATCH_UPDATE or STREAM_UPDATE. Empty lets GCP default
    # (BATCH_UPDATE). Immutable (ForceNew).
    index_update_method = optional(string, "")

    # gs:// DIRECTORY holding the initial/delta vector data. Write-only
    # on the provider (GCP never reports it back); changes travel in
    # their own single-field update.
    contents_delta_uri = optional(string, "")

    # Only meaningful together with contents_delta_uri: true replaces
    # the whole index contents, false applies the files as a delta.
    is_complete_overwrite = optional(bool, false)

    # Vector-search geometry. The whole block is immutable (ForceNew).
    config = object({
      dimensions                  = number
      approximate_neighbors_count = optional(number, 0)
      shard_size                  = optional(string, "")
      distance_measure_type       = optional(string, "")
      feature_norm_type           = optional(string, "")

      # At most one algorithm arm (CEL-enforced pre-deploy). Empty
      # object {} is the brute-force presence marker.
      tree_ah_config = optional(object({
        leaf_node_embedding_count    = optional(number, 0)
        leaf_nodes_to_search_percent = optional(number, 0)
      }), null)
      brute_force_config = optional(object({}), null)
    })

    # User labels; merged beneath the platform attribution labels.
    labels = optional(map(string), {})

    # CMEK key resource path (a GcpKmsKey reference resolves to it).
    # Empty means Google-managed encryption. Immutable (ForceNew).
    kms_key_name = optional(string, "")

    # Client-side destroy behavior: DELETE (default), PREVENT, ABANDON.
    deletion_policy = optional(string, "")
  })

  validation {
    condition     = var.spec.location != ""
    error_message = "location is required."
  }

  validation {
    condition     = var.spec.display_name != ""
    error_message = "display_name is required."
  }

  validation {
    condition     = var.spec.config.dimensions >= 1
    error_message = "config.dimensions must be at least 1."
  }
}
