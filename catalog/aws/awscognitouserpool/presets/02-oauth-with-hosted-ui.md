# OAuth with Hosted UI

**Rank**: 2

## When to Use

- Web applications using the Cognito-hosted sign-in page
- OAuth 2.0 Authorization Code flow with hosted `/oauth2/*` endpoints
- Teams that want login pages without building auth UI

## What It Provides

- Email sign-in with a strengthened password policy
- A Cognito-hosted prefix domain (`{name}-auth.auth.{region}.amazoncognito.com`) serving the hosted UI and OAuth2 endpoints
- Password recovery via verified email

## What You Might Add

- An `AwsCognitoUserPoolClient` with `allowedOauthFlows: [code]`, callback/logout URLs, and `allowedOauthScopes: [openid, email, profile]` -- the OAuth contract lives on the client resource
- A custom domain: change `domain.domain` to your FQDN and set `domain.certificateArn` (ACM certificate in us-east-1)
- `AwsCognitoIdentityProvider` resources for Google/OIDC/SAML federation on the hosted UI
