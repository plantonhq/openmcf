# Civo DNS Zone Examples

This document provides real-world examples of DNS configurations using the `CivoDnsZone` resource. Each example includes the YAML manifest and explanation of the use case.

## Table of Contents

1. [Simple Website (Basic Web Hosting)](#1-simple-website-basic-web-hosting)
2. [Email-Enabled Domain (Corporate Email Setup)](#2-email-enabled-domain-corporate-email-setup)
3. [Multi-Region API (Load Distribution)](#3-multi-region-api-load-distribution)
4. [SaaS Application (Full Production Setup)](#4-saas-application-full-production-setup)
5. [Wildcard Certificate Validation (cert-manager Integration)](#5-wildcard-certificate-validation-cert-manager-integration)

---

## 1. Simple Website (Basic Web Hosting)

**Use Case:** You're hosting a simple website on a Civo VM and want both `example.com` and `www.example.com` to work.

**Requirements:**
- Root domain points to server IP
- `www` subdomain redirects to root
- Standard 1-hour TTL for easy updates

```yaml
apiVersion: civo.openmcf.org/v1
kind: CivoDnsZone
metadata:
  name: simple-website
spec:
  domainName: example.com
  records:
    # Root domain A record
    - name: "@"
      type: A
      values:
        - value: "198.51.100.42"
      ttlSeconds: 3600
    
    # www subdomain as CNAME to root
    - name: "www"
      type: CNAME
      values:
        - value: "example.com"
      ttlSeconds: 3600
```

**After deployment:**
1. Update nameservers at your registrar to `ns0.civo.com`, `ns1.civo.com`, `ns2.civo.com`
2. Wait 1-48 hours for propagation
3. Test: `curl http://example.com` and `curl http://www.example.com`

**Notes:**
- The `@` symbol represents the root domain (apex)
- CNAME records must point to a domain, not an IP address
- Never use CNAME at `@` - always use an A record for the root

---

## 2. Email-Enabled Domain (Corporate Email Setup)

**Use Case:** Configure DNS for a domain that sends and receives email (e.g., with Google Workspace, Office 365, or self-hosted mail).

**Requirements:**
- MX records for mail routing
- SPF record to prevent spoofing
- DKIM record for message signing
- DMARC policy for handling failures
- A record for mail server hostname

```yaml
apiVersion: civo.openmcf.org/v1
kind: CivoDnsZone
metadata:
  name: corporate-email
spec:
  domainName: acme-corp.com
  records:
    # Website root
    - name: "@"
      type: A
      values:
        - value: "198.51.100.100"
      ttlSeconds: 3600
    
    # Primary and backup mail servers
    - name: "@"
      type: MX
      values:
        - value: "10 mail1.acme-corp.com"
        - value: "20 mail2.acme-corp.com"
      ttlSeconds: 3600
    
    # Mail server A records
    - name: "mail1"
      type: A
      values:
        - value: "198.51.100.50"
      ttlSeconds: 3600
    
    - name: "mail2"
      type: A
      values:
        - value: "198.51.100.51"
      ttlSeconds: 3600
    
    # SPF record - authorize mail servers
    - name: "@"
      type: TXT
      values:
        - value: "v=spf1 mx ip4:198.51.100.50 ip4:198.51.100.51 ~all"
      ttlSeconds: 3600
    
    # DKIM record - message signing (example key)
    - name: "default._domainkey"
      type: TXT
      values:
        - value: "v=DKIM1; k=rsa; p=MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQC3QEKyU1fSma0axspqYK5iAj+54lsAg4qRRCnpKK68hawSJwfvEiLZ9gPiY2aPh3LFtHZLaPGOHc+EgGBUv9rJGqT1o8m+p3X9sEhWLe1JxvzWIJf1wU4SjV2/5Vz9v6r6SBaC3hLyGXcY6Wr+TZ6KLvJtF7sG5g0XTg4m3D4JFQIDAQAB"
      ttlSeconds: 3600
    
    # DMARC policy - quarantine suspicious emails
    - name: "_dmarc"
      type: TXT
      values:
        - value: "v=DMARC1; p=quarantine; rua=mailto:dmarc-reports@acme-corp.com; ruf=mailto:dmarc-failures@acme-corp.com; pct=100"
      ttlSeconds: 3600
```

**Testing email configuration:**

```bash
# Test MX records
dig MX acme-corp.com

# Validate SPF
# Use https://mxtoolbox.com/spf.aspx

# Test DKIM
# Send test email and check headers

# Verify DMARC
# Use https://mxtoolbox.com/dmarc.aspx
```

**Notes:**
- MX priority (first number) determines preference - lower is higher priority
- SPF `~all` means soft fail (allow but mark suspicious)
- DKIM public key must match private key used by mail server
- DMARC reports go to `rua` (aggregate) and `ruf` (forensic) addresses

---

## 3. Multi-Region API (Load Distribution)

**Use Case:** Distribute API traffic across multiple servers in different regions using DNS round-robin.

**Requirements:**
- Single API endpoint (`api.example.com`)
- Multiple backend servers
- Balanced traffic distribution
- Quick failover with low TTL

```yaml
apiVersion: civo.openmcf.org/v1
kind: CivoDnsZone
metadata:
  name: multi-region-api
spec:
  domainName: example.com
  records:
    # Website on CDN
    - name: "@"
      type: A
      values:
        - value: "198.51.100.10"
      ttlSeconds: 3600
    
    # API endpoint with round-robin DNS (3 regions)
    - name: "api"
      type: A
      values:
        - value: "203.0.113.10"  # US East
        - value: "203.0.113.20"  # EU West
        - value: "203.0.113.30"  # Asia Pacific
      ttlSeconds: 300  # Low TTL for faster failover
    
    # Monitoring endpoint (single server)
    - name: "monitor"
      type: A
      values:
        - value: "203.0.113.100"
      ttlSeconds: 3600
    
    # WebSocket endpoint (separate from API)
    - name: "ws"
      type: A
      values:
        - value: "203.0.113.40"
        - value: "203.0.113.41"
      ttlSeconds: 300
```

**How round-robin works:**

1. Client requests `api.example.com`
2. DNS returns all 3 IPs in random order
3. Client typically connects to the first IP
4. Each client may receive different IP order
5. Load distributes across all servers

**Notes:**
- Round-robin is **not** true load balancing - clients pick one IP
- Use low TTL (300s) for faster manual failover
- For health-checked failover, use a proper load balancer or Route 53
- Consider adding health checks via external monitoring

---

## 4. SaaS Application (Full Production Setup)

**Use Case:** Complete DNS setup for a SaaS application with marketing site, application, API, admin panel, and email.

**Requirements:**
- Marketing site on CDN
- Application and API on separate subdomains
- Admin panel on dedicated subdomain
- Email with SPF/DKIM/DMARC
- TXT records for domain verification (Google, Microsoft, etc.)

```yaml
apiVersion: civo.openmcf.org/v1
kind: CivoDnsZone
metadata:
  name: saas-production
spec:
  domainName: myapp.io
  records:
    # Marketing website on CDN
    - name: "@"
      type: A
      values:
        - value: "192.0.2.10"
      ttlSeconds: 7200  # Higher TTL for stable CDN
    
    - name: "www"
      type: CNAME
      values:
        - value: "myapp.io"
      ttlSeconds: 7200
    
    # Application (main product)
    - name: "app"
      type: A
      values:
        - value: "192.0.2.20"
        - value: "192.0.2.21"  # Redundant app servers
      ttlSeconds: 3600
    
    # API (separate infrastructure)
    - name: "api"
      type: A
      values:
        - value: "192.0.2.30"
        - value: "192.0.2.31"
        - value: "192.0.2.32"  # 3 API servers
      ttlSeconds: 600  # Low TTL for rolling updates
    
    # Admin panel (restricted access)
    - name: "admin"
      type: A
      values:
        - value: "192.0.2.40"
      ttlSeconds: 3600
    
    # Documentation site (separate hosting)
    - name: "docs"
      type: CNAME
      values:
        - value: "myapp.gitbook.io"
      ttlSeconds: 3600
    
    # Status page (external service)
    - name: "status"
      type: CNAME
      values:
        - value: "myapp.statuspage.io"
      ttlSeconds: 3600
    
    # Email configuration
    - name: "@"
      type: MX
      values:
        - value: "10 mail.myapp.io"
      ttlSeconds: 7200
    
    - name: "mail"
      type: A
      values:
        - value: "192.0.2.50"
      ttlSeconds: 7200
    
    # SPF - authorize mail servers and SendGrid
    - name: "@"
      type: TXT
      values:
        - value: "v=spf1 mx include:sendgrid.net ~all"
      ttlSeconds: 7200
    
    # DKIM (SendGrid)
    - name: "s1._domainkey"
      type: CNAME
      values:
        - value: "s1.domainkey.u12345.wl.sendgrid.net"
      ttlSeconds: 7200
    
    - name: "s2._domainkey"
      type: CNAME
      values:
        - value: "s2.domainkey.u12345.wl.sendgrid.net"
      ttlSeconds: 7200
    
    # DMARC policy
    - name: "_dmarc"
      type: TXT
      values:
        - value: "v=DMARC1; p=reject; rua=mailto:dmarc@myapp.io; pct=100"
      ttlSeconds: 7200
    
    # Domain verification (Google Workspace)
    - name: "@"
      type: TXT
      values:
        - value: "google-site-verification=ABC123XYZ789"
      ttlSeconds: 3600
    
    # Domain verification (Microsoft 365)
    - name: "@"
      type: TXT
      values:
        - value: "MS=ms12345678"
      ttlSeconds: 3600
```

**Deployment strategy:**

```bash
# 1. Apply DNS configuration
planton apply -f saas-dns.yaml

# 2. Update nameservers at registrar
# (See planton outputs for nameserver list)

# 3. Wait for propagation (check with):
dig @ns0.civo.com myapp.io
dig @ns0.civo.com api.myapp.io

# 4. Test all endpoints
curl https://myapp.io
curl https://app.myapp.io
curl https://api.myapp.io/health
```

**Notes:**
- Use higher TTL (7200s) for stable records like email
- Use lower TTL (600s) for frequently-updated records like API
- Multiple TXT records at `@` are allowed and common
- Keep admin panel DNS separate for security isolation

---

## 5. Wildcard Certificate Validation (cert-manager Integration)

**Use Case:** Automatically issue wildcard SSL certificates using cert-manager with DNS-01 challenges.

**Requirements:**
- Kubernetes cluster with cert-manager installed
- Civo DNS webhook for cert-manager
- Wildcard certificate for `*.example.com`
- Automatic renewal

**DNS Configuration:**

```yaml
apiVersion: civo.openmcf.org/v1
kind: CivoDnsZone
metadata:
  name: cert-validation
spec:
  domainName: example.com
  records:
    # Main website
    - name: "@"
      type: A
      values:
        - value: "198.51.100.42"
      ttlSeconds: 3600
    
    # Note: _acme-challenge TXT records are created automatically
    # by cert-manager webhook during certificate issuance.
    # You don't need to manually create them here.
    # 
    # If you need to test manually, you can add:
    - name: "_acme-challenge"
      type: TXT
      values:
        - value: "manual-test-token-abc123"
      ttlSeconds: 600  # Low TTL for faster validation
```

**cert-manager Certificate Resource:**

```yaml
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: wildcard-example-com
  namespace: default
spec:
  secretName: wildcard-example-com-tls
  issuerRef:
    name: letsencrypt-prod
    kind: ClusterIssuer
  dnsNames:
    - "example.com"
    - "*.example.com"
  # DNS-01 challenge will use Civo DNS webhook
```

**cert-manager ClusterIssuer with Civo DNS:**

```yaml
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: letsencrypt-prod
spec:
  acme:
    server: https://acme-v02.api.letsencrypt.org/directory
    email: ssl-admin@example.com
    privateKeySecretRef:
      name: letsencrypt-prod
    solvers:
      - dns01:
          webhook:
            groupName: acme.civo.com
            solverName: civo
            config:
              apiTokenSecretRef:
                name: civo-api-token
                key: api-token
```

**Complete workflow:**

1. Deploy DNS zone with `CivoDnsZone`
2. Install cert-manager and Civo DNS webhook in Kubernetes
3. Create `ClusterIssuer` with Civo credentials
4. Create `Certificate` resource requesting wildcard cert
5. cert-manager automatically:
   - Creates `_acme-challenge.example.com` TXT record via Civo API
   - Let's Encrypt validates ownership
   - Issues certificate
   - Stores in Kubernetes Secret
   - Cleans up TXT record
6. Certificate auto-renews before expiry (every 60 days)

**Testing validation:**

```bash
# Check if cert-manager created the challenge record
dig TXT _acme-challenge.example.com

# Watch cert-manager logs
kubectl logs -n cert-manager deploy/cert-manager -f

# Check certificate status
kubectl describe certificate wildcard-example-com

# Verify issued certificate
kubectl get secret wildcard-example-com-tls -o yaml
```

**Notes:**
- DNS-01 is required for wildcard certificates (HTTP-01 doesn't support them)
- TTL of 600s balances validation speed with cache efficiency
- cert-manager webhook needs Civo API token with DNS write permissions
- Certificates are stored as Kubernetes TLS secrets
- Automatic renewal happens ~30 days before expiry

---

## Advanced Patterns

### Multi-Environment Setup

Use different zones or subdomains for dev/staging/prod:

```yaml
# Production
domainName: example.com

# Staging
domainName: staging.example.com

# Development
domainName: dev.example.com
```

### Blue-Green Deployments

Use low TTL and switch A records during deployment:

```yaml
# Blue environment (current)
- name: "app"
  type: A
  values:
    - value: "192.0.2.10"
  ttlSeconds: 300  # Low TTL for quick switch

# After verifying green, update to:
- name: "app"
  type: A
  values:
    - value: "192.0.2.20"  # Green environment
  ttlSeconds: 300
```

### Geo-Specific Subdomains

For manual geographic routing (without advanced DNS features):

```yaml
records:
  - name: "us"
    type: A
    values:
      - value: "192.0.2.10"
    ttlSeconds: 3600
  
  - name: "eu"
    type: A
    values:
      - value: "198.51.100.10"
    ttlSeconds: 3600
  
  - name: "asia"
    type: A
    values:
      - value: "203.0.113.10"
    ttlSeconds: 3600
```

---

## Additional Resources

- [Main README](README.md) - Component overview and quick start
- [Research Documentation](docs/README.md) - Deep dive into DNS management patterns
- [Pulumi Module](iac/pulumi/README.md) - Direct Pulumi usage
- [cert-manager webhook for Civo DNS](https://github.com/okteto/cert-manager-webhook-civo)
- [ExternalDNS Civo Tutorial](https://kubernetes-sigs.github.io/external-dns/v0.14.2/tutorials/civo/)

---

## Need Help?

- Check the [Troubleshooting section](README.md#troubleshooting) in the main README
- Open an issue on [GitHub](https://github.com/plantonhq/openmcf/issues)
- Contact Civo support: support@civo.com

