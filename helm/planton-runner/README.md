# Planton Runner Helm Chart

Deploy [Planton Runner](https://github.com/plantonhq/planton) to any Kubernetes cluster.

The Planton Runner is an agent that connects to the Planton control plane via a secure
mTLS reverse tunnel to execute cloud operations and IaC workflows on behalf of your
organization. It authenticates to the control plane using a ServiceAccount API key
embedded in the credentials file.

## Prerequisites

- Kubernetes 1.24+
- Helm 3.x
- A registered runner with generated credentials (via `planton runner generate-credentials`)

## Quick Start

The chart is published to GitHub Container Registry as an OCI artifact:

```bash
# Install the runner with your credentials file
helm install my-runner oci://ghcr.io/plantonhq/charts/planton-runner \
  --namespace planton-runner \
  --create-namespace \
  --set-file credentials.content=~/.planton/org/<org>/runner/<slug>/credentials.json
```

## Installation

### Via `planton runner deploy` (recommended)

The Planton CLI automates Helm repository setup, values generation, and installation:

```bash
planton runner deploy <runner-name>
```

### Manual Helm install

```bash
helm install my-runner oci://ghcr.io/plantonhq/charts/planton-runner \
  --namespace planton-runner \
  --create-namespace \
  --set-file credentials.content=/path/to/credentials.json
```

### Using an existing Kubernetes Secret

If you manage the credentials Secret externally (e.g., via sealed-secrets, external-secrets,
or a GitOps pipeline), reference it instead of providing the content directly:

```bash
# Create the Secret yourself
kubectl create secret generic runner-creds \
  --namespace planton-runner \
  --from-file=credentials.json=/path/to/credentials.json

# Install the chart referencing the existing Secret
helm install my-runner oci://ghcr.io/plantonhq/charts/planton-runner \
  --namespace planton-runner \
  --set credentials.existingSecret=runner-creds
```

## Configuration

### Credentials

| Parameter | Description | Default |
|-----------|-------------|---------|
| `credentials.content` | Raw JSON content of the RunnerCredentials file | `""` |
| `credentials.existingSecret` | Name of a pre-existing Secret containing the credentials | `""` |
| `credentials.existingSecretKey` | Key within the existing Secret | `"credentials.json"` |

Either `credentials.content` or `credentials.existingSecret` must be provided. If neither
is set, the chart will fail at render time with an error.

### Runner

| Parameter | Description | Default |
|-----------|-------------|---------|
| `runner.port` | gRPC server port (CloudOps routes via tunnel to this port) | `50051` |
| `runner.logLevel` | Log level: debug, info, warn, error | `"info"` |
| `runner.executionMode` | Execution mode: grpc, temporal, dual | `"grpc"` |

### Temporal

These settings are only relevant when `runner.executionMode` is `temporal` or `dual`.
The environment variables are only injected into the pod when `temporal.address` is set.

| Parameter | Description | Default |
|-----------|-------------|---------|
| `temporal.address` | Temporal server address (host:port) | `""` |
| `temporal.namespace` | Temporal namespace | `"default"` |
| `temporal.maxConcurrency` | Maximum concurrent Temporal activities | `10` |

### Build capability

Enable this to make the cluster a **build cluster**: the runner executes build
pipelines on it (Tekton PipelineRuns), streams task pod logs to the control
plane, receives Tekton CloudEvents on its webhook, and serves the readiness
checks a registered build connection reports.

Prerequisites:

- **Tekton Pipelines** installed on the cluster.
- `runner.executionMode` set to `temporal` or `dual` (the build worker is a
  Temporal worker; the chart fails at render time if builds are enabled in
  `grpc` mode).
- Exactly **one** build-capable runner per watched Tekton namespace (the build
  log streamer is a singleton; the chart already pins one replica).

```bash
helm install my-runner oci://ghcr.io/plantonhq/charts/planton-runner \
  --namespace planton-runner \
  --create-namespace \
  --set-file credentials.content=/path/to/credentials.json \
  --set runner.executionMode=dual \
  --set build.enabled=true \
  --set build.tektonNamespace=build-pipelines
```

| Parameter | Description | Default |
|-----------|-------------|---------|
| `build.enabled` | Run the pipeline-build worker and the Tekton CloudEvents webhook | `false` |
| `build.tektonNamespace` | Namespace where builds land and the log streamer watches; empty uses the runner's own namespace | `""` |
| `build.webhookPort` | Container port for the Tekton CloudEvents webhook | `8086` |
| `build.rbac.create` | Create the Role/RoleBinding the build capability needs | `true` |
| `build.tekton.installNamespace` | Namespace of the Tekton Pipelines installation | `"tekton-pipelines"` |

Two steps complete the build-cluster setup after install:

1. **Point Tekton's CloudEvents sink at the runner's webhook** so live build
   status flows without waiting on the reconciliation safety net. In the
   `tekton-pipelines/config-defaults` ConfigMap set:

   ```yaml
   default-cloud-events-sink: http://<release-fullname>.<namespace>.svc.cluster.local/service-hub/tekton/cloud-event
   ```

   (The chart's Service serves the webhook on port 80, so the URL needs no
   explicit port -- but the `/service-hub/tekton/cloud-event` path is
   required: the webhook serves only that route, and a sink pointed at the
   bare Service root gets a 404 on every event, silently degrading builds to
   the reconciliation safety net. The exact URL is printed in the
   post-install notes.)

2. **Register the cluster as a build connection** (console: Connections →
   Build, or `planton apply` a `TektonConnection` naming this runner), then
   verify it. The readiness check runs on this runner over the same queue real
   builds ride, so a passing verify also proves end-to-end routing.

The RBAC the chart creates mirrors exactly what the build capability performs:
create PipelineRuns and their per-build supporting resources, reconcile by
list, label-scoped cleanup, and follow task pod logs — plus a read grant on
Tekton's `config-defaults` ConfigMap so the events-sink readiness check can
report a conclusive verdict.

### Image

| Parameter | Description | Default |
|-----------|-------------|---------|
| `image.repository` | Container image repository | `ghcr.io/plantonhq/planton/runner` |
| `image.tag` | Container image tag | Chart `appVersion` |
| `image.pullPolicy` | Image pull policy | `IfNotPresent` |

### Resources and Scheduling

| Parameter | Description | Default |
|-----------|-------------|---------|
| `resources.requests.cpu` | CPU request | `100m` |
| `resources.requests.memory` | Memory request | `256Mi` |
| `resources.limits.cpu` | CPU limit | `1` |
| `resources.limits.memory` | Memory limit | `1Gi` |
| `serviceAccount.create` | Create a Kubernetes ServiceAccount | `true` |
| `serviceAccount.annotations` | ServiceAccount annotations (e.g., IRSA) | `{}` |
| `nodeSelector` | Node selector for pod scheduling | `{}` |
| `tolerations` | Tolerations for pod scheduling | `[]` |
| `affinity` | Affinity rules for pod scheduling | `{}` |

## Architecture

The runner makes only **outbound** connections. No Ingress or public Service is required.

```
Runner Pod (your cluster)
  ├─ mTLS tunnel ──► runner-tunnel.planton.live:443  (tunnel server)
  └─ gRPC + TLS  ──► api.planton.live:443            (control plane API)
```

The runner maintains a persistent reverse tunnel through which the Planton control plane
sends cloud operations requests. It also makes authenticated gRPC calls to the control
plane API to resolve secrets and variables at runtime.

With builds enabled there is one additional, **cluster-internal** listener: Tekton posts
pipeline CloudEvents to the runner's webhook through the chart's ClusterIP Service. No
Ingress and no public exposure — the traffic never leaves the cluster.

### Credentials File

The credentials file is a JSON document generated by the control plane's
`generateCredentials` RPC. It contains the runner's identity, mTLS certificates,
API key, and endpoint configuration in a single file. The chart stores this as a
Kubernetes Secret and mounts it into the pod at `/etc/planton/credentials/credentials.json`.

### Automatic Rollouts

When the chart manages the credentials Secret (i.e., `credentials.existingSecret` is not
set), the Deployment includes a `checksum/credentials` annotation that forces a pod rollout
whenever the credentials content changes. When using an existing Secret, you are responsible
for triggering rollouts (e.g., via `kubectl rollout restart`).

## Verification

After installation, verify the runner is running:

```bash
kubectl -n planton-runner get pods
kubectl -n planton-runner logs -l app.kubernetes.io/name=planton-runner
```

## Upgrading

To update the credentials or configuration:

```bash
helm upgrade my-runner oci://ghcr.io/plantonhq/charts/planton-runner \
  --namespace planton-runner \
  --set-file credentials.content=/path/to/new-credentials.json
```

## Uninstalling

```bash
helm uninstall my-runner --namespace planton-runner
```

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
