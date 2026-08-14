# Vaulted Credentials

This preset vaults the two credential shapes agent tools most commonly
need — an API key and a GitHub OAuth client — as managed-secret
references resolved just-in-time at deploy. Gateway targets and harness
tools reference the provider ARNs; the values never leave the vault.

## When to Use

- Gateway targets calling API-key or OAuth-protected backends
- Replacing credentials scattered through agent code with one rotation
  point

## What You Get

- AWS stores each value in Secrets Manager under the service token
  vault; the output maps carry both the provider ARNs (for consumers)
  and the secret ARNs (for audits)
- Rotation = re-apply with a new secret reference; consumers keep the
  same provider ARN

## Customize

- `vendor: CUSTOM` with `oauthDiscovery.discoveryUrl` for any OIDC
  provider; spell out `authorizationServerMetadata` when the vendor has
  no discovery document
- Add `workloadIdentities` when agents present their own identity in
  user-delegated OAuth flows
