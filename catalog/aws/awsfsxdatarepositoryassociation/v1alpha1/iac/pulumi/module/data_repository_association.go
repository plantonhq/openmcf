package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/fsx"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// dataRepositoryAssociation creates the FSx data repository association. The
// association's identity is the (file system, path, bucket) triple — all
// three are ForceNew; the sync policies update in place.
func dataRepositoryAssociation(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) (*fsx.DataRepositoryAssociation, error) {
	spec := locals.AwsFsxDataRepositoryAssociation.Spec
	name := locals.AwsFsxDataRepositoryAssociation.Metadata.Name

	args := &fsx.DataRepositoryAssociationArgs{
		FileSystemId:       pulumi.String(spec.FileSystemId.GetValue()),
		FileSystemPath:     pulumi.String(spec.FileSystemPath),
		DataRepositoryPath: pulumi.String(spec.DataRepositoryPath),
		Tags:               pulumi.ToStringMap(locals.AwsTags),
	}

	// Bidirectional sync policies. The provider wraps them in an `s3` block;
	// the spec exposes the two event lists directly (the wrapper carries no
	// information of its own), so the block is built only when at least one
	// policy has events.
	if len(spec.AutoImportEvents) > 0 || len(spec.AutoExportEvents) > 0 {
		s3Args := &fsx.DataRepositoryAssociationS3Args{}
		if len(spec.AutoImportEvents) > 0 {
			s3Args.AutoImportPolicy = &fsx.DataRepositoryAssociationS3AutoImportPolicyArgs{
				Events: pulumi.ToStringArray(spec.AutoImportEvents),
			}
		}
		if len(spec.AutoExportEvents) > 0 {
			s3Args.AutoExportPolicy = &fsx.DataRepositoryAssociationS3AutoExportPolicyArgs{
				Events: pulumi.ToStringArray(spec.AutoExportEvents),
			}
		}
		args.S3 = s3Args
	}

	// Stripe size for imported files (omitted keeps the AWS default of 1024
	// MiB).
	if spec.ImportedFileChunkSize != nil {
		args.ImportedFileChunkSize = pulumi.IntPtr(int(spec.GetImportedFileChunkSize()))
	}

	// Create-time batch import of the existing S3 metadata — without it,
	// only objects changing AFTER creation appear in the namespace.
	if spec.BatchImportMetaDataOnCreate {
		args.BatchImportMetaDataOnCreate = pulumi.BoolPtr(true)
	}

	// Delete-time cascade: remove the linked files from the file system when
	// the association is deleted (default keeps them).
	if spec.DeleteDataInFilesystem {
		args.DeleteDataInFilesystem = pulumi.BoolPtr(true)
	}

	createdAssociation, err := fsx.NewDataRepositoryAssociation(ctx, name, args, pulumi.Provider(provider))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create fsx data repository association")
	}

	return createdAssociation, nil
}
