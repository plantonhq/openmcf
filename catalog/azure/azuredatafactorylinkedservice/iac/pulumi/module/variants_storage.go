package module

import (
	"github.com/pkg/errors"
	azuredatafactorylinkedservicev1alpha1 "github.com/plantonhq/planton/catalog/azure/azuredatafactorylinkedservice/v1alpha1"
	"github.com/pulumi/pulumi-azure/sdk/v6/go/azure/datafactory"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// The storage-family builders: blob, file share, table, Data Lake
// Gen2, and the two Cosmos DB forms.

func createAzureBlobStorage(
	ctx *pulumi.Context,
	resourceName string,
	spec *azuredatafactorylinkedservicev1alpha1.AzureDataFactoryLinkedServiceSpec,
	azureProvider pulumi.ProviderResource,
) (pulumi.StringInput, pulumi.StringInput, error) {
	blob := spec.AzureBlobStorage

	args := &datafactory.LinkedServiceAzureBlobStorageArgs{
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
	if useManagedIdentityOrDefault(blob.UseManagedIdentity) {
		args.UseManagedIdentity = pulumi.BoolPtr(true)
	}

	if blob.ConnectionString != "" {
		args.ConnectionString = pulumi.StringPtr(blob.ConnectionString)
	}
	if blob.ConnectionStringInsecure != "" {
		args.ConnectionStringInsecure = pulumi.StringPtr(blob.ConnectionStringInsecure)
	}
	if blob.SasUri != "" {
		args.SasUri = pulumi.StringPtr(blob.SasUri)
	}
	if blob.ServiceEndpoint != nil && blob.ServiceEndpoint.GetValue() != "" {
		args.ServiceEndpoint = pulumi.StringPtr(blob.ServiceEndpoint.GetValue())
	}
	if blob.SasTokenLinkedKeyVaultKey != nil {
		args.SasTokenLinkedKeyVaultKey = datafactory.LinkedServiceAzureBlobStorageSasTokenLinkedKeyVaultKeyArgs{
			LinkedServiceName: pulumi.String(blob.SasTokenLinkedKeyVaultKey.LinkedServiceName.GetValue()),
			SecretName:        pulumi.String(blob.SasTokenLinkedKeyVaultKey.SecretName),
		}
	}
	if blob.ServicePrincipalLinkedKeyVaultKey != nil {
		args.ServicePrincipalLinkedKeyVaultKey = datafactory.LinkedServiceAzureBlobStorageServicePrincipalLinkedKeyVaultKeyArgs{
			LinkedServiceName: pulumi.String(blob.ServicePrincipalLinkedKeyVaultKey.LinkedServiceName.GetValue()),
			SecretName:        pulumi.String(blob.ServicePrincipalLinkedKeyVaultKey.SecretName),
		}
	}
	if blob.StorageKind != "" {
		args.StorageKind = pulumi.StringPtr(blob.StorageKind)
	}
	if blob.ServicePrincipalId != "" {
		args.ServicePrincipalId = pulumi.StringPtr(blob.ServicePrincipalId)
	}
	if blob.ServicePrincipalKey != "" {
		args.ServicePrincipalKey = pulumi.StringPtr(blob.ServicePrincipalKey)
	}
	if blob.TenantId != "" {
		args.TenantId = pulumi.StringPtr(blob.TenantId)
	}

	created, err := datafactory.NewLinkedServiceAzureBlobStorage(ctx, resourceName, args, pulumi.Provider(azureProvider))
	if err != nil {
		return nil, nil, errors.Wrapf(err, "failed to create blob storage linked service %s", resourceName)
	}
	return created.ID(), created.Name, nil
}

func createAzureFileStorage(
	ctx *pulumi.Context,
	resourceName string,
	spec *azuredatafactorylinkedservicev1alpha1.AzureDataFactoryLinkedServiceSpec,
	azureProvider pulumi.ProviderResource,
) (pulumi.StringInput, pulumi.StringInput, error) {
	fileStorage := spec.AzureFileStorage

	args := &datafactory.LinkedServiceAzureFileStorageArgs{
		Name:                   pulumi.String(spec.Name),
		DataFactoryId:          pulumi.String(spec.DataFactoryId.GetValue()),
		Description:            descriptionPtr(spec),
		Annotations:            annotationsArray(spec),
		Parameters:             parametersMap(spec),
		AdditionalProperties:   additionalPropertiesMap(spec),
		IntegrationRuntimeName: integrationRuntimeNamePtr(spec),
		ConnectionString:       pulumi.String(fileStorage.ConnectionString),
	}

	if fileStorage.FileShare != "" {
		args.FileShare = pulumi.StringPtr(fileStorage.FileShare)
	}
	if fileStorage.Host != "" {
		args.Host = pulumi.StringPtr(fileStorage.Host)
	}
	if fileStorage.UserId != "" {
		args.UserId = pulumi.StringPtr(fileStorage.UserId)
	}
	if fileStorage.Password != "" {
		args.Password = pulumi.StringPtr(fileStorage.Password)
	}
	if fileStorage.KeyVaultPassword != nil {
		args.KeyVaultPassword = datafactory.LinkedServiceAzureFileStorageKeyVaultPasswordArgs{
			LinkedServiceName: pulumi.String(fileStorage.KeyVaultPassword.LinkedServiceName.GetValue()),
			SecretName:        pulumi.String(fileStorage.KeyVaultPassword.SecretName),
		}
	}

	created, err := datafactory.NewLinkedServiceAzureFileStorage(ctx, resourceName, args, pulumi.Provider(azureProvider))
	if err != nil {
		return nil, nil, errors.Wrapf(err, "failed to create file storage linked service %s", resourceName)
	}
	return created.ID(), created.Name, nil
}

func createAzureTableStorage(
	ctx *pulumi.Context,
	resourceName string,
	spec *azuredatafactorylinkedservicev1alpha1.AzureDataFactoryLinkedServiceSpec,
	azureProvider pulumi.ProviderResource,
) (pulumi.StringInput, pulumi.StringInput, error) {
	args := &datafactory.LinkedServiceAzureTableStorageArgs{
		Name:                   pulumi.String(spec.Name),
		DataFactoryId:          pulumi.String(spec.DataFactoryId.GetValue()),
		Description:            descriptionPtr(spec),
		Annotations:            annotationsArray(spec),
		Parameters:             parametersMap(spec),
		AdditionalProperties:   additionalPropertiesMap(spec),
		IntegrationRuntimeName: integrationRuntimeNamePtr(spec),
		ConnectionString:       pulumi.String(spec.AzureTableStorage.ConnectionString),
	}

	created, err := datafactory.NewLinkedServiceAzureTableStorage(ctx, resourceName, args, pulumi.Provider(azureProvider))
	if err != nil {
		return nil, nil, errors.Wrapf(err, "failed to create table storage linked service %s", resourceName)
	}
	return created.ID(), created.Name, nil
}

func createDataLakeStorageGen2(
	ctx *pulumi.Context,
	resourceName string,
	spec *azuredatafactorylinkedservicev1alpha1.AzureDataFactoryLinkedServiceSpec,
	azureProvider pulumi.ProviderResource,
) (pulumi.StringInput, pulumi.StringInput, error) {
	gen2 := spec.DataLakeStorageGen2

	args := &datafactory.LinkedServiceDataLakeStorageGen2Args{
		Name:                   pulumi.String(spec.Name),
		DataFactoryId:          pulumi.String(spec.DataFactoryId.GetValue()),
		Description:            descriptionPtr(spec),
		Annotations:            annotationsArray(spec),
		Parameters:             parametersMap(spec),
		AdditionalProperties:   additionalPropertiesMap(spec),
		IntegrationRuntimeName: integrationRuntimeNamePtr(spec),
		Url:                    pulumi.String(gen2.Url.GetValue()),
	}

	// Sent only when true -- the provider's AtLeastOneOf group reads
	// an explicit false as "this mode declared", which would collide
	// with the mode actually chosen.
	if useManagedIdentityOrDefault(gen2.UseManagedIdentity) {
		args.UseManagedIdentity = pulumi.BoolPtr(true)
	}
	if gen2.StorageAccountKey != "" {
		args.StorageAccountKey = pulumi.StringPtr(gen2.StorageAccountKey)
	}
	if gen2.ServicePrincipalId != "" {
		args.ServicePrincipalId = pulumi.StringPtr(gen2.ServicePrincipalId)
	}
	if gen2.ServicePrincipalKey != "" {
		args.ServicePrincipalKey = pulumi.StringPtr(gen2.ServicePrincipalKey)
	}
	if gen2.Tenant != "" {
		args.Tenant = pulumi.StringPtr(gen2.Tenant)
	}

	created, err := datafactory.NewLinkedServiceDataLakeStorageGen2(ctx, resourceName, args, pulumi.Provider(azureProvider))
	if err != nil {
		return nil, nil, errors.Wrapf(err, "failed to create data lake gen2 linked service %s", resourceName)
	}
	return created.ID(), created.Name, nil
}

func createCosmosdb(
	ctx *pulumi.Context,
	resourceName string,
	spec *azuredatafactorylinkedservicev1alpha1.AzureDataFactoryLinkedServiceSpec,
	azureProvider pulumi.ProviderResource,
) (pulumi.StringInput, pulumi.StringInput, error) {
	cosmosdb := spec.Cosmosdb

	args := &datafactory.LinkedServiceCosmosDbArgs{
		Name:                   pulumi.String(spec.Name),
		DataFactoryId:          pulumi.String(spec.DataFactoryId.GetValue()),
		Description:            descriptionPtr(spec),
		Annotations:            annotationsArray(spec),
		Parameters:             parametersMap(spec),
		AdditionalProperties:   additionalPropertiesMap(spec),
		IntegrationRuntimeName: integrationRuntimeNamePtr(spec),
	}

	if cosmosdb.ConnectionString != "" {
		args.ConnectionString = pulumi.StringPtr(cosmosdb.ConnectionString)
	}
	if cosmosdb.AccountEndpoint != "" {
		args.AccountEndpoint = pulumi.StringPtr(cosmosdb.AccountEndpoint)
	}
	if cosmosdb.AccountKey != "" {
		args.AccountKey = pulumi.StringPtr(cosmosdb.AccountKey)
	}
	if cosmosdb.Database != "" {
		args.Database = pulumi.StringPtr(cosmosdb.Database)
	}

	created, err := datafactory.NewLinkedServiceCosmosDb(ctx, resourceName, args, pulumi.Provider(azureProvider))
	if err != nil {
		return nil, nil, errors.Wrapf(err, "failed to create cosmosdb linked service %s", resourceName)
	}
	return created.ID(), created.Name, nil
}

func createCosmosdbMongoapi(
	ctx *pulumi.Context,
	resourceName string,
	spec *azuredatafactorylinkedservicev1alpha1.AzureDataFactoryLinkedServiceSpec,
	azureProvider pulumi.ProviderResource,
) (pulumi.StringInput, pulumi.StringInput, error) {
	mongo := spec.CosmosdbMongoapi

	args := &datafactory.LinkedServiceCosmosDbMongoApiArgs{
		Name:                   pulumi.String(spec.Name),
		DataFactoryId:          pulumi.String(spec.DataFactoryId.GetValue()),
		Description:            descriptionPtr(spec),
		Annotations:            annotationsArray(spec),
		Parameters:             parametersMap(spec),
		AdditionalProperties:   additionalPropertiesMap(spec),
		IntegrationRuntimeName: integrationRuntimeNamePtr(spec),
		ConnectionString:       pulumi.StringPtr(mongo.ConnectionString),
		// The platform default, sent explicitly (mirrors the
		// provider's own schema default).
		ServerVersionIs32OrHigher: pulumi.Bool(serverVersion32OrDefault(mongo)),
	}

	if mongo.Database != "" {
		args.Database = pulumi.StringPtr(mongo.Database)
	}

	created, err := datafactory.NewLinkedServiceCosmosDbMongoApi(ctx, resourceName, args, pulumi.Provider(azureProvider))
	if err != nil {
		return nil, nil, errors.Wrapf(err, "failed to create cosmosdb mongo linked service %s", resourceName)
	}
	return created.ID(), created.Name, nil
}

func serverVersion32OrDefault(mongo *azuredatafactorylinkedservicev1alpha1.AzureDataFactoryLinkedServiceCosmosdbMongoapi) bool {
	if mongo.ServerVersionIs_32OrHigher != nil {
		return mongo.GetServerVersionIs_32OrHigher()
	}
	return false
}
