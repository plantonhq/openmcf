package module

import (
	"github.com/pkg/errors"
	azuredatafactorylinkedservicev1alpha1 "github.com/plantonhq/planton/catalog/azure/azuredatafactorylinkedservice/v1alpha1"
	"github.com/pulumi/pulumi-azure/sdk/v6/go/azure/datafactory"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// The database-family builders: Azure SQL Database, SQL Server, SQL
// Managed Instance, Synapse, PostgreSQL, MySQL, Snowflake, and Kusto.

func createAzureSqlDatabase(
	ctx *pulumi.Context,
	resourceName string,
	spec *azuredatafactorylinkedservicev1alpha1.AzureDataFactoryLinkedServiceSpec,
	azureProvider pulumi.ProviderResource,
) (pulumi.StringInput, pulumi.StringInput, error) {
	sqlDatabase := spec.AzureSqlDatabase

	args := &datafactory.LinkedServiceAzureSqlDatabaseArgs{
		Name:                   pulumi.String(spec.Name),
		DataFactoryId:          pulumi.String(spec.DataFactoryId.GetValue()),
		Description:            descriptionPtr(spec),
		Annotations:            annotationsArray(spec),
		Parameters:             parametersMap(spec),
		AdditionalProperties:   additionalPropertiesMap(spec),
		IntegrationRuntimeName: integrationRuntimeNamePtr(spec),
	}

	// Sent only when true -- the provider's conflict check fires on
	// the argument's PRESENCE, so an explicit false alongside a
	// service principal is rejected; unset means false anyway (the
	// provider's default).
	if useManagedIdentityOrDefault(sqlDatabase.UseManagedIdentity) {
		args.UseManagedIdentity = pulumi.BoolPtr(true)
	}

	if sqlDatabase.ConnectionString != "" {
		args.ConnectionString = pulumi.StringPtr(sqlDatabase.ConnectionString)
	}
	if sqlDatabase.KeyVaultConnectionString != nil {
		args.KeyVaultConnectionString = datafactory.LinkedServiceAzureSqlDatabaseKeyVaultConnectionStringArgs{
			LinkedServiceName: pulumi.String(sqlDatabase.KeyVaultConnectionString.LinkedServiceName.GetValue()),
			SecretName:        pulumi.String(sqlDatabase.KeyVaultConnectionString.SecretName),
		}
	}
	if sqlDatabase.KeyVaultPassword != nil {
		args.KeyVaultPassword = datafactory.LinkedServiceAzureSqlDatabaseKeyVaultPasswordArgs{
			LinkedServiceName: pulumi.String(sqlDatabase.KeyVaultPassword.LinkedServiceName.GetValue()),
			SecretName:        pulumi.String(sqlDatabase.KeyVaultPassword.SecretName),
		}
	}
	if sqlDatabase.ServicePrincipalId != "" {
		args.ServicePrincipalId = pulumi.StringPtr(sqlDatabase.ServicePrincipalId)
	}
	if sqlDatabase.ServicePrincipalKey != "" {
		args.ServicePrincipalKey = pulumi.StringPtr(sqlDatabase.ServicePrincipalKey)
	}
	if sqlDatabase.TenantId != "" {
		args.TenantId = pulumi.StringPtr(sqlDatabase.TenantId)
	}
	if sqlDatabase.CredentialName != "" {
		args.CredentialName = pulumi.StringPtr(sqlDatabase.CredentialName)
	}

	created, err := datafactory.NewLinkedServiceAzureSqlDatabase(ctx, resourceName, args, pulumi.Provider(azureProvider))
	if err != nil {
		return nil, nil, errors.Wrapf(err, "failed to create azure sql database linked service %s", resourceName)
	}
	return created.ID(), created.Name, nil
}

func createSqlServer(
	ctx *pulumi.Context,
	resourceName string,
	spec *azuredatafactorylinkedservicev1alpha1.AzureDataFactoryLinkedServiceSpec,
	azureProvider pulumi.ProviderResource,
) (pulumi.StringInput, pulumi.StringInput, error) {
	sqlServer := spec.SqlServer

	args := &datafactory.LinkedServiceSqlServerArgs{
		Name:                   pulumi.String(spec.Name),
		DataFactoryId:          pulumi.String(spec.DataFactoryId.GetValue()),
		Description:            descriptionPtr(spec),
		Annotations:            annotationsArray(spec),
		Parameters:             parametersMap(spec),
		AdditionalProperties:   additionalPropertiesMap(spec),
		IntegrationRuntimeName: integrationRuntimeNamePtr(spec),
	}

	if sqlServer.ConnectionString != "" {
		args.ConnectionString = pulumi.StringPtr(sqlServer.ConnectionString)
	}
	if sqlServer.KeyVaultConnectionString != nil {
		args.KeyVaultConnectionString = datafactory.LinkedServiceSqlServerKeyVaultConnectionStringArgs{
			LinkedServiceName: pulumi.String(sqlServer.KeyVaultConnectionString.LinkedServiceName.GetValue()),
			SecretName:        pulumi.String(sqlServer.KeyVaultConnectionString.SecretName),
		}
	}
	if sqlServer.KeyVaultPassword != nil {
		args.KeyVaultPassword = datafactory.LinkedServiceSqlServerKeyVaultPasswordArgs{
			LinkedServiceName: pulumi.String(sqlServer.KeyVaultPassword.LinkedServiceName.GetValue()),
			SecretName:        pulumi.String(sqlServer.KeyVaultPassword.SecretName),
		}
	}
	if sqlServer.UserName != "" {
		args.UserName = pulumi.StringPtr(sqlServer.UserName)
	}

	created, err := datafactory.NewLinkedServiceSqlServer(ctx, resourceName, args, pulumi.Provider(azureProvider))
	if err != nil {
		return nil, nil, errors.Wrapf(err, "failed to create sql server linked service %s", resourceName)
	}
	return created.ID(), created.Name, nil
}

func createSqlManagedInstance(
	ctx *pulumi.Context,
	resourceName string,
	spec *azuredatafactorylinkedservicev1alpha1.AzureDataFactoryLinkedServiceSpec,
	azureProvider pulumi.ProviderResource,
) (pulumi.StringInput, pulumi.StringInput, error) {
	sqlMi := spec.SqlManagedInstance

	// The SQL Managed Instance resource is the one linked service the
	// provider models WITHOUT an additional_properties argument -- the
	// spec's shared field does not apply here (noted on the spec
	// field).
	args := &datafactory.LinkedServiceSqlManagedInstanceArgs{
		Name:                   pulumi.String(spec.Name),
		DataFactoryId:          pulumi.String(spec.DataFactoryId.GetValue()),
		Description:            descriptionPtr(spec),
		Annotations:            annotationsArray(spec),
		Parameters:             parametersMap(spec),
		IntegrationRuntimeName: integrationRuntimeNamePtr(spec),
	}

	if sqlMi.ConnectionString != "" {
		args.ConnectionString = pulumi.StringPtr(sqlMi.ConnectionString)
	}
	if sqlMi.KeyVaultConnectionString != nil {
		args.KeyVaultConnectionString = datafactory.LinkedServiceSqlManagedInstanceKeyVaultConnectionStringArgs{
			LinkedServiceName: pulumi.String(sqlMi.KeyVaultConnectionString.LinkedServiceName.GetValue()),
			SecretName:        pulumi.String(sqlMi.KeyVaultConnectionString.SecretName),
		}
	}
	if sqlMi.KeyVaultPassword != nil {
		args.KeyVaultPassword = datafactory.LinkedServiceSqlManagedInstanceKeyVaultPasswordArgs{
			LinkedServiceName: pulumi.String(sqlMi.KeyVaultPassword.LinkedServiceName.GetValue()),
			SecretName:        pulumi.String(sqlMi.KeyVaultPassword.SecretName),
		}
	}
	if sqlMi.ServicePrincipalId != "" {
		args.ServicePrincipalId = pulumi.StringPtr(sqlMi.ServicePrincipalId)
	}
	if sqlMi.ServicePrincipalKey != "" {
		args.ServicePrincipalKey = pulumi.StringPtr(sqlMi.ServicePrincipalKey)
	}
	if sqlMi.Tenant != "" {
		args.Tenant = pulumi.StringPtr(sqlMi.Tenant)
	}

	created, err := datafactory.NewLinkedServiceSqlManagedInstance(ctx, resourceName, args, pulumi.Provider(azureProvider))
	if err != nil {
		return nil, nil, errors.Wrapf(err, "failed to create sql managed instance linked service %s", resourceName)
	}
	return created.ID(), created.Name, nil
}

func createSynapse(
	ctx *pulumi.Context,
	resourceName string,
	spec *azuredatafactorylinkedservicev1alpha1.AzureDataFactoryLinkedServiceSpec,
	azureProvider pulumi.ProviderResource,
) (pulumi.StringInput, pulumi.StringInput, error) {
	synapse := spec.Synapse

	args := &datafactory.LinkedServiceSynapseArgs{
		Name:                   pulumi.String(spec.Name),
		DataFactoryId:          pulumi.String(spec.DataFactoryId.GetValue()),
		Description:            descriptionPtr(spec),
		Annotations:            annotationsArray(spec),
		Parameters:             parametersMap(spec),
		AdditionalProperties:   additionalPropertiesMap(spec),
		IntegrationRuntimeName: integrationRuntimeNamePtr(spec),
		ConnectionString:       pulumi.String(synapse.ConnectionString),
	}

	if synapse.KeyVaultPassword != nil {
		args.KeyVaultPassword = datafactory.LinkedServiceSynapseKeyVaultPasswordArgs{
			LinkedServiceName: pulumi.String(synapse.KeyVaultPassword.LinkedServiceName.GetValue()),
			SecretName:        pulumi.String(synapse.KeyVaultPassword.SecretName),
		}
	}

	created, err := datafactory.NewLinkedServiceSynapse(ctx, resourceName, args, pulumi.Provider(azureProvider))
	if err != nil {
		return nil, nil, errors.Wrapf(err, "failed to create synapse linked service %s", resourceName)
	}
	return created.ID(), created.Name, nil
}

func createPostgresql(
	ctx *pulumi.Context,
	resourceName string,
	spec *azuredatafactorylinkedservicev1alpha1.AzureDataFactoryLinkedServiceSpec,
	azureProvider pulumi.ProviderResource,
) (pulumi.StringInput, pulumi.StringInput, error) {
	args := &datafactory.LinkedServicePostgresqlArgs{
		Name:                   pulumi.String(spec.Name),
		DataFactoryId:          pulumi.String(spec.DataFactoryId.GetValue()),
		Description:            descriptionPtr(spec),
		Annotations:            annotationsArray(spec),
		Parameters:             parametersMap(spec),
		AdditionalProperties:   additionalPropertiesMap(spec),
		IntegrationRuntimeName: integrationRuntimeNamePtr(spec),
		ConnectionString:       pulumi.String(spec.Postgresql.ConnectionString),
	}

	created, err := datafactory.NewLinkedServicePostgresql(ctx, resourceName, args, pulumi.Provider(azureProvider))
	if err != nil {
		return nil, nil, errors.Wrapf(err, "failed to create postgresql linked service %s", resourceName)
	}
	return created.ID(), created.Name, nil
}

func createMysql(
	ctx *pulumi.Context,
	resourceName string,
	spec *azuredatafactorylinkedservicev1alpha1.AzureDataFactoryLinkedServiceSpec,
	azureProvider pulumi.ProviderResource,
) (pulumi.StringInput, pulumi.StringInput, error) {
	// PARITY-EXCEPTION: the classic SDK's LinkedServiceMysqlArgs
	// (pulumi-azure v6.38.0) does not carry the provider's
	// driver_version argument -- this engine cannot send it, so Azure
	// applies its own driver default. Recorded in the kind's
	// provider-parity manifest; remove when the SDK catches up.
	args := &datafactory.LinkedServiceMysqlArgs{
		Name:                   pulumi.String(spec.Name),
		DataFactoryId:          pulumi.String(spec.DataFactoryId.GetValue()),
		Description:            descriptionPtr(spec),
		Annotations:            annotationsArray(spec),
		Parameters:             parametersMap(spec),
		AdditionalProperties:   additionalPropertiesMap(spec),
		IntegrationRuntimeName: integrationRuntimeNamePtr(spec),
		ConnectionString:       pulumi.String(spec.Mysql.ConnectionString),
	}

	created, err := datafactory.NewLinkedServiceMysql(ctx, resourceName, args, pulumi.Provider(azureProvider))
	if err != nil {
		return nil, nil, errors.Wrapf(err, "failed to create mysql linked service %s", resourceName)
	}
	return created.ID(), created.Name, nil
}

func createSnowflake(
	ctx *pulumi.Context,
	resourceName string,
	spec *azuredatafactorylinkedservicev1alpha1.AzureDataFactoryLinkedServiceSpec,
	azureProvider pulumi.ProviderResource,
) (pulumi.StringInput, pulumi.StringInput, error) {
	snowflake := spec.Snowflake

	args := &datafactory.LinkedServiceSnowflakeArgs{
		Name:                   pulumi.String(spec.Name),
		DataFactoryId:          pulumi.String(spec.DataFactoryId.GetValue()),
		Description:            descriptionPtr(spec),
		Annotations:            annotationsArray(spec),
		Parameters:             parametersMap(spec),
		AdditionalProperties:   additionalPropertiesMap(spec),
		IntegrationRuntimeName: integrationRuntimeNamePtr(spec),
		ConnectionString:       pulumi.String(snowflake.ConnectionString),
	}

	if snowflake.KeyVaultPassword != nil {
		args.KeyVaultPassword = datafactory.LinkedServiceSnowflakeKeyVaultPasswordArgs{
			LinkedServiceName: pulumi.String(snowflake.KeyVaultPassword.LinkedServiceName.GetValue()),
			SecretName:        pulumi.String(snowflake.KeyVaultPassword.SecretName),
		}
	}

	created, err := datafactory.NewLinkedServiceSnowflake(ctx, resourceName, args, pulumi.Provider(azureProvider))
	if err != nil {
		return nil, nil, errors.Wrapf(err, "failed to create snowflake linked service %s", resourceName)
	}
	return created.ID(), created.Name, nil
}

func createKusto(
	ctx *pulumi.Context,
	resourceName string,
	spec *azuredatafactorylinkedservicev1alpha1.AzureDataFactoryLinkedServiceSpec,
	azureProvider pulumi.ProviderResource,
) (pulumi.StringInput, pulumi.StringInput, error) {
	kusto := spec.Kusto

	args := &datafactory.LinkedServiceKustoArgs{
		Name:                   pulumi.String(spec.Name),
		DataFactoryId:          pulumi.String(spec.DataFactoryId.GetValue()),
		Description:            descriptionPtr(spec),
		Annotations:            annotationsArray(spec),
		Parameters:             parametersMap(spec),
		AdditionalProperties:   additionalPropertiesMap(spec),
		IntegrationRuntimeName: integrationRuntimeNamePtr(spec),
		KustoEndpoint:          pulumi.String(kusto.KustoEndpoint),
		KustoDatabaseName:      pulumi.String(kusto.KustoDatabaseName),
	}

	// Sent only when true -- the provider's ExactlyOneOf pair reads an
	// explicit false as "this mode declared", which would collide with
	// the service principal actually chosen.
	if useManagedIdentityOrDefault(kusto.UseManagedIdentity) {
		args.UseManagedIdentity = pulumi.BoolPtr(true)
	}
	if kusto.ServicePrincipalId != "" {
		args.ServicePrincipalId = pulumi.StringPtr(kusto.ServicePrincipalId)
	}
	if kusto.ServicePrincipalKey != "" {
		args.ServicePrincipalKey = pulumi.StringPtr(kusto.ServicePrincipalKey)
	}
	if kusto.Tenant != "" {
		args.Tenant = pulumi.StringPtr(kusto.Tenant)
	}

	created, err := datafactory.NewLinkedServiceKusto(ctx, resourceName, args, pulumi.Provider(azureProvider))
	if err != nil {
		return nil, nil, errors.Wrapf(err, "failed to create kusto linked service %s", resourceName)
	}
	return created.ID(), created.Name, nil
}
