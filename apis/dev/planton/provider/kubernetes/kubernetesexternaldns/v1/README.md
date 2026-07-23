# Kubernetes External DNS

## When NOT to Use This

**If you manage a handful of DNS records by hand, you do not need a controller.** ExternalDNS earns its keep when records must FOLLOW workloads — Services, Ingresses, and Gateway API routes appearing, moving, and disappearing faster than humans update zones. For static records that change a few times a year, manage them directly on the zone resource (AwsRoute53Zone, GcpDnsZone, AzureDnsZone, CloudflareDnsZone). And one installation manages ONE DNS provider: clusters publishing to several providers (or with different ownership boundaries) deploy multiple instances of this component, not one instance with more configuration.

## Overview

**KubernetesExternalDns** installs ExternalDNS — the controller that watches Kubernetes resources (Services, Ingresses, Gateway API routes, ...) and publishes matching DNS records into a DNS provider — from the official Helm chart (`external-dns` at `https://kubernetes-sigs.github.io/external-dns/`). The typed spec covers the chart's meaningful configuration surface — provider selection, sources, sync policy, ownership registry, filtering, tuning — with a `helm_values` escape hatch (merged last, Helm `-f` semantics, identical on both engines) for anything beyond it.

**Key design points:**

- **One installation, one provider**: `dns_provider` is a oneof — exactly one of `aws_route53`, `google_cloud_dns`, `azure_dns`, `cloudflare`, `webhook`, or `in_memory`. Multi-provider clusters run multiple instances; the release is named after `metadata.name` so instances coexist naturally
- **Ownership is first-class**: the TXT registry (`registry`, `txt_owner_id`, `txt_prefix`/`txt_suffix`) is what lets several instances — and humans — share a zone without deleting each other's records. Every instance sharing a zone needs a distinct `txt_owner_id`
- **The sync policy is a safety dial**: `upsert-only` (default) never deletes; `sync` fully reconciles the records this instance owns — the right choice for dedicated zones; `create-only` never touches a record twice
- **Keyless authentication is the production posture**: `workload_identity` binds the controller ServiceAccount to a cloud identity (EKS IRSA, GKE Workload Identity, AKS Workload Identity); static credentials are the fallback and materialize as Kubernetes Secrets, never as chart values
- **Out-of-tree providers ride the `webhook` arm**: upstream's extension architecture runs the provider's webhook image as a sidecar next to the controller — every provider that is not in-tree (Akamai, OVH, Scaleway, RFC2136, ...) works this way

## Environment Injection (where cloud identity flows in)

Kubernetes is a dual provider: the cluster runs IN one environment while the DNS zone often lives in ANOTHER. `dns_provider` selects WHERE records are written; `workload_identity` plus per-arm static credentials select HOW the controller authenticates. Any combination is expressible:

| Host cluster | Keyless (`workload_identity`) | Static credentials (per provider arm) | Ambient (node identity) |
|---|---|---|---|
| EKS | `workload_identity.eks.role_arn` → Route 53 (IRSA) | any arm: AWS keys, GCP key JSON, Azure SP, Cloudflare token | node role → Route 53 |
| GKE | `workload_identity.gke.service_account_email` → Cloud DNS | any arm | node service account → Cloud DNS |
| AKS | `workload_identity.aks.client_id` → Azure DNS (the module also stamps the required pod label) | any arm | kubelet / user-assigned managed identity (`azure_dns.managed_identity_client_id`) |
| Self-managed / kind / datacenter | — (no cloud federation) | required for any cloud DNS: static keys, SA key, SP, or token | — |

The cross-cloud combinations this table implies are first-class, not workarounds:

- **EKS cluster + Cloudflare DNS**: no workload identity at all — the Cloudflare token rides the `cloudflare` arm and materializes as a Kubernetes Secret (`CF_API_TOKEN`)
- **GKE cluster + Route 53**: GKE has no AWS federation, so the `aws_route53` arm carries static keys (also a Secret) — or `assume_role` on top for cross-account zones
- **Any cluster + any webhook provider**: the `webhook` arm's sidecar carries its own provider-specific configuration

The controller ServiceAccount name is pinned to `metadata.name` and exported (`status.outputs.service_account_name`) precisely so the cloud-side half of a keyless binding (IAM trust policy, Workload Identity binding, federated credential) can be composed in the same infra chart.

## Essential Configuration Fields

### Required

- **`spec.namespace`**: installation namespace (`external-dns` by convention) — literal or a KubernetesNamespace reference
- **`spec.dns_provider`** (oneof): exactly one provider arm

### Common

- **`spec.create_namespace`**: create (and own) the namespace with the release
- **`spec.chart_version`**: pinned chart version (defaults to the validated pin; chart releases are cut separately from the controller — chart 1.21.x ships controller v0.21.x)
- **`spec.sources`**: what the controller watches (chart default `service` + `ingress`; add Gateway API route sources or `crd` for DNSEndpoint objects)
- **`spec.policy`**: `upsert-only` (safe default) / `sync` (dedicated zones) / `create-only`
- **`spec.txt_owner_id`**: distinct per instance sharing a zone — what stops one instance from deleting another's records under `sync`
- **`spec.domain_filters`** + per-arm `zone_id_filters`: the guardrails against touching unrelated zones in a shared account
- **`spec.interval`** / `spec.trigger_loop_on_event`: reconcile cadence vs. reaction to source events
- **`spec.prometheus.service_monitor`**: opt-in ServiceMonitor (requires the Prometheus operator CRDs — the release fails without them)
- **`spec.helm_values`**: escape hatch for chart values beyond the typed fields — never the primary interface

## Stack Outputs

| Output | Purpose |
|---|---|
| `namespace` | Installation namespace |
| `release_name` | Helm release name (equals `metadata.name` — instances coexist) |
| `service_account_name` | Controller ServiceAccount — the identity to bind cloud-side for keyless provider access |

## Composing in Infra Charts

The standard chart wiring: the zone and this component deploy in one run, with the zone ID flowing into `zone_id_filters` as a reference (`AwsRoute53Zone` → `status.outputs.zone_id`, likewise for GCP/Azure/Cloudflare zone kinds) and the IAM role/GCP SA/managed identity flowing into `workload_identity`. A cluster publishing public records to Route 53 and internal records to Cloud DNS runs two instances of this component — each with its own provider arm, `txt_owner_id`, and filters. ExternalDNS pairs naturally with KubernetesCertManager: external-dns publishes the names, cert-manager proves ownership of them for certificates.

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesExternalDns
metadata:
  name: external-dns-route53
spec:
  namespace:
    value: external-dns
  createNamespace: true
  awsRoute53:
    region: us-east-1
    zoneIdFilters:
      - valueFrom:
          kind: AwsRoute53Zone
          name: my-zone
          fieldPath: status.outputs.zone_id
  workloadIdentity:
    eks:
      roleArn:
        valueFrom:
          kind: AwsIamRole
          name: external-dns-role
          fieldPath: status.outputs.role_arn
  policy: sync
  txtOwnerId: my-cluster
```
||||||| 178dc56b9

This creates an ExternalDNS deployment that:
- Uses Workload Identity to authenticate with Google Cloud DNS
- Is scoped to manage records only in the specified DNS zone
- Watches for Kubernetes Ingresses and Services with DNS annotations
- Automatically creates A/CNAME records pointing to LoadBalancer IPs

## Benefits

### For Platform Teams
- **Standardized deployment**: Same interface across all clusters and clouds
- **Security compliance**: No static credentials, all authentication via cloud-native IAM
- **Audit trail**: All DNS changes are logged through cloud provider audit logs
- **Multi-tenancy**: Deploy multiple ExternalDNS instances for different zones/domains

### For Application Teams  
- **Zero DNS overhead**: Just annotate your Ingress/Service, DNS is automatic
- **Faster deployments**: No waiting for manual DNS changes
- **Self-service**: No dependency on platform team for DNS records
- **GitOps-friendly**: DNS is declared in the same YAML as your service

### For Everyone
- **Fewer incidents**: No more "forgot to update DNS" outages
- **Faster iteration**: Deploy new services without DNS friction
- **Cost savings**: Reduce operational toil and manual processes
- **Better reliability**: DNS stays in sync with actual service state

## How It Works

1. **Deploy**: Apply a KubernetesExternalDns manifest to your cluster
2. **Watch**: ExternalDNS monitors Ingress/Service resources for annotations like `external-dns.alpha.kubernetes.io/hostname`
3. **Sync**: When it detects a hostname annotation, ExternalDNS creates/updates the DNS record in your cloud provider
4. **Maintain**: As services scale, move, or are deleted, DNS records are kept in sync automatically

## Use Cases

- **Public-facing services**: Automatically create DNS for LoadBalancer services
- **Ingress hostnames**: Auto-manage A/CNAME records for Ingress resources
- **Multi-region active-active**: Keep DNS updated as traffic shifts between regions
- **Blue-green deployments**: DNS points to new stack automatically when Ingress updates
- **Development environments**: Each dev cluster gets its own subdomain, managed automatically

## Documentation

- **[Research Documentation](docs/README.md)**: Deep dive into ExternalDNS landscape, authentication patterns, and best practices
- **[Examples](examples.md)**: Complete YAML manifests for GKE, EKS, AKS, and Cloudflare
- **[Pulumi Module](iac/pulumi/README.md)**: Using the Pulumi implementation directly
- **[Terraform Module](iac/tf/README.md)**: Using the Terraform implementation directly

## Configuration Reference

| Field | Description | Default |
|-------|-------------|---------|
| `namespace` | Kubernetes namespace for ExternalDNS | `kubernetes-external-dns` |
| `create_namespace` | Whether to create the namespace (if false, namespace must exist) | `false` |
| `kubernetes_external_dns_version` | ExternalDNS image version | `v0.19.0` |
| `helm_chart_version` | Helm chart version | `1.19.0` |
| `provider_config` | Cloud provider configuration (gke, eks, aks, or cloudflare) | Required |

## Namespace Management

The component provides flexible namespace management through the `create_namespace` flag:

### Automatic Namespace Creation

When `create_namespace` is set to `true`, the module creates the namespace:

```yaml
spec:
  namespace:
    value: kubernetes-external-dns
  create_namespace: true  # Module will create the namespace
```

### Using Existing Namespace

When `create_namespace` is `false` (default), the namespace must exist before deployment:

```yaml
spec:
  namespace:
    value: existing-namespace
  create_namespace: false  # Namespace must already exist
```

**Important**: If `create_namespace` is `false` and the namespace doesn't exist, deployment will fail with a clear error message.

### GKE Configuration

| Field | Description |
|-------|-------------|
| `project_id` | GCP project ID hosting the DNS zone |
| `dns_zone_id` | Cloud DNS zone ID to manage |

### EKS Configuration

| Field | Description |
|-------|-------------|
| `route53_zone_id` | Route53 hosted zone ID to manage |
| `irsa_role_arn_override` | Optional IAM role ARN for IRSA (auto-created if empty) |

### AKS Configuration

| Field | Description |
|-------|-------------|
| `dns_zone_id` | Azure DNS zone ID |
| `managed_identity_client_id` | Azure Managed Identity client ID |

### Cloudflare Configuration

| Field | Description |
|-------|-------------|
| `api_token` | Cloudflare API token with Zone:DNS:Edit permissions |
| `dns_zone_id` | Cloudflare zone ID to manage |
| `is_proxied` | Enable Cloudflare proxy (orange cloud) for records |

## Prerequisites

### For GKE
- GKE cluster with Workload Identity enabled
- Google Cloud DNS zone created
- GCP service account with `roles/dns.admin` on the DNS zone
- IAM binding for Workload Identity (created by module)

### For EKS
- EKS cluster with OIDC provider configured
- Route53 hosted zone created
- IAM role with Route53 permissions (created by module or provided)

### For AKS
- AKS cluster with Azure AD Workload Identity enabled
- Azure DNS zone created
- User-assigned Managed Identity with "DNS Zone Contributor" role

### For Cloudflare
- Cloudflare account and zone
- API token with Zone:Zone:Read and Zone:DNS:Edit permissions

## Next Steps

1. Review the [examples](examples.md) for your cloud provider
2. Read the [research documentation](docs/README.md) to understand architecture patterns
3. Create a KubernetesExternalDns manifest
4. Deploy it via `planton deploy`
5. Test by creating a Service or Ingress with DNS annotation

## Support

For issues, questions, or contributions, see the main [Planton repository](https://github.com/plantonhq/planton).

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
