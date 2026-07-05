# Terraform module for GcpCloudFunction

Provisions a `google_cloudfunctions2_function` (plus the public-invoker `google_cloud_run_service_iam_member` when `allowUnauthenticated` is set) from the validated protobuf spec. Enables the Cloud Functions, Cloud Build, Cloud Run, Artifact Registry, and Eventarc APIs automatically.

See the component [README](../../README.md) for the full configuration reference.
