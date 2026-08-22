package module

import (
	"github.com/pkg/errors"
	cloudflarezerotrustmcpportalv1alpha1 "github.com/plantonhq/planton/catalog/cloudflare/cloudflarezerotrustmcpportal/v1alpha1"
	"github.com/pulumi/pulumi-cloudflare/sdk/v6/go/cloudflare"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// mcpPortal publishes the MCP portal. The portal's id is user-supplied and
// forces replacement; hostname and name update in place. The servers rows
// are a SET at the provider (the backend ignores declaration order and
// returns its own canonical order), so reordering spec rows never plans a
// change.
func mcpPortal(
	ctx *pulumi.Context,
	locals *Locals,
	cloudflareProvider *cloudflare.Provider,
) error {
	spec := locals.CloudflareZeroTrustMcpPortal.Spec

	args := &cloudflare.ZeroTrustAccessAiControlsMcpPortalArgs{
		AccountId:                            pulumi.String(spec.AccountId),
		ZeroTrustAccessAiControlsMcpPortalId: pulumi.String(spec.PortalId),
		Hostname:                             pulumi.String(spec.Hostname),
		Name:                                 pulumi.String(spec.Name),
	}

	if spec.Description != "" {
		args.Description = pulumi.String(spec.Description)
	}
	if spec.AllowCodeMode != nil {
		args.AllowCodeMode = pulumi.BoolPtr(spec.GetAllowCodeMode())
	}
	if spec.SecureWebGateway != nil {
		args.SecureWebGateway = pulumi.BoolPtr(spec.GetSecureWebGateway())
	}

	if len(spec.Servers) > 0 {
		servers := cloudflare.ZeroTrustAccessAiControlsMcpPortalServerArray{}
		for _, server := range spec.Servers {
			servers = append(servers, buildServerRow(server))
		}
		args.Servers = servers
	}

	createdPortal, err := cloudflare.NewZeroTrustAccessAiControlsMcpPortal(
		ctx,
		"mcp_portal",
		args,
		pulumi.Provider(cloudflareProvider),
	)
	if err != nil {
		return errors.Wrap(err, "failed to create mcp portal")
	}

	ctx.Export(OpPortalId, createdPortal.ZeroTrustAccessAiControlsMcpPortalId)
	ctx.Export(OpHostname, createdPortal.Hostname)

	return nil
}

// buildServerRow maps one published-server row. Omitted booleans are not
// sent -- Cloudflare's defaults keep the server enabled with on-behalf
// authentication.
func buildServerRow(server *cloudflarezerotrustmcpportalv1alpha1.CloudflareZeroTrustMcpPortalServer) cloudflare.ZeroTrustAccessAiControlsMcpPortalServerArgs {
	row := cloudflare.ZeroTrustAccessAiControlsMcpPortalServerArgs{
		ServerId: pulumi.String(server.ServerId.GetValue()),
	}
	if server.DefaultDisabled != nil {
		row.DefaultDisabled = pulumi.BoolPtr(server.GetDefaultDisabled())
	}
	if server.OnBehalf != nil {
		row.OnBehalf = pulumi.BoolPtr(server.GetOnBehalf())
	}

	if len(server.UpdatedPrompts) > 0 {
		prompts := cloudflare.ZeroTrustAccessAiControlsMcpPortalServerUpdatedPromptArray{}
		for _, override := range server.UpdatedPrompts {
			prompt := cloudflare.ZeroTrustAccessAiControlsMcpPortalServerUpdatedPromptArgs{
				Name: pulumi.String(override.Name),
			}
			if override.Alias != "" {
				prompt.Alias = pulumi.String(override.Alias)
			}
			if override.Description != "" {
				prompt.Description = pulumi.String(override.Description)
			}
			if override.Enabled != nil {
				prompt.Enabled = pulumi.BoolPtr(override.GetEnabled())
			}
			prompts = append(prompts, prompt)
		}
		row.UpdatedPrompts = prompts
	}

	if len(server.UpdatedTools) > 0 {
		tools := cloudflare.ZeroTrustAccessAiControlsMcpPortalServerUpdatedToolArray{}
		for _, override := range server.UpdatedTools {
			tool := cloudflare.ZeroTrustAccessAiControlsMcpPortalServerUpdatedToolArgs{
				Name: pulumi.String(override.Name),
			}
			if override.Alias != "" {
				tool.Alias = pulumi.String(override.Alias)
			}
			if override.Description != "" {
				tool.Description = pulumi.String(override.Description)
			}
			if override.Enabled != nil {
				tool.Enabled = pulumi.BoolPtr(override.GetEnabled())
			}
			tools = append(tools, tool)
		}
		row.UpdatedTools = tools
	}

	return row
}
