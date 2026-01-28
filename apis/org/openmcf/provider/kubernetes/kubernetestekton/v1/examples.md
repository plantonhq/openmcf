# KubernetesTekton Examples

## Example 1: Minimal Pipelines Only

Deploy just Tekton Pipelines without the dashboard:

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesTekton
metadata:
  name: tekton-minimal
spec:
  pipeline_version: "latest"
```

## Example 2: Pipelines with Dashboard

Deploy Tekton Pipelines and Dashboard:

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesTekton
metadata:
  name: tekton-with-dashboard
spec:
  pipeline_version: "v0.65.2"
  dashboard:
    enabled: true
    version: "v0.53.0"
```

## Example 3: With CloudEvents Sink

Configure Tekton to send CloudEvents to an external receiver:

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesTekton
metadata:
  name: tekton-cloudevents
spec:
  pipeline_version: "latest"
  dashboard:
    enabled: true
  cloud_events:
    sink_url: "http://event-receiver.platform.svc.cluster.local/tekton/events"
```

**CloudEvents Format:**

Tekton sends CloudEvents with:
- `type: dev.tekton.event.taskrun.successful` (and similar for other states)
- `source: /tekton/pipeline/taskruns/{name}`
- Data payload with TaskRun/PipelineRun details

## Example 4: Full Production Setup

Complete setup with dashboard, CloudEvents, and external ingress:

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesTekton
metadata:
  name: tekton-prod
  org: my-org
  env: prod
spec:
  pipeline_version: "v0.65.2"
  dashboard:
    enabled: true
    version: "v0.53.0"
    ingress:
      enabled: true
      hostname: "tekton-dashboard.planton.live"
  cloud_events:
    sink_url: "http://service-hub.platform.svc.cluster.local/tekton/cloud-event"
```

**Prerequisites for Ingress:**
1. Istio ingress gateway in `istio-ingress` namespace
2. cert-manager with ClusterIssuer named `planton.live`
3. DNS pointing `tekton-dashboard.planton.live` to ingress LB

## Example 5: Development Environment

Lighter configuration for development:

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesTekton
metadata:
  name: tekton-dev
  env: dev
spec:
  pipeline_version: "latest"
  dashboard:
    enabled: true
    version: "latest"
```

Access via port-forward:
```bash
kubectl port-forward -n tekton-pipelines service/tekton-dashboard 9097:9097
open http://localhost:9097
```

## Verification Commands

After deployment, verify Tekton is running:

```bash
# Check pipeline controller
kubectl get pods -n tekton-pipelines -l app=tekton-pipelines-controller

# Check dashboard (if enabled)
kubectl get pods -n tekton-pipelines -l app=tekton-dashboard

# Verify CRDs are installed
kubectl get crds | grep tekton

# Check config-defaults ConfigMap (for cloud events)
kubectl get configmap config-defaults -n tekton-pipelines -o yaml

# Run a test TaskRun
kubectl apply -f - <<EOF
apiVersion: tekton.dev/v1
kind: TaskRun
metadata:
  generateName: hello-world-
  namespace: default
spec:
  taskSpec:
    steps:
      - name: echo
        image: alpine
        command: ["echo"]
        args: ["Hello from Tekton!"]
EOF

# Check TaskRun status
kubectl get taskrun -n default
```
