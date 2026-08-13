package module

import (
	"github.com/pkg/errors"
	azuredatafactorylinkedservicev1alpha1 "github.com/plantonhq/planton/catalog/azure/azuredatafactorylinkedservice/v1alpha1"
	"github.com/pulumi/pulumi-azure/sdk/v6/go/azure/datafactory"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// The service and protocol builders: Key Vault, Azure Search, Azure
// Function, Databricks, web, OData, ODBC, SFTP, and the raw-JSON
// custom form.

func createKeyVault(
	ctx *pulumi.Context,
	resourceName string,
	spec *azuredatafactorylinkedservicev1alpha1.AzureDataFactoryLinkedServiceSpec,
	azureProvider pulumi.ProviderResource,
) (pulumi.StringInput, pulumi.StringInput, error) {
	args := &datafactory.LinkedServiceKeyVaultArgs{
		Name:                   pulumi.String(spec.Name),
		DataFactoryId:          pulumi.String(spec.DataFactoryId.GetValue()),
		Description:            descriptionPtr(spec),
		Annotations:            annotationsArray(spec),
		Parameters:             parametersMap(spec),
		AdditionalProperties:   additionalPropertiesMap(spec),
		IntegrationRuntimeName: integrationRuntimeNamePtr(spec),
		KeyVaultId:             pulumi.String(spec.KeyVault.KeyVaultId.GetValue()),
	}

	created, err := datafactory.NewLinkedServiceKeyVault(ctx, resourceName, args, pulumi.Provider(azureProvider))
	if err != nil {
		return nil, nil, errors.Wrapf(err, "failed to create key vault linked service %s", resourceName)
	}
	return created.ID(), created.Name, nil
}

func createAzureSearch(
	ctx *pulumi.Context,
	resourceName string,
	spec *azuredatafactorylinkedservicev1alpha1.AzureDataFactoryLinkedServiceSpec,
	azureProvider pulumi.ProviderResource,
) (pulumi.StringInput, pulumi.StringInput, error) {
	search := spec.AzureSearch

	args := &datafactory.LinkedServiceAzureSearchArgs{
		Name:                   pulumi.String(spec.Name),
		DataFactoryId:          pulumi.String(spec.DataFactoryId.GetValue()),
		Description:            descriptionPtr(spec),
		Annotations:            annotationsArray(spec),
		Parameters:             parametersMap(spec),
		AdditionalProperties:   additionalPropertiesMap(spec),
		IntegrationRuntimeName: integrationRuntimeNamePtr(spec),
		Url:                    pulumi.String(search.Url.GetValue()),
		SearchServiceKey:       pulumi.String(search.SearchServiceKey),
	}

	created, err := datafactory.NewLinkedServiceAzureSearch(ctx, resourceName, args, pulumi.Provider(azureProvider))
	if err != nil {
		return nil, nil, errors.Wrapf(err, "failed to create azure search linked service %s", resourceName)
	}
	return created.ID(), created.Name, nil
}

func createAzureFunction(
	ctx *pulumi.Context,
	resourceName string,
	spec *azuredatafactorylinkedservicev1alpha1.AzureDataFactoryLinkedServiceSpec,
	azureProvider pulumi.ProviderResource,
) (pulumi.StringInput, pulumi.StringInput, error) {
	function := spec.AzureFunction

	args := &datafactory.LinkedServiceAzureFunctionArgs{
		Name:                   pulumi.String(spec.Name),
		DataFactoryId:          pulumi.String(spec.DataFactoryId.GetValue()),
		Description:            descriptionPtr(spec),
		Annotations:            annotationsArray(spec),
		Parameters:             parametersMap(spec),
		AdditionalProperties:   additionalPropertiesMap(spec),
		IntegrationRuntimeName: integrationRuntimeNamePtr(spec),
		Url:                    pulumi.String(function.Url),
	}

	if function.Key != "" {
		args.Key = pulumi.StringPtr(function.Key)
	}
	if function.KeyVaultKey != nil {
		args.KeyVaultKey = datafactory.LinkedServiceAzureFunctionKeyVaultKeyArgs{
			LinkedServiceName: pulumi.String(function.KeyVaultKey.LinkedServiceName.GetValue()),
			SecretName:        pulumi.String(function.KeyVaultKey.SecretName),
		}
	}

	created, err := datafactory.NewLinkedServiceAzureFunction(ctx, resourceName, args, pulumi.Provider(azureProvider))
	if err != nil {
		return nil, nil, errors.Wrapf(err, "failed to create azure function linked service %s", resourceName)
	}
	return created.ID(), created.Name, nil
}

func createAzureDatabricks(
	ctx *pulumi.Context,
	resourceName string,
	spec *azuredatafactorylinkedservicev1alpha1.AzureDataFactoryLinkedServiceSpec,
	azureProvider pulumi.ProviderResource,
) (pulumi.StringInput, pulumi.StringInput, error) {
	databricks := spec.AzureDatabricks

	args := &datafactory.LinkedServiceAzureDatabricksArgs{
		Name:                   pulumi.String(spec.Name),
		DataFactoryId:          pulumi.String(spec.DataFactoryId.GetValue()),
		Description:            descriptionPtr(spec),
		Annotations:            annotationsArray(spec),
		Parameters:             parametersMap(spec),
		AdditionalProperties:   additionalPropertiesMap(spec),
		IntegrationRuntimeName: integrationRuntimeNamePtr(spec),
		AdbDomain:              pulumi.String(databricks.AdbDomain),
	}

	if databricks.MsiWorkspaceId != "" {
		args.MsiWorkspaceId = pulumi.StringPtr(databricks.MsiWorkspaceId)
	}
	if databricks.AccessToken != "" {
		args.AccessToken = pulumi.StringPtr(databricks.AccessToken)
	}
	if databricks.KeyVaultPassword != nil {
		args.KeyVaultPassword = datafactory.LinkedServiceAzureDatabricksKeyVaultPasswordArgs{
			LinkedServiceName: pulumi.String(databricks.KeyVaultPassword.LinkedServiceName.GetValue()),
			SecretName:        pulumi.String(databricks.KeyVaultPassword.SecretName),
		}
	}
	if databricks.ExistingClusterId != "" {
		args.ExistingClusterId = pulumi.StringPtr(databricks.ExistingClusterId)
	}
	if newCluster := databricks.NewClusterConfig; newCluster != nil {
		clusterArgs := datafactory.LinkedServiceAzureDatabricksNewClusterConfigArgs{
			NodeType:       pulumi.String(newCluster.NodeType),
			ClusterVersion: pulumi.String(newCluster.ClusterVersion),
			// The platform default, sent explicitly (mirrors the
			// provider's own schema default); max 0 means "fixed
			// size" and is omitted.
			MinNumberOfWorkers: pulumi.IntPtr(minWorkersOrDefault(newCluster.MinNumberOfWorkers)),
		}
		if newCluster.MaxNumberOfWorkers > 0 {
			clusterArgs.MaxNumberOfWorkers = pulumi.IntPtr(int(newCluster.MaxNumberOfWorkers))
		}
		if newCluster.DriverNodeType != "" {
			clusterArgs.DriverNodeType = pulumi.StringPtr(newCluster.DriverNodeType)
		}
		if newCluster.LogDestination != "" {
			clusterArgs.LogDestination = pulumi.StringPtr(newCluster.LogDestination)
		}
		if len(newCluster.SparkConfig) > 0 {
			clusterArgs.SparkConfig = pulumi.ToStringMap(newCluster.SparkConfig)
		}
		if len(newCluster.SparkEnvironmentVariables) > 0 {
			clusterArgs.SparkEnvironmentVariables = pulumi.ToStringMap(newCluster.SparkEnvironmentVariables)
		}
		if len(newCluster.CustomTags) > 0 {
			clusterArgs.CustomTags = pulumi.ToStringMap(newCluster.CustomTags)
		}
		if len(newCluster.InitScripts) > 0 {
			clusterArgs.InitScripts = pulumi.ToStringArray(newCluster.InitScripts)
		}
		args.NewClusterConfig = clusterArgs
	}
	if pool := databricks.InstancePool; pool != nil {
		poolArgs := datafactory.LinkedServiceAzureDatabricksInstancePoolArgs{
			InstancePoolId:     pulumi.String(pool.InstancePoolId),
			ClusterVersion:     pulumi.String(pool.ClusterVersion),
			MinNumberOfWorkers: pulumi.IntPtr(minWorkersOrDefault(pool.MinNumberOfWorkers)),
		}
		if pool.MaxNumberOfWorkers > 0 {
			poolArgs.MaxNumberOfWorkers = pulumi.IntPtr(int(pool.MaxNumberOfWorkers))
		}
		args.InstancePool = poolArgs
	}

	created, err := datafactory.NewLinkedServiceAzureDatabricks(ctx, resourceName, args, pulumi.Provider(azureProvider))
	if err != nil {
		return nil, nil, errors.Wrapf(err, "failed to create databricks linked service %s", resourceName)
	}
	return created.ID(), created.Name, nil
}

func minWorkersOrDefault(value *int32) int {
	if value != nil {
		return int(*value)
	}
	return 1
}

func createWeb(
	ctx *pulumi.Context,
	resourceName string,
	spec *azuredatafactorylinkedservicev1alpha1.AzureDataFactoryLinkedServiceSpec,
	azureProvider pulumi.ProviderResource,
) (pulumi.StringInput, pulumi.StringInput, error) {
	web := spec.Web

	args := &datafactory.LinkedServiceWebArgs{
		Name:                   pulumi.String(spec.Name),
		DataFactoryId:          pulumi.String(spec.DataFactoryId.GetValue()),
		Description:            descriptionPtr(spec),
		Annotations:            annotationsArray(spec),
		Parameters:             parametersMap(spec),
		AdditionalProperties:   additionalPropertiesMap(spec),
		IntegrationRuntimeName: integrationRuntimeNamePtr(spec),
		Url:                    pulumi.String(web.Url),
		AuthenticationType:     pulumi.String(web.AuthenticationType),
	}

	if web.Username != "" {
		args.Username = pulumi.StringPtr(web.Username)
	}
	if web.Password != "" {
		args.Password = pulumi.StringPtr(web.Password)
	}

	created, err := datafactory.NewLinkedServiceWeb(ctx, resourceName, args, pulumi.Provider(azureProvider))
	if err != nil {
		return nil, nil, errors.Wrapf(err, "failed to create web linked service %s", resourceName)
	}
	return created.ID(), created.Name, nil
}

func createOdata(
	ctx *pulumi.Context,
	resourceName string,
	spec *azuredatafactorylinkedservicev1alpha1.AzureDataFactoryLinkedServiceSpec,
	azureProvider pulumi.ProviderResource,
) (pulumi.StringInput, pulumi.StringInput, error) {
	odata := spec.Odata

	args := &datafactory.LinkedServiceOdataArgs{
		Name:                   pulumi.String(spec.Name),
		DataFactoryId:          pulumi.String(spec.DataFactoryId.GetValue()),
		Description:            descriptionPtr(spec),
		Annotations:            annotationsArray(spec),
		Parameters:             parametersMap(spec),
		AdditionalProperties:   additionalPropertiesMap(spec),
		IntegrationRuntimeName: integrationRuntimeNamePtr(spec),
		Url:                    pulumi.String(odata.Url),
	}

	if odata.BasicAuthentication != nil {
		args.BasicAuthentication = datafactory.LinkedServiceOdataBasicAuthenticationArgs{
			Username: pulumi.String(odata.BasicAuthentication.Username),
			Password: pulumi.String(odata.BasicAuthentication.Password),
		}
	}

	created, err := datafactory.NewLinkedServiceOdata(ctx, resourceName, args, pulumi.Provider(azureProvider))
	if err != nil {
		return nil, nil, errors.Wrapf(err, "failed to create odata linked service %s", resourceName)
	}
	return created.ID(), created.Name, nil
}

func createOdbc(
	ctx *pulumi.Context,
	resourceName string,
	spec *azuredatafactorylinkedservicev1alpha1.AzureDataFactoryLinkedServiceSpec,
	azureProvider pulumi.ProviderResource,
) (pulumi.StringInput, pulumi.StringInput, error) {
	odbc := spec.Odbc

	args := &datafactory.LinkedServiceOdbcArgs{
		Name:                   pulumi.String(spec.Name),
		DataFactoryId:          pulumi.String(spec.DataFactoryId.GetValue()),
		Description:            descriptionPtr(spec),
		Annotations:            annotationsArray(spec),
		Parameters:             parametersMap(spec),
		AdditionalProperties:   additionalPropertiesMap(spec),
		IntegrationRuntimeName: integrationRuntimeNamePtr(spec),
		ConnectionString:       pulumi.String(odbc.ConnectionString),
	}

	if odbc.BasicAuthentication != nil {
		args.BasicAuthentication = datafactory.LinkedServiceOdbcBasicAuthenticationArgs{
			Username: pulumi.String(odbc.BasicAuthentication.Username),
			Password: pulumi.String(odbc.BasicAuthentication.Password),
		}
	}

	created, err := datafactory.NewLinkedServiceOdbc(ctx, resourceName, args, pulumi.Provider(azureProvider))
	if err != nil {
		return nil, nil, errors.Wrapf(err, "failed to create odbc linked service %s", resourceName)
	}
	return created.ID(), created.Name, nil
}

func createSftp(
	ctx *pulumi.Context,
	resourceName string,
	spec *azuredatafactorylinkedservicev1alpha1.AzureDataFactoryLinkedServiceSpec,
	azureProvider pulumi.ProviderResource,
) (pulumi.StringInput, pulumi.StringInput, error) {
	sftp := spec.Sftp

	args := &datafactory.LinkedServiceSftpArgs{
		Name:                   pulumi.String(spec.Name),
		DataFactoryId:          pulumi.String(spec.DataFactoryId.GetValue()),
		Description:            descriptionPtr(spec),
		Annotations:            annotationsArray(spec),
		Parameters:             parametersMap(spec),
		AdditionalProperties:   additionalPropertiesMap(spec),
		IntegrationRuntimeName: integrationRuntimeNamePtr(spec),
		AuthenticationType:     pulumi.String(sftp.AuthenticationType),
		Host:                   pulumi.String(sftp.Host),
		Port:                   pulumi.Int(int(sftp.Port)),
		Username:               pulumi.String(sftp.Username),
	}

	if sftp.Password != "" {
		args.Password = pulumi.StringPtr(sftp.Password)
	}
	if sftp.KeyVaultPassword != nil {
		args.KeyVaultPasswords = datafactory.LinkedServiceSftpKeyVaultPasswordArray{
			datafactory.LinkedServiceSftpKeyVaultPasswordArgs{
				LinkedServiceName: pulumi.String(sftp.KeyVaultPassword.LinkedServiceName.GetValue()),
				SecretName:        pulumi.String(sftp.KeyVaultPassword.SecretName),
			},
		}
	}
	if sftp.PrivateKeyContentBase64 != "" {
		args.PrivateKeyContentBase64 = pulumi.StringPtr(sftp.PrivateKeyContentBase64)
	}
	if sftp.KeyVaultPrivateKeyContentBase64 != nil {
		args.KeyVaultPrivateKeyContentBase64 = datafactory.LinkedServiceSftpKeyVaultPrivateKeyContentBase64Args{
			LinkedServiceName: pulumi.String(sftp.KeyVaultPrivateKeyContentBase64.LinkedServiceName.GetValue()),
			SecretName:        pulumi.String(sftp.KeyVaultPrivateKeyContentBase64.SecretName),
		}
	}
	if sftp.PrivateKeyPath != "" {
		args.PrivateKeyPath = pulumi.StringPtr(sftp.PrivateKeyPath)
	}
	if sftp.PrivateKeyPassphrase != "" {
		args.PrivateKeyPassphrase = pulumi.StringPtr(sftp.PrivateKeyPassphrase)
	}
	if sftp.KeyVaultPrivateKeyPassphrase != nil {
		args.KeyVaultPrivateKeyPassphrase = datafactory.LinkedServiceSftpKeyVaultPrivateKeyPassphraseArgs{
			LinkedServiceName: pulumi.String(sftp.KeyVaultPrivateKeyPassphrase.LinkedServiceName.GetValue()),
			SecretName:        pulumi.String(sftp.KeyVaultPrivateKeyPassphrase.SecretName),
		}
	}
	if sftp.SkipHostKeyValidation != nil {
		args.SkipHostKeyValidation = pulumi.BoolPtr(sftp.GetSkipHostKeyValidation())
	}
	if sftp.HostKeyFingerprint != "" {
		args.HostKeyFingerprint = pulumi.StringPtr(sftp.HostKeyFingerprint)
	}

	created, err := datafactory.NewLinkedServiceSftp(ctx, resourceName, args, pulumi.Provider(azureProvider))
	if err != nil {
		return nil, nil, errors.Wrapf(err, "failed to create sftp linked service %s", resourceName)
	}
	return created.ID(), created.Name, nil
}

func createCustom(
	ctx *pulumi.Context,
	resourceName string,
	spec *azuredatafactorylinkedservicev1alpha1.AzureDataFactoryLinkedServiceSpec,
	azureProvider pulumi.ProviderResource,
) (pulumi.StringInput, pulumi.StringInput, error) {
	custom := spec.Custom

	args := &datafactory.LinkedCustomServiceArgs{
		Name:                 pulumi.String(spec.Name),
		DataFactoryId:        pulumi.String(spec.DataFactoryId.GetValue()),
		Description:          descriptionPtr(spec),
		Annotations:          annotationsArray(spec),
		Parameters:           parametersMap(spec),
		AdditionalProperties: additionalPropertiesMap(spec),
		Type:                 pulumi.String(custom.Type),
		TypePropertiesJson:   pulumi.String(custom.TypePropertiesJson),
	}

	// The custom resource is the one whose integration runtime travels
	// as a block (name + per-use parameters) instead of a plain name.
	if spec.IntegrationRuntimeName != nil && spec.IntegrationRuntimeName.GetValue() != "" {
		runtimeArgs := datafactory.LinkedCustomServiceIntegrationRuntimeArgs{
			Name: pulumi.String(spec.IntegrationRuntimeName.GetValue()),
		}
		if len(custom.IntegrationRuntimeParameters) > 0 {
			runtimeArgs.Parameters = pulumi.ToStringMap(custom.IntegrationRuntimeParameters)
		}
		args.IntegrationRuntime = runtimeArgs
	}

	created, err := datafactory.NewLinkedCustomService(ctx, resourceName, args, pulumi.Provider(azureProvider))
	if err != nil {
		return nil, nil, errors.Wrapf(err, "failed to create custom linked service %s", resourceName)
	}
	return created.ID(), created.Name, nil
}
