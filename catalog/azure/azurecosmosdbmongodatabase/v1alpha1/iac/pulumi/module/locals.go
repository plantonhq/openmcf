package module

import (
	azurecosmosdbmongodatabasev1alpha1 "github.com/plantonhq/planton/catalog/azure/azurecosmosdbmongodatabase/v1alpha1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AzureCosmosdbMongoDatabase *azurecosmosdbmongodatabasev1alpha1.AzureCosmosdbMongoDatabase
	CosmosdbAccountId          string
}

func initializeLocals(ctx *pulumi.Context, stackInput *azurecosmosdbmongodatabasev1alpha1.AzureCosmosdbMongoDatabaseStackInput) *Locals {
	locals := &Locals{}

	locals.AzureCosmosdbMongoDatabase = stackInput.Target
	locals.CosmosdbAccountId = stackInput.Target.Spec.CosmosdbAccountId.GetValue()

	// No Azure tags: ARM does not support tags on Cosmos child
	// resources, so the platform's identity tags live on the account.

	return locals
}
