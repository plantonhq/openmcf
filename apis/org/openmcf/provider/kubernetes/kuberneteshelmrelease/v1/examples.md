# Kubernetes Helm Release - Example Configurations

This document provides comprehensive examples for deploying Helm charts using the KubernetesHelmRelease API resource.

## Table of Contents

1. [Basic Helm Release Deployment](#example-1-basic-helm-release-deployment)
2. [Helm Release with Custom Values](#example-2-helm-release-with-custom-values)
3. [Helm Release with Ingress Configuration](#example-3-helm-release-with-ingress-configuration)
4. [Helm Release with Multiple Values](#example-4-helm-release-with-multiple-values)
5. [Production PostgreSQL Deployment](#example-5-production-postgresql-deployment)
6. [Monitoring Stack with Prometheus](#example-6-monitoring-stack-with-prometheus)
7. [Private Helm Repository with Authentication](#example-7-private-helm-repository-with-authentication)
8. [OCI Registry Helm Chart](#example-8-oci-registry-helm-chart)
9. [Multi-Environment Configuration (Development)](#example-9-multi-environment-configuration-development)
10. [Multi-Environment Configuration (Production)](#example-10-multi-environment-configuration-production)

---

## Example 1: Basic Helm Release Deployment

A minimal example deploying NGINX from the Bitnami Helm repository with default values.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: HelmRelease
metadata:
  name: basic-nginx
spec:
  target_cluster:
    cluster_name: my-cluster
  namespace:
    value: default
  create_namespace: true
  repo: https://charts.bitnami.com/bitnami
  name: nginx
  version: 15.14.0
  values: {}
```

**Use Case:** Quick deployment for testing or development environments.

---

## Example 2: Helm Release with Custom Values

Deploy NGINX with a LoadBalancer service and 3 replicas for high availability.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: HelmRelease
metadata:
  name: nginx-ha
spec:
  target_cluster:
    cluster_name: my-cluster
  namespace:
    value: default
  create_namespace: true
  repo: https://charts.bitnami.com/bitnami
  name: nginx
  version: 15.14.0
  values:
    replicaCount: "3"
    service:
      type: "LoadBalancer"
    resources:
      requests:
        memory: "256Mi"
        cpu: "100m"
      limits:
        memory: "512Mi"
        cpu: "500m"
```

**Use Case:** Production deployment with resource limits and multiple replicas.

---

## Example 3: Helm Release with Ingress Configuration

Deploy WordPress with Ingress enabled for external access via a custom domain.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: HelmRelease
metadata:
  name: wordpress-blog
spec:
  target_cluster:
    cluster_name: my-cluster
  namespace:
    value: wordpress
  create_namespace: true
  repo: https://charts.bitnami.com/bitnami
  name: wordpress
  version: 19.0.0
  values:
    wordpressUsername: "admin"
    wordpressEmail: "admin@example.com"
    ingress:
      enabled: "true"
      hostname: "blog.example.com"
      pathType: "Prefix"
      ingressClassName: "nginx"
    service:
      type: "ClusterIP"
    persistence:
      enabled: "true"
      size: "10Gi"
```

**Use Case:** Public-facing blog or website with persistent storage.

---

## Example 4: Helm Release with Multiple Values

Deploy Redis with Sentinel for high availability, metrics enabled, and authentication.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: HelmRelease
metadata:
  name: redis-ha
spec:
  target_cluster:
    cluster_name: my-cluster
  namespace:
    value: redis
  create_namespace: true
  repo: https://charts.bitnami.com/bitnami
  name: redis
  version: 18.19.0
  values:
    architecture: "replication"
    auth:
      enabled: "true"
      password: "super-secret-password"
    sentinel:
      enabled: "true"
      quorum: "2"
    metrics:
      enabled: "true"
      serviceMonitor:
        enabled: "true"
    replica:
      replicaCount: "3"
    master:
      persistence:
        enabled: "true"
        size: "8Gi"
```

**Use Case:** Production Redis cluster with monitoring and high availability.

---

## Example 5: Using Existing Namespace

Deploy NGINX to an existing namespace without creating it.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: HelmRelease
metadata:
  name: nginx-existing-ns
spec:
  target_cluster:
    cluster_name: my-cluster
  namespace:
    value: existing-namespace
  create_namespace: false
  repo: https://charts.bitnami.com/bitnami
  name: nginx
  version: 15.14.0
  values: {}
```

**Use Case:** Deploy to a namespace that already exists and is managed separately.

---

## Example 6: Creating Namespace with Helm Release (Default)

Explicitly create the namespace along with the Helm release.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: HelmRelease
metadata:
  name: redis-new-ns
spec:
  target_cluster:
    cluster_name: my-cluster
  namespace:
    value: redis-app
  create_namespace: true
  repo: https://charts.bitnami.com/bitnami
  name: redis
  version: 18.19.0
  values: {}
```

**Use Case:** Create both namespace and Helm release together (default behavior).

---

## Example 7: Production PostgreSQL Deployment

Deploy PostgreSQL with replication, persistence, and backup configuration.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: HelmRelease
metadata:
  name: postgres-production
  labels:
    environment: production
    team: platform
spec:
  target_cluster:
    cluster_name: prod-cluster
  namespace:
    value: postgres
  create_namespace: true
  repo: https://charts.bitnami.com/bitnami
  name: postgresql
  version: 14.3.0
  values:
    global:
      postgresql:
        auth:
          username: "appuser"
          password: "changeme"
          database: "myapp"
    primary:
      persistence:
        enabled: "true"
        size: "50Gi"
        storageClass: "fast-ssd"
      resources:
        requests:
          memory: "1Gi"
          cpu: "500m"
        limits:
          memory: "2Gi"
          cpu: "2000m"
    readReplicas:
      replicaCount: "2"
      persistence:
        enabled: "true"
        size: "50Gi"
    metrics:
      enabled: "true"
      serviceMonitor:
        enabled: "true"
```

**Use Case:** Production database with read replicas and monitoring.

---

## Example 8: Monitoring Stack with Prometheus

Deploy Prometheus for cluster monitoring with persistent storage and retention.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: HelmRelease
metadata:
  name: prometheus-monitoring
spec:
  target_cluster:
    cluster_name: my-cluster
  namespace:
    value: monitoring
  create_namespace: true
  repo: https://prometheus-community.github.io/helm-charts
  name: prometheus
  version: 25.11.0
  values:
    server:
      persistentVolume:
        enabled: "true"
        size: "100Gi"
      retention: "30d"
      resources:
        requests:
          cpu: "500m"
          memory: "2Gi"
        limits:
          cpu: "2000m"
          memory: "4Gi"
    alertmanager:
      enabled: "true"
      persistentVolume:
        enabled: "true"
        size: "10Gi"
    nodeExporter:
      enabled: "true"
    pushgateway:
      enabled: "true"
    serverFiles:
      prometheus.yml:
        scrape_configs:
          - job_name: "kubernetes-pods"
            kubernetes_sd_configs:
              - role: "pod"
```

**Use Case:** Comprehensive cluster monitoring with alerts.

---

## Example 9: Private Helm Repository with Authentication

Deploy a chart from a private Helm repository requiring authentication.

**Note:** This example shows the structure. In practice, use Kubernetes secrets for credentials.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: HelmRelease
metadata:
  name: private-app
spec:
  target_cluster:
    cluster_name: my-cluster
  namespace:
    value: default
  create_namespace: true
  repo: https://private-charts.company.com/charts
  name: internal-app
  version: 2.1.0
  values:
    image:
      repository: "company-registry.com/internal-app"
      tag: "v2.1.0"
      pullSecrets:
        - "regcred"
    replicaCount: "2"
    env:
      - name: "DATABASE_URL"
        valueFrom:
          secretKeyRef:
            name: "app-secrets"
            key: "db-url"
```

**Use Case:** Deploying proprietary applications from private repositories.

**Setup Instructions:**
1. Create image pull secret: `kubectl create secret docker-registry regcred --docker-server=company-registry.com --docker-username=user --docker-password=pass`
2. Create app secrets: `kubectl create secret generic app-secrets --from-literal=db-url=postgresql://...`

---

## Example 10: OCI Registry Helm Chart

Deploy a Helm chart stored in an OCI-compliant registry (e.g., GitHub Container Registry, Amazon ECR).

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: HelmRelease
metadata:
  name: oci-chart-deployment
spec:
  target_cluster:
    cluster_name: my-cluster
  namespace:
    value: default
  create_namespace: true
  repo: oci://ghcr.io/company/charts
  name: myapp
  version: 1.5.0
  values:
    image:
      registry: "ghcr.io"
      repository: "company/myapp"
      tag: "1.5.0"
    service:
      type: "ClusterIP"
      port: "8080"
    ingress:
      enabled: "true"
      className: "nginx"
      hosts:
        - "myapp.example.com"
```

**Use Case:** Modern chart distribution using OCI registries.

---

## Example 11: Multi-Environment Configuration (Development)

Development environment configuration with minimal resources and debugging enabled.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: HelmRelease
metadata:
  name: myapp-dev
  labels:
    environment: development
    team: engineering
spec:
  target_cluster:
    cluster_name: dev-cluster
  namespace:
    value: myapp-dev
  create_namespace: true
  repo: https://charts.company.com/stable
  name: myapp
  version: 3.2.0
  values:
    environment: "development"
    replicaCount: "1"
    image:
      tag: "develop"
      pullPolicy: "Always"
    resources:
      requests:
        cpu: "100m"
        memory: "256Mi"
      limits:
        cpu: "500m"
        memory: "512Mi"
    debug:
      enabled: "true"
    autoscaling:
      enabled: "false"
    ingress:
      enabled: "true"
      hostname: "myapp-dev.internal.company.com"
      tls: "false"
    database:
      host: "postgres-dev.internal"
      name: "myapp_dev"
```

**Use Case:** Development environment with debugging and lower resource limits.

---

## Example 12: Multi-Environment Configuration (Production)

Production environment with high availability, autoscaling, and security features.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: HelmRelease
metadata:
  name: myapp-prod
  labels:
    environment: production
    team: engineering
    criticality: high
spec:
  target_cluster:
    cluster_name: prod-cluster
  namespace:
    value: myapp-prod
  create_namespace: true
  repo: https://charts.company.com/stable
  name: myapp
  version: 3.2.0
  values:
    environment: "production"
    replicaCount: "3"
    image:
      tag: "3.2.0"
      pullPolicy: "IfNotPresent"
    resources:
      requests:
        cpu: "1000m"
        memory: "2Gi"
      limits:
        cpu: "4000m"
        memory: "4Gi"
    autoscaling:
      enabled: "true"
      minReplicas: "3"
      maxReplicas: "10"
      targetCPUUtilizationPercentage: "70"
      targetMemoryUtilizationPercentage: "80"
    podDisruptionBudget:
      enabled: "true"
      minAvailable: "2"
    ingress:
      enabled: "true"
      hostname: "myapp.company.com"
      tls: "true"
      annotations:
        cert-manager.io/cluster-issuer: "letsencrypt-prod"
    database:
      host: "postgres-prod.internal"
      name: "myapp_prod"
      connectionPooling: "true"
    monitoring:
      enabled: "true"
      serviceMonitor: "true"
    securityContext:
      runAsNonRoot: "true"
      runAsUser: "1000"
      fsGroup: "1000"
```

**Use Case:** Production deployment with HA, autoscaling, and security hardening.

---

## Deployment Instructions

### Using OpenMCF CLI

```bash
# Deploy a Helm release
planton apply -f <example-file>.yaml

# Check deployment status
planton get helmrelease <name>

# View logs
kubectl logs -n <namespace> -l app=<app-name>

# Delete deployment
planton delete -f <example-file>.yaml
```

### Verify Deployment

```bash
# Check Helm releases
helm list -A

# Check pods in namespace
kubectl get pods -n <namespace>

# Check services
kubectl get svc -n <namespace>

# Check ingress (if enabled)
kubectl get ingress -n <namespace>
```

## Troubleshooting

### Common Issues

1. **Chart Not Found**
   - Verify the repo URL is accessible
   - Check that the chart name and version exist
   - For private repos, ensure authentication is configured

2. **Values Not Applied**
   - Verify values syntax (use quotes for boolean/numeric strings)
   - Check Helm chart's values.yaml for correct key names
   - Use `helm show values <repo/chart>` to see available options

3. **Resource Conflicts**
   - Ensure namespace doesn't have conflicting resources
   - Check for port conflicts with existing services
   - Verify ingress hostnames are unique

### Debug Commands

```bash
# Get Helm release status
helm status <release-name> -n <namespace>

# See rendered manifests
helm get manifest <release-name> -n <namespace>

# View Helm values
helm get values <release-name> -n <namespace>

# Helm history
helm history <release-name> -n <namespace>
```

## Best Practices

1. **Version Pinning**: Always specify exact chart versions in production
2. **Resource Limits**: Set appropriate CPU and memory limits
3. **Persistence**: Enable for stateful applications
4. **Monitoring**: Include metrics and health checks
5. **Security**: Run as non-root, use network policies
6. **High Availability**: Use multiple replicas and pod disruption budgets
7. **Secrets Management**: Never hardcode secrets; use Kubernetes secrets or external secret managers

## Additional Resources

- [Helm Official Documentation](https://helm.sh/docs/)
- [Artifact Hub](https://artifacthub.io/) - Search for Helm charts
- [Bitnami Charts](https://github.com/bitnami/charts) - Popular Helm charts
- [Helm Best Practices](https://helm.sh/docs/chart_best_practices/)
