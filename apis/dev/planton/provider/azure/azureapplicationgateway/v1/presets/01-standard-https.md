# Standard HTTPS Gateway with Key Vault TLS

This preset creates the production L7 baseline: a zone-redundant
Standard_v2 gateway that terminates TLS with a Key Vault certificate,
redirects all HTTP to HTTPS, autoscales 2-10 instances, health-checks
backends on `/healthz`, and pins the Microsoft strict TLS policy
(TLS 1.2+, strong ciphers).

## When to Use

- The standard entry point for web applications on VMs, scale sets, or
  AKS behind one public host name
- Any workload that needs SSL termination with certificates that renew
  without touching the gateway

## Key Configuration Choices

- **Key Vault certificate by versionless secret ID** -- referencing an
  `AzureKeyVaultCertificate`'s `versionless_secret_id` output means
  renewals propagate automatically; the user-assigned identity must hold
  GET on the vault's secrets BEFORE the gateway deploys
- **HTTP -> HTTPS as a redirect rule** -- priority 10 beats the HTTPS
  rule's 100, and `PERMANENT` (301) lets clients cache it
- **Autoscale over fixed capacity** -- 2 minimum for availability, 10
  ceiling for cost control; the two are mutually exclusive
- **`AppGwSslPolicy20220101S`** -- Microsoft's current strict predefined
  policy; prefer it over hand-rolled cipher lists
- **Member-side pool joining** -- leave `ipAddresses` empty and have NICs
  or scale sets reference `status.outputs.backend_address_pool_ids.web`

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<resource-group-name>` | The resource group to create the gateway in | The resource group's `status.outputs.resource_group_name` |
| `<gateway-name>` | The gateway's name, unique within the resource group | Your naming convention |
| `<gateway-subnet-arm-id>` | The DEDICATED gateway subnet (/24 recommended) | The subnet's `status.outputs.subnet_id` |
| `<gateway-identity-arm-id>` | The user-assigned identity with vault access | The identity's `status.outputs.identity_id` |
| `<public-ip-arm-id>` | A Standard static public IP | The public IP's `status.outputs.public_ip_id` |
| `<backend-ip>` | A backend server address (or join member-side) | Your workload manifests |
| `<public-host-name>` | The site's host name (e.g. www.contoso.com) | Your DNS zone |
| `<key-vault-certificate-secret-id>` | The certificate's versionless secret ID | The certificate's `status.outputs.versionless_secret_id` |
| `<cost-center>` | Your org's cost-attribution tag value | Your tagging convention |
