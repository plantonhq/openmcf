# Pulumi module for GcpCloudFunction

Provisions a `cloudfunctionsv2.Function` (plus the public-invoker `cloudrunv2.ServiceIamMember` when `allowUnauthenticated` is set) from the validated protobuf spec. Enables the Cloud Functions, Cloud Build, Cloud Run, Artifact Registry, and Eventarc APIs automatically.

See the component [README](../../README.md) for the full configuration reference.
