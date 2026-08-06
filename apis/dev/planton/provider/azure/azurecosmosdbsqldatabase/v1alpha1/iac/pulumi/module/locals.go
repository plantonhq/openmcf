package module

import (
	azurecosmosdbsqldatabasev1alpha1 "github.com/plantonhq/planton/apis/dev/planton/provider/azure/azurecosmosdbsqldatabase/v1alpha1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AzureCosmosdbSqlDatabase *azurecosmosdbsqldatabasev1alpha1.AzureCosmosdbSqlDatabase
	CosmosdbAccountId        string
}

func initializeLocals(ctx *pulumi.Context, stackInput *azurecosmosdbsqldatabasev1alpha1.AzureCosmosdbSqlDatabaseStackInput) *Locals {
	locals := &Locals{}

	locals.AzureCosmosdbSqlDatabase = stackInput.Target
	locals.CosmosdbAccountId = stackInput.Target.Spec.CosmosdbAccountId.GetValue()

	// No Azure tags: ARM does not support tags on Cosmos child
	// resources, so the platform's identity tags live on the account.

	return locals
}
