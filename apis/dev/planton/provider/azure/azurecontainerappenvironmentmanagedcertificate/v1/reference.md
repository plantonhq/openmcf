# AzureContainerAppEnvironmentManagedCertificate

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `azure.planton.dev/v1`

**AzureContainerAppEnvironmentManagedCertificateSpec** defines a TLS
certificate that Azure provisions and renews end to end for one custom
domain -- free, domain-validated, and rotated by Azure before expiry.
The managed certificate lives on the Container App Environment and is
bound to an app's domain through AzureContainerAppCustomDomain.

**Deployment blocks on domain validation.** Azure only issues the
certificate after proving you control the domain, and the create
operation polls until that proof succeeds (or times out around 30
minutes). The required DNS records must therefore exist BEFORE this
resource deploys:

 1. A TXT record at `asuid.{subject_name}` carrying the app's
    custom_domain_verification_id output (domain ownership), AND
 2. the routing record for `subject_name` itself -- a CNAME to the
    app's ingress_fqdn (for CNAME validation), or an HTTP-reachable
    A-record path (for HTTP validation).

With the domain's zone on Azure DNS, both records are AzureDnsRecord
resources composed in the same deployment; with external DNS, publish
them at your provider first.

The typical sequence: bind the domain with
AzureContainerAppCustomDomain first (certificate-less), then deploy
this certificate for the same subject_name -- Azure attaches it to the
binding asynchronously once issued.

**ForceNew fields**: everything except `tags` replaces the certificate
(Azure re-issues rather than mutating managed certificates).

## Example

```yaml
# CNAME validation exercises the domain_control_validation enum
# translation (unspecified would silently deploy HTTP on both engines).
apiVersion: azure.planton.dev/v1
kind: AzureContainerAppEnvironmentManagedCertificate
metadata:
  name: test-managed-certificate
spec:
  certificate_name: app-example-com
  container_app_environment_id:
    value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.App/managedEnvironments/test-env
  subject_name: app.example.com
  domain_control_validation: CNAME
  tags:
    team: platform
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.certificateName` | `string` | yes |  |  |
| `spec.containerAppEnvironmentId` | `string \| valueFrom` | yes |  | AzureContainerAppEnvironment (`status.outputs.environment_id`) |
| `spec.subjectName` | `string` | yes |  |  |
| `spec.domainControlValidation` | `enum` |  |  |  |
| `spec.tags` | `map<string, string>` |  |  |  |

## Field Details

### spec.certificateName

`string` · required

The certificate's name on the environment. Lowercase letters, digits,
hyphens, and dots (commonly derived from the domain, e.g.
"app-example-com"). Changing it replaces the certificate.

- rule: Certificate name must use lowercase letters, digits, hyphens, and dots, start and end with a letter or digit, and avoid consecutive hyphens -- e.g. app-example-com
- rule: {"required":true}

### spec.containerAppEnvironmentId

`string | valueFrom` · required

The Container App Environment the certificate is provisioned on, by
ARM ID. References an AzureContainerAppEnvironment's environment_id
output. Changing it replaces the certificate.

- references: AzureContainerAppEnvironment (`status.outputs.environment_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureContainerAppEnvironment, name: <that resource's name>, fieldPath: status.outputs.environment_id}} -- a bare string does not parse

### spec.subjectName

`string` · required

The domain the certificate is issued for (e.g. "app.example.com") --
the same hostname the AzureContainerAppCustomDomain binding uses.
Azure managed certificates cover exactly one name: no wildcards, no
additional SANs (bring your own certificate for those). Changing it
replaces the certificate.

- rule: Subject name must be a fully qualified domain name like app.example.com -- lowercase labels of letters, digits, and hyphens; wildcards are not supported by Azure managed certificates (use AzureContainerAppEnvironmentCertificate to bring a wildcard certificate)
- rule: {"required":true,"string":{"maxLen":"253"}}

### spec.domainControlValidation

`enum`

How Azure proves you control the domain before issuing:

- HTTP (the default): Azure serves a token over HTTP on the domain --
  works when the domain already routes to the app (an A record or
  apex domain).
- CNAME: Azure checks that the domain CNAMEs to the app's ingress
  FQDN -- the standard choice for subdomains.

Either way, the asuid TXT record must also be in place. Changing it
replaces the certificate.

- rule: {"enum":{"definedOnly":true}}

Allowed values (use exactly as shown):

- `azure_container_app_managed_certificate_domain_control_validation_unspecified` -- Not specified -- deploys HTTP, Azure's default.
- `HTTP` -- Azure serves a validation token over HTTP on the domain. Use when the domain already routes to the app (A record or apex).
- `CNAME` -- Azure checks the domain CNAMEs to the app's ingress FQDN. The standard choice for subdomains.

### spec.tags

`map<string, string>`

Free-form Azure resource tags applied to the certificate, merged over
the platform's metadata-derived tags (user tags win on key collision).
The only field Azure updates in place.

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AzureContainerAppEnvironmentManagedCertificate, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.certificate_id` | `string` | The Azure Resource Manager ID of the managed certificate. Format: /subscriptions/{sub}/resourceGroups/{rg}/providers/Microsoft.App/managedEnvironments/{env}/managedCertificates/{name} |
| `status.outputs.validation_token` | `string` | The domain-validation token Azure generated for this certificate. Relevant while validation is pending; once the certificate issues (which the deployment waits for), it is informational. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.containerAppEnvironmentId` | AzureContainerAppEnvironment | `status.outputs.environment_id` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
