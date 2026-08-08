package module

import (
	"github.com/pkg/errors"
	awsecsclusterv1alpha1 "github.com/plantonhq/planton/catalog/aws/awsecscluster/v1alpha1"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/ecs"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// capacityProviders materializes the cluster's capacity: one
// aws_ecs_capacity_provider per folded EC2 entry and per managed-instances
// entry (each keyed by name, so adding or removing an entry never disturbs
// its siblings) plus ONE association that puts the union of Fargate
// built-ins and folded provider names onto the cluster together with the
// default strategy.
//
// The association is a PUT of the whole set -- exactly why it is a single
// resource here: two association resources on one cluster would fight
// each other on every apply.
func capacityProviders(ctx *pulumi.Context, locals *Locals, provider *aws.Provider, createdCluster *ecs.Cluster) ([]*ecs.CapacityProvider, error) {
	spec := locals.AwsEcsCluster.Spec

	createdProviders := make([]*ecs.CapacityProvider, 0, len(spec.Ec2CapacityProviders)+len(spec.ManagedInstancesCapacityProviders))

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

	for _, providerSpec := range spec.ManagedInstancesCapacityProviders {
		created, err := managedInstancesCapacityProvider(ctx, locals, provider, createdCluster, providerSpec)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create managed-instances capacity provider %q", providerSpec.Name)
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

// managedInstancesCapacityProvider materializes one folded ECS Managed
// Instances capacity provider: ECS launches and retires the EC2 instances
// itself, so unlike the ASG-backed providers there is no group to wrap --
// the provider is bound to its cluster at create (the AWS API requires the
// pairing) and carries the whole launch template inline. Creating one
// launches nothing; instances appear only when a service's strategy
// schedules tasks onto it.
func managedInstancesCapacityProvider(ctx *pulumi.Context, locals *Locals, provider *aws.Provider, createdCluster *ecs.Cluster, providerSpec *awsecsclusterv1alpha1.AwsEcsClusterManagedInstancesCapacityProvider) (*ecs.CapacityProvider, error) {
	launchTemplate := providerSpec.InstanceLaunchTemplate

	launchTemplateArgs := ecs.CapacityProviderManagedInstancesProviderInstanceLaunchTemplateArgs{
		Ec2InstanceProfileArn: pulumi.String(launchTemplate.Ec2InstanceProfileArn.GetValue()),
		NetworkConfiguration:  managedInstancesNetworkConfiguration(launchTemplate.NetworkConfiguration),
	}

	// Changing the purchase model replaces the whole capacity provider
	// (ForceNew); the rest of the launch template updates in place.
	if launchTemplate.CapacityOptionType != "" {
		launchTemplateArgs.CapacityOptionType = pulumi.String(launchTemplate.CapacityOptionType)
	}

	if reservations := launchTemplate.CapacityReservations; reservations != nil {
		reservationArgs := ecs.CapacityProviderManagedInstancesProviderInstanceLaunchTemplateCapacityReservationsArgs{}
		if reservations.ReservationPreference != "" {
			reservationArgs.ReservationPreference = pulumi.String(reservations.ReservationPreference)
		}
		if reservations.ReservationGroupArn != "" {
			reservationArgs.ReservationGroupArn = pulumi.String(reservations.ReservationGroupArn)
		}
		launchTemplateArgs.CapacityReservations = reservationArgs.ToCapacityProviderManagedInstancesProviderInstanceLaunchTemplateCapacityReservationsPtrOutput()
	}

	if requirements := launchTemplate.InstanceRequirements; requirements != nil {
		launchTemplateArgs.InstanceRequirements = managedInstancesRequirements(requirements)
	}

	if launchTemplate.UseLocalStorage != nil {
		launchTemplateArgs.LocalStorageConfiguration = ecs.CapacityProviderManagedInstancesProviderInstanceLaunchTemplateLocalStorageConfigurationArgs{
			UseLocalStorage: pulumi.BoolPtr(*launchTemplate.UseLocalStorage),
		}.ToCapacityProviderManagedInstancesProviderInstanceLaunchTemplateLocalStorageConfigurationPtrOutput()
	}

	if launchTemplate.Monitoring != "" {
		launchTemplateArgs.Monitoring = pulumi.String(launchTemplate.Monitoring)
	}

	if launchTemplate.StorageSizeGib != 0 {
		launchTemplateArgs.StorageConfiguration = ecs.CapacityProviderManagedInstancesProviderInstanceLaunchTemplateStorageConfigurationArgs{
			StorageSizeGib: pulumi.Int(int(launchTemplate.StorageSizeGib)),
		}.ToCapacityProviderManagedInstancesProviderInstanceLaunchTemplateStorageConfigurationPtrOutput()
	}

	managedInstancesArgs := &ecs.CapacityProviderManagedInstancesProviderArgs{
		InfrastructureRoleArn:  pulumi.String(providerSpec.InfrastructureRoleArn.GetValue()),
		InstanceLaunchTemplate: launchTemplateArgs,
	}

	// -1 disables scale-in entirely; nil keeps AWS's default optimization
	// -- distinct values, so the unset sentinel is nil, never a number.
	if providerSpec.ScaleInAfterSeconds != nil {
		managedInstancesArgs.InfrastructureOptimization = ecs.CapacityProviderManagedInstancesProviderInfrastructureOptimizationArgs{
			ScaleInAfter: pulumi.IntPtr(int(*providerSpec.ScaleInAfterSeconds)),
		}.ToCapacityProviderManagedInstancesProviderInfrastructureOptimizationPtrOutput()
	}

	if providerSpec.PropagateTags != "" {
		managedInstancesArgs.PropagateTags = pulumi.String(providerSpec.PropagateTags)
	}

	createdProvider, err := ecs.NewCapacityProvider(ctx, providerSpec.Name, &ecs.CapacityProviderArgs{
		Name: pulumi.String(providerSpec.Name),
		// MI providers are cluster-bound at create -- the AWS API
		// requires the pairing (unlike ASG-backed providers, which
		// forbid it).
		Cluster:                  createdCluster.Name,
		ManagedInstancesProvider: managedInstancesArgs,
		Tags:                     pulumi.ToStringMap(locals.AwsTags),
	}, pulumi.Provider(provider))
	if err != nil {
		return nil, err
	}

	return createdProvider, nil
}

// managedInstancesNetworkConfiguration places the managed instances on the
// network: subnets (required) and the security groups applied to each
// instance.
func managedInstancesNetworkConfiguration(networkSpec *awsecsclusterv1alpha1.AwsEcsClusterManagedInstancesNetworkConfiguration) ecs.CapacityProviderManagedInstancesProviderInstanceLaunchTemplateNetworkConfigurationArgs {
	subnets := pulumi.StringArray{}
	for _, subnet := range networkSpec.Subnets {
		subnets = append(subnets, pulumi.String(subnet.GetValue()))
	}
	networkArgs := ecs.CapacityProviderManagedInstancesProviderInstanceLaunchTemplateNetworkConfigurationArgs{
		Subnets: subnets,
	}
	if len(networkSpec.SecurityGroups) > 0 {
		securityGroups := pulumi.StringArray{}
		for _, group := range networkSpec.SecurityGroups {
			securityGroups = append(securityGroups, pulumi.String(group.GetValue()))
		}
		networkArgs.SecurityGroups = securityGroups
	}
	return networkArgs
}

// managedInstancesRequirements translates the attribute-based instance
// requirements. The two required dimensions always render; every other
// field renders only when set (0 means "unset" on ranges -- ranges here
// never legitimately start at zero).
func managedInstancesRequirements(requirements *awsecsclusterv1alpha1.AwsEcsClusterManagedInstancesRequirements) ecs.CapacityProviderManagedInstancesProviderInstanceLaunchTemplateInstanceRequirementsPtrInput {
	memoryArgs := ecs.CapacityProviderManagedInstancesProviderInstanceLaunchTemplateInstanceRequirementsMemoryMibArgs{
		Min: pulumi.Int(int(requirements.MemoryMib.Min)),
	}
	if requirements.MemoryMib.Max != 0 {
		memoryArgs.Max = pulumi.IntPtr(int(requirements.MemoryMib.Max))
	}
	vcpuArgs := ecs.CapacityProviderManagedInstancesProviderInstanceLaunchTemplateInstanceRequirementsVcpuCountArgs{
		Min: pulumi.Int(int(requirements.VcpuCount.Min)),
	}
	if requirements.VcpuCount.Max != 0 {
		vcpuArgs.Max = pulumi.IntPtr(int(requirements.VcpuCount.Max))
	}

	requirementArgs := ecs.CapacityProviderManagedInstancesProviderInstanceLaunchTemplateInstanceRequirementsArgs{
		MemoryMib: memoryArgs,
		VcpuCount: vcpuArgs,
	}

	if len(requirements.AllowedInstanceTypes) > 0 {
		requirementArgs.AllowedInstanceTypes = pulumi.ToStringArray(requirements.AllowedInstanceTypes)
	}
	if len(requirements.ExcludedInstanceTypes) > 0 {
		requirementArgs.ExcludedInstanceTypes = pulumi.ToStringArray(requirements.ExcludedInstanceTypes)
	}
	if len(requirements.InstanceGenerations) > 0 {
		requirementArgs.InstanceGenerations = pulumi.ToStringArray(requirements.InstanceGenerations)
	}
	if len(requirements.CpuManufacturers) > 0 {
		requirementArgs.CpuManufacturers = pulumi.ToStringArray(requirements.CpuManufacturers)
	}
	if requirements.BareMetal != "" {
		requirementArgs.BareMetal = pulumi.String(requirements.BareMetal)
	}
	if requirements.BurstablePerformance != "" {
		requirementArgs.BurstablePerformance = pulumi.String(requirements.BurstablePerformance)
	}
	if requirements.RequireHibernateSupport {
		requirementArgs.RequireHibernateSupport = pulumi.Bool(true)
	}
	if requirements.SpotMaxPricePercentageOverLowestPrice != 0 {
		requirementArgs.SpotMaxPricePercentageOverLowestPrice = pulumi.IntPtr(int(requirements.SpotMaxPricePercentageOverLowestPrice))
	}
	if requirements.MaxSpotPriceAsPercentageOfOptimalOnDemandPrice != 0 {
		requirementArgs.MaxSpotPriceAsPercentageOfOptimalOnDemandPrice = pulumi.IntPtr(int(requirements.MaxSpotPriceAsPercentageOfOptimalOnDemandPrice))
	}
	if requirements.OnDemandMaxPricePercentageOverLowestPrice != 0 {
		requirementArgs.OnDemandMaxPricePercentageOverLowestPrice = pulumi.IntPtr(int(requirements.OnDemandMaxPricePercentageOverLowestPrice))
	}
	if requirements.LocalStorage != "" {
		requirementArgs.LocalStorage = pulumi.String(requirements.LocalStorage)
	}
	if len(requirements.LocalStorageTypes) > 0 {
		requirementArgs.LocalStorageTypes = pulumi.ToStringArray(requirements.LocalStorageTypes)
	}
	if r := requirements.TotalLocalStorageGb; r != nil && (r.Min != 0 || r.Max != 0) {
		rangeArgs := ecs.CapacityProviderManagedInstancesProviderInstanceLaunchTemplateInstanceRequirementsTotalLocalStorageGbArgs{}
		if r.Min != 0 {
			rangeArgs.Min = pulumi.Float64Ptr(r.Min)
		}
		if r.Max != 0 {
			rangeArgs.Max = pulumi.Float64Ptr(r.Max)
		}
		requirementArgs.TotalLocalStorageGb = rangeArgs.ToCapacityProviderManagedInstancesProviderInstanceLaunchTemplateInstanceRequirementsTotalLocalStorageGbPtrOutput()
	}
	if r := requirements.MemoryGibPerVcpu; r != nil && (r.Min != 0 || r.Max != 0) {
		rangeArgs := ecs.CapacityProviderManagedInstancesProviderInstanceLaunchTemplateInstanceRequirementsMemoryGibPerVcpuArgs{}
		if r.Min != 0 {
			rangeArgs.Min = pulumi.Float64Ptr(r.Min)
		}
		if r.Max != 0 {
			rangeArgs.Max = pulumi.Float64Ptr(r.Max)
		}
		requirementArgs.MemoryGibPerVcpu = rangeArgs.ToCapacityProviderManagedInstancesProviderInstanceLaunchTemplateInstanceRequirementsMemoryGibPerVcpuPtrOutput()
	}
	if r := requirements.NetworkInterfaceCount; r != nil && (r.Min != 0 || r.Max != 0) {
		rangeArgs := ecs.CapacityProviderManagedInstancesProviderInstanceLaunchTemplateInstanceRequirementsNetworkInterfaceCountArgs{}
		if r.Min != 0 {
			rangeArgs.Min = pulumi.IntPtr(int(r.Min))
		}
		if r.Max != 0 {
			rangeArgs.Max = pulumi.IntPtr(int(r.Max))
		}
		requirementArgs.NetworkInterfaceCount = rangeArgs.ToCapacityProviderManagedInstancesProviderInstanceLaunchTemplateInstanceRequirementsNetworkInterfaceCountPtrOutput()
	}
	if r := requirements.NetworkBandwidthGbps; r != nil && (r.Min != 0 || r.Max != 0) {
		rangeArgs := ecs.CapacityProviderManagedInstancesProviderInstanceLaunchTemplateInstanceRequirementsNetworkBandwidthGbpsArgs{}
		if r.Min != 0 {
			rangeArgs.Min = pulumi.Float64Ptr(r.Min)
		}
		if r.Max != 0 {
			rangeArgs.Max = pulumi.Float64Ptr(r.Max)
		}
		requirementArgs.NetworkBandwidthGbps = rangeArgs.ToCapacityProviderManagedInstancesProviderInstanceLaunchTemplateInstanceRequirementsNetworkBandwidthGbpsPtrOutput()
	}
	if r := requirements.BaselineEbsBandwidthMbps; r != nil && (r.Min != 0 || r.Max != 0) {
		rangeArgs := ecs.CapacityProviderManagedInstancesProviderInstanceLaunchTemplateInstanceRequirementsBaselineEbsBandwidthMbpsArgs{}
		if r.Min != 0 {
			rangeArgs.Min = pulumi.IntPtr(int(r.Min))
		}
		if r.Max != 0 {
			rangeArgs.Max = pulumi.IntPtr(int(r.Max))
		}
		requirementArgs.BaselineEbsBandwidthMbps = rangeArgs.ToCapacityProviderManagedInstancesProviderInstanceLaunchTemplateInstanceRequirementsBaselineEbsBandwidthMbpsPtrOutput()
	}
	if r := requirements.AcceleratorCount; r != nil && (r.Min != 0 || r.Max != 0) {
		rangeArgs := ecs.CapacityProviderManagedInstancesProviderInstanceLaunchTemplateInstanceRequirementsAcceleratorCountArgs{}
		if r.Min != 0 {
			rangeArgs.Min = pulumi.IntPtr(int(r.Min))
		}
		if r.Max != 0 {
			rangeArgs.Max = pulumi.IntPtr(int(r.Max))
		}
		requirementArgs.AcceleratorCount = rangeArgs.ToCapacityProviderManagedInstancesProviderInstanceLaunchTemplateInstanceRequirementsAcceleratorCountPtrOutput()
	}
	if len(requirements.AcceleratorManufacturers) > 0 {
		requirementArgs.AcceleratorManufacturers = pulumi.ToStringArray(requirements.AcceleratorManufacturers)
	}
	if len(requirements.AcceleratorNames) > 0 {
		requirementArgs.AcceleratorNames = pulumi.ToStringArray(requirements.AcceleratorNames)
	}
	if len(requirements.AcceleratorTypes) > 0 {
		requirementArgs.AcceleratorTypes = pulumi.ToStringArray(requirements.AcceleratorTypes)
	}
	if r := requirements.AcceleratorTotalMemoryMib; r != nil && (r.Min != 0 || r.Max != 0) {
		rangeArgs := ecs.CapacityProviderManagedInstancesProviderInstanceLaunchTemplateInstanceRequirementsAcceleratorTotalMemoryMibArgs{}
		if r.Min != 0 {
			rangeArgs.Min = pulumi.IntPtr(int(r.Min))
		}
		if r.Max != 0 {
			rangeArgs.Max = pulumi.IntPtr(int(r.Max))
		}
		requirementArgs.AcceleratorTotalMemoryMib = rangeArgs.ToCapacityProviderManagedInstancesProviderInstanceLaunchTemplateInstanceRequirementsAcceleratorTotalMemoryMibPtrOutput()
	}

	return requirementArgs.ToCapacityProviderManagedInstancesProviderInstanceLaunchTemplateInstanceRequirementsPtrOutput()
}
