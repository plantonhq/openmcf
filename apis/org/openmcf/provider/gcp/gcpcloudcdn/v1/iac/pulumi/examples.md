Here are a few examples for the `GcpCloudCdn` API resource, showcasing different configurations that demonstrate basic setups, CDN caching, and security configurations.

---

### Basic Example with Literal Project ID

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpCloudCdn
metadata:
  name: basic-cloud-cdn
spec:
  gcpProjectId:
    value: my-gcp-project
  backend:
    gcsBucket:
      bucketName: my-static-website-bucket
      enableUniformAccess: true
```

---

### Example with Cross-Resource Reference

Use `valueFrom` to dynamically reference a GcpProject resource:

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpCloudCdn
metadata:
  name: cross-ref-cloud-cdn
spec:
  gcpProjectId:
    valueFrom:
      kind: GcpProject
      name: main-project
      fieldPath: status.outputs.project_id
  backend:
    gcsBucket:
      bucketName: my-static-website-bucket
      enableUniformAccess: true
```

---

### Example with Custom Domain and SSL

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpCloudCdn
metadata:
  name: production-cdn
spec:
  gcpProjectId:
    value: my-production-project
  backend:
    gcsBucket:
      bucketName: prod-assets-bucket
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

---

### Note on ValueFrom References

**Current limitation**: Reference resolution is not yet fully implemented. Only literal `value` is used in the current implementation. References will be resolved once the shared reference resolution library is completed.

**Future work**: Implement reference resolution in a shared library that all modules can use.
