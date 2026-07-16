# Rotating Bring-Your-Own Certificate

This preset creates a Front Door secret wrapping a Key Vault
certificate by its VERSIONLESS id -- the rotation-follows-latest
posture: when Key Vault renews or you upload a new version, Front Door
picks it up automatically and every domain using this secret rotates
with zero redeploys.

## When to Use

- The default for every BYO certificate -- wildcard certs for
  multi-tenant platforms, EV/OV certs, org-CA certs
- Paired with an auto-renewing AzureKeyVaultCertificate (issuer-managed
  or self-signed with an AUTO_RENEW lifetime action), which makes the
  whole TLS chain hands-off

## Key Configuration Choices

- **`versionless_id`, not `certificate_id`** -- the versionless
  reference is what makes rotation propagate; see the pinned-version
  preset for the opposite trade-off
- **The secret is immutable** -- any field change replaces it; that is
  fine, because rotation happens in Key Vault, not here

## One-Time Prerequisite (per tenant)

Front Door reads Key Vault with Microsoft's own service principal (the
`Microsoft.AzureFrontDoor-Cdn` enterprise application). Grant it read
access on the vault -- e.g. the "Key Vault Secrets User" role on an
RBAC-mode vault -- before the first secret deploys.

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<front-door-profile-resource-name>` | The AzureFrontDoorProfile's Planton resource name | Your Front Door composition |
| `<key-vault-certificate-resource-name>` | The AzureKeyVaultCertificate holding the certificate | Your certificate composition |

## Downstream Wiring

Custom domains terminate TLS with this secret:

```yaml
# On an AzureFrontDoorCustomDomain
tls:
  certificateType: CUSTOMER_CERTIFICATE
  secretId:
    valueFrom:
      kind: AzureFrontDoorSecret
      name: my-wildcard-cert-secret
      fieldPath: status.outputs.secret_id
```
