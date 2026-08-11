# A logical database inside a Cloud SQL instance. Databases carry their own
# lifecycle: create and drop application databases freely without touching
# the instance node they live on.
#
# No API enablement here: the instance this database lives on cannot exist
# without sqladmin.googleapis.com already enabled (its own module enables
# it), so a database module enabling the API again would only add churn.
#
# Charset/collation semantics are engine-specific — MySQL accepts any
# supported pair, PostgreSQL requires UTF8 at creation with an OS-locale
# collation, SQL Server ignores charset entirely. The API validates the
# combination at deploy time.
resource "google_sql_database" "this" {
  name     = var.spec.database_name
  project  = local.project_id
  instance = var.spec.instance

  charset   = local.charset
  collation = local.collation

  # DELETE (default) drops the database; ABANDON removes it from IaC
  # management — the documented workaround when live connections block a
  # PostgreSQL drop; PREVENT fails destroying plans.
  deletion_policy = var.spec.deletion_policy != "" ? var.spec.deletion_policy : null
}
