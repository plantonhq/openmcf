package module

import (
	"github.com/pkg/errors"
	awsgav1 "github.com/plantonhq/planton/apis/dev/planton/provider/aws/awsglobalaccelerator/v1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/aws/pulumiawsprovider"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/globalaccelerator"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Resources is the primary entry point for the AwsGlobalAccelerator Pulumi
// module. It creates the accelerator, listeners, and endpoint groups, then
// exports all outputs for downstream consumption.
func Resources(ctx *pulumi.Context, stackInput *awsgav1.AwsGlobalAcceleratorStackInput) error {
	locals := initializeLocals(ctx, stackInput)

	// Build the AWS provider from the stack input via the shared builder, which resolves
	// the right credential mechanism (static keys, keyless web identity, or ambient chain).
	provider, err := pulumiawsprovider.Get(ctx, stackInput.ProviderConfig, locals.GlobalAccelerator.Spec.Region)
	if err != nil {
		return errors.Wrap(err, "failed to create AWS provider")
	}

	accel, err := accelerator(ctx, locals, provider)
	if err != nil {
		return errors.Wrap(err, "failed to create global accelerator")
	}

	listenerResult, err := listeners(ctx, locals, provider, accel)
	if err != nil {
		return errors.Wrap(err, "failed to create listeners")
	}

	endpointGroupResult, err := endpointGroups(ctx, locals, provider, listenerResult)
	if err != nil {
		return errors.Wrap(err, "failed to create endpoint groups")
	}

	// Export accelerator-level outputs.
	ctx.Export(OpAcceleratorArn, accel.Arn)
	ctx.Export(OpAcceleratorDnsName, accel.DnsName)
	ctx.Export(OpAcceleratorDualStackDns, accel.DualStackDnsName)
	ctx.Export(OpAcceleratorHostedZoneId, accel.HostedZoneId)

	// Flatten the static anycast addresses across all IP sets (IPv4, and IPv6
	// for dual-stack accelerators) — the same flatten the Terraform module
	// performs, so both engines export an identical list. The ApplyT callback
	// takes the SDK's concrete []globalaccelerator.AcceleratorIpSet element
	// type: a mistyped applier compiles but panics at deploy, so keep the
	// signature aligned with the SDK type.
	ctx.Export(OpAcceleratorIpAddresses, accel.IpSets.ApplyT(func(sets []globalaccelerator.AcceleratorIpSet) []string {
		addresses := []string{}
		for _, ipSet := range sets {
			addresses = append(addresses, ipSet.IpAddresses...)
		}
		return addresses
	}))

	// Build and export listener ARN map.
	listenerArnMap := pulumi.StringMap{}
	for name, l := range listenerResult.Listeners {
		listenerArnMap[name] = l.ID().ToStringOutput()
	}
	ctx.Export(OpListenerArns, listenerArnMap)

	// Build and export endpoint group ARN map.
	endpointGroupArnMap := pulumi.StringMap{}
	for key, eg := range endpointGroupResult.EndpointGroups {
		endpointGroupArnMap[key] = eg.Arn
	}
	ctx.Export(OpEndpointGroupArns, endpointGroupArnMap)

	return nil
}
