# Kubernetes Ingress NGINX - Example Configurations

This document provides comprehensive examples for deploying the NGINX Ingress Controller using the KubernetesIngressNginx API resource.

## Table of Contents

1. [Basic External Ingress Controller](#example-1-basic-external-ingress-controller)
2. [Internal Ingress Controller](#example-2-internal-ingress-controller)
3. [GKE with Static IP](#example-3-gke-with-static-ip)
4. [GKE Internal Load Balancer](#example-4-gke-internal-load-balancer)
5. [EKS with Network Load Balancer](#example-5-eks-with-network-load-balancer)
6. [EKS Internal with Specific Subnets](#example-6-eks-internal-with-specific-subnets)
7. [AKS with Managed Identity](#example-7-aks-with-managed-identity)
8. [AKS with Reserved Public IP](#example-8-aks-with-reserved-public-ip)
9. [Development Environment Setup](#example-9-development-environment-setup)
10. [Production Multi-Cloud Setup](#example-10-production-multi-cloud-setup)

---

## Example 1: Basic External Ingress Controller

The simplest deployment for a generic Kubernetes cluster with an external load balancer.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesIngressNginx
metadata:
  name: basic-ingress
spec:
  target_cluster:
    cluster_name: my-cluster
  namespace:
    value: ingress-nginx
  create_namespace: true
  chart_version: "4.11.1"
  internal: false
```

**Use Case:** Development or testing on any Kubernetes cluster (generic, on-prem, or cloud).

**Result:**
- Creates namespace: `kubernetes-ingress-nginx`
- Deploys NGINX Ingress Controller with external LoadBalancer
- Service accessible from internet

---

## Example 2: Internal Ingress Controller

Deploy an ingress controller with an internal load balancer for private applications.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesIngressNginx
metadata:
  name: internal-ingress
spec:
  target_cluster:
    cluster_name: my-cluster
  namespace:
    value: ingress-nginx
  create_namespace: true
  chart_version: "4.11.1"
  internal: true
```

**Use Case:** Internal applications, microservices, admin panels accessible only from VPN/VPC.

**Result:**
- Creates internal load balancer
- Not accessible from public internet
- Accessible from within VPC/VPN

---

## Example 2b: Using Existing Namespace

Deploy to a pre-existing namespace managed by your platform team.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesIngressNginx
metadata:
  name: ingress-with-existing-ns
spec:
  target_cluster:
    cluster_name: my-cluster
  namespace:
    value: platform-ingress
  create_namespace: false
  chart_version: "4.11.1"
  internal: false
```

**Prerequisites:**
1. Namespace must already exist:
   ```bash
   kubectl create namespace platform-ingress
   ```

**Use Case:** Organizations with centralized namespace management, shared infrastructure namespaces.

**Result:**
- Uses existing namespace `platform-ingress`
- No namespace creation attempted
- Works with existing RBAC, policies, quotas

---

## Example 3: GKE with Static IP

Deploy on Google Kubernetes Engine with a reserved static IP address.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesIngressNginx
metadata:
  name: gke-static-ingress
  labels:
    environment: production
    cloud: gcp
spec:
  target_cluster:
    cluster_name: gke-prod-cluster
  namespace:
    value: ingress-nginx
  create_namespace: true
  chart_version: "4.11.1"
  internal: false
  gke:
    static_ip_name: my-ingress-static-ip
```

**Prerequisites:**
1. Reserve a static IP in GCP:
   ```bash
   gcloud compute addresses create my-ingress-static-ip --global
   ```

**Use Case:** Production deployments requiring a stable IP for DNS configuration.

**Result:**
- External load balancer with static IP
- Consistent IP address for DNS records
- Automatic SSL certificate renewal works seamlessly

---

## Example 4: GKE Internal Load Balancer

Deploy an internal load balancer on GKE with a specific subnetwork.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesIngressNginx
metadata:
  name: gke-internal-ingress
spec:
  target_cluster:
    cluster_name: gke-prod-cluster
  namespace:
    value: ingress-nginx
  create_namespace: true
  chart_version: "4.11.1"
  internal: true
  gke:
    subnetwork_self_link: projects/my-project/regions/us-west1/subnetworks/private-subnet
```

**Use Case:** Internal applications accessible only from specific VPC subnetworks.

**Result:**
- Internal load balancer in specified subnet
- Accessible only from VPC
- No public IP assigned

---

## Example 5: EKS with Network Load Balancer

Deploy on Amazon EKS with Network Load Balancer and additional security groups.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesIngressNginx
metadata:
  name: eks-ingress
  labels:
    environment: production
    cloud: aws
spec:
  target_cluster:
    cluster_name: eks-prod-cluster
  namespace:
    value: ingress-nginx
  create_namespace: true
  chart_version: "4.11.1"
  internal: false
  eks:
    additional_security_group_ids:
      - value: sg-0123456789abcdef0
      - value: sg-fedcba9876543210
    irsa_role_arn_override: arn:aws:iam::123456789012:role/ingress-nginx-role
```

**Prerequisites:**
1. Create IAM role for IRSA (optional)
2. Create security groups for load balancer

**Use Case:** Production EKS deployment with custom security policies.

**Result:**
- NLB with specified security groups
- IRSA role attached for AWS API access
- Enhanced security posture

---

## Example 6: EKS Internal with Specific Subnets

Deploy an internal NLB on EKS in specific private subnets.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesIngressNginx
metadata:
  name: eks-internal-ingress
spec:
  target_cluster:
    cluster_name: eks-prod-cluster
  namespace:
    value: ingress-nginx
  create_namespace: true
  chart_version: "4.11.1"
  internal: true
  eks:
    subnet_ids:
      - value: subnet-private1
      - value: subnet-private2
      - value: subnet-private3
```

**Use Case:** Internal applications with precise subnet control for compliance.

**Result:**
- Internal NLB in specified subnets
- Accessible only from VPC
- Multi-AZ deployment for HA

---

## Example 7: AKS with Managed Identity

Deploy on Azure Kubernetes Service with Workload Identity.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesIngressNginx
metadata:
  name: aks-ingress
  labels:
    environment: production
    cloud: azure
spec:
  target_cluster:
    cluster_name: aks-prod-cluster
  namespace:
    value: ingress-nginx
  create_namespace: true
  chart_version: "4.11.1"
  internal: false
  aks:
    managed_identity_client_id: 12345678-1234-1234-1234-123456789012
```

**Prerequisites:**
1. Create user-assigned managed identity
2. Configure workload identity on AKS cluster

**Use Case:** Production AKS with Azure-native identity management.

**Result:**
- Azure Load Balancer
- Workload Identity integration
- Seamless Azure service access

---

## Example 8: AKS with Reserved Public IP

Deploy on AKS reusing an existing public IP resource.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesIngressNginx
metadata:
  name: aks-static-ingress
spec:
  target_cluster:
    cluster_name: aks-prod-cluster
  namespace:
    value: ingress-nginx
  create_namespace: true
  chart_version: "4.11.1"
  internal: false
  aks:
    public_ip_name: my-ingress-public-ip
```

**Prerequisites:**
1. Create public IP in Azure:
   ```bash
   az network public-ip create \
     --resource-group myResourceGroup \
     --name my-ingress-public-ip \
     --sku Standard \
     --allocation-method Static
   ```

**Use Case:** Consistent IP for DNS records on Azure.

**Result:**
- Load balancer uses specified public IP
- Stable IP for A records
- No IP changes on redeployment

---

## Example 9: Development Environment Setup

Lightweight setup for development clusters with automatic defaults.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesIngressNginx
metadata:
  name: dev-ingress
  labels:
    environment: development
spec:
  target_cluster:
    cluster_name: dev-cluster
  namespace:
    value: ingress-nginx
  create_namespace: true
  # chart_version omitted - uses default stable version
  internal: false
```

**Use Case:** Quick setup for development and testing.

**Result:**
- Uses default stable chart version
- External load balancer (may not be supported on local clusters like Minikube)
- Minimal configuration

---

## Example 10: Production Multi-Cloud Setup

Complete production example showing GKE, EKS, and AKS configurations side-by-side.

### GCP Production

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesIngressNginx
metadata:
  name: prod-gke-ingress
  labels:
    environment: production
    cloud: gcp
    region: us-west1
spec:
  target_cluster:
    cluster_name: prod-gke-cluster
  namespace:
    value: ingress-nginx
  create_namespace: true
  chart_version: "4.11.1"
  internal: false
  gke:
    static_ip_name: prod-gke-ingress-ip
```

### AWS Production

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesIngressNginx
metadata:
  name: prod-eks-ingress
  labels:
    environment: production
    cloud: aws
    region: us-west-2
spec:
  target_cluster:
    cluster_name: prod-eks-cluster
  namespace:
    value: ingress-nginx
  create_namespace: true
  chart_version: "4.11.1"
  internal: false
  eks:
    subnet_ids:
      - value: subnet-public-1a
      - value: subnet-public-1b
      - value: subnet-public-1c
    irsa_role_arn_override: arn:aws:iam::123456789012:role/prod-ingress-nginx
```

### Azure Production

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesIngressNginx
metadata:
  name: prod-aks-ingress
  labels:
    environment: production
    cloud: azure
    region: westus2
spec:
  target_cluster:
    cluster_name: prod-aks-cluster
  namespace:
    value: ingress-nginx
  create_namespace: true
  chart_version: "4.11.1"
  internal: false
  aks:
    managed_identity_client_id: 12345678-1234-1234-1234-123456789012
    public_ip_name: prod-aks-ingress-ip
```

**Use Case:** Consistent production deployments across multiple cloud providers.

---

## Namespace Management Patterns

### Pattern 1: Component Manages Namespace (Default)

The component creates and manages its own namespace:

```yaml
spec:
  namespace:
    value: my-ingress-namespace
  create_namespace: true
```

### Pattern 2: Platform Team Manages Namespace

Platform team creates namespace with policies, component just uses it:

```bash
# Platform team
kubectl create namespace shared-ingress
kubectl apply -f network-policies.yaml
kubectl apply -f resource-quotas.yaml
```

```yaml
# Application team
spec:
  namespace:
    value: shared-ingress
  create_namespace: false
```

### Pattern 3: Reference KubernetesNamespace Resource

Use a KubernetesNamespace resource reference:

```yaml
spec:
  namespace:
    ref: my-namespace-resource-id
  create_namespace: false  # Referenced namespace already exists
```

---

## Deployment Instructions

### Using OpenMCF CLI

```bash
# Deploy ingress controller
planton apply -f <example-file>.yaml

# Check deployment status
planton get kubernetesingressnginx <name>

# View ingress controller pods
kubectl get pods -n kubernetes-ingress-nginx

# Check load balancer service
kubectl get svc -n kubernetes-ingress-nginx

# View load balancer IP/hostname
kubectl get svc -n kubernetes-ingress-nginx kubernetes-ingress-nginx-controller

# Delete ingress controller
planton delete -f <example-file>.yaml
```

### Verify Deployment

```bash
# Check Helm release
helm list -n kubernetes-ingress-nginx

# Check all resources
kubectl get all -n kubernetes-ingress-nginx

# View ingress controller logs
kubectl logs -n kubernetes-ingress-nginx -l app.kubernetes.io/name=kubernetes-ingress-nginx

# Test ingress (create a test ingress resource)
kubectl apply -f - <<EOF
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: test-ingress
  namespace: default
spec:
  ingressClassName: nginx
  rules:
  - host: test.example.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: test-service
            port:
              number: 80
EOF
```

## Troubleshooting

### Load Balancer Not Getting External IP

**GKE:**
```bash
# Check if static IP is created
gcloud compute addresses list

# Verify GKE has necessary permissions
gcloud projects get-iam-policy <project-id>
```

**EKS:**
```bash
# Check subnet configuration
aws ec2 describe-subnets --subnet-ids <subnet-id>

# Verify security groups
aws ec2 describe-security-groups --group-ids <sg-id>
```

**AKS:**
```bash
# Check public IP
az network public-ip show --name <ip-name> --resource-group <rg>

# Verify load balancer
az network lb list --resource-group <rg>
```

### Ingress Not Routing Traffic

```bash
# Check ingress controller is running
kubectl get pods -n kubernetes-ingress-nginx

# View controller logs
kubectl logs -n kubernetes-ingress-nginx -l app.kubernetes.io/component=controller

# Check ingress resource
kubectl describe ingress <ingress-name>

# Verify ingress class
kubectl get ingressclass
```

### Certificate Issues

```bash
# Check cert-manager (if using)
kubectl get certificates -A

# View certificate request status
kubectl describe certificaterequest <request-name>
```

## Best Practices

1. **Production Deployments**
   - Use specific chart versions (not "latest")
   - Configure resource limits on the controller pods
   - Set up monitoring and alerting
   - Use cloud-native load balancer features (static IPs, security groups)

2. **Security**
   - Use internal load balancers for non-public services
   - Configure network policies
   - Enable TLS/SSL termination
   - Restrict ingress sources with security groups/firewall rules

3. **High Availability**
   - Deploy controller with multiple replicas
   - Use pod anti-affinity rules
   - Configure pod disruption budgets
   - Spread across availability zones

4. **Cost Optimization**
   - Use internal load balancers where possible (cheaper than external)
   - Share ingress controllers across namespaces
   - Monitor and optimize resource requests

5. **Monitoring**
   - Enable Prometheus metrics
   - Set up alerts for controller health
   - Monitor load balancer health checks
   - Track request rates and error rates

## Common Patterns

### Separate Public and Internal Ingress

Deploy two ingress controllers for different traffic types:

```yaml
# Public ingress
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesIngressNginx
metadata:
  name: public-ingress
spec:
  target_cluster:
    cluster_name: my-cluster
  namespace:
    value: ingress-nginx
  create_namespace: true
  chart_version: "4.11.1"
  internal: false
---
# Internal ingress
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesIngressNginx
metadata:
  name: internal-ingress
spec:
  target_cluster:
    cluster_name: my-cluster
  namespace:
    value: ingress-nginx
  create_namespace: true
  chart_version: "4.11.1"
  internal: true
```

### Using Cluster Selector

Reference a cluster in the same environment:

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesIngressNginx
metadata:
  name: ingress-with-selector
  env: production
spec:
  target_cluster:
    cluster_name: my-prod-cluster
  namespace:
    value: ingress-nginx
  create_namespace: true
  chart_version: "4.11.1"
  internal: false
```

## Configuration Reference

### Cloud-Specific Fields

#### GKE Configuration

- `static_ip_name`: Name of reserved static IP address (global or regional)
- `subnetwork_self_link`: Subnetwork for internal load balancers

#### EKS Configuration

- `additional_security_group_ids`: Security groups to attach to NLB
- `subnet_ids`: Specific subnets for load balancer placement
- `irsa_role_arn_override`: IAM role ARN for IRSA

#### AKS Configuration

- `managed_identity_client_id`: Client ID for Workload Identity
- `public_ip_name`: Pre-existing public IP resource name

## Accessing Applications

Once deployed, create Ingress resources to route traffic:

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: myapp-ingress
  namespace: myapp
  annotations:
    nginx.ingress.kubernetes.io/rewrite-target: /
spec:
  ingressClassName: nginx
  rules:
  - host: myapp.example.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: myapp
            port:
              number: 80
```

## Upgrading

### Chart Version Upgrade

Update the `chart_version` field:

```yaml
spec:
  chart_version: "4.12.0"  # Updated from 4.11.1
```

Then apply:

```bash
planton apply -f ingress.yaml
```

## Additional Resources

- [NGINX Ingress Controller Docs](https://kubernetes.github.io/ingress-nginx/)
- [Ingress NGINX Annotations](https://kubernetes.github.io/ingress-nginx/user-guide/nginx-configuration/annotations/)
- [Helm Chart Values](https://github.com/kubernetes/ingress-nginx/blob/main/charts/ingress-nginx/values.yaml)
- [AWS Load Balancer Controller](https://kubernetes-sigs.github.io/aws-load-balancer-controller/)
- [GKE Ingress Features](https://cloud.google.com/kubernetes-engine/docs/concepts/ingress)
- [AKS Ingress Configuration](https://learn.microsoft.com/en-us/azure/aks/ingress-basic)

