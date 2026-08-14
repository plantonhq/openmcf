package module

import (
	"github.com/pkg/errors"
	azuredatafactorydatasetv1alpha1 "github.com/plantonhq/planton/catalog/azure/azuredatafactorydataset/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/azure/pulumiazureprovider"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// One kind, 13 provider resources: Azure stores every dataset shape
// in the SAME factory-scoped dataset namespace, so the spec's variant
// block selects which resource is created. Shared fields (name,
// factory, the linked service reference, description, annotations,
// parameters, additional_properties, folder) travel identically on
// every shape; each builder adds only its variant's own arguments.
func Resources(ctx *pulumi.Context, stackInput *azuredatafactorydatasetv1alpha1.AzureDataFactoryDatasetStackInput) error {
	locals := initializeLocals(ctx, stackInput)

	// Build the Azure provider from the stack input via the shared builder, which resolves
	// the right credential mechanism (static client secret, keyless web identity, or ambient chain).
	azureProvider, err := pulumiazureprovider.Get(ctx, stackInput.ProviderConfig)
	if err != nil {
		return errors.Wrap(err, "failed to create azure provider")
	}

	spec := locals.AzureDataFactoryDataset.Spec
	resourceName := locals.AzureDataFactoryDataset.Metadata.Name

	var datasetId pulumi.StringInput
	var datasetName pulumi.StringInput

	switch {
	case spec.AzureBlob != nil:
		datasetId, datasetName, err = createAzureBlob(ctx, resourceName, spec, azureProvider)
	case spec.AzureSqlTable != nil:
		datasetId, datasetName, err = createAzureSqlTable(ctx, resourceName, spec, azureProvider)
	case spec.Binary != nil:
		datasetId, datasetName, err = createBinary(ctx, resourceName, spec, azureProvider)
	case spec.CosmosdbSqlapi != nil:
		datasetId, datasetName, err = createCosmosdbSqlapi(ctx, resourceName, spec, azureProvider)
	case spec.Custom != nil:
		datasetId, datasetName, err = createCustom(ctx, resourceName, spec, azureProvider)
	case spec.DelimitedText != nil:
		datasetId, datasetName, err = createDelimitedText(ctx, resourceName, spec, azureProvider)
	case spec.Http != nil:
		datasetId, datasetName, err = createHttp(ctx, resourceName, spec, azureProvider)
	case spec.Json != nil:
		datasetId, datasetName, err = createJson(ctx, resourceName, spec, azureProvider)
	case spec.Mysql != nil:
		datasetId, datasetName, err = createMysql(ctx, resourceName, spec, azureProvider)
	case spec.Parquet != nil:
		datasetId, datasetName, err = createParquet(ctx, resourceName, spec, azureProvider)
	case spec.Postgresql != nil:
		datasetId, datasetName, err = createPostgresql(ctx, resourceName, spec, azureProvider)
	case spec.Snowflake != nil:
		datasetId, datasetName, err = createSnowflake(ctx, resourceName, spec, azureProvider)
	case spec.SqlServerTable != nil:
		datasetId, datasetName, err = createSqlServerTable(ctx, resourceName, spec, azureProvider)
	default:
		// The spec's exactly-one CEL makes this unreachable; the guard
		// keeps a broken input loud instead of silently exporting nothing.
		return errors.New("exactly one dataset variant block must be set")
	}
	if err != nil {
		return err
	}

	ctx.Export(OpDatasetId, datasetId)
	ctx.Export(OpDatasetName, datasetName)

	return nil
}
