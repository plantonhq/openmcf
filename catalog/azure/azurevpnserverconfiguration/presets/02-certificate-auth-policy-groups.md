# Certificate Auth with Policy Groups

This preset authenticates clients against your own root certificate and segments them into policy groups by the certificate's common name -- engineering and contractors here. A point-to-site gateway maps each group to its own address pool, so network policy can tell the two populations apart.

## When to Use

- Unmanaged or offline-capable endpoints (labs, contractors, appliances) that cannot use Entra ID sign-in
- Rollouts that need different firewall/routing treatment per user population

## Key Configuration Choices

- **Certificate authentication** -- clients present certificates chaining to `corp-root`; you own issuance and revocation (`clientRevokedCertificates` blocks individual thumbprints in place)
- **Policy groups match the certificate CN** -- `CertificateGroupId` compares the client certificate's common name; the group's name is the key the `policy_group_ids` output and the gateway's mapping use
- **One default group** -- `isDefault: true` catches members matching no other group (fixed at group creation)
- **OpenVPN** -- policy groups require the OpenVPN tunnel protocol

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-resource-group-name>` | Resource group for the configuration object | `AzureResourceGroup` status outputs (`resource_group_name`), or reference it with valueFrom |
| `<your-root-certificate-base64-body>` | The root certificate's base64 X.509 body (PEM payload without the BEGIN/END lines or line breaks) | Your certificate authority; for a self-signed root: `openssl x509 -in root.pem -outform der \| base64` |

Issue client certificates with the common names the policy groups match ("engineering", "contractor").
