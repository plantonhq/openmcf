package module

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/compute"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/projects"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// subnetwork provisions a subnetwork in a custom-mode VPC — the regional
// address space workloads live in: primary IPv4 range for VM interfaces,
// secondary ranges for alias IPs (GKE pods/services), optional IPv6, and
// VPC Flow Logs.
//
// name, project, region, network, and description are immutable (ForceNew in
// the provider): changing any of them destroys and recreates the subnet — an
// outage for everything addressed in it. The primary range is the one
// asymmetric knob: EXPANDING ip_cidr_range updates in place, shrinking
// recreates.
//
// purpose creates the special-role subnets other features depend on:
// REGIONAL_MANAGED_PROXY reserves Envoy address space for the region's
// regional load balancers; PRIVATE_SERVICE_CONNECT backs published PSC
// services.
func subnetwork(ctx *pulumi.Context, locals *Locals, gcpProvider *gcp.Provider) (*compute.Subnetwork, error) {
	spec := locals.GcpSubnetwork.Spec

	// Enable the Compute Engine API so a fresh project can host the subnet.
	// disable_on_destroy stays false (the provider default): tearing down one
	// subnet must never disable the API for everything else in the project.
	serviceArgs := &projects.ServiceArgs{
		DisableDependentServices: pulumi.BoolPtr(true),
		Service:                  pulumi.String("compute.googleapis.com"),
	}
	// Honor the spec contract: an empty project_id falls back to the
	// provider's default project. Leaving Project unset lets the gcp provider
	// resolve its own project (configuration or the GOOGLE_PROJECT /
	// GOOGLE_CLOUD_PROJECT environment chain); an empty string would be sent
	// verbatim and rejected.
	if spec.ProjectId.GetValue() != "" {
		serviceArgs.Project = pulumi.String(spec.ProjectId.GetValue())
	}
	createdProjectService, err := projects.NewService(ctx,
		"subnetwork-compute.googleapis.com", serviceArgs, pulumi.Provider(gcpProvider))
	if err != nil {
		return nil, errors.Wrap(err, "failed to enable compute.googleapis.com api")
	}

	// Drop entries that arrive as empty objects so the API never sees a blank
	// secondary range — identical to the Terraform module's filter.
	var secondaryRanges compute.SubnetworkSecondaryIpRangeArray
	for _, secondaryRange := range spec.SecondaryIpRanges {
		if secondaryRange.RangeName == "" || secondaryRange.IpCidrRange == "" {
			continue
		}
		secondaryRanges = append(secondaryRanges, &compute.SubnetworkSecondaryIpRangeArgs{
			RangeName:   pulumi.String(secondaryRange.RangeName),
			IpCidrRange: pulumi.String(secondaryRange.IpCidrRange),
		})
	}

	args := &compute.SubnetworkArgs{
		Name:    pulumi.String(spec.SubnetworkName),
		Region:  pulumi.String(spec.Region),
		Network: pulumi.String(spec.VpcSelfLink.GetValue()),

		// Planton middleware applies the spec's proto defaults (PRIVATE /
		// IPV4_ONLY) before the module runs; GetPurpose/GetStackType echo them
		// here so direct invocations behave identically.
		Purpose:   pulumi.String(spec.GetPurpose()),
		StackType: pulumi.String(spec.GetStackType()),

		PrivateIpGoogleAccess: pulumi.BoolPtr(spec.PrivateIpGoogleAccess),

		// Safety latch: by default an empty secondary-range list is NOT sent
		// on update, so a partial manifest cannot silently wipe GKE pod
		// ranges.
		SendSecondaryIpRangeIfEmpty: pulumi.BoolPtr(spec.SendSecondaryIpRangeIfEmpty),

		// Deliberate address-space reclaims only: subnet routes still win
		// over the overlapping peer/on-prem routes this permits.
		AllowSubnetCidrRoutesOverlap: pulumi.BoolPtr(spec.AllowSubnetCidrRoutesOverlap),

		SecondaryIpRanges: secondaryRanges,
	}

	// Omitted optionals stay unset (matching the Terraform module's null)
	// rather than being sent as empty strings the API would reject.
	if spec.IpCidrRange != "" {
		args.IpCidrRange = pulumi.String(spec.IpCidrRange)
	}
	if spec.Description != "" {
		args.Description = pulumi.String(spec.Description)
	}
	if spec.Role != "" {
		args.Role = pulumi.String(spec.Role)
	}
	if spec.PrivateIpv6GoogleAccess != "" {
		args.PrivateIpv6GoogleAccess = pulumi.String(spec.PrivateIpv6GoogleAccess)
	}
	// Dual-stack wiring: ipv6_access_type decides whether the assigned prefix
	// is internet-routable (EXTERNAL GUAs) or VPC-internal (INTERNAL ULAs —
	// requires the VPC's ULA range enabled).
	if spec.Ipv6AccessType != "" {
		args.Ipv6AccessType = pulumi.String(spec.Ipv6AccessType)
	}
	if spec.ExternalIpv6Prefix != "" {
		args.ExternalIpv6Prefix = pulumi.String(spec.ExternalIpv6Prefix)
	}
	if spec.ProjectId.GetValue() != "" {
		args.Project = pulumi.String(spec.ProjectId.GetValue())
	}

	// VPC Flow Logs: presence of the message enables logging. Defaults mirror
	// the API's own (5s aggregation, 50% sampling, all metadata) so an empty
	// spec object behaves sanely — identical to the Terraform module.
	if spec.LogConfig != nil {
		logConfig := &compute.SubnetworkLogConfigArgs{
			AggregationInterval: pulumi.String(spec.LogConfig.GetAggregationInterval()),
			FlowSampling:        pulumi.Float64(spec.LogConfig.GetFlowSampling()),
			Metadata:            pulumi.String(spec.LogConfig.GetMetadata()),
		}
		// metadata_fields only accompanies CUSTOM_METADATA (spec-enforced).
		if spec.LogConfig.GetMetadata() == "CUSTOM_METADATA" && len(spec.LogConfig.MetadataFields) > 0 {
			metadataFields := pulumi.StringArray{}
			for _, metadataField := range spec.LogConfig.MetadataFields {
				metadataFields = append(metadataFields, pulumi.String(metadataField))
			}
			logConfig.MetadataFields = metadataFields
		}
		if spec.LogConfig.FilterExpr != "" {
			logConfig.FilterExpr = pulumi.String(spec.LogConfig.FilterExpr)
		}
		args.LogConfig = logConfig
	}

	createdSubnetwork, err := compute.NewSubnetwork(ctx, "subnetwork", args,
		pulumi.Provider(gcpProvider),
		pulumi.DependsOn([]pulumi.Resource{createdProjectService}))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create subnetwork")
	}

	ctx.Export(OpSubnetworkSelfLink, createdSubnetwork.SelfLink)
	ctx.Export(OpSubnetworkName, createdSubnetwork.Name)
	ctx.Export(OpRegion, createdSubnetwork.Region)
	ctx.Export(OpIpCidrRange, createdSubnetwork.IpCidrRange)
	ctx.Export(OpGatewayAddress, createdSubnetwork.GatewayAddress)
	// Exported as a string for a stable cross-engine output shape (the
	// Terraform module does the same with tostring()).
	ctx.Export(OpSubnetworkId, createdSubnetwork.SubnetworkId.ApplyT(func(id int) string {
		return fmt.Sprintf("%d", id)
	}).(pulumi.StringOutput))
	ctx.Export(OpInternalIpv6Prefix, createdSubnetwork.InternalIpv6Prefix)
	ctx.Export(OpExternalIpv6Prefix, createdSubnetwork.ExternalIpv6Prefix)

	// Export secondary ranges with individual fields so the outputs
	// transformer can map them onto the repeated proto message.
	createdSubnetwork.SecondaryIpRanges.ApplyT(func(secondaryIpRanges []compute.SubnetworkSecondaryIpRange) error {
		for index, secondaryIpRange := range secondaryIpRanges {
			ctx.Export(fmt.Sprintf("%s.%d.%s", OpSecondaryRanges, index, "range_name"), pulumi.String(secondaryIpRange.RangeName))
			if secondaryIpRange.IpCidrRange != nil {
				ctx.Export(fmt.Sprintf("%s.%d.%s", OpSecondaryRanges, index, "ip_cidr_range"), pulumi.String(*secondaryIpRange.IpCidrRange))
			}
		}
		return nil
	})

	return createdSubnetwork, nil
}
