---
title: "OIDC-Protected HTTPS"
description: "This preset creates an HTTPS listener that requires a login before any request reaches a target: an `authenticate-oidc` action redirects unauthenticated browsers to your identity provider (Okta,..."
type: "preset"
rank: "03"
presetSlug: "03-oidc-protected"
componentSlug: "lb-listener"
componentTitle: "LB Listener"
provider: "aws"
icon: "package"
order: 3
---

# OIDC-Protected HTTPS

This preset creates an HTTPS listener that requires a login before any
request reaches a target: an `authenticate-oidc` action redirects
unauthenticated browsers to your identity provider (Okta, Auth0, Entra ID,
Google, or any OIDC-compliant IdP), then the forward action delivers the
authenticated request. The ALB manages the session cookie afterward -- the
application never implements a login flow.

## When to Use

- Internal tools, dashboards, and admin UIs that should sit behind SSO
- Putting authentication in front of applications that have none
- Enforcing one login policy at the edge instead of per application

## Key Configuration Choices

- **Action order matters** -- the authenticate action is listed before the
  forward; the chain runs top to bottom and must end in a terminal action
- **`clientSecret` is sensitive end to end** -- the field is masked in plans,
  logs, and the console; replace the placeholder with the real secret when
  deploying
- **`scope: openid email`** -- requests the identity token plus the email
  claim, which the ALB passes to targets in `x-amzn-oidc-data`; trim to
  `openid` if targets do not need it
- **Default session length** -- 7 days (`sessionTimeoutSeconds` unset); set
  it lower for sensitive tools

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<aws-region>` | AWS region code (e.g., `us-east-1`) | Your deployment region |
| `<alb-resource-name>` | Name of the AwsAlb resource to attach to | Your AwsAlb manifest's `metadata.name` |
| `<certificate-resource-name>` | Name of the AwsCertManagerCert for the domain | Your AwsCertManagerCert manifest's `metadata.name` |
| `<oidc-issuer-url>` | Issuer identifier (e.g., `https://accounts.google.com`) | Your IdP's OIDC discovery document (`/.well-known/openid-configuration`) |
| `<oidc-authorization-endpoint>` | Authorization endpoint URL | Same discovery document |
| `<oidc-token-endpoint>` | Token endpoint URL | Same discovery document |
| `<oidc-user-info-endpoint>` | UserInfo endpoint URL | Same discovery document |
| `<oidc-client-id>` | OAuth client ID registered with the IdP | Your IdP's application registration |
| `replace-with-client-secret` | OAuth client secret | Your IdP's application registration |
| `<target-group-resource-name>` | Name of the AwsLbTargetGroup receiving traffic | Your AwsLbTargetGroup manifest's `metadata.name` |

## Common Additions

- Register the callback URL `https://<your-domain>/oauth2/idpresponse` with
  the IdP -- the ALB handles that path itself
- Use `authenticate-cognito` with an `AwsCognitoUserPool` reference instead
  when Cognito is the IdP
- For APIs (no browser), use a `jwt-validation` action instead -- stateless
  bearer-token checks with no redirects or cookies

## Related Presets

- **01-https-forward** -- the same listener without the login step
- **02-http-redirect-to-https** -- the port-80 companion redirect
