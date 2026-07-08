---
title: "Custom Domain WAF Association"
description: "This preset attaches a Front Door WAF policy to both the endpoint's default hostname and a validated custom domain -- the production shape once real domains serve traffic."
type: "preset"
rank: "02"
presetSlug: "02-custom-domain-association"
componentSlug: "front-door-security-policy"
componentTitle: "Front Door Security Policy"
provider: "azure"
icon: "package"
order: 2
---

# Custom Domain WAF Association

This preset attaches a Front Door WAF policy to both the endpoint's
default hostname and a validated custom domain -- the production shape
once real domains serve traffic.

## When to Use

- Production deployments whose routes serve custom domains
  (www.example.com) alongside or instead of the generated
  *.azurefd.net hostname
- Growing deployments: the `domainIds` list updates in place, so each
  new custom domain joins the association without touching the others

## Key Configuration Choices

- **Mixed ID types in one list** -- endpoint IDs protect the default
  hostname, custom-domain IDs protect their hostname; Azure accepts
  both in the same association
- **Domain caps ride the profile's sku** -- up to 100 domains on a
  STANDARD profile, 500 on PREMIUM (checked at deploy time)
- **Domains should be validated first** -- a pending custom domain can
  be associated, but the WAF only sees traffic once DNS validation
  passes and the CNAME points at Front Door

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<front-door-profile-resource-name>` | The AzureFrontDoorProfile's Planton resource name | Your Front Door composition |
| `<firewall-policy-resource-name>` | The AzureFrontDoorFirewallPolicy's Planton resource name | Your Front Door composition |
| `<front-door-endpoint-resource-name>` | The AzureFrontDoorEndpoint's Planton resource name | Your Front Door composition |
| `<custom-domain-resource-name>` | The AzureFrontDoorCustomDomain's Planton resource name | Your Front Door composition |
