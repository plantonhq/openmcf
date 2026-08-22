package module

import (
	"github.com/pkg/errors"
	cloudflareaigatewayv1alpha1 "github.com/plantonhq/planton/catalog/cloudflare/cloudflareaigateway/v1alpha1"
	"github.com/pulumi/pulumi-cloudflare/sdk/v6/go/cloudflare"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// dynamicRoutes creates one ai_gateway_dynamic_routing object per
// dynamic_routes row, attached to the created gateway. The provider forces
// replacement on any change to a route's elements list -- editing a graph
// recreates that route object (requests re-resolve on the next call), while
// renaming a route is its only in-place update. Route resources are keyed by
// route name so a graph edit replaces only its own route.
func dynamicRoutes(
	ctx *pulumi.Context,
	locals *Locals,
	cloudflareProvider *cloudflare.Provider,
	createdGateway *cloudflare.AiGateway,
) error {
	spec := locals.CloudflareAiGateway.Spec

	if len(spec.DynamicRoutes) == 0 {
		// Export the empty map anyway so both engines emit the same output
		// set for the same contract.
		ctx.Export(OpDynamicRouteIds, pulumi.StringMap{})
		return nil
	}

	routeIds := pulumi.StringMap{}
	for _, route := range spec.DynamicRoutes {
		createdRoute, err := cloudflare.NewAiGatewayDynamicRouting(
			ctx,
			"dynamic_route_"+route.Name,
			&cloudflare.AiGatewayDynamicRoutingArgs{
				AccountId: pulumi.String(spec.AccountId),
				// Attaching via the created gateway's id (not the raw spec
				// value) makes the dependency explicit: routes create after
				// the gateway and die with it.
				GatewayId: createdGateway.AiGatewayId,
				Name:      pulumi.String(route.Name),
				Elements:  buildRouteElements(route.Elements),
			},
			pulumi.Provider(cloudflareProvider),
		)
		if err != nil {
			return errors.Wrapf(err, "failed to create dynamic route %s", route.Name)
		}
		routeIds[route.Name] = createdRoute.ID().ToStringOutput()
	}

	ctx.Export(OpDynamicRouteIds, routeIds)

	return nil
}

func buildRouteElements(elements []*cloudflareaigatewayv1alpha1.CloudflareAiGatewayRouteElement) cloudflare.AiGatewayDynamicRoutingElementArray {
	built := cloudflare.AiGatewayDynamicRoutingElementArray{}
	for _, element := range elements {
		elementArgs := cloudflare.AiGatewayDynamicRoutingElementArgs{
			Id:      pulumi.String(element.Id),
			Type:    pulumi.String(element.Type),
			Outputs: buildRouteOutputs(element.Outputs),
		}
		if element.Properties != nil {
			elementArgs.Properties = buildRouteProperties(element.Properties)
		}
		built = append(built, elementArgs)
	}
	return built
}

// buildRouteOutputs maps the spec's edge names to the provider's: on_true
// and on_false are the wire's `true` and `false` (renamed in the spec
// because proto cannot use boolean literals as field names).
func buildRouteOutputs(outputs *cloudflareaigatewayv1alpha1.CloudflareAiGatewayRouteElementOutputs) cloudflare.AiGatewayDynamicRoutingElementOutputsArgs {
	args := cloudflare.AiGatewayDynamicRoutingElementOutputsArgs{}
	if outputs.Next != nil {
		args.Next = cloudflare.AiGatewayDynamicRoutingElementOutputsNextArgs{
			ElementId: pulumi.String(outputs.Next.ElementId),
		}
	}
	if outputs.OnTrue != nil {
		args.True = cloudflare.AiGatewayDynamicRoutingElementOutputsTrueArgs{
			ElementId: pulumi.String(outputs.OnTrue.ElementId),
		}
	}
	if outputs.OnFalse != nil {
		args.False = cloudflare.AiGatewayDynamicRoutingElementOutputsFalseArgs{
			ElementId: pulumi.String(outputs.OnFalse.ElementId),
		}
	}
	if outputs.Success != nil {
		args.Success = cloudflare.AiGatewayDynamicRoutingElementOutputsSuccessArgs{
			ElementId: pulumi.String(outputs.Success.ElementId),
		}
	}
	if outputs.Fallback != nil {
		args.Fallback = cloudflare.AiGatewayDynamicRoutingElementOutputsFallbackArgs{
			ElementId: pulumi.String(outputs.Fallback.ElementId),
		}
	}
	if outputs.ElementId != "" {
		args.ElementId = pulumi.StringPtr(outputs.ElementId)
	}
	return args
}

func buildRouteProperties(properties *cloudflareaigatewayv1alpha1.CloudflareAiGatewayRouteElementProperties) cloudflare.AiGatewayDynamicRoutingElementPropertiesPtrInput {
	args := cloudflare.AiGatewayDynamicRoutingElementPropertiesArgs{}
	if properties.Conditions != "" {
		args.Conditions = pulumi.StringPtr(properties.Conditions)
	}
	if properties.Key != "" {
		args.Key = pulumi.StringPtr(properties.Key)
	}
	if properties.Limit != nil {
		args.Limit = pulumi.Float64Ptr(properties.GetLimit())
	}
	if properties.LimitType != "" {
		args.LimitType = pulumi.StringPtr(properties.LimitType)
	}
	if properties.Window != nil {
		args.Window = pulumi.Float64Ptr(properties.GetWindow())
	}
	if properties.Model != "" {
		args.Model = pulumi.StringPtr(properties.Model)
	}
	// The spec calls this field `provider` (the wire name); the Terraform
	// provider renamed it ai_gateway_dynamic_routing_provider.
	if properties.Provider != "" {
		args.AiGatewayDynamicRoutingProvider = pulumi.StringPtr(properties.Provider)
	}
	if properties.Retries != nil {
		args.Retries = pulumi.Float64Ptr(properties.GetRetries())
	}
	if properties.Timeout != nil {
		args.Timeout = pulumi.Float64Ptr(properties.GetTimeout())
	}
	return args
}
