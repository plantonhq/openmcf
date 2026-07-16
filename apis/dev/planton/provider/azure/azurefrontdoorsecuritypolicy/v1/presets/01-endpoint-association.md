# Endpoint WAF Association

This preset attaches a Front Door WAF policy to an endpoint's default
`*.azurefd.net` hostname -- the association that actually turns the WAF
on. Without a security policy, a WAF policy sits idle.

## When to Use

- Immediately after deploying an AzureFrontDoorFirewallPolicy -- this
  is the second half of every Front Door WAF rollout
- Deployments serving traffic on the endpoint's generated hostname
  (custom domains join the same list as they come online)

## Key Configuration Choices

- **The endpoint reference protects the default domain** -- Azure
  accepts endpoint IDs and custom-domain IDs interchangeably in
  `domainIds`, and one list can mix both
- **The WAF applies to all paths (`/*`)** -- Azure's security policies
  accept no other pattern; scope enforcement by choosing WHICH domains
  to associate
- **The skus must match** -- a STANDARD profile pairs with a STANDARD
  WAF policy, PREMIUM with PREMIUM; Azure rejects a mismatch at deploy
  time

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<front-door-profile-resource-name>` | The AzureFrontDoorProfile's Planton resource name | Your Front Door composition |
| `<firewall-policy-resource-name>` | The AzureFrontDoorFirewallPolicy's Planton resource name | Your Front Door composition |
| `<front-door-endpoint-resource-name>` | The AzureFrontDoorEndpoint's Planton resource name | Your Front Door composition |
