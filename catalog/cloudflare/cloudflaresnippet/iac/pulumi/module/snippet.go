package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-cloudflare/sdk/v6/go/cloudflare"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// snippet deploys the snippet's files and entry module. The snippet NAME is the
// identity and Cloudflare's create is an upsert -- deploying a name that already
// exists in the zone silently adopts and overwrites it; renames replace the
// resource. The provider refetches stored content on refresh, so server-side
// normalization of the source can read back as drift -- keep content byte-stable.
func snippet(
	ctx *pulumi.Context,
	locals *Locals,
	cloudflareProvider *cloudflare.Provider,
) error {
	spec := locals.CloudflareSnippet.Spec

	files := make(cloudflare.SnippetFileArray, 0, len(spec.Files))
	for _, file := range spec.Files {
		files = append(files, &cloudflare.SnippetFileArgs{
			Name:    pulumi.String(file.Name),
			Content: pulumi.String(file.Content),
		})
	}

	args := &cloudflare.SnippetArgs{
		ZoneId:      pulumi.String(spec.ZoneId.GetValue()),
		SnippetName: pulumi.String(spec.SnippetName),
		Files:       files,
		Metadata: &cloudflare.SnippetMetadataArgs{
			MainModule: pulumi.String(spec.MainModule),
		},
	}

	createdSnippet, err := cloudflare.NewSnippet(
		ctx,
		"snippet",
		args,
		pulumi.Provider(cloudflareProvider),
	)
	if err != nil {
		return errors.Wrap(err, "failed to deploy snippet")
	}

	ctx.Export(OpSnippetName, createdSnippet.SnippetName)
	ctx.Export(OpZoneId, pulumi.String(spec.ZoneId.GetValue()))

	return nil
}
