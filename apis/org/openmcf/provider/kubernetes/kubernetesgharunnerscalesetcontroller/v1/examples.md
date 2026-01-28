# KubernetesGhaRunnerScaleSetController Examples

This document provides practical examples for deploying the GitHub Actions Runner Scale Set Controller using OpenMCF.

## Basic Examples

### Minimal Configuration

Deploy the controller with default settings:

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesGhaRunnerScaleSetController
metadata:
  name: arc-controller
spec:
  namespace:
    value: arc-system
  createNamespace: true
  container: {}
```

**Deploy with Pulumi:**
```bash
openmcf pulumi up --manifest arc-controller.yaml --stack org/project/env
```

**Deploy with Terraform:**
```bash
openmcf tofu apply --manifest arc-controller.yaml
```

### Standard Configuration

Deploy with explicit resource configuration:

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesGhaRunnerScaleSetController
metadata:
  name: arc-controller
spec:
  namespace:
    value: arc-system
  createNamespace: true
  helmChartVersion: "0.13.1"
  container:
    resources:
      requests:
        cpu: 100m
        memory: 128Mi
      limits:
        cpu: 500m
        memory: 512Mi
```

## Production Examples

### High Availability Setup

Deploy with multiple replicas for production reliability:

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesGhaRunnerScaleSetController
metadata:
  name: arc-controller-ha
spec:
  namespace:
    value: arc-system
  createNamespace: true
  replicaCount: 3
  priorityClassName: system-cluster-critical
  container:
    resources:
      requests:
        cpu: 200m
        memory: 256Mi
      limits:
        cpu: 1000m
        memory: 1Gi
  flags:
    logLevel: info
    logFormat: json
```

### With Metrics Enabled

Enable Prometheus metrics for monitoring:

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesGhaRunnerScaleSetController
metadata:
  name: arc-controller-monitored
spec:
  namespace:
    value: arc-system
  createNamespace: true
  container:
    resources:
      requests:
        cpu: 100m
        memory: 128Mi
      limits:
        cpu: 500m
        memory: 512Mi
  metrics:
    controllerManagerAddr: ":8080"
    listenerAddr: ":8080"
    listenerEndpoint: "/metrics"
```

### With Rate Limiting

Configure API client rate limits for large clusters:

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesGhaRunnerScaleSetController
metadata:
  name: arc-controller-tuned
spec:
  namespace:
    value: arc-system
  createNamespace: true
  container:
    resources:
      requests:
        cpu: 200m
        memory: 256Mi
      limits:
        cpu: 1000m
        memory: 1Gi
  flags:
    logLevel: info
    runnerMaxConcurrentReconciles: 10
    k8sClientRateLimiterQps: 50
    k8sClientRateLimiterBurst: 100
```

## Namespace-Scoped Examples

### Watch Single Namespace

Restrict the controller to watch only a specific namespace:

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesGhaRunnerScaleSetController
metadata:
  name: arc-controller-scoped
spec:
  namespace:
    value: arc-system
  createNamespace: true
  container:
    resources:
      requests:
        cpu: 100m
        memory: 128Mi
      limits:
        cpu: 500m
        memory: 512Mi
  flags:
    watchSingleNamespace: arc-runners
```

## Custom Image Examples

### Using a Custom Controller Image

Deploy with a custom or mirrored controller image:

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesGhaRunnerScaleSetController
metadata:
  name: arc-controller-custom
spec:
  namespace:
    value: arc-system
  createNamespace: true
  container:
    resources:
      requests:
        cpu: 100m
        memory: 128Mi
      limits:
        cpu: 500m
        memory: 512Mi
    image:
      repository: my-registry.example.com/actions/gha-runner-scale-set-controller
      tag: v0.13.1
      pullPolicy: IfNotPresent
  imagePullSecrets:
    - my-registry-secret
```

## Update Strategy Examples

### Eventual Update Strategy

Wait for running jobs before applying changes:

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesGhaRunnerScaleSetController
metadata:
  name: arc-controller-safe-updates
spec:
  namespace:
    value: arc-system
  createNamespace: true
  container:
    resources:
      requests:
        cpu: 100m
        memory: 128Mi
      limits:
        cpu: 500m
        memory: 512Mi
  flags:
    updateStrategy: eventual
```

## ArgoCD Integration Example

### Excluding ArgoCD Labels

Prevent ArgoCD labels from being propagated to internal resources:

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesGhaRunnerScaleSetController
metadata:
  name: arc-controller-argocd
spec:
  namespace:
    value: arc-system
  createNamespace: true
  container:
    resources:
      requests:
        cpu: 100m
        memory: 128Mi
      limits:
        cpu: 500m
        memory: 512Mi
  flags:
    excludeLabelPropagationPrefixes:
      - argocd.argoproj.io/instance
```

## Target Cluster Examples

### Deploying to a Remote Cluster

Deploy the controller to a specific Kubernetes cluster:

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesGhaRunnerScaleSetController
metadata:
  name: arc-controller
spec:
  targetCluster:
    clusterName: production-cluster
  namespace:
    value: arc-system
  createNamespace: true
  container:
    resources:
      requests:
        cpu: 100m
        memory: 128Mi
      limits:
        cpu: 500m
        memory: 512Mi
```

## Complete Production Example

### Full-Featured Production Setup

A comprehensive production configuration:

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesGhaRunnerScaleSetController
metadata:
  name: arc-controller-production
spec:
  targetCluster:
    clusterName: production-gke-cluster
  namespace:
    value: arc-system
  createNamespace: true
  helmChartVersion: "0.13.1"
  replicaCount: 3
  priorityClassName: system-cluster-critical
  container:
    resources:
      requests:
        cpu: 200m
        memory: 256Mi
      limits:
        cpu: 1000m
        memory: 1Gi
    image:
      repository: ghcr.io/actions/gha-runner-scale-set-controller
      pullPolicy: IfNotPresent
  flags:
    logLevel: info
    logFormat: json
    runnerMaxConcurrentReconciles: 10
    updateStrategy: eventual
    k8sClientRateLimiterQps: 50
    k8sClientRateLimiterBurst: 100
    excludeLabelPropagationPrefixes:
      - argocd.argoproj.io/instance
  metrics:
    controllerManagerAddr: ":8080"
    listenerAddr: ":8080"
    listenerEndpoint: "/metrics"
```

## Post-Installation: Creating Runner Scale Sets

After deploying the controller, create AutoScalingRunnerSets:

### Example: Organization-Level Runners

```yaml
apiVersion: actions.github.com/v1alpha1
kind: AutoScalingRunnerSet
metadata:
  name: org-runners
  namespace: arc-runners
spec:
  githubConfigUrl: "https://github.com/myorg"
  githubConfigSecret: github-auth
  maxRunners: 20
  minRunners: 0
  template:
    spec:
      containers:
        - name: runner
          image: ghcr.io/actions/actions-runner:latest
          resources:
            requests:
              cpu: "1000m"
              memory: "2Gi"
            limits:
              cpu: "2000m"
              memory: "4Gi"
```

### Example: Repository-Level Runners

```yaml
apiVersion: actions.github.com/v1alpha1
kind: AutoScalingRunnerSet
metadata:
  name: repo-runners
  namespace: arc-runners
spec:
  githubConfigUrl: "https://github.com/myorg/myrepo"
  githubConfigSecret: github-auth
  maxRunners: 5
  minRunners: 0
  template:
    spec:
      containers:
        - name: runner
          image: ghcr.io/actions/actions-runner:latest
          resources:
            requests:
              cpu: "500m"
              memory: "1Gi"
```

