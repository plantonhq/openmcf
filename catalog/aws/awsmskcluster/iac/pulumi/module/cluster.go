package module

import (
	"github.com/pkg/errors"
	awsmskclusterv1alpha1 "github.com/plantonhq/planton/catalog/aws/awsmskcluster/v1alpha1"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/msk"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// cluster creates the MSK cluster. Broker networking (subnets, security
// groups) and encryption_info are create-time in AWS: changing them replaces
// the cluster. Compute (instance type), storage size, broker count (increase
// only), monitoring, configuration, and connectivity are all granular in-place
// update operations the provider drives.
func cluster(
	ctx *pulumi.Context,
	locals *Locals,
	provider *aws.Provider,
	createdConfig *msk.Configuration,
) (*msk.Cluster, error) {
	spec := locals.AwsMskCluster.Spec

	var subnets pulumi.StringArray
	for _, s := range spec.SubnetIds {
		subnets = append(subnets, pulumi.String(s.GetValue()))
	}

	// Attached directly -- ingress rules live on the referenced first-class
	// security-group nodes, never on a module-managed shadow group.
	var securityGroups pulumi.StringArray
	for _, sg := range spec.SecurityGroupIds {
		securityGroups = append(securityGroups, pulumi.String(sg.GetValue()))
	}

	// Storage info
	var storageInfo msk.ClusterBrokerNodeGroupInfoStorageInfoPtrInput
	if spec.EbsVolumeSizeGib != nil || spec.ProvisionedThroughputEnabled {
		ebsArgs := &msk.ClusterBrokerNodeGroupInfoStorageInfoEbsStorageInfoArgs{}
		if spec.EbsVolumeSizeGib != nil {
			ebsArgs.VolumeSize = pulumi.Int(int(spec.GetEbsVolumeSizeGib()))
		}
		if spec.ProvisionedThroughputEnabled {
			ebsArgs.ProvisionedThroughput = &msk.ClusterBrokerNodeGroupInfoStorageInfoEbsStorageInfoProvisionedThroughputArgs{
				Enabled:          pulumi.Bool(true),
				VolumeThroughput: pulumi.Int(int(spec.ProvisionedThroughputMbs)),
			}
		}
		storageInfo = &msk.ClusterBrokerNodeGroupInfoStorageInfoArgs{
			EbsStorageInfo: ebsArgs,
		}
	}

	// Connectivity: AWS activates public access, PrivateLink auth schemes, and
	// dual-stack addressing as follow-up updates after the cluster is created;
	// the provider drives that create-then-update flow from this one
	// declarative block.
	connectivityInfo := buildConnectivityInfo(spec)

	args := &msk.ClusterArgs{
		ClusterName:         pulumi.String(locals.ClusterName),
		KafkaVersion:        pulumi.String(spec.KafkaVersion),
		NumberOfBrokerNodes: pulumi.Int(int(spec.NumberOfBrokerNodes)),
		BrokerNodeGroupInfo: &msk.ClusterBrokerNodeGroupInfoArgs{
			InstanceType:     pulumi.String(spec.InstanceType),
			ClientSubnets:    subnets,
			SecurityGroups:   securityGroups,
			StorageInfo:      storageInfo,
			ConnectivityInfo: connectivityInfo,
		},
		Tags: pulumi.ToStringMap(locals.Labels),
	}

	if spec.StorageMode != "" {
		args.StorageMode = pulumi.String(spec.StorageMode)
	}

	// Encryption. encryption_info is create-time: the at-rest KMS key and
	// in-cluster TLS cannot be changed after creation.
	encryptionInTransit := &msk.ClusterEncryptionInfoEncryptionInTransitArgs{}
	hasEncryptionInTransit := false

	if spec.ClientBrokerEncryption != nil && spec.GetClientBrokerEncryption() != "" {
		encryptionInTransit.ClientBroker = pulumi.String(spec.GetClientBrokerEncryption())
		hasEncryptionInTransit = true
	}
	if spec.InClusterEncryption != nil {
		encryptionInTransit.InCluster = pulumi.Bool(spec.GetInClusterEncryption())
		hasEncryptionInTransit = true
	}

	encryptionArgs := &msk.ClusterEncryptionInfoArgs{}
	hasEncryption := false

	if spec.KmsKeyArn.GetValue() != "" {
		encryptionArgs.EncryptionAtRestKmsKeyArn = pulumi.String(spec.KmsKeyArn.GetValue())
		hasEncryption = true
	}
	if hasEncryptionInTransit {
		encryptionArgs.EncryptionInTransit = encryptionInTransit
		hasEncryption = true
	}
	if hasEncryption {
		args.EncryptionInfo = encryptionArgs
	}

	// Authentication
	if spec.Authentication != nil {
		auth := spec.Authentication
		authArgs := &msk.ClusterClientAuthenticationArgs{}
		hasAuth := false

		if auth.SaslIamEnabled || auth.SaslScramEnabled {
			authArgs.Sasl = &msk.ClusterClientAuthenticationSaslArgs{
				Iam:   pulumi.Bool(auth.SaslIamEnabled),
				Scram: pulumi.Bool(auth.SaslScramEnabled),
			}
			hasAuth = true
		}

		if auth.TlsEnabled {
			tlsArgs := &msk.ClusterClientAuthenticationTlsArgs{}
			if len(auth.TlsCertificateAuthorityArns) > 0 {
				var caArns pulumi.StringArray
				for _, ca := range auth.TlsCertificateAuthorityArns {
					caArns = append(caArns, pulumi.String(ca.GetValue()))
				}
				tlsArgs.CertificateAuthorityArns = caArns
			}
			authArgs.Tls = tlsArgs
			hasAuth = true
		}

		if auth.Unauthenticated {
			authArgs.Unauthenticated = pulumi.Bool(true)
			hasAuth = true
		}

		if hasAuth {
			args.ClientAuthentication = authArgs
		}
	}

	// Configuration: the module-managed configuration (from server_properties)
	// or an externally managed one (configuration_arn + revision). An update to
	// inline properties bumps the configuration revision and the cluster follows
	// -- a rolling broker restart, never a replacement.
	if createdConfig != nil {
		args.ConfigurationInfo = &msk.ClusterConfigurationInfoArgs{
			Arn:      createdConfig.Arn,
			Revision: createdConfig.LatestRevision,
		}
	} else if spec.ConfigurationArn != "" {
		args.ConfigurationInfo = &msk.ClusterConfigurationInfoArgs{
			Arn:      pulumi.String(spec.ConfigurationArn),
			Revision: pulumi.Int(int(spec.ConfigurationRevision)),
		}
	}

	// Logging
	if spec.Logging != nil {
		loggingArgs := buildLogging(spec)
		if loggingArgs != nil {
			args.LoggingInfo = loggingArgs
		}
	}

	// Monitoring
	if spec.EnhancedMonitoring != "" {
		args.EnhancedMonitoring = pulumi.String(spec.EnhancedMonitoring)
	}

	if spec.JmxExporterEnabled || spec.NodeExporterEnabled {
		args.OpenMonitoring = &msk.ClusterOpenMonitoringArgs{
			Prometheus: &msk.ClusterOpenMonitoringPrometheusArgs{
				JmxExporter: &msk.ClusterOpenMonitoringPrometheusJmxExporterArgs{
					EnabledInBroker: pulumi.Bool(spec.JmxExporterEnabled),
				},
				NodeExporter: &msk.ClusterOpenMonitoringPrometheusNodeExporterArgs{
					EnabledInBroker: pulumi.Bool(spec.NodeExporterEnabled),
				},
			},
		}
	}

	// Intelligent rebalancing is an Express-broker (express.* instance type)
	// capability; AWS rejects it on standard kafka.* clusters, so it is set
	// only when the spec opts in.
	if spec.RebalancingStatus != "" {
		args.Rebalancing = &msk.ClusterRebalancingArgs{
			Status: pulumi.String(spec.RebalancingStatus),
		}
	}

	mskCluster, err := msk.NewCluster(ctx, "msk-cluster", args, pulumi.Provider(provider))
	if err != nil {
		return nil, errors.Wrap(err, "create msk cluster")
	}

	return mskCluster, nil
}

// buildConnectivityInfo assembles the connectivity_info block. It is emitted
// only when some connectivity surface is configured; an empty block would
// still trigger AWS's create-then-update connectivity flow for no reason.
func buildConnectivityInfo(spec *awsmskclusterv1alpha1.AwsMskClusterSpec) msk.ClusterBrokerNodeGroupInfoConnectivityInfoPtrInput {
	vpcConnectivityEnabled := spec.VpcConnectivity != nil &&
		(spec.VpcConnectivity.SaslIamEnabled || spec.VpcConnectivity.SaslScramEnabled || spec.VpcConnectivity.TlsEnabled)

	if spec.PublicAccessType == "" && !vpcConnectivityEnabled && spec.NetworkType == "" {
		return nil
	}

	connectivityArgs := &msk.ClusterBrokerNodeGroupInfoConnectivityInfoArgs{}

	if spec.NetworkType != "" {
		connectivityArgs.NetworkType = pulumi.String(spec.NetworkType)
	}

	if spec.PublicAccessType != "" {
		connectivityArgs.PublicAccess = &msk.ClusterBrokerNodeGroupInfoConnectivityInfoPublicAccessArgs{
			Type: pulumi.String(spec.PublicAccessType),
		}
	}

	if vpcConnectivityEnabled {
		clientAuthArgs := &msk.ClusterBrokerNodeGroupInfoConnectivityInfoVpcConnectivityClientAuthenticationArgs{
			Tls: pulumi.Bool(spec.VpcConnectivity.TlsEnabled),
		}
		if spec.VpcConnectivity.SaslIamEnabled || spec.VpcConnectivity.SaslScramEnabled {
			clientAuthArgs.Sasl = &msk.ClusterBrokerNodeGroupInfoConnectivityInfoVpcConnectivityClientAuthenticationSaslArgs{
				Iam:   pulumi.Bool(spec.VpcConnectivity.SaslIamEnabled),
				Scram: pulumi.Bool(spec.VpcConnectivity.SaslScramEnabled),
			}
		}
		connectivityArgs.VpcConnectivity = &msk.ClusterBrokerNodeGroupInfoConnectivityInfoVpcConnectivityArgs{
			ClientAuthentication: clientAuthArgs,
		}
	}

	return connectivityArgs
}

// buildLogging constructs the logging configuration from the spec.
func buildLogging(spec *awsmskclusterv1alpha1.AwsMskClusterSpec) *msk.ClusterLoggingInfoArgs {
	logging := spec.Logging
	if logging == nil {
		return nil
	}

	brokerLogsArgs := &msk.ClusterLoggingInfoBrokerLogsArgs{}
	hasBrokerLogs := false

	if logging.CloudwatchLogs != nil {
		cwArgs := &msk.ClusterLoggingInfoBrokerLogsCloudwatchLogsArgs{
			Enabled: pulumi.Bool(logging.CloudwatchLogs.Enabled),
		}
		if logging.CloudwatchLogs.LogGroup.GetValue() != "" {
			cwArgs.LogGroup = pulumi.String(logging.CloudwatchLogs.LogGroup.GetValue())
		}
		brokerLogsArgs.CloudwatchLogs = cwArgs
		hasBrokerLogs = true
	}

	if logging.Firehose != nil {
		fhArgs := &msk.ClusterLoggingInfoBrokerLogsFirehoseArgs{
			Enabled: pulumi.Bool(logging.Firehose.Enabled),
		}
		if logging.Firehose.DeliveryStream.GetValue() != "" {
			fhArgs.DeliveryStream = pulumi.String(logging.Firehose.DeliveryStream.GetValue())
		}
		brokerLogsArgs.Firehose = fhArgs
		hasBrokerLogs = true
	}

	if logging.S3 != nil {
		s3Args := &msk.ClusterLoggingInfoBrokerLogsS3Args{
			Enabled: pulumi.Bool(logging.S3.Enabled),
		}
		if logging.S3.Bucket.GetValue() != "" {
			s3Args.Bucket = pulumi.String(logging.S3.Bucket.GetValue())
		}
		if logging.S3.Prefix != "" {
			s3Args.Prefix = pulumi.String(logging.S3.Prefix)
		}
		brokerLogsArgs.S3 = s3Args
		hasBrokerLogs = true
	}

	if !hasBrokerLogs {
		return nil
	}

	return &msk.ClusterLoggingInfoArgs{
		BrokerLogs: brokerLogsArgs,
	}
}
