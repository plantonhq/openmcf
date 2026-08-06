# Standard Docker Repository

The workhorse of container delivery: a private Docker repository with
immutable tags, self-cleaning storage, and an additive pull grant for the
runtime service account.

## What this preset creates

A `STANDARD_REPOSITORY` in `us-central1` storing Docker/OCI images at
`us-central1-docker.pkg.dev/{project}/app-images`. Immutable tags make
every deployment reproducible; the paired cleanup policies delete
superseded untagged digests after 30 days while always protecting the 10
most recent versions of every image.

## Prerequisites

- A `GcpServiceAccount` named `app-runtime` (the workload identity that
  pulls images). Replace the reference with your own service account, or
  drop the grant and manage access elsewhere.

## Composing delivery

- CI pushes to the `registry_uri` output
  (`us-central1-docker.pkg.dev/{project}/app-images/api-server:v1.2.3`).
- GKE workloads, Cloud Run services, and Cloud Functions pull the same
  reference; a Cloud Function's `dockerRepository` takes this
  repository's `repository_path` output directly.

## Remix ideas

- Set `cleanupPolicyDryRun: true` first and watch Cloud Audit Logs to
  validate the policies against real traffic before they delete.
- Add `kmsKeyName` (a `GcpKmsKey` reference) for CMEK-protected storage.
- Grant `roles/artifactregistry.writer` to your CI service account as a
  second `iamMembers` entry.
