# WAF Gateway with Path-Based Routing

This preset creates the protected microservice front door: a WAF_v2
gateway enforcing a referenced Web Application Firewall policy, with
path-based routing that sends `/api/*` to a dedicated backend pool over
end-to-end TLS and everything else to the web pool.

## When to Use

- Public applications that must sit behind OWASP protection (the WAF
  policy composes as its own `AzureWebApplicationFirewallPolicy` resource)
- One host name fronting multiple backend services split by path

## Key Configuration Choices

- **The WAF policy is a reference, not configuration** -- tuning the
  policy never redeploys the 15-25-minute gateway; one org-standard
  policy can govern many gateways
- **`PATH_BASED_ROUTING` + url path map** -- `/api/*` routes to the API
  pool with HTTPS backend settings (end-to-end encryption); unmatched
  paths fall through to the web pool's defaults
- **Per-route WAF overrides are available** -- a path rule can carry its
  own `firewallPolicyId` when one route needs stricter or looser rules
- **Pools left empty** -- both pools are joined member-side through
  `status.outputs.backend_address_pool_ids.<pool>` from NICs or scale
  sets, keeping workload membership out of the gateway spec

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<resource-group-name>` | The resource group to create the gateway in | The resource group's `status.outputs.resource_group_name` |
| `waf-gateway` | The gateway's name, unique within the resource group | Your naming convention |
| `<gateway-subnet-arm-id>` | The DEDICATED gateway subnet | The subnet's `status.outputs.subnet_id` |
| `<gateway-identity-arm-id>` | The user-assigned identity with vault access | The identity's `status.outputs.identity_id` |
| `<waf-policy-resource-name>` | The AzureWebApplicationFirewallPolicy resource | Your security manifests |
| `<public-ip-arm-id>` | A Standard static public IP | The public IP's `status.outputs.public_ip_id` |
| `<public-host-name>` | The site's host name | Your DNS zone |
| `<key-vault-certificate-secret-id>` | The certificate's versionless secret ID | The certificate's `status.outputs.versionless_secret_id` |
| `<cost-center>` | Your org's cost-attribution tag value | Your tagging convention |
