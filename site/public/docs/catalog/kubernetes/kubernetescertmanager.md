---
title: "Cert Manager"
description: "Cert Manager deployment documentation"
icon: "package"
order: 100
componentName: "kubernetescertmanager"
---

# Deploying kubernetes-cert-manager on Kubernetes: From Manual Manifests to Production Automation

## The TLS Certificate Challenge

For years, the conventional wisdom was simple: managing TLS certificates in Kubernetes meant either manually uploading them as secrets or building custom automation. Then Let's Encrypt democratized certificate issuance with its free, automated ACME protocol, and kubernetes-cert-manager emerged as the Kubernetes-native way to bridge these worlds.

But here's what changed: **deploying kubernetes-cert-manager correctly** has become just as important as having it installed. A kubernetes-cert-manager deployment without proper DNS provider integration, high availability, or secure credential management isn't production-ready—it's a time bomb waiting for certificates to expire during an outage.

This document explains how deployment methods for kubernetes-cert-manager evolved from simple `kubectl apply` commands to production-grade, multi-cloud automation patterns. More importantly, it explains **why** OpenMCF chose to abstract away the complexity of ClusterIssuer creation and DNS provider integration while maintaining the flexibility production environments demand.

## The Deployment Maturity Spectrum

### Level 0: Static Manifests (Quick Start, Not Production)

The simplest path is applying kubernetes-cert-manager's official YAML bundle directly:

```bash
kubectl apply -f https://github.com/kubernetes-cert-manager/kubernetes-cert-manager/releases/download/v1.19.1/kubernetes-cert-manager.yaml
```

This works. It installs kubernetes-cert-manager with all necessary CRDs, RBAC, and components in one command. But it's the deployment equivalent of hardcoding configuration: **what you gain in simplicity, you lose in maintainability**.

**The reality**: Static manifests become technical debt. You're responsible for:
- Manually tracking upstream releases and applying upgrades sequentially
- Ensuring CRDs are updated *before* kubernetes-cert-manager components (a common failure mode)
- Managing custom configurations by forking and editing the manifest (version control nightmare)
- Coordinating upgrades across multiple clusters without tooling

**The verdict**: Fine for experimentation and local development. Avoid in production unless you have strong operational discipline and very few clusters.

### Level 1: Helm Chart (Production Starting Point)

Using the official Helm chart is what the kubernetes-cert-manager project calls a "first-class" installation method:

```bash
helm repo add jetstack https://charts.jetstack.io
helm install kubernetes-cert-manager jetstack/kubernetes-cert-manager \
  --namespace kubernetes-cert-manager --create-namespace \
  --set crds.enabled=true
```

Helm brings structure: versioned releases, values-based configuration, atomic upgrades, and rollback capability. The official chart (published as OCI at `quay.io/jetstack/charts/kubernetes-cert-manager`) is well-maintained and handles CRD lifecycle properly when configured correctly.

**What you gain**:
- Customizable deployment via values (replicas, resource limits, feature flags)
- Upgrade automation with version constraints
- Ability to diff changes before applying
- Package management consistency across clusters

**What you still manage**:
- Creating and maintaining ClusterIssuer resources separately
- DNS provider credential management (secrets, workload identity bindings)
- High availability configuration (replica counts, PodDisruptionBudgets, anti-affinity)
- DNS-01 solver configuration for each provider

**The reality**: This is the minimum bar for production. Most teams stop here and layer additional tooling on top for credential injection and ClusterIssuer templating.

**The verdict**: Production-viable with operational overhead. The Helm chart solves deployment, but not the integration complexity.

### Level 2: GitOps with Helm (Continuous Deployment)

Storing kubernetes-cert-manager's Helm release configuration in Git and using Argo CD or Flux to apply it brings declarative, auditable deployments:

```yaml
# ArgoCD Application
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: kubernetes-cert-manager
spec:
  source:
    chart: kubernetes-cert-manager
    repoURL: https://charts.jetstack.io
    targetRevision: v1.19.1
    helm:
      values: |
        crds:
          enabled: true
        replicaCount: 2
```

GitOps operators handle the reconciliation loop: drift detection, automated updates, and multi-cluster synchronization.

**The critical pitfall**: **CRD management in GitOps**. Kubernetes applies admission webhooks before resources are created. During upgrades, kubernetes-cert-manager's webhook can reject new CRD schemas if it's running an older version. The solution is carefully orchestrated sync waves or temporarily disabling webhook validation—an operational hazard most teams discover the hard way.

**What you gain**:
- Declarative infrastructure with Git as the source of truth
- Automated drift correction
- Multi-cluster deployment consistency
- Audit trail of all changes

**What you still manage**:
- CRD upgrade orchestration (sync waves, health checks)
- Secret management (often via external-secrets operator or sealed secrets)
- ClusterIssuer resources (separate GitOps manifests)
- Version compatibility across the toolchain (Helm, Argo/Flux, kubernetes-cert-manager)

**The verdict**: The production standard for multi-cluster environments, but requires investment in GitOps tooling and expertise.

### Level 3: Integrated Automation (Production-Grade)

This is where deployment meets integration. Instead of deploying kubernetes-cert-manager as a standalone component and then manually wiring up ClusterIssuers, DNS credentials, and cloud provider integrations, **automation handles the entire lifecycle**.

The pattern: a higher-level API that accepts intent ("I want kubernetes-cert-manager with Cloudflare DNS-01 for these domains") and produces a complete, working system—deployment, credentials, ClusterIssuers, and all.

**Characteristics**:
- Single source of configuration capturing both deployment and integration
- Automatic ClusterIssuer creation based on DNS provider configuration
- Secure credential management (cloud workload identity or secret injection)
- High availability by default (multiple replicas, PodDisruptionBudgets)
- Version enforcement (minimum versions for stability)
- Best practices encoded (recursive nameservers, DNS propagation settings)

**What you specify**:
- ACME email and server (staging vs. production)
- DNS providers with their authentication methods
- Target cluster and namespace

**What you don't manage**:
- ClusterIssuer YAML files
- Kubernetes Secret creation for DNS provider credentials
- ServiceAccount annotations for workload identity
- Solver configuration and DNS zone mapping
- Helm values templating for high availability

**The reality**: This is the automation layer most enterprises eventually build in-house. Some use Terraform modules, others build Kubernetes operators, many use internal platforms. The goal is consistent: **abstract away the undifferentiated heavy lifting** while preserving flexibility.

**The verdict**: The production endgame. Infrastructure teams own the automation. Application teams consume the abstraction.

## The DNS Provider Integration Challenge

kubernetes-cert-manager's power lies in DNS-01 ACME challenges: prove you control a domain by creating a specific TXT record. But that proof requires giving kubernetes-cert-manager permission to create DNS records, and each cloud provider has different authentication patterns.

### Authentication Evolution: From Static Keys to Ambient Identity

**The old way (still supported)**:
- Generate API keys or service account JSON files
- Store them in Kubernetes Secrets
- Reference the secret in ClusterIssuer configuration
- Rotate credentials manually (rarely happens)

**The new way (production standard)**:
- **GCP**: Workload Identity—link Kubernetes ServiceAccount to Google Cloud Service Account
- **AWS**: IRSA (IAM Roles for Service Accounts)—OIDC-based IAM role assumption
- **Azure**: Workload Identity Federation—link to Managed Identity
- **Cloudflare**: Still requires API tokens, but scoped to specific zones with minimal permissions

This shift from static credentials to ambient identity **eliminates secret sprawl** and reduces the blast radius of compromised credentials. The kubernetes-cert-manager pod authenticates using its Kubernetes ServiceAccount, which cloud metadata services map to cloud-provider IAM roles.

### Multi-Provider Patterns

Production environments rarely live in a single cloud or use a single DNS provider. Common patterns:

**Pattern 1: Multi-cloud domains**
- Public domains managed in Cloudflare (centralized, multi-cloud)
- Internal domains per cloud (GCP Cloud DNS, AWS Route53, Azure DNS)
- One ClusterIssuer with multiple solvers, each solver bound to specific DNS zones

**Pattern 2: Environment segregation**
- Staging uses Cloudflare with scoped API tokens
- Production uses cloud-native DNS with workload identity
- Separate ClusterIssuers per environment

**Pattern 3: Cross-account DNS**
- Kubernetes cluster in Account A
- DNS zones in Account B (security boundary)
- Cross-account IAM role assumption configured in ClusterIssuer

The key insight: **kubernetes-cert-manager supports solver multiplexing**. A single ClusterIssuer can have multiple DNS-01 solvers, each with different authentication and DNS zones. When a Certificate is requested, kubernetes-cert-manager selects the appropriate solver based on domain matching.

### The Configuration Complexity Tax

Here's what a production-grade multi-provider ClusterIssuer looks like:

```yaml
apiVersion: kubernetes-cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: multi-cloud-issuer
spec:
  acme:
    server: https://acme-v02.api.letsencrypt.org/directory
    email: security@example.com
    privateKeySecretRef:
      name: acme-account-key
    solvers:
    # Cloudflare solver for public domains
    - selector:
        dnsZones:
          - example.com
          - example.org
      dns01:
        cloudflare:
          apiTokenSecretRef:
            name: cloudflare-api-token
            key: api-token
    
    # GCP Cloud DNS solver with Workload Identity
    - selector:
        dnsZones:
          - internal.example.net
      dns01:
        cloudDNS:
          project: my-gcp-project
          # No serviceAccountSecretRef = use ambient credentials
    
    # AWS Route53 solver with IRSA
    - selector:
        dnsZones:
          - aws.example.com
      dns01:
        route53:
          region: us-east-1
          # No accessKeyID/secret = use IRSA
```

Plus the corresponding Secrets, ServiceAccount annotations, IAM policies, and workload identity bindings. **That's a lot of YAML to maintain across environments and clusters.**

## What OpenMCF Automates

OpenMCF chose kubernetes-cert-manager as the certificate automation engine—there's no viable alternative with kubernetes-cert-manager's maturity, community, and CNCF backing. But deploying kubernetes-cert-manager is table stakes. **The differentiation is integration automation.**

### The 80/20 Principle in Action

Most production kubernetes-cert-manager deployments configure:
- **Version**: Usually latest stable (v1.19.1 as of this writing, minimum v1.16.4 for Cloudflare compatibility)
- **Namespace**: `kubernetes-cert-manager` (convention)
- **ACME server**: Let's Encrypt production or staging
- **ACME email**: Operations team email
- **DNS providers**: 1-3 providers (Cloudflare, GCP Cloud DNS, AWS Route53, Azure DNS)
- **High availability**: 2 controller replicas, 3 webhook replicas (best practice)

What most teams **don't** customize:
- CRD installation (always enabled)
- Webhook configuration (defaults are secure)
- CA injector settings (auto-configured)
- Custom ACME solver configurations beyond DNS-01
- Prometheus metrics (enabled by default)

**The OpenMCF API design**: Capture the 20% of configuration that 80% of users need. Default everything else to production best practices.

### Automatic ClusterIssuer Creation

Instead of requiring users to write and maintain ClusterIssuer YAML, OpenMCF generates it from DNS provider configuration. You specify:

```yaml
acme:
  email: "admin@example.com"
  server: "https://acme-v02.api.letsencrypt.org/directory"

dnsProviders:
  - name: cloudflare-prod
    dnsZones:
      - example.com
      - example.org
    cloudflare:
      apiToken: "your-token"
```

OpenMCF generates:
- Kubernetes Secret with the Cloudflare API token
- ClusterIssuer with DNS-01 solver configured for those zones
- ServiceAccount with appropriate cloud provider annotations (for GCP/AWS/Azure)

**The benefit**: One source of truth. Updating DNS providers updates the ClusterIssuer automatically. No manual YAML maintenance.

### Credential Management Abstraction

For **Cloudflare**: API tokens are stored as Kubernetes Secrets (no avoiding this). OpenMCF creates the secret in the correct namespace with the correct key name that kubernetes-cert-manager expects.

For **cloud providers** (GCP/AWS/Azure): OpenMCF configures the kubernetes-cert-manager ServiceAccount with the appropriate workload identity annotations. The cloud provider's metadata service handles authentication. No static credentials stored in Kubernetes.

This pattern follows the **principle of least surprise**: credentials are handled the way each provider expects, abstracting away the provider-specific details.

### High Availability by Default

Production kubernetes-cert-manager requires:
- **Controller**: 2 replicas (leader election, hot standby)
- **Webhook**: 3 replicas (no leader election, needs quorum for availability)
- **CA Injector**: 2 replicas (leader election)
- **PodDisruptionBudgets**: Ensure at least one replica survives node drains

OpenMCF's Helm configuration sets these defaults. You don't declare them; they're encoded as operational best practices.

### DNS Propagation Reliability

A common kubernetes-cert-manager failure mode: **DNS propagation timeouts**. kubernetes-cert-manager creates a TXT record via the DNS provider API, but then can't verify it exists because it's querying the wrong nameservers (often in-cluster CoreDNS with stale caches).

The fix: Configure kubernetes-cert-manager to use **public recursive nameservers** (Cloudflare 1.1.1.1 and Google 8.8.8.8) for DNS propagation checks. This is a flag in the kubernetes-cert-manager deployment:

```yaml
--dns01-recursive-nameservers=1.1.1.1:53,8.8.8.8:53
--dns01-recursive-nameservers-only=true
```

OpenMCF sets this automatically. Users don't need to know about DNS propagation nuances.

### Version Enforcement

Not all kubernetes-cert-manager versions are created equal. **Version v1.16.4** fixed critical Cloudflare API compatibility issues where DNS records weren't properly cleaned up after challenges, causing subsequent challenges to fail.

OpenMCF enforces a minimum version. You can specify newer versions, but you can't accidentally deploy a broken older one.

## Production Patterns Worth Knowing

Even with automation, understanding how kubernetes-cert-manager works in production prevents surprises.

### ClusterIssuer vs. Issuer: Trust Boundaries

- **ClusterIssuer**: Cluster-scoped, usable by Certificates in any namespace. Default for most deployments.
- **Issuer**: Namespace-scoped, only usable by Certificates in the same namespace.

Use **ClusterIssuer** when: You operate a trusted, single-tenant cluster and want a single point of configuration.

Use **Issuer** when: You need namespace-level trust isolation (multi-tenant clusters, per-team credentials).

OpenMCF defaults to ClusterIssuer because most production clusters are single-tenant or have environment-level trust boundaries, not namespace-level.

### Wildcard vs. Individual Certificates

**Wildcard certificates** (`*.example.com`) simplify management for many subdomains but:
- **Require DNS-01 challenges** (HTTP-01 doesn't support wildcards)
- Don't cover the apex domain (`example.com` needs a separate cert)
- If compromised, impact is wider (all subdomains)

**Individual certificates** offer:
- Least privilege (each cert covers only specific hosts)
- Can use HTTP-01 challenges (easier in some cases)
- Better blast radius containment

**The production pattern**: Wildcard for dynamic environments (e.g., `*.dev.example.com`) where subdomains are ephemeral. Individual certificates for production services with stable hostnames.

### Certificate Lifecycle and Auto-Renewal

kubernetes-cert-manager automatically renews certificates **30 days before expiration** (configurable via `renewBefore`). The renewal process:
1. Creates a new ACME order
2. Completes DNS-01 challenge (same flow as initial issuance)
3. Receives new certificate from Let's Encrypt
4. Updates the Kubernetes Secret with new certificate and (by default) a new private key

**Critical insight**: Applications must reload certificates when the Secret changes. Ingress controllers (nginx-ingress, Traefik) handle this automatically. Custom applications need to watch the Secret and reload, or use kubernetes-cert-manager's CSI driver for automatic file updates.

**Monitoring**: Track certificate expiration using Prometheus metrics (`kubernetescertmanager_certificate_expiration_timestamp_seconds`) or external monitoring. While auto-renewal should work, monitoring ensures you catch failures before outages.

## Security Best Practices

### Least Privilege DNS Access

Each DNS provider has scoping mechanisms. Use them:

**Cloudflare**:
- Create API tokens scoped to specific zones (not "All zones")
- Grant only `Zone:Zone:Read` + `Zone:DNS:Edit` (not full zone edit)
- Set token expiration and rotate quarterly

**GCP**:
- Grant `dns.admin` role only on specific DNS zones (not project-wide)
- Use Workload Identity (no service account keys in Secrets)
- Enable DNS audit logs

**AWS**:
- IAM policy with `route53:ChangeResourceRecordSets` only on specific hosted zone ARNs
- Use IRSA (no access keys)
- Enable CloudTrail for Route53 API calls

**Azure**:
- Managed Identity with DNS Zone Contributor role only on specific resource groups
- Use Workload Identity Federation
- Enable Azure Monitor for DNS changes

### Secret Rotation

**Cloudflare API tokens**: Must be rotated manually. Strategy:
1. Create new token with same permissions
2. Update OpenMCF spec with new token
3. Redeploy addon (updates Secret, ClusterIssuer is unchanged)
4. Revoke old token after verification

**Cloud provider credentials**: Workload identity eliminates long-lived credentials. IAM role/policy changes take effect immediately without redeployment.

### Rate Limits and ACME Hygiene

Let's Encrypt production has rate limits:
- **50 certificates per registered domain per week**
- **5 duplicate certificates per week**
- **300 new orders per account per 3 hours**

**Strategies**:
- Use **Let's Encrypt staging** for testing (much higher limits)
- Consolidate multiple hostnames into SAN certificates where possible
- Avoid reissuing certificates unnecessarily (kubernetes-cert-manager's auto-renewal handles this)

## Troubleshooting Patterns

### Follow the Resource Chain

When certificates aren't being issued, trace the dependency chain:

```
Certificate → CertificateRequest → Order → Challenge
```

Each resource has status conditions. Start with the Certificate, check its status, follow to CertificateRequest, then Order, then Challenge. The Challenge resource shows DNS-01-specific details like which solver was selected and DNS propagation status.

### Common Failure Modes

**"Waiting for DNS propagation"**: The TXT record exists but kubernetes-cert-manager can't verify it. Usually resolves in 1-3 minutes. If stuck longer, check:
- DNS provider API logs (rate limiting, authentication errors)
- Manual TXT record verification: `dig _acme-challenge.example.com TXT @1.1.1.1`
- kubernetes-cert-manager logs for DNS resolver issues

**Authentication errors**: Provider-specific, but the pattern is consistent:
1. Check kubernetes-cert-manager logs for specific error messages
2. Verify credentials in the provider's console/CLI
3. For cloud providers, test ServiceAccount annotation and workload identity binding
4. Regenerate credentials if needed

**"No configured challenge solvers"**: The domain in the Certificate doesn't match any `dnsZones` in the ClusterIssuer's solvers. Either add the domain to an existing solver or create a new solver configuration.

## The Licensing and Distribution Landscape

kubernetes-cert-manager is **Apache 2.0 licensed**, free to use, modify, and distribute. It's a CNCF incubating project (since 2020), governed by The Linux Foundation.

**Official container images**: Published at `quay.io/jetstack/kubernetes-cert-manager-*` (controller, webhook, cainjector). Also available as OCI artifacts.

**Enterprise offerings**: Jetstack (now Venafi) offers **TLS Protect for Kubernetes**, which adds:
- UI and dashboards for certificate visibility
- Policy enforcement and compliance reporting
- Long-term support (LTS) versions with extended backports
- FIPS-compliant images

CyberArk provides **LTS kubernetes-cert-manager on AWS Marketplace** as an EKS add-on, offering:
- Subscription-based support
- FIPS-compliant builds
- Integration with AWS billing

**Alternative solutions**: Before kubernetes-cert-manager, projects like **kube-lego** and **kube-kubernetes-cert-manager** existed but are now deprecated. Today, kubernetes-cert-manager has no serious competition for Kubernetes certificate automation. External tools like Certbot or cloud-native certificate managers (AWS ACM, GCP Certificate Manager) serve different use cases (external automation or cloud load balancer integration, respectively).

## Why OpenMCF Chose This Approach

OpenMCF is opinionated infrastructure automation, not a configuration management tool. The philosophy: **encode production best practices, expose essential configuration, hide operational complexity**.

For kubernetes-cert-manager:

**What we abstracted away**:
- ClusterIssuer creation and maintenance
- DNS provider credential Secret creation
- ServiceAccount annotation for workload identity
- High availability Helm values
- DNS propagation resolver configuration
- Version compatibility enforcement

**What we preserved**:
- Choice of DNS providers (users specify their providers)
- ACME server selection (staging vs. production)
- DNS zone mapping (users control which domains use which providers)
- Multi-provider flexibility (mix Cloudflare, GCP, AWS, Azure)

**The result**: You declare intent ("I want kubernetes-cert-manager with these DNS providers"). OpenMCF produces a complete, working system. When you need to add a domain or rotate credentials, you update the spec and redeploy. The integration layer handles the rest.

This is the pattern OpenMCF applies across all cloud resources: **deploy once, integrate automatically, maintain declaratively**.

## Conclusion: Automation as a Strategic Choice

Deploying kubernetes-cert-manager is straightforward. Integrating it correctly—with secure DNS provider authentication, high availability, multi-cloud support, and operational best practices—is not. That integration complexity is undifferentiated heavy lifting: every team solves the same problems, writes the same ClusterIssuer YAML, makes the same mistakes with DNS propagation and credential management.

OpenMCF's choice is to solve it once, correctly, and expose a clean abstraction. You focus on **which domains need certificates**. The platform focuses on **how to get them reliably**.

That's the paradigm shift from deployment to automation: moving from "how do I install kubernetes-cert-manager?" to "how do I enable certificate automation for my domains across my infrastructure?". OpenMCF answers the latter by handling the former as an implementation detail.

