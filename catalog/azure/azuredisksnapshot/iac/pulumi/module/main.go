package module

import (
	"github.com/pkg/errors"
	azuredisksnapshotv1alpha1 "github.com/plantonhq/planton/catalog/azure/azuredisksnapshot/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/azure/pulumiazureprovider"
	"github.com/pulumi/pulumi-azure/sdk/v6/go/azure/compute"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func Resources(ctx *pulumi.Context, stackInput *azuredisksnapshotv1alpha1.AzureDiskSnapshotStackInput) error {
	locals := initializeLocals(ctx, stackInput)

	// Build the Azure provider from the stack input via the shared builder, which resolves
	// the right credential mechanism (static client secret, keyless web identity, or ambient chain).
	azureProvider, err := pulumiazureprovider.Get(ctx, stackInput.ProviderConfig)
	if err != nil {
		return errors.Wrap(err, "failed to create azure provider")
	}

	spec := locals.AzureDiskSnapshot.Spec

	// Create the snapshot. The source fields pair with create_option
	// ("Copy" reads source_resource_id; "Import" reads source_uri +
	// storage_account_id) -- the provider's own schema does not tie
	// them together and Azure validates the pairing at create time, so
	// the module sends each source field only when set. Removing
	// encryption settings from a snapshot that had them forces
	// replacement (Azure cannot disable encryption in place).
	args := &compute.SnapshotArgs{
		Name:              pulumi.String(spec.Name),
		ResourceGroupName: pulumi.String(spec.ResourceGroup.GetValue()),
		Location:          pulumi.String(spec.Region),
		CreateOption:      pulumi.String(spec.CreateOption),
		Tags:              pulumi.ToStringMap(locals.AzureTags),
	}

	if spec.SourceResourceId.GetValue() != "" {
		args.SourceResourceId = pulumi.String(spec.SourceResourceId.GetValue())
	}
	if spec.SourceUri != "" {
		args.SourceUri = pulumi.String(spec.SourceUri)
	}
	if spec.StorageAccountId.GetValue() != "" {
		args.StorageAccountId = pulumi.String(spec.StorageAccountId.GetValue())
	}
	if spec.IncrementalEnabled {
		args.IncrementalEnabled = pulumi.Bool(true)
	}
	// Unset inherits the source's size (Azure computes it).
	if spec.DiskSizeGb != nil {
		args.DiskSizeGb = pulumi.Int(int(spec.GetDiskSizeGb()))
	}
	// Unset rides the provider default, "AllowAll".
	if spec.NetworkAccessPolicy != "" {
		args.NetworkAccessPolicy = pulumi.String(spec.NetworkAccessPolicy)
	}
	if spec.DiskAccessId.GetValue() != "" {
		args.DiskAccessId = pulumi.String(spec.DiskAccessId.GetValue())
	}
	// Unset rides the provider default, true.
	if spec.PublicNetworkAccessEnabled != nil {
		args.PublicNetworkAccessEnabled = pulumi.Bool(spec.GetPublicNetworkAccessEnabled())
	}

	if encryption := spec.EncryptionSettings; encryption != nil {
		encryptionArgs := &compute.SnapshotEncryptionSettingsArgs{
			DiskEncryptionKey: &compute.SnapshotEncryptionSettingsDiskEncryptionKeyArgs{
				SecretUrl:     pulumi.String(encryption.DiskEncryptionKey.SecretUrl),
				SourceVaultId: pulumi.String(encryption.DiskEncryptionKey.SourceVaultId.GetValue()),
			},
		}
		if kek := encryption.KeyEncryptionKey; kek != nil {
			encryptionArgs.KeyEncryptionKey = &compute.SnapshotEncryptionSettingsKeyEncryptionKeyArgs{
				KeyUrl:        pulumi.String(kek.KeyUrl),
				SourceVaultId: pulumi.String(kek.SourceVaultId.GetValue()),
			}
		}
		args.EncryptionSettings = encryptionArgs
	}

	createdSnapshot, err := compute.NewSnapshot(ctx,
		locals.AzureDiskSnapshot.Metadata.Name,
		args,
		pulumi.Provider(azureProvider))
	if err != nil {
		return errors.Wrapf(err, "failed to create snapshot %s",
			locals.AzureDiskSnapshot.Metadata.Name)
	}

	ctx.Export(OpSnapshotId, createdSnapshot.ID())
	ctx.Export(OpSnapshotName, createdSnapshot.Name)

	return nil
}
