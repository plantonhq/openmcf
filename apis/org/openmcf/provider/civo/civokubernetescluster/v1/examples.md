# Civo Kubernetes Cluster Examples

This document provides real-world examples of Kubernetes cluster configurations using the `CivoKubernetesCluster` resource. Each example includes complete YAML manifests, cost estimates, and deployment guidance.

## Table of Contents

1. [Development Cluster (Minimal Cost)](#1-development-cluster-minimal-cost)
2. [Staging Cluster (Production-Like)](#2-staging-cluster-production-like)
3. [Production Cluster (High Availability)](#3-production-cluster-high-availability)
4. [Multi-Region Deployment](#4-multi-region-deployment)
5. [Cost-Optimized CI/CD Cluster](#5-cost-optimized-cicd-cluster)

---

## 1. Development Cluster (Minimal Cost)

**Use Case:** Single-node cluster for local development, testing, and experimentation.

**Requirements:**
- Lowest possible cost (~$11/month)
- Fast provisioning (< 2 minutes)
- Auto-upgrades enabled (convenience over stability)
- Suitable for throwaway workloads

```yaml
apiVersion: civo.openmcf.org/v1
kind: CivoKubernetesCluster
metadata:
  name: dev-cluster
spec:
  clusterName: dev-k8s
  region: nyc1  # Choose nearest region
  kubernetesVersion: "1.29.0+k3s1"
  network:
    value: "network-dev-123"
  autoUpgrade: true  # Auto-patch for convenience
  tags:
    - environment:dev
    - managed-by:planton
    - cost-center:engineering
  defaultNodePool:
    size: g4s.kube.small  # 1 core, 2 GB RAM
    nodeCount: 1
```

**Cost estimate:**
- Nodes: 1 × $10.86 = **$10.86/month**
- Total: **~$11/month**

**After deployment:**

```bash
# Deploy cluster
planton apply -f dev-cluster.yaml

# Wait for ready (~90 seconds)
planton get civokubernetesclusters/dev-cluster

# Get kubeconfig
planton kubeconfig civokubernetesclusters/dev-cluster > ~/.kube/dev-config
export KUBECONFIG=~/.kube/dev-config

# Test cluster
kubectl get nodes
kubectl get pods --all-namespaces

# Deploy test application
kubectl create deployment nginx --image=nginx
kubectl expose deployment nginx --type=LoadBalancer --port=80

# Get LoadBalancer IP
kubectl get service nginx
```

**Capabilities:**
- ✅ Full Kubernetes API
- ✅ Traefik ingress controller (pre-installed)
- ✅ Suitable for learning and development
- ❌ Not suitable for production
- ❌ No high availability
- ❌ Limited resources (single small node)

**Cleanup:**

```bash
planton delete civokubernetesclusters/dev-cluster
```

---

## 2. Staging Cluster (Production-Like)

**Use Case:** Multi-node cluster for pre-production testing, QA, and integration validation.

**Requirements:**
- Mirror production architecture (3 nodes)
- Medium-sized nodes for realistic testing
- Manual upgrades (controlled version changes)
- Support for production-like workloads

```yaml
apiVersion: civo.openmcf.org/v1
kind: CivoKubernetesCluster
metadata:
  name: staging-cluster
spec:
  clusterName: staging-k8s
  region: lon1
  kubernetesVersion: "1.29.0+k3s1"
  network:
    value: "network-staging-456"
  autoUpgrade: false  # Manual upgrades like production
  tags:
    - environment:staging
    - team:platform
    - backup:enabled
  defaultNodePool:
    size: g4s.kube.medium  # 2 cores, 4 GB RAM
    nodeCount: 3
```

**Cost estimate:**
- Nodes: 3 × $21.73 = **$65.19/month**
- Total: **~$65/month**

**Testing workflow:**

```bash
# 1. Deploy staging cluster
planton apply -f staging-cluster.yaml

# 2. Configure kubectl
planton kubeconfig civokubernetesclusters/staging-cluster > ~/.kube/staging-config
export KUBECONFIG=~/.kube/staging-config

# 3. Install monitoring (optional)
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm install prometheus prometheus-community/kube-prometheus-stack

# 4. Deploy application
kubectl apply -f manifests/

# 5. Run integration tests
./run-integration-tests.sh

# 6. Verify scaling
kubectl scale deployment myapp --replicas=10
kubectl get pods -w
```

**Capabilities:**
- ✅ Multi-node for HA testing
- ✅ Sufficient resources for realistic workloads
- ✅ Test pod distribution and affinity rules
- ✅ Validate network policies (if using Cilium)
- ⚠️ Not HA control plane (single master)

---

## 3. Production Cluster (High Availability)

**Use Case:** Mission-critical cluster for customer-facing production workloads.

**Requirements:**
- High availability (multiple control plane nodes)
- 5+ worker nodes for graceful rolling updates
- Large nodes for performance
- Manual upgrade control
- Production-grade monitoring and security

```yaml
apiVersion: civo.openmcf.org/v1
kind: CivoKubernetesCluster
metadata:
  name: production-cluster
spec:
  clusterName: prod-k8s
  region: fra1  # Frankfurt for EU users
  kubernetesVersion: "1.29.0+k3s1"
  network:
    value: "network-production-789"
  highlyAvailable: true  # HA control plane
  autoUpgrade: false  # Manual upgrades only
  disableSurgeUpgrade: false  # Allow extra resources during upgrades
  tags:
    - environment:production
    - critical:true
    - compliance:required
    - backup:daily
  defaultNodePool:
    size: g4s.kube.large  # 4 cores, 8 GB RAM
    nodeCount: 5
```

**Cost estimate:**
- Nodes: 5 × $43.45 = **$217.25/month**
- HA control plane: +$20-30/month (estimated)
- Total: **~$240/month**

**Production checklist:**

```bash
# 1. Deploy cluster
planton apply -f production-cluster.yaml

# 2. Configure kubectl
planton kubeconfig civokubernetesclusters/production-cluster > ~/.kube/prod-config
export KUBECONFIG=~/.kube/prod-config

# 3. Verify HA
kubectl get nodes
kubectl get componentstatuses  # Check control plane health

# 4. Install cert-manager for TLS
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.13.0/cert-manager.yaml

# 5. Install ingress-nginx (alternative to Traefik)
helm install ingress-nginx ingress-nginx/ingress-nginx

# 6. Set up monitoring
helm install prometheus prometheus-community/kube-prometheus-stack

# 7. Configure alerting
kubectl apply -f alertmanager-config.yaml

# 8. Deploy applications
kubectl apply -f production/

# 9. Verify readiness
kubectl get pods --all-namespaces
kubectl get ingress
```

**Production best practices:**

1. **Network policies** - Use Cilium CNI and define policies:

```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: default-deny-ingress
spec:
  podSelector: {}
  policyTypes:
    - Ingress
```

2. **Pod disruption budgets**:

```yaml
apiVersion: policy/v1
kind: PodDisruptionBudget
metadata:
  name: myapp-pdb
spec:
  minAvailable: 2
  selector:
    matchLabels:
      app: myapp
```

3. **Resource quotas**:

```yaml
apiVersion: v1
kind: ResourceQuota
metadata:
  name: namespace-quota
spec:
  hard:
    requests.cpu: "20"
    requests.memory: 40Gi
    limits.cpu: "40"
    limits.memory: 80Gi
```

---

## 4. Multi-Region Deployment

**Use Case:** Deploy clusters in multiple regions for geographic redundancy and low-latency access.

### EU Cluster (Frankfurt)

```yaml
apiVersion: civo.openmcf.org/v1
kind: CivoKubernetesCluster
metadata:
  name: eu-prod-cluster
spec:
  clusterName: eu-prod-k8s
  region: fra1
  kubernetesVersion: "1.29.0+k3s1"
  network:
    value: "network-eu-prod"
  highlyAvailable: true
  autoUpgrade: false
  tags:
    - region:eu
    - environment:production
  defaultNodePool:
    size: g4s.kube.large
    nodeCount: 5
```

### US Cluster (New York)

```yaml
apiVersion: civo.openmcf.org/v1
kind: CivoKubernetesCluster
metadata:
  name: us-prod-cluster
spec:
  clusterName: us-prod-k8s
  region: nyc1
  kubernetesVersion: "1.29.0+k3s1"
  network:
    value: "network-us-prod"
  highlyAvailable: true
  autoUpgrade: false
  tags:
    - region:us
    - environment:production
  defaultNodePool:
    size: g4s.kube.large
    nodeCount: 5
```

**Multi-region strategy:**

1. **Deploy both clusters**:

```bash
planton apply -f eu-cluster.yaml
planton apply -f us-cluster.yaml
```

2. **Configure geo-routing** (use DNS or CDN):

```yaml
# CivoDnsZone with geo-specific subdomains
apiVersion: civo.openmcf.org/v1
kind: CivoDnsZone
metadata:
  name: multi-region-dns
spec:
  domainName: myapp.com
  records:
    - name: "eu"
      type: A
      values:
        - value: "<eu-cluster-lb-ip>"
      ttlSeconds: 300
    - name: "us"
      type: A
      values:
        - value: "<us-cluster-lb-ip>"
      ttlSeconds: 300
```

3. **Sync application deployments** across regions (GitOps):

```bash
# Use ArgoCD or Flux to deploy to both clusters
argocd app create myapp-eu --repo ... --dest-server eu-cluster
argocd app create myapp-us --repo ... --dest-server us-cluster
```

**Cost estimate:**
- 2 clusters × $217/month = **$434/month**

---

## 5. Cost-Optimized CI/CD Cluster

**Use Case:** Dedicated cluster for CI/CD pipelines (GitLab Runner, Argo Workflows, Tekton).

**Requirements:**
- Cost-effective (single small node)
- Auto-upgrades enabled (low-risk workloads)
- Ephemeral workloads (can tolerate restarts)

```yaml
apiVersion: civo.openmcf.org/v1
kind: CivoKubernetesCluster
metadata:
  name: cicd-cluster
spec:
  clusterName: cicd-k8s
  region: lon1
  kubernetesVersion: "1.29.0+k3s1"
  network:
    value: "network-cicd"
  autoUpgrade: true
  tags:
    - purpose:cicd
    - environment:shared
  defaultNodePool:
    size: g4s.kube.medium  # Sufficient for CI runners
    nodeCount: 2  # One for runners, one for resilience
```

**Cost estimate:**
- Nodes: 2 × $21.73 = **$43.46/month**

**Setup CI/CD runners:**

```bash
# 1. Deploy cluster
planton apply -f cicd-cluster.yaml

# 2. Get kubeconfig
planton kubeconfig civokubernetesclusters/cicd-cluster > ~/.kube/cicd-config
export KUBECONFIG=~/.kube/cicd-config

# 3. Install GitLab Runner
helm repo add gitlab https://charts.gitlab.io
helm install gitlab-runner gitlab/gitlab-runner \
  --set gitlabUrl=https://gitlab.com \
  --set runnerRegistrationToken=$REGISTRATION_TOKEN

# 4. Or install Argo Workflows
kubectl create namespace argo
kubectl apply -n argo -f https://raw.githubusercontent.com/argoproj/argo-workflows/stable/manifests/install.yaml

# 5. Verify installation
kubectl get pods -n default  # GitLab Runner
kubectl get pods -n argo     # Argo Workflows
```

**Runner configuration:**

```yaml
# GitLab Runner values.yaml
runners:
  config: |
    [[runners]]
      [runners.kubernetes]
        namespace = "default"
        image = "alpine:latest"
        cpu_request = "100m"
        memory_request = "128Mi"
```

---

## Additional Patterns

### Reference Network from CivoVpc

Instead of hardcoding network ID:

```yaml
apiVersion: civo.openmcf.org/v1
kind: CivoVpc
metadata:
  name: k8s-network
spec:
  label: Kubernetes network
  region: lon1
---
apiVersion: civo.openmcf.org/v1
kind: CivoKubernetesCluster
metadata:
  name: my-cluster
spec:
  clusterName: my-k8s
  region: lon1
  kubernetesVersion: "1.29.0+k3s1"
  network:
    valueFrom:
      kind: CivoVpc
      name: k8s-network
      env: production
      fieldPath: "status.outputs.network_id"
  defaultNodePool:
    size: g4s.kube.medium
    nodeCount: 3
```

### Cluster with All Optional Features

```yaml
apiVersion: civo.openmcf.org/v1
kind: CivoKubernetesCluster
metadata:
  name: full-featured-cluster
spec:
  clusterName: full-featured-k8s
  region: lon1
  kubernetesVersion: "1.29.0+k3s1"
  network:
    value: "network-full"
  highlyAvailable: true
  autoUpgrade: false
  disableSurgeUpgrade: false
  tags:
    - environment:production
    - ha:enabled
    - monitoring:prometheus
    - backup:velero
  defaultNodePool:
    size: g4s.kube.xlarge
    nodeCount: 7
```

---

## Deployment Workflows

### Blue-Green Cluster Deployment

Deploy new cluster alongside existing one for zero-downtime upgrades:

```bash
# 1. Deploy new "green" cluster
planton apply -f green-cluster.yaml

# 2. Deploy applications to green
kubectl --context=green apply -f manifests/

# 3. Test green cluster
./smoke-tests.sh green-cluster

# 4. Switch DNS/load balancer to green
kubectl --context=green get service ingress-nginx -o jsonpath='{.status.loadBalancer.ingress[0].ip}'
# Update DNS A record to green cluster IP

# 5. Wait for traffic migration (monitor)

# 6. Delete old "blue" cluster
planton delete civokubernetesclusters/blue-cluster
```

### Rolling Upgrade Strategy

For in-place upgrades:

```bash
# 1. Lower TTL on DNS records (1 hour before upgrade)
# Update DNS TTL to 300 seconds

# 2. Test new version in dev/staging
planton apply -f cluster-v1.30.yaml --env=dev

# 3. Schedule maintenance window

# 4. Upgrade production
# Edit kubernetesVersion in cluster.yaml to "1.30.0+k3s1"
planton apply -f cluster.yaml --env=production

# 5. Monitor upgrade
kubectl get nodes -w
kubectl get pods --all-namespaces -w

# 6. Verify application health
kubectl exec -it <pod> -- curl http://localhost:8080/health

# 7. Restore DNS TTL after success
```

---

## Post-Deployment Configuration

### Install Essential Add-Ons

```bash
# cert-manager for automated TLS
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.13.0/cert-manager.yaml

# Metrics Server (if not included)
kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml

# Kubernetes Dashboard
kubectl apply -f https://raw.githubusercontent.com/kubernetes/dashboard/v2.7.0/aio/deploy/recommended.yaml
```

### Set Up Cluster Autoscaler

For dynamic node scaling:

```bash
# Install from Civo marketplace during cluster creation
# Or deploy manually
helm repo add autoscaler https://kubernetes.github.io/autoscaler
helm install cluster-autoscaler autoscaler/cluster-autoscaler \
  --set autoDiscovery.clusterName=prod-k8s \
  --set cloudProvider=civo
```

### Configure Backups with Velero

```bash
# Install Velero
velero install \
  --provider civo \
  --bucket k8s-backups \
  --backup-location-config region=lon1

# Create backup schedule
velero schedule create daily-backup --schedule="@daily"

# Test restore
velero restore create --from-backup daily-backup-20251116
```

---

## Troubleshooting

### Cluster Stuck in "Building"

```bash
# Check status via API
civo kubernetes show prod-k8s

# Common causes:
# - Network quota exceeded
# - Insufficient IP addresses in VPC
# - Region capacity issues

# Solution: Check quotas
civo quota show
```

### Nodes Not Ready

```bash
# Describe node
kubectl describe node <node-name>

# Check events
kubectl get events --sort-by='.lastTimestamp'

# SSH to node (if accessible)
civo kubernetes show prod-k8s  # Get node IPs
ssh civo@<node-ip>
sudo journalctl -u k3s-agent -f
```

### LoadBalancer Stuck in Pending

```bash
# Check service
kubectl describe service my-service

# Verify Civo CCM is running
kubectl get pods -n kube-system | grep civo

# Check quotas
civo quota show  # Load balancer quota
```

---

## Cost Optimization Tips

1. **Right-size nodes** - Start with `medium`, scale up if needed
2. **Use single node for dev** - No need for multi-node in development
3. **Delete unused clusters** - Don't leave staging running 24/7
4. **Use spot instances** (when available) - 60-70% savings
5. **Set resource requests/limits** - Prevent resource waste
6. **Use namespace quotas** - Limit runaway pods

---

## Additional Resources

- [Main README](README.md) - Component overview and quick start
- [Research Documentation](docs/README.md) - Deep dive into K3s, CNI decisions, and production architecture
- [Pulumi Module](iac/pulumi/README.md) - Direct Pulumi usage
- [Civo Kubernetes Docs](https://www.civo.com/docs/kubernetes)
- [K3s Documentation](https://docs.k3s.io/)

---

## Need Help?

- Check the [Troubleshooting section](README.md#troubleshooting) in the main README
- Open an issue on [GitHub](https://github.com/plantonhq/openmcf/issues)
- Contact Civo support: support@civo.com

