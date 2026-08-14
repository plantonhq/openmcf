# AwsConfigRule

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `aws.planton.dev/v1alpha1`

**Guide**: [GUIDE.md](../GUIDE.md) -- authored operational judgment for this component: conventions, trade-offs, and what pairs well with it.

## Example

```yaml
# Canonical AwsConfigRule example (hack/dev manifest and refgen Example
# source): a managed rule scoped to S3 buckets with non-automatic SSM
# remediation. Literal values stand in for composed references so the
# offline `tofu plan` renders every arm.
apiVersion: aws.planton.dev/v1alpha1
kind: AwsConfigRule
metadata:
  name: s3-bucket-versioning-enabled
  id: s3-bucket-versioning-enabled
  org: test-org
  env: dev
spec:
  region: us-west-2
  description: Every S3 bucket must have versioning enabled
  managed:
    ruleIdentifier: S3_BUCKET_VERSIONING_ENABLED
  scope:
    complianceResourceTypes:
      - AWS::S3::Bucket
  evaluationModes:
    - DETECTIVE
  remediation:
    automatic: false
    targetId: AWS-ConfigureS3BucketVersioning
    resourceType: AWS::S3::Bucket
    parameters:
      - name: BucketName
        resourceValue: RESOURCE_ID
      - name: AutomationAssumeRole
        staticValue: arn:aws:iam::123456789012:role/config-remediation
    concurrentExecutionRatePercentage: 10
    errorPercentage: 50
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.description` | `string` |  |  |  |
| `spec.inputParameters` | `string` |  |  |  |
| `spec.maximumExecutionFrequency` | `string` |  |  |  |
| `spec.managed` | `AwsConfigRuleManagedSource` |  |  |  |
| `spec.managed.ruleIdentifier` | `string` | yes |  |  |
| `spec.customLambda` | `AwsConfigRuleCustomLambdaSource` |  |  |  |
| `spec.customLambda.functionArn` | `string \| valueFrom` | yes |  | AwsLambda (`status.outputs.function_arn`) |
| `spec.customLambda.sourceDetails` | `[]AwsConfigRuleSourceDetail` |  |  |  |
| `spec.customLambda.sourceDetails[].messageType` | `string` |  |  |  |
| `spec.customLambda.sourceDetails[].maximumExecutionFrequency` | `string` |  |  |  |
| `spec.customPolicy` | `AwsConfigRuleCustomPolicySource` |  |  |  |
| `spec.customPolicy.policyRuntime` | `string` |  |  |  |
| `spec.customPolicy.policyText` | `string` | yes |  |  |
| `spec.customPolicy.enableDebugLogDelivery` | `bool` |  |  |  |
| `spec.scope` | `AwsConfigRuleScope` |  |  |  |
| `spec.scope.complianceResourceId` | `string` |  |  |  |
| `spec.scope.complianceResourceTypes` | `[]string` |  |  |  |
| `spec.scope.tagKey` | `string` |  |  |  |
| `spec.scope.tagValue` | `string` |  |  |  |
| `spec.evaluationModes` | `[]string` |  |  |  |
| `spec.organization` | `AwsConfigRuleOrganization` |  |  |  |
| `spec.organization.excludedAccounts` | `[]string` |  |  |  |
| `spec.organization.triggerTypes` | `[]string` |  |  |  |
| `spec.organization.debugLogDeliveryAccounts` | `[]string` |  |  |  |
| `spec.remediation` | `AwsConfigRuleRemediation` |  |  |  |
| `spec.remediation.automatic` | `bool` |  |  |  |
| `spec.remediation.targetId` | `string` | yes |  |  |
| `spec.remediation.targetVersion` | `string` |  |  |  |
| `spec.remediation.resourceType` | `string` |  |  |  |
| `spec.remediation.parameters` | `[]AwsConfigRuleRemediationParameter` |  |  |  |
| `spec.remediation.parameters[].name` | `string` | yes |  |  |
| `spec.remediation.parameters[].resourceValue` | `string` |  |  |  |
| `spec.remediation.parameters[].staticValue` | `string` |  |  |  |
| `spec.remediation.parameters[].staticValues` | `[]string` |  |  |  |
| `spec.remediation.maximumAutomaticAttempts` | `int32` |  |  |  |
| `spec.remediation.retryAttemptSeconds` | `int64` |  |  |  |
| `spec.remediation.concurrentExecutionRatePercentage` | `int32` |  |  |  |
| `spec.remediation.errorPercentage` | `int32` |  |  |  |

## Field Details

### spec.region

`string` · required

- rule: {"string":{"minLen":"1"}}

### spec.description

`string`

- rule: {"string":{"maxLen":"256"}}

### spec.inputParameters

`string`

- rule: {"string":{"maxLen":"2048"}}

### spec.maximumExecutionFrequency

`string`

- rule: {"string":{"in":["","One_Hour","Three_Hours","Six_Hours","Twelve_Hours","TwentyFour_Hours"]}}

### spec.managed

`AwsConfigRuleManagedSource`

### spec.managed.ruleIdentifier

`string` · required

- rule: {"string":{"minLen":"1","maxLen":"256"}}

### spec.customLambda

`AwsConfigRuleCustomLambdaSource`

### spec.customLambda.functionArn

`string | valueFrom` · required

- references: AwsLambda (`status.outputs.function_arn`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsLambda, name: <that resource's name>, fieldPath: status.outputs.function_arn}} -- a bare string does not parse

### spec.customLambda.sourceDetails

`[]AwsConfigRuleSourceDetail`

- rule: {"repeated":{"maxItems":"25"}}

### spec.customLambda.sourceDetails[].messageType

`string`

- rule: {"string":{"in":["ConfigurationItemChangeNotification","OversizedConfigurationItemChangeNotification","ScheduledNotification","ConfigurationSnapshotDeliveryCompleted"]}}

### spec.customLambda.sourceDetails[].maximumExecutionFrequency

`string`

- rule: {"string":{"in":["","One_Hour","Three_Hours","Six_Hours","Twelve_Hours","TwentyFour_Hours"]}}

### spec.customPolicy

`AwsConfigRuleCustomPolicySource`

### spec.customPolicy.policyRuntime

`string`

- rule: {"string":{"pattern":"^guard\\-2\\.x\\.x$"}}

### spec.customPolicy.policyText

`string` · required

- rule: {"string":{"minLen":"1","maxLen":"10000"}}

### spec.customPolicy.enableDebugLogDelivery

`bool`

### spec.scope

`AwsConfigRuleScope`

- rule: compliance_resource_id requires exactly one compliance_resource_types entry
- rule: tag_value requires tag_key

### spec.scope.complianceResourceId

`string`

- rule: {"string":{"maxLen":"768"}}

### spec.scope.complianceResourceTypes

`[]string`

- rule: {"repeated":{"maxItems":"100","unique":true}}

### spec.scope.tagKey

`string`

- rule: {"string":{"maxLen":"128"}}

### spec.scope.tagValue

`string`

- rule: {"string":{"maxLen":"256"}}

### spec.evaluationModes

`[]string`

- rule: {"repeated":{"unique":true,"items":{"string":{"in":["DETECTIVE","PROACTIVE"]}}}}

### spec.organization

`AwsConfigRuleOrganization`

### spec.organization.excludedAccounts

`[]string`

- rule: {"repeated":{"maxItems":"1000","unique":true,"items":{"string":{"pattern":"^[0-9]{12}$"}}}}

### spec.organization.triggerTypes

`[]string`

- rule: {"repeated":{"maxItems":"3","unique":true,"items":{"string":{"in":["ConfigurationItemChangeNotification","OversizedConfigurationItemChangeNotification","ScheduledNotification"]}}}}

### spec.organization.debugLogDeliveryAccounts

`[]string`

- rule: {"repeated":{"maxItems":"1000","unique":true,"items":{"string":{"pattern":"^[0-9]{12}$"}}}}

### spec.remediation

`AwsConfigRuleRemediation`

- rule: automatic remediation requires maximum_automatic_attempts and retry_attempt_seconds

### spec.remediation.automatic

`bool`

### spec.remediation.targetId

`string` · required

- rule: {"string":{"minLen":"1","maxLen":"256"}}

### spec.remediation.targetVersion

`string`

### spec.remediation.resourceType

`string`

### spec.remediation.parameters

`[]AwsConfigRuleRemediationParameter`

- rule: {"repeated":{"maxItems":"25"}}
- rule: set exactly one of resource_value, static_value, static_values

### spec.remediation.parameters[].name

`string` · required

- rule: {"string":{"minLen":"1"}}

### spec.remediation.parameters[].resourceValue

`string`

- rule: {"string":{"in":["","RESOURCE_ID"]}}

### spec.remediation.parameters[].staticValue

`string`

### spec.remediation.parameters[].staticValues

`[]string`

### spec.remediation.maximumAutomaticAttempts

`int32`

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","int32":{"lte":25,"gte":1}}

### spec.remediation.retryAttemptSeconds

`int64`

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","int64":{"lte":"2678000","gte":"1"}}

### spec.remediation.concurrentExecutionRatePercentage

`int32`

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","int32":{"lte":100,"gte":1}}

### spec.remediation.errorPercentage

`int32`

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","int32":{"lte":100,"gte":1}}

## Validation Rules

- `exactly_one_source`: set exactly one of managed, custom_lambda, custom_policy
- `org_trigger_types_by_source`: organization.trigger_types is required for custom_lambda/custom_policy organization rules and must be empty for managed ones
- `org_custom_policy_no_scheduled_trigger`: organization custom_policy rules do not accept the ScheduledNotification trigger
- `debug_accounts_custom_policy_only`: organization.debug_log_delivery_accounts applies only to custom_policy rules
- `org_lambda_uses_trigger_types`: organization custom_lambda rules declare triggers via organization.trigger_types, not custom_lambda.source_details
- `remediation_account_scope_only`: remediation is not supported on organization rules
- `evaluation_modes_account_scope_only`: evaluation_modes apply only to account-scoped rules

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AwsConfigRule, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.rule_arn` | `string` |  |
| `status.outputs.rule_name` | `string` |  |
| `status.outputs.rule_id` | `string` |  |
| `status.outputs.remediation_arn` | `string` |  |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.customLambda.functionArn` | AwsLambda | `status.outputs.function_arn` |

## See Also

- [Overview](../README.md)
