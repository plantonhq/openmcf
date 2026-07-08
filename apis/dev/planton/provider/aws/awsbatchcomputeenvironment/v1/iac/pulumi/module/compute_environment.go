package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/batch"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// computeEnvironment creates the MANAGED Batch compute environment.
//
// The spec's CEL rules guarantee the EC2/SPOT-only fields are absent for
// Fargate types, so the type switches below exist to build the right
// provider payload (AWS rejects EC2 knobs on Fargate requests), not to
// silently discard user intent.
func computeEnvironment(
	ctx *pulumi.Context,
	locals *Locals,
	provider *aws.Provider,
) (*batch.ComputeEnvironment, error) {
	spec := locals.AwsBatchComputeEnvironment.Spec
	cr := spec.ComputeResources

	computeResources := &batch.ComputeEnvironmentComputeResourcesArgs{
		Type:     pulumi.String(cr.Type),
		MaxVcpus: pulumi.Int(cr.MaxVcpus),
	}

	// Subnets are required for every type: even "serverless" Fargate tasks
	// get ENIs placed into these subnets.
	var subnetIds pulumi.StringArray
	for _, s := range cr.SubnetIds {
		if s.GetValue() != "" {
			subnetIds = append(subnetIds, pulumi.String(s.GetValue()))
		}
	}
	computeResources.Subnets = subnetIds

	var sgIds pulumi.StringArray
	for _, sg := range cr.SecurityGroupIds {
		if sg.GetValue() != "" {
			sgIds = append(sgIds, pulumi.String(sg.GetValue()))
		}
	}
	if len(sgIds) > 0 {
		computeResources.SecurityGroupIds = sgIds
	}

	isEc2Family := cr.Type == "EC2" || cr.Type == "SPOT"

	if isEc2Family {
		// min_vcpus is platform-defaulted to 0, so the getter is always
		// meaningful here; Fargate environments must not send it at all.
		computeResources.MinVcpus = pulumi.IntPtr(int(cr.GetMinVcpus()))
		if cr.DesiredVcpus > 0 {
			computeResources.DesiredVcpus = pulumi.IntPtr(int(cr.DesiredVcpus))
		}
		if len(cr.InstanceTypes) > 0 {
			computeResources.InstanceTypes = pulumi.ToStringArray(cr.InstanceTypes)
		}
		if cr.AllocationStrategy != "" {
			computeResources.AllocationStrategy = pulumi.StringPtr(cr.AllocationStrategy)
		}
		if cr.InstanceRole.GetValue() != "" {
			computeResources.InstanceRole = pulumi.StringPtr(cr.InstanceRole.GetValue())
		}
		if cr.Ec2KeyPair != "" {
			computeResources.Ec2KeyPair = pulumi.StringPtr(cr.Ec2KeyPair)
		}
		if cr.PlacementGroup != "" {
			computeResources.PlacementGroup = pulumi.StringPtr(cr.PlacementGroup)
		}
		// These tags land on the EC2 instances / Spot requests Batch
		// launches -- deliberately NOT merged with the environment's own
		// identity tags (locals.AwsTags), which tag the CE resource itself.
		if len(cr.ResourceTags) > 0 {
			computeResources.Tags = pulumi.ToStringMap(cr.ResourceTags)
		}

		if cr.LaunchTemplate != nil {
			lt := &batch.ComputeEnvironmentComputeResourcesLaunchTemplateArgs{
				LaunchTemplateId: pulumi.StringPtr(cr.LaunchTemplate.LaunchTemplateId.GetValue()),
			}
			if cr.LaunchTemplate.Version != "" {
				lt.Version = pulumi.StringPtr(cr.LaunchTemplate.Version)
			}
			computeResources.LaunchTemplate = lt
		}

		if len(cr.Ec2Configurations) > 0 {
			var ec2Configs batch.ComputeEnvironmentComputeResourcesEc2ConfigurationArray
			for _, ec2Cfg := range cr.Ec2Configurations {
				cfg := &batch.ComputeEnvironmentComputeResourcesEc2ConfigurationArgs{}
				if ec2Cfg.ImageType != "" {
					cfg.ImageType = pulumi.StringPtr(ec2Cfg.ImageType)
				}
				if ec2Cfg.ImageIdOverride != "" {
					cfg.ImageIdOverride = pulumi.StringPtr(ec2Cfg.ImageIdOverride)
				}
				if ec2Cfg.ImageKubernetesVersion != "" {
					cfg.ImageKubernetesVersion = pulumi.StringPtr(ec2Cfg.ImageKubernetesVersion)
				}
				ec2Configs = append(ec2Configs, cfg)
			}
			computeResources.Ec2Configurations = ec2Configs
		}
	}

	if cr.Type == "SPOT" {
		if cr.BidPercentage != nil {
			computeResources.BidPercentage = pulumi.IntPtr(int(cr.GetBidPercentage()))
		}
		if cr.SpotIamFleetRole.GetValue() != "" {
			computeResources.SpotIamFleetRole = pulumi.StringPtr(cr.SpotIamFleetRole.GetValue())
		}
	}

	args := &batch.ComputeEnvironmentArgs{
		// The cloud name comes from metadata.name (the catalog naming
		// basis) -- set explicitly so both engines create the same
		// environment name and Pulumi never auto-names.
		Name: pulumi.StringPtr(locals.AwsBatchComputeEnvironment.Metadata.Name),
		// Only MANAGED environments are modeled: Batch owns the instance
		// lifecycle. (UNMANAGED means bring-your-own ECS container
		// instances -- a different operating model.)
		Type:             pulumi.String("MANAGED"),
		State:            pulumi.StringPtr(spec.GetState()),
		ComputeResources: computeResources,
		Tags:             pulumi.ToStringMap(locals.AwsTags),
	}

	// Leaving service_role unset lets AWS use (and auto-create) the Batch
	// service-linked role -- which is also what keeps the environment
	// eligible for in-place infrastructure updates.
	if spec.ServiceRole.GetValue() != "" {
		args.ServiceRole = pulumi.StringPtr(spec.ServiceRole.GetValue())
	}

	// Batch-on-EKS attachment: create-time only (the provider replaces the
	// environment on any change here).
	if spec.EksConfiguration != nil {
		args.EksConfiguration = &batch.ComputeEnvironmentEksConfigurationArgs{
			EksClusterArn:       pulumi.String(spec.EksConfiguration.EksClusterArn.GetValue()),
			KubernetesNamespace: pulumi.String(spec.EksConfiguration.KubernetesNamespace),
		}
	}

	if spec.UpdatePolicy != nil {
		up := &batch.ComputeEnvironmentUpdatePolicyArgs{
			TerminateJobsOnUpdate: pulumi.Bool(spec.UpdatePolicy.TerminateJobsOnUpdate),
		}
		if spec.UpdatePolicy.JobExecutionTimeoutMinutes != nil {
			up.JobExecutionTimeoutMinutes = pulumi.Int(int(spec.UpdatePolicy.GetJobExecutionTimeoutMinutes()))
		}
		args.UpdatePolicy = up
	}

	ce, err := batch.NewComputeEnvironment(ctx, "compute-environment", args, pulumi.Provider(provider))
	if err != nil {
		return nil, errors.Wrap(err, "create batch compute environment")
	}

	return ce, nil
}
