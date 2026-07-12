# CNAME-Validated Certificate

The standard managed certificate for a subdomain: Azure validates ownership through the CNAME that already routes the hostname to the app, then issues and renews the certificate for free.

## When to Use

- Any subdomain custom domain (app.example.com) whose routing record is a CNAME to the app's ingress FQDN -- the overwhelmingly common case

## Key Configuration Choices

- Publish BOTH DNS records before deploying (the deployment blocks on validation): the `asuid.{hostname}` TXT with the app's `custom_domain_verification_id`, and the CNAME to the app's `ingress_fqdn`
- Deploy the `AzureContainerAppCustomDomain` binding (certificate-less) first; Azure attaches this certificate to it once issued
- `subjectName` covers exactly one hostname -- no wildcards or extra SANs (bring your own certificate for those)

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `app.example.com` | Replace with the domain the certificate covers | The custom domain you are binding |
| `app-example-com` | Replace with the certificate's resource name (derive from the hostname) | -- |
| `<your-environment>` | The AzureContainerAppEnvironment resource's metadata name | Your Planton resource inventory |

## Related Presets

- `02-http-validated` -- for apex domains routed by A record
