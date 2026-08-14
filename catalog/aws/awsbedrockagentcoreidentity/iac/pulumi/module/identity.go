package module

import (
	"sort"

	"github.com/pkg/errors"
	awsbedrockagentcoreidentityv1alpha1 "github.com/plantonhq/planton/catalog/aws/awsbedrockagentcoreidentity/v1alpha1"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/bedrock"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// identity creates the AgentCore identity bundle's arms and exports
// outputs.
//
// Lifecycle facts the renders below depend on:
//   - AWS vaults the credential secrets in Secrets Manager under the
//     service's token vault -- consumers reference the provider ARN,
//     never the secret;
//   - the write-only credential argument variants (api_key_wo,
//     client_id_wo, client_secret_wo) are excluded by design: the spec's
//     sensitive fields arrive just-in-time resolved, and the plain
//     arguments let the provider detect rotation;
//   - a Cedar policy is a structural child of its engine (created after,
//     destroyed before).
func identity(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) error {
	spec := locals.Spec

	// Named identities AgentCore workloads present when calling other
	// services. Iteration is name-sorted for deterministic previews.
	workloadIdentityArns := pulumi.StringMap{}
	for _, w := range sortedWorkloadIdentities(spec.WorkloadIdentities) {
		args := &bedrock.AgentcoreWorkloadIdentityArgs{
			Name: pulumi.String(w.Name),
		}
		if len(w.AllowedResourceOauth2ReturnUrls) > 0 {
			args.AllowedResourceOauth2ReturnUrls = pulumi.ToStringArray(w.AllowedResourceOauth2ReturnUrls)
		}
		created, err := bedrock.NewAgentcoreWorkloadIdentity(ctx, "workload-identity-"+w.Name, args, pulumi.Provider(provider))
		if err != nil {
			return errors.Wrapf(err, "create workload identity %q", w.Name)
		}
		workloadIdentityArns[w.Name] = created.WorkloadIdentityArn
	}
	ctx.Export(OpWorkloadIdentityArns, workloadIdentityArns)

	// Vaulted API keys for outbound calls.
	apiKeyProviderArns := pulumi.StringMap{}
	apiKeySecretArns := pulumi.StringMap{}
	for _, p := range sortedApiKeyProviders(spec.ApiKeyCredentialProviders) {
		created, err := bedrock.NewAgentcoreApiKeyCredentialProvider(ctx, "api-key-provider-"+p.Name,
			&bedrock.AgentcoreApiKeyCredentialProviderArgs{
				Name: pulumi.String(p.Name),
				// The spec value is sensitive end to end (JIT-resolved
				// secret reference); the plain argument lets the provider
				// detect rotation.
				ApiKey: pulumi.String(p.ApiKey),
				Tags:   pulumi.ToStringMap(locals.AwsTags),
			}, pulumi.Provider(provider))
		if err != nil {
			return errors.Wrapf(err, "create api key provider %q", p.Name)
		}
		apiKeyProviderArns[p.Name] = created.CredentialProviderArn
		apiKeySecretArns[p.Name] = created.ApiKeySecretArns.Index(pulumi.Int(0)).SecretArn()
	}
	ctx.Export(OpApiKeyProviderArns, apiKeyProviderArns)
	ctx.Export(OpApiKeySecretArns, apiKeySecretArns)

	// Vaulted OAuth2 clients. The spec's vendor field selects which of
	// the provider's six structurally-identical vendor members renders --
	// the vendor field IS the discriminator.
	oauth2ProviderArns := pulumi.StringMap{}
	oauth2ClientSecretArns := pulumi.StringMap{}
	for _, p := range sortedOauth2Providers(spec.Oauth2CredentialProviders) {
		config := &bedrock.AgentcoreOauth2CredentialProviderOauth2ProviderConfigArgs{}
		vendorValue := ""
		switch p.Vendor {
		case "CUSTOM":
			vendorValue = "CustomOauth2"
			custom := &bedrock.AgentcoreOauth2CredentialProviderOauth2ProviderConfigCustomOauth2ProviderConfigArgs{
				ClientId:     pulumi.String(p.ClientId),
				ClientSecret: pulumi.String(p.ClientSecret),
			}
			// Required for CUSTOM vendors (spec-validated): exactly one
			// of a discovery URL or spelled-out endpoints.
			discovery := &bedrock.AgentcoreOauth2CredentialProviderOauth2ProviderConfigCustomOauth2ProviderConfigOauthDiscoveryArgs{}
			if p.OauthDiscovery.DiscoveryUrl != "" {
				discovery.DiscoveryUrl = pulumi.String(p.OauthDiscovery.DiscoveryUrl)
			}
			if p.OauthDiscovery.AuthorizationServerMetadata != nil {
				metadata := p.OauthDiscovery.AuthorizationServerMetadata
				serverMetadata := &bedrock.AgentcoreOauth2CredentialProviderOauth2ProviderConfigCustomOauth2ProviderConfigOauthDiscoveryAuthorizationServerMetadataArgs{
					Issuer:                pulumi.String(metadata.Issuer),
					AuthorizationEndpoint: pulumi.String(metadata.AuthorizationEndpoint),
					TokenEndpoint:         pulumi.String(metadata.TokenEndpoint),
				}
				if len(metadata.ResponseTypes) > 0 {
					serverMetadata.ResponseTypes = pulumi.ToStringArray(metadata.ResponseTypes)
				}
				discovery.AuthorizationServerMetadata = serverMetadata
			}
			custom.OauthDiscovery = discovery
			config.CustomOauth2ProviderConfig = custom
		case "GITHUB":
			vendorValue = "GithubOauth2"
			config.GithubOauth2ProviderConfig = &bedrock.AgentcoreOauth2CredentialProviderOauth2ProviderConfigGithubOauth2ProviderConfigArgs{
				ClientId:     pulumi.String(p.ClientId),
				ClientSecret: pulumi.String(p.ClientSecret),
			}
		case "GOOGLE":
			vendorValue = "GoogleOauth2"
			config.GoogleOauth2ProviderConfig = &bedrock.AgentcoreOauth2CredentialProviderOauth2ProviderConfigGoogleOauth2ProviderConfigArgs{
				ClientId:     pulumi.String(p.ClientId),
				ClientSecret: pulumi.String(p.ClientSecret),
			}
		case "MICROSOFT":
			vendorValue = "MicrosoftOauth2"
			config.MicrosoftOauth2ProviderConfig = &bedrock.AgentcoreOauth2CredentialProviderOauth2ProviderConfigMicrosoftOauth2ProviderConfigArgs{
				ClientId:     pulumi.String(p.ClientId),
				ClientSecret: pulumi.String(p.ClientSecret),
			}
		case "SALESFORCE":
			vendorValue = "SalesforceOauth2"
			config.SalesforceOauth2ProviderConfig = &bedrock.AgentcoreOauth2CredentialProviderOauth2ProviderConfigSalesforceOauth2ProviderConfigArgs{
				ClientId:     pulumi.String(p.ClientId),
				ClientSecret: pulumi.String(p.ClientSecret),
			}
		case "SLACK":
			vendorValue = "SlackOauth2"
			config.SlackOauth2ProviderConfig = &bedrock.AgentcoreOauth2CredentialProviderOauth2ProviderConfigSlackOauth2ProviderConfigArgs{
				ClientId:     pulumi.String(p.ClientId),
				ClientSecret: pulumi.String(p.ClientSecret),
			}
		}
		created, err := bedrock.NewAgentcoreOauth2CredentialProvider(ctx, "oauth2-provider-"+p.Name,
			&bedrock.AgentcoreOauth2CredentialProviderArgs{
				Name:                     pulumi.String(p.Name),
				CredentialProviderVendor: pulumi.String(vendorValue),
				Oauth2ProviderConfig:     config,
				Tags:                     pulumi.ToStringMap(locals.AwsTags),
			}, pulumi.Provider(provider))
		if err != nil {
			return errors.Wrapf(err, "create oauth2 provider %q", p.Name)
		}
		oauth2ProviderArns[p.Name] = created.CredentialProviderArn
		oauth2ClientSecretArns[p.Name] = created.ClientSecretArns.Index(pulumi.Int(0)).SecretArn()
	}
	ctx.Export(OpOauth2ProviderArns, oauth2ProviderArns)
	ctx.Export(OpOauth2ClientSecretArns, oauth2ClientSecretArns)

	// The Cedar authorization engine and its policies.
	policyIds := pulumi.StringMap{}
	if spec.PolicyEngine != nil {
		engineArgs := &bedrock.AgentcorePolicyEngineArgs{
			Name: pulumi.String(spec.PolicyEngine.EngineName),
			Tags: pulumi.ToStringMap(locals.AwsTags),
		}
		if spec.PolicyEngine.Description != "" {
			engineArgs.Description = pulumi.String(spec.PolicyEngine.Description)
		}
		// Changing the key replaces the engine (provider-enforced).
		if spec.PolicyEngine.EncryptionKeyArn.GetValue() != "" {
			engineArgs.EncryptionKeyArn = pulumi.String(spec.PolicyEngine.EncryptionKeyArn.GetValue())
		}
		createdEngine, err := bedrock.NewAgentcorePolicyEngine(ctx, spec.PolicyEngine.EngineName, engineArgs, pulumi.Provider(provider))
		if err != nil {
			return errors.Wrap(err, "create policy engine")
		}
		ctx.Export(OpPolicyEngineId, createdEngine.PolicyEngineId)
		ctx.Export(OpPolicyEngineArn, createdEngine.PolicyEngineArn)

		for _, p := range sortedPolicies(spec.PolicyEngine.Policies) {
			policyArgs := &bedrock.AgentcorePolicyArgs{
				PolicyEngineId: createdEngine.PolicyEngineId,
				Name:           pulumi.String(p.Name),
				Definition: &bedrock.AgentcorePolicyDefinitionArgs{
					Cedar: &bedrock.AgentcorePolicyDefinitionCedarArgs{
						Statement: pulumi.String(p.CedarStatement),
					},
				},
			}
			if p.Description != "" {
				policyArgs.Description = pulumi.String(p.Description)
			}
			if p.ValidationMode != "" {
				policyArgs.ValidationMode = pulumi.String(p.ValidationMode)
			}
			createdPolicy, err := bedrock.NewAgentcorePolicy(ctx, "policy-"+p.Name, policyArgs,
				pulumi.Provider(provider), pulumi.DependsOn([]pulumi.Resource{createdEngine}))
			if err != nil {
				return errors.Wrapf(err, "create policy %q", p.Name)
			}
			policyIds[p.Name] = createdPolicy.PolicyId
		}
	} else {
		ctx.Export(OpPolicyEngineId, pulumi.String(""))
		ctx.Export(OpPolicyEngineArn, pulumi.String(""))
	}
	ctx.Export(OpPolicyIds, policyIds)

	return nil
}

func sortedWorkloadIdentities(in []*awsbedrockagentcoreidentityv1alpha1.AwsBedrockAgentCoreWorkloadIdentity) []*awsbedrockagentcoreidentityv1alpha1.AwsBedrockAgentCoreWorkloadIdentity {
	out := append([]*awsbedrockagentcoreidentityv1alpha1.AwsBedrockAgentCoreWorkloadIdentity{}, in...)
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}

func sortedApiKeyProviders(in []*awsbedrockagentcoreidentityv1alpha1.AwsBedrockAgentCoreApiKeyProvider) []*awsbedrockagentcoreidentityv1alpha1.AwsBedrockAgentCoreApiKeyProvider {
	out := append([]*awsbedrockagentcoreidentityv1alpha1.AwsBedrockAgentCoreApiKeyProvider{}, in...)
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}

func sortedOauth2Providers(in []*awsbedrockagentcoreidentityv1alpha1.AwsBedrockAgentCoreOauth2Provider) []*awsbedrockagentcoreidentityv1alpha1.AwsBedrockAgentCoreOauth2Provider {
	out := append([]*awsbedrockagentcoreidentityv1alpha1.AwsBedrockAgentCoreOauth2Provider{}, in...)
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}

func sortedPolicies(in []*awsbedrockagentcoreidentityv1alpha1.AwsBedrockAgentCoreCedarPolicy) []*awsbedrockagentcoreidentityv1alpha1.AwsBedrockAgentCoreCedarPolicy {
	out := append([]*awsbedrockagentcoreidentityv1alpha1.AwsBedrockAgentCoreCedarPolicy{}, in...)
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}
