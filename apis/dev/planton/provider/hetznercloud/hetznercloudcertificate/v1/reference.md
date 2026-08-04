# HetznerCloudCertificate

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `hetzner-cloud.planton.dev/v1`

HetznerCloudCertificateSpec defines the specification for a Hetzner Cloud
TLS certificate used by load balancer HTTPS services.

Hetzner Cloud supports two distinct certificate types:

  - **Uploaded**: You provide a PEM-encoded certificate chain and private key.
    Hetzner Cloud stores them and makes the certificate available for load
    balancer HTTPS listeners. You are responsible for renewal.

  - **Managed**: You specify one or more domain names and Hetzner Cloud
    automatically obtains and renews a Let's Encrypt certificate. The domains
    must have DNS records pointing to a Hetzner Cloud load balancer before
    provisioning so that the ACME HTTP-01 challenge can succeed.

Exactly one of `uploaded` or `managed` must be set. The proto oneof enforces
mutual exclusivity — it is structurally impossible to mix fields from both
types in the same manifest.

Fields not exposed in this spec (derived in IaC modules):
  - name:   Derived from metadata.name.
  - labels: Derived from metadata (CG01 pattern). Standard labels take
            precedence over user-specified metadata.labels.

## Example

```yaml
apiVersion: hetzner-cloud.planton.dev/v1
kind: HetznerCloudCertificate
metadata:
  name: hetznercloudcertificate-uploaded-demo
spec:
  uploaded:
    certificate: |
      -----BEGIN CERTIFICATE-----
      MIIBkTCB+wIJALExample...
      -----END CERTIFICATE-----
    privateKey: |
      -----BEGIN PRIVATE KEY-----
      MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQDExample...
      -----END PRIVATE KEY-----
---
apiVersion: hetzner-cloud.planton.dev/v1
kind: HetznerCloudCertificate
metadata:
  name: hetznercloudcertificate-managed-demo
spec:
  managed:
    domainNames:
      - example.com
      - www.example.com
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.uploaded` | `UploadedCertificateConfig` |  |  |  |
| `spec.uploaded.certificate` | `string` | yes |  |  |
| `spec.uploaded.privateKey` | `string` (sensitive) | yes |  |  |
| `spec.managed` | `ManagedCertificateConfig` |  |  |  |
| `spec.managed.domainNames` | `[]string` | yes |  |  |

## Field Details

### spec.uploaded

`UploadedCertificateConfig`

Upload your own TLS certificate and private key.

### spec.uploaded.certificate

`string` · required

PEM-encoded TLS certificate chain.

Must include the server certificate and, if applicable, intermediate CA
certificates in order (server cert first, root last). Changing this value
forces replacement of the certificate resource.

- rule: {"string":{"minLen":"1"}}

### spec.uploaded.privateKey

`string` · required · sensitive

PEM-encoded private key belonging to the certificate.

SENSITIVE — this field contains secret cryptographic material.
Changing this value forces replacement of the certificate resource.

- rule: {"string":{"minLen":"1"}}

### spec.managed

`ManagedCertificateConfig`

Let Hetzner Cloud obtain and manage a Let's Encrypt certificate.

### spec.managed.domainNames

`[]string` · required

Domain names for which a certificate should be obtained.

At least one domain is required. Hetzner Cloud issues a single certificate
covering all listed domains (SAN certificate). Changing this list forces
replacement of the certificate resource.

- rule: {"repeated":{"minItems":"1"}}

## Outputs

Reference an output from another manifest as `valueFrom: {kind: HetznerCloudCertificate, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.certificate_id` | `string` | The Hetzner Cloud numeric ID of the created certificate (as a string). Referenced by HetznerCloudLoadBalancer HTTPS services via StringValueOrRef. |
| `status.outputs.type` | `string` | Certificate type: "uploaded" or "managed". Computed by Hetzner Cloud. |
| `status.outputs.fingerprint` | `string` | SHA256 fingerprint of the certificate. Computed by Hetzner Cloud. |
| `status.outputs.not_valid_before` | `string` | Point in time when the certificate becomes valid (ISO-8601). Computed by Hetzner Cloud. |
| `status.outputs.not_valid_after` | `string` | Point in time when the certificate stops being valid (ISO-8601). Computed by Hetzner Cloud. |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| HetznerCloudLoadBalancer | `spec.services[].http.certificateIds` | `status.outputs.certificate_id` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
