# main.tf

# Create the Cloudflare D1 database (a serverless SQLite database a Worker
# reaches via a d1 binding). Placement is fixed at creation by an optional
# region hint.
resource "cloudflare_d1_database" "main" {
  account_id = var.spec.account_id
  name       = var.spec.database_name

  # Region hint (omitted when unspecified so Cloudflare selects a default).
  primary_location_hint = local.primary_location_hint

  # Data-residency jurisdiction (omitted when unset; mutually exclusive with region).
  jurisdiction = local.jurisdiction

  # Cloudflare always reports read replication on read (mode "disabled" when
  # never configured) while the provider models the attribute as
  # Optional-not-Computed, so an omitted spec block must coalesce to the
  # server default -- sending null instead leaves a refresh-vs-config diff
  # that never converges (measured live 2026-08-26).
  read_replication = {
    mode = var.spec.read_replication != null ? var.spec.read_replication.mode : "disabled"
  }

  lifecycle {
    # primary_location_hint is a create-time placement hint Cloudflare never
    # returns (no_refresh + RequiresReplace upstream at v5.23.0). Without the
    # ignore, adopting an existing database restores the hint as null and a
    # post-create region edit plans a REPLACE -- destroying the database to
    # move a hint. Placement is decided at birth; post-create edits are
    # deliberately inert.
    ignore_changes = [primary_location_hint]
  }
}
