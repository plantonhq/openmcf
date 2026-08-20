package module

import (
	"github.com/pkg/errors"
	cloudflarezerotrustmcpserverv1alpha1 "github.com/plantonhq/planton/catalog/cloudflare/cloudflarezerotrustmcpserver/v1alpha1"
	"github.com/pulumi/pulumi-cloudflare/sdk/v6/go/cloudflare"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// mcpServer registers the MCP server with Access AI Controls. Identity is
// user-supplied: server_id, hostname, and auth_type all force replacement.
// The two credential fields (auth_credentials, client_secret) are WRITE-ONLY
// at Cloudflare -- stored encrypted, never returned by any read -- so
// out-of-band rotation is invisible to IaC; rotate by changing the value
// here. Both are marked secret on the Pulumi side so they never appear in
// plaintext state.
func mcpServer(
	ctx *pulumi.Context,
	locals *Locals,
	cloudflareProvider *cloudflare.Provider,
) error {
	spec := locals.CloudflareZeroTrustMcpServer.Spec

	args := &cloudflare.ZeroTrustAccessAiControlsMcpServerArgs{
		AccountId:                            pulumi.String(spec.AccountId),
		ZeroTrustAccessAiControlsMcpServerId: pulumi.String(spec.ServerId),
		Name:                                 pulumi.String(spec.Name),
		Hostname:                             pulumi.String(spec.Hostname),
		AuthType:                             pulumi.String(spec.AuthType),
	}

	if spec.AuthCredentials.GetValue() != "" {
		args.AuthCredentials = pulumi.ToSecret(pulumi.String(spec.AuthCredentials.GetValue())).(pulumi.StringOutput)
	}
	if spec.ClientSecret.GetValue() != "" {
		args.ClientSecret = pulumi.ToSecret(pulumi.String(spec.ClientSecret.GetValue())).(pulumi.StringOutput)
	}
	if spec.Description != "" {
		args.Description = pulumi.String(spec.Description)
	}
	if spec.IsSharedOauthCallbackEnabled != nil {
		args.IsSharedOauthCallbackEnabled = pulumi.BoolPtr(spec.GetIsSharedOauthCallbackEnabled())
	}
	if spec.SecureWebGateway != nil {
		args.SecureWebGateway = pulumi.BoolPtr(spec.GetSecureWebGateway())
	}

	if len(spec.UpdatedPrompts) > 0 {
		prompts := cloudflare.ZeroTrustAccessAiControlsMcpServerUpdatedPromptArray{}
		for _, override := range spec.UpdatedPrompts {
			prompts = append(prompts, buildPromptOverride(override))
		}
		args.UpdatedPrompts = prompts
	}

	if len(spec.UpdatedTools) > 0 {
		tools := cloudflare.ZeroTrustAccessAiControlsMcpServerUpdatedToolArray{}
		for _, override := range spec.UpdatedTools {
			tools = append(tools, buildToolOverride(override))
		}
		args.UpdatedTools = tools
	}

	createdServer, err := cloudflare.NewZeroTrustAccessAiControlsMcpServer(
		ctx,
		"mcp_server",
		args,
		pulumi.Provider(cloudflareProvider),
	)
	if err != nil {
		return errors.Wrap(err, "failed to create mcp server")
	}

	ctx.Export(OpServerId, createdServer.ZeroTrustAccessAiControlsMcpServerId)

	return nil
}

// buildPromptOverride maps one prompt override row. An omitted enabled is
// not sent -- Cloudflare's default keeps the prompt available.
func buildPromptOverride(override *cloudflarezerotrustmcpserverv1alpha1.CloudflareZeroTrustMcpServerItemOverride) cloudflare.ZeroTrustAccessAiControlsMcpServerUpdatedPromptArgs {
	row := cloudflare.ZeroTrustAccessAiControlsMcpServerUpdatedPromptArgs{
		Name: pulumi.String(override.Name),
	}
	if override.Alias != "" {
		row.Alias = pulumi.String(override.Alias)
	}
	if override.Description != "" {
		row.Description = pulumi.String(override.Description)
	}
	if override.Enabled != nil {
		row.Enabled = pulumi.BoolPtr(override.GetEnabled())
	}
	return row
}

// buildToolOverride maps one tool override row (same shape as prompts; the
// SDK types differ).
func buildToolOverride(override *cloudflarezerotrustmcpserverv1alpha1.CloudflareZeroTrustMcpServerItemOverride) cloudflare.ZeroTrustAccessAiControlsMcpServerUpdatedToolArgs {
	row := cloudflare.ZeroTrustAccessAiControlsMcpServerUpdatedToolArgs{
		Name: pulumi.String(override.Name),
	}
	if override.Alias != "" {
		row.Alias = pulumi.String(override.Alias)
	}
	if override.Description != "" {
		row.Description = pulumi.String(override.Description)
	}
	if override.Enabled != nil {
		row.Enabled = pulumi.BoolPtr(override.GetEnabled())
	}
	return row
}
