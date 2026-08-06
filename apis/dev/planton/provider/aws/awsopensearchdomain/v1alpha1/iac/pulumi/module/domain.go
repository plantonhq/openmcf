package module

import (
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/opensearch"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// domain creates the OpenSearch Service domain and exports outputs.
//
// Create-time (ForceNew) surfaces: the domain name, the at-rest KMS key,
// adding/removing vpc_options, and disabling either encryption toggle or FGAC
// once enabled. Everything else -- topology, storage, endpoint policy,
// Auto-Tune, log publishing -- updates in place, though topology changes
// usually run as a blue/green deployment behind the endpoint
// (deployment_strategy tunes how that migration provisions capacity).
func domain(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) error {
	spec := locals.Spec

	args := &opensearch.DomainArgs{
		DomainName:    pulumi.String(locals.DomainName),
		EngineVersion: pulumi.String(spec.EngineVersion),
		Tags:          pulumi.ToStringMap(locals.AwsTags),
	}

	// -------------------------------------------------------------------
	// Cluster configuration
	// -------------------------------------------------------------------

	if cc := spec.ClusterConfig; cc != nil {
		clusterCfg := &opensearch.DomainClusterConfigArgs{
			InstanceType: pulumi.String(cc.InstanceType),
		}

		// Instance count (optional proto field — uses *int32)
		if cc.InstanceCount != nil {
			clusterCfg.InstanceCount = pulumi.Int(int(*cc.InstanceCount))
		}

		// Dedicated master nodes
		if cc.DedicatedMasterEnabled {
			clusterCfg.DedicatedMasterEnabled = pulumi.Bool(true)
			if cc.DedicatedMasterType != "" {
				clusterCfg.DedicatedMasterType = pulumi.String(cc.DedicatedMasterType)
			}
			if cc.DedicatedMasterCount > 0 {
				clusterCfg.DedicatedMasterCount = pulumi.Int(int(cc.DedicatedMasterCount))
			}
		}

		// Coordinator node pools: request routing, query fan-out, and response
		// aggregation offloaded from the data nodes.
		if len(cc.NodeOptions) > 0 {
			var nodeOptions opensearch.DomainClusterConfigNodeOptionArray
			for _, pool := range cc.NodeOptions {
				nodeConfig := &opensearch.DomainClusterConfigNodeOptionNodeConfigArgs{
					Enabled: pulumi.Bool(pool.Enabled),
				}
				if pool.InstanceType != "" {
					nodeConfig.Type = pulumi.String(pool.InstanceType)
				}
				if pool.Count > 0 {
					nodeConfig.Count = pulumi.Int(int(pool.Count))
				}
				nodeOptions = append(nodeOptions, &opensearch.DomainClusterConfigNodeOptionArgs{
					NodeType:   pulumi.String(pool.NodeType),
					NodeConfig: nodeConfig,
				})
			}
			clusterCfg.NodeOptions = nodeOptions
		}

		// Zone awareness
		if cc.ZoneAwarenessEnabled {
			clusterCfg.ZoneAwarenessEnabled = pulumi.Bool(true)
			if cc.AvailabilityZoneCount > 0 {
				clusterCfg.ZoneAwarenessConfig = &opensearch.DomainClusterConfigZoneAwarenessConfigArgs{
					AvailabilityZoneCount: pulumi.Int(int(cc.AvailabilityZoneCount)),
				}
			}
		}

		// UltraWarm storage tier
		if cc.WarmEnabled {
			clusterCfg.WarmEnabled = pulumi.Bool(true)
			if cc.WarmType != "" {
				clusterCfg.WarmType = pulumi.String(cc.WarmType)
			}
			if cc.WarmCount > 0 {
				clusterCfg.WarmCount = pulumi.Int(int(cc.WarmCount))
			}
		}

		// Cold storage (requires UltraWarm)
		if cc.ColdStorageEnabled {
			clusterCfg.ColdStorageOptions = &opensearch.DomainClusterConfigColdStorageOptionsArgs{
				Enabled: pulumi.Bool(true),
			}
		}

		// Multi-AZ with Standby
		if cc.MultiAzWithStandbyEnabled {
			clusterCfg.MultiAzWithStandbyEnabled = pulumi.Bool(true)
		}

		args.ClusterConfig = clusterCfg
	}

	// -------------------------------------------------------------------
	// EBS options
	// -------------------------------------------------------------------

	if ebs := spec.EbsOptions; ebs != nil {
		ebsOpts := &opensearch.DomainEbsOptionsArgs{
			EbsEnabled: pulumi.Bool(ebs.EbsEnabled),
		}
		if ebs.VolumeType != "" {
			ebsOpts.VolumeType = pulumi.String(ebs.VolumeType)
		}
		if ebs.VolumeSize > 0 {
			ebsOpts.VolumeSize = pulumi.Int(int(ebs.VolumeSize))
		}
		if ebs.Iops > 0 {
			ebsOpts.Iops = pulumi.Int(int(ebs.Iops))
		}
		if ebs.Throughput > 0 {
			ebsOpts.Throughput = pulumi.Int(int(ebs.Throughput))
		}
		args.EbsOptions = ebsOpts
	}

	// -------------------------------------------------------------------
	// Encryption at rest -- one-way, and the KMS key is fixed at creation
	// -------------------------------------------------------------------

	if spec.EncryptAtRestEnabled {
		encryptCfg := &opensearch.DomainEncryptAtRestArgs{
			Enabled: pulumi.Bool(true),
		}
		if spec.KmsKeyId.GetValue() != "" {
			encryptCfg.KmsKeyId = pulumi.String(spec.KmsKeyId.GetValue())
		}
		args.EncryptAtRest = encryptCfg
	}

	// -------------------------------------------------------------------
	// Node-to-node encryption -- one-way
	// -------------------------------------------------------------------

	if spec.NodeToNodeEncryptionEnabled {
		args.NodeToNodeEncryption = &opensearch.DomainNodeToNodeEncryptionArgs{
			Enabled: pulumi.Bool(true),
		}
	}

	// -------------------------------------------------------------------
	// VPC options (ForceNew — choose carefully)
	// -------------------------------------------------------------------

	if vpc := spec.VpcOptions; vpc != nil {
		vpcOpts := &opensearch.DomainVpcOptionsArgs{}

		var subnetIds pulumi.StringArray
		for _, s := range vpc.SubnetIds {
			if s.GetValue() != "" {
				subnetIds = append(subnetIds, pulumi.String(s.GetValue()))
			}
		}
		if len(subnetIds) > 0 {
			vpcOpts.SubnetIds = subnetIds
		}

		var sgIds pulumi.StringArray
		for _, sg := range vpc.SecurityGroupIds {
			if sg.GetValue() != "" {
				sgIds = append(sgIds, pulumi.String(sg.GetValue()))
			}
		}
		if len(sgIds) > 0 {
			vpcOpts.SecurityGroupIds = sgIds
		}

		args.VpcOptions = vpcOpts
	}

	// -------------------------------------------------------------------
	// Domain endpoint options (HTTPS, TLS, custom endpoint)
	// -------------------------------------------------------------------

	if deo := spec.DomainEndpointOptions; deo != nil {
		endpointOpts := &opensearch.DomainDomainEndpointOptionsArgs{}

		if deo.EnforceHttps != nil {
			endpointOpts.EnforceHttps = pulumi.Bool(*deo.EnforceHttps)
		}
		if deo.TlsSecurityPolicy != "" {
			endpointOpts.TlsSecurityPolicy = pulumi.String(deo.TlsSecurityPolicy)
		}
		if deo.CustomEndpointEnabled {
			endpointOpts.CustomEndpointEnabled = pulumi.Bool(true)
			if deo.CustomEndpoint != "" {
				endpointOpts.CustomEndpoint = pulumi.String(deo.CustomEndpoint)
			}
			if deo.CustomEndpointCertificateArn.GetValue() != "" {
				endpointOpts.CustomEndpointCertificateArn = pulumi.String(deo.CustomEndpointCertificateArn.GetValue())
			}
		}

		args.DomainEndpointOptions = endpointOpts
	}

	// -------------------------------------------------------------------
	// Advanced security options (FGAC). Emitted only when enabled -- FGAC is
	// one-way in AWS. JWT bearer auth and anonymous auth ride inside it.
	// -------------------------------------------------------------------

	if aso := spec.AdvancedSecurityOptions; aso != nil && aso.Enabled {
		secOpts := &opensearch.DomainAdvancedSecurityOptionsArgs{
			Enabled:                     pulumi.Bool(true),
			InternalUserDatabaseEnabled: pulumi.Bool(aso.InternalUserDatabaseEnabled),
		}

		if aso.AnonymousAuthEnabled {
			secOpts.AnonymousAuthEnabled = pulumi.Bool(true)
		}

		masterUserOpts := &opensearch.DomainAdvancedSecurityOptionsMasterUserOptionsArgs{}
		hasMasterUser := false

		if aso.MasterUserArn.GetValue() != "" {
			masterUserOpts.MasterUserArn = pulumi.String(aso.MasterUserArn.GetValue())
			hasMasterUser = true
		}
		if aso.MasterUserName != "" {
			masterUserOpts.MasterUserName = pulumi.String(aso.MasterUserName)
			hasMasterUser = true
		}
		if aso.MasterUserPassword.GetValue() != "" {
			masterUserOpts.MasterUserPassword = pulumi.String(aso.MasterUserPassword.GetValue())
			hasMasterUser = true
		}

		if hasMasterUser {
			secOpts.MasterUserOptions = masterUserOpts
		}

		if jwt := aso.JwtOptions; jwt != nil {
			jwtOpts := &opensearch.DomainAdvancedSecurityOptionsJwtOptionsArgs{
				Enabled: pulumi.Bool(jwt.Enabled),
			}
			if jwt.JwksUrl != "" {
				jwtOpts.JwksUrl = pulumi.String(jwt.JwksUrl)
			}
			if jwt.PublicKey != "" {
				jwtOpts.PublicKey = pulumi.String(jwt.PublicKey)
			}
			if jwt.RolesKey != "" {
				jwtOpts.RolesKey = pulumi.String(jwt.RolesKey)
			}
			if jwt.SubjectKey != "" {
				jwtOpts.SubjectKey = pulumi.String(jwt.SubjectKey)
			}
			secOpts.JwtOptions = jwtOpts
		}

		args.AdvancedSecurityOptions = secOpts
	}

	// -------------------------------------------------------------------
	// Cognito authentication for OpenSearch Dashboards
	// -------------------------------------------------------------------

	if cog := spec.CognitoOptions; cog != nil && cog.Enabled {
		args.CognitoOptions = &opensearch.DomainCognitoOptionsArgs{
			Enabled:        pulumi.Bool(true),
			UserPoolId:     pulumi.String(cog.UserPoolId.GetValue()),
			IdentityPoolId: pulumi.String(cog.IdentityPoolId),
			RoleArn:        pulumi.String(cog.RoleArn.GetValue()),
		}
	}

	// -------------------------------------------------------------------
	// Log publishing options
	// -------------------------------------------------------------------

	if len(spec.LogPublishingOptions) > 0 {
		var logOpts opensearch.DomainLogPublishingOptionArray
		for _, lpo := range spec.LogPublishingOptions {
			entry := &opensearch.DomainLogPublishingOptionArgs{
				LogType: pulumi.String(lpo.LogType),
			}
			if lpo.CloudwatchLogGroupArn.GetValue() != "" {
				entry.CloudwatchLogGroupArn = pulumi.String(lpo.CloudwatchLogGroupArn.GetValue())
			}
			if lpo.Enabled != nil {
				entry.Enabled = pulumi.Bool(*lpo.Enabled)
			}
			logOpts = append(logOpts, entry)
		}
		args.LogPublishingOptions = logOpts
	}

	// -------------------------------------------------------------------
	// Access policies (IAM policy document as JSON)
	// -------------------------------------------------------------------

	if spec.AccessPolicies != nil {
		policyJSON, err := json.Marshal(spec.AccessPolicies.AsMap())
		if err != nil {
			return errors.Wrap(err, "failed to serialize access_policies to JSON")
		}
		args.AccessPolicies = pulumi.String(string(policyJSON))
	}

	// -------------------------------------------------------------------
	// Auto-Tune. Non-disruptive JVM tuning applies immediately; blue/green
	// optimizations wait for a maintenance schedule or the off-peak window.
	// Not supported on t2/t3 (burstable) instance types.
	// -------------------------------------------------------------------

	if at := spec.AutoTuneOptions; at != nil {
		autoTune := &opensearch.DomainAutoTuneOptionsArgs{
			DesiredState: pulumi.String(at.DesiredState),
		}
		if at.RollbackOnDisable != "" {
			autoTune.RollbackOnDisable = pulumi.String(at.RollbackOnDisable)
		}
		if at.UseOffPeakWindow {
			autoTune.UseOffPeakWindow = pulumi.Bool(true)
		}
		if len(at.MaintenanceSchedules) > 0 {
			var schedules opensearch.DomainAutoTuneOptionsMaintenanceScheduleArray
			for _, schedule := range at.MaintenanceSchedules {
				schedules = append(schedules, &opensearch.DomainAutoTuneOptionsMaintenanceScheduleArgs{
					StartAt:                     pulumi.String(schedule.StartAt),
					CronExpressionForRecurrence: pulumi.String(schedule.CronExpressionForRecurrence),
					Duration: &opensearch.DomainAutoTuneOptionsMaintenanceScheduleDurationArgs{
						// HOURS is the only duration unit the AWS API supports.
						Unit:  pulumi.String("HOURS"),
						Value: pulumi.Int(int(schedule.DurationHours)),
					},
				})
			}
			autoTune.MaintenanceSchedules = schedules
		}
		args.AutoTuneOptions = autoTune
	}

	// -------------------------------------------------------------------
	// Snapshots, off-peak window, software updates
	// -------------------------------------------------------------------

	if spec.AutomatedSnapshotStartHour != nil {
		args.SnapshotOptions = &opensearch.DomainSnapshotOptionsArgs{
			AutomatedSnapshotStartHour: pulumi.Int(int(*spec.AutomatedSnapshotStartHour)),
		}
	}

	if opw := spec.OffPeakWindowOptions; opw != nil {
		offPeak := &opensearch.DomainOffPeakWindowOptionsArgs{
			Enabled: pulumi.Bool(opw.Enabled),
		}
		if opw.WindowStartHour != nil {
			minutes := 0
			if opw.WindowStartMinute != nil {
				minutes = int(*opw.WindowStartMinute)
			}
			offPeak.OffPeakWindow = &opensearch.DomainOffPeakWindowOptionsOffPeakWindowArgs{
				WindowStartTime: &opensearch.DomainOffPeakWindowOptionsOffPeakWindowWindowStartTimeArgs{
					Hours:   pulumi.Int(int(*opw.WindowStartHour)),
					Minutes: pulumi.Int(minutes),
				},
			}
		}
		args.OffPeakWindowOptions = offPeak
	}

	if spec.AutoSoftwareUpdateEnabled {
		args.SoftwareUpdateOptions = &opensearch.DomainSoftwareUpdateOptionsArgs{
			AutoSoftwareUpdateEnabled: pulumi.Bool(true),
		}
	}

	// -------------------------------------------------------------------
	// Blue/green deployment strategy for config changes that require one
	// -------------------------------------------------------------------

	if spec.DeploymentStrategy != "" {
		args.DeploymentStrategyOptions = &opensearch.DomainDeploymentStrategyOptionsArgs{
			DeploymentStrategy: pulumi.String(spec.DeploymentStrategy),
		}
	}

	// -------------------------------------------------------------------
	// IP address type (one-way: dualstack -> ipv4 replaces the domain)
	// -------------------------------------------------------------------

	if spec.IpAddressType != "" {
		args.IpAddressType = pulumi.String(spec.IpAddressType)
	}

	// -------------------------------------------------------------------
	// Advanced options (low-level key-value configuration)
	// -------------------------------------------------------------------

	if len(spec.AdvancedOptions) > 0 {
		args.AdvancedOptions = pulumi.ToStringMap(spec.AdvancedOptions)
	}

	// -------------------------------------------------------------------
	// AI/ML capabilities
	// -------------------------------------------------------------------

	if aiml := spec.AimlOptions; aiml != nil {
		aimlArgs := &opensearch.DomainAimlOptionsArgs{}
		if aiml.NaturalLanguageQueryGenerationDesiredState != "" {
			aimlArgs.NaturalLanguageQueryGenerationOptions = &opensearch.DomainAimlOptionsNaturalLanguageQueryGenerationOptionsArgs{
				DesiredState: pulumi.String(aiml.NaturalLanguageQueryGenerationDesiredState),
			}
		}
		if aiml.S3VectorsEngineEnabled {
			aimlArgs.S3VectorsEngine = &opensearch.DomainAimlOptionsS3VectorsEngineArgs{
				Enabled: pulumi.Bool(true),
			}
		}
		if aiml.ServerlessVectorAccelerationEnabled {
			aimlArgs.ServerlessVectorAcceleration = &opensearch.DomainAimlOptionsServerlessVectorAccelerationArgs{
				Enabled: pulumi.Bool(true),
			}
		}
		args.AimlOptions = aimlArgs
	}

	// -------------------------------------------------------------------
	// IAM Identity Center
	// -------------------------------------------------------------------

	if idc := spec.IdentityCenterOptions; idc != nil {
		idcArgs := &opensearch.DomainIdentityCenterOptionsArgs{
			EnabledApiAccess: pulumi.Bool(idc.EnabledApiAccess),
		}
		if idc.IdentityCenterInstanceArn != "" {
			idcArgs.IdentityCenterInstanceArn = pulumi.String(idc.IdentityCenterInstanceArn)
		}
		if idc.RolesKey != "" {
			idcArgs.RolesKey = pulumi.String(idc.RolesKey)
		}
		if idc.SubjectKey != "" {
			idcArgs.SubjectKey = pulumi.String(idc.SubjectKey)
		}
		args.IdentityCenterOptions = idcArgs
	}

	// -------------------------------------------------------------------
	// Create the OpenSearch domain
	// -------------------------------------------------------------------

	osDomain, err := opensearch.NewDomain(ctx, "opensearch-domain", args, pulumi.Provider(provider))
	if err != nil {
		return errors.Wrap(err, "create opensearch domain")
	}

	// -------------------------------------------------------------------
	// Export outputs
	// -------------------------------------------------------------------

	ctx.Export(OpDomainId, osDomain.DomainId)
	ctx.Export(OpDomainName, osDomain.DomainName)
	ctx.Export(OpDomainArn, osDomain.Arn)
	ctx.Export(OpEndpoint, osDomain.Endpoint)
	ctx.Export(OpDashboardEndpoint, osDomain.DashboardEndpoint)
	ctx.Export(OpEndpointV2, osDomain.EndpointV2)
	ctx.Export(OpDashboardEndpointV2, osDomain.DashboardEndpointV2)
	ctx.Export(OpDomainEndpointV2HostedZoneId, osDomain.DomainEndpointV2HostedZoneId)

	return nil
}
