package module

import (
	"github.com/pkg/errors"
	mongodbatlasv1 "github.com/plantonhq/openmcf/apis/org/openmcf/provider/atlas/mongodbatlas/v1"
	"github.com/pulumi/pulumi-mongodbatlas/sdk/v3/go/mongodbatlas"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Resources creates a MongoDB Atlas cluster with all configured parameters
func Resources(ctx *pulumi.Context, stackInput *mongodbatlasv1.MongodbAtlasStackInput) error {
	locals := initializeLocals(ctx, stackInput)

	// Setup MongoDB Atlas provider with credentials from provider config
	var provider *mongodbatlas.Provider
	var err error
	providerConfig := stackInput.ProviderConfig

	if providerConfig == nil {
		// Use default provider (assumes credentials from environment)
		provider, err = mongodbatlas.NewProvider(ctx, "mongodbatlas-provider", &mongodbatlas.ProviderArgs{})
		if err != nil {
			return errors.Wrap(err, "failed to create default MongoDB Atlas provider")
		}
	} else {
		// Create provider with explicit credentials
		provider, err = mongodbatlas.NewProvider(ctx, "mongodbatlas-provider", &mongodbatlas.ProviderArgs{
			PublicKey:  pulumi.String(providerConfig.PublicKey),
			PrivateKey: pulumi.String(providerConfig.PrivateKey),
		})
		if err != nil {
			return errors.Wrap(err, "failed to create MongoDB Atlas provider with credentials")
		}
	}

	// Create the MongoDB Atlas cluster
	createdCluster, err := createCluster(ctx, locals, provider)
	if err != nil {
		return errors.Wrap(err, "failed to create MongoDB Atlas cluster")
	}

	// Export stack outputs
	return exportOutputs(ctx, createdCluster, locals)
}
