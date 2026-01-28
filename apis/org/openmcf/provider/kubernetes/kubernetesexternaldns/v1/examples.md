# KubernetesExternalDns Examples

Complete YAML manifests for deploying ExternalDNS on different cloud providers.

---

## Example 1: ExternalDNS for GKE with Cloud DNS

Deploy ExternalDNS on Google Kubernetes Engine (GKE) to manage Google Cloud DNS records.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesExternalDns
metadata:
  name: external-dns-gke-prod
spec:
  target_cluster:
    cluster_name: prod-gke-cluster
  namespace:
    value: kubernetes-external-dns
  create_namespace: true  # Creates namespace if it doesn't exist
  gke:
    project_id:
      value: my-gcp-project
    dns_zone_id:
      value: my-cloud-dns-zone-id
```

**How it works:**
- Creates ExternalDNS deployment with Workload Identity
- Scoped to manage records only in the specified Cloud DNS zone
- Automatically creates A/AAAA records for Ingress and LoadBalancer services
- Uses GCP service account `external-dns-gke-prod@my-gcp-project.iam.gserviceaccount.com`

**Prerequisites:**
- GKE cluster with Workload Identity enabled
- Cloud DNS zone created
- GCP service account with `roles/dns.admin` on the zone

---

## Example 2: ExternalDNS for GKE with Foreign Key References

Use foreign keys to reference GCP project and DNS zone from other OpenMCF resources.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesExternalDns
metadata:
  name: external-dns-gke-ref
spec:
  target_cluster:
    cluster_name: prod-gke-cluster
  namespace:
    value: kubernetes-external-dns
  create_namespace: true  # Creates namespace if it doesn't exist
  gke:
    project_id:
      ref: gcp-project-platform
    dns_zone_id:
      ref: dns-zone-prod
```

**Benefits:**
- Project ID and DNS zone ID are pulled from existing OpenMCF resources
- Single source of truth for resource IDs
- Easier to manage across environments

---

## Example 3: ExternalDNS for EKS with Route53

Deploy ExternalDNS on Amazon EKS to manage AWS Route53 DNS records.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesExternalDns
metadata:
  name: external-dns-eks-prod
spec:
  target_cluster:
    cluster_name: prod-eks-cluster
  namespace:
    value: kubernetes-external-dns
  create_namespace: true  # Creates namespace if it doesn't exist
  eks:
    route53_zone_id:
      value: Z1234567890ABC
```

**How it works:**
- Creates ExternalDNS deployment with IRSA (IAM Roles for Service Accounts)
- Auto-creates IAM role with Route53 permissions for the specified zone
- Scoped to manage records only in the specified hosted zone
- No static AWS credentials needed

**Prerequisites:**
- EKS cluster with OIDC provider configured
- Route53 hosted zone created

---

## Example 4: ExternalDNS for EKS with Custom IRSA Role

Provide your own IAM role ARN for IRSA instead of auto-creating.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesExternalDns
metadata:
  name: external-dns-eks-custom
spec:
  target_cluster:
    cluster_name: prod-eks-cluster
  namespace:
    value: kubernetes-external-dns
  create_namespace: true  # Creates namespace if it doesn't exist
  eks:
    route53_zone_id:
      value: Z1234567890ABC
    irsa_role_arn_override: arn:aws:iam::123456789012:role/external-dns-role
```

**Use case:**
- When you manage IAM roles separately (e.g., via Terraform)
- For compliance requirements around IAM role naming
- For cross-account DNS management scenarios

---

## Example 5: ExternalDNS for AKS with Azure DNS

Deploy ExternalDNS on Azure Kubernetes Service (AKS) to manage Azure DNS records.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesExternalDns
metadata:
  name: external-dns-aks-prod
spec:
  target_cluster:
    cluster_name: prod-aks-cluster
  namespace:
    value: kubernetes-external-dns
  create_namespace: true  # Creates namespace if it doesn't exist
  aks:
    dns_zone_id:
      value: /subscriptions/12345678-1234-1234-1234-123456789012/resourceGroups/my-rg/providers/Microsoft.Network/dnszones/example.com
    managed_identity_client_id: 12345678-1234-1234-1234-123456789012
```

**How it works:**
- Uses Azure Workload Identity (federated credentials)
- Managed Identity has "DNS Zone Contributor" role on the DNS zone
- No client secrets or passwords needed

**Prerequisites:**
- AKS cluster with Azure AD Workload Identity enabled
- Azure DNS zone created
- User-assigned Managed Identity with DNS permissions
- Federated credential linking Managed Identity to Kubernetes ServiceAccount

---

## Example 6: ExternalDNS for Cloudflare (Any Kubernetes)

Deploy ExternalDNS to manage Cloudflare DNS records (works with any Kubernetes cluster).

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesExternalDns
metadata:
  name: external-dns-cloudflare
spec:
  target_cluster:
    cluster_name: my-k8s-cluster
  namespace:
    value: kubernetes-external-dns
  create_namespace: true  # Creates namespace if it doesn't exist
  cloudflare:
    api_token: your-cloudflare-api-token-here
    dns_zone_id:
      value: 1234567890abcdef1234567890abcdef
```

**How it works:**
- Creates a Kubernetes Secret with the Cloudflare API token
- Scoped to manage records only in the specified Cloudflare zone
- Works with any Kubernetes cluster (not cloud-specific)

**Prerequisites:**
- Cloudflare account with a zone
- API token with `Zone:Zone:Read` and `Zone:DNS:Edit` permissions

**Security note:** Store `api_token` securely. Consider using external secrets management (e.g., Vault, AWS Secrets Manager) instead of plain text.

---

## Example 7: Cloudflare with Proxy (Orange Cloud)

Enable Cloudflare's proxy/CDN for all managed DNS records.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesExternalDns
metadata:
  name: external-dns-cloudflare-proxy
spec:
  target_cluster:
    cluster_name: my-k8s-cluster
  namespace:
    value: kubernetes-external-dns
  create_namespace: true  # Creates namespace if it doesn't exist
  cloudflare:
    api_token: your-cloudflare-api-token-here
    dns_zone_id:
      value: 1234567890abcdef1234567890abcdef
    is_proxied: true
```

**Effect:**
- All DNS records created by ExternalDNS will have Cloudflare proxy enabled (orange cloud icon)
- Traffic routes through Cloudflare's edge network
- Benefits: DDoS protection, WAF, CDN, SSL/TLS termination
- Use case: Public-facing web services that benefit from Cloudflare's edge features

---

## Example 8: Custom Namespace and Versions

Override default namespace and specify exact ExternalDNS and Helm chart versions.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesExternalDns
metadata:
  name: external-dns-custom
spec:
  target_cluster:
    cluster_name: my-gke-cluster
  namespace:
    value: dns-automation
  create_namespace: true  # Creates namespace if it doesn't exist
  kubernetes_external_dns_version: v0.14.0
  helm_chart_version: 1.14.0
  gke:
    project_id:
      value: my-gcp-project
    dns_zone_id:
      value: my-dns-zone-id
```

**Use case:**
- Pin to specific versions for reproducibility
- Use different namespace for organizational reasons
- Upgrade/downgrade ExternalDNS independently

**Defaults (if not specified):**
- `namespace`: `kubernetes-external-dns`
- `kubernetes_external_dns_version`: `v0.19.0`
- `helm_chart_version`: `1.19.0`

---

## Example 9: Multiple ExternalDNS Instances (Different Zones)

Deploy multiple ExternalDNS instances in the same cluster to manage different DNS zones.

**Instance 1: Production domain**
```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesExternalDns
metadata:
  name: external-dns-prod-domain
spec:
  target_cluster:
    cluster_name: shared-cluster
  namespace:
    value: kubernetes-external-dns
  create_namespace: true  # Creates namespace if it doesn't exist
  gke:
    project_id:
      value: my-gcp-project
    dns_zone_id:
      value: prod-example-com-zone
```

**Instance 2: Staging domain**
```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesExternalDns
metadata:
  name: external-dns-staging-domain
spec:
  target_cluster:
    cluster_name: shared-cluster
  namespace:
    value: kubernetes-external-dns
  create_namespace: true  # Creates namespace if it doesn't exist
  gke:
    project_id:
      value: my-gcp-project
    dns_zone_id:
      value: staging-example-com-zone
```

**How it works:**
- Each ExternalDNS instance has its own Helm release name (matches `metadata.name`)
- Each is scoped to its own DNS zone via zone filtering
- Both can run in the same namespace or different namespaces
- Ingresses can specify which zone to use via annotations

**Use case:**
- Managing multiple domains/subdomains in one cluster
- Isolating production and staging DNS management
- Multi-tenant clusters with separate DNS zones per tenant

---

## Example 10: Foreign Key with Cluster Reference

Reference the target cluster from a OpenMCF Kubernetes cluster resource.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesExternalDns
metadata:
  name: external-dns-ref-cluster
spec:
  target_cluster:
    cluster_name: my-gke-cluster-resource
  namespace:
    value: kubernetes-external-dns
  create_namespace: true  # Creates namespace if it doesn't exist
  gke:
    project_id:
      ref: gcp-project-platform
    dns_zone_id:
      ref: dns-zone-prod
```

**Benefits:**
- All resource IDs are managed centrally in OpenMCF
- Safer than hardcoding values
- Easier to manage across multiple environments

---

## Example 11: Using Existing Namespace

Deploy ExternalDNS into a pre-existing namespace managed separately.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesExternalDns
metadata:
  name: external-dns-existing-ns
spec:
  target_cluster:
    cluster_name: my-gke-cluster
  namespace:
    value: shared-dns-automation  # Must exist before deployment
  create_namespace: false  # Use existing namespace
  gke:
    project_id:
      value: my-gcp-project
    dns_zone_id:
      value: my-dns-zone-id
```

**Use case:**
- Namespace is managed by a separate tool (Terraform, kubectl, etc.)
- Namespace requires special labels, annotations, or resource quotas
- Shared namespace used by multiple components
- GitOps workflow where namespace is created separately

**Important**: The namespace must exist before deploying this resource, or the deployment will fail with an error message indicating the namespace was not found.

---

## Using ExternalDNS After Deployment

Once ExternalDNS is deployed, annotate your Kubernetes resources to create DNS records:

### For Ingress:
```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: my-app
  annotations:
    external-dns.alpha.kubernetes.io/hostname: myapp.example.com
spec:
  rules:
  - host: myapp.example.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: my-app
            port:
              number: 80
```

### For LoadBalancer Service:
```yaml
apiVersion: v1
kind: Service
metadata:
  name: my-app
  annotations:
    external-dns.alpha.kubernetes.io/hostname: myapp.example.com
spec:
  type: LoadBalancer
  ports:
  - port: 80
    targetPort: 8080
  selector:
    app: my-app
```

ExternalDNS will automatically:
1. Detect the annotation
2. Get the LoadBalancer IP or Ingress endpoint
3. Create/update an A record in your DNS provider
4. Keep it in sync as the service changes

---

## Troubleshooting

**ExternalDNS not creating records?**
- Check ExternalDNS pod logs: `kubectl logs -n kubernetes-external-dns <pod-name>`
- Verify zone filtering matches your DNS zone
- Ensure cloud IAM permissions are correct (IRSA, Workload Identity, Managed Identity)

**Authentication errors?**
- **GKE**: Verify Workload Identity binding and GSA has `roles/dns.admin`
- **EKS**: Verify IRSA role has Route53 permissions and trust policy
- **AKS**: Verify Managed Identity has "DNS Zone Contributor" role and federated credential exists
- **Cloudflare**: Verify API token has correct permissions and is not expired

**DNS records not updating?**
- ExternalDNS uses `upsert-only` policy by default (doesn't delete)
- Check TTL - DNS changes may take time to propagate
- Verify the hostname annotation matches the zone ExternalDNS is managing

---

## Next Steps

1. Choose the example that matches your cloud provider
2. Customize with your cluster ID, project/zone IDs
3. Deploy via `openmcf deploy <manifest-file>`
4. Verify ExternalDNS pod is running
5. Test by creating an Ingress with hostname annotation
6. Check DNS records were created in your DNS provider

For more details, see the [README](README.md) and [research documentation](docs/README.md).

