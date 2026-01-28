# CivoCertificate Examples

This document provides real-world examples of `CivoCertificate` configurations for common use cases.

**Important Note**: The Civo provider does not currently support certificate resources in Terraform/Pulumi. These examples show the API specification structure but cannot be automatically provisioned. Use the Civo dashboard or API for actual certificate management.

## Table of Contents

1. [Single Domain Let's Encrypt](#1-single-domain-lets-encrypt)
2. [Wildcard Let's Encrypt](#2-wildcard-lets-encrypt)
3. [Multi-Domain Let's Encrypt](#3-multi-domain-lets-encrypt)
4. [Let's Encrypt Without Auto-Renewal](#4-lets-encrypt-without-auto-renewal)
5. [Custom Certificate with Chain](#5-custom-certificate-with-chain)
6. [Custom Certificate for Internal Services](#6-custom-certificate-for-internal-services)
7. [Multi-Environment Setup](#7-multi-environment-setup)

---

## 1. Single Domain Let's Encrypt

**Use Case**: Simple HTTPS for a single domain.

**Scenario**: Production website at `www.example.com`.

```yaml
apiVersion: civo.openmcf.org/v1
kind: CivoCertificate
metadata:
  name: simple-web-cert
  description: Single domain certificate for website
spec:
  certificateName: simple-web-cert
  type: letsEncrypt
  letsEncrypt:
    domains:
      - www.example.com
  tags:
    - env:prod
    - service:web
```

**Key Points**:
- Auto-renewal enabled by default
- Single domain validation via HTTP-01 challenge
- 90-day lifespan with automatic renewal

---

## 2. Wildcard Let's Encrypt

**Use Case**: Certificate covering all subdomains.

**Scenario**: Microservices architecture with multiple subdomains (`api.example.com`, `admin.example.com`, etc.).

```yaml
apiVersion: civo.openmcf.org/v1
kind: CivoCertificate
metadata:
  name: wildcard-cert
spec:
  certificateName: wildcard-cert
  type: letsEncrypt
  letsEncrypt:
    domains:
      - "*.example.com"
      - "example.com"  # Include apex if needed
  description: Wildcard certificate for all subdomains
  tags:
    - env:prod
    - type:wildcard
```

**Key Points**:
- Wildcard format: `*.domain.com`
- Requires DNS-01 challenge (Civo DNS or manual TXT record)
- Include apex domain separately if needed
- Covers infinite subdomains

---

## 3. Multi-Domain Let's Encrypt

**Use Case**: Single certificate for multiple related domains.

**Scenario**: Main website and CDN domain.

```yaml
apiVersion: civo.openmcf.org/v1
kind: CivoCertificate
metadata:
  name: multi-domain-cert
spec:
  certificateName: multi-domain-cert
  type: letsEncrypt
  letsEncrypt:
    domains:
      - example.com
      - www.example.com
      - cdn.example.com
      - assets.example.com
  description: Certificate for main and CDN domains
  tags:
    - env:prod
    - scope:multi-domain
```

**Key Points**:
- Up to 100 domains per certificate (Let's Encrypt limit)
- All domains validated independently
- Single certificate renewal for all domains

---

## 4. Let's Encrypt Without Auto-Renewal

**Use Case**: Testing or temporary certificates.

**Scenario**: Staging environment where manual control is preferred.

```yaml
apiVersion: civo.openmcf.org/v1
kind: CivoCertificate
metadata:
  name: staging-cert-no-renew
spec:
  certificateName: staging-cert-no-renew
  type: letsEncrypt
  letsEncrypt:
    domains:
      - staging.example.com
    disableAutoRenew: true
  description: Staging certificate with manual renewal
  tags:
    - env:staging
    - auto-renew:false
```

**Key Points**:
- Set `disable_auto_renew: true`
- Certificate expires after 90 days without renewal
- Must manually renew before expiration
- Useful for controlled testing environments

---

## 5. Custom Certificate with Chain

**Use Case**: Extended Validation (EV) or Organization Validation (OV) certificate.

**Scenario**: E-commerce site requiring EV certificate from commercial CA.

```yaml
apiVersion: civo.openmcf.org/v1
kind: CivoCertificate
metadata:
  name: ev-cert
spec:
  certificateName: ev-cert
  type: custom
  custom:
    leafCertificate: |
      -----BEGIN CERTIFICATE-----
      MIIFbjCCBFagAwIBAgIQDvO...
      (your EV certificate from Digicert/GlobalSign)
      -----END CERTIFICATE-----
    privateKey: |
      -----BEGIN RSA PRIVATE KEY-----
      MIIEpAIBAAKCAQEAv0...
      (private key - keep secure!)
      -----END RSA PRIVATE KEY-----
    certificateChain: |
      -----BEGIN CERTIFICATE-----
      MIIEADCCAuigAwIBAgISA...
      (intermediate certificate 1)
      -----END CERTIFICATE-----
      -----BEGIN CERTIFICATE-----
      MIIDdzCCAl+gAwIBAgIJALU...
      (intermediate certificate 2, if applicable)
      -----END CERTIFICATE-----
  description: EV certificate for e-commerce
  tags:
    - env:prod
    - cert-type:ev
    - ca:digicert
```

**Key Points**:
- Include full certificate chain (intermediates only, not root)
- Private key must be unencrypted (PEM format)
- Manual renewal responsibility
- Typically 1-2 year lifespan

---

## 6. Custom Certificate for Internal Services

**Use Case**: Internal Certificate Authority for private services.

**Scenario**: Internal API using company PKI.

```yaml
apiVersion: civo.openmcf.org/v1
kind: CivoCertificate
metadata:
  name: internal-api-cert
spec:
  certificateName: internal-api-cert
  type: custom
  custom:
    leafCertificate: |
      -----BEGIN CERTIFICATE-----
      MIIDXTCCAkWgAwIBAgIBADANBgkqhkiG9w0BA...
      (certificate from internal CA)
      -----END CERTIFICATE-----
    privateKey: |
      -----BEGIN RSA PRIVATE KEY-----
      MIIEowIBAAKCAQEA0Z...
      (private key)
      -----END RSA PRIVATE KEY-----
    certificateChain: |
      -----BEGIN CERTIFICATE-----
      MIIDXTCCAkWgAwIBAgIBADANBgkqhkiG9w0BA...
      (internal CA intermediate)
      -----END CERTIFICATE-----
  description: Internal API certificate from company PKI
  tags:
    - internal:true
    - ca:company-pki
    - service:api
```

**Key Points**:
- Internal CA certificates not publicly trusted
- Clients must trust internal CA root
- Shorter lifespans typical (30-90 days)
- Full control over issuance and revocation

---

## 7. Multi-Environment Setup

**Use Case**: Separate certificates for dev, staging, and production.

**Scenario**: SaaS application with environment-specific domains.

### Development

```yaml
apiVersion: civo.openmcf.org/v1
kind: CivoCertificate
metadata:
  name: dev-cert
spec:
  certificateName: dev-cert
  type: letsEncrypt
  letsEncrypt:
    domains:
      - dev.example.com
  description: Development environment certificate
  tags:
    - env:dev
```

### Staging

```yaml
apiVersion: civo.openmcf.org/v1
kind: CivoCertificate
metadata:
  name: staging-cert
spec:
  certificateName: staging-cert
  type: letsEncrypt
  letsEncrypt:
    domains:
      - staging.example.com
  description: Staging environment certificate
  tags:
    - env:staging
```

### Production

```yaml
apiVersion: civo.openmcf.org/v1
kind: CivoCertificate
metadata:
  name: prod-cert
spec:
  certificateName: prod-cert
  type: letsEncrypt
  letsEncrypt:
    domains:
      - example.com
      - www.example.com
  description: Production environment certificate
  tags:
    - env:prod
    - criticality:high
```

**Key Points**:
- Separate certificates per environment
- Different domain names prevent conflicts
- Production uses apex + www
- Dev/staging use subdomains

---

## Best Practices Summary

1. **Use Let's Encrypt by default** - Free, automated, secure
2. **Enable auto-renewal** - Don't disable unless you have a specific reason
3. **Include certificate chain** - Always for custom certificates
4. **Monitor expiration** - Even with auto-renewal, things can fail
5. **Separate environments** - Don't reuse production certs in dev/staging
6. **Secure private keys** - Never commit to Git, use secret managers
7. **Test with staging** - Use Let's Encrypt staging for testing to avoid rate limits

---

## Current Provider Limitations

⚠️ **Important**: These examples show valid API specifications, but actual provisioning requires manual steps or Civo API calls until provider support is added.

**To provision certificates manually**:
1. Use Civo dashboard: https://dashboard.civo.com/
2. Use Civo API: https://www.civo.com/api/certificates
3. Use civo CLI: `civo certificate create`

**Check for provider updates**:
- Terraform: https://registry.terraform.io/providers/civo/civo/latest/docs
- Pulumi: https://www.pulumi.com/registry/packages/civo/

---

## Related Documentation

- **API Reference**: [README.md](./README.md)
- **Research**: [docs/README.md](./docs/README.md)
- **Pulumi Module**: [iac/pulumi/](./iac/pulumi/)
