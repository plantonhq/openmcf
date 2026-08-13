# Enable the Pub/Sub API — the control plane that owns the schema.
# disable_on_destroy is false: tearing down one schema must never disable
# the API for everything else in the project.
resource "google_project_service" "pubsub_api" {
  project = local.project_id
  service = "pubsub.googleapis.com"

  disable_dependent_services = true
  disable_on_destroy         = false
}

# The Pub/Sub schema. A shareable resource: one schema can validate
# messages on many topics (each attaches it by reference), so the event
# contract is evolved in one place. Definition changes commit a new
# revision in place (up to 20 revisions per schema) rather than replacing
# the resource; only renaming replaces it.
resource "google_pubsub_schema" "schema" {
  name    = local.schema_name
  project = local.project_id

  # Type and definition travel together: the definition text is parsed as
  # the declared language (AVRO JSON or a protobuf message), and later
  # revisions must keep the same type.
  type       = var.spec.type
  definition = var.spec.definition

  # Client-side destroy behavior: DELETE (default), PREVENT, or ABANDON.
  # Sent only when set so the provider default stays in charge otherwise.
  deletion_policy = var.spec.deletion_policy != "" ? var.spec.deletion_policy : null

  depends_on = [
    google_project_service.pubsub_api,
  ]
}
