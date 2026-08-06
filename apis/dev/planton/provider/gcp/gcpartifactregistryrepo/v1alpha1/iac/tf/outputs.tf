# Short name of the repository (the repository ID).
output "name" {
  description = "Short name of the repository"
  value       = google_artifact_registry_repository.this.name
}

# The resource ID is the fully qualified repository path
# (projects/{project}/locations/{location}/repositories/{repository_id}) —
# the exact string every composing resource consumes: a Cloud Function's
# docker_repository, a virtual repository's upstream policy, and a remote
# repository's common upstream.
output "repository_path" {
  description = "Fully qualified repository resource path"
  value       = google_artifact_registry_repository.this.id
}

# The registry endpoint clients push to and pull from. Constructed from
# resolved attributes ({location}-{format}.pkg.dev/{project}/{repo}) because
# the released 6.x provider does not export a registry URI attribute — the
# Pulumi module builds the identical string.
output "registry_uri" {
  description = "Registry endpoint (e.g. us-central1-docker.pkg.dev/my-project/my-repo)"
  value = format(
    "%s-%s.pkg.dev/%s/%s",
    lower(google_artifact_registry_repository.this.location),
    lower(google_artifact_registry_repository.this.format),
    google_artifact_registry_repository.this.project,
    google_artifact_registry_repository.this.name,
  )
}

output "location" {
  description = "Location of the repository (region or multi-region)"
  value       = google_artifact_registry_repository.this.location
}
