# GcpCloudTasksQueue Pulumi Module

Pulumi Go implementation for provisioning a GCP Cloud Tasks queue.

## Architecture

```
module/
  main.go                 # Entry point (Resources function)
  locals.go               # Variable transformations
  cloud_tasks_queue.go    # Queue resource creation
  outputs.go              # Output constant definitions
```

## SDK

Uses `pulumi-gcp/sdk/v9/go/gcp/cloudtasks` for the `cloudtasks.NewQueue` resource.

## Features

- Cloud Tasks API enablement (fresh projects work first try)
- Queue-level HTTP target with OIDC/OAuth authentication
- URI overrides (scheme, host, port, path, query params)
- Header overrides for all tasks
- App Engine routing override (service/version/instance pinning)
- Configurable rate limits and retry behavior
- Stackdriver logging with sampling ratio

## Notes

- Cloud Tasks queues do NOT support GCP labels. No labels are computed or applied.
- If `project_id` is empty, the queue lands in the provider's default project.
- `max_burst_size` is computed by GCP from the dispatch rate and exported as an output.
- Flattened URI path/query overrides are mapped back to the SDK's nested structure.
- The bridged provider's client-side `deletion_policy` is pinned to `DELETE`, and its
  `desired_state` (pause/resume) surface is deliberately unused — see the PARITY
  comments in `cloud_tasks_queue.go`.

## Debug

```bash
cd ~/scm/github.com/plantonhq/planton/apis/dev/planton/provider/gcp/gcpcloudtasksqueue/v1/iac/pulumi
make build
PULUMI_CONFIG_PASSPHRASE="" pulumi login --local
pulumi stack init dev
pulumi config set --path 'stack_input' --secret < ../../hack/manifest.yaml
make preview
```
