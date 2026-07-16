package module

import (
	"github.com/pkg/errors"
	awsserverlesselasticachev1 "github.com/plantonhq/planton/apis/dev/planton/provider/aws/awsserverlesselasticache/v1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/aws/pulumiawsprovider"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Resources orchestrates creation of the AWS ElastiCache Serverless cache.
// Subnets, security groups, and KMS keys attach by reference; this module
// provisions only the serverless cache resource and exports connection
// endpoints.
func Resources(ctx *pulumi.Context, stackInput *awsserverlesselasticachev1.AwsServerlessElasticacheStackInput) error {
	locals := initializeLocals(ctx, stackInput)

	provider, err := pulumiawsprovider.Get(ctx, stackInput.ProviderConfig, locals.Spec.Region)
	if err != nil {
		return errors.Wrap(err, "failed to create AWS provider")
	}

	if err := serverlessCache(ctx, locals, provider); err != nil {
		return errors.Wrap(err, "serverless cache")
	}

	return nil
}
