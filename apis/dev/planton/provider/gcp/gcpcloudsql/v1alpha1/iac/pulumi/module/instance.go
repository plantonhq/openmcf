package module

import (
	"github.com/pkg/errors"
	gcpcloudsqlv1alpha1 "github.com/plantonhq/planton/apis/dev/planton/provider/gcp/gcpcloudsql/v1alpha1"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/projects"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/sql"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// databaseInstance provisions a Cloud SQL instance — a managed MySQL,
// PostgreSQL, or SQL Server server. One resource is one instance: a primary,
// or (with master_instance_name) a read replica. Databases and users inside
// it are separate composable resources referencing the instance by name.
//
// Lifecycle notes the API enforces:
//   - name/region/CMEK/disk type are immutable; deleted names stay reserved
//     for ~1 week.
//   - database_version upgrades and tier/edition changes apply IN PLACE
//     (with a restart) — no replacement.
//   - disks grow but never shrink; private_network can be set or changed but
//     never removed in place.
//   - a private-IP instance requires the VPC to already carry a service
//     networking connection; the provider prechecks and fails fast because a
//     failed create still burns the reserved instance name.
func databaseInstance(
	ctx *pulumi.Context,
	locals *Locals,
	gcpProvider *gcp.Provider,
) (*sql.DatabaseInstance, error) {
	spec := locals.GcpCloudSql.Spec

	// Enable the Cloud SQL Admin API first so a fresh project works on the
	// first deploy. disable_on_destroy stays false: tearing down one
	// instance must never disable the API for everything else in the
	// project.
	serviceArgs := &projects.ServiceArgs{
		Service:                  pulumi.String("sqladmin.googleapis.com"),
		DisableDependentServices: pulumi.BoolPtr(true),
		DisableOnDestroy:         pulumi.BoolPtr(false),
	}
	if spec.ProjectId.GetValue() != "" {
		serviceArgs.Project = pulumi.String(spec.ProjectId.GetValue())
	}
	createdProjectService, err := projects.NewService(ctx,
		"cloudsql-sqladmin.googleapis.com", serviceArgs, pulumi.Provider(gcpProvider))
	if err != nil {
		return nil, errors.Wrap(err, "failed to enable sqladmin.googleapis.com api")
	}

	settings := buildSettings(locals)

	args := &sql.DatabaseInstanceArgs{
		Name:            pulumi.String(spec.InstanceName),
		Region:          pulumi.String(spec.Region),
		DatabaseVersion: pulumi.String(spec.DatabaseVersion),
		Settings:        settings,
		// Engine-side destroy guard: refuses `destroy` at preview time while
		// true. Distinct from settings.deletion_protection_enabled (the
		// API-side guard).
		DeletionProtection: pulumi.Bool(spec.DeletionProtection),
	}

	// An empty project falls back to the provider's default project — the
	// ambient-project contract every GCP kind honors.
	if spec.ProjectId.GetValue() != "" {
		args.Project = pulumi.String(spec.ProjectId.GetValue())
	}

	// Set only on read replicas: the primary this instance replicates from.
	if spec.MasterInstanceName.GetValue() != "" {
		args.MasterInstanceName = pulumi.StringPtr(spec.MasterInstanceName.GetValue())
	}

	// CMEK: the key must live in the instance's region. Immutable.
	if spec.EncryptionKeyName.GetValue() != "" {
		args.EncryptionKeyName = pulumi.StringPtr(spec.EncryptionKeyName.GetValue())
	}

	// Write-only in GCP — never readable back, never in outputs. Required
	// for SQL Server (spec CEL enforces pre-deploy). Marked secret so the
	// value is encrypted in Pulumi state.
	if spec.RootPassword != "" {
		args.RootPassword = pulumi.ToSecret(pulumi.String(spec.RootPassword)).(pulumi.StringOutput)
	}

	if rc := buildReplicaConfiguration(spec.ReplicaConfiguration); rc != nil {
		args.ReplicaConfiguration = rc
	}

	createdInstance, err := sql.NewDatabaseInstance(ctx, "instance", args,
		pulumi.Provider(gcpProvider), pulumi.DependsOn([]pulumi.Resource{createdProjectService}))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create cloud sql instance")
	}

	ctx.Export(OpInstanceName, createdInstance.Name)
	ctx.Export(OpConnectionName, createdInstance.ConnectionName)
	ctx.Export(OpPrivateIp, createdInstance.PrivateIpAddress)
	ctx.Export(OpPublicIp, createdInstance.PublicIpAddress)
	ctx.Export(OpSelfLink, createdInstance.SelfLink)
	ctx.Export(OpServiceAccountEmail, createdInstance.ServiceAccountEmailAddress)
	ctx.Export(OpDnsName, createdInstance.DnsName)
	ctx.Export(OpPscServiceAttachmentLink, createdInstance.PscServiceAttachmentLink)

	return createdInstance, nil
}

// buildSettings translates the spec into the provider's settings block,
// applying the same omitted-block defaults as the Terraform module: no disk
// block means 10 GB PD_SSD with auto-resize; no network block means public
// IPv4 with no authorized networks (Auth Proxy / connector access only).
func buildSettings(locals *Locals) *sql.DatabaseInstanceSettingsArgs {
	spec := locals.GcpCloudSql.Spec

	settings := &sql.DatabaseInstanceSettingsArgs{
		Tier:             pulumi.String(spec.Tier),
		Edition:          pulumi.StringPtr(spec.GetEdition()),
		AvailabilityType: pulumi.StringPtr(spec.GetAvailabilityType()),
		ActivationPolicy: pulumi.StringPtr(spec.GetActivationPolicy()),
		UserLabels:       pulumi.ToStringMap(locals.GcpLabels),
		// API-side delete guard — blocks deletion from console/gcloud/API too.
		DeletionProtectionEnabled: pulumi.BoolPtr(spec.DeletionProtectionEnabled),
		// Keep automated backups (and PITR logs) after instance deletion.
		RetainBackupsOnDelete:     pulumi.BoolPtr(spec.RetainBackupsOnDelete),
		EnableGoogleMlIntegration: pulumi.BoolPtr(spec.EnableGoogleMlIntegration),
		EnableDataplexIntegration: pulumi.BoolPtr(spec.EnableDataplexIntegration),
	}

	// Disk defaults when the block is omitted: 10 GB PD_SSD, auto-resize on.
	diskType, diskSizeGb, diskAutoResize, diskAutoResizeLimit := "PD_SSD", 10, true, 0
	if spec.Disk != nil {
		diskType = spec.Disk.GetType()
		diskSizeGb = int(spec.Disk.GetSizeGb())
		diskAutoResize = spec.Disk.GetAutoResize()
		diskAutoResizeLimit = int(spec.Disk.GetAutoResizeLimit())
	}
	settings.DiskType = pulumi.StringPtr(diskType)
	settings.DiskSize = pulumi.IntPtr(diskSizeGb)
	settings.DiskAutoresize = pulumi.BoolPtr(diskAutoResize)
	settings.DiskAutoresizeLimit = pulumi.IntPtr(diskAutoResizeLimit)

	if spec.ConnectorEnforcement != "" {
		settings.ConnectorEnforcement = pulumi.StringPtr(spec.ConnectorEnforcement)
	}

	// SQL Server-only knobs (spec CEL restricts them to SQLSERVER engines).
	if spec.TimeZone != "" {
		settings.TimeZone = pulumi.StringPtr(spec.TimeZone)
	}
	if spec.Collation != "" {
		settings.Collation = pulumi.StringPtr(spec.Collation)
	}
	if spec.ThreadsPerCore != nil {
		settings.AdvancedMachineFeatures = &sql.DatabaseInstanceSettingsAdvancedMachineFeaturesArgs{
			ThreadsPerCore: pulumi.IntPtr(int(*spec.ThreadsPerCore)),
		}
	}
	if spec.SqlServerAuditConfig != nil {
		auditArgs := &sql.DatabaseInstanceSettingsSqlServerAuditConfigArgs{}
		if spec.SqlServerAuditConfig.Bucket != "" {
			auditArgs.Bucket = pulumi.StringPtr(spec.SqlServerAuditConfig.Bucket)
		}
		if spec.SqlServerAuditConfig.RetentionInterval != "" {
			auditArgs.RetentionInterval = pulumi.StringPtr(spec.SqlServerAuditConfig.RetentionInterval)
		}
		if spec.SqlServerAuditConfig.UploadInterval != "" {
			auditArgs.UploadInterval = pulumi.StringPtr(spec.SqlServerAuditConfig.UploadInterval)
		}
		settings.SqlServerAuditConfig = auditArgs
	}
	if spec.ActiveDirectoryDomain != "" {
		settings.ActiveDirectoryConfig = &sql.DatabaseInstanceSettingsActiveDirectoryConfigArgs{
			Domain: pulumi.String(spec.ActiveDirectoryDomain),
		}
	}

	// Emitted only when enabled: the API rejects a data-cache stanza on
	// ENTERPRISE instances (spec CEL already forces ENTERPRISE_PLUS).
	if spec.DataCacheEnabled {
		settings.DataCacheConfig = &sql.DatabaseInstanceSettingsDataCacheConfigArgs{
			DataCacheEnabled: pulumi.BoolPtr(true),
		}
	}

	settings.IpConfiguration = buildIpConfiguration(spec.Network)

	if spec.LocationPreference != nil &&
		(spec.LocationPreference.Zone != "" || spec.LocationPreference.SecondaryZone != "") {
		lpArgs := &sql.DatabaseInstanceSettingsLocationPreferenceArgs{}
		if spec.LocationPreference.Zone != "" {
			lpArgs.Zone = pulumi.StringPtr(spec.LocationPreference.Zone)
		}
		if spec.LocationPreference.SecondaryZone != "" {
			lpArgs.SecondaryZone = pulumi.StringPtr(spec.LocationPreference.SecondaryZone)
		}
		settings.LocationPreference = lpArgs
	}

	if spec.Backup != nil {
		backupArgs := &sql.DatabaseInstanceSettingsBackupConfigurationArgs{
			Enabled: pulumi.BoolPtr(spec.Backup.Enabled),
			// MySQL PITR mechanism; also required for MySQL replicas and HA.
			BinaryLogEnabled: pulumi.BoolPtr(spec.Backup.BinaryLogEnabled),
			// PostgreSQL / SQL Server PITR mechanism.
			PointInTimeRecoveryEnabled: pulumi.BoolPtr(spec.Backup.PointInTimeRecoveryEnabled),
		}
		if spec.Backup.StartTime != "" {
			backupArgs.StartTime = pulumi.StringPtr(spec.Backup.StartTime)
		}
		if spec.Backup.Location != "" {
			backupArgs.Location = pulumi.StringPtr(spec.Backup.Location)
		}
		if spec.Backup.TransactionLogRetentionDays != nil {
			backupArgs.TransactionLogRetentionDays = pulumi.IntPtr(int(*spec.Backup.TransactionLogRetentionDays))
		}
		if spec.Backup.RetainedBackups != nil {
			backupArgs.BackupRetentionSettings = &sql.DatabaseInstanceSettingsBackupConfigurationBackupRetentionSettingsArgs{
				RetainedBackups: pulumi.Int(int(*spec.Backup.RetainedBackups)),
				RetentionUnit:   pulumi.StringPtr("COUNT"),
			}
		}
		settings.BackupConfiguration = backupArgs
	}

	if spec.MaintenanceWindow != nil {
		mwArgs := &sql.DatabaseInstanceSettingsMaintenanceWindowArgs{
			Day: pulumi.IntPtr(int(spec.MaintenanceWindow.Day)),
		}
		if spec.MaintenanceWindow.Hour != nil {
			mwArgs.Hour = pulumi.IntPtr(int(*spec.MaintenanceWindow.Hour))
		}
		if spec.MaintenanceWindow.UpdateTrack != "" {
			mwArgs.UpdateTrack = pulumi.StringPtr(spec.MaintenanceWindow.UpdateTrack)
		}
		settings.MaintenanceWindow = mwArgs
	}

	if spec.DenyMaintenancePeriod != nil {
		settings.DenyMaintenancePeriod = &sql.DatabaseInstanceSettingsDenyMaintenancePeriodArgs{
			StartDate: pulumi.String(spec.DenyMaintenancePeriod.StartDate),
			EndDate:   pulumi.String(spec.DenyMaintenancePeriod.EndDate),
			Time:      pulumi.String(spec.DenyMaintenancePeriod.Time),
		}
	}

	if spec.InsightsConfig != nil && spec.InsightsConfig.QueryInsightsEnabled {
		insightsArgs := &sql.DatabaseInstanceSettingsInsightsConfigArgs{
			QueryInsightsEnabled:  pulumi.BoolPtr(true),
			RecordApplicationTags: pulumi.BoolPtr(spec.InsightsConfig.RecordApplicationTags),
			RecordClientAddress:   pulumi.BoolPtr(spec.InsightsConfig.RecordClientAddress),
		}
		if spec.InsightsConfig.QueryStringLength != nil {
			insightsArgs.QueryStringLength = pulumi.IntPtr(int(*spec.InsightsConfig.QueryStringLength))
		}
		if spec.InsightsConfig.QueryPlansPerMinute != nil {
			insightsArgs.QueryPlansPerMinute = pulumi.IntPtr(int(*spec.InsightsConfig.QueryPlansPerMinute))
		}
		settings.InsightsConfig = insightsArgs
	}

	if spec.PasswordValidationPolicy != nil {
		pvpArgs := &sql.DatabaseInstanceSettingsPasswordValidationPolicyArgs{
			EnablePasswordPolicy:      pulumi.Bool(spec.PasswordValidationPolicy.EnablePasswordPolicy),
			DisallowUsernameSubstring: pulumi.BoolPtr(spec.PasswordValidationPolicy.DisallowUsernameSubstring),
		}
		if spec.PasswordValidationPolicy.MinLength != nil {
			pvpArgs.MinLength = pulumi.IntPtr(int(*spec.PasswordValidationPolicy.MinLength))
		}
		if spec.PasswordValidationPolicy.Complexity != "" {
			pvpArgs.Complexity = pulumi.StringPtr(spec.PasswordValidationPolicy.Complexity)
		}
		if spec.PasswordValidationPolicy.ReuseInterval != nil {
			pvpArgs.ReuseInterval = pulumi.IntPtr(int(*spec.PasswordValidationPolicy.ReuseInterval))
		}
		if spec.PasswordValidationPolicy.PasswordChangeInterval != "" {
			pvpArgs.PasswordChangeInterval = pulumi.StringPtr(spec.PasswordValidationPolicy.PasswordChangeInterval)
		}
		settings.PasswordValidationPolicy = pvpArgs
	}

	if spec.ConnectionPooling != nil && spec.ConnectionPooling.Enabled {
		poolFlags := sql.DatabaseInstanceSettingsConnectionPoolConfigFlagArray{}
		for name, value := range spec.ConnectionPooling.Flags {
			poolFlags = append(poolFlags, &sql.DatabaseInstanceSettingsConnectionPoolConfigFlagArgs{
				Name:  pulumi.String(name),
				Value: pulumi.String(value),
			})
		}
		settings.ConnectionPoolConfigs = sql.DatabaseInstanceSettingsConnectionPoolConfigArray{
			&sql.DatabaseInstanceSettingsConnectionPoolConfigArgs{
				ConnectionPoolingEnabled: pulumi.BoolPtr(true),
				Flags:                    poolFlags,
			},
		}
	}

	if len(spec.DatabaseFlags) > 0 {
		databaseFlags := sql.DatabaseInstanceSettingsDatabaseFlagArray{}
		for name, value := range spec.DatabaseFlags {
			databaseFlags = append(databaseFlags, &sql.DatabaseInstanceSettingsDatabaseFlagArgs{
				Name:  pulumi.String(name),
				Value: pulumi.String(value),
			})
		}
		settings.DatabaseFlags = databaseFlags
	}

	return settings
}

// buildIpConfiguration translates the spec's network block. An omitted block
// resolves to public IPv4 with no authorized networks — reachable only
// through the Cloud SQL Auth Proxy or connectors (IAM-authenticated).
func buildIpConfiguration(network *gcpcloudsqlv1alpha1.GcpCloudSqlNetwork) *sql.DatabaseInstanceSettingsIpConfigurationArgs {
	if network == nil {
		return &sql.DatabaseInstanceSettingsIpConfigurationArgs{
			Ipv4Enabled: pulumi.BoolPtr(true),
		}
	}

	ipArgs := &sql.DatabaseInstanceSettingsIpConfigurationArgs{
		Ipv4Enabled:                             pulumi.BoolPtr(network.Ipv4Enabled),
		EnablePrivatePathForGoogleCloudServices: pulumi.BoolPtr(network.EnablePrivatePathForGoogleCloudServices),
	}

	// Setting a network enables private IP; the VPC must already carry a
	// service networking connection (see the module-level comment).
	if network.PrivateNetwork.GetValue() != "" {
		ipArgs.PrivateNetwork = pulumi.StringPtr(network.PrivateNetwork.GetValue())
	}
	if network.AllocatedIpRange != "" {
		ipArgs.AllocatedIpRange = pulumi.StringPtr(network.AllocatedIpRange)
	}
	if network.SslMode != "" {
		ipArgs.SslMode = pulumi.StringPtr(network.SslMode)
	}
	if network.ServerCaMode != "" {
		ipArgs.ServerCaMode = pulumi.StringPtr(network.ServerCaMode)
	}
	if network.ServerCaPool != "" {
		ipArgs.ServerCaPool = pulumi.StringPtr(network.ServerCaPool)
	}
	if len(network.CustomSubjectAlternativeNames) > 0 {
		ipArgs.CustomSubjectAlternativeNames = pulumi.ToStringArray(network.CustomSubjectAlternativeNames)
	}

	if len(network.AuthorizedNetworks) > 0 {
		authorized := sql.DatabaseInstanceSettingsIpConfigurationAuthorizedNetworkArray{}
		for _, an := range network.AuthorizedNetworks {
			anArgs := &sql.DatabaseInstanceSettingsIpConfigurationAuthorizedNetworkArgs{
				Value: pulumi.String(an.Value),
			}
			if an.Name != "" {
				anArgs.Name = pulumi.StringPtr(an.Name)
			}
			if an.ExpirationTime != "" {
				anArgs.ExpirationTime = pulumi.StringPtr(an.ExpirationTime)
			}
			authorized = append(authorized, anArgs)
		}
		ipArgs.AuthorizedNetworks = authorized
	}

	if network.Psc != nil && network.Psc.Enabled {
		pscArgs := &sql.DatabaseInstanceSettingsIpConfigurationPscConfigArgs{
			PscEnabled: pulumi.BoolPtr(true),
		}
		if len(network.Psc.AllowedConsumerProjects) > 0 {
			pscArgs.AllowedConsumerProjects = pulumi.ToStringArray(network.Psc.AllowedConsumerProjects)
		}
		if network.Psc.NetworkAttachmentUri != "" {
			pscArgs.NetworkAttachmentUri = pulumi.StringPtr(network.Psc.NetworkAttachmentUri)
		}
		if len(network.Psc.AutoConnections) > 0 {
			autoConnections := sql.DatabaseInstanceSettingsIpConfigurationPscConfigPscAutoConnectionArray{}
			for _, ac := range network.Psc.AutoConnections {
				acArgs := &sql.DatabaseInstanceSettingsIpConfigurationPscConfigPscAutoConnectionArgs{
					ConsumerNetwork: pulumi.String(ac.ConsumerNetwork),
				}
				if ac.ConsumerServiceProjectId != "" {
					acArgs.ConsumerServiceProjectId = pulumi.StringPtr(ac.ConsumerServiceProjectId)
				}
				autoConnections = append(autoConnections, acArgs)
			}
			pscArgs.PscAutoConnections = autoConnections
		}
		ipArgs.PscConfigs = sql.DatabaseInstanceSettingsIpConfigurationPscConfigArray{pscArgs}
	}

	return ipArgs
}

// buildReplicaConfiguration translates the replica block. All fields are
// immutable (ForceNew) in the provider; the username/password/certificate
// fields configure the replication channel from an EXTERNAL source, and the
// password and client key are secret material kept encrypted in state.
func buildReplicaConfiguration(rc *gcpcloudsqlv1alpha1.GcpCloudSqlReplicaConfiguration) *sql.DatabaseInstanceReplicaConfigurationArgs {
	if rc == nil {
		return nil
	}

	rcArgs := &sql.DatabaseInstanceReplicaConfigurationArgs{
		FailoverTarget:          pulumi.BoolPtr(rc.FailoverTarget),
		CascadableReplica:       pulumi.BoolPtr(rc.CascadableReplica),
		VerifyServerCertificate: pulumi.BoolPtr(rc.VerifyServerCertificate),
	}
	if rc.Username != "" {
		rcArgs.Username = pulumi.StringPtr(rc.Username)
	}
	if rc.Password != "" {
		// Secret replication-channel material — encrypted in Pulumi state.
		rcArgs.Password = pulumi.ToSecret(pulumi.String(rc.Password)).(pulumi.StringOutput)
	}
	if rc.CaCertificate != "" {
		rcArgs.CaCertificate = pulumi.StringPtr(rc.CaCertificate)
	}
	if rc.ClientCertificate != "" {
		rcArgs.ClientCertificate = pulumi.StringPtr(rc.ClientCertificate)
	}
	if rc.ClientKey != "" {
		// Secret private-key material — encrypted in Pulumi state.
		rcArgs.ClientKey = pulumi.ToSecret(pulumi.String(rc.ClientKey)).(pulumi.StringOutput)
	}
	if rc.DumpFilePath != "" {
		rcArgs.DumpFilePath = pulumi.StringPtr(rc.DumpFilePath)
	}
	if rc.ConnectRetryInterval != nil {
		rcArgs.ConnectRetryInterval = pulumi.IntPtr(int(*rc.ConnectRetryInterval))
	}
	if rc.MasterHeartbeatPeriod != nil {
		rcArgs.MasterHeartbeatPeriod = pulumi.IntPtr(int(*rc.MasterHeartbeatPeriod))
	}
	if rc.SslCipher != "" {
		rcArgs.SslCipher = pulumi.StringPtr(rc.SslCipher)
	}

	return rcArgs
}
