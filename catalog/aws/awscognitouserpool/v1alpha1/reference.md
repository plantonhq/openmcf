# AwsCognitoUserPool

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `aws.planton.dev/v1alpha1`

AwsCognitoUserPoolSpec defines the desired configuration for an AWS Cognito
User Pool -- the user directory and authentication service for web and
mobile applications.

The user pool owns everything that is pool-scoped configuration: the
identity model, password and sign-in policies, MFA, email/SMS delivery,
verification and invitation messaging, custom schema attributes, Lambda
triggers, threat protection, log delivery, and the hosted-UI domain.

Resources with their own AWS lifecycle compose onto the pool as separate
kinds rather than being buried inside it:
- **App clients** (`AwsCognitoUserPoolClient`): many per pool, each with its
  own OAuth contract and its own client ID that downstream systems (JWT
  authorizers, ALB authentication actions) reference.
- **Identity providers** (`AwsCognitoIdentityProvider`): many per pool,
  federating social/OIDC/SAML sign-in.
- **Resource servers** (`AwsCognitoResourceServer`): many per pool, defining
  custom OAuth scopes for machine-to-machine clients.

The hosted-UI **domain** stays folded here: AWS allows one domain per pool
and the domain string is its identity, so it shares the pool's lifecycle.

Key design notes:
- Identity model fields (`username_attributes`, `alias_attributes`,
  `username_case_sensitive`) are **ForceNew** in AWS -- changing them
  destroys and recreates the pool along with every user in it.
- Custom schema attributes are **append-only**: new attributes can be added
  in place, but an existing attribute can never be modified or removed
  (AWS has no API for it).
- Threat protection (`user_pool_add_ons`) beyond audit-only requires the
  PLUS feature tier (`user_pool_tier`).

Credentials, region, and deployment workflow live outside this spec in
stack inputs.

## Example

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsCognitoUserPool
metadata:
  name: test-auth-pool
  org: test-org
  env: dev
  id: awscog-test-001
spec:
  region: us-west-2
  usernameAttributes:
    - email
  passwordPolicy:
    minimumLength: 8
    requireLowercase: true
    requireUppercase: true
    requireNumbers: true
    requireSymbols: false
    passwordHistorySize: 3
  mfaConfiguration: OPTIONAL
  softwareTokenMfaEnabled: true
  autoVerifiedAttributes:
    - email
  attributesRequireVerificationBeforeUpdate:
    - email
  accountRecoveryMechanisms:
    - name: verified_email
      priority: 1
  deletionProtection: false
  domain:
    domain: test-auth-pool
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.usernameAttributes` | `[]string` |  |  |  |
| `spec.aliasAttributes` | `[]string` |  |  |  |
| `spec.usernameCaseSensitive` | `bool` |  |  |  |
| `spec.deletionProtection` | `bool` |  |  |  |
| `spec.userPoolTier` | `string` |  |  |  |
| `spec.passwordPolicy` | `AwsCognitoUserPoolPasswordPolicy` |  |  |  |
| `spec.passwordPolicy.minimumLength` | `int32` |  |  |  |
| `spec.passwordPolicy.requireLowercase` | `bool` |  |  |  |
| `spec.passwordPolicy.requireUppercase` | `bool` |  |  |  |
| `spec.passwordPolicy.requireNumbers` | `bool` |  |  |  |
| `spec.passwordPolicy.requireSymbols` | `bool` |  |  |  |
| `spec.passwordPolicy.passwordHistorySize` | `int32` |  |  |  |
| `spec.passwordPolicy.temporaryPasswordValidityDays` | `int32` |  |  |  |
| `spec.allowedFirstAuthFactors` | `[]string` |  |  |  |
| `spec.mfaConfiguration` | `string` |  |  |  |
| `spec.softwareTokenMfaEnabled` | `bool` |  |  |  |
| `spec.emailMfa` | `AwsCognitoUserPoolEmailMfaConfig` |  |  |  |
| `spec.emailMfa.message` | `string` |  |  |  |
| `spec.emailMfa.subject` | `string` |  |  |  |
| `spec.webAuthn` | `AwsCognitoUserPoolWebAuthnConfig` |  |  |  |
| `spec.webAuthn.relyingPartyId` | `string` |  |  |  |
| `spec.webAuthn.userVerification` | `string` |  |  |  |
| `spec.smsConfiguration` | `AwsCognitoUserPoolSmsConfig` |  |  |  |
| `spec.smsConfiguration.snsCallerArn` | `string \| valueFrom` | yes |  | AwsIamRole (`status.outputs.role_arn`) |
| `spec.smsConfiguration.externalId` | `string` | yes |  |  |
| `spec.smsConfiguration.snsRegion` | `string` |  |  |  |
| `spec.smsAuthenticationMessage` | `string` |  |  |  |
| `spec.autoVerifiedAttributes` | `[]string` |  |  |  |
| `spec.attributesRequireVerificationBeforeUpdate` | `[]string` |  |  |  |
| `spec.accountRecoveryMechanisms` | `[]AwsCognitoUserPoolRecoveryMechanism` |  |  |  |
| `spec.accountRecoveryMechanisms[].name` | `string` | yes |  |  |
| `spec.accountRecoveryMechanisms[].priority` | `int32` | yes |  |  |
| `spec.emailConfiguration` | `AwsCognitoUserPoolEmailConfig` |  |  |  |
| `spec.emailConfiguration.emailSendingAccount` | `string` |  |  |  |
| `spec.emailConfiguration.sourceArn` | `string \| valueFrom` |  |  |  |
| `spec.emailConfiguration.fromEmailAddress` | `string` |  |  |  |
| `spec.emailConfiguration.replyToEmailAddress` | `string` |  |  |  |
| `spec.emailConfiguration.configurationSet` | `string` |  |  |  |
| `spec.verificationMessageTemplate` | `AwsCognitoUserPoolVerificationMessageTemplate` |  |  |  |
| `spec.verificationMessageTemplate.defaultEmailOption` | `string` |  |  |  |
| `spec.verificationMessageTemplate.emailMessage` | `string` |  |  |  |
| `spec.verificationMessageTemplate.emailSubject` | `string` |  |  |  |
| `spec.verificationMessageTemplate.emailMessageByLink` | `string` |  |  |  |
| `spec.verificationMessageTemplate.emailSubjectByLink` | `string` |  |  |  |
| `spec.verificationMessageTemplate.smsMessage` | `string` |  |  |  |
| `spec.allowAdminCreateUserOnly` | `bool` |  |  |  |
| `spec.inviteMessageTemplate` | `AwsCognitoUserPoolInviteMessageTemplate` |  |  |  |
| `spec.inviteMessageTemplate.emailMessage` | `string` |  |  |  |
| `spec.inviteMessageTemplate.emailSubject` | `string` |  |  |  |
| `spec.inviteMessageTemplate.smsMessage` | `string` |  |  |  |
| `spec.deviceConfiguration` | `AwsCognitoUserPoolDeviceConfig` |  |  |  |
| `spec.deviceConfiguration.challengeRequiredOnNewDevice` | `bool` |  |  |  |
| `spec.deviceConfiguration.deviceOnlyRememberedOnUserPrompt` | `bool` |  |  |  |
| `spec.customAttributes` | `[]AwsCognitoUserPoolSchemaAttribute` |  |  |  |
| `spec.customAttributes[].name` | `string` | yes |  |  |
| `spec.customAttributes[].attributeDataType` | `string` | yes |  |  |
| `spec.customAttributes[].mutable` | `bool` |  |  |  |
| `spec.customAttributes[].required` | `bool` |  |  |  |
| `spec.customAttributes[].developerOnlyAttribute` | `bool` |  |  |  |
| `spec.customAttributes[].stringMinLength` | `string` |  |  |  |
| `spec.customAttributes[].stringMaxLength` | `string` |  |  |  |
| `spec.customAttributes[].numberMinValue` | `string` |  |  |  |
| `spec.customAttributes[].numberMaxValue` | `string` |  |  |  |
| `spec.lambdaConfig` | `AwsCognitoUserPoolLambdaConfig` |  |  |  |
| `spec.lambdaConfig.preSignUp` | `string \| valueFrom` |  |  | AwsLambda (`status.outputs.function_arn`) |
| `spec.lambdaConfig.preAuthentication` | `string \| valueFrom` |  |  | AwsLambda (`status.outputs.function_arn`) |
| `spec.lambdaConfig.postAuthentication` | `string \| valueFrom` |  |  | AwsLambda (`status.outputs.function_arn`) |
| `spec.lambdaConfig.postConfirmation` | `string \| valueFrom` |  |  | AwsLambda (`status.outputs.function_arn`) |
| `spec.lambdaConfig.preTokenGeneration` | `string \| valueFrom` |  |  | AwsLambda (`status.outputs.function_arn`) |
| `spec.lambdaConfig.preTokenGenerationConfig` | `AwsCognitoUserPoolPreTokenGenerationConfig` |  |  |  |
| `spec.lambdaConfig.preTokenGenerationConfig.lambdaArn` | `string \| valueFrom` | yes |  | AwsLambda (`status.outputs.function_arn`) |
| `spec.lambdaConfig.preTokenGenerationConfig.lambdaVersion` | `string` | yes |  |  |
| `spec.lambdaConfig.customMessage` | `string \| valueFrom` |  |  | AwsLambda (`status.outputs.function_arn`) |
| `spec.lambdaConfig.userMigration` | `string \| valueFrom` |  |  | AwsLambda (`status.outputs.function_arn`) |
| `spec.lambdaConfig.defineAuthChallenge` | `string \| valueFrom` |  |  | AwsLambda (`status.outputs.function_arn`) |
| `spec.lambdaConfig.createAuthChallenge` | `string \| valueFrom` |  |  | AwsLambda (`status.outputs.function_arn`) |
| `spec.lambdaConfig.verifyAuthChallengeResponse` | `string \| valueFrom` |  |  | AwsLambda (`status.outputs.function_arn`) |
| `spec.lambdaConfig.customEmailSender` | `AwsCognitoUserPoolCustomSenderConfig` |  |  |  |
| `spec.lambdaConfig.customEmailSender.lambdaArn` | `string \| valueFrom` | yes |  | AwsLambda (`status.outputs.function_arn`) |
| `spec.lambdaConfig.customEmailSender.lambdaVersion` | `string` | yes |  |  |
| `spec.lambdaConfig.customSmsSender` | `AwsCognitoUserPoolCustomSenderConfig` |  |  |  |
| `spec.lambdaConfig.customSmsSender.lambdaArn` | `string \| valueFrom` | yes |  | AwsLambda (`status.outputs.function_arn`) |
| `spec.lambdaConfig.customSmsSender.lambdaVersion` | `string` | yes |  |  |
| `spec.lambdaConfig.kmsKeyId` | `string \| valueFrom` |  |  | AwsKmsKey (`status.outputs.key_arn`) |
| `spec.userPoolAddOns` | `AwsCognitoUserPoolAddOns` |  |  |  |
| `spec.userPoolAddOns.advancedSecurityMode` | `string` | yes |  |  |
| `spec.userPoolAddOns.customAuthMode` | `string` |  |  |  |
| `spec.logConfigurations` | `[]AwsCognitoUserPoolLogConfiguration` |  |  |  |
| `spec.logConfigurations[].eventSource` | `string` | yes |  |  |
| `spec.logConfigurations[].logLevel` | `string` | yes |  |  |
| `spec.logConfigurations[].cloudwatchLogGroupArn` | `string \| valueFrom` |  |  | AwsCloudwatchLogGroup (`status.outputs.log_group_arn`) |
| `spec.logConfigurations[].firehoseStreamArn` | `string \| valueFrom` |  |  | AwsKinesisFirehose (`status.outputs.delivery_stream_arn`) |
| `spec.logConfigurations[].s3BucketArn` | `string \| valueFrom` |  |  | AwsS3Bucket (`status.outputs.bucket_arn`) |
| `spec.domain` | `AwsCognitoUserPoolDomainConfig` |  |  |  |
| `spec.domain.domain` | `string` | yes |  |  |
| `spec.domain.certificateArn` | `string \| valueFrom` |  |  | AwsCertManagerCert (`status.outputs.cert_arn`) |
| `spec.domain.managedLoginVersion` | `int32` |  |  |  |

## Field Details

### spec.region

`string` · required

The AWS region where the resource will be created.
Example: "us-west-2", "eu-west-1"

- rule: {"string":{"minLen":"1"}}

### spec.usernameAttributes

`[]string`

Attributes that users can use as their username when signing in. When set,
the chosen attribute(s) become the sign-in identifier. Common choice: ["email"]
to let users sign in with their email address.

Valid values: "email", "phone_number".
Mutually exclusive with `alias_attributes`. ForceNew.

### spec.aliasAttributes

`[]string`

Attributes that can be used as aliases for the username. Unlike `username_attributes`,
users always have a separate username and the alias attributes are alternative sign-in
identifiers. Common choice: ["email", "preferred_username"].

Valid values: "email", "phone_number", "preferred_username".
Mutually exclusive with `username_attributes`. ForceNew.

### spec.usernameCaseSensitive

`bool`

Whether usernames are case-sensitive. When false (default in AWS), "User" and "user"
are treated as the same username. ForceNew -- cannot be changed after pool creation.

### spec.deletionProtection

`bool`

Enable deletion protection. When true, the user pool cannot be deleted
without first disabling this setting. Recommended for production pools --
a deleted pool takes every user account with it, unrecoverably.

### spec.userPoolTier

`string`

The feature tier of the user pool. Governs which capabilities AWS makes
available (and how the pool is billed):
- "LITE": lowest-cost tier; core sign-in only.
- "ESSENTIALS": the AWS default for new pools; adds passwordless
  sign-in (email/SMS OTP, passkeys) and managed login branding.
- "PLUS": adds threat protection (advanced security features).
When omitted, AWS creates the pool on the ESSENTIALS tier. Downgrading a
tier is allowed but AWS rejects it while a feature exclusive to the
higher tier is still configured (e.g. PLUS -> ESSENTIALS with threat
protection enforced).

### spec.passwordPolicy

`AwsCognitoUserPoolPasswordPolicy`

Password policy configuration. Controls password strength requirements for
user self-registration and admin-created passwords.

### spec.passwordPolicy.minimumLength

`int32`

Minimum password length. Range: 6-99. AWS default: 8.

- rule: {"int32":{"lte":99,"gte":6}}

### spec.passwordPolicy.requireLowercase

`bool`

Require at least one lowercase letter.

### spec.passwordPolicy.requireUppercase

`bool`

Require at least one uppercase letter.

### spec.passwordPolicy.requireNumbers

`bool`

Require at least one digit.

### spec.passwordPolicy.requireSymbols

`bool`

Require at least one special character.

### spec.passwordPolicy.passwordHistorySize

`int32`

Number of previous passwords a new password must differ from. Range: 0-24
(0 disables history checking). Password history requires the pool to be on
the ESSENTIALS tier or higher.

- rule: {"int32":{"lte":24,"gte":0}}

### spec.passwordPolicy.temporaryPasswordValidityDays

`int32`

Number of days temporary passwords (admin-created) are valid.
Range: 0-365 (0 means no expiration). AWS default: 7.

- rule: {"int32":{"lte":365,"gte":0}}

### spec.allowedFirstAuthFactors

`[]string`

The authentication factors users may present as their FIRST factor when
signing in. This is the passwordless dial: leaving it empty keeps the
classic password-first behavior, while setting it enables the choice-based
sign-in flow (USER_AUTH) with the listed factors.

Valid values:
- "PASSWORD": classic password sign-in.
- "EMAIL_OTP": one-time code delivered by email (ESSENTIALS tier or higher).
- "SMS_OTP": one-time code delivered by SMS (requires `sms_configuration`).
- "WEB_AUTHN": passkeys/security keys (configure `web_authn_configuration`
  to pin the relying-party ID your apps expect).

### spec.mfaConfiguration

`string`

Multi-factor authentication enforcement level.
- "OFF": MFA is not used (default).
- "OPTIONAL": Users can opt in to MFA.
- "ON": MFA is required for all users.

### spec.softwareTokenMfaEnabled

`bool`

Enable TOTP-based (time-based one-time password) software token MFA. When
true, users can configure authenticator apps like Google Authenticator or
Authy. Requires `mfa_configuration` to be "OPTIONAL" or "ON".

### spec.emailMfa

`AwsCognitoUserPoolEmailMfaConfig`

Email-based MFA: Cognito emails a one-time code as the second factor.
Requires `mfa_configuration` to be "OPTIONAL" or "ON", and AWS requires
the pool to send email through SES (`email_configuration` in "DEVELOPER"
mode) because MFA codes exceed the built-in sender's delivery guarantees.

- rule: email_mfa message must contain the '{####}' placeholder where Cognito injects the code

### spec.emailMfa.message

`string`

The email body for MFA codes. Must contain the "{####}" placeholder where
Cognito injects the code. 6-20000 characters. When omitted, AWS uses its
default message.

### spec.emailMfa.subject

`string`

The email subject for MFA codes. 1-140 characters. When omitted, AWS uses
its default subject.

- rule: {"string":{"maxLen":"140"}}

### spec.webAuthn

`AwsCognitoUserPoolWebAuthnConfig`

WebAuthn (passkey / security key) relying-party configuration. Configure
this when "WEB_AUTHN" is an allowed first factor so registered passkeys
are bound to the domain your applications serve.

- rule: web_authn user_verification must be 'required' or 'preferred' when set

### spec.webAuthn.relyingPartyId

`string`

The relying-party ID passkeys are registered against -- the domain your
applications serve (e.g. "auth.example.com"). A passkey registered for one
relying-party ID does not work for another, so set this BEFORE users start
registering passkeys; changing it later invalidates existing registrations.
When omitted, Cognito uses the pool's own domain.

### spec.webAuthn.userVerification

`string`

Whether the authenticator must verify the user (PIN, biometric):
- "required": authenticators must perform user verification.
- "preferred": verification is requested but not required (AWS default).

### spec.smsConfiguration

`AwsCognitoUserPoolSmsConfig`

SMS delivery configuration. Cognito publishes SMS messages (verification
codes, MFA codes, invitations) through Amazon SNS by assuming the IAM role
referenced here. Required whenever any SMS-dependent feature is enabled:
phone_number in `auto_verified_attributes`, "SMS_OTP" as a first factor,
or SMS MFA.

### spec.smsConfiguration.snsCallerArn

`string | valueFrom` · required

The IAM role Cognito assumes to send SMS via SNS. Accepts a direct role
ARN or a reference to an AwsIamRole resource.

- references: AwsIamRole (`status.outputs.role_arn`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsIamRole, name: <that resource's name>, fieldPath: status.outputs.role_arn}} -- a bare string does not parse

### spec.smsConfiguration.externalId

`string` · required

The external ID Cognito presents when assuming the role -- the standard
confused-deputy guard. Must match the sts:ExternalId condition in the
role's trust policy.

- rule: {"string":{"minLen":"1"}}

### spec.smsConfiguration.snsRegion

`string`

The AWS region of the SNS topic-less publish (where SMS messages
originate). When omitted, AWS uses the pool's region. Set this when the
pool's region does not support SMS sending.

### spec.smsAuthenticationMessage

`string`

The SMS message sent for sign-in authentication codes. Must contain the
"{####}" placeholder where Cognito injects the code. 6-140 characters.
When omitted, AWS uses its default message.

### spec.autoVerifiedAttributes

`[]string`

Attributes to auto-verify when users sign up. Cognito sends a verification
code to these attributes. Common values: ["email"].
Valid values: "email", "phone_number".

### spec.attributesRequireVerificationBeforeUpdate

`[]string`

Attributes that must be verified before an update to them takes effect.
While the new value is pending verification, Cognito keeps the previous
value active -- without this, an unverified typo in an email update can
lock a user out of account recovery. Valid values: "email", "phone_number".

### spec.accountRecoveryMechanisms

`[]AwsCognitoUserPoolRecoveryMechanism`

Account recovery mechanisms in priority order. Each mechanism defines a
fallback method for users who forget their password.

### spec.accountRecoveryMechanisms[].name

`string` · required

Recovery method name. Valid values:
- "verified_email": send a recovery code to the verified email
- "verified_phone_number": send a recovery code via SMS
- "admin_only": only administrators can reset passwords

- rule: {"required":true}

### spec.accountRecoveryMechanisms[].priority

`int32` · required

Priority of this recovery method. 1 = primary, 2 = fallback.

- rule: {"required":true}

### spec.emailConfiguration

`AwsCognitoUserPoolEmailConfig`

Email sending configuration. Controls whether Cognito sends emails using
its built-in service (limited to 50/day in sandbox) or your verified SES
identity (production sending).

- rule: email_sending_account must be 'COGNITO_DEFAULT' or 'DEVELOPER' when set
- rule: source_arn is required when email_sending_account is 'DEVELOPER'

### spec.emailConfiguration.emailSendingAccount

`string`

Email sending mode. Valid values:
- "COGNITO_DEFAULT": Cognito sends emails using its built-in service.
  Limited to 50 emails/day in sandbox mode. No SES setup required.
- "DEVELOPER": Cognito sends emails through your verified SES identity.
  Requires `source_arn`. Supports production-level sending volumes.

### spec.emailConfiguration.sourceArn

`string | valueFrom`

SES verified identity ARN for sending emails. Required when
`email_sending_account` is "DEVELOPER". Example:
"arn:aws:ses:us-east-1:123456789012:identity/noreply@example.com"

- rule: write as {value: <literal>} or {valueFrom: {kind: <Kind>, name: <that resource's name>, fieldPath: status.outputs.<output>}} -- a bare string does not parse

### spec.emailConfiguration.fromEmailAddress

`string`

"From" email address shown to recipients. Only applicable when using
DEVELOPER mode with SES. Example: "No Reply <noreply@example.com>"

### spec.emailConfiguration.replyToEmailAddress

`string`

Reply-to email address. When set, user replies go to this address
instead of the "from" address.

### spec.emailConfiguration.configurationSet

`string`

SES configuration set name for tracking email delivery metrics.

### spec.verificationMessageTemplate

`AwsCognitoUserPoolVerificationMessageTemplate`

Templates for the verification message Cognito sends when users sign up
or change a verified attribute. Also selects between code-based and
link-based email verification.

- rule: default_email_option must be 'CONFIRM_WITH_CODE' or 'CONFIRM_WITH_LINK' when set
- rule: email_message must contain the '{####}' placeholder where Cognito injects the verification code
- rule: email_message_by_link must contain a '{##...##}' placeholder pair wrapping the link text (e.g. 'Click {##here##} to verify')
- rule: sms_message must contain the '{####}' placeholder where Cognito injects the verification code

### spec.verificationMessageTemplate.defaultEmailOption

`string`

How email verification is performed:
- "CONFIRM_WITH_CODE": the email carries a code the user types back
  (AWS default; uses `email_message`/`email_subject`).
- "CONFIRM_WITH_LINK": the email carries a click-through confirmation
  link (uses `email_message_by_link`/`email_subject_by_link`).

### spec.verificationMessageTemplate.emailMessage

`string`

Email body for code-based verification. Must contain the "{####}"
placeholder where Cognito injects the code. 6-20000 characters.

### spec.verificationMessageTemplate.emailSubject

`string`

Email subject for code-based verification. 1-140 characters.

- rule: {"string":{"maxLen":"140"}}

### spec.verificationMessageTemplate.emailMessageByLink

`string`

Email body for link-based verification. Must contain the "{##...##}"
placeholder pair wrapping the link text, e.g.
"Click {##here##} to verify your address." 6-20000 characters.

### spec.verificationMessageTemplate.emailSubjectByLink

`string`

Email subject for link-based verification. 1-140 characters.

- rule: {"string":{"maxLen":"140"}}

### spec.verificationMessageTemplate.smsMessage

`string`

SMS body for phone verification. Must contain the "{####}" placeholder.
6-140 characters.

- rule: {"string":{"maxLen":"140"}}

### spec.allowAdminCreateUserOnly

`bool`

When true, only administrators can create users -- self-registration is
disabled. Users must be created via the admin API or AWS console.

### spec.inviteMessageTemplate

`AwsCognitoUserPoolInviteMessageTemplate`

Templates for the invitation message sent to admin-created users with
their temporary credentials.

- rule: invite email_message must contain both '{username}' and '{####}' placeholders
- rule: invite sms_message must contain both '{username}' and '{####}' placeholders

### spec.inviteMessageTemplate.emailMessage

`string`

Email body for invitations. Must contain BOTH the "{username}" and
"{####}" placeholders (Cognito injects the username and the temporary
password). 6-20000 characters.

### spec.inviteMessageTemplate.emailSubject

`string`

Email subject for invitations. 1-140 characters.

- rule: {"string":{"maxLen":"140"}}

### spec.inviteMessageTemplate.smsMessage

`string`

SMS body for invitations. Must contain BOTH the "{username}" and "{####}"
placeholders. 6-140 characters.

- rule: {"string":{"maxLen":"140"}}

### spec.deviceConfiguration

`AwsCognitoUserPoolDeviceConfig`

Remembered-device configuration. When devices are remembered, users can
skip MFA on trusted devices, and sign-in events carry a device key.

### spec.deviceConfiguration.challengeRequiredOnNewDevice

`bool`

When true, a remembered device still requires a challenge (MFA) the first
time it is seen -- remembering only suppresses challenges afterwards.

### spec.deviceConfiguration.deviceOnlyRememberedOnUserPrompt

`bool`

When true, devices are remembered only after the user opts in when
prompted. When false, every device is remembered automatically.

### spec.customAttributes

`[]AwsCognitoUserPoolSchemaAttribute`

Custom user attributes beyond the standard set (email, phone, name, etc.).
Each attribute is added to the pool's schema. The schema is APPEND-ONLY:
new attributes can be added at any time, but an existing attribute can
never be modified or removed (AWS has no API for it -- removing one from
this list errors instead of recreating the pool). The `mutable` and
`required` flags are likewise fixed at the moment the attribute is added.
Maximum 50 custom attributes per pool.

- rule: attribute_data_type must be 'String', 'Number', 'DateTime', or 'Boolean'

### spec.customAttributes[].name

`string` · required

Attribute name (1-20 characters). Cognito auto-prefixes with "custom:".

- rule: {"string":{"minLen":"1","maxLen":"20"}}

### spec.customAttributes[].attributeDataType

`string` · required

Data type. Valid values: "String", "Number", "DateTime", "Boolean".

- rule: {"required":true}

### spec.customAttributes[].mutable

`bool`

Whether the attribute value can be changed after creation. Fixed at the
moment the attribute is added to the pool (the schema is append-only).

### spec.customAttributes[].required

`bool`

Whether the attribute is required during user registration. Fixed at the
moment the attribute is added to the pool.

### spec.customAttributes[].developerOnlyAttribute

`bool`

When true, only administrators (not the user themselves) can read and
write this attribute -- for backend bookkeeping like an external system's
record ID. Fixed at the moment the attribute is added.

### spec.customAttributes[].stringMinLength

`string`

Minimum string length (for String attributes). Leave at 0 for no minimum.

### spec.customAttributes[].stringMaxLength

`string`

Maximum string length (for String attributes). Leave empty for no maximum.

### spec.customAttributes[].numberMinValue

`string`

Minimum numeric value (for Number attributes). Leave empty for no minimum.

### spec.customAttributes[].numberMaxValue

`string`

Maximum numeric value (for Number attributes). Leave empty for no maximum.

### spec.lambdaConfig

`AwsCognitoUserPoolLambdaConfig`

Lambda trigger configuration for custom authentication and user lifecycle
hooks. Each trigger is an optional Lambda function ARN that Cognito invokes
at the corresponding lifecycle point.

- rule: set either pre_token_generation (V1_0 event) or pre_token_generation_config (versioned event), not both
- rule: kms_key_id is required when custom_email_sender or custom_sms_sender is set -- Cognito encrypts the delivered code with it

### spec.lambdaConfig.preSignUp

`string | valueFrom`

Invoked before user sign-up to perform custom validation or auto-confirm.

- references: AwsLambda (`status.outputs.function_arn`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsLambda, name: <that resource's name>, fieldPath: status.outputs.function_arn}} -- a bare string does not parse

### spec.lambdaConfig.preAuthentication

`string | valueFrom`

Invoked before authentication to perform custom validation.

- references: AwsLambda (`status.outputs.function_arn`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsLambda, name: <that resource's name>, fieldPath: status.outputs.function_arn}} -- a bare string does not parse

### spec.lambdaConfig.postAuthentication

`string | valueFrom`

Invoked after successful authentication for logging or custom logic.

- references: AwsLambda (`status.outputs.function_arn`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsLambda, name: <that resource's name>, fieldPath: status.outputs.function_arn}} -- a bare string does not parse

### spec.lambdaConfig.postConfirmation

`string | valueFrom`

Invoked after user confirmation to trigger welcome emails or provisioning.

- references: AwsLambda (`status.outputs.function_arn`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsLambda, name: <that resource's name>, fieldPath: status.outputs.function_arn}} -- a bare string does not parse

### spec.lambdaConfig.preTokenGeneration

`string | valueFrom`

Invoked before token generation to add, remove, or modify claims. This is
the V1_0 trigger event; use `pre_token_generation_config` instead when the
function needs the V2_0/V3_0 event shape (access-token customization).

- references: AwsLambda (`status.outputs.function_arn`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsLambda, name: <that resource's name>, fieldPath: status.outputs.function_arn}} -- a bare string does not parse

### spec.lambdaConfig.preTokenGenerationConfig

`AwsCognitoUserPoolPreTokenGenerationConfig`

Versioned pre-token-generation trigger. Unlike the plain
`pre_token_generation` field (which pins the V1_0 event), this selects the
trigger event version -- V2_0/V3_0 deliver the richer event that can also
customize ACCESS tokens, not just identity tokens. Set one or the other,
never both.

- rule: lambda_version must be 'V1_0', 'V2_0', or 'V3_0'

### spec.lambdaConfig.preTokenGenerationConfig.lambdaArn

`string | valueFrom` · required

The Lambda function to invoke. Accepts a direct ARN or a reference to an
AwsLambda resource.

- references: AwsLambda (`status.outputs.function_arn`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsLambda, name: <that resource's name>, fieldPath: status.outputs.function_arn}} -- a bare string does not parse

### spec.lambdaConfig.preTokenGenerationConfig.lambdaVersion

`string` · required

The trigger event version:
- "V1_0": identity-token customization only.
- "V2_0": adds access-token customization.
- "V3_0": V2_0 plus group/role overrides for machine-to-machine flows.
V2_0/V3_0 require the ESSENTIALS tier or higher.

- rule: {"required":true}

### spec.lambdaConfig.customMessage

`string | valueFrom`

Invoked to customize verification and invitation messages.

- references: AwsLambda (`status.outputs.function_arn`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsLambda, name: <that resource's name>, fieldPath: status.outputs.function_arn}} -- a bare string does not parse

### spec.lambdaConfig.userMigration

`string | valueFrom`

Invoked during user migration from an external identity provider.

- references: AwsLambda (`status.outputs.function_arn`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsLambda, name: <that resource's name>, fieldPath: status.outputs.function_arn}} -- a bare string does not parse

### spec.lambdaConfig.defineAuthChallenge

`string | valueFrom`

Invoked to define a custom authentication challenge (custom auth flow).

- references: AwsLambda (`status.outputs.function_arn`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsLambda, name: <that resource's name>, fieldPath: status.outputs.function_arn}} -- a bare string does not parse

### spec.lambdaConfig.createAuthChallenge

`string | valueFrom`

Invoked to create a custom authentication challenge (custom auth flow).

- references: AwsLambda (`status.outputs.function_arn`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsLambda, name: <that resource's name>, fieldPath: status.outputs.function_arn}} -- a bare string does not parse

### spec.lambdaConfig.verifyAuthChallengeResponse

`string | valueFrom`

Invoked to verify the response to a custom authentication challenge.

- references: AwsLambda (`status.outputs.function_arn`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsLambda, name: <that resource's name>, fieldPath: status.outputs.function_arn}} -- a bare string does not parse

### spec.lambdaConfig.customEmailSender

`AwsCognitoUserPoolCustomSenderConfig`

Custom email sender: Cognito calls this function to deliver email itself
(instead of SES). The event carries the message encrypted with
`kms_key_id`, so that key is required.

- rule: lambda_version must be 'V1_0'

### spec.lambdaConfig.customEmailSender.lambdaArn

`string | valueFrom` · required

The Lambda function Cognito invokes to deliver the message. Accepts a
direct ARN or a reference to an AwsLambda resource.

- references: AwsLambda (`status.outputs.function_arn`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsLambda, name: <that resource's name>, fieldPath: status.outputs.function_arn}} -- a bare string does not parse

### spec.lambdaConfig.customEmailSender.lambdaVersion

`string` · required

The sender event version. "V1_0" is the only version AWS currently
defines; modeled as a field so new versions are additive, not breaking.

- rule: {"required":true}

### spec.lambdaConfig.customSmsSender

`AwsCognitoUserPoolCustomSenderConfig`

Custom SMS sender: Cognito calls this function to deliver SMS itself
(instead of SNS). The event carries the message encrypted with
`kms_key_id`, so that key is required.

- rule: lambda_version must be 'V1_0'

### spec.lambdaConfig.customSmsSender.lambdaArn

`string | valueFrom` · required

The Lambda function Cognito invokes to deliver the message. Accepts a
direct ARN or a reference to an AwsLambda resource.

- references: AwsLambda (`status.outputs.function_arn`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsLambda, name: <that resource's name>, fieldPath: status.outputs.function_arn}} -- a bare string does not parse

### spec.lambdaConfig.customSmsSender.lambdaVersion

`string` · required

The sender event version. "V1_0" is the only version AWS currently
defines; modeled as a field so new versions are additive, not breaking.

- rule: {"required":true}

### spec.lambdaConfig.kmsKeyId

`string | valueFrom`

The KMS key Cognito uses to encrypt the code/message payload delivered to
custom sender functions. Required when either custom sender is set.
Accepts a direct key ARN or a reference to an AwsKmsKey resource.

- references: AwsKmsKey (`status.outputs.key_arn`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsKmsKey, name: <that resource's name>, fieldPath: status.outputs.key_arn}} -- a bare string does not parse

### spec.userPoolAddOns

`AwsCognitoUserPoolAddOns`

Threat protection posture (AWS "advanced security features"). Beyond
audit-only logging this requires the PLUS feature tier.

- rule: advanced_security_mode must be 'OFF', 'AUDIT', or 'ENFORCED'
- rule: custom_auth_mode must be 'AUDIT' or 'ENFORCED' when set
- rule: custom_auth_mode requires advanced_security_mode to be 'AUDIT' or 'ENFORCED'

### spec.userPoolAddOns.advancedSecurityMode

`string` · required

Threat protection mode for standard authentication flows:
- "OFF": no threat protection.
- "AUDIT": Cognito gathers risk signals and metrics without acting.
- "ENFORCED": Cognito acts on risk (blocks, requires MFA) per the pool's
  risk configuration.
AUDIT and ENFORCED require the PLUS feature tier.

- rule: {"required":true}

### spec.userPoolAddOns.customAuthMode

`string`

Extends threat protection to CUSTOM authentication flows (Lambda-driven
challenges), which standard mode does not cover. Valid values: "AUDIT",
"ENFORCED". When omitted, custom auth flows are not covered.

### spec.logConfigurations

`[]AwsCognitoUserPoolLogConfiguration`

Where Cognito delivers detailed event logs. Each entry routes one event
source (user notifications, or auth events from threat protection) to
exactly one destination: a CloudWatch log group, a Firehose stream, or an
S3 bucket. "userAuthEvents" logging requires threat protection (PLUS tier).

- rule: event_source must be 'userNotification' or 'userAuthEvents'
- rule: log_level must be 'ERROR' or 'INFO'
- rule: set exactly one destination per log configuration: cloudwatch_log_group_arn, firehose_stream_arn, or s3_bucket_arn

### spec.logConfigurations[].eventSource

`string` · required

The event source to deliver:
- "userNotification": message-delivery errors (email/SMS send failures).
  Supports the ERROR log level.
- "userAuthEvents": detailed auth events from threat protection. Requires
  the PLUS tier with advanced security enabled; supports the INFO level.

- rule: {"required":true}

### spec.logConfigurations[].logLevel

`string` · required

The log level to deliver: "ERROR" (userNotification) or "INFO"
(userAuthEvents).

- rule: {"required":true}

### spec.logConfigurations[].cloudwatchLogGroupArn

`string | valueFrom`

CloudWatch Logs destination. Accepts a direct log group ARN or a
reference to an AwsCloudwatchLogGroup resource. Exactly one destination
must be set per entry.

- references: AwsCloudwatchLogGroup (`status.outputs.log_group_arn`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsCloudwatchLogGroup, name: <that resource's name>, fieldPath: status.outputs.log_group_arn}} -- a bare string does not parse

### spec.logConfigurations[].firehoseStreamArn

`string | valueFrom`

Kinesis Data Firehose destination. Accepts a direct delivery-stream ARN
or a reference to an AwsKinesisFirehose resource.

- references: AwsKinesisFirehose (`status.outputs.delivery_stream_arn`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsKinesisFirehose, name: <that resource's name>, fieldPath: status.outputs.delivery_stream_arn}} -- a bare string does not parse

### spec.logConfigurations[].s3BucketArn

`string | valueFrom`

S3 destination. Accepts a direct bucket ARN or a reference to an
AwsS3Bucket resource.

- references: AwsS3Bucket (`status.outputs.bucket_arn`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsS3Bucket, name: <that resource's name>, fieldPath: status.outputs.bucket_arn}} -- a bare string does not parse

### spec.domain

`AwsCognitoUserPoolDomainConfig`

Domain configuration for the hosted sign-in UI and OAuth2 endpoints. When
set, Cognito provides a hosted login page at:
  https://{domain}.auth.{region}.amazoncognito.com
(for Cognito prefix domains) or at the custom domain you specify.

Required for OAuth flows that use the Authorization Code grant with a
hosted UI redirect.

### spec.domain.domain

`string` · required

Domain prefix (for Cognito-hosted domains) or fully-qualified domain name
(for custom domains). For Cognito-hosted: provide a unique prefix like
"myapp-auth" which creates "myapp-auth.auth.{region}.amazoncognito.com" --
AWS rejects prefixes containing the reserved words "aws", "amazon", or
"cognito". For custom domains: provide the FQDN like "auth.example.com".
ForceNew -- the domain cannot be changed after creation.

- rule: {"string":{"minLen":"1"}}

### spec.domain.certificateArn

`string | valueFrom`

ACM certificate ARN for custom domains. Required when `domain` contains
a dot (custom domain). The certificate must be in us-east-1 regardless
of the user pool's region (Cognito uses CloudFront for custom domains).
Not applicable for Cognito-hosted prefix domains. Adding or removing the
certificate replaces the domain (switching custom <-> prefix); rotating
one certificate ARN to another updates in place.

- references: AwsCertManagerCert (`status.outputs.cert_arn`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsCertManagerCert, name: <that resource's name>, fieldPath: status.outputs.cert_arn}} -- a bare string does not parse

### spec.domain.managedLoginVersion

`int32` · optional (explicit presence)

The sign-in page generation served on this domain:
- 1: classic hosted UI.
- 2: managed login (the newer, branding-designer experience; requires a
     branding style to be assigned before sign-in pages render).
When omitted, AWS defaults new domains to managed login (2).

- rule: {"int32":{"lte":2,"gte":1}}

## Validation Rules

- `username_or_alias_attributes`: username_attributes and alias_attributes are mutually exclusive; set one or neither, not both
- `username_attributes_valid`: username_attributes must contain only 'email' and/or 'phone_number'
- `alias_attributes_valid`: alias_attributes must contain only 'email', 'phone_number', and/or 'preferred_username'
- `user_pool_tier_valid`: user_pool_tier must be 'LITE', 'ESSENTIALS', or 'PLUS' when set
- `allowed_first_auth_factors_valid`: allowed_first_auth_factors must contain only 'PASSWORD', 'EMAIL_OTP', 'SMS_OTP', and/or 'WEB_AUTHN'
- `auto_verified_attributes_valid`: auto_verified_attributes must contain only 'email' and/or 'phone_number'
- `attributes_require_verification_valid`: attributes_require_verification_before_update must contain only 'email' and/or 'phone_number'
- `mfa_configuration_valid`: mfa_configuration must be 'OFF', 'OPTIONAL', or 'ON' when set
- `software_token_mfa_requires_mfa`: software_token_mfa_enabled requires mfa_configuration to be 'OPTIONAL' or 'ON'
- `email_mfa_requires_mfa`: email_mfa requires mfa_configuration to be 'OPTIONAL' or 'ON'
- `sms_otp_requires_sms_configuration`: allowing 'SMS_OTP' as a first auth factor requires sms_configuration so Cognito can deliver the codes
- `sms_authentication_message_placeholder`: sms_authentication_message must contain the '{####}' placeholder where Cognito injects the code
- `account_recovery_name_valid`: account_recovery_mechanisms name must be 'verified_email', 'verified_phone_number', or 'admin_only'
- `account_recovery_priority_range`: account_recovery_mechanisms priority must be 1 or 2
- `account_recovery_admin_only_exclusive`: 'admin_only' cannot be combined with other account recovery mechanisms
- `custom_domain_requires_certificate`: a custom domain (containing '.') requires certificate_arn to be set
- `prefix_domain_no_reserved_words`: a Cognito-hosted prefix domain cannot contain the reserved words 'aws', 'amazon', or 'cognito'

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AwsCognitoUserPool, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.user_pool_id` | `string` | The user pool identifier. Primary reference used in SDK calls, IAM policies, and as the user_pool_id input for app clients, identity providers, and resource servers. Format: "{region}_{poolId}" (e.g., "us-east-1_Ab1Cd2EfG"). |
| `status.outputs.user_pool_arn` | `string` | The Amazon Resource Name of the user pool. Used in IAM policies, cross-service permissions, and ALB authenticate-cognito actions. |
| `status.outputs.user_pool_endpoint` | `string` | The pool's endpoint as AWS reports it -- host and path WITHOUT a scheme: "cognito-idp.{region}.amazonaws.com/{user_pool_id}". Use `issuer` when a consumer needs the full OIDC issuer URL. |
| `status.outputs.issuer` | `string` | The OIDC issuer URL for tokens minted by this pool: "https://cognito-idp.{region}.amazonaws.com/{user_pool_id}". This is the value JWT authorizers (API Gateway, ALB) and OIDC client libraries validate the token's "iss" claim against. |
| `status.outputs.user_pool_domain` | `string` | The hosted-UI domain exactly as configured on the pool -- the prefix of a Cognito-hosted domain (e.g. "myapp-auth") or the full custom domain (e.g. "auth.example.com"), with no scheme. This is the join key ALB authenticate-cognito actions take as their user_pool_domain. Empty when no domain is configured. |
| `status.outputs.hosted_ui_url` | `string` | The full https:// URL of the hosted sign-in UI: "https://{prefix}.auth.{region}.amazoncognito.com" for Cognito-hosted domains, or "https://{custom_domain}" for custom domains. Empty when no domain is configured. |
| `status.outputs.cloudfront_distribution` | `string` | The CloudFront distribution DOMAIN NAME fronting a custom domain (e.g. "d111abcdef8.cloudfront.net"). Point the custom domain's DNS at this value -- with AwsRoute53DnsRecord, an alias A record composed from this output and `cloudfront_hosted_zone_id`. Empty for Cognito-hosted prefix domains or when no domain is configured. |
| `status.outputs.cloudfront_distribution_arn` | `string` | The ARN of the CloudFront distribution fronting a custom domain. Useful for IAM policies and tag-based governance. Empty for prefix domains. |
| `status.outputs.cloudfront_hosted_zone_id` | `string` | The Route53 hosted-zone ID of the CloudFront distribution (the fixed CloudFront zone) -- the alias target zone for DNS alias records pointing the custom domain at CloudFront. Empty for prefix domains. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.smsConfiguration.snsCallerArn` | AwsIamRole | `status.outputs.role_arn` |
| `spec.lambdaConfig.preSignUp` | AwsLambda | `status.outputs.function_arn` |
| `spec.lambdaConfig.preAuthentication` | AwsLambda | `status.outputs.function_arn` |
| `spec.lambdaConfig.postAuthentication` | AwsLambda | `status.outputs.function_arn` |
| `spec.lambdaConfig.postConfirmation` | AwsLambda | `status.outputs.function_arn` |
| `spec.lambdaConfig.preTokenGeneration` | AwsLambda | `status.outputs.function_arn` |
| `spec.lambdaConfig.preTokenGenerationConfig.lambdaArn` | AwsLambda | `status.outputs.function_arn` |
| `spec.lambdaConfig.customMessage` | AwsLambda | `status.outputs.function_arn` |
| `spec.lambdaConfig.userMigration` | AwsLambda | `status.outputs.function_arn` |
| `spec.lambdaConfig.defineAuthChallenge` | AwsLambda | `status.outputs.function_arn` |
| `spec.lambdaConfig.createAuthChallenge` | AwsLambda | `status.outputs.function_arn` |
| `spec.lambdaConfig.verifyAuthChallengeResponse` | AwsLambda | `status.outputs.function_arn` |
| `spec.lambdaConfig.customEmailSender.lambdaArn` | AwsLambda | `status.outputs.function_arn` |
| `spec.lambdaConfig.customSmsSender.lambdaArn` | AwsLambda | `status.outputs.function_arn` |
| `spec.lambdaConfig.kmsKeyId` | AwsKmsKey | `status.outputs.key_arn` |
| `spec.logConfigurations[].cloudwatchLogGroupArn` | AwsCloudwatchLogGroup | `status.outputs.log_group_arn` |
| `spec.logConfigurations[].firehoseStreamArn` | AwsKinesisFirehose | `status.outputs.delivery_stream_arn` |
| `spec.logConfigurations[].s3BucketArn` | AwsS3Bucket | `status.outputs.bucket_arn` |
| `spec.domain.certificateArn` | AwsCertManagerCert | `status.outputs.cert_arn` |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| AwsCognitoIdentityProvider | `spec.userPoolId` | `status.outputs.user_pool_id` |
| AwsCognitoResourceServer | `spec.userPoolId` | `status.outputs.user_pool_id` |
| AwsCognitoUserPoolClient | `spec.userPoolId` | `status.outputs.user_pool_id` |
| AwsHttpApiGateway | `spec.authorizers[].jwtConfiguration.issuer` | `status.outputs.issuer` |
| AwsLbListener | `spec.defaultActions[].authenticateCognito.userPoolArn` | `status.outputs.user_pool_arn` |
| AwsLbListener | `spec.defaultActions[].authenticateCognito.userPoolDomain` | `status.outputs.user_pool_domain` |
| AwsLbListenerRule | `spec.actions[].authenticateCognito.userPoolArn` | `status.outputs.user_pool_arn` |
| AwsLbListenerRule | `spec.actions[].authenticateCognito.userPoolDomain` | `status.outputs.user_pool_domain` |
| AwsOpenSearchDomain | `spec.cognitoOptions.userPoolId` | `status.outputs.user_pool_id` |

## See Also

- [Overview](../README.md)
