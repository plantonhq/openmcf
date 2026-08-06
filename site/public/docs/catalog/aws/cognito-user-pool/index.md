---
title: "Cognito User Pool"
description: "Cognito User Pool deployment documentation"
icon: "package"
order: 100
componentName: "awscognitouserpool"
---

# AWS Cognito User Pool

Deploys an Amazon Cognito user pool -- the managed user directory and OIDC token issuer for web and mobile applications. The pool owns everything pool-scoped: the identity model, password and passwordless sign-in policies, MFA (TOTP, email, SMS, passkeys), email/SMS delivery, verification and invitation messaging, custom schema attributes, Lambda triggers, threat protection, log delivery, and the hosted-UI domain. App clients, identity providers, and resource servers compose onto the pool as their own resources.

## What Gets Created

When you deploy an AwsCognitoUserPool resource, Planton provisions:

- **Cognito User Pool** -- an `aws_cognito_user_pool` with the configured identity model, policies, MFA, delivery, triggers, and threat protection
- **User Pool Domain** -- created only when `spec.domain` is set: an `aws_cognito_user_pool_domain` serving the hosted sign-in UI and OAuth2 endpoints (Authorization, Token, UserInfo)
- **Log delivery configuration** -- created only when `spec.logConfigurations` is set, routing Cognito event logs to CloudWatch, Firehose, or S3

App clients (`AwsCognitoUserPoolClient`), federated identity providers (`AwsCognitoIdentityProvider`), and resource servers (`AwsCognitoResourceServer`) are separate resources that reference this pool.

## Prerequisites

- **AWS credentials** configured via environment variables or Planton provider config
- **An ACM certificate in us-east-1** if configuring a custom domain (Cognito fronts custom domains with CloudFront)
- **A verified SES identity** if using `emailConfiguration.emailSendingAccount: DEVELOPER` (production email volumes; required for email MFA)
- **An IAM role assumable by cognito-idp.amazonaws.com** (with an `sts:ExternalId` condition) for any SMS feature
- **Lambda function(s)** with `cognito-idp.amazonaws.com` invoke permission if configuring triggers

## Quick Start

Create a file `cognito.yaml`:

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsCognitoUserPool
metadata:
  name: my-auth
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AwsCognitoUserPool.my-auth
spec:
  region: us-east-1
  usernameAttributes:
    - email
  autoVerifiedAttributes:
    - email
  accountRecoveryMechanisms:
    - name: verified_email
      priority: 1
```

Deploy:

```shell
planton apply -f cognito.yaml
```

This creates a user pool where users sign in with their email address and email is auto-verified on sign-up. Add an `AwsCognitoUserPoolClient` so an application can authenticate against it.

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `region` | `string` | The AWS region where the user pool is created (e.g., `us-east-1`). | Required |

### Identity Model (ForceNew)

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `usernameAttributes` | `string[]` | `[]` | Sign-in identifiers: `"email"`, `"phone_number"`. Mutually exclusive with `aliasAttributes`. Changing REPLACES the pool and its users. |
| `aliasAttributes` | `string[]` | `[]` | Alias identifiers: `"email"`, `"phone_number"`, `"preferred_username"`. Mutually exclusive with `usernameAttributes`. |
| `usernameCaseSensitive` | `bool` | `false` | Case-sensitive usernames. |

### Pool Posture

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `deletionProtection` | `bool` | `false` | Pool cannot be deleted while true. |
| `userPoolTier` | `string` | `"ESSENTIALS"` (AWS) | Feature tier: `"LITE"`, `"ESSENTIALS"`, `"PLUS"` (threat protection). |

### Password and Sign-In Policy

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `passwordPolicy.minimumLength` | `int` | 8 (AWS) | Minimum password length. Range: 6-99. |
| `passwordPolicy.requireLowercase` / `requireUppercase` / `requireNumbers` / `requireSymbols` | `bool` | `false` | Character-class requirements. |
| `passwordPolicy.passwordHistorySize` | `int` | 0 | Previous passwords a new one must differ from (0-24; ESSENTIALS+). |
| `passwordPolicy.temporaryPasswordValidityDays` | `int` | 7 (AWS) | Days until admin-created temporary passwords expire (0-365). |
| `allowedFirstAuthFactors` | `string[]` | `[]` | Passwordless first factors: `"PASSWORD"`, `"EMAIL_OTP"`, `"SMS_OTP"`, `"WEB_AUTHN"`. Empty keeps classic password-first sign-in. |

### MFA

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `mfaConfiguration` | `string` | `"OFF"` | MFA enforcement: `"OFF"`, `"OPTIONAL"`, `"ON"`. |
| `softwareTokenMfaEnabled` | `bool` | `false` | TOTP authenticator apps. Requires MFA not OFF. |
| `emailMfa.message` / `emailMfa.subject` | `string` | AWS defaults | Email-OTP second factor templates (`{####}` placeholder). Requires SES DEVELOPER email. |
| `webAuthn.relyingPartyId` | `string` | pool domain | The domain passkeys are registered against. |
| `webAuthn.userVerification` | `string` | `"preferred"` | `"required"` or `"preferred"`. |

### SMS Delivery

| Field | Type | Description |
|-------|------|-------------|
| `smsConfiguration.snsCallerArn` | `StringValueOrRef` | IAM role Cognito assumes to publish SMS via SNS. Reference an AwsIamRole via `valueFrom`. |
| `smsConfiguration.externalId` | `string` | The `sts:ExternalId` in the role's trust policy. |
| `smsConfiguration.snsRegion` | `string` | Region SMS originates from when the pool's region cannot send. |
| `smsAuthenticationMessage` | `string` | Sign-in code SMS body; must contain `{####}`. |

### Verification and Recovery

| Field | Type | Description |
|-------|------|-------------|
| `autoVerifiedAttributes` | `string[]` | Auto-verify on sign-up: `"email"`, `"phone_number"`. |
| `attributesRequireVerificationBeforeUpdate` | `string[]` | Keep the previous value active until the updated one is verified. |
| `accountRecoveryMechanisms` | `object[]` | Recovery methods with `.name` (`"verified_email"`, `"verified_phone_number"`, `"admin_only"`) and `.priority` (1-2). |
| `verificationMessageTemplate.defaultEmailOption` | `string` | `"CONFIRM_WITH_CODE"` (AWS default) or `"CONFIRM_WITH_LINK"`. |
| `verificationMessageTemplate.emailMessage` / `emailSubject` | `string` | Code-based templates (`{####}`). |
| `verificationMessageTemplate.emailMessageByLink` / `emailSubjectByLink` | `string` | Link-based templates (`{##link text##}`). |
| `verificationMessageTemplate.smsMessage` | `string` | Phone verification SMS (`{####}`). |

### Email Configuration

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `emailConfiguration.emailSendingAccount` | `string` | `"COGNITO_DEFAULT"` | `"COGNITO_DEFAULT"` (50/day sandbox) or `"DEVELOPER"` (SES). |
| `emailConfiguration.sourceArn` | `StringValueOrRef` | — | SES identity ARN. Required for `"DEVELOPER"`. |
| `emailConfiguration.fromEmailAddress` | `string` | — | "From" address for DEVELOPER mode. |
| `emailConfiguration.replyToEmailAddress` | `string` | — | Reply-to address. |
| `emailConfiguration.configurationSet` | `string` | — | SES configuration set for delivery metrics. |

### Admin and User Management

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `allowAdminCreateUserOnly` | `bool` | `false` | Disable self-registration. |
| `inviteMessageTemplate.emailMessage` / `emailSubject` / `smsMessage` | `string` | AWS defaults | Invitation templates; must contain `{username}` and `{####}`. |
| `deviceConfiguration.challengeRequiredOnNewDevice` | `bool` | `false` | New devices still require a challenge the first time. |
| `deviceConfiguration.deviceOnlyRememberedOnUserPrompt` | `bool` | `false` | Remember devices only when the user opts in. |

### Custom Attributes (Append-Only)

| Field | Type | Description |
|-------|------|-------------|
| `customAttributes[].name` | `string` | 1-20 chars, auto-prefixed `custom:`. |
| `customAttributes[].attributeDataType` | `string` | `"String"`, `"Number"`, `"DateTime"`, `"Boolean"`. |
| `customAttributes[].mutable` / `required` / `developerOnlyAttribute` | `bool` | Fixed at the moment the attribute is added. |
| `customAttributes[].stringMinLength` / `stringMaxLength` / `numberMinValue` / `numberMaxValue` | `string` | Type constraints. |

### Lambda Triggers

All trigger fields accept a Lambda ARN or a `valueFrom` reference to an AwsLambda resource.

| Field | Description |
|-------|-------------|
| `lambdaConfig.preSignUp` / `preAuthentication` / `postAuthentication` / `postConfirmation` | Lifecycle hooks. |
| `lambdaConfig.preTokenGeneration` | Claim customization (V1_0 event). |
| `lambdaConfig.preTokenGenerationConfig.lambdaArn` + `.lambdaVersion` | Versioned claim customization -- `"V2_0"`/`"V3_0"` also customize ACCESS tokens. Set this or `preTokenGeneration`, not both. |
| `lambdaConfig.customMessage` / `userMigration` | Message customization; on-the-fly user import. |
| `lambdaConfig.defineAuthChallenge` / `createAuthChallenge` / `verifyAuthChallengeResponse` | Custom auth challenge flow. |
| `lambdaConfig.customEmailSender` / `customSmsSender` | Self-managed delivery (`lambdaArn` + `lambdaVersion: "V1_0"`); requires `lambdaConfig.kmsKeyId`. |
| `lambdaConfig.kmsKeyId` | KMS key encrypting custom-sender payloads. Reference an AwsKmsKey via `valueFrom`. |

### Threat Protection and Log Delivery

| Field | Type | Description |
|-------|------|-------------|
| `userPoolAddOns.advancedSecurityMode` | `string` | `"OFF"`, `"AUDIT"`, `"ENFORCED"`. AUDIT/ENFORCED require the PLUS tier. |
| `userPoolAddOns.customAuthMode` | `string` | Extends threat protection to custom auth flows: `"AUDIT"`, `"ENFORCED"`. |
| `logConfigurations[].eventSource` | `string` | `"userNotification"` (ERROR) or `"userAuthEvents"` (INFO; PLUS tier). |
| `logConfigurations[].logLevel` | `string` | `"ERROR"` or `"INFO"`. |
| `logConfigurations[].cloudwatchLogGroupArn` / `firehoseStreamArn` / `s3BucketArn` | `StringValueOrRef` | Exactly one destination per entry; each accepts a `valueFrom` reference. |

### Hosted-UI Domain

| Field | Type | Description |
|-------|------|-------------|
| `domain.domain` | `string` | Cognito prefix (e.g., `"myapp-auth"`) or custom FQDN (e.g., `"auth.example.com"`). ForceNew. |
| `domain.certificateArn` | `StringValueOrRef` | ACM cert for custom domains (us-east-1). Reference an AwsCertManagerCert via `valueFrom`. |
| `domain.managedLoginVersion` | `int` | 1 = classic hosted UI, 2 = managed login (AWS default for new domains). |

## Examples

### Passwordless Sign-In (Email OTP + Passkeys)

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsCognitoUserPool
metadata:
  name: passwordless-auth
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: acme
    pulumi.planton.dev/project: platform
    pulumi.planton.dev/stack.name: prod.AwsCognitoUserPool.passwordless-auth
spec:
  region: us-east-1
  usernameAttributes:
    - email
  autoVerifiedAttributes:
    - email
  allowedFirstAuthFactors:
    - PASSWORD
    - EMAIL_OTP
    - WEB_AUTHN
  webAuthn:
    relyingPartyId: auth.example.com
    userVerification: required
  emailConfiguration:
    emailSendingAccount: DEVELOPER
    sourceArn:
      value: "arn:aws:ses:us-east-1:123456789012:identity/noreply@example.com"
    fromEmailAddress: "Acme <noreply@example.com>"
  accountRecoveryMechanisms:
    - name: verified_email
      priority: 1
```

### Production Pool with MFA, Threat Protection, and a Hosted UI

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsCognitoUserPool
metadata:
  name: prod-auth
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: acme
    pulumi.planton.dev/project: platform
    pulumi.planton.dev/stack.name: prod.AwsCognitoUserPool.prod-auth
spec:
  region: us-east-1
  usernameAttributes:
    - email
  userPoolTier: PLUS
  passwordPolicy:
    minimumLength: 12
    requireLowercase: true
    requireUppercase: true
    requireNumbers: true
    requireSymbols: true
    passwordHistorySize: 5
    temporaryPasswordValidityDays: 3
  mfaConfiguration: OPTIONAL
  softwareTokenMfaEnabled: true
  autoVerifiedAttributes:
    - email
  attributesRequireVerificationBeforeUpdate:
    - email
  accountRecoveryMechanisms:
    - name: verified_email
      priority: 1
  emailConfiguration:
    emailSendingAccount: DEVELOPER
    sourceArn:
      value: "arn:aws:ses:us-east-1:123456789012:identity/noreply@example.com"
    fromEmailAddress: "Acme <noreply@example.com>"
  userPoolAddOns:
    advancedSecurityMode: ENFORCED
  deletionProtection: true
  domain:
    domain: acme-prod-auth
```

### API Gateway JWT Integration

The pool's `issuer` output and an app client's `client_id` output wire directly into an AwsHttpApiGateway JWT authorizer:

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsHttpApiGateway
metadata:
  name: my-api
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: acme
    pulumi.planton.dev/project: platform
    pulumi.planton.dev/stack.name: prod.AwsHttpApiGateway.my-api
spec:
  region: us-east-1
  routes:
    - routeKey: "GET /users"
      integration:
        integrationType: AWS_PROXY
        integrationUri:
          valueFrom:
            kind: AwsLambda
            name: get-users
            fieldPath: status.outputs.function_arn
      authorizationType: JWT
      authorizerName: cognito
  authorizers:
    - name: cognito
      authorizerType: JWT
      jwtConfiguration:
        issuer:
          valueFrom:
            kind: AwsCognitoUserPool
            name: prod-auth
            fieldPath: status.outputs.issuer
        audiences:
          - valueFrom:
              kind: AwsCognitoUserPoolClient
              name: web-spa
              fieldPath: status.outputs.client_id
```

## Stack Outputs

| Output | Type | Description |
|--------|------|-------------|
| `user_pool_id` | `string` | Pool identifier (`{region}_{poolId}`) -- the join key app clients, identity providers, and resource servers reference. |
| `user_pool_arn` | `string` | Pool ARN -- IAM policies and ALB authenticate-cognito actions. |
| `user_pool_endpoint` | `string` | The endpoint as AWS reports it, WITHOUT a scheme (`cognito-idp.{region}.amazonaws.com/{id}`). |
| `issuer` | `string` | The full OIDC issuer URL -- what JWT authorizers validate the token's `iss` claim against. |
| `user_pool_domain` | `string` | The hosted-UI domain exactly as configured (prefix or FQDN, no scheme) -- what ALB actions take as `user_pool_domain`. Empty when no domain. |
| `hosted_ui_url` | `string` | The full `https://` hosted sign-in URL. Empty when no domain. |
| `cloudfront_distribution` | `string` | CloudFront domain name fronting a custom domain -- the DNS alias target. Empty for prefix domains. |
| `cloudfront_distribution_arn` | `string` | ARN of that CloudFront distribution. Empty for prefix domains. |
| `cloudfront_hosted_zone_id` | `string` | The CloudFront alias-target zone ID for Route53 alias records. Empty for prefix domains. |

## Related Components

- [AWS Cognito User Pool Client](/docs/catalog/aws/cognito-user-pool-client) -- app clients that authenticate against this pool; their `client_id` is the JWT audience
- [AWS Cognito Identity Provider](/docs/catalog/aws/cognito-identity-provider) -- federated Google/Facebook/Amazon/Apple/OIDC/SAML sign-in for this pool
- [AWS Cognito Resource Server](/docs/catalog/aws/cognito-resource-server) -- custom OAuth scopes for machine-to-machine clients
- [AWS HTTP API Gateway](/docs/catalog/aws/http-api-gateway) -- JWT authorizers built from this pool's `issuer` output
- [AWS Lambda](/docs/catalog/aws/lambda) -- trigger functions for the pool's lifecycle hooks
- [AWS ACM Certificate](/docs/catalog/aws/cert-manager-cert) -- the us-east-1 certificate for custom hosted-UI domains
- [AWS Route53 DNS Record](/docs/catalog/aws/route53-dns-record) -- the alias record pointing a custom domain at the pool's CloudFront distribution
