package module

import (
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/globalaccelerator"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// accelerator creates the AWS Global Accelerator resource.
//
// Global Accelerator is a global service homed in us-west-2; the provider
// routes API calls there regardless of the configured region. Creates and
// updates wait for the accelerator to reach the DEPLOYED state, and deletes
// are preceded by an automatic disable (an AWS requirement) — expect minutes
// per operation.
func accelerator(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) (*globalaccelerator.Accelerator, error) {
	spec := locals.GlobalAccelerator.Spec
	name := locals.GlobalAccelerator.Metadata.Name

	args := &globalaccelerator.AcceleratorArgs{
		// The cloud name is set explicitly from metadata.name — never Pulumi
		// auto-naming — so both engines produce the same console identity.
		Name: pulumi.String(name),
		Tags: pulumi.ToStringMap(locals.AwsTags),
	}

	// Presence-honest optional scalars: omit when unset so the provider
	// defaults apply (enabled=true, IPV4).
	if spec.Enabled != nil {
		args.Enabled = pulumi.Bool(spec.GetEnabled())
	}
	if spec.IpAddressType != nil {
		args.IpAddressType = pulumi.StringPtr(spec.GetIpAddressType())
	}

	// BYOIP static addresses (ForceNew). Empty means AWS-allocated anycast IPs.
	if len(spec.IpAddresses) > 0 {
		args.IpAddresses = pulumi.ToStringArray(spec.IpAddresses)
	}

	// The attributes block is ALWAYS materialized with an explicit
	// flow_logs_enabled value. Flow-log settings live on a separate
	// accelerator-attributes API, and an omitted block is diff-suppressed —
	// so dropping the block after flow logs were enabled would silently
	// leave AWS logging forever. Sending the explicit disabled state makes
	// the manifest the single source of truth.
	// The VALUE form of AcceleratorAttributesArgs satisfies the PtrInput
	// interface directly; the AcceleratorAttributesPtr(...) wrapper trips a
	// runtime marshaling assertion in the engine (compiles clean, panics only
	// at deploy), so it is deliberately not used here.
	attributes := globalaccelerator.AcceleratorAttributesArgs{
		FlowLogsEnabled: pulumi.Bool(spec.FlowLogs != nil && spec.FlowLogs.Enabled),
	}
	if spec.FlowLogs != nil && spec.FlowLogs.Enabled {
		attributes.FlowLogsS3Bucket = pulumi.String(spec.FlowLogs.S3Bucket.GetValue())
		if spec.FlowLogs.S3Prefix != "" {
			attributes.FlowLogsS3Prefix = pulumi.String(spec.FlowLogs.S3Prefix)
		}
	}
	args.Attributes = attributes

	return globalaccelerator.NewAccelerator(ctx, name, args, pulumi.Provider(provider))
}
