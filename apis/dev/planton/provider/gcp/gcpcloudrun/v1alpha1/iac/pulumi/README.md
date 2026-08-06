# GcpCloudRun - Pulumi Module

This Pulumi (Go) module provisions a Cloud Run service (`cloudrunv2.Service`). It is the Pulumi-side implementation of the Planton `GcpCloudRun` resource kind and has feature parity with the Terraform module.

## Overview

The module enables the Cloud Run Admin API (`run.googleapis.com`) so a fresh project works first try, then creates the service with the full released-provider surface: multiple containers (the serving container plus sidecars with `dependsOn` startup ordering), literal and Secret Manager environment variables, startup and liveness probes (HTTP/TCP/gRPC), all five volume types (Cloud SQL sockets, Secret Manager files, in-memory or disk scratch, GCS FUSE, NFS), per-revision and service-level scaling (including MANUAL mode), revision traffic splitting with tags, direct VPC egress or a Serverless VPC Access connector, GPU node selection, CMEK image encryption, Binary Authorization, session affinity, and custom audiences.

Public access is the additive-IAM path: when `allow_unauthenticated` is true a `cloudrunv2.ServiceIamMember` grants `roles/run.invoker` to `allUsers`; `invoker_iam_disabled` is the org-policy alternative that switches the IAM check off instead (the spec rejects setting both). `deletion_protection` defaults to true — a destroy fails until the manifest opts out.

An empty `project_id` falls back to the provider's default project — the ambient-project contract every GCP kind honors. Enum values pass through as the API's own names, untranslated.

## Usage with Planton CLI

```shell
planton pulumi up --manifest ../hack/manifest.yaml --module-dir .
planton pulumi destroy --manifest ../hack/manifest.yaml --module-dir .
```

Credentials are provided via stack input (by the CLI), not in the manifest `spec`. Manifest file: `../hack/manifest.yaml`.

## Build

```bash
cd apis/dev/planton/provider/gcp/gcpcloudrun/v1alpha1/iac/pulumi
make build
```

## Outputs

| Name | Description |
|------|-------------|
| `url` | Canonical serving URL of the service |
| `service_name` | Name of the service as created in GCP (the serverless-NEG handle) |
| `revision` | Latest ready revision name |
| `location` | Region the service is deployed in |
| `uid` | Server-assigned unique identifier |
| `urls` | Every URL serving this service |
