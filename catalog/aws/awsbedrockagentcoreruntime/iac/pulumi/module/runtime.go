package module

import (
	"encoding/json"
	"sort"

	"github.com/pkg/errors"
	awsbedrockagentcoreruntimev1alpha1 "github.com/plantonhq/planton/catalog/aws/awsbedrockagentcoreruntime/v1alpha1"
	foreignkeyv1 "github.com/plantonhq/planton/shared/foreignkey/v1"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/bedrock"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// runtime creates the AgentCore agent runtime, its named endpoints, and
// the runtime's resource-based policy, and exports outputs.
//
// Lifecycle facts the renders below depend on:
//   - every configuration change creates a new runtime VERSION in place;
//     switching the artifact arm (code <-> container) replaces the
//     runtime (provider-enforced);
//   - endpoints pin or float across versions -- an endpoint without an
//     explicit version tracks the latest;
//   - the resource policy attaches to the runtime's own ARN (the
//     provider resource accepts any AgentCore ARN; this module scopes it
//     to the runtime it deploys).
func runtime(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) error {
	spec := locals.Spec

	args := &bedrock.AgentcoreAgentRuntimeArgs{
		// AWS's runtime-name charset (letter first, then letters/digits/_)
		// is stricter than metadata.name conventions, so the name is an
		// explicit spec field. Changing it replaces the runtime.
		AgentRuntimeName: pulumi.String(spec.RuntimeName),
		// Required by AWS: the role the AgentCore service assumes to pull
		// the image / read the code bundle and run the agent.
		RoleArn: pulumi.String(spec.RoleArn.GetValue()),
		Tags:    pulumi.ToStringMap(locals.AwsTags),
	}

	if spec.Description != "" {
		args.Description = pulumi.String(spec.Description)
	}
	if len(spec.EnvironmentVariables) > 0 {
		args.EnvironmentVariables = pulumi.ToStringMap(spec.EnvironmentVariables)
	}

	// Optional+Computed single-entry list at the provider: sent only when
	// set so the module never fights AWS's defaults.
	if spec.Lifecycle != nil {
		lifecycle := &bedrock.AgentcoreAgentRuntimeLifecycleConfigurationArgs{}
		if spec.Lifecycle.IdleRuntimeSessionTimeoutSeconds != 0 {
			lifecycle.IdleRuntimeSessionTimeout = pulumi.Int(int(spec.Lifecycle.IdleRuntimeSessionTimeoutSeconds))
		}
		if spec.Lifecycle.MaxLifetimeSeconds != 0 {
			lifecycle.MaxLifetime = pulumi.Int(int(spec.Lifecycle.MaxLifetimeSeconds))
		}
		args.LifecycleConfigurations = bedrock.AgentcoreAgentRuntimeLifecycleConfigurationArray{lifecycle}
	}

	// Exactly one artifact arm (spec-validated). Switching arms replaces
	// the runtime.
	artifact := &bedrock.AgentcoreAgentRuntimeAgentRuntimeArtifactArgs{}
	if spec.Artifact.Container != nil {
		artifact.ContainerConfiguration = &bedrock.AgentcoreAgentRuntimeAgentRuntimeArtifactContainerConfigurationArgs{
			ContainerUri: pulumi.String(spec.Artifact.Container.ImageUri),
		}
	}
	if spec.Artifact.Code != nil {
		code := spec.Artifact.Code
		s3 := &bedrock.AgentcoreAgentRuntimeAgentRuntimeArtifactCodeConfigurationCodeS3Args{
			Bucket: pulumi.String(code.S3.Bucket.GetValue()),
			Prefix: pulumi.String(code.S3.Prefix),
		}
		if code.S3.VersionId != "" {
			s3.VersionId = pulumi.String(code.S3.VersionId)
		}
		artifact.CodeConfiguration = &bedrock.AgentcoreAgentRuntimeAgentRuntimeArtifactCodeConfigurationArgs{
			Runtime:     pulumi.String(code.Runtime),
			EntryPoints: pulumi.ToStringArray(code.EntryPoint),
			Code: &bedrock.AgentcoreAgentRuntimeAgentRuntimeArtifactCodeConfigurationCodeArgs{
				S3: s3,
			},
		}
	}
	args.AgentRuntimeArtifact = artifact

	// Required by AWS. VPC mode carries the placement (spec-validated
	// pairing).
	network := &bedrock.AgentcoreAgentRuntimeNetworkConfigurationArgs{
		NetworkMode: pulumi.String(spec.Network.Mode),
	}
	if spec.Network.VpcConfig != nil {
		network.NetworkModeConfig = &bedrock.AgentcoreAgentRuntimeNetworkConfigurationNetworkModeConfigArgs{
			Subnets:        svrsToStringArray(spec.Network.VpcConfig.Subnets),
			SecurityGroups: svrsToStringArray(spec.Network.VpcConfig.SecurityGroups),
		}
	}
	args.NetworkConfiguration = network

	// HTTP is AWS's default protocol; sent only on an explicit choice so
	// the module never fights the default.
	if spec.ServerProtocol != "" {
		args.ProtocolConfiguration = &bedrock.AgentcoreAgentRuntimeProtocolConfigurationArgs{
			ServerProtocol: pulumi.String(spec.ServerProtocol),
		}
	}

	if len(spec.RequestHeaderAllowlist) > 0 {
		args.RequestHeaderConfiguration = &bedrock.AgentcoreAgentRuntimeRequestHeaderConfigurationArgs{
			RequestHeaderAllowlists: pulumi.ToStringArray(spec.RequestHeaderAllowlist),
		}
	}

	// Inbound JWT authorization. Omitted = AWS IAM (SigV4) callers only.
	if spec.CustomJwtAuthorizer != nil {
		jwt := spec.CustomJwtAuthorizer
		authorizer := &bedrock.AgentcoreAgentRuntimeAuthorizerConfigurationCustomJwtAuthorizerArgs{
			DiscoveryUrl: pulumi.String(jwt.DiscoveryUrl),
		}
		if len(jwt.AllowedAudience) > 0 {
			authorizer.AllowedAudiences = pulumi.ToStringArray(jwt.AllowedAudience)
		}
		if len(jwt.AllowedClients) > 0 {
			authorizer.AllowedClients = pulumi.ToStringArray(jwt.AllowedClients)
		}
		if len(jwt.AllowedScopes) > 0 {
			authorizer.AllowedScopes = pulumi.ToStringArray(jwt.AllowedScopes)
		}
		if jwt.AllowedWorkloads != nil {
			workloads := &bedrock.AgentcoreAgentRuntimeAuthorizerConfigurationCustomJwtAuthorizerAllowedWorkloadConfigurationArgs{}
			if len(jwt.AllowedWorkloads.WorkloadIdentities) > 0 {
				workloads.WorkloadIdentities = pulumi.ToStringArray(jwt.AllowedWorkloads.WorkloadIdentities)
			}
			var environments bedrock.AgentcoreAgentRuntimeAuthorizerConfigurationCustomJwtAuthorizerAllowedWorkloadConfigurationHostingEnvironmentArray
			for _, arn := range jwt.AllowedWorkloads.HostingEnvironmentArns {
				environments = append(environments, &bedrock.AgentcoreAgentRuntimeAuthorizerConfigurationCustomJwtAuthorizerAllowedWorkloadConfigurationHostingEnvironmentArgs{
					Arn: pulumi.String(arn),
				})
			}
			if len(environments) > 0 {
				workloads.HostingEnvironments = environments
			}
			authorizer.AllowedWorkloadConfiguration = workloads
		}
		var claims bedrock.AgentcoreAgentRuntimeAuthorizerConfigurationCustomJwtAuthorizerCustomClaimArray
		for _, c := range jwt.CustomClaims {
			matchValue := &bedrock.AgentcoreAgentRuntimeAuthorizerConfigurationCustomJwtAuthorizerCustomClaimAuthorizingClaimMatchValueClaimMatchValueArgs{}
			if c.MatchValue != "" {
				matchValue.MatchValueString = pulumi.String(c.MatchValue)
			}
			if len(c.MatchValues) > 0 {
				matchValue.MatchValueStringLists = pulumi.ToStringArray(c.MatchValues)
			}
			claims = append(claims, &bedrock.AgentcoreAgentRuntimeAuthorizerConfigurationCustomJwtAuthorizerCustomClaimArgs{
				InboundTokenClaimName:      pulumi.String(c.ClaimName),
				InboundTokenClaimValueType: pulumi.String(c.ValueType),
				AuthorizingClaimMatchValue: &bedrock.AgentcoreAgentRuntimeAuthorizerConfigurationCustomJwtAuthorizerCustomClaimAuthorizingClaimMatchValueArgs{
					ClaimMatchOperator: pulumi.String(c.MatchOperator),
					ClaimMatchValue:    matchValue,
				},
			})
		}
		if len(claims) > 0 {
			authorizer.CustomClaims = claims
		}
		if jwt.PrivateEndpoint != nil {
			authorizer.PrivateEndpoint = privateEndpointArgs(jwt.PrivateEndpoint)
		}
		var overrides bedrock.AgentcoreAgentRuntimeAuthorizerConfigurationCustomJwtAuthorizerPrivateEndpointOverrideArray
		for _, o := range jwt.PrivateEndpointOverrides {
			overrides = append(overrides, &bedrock.AgentcoreAgentRuntimeAuthorizerConfigurationCustomJwtAuthorizerPrivateEndpointOverrideArgs{
				Domain:          pulumi.String(o.Domain),
				PrivateEndpoint: overridePrivateEndpointArgs(o.PrivateEndpoint),
			})
		}
		if len(overrides) > 0 {
			authorizer.PrivateEndpointOverrides = overrides
		}
		args.AuthorizerConfiguration = &bedrock.AgentcoreAgentRuntimeAuthorizerConfigurationArgs{
			CustomJwtAuthorizer: authorizer,
		}
	}

	// Filesystem mounts (max 5; exactly one source arm each,
	// spec-validated). session_storage's only argument is its mount path.
	var filesystems bedrock.AgentcoreAgentRuntimeFilesystemConfigurationArray
	for _, f := range spec.Filesystems {
		filesystem := &bedrock.AgentcoreAgentRuntimeFilesystemConfigurationArgs{}
		if f.EfsAccessPointArn.GetValue() != "" {
			filesystem.EfsAccessPoint = &bedrock.AgentcoreAgentRuntimeFilesystemConfigurationEfsAccessPointArgs{
				AccessPointArn: pulumi.String(f.EfsAccessPointArn.GetValue()),
				MountPath:      pulumi.String(f.MountPath),
			}
		}
		if f.S3FilesAccessPointArn.GetValue() != "" {
			filesystem.S3FilesAccessPoint = &bedrock.AgentcoreAgentRuntimeFilesystemConfigurationS3FilesAccessPointArgs{
				AccessPointArn: pulumi.String(f.S3FilesAccessPointArn.GetValue()),
				MountPath:      pulumi.String(f.MountPath),
			}
		}
		if f.SessionStorage {
			filesystem.SessionStorage = &bedrock.AgentcoreAgentRuntimeFilesystemConfigurationSessionStorageArgs{
				MountPath: pulumi.String(f.MountPath),
			}
		}
		filesystems = append(filesystems, filesystem)
	}
	if len(filesystems) > 0 {
		args.FilesystemConfigurations = filesystems
	}

	createdRuntime, err := bedrock.NewAgentcoreAgentRuntime(ctx, spec.RuntimeName, args, pulumi.Provider(provider))
	if err != nil {
		return errors.Wrap(err, "create agent runtime")
	}

	ctx.Export(OpAgentRuntimeId, createdRuntime.AgentRuntimeId)
	ctx.Export(OpAgentRuntimeArn, createdRuntime.AgentRuntimeArn)
	ctx.Export(OpAgentRuntimeVersion, createdRuntime.AgentRuntimeVersion)
	ctx.Export(OpWorkloadIdentityArn, createdRuntime.WorkloadIdentityDetails.Index(pulumi.Int(0)).WorkloadIdentityArn())

	// Named serving endpoints, keyed by their stable entry names.
	// Iteration is name-sorted for deterministic previews.
	endpointArns := pulumi.StringMap{}
	for _, e := range sortedEndpoints(spec.Endpoints) {
		endpointArgs := &bedrock.AgentcoreAgentRuntimeEndpointArgs{
			AgentRuntimeId: createdRuntime.AgentRuntimeId,
			Name:           pulumi.String(e.Name),
			Tags:           pulumi.ToStringMap(locals.AwsTags),
		}
		if e.Description != "" {
			endpointArgs.Description = pulumi.String(e.Description)
		}
		if e.AgentRuntimeVersion != "" {
			endpointArgs.AgentRuntimeVersion = pulumi.String(e.AgentRuntimeVersion)
		}
		createdEndpoint, err := bedrock.NewAgentcoreAgentRuntimeEndpoint(ctx, "endpoint-"+e.Name, endpointArgs,
			pulumi.Provider(provider), pulumi.DependsOn([]pulumi.Resource{createdRuntime}))
		if err != nil {
			return errors.Wrapf(err, "create endpoint %q", e.Name)
		}
		endpointArns[e.Name] = createdEndpoint.AgentRuntimeEndpointArn
	}
	ctx.Export(OpEndpointArns, endpointArns)

	// Resource-based policy on the runtime's own ARN.
	if spec.ResourcePolicy != nil {
		policyBytes, err := json.Marshal(spec.ResourcePolicy.AsMap())
		if err != nil {
			return errors.Wrap(err, "marshal resource policy")
		}
		if _, err := bedrock.NewAgentcoreResourcePolicy(ctx, spec.RuntimeName+"-policy",
			&bedrock.AgentcoreResourcePolicyArgs{
				ResourceArn: createdRuntime.AgentRuntimeArn,
				Policy:      pulumi.String(string(policyBytes)),
			},
			pulumi.Provider(provider), pulumi.DependsOn([]pulumi.Resource{createdRuntime})); err != nil {
			return errors.Wrap(err, "create resource policy")
		}
	}

	return nil
}

func privateEndpointArgs(in *awsbedrockagentcoreruntimev1alpha1.AwsBedrockAgentCorePrivateEndpoint) *bedrock.AgentcoreAgentRuntimeAuthorizerConfigurationCustomJwtAuthorizerPrivateEndpointArgs {
	out := &bedrock.AgentcoreAgentRuntimeAuthorizerConfigurationCustomJwtAuthorizerPrivateEndpointArgs{}
	if in.ManagedVpc != nil {
		managed := &bedrock.AgentcoreAgentRuntimeAuthorizerConfigurationCustomJwtAuthorizerPrivateEndpointManagedVpcResourceArgs{
			VpcIdentifier:         pulumi.String(in.ManagedVpc.VpcId.GetValue()),
			SubnetIds:             svrsToStringArray(in.ManagedVpc.SubnetIds),
			EndpointIpAddressType: pulumi.String(in.ManagedVpc.EndpointIpAddressType),
		}
		if len(in.ManagedVpc.SecurityGroupIds) > 0 {
			managed.SecurityGroupIds = svrsToStringArray(in.ManagedVpc.SecurityGroupIds)
		}
		if in.ManagedVpc.RoutingDomain != "" {
			managed.RoutingDomain = pulumi.String(in.ManagedVpc.RoutingDomain)
		}
		if len(in.ManagedVpc.Tags) > 0 {
			managed.Tags = pulumi.ToStringMap(in.ManagedVpc.Tags)
		}
		out.ManagedVpcResource = managed
	}
	if in.SelfManagedLattice != nil {
		out.SelfManagedLatticeResource = &bedrock.AgentcoreAgentRuntimeAuthorizerConfigurationCustomJwtAuthorizerPrivateEndpointSelfManagedLatticeResourceArgs{
			ResourceConfigurationIdentifier: pulumi.String(in.SelfManagedLattice.ResourceConfigurationId),
		}
	}
	return out
}

// overridePrivateEndpointArgs mirrors privateEndpointArgs for the
// override entries' distinct bridge type.
func overridePrivateEndpointArgs(in *awsbedrockagentcoreruntimev1alpha1.AwsBedrockAgentCorePrivateEndpoint) *bedrock.AgentcoreAgentRuntimeAuthorizerConfigurationCustomJwtAuthorizerPrivateEndpointOverridePrivateEndpointArgs {
	out := &bedrock.AgentcoreAgentRuntimeAuthorizerConfigurationCustomJwtAuthorizerPrivateEndpointOverridePrivateEndpointArgs{}
	if in.ManagedVpc != nil {
		managed := &bedrock.AgentcoreAgentRuntimeAuthorizerConfigurationCustomJwtAuthorizerPrivateEndpointOverridePrivateEndpointManagedVpcResourceArgs{
			VpcIdentifier:         pulumi.String(in.ManagedVpc.VpcId.GetValue()),
			SubnetIds:             svrsToStringArray(in.ManagedVpc.SubnetIds),
			EndpointIpAddressType: pulumi.String(in.ManagedVpc.EndpointIpAddressType),
		}
		if len(in.ManagedVpc.SecurityGroupIds) > 0 {
			managed.SecurityGroupIds = svrsToStringArray(in.ManagedVpc.SecurityGroupIds)
		}
		if in.ManagedVpc.RoutingDomain != "" {
			managed.RoutingDomain = pulumi.String(in.ManagedVpc.RoutingDomain)
		}
		if len(in.ManagedVpc.Tags) > 0 {
			managed.Tags = pulumi.ToStringMap(in.ManagedVpc.Tags)
		}
		out.ManagedVpcResource = managed
	}
	if in.SelfManagedLattice != nil {
		out.SelfManagedLatticeResource = &bedrock.AgentcoreAgentRuntimeAuthorizerConfigurationCustomJwtAuthorizerPrivateEndpointOverridePrivateEndpointSelfManagedLatticeResourceArgs{
			ResourceConfigurationIdentifier: pulumi.String(in.SelfManagedLattice.ResourceConfigurationId),
		}
	}
	return out
}

func svrsToStringArray(in []*foreignkeyv1.StringValueOrRef) pulumi.StringArray {
	out := pulumi.StringArray{}
	for _, ref := range in {
		out = append(out, pulumi.String(ref.GetValue()))
	}
	return out
}

func sortedEndpoints(in []*awsbedrockagentcoreruntimev1alpha1.AwsBedrockAgentCoreRuntimeEndpoint) []*awsbedrockagentcoreruntimev1alpha1.AwsBedrockAgentCoreRuntimeEndpoint {
	out := append([]*awsbedrockagentcoreruntimev1alpha1.AwsBedrockAgentCoreRuntimeEndpoint{}, in...)
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}
