# KubernetesTektonOperator - Usage Examples

This document provides practical examples for deploying the Tekton Operator using the **KubernetesTektonOperator** API resource.

> **Note:** After creating a YAML file with your configuration, apply it using:
> ```shell
> planton apply -f <yaml-file>
> ```

---

## ⚠️ Important: Fixed Namespace Architecture

**Unlike other Kubernetes components in OpenMCF, the Tekton Operator uses fixed namespaces** that are managed by the operator itself:

| Component | Namespace | Managed By |
|-----------|-----------|------------|
| Tekton Operator | `tekton-operator` | Operator Release |
| Tekton Pipelines | `tekton-pipelines` | TektonConfig |
| Tekton Triggers | `tekton-pipelines` | TektonConfig |
| Tekton Dashboard | `tekton-pipelines` | TektonConfig |

**These namespaces are automatically created and managed by the Tekton Operator and cannot be customized.**

See: https://tekton.dev/docs/operator/tektonconfig/

---

## 1. Full Installation (All Components)

Deploy Tekton Operator with all components enabled for a complete CI/CD platform.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesTektonOperator
metadata:
  name: tekton-operator
spec:
  target_cluster:
    cluster_name: "my-k8s-cluster"
  container:
    resources:
      requests:
        cpu: "100m"
        memory: "128Mi"
      limits:
        cpu: "500m"
        memory: "512Mi"
  components:
    pipelines: true
    triggers: true
    dashboard: true
  operator_version: "v0.78.0"  # Specify the operator version (default: v0.78.0)
```

**What this does:**
- Installs Tekton Operator in `tekton-operator` namespace
- Deploys Tekton components in `tekton-pipelines` namespace
- Enables Tekton Pipelines for CI/CD pipeline execution
- Enables Tekton Triggers for event-driven pipeline execution
- Enables Tekton Dashboard for web-based UI
- Allocates 100m CPU and 128Mi memory as baseline

**When to use:** Production environments requiring a complete CI/CD platform with webhook integration and visual management.

---

## 2. Specific Operator Version

Deploy a specific version of the Tekton Operator for compatibility or testing.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesTektonOperator
metadata:
  name: tekton-operator-pinned
spec:
  target_cluster:
    cluster_name: "my-k8s-cluster"
  container:
    resources:
      requests:
        cpu: "100m"
        memory: "128Mi"
      limits:
        cpu: "500m"
        memory: "512Mi"
  components:
    pipelines: true
    triggers: true
    dashboard: false
  operator_version: "v0.75.0"  # Pin to specific version
```

**What this does:**
- Installs Tekton Operator version v0.75.0
- Useful for testing new versions before production rollout
- Ensures consistency across environments

**When to use:**
- Testing compatibility with specific Tekton versions
- Environments requiring version pinning for compliance
- Gradual rollout of new operator versions

**Available versions:** https://github.com/tektoncd/operator/releases

---

## 3. Minimal Installation (Pipelines Only, Default Version)

Deploy only Tekton Pipelines for basic CI/CD functionality.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesTektonOperator
metadata:
  name: tekton-operator-minimal
spec:
  target_cluster:
    cluster_name: "my-k8s-cluster"
  container:
    resources:
      requests:
        cpu: "50m"
        memory: "64Mi"
      limits:
        cpu: "250m"
        memory: "256Mi"
  components:
    pipelines: true
    triggers: false
    dashboard: false
```

**What this does:**
- Installs only Tekton Pipelines
- Minimal resource allocation
- No event-driven capabilities (manual pipeline runs only)
- No web UI

**When to use:**
- Development/testing environments
- When pipelines are triggered programmatically (e.g., from other CI systems)
- Resource-constrained environments

---

## 4. Pipelines with Triggers (No Dashboard)

Enable event-driven pipeline execution without the dashboard.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesTektonOperator
metadata:
  name: tekton-operator-headless
spec:
  target_cluster:
    cluster_name: "production-cluster"
  container:
    resources:
      requests:
        cpu: "100m"
        memory: "128Mi"
      limits:
        cpu: "500m"
        memory: "512Mi"
  components:
    pipelines: true
    triggers: true
    dashboard: false
```

**What this does:**
- Enables pipeline execution
- Enables event listeners for webhooks
- No web UI (manage via kubectl or GitOps)

**When to use:**
- GitOps workflows where UI is not needed
- Environments with strict security requirements (reduced attack surface)
- When using external monitoring tools (e.g., Grafana)

---

## 5. Production Configuration

High-resource configuration for production workloads.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesTektonOperator
metadata:
  name: tekton-operator-prod
  labels:
    environment: production
    team: platform
spec:
  target_cluster:
    cluster_name: "production-cluster"
  container:
    resources:
      requests:
        cpu: "200m"
        memory: "256Mi"
      limits:
        cpu: "1000m"
        memory: "1Gi"
  components:
    pipelines: true
    triggers: true
    dashboard: true
```

**What this does:**
- Higher resource allocation for production stability
- All components enabled
- Labels for organizational tracking

**When to use:** Production environments with high pipeline throughput.

---

## 6. Development Environment

Minimal resources for local development with kind or minikube.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesTektonOperator
metadata:
  name: tekton-operator-dev
spec:
  target_cluster:
    cluster_name: "dev-cluster"
  container:
    resources:
      requests:
        cpu: "25m"
        memory: "32Mi"
      limits:
        cpu: "200m"
        memory: "256Mi"
  components:
    pipelines: true
    triggers: false
    dashboard: true
```

**What this does:**
- Minimal resource footprint
- Pipelines and Dashboard for development
- No Triggers (manual runs only)

**When to use:**
- Local development with minikube or kind
- CI/CD for development branches
- Learning and experimentation

---

## 7. Multi-Cluster Deployment Pattern

> **Note:** Each cluster can only have one Tekton Operator installation since it uses fixed namespaces. For multi-environment setups, deploy to separate clusters.

### Production Cluster

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesTektonOperator
metadata:
  name: tekton-operator-prod
  labels:
    environment: production
spec:
  target_cluster:
    cluster_name: "prod-cluster"
  container:
    resources:
      requests:
        cpu: "200m"
        memory: "256Mi"
      limits:
        cpu: "1000m"
        memory: "1Gi"
  components:
    pipelines: true
    triggers: true
    dashboard: true
```

### Staging Cluster

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesTektonOperator
metadata:
  name: tekton-operator-staging
  labels:
    environment: staging
spec:
  target_cluster:
    cluster_name: "staging-cluster"
  container:
    resources:
      requests:
        cpu: "100m"
        memory: "128Mi"
      limits:
        cpu: "500m"
        memory: "512Mi"
  components:
    pipelines: true
    triggers: true
    dashboard: true
```

### Development Cluster

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesTektonOperator
metadata:
  name: tekton-operator-dev
  labels:
    environment: development
spec:
  target_cluster:
    cluster_name: "dev-cluster"
  container:
    resources:
      requests:
        cpu: "50m"
        memory: "64Mi"
      limits:
        cpu: "250m"
        memory: "256Mi"
  components:
    pipelines: true
    triggers: false
    dashboard: true
```

---

## Post-Installation: Creating Tekton Resources

After the Tekton Operator is installed, you can create CI/CD resources.

### Example: Simple Task

```yaml
apiVersion: tekton.dev/v1beta1
kind: Task
metadata:
  name: hello-world
spec:
  steps:
    - name: say-hello
      image: alpine
      command:
        - echo
      args:
        - "Hello, World!"
```

### Example: Build Pipeline

```yaml
apiVersion: tekton.dev/v1beta1
kind: Pipeline
metadata:
  name: build-pipeline
spec:
  params:
    - name: repo-url
      type: string
  tasks:
    - name: clone
      taskRef:
        name: git-clone
      params:
        - name: url
          value: $(params.repo-url)
    - name: build
      taskRef:
        name: kaniko
      runAfter:
        - clone
```

### Example: GitHub Event Listener

```yaml
apiVersion: triggers.tekton.dev/v1beta1
kind: EventListener
metadata:
  name: github-listener
spec:
  serviceAccountName: tekton-triggers-sa
  triggers:
    - name: github-push
      interceptors:
        - ref:
            name: github
          params:
            - name: secretRef
              value:
                secretName: github-secret
                secretKey: token
            - name: eventTypes
              value:
                - push
      bindings:
        - ref: github-push-binding
      template:
        ref: pipeline-template
```

---

## Resource Sizing Guidelines

### Operator Resource Needs

| Scenario | Pipelines/Day | Recommended CPU | Recommended Memory |
|----------|---------------|-----------------|-------------------|
| Small (Dev/Test) | 1-10 | 25-50m | 32-64Mi |
| Medium | 10-50 | 100-200m | 128-256Mi |
| Large (Production) | 50-200 | 200-500m | 256-512Mi |
| Enterprise | 200+ | 500m-1000m | 512Mi-1Gi |

### Factors Affecting Resource Usage

- **Pipeline Concurrency**: More concurrent pipelines = higher resource needs
- **Trigger Volume**: High webhook traffic increases Triggers controller load
- **Dashboard Usage**: Active Dashboard users increase memory usage

---

## Verification

After deployment, verify the Tekton Operator is running:

```bash
# Check operator pod
kubectl get pods -n tekton-operator

# View operator logs
kubectl logs -n tekton-operator -l app=tekton-operator

# Verify TektonConfig
kubectl get tektonconfig

# Check Tekton Pipelines
kubectl get pods -n tekton-pipelines

# Verify CRDs are installed
kubectl get crds | grep tekton
```

Expected CRDs installed:
- `pipelines.tekton.dev`
- `tasks.tekton.dev`
- `pipelineruns.tekton.dev`
- `taskruns.tekton.dev`
- `eventlisteners.triggers.tekton.dev` (if Triggers enabled)
- `triggerbindings.triggers.tekton.dev` (if Triggers enabled)

---

## Troubleshooting

### Operator Pod Not Starting

```bash
# Check pod events
kubectl describe pod -n tekton-operator -l app=tekton-operator

# Check resource availability
kubectl top nodes
```

### Pipelines Not Running

```bash
# Check pipeline controller
kubectl logs -n tekton-pipelines -l app=tekton-pipelines-controller

# Check PipelineRun status
kubectl describe pipelinerun <name>
```

### Dashboard Not Accessible

```bash
# Port-forward to dashboard
kubectl port-forward -n tekton-pipelines svc/tekton-dashboard 9097:9097

# Check dashboard pod
kubectl logs -n tekton-pipelines -l app=tekton-dashboard
```

---

## Next Steps

1. **Create Tasks**: Define reusable task templates
2. **Build Pipelines**: Compose tasks into pipelines
3. **Set Up Triggers**: Connect webhooks to pipelines
4. **Configure RBAC**: Set up proper access control
5. **Implement Monitoring**: Add pipeline observability

For more information:
- [Tekton Documentation](https://tekton.dev/docs/)
- [Tekton Pipelines](https://github.com/tektoncd/pipeline)
- [Tekton Triggers](https://github.com/tektoncd/triggers)
- [Tekton Dashboard](https://github.com/tektoncd/dashboard)
