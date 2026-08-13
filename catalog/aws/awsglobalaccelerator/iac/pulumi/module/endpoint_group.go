package module

import (
	"fmt"

	awsgav1 "github.com/plantonhq/planton/catalog/aws/awsglobalaccelerator/v1alpha1"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/globalaccelerator"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// EndpointGroupResult holds the created endpoint groups keyed by their
// composite key ("listener_name/group_name").
type EndpointGroupResult struct {
	EndpointGroups map[string]*globalaccelerator.EndpointGroup
}

// endpointGroups creates endpoint groups for all listeners, iterating over
// the nested spec structure. Each endpoint group is keyed by
// "listener_name/group_name" for output map construction.
func endpointGroups(
	ctx *pulumi.Context,
	locals *Locals,
	provider *aws.Provider,
	listenerResult *ListenerResult,
) (*EndpointGroupResult, error) {
	result := &EndpointGroupResult{
		EndpointGroups: make(map[string]*globalaccelerator.EndpointGroup),
	}

	for _, listenerSpec := range locals.GlobalAccelerator.Spec.Listeners {
		listener, ok := listenerResult.Listeners[listenerSpec.Name]
		if !ok {
			continue
		}

		for _, groupSpec := range listenerSpec.EndpointGroups {
			compositeKey := fmt.Sprintf("%s/%s", listenerSpec.Name, groupSpec.Name)
			resourceName := fmt.Sprintf("%s-%s", listenerSpec.Name, groupSpec.Name)

			group, err := createEndpointGroup(ctx, provider, listener, groupSpec, resourceName)
			if err != nil {
				return nil, fmt.Errorf("failed to create endpoint group %s: %w", compositeKey, err)
			}
			result.EndpointGroups[compositeKey] = group
		}
	}

	return result, nil
}

// createEndpointGroup creates a single Global Accelerator endpoint group resource.
func createEndpointGroup(
	ctx *pulumi.Context,
	provider *aws.Provider,
	listener *globalaccelerator.Listener,
	spec *awsgav1.AwsGlobalAcceleratorEndpointGroup,
	resourceName string,
) (*globalaccelerator.EndpointGroup, error) {
	args := &globalaccelerator.EndpointGroupArgs{
		ListenerArn: listener.ID().ToStringOutput(),
	}

	// Empty string means "inherit the provider region" (the spec's region).
	if spec.EndpointGroupRegion != "" {
		args.EndpointGroupRegion = pulumi.StringPtr(spec.EndpointGroupRegion)
	}

	// Presence-honest health-check dials: omit when unset so the provider/AWS
	// defaults apply (listener port, TCP, interval 30, threshold 3).
	if spec.HealthCheckPort != nil {
		args.HealthCheckPort = pulumi.IntPtr(int(spec.GetHealthCheckPort()))
	}
	if spec.HealthCheckProtocol != nil {
		args.HealthCheckProtocol = pulumi.StringPtr(spec.GetHealthCheckProtocol())
	}
	// health_check_path is meaningful only for HTTP/HTTPS checks (CEL-enforced
	// in the spec).
	if spec.HealthCheckPath != "" {
		args.HealthCheckPath = pulumi.StringPtr(spec.HealthCheckPath)
	}
	if spec.HealthCheckIntervalSeconds != nil {
		args.HealthCheckIntervalSeconds = pulumi.IntPtr(int(spec.GetHealthCheckIntervalSeconds()))
	}
	if spec.ThresholdCount != nil {
		args.ThresholdCount = pulumi.IntPtr(int(spec.GetThresholdCount()))
	}

	// Omitted means 100% (the AWS default). An explicit 0 is a real value —
	// it drains the region while keeping its endpoints registered.
	if spec.TrafficDialPercentage != nil {
		args.TrafficDialPercentage = pulumi.Float64Ptr(spec.GetTrafficDialPercentage())
	}

	if len(spec.Endpoints) > 0 {
		endpointConfigs := make(globalaccelerator.EndpointGroupEndpointConfigurationArray, len(spec.Endpoints))
		for i, ep := range spec.Endpoints {
			epArgs := &globalaccelerator.EndpointGroupEndpointConfigurationArgs{
				EndpointId: pulumi.StringPtr(ep.EndpointId.GetValue()),
			}
			// AWS's documented default weight is 128, but the provider has no
			// default of its own and materializes an omitted weight as 0 —
			// the "route no traffic" value, silently draining the endpoint.
			// The module therefore materializes 128 itself; an explicit 0
			// still stops routing to the endpoint without removing it.
			weight := 128
			if ep.Weight != nil {
				weight = int(ep.GetWeight())
			}
			epArgs.Weight = pulumi.IntPtr(weight)
			// Tri-state: omitted lets AWS apply its per-endpoint-type default;
			// an explicit value pins it. Only meaningful for ALB and EC2
			// endpoints.
			if ep.ClientIpPreservationEnabled != nil {
				epArgs.ClientIpPreservationEnabled = pulumi.BoolPtr(ep.GetClientIpPreservationEnabled())
			}
			// Cross-account endpoints authorize through a Global Accelerator
			// cross-account attachment created in the endpoint-owning account.
			if ep.AttachmentArn != "" {
				epArgs.AttachmentArn = pulumi.StringPtr(ep.AttachmentArn)
			}
			endpointConfigs[i] = epArgs
		}
		args.EndpointConfigurations = endpointConfigs
	}

	if len(spec.PortOverrides) > 0 {
		overrides := make(globalaccelerator.EndpointGroupPortOverrideArray, len(spec.PortOverrides))
		for i, po := range spec.PortOverrides {
			overrides[i] = &globalaccelerator.EndpointGroupPortOverrideArgs{
				ListenerPort: pulumi.Int(int(po.ListenerPort)),
				EndpointPort: pulumi.Int(int(po.EndpointPort)),
			}
		}
		args.PortOverrides = overrides
	}

	return globalaccelerator.NewEndpointGroup(ctx, resourceName, args, pulumi.Provider(provider))
}
