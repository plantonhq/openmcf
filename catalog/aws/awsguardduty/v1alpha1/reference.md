# AwsGuardDuty

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `aws.planton.dev/v1alpha1`

**Guide**: [GUIDE.md](../GUIDE.md) -- authored operational judgment for this component: conventions, trade-offs, and what pairs well with it.

## Example

```yaml
# Canonical AwsGuardDuty example (hack/dev manifest and refgen Example
# source): the region's threat-detection posture -- detector, protection
# plans, a noise filter, trusted/threat lists, and findings export.
# Literal values stand in for composed references so the offline `tofu
# plan` renders every arm.
apiVersion: aws.planton.dev/v1alpha1
kind: AwsGuardDuty
metadata:
  name: guardduty-us-west-2
  id: guardduty-us-west-2
  org: test-org
  env: dev
spec:
  region: us-west-2
  findingPublishingFrequency: FIFTEEN_MINUTES
  features:
    - name: S3_DATA_EVENTS
    - name: RUNTIME_MONITORING
      additionalConfiguration:
        - name: EC2_AGENT_MANAGEMENT
          enabled: false
  filters:
    - name: archive-sandbox-low-severity
      description: Archive low-severity findings from the sandbox account
      action: ARCHIVE
      rank: 1
      criteria:
        - field: severity
          lessThan: "4"
        - field: accountId
          equals: ["123456789012"]
  ipSets:
    - name: trusted-office
      format: TXT
      location: https://s3.amazonaws.com/my-security-lists/trusted.txt
      activate: true
  threatIntelSets:
    - name: known-bad
      format: TXT
      location: https://s3.amazonaws.com/my-security-lists/bad.txt
      activate: true
  publishingDestination:
    bucketArn:
      value: arn:aws:s3:::my-findings-archive
    kmsKeyArn:
      value: arn:aws:kms:us-west-2:123456789012:key/1234abcd-12ab-34cd-56ef-1234567890ab
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.enable` | `bool` |  |  |  |
| `spec.findingPublishingFrequency` | `string` |  |  |  |
| `spec.features` | `[]AwsGuardDutyFeature` |  |  |  |
| `spec.features[].name` | `string` |  |  |  |
| `spec.features[].enabled` | `bool` |  |  |  |
| `spec.features[].additionalConfiguration` | `[]AwsGuardDutyFeatureAdditionalConfiguration` |  |  |  |
| `spec.features[].additionalConfiguration[].name` | `string` |  |  |  |
| `spec.features[].additionalConfiguration[].enabled` | `bool` |  |  |  |
| `spec.filters` | `[]AwsGuardDutyFilter` |  |  |  |
| `spec.filters[].name` | `string` | yes |  |  |
| `spec.filters[].description` | `string` |  |  |  |
| `spec.filters[].action` | `string` |  |  |  |
| `spec.filters[].rank` | `int32` |  |  |  |
| `spec.filters[].criteria` | `[]AwsGuardDutyFilterCriterion` | yes |  |  |
| `spec.filters[].criteria[].field` | `string` | yes |  |  |
| `spec.filters[].criteria[].equals` | `[]string` |  |  |  |
| `spec.filters[].criteria[].notEquals` | `[]string` |  |  |  |
| `spec.filters[].criteria[].matches` | `[]string` |  |  |  |
| `spec.filters[].criteria[].notMatches` | `[]string` |  |  |  |
| `spec.filters[].criteria[].greaterThan` | `string` |  |  |  |
| `spec.filters[].criteria[].greaterThanOrEqual` | `string` |  |  |  |
| `spec.filters[].criteria[].lessThan` | `string` |  |  |  |
| `spec.filters[].criteria[].lessThanOrEqual` | `string` |  |  |  |
| `spec.ipSets` | `[]AwsGuardDutyIpSet` |  |  |  |
| `spec.ipSets[].name` | `string` | yes |  |  |
| `spec.ipSets[].format` | `string` |  |  |  |
| `spec.ipSets[].location` | `string` | yes |  |  |
| `spec.ipSets[].activate` | `bool` |  |  |  |
| `spec.threatIntelSets` | `[]AwsGuardDutyThreatIntelSet` |  |  |  |
| `spec.threatIntelSets[].name` | `string` | yes |  |  |
| `spec.threatIntelSets[].format` | `string` |  |  |  |
| `spec.threatIntelSets[].location` | `string` | yes |  |  |
| `spec.threatIntelSets[].activate` | `bool` |  |  |  |
| `spec.publishingDestination` | `AwsGuardDutyPublishingDestination` |  |  |  |
| `spec.publishingDestination.bucketArn` | `string \| valueFrom` | yes |  | AwsS3Bucket (`status.outputs.bucket_arn`) |
| `spec.publishingDestination.kmsKeyArn` | `string \| valueFrom` | yes |  | AwsKmsKey (`status.outputs.key_arn`) |
| `spec.organization` | `AwsGuardDutyOrganization` |  |  |  |
| `spec.organization.adminAccountId` | `string` |  |  |  |
| `spec.organization.autoEnableOrganizationMembers` | `string` |  |  |  |
| `spec.organization.features` | `[]AwsGuardDutyOrganizationFeature` |  |  |  |
| `spec.organization.features[].name` | `string` |  |  |  |
| `spec.organization.features[].autoEnable` | `string` |  |  |  |
| `spec.organization.features[].additionalConfiguration` | `[]AwsGuardDutyOrganizationFeatureAdditionalConfiguration` |  |  |  |
| `spec.organization.features[].additionalConfiguration[].name` | `string` |  |  |  |
| `spec.organization.features[].additionalConfiguration[].autoEnable` | `string` |  |  |  |
| `spec.members` | `[]AwsGuardDutyMember` |  |  |  |
| `spec.members[].accountId` | `string` |  |  |  |
| `spec.members[].email` | `string` | yes |  |  |
| `spec.members[].invite` | `bool` |  |  |  |
| `spec.members[].invitationMessage` | `string` |  |  |  |
| `spec.members[].disableEmailNotification` | `bool` |  |  |  |
| `spec.members[].features` | `[]AwsGuardDutyMemberFeature` |  |  |  |
| `spec.members[].features[].name` | `string` |  |  |  |
| `spec.members[].features[].enabled` | `bool` |  |  |  |
| `spec.members[].features[].additionalConfiguration` | `[]AwsGuardDutyFeatureAdditionalConfiguration` |  |  |  |
| `spec.members[].features[].additionalConfiguration[].name` | `string` |  |  |  |
| `spec.members[].features[].additionalConfiguration[].enabled` | `bool` |  |  |  |
| `spec.acceptInvitationFromAccountId` | `string` |  |  |  |

## Field Details

### spec.region

`string` · required

- rule: {"string":{"minLen":"1"}}

### spec.enable

`bool` · optional (explicit presence)

### spec.findingPublishingFrequency

`string`

- rule: {"string":{"in":["","FIFTEEN_MINUTES","ONE_HOUR","SIX_HOURS"]}}

### spec.features

`[]AwsGuardDutyFeature`

- rule: additional_configuration entries must have unique names

### spec.features[].name

`string`

- rule: {"string":{"in":["S3_DATA_EVENTS","EKS_AUDIT_LOGS","EBS_MALWARE_PROTECTION","RDS_LOGIN_EVENTS","LAMBDA_NETWORK_LOGS","EKS_RUNTIME_MONITORING","RUNTIME_MONITORING","AI_PROTECTION","AI_ANALYST"]}}

### spec.features[].enabled

`bool` · optional (explicit presence)

### spec.features[].additionalConfiguration

`[]AwsGuardDutyFeatureAdditionalConfiguration`

### spec.features[].additionalConfiguration[].name

`string`

- rule: {"string":{"in":["EKS_ADDON_MANAGEMENT","ECS_FARGATE_AGENT_MANAGEMENT","EC2_AGENT_MANAGEMENT"]}}

### spec.features[].additionalConfiguration[].enabled

`bool` · optional (explicit presence)

### spec.filters

`[]AwsGuardDutyFilter`

- rule: criteria entries must have unique field values

### spec.filters[].name

`string` · required

- rule: {"string":{"minLen":"3","maxLen":"64","pattern":"^[a-zA-Z0-9_.-]+$"}}

### spec.filters[].description

`string`

- rule: {"string":{"maxLen":"512"}}

### spec.filters[].action

`string`

- rule: {"string":{"in":["NOOP","ARCHIVE"]}}

### spec.filters[].rank

`int32`

- rule: {"int32":{"lte":100,"gte":1}}

### spec.filters[].criteria

`[]AwsGuardDutyFilterCriterion` · required

- rule: {"repeated":{"minItems":"1"}}
- rule: set at least one condition (equals, not_equals, matches, not_matches, or a bound)

### spec.filters[].criteria[].field

`string` · required

- rule: {"string":{"minLen":"1"}}

### spec.filters[].criteria[].equals

`[]string`

### spec.filters[].criteria[].notEquals

`[]string`

### spec.filters[].criteria[].matches

`[]string`

- rule: {"repeated":{"maxItems":"5"}}

### spec.filters[].criteria[].notMatches

`[]string`

- rule: {"repeated":{"maxItems":"5"}}

### spec.filters[].criteria[].greaterThan

`string`

### spec.filters[].criteria[].greaterThanOrEqual

`string`

### spec.filters[].criteria[].lessThan

`string`

### spec.filters[].criteria[].lessThanOrEqual

`string`

### spec.ipSets

`[]AwsGuardDutyIpSet`

### spec.ipSets[].name

`string` · required

- rule: {"string":{"minLen":"1","maxLen":"300"}}

### spec.ipSets[].format

`string`

- rule: {"string":{"in":["TXT","STIX","OTX_CSV","ALIEN_VAULT","PROOF_POINT","FIRE_EYE"]}}

### spec.ipSets[].location

`string` · required

- rule: {"string":{"minLen":"1"}}

### spec.ipSets[].activate

`bool`

### spec.threatIntelSets

`[]AwsGuardDutyThreatIntelSet`

### spec.threatIntelSets[].name

`string` · required

- rule: {"string":{"minLen":"1","maxLen":"300"}}

### spec.threatIntelSets[].format

`string`

- rule: {"string":{"in":["TXT","STIX","OTX_CSV","ALIEN_VAULT","PROOF_POINT","FIRE_EYE"]}}

### spec.threatIntelSets[].location

`string` · required

- rule: {"string":{"minLen":"1"}}

### spec.threatIntelSets[].activate

`bool`

### spec.publishingDestination

`AwsGuardDutyPublishingDestination`

### spec.publishingDestination.bucketArn

`string | valueFrom` · required

- references: AwsS3Bucket (`status.outputs.bucket_arn`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsS3Bucket, name: <that resource's name>, fieldPath: status.outputs.bucket_arn}} -- a bare string does not parse

### spec.publishingDestination.kmsKeyArn

`string | valueFrom` · required

- references: AwsKmsKey (`status.outputs.key_arn`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsKmsKey, name: <that resource's name>, fieldPath: status.outputs.key_arn}} -- a bare string does not parse

### spec.organization

`AwsGuardDutyOrganization`

- rule: features entries must have unique names

### spec.organization.adminAccountId

`string`

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"pattern":"^[0-9]{12}$"}}

### spec.organization.autoEnableOrganizationMembers

`string`

- rule: {"string":{"in":["NEW","ALL","NONE"]}}

### spec.organization.features

`[]AwsGuardDutyOrganizationFeature`

- rule: additional_configuration entries must have unique names

### spec.organization.features[].name

`string`

- rule: {"string":{"in":["S3_DATA_EVENTS","EKS_AUDIT_LOGS","EBS_MALWARE_PROTECTION","RDS_LOGIN_EVENTS","LAMBDA_NETWORK_LOGS","EKS_RUNTIME_MONITORING","RUNTIME_MONITORING","AI_PROTECTION"]}}

### spec.organization.features[].autoEnable

`string`

- rule: {"string":{"in":["NEW","ALL","NONE"]}}

### spec.organization.features[].additionalConfiguration

`[]AwsGuardDutyOrganizationFeatureAdditionalConfiguration`

### spec.organization.features[].additionalConfiguration[].name

`string`

- rule: {"string":{"in":["EKS_ADDON_MANAGEMENT","ECS_FARGATE_AGENT_MANAGEMENT","EC2_AGENT_MANAGEMENT"]}}

### spec.organization.features[].additionalConfiguration[].autoEnable

`string`

- rule: {"string":{"in":["NEW","ALL","NONE"]}}

### spec.members

`[]AwsGuardDutyMember`

- rule: features entries must have unique names

### spec.members[].accountId

`string`

- rule: {"string":{"pattern":"^[0-9]{12}$"}}

### spec.members[].email

`string` · required

- rule: {"string":{"minLen":"1","maxLen":"64"}}

### spec.members[].invite

`bool` · optional (explicit presence)

### spec.members[].invitationMessage

`string`

### spec.members[].disableEmailNotification

`bool`

### spec.members[].features

`[]AwsGuardDutyMemberFeature`

- rule: additional_configuration entries must have unique names

### spec.members[].features[].name

`string`

- rule: {"string":{"in":["S3_DATA_EVENTS","EKS_AUDIT_LOGS","EBS_MALWARE_PROTECTION","RDS_LOGIN_EVENTS","LAMBDA_NETWORK_LOGS","EKS_RUNTIME_MONITORING","RUNTIME_MONITORING","AI_PROTECTION"]}}

### spec.members[].features[].enabled

`bool` · optional (explicit presence)

### spec.members[].features[].additionalConfiguration

`[]AwsGuardDutyFeatureAdditionalConfiguration`

### spec.members[].features[].additionalConfiguration[].name

`string`

- rule: {"string":{"in":["EKS_ADDON_MANAGEMENT","ECS_FARGATE_AGENT_MANAGEMENT","EC2_AGENT_MANAGEMENT"]}}

### spec.members[].features[].additionalConfiguration[].enabled

`bool` · optional (explicit presence)

### spec.acceptInvitationFromAccountId

`string`

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"pattern":"^[0-9]{12}$"}}

## Validation Rules

- `feature_names_unique`: features entries must have unique names
- `filter_names_unique`: filters entries must have unique names
- `ip_set_names_unique`: ip_sets entries must have unique names
- `threat_intel_set_names_unique`: threat_intel_sets entries must have unique names
- `member_accounts_unique`: members entries must have unique account_id values
- `member_xor_admin_posture`: accept_invitation_from_account_id is mutually exclusive with organization and members (an account is either a member or an administrator)

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AwsGuardDuty, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.detector_id` | `string` |  |
| `status.outputs.detector_arn` | `string` |  |
| `status.outputs.account_id` | `string` |  |
| `status.outputs.ip_set_ids` | `map<string, string>` |  |
| `status.outputs.threat_intel_set_ids` | `map<string, string>` |  |
| `status.outputs.publishing_destination_id` | `string` |  |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.publishingDestination.bucketArn` | AwsS3Bucket | `status.outputs.bucket_arn` |
| `spec.publishingDestination.kmsKeyArn` | AwsKmsKey | `status.outputs.key_arn` |

## See Also

- [Overview](../README.md)
