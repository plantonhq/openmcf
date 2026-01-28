# Azure DNS Zone Examples

This document provides comprehensive examples for the `AzureDnsZone` API resource, demonstrating various DNS management scenarios in Microsoft Azure. Each example is designed to showcase different features and use cases, from minimal configurations to production-ready setups.

## Table of Contents

1. [Minimal Configuration](#minimal-configuration)
2. [Standard Production Setup](#standard-production-setup)
3. [Complete Email Configuration](#complete-email-configuration)
4. [Multi-Environment Setup](#multi-environment-setup)
5. [Comprehensive Production Zone](#comprehensive-production-zone)
6. [Wildcard Domain Configuration](#wildcard-domain-configuration)
7. [Subdomain Delegation](#subdomain-delegation)

---

## Minimal Configuration

The simplest possible DNS zone configuration with only required fields. This creates a public DNS zone in Azure.

```yaml
apiVersion: azure.openmcf.org/v1
kind: AzureDnsZone
metadata:
  name: example.com
spec:
  zoneName: example.com
  resourceGroup: prod-network-rg
```

**Use Case**: Quick setup for a domain where DNS records will be managed dynamically by tools like external-dns or manually through the Azure portal.

**What Gets Created**:
- A public DNS zone for `example.com`
- Azure-assigned nameservers (exported as stack outputs)
- No DNS records (zone is empty until records are added)

**Next Steps**:
1. Retrieve nameservers from stack outputs
2. Update your domain registrar with Azure's nameservers
3. Add DNS records either through the API or using external-dns

---

## Standard Production Setup

A typical production configuration with foundational DNS records for a web application.

```yaml
apiVersion: azure.openmcf.org/v1
kind: AzureDnsZone
metadata:
  name: mycompany.com
  org: mycompany
  env: production
spec:
  zoneName: mycompany.com
  resourceGroup: prod-network-rg
  records:
    # Root domain A record
    - recordType: A
      name: mycompany.com.
      values:
        - 20.42.15.123
      ttlSeconds: 300
    
    # WWW subdomain A record
    - recordType: A
      name: www.mycompany.com.
      values:
        - 20.42.15.123
      ttlSeconds: 300
    
    # API subdomain CNAME
    - recordType: CNAME
      name: api.mycompany.com.
      values:
        - mycompany.com.
      ttlSeconds: 300
```

**Use Case**: Standard web application with root domain and www subdomain pointing to an Azure load balancer IP, plus an API subdomain using a CNAME.

**Important Notes**:
- All DNS names must end with a dot (`.`) to signify fully qualified domain names (FQDN)
- CNAME values must also end with a dot when pointing to another domain
- TTL of 300 seconds (5 minutes) provides a good balance between caching and flexibility
- The load balancer IP (20.42.15.123) should be replaced with your actual Azure Load Balancer or Application Gateway IP

**Azure-Specific Considerations**:
- Consider using Azure Front Door or Application Gateway for global load balancing
- For AKS workloads, the IP would be the ingress controller's external IP
- Azure CDN can be integrated using CNAME records to Azure CDN endpoints

---

## Complete Email Configuration

Demonstrates comprehensive email setup with MX, SPF, DKIM, DMARC, and CAA records.

```yaml
apiVersion: azure.openmcf.org/v1
kind: AzureDnsZone
metadata:
  name: contoso.com
  org: contoso
  env: production
spec:
  zoneName: contoso.com
  resourceGroup: prod-network-rg
  records:
    # Root domain A record
    - recordType: A
      name: contoso.com.
      values:
        - 198.51.100.45
      ttlSeconds: 300
    
    # WWW subdomain
    - recordType: CNAME
      name: www.contoso.com.
      values:
        - contoso.com.
      ttlSeconds: 300
    
    # MX records for email routing
    - recordType: MX
      name: contoso.com.
      values:
        - mail.contoso.com.
        - backup-mail.contoso.com.
      ttlSeconds: 3600
    
    # Mail server A records
    - recordType: A
      name: mail.contoso.com.
      values:
        - 198.51.100.100
      ttlSeconds: 300
    
    - recordType: A
      name: backup-mail.contoso.com.
      values:
        - 198.51.100.101
      ttlSeconds: 300
    
    # SPF record for email authentication
    - recordType: TXT
      name: contoso.com.
      values:
        - v=spf1 include:mail.contoso.com -all
      ttlSeconds: 3600
    
    # DMARC record for email policy
    - recordType: TXT
      name: _dmarc.contoso.com.
      values:
        - v=DMARC1; p=reject; rua=mailto:dmarc@contoso.com
      ttlSeconds: 3600
    
    # DKIM selector record
    - recordType: TXT
      name: selector1._domainkey.contoso.com.
      values:
        - v=DKIM1; k=rsa; p=MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQC...
      ttlSeconds: 3600
    
    # Certificate Authority Authorization
    - recordType: CAA
      name: contoso.com.
      values:
        - letsencrypt.org
      ttlSeconds: 86400
```

**Use Case**: Production domain with complete email infrastructure including authentication and security records.

**Key Points**:
- **MX Records**: Configure mail routing with primary and backup servers
- **SPF (Sender Policy Framework)**: Prevents email spoofing by specifying authorized mail servers
- **DKIM (DomainKeys Identified Mail)**: Cryptographic signing of emails for authenticity
- **DMARC (Domain-based Message Authentication)**: Policy for handling failed SPF/DKIM checks
- **CAA Records**: Control which certificate authorities can issue TLS certificates for your domain

**Email Provider Integration**:
- **Microsoft 365**: Use Microsoft's MX records like `contoso-com.mail.protection.outlook.com`
- **Google Workspace**: Use Google's MX records like `aspmx.l.google.com`
- **SendGrid/Mailgun**: Include their domains in SPF and add their DKIM records

---

## Multi-Environment Setup

Demonstrates DNS configuration for multiple environments (dev, staging, production) using the same domain.

```yaml
apiVersion: azure.openmcf.org/v1
kind: AzureDnsZone
metadata:
  name: company.io
  org: company
  env: shared
spec:
  zoneName: company.io
  resourceGroup: shared-network-rg
  records:
    # Production environment (root domain)
    - recordType: A
      name: company.io.
      values:
        - 20.42.10.100
      ttlSeconds: 300
    
    - recordType: A
      name: www.company.io.
      values:
        - 20.42.10.100
      ttlSeconds: 300
    
    # API endpoints per environment
    - recordType: A
      name: api.company.io.
      values:
        - 20.42.10.101
      ttlSeconds: 300
    
    - recordType: A
      name: api.staging.company.io.
      values:
        - 20.42.10.201
      ttlSeconds: 300
    
    - recordType: A
      name: api.dev.company.io.
      values:
        - 20.42.10.301
      ttlSeconds: 300
    
    # Application endpoints per environment
    - recordType: A
      name: app.company.io.
      values:
        - 20.42.10.102
      ttlSeconds: 300
    
    - recordType: A
      name: app.staging.company.io.
      values:
        - 20.42.10.202
      ttlSeconds: 300
    
    - recordType: A
      name: app.dev.company.io.
      values:
        - 20.42.10.302
      ttlSeconds: 300
    
    # Shared services
    - recordType: A
      name: admin.company.io.
      values:
        - 20.42.10.150
      ttlSeconds: 300
    
    - recordType: A
      name: monitoring.company.io.
      values:
        - 20.42.10.151
      ttlSeconds: 300
```

**Use Case**: Organization with multiple environments using subdomain prefixes to separate dev, staging, and production traffic.

**Architecture Pattern**:
- **Production**: Root domain and `www` subdomain
- **Staging**: `*.staging.company.io` subdomains
- **Development**: `*.dev.company.io` subdomains
- **Shared Services**: Common services like admin panels and monitoring

**Best Practices**:
- Use separate Azure resource groups for each environment
- Consider using Azure Traffic Manager for environment-specific routing
- Lower TTLs for dev/staging environments for faster iteration
- Higher TTLs for production to reduce DNS query load

---

## Comprehensive Production Zone

A complete production-ready configuration demonstrating all major features and best practices.

```yaml
apiVersion: azure.openmcf.org/v1
kind: AzureDnsZone
metadata:
  name: enterprise.com
  id: dns-enterprise-prod
  org: enterprise
  env: production
spec:
  zoneName: enterprise.com
  resourceGroup: prod-network-rg
  records:
    # === Root Domain Configuration ===
    
    # Primary A record pointing to Azure Front Door
    - recordType: A
      name: enterprise.com.
      values:
        - 20.50.120.45
      ttlSeconds: 300
    
    # IPv6 support
    - recordType: AAAA
      name: enterprise.com.
      values:
        - 2603:1030:20::1
      ttlSeconds: 300
    
    # === Web Application Subdomains ===
    
    # WWW subdomain (CNAME to root)
    - recordType: CNAME
      name: www.enterprise.com.
      values:
        - enterprise.com.
      ttlSeconds: 300
    
    # API gateway
    - recordType: A
      name: api.enterprise.com.
      values:
        - 20.50.120.46
      ttlSeconds: 300
    
    # Mobile API
    - recordType: A
      name: mobile-api.enterprise.com.
      values:
        - 20.50.120.47
      ttlSeconds: 300
    
    # Admin portal
    - recordType: A
      name: admin.enterprise.com.
      values:
        - 20.50.120.48
      ttlSeconds: 300
    
    # === Email Configuration ===
    
    # MX records (Microsoft 365)
    - recordType: MX
      name: enterprise.com.
      values:
        - enterprise-com.mail.protection.outlook.com.
      ttlSeconds: 3600
    
    # SPF record for Microsoft 365 and SendGrid
    - recordType: TXT
      name: enterprise.com.
      values:
        - v=spf1 include:spf.protection.outlook.com include:sendgrid.net ~all
      ttlSeconds: 3600
    
    # DMARC policy
    - recordType: TXT
      name: _dmarc.enterprise.com.
      values:
        - v=DMARC1; p=quarantine; rua=mailto:dmarc-reports@enterprise.com; pct=100
      ttlSeconds: 3600
    
    # DKIM selectors for Microsoft 365
    - recordType: TXT
      name: selector1._domainkey.enterprise.com.
      values:
        - v=DKIM1; k=rsa; p=MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA...
      ttlSeconds: 3600
    
    - recordType: TXT
      name: selector2._domainkey.enterprise.com.
      values:
        - v=DKIM1; k=rsa; p=MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA...
      ttlSeconds: 3600
    
    # === Domain Verification ===
    
    # Microsoft domain verification
    - recordType: TXT
      name: enterprise.com.
      values:
        - MS=ms12345678
      ttlSeconds: 3600
    
    # Azure AD domain verification
    - recordType: TXT
      name: _azuread.enterprise.com.
      values:
        - azuread-verification-abc123xyz789
      ttlSeconds: 3600
    
    # === Azure Service Integration ===
    
    # Azure CDN endpoint
    - recordType: CNAME
      name: cdn.enterprise.com.
      values:
        - enterprise-cdn.azureedge.net.
      ttlSeconds: 3600
    
    # Azure Front Door (custom domain)
    - recordType: CNAME
      name: global.enterprise.com.
      values:
        - enterprise-afd.azurefd.net.
      ttlSeconds: 300
    
    # Azure Storage static website
    - recordType: CNAME
      name: docs.enterprise.com.
      values:
        - enterprisedocs.z13.web.core.windows.net.
      ttlSeconds: 3600
    
    # === Monitoring and Observability ===
    
    # Application Insights endpoint
    - recordType: CNAME
      name: telemetry.enterprise.com.
      values:
        - enterprise-appinsights.azurewebsites.net.
      ttlSeconds: 300
    
    # Status page (hosted on Azure Static Web Apps)
    - recordType: CNAME
      name: status.enterprise.com.
      values:
        - wonderful-bay-12345.azurestaticapps.net.
      ttlSeconds: 300
    
    # === Security Records ===
    
    # CAA record (Let's Encrypt and DigiCert)
    - recordType: CAA
      name: enterprise.com.
      values:
        - letsencrypt.org
        - digicert.com
      ttlSeconds: 86400
    
    # === Regional Endpoints ===
    
    # US region
    - recordType: A
      name: us.enterprise.com.
      values:
        - 20.50.120.50
      ttlSeconds: 300
    
    # Europe region
    - recordType: A
      name: eu.enterprise.com.
      values:
        - 20.76.45.30
      ttlSeconds: 300
    
    # Asia region
    - recordType: A
      name: asia.enterprise.com.
      values:
        - 20.195.65.20
      ttlSeconds: 300
```

**Use Case**: Enterprise production environment with comprehensive DNS configuration including:
- Multi-region application endpoints
- Complete email authentication setup (Microsoft 365 + SendGrid)
- Azure service integrations (CDN, Front Door, Static Web Apps)
- Security records (CAA)
- Monitoring and observability endpoints
- IPv4 and IPv6 support

**Best Practices Demonstrated**:

1. **TTL Strategy**:
   - Short TTLs (300s) for application endpoints that may change
   - Medium TTLs (3600s = 1 hour) for stable configuration like email and CDN
   - Long TTLs (86400s = 1 day) for security records like CAA

2. **Azure Service Integration**:
   - CNAMEs for Azure CDN, Front Door, and Static Web Apps
   - Proper domain verification for Azure AD and Microsoft 365
   - Regional A records for multi-region deployments

3. **Email Security**:
   - SPF includes both Microsoft 365 and third-party services (SendGrid)
   - DMARC policy with quarantine and reporting
   - Dual DKIM selectors for Microsoft 365 (recommended practice)

4. **Security**:
   - CAA records limit certificate issuance to trusted CAs
   - Multiple verification TXT records for Azure services
   - Separate endpoints for different security zones

---

## Wildcard Domain Configuration

Demonstrates wildcard DNS records for dynamic subdomain routing, commonly used in multi-tenant SaaS applications.

```yaml
apiVersion: azure.openmcf.org/v1
kind: AzureDnsZone
metadata:
  name: saas-app.io
  org: saas-company
  env: production
spec:
  zoneName: saas-app.io
  resourceGroup: prod-network-rg
  records:
    # Root domain
    - recordType: A
      name: saas-app.io.
      values:
        - 20.85.45.100
      ttlSeconds: 300
    
    # Marketing site
    - recordType: A
      name: www.saas-app.io.
      values:
        - 20.85.45.100
      ttlSeconds: 300
    
    # Wildcard subdomain for customer tenants
    # Matches: customer1.saas-app.io, customer2.saas-app.io, etc.
    - recordType: A
      name: *.saas-app.io.
      values:
        - 20.85.45.101
      ttlSeconds: 60
    
    # Wildcard subdomain for staging environments
    - recordType: A
      name: *.staging.saas-app.io.
      values:
        - 20.85.45.201
      ttlSeconds: 60
    
    # Wildcard for preview environments
    - recordType: A
      name: *.preview.saas-app.io.
      values:
        - 20.85.45.202
      ttlSeconds: 60
    
    # Specific subdomain (takes precedence over wildcard)
    - recordType: A
      name: app.saas-app.io.
      values:
        - 20.85.45.102
      ttlSeconds: 300
    
    # Admin portal (specific, not wildcard)
    - recordType: A
      name: admin.saas-app.io.
      values:
        - 20.85.45.103
      ttlSeconds: 300
    
    # API endpoint
    - recordType: A
      name: api.saas-app.io.
      values:
        - 20.85.45.104
      ttlSeconds: 300
```

**Use Case**: Multi-tenant SaaS platform where each customer gets a unique subdomain (e.g., `acme-corp.saas-app.io`, `widgets-inc.saas-app.io`), all pointing to the same Azure Application Gateway or AKS ingress controller. The application backend uses the hostname to identify the tenant.

**How Wildcard DNS Works**:
- The `*.saas-app.io` record matches **any** subdomain under `saas-app.io`
- Specific records (like `app.saas-app.io`) take precedence over the wildcard
- Nested wildcards (like `*.staging.saas-app.io`) work for second-level subdomains
- Very low TTL (60 seconds) allows for rapid tenant provisioning changes

**Azure-Specific Implementation**:

1. **Azure Application Gateway with Wildcard TLS**:
   - Configure Application Gateway with wildcard certificate (`*.saas-app.io`)
   - Use hostname-based routing rules to identify tenants
   - Backend pool can be Azure Kubernetes Service (AKS)

2. **AKS with Cert-Manager**:
   - Use cert-manager to provision wildcard certificates via DNS-01 challenge
   - Configure ingress-nginx or Azure Application Gateway Ingress Controller
   - Each tenant ingress references the wildcard certificate

**Important for TLS Certificates**:
- Wildcard DNS requires a **wildcard TLS certificate** (e.g., `*.saas-app.io`)
- Use Azure Key Vault to store wildcard certificates
- cert-manager can automate wildcard certificate renewal via DNS-01 challenges

**Caveats**:
- Wildcard records **do not** match the root domain itself (e.g., `*.example.com` does not match `example.com`)
- Wildcard records only match one level (e.g., `*.example.com` matches `sub.example.com` but not `deep.sub.example.com`)

---

## Subdomain Delegation

Demonstrates subdomain delegation using NS records, allowing different DNS zones to manage different parts of the domain hierarchy.

```yaml
apiVersion: azure.openmcf.org/v1
kind: AzureDnsZone
metadata:
  name: bigcorp.net
  org: bigcorp
  env: shared
spec:
  zoneName: bigcorp.net
  resourceGroup: shared-network-rg
  records:
    # Root domain A record
    - recordType: A
      name: bigcorp.net.
      values:
        - 203.0.113.100
      ttlSeconds: 300
    
    # WWW subdomain
    - recordType: CNAME
      name: www.bigcorp.net.
      values:
        - bigcorp.net.
      ttlSeconds: 300
    
    # Delegate 'dev.bigcorp.net' subdomain to separate Azure DNS zone
    # (Managed by development team in a different resource group/subscription)
    - recordType: NS
      name: dev.bigcorp.net.
      values:
        - ns1-01.azure-dns.com.
        - ns2-01.azure-dns.net.
        - ns3-01.azure-dns.org.
        - ns4-01.azure-dns.info.
      ttlSeconds: 300
    
    # Delegate 'staging.bigcorp.net' subdomain to separate DNS zone
    - recordType: NS
      name: staging.bigcorp.net.
      values:
        - ns1-02.azure-dns.com.
        - ns2-02.azure-dns.net.
        - ns3-02.azure-dns.org.
        - ns4-02.azure-dns.info.
      ttlSeconds: 300
    
    # Delegate 'partner.bigcorp.net' to external DNS provider (Cloudflare)
    - recordType: NS
      name: partner.bigcorp.net.
      values:
        - dante.ns.cloudflare.com.
        - gail.ns.cloudflare.com.
      ttlSeconds: 3600
    
    # Delegate 'internal.bigcorp.net' to on-premises DNS servers
    - recordType: NS
      name: internal.bigcorp.net.
      values:
        - ns1.onprem.bigcorp.net.
        - ns2.onprem.bigcorp.net.
      ttlSeconds: 3600
```

**Use Case**: Large organization where different teams or departments manage their own DNS zones independently:
- **Infrastructure team** manages the root `bigcorp.net` zone (this configuration)
- **Development team** manages `dev.bigcorp.net` in a separate Azure DNS zone
- **Staging team** manages `staging.bigcorp.net` in another Azure DNS zone
- **Partner/external team** manages `partner.bigcorp.net` using Cloudflare
- **On-premises team** manages `internal.bigcorp.net` using on-prem DNS servers

**How Subdomain Delegation Works**:

1. The parent zone (`bigcorp.net`) contains NS records pointing to the nameservers of the child zone
2. When a DNS resolver queries `app.dev.bigcorp.net`:
   - It first queries the root nameservers for `bigcorp.net`
   - Receives the NS records pointing to the nameservers for `dev.bigcorp.net`
   - Queries those nameservers directly for `app.dev.bigcorp.net`
3. The child zone (`dev.bigcorp.net`) has full autonomy to manage its own DNS records

**Setup Steps**:

1. **Create the child DNS zone** as a separate `AzureDnsZone` resource:
   ```yaml
   apiVersion: azure.openmcf.org/v1
   kind: AzureDnsZone
   metadata:
     name: dev.bigcorp.net
   spec:
     zoneName: dev.bigcorp.net
     resourceGroup: dev-network-rg
   ```

2. **Retrieve nameservers** from the child zone's stack outputs

3. **Add NS records** in the parent zone (this configuration) pointing to the child zone's nameservers

4. **Verify delegation**:
   ```bash
   dig NS dev.bigcorp.net
   # or
   nslookup -type=NS dev.bigcorp.net
   ```

**Benefits**:
- **Team Autonomy**: Each team manages their own DNS records without needing access to the root zone
- **Blast Radius Reduction**: Mistakes in one zone don't affect other zones
- **Cross-Subscription Isolation**: Different Azure subscriptions can manage different parts of the domain hierarchy
- **Hybrid Cloud**: Subdomains can be delegated to on-premises DNS or external providers (AWS Route 53, Cloudflare, etc.)

**Important Notes**:
- The nameservers listed in the NS records must be the **actual nameservers assigned by Azure** when the child zone is created
- You cannot create records in the parent zone that overlap with a delegated subdomain
- NS records typically require 4 nameservers for redundancy (Azure provides 4 by default)
- For cross-subscription scenarios, use Azure Lighthouse or appropriate RBAC policies

---

## Azure-Specific Integration Patterns

### Integration with Azure Kubernetes Service (AKS)

Use **external-dns** to automatically manage DNS records based on Kubernetes Ingresses:

```yaml
# Ingress with automatic DNS management
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: myapp
  annotations:
    external-dns.alpha.kubernetes.io/hostname: myapp.company.io
spec:
  rules:
  - host: myapp.company.io
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

**Requirements**:
- external-dns deployed in AKS cluster
- Managed identity or service principal with DNS Zone Contributor role
- DNS zone must exist (created by this API)

### Integration with Azure App Service

Custom domain configuration for Azure App Service:

```yaml
# DNS records needed for App Service custom domain
records:
  # A record pointing to App Service IP
  - record_type: A
    name: webapp.company.io.
    values:
      - 20.42.15.123
    ttl_seconds: 300
  
  # TXT record for domain verification
  - record_type: TXT
    name: asuid.webapp.company.io.
    values:
      - verification-string-from-azure
    ttl_seconds: 300
```

### Integration with Azure Front Door

Configure custom domains with Azure Front Door:

```yaml
# CNAME record for Azure Front Door
records:
  - record_type: CNAME
    name: www.company.io.
    values:
      - company-afd.azurefd.net.
    ttl_seconds: 3600
  
  # Validation record for Front Door
  - record_type: TXT
    name: _afdverify.www.company.io.
    values:
      - _afdverify.company-afd.azurefd.net
    ttl_seconds: 300
```

---

## Best Practices Summary

### TTL Strategy
- **60 seconds**: Rapidly changing records (load testing, blue-green deployments)
- **300 seconds (5 minutes)**: Standard application endpoints
- **3600 seconds (1 hour)**: Stable records (email, verification)
- **86400 seconds (1 day)**: Security records (CAA)

### Migration Strategy
When migrating from another DNS provider:

1. **Lower TTLs** at old provider 48 hours before migration
2. **Create zone** and records in Azure (this API)
3. **Verify all records** with `dig @ns1-01.azure-dns.com domain.com`
4. **Update nameservers** at domain registrar
5. **Monitor** for 24-48 hours
6. **Raise TTLs** back to normal values

### Security
- **Always use CAA records** to restrict certificate issuance
- **Implement SPF, DKIM, DMARC** for email domains
- **Use separate zones** for different environments (dev/staging/prod)
- **Apply Azure RBAC** to limit who can modify DNS records

### Resource Organization
- **Dedicated resource group** for shared network resources
- **Consistent naming**: `{org}-{env}-network-rg`
- **Tagging strategy**: Use metadata.org, metadata.env for consistent tagging
- **Separate subscriptions** for production vs non-production when required

---

## Troubleshooting

### DNS Records Not Resolving

**Check**:
1. Verify zone was created successfully (check stack outputs for nameservers)
2. Ensure domain registrar has correct NS records pointing to Azure nameservers
3. DNS propagation can take up to 48 hours (usually 5-10 minutes)
4. Test with: `dig @ns1-01.azure-dns.com yourdomain.com`
5. Check for typos in record names (must end with dot for FQDN)

### CNAME Conflicts

**Error**: "CNAME record conflicts with other record types"
- CNAME records cannot coexist with other record types at the same name
- CNAME cannot be created at zone apex (root domain)
- Use A records at apex, CNAME for subdomains

### Azure Service Integration Failures

**Issue**: Custom domain not validating in Azure App Service/Front Door
- Verify TXT verification records are created correctly
- Wait for DNS propagation (5-10 minutes)
- Check Azure service documentation for exact verification record format

---

## Summary

These examples demonstrate the full range of configurations supported by the `AzureDnsZone` API resource:

- **Minimal**: Quick setup for domains with dynamic records managed elsewhere
- **Standard**: Production web applications with foundational A/CNAME records
- **Comprehensive**: Complete DNS setup including email, verification, and Azure service integration
- **Wildcard**: Multi-tenant SaaS platforms with dynamic subdomains
- **Delegation**: Large organizations with decentralized DNS management

All configurations follow Azure best practices and integrate seamlessly with Azure services like AKS, App Service, Front Door, and CDN.

