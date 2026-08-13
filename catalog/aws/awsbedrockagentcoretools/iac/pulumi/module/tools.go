package module

import (
	"sort"

	"github.com/pkg/errors"
	awsbedrockagentcoretoolsv1alpha1 "github.com/plantonhq/planton/catalog/aws/awsbedrockagentcoretools/v1alpha1"
	foreignkeyv1 "github.com/plantonhq/planton/shared/foreignkey/v1"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/bedrock"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// tools creates the AgentCore tools bundle's arms and exports outputs.
//
// Lifecycle facts the renders below depend on:
//   - AWS exposes NO update for any of these resources -- every field
//     change recreates the tool. That is cheap: the tools are
//     configuration shells; AWS bills only per session at runtime;
//   - the browser and code interpreter share the certificate shape (a
//     Secrets Manager location); the enterprise policy and recording
//     shapes are browser-only.
func tools(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) error {
	spec := locals.Spec

	// Managed cloud browsers. Iteration is name-sorted for deterministic
	// previews.
	browserIds := pulumi.StringMap{}
	browserArns := pulumi.StringMap{}
	for _, b := range sortedBrowsers(spec.Browsers) {
		args := &bedrock.AgentcoreBrowserArgs{
			Name: pulumi.String(b.Name),
			Tags: pulumi.ToStringMap(locals.AwsTags),
		}
		if b.Description != "" {
			args.Description = pulumi.String(b.Description)
		}
		if b.ExecutionRoleArn.GetValue() != "" {
			args.ExecutionRoleArn = pulumi.String(b.ExecutionRoleArn.GetValue())
		}
		// Required by AWS. VPC mode carries the placement (spec-validated
		// pairing).
		network := &bedrock.AgentcoreBrowserNetworkConfigurationArgs{
			NetworkMode: pulumi.String(b.Network.Mode),
		}
		if b.Network.VpcConfig != nil {
			network.VpcConfig = &bedrock.AgentcoreBrowserNetworkConfigurationVpcConfigArgs{
				Subnets:        svrsToStringArray(b.Network.VpcConfig.Subnets),
				SecurityGroups: svrsToStringArray(b.Network.VpcConfig.SecurityGroups),
			}
		}
		args.NetworkConfiguration = network

		// Cryptographic traffic signing -- rendered only on an explicit
		// choice so the module never fights AWS's default.
		if b.SigningEnabled != nil {
			args.BrowserSigning = &bedrock.AgentcoreBrowserBrowserSigningArgs{
				Enabled: pulumi.Bool(*b.SigningEnabled),
			}
		}

		// Session recording to S3.
		if b.Recording != nil {
			recording := &bedrock.AgentcoreBrowserRecordingArgs{}
			if b.Recording.Enabled != nil {
				recording.Enabled = pulumi.Bool(*b.Recording.Enabled)
			}
			if b.Recording.S3Location != nil {
				recording.S3Location = &bedrock.AgentcoreBrowserRecordingS3LocationArgs{
					Bucket: pulumi.String(b.Recording.S3Location.Bucket.GetValue()),
					Prefix: pulumi.String(b.Recording.S3Location.Prefix),
				}
			}
			args.Recording = recording
		}

		// Chrome enterprise policy files (max 100).
		var policies bedrock.AgentcoreBrowserEnterprisePolicyArray
		for _, p := range b.EnterprisePolicies {
			s3 := &bedrock.AgentcoreBrowserEnterprisePolicyLocationS3Args{
				Bucket: pulumi.String(p.S3.Bucket.GetValue()),
				Prefix: pulumi.String(p.S3.Prefix),
			}
			if p.S3.VersionId != "" {
				s3.VersionId = pulumi.String(p.S3.VersionId)
			}
			policy := &bedrock.AgentcoreBrowserEnterprisePolicyArgs{
				Location: &bedrock.AgentcoreBrowserEnterprisePolicyLocationArgs{
					S3: s3,
				},
			}
			if p.Type != "" {
				policy.Type = pulumi.String(p.Type)
			}
			policies = append(policies, policy)
		}
		if len(policies) > 0 {
			args.EnterprisePolicies = policies
		}

		// Client certificates for mTLS-protected sites (max 200).
		var certificates bedrock.AgentcoreBrowserCertificateArray
		for _, c := range b.Certificates {
			certificates = append(certificates, &bedrock.AgentcoreBrowserCertificateArgs{
				Location: &bedrock.AgentcoreBrowserCertificateLocationArgs{
					SecretsManager: &bedrock.AgentcoreBrowserCertificateLocationSecretsManagerArgs{
						SecretArn: pulumi.String(c.SecretArn.GetValue()),
					},
				},
			})
		}
		if len(certificates) > 0 {
			args.Certificates = certificates
		}

		created, err := bedrock.NewAgentcoreBrowser(ctx, "browser-"+b.Name, args, pulumi.Provider(provider))
		if err != nil {
			return errors.Wrapf(err, "create browser %q", b.Name)
		}
		browserIds[b.Name] = created.BrowserId
		browserArns[b.Name] = created.BrowserArn
	}
	ctx.Export(OpBrowserIds, browserIds)
	ctx.Export(OpBrowserArns, browserArns)

	// Reusable saved browser state (cookies, logins).
	profileIds := pulumi.StringMap{}
	profileArns := pulumi.StringMap{}
	for _, p := range sortedProfiles(spec.BrowserProfiles) {
		args := &bedrock.AgentcoreBrowserProfileArgs{
			Name: pulumi.String(p.Name),
			Tags: pulumi.ToStringMap(locals.AwsTags),
		}
		if p.Description != "" {
			args.Description = pulumi.String(p.Description)
		}
		created, err := bedrock.NewAgentcoreBrowserProfile(ctx, "browser-profile-"+p.Name, args, pulumi.Provider(provider))
		if err != nil {
			return errors.Wrapf(err, "create browser profile %q", p.Name)
		}
		profileIds[p.Name] = created.ProfileId
		profileArns[p.Name] = created.ProfileArn
	}
	ctx.Export(OpBrowserProfileIds, profileIds)
	ctx.Export(OpBrowserProfileArns, profileArns)

	// Managed code-execution sandboxes.
	interpreterIds := pulumi.StringMap{}
	interpreterArns := pulumi.StringMap{}
	for _, c := range sortedInterpreters(spec.CodeInterpreters) {
		args := &bedrock.AgentcoreCodeInterpreterArgs{
			Name: pulumi.String(c.Name),
			Tags: pulumi.ToStringMap(locals.AwsTags),
		}
		if c.Description != "" {
			args.Description = pulumi.String(c.Description)
		}
		if c.ExecutionRoleArn.GetValue() != "" {
			args.ExecutionRoleArn = pulumi.String(c.ExecutionRoleArn.GetValue())
		}
		// Required by AWS. SANDBOX blocks all network access (the safest
		// for untrusted code); VPC mode carries the placement
		// (spec-validated pairing).
		network := &bedrock.AgentcoreCodeInterpreterNetworkConfigurationArgs{
			NetworkMode: pulumi.String(c.Network.Mode),
		}
		if c.Network.VpcConfig != nil {
			network.VpcConfig = &bedrock.AgentcoreCodeInterpreterNetworkConfigurationVpcConfigArgs{
				Subnets:        svrsToStringArray(c.Network.VpcConfig.Subnets),
				SecurityGroups: svrsToStringArray(c.Network.VpcConfig.SecurityGroups),
			}
		}
		args.NetworkConfiguration = network

		// Client certificates for mTLS-protected endpoints the code
		// calls (max 200).
		var certificates bedrock.AgentcoreCodeInterpreterCertificateArray
		for _, cert := range c.Certificates {
			certificates = append(certificates, &bedrock.AgentcoreCodeInterpreterCertificateArgs{
				Location: &bedrock.AgentcoreCodeInterpreterCertificateLocationArgs{
					SecretsManager: &bedrock.AgentcoreCodeInterpreterCertificateLocationSecretsManagerArgs{
						SecretArn: pulumi.String(cert.SecretArn.GetValue()),
					},
				},
			})
		}
		if len(certificates) > 0 {
			args.Certificates = certificates
		}

		created, err := bedrock.NewAgentcoreCodeInterpreter(ctx, "code-interpreter-"+c.Name, args, pulumi.Provider(provider))
		if err != nil {
			return errors.Wrapf(err, "create code interpreter %q", c.Name)
		}
		interpreterIds[c.Name] = created.CodeInterpreterId
		interpreterArns[c.Name] = created.CodeInterpreterArn
	}
	ctx.Export(OpCodeInterpreterIds, interpreterIds)
	ctx.Export(OpCodeInterpreterArns, interpreterArns)

	return nil
}

func svrsToStringArray(in []*foreignkeyv1.StringValueOrRef) pulumi.StringArray {
	out := pulumi.StringArray{}
	for _, ref := range in {
		out = append(out, pulumi.String(ref.GetValue()))
	}
	return out
}

func sortedBrowsers(in []*awsbedrockagentcoretoolsv1alpha1.AwsBedrockAgentCoreBrowser) []*awsbedrockagentcoretoolsv1alpha1.AwsBedrockAgentCoreBrowser {
	out := append([]*awsbedrockagentcoretoolsv1alpha1.AwsBedrockAgentCoreBrowser{}, in...)
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}

func sortedProfiles(in []*awsbedrockagentcoretoolsv1alpha1.AwsBedrockAgentCoreBrowserProfile) []*awsbedrockagentcoretoolsv1alpha1.AwsBedrockAgentCoreBrowserProfile {
	out := append([]*awsbedrockagentcoretoolsv1alpha1.AwsBedrockAgentCoreBrowserProfile{}, in...)
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}

func sortedInterpreters(in []*awsbedrockagentcoretoolsv1alpha1.AwsBedrockAgentCoreCodeInterpreter) []*awsbedrockagentcoretoolsv1alpha1.AwsBedrockAgentCoreCodeInterpreter {
	out := append([]*awsbedrockagentcoretoolsv1alpha1.AwsBedrockAgentCoreCodeInterpreter{}, in...)
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}
