# AWS Cognito User Pool -- Architecture and Design

## Overview

Amazon Cognito user pools are a fully managed user directory that provides sign-up, sign-in, and token-based authentication for web and mobile applications. A pool implements OAuth 2.0 and OpenID Connect: applications authenticate users against it and receive JWTs without running any identity infrastructure.

The pool is the ROOT of a family of resources. Everything that is pool-scoped configuration lives on this kind; everything with its own AWS lifecycle is its own kind that references the pool:

- **App clients** (`AwsCognitoUserPoolClient`) -- many per pool. Each is one application's OAuth contract, and its client ID is what JWT authorizers validate as the token audience.
- **Identity providers** (`AwsCognitoIdentityProvider`) -- many per pool, federating Google/Facebook/Amazon/Apple/OIDC/SAML sign-in.
- **Resource servers** (`AwsCognitoResourceServer`) -- many per pool, minting the custom OAuth scopes machine-to-machine clients request.

The hosted-UI **domain** is the deliberate exception: AWS allows one domain per pool and the domain string is its identity, so it folds into the pool's spec and lifecycle.

## Identity Model

Cognito offers two mutually exclusive identity models, and the choice is **permanent** -- changing it replaces the pool and every user in it:

- **Username attributes** (`username_attributes`): the chosen attribute IS the username. `["email"]` is the common consumer-app model; Cognito generates the immutable `sub` claim internally.
- **Alias attributes** (`alias_attributes`): users have a real username, and the aliases are alternative sign-in identifiers -- common for gaming and enterprise apps where usernames matter.
- **Neither**: users pick a unique username at sign-up.

## Sign-In Policy: the Passwordless Dial

`allowed_first_auth_factors` selects what users may present as their FIRST factor: `PASSWORD`, `EMAIL_OTP`, `SMS_OTP`, `WEB_AUTHN` (passkeys). Leaving it empty keeps classic password-first behavior; setting it enables the choice-based sign-in flow (clients opt in with the `ALLOW_USER_AUTH` auth flow). Passkeys bind to the WebAuthn relying-party ID -- pin `web_authn.relying_party_id` to the domain your applications serve BEFORE users register passkeys, because changing it invalidates existing registrations.

## Feature Tiers

`user_pool_tier` selects the AWS feature plan and billing: `LITE` (core sign-in), `ESSENTIALS` (the AWS default -- passwordless factors, managed login, password history), and `PLUS` (threat protection). Downgrades are allowed but AWS rejects them while a higher-tier feature is still configured.

## MFA Architecture

`mfa_configuration` sets enforcement (`OFF`/`OPTIONAL`/`ON`); the factor set is configured independently:

- **TOTP** (`software_token_mfa_enabled`) -- authenticator apps; no delivery dependency.
- **Email OTP** (`email_mfa`) -- requires SES `DEVELOPER` email (codes exceed the built-in sender's guarantees).
- **SMS OTP** -- requires `sms_configuration` (below).
- **Passkeys** participate as a first factor rather than a second.

## Delivery Channels

**Email** has two modes: `COGNITO_DEFAULT` (built-in, ~50 emails/day -- development only) and `DEVELOPER` (your verified SES identity -- production). **SMS** always routes through SNS by role assumption: `sms_configuration.sns_caller_arn` references the IAM role, `external_id` is the confused-deputy guard the role's trust policy must match, and every SMS-dependent feature (phone auto-verification, SMS OTP, SMS MFA) needs it.

For full control, `custom_email_sender` / `custom_sms_sender` hand delivery to your own Lambda; Cognito encrypts the code payload with `lambda_config.kms_key_id`, which is why the key is required with either sender.

## Custom Attributes are Append-Only

The pool schema accepts up to 50 custom attributes (auto-prefixed `custom:`). AWS has no API to modify or remove an attribute once added -- removing one from the spec errors rather than recreating the pool, and the `mutable`/`required`/`developer_only_attribute` flags are fixed at the moment of addition. Model attribute changes as NEW attributes.

## Lambda Triggers

Cognito invokes Lambda at every lifecycle point: sign-up validation, pre/post authentication, post-confirmation provisioning, message customization, on-the-fly user migration, and the three-step custom-challenge flow. Claim customization has two spellings: the legacy `pre_token_generation` (V1_0 event, identity token only) and `pre_token_generation_config` with an explicit `lambda_version` -- `V2_0`/`V3_0` deliver the richer event that also customizes ACCESS tokens. Exactly one may be set. Every trigger function must grant `cognito-idp.amazonaws.com` invoke permission.

## Threat Protection and Log Delivery

`user_pool_add_ons.advanced_security_mode` turns on risk-based protection (`AUDIT` gathers signals; `ENFORCED` acts on them) -- both require the PLUS tier. `custom_auth_mode` extends coverage to Lambda-driven custom auth flows, which standard mode does not inspect. `log_configurations` routes two event sources to destinations you own: `userNotification` (message-delivery errors, ERROR level) and `userAuthEvents` (detailed auth telemetry from threat protection, INFO level) -- each entry targets exactly one CloudWatch log group, Firehose stream, or S3 bucket by reference.

## Domain and Hosted UI

The folded domain serves the hosted sign-in pages and the standard OAuth2 endpoints (`/oauth2/authorize`, `/oauth2/token`, `/oauth2/userInfo`):

- **Prefix domain**: `{prefix}.auth.{region}.amazoncognito.com` -- free, resolves in about a minute.
- **Custom domain**: your FQDN, fronted by CloudFront, requiring an ACM certificate in **us-east-1** regardless of the pool's region. The pool exports the CloudFront distribution domain and its fixed alias zone ID so a Route53 alias record composes directly from outputs. Adding or removing the certificate switches the domain type and replaces the domain; rotating one certificate ARN to another updates in place.
- `managed_login_version` selects the page generation: 1 is the classic hosted UI; 2 is managed login (the AWS default for new domains), which requires a branding style to be assigned in the console before pages render.

## The Three Join Keys

Downstream composition hangs off three outputs with distinct shapes -- they are deliberately separate:

- **`issuer`** -- the full `https://` OIDC issuer URL. JWT authorizers (API Gateway, ALB) validate the token's `iss` claim against exactly this string.
- **`user_pool_domain`** -- the RAW domain (prefix or FQDN, no scheme). ALB `authenticate_cognito` actions take exactly this shape.
- **`hosted_ui_url`** -- the full `https://` URL applications link users to.

(`user_pool_endpoint` remains the scheme-less host/path AWS reports, for SDK configuration.)

## Cost Model

Cognito bills on Monthly Active Users, by tier (LITE/ESSENTIALS/PLUS have different free allowances and per-MAU rates). Threat protection is part of the PLUS tier. A pool with no users costs nothing meaningful; custom domains add no Cognito cost (the ACM certificate is free).

## Security Defaults Worth Setting

1. `ALLOW_USER_SRP_AUTH` on clients -- the password never leaves the device.
2. `prevent_user_existence_errors: ENABLED` on clients -- blocks user enumeration.
3. `deletion_protection: true` for production -- a deleted pool takes every user with it, unrecoverably.
4. `DEVELOPER` email mode for anything beyond development.
5. `attributes_require_verification_before_update: ["email"]` -- an unverified typo in an email update cannot lock the user out of recovery.
6. MFA `OPTIONAL` at minimum; `ON` for sensitive applications.
