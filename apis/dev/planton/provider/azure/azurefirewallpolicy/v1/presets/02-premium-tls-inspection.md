# Premium TLS Inspection + IDPS

This preset creates the PREMIUM-tier policy for regulated or
high-security environments: outbound TLS is decrypted, inspected, and
re-encrypted using your intermediate CA from Key Vault, and the IDPS
engine blocks traffic matching known attack signatures.

## When to Use

- Compliance regimes requiring outbound content inspection
- Environments where URL (path-level) filtering or web categories are
  needed -- both require TLS termination for HTTPS traffic
- Networks that want signature-based intrusion prevention at the egress
  chokepoint

## Key Configuration Choices

- **`sku: PREMIUM`** -- TLS inspection and IDPS are Premium features;
  the tier is fixed at creation and every attached firewall must also be
  PREMIUM
- **`identity` + `tlsCertificate`** -- the firewall reads the CA's
  secret face from Key Vault through the user-assigned identity; grant
  it "Key Vault Secrets User" on the vault. The versionless secret
  reference (the default) follows CA renewals automatically
- **`intrusionDetection.mode: IDPS_DENY`** -- start at IDPS_ALERT while
  tuning, then move to deny; per-signature overrides handle false
  positives without lowering the engine mode

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<resource-group-name>` | The resource group to create the policy in | The resource group's `status.outputs.resource_group_name` |
| `<tls-identity-name>` | The AzureUserAssignedIdentity reading the CA | That identity's Planton resource name |
| `<ca-certificate-name>` | The AzureKeyVaultCertificate holding the CA | That certificate's Planton resource name |
| `<cost-center>` | Your org's cost-attribution tag value | Your tagging convention |

## Downstream Wiring

The attached firewall must also be PREMIUM:

```yaml
# AzureFirewall
spec:
  skuTier: PREMIUM
  firewallPolicyId:
    valueFrom:
      name: inspected-egress
```
