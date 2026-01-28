# GCP Cloud CDN Examples

This document provides comprehensive examples for deploying Google Cloud CDN using OpenMCF. These examples demonstrate the 80/20 configuration principle, covering common use cases with simple configurations and advanced scenarios with detailed tuning.

## Project ID Configuration

The `gcpProjectId` field supports two configuration patterns:

### Literal Value (Direct)
```yaml
gcpProjectId:
  value: my-project-123456
```

### ValueFrom Reference (Cross-Resource)
```yaml
gcpProjectId:
  valueFrom:
    kind: GcpProject
    name: main-project
    fieldPath: status.outputs.project_id
```

This enables dynamic cross-resource references, allowing you to automatically use the project ID from a managed GcpProject resource.

## Table of Contents

1. [Basic Static Website CDN (GCS Backend)](#1-basic-static-website-cdn-gcs-backend)
2. [Static Website with Custom Domain and SSL](#2-static-website-with-custom-domain-and-ssl)
3. [CDN with Custom Cache Configuration](#3-cdn-with-custom-cache-configuration)
4. [CDN with Cache Key Optimization](#4-cdn-with-cache-key-optimization)
5. [CDN with Signed URLs for Private Content](#5-cdn-with-signed-urls-for-private-content)
6. [Cloud Run Backend with CDN](#6-cloud-run-backend-with-cdn)
7. [Compute Engine Backend with Health Checks](#7-compute-engine-backend-with-health-checks)
8. [External Origin (Hybrid/Multi-Cloud)](#8-external-origin-hybridmulti-cloud)
9. [Production CDN with Cloud Armor](#9-production-cdn-with-cloud-armor)
10. [Advanced Configuration with Negative Caching](#10-advanced-configuration-with-negative-caching)
11. [Cross-Resource Reference Example](#11-cross-resource-reference-example)

---

## 1. Basic Static Website CDN (GCS Backend)

The simplest use case: serve a static website from Google Cloud Storage with Cloud CDN enabled.

**Use case:** Personal blog, marketing site, documentation portal

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpCloudCdn
metadata:
  name: my-static-site-cdn
spec:
  gcpProjectId:
    value: my-project-123456
  backend:
    gcsBucket:
      bucketName: my-static-website-bucket
      enableUniformAccess: true
  # Optional: Using defaults for cache_mode (CACHE_ALL_STATIC), default_ttl_seconds (3600)
```

**What this creates:**
- BackendBucket pointing to GCS bucket
- Cloud CDN enabled with default cache settings
- Global HTTPS load balancer with auto-generated IP

---

## 2. Static Website with Custom Domain and SSL

Same as above, but with custom domain and Google-managed SSL certificate.

**Use case:** Production website with branded domain

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpCloudCdn
metadata:
  name: production-website-cdn
spec:
  gcpProjectId:
    value: my-company-prod
  backend:
    gcsBucket:
      bucketName: www-mycompany-com-bucket
      enableUniformAccess: true
  cacheMode: CACHE_ALL_STATIC
  defaultTtlSeconds: 3600  # 1 hour
  maxTtlSeconds: 86400     # 1 day
  enableNegativeCaching: true
  frontendConfig:
    customDomains:
      - www.mycompany.com
      - mycompany.com
    sslCertificate:
      googleManaged:
        domains:
          - www.mycompany.com
          - mycompany.com
    enableHttpsRedirect: true
```

**After deployment:**
1. Point DNS A records for `www.mycompany.com` and `mycompany.com` to the `global_ip_address` output
2. Google will automatically provision SSL certificates (takes 10-15 minutes)
3. Users accessing `http://mycompany.com` will be redirected to `https://mycompany.com`

---

## 3. CDN with Custom Cache Configuration

Fine-tune caching behavior with custom TTLs and cache mode.

**Use case:** API responses or dynamic content with specific caching requirements

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpCloudCdn
metadata:
  name: api-cdn
spec:
  gcpProjectId:
    value: my-project-123456
  backend:
    gcsBucket:
      bucketName: api-static-responses-bucket
  cacheMode: USE_ORIGIN_HEADERS  # Only cache if origin sends Cache-Control headers
  defaultTtlSeconds: 1800        # 30 minutes default
  maxTtlSeconds: 3600            # 1 hour maximum (overrides origin headers)
  clientTtlSeconds: 900          # 15 minutes for browser cache
  enableNegativeCaching: true    # Cache 404 errors
```

**When to use `USE_ORIGIN_HEADERS`:**
- Your origin (GCS bucket or backend service) sets proper `Cache-Control` headers on all responses
- You want fine-grained control per-file (different TTLs for different content)
- You're confident all cacheable content has appropriate headers

---

## 4. CDN with Cache Key Optimization

Prevent "cache shattering" by controlling what goes into the cache key.

**Use case:** Website with analytics parameters, user tracking, or session IDs in query strings

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpCloudCdn
metadata:
  name: optimized-cache-cdn
spec:
  gcpProjectId:
    value: my-project-123456
  backend:
    gcsBucket:
      bucketName: webapp-assets-bucket
  cacheMode: CACHE_ALL_STATIC
  defaultTtlSeconds: 3600
  maxTtlSeconds: 86400
  enableNegativeCaching: true
  advancedConfig:
    cacheKeyPolicy:
      includeQueryString: true
      # Only these query params affect caching; all others are ignored
      queryStringWhitelist:
        - version
        - lang
        - page
      includeProtocol: true
      includeHost: true
    enableRequestCoalescing: true
```

**Problem this solves:**

Without cache key optimization:
```
Request 1: /app.js?version=1.0&user_id=123&session=abc → Cache entry 1
Request 2: /app.js?version=1.0&user_id=456&session=xyz → Cache entry 2 (cache MISS!)
```

With query_string_whitelist = ["version"]:
```
Request 1: /app.js?version=1.0&user_id=123&session=abc → Cache entry (version=1.0)
Request 2: /app.js?version=1.0&user_id=456&session=xyz → Cache HIT! (same cache entry)
```

---

## 5. CDN with Signed URLs for Private Content

Serve paid or user-specific content with time-limited, cryptographically signed URLs.

**Use case:** Paid video courses, premium downloads, user-uploaded files

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpCloudCdn
metadata:
  name: private-content-cdn
spec:
  gcpProjectId:
    value: my-project-123456
  backend:
    gcsBucket:
      bucketName: premium-videos-bucket
  cacheMode: CACHE_ALL_STATIC
  defaultTtlSeconds: 7200  # 2 hours (videos change infrequently)
  maxTtlSeconds: 86400     # 1 day
  enableNegativeCaching: true
  advancedConfig:
    signedUrlConfig:
      enabled: true
      keys:
        - keyName: primary-key-2024
          keyValue: aGVsbG93b3JsZHRoaXNpc2Fwcm9kdWN0aW9ua2V5  # Base64-encoded 128-bit key
        - keyName: backup-key-2024
          keyValue: YmFja3Vwa2V5Zm9yY2RucmVkaXN0cmlidXRpb24  # For key rotation
```

**Generate signing keys:**
```bash
# Generate a secure 128-bit key
openssl rand -base64 16
```

**Application integration:**
Your application generates signed URLs with expiration time:
```python
# Python example (using google-auth library)
from google.auth.transport import requests
from google.auth import compute_engine
import base64
import time

def sign_url(url, key_name, key_value, expires_in=3600):
    expiration = int(time.time()) + expires_in
    url_to_sign = f"{url}?Expires={expiration}&KeyName={key_name}"
    # Create signature...
    return url_to_sign
```

---

## 6. Cloud Run Backend with CDN

Accelerate a serverless Cloud Run application with Cloud CDN.

**Use case:** API with cacheable responses, dynamic web app with static assets

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpCloudCdn
metadata:
  name: cloud-run-api-cdn
spec:
  gcpProjectId:
    value: my-project-123456
  backend:
    cloudRunService:
      serviceName: my-api-service
      region: us-central1
  cacheMode: USE_ORIGIN_HEADERS  # Cloud Run app sets Cache-Control headers
  defaultTtlSeconds: 300         # 5 minutes for uncached responses
  maxTtlSeconds: 3600            # 1 hour max
  clientTtlSeconds: 600          # 10 minutes for browsers
  enableNegativeCaching: true
  frontendConfig:
    customDomains:
      - api.mycompany.com
    sslCertificate:
      googleManaged:
        domains:
          - api.mycompany.com
    enableHttpsRedirect: true
```

**Cloud Run application must:**
1. Send `Cache-Control` headers for cacheable responses:
   ```python
   # Python Flask example
   @app.route('/api/data')
   def get_data():
       resp = make_response(jsonify(data))
       resp.headers['Cache-Control'] = 'public, max-age=300'  # Cache for 5 minutes
       return resp
   ```

2. Send `Cache-Control: no-store` for dynamic/user-specific responses:
   ```python
   @app.route('/api/user/profile')
   def get_profile():
       resp = make_response(jsonify(user_data))
       resp.headers['Cache-Control'] = 'no-store, private'  # Never cache
       return resp
   ```

---

## 7. Compute Engine Backend with Health Checks

CDN in front of Compute Engine VMs (Managed Instance Group).

**Use case:** Traditional web application on VMs, monolithic application

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpCloudCdn
metadata:
  name: webapp-cdn
spec:
  gcpProjectId:
    value: my-project-123456
  backend:
    computeService:
      instanceGroupName: webapp-mig-us-central1
      protocol: HTTP
      port: 8080
      healthCheck:
        path: /healthz
        port: 8080
        checkIntervalSeconds: 10
        timeoutSeconds: 5
        healthyThreshold: 2
        unhealthyThreshold: 3
  cacheMode: CACHE_ALL_STATIC
  defaultTtlSeconds: 3600
  maxTtlSeconds: 86400
  enableNegativeCaching: true
  advancedConfig:
    cacheKeyPolicy:
      includeQueryString: true
      queryStringWhitelist:
        - page
        - category
        - sort
      includeProtocol: true
      includeHost: true
  frontendConfig:
    customDomains:
      - www.mywebapp.com
    sslCertificate:
      googleManaged:
        domains:
          - www.mywebapp.com
    enableHttpsRedirect: true
```

**Prerequisites:**
- Managed Instance Group `webapp-mig-us-central1` must exist in the project
- VMs must respond to health checks at `/healthz` on port 8080
- VMs must serve HTTP traffic on port 8080

---

## 8. External Origin (Hybrid/Multi-Cloud)

Use Cloud CDN to accelerate content from an origin outside GCP (AWS S3, on-prem server, another cloud).

**Use case:** Gradual migration to GCP, multi-cloud architecture, hybrid deployments

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpCloudCdn
metadata:
  name: hybrid-cdn
spec:
  gcpProjectId:
    value: my-project-123456
  backend:
    externalOrigin:
      hostname: assets.mycompany.com  # Could be AWS S3, on-prem, etc.
      port: 443
      protocol: HTTPS
  cacheMode: CACHE_ALL_STATIC
  defaultTtlSeconds: 3600
  maxTtlSeconds: 86400
  enableNegativeCaching: true
  frontendConfig:
    customDomains:
      - cdn.mycompany.com
    sslCertificate:
      googleManaged:
        domains:
          - cdn.mycompany.com
    enableHttpsRedirect: true
```

**What this enables:**
- External origin (e.g., AWS S3 bucket at `assets.mycompany.com`) remains the source of truth
- Google Cloud CDN caches content at its global edge network
- End users benefit from Google's edge locations without migrating data

**Use case:** Gradual AWS → GCP migration:
1. Phase 1: Keep data in AWS S3, use GCP Cloud CDN (this example)
2. Phase 2: Migrate data to GCS, switch backend to `gcs_bucket`
3. Users see no downtime during migration

---

## 9. Production CDN with Cloud Armor

Full production deployment with WAF/DDoS protection via Cloud Armor.

**Use case:** E-commerce site, high-value application requiring security

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpCloudCdn
metadata:
  name: production-secure-cdn
spec:
  gcpProjectId:
    value: my-company-prod
  backend:
    gcsBucket:
      bucketName: ecommerce-static-assets
      enableUniformAccess: true
  cacheMode: CACHE_ALL_STATIC
  defaultTtlSeconds: 3600
  maxTtlSeconds: 86400
  clientTtlSeconds: 1800
  enableNegativeCaching: true
  advancedConfig:
    cacheKeyPolicy:
      includeQueryString: true
      queryStringWhitelist:
        - product_id
        - category
        - page
      includeProtocol: true
      includeHost: true
    enableRequestCoalescing: true
    serveWhileStaleSeconds: 60  # Serve stale content for 1 min while revalidating
  frontendConfig:
    customDomains:
      - www.mystore.com
      - mystore.com
    sslCertificate:
      googleManaged:
        domains:
          - www.mystore.com
          - mystore.com
    cloudArmor:
      enabled: true
      securityPolicyName: ecommerce-security-policy  # Pre-existing Cloud Armor policy
    enableHttpsRedirect: true
```

**Cloud Armor security policy must exist:**
```bash
# Create a Cloud Armor policy (if not already exists)
gcloud compute security-policies create ecommerce-security-policy \
  --description "WAF and DDoS protection for e-commerce site"

# Add rules (example: block SQL injection patterns)
gcloud compute security-policies rules create 1000 \
  --security-policy ecommerce-security-policy \
  --expression "evaluatePreconfiguredExpr('sqli-stable')" \
  --action deny-403
```

---

## 10. Advanced Configuration with Negative Caching

Granular control over error response caching to optimize origin load and user experience.

**Use case:** API with different error handling strategies for different status codes

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpCloudCdn
metadata:
  name: advanced-api-cdn
spec:
  gcpProjectId:
    value: my-project-123456
  backend:
    cloudRunService:
      serviceName: advanced-api
      region: us-central1
  cacheMode: USE_ORIGIN_HEADERS
  defaultTtlSeconds: 600
  maxTtlSeconds: 3600
  clientTtlSeconds: 300
  enableNegativeCaching: true
  advancedConfig:
    cacheKeyPolicy:
      includeQueryString: true
      queryStringWhitelist:
        - api_version
        - format
      includeProtocol: true
      includeHost: true
      includedHeaders:
        - Accept
        - Accept-Encoding
    negativeCachingPolicies:
      - code: 404
        ttlSeconds: 600  # Cache "Not Found" for 10 minutes (reduces origin load)
      - code: 500
        ttlSeconds: 10   # Cache "Internal Server Error" for 10 seconds only
      - code: 502
        ttlSeconds: 10   # Cache "Bad Gateway" for 10 seconds only
      - code: 503
        ttlSeconds: 30   # Cache "Service Unavailable" for 30 seconds
    serveWhileStaleSeconds: 120  # Serve stale content for 2 minutes during origin outage
    enableRequestCoalescing: true
  frontendConfig:
    customDomains:
      - api.example.com
    sslCertificate:
      googleManaged:
        domains:
          - api.example.com
    enableHttpsRedirect: true
```

**Why this configuration?**

- **404 cached for 10 minutes:** Invalid API endpoints won't hammer origin (e.g., bot scanning)
- **500/502 cached for 10 seconds:** Transient errors are cached briefly to protect origin, but not so long that legitimate users are affected
- **503 cached for 30 seconds:** Service temporarily unavailable (e.g., during deployment) - cache longer than 500 errors
- **Serve-while-stale 120 seconds:** If origin goes down, CDN serves cached content for 2 minutes, giving origin time to recover

---

## Deployment Instructions

### Using Pulumi

```bash
# Navigate to the component directory
cd apis/org/openmcf/provider/gcp/gcpcloudcdn/v1/iac/pulumi

# Set your GCP project
export GCP_PROJECT_ID=my-project-123456

# Create a manifest file (e.g., manifest.yaml) with your desired configuration
# Then deploy
planton pulumi up --manifest manifest.yaml
```

### Using Terraform

```bash
# Navigate to the Terraform module directory
cd apis/org/openmcf/provider/gcp/gcpcloudcdn/v1/iac/tf

# Initialize Terraform
terraform init

# Create a tfvars file or pass variables
terraform plan -var-file=dev.tfvars

# Apply the configuration
terraform apply -var-file=dev.tfvars
```

---

## Post-Deployment Verification

### Check CDN Status

```bash
# Get backend bucket details
gcloud compute backend-buckets describe <backend-name> --global

# Check if CDN is enabled
gcloud compute backend-buckets describe <backend-name> --global \
  --format="get(cdnPolicy.cacheMode)"
```

### Test Cache Behavior

```bash
# First request (cache MISS)
curl -I https://your-cdn-url.com/test.jpg
# Look for: X-Cache: MISS

# Second request (cache HIT)
curl -I https://your-cdn-url.com/test.jpg
# Look for: X-Cache: HIT
```

### Monitor Cache Hit Ratio

Navigate to Cloud Console → Network Services → Cloud CDN

Or query via gcloud:
```bash
gcloud logging read "resource.type=http_load_balancer
  AND jsonPayload.statusDetails=response_from_cache" \
  --limit 100 --format json
```

---

## Common Patterns

### Pattern: Versioned URLs (Best Practice)

Instead of cache invalidation, use versioned filenames:

```yaml
# Configure long TTLs
spec:
  cacheMode: CACHE_ALL_STATIC
  defaultTtlSeconds: 31536000  # 1 year
  maxTtlSeconds: 31536000      # 1 year
```

Your build process generates hashed filenames:
```
# Build output
app.a83f5b2c.js
app.f19d84e1.css
logo.9c3f4a7b.png
```

HTML references the versioned files:
```html
<script src="/app.a83f5b2c.js"></script>
<link rel="stylesheet" href="/app.f19d84e1.css">
```

**Benefit:** No cache invalidation needed. New deployment = new filenames = instant updates.

---

## Troubleshooting

### Low Cache Hit Ratio

**Symptoms:** High origin load, slow response times, high costs

**Common causes:**
1. **Cache shattering:** Too many query parameters in cache key
   - **Fix:** Use `query_string_whitelist` in `cache_key_policy`
2. **Unsupported Vary headers:** Origin sends `Vary: Cookie`
   - **Fix:** Remove unsupported Vary headers from origin (Cloud CDN only supports Vary: Accept, Accept-Encoding, Origin)
3. **Incorrect cache mode:** Using `USE_ORIGIN_HEADERS` without sending Cache-Control
   - **Fix:** Switch to `CACHE_ALL_STATIC` or ensure origin sends proper headers

### SSL Certificate Not Provisioning

**Symptoms:** `https://` doesn't work, certificate pending

**Common causes:**
1. DNS not pointed to load balancer IP
   - **Fix:** Point A record to `global_ip_address` output
2. Certificate validation in progress (takes 10-15 minutes)
   - **Fix:** Wait for Google to complete validation

### Origin Errors Not Cached

**Symptoms:** 404/500 errors hitting origin repeatedly

**Fix:** Enable `enable_negative_caching` and configure `negative_caching_policies`

---

## 11. Cross-Resource Reference Example

Use `valueFrom` to dynamically reference a GcpProject resource for the project ID.

**Use case:** Multi-environment deployments where the project ID comes from a managed GcpProject resource

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpCloudCdn
metadata:
  name: cross-ref-cdn
  org: my-org
  env:
    id: production
spec:
  # Reference the project ID from a GcpProject resource
  gcpProjectId:
    valueFrom:
      kind: GcpProject
      name: main-project
      fieldPath: status.outputs.project_id
  backend:
    gcsBucket:
      bucketName: my-static-assets-bucket
      enableUniformAccess: true
  cacheMode: CACHE_ALL_STATIC
  defaultTtlSeconds: 3600
  maxTtlSeconds: 86400
  enableNegativeCaching: true
  frontendConfig:
    customDomains:
      - cdn.example.com
    sslCertificate:
      googleManaged:
        domains:
          - cdn.example.com
    enableHttpsRedirect: true
```

**Benefits of Cross-Resource References:**
- **Dynamic Configuration**: Project ID is automatically resolved at deployment time
- **Environment Consistency**: Ensures CDN is deployed in the correct project without manual configuration
- **Dependency Management**: Creates implicit dependency between resources
- **Reduced Errors**: Eliminates typos and misconfigurations from hard-coded project IDs

---

*For more information, see the [Research Document](docs/README.md) for architectural deep-dive and best practices.*
