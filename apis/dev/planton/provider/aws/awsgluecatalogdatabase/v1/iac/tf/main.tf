# AWS Glue Data Catalog database.
#
# A single metadata resource with three creation shapes: a regular database,
# a resource link to a database shared from another account/region
# (target_database), or a federated projection of an external source
# (federated_database). The spec's CEL keeps the shapes exclusive, so the
# module simply maps each optional message onto its provider block.
resource "aws_glue_catalog_database" "this" {
  name = local.database_name

  # Defaults to the deploying account's own catalog; set only for the
  # cross-account create-in-another-catalog governance pattern. ForceNew.
  catalog_id = var.spec.catalog_id != "" ? var.spec.catalog_id : null

  description  = var.spec.description != "" ? var.spec.description : null
  location_uri = var.spec.location_uri != "" ? var.spec.location_uri : null

  # Catalog metadata properties read by engines and governance tooling --
  # distinct from AWS resource tags (local.tags below).
  parameters = length(var.spec.parameters) > 0 ? var.spec.parameters : null

  # Default Lake Formation grants applied to tables CREATED in this database.
  # An entry with an empty permissions list (and no principal) is meaningful:
  # it stops granting ALL to IAM_ALLOWED_PRINCIPALS on new tables -- the
  # hardening step when moving to Lake Formation-managed permissions.
  dynamic "create_table_default_permission" {
    for_each = var.spec.create_table_default_permissions
    content {
      permissions = create_table_default_permission.value.permissions

      dynamic "principal" {
        for_each = create_table_default_permission.value.principal != "" ? [1] : []
        content {
          data_lake_principal_identifier = create_table_default_permission.value.principal
        }
      }
    }
  }

  # Resource link: a local pointer to a database shared via AWS RAM /
  # Lake Formation. All coordinates are ForceNew.
  dynamic "target_database" {
    for_each = local.has_target_database ? [1] : []
    content {
      catalog_id    = var.spec.target_database.catalog_id
      database_name = var.spec.target_database.database_name
      region        = var.spec.target_database.region != "" ? var.spec.target_database.region : null
    }
  }

  # Federated database: projects an external source (e.g. a Redshift
  # datashare) into the catalog through a Glue connection.
  dynamic "federated_database" {
    for_each = local.has_federated_database ? [1] : []
    content {
      identifier      = var.spec.federated_database.identifier != "" ? var.spec.federated_database.identifier : null
      connection_name = var.spec.federated_database.connection_name != "" ? var.spec.federated_database.connection_name : null
    }
  }

  tags = local.tags
}
