package module

import (
	azurecosmosdbmongocollectionv1 "github.com/plantonhq/planton/apis/dev/planton/provider/azure/azurecosmosdbmongocollection/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AzureCosmosdbMongoCollection *azurecosmosdbmongocollectionv1.AzureCosmosdbMongoCollection
	MongoDatabaseId              string
}

func initializeLocals(ctx *pulumi.Context, stackInput *azurecosmosdbmongocollectionv1.AzureCosmosdbMongoCollectionStackInput) *Locals {
	locals := &Locals{}

	locals.AzureCosmosdbMongoCollection = stackInput.Target
	locals.MongoDatabaseId = stackInput.Target.Spec.MongoDatabaseId.GetValue()

	// No Azure tags: ARM does not support tags on Cosmos child
	// resources, so the platform's identity tags live on the account.

	return locals
}
