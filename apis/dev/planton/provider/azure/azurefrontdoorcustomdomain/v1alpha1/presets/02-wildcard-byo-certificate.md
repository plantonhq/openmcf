# Wildcard Domain with Bring-Your-Own Certificate

This preset creates a wildcard custom domain served by a customer
certificate wrapped in an AzureFrontDoorSecret -- the shape for
multi-tenant platforms where every tenant gets a subdomain.

## When to Use

- Wildcard hostnames (`*.example.com`) -- Azure's managed certificates
  do not cover them, so a BYO certificate is mandatory
- EV/OV certificates or an org-mandated CA on any hostname

## Key Configuration Choices

- **`certificateType: CUSTOMER_CERTIFICATE` pairs with `secretId`** --
  the spec enforces the pairing in both directions (a managed domain
  must not reference a secret)
- **The secret, not the vault, is referenced** -- the
  AzureFrontDoorSecret wraps the Key Vault certificate and owns the
  rotation posture (versionless = follows renewals automatically)
- **One wildcard certificate serves many domains** -- point every
  tenant subdomain's `secretId` at the same secret; rotation is then
  one Key Vault operation

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<front-door-profile-resource-name>` | The AzureFrontDoorProfile's Planton resource name | Your Front Door composition |
| `<front-door-secret-resource-name>` | The AzureFrontDoorSecret wrapping your wildcard certificate | Your certificate composition |
| `hostName` (example value) | Your real wildcard hostname | Your DNS zone |

## Downstream Wiring

Validate the domain (TXT record from `validation_token`), then attach it
to routes; disable `linkToDefaultDomain` on routes that should answer
only on the custom domain.
