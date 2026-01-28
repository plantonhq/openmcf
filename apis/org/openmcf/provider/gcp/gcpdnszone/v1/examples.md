# GCP DNS Zone Examples

This document provides comprehensive examples for the `GcpDnsZone` API resource, demonstrating various DNS management scenarios in Google Cloud. Each example is designed to showcase different features and use cases, from minimal configurations to production-ready setups.

## Table of Contents

1. [Minimal Configuration](#minimal-configuration)
2. [Standard Production Setup](#standard-production-setup)
3. [Multiple Record Types](#multiple-record-types)
4. [Integration with IAM Service Accounts](#integration-with-iam-service-accounts)
5. [Comprehensive Production Zone](#comprehensive-production-zone)
6. [Wildcard Domain Configuration](#wildcard-domain-configuration)
7. [Subdomain Delegation](#subdomain-delegation)
8. [Using Foreign Key References](#using-foreign-key-references)

---

## Minimal Configuration

The simplest possible DNS zone configuration with only required fields. This creates a public DNS zone in GCP with DNSSEC enabled by default (managed by the underlying infrastructure module).

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpDnsZone
metadata:
  name: example.com
spec:
  projectId:
    value: my-gcp-project
```

**Use Case**: Quick setup for a domain where DNS records will be managed dynamically by tools like external-dns or manually through the GCP console.

**What Gets Created**:
- A public managed DNS zone for `example.com`
- DNSSEC automatically enabled with default configuration
- GCP-assigned nameservers (exported as stack outputs)

---

## Standard Production Setup

A typical production configuration with foundational DNS records for a web application.

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpDnsZone
metadata:
  name: mycompany.com
spec:
  projectId:
    value: production-gcp-project
  records:
    - recordType: A
      name: mycompany.com.
      values:
        - 104.198.14.52
      ttlSeconds: 300
    
    - recordType: A
      name: www.mycompany.com.
      values:
        - 104.198.14.52
      ttlSeconds: 300
    
    - recordType: CNAME
      name: api.mycompany.com.
      values:
        - mycompany.com.
      ttlSeconds: 300
```

**Use Case**: Standard web application with root domain and www subdomain pointing to a load balancer IP, plus an API subdomain using a CNAME.

**Important Notes**:
- All DNS names must end with a dot (`.`) to signify fully qualified domain names (FQDN)
- CNAME values must also end with a dot when pointing to another domain
- TTL of 300 seconds (5 minutes) provides a good balance between caching and flexibility
- The load balancer IP (104.198.14.52) should be replaced with your actual GKE/Cloud Load Balancer IP

---

## Multiple Record Types

Demonstrates various DNS record types supported by the API, including mail server configuration and domain verification records.

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpDnsZone
metadata:
  name: example.com
spec:
  projectId:
    value: my-gcp-project
  records:
    # A record for root domain
    - recordType: A
      name: example.com.
      values:
        - 203.0.113.10
      ttlSeconds: 300
    
    # CNAME for www subdomain
    - recordType: CNAME
      name: www.example.com.
      values:
        - example.com.
      ttlSeconds: 300
    
    # MX records for email routing
    - recordType: MX
      name: example.com.
      values:
        - 10 mail.example.com.
        - 20 backup-mail.example.com.
      ttlSeconds: 3600
    
    # A records for mail servers
    - recordType: A
      name: mail.example.com.
      values:
        - 203.0.113.20
      ttlSeconds: 300
    
    - recordType: A
      name: backup-mail.example.com.
      values:
        - 203.0.113.21
      ttlSeconds: 300
    
    # SPF record for email authentication
    - recordType: TXT
      name: example.com.
      values:
        - v=spf1 include:_spf.google.com ~all
      ttlSeconds: 3600
    
    # DKIM record for email signing
    - recordType: TXT
      name: default._domainkey.example.com.
      values:
        - v=DKIM1; k=rsa; p=MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQC...
      ttlSeconds: 3600
    
    # Google domain verification
    - recordType: TXT
      name: example.com.
      values:
        - google-site-verification=abc123xyz789
      ttlSeconds: 3600
    
    # AAAA record for IPv6
    - recordType: AAAA
      name: example.com.
      values:
        - 2001:0db8:85a3:0000:0000:8a2e:0370:7334
      ttlSeconds: 300
```

**Use Case**: Production domain with email services, domain verification, and IPv6 support.

**Key Points**:
- MX records require priority values (10, 20) before the mail server hostname
- Multiple values can be specified for redundancy (MX records have primary and backup)
- TXT records are commonly used for SPF (email authentication), DKIM (email signing), and domain verification
- Longer TTLs (3600 seconds = 1 hour) are appropriate for rarely-changing records like MX and TXT
- DKIM public keys in TXT records are typically very long; this example is truncated for readability

---

## Integration with IAM Service Accounts

Demonstrates IAM integration for automated DNS management by Kubernetes controllers like cert-manager and external-dns.

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpDnsZone
metadata:
  name: production.example.com
spec:
  projectId:
    value: prod-gcp-project
  
  # Grant DNS management permissions to Kubernetes workload identities
  iamServiceAccounts:
    - cert-manager@prod-gcp-project.iam.gserviceaccount.com
    - external-dns@prod-gcp-project.iam.gserviceaccount.com
  
  records:
    # Root domain A record
    - recordType: A
      name: production.example.com.
      values:
        - 104.198.14.52
      ttlSeconds: 300
    
    # TXT record for domain verification (static, managed by IaC)
    - recordType: TXT
      name: production.example.com.
      values:
        - google-site-verification=xyz789abc123
      ttlSeconds: 3600
```

**Use Case**: Production environment where:
- **cert-manager** needs to create temporary TXT records to solve DNS-01 ACME challenges for wildcard TLS certificates
- **external-dns** needs to create A/CNAME records automatically when Kubernetes Ingress resources are created

**How It Works**:
1. The `iamServiceAccounts` field grants `roles/dns.admin` to the specified service accounts
2. cert-manager can now create `_acme-challenge.production.example.com` TXT records during certificate issuance
3. external-dns can create records like `app.production.example.com` when Ingress resources are deployed
4. Static foundational records (like domain verification) remain managed by this IaC configuration

**Important**: This follows the "split state" pattern described in the research documentation:
- **IaC manages**: DNS zones and static foundational records
- **Automation tools manage**: Dynamic application records (Ingress A/CNAME records, ACME challenge TXT records)

---

## Comprehensive Production Zone

A complete production-ready configuration demonstrating all major features and best practices.

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpDnsZone
metadata:
  name: company.io
spec:
  projectId:
    value: production-infrastructure
  
  # Service accounts for automation tools
  iamServiceAccounts:
    - cert-manager@production-infrastructure.iam.gserviceaccount.com
    - external-dns@production-infrastructure.iam.gserviceaccount.com
    - cloudflare-dns-sync@production-infrastructure.iam.gserviceaccount.com
  
  records:
    # === Root Domain Configuration ===
    
    # Primary A record pointing to global load balancer
    - recordType: A
      name: company.io.
      values:
        - 35.201.85.123
      ttlSeconds: 300
    
    # IPv6 support
    - recordType: AAAA
      name: company.io.
      values:
        - 2600:1901:0:b::1
      ttlSeconds: 300
    
    # === Web Application Subdomains ===
    
    # WWW subdomain (CNAME to root)
    - recordType: CNAME
      name: www.company.io.
      values:
        - company.io.
      ttlSeconds: 300
    
    # API gateway
    - recordType: A
      name: api.company.io.
      values:
        - 35.201.85.124
      ttlSeconds: 300
    
    # Application environments (static records for infrastructure endpoints)
    - recordType: A
      name: staging.company.io.
      values:
        - 35.201.85.125
      ttlSeconds: 300
    
    # === Email Configuration ===
    
    # MX records with priority (10 = primary, 20 = backup)
    - recordType: MX
      name: company.io.
      values:
        - 10 aspmx.l.google.com.
        - 20 alt1.aspmx.l.google.com.
        - 30 alt2.aspmx.l.google.com.
      ttlSeconds: 3600
    
    # SPF record for email sender authentication
    - recordType: TXT
      name: company.io.
      values:
        - v=spf1 include:_spf.google.com include:sendgrid.net ~all
      ttlSeconds: 3600
    
    # DMARC policy for email authentication reporting
    - recordType: TXT
      name: _dmarc.company.io.
      values:
        - v=DMARC1; p=quarantine; rua=mailto:dmarc-reports@company.io; pct=100
      ttlSeconds: 3600
    
    # DKIM selector for Google Workspace
    - recordType: TXT
      name: google._domainkey.company.io.
      values:
        - v=DKIM1; k=rsa; p=MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA...
      ttlSeconds: 3600
    
    # === Domain Verification ===
    
    # Google Search Console verification
    - recordType: TXT
      name: company.io.
      values:
        - google-site-verification=AbC123XyZ789...
      ttlSeconds: 3600
    
    # === Service Discovery ===
    
    # Database endpoint (private, but using public DNS for simplicity)
    - recordType: CNAME
      name: db.company.io.
      values:
        - prod-postgres.c.production-infrastructure.internal.
      ttlSeconds: 60
    
    # Redis cache endpoint
    - recordType: CNAME
      name: cache.company.io.
      values:
        - prod-redis.c.production-infrastructure.internal.
      ttlSeconds: 60
    
    # === Monitoring and Observability ===
    
    # Grafana dashboard
    - recordType: A
      name: metrics.company.io.
      values:
        - 35.201.85.126
      ttlSeconds: 300
    
    # Status page
    - recordType: CNAME
      name: status.company.io.
      values:
        - statuspage.company.io.herokudns.com.
      ttlSeconds: 300
```

**Use Case**: Production environment with comprehensive DNS configuration including:
- Web application endpoints with IPv4 and IPv6 support
- Complete email authentication setup (MX, SPF, DKIM, DMARC)
- Domain verification records
- Service discovery endpoints
- Monitoring and status page endpoints
- IAM integration for automated DNS management

**Best Practices Demonstrated**:
1. **TTL Strategy**:
   - Short TTLs (60-300 seconds) for frequently-changing records or infrastructure endpoints
   - Long TTLs (3600 seconds = 1 hour) for stable configuration like email and domain verification
   
2. **Email Security**:
   - SPF record includes both Google Workspace and SendGrid (third-party email service)
   - DMARC policy set to "quarantine" with 100% enforcement and reporting enabled
   - DKIM keys properly configured for email signing
   
3. **CNAME Usage**:
   - External services (like Heroku status page) use CNAME to maintain their own IP management
   - Internal GCP services can use CNAME to reference Cloud SQL or Memorystore internal hostnames
   
4. **Split State Pattern**:
   - Static, foundational records are defined here
   - Dynamic application records (for Kubernetes Ingresses) will be managed by external-dns
   - Temporary ACME challenge records will be managed by cert-manager

**Security Note**: The IAM service accounts granted access should follow the principle of least privilege. In this example, we grant `roles/dns.admin` to the entire project. In a more restricted setup, you would use per-zone IAM policies (when available in the GCP provider) or create custom roles with limited permissions.

---

## Wildcard Domain Configuration

Demonstrates wildcard DNS records for dynamic subdomain routing, commonly used in multi-tenant SaaS applications.

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpDnsZone
metadata:
  name: saas-platform.io
spec:
  projectId:
    value: saas-production
  
  iamServiceAccounts:
    - cert-manager@saas-production.iam.gserviceaccount.com
  
  records:
    # Root domain
    - recordType: A
      name: saas-platform.io.
      values:
        - 34.120.45.67
      ttlSeconds: 300
    
    # Wildcard subdomain for customer tenants
    # Matches: customer1.saas-platform.io, customer2.saas-platform.io, etc.
    - recordType: A
      name: *.saas-platform.io.
      values:
        - 34.120.45.67
      ttlSeconds: 300
    
    # Wildcard subdomain for staging environments
    - recordType: A
      name: *.staging.saas-platform.io.
      values:
        - 34.120.45.68
      ttlSeconds: 300
    
    # Specific subdomain (takes precedence over wildcard)
    - recordType: A
      name: app.saas-platform.io.
      values:
        - 34.120.45.69
      ttlSeconds: 300
```

**Use Case**: Multi-tenant SaaS platform where each customer gets a unique subdomain (e.g., `acme-corp.saas-platform.io`, `widgets-inc.saas-platform.io`), all pointing to the same load balancer IP. The application backend uses the hostname to identify the tenant.

**How Wildcard DNS Works**:
- The `*.saas-platform.io` record matches **any** subdomain under `saas-platform.io`
- Specific records (like `app.saas-platform.io`) take precedence over the wildcard
- Nested wildcards (like `*.staging.saas-platform.io`) work for second-level subdomains

**Important for TLS Certificates**:
- Wildcard DNS requires a **wildcard TLS certificate** (e.g., `*.saas-platform.io`)
- cert-manager can provision wildcard certificates using DNS-01 ACME challenges
- The `iam_service_accounts` field grants cert-manager the necessary permissions

**Caveats**:
- Wildcard records **do not** match the root domain itself (e.g., `*.example.com` does not match `example.com`)
- Wildcard records only match one level (e.g., `*.example.com` matches `sub.example.com` but not `deep.sub.example.com`)

---

## Subdomain Delegation

Demonstrates subdomain delegation using NS records, allowing different DNS zones to manage different parts of the domain hierarchy.

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpDnsZone
metadata:
  name: company.net
spec:
  projectId:
    value: infrastructure-project
  
  records:
    # Root domain A record
    - recordType: A
      name: company.net.
      values:
        - 203.0.113.50
      ttlSeconds: 300
    
    # WWW subdomain
    - recordType: CNAME
      name: www.company.net.
      values:
        - company.net.
      ttlSeconds: 300
    
    # Delegate 'dev.company.net' subdomain to separate DNS zone
    # (Managed by development team in a different GCP project)
    - recordType: NS
      name: dev.company.net.
      values:
        - ns-cloud-d1.googledomains.com.
        - ns-cloud-d2.googledomains.com.
        - ns-cloud-d3.googledomains.com.
        - ns-cloud-d4.googledomains.com.
      ttlSeconds: 300
    
    # Delegate 'staging.company.net' subdomain to separate DNS zone
    - recordType: NS
      name: staging.company.net.
      values:
        - ns-cloud-e1.googledomains.com.
        - ns-cloud-e2.googledomains.com.
        - ns-cloud-e3.googledomains.com.
        - ns-cloud-e4.googledomains.com.
      ttlSeconds: 300
    
    # Delegate 'partner.company.net' to external DNS provider (e.g., Cloudflare)
    - recordType: NS
      name: partner.company.net.
      values:
        - dante.ns.cloudflare.com.
        - gail.ns.cloudflare.com.
      ttlSeconds: 3600
```

**Use Case**: Large organization where different teams or departments manage their own DNS zones independently:
- **Infrastructure team** manages the root `company.net` zone (this configuration)
- **Development team** manages `dev.company.net` in a separate GCP project
- **Staging team** manages `staging.company.net` in another GCP project
- **Partner/external team** manages `partner.company.net` using Cloudflare

**How Subdomain Delegation Works**:
1. The parent zone (`company.net`) contains NS records pointing to the nameservers of the child zone
2. When a DNS resolver queries `app.dev.company.net`:
   - It first queries the root nameservers for `company.net`
   - Receives the NS records pointing to the nameservers for `dev.company.net`
   - Queries those nameservers directly for `app.dev.company.net`
3. The child zone (`dev.company.net`) has full autonomy to manage its own DNS records

**Setup Steps**:
1. Create the child DNS zone (e.g., `dev.company.net`) as a separate `GcpDnsZone` resource
2. Note the nameservers assigned to the child zone (available in stack outputs)
3. Add NS records in the parent zone (this configuration) pointing to the child zone's nameservers
4. Verify delegation using `dig NS dev.company.net` or `nslookup -type=NS dev.company.net`

**Benefits**:
- **Team Autonomy**: Each team manages their own DNS records without needing access to the root zone
- **Blast Radius Reduction**: Mistakes in one zone don't affect other zones
- **Cross-Project Isolation**: Different GCP projects can manage different parts of the domain hierarchy
- **Hybrid Cloud**: Subdomains can be delegated to external DNS providers (AWS Route 53, Cloudflare, etc.)

**Important Notes**:
- The nameservers listed in the NS records must be the **actual nameservers assigned by GCP** when the child zone is created
- You cannot create records in the parent zone that overlap with a delegated subdomain (e.g., after delegating `dev.company.net`, you cannot create `app.dev.company.net` in the parent zone)
- NS records typically require 4 nameservers for redundancy (GCP provides 4 by default)

---

## Notes on DNSSEC

All DNS zones created with this API resource have **DNSSEC automatically enabled by default** (configured in the underlying Pulumi/Terraform module). This provides cryptographic authentication of DNS responses, protecting against cache poisoning and man-in-the-middle attacks.

**What You Need to Do**:
1. After creating the zone, retrieve the **DS (Delegation Signer) record** from the stack outputs or GCP console
2. Add the DS record to your domain registrar (the entity where you purchased the domain)
3. Verify the DNSSEC chain using tools like:
   - [Verisign DNSSEC Debugger](https://dnssec-debugger.verisignlabs.com/)
   - [DNSViz](https://dnsviz.net/)

**Do Not**:
- Disable DNSSEC in the zone without removing the DS record from your registrar (or vice versa) — this will break DNS resolution
- Forget to update the DS record if you ever recreate the zone (new keys are generated)

---

## Dynamic Records and the "Split State" Pattern

The examples above demonstrate **static DNS records** that are explicitly defined in the infrastructure configuration. However, in Kubernetes environments, many DNS records are created **dynamically** by automation tools.

### Records Managed Outside of This Configuration

The following types of records are typically **not** defined in these YAML configurations:

1. **Kubernetes Ingress A/CNAME Records**:
   - Managed automatically by **external-dns**
   - Created when Ingress resources are deployed
   - Example: `app.company.io` → load balancer IP (created by external-dns, not IaC)

2. **ACME DNS-01 Challenge TXT Records**:
   - Managed temporarily by **cert-manager**
   - Created during TLS certificate issuance
   - Example: `_acme-challenge.company.io` → token (created by cert-manager, deleted after verification)

3. **Service Discovery Records**:
   - Some organizations use external-dns to create records for Kubernetes Services (not just Ingresses)
   - Example: `postgres.company.io` → internal service endpoint

### Why This Separation Matters

If you define application DNS records in this IaC configuration **and** use external-dns, you create a conflict:
- IaC will try to enforce its view of what records should exist
- external-dns will create its own records dynamically
- Result: perpetual "war" between the two systems, with records being created and deleted repeatedly

**The Correct Pattern** (as described in the research documentation):
- Use this API resource for **zones** and **foundational static records** (MX, TXT, NS, verification records)
- Use **external-dns** for dynamic application records (Ingress A/CNAME records)
- Use **cert-manager** for temporary ACME challenge records

The `iamServiceAccounts` field enables this pattern by granting the necessary permissions to external-dns and cert-manager to manage their own records within the zone you create.

---

## Advanced Features Not Supported in This API

The following advanced Cloud DNS features are **intentionally omitted** from this API (representing the "20%" of features used by a minority of users):

1. **Private DNS Zones**: Zones that resolve only within GCP VPC networks
2. **Forwarding Zones**: Forward queries for specific domains to on-premises DNS servers
3. **Peering Zones**: Share DNS zones across VPC networks
4. **Advanced Routing Policies**:
   - Geolocation routing (route based on user's geographic location)
   - Weighted round-robin (distribute traffic across multiple IPs by percentage)
   - Failover routing (automatic failover based on health checks)

**If You Need These Features**:
- Manage them using direct Terraform resources (`google_dns_managed_zone` with additional arguments)
- Or use GCP Config Connector CRDs (`DNSManagedZone` with full feature support)
- The "80/20 principle" in OpenMCF means the API covers the most common use cases; advanced features remain accessible through underlying IaC tools

---

## Troubleshooting

### DNS Records Not Resolving

**Check**:
1. Verify the zone was created successfully (check stack outputs for nameservers)
2. Ensure your domain registrar has the correct NS records pointing to GCP's nameservers
3. DNS propagation can take up to 48 hours (but usually 5-10 minutes)
4. Use `dig @8.8.8.8 yourdomain.com` to query Google's public DNS directly
5. Use `dig +trace yourdomain.com` to see the full DNS resolution path

### DNSSEC Validation Failures

**Check**:
1. Verify the DS record in your registrar matches the one provided by Cloud DNS
2. Ensure DNSSEC is enabled in both Cloud DNS and your registrar
3. Use [DNSSEC Debugger](https://dnssec-debugger.verisignlabs.com/) to identify chain-of-trust issues

### IAM Permission Errors

**Check**:
1. Verify the service account emails are correctly specified
2. Ensure the service accounts exist in the GCP project
3. Check that Workload Identity is configured (if using GKE with external-dns/cert-manager)
4. Review GCP IAM audit logs for permission denial events

### Terraform/Pulumi State Conflicts

**Check**:
1. If external-dns is managing records, do **not** define those records in this configuration
2. Review the "split state" pattern in this document
3. external-dns creates TXT records to claim ownership; respect those claims in your IaC

---

## Using Foreign Key References

Instead of hardcoding the GCP project ID, you can reference another OpenMCF resource (like a `GcpProject`) using the `valueFrom` pattern. This enables clean dependency management and ensures the DNS zone is created only after the referenced project exists.

### GcpProject Resource

First, define the GCP project that will host the DNS zone:

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpProject
metadata:
  name: dns-project
spec:
  projectId: my-dns-project-123
  billingAccountId: 012345-ABCDEF-678910
  folderId: "123456789012"
```

### DNS Zone with Project Reference

Then, reference the project in your DNS zone configuration:

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpDnsZone
metadata:
  name: example.com
spec:
  projectId:
    valueFrom:
      kind: GcpProject
      name: dns-project
      fieldPath: status.outputs.project_id
  records:
    - recordType: A
      name: example.com.
      values:
        - 104.198.14.52
      ttlSeconds: 300
```

### How Foreign Key References Work

1. OpenMCF creates the `GcpProject` resource first
2. Once the project exists, it extracts `status.outputs.project_id`
3. The `GcpDnsZone` resource uses that project ID automatically
4. This ensures proper dependency ordering (DNS zone created after project)

### When to Use Foreign Key References

**Use references when:**
- You're managing both the project and DNS zone via OpenMCF
- You want explicit dependency ordering in your infrastructure code
- You're building reusable templates that reference other resources
- You want to avoid hardcoding project IDs across multiple manifests

**Use direct values when:**
- The project already exists and is managed outside OpenMCF
- You need to deploy the DNS zone independently without dependency tracking
- You're working with a simple, single-environment setup

### Combined Example with Multiple References

This example shows a production DNS zone that references a GcpProject for its project ID and grants permissions to service accounts:

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpDnsZone
metadata:
  name: production.mycompany.io
spec:
  projectId:
    valueFrom:
      kind: GcpProject
      name: production-infra
      fieldPath: status.outputs.project_id
  iamServiceAccounts:
    - cert-manager@production-infra-123.iam.gserviceaccount.com
    - external-dns@production-infra-123.iam.gserviceaccount.com
  records:
    - recordType: A
      name: production.mycompany.io.
      values:
        - 35.201.85.123
      ttlSeconds: 300
    - recordType: TXT
      name: production.mycompany.io.
      values:
        - "v=spf1 include:_spf.google.com ~all"
      ttlSeconds: 3600
```

---

## Summary

These examples demonstrate the full range of configurations supported by the `GcpDnsZone` API resource:

- **Minimal**: Quick setup for domains with dynamic records managed elsewhere
- **Standard**: Production web applications with foundational A/CNAME records
- **Comprehensive**: Complete DNS setup including email, verification, and service discovery
- **Wildcard**: Multi-tenant SaaS platforms with dynamic subdomains
- **Delegation**: Large organizations with decentralized DNS management
- **Foreign Key References**: Dynamic project ID resolution from other resources

All configurations follow the "split state" pattern, ensuring harmony between IaC-managed static records and dynamically-managed application records.
