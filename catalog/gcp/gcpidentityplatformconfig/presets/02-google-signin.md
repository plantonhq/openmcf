# Google Sign-in

Email/password plus "Sign in with Google" — the standard consumer-app
pair. Google handles the second flow's credentials; your app never sees
a password for those users.

## What it configures

- `signIn.email` enabled with passwords required.
- One `defaultSupportedIdps` entry with `idpId: google.com` and the
  OAuth client obtained from the Google Cloud console.

## Adjust before deploying

- **clientId / clientSecret** — replace the placeholders with the real
  OAuth client from the Google Cloud console's OAuth consent screen
  (APIs & Services → Credentials). There is no API that creates
  consent-screen clients — this is a one-time console step. The secret
  is handled as a managed secret end to end.
- **authorizedDomains** — add your app's domain so the OAuth redirect
  completes from it.

## When to choose something else

For other well-known providers, change `idpId` (facebook.com,
apple.com, github.com, microsoft.com, ...) and supply that provider's
console-issued client. For a corporate IdP, use `oauthIdpConfigs`
(custom OIDC) or the **Enterprise SAML** preset instead.
