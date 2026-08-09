# Exactly one of the four sink resources exists (see the count guards in
# main.tf), so each output selects whichever branch was created with
# one(concat(...)).

# The sink name as it exists in GCP.
output "sink_name" {
  description = "Name of the logging sink"
  value       = one(concat(google_logging_project_sink.this[*].name, google_logging_folder_sink.this[*].name, google_logging_organization_sink.this[*].name, google_logging_billing_account_sink.this[*].name))
}

# The identity GCP minted (or adopted) for this sink — grant it write
# access on the destination or the sink silently exports nothing.
output "writer_identity" {
  description = "serviceAccount:{email} identity that must be granted write access on the destination"
  value       = one(concat(google_logging_project_sink.this[*].writer_identity, google_logging_folder_sink.this[*].writer_identity, google_logging_organization_sink.this[*].writer_identity, google_logging_billing_account_sink.this[*].writer_identity))
}
