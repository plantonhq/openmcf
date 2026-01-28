# Kubernetes Cert-Manager Examples

This document provides complete, copy-paste ready examples for deploying kubernetes-cert-manager with different DNS providers.

---

## Example 1: Minimal Cloudflare Setup

Simple deployment with Cloudflare DNS for a single domain.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesCertManager
metadata:
  name: cert-manager
spec:
  targetCluster:
    clusterName: "my-gke-cluster"
  namespace:
    value: "kubernetes-cert-manager"
  createNamespace: true
  
  acme:
    email: "admin@example.com"
    server: "https://acme-v02.api.letsencrypt.org/directory"
  
  dnsProviders:
    - name: cloudflare-prod
      dnsZones:
        - example.com
      cloudflare:
        apiToken: "your-cloudflare-api-token"
```

**Use Case:** Basic production setup for a single domain managed in Cloudflare.

---

## Example 2: Cloudflare with Multiple Domains

Managing multiple domains with a single Cloudflare account.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesCertManager
metadata:
  name: cert-manager-multi
spec:
  targetCluster:
    clusterName: "my-gke-cluster"
  namespace:
    value: "kubernetes-cert-manager"
  createNamespace: true
  
  acme:
    email: "certs@acme-corp.com"
  
  dnsProviders:
    - name: cloudflare-primary
      dnsZones:
        - example.com
        - example.org
        - example.net
      cloudflare:
        apiToken: "cloudflare-api-token-here"
```

**Use Case:** Organization with multiple domains all managed in Cloudflare.

---

## Example 3: GCP Cloud DNS

Using Google Cloud DNS with Workload Identity (GKE).

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesCertManager
metadata:
  name: cert-manager-gcp
spec:
  targetCluster:
    clusterName: "my-gke-cluster"
  namespace:
    value: "kubernetes-cert-manager"
  createNamespace: true
  
  acme:
    email: "platform@example.com"
  
  dnsProviders:
    - name: gcp-internal
      dnsZones:
        - internal.example.net
      gcpCloudDns:
        projectId: "my-gcp-project"
        serviceAccountEmail: "cert-manager@my-project.iam.gserviceaccount.com"
```

**Use Case:** Internal domains managed in Google Cloud DNS, running on GKE with Workload Identity.

**Prerequisites:**
- GCP service account with `dns.admin` role
- Workload Identity binding configured

---

## Example 4: AWS Route53

Using AWS Route53 with IRSA (IAM Roles for Service Accounts).

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesCertManager
metadata:
  name: cert-manager-aws
spec:
  targetCluster:
    clusterName: "my-eks-cluster"
  namespace:
    value: "kubernetes-cert-manager"
  createNamespace: true
  
  acme:
    email: "devops@example.com"
  
  dnsProviders:
    - name: aws-route53
      dnsZones:
        - aws.example.com
        - api.example.com
      awsRoute53:
        region: "us-east-1"
        roleArn: "arn:aws:iam::123456789012:role/cert-manager-dns-role"
```

**Use Case:** Domains managed in AWS Route53, running on EKS with IRSA.

**Prerequisites:**
- IAM role with Route53 permissions
- IRSA trust relationship configured

---

## Example 5: Azure DNS

Using Azure DNS with Managed Identity.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesCertManager
metadata:
  name: cert-manager-azure
spec:
  targetCluster:
    clusterName: "my-aks-cluster"
  namespace:
    value: "kubernetes-cert-manager"
  createNamespace: true
  
  acme:
    email: "it@example.com"
  
  dnsProviders:
    - name: azure-dns
      dnsZones:
        - azure.example.com
      azureDns:
        subscriptionId: "12345678-1234-1234-1234-123456789012"
        resourceGroup: "dns-resources"
        clientId: "87654321-4321-4321-4321-210987654321"
```

**Use Case:** Domains managed in Azure DNS, running on AKS with Managed Identity.

**Prerequisites:**
- Managed Identity with DNS Zone Contributor role
- Workload identity federation configured for AKS

---

## Example 6: Multi-Provider Setup

Combining multiple DNS providers in a single deployment.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesCertManager
metadata:
  name: cert-manager-hybrid
spec:
  targetCluster:
    clusterName: "my-gke-cluster"
  namespace:
    value: "kubernetes-cert-manager"
  createNamespace: true
  
  acme:
    email: "certificates@multi-cloud.com"
  
  dnsProviders:
    # Cloudflare for public domains
    - name: cloudflare-public
      dnsZones:
        - example.com
        - example.org
      cloudflare:
        apiToken: "cloudflare-token"
    
    # GCP Cloud DNS for internal domains
    - name: gcp-internal
      dnsZones:
        - internal.example.net
      gcpCloudDns:
        projectId: "gcp-project-123"
        serviceAccountEmail: "cert-manager@gcp-project-123.iam.gserviceaccount.com"
    
    # AWS Route53 for AWS-specific domains
    - name: aws-services
      dnsZones:
        - aws.example.com
      awsRoute53:
        region: "us-west-2"
        roleArn: "arn:aws:iam::123456789012:role/cert-manager"
```

**Use Case:** Multi-cloud organization with domains split across different DNS providers.

**What Gets Created:**
- kubernetes-cert-manager deployed once
- ServiceAccount with annotations for GCP and AWS workload identity
- 3 ClusterIssuers (one per domain): `example.com`, `example.org`, `internal.example.net`, `aws.example.com`
- Kubernetes Secret for Cloudflare credentials

---

## Example 7: Staging Environment

Using Let's Encrypt staging for testing.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesCertManager
metadata:
  name: cert-manager-staging
spec:
  targetCluster:
    clusterName: "my-staging-cluster"
  namespace:
    value: "cert-manager-staging"
  createNamespace: true
  
  acme:
    email: "staging@example.com"
    # Use staging server for testing
    server: "https://acme-staging-v02.api.letsencrypt.org/directory"
  
  dnsProviders:
    - name: cloudflare-test
      dnsZones:
        - staging.example.com
      cloudflare:
        apiToken: "cloudflare-staging-token"
```

**Use Case:** Testing certificate issuance without hitting Let's Encrypt rate limits.

**Important:** Staging certificates are not trusted by browsers. Use for testing only.

---

## Example 8: Custom Namespace and Version

Deploying to a custom namespace with specific version.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesCertManager
metadata:
  name: cert-manager-custom
spec:
  targetCluster:
    clusterName: "my-production-cluster"
  namespace:
    value: "security-cert-manager"
  createNamespace: false  # Use existing namespace with custom security policies
  kubernetesCertManagerVersion: "v1.19.1"
  helmChartVersion: "v1.19.1"
  
  acme:
    email: "security@example.com"
  
  dnsProviders:
    - name: cloudflare-security
      dnsZones:
        - secure.example.com
      cloudflare:
        apiToken: "cloudflare-api-token"
```

**Use Case:** Organizations with specific namespace requirements or version pinning.

---

## Requesting a Certificate

After deploying kubernetes-cert-manager, request certificates using Certificate resources:

```yaml
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: example-com-tls
  namespace: default
spec:
  secretName: example-com-tls
  issuerRef:
    name: example.com  # References the ClusterIssuer created by the addon
    kind: ClusterIssuer
  dnsNames:
    - example.com
    - www.example.com
  # For wildcard certificates:
  # dnsNames:
  #   - "*.example.com"
```

---

## Verification

After deployment, verify kubernetes-cert-manager is working:

```bash
# Check cert-manager pods
kubectl get pods -n kubernetes-cert-manager

# List ClusterIssuers
kubectl get clusterissuers

# Check a specific ClusterIssuer status
kubectl describe clusterissuer example.com

# Request a test certificate
kubectl apply -f test-certificate.yaml

# Watch certificate creation
kubectl get certificate -n default -w

# Check certificate details
kubectl describe certificate example-com-tls -n default
```

---

## Troubleshooting

### Check ClusterIssuer Status

```bash
kubectl describe clusterissuer <domain-name>
```

Look for `Ready` condition in the status.

### Check Cert-Manager Logs

```bash
kubectl logs -n kubernetes-cert-manager -l app=cert-manager --tail=100 -f
```

### Verify DNS Propagation

```bash
# Check if DNS TXT record was created
dig _acme-challenge.example.com TXT

# Use DNS checker
dig @1.1.1.1 _acme-challenge.example.com TXT
```

### Common Issues

1. **Challenge Failed:** Check DNS provider credentials and permissions
2. **Timeout:** DNS propagation can take 60-120 seconds
3. **Rate Limit:** Use staging server for testing
4. **Wrong Secret:** Verify Cloudflare secret name matches provider name

---

## Next Steps

- Configure ingress to use TLS certificates
- Set up certificate renewal monitoring
- Implement certificate backup strategy
- Configure alerts for certificate expiry

