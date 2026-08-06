package module

import (
	"github.com/pkg/errors"
	awsecsclusterv1alpha1 "github.com/plantonhq/planton/apis/dev/planton/provider/aws/awsecscluster/v1alpha1"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/ecs"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// capacityProviders materializes the cluster's capacity: one
// aws_ecs_capacity_provider per folded EC2 entry (keyed by name, so
// adding or removing an entry never disturbs its siblings) plus ONE
// association that puts the union of Fargate built-ins and EC2 provider
// names onto the cluster together with the default strategy.
//
// The association is a PUT of the whole set -- exactly why it is a single
// resource here: two association resources on one cluster would fight
// each other on every apply.
func capacityProviders(ctx *pulumi.Context, locals *Locals, provider *aws.Provider, createdCluster *ecs.Cluster) ([]*ecs.CapacityProvider, error) {
	spec := locals.AwsEcsCluster.Spec

	createdProviders := make([]*ecs.CapacityProvider, 0, len(spec.Ec2CapacityProviders))

	// The union starts with the built-ins (FARGATE/FARGATE_SPOT are
	// AWS-managed -- never created, only associated) and grows one name
	// per created EC2 provider. Using the created resources' Name
	// outputs gives the association its dependency edges for free.
	associatedNames := pulumi.StringArray{}
	for _, builtin := range spec.CapacityProviders {
		associatedNames = append(associatedNames, pulumi.String(builtin))
	}

	for _, providerSpec := range spec.Ec2CapacityProviders {
		created, err := ec2CapacityProvider(ctx, locals, provider, providerSpec)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create capacity provider %q", providerSpec.Name)
		}
		createdProviders = append(createdProviders, created)
		associatedNames = append(associatedNames, created.Name)
	}

	// Nothing to associate: a bare cluster (services name a launch_type
	// directly) legitimately has no providers at all.
	if len(associatedNames) == 0 {
		return createdProviders, nil
	}

	associationArgs := &ecs.ClusterCapacityProvidersArgs{
		ClusterName:       createdCluster.Name,
		CapacityProviders: associatedNames,
	}

	if len(spec.DefaultCapacityProviderStrategy) > 0 {
		strategies := ecs.ClusterCapacityProvidersDefaultCapacityProviderStrategyArray{}
		for _, strategy := range spec.DefaultCapacityProviderStrategy {
			strategyArgs := ecs.ClusterCapacityProvidersDefaultCapacityProviderStrategyArgs{
				CapacityProvider: pulumi.String(strategy.CapacityProvider),
			}
			if strategy.Base != 0 {
				strategyArgs.Base = pulumi.Int(int(strategy.Base))
			}
			if strategy.Weight != 0 {
				strategyArgs.Weight = pulumi.Int(int(strategy.Weight))
			}
			strategies = append(strategies, strategyArgs)
		}
		associationArgs.DefaultCapacityProviderStrategies = strategies
	}

	if _, err := ecs.NewClusterCapacityProviders(ctx, "capacity-provider-association",
		associationArgs, pulumi.Provider(provider), pulumi.Parent(createdCluster)); err != nil {
		return nil, errors.Wrap(err, "failed to associate capacity providers with the cluster")
	}

	return createdProviders, nil
}

// ec2CapacityProvider materializes one folded EC2 capacity provider. The
// wrapped auto-scaling group and the provider name are create-time
// (ForceNew); the managed scaling/draining/protection knobs update in
// place.
func ec2CapacityProvider(ctx *pulumi.Context, locals *Locals, provider *aws.Provider, providerSpec *awsecsclusterv1alpha1.AwsEcsClusterEc2CapacityProvider) (*ecs.CapacityProvider, error) {
	asgProviderArgs := &ecs.CapacityProviderAutoScalingGroupProviderArgs{
		AutoScalingGroupArn: pulumi.String(providerSpec.AutoScalingGroupArn.GetValue()),
	}

	if scaling := providerSpec.ManagedScaling; scaling != nil {
		scalingArgs := &ecs.CapacityProviderAutoScalingGroupProviderManagedScalingArgs{}
		if scaling.Status != "" {
			scalingArgs.Status = pulumi.String(scaling.Status)
		}
		if scaling.TargetCapacity != 0 {
			scalingArgs.TargetCapacity = pulumi.Int(int(scaling.TargetCapacity))
		}
		if scaling.MinimumScalingStepSize != 0 {
			scalingArgs.MinimumScalingStepSize = pulumi.Int(int(scaling.MinimumScalingStepSize))
		}
		if scaling.MaximumScalingStepSize != 0 {
			scalingArgs.MaximumScalingStepSize = pulumi.Int(int(scaling.MaximumScalingStepSize))
		}
		if scaling.InstanceWarmupPeriodSeconds != 0 {
			scalingArgs.InstanceWarmupPeriod = pulumi.Int(int(scaling.InstanceWarmupPeriodSeconds))
		}
		asgProviderArgs.ManagedScaling = scalingArgs
	}

	// Requires the group's own new-instance scale-in protection when
	// ENABLED -- AWS validates the pairing at create.
	if providerSpec.ManagedTerminationProtection != "" {
		asgProviderArgs.ManagedTerminationProtection = pulumi.String(providerSpec.ManagedTerminationProtection)
	}
	if providerSpec.ManagedDraining != "" {
		asgProviderArgs.ManagedDraining = pulumi.String(providerSpec.ManagedDraining)
	}

	// The Pulumi resource name is the folded entry's name, so each
	// entry keeps its own state identity (adding/removing entries never
	// replaces siblings); the same name is the cloud identity.
	createdProvider, err := ecs.NewCapacityProvider(ctx, providerSpec.Name, &ecs.CapacityProviderArgs{
		Name:                     pulumi.String(providerSpec.Name),
		AutoScalingGroupProvider: asgProviderArgs,
		Tags:                     pulumi.ToStringMap(locals.AwsTags),
	}, pulumi.Provider(provider))
	if err != nil {
		return nil, err
	}

	return createdProvider, nil
}
