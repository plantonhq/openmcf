# AzureFrontDoorSecret

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `azure.planton.dev/v1alpha1`

**AzureFrontDoorSecretSpec** defines the configuration for creating a
secret inside an Azure Front Door (Standard/Premium) profile -- the
bring-your-own TLS certificate node. A secret wraps a Key Vault
certificate so custom domains can terminate TLS with it: the domain's
tls.secret_id references this secret, and this secret references the
AzureKeyVaultCertificate that actually holds the key material.

The secret is a first-class resource (not a field on the domain)
because one certificate -- typically a wildcard or multi-SAN cert --
serves many custom domains, and rotating it must be a single
operation, not a per-domain edit.

**The whole resource is immutable**: Azure exposes no update on Front
Door secrets, so changing ANY field replaces the secret. That is safe
in practice -- certificate ROTATION happens inside Key Vault (new
certificate versions), not by editing the secret, when the reference
is versionless (see key_vault_certificate_id).

**Operational prerequisite (one-time, per tenant)**: Front Door reads
Key Vault with its own service principal (the "Microsoft.AzureFrontDoor-Cdn"
enterprise application). That principal must be granted read access to
the vault's certificates and secrets -- e.g. the "Key Vault Secrets
User" role on an RBAC-mode vault -- before a secret can deploy. This
is a grant on YOUR vault to Microsoft's principal; no module can
create the principal itself.

**Certificate content requirement**: Azure rejects SELF-SIGNED
certificates ("the certificate chain includes an invalid number of
certificates") -- the wrapped certificate must be CA-issued with a
complete chain (leaf plus issuer, at least two certificates). Use a
Key Vault certificate enrolled through a CA integration or an
imported PKCS#12 that carries its full chain.

## Example

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureFrontDoorSecret
metadata:
  name: test-front-door-secret
spec:
  profileId:
    value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.Cdn/profiles/test-frontdoor
  secretName: wildcard-example-com
  # Versionless certificate id -- Front Door follows the latest version
  # (rotation propagates without redeploying).
  keyVaultCertificateId:
    value: https://test-vault.vault.azure.net/certificates/wildcard-example-com
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.profileId` | `string \| valueFrom` | yes |  | AzureFrontDoorProfile (`status.outputs.profile_id`) |
| `spec.secretName` | `string` | yes |  |  |
| `spec.keyVaultCertificateId` | `string \| valueFrom` | yes |  | AzureKeyVaultCertificate (`status.outputs.versionless_id`) |

## Field Details

### spec.profileId

`string | valueFrom` · required

The Front Door profile the secret lives in, by ARM ID. References
an AzureFrontDoorProfile's profile_id output so the profile and its
secrets compose in one manifest set. Fixed at creation.

- references: AzureFrontDoorProfile (`status.outputs.profile_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureFrontDoorProfile, name: <that resource's name>, fieldPath: status.outputs.profile_id}} -- a bare string does not parse

### spec.secretName

`string` · required

The secret's name -- unique within the profile. Custom domains
reference the secret by ARM ID (see the secret_id output), so the
name is mostly a human-facing label in the portal.

2-260 characters; letters, digits, and hyphens; must start and end
with a letter or digit.

- rule: secret_name must be 2-260 characters, start and end with a letter or digit, and contain only letters, digits, and hyphens
- rule: {"required":true}

### spec.keyVaultCertificateId

`string | valueFrom` · required

The Key Vault CERTIFICATE the secret wraps, by its Key Vault
certificate identifier (a vault data-plane URL, not an ARM ID).
References an AzureKeyVaultCertificate's versionless_id output --
the versionless form (no trailing version segment) tells Front Door
to follow the certificate's LATEST version, so Key Vault rotation
and auto-renewal propagate to every domain using this secret
without redeploying anything. Reference the versioned certificate_id
output instead only when a domain must pin one exact certificate
version (rotation then requires replacing this secret).

- references: AzureKeyVaultCertificate (`status.outputs.versionless_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureKeyVaultCertificate, name: <that resource's name>, fieldPath: status.outputs.versionless_id}} -- a bare string does not parse

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AzureFrontDoorSecret, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.secret_id` | `string` | The Azure Resource Manager ID of the secret -- what AzureFrontDoorCustomDomain's tls.secret_id references to terminate TLS with the wrapped certificate. Format: /subscriptions/{sub}/resourceGroups/{rg}/providers/Microsoft.Cdn/profiles/{profile}/secrets/{name} |
| `status.outputs.secret_name` | `string` | The secret's name -- unique within its profile. |
| `status.outputs.subject_alternative_names` | `[]string` | The DNS names (subject + subject alternative names) the wrapped certificate covers -- read back from the certificate so operators can confirm a domain's host_name is actually covered before attaching the secret to it. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.profileId` | AzureFrontDoorProfile | `status.outputs.profile_id` |
| `spec.keyVaultCertificateId` | AzureKeyVaultCertificate | `status.outputs.versionless_id` |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| AzureFrontDoorCustomDomain | `spec.tls.secretId` | `status.outputs.secret_id` |

## See Also

- [Overview](../README.md)
