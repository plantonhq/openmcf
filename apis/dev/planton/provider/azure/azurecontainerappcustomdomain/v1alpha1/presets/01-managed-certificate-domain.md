# Managed-Certificate Domain

The standard custom-domain flow: bind the hostname certificate-less, then let a free Azure-managed certificate attach automatically once issued.

## When to Use

- Any app custom domain where standard domain-validated TLS is enough -- the default choice

## Key Configuration Choices

- Publish the DNS records BEFORE deploying (Azure validates during creation): the `asuid.{hostname}` TXT with the app's `custom_domain_verification_id` output, and the CNAME to the app's `ingress_fqdn` output -- both compose as `AzureDnsRecord` resources when the zone is on Azure DNS
- Follow this binding with an `AzureContainerAppEnvironmentManagedCertificate` for the same hostname
- The certificate fields stay unset deliberately -- Azure fills the binding in asynchronously, and the modules treat that attachment as expected state, never drift

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `app.example.com` | Replace with the domain to bind | Your DNS zone |
| `<your-app>` | The AzureContainerApp resource's metadata name | Your Planton resource inventory |

## Related Presets

- `02-byo-certificate-domain` -- serve TLS from a certificate you bring
