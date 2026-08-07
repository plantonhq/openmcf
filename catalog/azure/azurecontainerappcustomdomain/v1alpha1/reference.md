# AzureContainerAppCustomDomain

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `azure.planton.dev/v1alpha1`

**AzureContainerAppCustomDomainSpec** binds a custom domain to a
Container App: your own hostname (app.example.com) serving the app
instead of the generated *.azurecontainerapps.io address.

**DNS must prove ownership before the binding deploys.** Azure
validates the domain during creation and the deployment fails without
these records in place:

 1. A TXT record at `asuid.{domain_name}` carrying the app's
    custom_domain_verification_id output (domain ownership), AND
 2. a CNAME from `domain_name` to the app's ingress_fqdn output (or an
    A record to the environment's static IP for apex domains).

With the domain's zone on Azure DNS, both are AzureDnsRecord resources
composed in the same deployment; with external DNS, publish them at
your provider first. The app must have ingress enabled -- a binding
has nothing to attach to otherwise.

**TLS comes one of two ways:**

- **Azure-managed certificate** (the common path): leave the
  certificate fields unset. The binding deploys certificate-less, then
  an AzureContainerAppEnvironmentManagedCertificate for the same
  hostname issues and Azure attaches it to this binding asynchronously.
  Both engines ignore Azure's out-of-band certificate attachment so it
  never reads as drift.
- **Bring-your-own certificate**: set both
  container_app_environment_certificate_id (an
  AzureContainerAppEnvironmentCertificate on the app's environment)
  and certificate_binding_type (SNI_ENABLED for TLS serving).

**All fields replace the binding when changed** -- Azure models the
binding as an entry in the app's ingress configuration with no update
surface.

## Example

```yaml
# The BYO shape exercises the binding-type enum translation (SNI_ENABLED
# must map to Azure's SniEnabled wire form) and the byo/managed resource
# dispatch in the Terraform module.
apiVersion: azure.planton.dev/v1alpha1
kind: AzureContainerAppCustomDomain
metadata:
  name: test-custom-domain
spec:
  domain_name: app.example.com
  container_app_id:
    value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.App/containerApps/test-app
  container_app_environment_certificate_id:
    value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.App/managedEnvironments/test-env/certificates/app.example.com
  certificate_binding_type: SNI_ENABLED
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.domainName` | `string` | yes |  |  |
| `spec.containerAppId` | `string \| valueFrom` | yes |  | AzureContainerApp (`status.outputs.container_app_id`) |
| `spec.containerAppEnvironmentCertificateId` | `string \| valueFrom` |  |  | AzureContainerAppEnvironmentCertificate (`status.outputs.certificate_id`) |
| `spec.certificateBindingType` | `enum` |  |  |  |

## Field Details

### spec.domainName

`string` · required

The custom hostname to bind, e.g. "app.example.com". Changing it
replaces the binding.

- rule: Domain name must be a fully qualified hostname like app.example.com -- lowercase labels of letters, digits, and hyphens; bind one concrete hostname per resource (the environment's custom DNS suffix is the wildcard mechanism)
- rule: {"required":true,"string":{"maxLen":"253"}}

### spec.containerAppId

`string | valueFrom` · required

The Container App the domain binds to, by ARM ID. References an
AzureContainerApp's container_app_id output. The app must have
ingress enabled. Changing it replaces the binding.

- references: AzureContainerApp (`status.outputs.container_app_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureContainerApp, name: <that resource's name>, fieldPath: status.outputs.container_app_id}} -- a bare string does not parse

### spec.containerAppEnvironmentCertificateId

`string | valueFrom`

The bring-your-own certificate to serve TLS with, by ARM ID.
References an AzureContainerAppEnvironmentCertificate's
certificate_id output (the certificate must live on the app's own
environment). Set together with certificate_binding_type; leave BOTH
unset for the Azure-managed-certificate flow. Changing it replaces
the binding.

- references: AzureContainerAppEnvironmentCertificate (`status.outputs.certificate_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureContainerAppEnvironmentCertificate, name: <that resource's name>, fieldPath: status.outputs.certificate_id}} -- a bare string does not parse

### spec.certificateBindingType

`enum`

How the certificate binds to the domain. Set together with
container_app_environment_certificate_id; leave BOTH unset for the
Azure-managed-certificate flow (Azure attaches the issued
certificate itself).

- rule: {"enum":{"definedOnly":true}}

Allowed values (use exactly as shown):

- `azure_container_app_custom_domain_certificate_binding_type_unspecified` -- Not specified -- the Azure-managed-certificate flow (no certificate reference; Azure attaches the issued certificate out of band).
- `SNI_ENABLED` -- Serve TLS for this hostname via SNI with the referenced certificate -- the standard choice for every bring-your-own binding.
- `DISABLED` -- Attach the domain without serving TLS from the referenced certificate. A transitional state -- browsers reject the hostname over HTTPS until a certificate binds.
- `AUTO` -- Let Azure pick the binding mode for the referenced certificate.

## Validation Rules

- `azure_container_app_custom_domain_certificate_pairing`: Set container_app_environment_certificate_id and certificate_binding_type together for a bring-your-own certificate, or leave both unset to let an Azure-managed certificate attach to the binding

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AzureContainerAppCustomDomain, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.custom_domain_id` | `string` | The binding's resource ID. Azure models the binding as an entry in the app's ingress configuration rather than a standalone ARM resource, so this is the providers' synthetic identifier: {container-app-id}/customDomainName/{domain} |
| `status.outputs.managed_certificate_id` | `string` | The ARM ID of the Azure-managed certificate attached to this binding, once one issues for the hostname. Empty for bring-your-own-certificate bindings -- and empty until Azure attaches the managed certificate, which happens asynchronously after the certificate resource issues. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.containerAppId` | AzureContainerApp | `status.outputs.container_app_id` |
| `spec.containerAppEnvironmentCertificateId` | AzureContainerAppEnvironmentCertificate | `status.outputs.certificate_id` |

## See Also

- [Overview](../README.md)
