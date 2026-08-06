# CivoCertificate

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `civo.planton.dev/v1alpha1`

CivoCertificateSpec defines the fields required to create a TLS certificate in Civo.

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.certificateName` | `string` | yes |  |  |
| `spec.type` | `enum` | yes |  |  |
| `spec.letsEncrypt` | `CivoCertificateLetsEncryptParams` |  |  |  |
| `spec.letsEncrypt.domains` | `[]string` | yes |  |  |
| `spec.letsEncrypt.disableAutoRenew` | `bool` |  |  |  |
| `spec.custom` | `CivoCertificateCustomParams` |  |  |  |
| `spec.custom.leafCertificate` | `string` | yes |  |  |
| `spec.custom.privateKey` | `string` (sensitive) | yes |  |  |
| `spec.custom.certificateChain` | `string` |  |  |  |
| `spec.description` | `string` |  |  |  |
| `spec.tags` | `[]string` |  |  |  |

## Field Details

### spec.certificateName

`string` · required

certificate_name is a unique, human‑readable identifier (≤ 64 chars).

- rule: {"required":true,"string":{"minLen":"1","maxLen":"64"}}

### spec.type

`enum` · required

type must align with the branch chosen in certificate_source.

- rule: {"required":true,"enum":{"definedOnly":true}}

Allowed values (use exactly as shown):

- `civo_certificate_type_unspecified` -- Unspecified type (invalid).
- `letsEncrypt` -- A free, auto‑managed Let's Encrypt certificate.
- `custom` -- A user‑provided custom certificate.

### spec.letsEncrypt

`CivoCertificateLetsEncryptParams`

### spec.letsEncrypt.domains

`[]string` · required

domains is the list of FQDNs (or wildcard domains) to include.
At least one domain is required.

- rule: {"required":true,"repeated":{"unique":true,"items":{"string":{"pattern":"^(?:\\*\\.[A-Za-z0-9\\-\\.]+|[A-Za-z0-9\\-\\.]+\\.[A-Za-z]{2,})$"}}}}

### spec.letsEncrypt.disableAutoRenew

`bool`

disable_auto_renew controls automatic renewal of the Let's Encrypt certificate.

### spec.custom

`CivoCertificateCustomParams`

### spec.custom.leafCertificate

`string` · required

leaf_certificate is the PEM‑encoded public certificate.

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.custom.privateKey

`string` · required · sensitive

private_key is the PEM‑encoded private key.

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.custom.certificateChain

`string`

certificate_chain is an optional PEM‑encoded intermediate chain.

### spec.description

`string`

Optional free‑form description (≤ 128 chars).

- rule: {"string":{"maxLen":"128"}}

### spec.tags

`[]string`

Optional tags; must be unique and lowercase kebab.

- rule: {"repeated":{"unique":true}}

## Outputs

Reference an output from another manifest as `valueFrom: {kind: CivoCertificate, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.certificate_id` | `string` | certificate_id is the unique identifier of the certificate in Civo. |
| `status.outputs.expiry_rfc3339` | `string` | expiry_rfc3339 is the expiration timestamp of the certificate in RFC 3339 format. |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
