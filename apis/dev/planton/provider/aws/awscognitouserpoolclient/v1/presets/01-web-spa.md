# Preset: Web SPA

**Rank**: 1 (most common starting point)

## When to Use

- Single-page applications and mobile apps (public clients -- no secret)
- Authorization Code + PKCE sign-in, with or without the hosted UI

## What It Provides

- Authorization Code grant with the standard OIDC scopes (`openid`, `email`, `profile`)
- SRP authentication (the password never leaves the device) plus refresh tokens
- Token revocation on sign-out and user-enumeration protection
- The pool's own directory (`COGNITO`) as the sign-in provider

## What You Might Add

- Federated providers in `supportedIdentityProviders` -- reference `AwsCognitoIdentityProvider` resources
- `refreshTokenRotation: {feature: ENABLED}` to make refresh tokens one-time-use (drop `ALLOW_REFRESH_TOKEN_AUTH` from `explicitAuthFlows` when you do -- AWS rejects the combination)
- Tighter `accessTokenValidity` / `idTokenValidity` for sensitive applications
