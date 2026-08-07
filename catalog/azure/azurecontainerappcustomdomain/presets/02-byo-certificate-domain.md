# BYO-Certificate Domain

A custom domain served with a certificate you brought: an EV/OV chain, an org-mandated CA, or a wildcard certificate stored on the app's environment.

## When to Use

- Compliance-mandated certificate authorities or extended-validation chains
- A wildcard environment certificate covering many app hostnames (managed certificates cannot do wildcards)

## Key Configuration Choices

- The certificate must live on the SAME environment as the app, and the hostname must match its CN or a SAN
- `certificateBindingType: SNI_ENABLED` is the value for serving TLS; the certificate id and binding type always travel together
- The DNS prerequisites are unchanged from the managed flow: the `asuid.{hostname}` TXT and the routing CNAME must resolve publicly before deployment

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `app.example.com` | Replace with the domain to bind | Your DNS zone (must match the certificate's CN/SAN) |
| `<your-app>` | The AzureContainerApp resource's metadata name | Your Planton resource inventory |
| `<your-certificate>` | The AzureContainerAppEnvironmentCertificate resource's metadata name | Your Planton resource inventory |

## Related Presets

- `01-managed-certificate-domain` -- the free Azure-managed flow
