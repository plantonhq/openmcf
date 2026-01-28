package module

import (
	"fmt"

	"github.com/pkg/errors"
	cloudflareworkerv1 "github.com/plantonhq/openmcf/apis/org/openmcf/provider/cloudflare/cloudflareworker/v1"
	"github.com/plantonhq/openmcf/pkg/iac/pulumi/pulumimodule/provider/cloudflare/pulumicloudflareprovider"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Resources is the entry‑point expected by the OpenMCF CLI.
func Resources(
	ctx *pulumi.Context,
	stackInput *cloudflareworkerv1.CloudflareWorkerStackInput,
) error {
	// 1. Prepare handy references.
	locals := initializeLocals(ctx, stackInput)

	// 2. Stand‑up a Cloudflare provider from the supplied credential.
	createdProvider, err := pulumicloudflareprovider.Get(
		ctx,
		locals.CloudflareProviderConfig,
	)
	if err != nil {
		return errors.Wrap(err, "failed to set up cloudflare provider")
	}

	// 3. Create AWS provider for R2 (with explicit creds or default from env)
	var r2Provider *aws.Provider
	if locals.CloudflareProviderConfig != nil && locals.CloudflareProviderConfig.R2 != nil {
		// Use explicit R2 credentials
		r2Creds := locals.CloudflareProviderConfig.R2

		// Determine R2 endpoint (use custom if provided, otherwise derive from account ID)
		var endpoint string
		if r2Creds.Endpoint != "" {
			endpoint = r2Creds.Endpoint
		} else {
			endpoint = fmt.Sprintf("https://%s.r2.cloudflarestorage.com", locals.CloudflareWorker.Spec.AccountId)
		}

		r2Provider, err = aws.NewProvider(ctx, "r2-provider", &aws.ProviderArgs{
			Region:                    pulumi.String("auto"),
			AccessKey:                 pulumi.String(r2Creds.AccessKeyId),
			SecretKey:                 pulumi.String(r2Creds.SecretAccessKey),
			S3UsePathStyle:            pulumi.Bool(true),
			SkipCredentialsValidation: pulumi.Bool(true),
			SkipMetadataApiCheck:      pulumi.Bool(true),
			SkipRegionValidation:      pulumi.Bool(true),
			SkipRequestingAccountId:   pulumi.Bool(true),
			Endpoints: aws.ProviderEndpointArray{
				aws.ProviderEndpointArgs{
					S3: pulumi.String(endpoint),
				},
			},
		})
	} else {
		// No explicit credentials - use AWS env vars but configure for R2
		endpoint := fmt.Sprintf("https://%s.r2.cloudflarestorage.com", locals.CloudflareWorker.Spec.AccountId)

		r2Provider, err = aws.NewProvider(ctx, "r2-provider", &aws.ProviderArgs{
			Region:                    pulumi.String("auto"),
			S3UsePathStyle:            pulumi.Bool(true),
			SkipCredentialsValidation: pulumi.Bool(true),
			SkipMetadataApiCheck:      pulumi.Bool(true),
			SkipRegionValidation:      pulumi.Bool(true),
			SkipRequestingAccountId:   pulumi.Bool(true),
			Endpoints: aws.ProviderEndpointArray{
				aws.ProviderEndpointArgs{
					S3: pulumi.String(endpoint),
				},
			},
		})
	}
	if err != nil {
		return errors.Wrap(err, "failed to create R2 provider")
	}

	// 4. Create the Worker script with content from R2 bundle.
	createdWorkerScript, err := createWorkerScript(ctx, locals, createdProvider, r2Provider)
	if err != nil {
		return errors.Wrap(err, "failed to create worker script")
	}

	// 5. Secrets upload - REMOVED
	// TODO: Implement secrets management as a separate feature in the future
	// Secrets should be managed via Cloudflare Workers Secrets API or separate resource

	// 6. Optionally attach to a route.
	// TODO: Temporarily commented out to test worker script deployment
	// The API token needs "Workers Routes: Edit" permission for this to work
	// createdWorkerScript variable would be used here when route creation is re-enabled
	if _, err := route(ctx, locals, createdProvider, createdWorkerScript); err != nil {
		return errors.Wrap(err, "failed to create cloudflare worker route")
	}

	return nil
}
