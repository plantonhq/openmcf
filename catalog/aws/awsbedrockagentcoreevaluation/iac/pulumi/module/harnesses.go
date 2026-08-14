package module

import (
	"github.com/pkg/errors"
	awsbedrockagentcoreevaluationv1alpha1 "github.com/plantonhq/planton/catalog/aws/awsbedrockagentcoreevaluation/v1alpha1"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/bedrock"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// harnesses creates the bundle's agent test benches and exports
// outputs.
//
// Lifecycle facts the renders below depend on:
//   - the harness name has no rename (the provider replaces on change);
//     everything else updates in place;
//   - creation waits CREATING -> READY (the harness's terminal success
//     state -- the other two resources use ACTIVE), retrying briefly
//     upstream while the execution role propagates;
//   - the provider models several single-member wrappers as PLURALIZED
//     arrays (environments, network_mode_configs, truncation configs,
//     filesystem arms) -- this engine's SDK inherits those plurals, so
//     the singular spec arms render as one-element arrays here
//     (bridge-shape quirk, same class as the SageMaker
//     S3DataSources/DefaultDomainIdLists renders);
//   - non-Bedrock model vendors authenticate via a Secrets Manager ARN
//     the harness reads at run time -- the module sends the reference,
//     never a key.
func harnesses(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) error {
	spec := locals.Spec

	harnessIds := pulumi.StringMap{}
	harnessArns := pulumi.StringMap{}

	// Iteration is name-sorted for deterministic previews.
	for _, h := range sortedHarnesses(spec.Harnesses) {
		args := &bedrock.AgentcoreHarnessArgs{
			HarnessName:      pulumi.String(h.Name),
			ExecutionRoleArn: pulumi.String(h.ExecutionRoleArn.GetValue()),
			Model:            harnessModel(h.Model),
			Tags:             pulumi.ToStringMap(locals.AwsTags),
		}

		if len(h.SystemPrompts) > 0 {
			var prompts bedrock.AgentcoreHarnessSystemPromptArray
			for _, p := range h.SystemPrompts {
				prompts = append(prompts, &bedrock.AgentcoreHarnessSystemPromptArgs{
					Text: pulumi.String(p.Text),
				})
			}
			args.SystemPrompts = prompts
		}

		if len(h.Tools) > 0 {
			tools, err := harnessTools(h.Tools)
			if err != nil {
				return errors.Wrapf(err, "harness %q tools", h.Name)
			}
			args.Tools = tools
		}

		if len(h.SkillPaths) > 0 {
			var skills bedrock.AgentcoreHarnessSkillArray
			for _, s := range h.SkillPaths {
				skills = append(skills, &bedrock.AgentcoreHarnessSkillArgs{Path: pulumi.String(s)})
			}
			args.Skills = skills
		}

		if h.Memory != nil {
			memoryConfig := &bedrock.AgentcoreHarnessMemoryAgentcoreMemoryConfigurationArgs{
				Arn: pulumi.String(h.Memory.MemoryArn.GetValue()),
			}
			if h.Memory.ActorId != "" {
				memoryConfig.ActorId = pulumi.String(h.Memory.ActorId)
			}
			if h.Memory.MessagesCount > 0 {
				memoryConfig.MessagesCount = pulumi.Int(int(h.Memory.MessagesCount))
			}
			if h.Memory.Retrieval != nil {
				retrieval := &bedrock.AgentcoreHarnessMemoryAgentcoreMemoryConfigurationRetrievalConfigArgs{
					MapBlockKey: pulumi.String(h.Memory.Retrieval.Namespace),
				}
				if h.Memory.Retrieval.RelevanceScore != nil {
					retrieval.RelevanceScore = pulumi.Float64(*h.Memory.Retrieval.RelevanceScore)
				}
				if h.Memory.Retrieval.StrategyId != "" {
					retrieval.StrategyId = pulumi.String(h.Memory.Retrieval.StrategyId)
				}
				if h.Memory.Retrieval.TopK > 0 {
					retrieval.TopK = pulumi.Int(int(h.Memory.Retrieval.TopK))
				}
				memoryConfig.RetrievalConfig = retrieval
			}
			args.Memory = &bedrock.AgentcoreHarnessMemoryArgs{
				AgentcoreMemoryConfiguration: memoryConfig,
			}
		}

		if len(h.EnvironmentVariables) > 0 {
			args.EnvironmentVariables = pulumi.ToStringMap(h.EnvironmentVariables)
		}

		if h.RuntimeEnvironment != nil {
			args.Environments = bedrock.AgentcoreHarnessEnvironmentArray{
				&bedrock.AgentcoreHarnessEnvironmentArgs{
					AgentcoreRuntimeEnvironments: bedrock.AgentcoreHarnessEnvironmentAgentcoreRuntimeEnvironmentArray{
						runtimeEnvironment(h.RuntimeEnvironment),
					},
				},
			}
		}

		if h.ContainerImageUri != "" {
			args.EnvironmentArtifact = &bedrock.AgentcoreHarnessEnvironmentArtifactArgs{
				ContainerConfiguration: &bedrock.AgentcoreHarnessEnvironmentArtifactContainerConfigurationArgs{
					ContainerUri: pulumi.String(h.ContainerImageUri),
				},
			}
		}

		if h.CustomJwtAuthorizer != nil {
			args.AuthorizerConfiguration = &bedrock.AgentcoreHarnessAuthorizerConfigurationArgs{
				CustomJwtAuthorizer: jwtAuthorizer(h.CustomJwtAuthorizer),
			}
		}

		if len(h.AllowedTools) > 0 {
			args.AllowedTools = pulumi.ToStringArray(h.AllowedTools)
		}
		if h.MaxIterations > 0 {
			args.MaxIterations = pulumi.Int(int(h.MaxIterations))
		}
		if h.MaxTokens > 0 {
			args.MaxTokens = pulumi.Int(int(h.MaxTokens))
		}
		if h.TimeoutSeconds > 0 {
			args.TimeoutSeconds = pulumi.Int(int(h.TimeoutSeconds))
		}

		if h.Truncation != nil {
			truncation := &bedrock.AgentcoreHarnessTruncationArgs{
				Strategy: pulumi.String(h.Truncation.Strategy),
			}
			config := &bedrock.AgentcoreHarnessTruncationConfigArgs{}
			hasConfig := false
			if h.Truncation.SlidingWindow != nil {
				config.SlidingWindows = bedrock.AgentcoreHarnessTruncationConfigSlidingWindowArray{
					&bedrock.AgentcoreHarnessTruncationConfigSlidingWindowArgs{
						MessagesCount: pulumi.Int(int(h.Truncation.SlidingWindow.MessagesCount)),
					},
				}
				hasConfig = true
			}
			if h.Truncation.Summarization != nil {
				s := h.Truncation.Summarization
				summarization := &bedrock.AgentcoreHarnessTruncationConfigSummarizationArgs{}
				if s.SummaryRatio != nil {
					summarization.SummaryRatio = pulumi.Float64(*s.SummaryRatio)
				}
				if s.PreserveRecentMessages > 0 {
					summarization.PreserveRecentMessages = pulumi.Int(int(s.PreserveRecentMessages))
				}
				if s.SummarizationSystemPrompt != "" {
					summarization.SummarizationSystemPrompt = pulumi.String(s.SummarizationSystemPrompt)
				}
				config.Summarizations = bedrock.AgentcoreHarnessTruncationConfigSummarizationArray{summarization}
				hasConfig = true
			}
			if hasConfig {
				truncation.Configs = bedrock.AgentcoreHarnessTruncationConfigArray{config}
			}
			args.Truncations = bedrock.AgentcoreHarnessTruncationArray{truncation}
		}

		resource, err := bedrock.NewAgentcoreHarness(ctx, "harness-"+h.Name, args, pulumi.Provider(provider))
		if err != nil {
			return errors.Wrapf(err, "create harness %q", h.Name)
		}
		harnessIds[h.Name] = resource.HarnessId
		harnessArns[h.Name] = resource.Arn
	}

	ctx.Export(OpHarnessIds, harnessIds)
	ctx.Export(OpHarnessArns, harnessArns)
	return nil
}

// harnessModel renders the exactly-one vendor arm (spec-validated).
func harnessModel(m *awsbedrockagentcoreevaluationv1alpha1.AwsBedrockAgentCoreHarnessModel) *bedrock.AgentcoreHarnessModelArgs {
	args := &bedrock.AgentcoreHarnessModelArgs{}
	if m.Bedrock != nil {
		config := &bedrock.AgentcoreHarnessModelBedrockModelConfigArgs{
			ModelId: pulumi.String(m.Bedrock.ModelId),
		}
		if m.Bedrock.MaxTokens > 0 {
			config.MaxTokens = pulumi.Int(int(m.Bedrock.MaxTokens))
		}
		if m.Bedrock.Temperature != nil {
			config.Temperature = pulumi.Float64(*m.Bedrock.Temperature)
		}
		if m.Bedrock.TopP != nil {
			config.TopP = pulumi.Float64(*m.Bedrock.TopP)
		}
		args.BedrockModelConfig = config
	}
	if m.Gemini != nil {
		config := &bedrock.AgentcoreHarnessModelGeminiModelConfigArgs{
			ApiKeyArn: pulumi.String(m.Gemini.ApiKeyArn.GetValue()),
			ModelId:   pulumi.String(m.Gemini.ModelId),
		}
		if m.Gemini.MaxTokens > 0 {
			config.MaxTokens = pulumi.Int(int(m.Gemini.MaxTokens))
		}
		if m.Gemini.Temperature != nil {
			config.Temperature = pulumi.Float64(*m.Gemini.Temperature)
		}
		if m.Gemini.TopP != nil {
			config.TopP = pulumi.Float64(*m.Gemini.TopP)
		}
		if m.Gemini.TopK != nil {
			config.TopK = pulumi.Int(int(*m.Gemini.TopK))
		}
		args.GeminiModelConfig = config
	}
	if m.Openai != nil {
		config := &bedrock.AgentcoreHarnessModelOpenaiModelConfigArgs{
			ApiKeyArn: pulumi.String(m.Openai.ApiKeyArn.GetValue()),
			ModelId:   pulumi.String(m.Openai.ModelId),
		}
		if m.Openai.MaxTokens > 0 {
			config.MaxTokens = pulumi.Int(int(m.Openai.MaxTokens))
		}
		if m.Openai.Temperature != nil {
			config.Temperature = pulumi.Float64(*m.Openai.Temperature)
		}
		if m.Openai.TopP != nil {
			config.TopP = pulumi.Float64(*m.Openai.TopP)
		}
		args.OpenaiModelConfig = config
	}
	return args
}

// harnessTools renders the tool list; each entry's arm was already
// validated against its declared type (stricter than the provider,
// which silently takes the first configured arm).
func harnessTools(in []*awsbedrockagentcoreevaluationv1alpha1.AwsBedrockAgentCoreHarnessTool) (bedrock.AgentcoreHarnessToolArray, error) {
	var out bedrock.AgentcoreHarnessToolArray
	for _, t := range in {
		tool := &bedrock.AgentcoreHarnessToolArgs{
			Type: pulumi.String(t.Type),
		}
		if t.Name != "" {
			tool.Name = pulumi.String(t.Name)
		}
		config := &bedrock.AgentcoreHarnessToolConfigArgs{}
		hasConfig := false
		if t.RemoteMcp != nil {
			mcp := &bedrock.AgentcoreHarnessToolConfigRemoteMcpArgs{
				Url: pulumi.String(t.RemoteMcp.Url),
			}
			if len(t.RemoteMcp.Headers) > 0 {
				mcp.Headers = pulumi.ToStringMap(t.RemoteMcp.Headers)
			}
			config.RemoteMcp = mcp
			hasConfig = true
		}
		if t.AgentcoreBrowser != nil {
			browser := &bedrock.AgentcoreHarnessToolConfigAgentcoreBrowserArgs{}
			if t.AgentcoreBrowser.BrowserArn.GetValue() != "" {
				browser.BrowserArn = pulumi.String(t.AgentcoreBrowser.BrowserArn.GetValue())
			}
			config.AgentcoreBrowser = browser
			hasConfig = true
		}
		if t.AgentcoreCodeInterpreter != nil {
			interpreter := &bedrock.AgentcoreHarnessToolConfigAgentcoreCodeInterpreterArgs{}
			if t.AgentcoreCodeInterpreter.CodeInterpreterArn.GetValue() != "" {
				interpreter.CodeInterpreterArn = pulumi.String(t.AgentcoreCodeInterpreter.CodeInterpreterArn.GetValue())
			}
			config.AgentcoreCodeInterpreter = interpreter
			hasConfig = true
		}
		if t.AgentcoreGateway != nil {
			gateway := &bedrock.AgentcoreHarnessToolConfigAgentcoreGatewayArgs{
				GatewayArn: pulumi.String(t.AgentcoreGateway.GatewayArn.GetValue()),
			}
			if t.AgentcoreGateway.OutboundAuth != nil {
				auth := &bedrock.AgentcoreHarnessToolConfigAgentcoreGatewayOutboundAuthArgs{}
				if t.AgentcoreGateway.OutboundAuth.AwsIam {
					auth.AwsIam = pulumi.Bool(true)
				}
				if t.AgentcoreGateway.OutboundAuth.NoAuth {
					auth.None = pulumi.Bool(true)
				}
				if t.AgentcoreGateway.OutboundAuth.Oauth != nil {
					oauth := t.AgentcoreGateway.OutboundAuth.Oauth
					oauthArgs := &bedrock.AgentcoreHarnessToolConfigAgentcoreGatewayOutboundAuthOauthArgs{
						ProviderArn: pulumi.String(oauth.ProviderArn.GetValue()),
						Scopes:      pulumi.ToStringArray(oauth.Scopes),
					}
					if len(oauth.CustomParameters) > 0 {
						oauthArgs.CustomParameters = pulumi.ToStringMap(oauth.CustomParameters)
					}
					if oauth.DefaultReturnUrl != "" {
						oauthArgs.DefaultReturnUrl = pulumi.String(oauth.DefaultReturnUrl)
					}
					if oauth.GrantType != "" {
						oauthArgs.GrantType = pulumi.String(oauth.GrantType)
					}
					auth.Oauth = oauthArgs
				}
				gateway.OutboundAuth = auth
			}
			config.AgentcoreGateway = gateway
			hasConfig = true
		}
		if t.InlineFunction != nil {
			config.InlineFunction = &bedrock.AgentcoreHarnessToolConfigInlineFunctionArgs{
				Description: pulumi.String(t.InlineFunction.Description),
				InputSchema: pulumi.String(t.InlineFunction.InputSchema),
			}
			hasConfig = true
		}
		if hasConfig {
			tool.Config = config
		}
		out = append(out, tool)
	}
	return out, nil
}

// runtimeEnvironment renders the explicit runtime-environment arm.
func runtimeEnvironment(env *awsbedrockagentcoreevaluationv1alpha1.AwsBedrockAgentCoreHarnessRuntimeEnvironment) *bedrock.AgentcoreHarnessEnvironmentAgentcoreRuntimeEnvironmentArgs {
	args := &bedrock.AgentcoreHarnessEnvironmentAgentcoreRuntimeEnvironmentArgs{}
	if env.AgentRuntimeArn.GetValue() != "" {
		args.AgentRuntimeArn = pulumi.String(env.AgentRuntimeArn.GetValue())
	}
	if len(env.Filesystems) > 0 {
		var filesystems bedrock.AgentcoreHarnessEnvironmentAgentcoreRuntimeEnvironmentFilesystemConfigurationArray
		for _, f := range env.Filesystems {
			fs := &bedrock.AgentcoreHarnessEnvironmentAgentcoreRuntimeEnvironmentFilesystemConfigurationArgs{}
			// Exactly one source arm per mount (spec-validated); each arm
			// is a bridge-pluralized single-member wrapper.
			if f.EfsAccessPointArn.GetValue() != "" {
				fs.EfsAccessPoints = bedrock.AgentcoreHarnessEnvironmentAgentcoreRuntimeEnvironmentFilesystemConfigurationEfsAccessPointArray{
					&bedrock.AgentcoreHarnessEnvironmentAgentcoreRuntimeEnvironmentFilesystemConfigurationEfsAccessPointArgs{
						AccessPointArn: pulumi.String(f.EfsAccessPointArn.GetValue()),
						MountPath:      pulumi.String(f.MountPath),
					},
				}
			}
			if f.S3AccessPointArn.GetValue() != "" {
				fs.S3FilesAccessPoints = bedrock.AgentcoreHarnessEnvironmentAgentcoreRuntimeEnvironmentFilesystemConfigurationS3FilesAccessPointArray{
					&bedrock.AgentcoreHarnessEnvironmentAgentcoreRuntimeEnvironmentFilesystemConfigurationS3FilesAccessPointArgs{
						AccessPointArn: pulumi.String(f.S3AccessPointArn.GetValue()),
						MountPath:      pulumi.String(f.MountPath),
					},
				}
			}
			if f.SessionStorage {
				fs.SessionStorages = bedrock.AgentcoreHarnessEnvironmentAgentcoreRuntimeEnvironmentFilesystemConfigurationSessionStorageArray{
					&bedrock.AgentcoreHarnessEnvironmentAgentcoreRuntimeEnvironmentFilesystemConfigurationSessionStorageArgs{
						MountPath: pulumi.String(f.MountPath),
					},
				}
			}
			filesystems = append(filesystems, fs)
		}
		args.FilesystemConfigurations = filesystems
	}
	if env.Lifecycle != nil {
		lifecycle := &bedrock.AgentcoreHarnessEnvironmentAgentcoreRuntimeEnvironmentLifecycleConfigurationArgs{}
		if env.Lifecycle.IdleRuntimeSessionTimeoutSeconds > 0 {
			lifecycle.IdleRuntimeSessionTimeout = pulumi.Int(int(env.Lifecycle.IdleRuntimeSessionTimeoutSeconds))
		}
		if env.Lifecycle.MaxLifetimeSeconds > 0 {
			lifecycle.MaxLifetime = pulumi.Int(int(env.Lifecycle.MaxLifetimeSeconds))
		}
		args.LifecycleConfigurations = bedrock.AgentcoreHarnessEnvironmentAgentcoreRuntimeEnvironmentLifecycleConfigurationArray{lifecycle}
	}
	if env.Network != nil {
		network := &bedrock.AgentcoreHarnessEnvironmentAgentcoreRuntimeEnvironmentNetworkConfigurationArgs{
			NetworkMode: pulumi.String(env.Network.Mode),
		}
		if env.Network.VpcConfig != nil {
			network.NetworkModeConfigs = bedrock.AgentcoreHarnessEnvironmentAgentcoreRuntimeEnvironmentNetworkConfigurationNetworkModeConfigArray{
				&bedrock.AgentcoreHarnessEnvironmentAgentcoreRuntimeEnvironmentNetworkConfigurationNetworkModeConfigArgs{
					Subnets:                  svrsToStringArray(env.Network.VpcConfig.Subnets),
					SecurityGroups:           svrsToStringArray(env.Network.VpcConfig.SecurityGroups),
					RequireServiceS3Endpoint: pulumi.Bool(env.Network.VpcConfig.RequireServiceS3Endpoint),
				},
			}
		}
		args.NetworkConfigurations = bedrock.AgentcoreHarnessEnvironmentAgentcoreRuntimeEnvironmentNetworkConfigurationArray{network}
	}
	return args
}

// jwtAuthorizer renders the shared AgentCore JWT authorizer shape.
func jwtAuthorizer(a *awsbedrockagentcoreevaluationv1alpha1.AwsBedrockAgentCoreJwtAuthorizer) *bedrock.AgentcoreHarnessAuthorizerConfigurationCustomJwtAuthorizerArgs {
	args := &bedrock.AgentcoreHarnessAuthorizerConfigurationCustomJwtAuthorizerArgs{
		DiscoveryUrl: pulumi.String(a.DiscoveryUrl),
	}
	if len(a.AllowedAudience) > 0 {
		args.AllowedAudiences = pulumi.ToStringArray(a.AllowedAudience)
	}
	if len(a.AllowedClients) > 0 {
		args.AllowedClients = pulumi.ToStringArray(a.AllowedClients)
	}
	if len(a.AllowedScopes) > 0 {
		args.AllowedScopes = pulumi.ToStringArray(a.AllowedScopes)
	}
	if a.AllowedWorkloads != nil {
		workloads := &bedrock.AgentcoreHarnessAuthorizerConfigurationCustomJwtAuthorizerAllowedWorkloadConfigurationArgs{}
		if len(a.AllowedWorkloads.WorkloadIdentities) > 0 {
			workloads.WorkloadIdentities = pulumi.ToStringArray(a.AllowedWorkloads.WorkloadIdentities)
		}
		if len(a.AllowedWorkloads.HostingEnvironmentArns) > 0 {
			var environments bedrock.AgentcoreHarnessAuthorizerConfigurationCustomJwtAuthorizerAllowedWorkloadConfigurationHostingEnvironmentArray
			for _, arn := range a.AllowedWorkloads.HostingEnvironmentArns {
				environments = append(environments, &bedrock.AgentcoreHarnessAuthorizerConfigurationCustomJwtAuthorizerAllowedWorkloadConfigurationHostingEnvironmentArgs{
					Arn: pulumi.String(arn),
				})
			}
			workloads.HostingEnvironments = environments
		}
		args.AllowedWorkloadConfiguration = workloads
	}
	if len(a.CustomClaims) > 0 {
		var claims bedrock.AgentcoreHarnessAuthorizerConfigurationCustomJwtAuthorizerCustomClaimArray
		for _, c := range a.CustomClaims {
			matchValue := &bedrock.AgentcoreHarnessAuthorizerConfigurationCustomJwtAuthorizerCustomClaimAuthorizingClaimMatchValueClaimMatchValueArgs{}
			// Exactly one match-value shape (spec-validated); the list
			// arm is a bridge-pluralized wrapper.
			if c.MatchValue != "" {
				matchValue.MatchValueString = pulumi.String(c.MatchValue)
			}
			if len(c.MatchValues) > 0 {
				matchValue.MatchValueStringLists = pulumi.ToStringArray(c.MatchValues)
			}
			claims = append(claims, &bedrock.AgentcoreHarnessAuthorizerConfigurationCustomJwtAuthorizerCustomClaimArgs{
				InboundTokenClaimName:      pulumi.String(c.ClaimName),
				InboundTokenClaimValueType: pulumi.String(c.ValueType),
				AuthorizingClaimMatchValue: &bedrock.AgentcoreHarnessAuthorizerConfigurationCustomJwtAuthorizerCustomClaimAuthorizingClaimMatchValueArgs{
					ClaimMatchOperator: pulumi.String(c.MatchOperator),
					ClaimMatchValue:    matchValue,
				},
			})
		}
		args.CustomClaims = claims
	}
	if a.PrivateEndpoint != nil {
		args.PrivateEndpoint = privateEndpoint(a.PrivateEndpoint)
	}
	if len(a.PrivateEndpointOverrides) > 0 {
		var overrides bedrock.AgentcoreHarnessAuthorizerConfigurationCustomJwtAuthorizerPrivateEndpointOverrideArray
		for _, o := range a.PrivateEndpointOverrides {
			overrides = append(overrides, &bedrock.AgentcoreHarnessAuthorizerConfigurationCustomJwtAuthorizerPrivateEndpointOverrideArgs{
				Domain:          pulumi.String(o.Domain),
				PrivateEndpoint: overridePrivateEndpoint(o.PrivateEndpoint),
			})
		}
		args.PrivateEndpointOverrides = overrides
	}
	return args
}

func privateEndpoint(p *awsbedrockagentcoreevaluationv1alpha1.AwsBedrockAgentCorePrivateEndpoint) *bedrock.AgentcoreHarnessAuthorizerConfigurationCustomJwtAuthorizerPrivateEndpointArgs {
	args := &bedrock.AgentcoreHarnessAuthorizerConfigurationCustomJwtAuthorizerPrivateEndpointArgs{}
	if p.ManagedVpc != nil {
		// AWS requires the address type on the harness's managed
		// endpoint (unlike the runtime's otherwise-identical shape).
		managed := &bedrock.AgentcoreHarnessAuthorizerConfigurationCustomJwtAuthorizerPrivateEndpointManagedVpcResourceArgs{
			VpcIdentifier:         pulumi.String(p.ManagedVpc.VpcId.GetValue()),
			SubnetIds:             svrsToStringArray(p.ManagedVpc.SubnetIds),
			EndpointIpAddressType: pulumi.String(p.ManagedVpc.EndpointIpAddressType),
		}
		if len(p.ManagedVpc.SecurityGroupIds) > 0 {
			managed.SecurityGroupIds = svrsToStringArray(p.ManagedVpc.SecurityGroupIds)
		}
		if p.ManagedVpc.RoutingDomain != "" {
			managed.RoutingDomain = pulumi.String(p.ManagedVpc.RoutingDomain)
		}
		if len(p.ManagedVpc.Tags) > 0 {
			managed.Tags = pulumi.ToStringMap(p.ManagedVpc.Tags)
		}
		args.ManagedVpcResource = managed
	}
	if p.SelfManagedLattice != nil {
		args.SelfManagedLatticeResource = &bedrock.AgentcoreHarnessAuthorizerConfigurationCustomJwtAuthorizerPrivateEndpointSelfManagedLatticeResourceArgs{
			ResourceConfigurationIdentifier: pulumi.String(p.SelfManagedLattice.ResourceConfigurationId),
		}
	}
	return args
}

func overridePrivateEndpoint(p *awsbedrockagentcoreevaluationv1alpha1.AwsBedrockAgentCorePrivateEndpoint) *bedrock.AgentcoreHarnessAuthorizerConfigurationCustomJwtAuthorizerPrivateEndpointOverridePrivateEndpointArgs {
	args := &bedrock.AgentcoreHarnessAuthorizerConfigurationCustomJwtAuthorizerPrivateEndpointOverridePrivateEndpointArgs{}
	if p.ManagedVpc != nil {
		// Same shape as the top-level private endpoint.
		managed := &bedrock.AgentcoreHarnessAuthorizerConfigurationCustomJwtAuthorizerPrivateEndpointOverridePrivateEndpointManagedVpcResourceArgs{
			VpcIdentifier:         pulumi.String(p.ManagedVpc.VpcId.GetValue()),
			SubnetIds:             svrsToStringArray(p.ManagedVpc.SubnetIds),
			EndpointIpAddressType: pulumi.String(p.ManagedVpc.EndpointIpAddressType),
		}
		if len(p.ManagedVpc.SecurityGroupIds) > 0 {
			managed.SecurityGroupIds = svrsToStringArray(p.ManagedVpc.SecurityGroupIds)
		}
		if p.ManagedVpc.RoutingDomain != "" {
			managed.RoutingDomain = pulumi.String(p.ManagedVpc.RoutingDomain)
		}
		if len(p.ManagedVpc.Tags) > 0 {
			managed.Tags = pulumi.ToStringMap(p.ManagedVpc.Tags)
		}
		args.ManagedVpcResource = managed
	}
	if p.SelfManagedLattice != nil {
		args.SelfManagedLatticeResource = &bedrock.AgentcoreHarnessAuthorizerConfigurationCustomJwtAuthorizerPrivateEndpointOverridePrivateEndpointSelfManagedLatticeResourceArgs{
			ResourceConfigurationIdentifier: pulumi.String(p.SelfManagedLattice.ResourceConfigurationId),
		}
	}
	return args
}
