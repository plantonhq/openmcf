# AwsBedrockAgentCoreEvaluation

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `aws.planton.dev/v1alpha1`

**Guide**: [GUIDE.md](../GUIDE.md) -- authored operational judgment for this component: conventions, trade-offs, and what pairs well with it.

## Example

```yaml
# Canonical AwsBedrockAgentCoreEvaluation example (hack/dev manifest and
# refgen Example source): a bundle exercising every arm -- an LLM-judge
# evaluator with a numerical scale, a code-based evaluator, a Bedrock
# harness, and an online evaluation config sampling CloudWatch logs
# with a builtin evaluator. Literal ARNs stand in for composed
# references so the offline `tofu plan` renders every arm.
apiVersion: aws.planton.dev/v1alpha1
kind: AwsBedrockAgentCoreEvaluation
metadata:
  name: support-eval
  id: support-eval
  org: test-org
  env: dev
spec:
  region: us-west-2
  evaluators:
    - name: helpfulness
      description: Scores whether the agent was helpful
      level: SESSION
      llmAsAJudge:
        instructions: Score how helpful the agent's response was.
        model:
          modelId: anthropic.claude-3-haiku-20240307-v1:0
          inference:
            temperature: 0
            maxTokens: 256
        ratingScale:
          numerical:
            - label: poor
              definition: Unhelpful or incorrect
              value: 1
            - label: excellent
              definition: Fully helpful and correct
              value: 5
    - name: custom_score
      description: Lambda scorer
      level: TRACE
      codeBased:
        lambdaArn:
          value: arn:aws:lambda:us-west-2:123456789012:function:score-trace
        timeoutSeconds: 60
  harnesses:
    - name: support-bench
      executionRoleArn:
        value: arn:aws:iam::123456789012:role/agentcore-eval-role
      model:
        bedrock:
          modelId: anthropic.claude-3-haiku-20240307-v1:0
          temperature: 0.2
          maxTokens: 512
      systemPrompts:
        - text: You are a support agent under evaluation.
      tools:
        - name: search
          type: inline_function
          inlineFunction:
            description: Search the knowledge base
            inputSchema: '{"type":"object","properties":{"query":{"type":"string"}}}'
  onlineEvaluationConfigs:
    - name: prod_sampling
      description: Sample 5 percent of production sessions
      executionRoleArn:
        value: arn:aws:iam::123456789012:role/agentcore-eval-role
      enabled: true
      dataSource:
        logGroupNames:
          - value: /aws/bedrock/agentcore/support
        serviceNames:
          - support-agent
      evaluatorIds:
        - Builtin.Helpfulness
        - helpfulness
      rule:
        samplingPercentage: 5
        sessionTimeoutMinutes: 15
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.evaluators` | `[]AwsBedrockAgentCoreEvaluator` |  |  |  |
| `spec.evaluators[].name` | `string` | yes |  |  |
| `spec.evaluators[].description` | `string` |  |  |  |
| `spec.evaluators[].level` | `string` |  |  |  |
| `spec.evaluators[].kmsKeyArn` | `string \| valueFrom` |  |  | AwsKmsKey (`status.outputs.key_arn`) |
| `spec.evaluators[].llmAsAJudge` | `AwsBedrockAgentCoreLlmJudge` |  |  |  |
| `spec.evaluators[].llmAsAJudge.instructions` | `string` (sensitive) | yes |  |  |
| `spec.evaluators[].llmAsAJudge.model` | `AwsBedrockAgentCoreJudgeModel` | yes |  |  |
| `spec.evaluators[].llmAsAJudge.model.modelId` | `string` | yes |  |  |
| `spec.evaluators[].llmAsAJudge.model.additionalModelRequestFields` | `object` |  |  |  |
| `spec.evaluators[].llmAsAJudge.model.inference` | `AwsBedrockAgentCoreJudgeInference` |  |  |  |
| `spec.evaluators[].llmAsAJudge.model.inference.maxTokens` | `int32` |  |  |  |
| `spec.evaluators[].llmAsAJudge.model.inference.stopSequences` | `[]string` |  |  |  |
| `spec.evaluators[].llmAsAJudge.model.inference.temperature` | `double` |  |  |  |
| `spec.evaluators[].llmAsAJudge.model.inference.topP` | `double` |  |  |  |
| `spec.evaluators[].llmAsAJudge.ratingScale` | `AwsBedrockAgentCoreRatingScale` | yes |  |  |
| `spec.evaluators[].llmAsAJudge.ratingScale.categorical` | `[]AwsBedrockAgentCoreCategoricalRating` |  |  |  |
| `spec.evaluators[].llmAsAJudge.ratingScale.categorical[].label` | `string` | yes |  |  |
| `spec.evaluators[].llmAsAJudge.ratingScale.categorical[].definition` | `string` | yes |  |  |
| `spec.evaluators[].llmAsAJudge.ratingScale.numerical` | `[]AwsBedrockAgentCoreNumericalRating` |  |  |  |
| `spec.evaluators[].llmAsAJudge.ratingScale.numerical[].label` | `string` | yes |  |  |
| `spec.evaluators[].llmAsAJudge.ratingScale.numerical[].definition` | `string` | yes |  |  |
| `spec.evaluators[].llmAsAJudge.ratingScale.numerical[].value` | `double` |  |  |  |
| `spec.evaluators[].codeBased` | `AwsBedrockAgentCoreCodeEvaluator` |  |  |  |
| `spec.evaluators[].codeBased.lambdaArn` | `string \| valueFrom` | yes |  | AwsLambda (`status.outputs.function_arn`) |
| `spec.evaluators[].codeBased.timeoutSeconds` | `int32` |  |  |  |
| `spec.harnesses` | `[]AwsBedrockAgentCoreHarness` |  |  |  |
| `spec.harnesses[].name` | `string` | yes |  |  |
| `spec.harnesses[].executionRoleArn` | `string \| valueFrom` | yes |  | AwsIamRole (`status.outputs.role_arn`) |
| `spec.harnesses[].model` | `AwsBedrockAgentCoreHarnessModel` | yes |  |  |
| `spec.harnesses[].model.bedrock` | `AwsBedrockAgentCoreHarnessBedrockModel` |  |  |  |
| `spec.harnesses[].model.bedrock.modelId` | `string` | yes |  |  |
| `spec.harnesses[].model.bedrock.maxTokens` | `int32` |  |  |  |
| `spec.harnesses[].model.bedrock.temperature` | `double` |  |  |  |
| `spec.harnesses[].model.bedrock.topP` | `double` |  |  |  |
| `spec.harnesses[].model.gemini` | `AwsBedrockAgentCoreHarnessGeminiModel` |  |  |  |
| `spec.harnesses[].model.gemini.apiKeyArn` | `string \| valueFrom` | yes |  | AwsSecretsManagerSecret (`status.outputs.secret_arn`) |
| `spec.harnesses[].model.gemini.modelId` | `string` | yes |  |  |
| `spec.harnesses[].model.gemini.maxTokens` | `int32` |  |  |  |
| `spec.harnesses[].model.gemini.temperature` | `double` |  |  |  |
| `spec.harnesses[].model.gemini.topP` | `double` |  |  |  |
| `spec.harnesses[].model.gemini.topK` | `int32` |  |  |  |
| `spec.harnesses[].model.openai` | `AwsBedrockAgentCoreHarnessOpenAiModel` |  |  |  |
| `spec.harnesses[].model.openai.apiKeyArn` | `string \| valueFrom` | yes |  | AwsSecretsManagerSecret (`status.outputs.secret_arn`) |
| `spec.harnesses[].model.openai.modelId` | `string` | yes |  |  |
| `spec.harnesses[].model.openai.maxTokens` | `int32` |  |  |  |
| `spec.harnesses[].model.openai.temperature` | `double` |  |  |  |
| `spec.harnesses[].model.openai.topP` | `double` |  |  |  |
| `spec.harnesses[].systemPrompts` | `[]AwsBedrockAgentCoreHarnessSystemPrompt` |  |  |  |
| `spec.harnesses[].systemPrompts[].text` | `string` (sensitive) | yes |  |  |
| `spec.harnesses[].tools` | `[]AwsBedrockAgentCoreHarnessTool` |  |  |  |
| `spec.harnesses[].tools[].name` | `string` |  |  |  |
| `spec.harnesses[].tools[].type` | `string` |  |  |  |
| `spec.harnesses[].tools[].remoteMcp` | `AwsBedrockAgentCoreHarnessRemoteMcpTool` |  |  |  |
| `spec.harnesses[].tools[].remoteMcp.url` | `string` (sensitive) | yes |  |  |
| `spec.harnesses[].tools[].remoteMcp.headers` | `map<string, string>` (sensitive) |  |  |  |
| `spec.harnesses[].tools[].agentcoreBrowser` | `AwsBedrockAgentCoreHarnessBrowserTool` |  |  |  |
| `spec.harnesses[].tools[].agentcoreBrowser.browserArn` | `string \| valueFrom` |  |  |  |
| `spec.harnesses[].tools[].agentcoreGateway` | `AwsBedrockAgentCoreHarnessGatewayTool` |  |  |  |
| `spec.harnesses[].tools[].agentcoreGateway.gatewayArn` | `string \| valueFrom` | yes |  | AwsBedrockAgentCoreGateway (`status.outputs.gateway_arn`) |
| `spec.harnesses[].tools[].agentcoreGateway.outboundAuth` | `AwsBedrockAgentCoreHarnessGatewayOutboundAuth` |  |  |  |
| `spec.harnesses[].tools[].agentcoreGateway.outboundAuth.awsIam` | `bool` |  |  |  |
| `spec.harnesses[].tools[].agentcoreGateway.outboundAuth.noAuth` | `bool` |  |  |  |
| `spec.harnesses[].tools[].agentcoreGateway.outboundAuth.oauth` | `AwsBedrockAgentCoreHarnessGatewayOAuth` |  |  |  |
| `spec.harnesses[].tools[].agentcoreGateway.outboundAuth.oauth.providerArn` | `string \| valueFrom` | yes |  |  |
| `spec.harnesses[].tools[].agentcoreGateway.outboundAuth.oauth.scopes` | `[]string` | yes |  |  |
| `spec.harnesses[].tools[].agentcoreGateway.outboundAuth.oauth.customParameters` | `map<string, string>` |  |  |  |
| `spec.harnesses[].tools[].agentcoreGateway.outboundAuth.oauth.defaultReturnUrl` | `string` |  |  |  |
| `spec.harnesses[].tools[].agentcoreGateway.outboundAuth.oauth.grantType` | `string` |  |  |  |
| `spec.harnesses[].tools[].inlineFunction` | `AwsBedrockAgentCoreHarnessInlineFunctionTool` |  |  |  |
| `spec.harnesses[].tools[].inlineFunction.description` | `string` | yes |  |  |
| `spec.harnesses[].tools[].inlineFunction.inputSchema` | `string` (sensitive) | yes |  |  |
| `spec.harnesses[].tools[].agentcoreCodeInterpreter` | `AwsBedrockAgentCoreHarnessCodeInterpreterTool` |  |  |  |
| `spec.harnesses[].tools[].agentcoreCodeInterpreter.codeInterpreterArn` | `string \| valueFrom` |  |  |  |
| `spec.harnesses[].skillPaths` | `[]string` |  |  |  |
| `spec.harnesses[].memory` | `AwsBedrockAgentCoreHarnessMemory` |  |  |  |
| `spec.harnesses[].memory.memoryArn` | `string \| valueFrom` | yes |  | AwsBedrockAgentCoreMemory (`status.outputs.memory_arn`) |
| `spec.harnesses[].memory.actorId` | `string` |  |  |  |
| `spec.harnesses[].memory.messagesCount` | `int32` |  |  |  |
| `spec.harnesses[].memory.retrieval` | `AwsBedrockAgentCoreHarnessMemoryRetrieval` |  |  |  |
| `spec.harnesses[].memory.retrieval.namespace` | `string` | yes |  |  |
| `spec.harnesses[].memory.retrieval.relevanceScore` | `double` |  |  |  |
| `spec.harnesses[].memory.retrieval.strategyId` | `string` |  |  |  |
| `spec.harnesses[].memory.retrieval.topK` | `int32` |  |  |  |
| `spec.harnesses[].environmentVariables` | `map<string, string>` (sensitive) |  |  |  |
| `spec.harnesses[].runtimeEnvironment` | `AwsBedrockAgentCoreHarnessRuntimeEnvironment` |  |  |  |
| `spec.harnesses[].runtimeEnvironment.agentRuntimeArn` | `string \| valueFrom` |  |  | AwsBedrockAgentCoreRuntime (`status.outputs.agent_runtime_arn`) |
| `spec.harnesses[].runtimeEnvironment.filesystems` | `[]AwsBedrockAgentCoreHarnessFilesystem` |  |  |  |
| `spec.harnesses[].runtimeEnvironment.filesystems[].mountPath` | `string` | yes |  |  |
| `spec.harnesses[].runtimeEnvironment.filesystems[].efsAccessPointArn` | `string \| valueFrom` |  |  | AwsEfsAccessPoint (`status.outputs.access_point_arn`) |
| `spec.harnesses[].runtimeEnvironment.filesystems[].s3AccessPointArn` | `string \| valueFrom` |  |  |  |
| `spec.harnesses[].runtimeEnvironment.filesystems[].sessionStorage` | `bool` |  |  |  |
| `spec.harnesses[].runtimeEnvironment.lifecycle` | `AwsBedrockAgentCoreHarnessLifecycle` |  |  |  |
| `spec.harnesses[].runtimeEnvironment.lifecycle.idleRuntimeSessionTimeoutSeconds` | `int32` |  |  |  |
| `spec.harnesses[].runtimeEnvironment.lifecycle.maxLifetimeSeconds` | `int32` |  |  |  |
| `spec.harnesses[].runtimeEnvironment.network` | `AwsBedrockAgentCoreHarnessNetwork` |  |  |  |
| `spec.harnesses[].runtimeEnvironment.network.mode` | `string` |  |  |  |
| `spec.harnesses[].runtimeEnvironment.network.vpcConfig` | `AwsBedrockAgentCoreHarnessVpcConfig` |  |  |  |
| `spec.harnesses[].runtimeEnvironment.network.vpcConfig.subnets` | `[]string \| valueFrom` | yes |  | AwsSubnet (`status.outputs.subnet_id`) |
| `spec.harnesses[].runtimeEnvironment.network.vpcConfig.securityGroups` | `[]string \| valueFrom` | yes |  | AwsSecurityGroup (`status.outputs.security_group_id`) |
| `spec.harnesses[].runtimeEnvironment.network.vpcConfig.requireServiceS3Endpoint` | `bool` |  |  |  |
| `spec.harnesses[].containerImageUri` | `string` |  |  |  |
| `spec.harnesses[].customJwtAuthorizer` | `AwsBedrockAgentCoreJwtAuthorizer` |  |  |  |
| `spec.harnesses[].customJwtAuthorizer.discoveryUrl` | `string` | yes |  |  |
| `spec.harnesses[].customJwtAuthorizer.allowedAudience` | `[]string` |  |  |  |
| `spec.harnesses[].customJwtAuthorizer.allowedClients` | `[]string` |  |  |  |
| `spec.harnesses[].customJwtAuthorizer.allowedScopes` | `[]string` |  |  |  |
| `spec.harnesses[].customJwtAuthorizer.allowedWorkloads` | `AwsBedrockAgentCoreAllowedWorkloads` |  |  |  |
| `spec.harnesses[].customJwtAuthorizer.allowedWorkloads.workloadIdentities` | `[]string` | yes |  |  |
| `spec.harnesses[].customJwtAuthorizer.allowedWorkloads.hostingEnvironmentArns` | `[]string` | yes |  |  |
| `spec.harnesses[].customJwtAuthorizer.customClaims` | `[]AwsBedrockAgentCoreCustomClaim` |  |  |  |
| `spec.harnesses[].customJwtAuthorizer.customClaims[].claimName` | `string` | yes |  |  |
| `spec.harnesses[].customJwtAuthorizer.customClaims[].valueType` | `string` |  |  |  |
| `spec.harnesses[].customJwtAuthorizer.customClaims[].matchOperator` | `string` |  |  |  |
| `spec.harnesses[].customJwtAuthorizer.customClaims[].matchValue` | `string` |  |  |  |
| `spec.harnesses[].customJwtAuthorizer.customClaims[].matchValues` | `[]string` |  |  |  |
| `spec.harnesses[].customJwtAuthorizer.privateEndpoint` | `AwsBedrockAgentCorePrivateEndpoint` |  |  |  |
| `spec.harnesses[].customJwtAuthorizer.privateEndpoint.managedVpc` | `AwsBedrockAgentCoreManagedVpcEndpoint` |  |  |  |
| `spec.harnesses[].customJwtAuthorizer.privateEndpoint.managedVpc.vpcId` | `string \| valueFrom` | yes |  | AwsVpc (`status.outputs.vpc_id`) |
| `spec.harnesses[].customJwtAuthorizer.privateEndpoint.managedVpc.subnetIds` | `[]string \| valueFrom` | yes |  | AwsSubnet (`status.outputs.subnet_id`) |
| `spec.harnesses[].customJwtAuthorizer.privateEndpoint.managedVpc.securityGroupIds` | `[]string \| valueFrom` |  |  | AwsSecurityGroup (`status.outputs.security_group_id`) |
| `spec.harnesses[].customJwtAuthorizer.privateEndpoint.managedVpc.endpointIpAddressType` | `string` |  |  |  |
| `spec.harnesses[].customJwtAuthorizer.privateEndpoint.managedVpc.routingDomain` | `string` | yes |  |  |
| `spec.harnesses[].customJwtAuthorizer.privateEndpoint.managedVpc.tags` | `map<string, string>` |  |  |  |
| `spec.harnesses[].customJwtAuthorizer.privateEndpoint.selfManagedLattice` | `AwsBedrockAgentCoreLatticeEndpoint` |  |  |  |
| `spec.harnesses[].customJwtAuthorizer.privateEndpoint.selfManagedLattice.resourceConfigurationId` | `string` | yes |  |  |
| `spec.harnesses[].customJwtAuthorizer.privateEndpointOverrides` | `[]AwsBedrockAgentCorePrivateEndpointOverride` |  |  |  |
| `spec.harnesses[].customJwtAuthorizer.privateEndpointOverrides[].domain` | `string` | yes |  |  |
| `spec.harnesses[].customJwtAuthorizer.privateEndpointOverrides[].privateEndpoint` | `AwsBedrockAgentCorePrivateEndpoint` | yes |  |  |
| `spec.harnesses[].customJwtAuthorizer.privateEndpointOverrides[].privateEndpoint.managedVpc` | `AwsBedrockAgentCoreManagedVpcEndpoint` |  |  |  |
| `spec.harnesses[].customJwtAuthorizer.privateEndpointOverrides[].privateEndpoint.managedVpc.vpcId` | `string \| valueFrom` | yes |  | AwsVpc (`status.outputs.vpc_id`) |
| `spec.harnesses[].customJwtAuthorizer.privateEndpointOverrides[].privateEndpoint.managedVpc.subnetIds` | `[]string \| valueFrom` | yes |  | AwsSubnet (`status.outputs.subnet_id`) |
| `spec.harnesses[].customJwtAuthorizer.privateEndpointOverrides[].privateEndpoint.managedVpc.securityGroupIds` | `[]string \| valueFrom` |  |  | AwsSecurityGroup (`status.outputs.security_group_id`) |
| `spec.harnesses[].customJwtAuthorizer.privateEndpointOverrides[].privateEndpoint.managedVpc.endpointIpAddressType` | `string` |  |  |  |
| `spec.harnesses[].customJwtAuthorizer.privateEndpointOverrides[].privateEndpoint.managedVpc.routingDomain` | `string` | yes |  |  |
| `spec.harnesses[].customJwtAuthorizer.privateEndpointOverrides[].privateEndpoint.managedVpc.tags` | `map<string, string>` |  |  |  |
| `spec.harnesses[].customJwtAuthorizer.privateEndpointOverrides[].privateEndpoint.selfManagedLattice` | `AwsBedrockAgentCoreLatticeEndpoint` |  |  |  |
| `spec.harnesses[].customJwtAuthorizer.privateEndpointOverrides[].privateEndpoint.selfManagedLattice.resourceConfigurationId` | `string` | yes |  |  |
| `spec.harnesses[].allowedTools` | `[]string` |  |  |  |
| `spec.harnesses[].maxIterations` | `int32` |  |  |  |
| `spec.harnesses[].maxTokens` | `int32` |  |  |  |
| `spec.harnesses[].timeoutSeconds` | `int32` |  |  |  |
| `spec.harnesses[].truncation` | `AwsBedrockAgentCoreHarnessTruncation` |  |  |  |
| `spec.harnesses[].truncation.strategy` | `string` |  |  |  |
| `spec.harnesses[].truncation.slidingWindow` | `AwsBedrockAgentCoreHarnessSlidingWindow` |  |  |  |
| `spec.harnesses[].truncation.slidingWindow.messagesCount` | `int32` |  |  |  |
| `spec.harnesses[].truncation.summarization` | `AwsBedrockAgentCoreHarnessSummarization` |  |  |  |
| `spec.harnesses[].truncation.summarization.summaryRatio` | `double` |  |  |  |
| `spec.harnesses[].truncation.summarization.preserveRecentMessages` | `int32` |  |  |  |
| `spec.harnesses[].truncation.summarization.summarizationSystemPrompt` | `string` (sensitive) |  |  |  |
| `spec.onlineEvaluationConfigs` | `[]AwsBedrockAgentCoreOnlineEvaluationConfig` |  |  |  |
| `spec.onlineEvaluationConfigs[].name` | `string` | yes |  |  |
| `spec.onlineEvaluationConfigs[].description` | `string` |  |  |  |
| `spec.onlineEvaluationConfigs[].executionRoleArn` | `string \| valueFrom` | yes |  | AwsIamRole (`status.outputs.role_arn`) |
| `spec.onlineEvaluationConfigs[].enabled` | `bool` |  |  |  |
| `spec.onlineEvaluationConfigs[].dataSource` | `AwsBedrockAgentCoreOnlineEvaluationDataSource` | yes |  |  |
| `spec.onlineEvaluationConfigs[].dataSource.logGroupNames` | `[]string \| valueFrom` | yes |  | AwsCloudwatchLogGroup (`status.outputs.log_group_name`) |
| `spec.onlineEvaluationConfigs[].dataSource.serviceNames` | `[]string` | yes |  |  |
| `spec.onlineEvaluationConfigs[].evaluatorIds` | `[]string` | yes |  |  |
| `spec.onlineEvaluationConfigs[].rule` | `AwsBedrockAgentCoreOnlineEvaluationRule` | yes |  |  |
| `spec.onlineEvaluationConfigs[].rule.samplingPercentage` | `double` |  |  |  |
| `spec.onlineEvaluationConfigs[].rule.filters` | `[]AwsBedrockAgentCoreOnlineEvaluationFilter` |  |  |  |
| `spec.onlineEvaluationConfigs[].rule.filters[].key` | `string` | yes |  |  |
| `spec.onlineEvaluationConfigs[].rule.filters[].operator` | `string` |  |  |  |
| `spec.onlineEvaluationConfigs[].rule.filters[].stringValue` | `string` |  |  |  |
| `spec.onlineEvaluationConfigs[].rule.filters[].booleanValue` | `bool` |  |  |  |
| `spec.onlineEvaluationConfigs[].rule.filters[].doubleValue` | `double` |  |  |  |
| `spec.onlineEvaluationConfigs[].rule.sessionTimeoutMinutes` | `int32` |  |  |  |

## Field Details

### spec.region

`string` · required

- rule: {"string":{"minLen":"1"}}

### spec.evaluators

`[]AwsBedrockAgentCoreEvaluator`

- rule: evaluator must set exactly one of llm_as_a_judge or code_based

### spec.evaluators[].name

`string` · required

- rule: {"string":{"minLen":"1","pattern":"^[a-zA-Z][a-zA-Z0-9_]{0,47}$"}}

### spec.evaluators[].description

`string`

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"maxLen":"200"}}

### spec.evaluators[].level

`string`

- rule: {"string":{"in":["TOOL_CALL","TRACE","SESSION"]}}

### spec.evaluators[].kmsKeyArn

`string | valueFrom`

- references: AwsKmsKey (`status.outputs.key_arn`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsKmsKey, name: <that resource's name>, fieldPath: status.outputs.key_arn}} -- a bare string does not parse

### spec.evaluators[].llmAsAJudge

`AwsBedrockAgentCoreLlmJudge`

### spec.evaluators[].llmAsAJudge.instructions

`string` · required · sensitive

- rule: {"required":true}

### spec.evaluators[].llmAsAJudge.model

`AwsBedrockAgentCoreJudgeModel` · required

- rule: {"required":true}

### spec.evaluators[].llmAsAJudge.model.modelId

`string` · required

- rule: {"string":{"minLen":"1"}}

### spec.evaluators[].llmAsAJudge.model.additionalModelRequestFields

`object`

### spec.evaluators[].llmAsAJudge.model.inference

`AwsBedrockAgentCoreJudgeInference`

### spec.evaluators[].llmAsAJudge.model.inference.maxTokens

`int32`

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","int32":{"gte":1}}

### spec.evaluators[].llmAsAJudge.model.inference.stopSequences

`[]string`

- rule: {"repeated":{"maxItems":"2500","items":{"string":{"minLen":"1"}}}}

### spec.evaluators[].llmAsAJudge.model.inference.temperature

`double` · optional (explicit presence)

- rule: {"double":{"lte":1,"gte":0}}

### spec.evaluators[].llmAsAJudge.model.inference.topP

`double` · optional (explicit presence)

- rule: {"double":{"lte":1,"gte":0}}

### spec.evaluators[].llmAsAJudge.ratingScale

`AwsBedrockAgentCoreRatingScale` · required

- rule: {"required":true}
- rule: rating scale must define exactly one of categorical or numerical entries

### spec.evaluators[].llmAsAJudge.ratingScale.categorical

`[]AwsBedrockAgentCoreCategoricalRating`

### spec.evaluators[].llmAsAJudge.ratingScale.categorical[].label

`string` · required

- rule: {"string":{"minLen":"1","maxLen":"100"}}

### spec.evaluators[].llmAsAJudge.ratingScale.categorical[].definition

`string` · required

- rule: {"string":{"minLen":"1"}}

### spec.evaluators[].llmAsAJudge.ratingScale.numerical

`[]AwsBedrockAgentCoreNumericalRating`

### spec.evaluators[].llmAsAJudge.ratingScale.numerical[].label

`string` · required

- rule: {"string":{"minLen":"1","maxLen":"100"}}

### spec.evaluators[].llmAsAJudge.ratingScale.numerical[].definition

`string` · required

- rule: {"string":{"minLen":"1"}}

### spec.evaluators[].llmAsAJudge.ratingScale.numerical[].value

`double`

- rule: {"double":{"gte":0}}

### spec.evaluators[].codeBased

`AwsBedrockAgentCoreCodeEvaluator`

### spec.evaluators[].codeBased.lambdaArn

`string | valueFrom` · required

- references: AwsLambda (`status.outputs.function_arn`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsLambda, name: <that resource's name>, fieldPath: status.outputs.function_arn}} -- a bare string does not parse

### spec.evaluators[].codeBased.timeoutSeconds

`int32`

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","int32":{"lte":300,"gte":1}}

### spec.harnesses

`[]AwsBedrockAgentCoreHarness`

- rule: each allowed_tools entry must match the name of a tool configured in tools

### spec.harnesses[].name

`string` · required

- rule: {"string":{"minLen":"1","maxLen":"40"}}

### spec.harnesses[].executionRoleArn

`string | valueFrom` · required

- references: AwsIamRole (`status.outputs.role_arn`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsIamRole, name: <that resource's name>, fieldPath: status.outputs.role_arn}} -- a bare string does not parse

### spec.harnesses[].model

`AwsBedrockAgentCoreHarnessModel` · required

- rule: {"required":true}
- rule: harness model must set exactly one of bedrock, gemini, or openai

### spec.harnesses[].model.bedrock

`AwsBedrockAgentCoreHarnessBedrockModel`

### spec.harnesses[].model.bedrock.modelId

`string` · required

- rule: {"string":{"minLen":"1"}}

### spec.harnesses[].model.bedrock.maxTokens

`int32`

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","int32":{"gte":1}}

### spec.harnesses[].model.bedrock.temperature

`double` · optional (explicit presence)

- rule: {"double":{"lte":2,"gte":0}}

### spec.harnesses[].model.bedrock.topP

`double` · optional (explicit presence)

- rule: {"double":{"lte":1,"gte":0}}

### spec.harnesses[].model.gemini

`AwsBedrockAgentCoreHarnessGeminiModel`

### spec.harnesses[].model.gemini.apiKeyArn

`string | valueFrom` · required

- references: AwsSecretsManagerSecret (`status.outputs.secret_arn`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsSecretsManagerSecret, name: <that resource's name>, fieldPath: status.outputs.secret_arn}} -- a bare string does not parse

### spec.harnesses[].model.gemini.modelId

`string` · required

- rule: {"string":{"minLen":"1"}}

### spec.harnesses[].model.gemini.maxTokens

`int32`

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","int32":{"gte":1}}

### spec.harnesses[].model.gemini.temperature

`double` · optional (explicit presence)

- rule: {"double":{"lte":2,"gte":0}}

### spec.harnesses[].model.gemini.topP

`double` · optional (explicit presence)

- rule: {"double":{"lte":1,"gte":0}}

### spec.harnesses[].model.gemini.topK

`int32` · optional (explicit presence)

- rule: {"int32":{"lte":500,"gte":0}}

### spec.harnesses[].model.openai

`AwsBedrockAgentCoreHarnessOpenAiModel`

### spec.harnesses[].model.openai.apiKeyArn

`string | valueFrom` · required

- references: AwsSecretsManagerSecret (`status.outputs.secret_arn`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsSecretsManagerSecret, name: <that resource's name>, fieldPath: status.outputs.secret_arn}} -- a bare string does not parse

### spec.harnesses[].model.openai.modelId

`string` · required

- rule: {"string":{"minLen":"1"}}

### spec.harnesses[].model.openai.maxTokens

`int32`

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","int32":{"gte":1}}

### spec.harnesses[].model.openai.temperature

`double` · optional (explicit presence)

- rule: {"double":{"lte":2,"gte":0}}

### spec.harnesses[].model.openai.topP

`double` · optional (explicit presence)

- rule: {"double":{"lte":1,"gte":0}}

### spec.harnesses[].systemPrompts

`[]AwsBedrockAgentCoreHarnessSystemPrompt`

### spec.harnesses[].systemPrompts[].text

`string` · required · sensitive

- rule: {"required":true}

### spec.harnesses[].tools

`[]AwsBedrockAgentCoreHarnessTool`

- rule: set exactly the config arm matching type (remote_mcp, agentcore_gateway, and inline_function require theirs; agentcore_browser and agentcore_code_interpreter may omit theirs to use the AWS default tool)

### spec.harnesses[].tools[].name

`string`

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"maxLen":"64"}}

### spec.harnesses[].tools[].type

`string`

- rule: {"string":{"in":["remote_mcp","agentcore_browser","agentcore_gateway","inline_function","agentcore_code_interpreter"]}}

### spec.harnesses[].tools[].remoteMcp

`AwsBedrockAgentCoreHarnessRemoteMcpTool`

### spec.harnesses[].tools[].remoteMcp.url

`string` · required · sensitive

- rule: {"required":true}

### spec.harnesses[].tools[].remoteMcp.headers

`map<string, string>` · sensitive

### spec.harnesses[].tools[].agentcoreBrowser

`AwsBedrockAgentCoreHarnessBrowserTool`

### spec.harnesses[].tools[].agentcoreBrowser.browserArn

`string | valueFrom`

- rule: write as {value: <literal>} or {valueFrom: {kind: <Kind>, name: <that resource's name>, fieldPath: status.outputs.<output>}} -- a bare string does not parse

### spec.harnesses[].tools[].agentcoreGateway

`AwsBedrockAgentCoreHarnessGatewayTool`

### spec.harnesses[].tools[].agentcoreGateway.gatewayArn

`string | valueFrom` · required

- references: AwsBedrockAgentCoreGateway (`status.outputs.gateway_arn`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsBedrockAgentCoreGateway, name: <that resource's name>, fieldPath: status.outputs.gateway_arn}} -- a bare string does not parse

### spec.harnesses[].tools[].agentcoreGateway.outboundAuth

`AwsBedrockAgentCoreHarnessGatewayOutboundAuth`

- rule: set exactly one of aws_iam, no_auth, or oauth

### spec.harnesses[].tools[].agentcoreGateway.outboundAuth.awsIam

`bool`

### spec.harnesses[].tools[].agentcoreGateway.outboundAuth.noAuth

`bool`

### spec.harnesses[].tools[].agentcoreGateway.outboundAuth.oauth

`AwsBedrockAgentCoreHarnessGatewayOAuth`

### spec.harnesses[].tools[].agentcoreGateway.outboundAuth.oauth.providerArn

`string | valueFrom` · required

- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: <Kind>, name: <that resource's name>, fieldPath: status.outputs.<output>}} -- a bare string does not parse

### spec.harnesses[].tools[].agentcoreGateway.outboundAuth.oauth.scopes

`[]string` · required

- rule: {"repeated":{"minItems":"1","items":{"string":{"minLen":"1"}}}}

### spec.harnesses[].tools[].agentcoreGateway.outboundAuth.oauth.customParameters

`map<string, string>`

### spec.harnesses[].tools[].agentcoreGateway.outboundAuth.oauth.defaultReturnUrl

`string`

### spec.harnesses[].tools[].agentcoreGateway.outboundAuth.oauth.grantType

`string`

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"in":["CLIENT_CREDENTIALS","AUTHORIZATION_CODE","TOKEN_EXCHANGE"]}}

### spec.harnesses[].tools[].inlineFunction

`AwsBedrockAgentCoreHarnessInlineFunctionTool`

### spec.harnesses[].tools[].inlineFunction.description

`string` · required

- rule: {"string":{"minLen":"1"}}

### spec.harnesses[].tools[].inlineFunction.inputSchema

`string` · required · sensitive

- rule: {"required":true}

### spec.harnesses[].tools[].agentcoreCodeInterpreter

`AwsBedrockAgentCoreHarnessCodeInterpreterTool`

### spec.harnesses[].tools[].agentcoreCodeInterpreter.codeInterpreterArn

`string | valueFrom`

- rule: write as {value: <literal>} or {valueFrom: {kind: <Kind>, name: <that resource's name>, fieldPath: status.outputs.<output>}} -- a bare string does not parse

### spec.harnesses[].skillPaths

`[]string`

- rule: {"repeated":{"items":{"string":{"minLen":"1"}}}}

### spec.harnesses[].memory

`AwsBedrockAgentCoreHarnessMemory`

### spec.harnesses[].memory.memoryArn

`string | valueFrom` · required

- references: AwsBedrockAgentCoreMemory (`status.outputs.memory_arn`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsBedrockAgentCoreMemory, name: <that resource's name>, fieldPath: status.outputs.memory_arn}} -- a bare string does not parse

### spec.harnesses[].memory.actorId

`string`

### spec.harnesses[].memory.messagesCount

`int32`

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","int32":{"gte":1}}

### spec.harnesses[].memory.retrieval

`AwsBedrockAgentCoreHarnessMemoryRetrieval`

### spec.harnesses[].memory.retrieval.namespace

`string` · required

- rule: {"string":{"minLen":"1"}}

### spec.harnesses[].memory.retrieval.relevanceScore

`double` · optional (explicit presence)

- rule: {"double":{"lte":1,"gte":0}}

### spec.harnesses[].memory.retrieval.strategyId

`string`

### spec.harnesses[].memory.retrieval.topK

`int32`

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","int32":{"gte":1}}

### spec.harnesses[].environmentVariables

`map<string, string>` · sensitive

### spec.harnesses[].runtimeEnvironment

`AwsBedrockAgentCoreHarnessRuntimeEnvironment`

### spec.harnesses[].runtimeEnvironment.agentRuntimeArn

`string | valueFrom`

- references: AwsBedrockAgentCoreRuntime (`status.outputs.agent_runtime_arn`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsBedrockAgentCoreRuntime, name: <that resource's name>, fieldPath: status.outputs.agent_runtime_arn}} -- a bare string does not parse

### spec.harnesses[].runtimeEnvironment.filesystems

`[]AwsBedrockAgentCoreHarnessFilesystem`

- rule: each filesystem must set exactly one of efs_access_point_arn, s3_access_point_arn, or session_storage

### spec.harnesses[].runtimeEnvironment.filesystems[].mountPath

`string` · required

- rule: {"string":{"minLen":"1"}}

### spec.harnesses[].runtimeEnvironment.filesystems[].efsAccessPointArn

`string | valueFrom`

- references: AwsEfsAccessPoint (`status.outputs.access_point_arn`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsEfsAccessPoint, name: <that resource's name>, fieldPath: status.outputs.access_point_arn}} -- a bare string does not parse

### spec.harnesses[].runtimeEnvironment.filesystems[].s3AccessPointArn

`string | valueFrom`

- rule: write as {value: <literal>} or {valueFrom: {kind: <Kind>, name: <that resource's name>, fieldPath: status.outputs.<output>}} -- a bare string does not parse

### spec.harnesses[].runtimeEnvironment.filesystems[].sessionStorage

`bool`

### spec.harnesses[].runtimeEnvironment.lifecycle

`AwsBedrockAgentCoreHarnessLifecycle`

### spec.harnesses[].runtimeEnvironment.lifecycle.idleRuntimeSessionTimeoutSeconds

`int32`

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","int32":{"gte":1}}

### spec.harnesses[].runtimeEnvironment.lifecycle.maxLifetimeSeconds

`int32`

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","int32":{"gte":1}}

### spec.harnesses[].runtimeEnvironment.network

`AwsBedrockAgentCoreHarnessNetwork`

- rule: vpc_config is required when mode is VPC and forbidden otherwise

### spec.harnesses[].runtimeEnvironment.network.mode

`string`

- rule: {"string":{"in":["PUBLIC","VPC"]}}

### spec.harnesses[].runtimeEnvironment.network.vpcConfig

`AwsBedrockAgentCoreHarnessVpcConfig`

### spec.harnesses[].runtimeEnvironment.network.vpcConfig.subnets

`[]string | valueFrom` · required

- references: AwsSubnet (`status.outputs.subnet_id`)
- rule: {"required":true,"repeated":{"minItems":"1"}}
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsSubnet, name: <that resource's name>, fieldPath: status.outputs.subnet_id}} -- a bare string does not parse

### spec.harnesses[].runtimeEnvironment.network.vpcConfig.securityGroups

`[]string | valueFrom` · required

- references: AwsSecurityGroup (`status.outputs.security_group_id`)
- rule: {"required":true,"repeated":{"minItems":"1"}}
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsSecurityGroup, name: <that resource's name>, fieldPath: status.outputs.security_group_id}} -- a bare string does not parse

### spec.harnesses[].runtimeEnvironment.network.vpcConfig.requireServiceS3Endpoint

`bool`

### spec.harnesses[].containerImageUri

`string`

### spec.harnesses[].customJwtAuthorizer

`AwsBedrockAgentCoreJwtAuthorizer`

### spec.harnesses[].customJwtAuthorizer.discoveryUrl

`string` · required

- rule: {"string":{"minLen":"1"}}

### spec.harnesses[].customJwtAuthorizer.allowedAudience

`[]string`

- rule: {"repeated":{"items":{"string":{"minLen":"1"}}}}

### spec.harnesses[].customJwtAuthorizer.allowedClients

`[]string`

- rule: {"repeated":{"items":{"string":{"minLen":"1"}}}}

### spec.harnesses[].customJwtAuthorizer.allowedScopes

`[]string`

- rule: {"repeated":{"items":{"string":{"minLen":"1"}}}}

### spec.harnesses[].customJwtAuthorizer.allowedWorkloads

`AwsBedrockAgentCoreAllowedWorkloads`

### spec.harnesses[].customJwtAuthorizer.allowedWorkloads.workloadIdentities

`[]string` · required

- rule: {"repeated":{"minItems":"1","maxItems":"10","items":{"string":{"minLen":"1"}}}}

### spec.harnesses[].customJwtAuthorizer.allowedWorkloads.hostingEnvironmentArns

`[]string` · required

- rule: {"repeated":{"minItems":"1","maxItems":"10","items":{"string":{"minLen":"1"}}}}

### spec.harnesses[].customJwtAuthorizer.customClaims

`[]AwsBedrockAgentCoreCustomClaim`

- rule: custom claim must set exactly one of match_value or match_values

### spec.harnesses[].customJwtAuthorizer.customClaims[].claimName

`string` · required

- rule: {"string":{"minLen":"1","maxLen":"255","pattern":"^[A-Za-z0-9_.:-]+$"}}

### spec.harnesses[].customJwtAuthorizer.customClaims[].valueType

`string`

- rule: {"string":{"in":["STRING","STRING_ARRAY"]}}

### spec.harnesses[].customJwtAuthorizer.customClaims[].matchOperator

`string`

- rule: {"string":{"in":["EQUALS","CONTAINS","CONTAINS_ANY"]}}

### spec.harnesses[].customJwtAuthorizer.customClaims[].matchValue

`string`

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"maxLen":"255","pattern":"^[A-Za-z0-9_.:-]+$"}}

### spec.harnesses[].customJwtAuthorizer.customClaims[].matchValues

`[]string`

- rule: {"repeated":{"items":{"string":{"minLen":"1","maxLen":"255","pattern":"^[A-Za-z0-9_.:-]+$"}}}}

### spec.harnesses[].customJwtAuthorizer.privateEndpoint

`AwsBedrockAgentCorePrivateEndpoint`

- rule: private endpoint must set exactly one of managed_vpc or self_managed_lattice

### spec.harnesses[].customJwtAuthorizer.privateEndpoint.managedVpc

`AwsBedrockAgentCoreManagedVpcEndpoint`

### spec.harnesses[].customJwtAuthorizer.privateEndpoint.managedVpc.vpcId

`string | valueFrom` · required

- references: AwsVpc (`status.outputs.vpc_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsVpc, name: <that resource's name>, fieldPath: status.outputs.vpc_id}} -- a bare string does not parse

### spec.harnesses[].customJwtAuthorizer.privateEndpoint.managedVpc.subnetIds

`[]string | valueFrom` · required

- references: AwsSubnet (`status.outputs.subnet_id`)
- rule: {"required":true,"repeated":{"minItems":"1"}}
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsSubnet, name: <that resource's name>, fieldPath: status.outputs.subnet_id}} -- a bare string does not parse

### spec.harnesses[].customJwtAuthorizer.privateEndpoint.managedVpc.securityGroupIds

`[]string | valueFrom`

- references: AwsSecurityGroup (`status.outputs.security_group_id`)
- rule: {"repeated":{"maxItems":"5"}}
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsSecurityGroup, name: <that resource's name>, fieldPath: status.outputs.security_group_id}} -- a bare string does not parse

### spec.harnesses[].customJwtAuthorizer.privateEndpoint.managedVpc.endpointIpAddressType

`string`

- rule: {"string":{"in":["IPV4","IPV6"]}}

### spec.harnesses[].customJwtAuthorizer.privateEndpoint.managedVpc.routingDomain

`string` · required

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"minLen":"3","maxLen":"255"}}

### spec.harnesses[].customJwtAuthorizer.privateEndpoint.managedVpc.tags

`map<string, string>`

### spec.harnesses[].customJwtAuthorizer.privateEndpoint.selfManagedLattice

`AwsBedrockAgentCoreLatticeEndpoint`

### spec.harnesses[].customJwtAuthorizer.privateEndpoint.selfManagedLattice.resourceConfigurationId

`string` · required

- rule: {"string":{"minLen":"1"}}

### spec.harnesses[].customJwtAuthorizer.privateEndpointOverrides

`[]AwsBedrockAgentCorePrivateEndpointOverride`

- rule: {"repeated":{"maxItems":"5"}}

### spec.harnesses[].customJwtAuthorizer.privateEndpointOverrides[].domain

`string` · required

- rule: {"string":{"minLen":"1","maxLen":"253"}}

### spec.harnesses[].customJwtAuthorizer.privateEndpointOverrides[].privateEndpoint

`AwsBedrockAgentCorePrivateEndpoint` · required

- rule: {"required":true}
- rule: private endpoint must set exactly one of managed_vpc or self_managed_lattice

### spec.harnesses[].customJwtAuthorizer.privateEndpointOverrides[].privateEndpoint.managedVpc

`AwsBedrockAgentCoreManagedVpcEndpoint`

### spec.harnesses[].customJwtAuthorizer.privateEndpointOverrides[].privateEndpoint.managedVpc.vpcId

`string | valueFrom` · required

- references: AwsVpc (`status.outputs.vpc_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsVpc, name: <that resource's name>, fieldPath: status.outputs.vpc_id}} -- a bare string does not parse

### spec.harnesses[].customJwtAuthorizer.privateEndpointOverrides[].privateEndpoint.managedVpc.subnetIds

`[]string | valueFrom` · required

- references: AwsSubnet (`status.outputs.subnet_id`)
- rule: {"required":true,"repeated":{"minItems":"1"}}
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsSubnet, name: <that resource's name>, fieldPath: status.outputs.subnet_id}} -- a bare string does not parse

### spec.harnesses[].customJwtAuthorizer.privateEndpointOverrides[].privateEndpoint.managedVpc.securityGroupIds

`[]string | valueFrom`

- references: AwsSecurityGroup (`status.outputs.security_group_id`)
- rule: {"repeated":{"maxItems":"5"}}
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsSecurityGroup, name: <that resource's name>, fieldPath: status.outputs.security_group_id}} -- a bare string does not parse

### spec.harnesses[].customJwtAuthorizer.privateEndpointOverrides[].privateEndpoint.managedVpc.endpointIpAddressType

`string`

- rule: {"string":{"in":["IPV4","IPV6"]}}

### spec.harnesses[].customJwtAuthorizer.privateEndpointOverrides[].privateEndpoint.managedVpc.routingDomain

`string` · required

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"minLen":"3","maxLen":"255"}}

### spec.harnesses[].customJwtAuthorizer.privateEndpointOverrides[].privateEndpoint.managedVpc.tags

`map<string, string>`

### spec.harnesses[].customJwtAuthorizer.privateEndpointOverrides[].privateEndpoint.selfManagedLattice

`AwsBedrockAgentCoreLatticeEndpoint`

### spec.harnesses[].customJwtAuthorizer.privateEndpointOverrides[].privateEndpoint.selfManagedLattice.resourceConfigurationId

`string` · required

- rule: {"string":{"minLen":"1"}}

### spec.harnesses[].allowedTools

`[]string`

- rule: {"repeated":{"items":{"string":{"minLen":"1"}}}}

### spec.harnesses[].maxIterations

`int32`

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","int32":{"gte":1}}

### spec.harnesses[].maxTokens

`int32`

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","int32":{"gte":1}}

### spec.harnesses[].timeoutSeconds

`int32`

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","int32":{"gte":1}}

### spec.harnesses[].truncation

`AwsBedrockAgentCoreHarnessTruncation`

- rule: set exactly the config matching strategy: sliding_window with 'sliding_window', summarization with 'summarization', neither with 'none'

### spec.harnesses[].truncation.strategy

`string`

- rule: {"string":{"in":["sliding_window","summarization","none"]}}

### spec.harnesses[].truncation.slidingWindow

`AwsBedrockAgentCoreHarnessSlidingWindow`

### spec.harnesses[].truncation.slidingWindow.messagesCount

`int32`

- rule: {"int32":{"gte":1}}

### spec.harnesses[].truncation.summarization

`AwsBedrockAgentCoreHarnessSummarization`

### spec.harnesses[].truncation.summarization.summaryRatio

`double` · optional (explicit presence)

- rule: {"double":{"lte":1,"gte":0}}

### spec.harnesses[].truncation.summarization.preserveRecentMessages

`int32`

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","int32":{"gte":1}}

### spec.harnesses[].truncation.summarization.summarizationSystemPrompt

`string` · sensitive

### spec.onlineEvaluationConfigs

`[]AwsBedrockAgentCoreOnlineEvaluationConfig`

### spec.onlineEvaluationConfigs[].name

`string` · required

- rule: {"string":{"minLen":"1","pattern":"^[a-zA-Z][a-zA-Z0-9_]{0,47}$"}}

### spec.onlineEvaluationConfigs[].description

`string`

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"maxLen":"200"}}

### spec.onlineEvaluationConfigs[].executionRoleArn

`string | valueFrom` · required

- references: AwsIamRole (`status.outputs.role_arn`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsIamRole, name: <that resource's name>, fieldPath: status.outputs.role_arn}} -- a bare string does not parse

### spec.onlineEvaluationConfigs[].enabled

`bool` · optional (explicit presence)

### spec.onlineEvaluationConfigs[].dataSource

`AwsBedrockAgentCoreOnlineEvaluationDataSource` · required

- rule: {"required":true}

### spec.onlineEvaluationConfigs[].dataSource.logGroupNames

`[]string | valueFrom` · required

- references: AwsCloudwatchLogGroup (`status.outputs.log_group_name`)
- rule: {"required":true,"repeated":{"minItems":"1"}}
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsCloudwatchLogGroup, name: <that resource's name>, fieldPath: status.outputs.log_group_name}} -- a bare string does not parse

### spec.onlineEvaluationConfigs[].dataSource.serviceNames

`[]string` · required

- rule: {"repeated":{"minItems":"1","items":{"string":{"minLen":"1"}}}}

### spec.onlineEvaluationConfigs[].evaluatorIds

`[]string` · required

- rule: {"repeated":{"minItems":"1","maxItems":"10","unique":true,"items":{"string":{"minLen":"1"}}}}

### spec.onlineEvaluationConfigs[].rule

`AwsBedrockAgentCoreOnlineEvaluationRule` · required

- rule: {"required":true}

### spec.onlineEvaluationConfigs[].rule.samplingPercentage

`double`

- rule: {"double":{"lte":100,"gte":0.01}}

### spec.onlineEvaluationConfigs[].rule.filters

`[]AwsBedrockAgentCoreOnlineEvaluationFilter`

- rule: {"repeated":{"maxItems":"5"}}
- rule: each filter must set exactly one of string_value, boolean_value, or double_value

### spec.onlineEvaluationConfigs[].rule.filters[].key

`string` · required

- rule: {"string":{"minLen":"1","maxLen":"256"}}

### spec.onlineEvaluationConfigs[].rule.filters[].operator

`string`

- rule: {"string":{"in":["Equals","NotEquals","GreaterThan","LessThan","GreaterThanOrEqual","LessThanOrEqual","Contains","NotContains"]}}

### spec.onlineEvaluationConfigs[].rule.filters[].stringValue

`string`

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"maxLen":"1024"}}

### spec.onlineEvaluationConfigs[].rule.filters[].booleanValue

`bool` · optional (explicit presence)

### spec.onlineEvaluationConfigs[].rule.filters[].doubleValue

`double` · optional (explicit presence)

### spec.onlineEvaluationConfigs[].rule.sessionTimeoutMinutes

`int32`

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","int32":{"lte":60,"gte":1}}

## Validation Rules

- `at_least_one_evaluation_resource`: set at least one of evaluators, harnesses, or online_evaluation_configs
- `evaluator_names_unique`: evaluators entries must have unique names
- `harness_names_unique`: harnesses entries must have unique names
- `online_evaluation_config_names_unique`: online_evaluation_configs entries must have unique names
- `online_config_evaluators_resolvable`: each online_evaluation_configs evaluator entry must be an AWS builtin (Builtin.*), a full custom evaluator ID (name-XXXXXXXXXX), or the name of an evaluator defined in this spec

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AwsBedrockAgentCoreEvaluation, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.evaluator_ids` | `map<string, string>` |  |
| `status.outputs.evaluator_arns` | `map<string, string>` |  |
| `status.outputs.harness_ids` | `map<string, string>` |  |
| `status.outputs.harness_arns` | `map<string, string>` |  |
| `status.outputs.online_evaluation_config_ids` | `map<string, string>` |  |
| `status.outputs.online_evaluation_config_arns` | `map<string, string>` |  |
| `status.outputs.online_evaluation_output_log_groups` | `map<string, string>` |  |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.evaluators[].kmsKeyArn` | AwsKmsKey | `status.outputs.key_arn` |
| `spec.evaluators[].codeBased.lambdaArn` | AwsLambda | `status.outputs.function_arn` |
| `spec.harnesses[].executionRoleArn` | AwsIamRole | `status.outputs.role_arn` |
| `spec.harnesses[].model.gemini.apiKeyArn` | AwsSecretsManagerSecret | `status.outputs.secret_arn` |
| `spec.harnesses[].model.openai.apiKeyArn` | AwsSecretsManagerSecret | `status.outputs.secret_arn` |
| `spec.harnesses[].tools[].agentcoreGateway.gatewayArn` | AwsBedrockAgentCoreGateway | `status.outputs.gateway_arn` |
| `spec.harnesses[].memory.memoryArn` | AwsBedrockAgentCoreMemory | `status.outputs.memory_arn` |
| `spec.harnesses[].runtimeEnvironment.agentRuntimeArn` | AwsBedrockAgentCoreRuntime | `status.outputs.agent_runtime_arn` |
| `spec.harnesses[].runtimeEnvironment.filesystems[].efsAccessPointArn` | AwsEfsAccessPoint | `status.outputs.access_point_arn` |
| `spec.harnesses[].runtimeEnvironment.network.vpcConfig.subnets` | AwsSubnet | `status.outputs.subnet_id` |
| `spec.harnesses[].runtimeEnvironment.network.vpcConfig.securityGroups` | AwsSecurityGroup | `status.outputs.security_group_id` |
| `spec.harnesses[].customJwtAuthorizer.privateEndpoint.managedVpc.vpcId` | AwsVpc | `status.outputs.vpc_id` |
| `spec.harnesses[].customJwtAuthorizer.privateEndpoint.managedVpc.subnetIds` | AwsSubnet | `status.outputs.subnet_id` |
| `spec.harnesses[].customJwtAuthorizer.privateEndpoint.managedVpc.securityGroupIds` | AwsSecurityGroup | `status.outputs.security_group_id` |
| `spec.harnesses[].customJwtAuthorizer.privateEndpointOverrides[].privateEndpoint.managedVpc.vpcId` | AwsVpc | `status.outputs.vpc_id` |
| `spec.harnesses[].customJwtAuthorizer.privateEndpointOverrides[].privateEndpoint.managedVpc.subnetIds` | AwsSubnet | `status.outputs.subnet_id` |
| `spec.harnesses[].customJwtAuthorizer.privateEndpointOverrides[].privateEndpoint.managedVpc.securityGroupIds` | AwsSecurityGroup | `status.outputs.security_group_id` |
| `spec.onlineEvaluationConfigs[].executionRoleArn` | AwsIamRole | `status.outputs.role_arn` |
| `spec.onlineEvaluationConfigs[].dataSource.logGroupNames` | AwsCloudwatchLogGroup | `status.outputs.log_group_name` |

## See Also

- [Overview](../README.md)
