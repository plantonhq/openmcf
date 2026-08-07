package module

import (
	"github.com/pkg/errors"
	ocifilesystemv1alpha1 "github.com/plantonhq/planton/catalog/oci/ocifilesystem/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/oci/pulumiociprovider"
	"github.com/pulumi/pulumi-oci/sdk/v4/go/oci"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

var accessMap = map[ocifilesystemv1alpha1.OciFileSystemSpec_Access]string{
	ocifilesystemv1alpha1.OciFileSystemSpec_read_write: "READ_WRITE",
	ocifilesystemv1alpha1.OciFileSystemSpec_read_only:  "READ_ONLY",
}

var identitySquashMap = map[ocifilesystemv1alpha1.OciFileSystemSpec_IdentitySquash]string{
	ocifilesystemv1alpha1.OciFileSystemSpec_no_squash:   "NONE",
	ocifilesystemv1alpha1.OciFileSystemSpec_root_squash: "ROOT",
	ocifilesystemv1alpha1.OciFileSystemSpec_all_squash:  "ALL",
}

func Resources(ctx *pulumi.Context, stackInput *ocifilesystemv1alpha1.OciFileSystemStackInput) error {
	locals := initializeLocals(ctx, stackInput)

	ociProvider, err := pulumiociprovider.Get(ctx, stackInput.ProviderConfig)
	if err != nil {
		return errors.Wrap(err, "failed to setup oci provider")
	}

	createdFileSystem, err := fileSystem(ctx, locals, ociProvider)
	if err != nil {
		return errors.Wrap(err, "failed to create file system")
	}

	createdMountTarget, err := mountTarget(ctx, locals, ociProvider)
	if err != nil {
		return errors.Wrap(err, "failed to create mount target")
	}

	if err := exports(ctx, locals, ociProvider, createdFileSystem, createdMountTarget); err != nil {
		return errors.Wrap(err, "failed to create exports")
	}

	return nil
}

func pulumiOciOpt(provider *oci.Provider) pulumi.ResourceOption {
	return pulumi.Provider(provider)
}
