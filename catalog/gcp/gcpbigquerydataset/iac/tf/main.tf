# Enable the BigQuery API so a fresh project can host datasets.
resource "google_project_service" "bigquery_api" {
  project = local.project_id
  service = "bigquery.googleapis.com"

  disable_dependent_services = true
  disable_on_destroy         = false
}

# The BigQuery dataset — the container that pins data location and owns the
# defaults every contained table inherits (expiration, CMEK, collation) plus
# the dataset-level ACL. dataset_id, location, and is_case_insensitive are
# immutable; everything else updates in place.
resource "google_bigquery_dataset" "this" {
  dataset_id = var.spec.dataset_id
  project    = local.project_id
  location   = var.spec.location

  friendly_name                   = local.friendly_name
  description                     = local.description
  default_table_expiration_ms     = local.default_table_expiration_ms
  default_partition_expiration_ms = local.default_partition_expiration_ms
  max_time_travel_hours           = local.max_time_travel_hours
  is_case_insensitive             = var.spec.is_case_insensitive
  default_collation               = local.default_collation
  storage_billing_model           = local.storage_billing_model

  # When false (the safe default), destroy fails while the dataset contains
  # tables — the guard against deleting data with its container.
  delete_contents_on_destroy = var.spec.delete_contents_on_destroy

  labels        = local.final_labels
  resource_tags = local.resource_tags

  # CMEK default for all new tables. The BigQuery service agent must hold
  # cryptoKeyEncrypterDecrypter on the key before the first table write.
  dynamic "default_encryption_configuration" {
    for_each = local.kms_key_name != null ? [local.kms_key_name] : []
    content {
      kms_key_name = default_encryption_configuration.value
    }
  }

  # The spec's access list is AUTHORITATIVE: these entries become the
  # dataset's complete ACL. An entry is either a principal grant (role +
  # one identity) or a resource authorization (view/routine/dataset with
  # implicit read access and no role) — the spec's CEL rules guarantee the
  # shape before this module ever runs.
  dynamic "access" {
    for_each = var.spec.access
    content {
      role           = access.value.role != "" ? access.value.role : null
      user_by_email  = access.value.user_by_email != "" ? access.value.user_by_email : null
      group_by_email = access.value.group_by_email != "" ? access.value.group_by_email : null
      domain         = access.value.domain != "" ? access.value.domain : null
      special_group  = access.value.special_group != "" ? access.value.special_group : null
      iam_member     = access.value.iam_member != "" ? access.value.iam_member : null

      dynamic "view" {
        for_each = access.value.view != null ? [access.value.view] : []
        content {
          project_id = view.value.project_id
          dataset_id = view.value.dataset_id
          table_id   = view.value.table_id
        }
      }

      dynamic "routine" {
        for_each = access.value.routine != null ? [access.value.routine] : []
        content {
          project_id = routine.value.project_id
          dataset_id = routine.value.dataset_id
          routine_id = routine.value.routine_id
        }
      }

      # The provider nests the grantee dataset reference one level deeper
      # (dataset { dataset { ... } target_types }); the spec flattens that
      # single-purpose wrapper.
      dynamic "dataset" {
        for_each = access.value.dataset != null ? [access.value.dataset] : []
        content {
          dataset {
            project_id = dataset.value.project_id
            dataset_id = dataset.value.dataset_id
          }
          target_types = dataset.value.target_types
        }
      }

      dynamic "condition" {
        for_each = access.value.condition != null ? [access.value.condition] : []
        content {
          expression  = condition.value.expression
          title       = condition.value.title != "" ? condition.value.title : null
          description = condition.value.description != "" ? condition.value.description : null
          location    = condition.value.location != "" ? condition.value.location : null
        }
      }
    }
  }

  # Immutable: converts the dataset into a read-only projection of an
  # external source (e.g. AWS Glue) through a BigQuery Omni connection.
  dynamic "external_dataset_reference" {
    for_each = var.spec.external_dataset_reference != null ? [var.spec.external_dataset_reference] : []
    content {
      external_source = external_dataset_reference.value.external_source
      connection      = external_dataset_reference.value.connection
    }
  }

  # Hive Metastore compatibility metadata for open-source engines.
  dynamic "external_catalog_dataset_options" {
    for_each = var.spec.external_catalog_options != null ? [var.spec.external_catalog_options] : []
    content {
      default_storage_location_uri = external_catalog_dataset_options.value.default_storage_location_uri != "" ? external_catalog_dataset_options.value.default_storage_location_uri : null
      parameters                   = length(external_catalog_dataset_options.value.parameters) > 0 ? external_catalog_dataset_options.value.parameters : null
    }
  }

  depends_on = [google_project_service.bigquery_api]
}
