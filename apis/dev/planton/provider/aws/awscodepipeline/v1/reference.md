# AwsCodePipeline

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `aws.planton.dev/v1`

AwsCodePipelineSpec defines the specification for an AWS CodePipeline
continuous delivery pipeline.

AWS CodePipeline is a fully managed continuous delivery service that
orchestrates build, test, and deploy phases of your release process.
A pipeline is an ordered sequence of stages, each containing one or more
actions that perform tasks such as fetching source code, running builds,
executing tests, or deploying to production.

This component creates a pipeline with artifact stores, stages, actions,
and the V2 feature set: git-based triggers for automatic execution,
pipeline-level variables for parameterization, and stage-level conditions
(entry gates, success checks, and failure handling with automatic
rollback or retry).

Bundling rationale: A pipeline is the top-level orchestration unit. Stages
and actions are integral parts of the pipeline definition and cannot exist
independently. Webhooks are excluded because V2 pipelines use native
triggers (superior to the legacy webhook mechanism). Custom action types
are account-level resources with independent lifecycles and are excluded.

## Example

```yaml
# Full-surface offline-plan manifest. This exercises every optional block the
# module renders -- V2 triggers with push + pull-request filters, pipeline
# variables, per-action timeouts/regions/roles, KMS-encrypted artifact store,
# and all three stage-condition arms (entry gate, success check, failure
# rollback) -- so the `tofu plan` / `pulumi preview` proofs cover the whole
# contract.
apiVersion: aws.planton.dev/v1
kind: AwsCodePipeline
metadata:
  name: test-pipeline
  id: test-pipeline
  org: test-org
  env: dev
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: test
    pulumi.planton.dev/project: test
    pulumi.planton.dev/stack.name: dev.AwsCodePipeline.test-pipeline
spec:
  region: us-west-2
  pipelineType: V2
  executionMode: QUEUED
  roleArn:
    value: arn:aws:iam::123456789012:role/codepipeline-role
  artifactStores:
    - location:
        value: my-artifact-bucket
      encryptionKeyId:
        value: arn:aws:kms:us-west-2:123456789012:key/1234abcd-12ab-34cd-56ef-1234567890ab
  variables:
    - name: DeployEnvironment
      defaultValue: staging
      description: "Target environment consumed by deploy actions as #{variables.DeployEnvironment}"
  triggers:
    - providerType: CodeStarSourceConnection
      gitConfiguration:
        sourceActionName: SourceAction
        push:
          - branches:
              includes:
                - main
                - release/*
              excludes:
                - release/*-hotfix
            filePaths:
              includes:
                - src/**
            tags:
              includes:
                - v*
        pullRequest:
          - events:
              - OPEN
              - UPDATED
            branches:
              includes:
                - main
            filePaths:
              includes:
                - src/**
  stages:
    - name: Source
      actions:
        - name: SourceAction
          category: Source
          owner: AWS
          provider: CodeStarSourceConnection
          version: "1"
          namespace: SourceVars
          outputArtifacts:
            - SourceOutput
          configuration:
            ConnectionArn: arn:aws:codeconnections:us-west-2:123456789012:connection/aaaa1111-22bb-33cc-44dd-5555eeee6666
            FullRepositoryId: my-org/my-repo
            BranchName: main
            DetectChanges: "false"
    - name: Build
      beforeEntry:
        rules:
          - name: BusinessHoursOnly
            ruleTypeId:
              provider: DeploymentWindow
            configuration:
              Cron: "* 9-17 ? * MON-FRI *"
              TimeZone: America/Los_Angeles
      actions:
        - name: BuildAction
          category: Build
          owner: AWS
          provider: CodeBuild
          version: "1"
          runOrder: 1
          inputArtifacts:
            - SourceOutput
          outputArtifacts:
            - BuildOutput
          configuration:
            ProjectName: my-build-project
      onSuccess:
        rules:
          - name: VerifyNoAlarm
            ruleTypeId:
              provider: CloudWatchAlarm
            configuration:
              AlarmName: build-error-rate
              WaitTime: "300"
            region: us-west-2
    - name: Deploy
      actions:
        # A Manual Approval gate exercises the action-timeout override --
        # the ONLY action type AWS allows it on.
        - name: ApproveDeploy
          category: Approval
          owner: AWS
          provider: Manual
          version: "1"
          runOrder: 1
          timeoutInMinutes: 2880
          configuration:
            CustomData: Approve to release to production
        - name: DeployAction
          category: Deploy
          owner: AWS
          provider: ECS
          version: "1"
          runOrder: 2
          region: us-west-2
          roleArn:
            value: arn:aws:iam::123456789012:role/codepipeline-deploy-role
          inputArtifacts:
            - BuildOutput
          configuration:
            ClusterName: my-cluster
            ServiceName: my-service
            FileName: imagedefinitions.json
      onFailure:
        result: ROLLBACK
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.pipelineType` | `string` |  | `V2` |  |
| `spec.executionMode` | `string` |  | `SUPERSEDED` |  |
| `spec.roleArn` | `string \| valueFrom` | yes |  | AwsIamRole (`status.outputs.role_arn`) |
| `spec.artifactStores` | `[]AwsCodePipelineArtifactStore` | yes |  |  |
| `spec.artifactStores[].location` | `string \| valueFrom` | yes |  | AwsS3Bucket (`status.outputs.bucket_id`) |
| `spec.artifactStores[].region` | `string` |  |  |  |
| `spec.artifactStores[].encryptionKeyId` | `string \| valueFrom` |  |  | AwsKmsKey (`status.outputs.key_arn`) |
| `spec.stages` | `[]AwsCodePipelineStage` | yes |  |  |
| `spec.stages[].name` | `string` | yes |  |  |
| `spec.stages[].actions` | `[]AwsCodePipelineAction` | yes |  |  |
| `spec.stages[].actions[].name` | `string` | yes |  |  |
| `spec.stages[].actions[].category` | `string` | yes |  |  |
| `spec.stages[].actions[].owner` | `string` | yes |  |  |
| `spec.stages[].actions[].provider` | `string` | yes |  |  |
| `spec.stages[].actions[].version` | `string` | yes |  |  |
| `spec.stages[].actions[].configuration` | `map<string, string>` |  |  |  |
| `spec.stages[].actions[].inputArtifacts` | `[]string` |  |  |  |
| `spec.stages[].actions[].outputArtifacts` | `[]string` |  |  |  |
| `spec.stages[].actions[].namespace` | `string` |  |  |  |
| `spec.stages[].actions[].region` | `string` |  |  |  |
| `spec.stages[].actions[].roleArn` | `string \| valueFrom` |  |  | AwsIamRole (`status.outputs.role_arn`) |
| `spec.stages[].actions[].runOrder` | `int32` |  |  |  |
| `spec.stages[].actions[].timeoutInMinutes` | `int32` |  |  |  |
| `spec.stages[].beforeEntry` | `AwsCodePipelineStageCondition` |  |  |  |
| `spec.stages[].beforeEntry.result` | `string` |  |  |  |
| `spec.stages[].beforeEntry.rules` | `[]AwsCodePipelineRule` | yes |  |  |
| `spec.stages[].beforeEntry.rules[].name` | `string` | yes |  |  |
| `spec.stages[].beforeEntry.rules[].ruleTypeId` | `AwsCodePipelineRuleTypeId` | yes |  |  |
| `spec.stages[].beforeEntry.rules[].ruleTypeId.provider` | `string` | yes |  |  |
| `spec.stages[].beforeEntry.rules[].ruleTypeId.category` | `string` |  | `Rule` |  |
| `spec.stages[].beforeEntry.rules[].ruleTypeId.owner` | `string` |  | `AWS` |  |
| `spec.stages[].beforeEntry.rules[].ruleTypeId.version` | `string` |  |  |  |
| `spec.stages[].beforeEntry.rules[].configuration` | `map<string, string>` |  |  |  |
| `spec.stages[].beforeEntry.rules[].commands` | `[]string` |  |  |  |
| `spec.stages[].beforeEntry.rules[].inputArtifacts` | `[]string` |  |  |  |
| `spec.stages[].beforeEntry.rules[].region` | `string` |  |  |  |
| `spec.stages[].beforeEntry.rules[].roleArn` | `string \| valueFrom` |  |  | AwsIamRole (`status.outputs.role_arn`) |
| `spec.stages[].beforeEntry.rules[].timeoutInMinutes` | `int32` |  |  |  |
| `spec.stages[].onSuccess` | `AwsCodePipelineStageCondition` |  |  |  |
| `spec.stages[].onSuccess.result` | `string` |  |  |  |
| `spec.stages[].onSuccess.rules` | `[]AwsCodePipelineRule` | yes |  |  |
| `spec.stages[].onSuccess.rules[].name` | `string` | yes |  |  |
| `spec.stages[].onSuccess.rules[].ruleTypeId` | `AwsCodePipelineRuleTypeId` | yes |  |  |
| `spec.stages[].onSuccess.rules[].ruleTypeId.provider` | `string` | yes |  |  |
| `spec.stages[].onSuccess.rules[].ruleTypeId.category` | `string` |  | `Rule` |  |
| `spec.stages[].onSuccess.rules[].ruleTypeId.owner` | `string` |  | `AWS` |  |
| `spec.stages[].onSuccess.rules[].ruleTypeId.version` | `string` |  |  |  |
| `spec.stages[].onSuccess.rules[].configuration` | `map<string, string>` |  |  |  |
| `spec.stages[].onSuccess.rules[].commands` | `[]string` |  |  |  |
| `spec.stages[].onSuccess.rules[].inputArtifacts` | `[]string` |  |  |  |
| `spec.stages[].onSuccess.rules[].region` | `string` |  |  |  |
| `spec.stages[].onSuccess.rules[].roleArn` | `string \| valueFrom` |  |  | AwsIamRole (`status.outputs.role_arn`) |
| `spec.stages[].onSuccess.rules[].timeoutInMinutes` | `int32` |  |  |  |
| `spec.stages[].onFailure` | `AwsCodePipelineFailureHandling` |  |  |  |
| `spec.stages[].onFailure.result` | `string` |  |  |  |
| `spec.stages[].onFailure.retryConfiguration` | `AwsCodePipelineRetryConfiguration` |  |  |  |
| `spec.stages[].onFailure.retryConfiguration.retryMode` | `string` | yes |  |  |
| `spec.stages[].onFailure.condition` | `AwsCodePipelineStageCondition` |  |  |  |
| `spec.stages[].onFailure.condition.result` | `string` |  |  |  |
| `spec.stages[].onFailure.condition.rules` | `[]AwsCodePipelineRule` | yes |  |  |
| `spec.stages[].onFailure.condition.rules[].name` | `string` | yes |  |  |
| `spec.stages[].onFailure.condition.rules[].ruleTypeId` | `AwsCodePipelineRuleTypeId` | yes |  |  |
| `spec.stages[].onFailure.condition.rules[].ruleTypeId.provider` | `string` | yes |  |  |
| `spec.stages[].onFailure.condition.rules[].ruleTypeId.category` | `string` |  | `Rule` |  |
| `spec.stages[].onFailure.condition.rules[].ruleTypeId.owner` | `string` |  | `AWS` |  |
| `spec.stages[].onFailure.condition.rules[].ruleTypeId.version` | `string` |  |  |  |
| `spec.stages[].onFailure.condition.rules[].configuration` | `map<string, string>` |  |  |  |
| `spec.stages[].onFailure.condition.rules[].commands` | `[]string` |  |  |  |
| `spec.stages[].onFailure.condition.rules[].inputArtifacts` | `[]string` |  |  |  |
| `spec.stages[].onFailure.condition.rules[].region` | `string` |  |  |  |
| `spec.stages[].onFailure.condition.rules[].roleArn` | `string \| valueFrom` |  |  | AwsIamRole (`status.outputs.role_arn`) |
| `spec.stages[].onFailure.condition.rules[].timeoutInMinutes` | `int32` |  |  |  |
| `spec.triggers` | `[]AwsCodePipelineTrigger` |  |  |  |
| `spec.triggers[].providerType` | `string` | yes |  |  |
| `spec.triggers[].gitConfiguration` | `AwsCodePipelineGitConfiguration` | yes |  |  |
| `spec.triggers[].gitConfiguration.sourceActionName` | `string` | yes |  |  |
| `spec.triggers[].gitConfiguration.push` | `[]AwsCodePipelineGitPush` |  |  |  |
| `spec.triggers[].gitConfiguration.push[].branches` | `AwsCodePipelineGitFilter` |  |  |  |
| `spec.triggers[].gitConfiguration.push[].branches.includes` | `[]string` |  |  |  |
| `spec.triggers[].gitConfiguration.push[].branches.excludes` | `[]string` |  |  |  |
| `spec.triggers[].gitConfiguration.push[].filePaths` | `AwsCodePipelineGitFilter` |  |  |  |
| `spec.triggers[].gitConfiguration.push[].filePaths.includes` | `[]string` |  |  |  |
| `spec.triggers[].gitConfiguration.push[].filePaths.excludes` | `[]string` |  |  |  |
| `spec.triggers[].gitConfiguration.push[].tags` | `AwsCodePipelineGitFilter` |  |  |  |
| `spec.triggers[].gitConfiguration.push[].tags.includes` | `[]string` |  |  |  |
| `spec.triggers[].gitConfiguration.push[].tags.excludes` | `[]string` |  |  |  |
| `spec.triggers[].gitConfiguration.pullRequest` | `[]AwsCodePipelineGitPullRequest` |  |  |  |
| `spec.triggers[].gitConfiguration.pullRequest[].branches` | `AwsCodePipelineGitFilter` |  |  |  |
| `spec.triggers[].gitConfiguration.pullRequest[].branches.includes` | `[]string` |  |  |  |
| `spec.triggers[].gitConfiguration.pullRequest[].branches.excludes` | `[]string` |  |  |  |
| `spec.triggers[].gitConfiguration.pullRequest[].filePaths` | `AwsCodePipelineGitFilter` |  |  |  |
| `spec.triggers[].gitConfiguration.pullRequest[].filePaths.includes` | `[]string` |  |  |  |
| `spec.triggers[].gitConfiguration.pullRequest[].filePaths.excludes` | `[]string` |  |  |  |
| `spec.triggers[].gitConfiguration.pullRequest[].events` | `[]string` |  |  |  |
| `spec.variables` | `[]AwsCodePipelineVariable` |  |  |  |
| `spec.variables[].name` | `string` | yes |  |  |
| `spec.variables[].defaultValue` | `string` |  |  |  |
| `spec.variables[].description` | `string` |  |  |  |

## Field Details

### spec.region

`string` · required

The AWS region where the resource will be created.
Example: "us-west-2", "eu-west-1"

- rule: {"string":{"minLen":"1"}}

### spec.pipelineType

`string` · optional (explicit presence)

pipeline_type selects the pipeline version.
  V1: Legacy pipeline (no triggers, no variables, SUPERSEDED execution only)
  V2: Modern pipeline with triggers, variables, and advanced execution modes
Default: V2 (recommended for all new pipelines).

- default: `V2`
- rule: {"string":{"in":["V1","V2"]}}

### spec.executionMode

`string` · optional (explicit presence)

execution_mode controls how concurrent pipeline executions are handled.
  SUPERSEDED: New execution supersedes any in-progress execution (default)
  QUEUED:     New executions queue behind the current one (V2 only)
  PARALLEL:   Executions run simultaneously without waiting (V2 only)
Default: SUPERSEDED.

- default: `SUPERSEDED`
- rule: {"string":{"in":["SUPERSEDED","QUEUED","PARALLEL"]}}

### spec.roleArn

`string | valueFrom` · required

role_arn is the IAM role ARN that grants CodePipeline permission to
access source providers, invoke build/deploy actions, and manage
artifacts in S3. This role must have policies for every action provider
used in the pipeline (e.g., CodeBuild, S3, ECS, Lambda).

- references: AwsIamRole (`status.outputs.role_arn`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsIamRole, name: <that resource's name>, fieldPath: status.outputs.role_arn}} -- a bare string does not parse

### spec.artifactStores

`[]AwsCodePipelineArtifactStore` · required

artifact_stores define S3 buckets where pipeline artifacts are stored.
For single-region pipelines, provide exactly one store without a region.
For cross-region pipelines, provide one store per region (each with a
region field) so that actions in different regions have local artifact
access.

- rule: {"required":true,"repeated":{"minItems":"1"}}

### spec.artifactStores[].location

`string | valueFrom` · required

location is the S3 bucket name for artifact storage.

- references: AwsS3Bucket (`status.outputs.bucket_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsS3Bucket, name: <that resource's name>, fieldPath: status.outputs.bucket_id}} -- a bare string does not parse

### spec.artifactStores[].region

`string`

region is the AWS region for this artifact store. Required only for
cross-region pipelines. For single-region pipelines, omit this field.

### spec.artifactStores[].encryptionKeyId

`string | valueFrom`

encryption_key_id is the KMS key ARN or ID used to encrypt artifacts.
If omitted, the default AWS-managed S3 encryption key is used.

- references: AwsKmsKey (`status.outputs.key_arn`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsKmsKey, name: <that resource's name>, fieldPath: status.outputs.key_arn}} -- a bare string does not parse

### spec.stages

`[]AwsCodePipelineStage` · required

stages define the ordered sequence of pipeline stages. Each stage
contains one or more actions that run in parallel (same run_order)
or sequentially (different run_order values).
A pipeline requires at minimum two stages: a source stage and at
least one build, test, deploy, or approval stage.

- rule: {"required":true,"repeated":{"minItems":"2"}}

### spec.stages[].name

`string` · required

name is the stage name. Must be unique within the pipeline.
Pattern: alphanumeric, dots, at-signs, hyphens, underscores (1-100 chars).

- rule: {"required":true,"string":{"minLen":"1","maxLen":"100","pattern":"^[0-9A-Za-z_.@\\-]+$"}}

### spec.stages[].actions

`[]AwsCodePipelineAction` · required

actions define the operations performed in this stage. Actions with
the same run_order execute in parallel; different run_order values
execute sequentially within the stage.

- rule: {"required":true,"repeated":{"minItems":"1"}}
- rule: timeout_in_minutes is only supported on Manual Approval actions (category Approval)

### spec.stages[].actions[].name

`string` · required

name is the action name. Must be unique within the stage.

- rule: {"required":true,"string":{"minLen":"1","maxLen":"100","pattern":"^[0-9A-Za-z_.@\\-]+$"}}

### spec.stages[].actions[].category

`string` · required

category classifies the action type.
  Source:   Fetches source code or artifacts from a provider
  Build:    Compiles, tests, or transforms code
  Test:     Runs test suites
  Deploy:   Deploys artifacts to a target environment
  Approval: Requires manual approval before proceeding
  Invoke:   Invokes a Lambda function or other compute
  Compute:  Runs commands in a managed compute environment (V2)

- rule: {"required":true,"string":{"in":["Source","Build","Test","Deploy","Approval","Invoke","Compute"]}}

### spec.stages[].actions[].owner

`string` · required

owner identifies who created the action type.
  AWS:        Built-in AWS actions (CodeBuild, S3, ECS, Lambda, etc.)
  ThirdParty: Third-party integrations (GitHub v1, etc.)
  Custom:     User-defined custom action types

- rule: {"required":true,"string":{"in":["AWS","ThirdParty","Custom"]}}

### spec.stages[].actions[].provider

`string` · required

provider is the service that performs the action. The valid values
depend on the category and owner combination. Common providers:
  Source:   CodeStarSourceConnection, S3, ECR, CodeCommit
  Build:    CodeBuild
  Test:     CodeBuild, DeviceFarm
  Deploy:   S3, CodeDeploy, CloudFormation, ECS, ElasticBeanstalk, Lambda
  Approval: Manual
  Invoke:   Lambda
  Compute:  EC2 (V2)

- rule: {"required":true,"string":{"minLen":"1","maxLen":"35"}}

### spec.stages[].actions[].version

`string` · required

version is the action type version. Typically "1" for all built-in actions.

- rule: {"required":true,"string":{"minLen":"1","maxLen":"9"}}

### spec.stages[].actions[].configuration

`map<string, string>`

configuration contains provider-specific key-value pairs that control
the action's behavior. Each provider expects different keys.

Common examples:
  CodeStarSourceConnection: ConnectionArn, FullRepositoryId, BranchName,
                            OutputArtifactFormat, DetectChanges
  CodeBuild:                ProjectName, PrimarySource, EnvironmentVariables
  S3 (source):              S3Bucket, S3ObjectKey, PollForSourceChanges
  S3 (deploy):              BucketName, Extract, ObjectKey, CannedACL
  ECS:                      ClusterName, ServiceName, FileName
  Lambda:                   FunctionName, UserParameters
  Manual (approval):        CustomData, ExternalEntityLink, NotificationArn
  CloudFormation:           ActionMode, StackName, TemplatePath, RoleArn

### spec.stages[].actions[].inputArtifacts

`[]string`

input_artifacts are artifact names from previous stages or actions
that this action consumes. For Source actions, this is typically empty.

### spec.stages[].actions[].outputArtifacts

`[]string`

output_artifacts are artifact names that this action produces for
consumption by downstream stages or actions.

### spec.stages[].actions[].namespace

`string`

namespace defines a variable namespace for this action's output variables.
Other actions can reference these variables as #{namespace.VariableName}.
Only meaningful for actions that produce output variables (e.g., source actions).

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"maxLen":"100","pattern":"^[0-9A-Za-z_@\\-]+$"}}

### spec.stages[].actions[].region

`string`

region is the AWS region where this action executes. Required for
cross-region actions. If omitted, the action runs in the pipeline's region.

### spec.stages[].actions[].roleArn

`string | valueFrom`

role_arn is an IAM role ARN that the action assumes instead of the
pipeline's role. Useful for cross-account deployments or when an action
needs different permissions than the pipeline role.

- references: AwsIamRole (`status.outputs.role_arn`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsIamRole, name: <that resource's name>, fieldPath: status.outputs.role_arn}} -- a bare string does not parse

### spec.stages[].actions[].runOrder

`int32`

run_order controls execution order within a stage. Actions with the
same run_order execute in parallel. Lower values run first.
Range: 1-999. Default: 1 (all actions parallel).

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","int32":{"lte":999,"gte":1}}

### spec.stages[].actions[].timeoutInMinutes

`int32`

timeout_in_minutes overrides the action type's default timeout.
Range: 5-86400 (60 days). AWS supports this override ONLY on Manual
Approval actions (category Approval, provider Manual) — every other
action type is rejected at pipeline creation.

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","int32":{"lte":86400,"gte":5}}

### spec.stages[].beforeEntry

`AwsCodePipelineStageCondition`

before_entry is an entry gate: its rules must pass before the stage
starts (e.g., a DeploymentWindow rule that only admits executions
during business hours, or a CloudWatchAlarm rule that blocks deploys
while an alarm is firing). V2 pipelines only.

### spec.stages[].beforeEntry.result

`string`

result is the outcome applied when the condition's rules do NOT pass.
  FAIL: Fail the stage (default behavior when omitted)
  SKIP: Skip the stage and continue the pipeline (before_entry gates)

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"in":["FAIL","SKIP"]}}

### spec.stages[].beforeEntry.rules

`[]AwsCodePipelineRule` · required

rules are the checks that make up this condition (1-5, AND'd together).

- rule: {"required":true,"repeated":{"minItems":"1","maxItems":"5"}}

### spec.stages[].beforeEntry.rules[].name

`string` · required

name is the rule name, unique within the condition.

- rule: {"required":true,"string":{"minLen":"1","maxLen":"100","pattern":"^[A-Za-z0-9.@\\-_]+$"}}

### spec.stages[].beforeEntry.rules[].ruleTypeId

`AwsCodePipelineRuleTypeId` · required

rule_type_id identifies which managed rule runs.

- rule: {"required":true}

### spec.stages[].beforeEntry.rules[].ruleTypeId.provider

`string` · required

provider is the managed rule provider: DeploymentWindow,
CloudWatchAlarm, LambdaInvoke, VariableCheck, or Commands.

- rule: {"required":true,"string":{"minLen":"1","maxLen":"35"}}

### spec.stages[].beforeEntry.rules[].ruleTypeId.category

`string` · optional (explicit presence)

category is the rule category. AWS currently supports only "Rule".

- default: `Rule`
- rule: {"string":{"const":"Rule"}}

### spec.stages[].beforeEntry.rules[].ruleTypeId.owner

`string` · optional (explicit presence)

owner identifies who provides the rule. AWS currently supports only
"AWS" (managed rules).

- default: `AWS`
- rule: {"string":{"const":"AWS"}}

### spec.stages[].beforeEntry.rules[].ruleTypeId.version

`string`

version is the rule type version. Typically "1"; AWS resolves the
current version when omitted.

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"maxLen":"9","pattern":"^[0-9A-Za-z_\\-]+$"}}

### spec.stages[].beforeEntry.rules[].configuration

`map<string, string>`

configuration contains rule-provider-specific key-value pairs (e.g.,
DeploymentWindow: Cron, TimeZone; CloudWatchAlarm: AlarmName,
WaitTime; LambdaInvoke: FunctionName).

### spec.stages[].beforeEntry.rules[].commands

`[]string`

commands are shell commands executed by the Commands rule provider
(max 50, each 1-1000 characters).

- rule: {"repeated":{"maxItems":"50","items":{"string":{"minLen":"1","maxLen":"1000"}}}}

### spec.stages[].beforeEntry.rules[].inputArtifacts

`[]string`

input_artifacts are artifact names available to the rule (e.g., for
Commands rules operating on build output).

### spec.stages[].beforeEntry.rules[].region

`string`

region is the AWS region where the rule executes. Omit to run in the
pipeline's region.

### spec.stages[].beforeEntry.rules[].roleArn

`string | valueFrom`

role_arn is an IAM role the rule assumes instead of the pipeline's
role (e.g., a scoped role for the LambdaInvoke rule).

- references: AwsIamRole (`status.outputs.role_arn`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsIamRole, name: <that resource's name>, fieldPath: status.outputs.role_arn}} -- a bare string does not parse

### spec.stages[].beforeEntry.rules[].timeoutInMinutes

`int32`

timeout_in_minutes is the maximum time the rule can run (5-86400).

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","int32":{"lte":86400,"gte":5}}

### spec.stages[].onSuccess

`AwsCodePipelineStageCondition`

on_success runs verification rules after the stage's actions succeed;
a failing rule fails the stage despite the successful actions (e.g., a
post-deploy CloudWatchAlarm check). V2 pipelines only.

### spec.stages[].onSuccess.result

`string`

result is the outcome applied when the condition's rules do NOT pass.
  FAIL: Fail the stage (default behavior when omitted)
  SKIP: Skip the stage and continue the pipeline (before_entry gates)

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"in":["FAIL","SKIP"]}}

### spec.stages[].onSuccess.rules

`[]AwsCodePipelineRule` · required

rules are the checks that make up this condition (1-5, AND'd together).

- rule: {"required":true,"repeated":{"minItems":"1","maxItems":"5"}}

### spec.stages[].onSuccess.rules[].name

`string` · required

name is the rule name, unique within the condition.

- rule: {"required":true,"string":{"minLen":"1","maxLen":"100","pattern":"^[A-Za-z0-9.@\\-_]+$"}}

### spec.stages[].onSuccess.rules[].ruleTypeId

`AwsCodePipelineRuleTypeId` · required

rule_type_id identifies which managed rule runs.

- rule: {"required":true}

### spec.stages[].onSuccess.rules[].ruleTypeId.provider

`string` · required

provider is the managed rule provider: DeploymentWindow,
CloudWatchAlarm, LambdaInvoke, VariableCheck, or Commands.

- rule: {"required":true,"string":{"minLen":"1","maxLen":"35"}}

### spec.stages[].onSuccess.rules[].ruleTypeId.category

`string` · optional (explicit presence)

category is the rule category. AWS currently supports only "Rule".

- default: `Rule`
- rule: {"string":{"const":"Rule"}}

### spec.stages[].onSuccess.rules[].ruleTypeId.owner

`string` · optional (explicit presence)

owner identifies who provides the rule. AWS currently supports only
"AWS" (managed rules).

- default: `AWS`
- rule: {"string":{"const":"AWS"}}

### spec.stages[].onSuccess.rules[].ruleTypeId.version

`string`

version is the rule type version. Typically "1"; AWS resolves the
current version when omitted.

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"maxLen":"9","pattern":"^[0-9A-Za-z_\\-]+$"}}

### spec.stages[].onSuccess.rules[].configuration

`map<string, string>`

configuration contains rule-provider-specific key-value pairs (e.g.,
DeploymentWindow: Cron, TimeZone; CloudWatchAlarm: AlarmName,
WaitTime; LambdaInvoke: FunctionName).

### spec.stages[].onSuccess.rules[].commands

`[]string`

commands are shell commands executed by the Commands rule provider
(max 50, each 1-1000 characters).

- rule: {"repeated":{"maxItems":"50","items":{"string":{"minLen":"1","maxLen":"1000"}}}}

### spec.stages[].onSuccess.rules[].inputArtifacts

`[]string`

input_artifacts are artifact names available to the rule (e.g., for
Commands rules operating on build output).

### spec.stages[].onSuccess.rules[].region

`string`

region is the AWS region where the rule executes. Omit to run in the
pipeline's region.

### spec.stages[].onSuccess.rules[].roleArn

`string | valueFrom`

role_arn is an IAM role the rule assumes instead of the pipeline's
role (e.g., a scoped role for the LambdaInvoke rule).

- references: AwsIamRole (`status.outputs.role_arn`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsIamRole, name: <that resource's name>, fieldPath: status.outputs.role_arn}} -- a bare string does not parse

### spec.stages[].onSuccess.rules[].timeoutInMinutes

`int32`

timeout_in_minutes is the maximum time the rule can run (5-86400).

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","int32":{"lte":86400,"gte":5}}

### spec.stages[].onFailure

`AwsCodePipelineFailureHandling`

on_failure controls what happens when the stage fails: automatically
roll back to the last successful state, retry the stage, or
conditionally decide via rules. V2 pipelines only.

- rule: retry_configuration requires result to be RETRY

### spec.stages[].onFailure.result

`string`

result is the action taken when the stage fails.
  ROLLBACK: Automatically roll the stage back to its last successful
            execution's state
  RETRY:    Automatically retry per retry_configuration
  FAIL:     Fail the pipeline execution (default behavior when omitted)

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"in":["ROLLBACK","RETRY","FAIL"]}}

### spec.stages[].onFailure.retryConfiguration

`AwsCodePipelineRetryConfiguration`

retry_configuration tunes automatic retry when result is RETRY.

### spec.stages[].onFailure.retryConfiguration.retryMode

`string` · required

retry_mode selects what is retried.
  FAILED_ACTIONS: Only the failed actions re-run (default)
  ALL_ACTIONS:    The whole stage re-runs from its first action

- rule: {"required":true,"string":{"in":["FAILED_ACTIONS","ALL_ACTIONS"]}}

### spec.stages[].onFailure.condition

`AwsCodePipelineStageCondition`

condition optionally gates the failure handling behind rules — the
result applies only when the rules pass.

### spec.stages[].onFailure.condition.result

`string`

result is the outcome applied when the condition's rules do NOT pass.
  FAIL: Fail the stage (default behavior when omitted)
  SKIP: Skip the stage and continue the pipeline (before_entry gates)

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"in":["FAIL","SKIP"]}}

### spec.stages[].onFailure.condition.rules

`[]AwsCodePipelineRule` · required

rules are the checks that make up this condition (1-5, AND'd together).

- rule: {"required":true,"repeated":{"minItems":"1","maxItems":"5"}}

### spec.stages[].onFailure.condition.rules[].name

`string` · required

name is the rule name, unique within the condition.

- rule: {"required":true,"string":{"minLen":"1","maxLen":"100","pattern":"^[A-Za-z0-9.@\\-_]+$"}}

### spec.stages[].onFailure.condition.rules[].ruleTypeId

`AwsCodePipelineRuleTypeId` · required

rule_type_id identifies which managed rule runs.

- rule: {"required":true}

### spec.stages[].onFailure.condition.rules[].ruleTypeId.provider

`string` · required

provider is the managed rule provider: DeploymentWindow,
CloudWatchAlarm, LambdaInvoke, VariableCheck, or Commands.

- rule: {"required":true,"string":{"minLen":"1","maxLen":"35"}}

### spec.stages[].onFailure.condition.rules[].ruleTypeId.category

`string` · optional (explicit presence)

category is the rule category. AWS currently supports only "Rule".

- default: `Rule`
- rule: {"string":{"const":"Rule"}}

### spec.stages[].onFailure.condition.rules[].ruleTypeId.owner

`string` · optional (explicit presence)

owner identifies who provides the rule. AWS currently supports only
"AWS" (managed rules).

- default: `AWS`
- rule: {"string":{"const":"AWS"}}

### spec.stages[].onFailure.condition.rules[].ruleTypeId.version

`string`

version is the rule type version. Typically "1"; AWS resolves the
current version when omitted.

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"maxLen":"9","pattern":"^[0-9A-Za-z_\\-]+$"}}

### spec.stages[].onFailure.condition.rules[].configuration

`map<string, string>`

configuration contains rule-provider-specific key-value pairs (e.g.,
DeploymentWindow: Cron, TimeZone; CloudWatchAlarm: AlarmName,
WaitTime; LambdaInvoke: FunctionName).

### spec.stages[].onFailure.condition.rules[].commands

`[]string`

commands are shell commands executed by the Commands rule provider
(max 50, each 1-1000 characters).

- rule: {"repeated":{"maxItems":"50","items":{"string":{"minLen":"1","maxLen":"1000"}}}}

### spec.stages[].onFailure.condition.rules[].inputArtifacts

`[]string`

input_artifacts are artifact names available to the rule (e.g., for
Commands rules operating on build output).

### spec.stages[].onFailure.condition.rules[].region

`string`

region is the AWS region where the rule executes. Omit to run in the
pipeline's region.

### spec.stages[].onFailure.condition.rules[].roleArn

`string | valueFrom`

role_arn is an IAM role the rule assumes instead of the pipeline's
role (e.g., a scoped role for the LambdaInvoke rule).

- references: AwsIamRole (`status.outputs.role_arn`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsIamRole, name: <that resource's name>, fieldPath: status.outputs.role_arn}} -- a bare string does not parse

### spec.stages[].onFailure.condition.rules[].timeoutInMinutes

`int32`

timeout_in_minutes is the maximum time the rule can run (5-86400).

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","int32":{"lte":86400,"gte":5}}

### spec.triggers

`[]AwsCodePipelineTrigger`

triggers define automatic pipeline execution rules based on git events.
V2 pipelines only. Triggers use CodeStar Connections to listen for
push or pull request events on source repositories with branch, tag,
and file path filtering.

- rule: {"repeated":{"maxItems":"50"}}

### spec.triggers[].providerType

`string` · required

provider_type is the trigger provider. Currently only
CodeStarSourceConnection is supported.

- rule: {"required":true,"string":{"const":"CodeStarSourceConnection"}}

### spec.triggers[].gitConfiguration

`AwsCodePipelineGitConfiguration` · required

git_configuration defines the git event filters that trigger the pipeline.

- rule: {"required":true}
- rule: git_configuration requires at least one push or pull_request filter

### spec.triggers[].gitConfiguration.sourceActionName

`string` · required

source_action_name must match the name of a Source action in the first
stage that uses a CodeStarSourceConnection provider.

- rule: {"required":true,"string":{"minLen":"1","maxLen":"100"}}

### spec.triggers[].gitConfiguration.push

`[]AwsCodePipelineGitPush`

push defines filters for git push events (branch pushes, tag pushes).
Multiple push filters are OR'd — a push triggers the pipeline if ANY
filter matches.

- rule: {"repeated":{"maxItems":"3"}}

### spec.triggers[].gitConfiguration.push[].branches

`AwsCodePipelineGitFilter`

branches filters by branch name patterns (glob syntax).

- rule: a git filter requires at least one includes or excludes pattern

### spec.triggers[].gitConfiguration.push[].branches.includes

`[]string`

includes are glob patterns that must match for the filter to pass.
At least one include pattern must match. Each pattern 1-255 characters.

- rule: {"repeated":{"maxItems":"8","items":{"string":{"minLen":"1","maxLen":"255"}}}}

### spec.triggers[].gitConfiguration.push[].branches.excludes

`[]string`

excludes are glob patterns that cause the filter to reject a match.
Exclusions take precedence over inclusions. Each pattern 1-255 characters.

- rule: {"repeated":{"maxItems":"8","items":{"string":{"minLen":"1","maxLen":"255"}}}}

### spec.triggers[].gitConfiguration.push[].filePaths

`AwsCodePipelineGitFilter`

file_paths filters by changed file path patterns (glob syntax).

- rule: a git filter requires at least one includes or excludes pattern

### spec.triggers[].gitConfiguration.push[].filePaths.includes

`[]string`

includes are glob patterns that must match for the filter to pass.
At least one include pattern must match. Each pattern 1-255 characters.

- rule: {"repeated":{"maxItems":"8","items":{"string":{"minLen":"1","maxLen":"255"}}}}

### spec.triggers[].gitConfiguration.push[].filePaths.excludes

`[]string`

excludes are glob patterns that cause the filter to reject a match.
Exclusions take precedence over inclusions. Each pattern 1-255 characters.

- rule: {"repeated":{"maxItems":"8","items":{"string":{"minLen":"1","maxLen":"255"}}}}

### spec.triggers[].gitConfiguration.push[].tags

`AwsCodePipelineGitFilter`

tags filters by tag name patterns (glob syntax).

- rule: a git filter requires at least one includes or excludes pattern

### spec.triggers[].gitConfiguration.push[].tags.includes

`[]string`

includes are glob patterns that must match for the filter to pass.
At least one include pattern must match. Each pattern 1-255 characters.

- rule: {"repeated":{"maxItems":"8","items":{"string":{"minLen":"1","maxLen":"255"}}}}

### spec.triggers[].gitConfiguration.push[].tags.excludes

`[]string`

excludes are glob patterns that cause the filter to reject a match.
Exclusions take precedence over inclusions. Each pattern 1-255 characters.

- rule: {"repeated":{"maxItems":"8","items":{"string":{"minLen":"1","maxLen":"255"}}}}

### spec.triggers[].gitConfiguration.pullRequest

`[]AwsCodePipelineGitPullRequest`

pull_request defines filters for pull request events (open, update, close).
Multiple PR filters are OR'd.

- rule: {"repeated":{"maxItems":"3"}}

### spec.triggers[].gitConfiguration.pullRequest[].branches

`AwsCodePipelineGitFilter`

branches filters by target branch name patterns (glob syntax).

- rule: a git filter requires at least one includes or excludes pattern

### spec.triggers[].gitConfiguration.pullRequest[].branches.includes

`[]string`

includes are glob patterns that must match for the filter to pass.
At least one include pattern must match. Each pattern 1-255 characters.

- rule: {"repeated":{"maxItems":"8","items":{"string":{"minLen":"1","maxLen":"255"}}}}

### spec.triggers[].gitConfiguration.pullRequest[].branches.excludes

`[]string`

excludes are glob patterns that cause the filter to reject a match.
Exclusions take precedence over inclusions. Each pattern 1-255 characters.

- rule: {"repeated":{"maxItems":"8","items":{"string":{"minLen":"1","maxLen":"255"}}}}

### spec.triggers[].gitConfiguration.pullRequest[].filePaths

`AwsCodePipelineGitFilter`

file_paths filters by changed file path patterns (glob syntax).

- rule: a git filter requires at least one includes or excludes pattern

### spec.triggers[].gitConfiguration.pullRequest[].filePaths.includes

`[]string`

includes are glob patterns that must match for the filter to pass.
At least one include pattern must match. Each pattern 1-255 characters.

- rule: {"repeated":{"maxItems":"8","items":{"string":{"minLen":"1","maxLen":"255"}}}}

### spec.triggers[].gitConfiguration.pullRequest[].filePaths.excludes

`[]string`

excludes are glob patterns that cause the filter to reject a match.
Exclusions take precedence over inclusions. Each pattern 1-255 characters.

- rule: {"repeated":{"maxItems":"8","items":{"string":{"minLen":"1","maxLen":"255"}}}}

### spec.triggers[].gitConfiguration.pullRequest[].events

`[]string`

events specifies which PR lifecycle events trigger the pipeline
(1-3 values). Valid values: OPEN, UPDATED, CLOSED.

- rule: {"repeated":{"maxItems":"3","items":{"string":{"in":["OPEN","UPDATED","CLOSED"]}}}}

### spec.variables

`[]AwsCodePipelineVariable`

variables define pipeline-level parameters that can be referenced in
action configurations using #{variables.VariableName} syntax.
V2 pipelines only.

### spec.variables[].name

`string` · required

name is the variable name. Referenced in action configurations as
#{variables.Name}.

- rule: {"required":true}

### spec.variables[].defaultValue

`string`

default_value is used when no value is supplied at execution time.

### spec.variables[].description

`string`

description is a human-readable explanation of the variable's purpose.

## Validation Rules

- `triggers_require_v2`: triggers are only supported on V2 pipelines; set pipeline_type to V2 or remove triggers
- `variables_require_v2`: variables are only supported on V2 pipelines; set pipeline_type to V2 or remove variables
- `advanced_execution_mode_requires_v2`: execution_mode QUEUED and PARALLEL are only supported on V2 pipelines
- `stage_conditions_require_v2`: stage before_entry, on_success, and on_failure conditions are only supported on V2 pipelines
- `artifact_stores_region_shape`: a single artifact store must omit region; multiple (cross-region) artifact stores must each set region

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AwsCodePipeline, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.pipelineArn` | `string` | The Amazon Resource Name (ARN) of the pipeline. Use this for IAM policies, EventBridge targets, and cross-resource references. |
| `status.outputs.pipelineName` | `string` | The name of the pipeline. Use this when referencing the pipeline in action configurations of other pipelines or in CLI commands. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.roleArn` | AwsIamRole | `status.outputs.role_arn` |
| `spec.artifactStores[].location` | AwsS3Bucket | `status.outputs.bucket_id` |
| `spec.artifactStores[].encryptionKeyId` | AwsKmsKey | `status.outputs.key_arn` |
| `spec.stages[].actions[].roleArn` | AwsIamRole | `status.outputs.role_arn` |
| `spec.stages[].beforeEntry.rules[].roleArn` | AwsIamRole | `status.outputs.role_arn` |
| `spec.stages[].onSuccess.rules[].roleArn` | AwsIamRole | `status.outputs.role_arn` |
| `spec.stages[].onFailure.condition.rules[].roleArn` | AwsIamRole | `status.outputs.role_arn` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
