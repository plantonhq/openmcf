# Okta Federation

This preset creates the trust anchor for the most common IdP pairing:
Okta's AWS Account Federation app signing users into IAM roles.

## When to Use

- Okta as the corporate identity provider with the "AWS Account
  Federation" app
- Per-group role mapping managed on the Okta side

## What You Get

- The IAM SAML provider created from the Okta app's published
  metadata
- The provider ARN output that Okta's app settings and every role
  trust policy reference

## Customize

- In Okta, set the app's Identity Provider ARN to this provider's
  ARN output and map groups to `role_arn,provider_arn` pairs
- Okta rotates its signing certificate on a schedule — re-paste the
  metadata before `valid_until` (an in-place update)
- Session length is the ROLE's max_session_duration; align Okta's
  session settings with it
