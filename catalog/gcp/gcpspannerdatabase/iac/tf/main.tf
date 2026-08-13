# Enable the Spanner API so a fresh project can host databases.
resource "google_project_service" "spanner_api" {
  project = local.project_id
  service = "spanner.googleapis.com"

  disable_dependent_services = true
  disable_on_destroy         = false
}

# The Spanner database — schemas, data, encryption, and retention on a
# Spanner instance. name, dialect, encryption, and instance are immutable;
# retention, drop protection, and appended DDL update in place.
resource "google_spanner_database" "this" {
  instance = var.spec.instance
  name     = local.database_name
  project  = local.project_id

  database_dialect         = local.database_dialect
  version_retention_period = local.version_retention_period
  default_time_zone        = local.default_time_zone

  # GCP API-side lock: while true, no interface (console, gcloud, IaC) can
  # delete the database, and the PARENT INSTANCE cannot be deleted either.
  enable_drop_protection = var.spec.enable_drop_protection

  # DDL is append-only after creation: new statements apply via UpdateDDL;
  # editing or removing an existing entry forces database recreation. The
  # provider quotes identifiers per dialect (backticks vs double quotes).
  ddl = var.spec.ddl

  # IaC-side deletion guard (spec default TRUE): a destroy plan fails while
  # set, before touching GCP. Identical semantics on both engines.
  deletion_protection = var.spec.deletion_protection

  # What a PERMITTED destroy does once the guards above allow one:
  # DELETE (default), PREVENT (destroy fails), or ABANDON (drop from
  # state, keep the database and its data in GCP).
  deletion_policy = var.spec.deletion_policy != "" ? var.spec.deletion_policy : null

  # CMEK: immutable. kms_key_name for regional instance configs,
  # kms_key_names (one key per region) for multi-region configs — the spec
  # enforces exactly one shape pre-deploy.
  dynamic "encryption_config" {
    for_each = var.spec.encryption_config != null ? [var.spec.encryption_config] : []
    content {
      kms_key_name  = encryption_config.value.kms_key_name != "" ? encryption_config.value.kms_key_name : null
      kms_key_names = length(encryption_config.value.kms_key_names) > 0 ? encryption_config.value.kms_key_names : null
    }
  }

  depends_on = [google_project_service.spanner_api]
}
