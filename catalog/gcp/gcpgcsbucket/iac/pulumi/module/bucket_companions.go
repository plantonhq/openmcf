package module

import (
	"fmt"
	"sort"
	"strings"

	"github.com/pkg/errors"
	gcpgcsbucketv1alpha1 "github.com/plantonhq/planton/catalog/gcp/gcpgcsbucket/v1alpha1"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/storage"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// bucketCompanions creates the bucket's structural companion resources:
// folders (hierarchical-namespace directories), managed folders
// (prefix-scoped IAM anchors), and Pub/Sub notification configs. All three
// are children of the bucket; resource names key on each entry's identity
// (the folder path, or the notification's list index) — matching the
// Terraform module's for_each keys.
func bucketCompanions(ctx *pulumi.Context, locals *Locals, gcpProvider *gcp.Provider, createdBucket *storage.Bucket) error {
	spec := locals.GcpGcsBucket.Spec

	// Folders. The Storage API does not auto-create missing parents, so
	// creation must run parents-before-children (and destroy runs the
	// reverse, which a non-empty parent needs). Sorting by depth and
	// depending each folder on its nearest managed ancestor gives both
	// orders; the Terraform module expresses the same contract as
	// depth-tiered resource groups chained with depends_on.
	folders := append([]*gcpgcsbucketv1alpha1.GcpGcsBucketFolder{}, spec.Folders...)
	sort.Slice(folders, func(i, j int) bool {
		if d1, d2 := strings.Count(folders[i].Name, "/"), strings.Count(folders[j].Name, "/"); d1 != d2 {
			return d1 < d2
		}
		return folders[i].Name < folders[j].Name
	})
	createdFolders := map[string]*storage.Folder{}
	for _, folder := range folders {
		args := &storage.FolderArgs{
			Bucket: createdBucket.Name,
			Name:   pulumi.String(folder.Name),
		}
		// Per-folder destroy semantics: when true, the provider sweeps
		// every object under the prefix client-side before deleting the
		// folder. False fails the destroy of a non-empty folder instead
		// of erasing data. Always sent — the bucket's own force_destroy
		// idiom.
		args.ForceDestroy = pulumi.BoolPtr(folder.ForceDestroy)
		// The kind-level deletion_policy guards folders exactly like the
		// bucket (PREVENT fails destroy; ABANDON drops from state).
		if spec.DeletionPolicy != "" {
			args.DeletionPolicy = pulumi.StringPtr(spec.DeletionPolicy)
		}
		opts := []pulumi.ResourceOption{
			pulumi.Provider(gcpProvider),
			pulumi.Parent(createdBucket),
		}
		if parent := nearestManagedAncestor(createdFolders, folder.Name); parent != nil {
			opts = append(opts, pulumi.DependsOn([]pulumi.Resource{parent}))
		}
		created, err := storage.NewFolder(ctx, "folder-"+folder.Name, args, opts...)
		if err != nil {
			return errors.Wrapf(err, "failed to create folder %s", folder.Name)
		}
		createdFolders[folder.Name] = created
	}

	// Managed folders — independent prefix anchors, no ordering concerns.
	// Their IAM grants live on the managed-folder IAM surface (composed
	// outside this kind), keeping the same additive-grants doctrine as
	// the bucket's own iam_members.
	for _, managedFolder := range spec.ManagedFolders {
		args := &storage.ManagedFolderArgs{
			Bucket: createdBucket.Name,
			Name:   pulumi.String(managedFolder.Name),
		}
		// Server-side allowNonEmpty: deletion succeeds while objects
		// live under the prefix, and the objects SURVIVE (they simply
		// stop being covered by the managed folder's IAM).
		args.ForceDestroy = pulumi.BoolPtr(managedFolder.ForceDestroy)
		if spec.DeletionPolicy != "" {
			args.DeletionPolicy = pulumi.StringPtr(spec.DeletionPolicy)
		}
		_, err := storage.NewManagedFolder(ctx,
			"managed-folder-"+managedFolder.Name, args,
			pulumi.Provider(gcpProvider),
			pulumi.Parent(createdBucket))
		if err != nil {
			return errors.Wrapf(err, "failed to create managed folder %s", managedFolder.Name)
		}
	}

	// Notification configs. The resource is IMMUTABLE end to end (every
	// argument replaces), so entries key on their list index — reordering
	// the spec list therefore replaces configs, which is exactly the
	// resource's only change mode anyway. The GCS service agent must
	// already hold roles/pubsub.publisher on the topic (composed grant —
	// see the spec's notifications comment); a missing grant surfaces as
	// a create-time API error, not silent misconfiguration.
	for index, notification := range spec.Notifications {
		args := &storage.NotificationArgs{
			Bucket: createdBucket.Name,
			// The API accepts only the fully-qualified
			// projects/{project}/topics/{name} form — the topic_id output
			// of GcpPubSubTopic is exactly this value.
			Topic:         pulumi.String(notification.Topic.GetValue()),
			PayloadFormat: pulumi.String(notification.PayloadFormat),
		}
		// Empty event_types means ALL event types (API default) — send
		// only when narrowed.
		if len(notification.EventTypes) > 0 {
			args.EventTypes = pulumi.ToStringArray(notification.EventTypes)
		}
		if notification.ObjectNamePrefix != "" {
			args.ObjectNamePrefix = pulumi.StringPtr(notification.ObjectNamePrefix)
		}
		if len(notification.CustomAttributes) > 0 {
			args.CustomAttributes = pulumi.ToStringMap(notification.CustomAttributes)
		}
		_, err := storage.NewNotification(ctx,
			fmt.Sprintf("notification-%d", index), args,
			pulumi.Provider(gcpProvider),
			pulumi.Parent(createdBucket))
		if err != nil {
			return errors.Wrapf(err, "failed to create notification %d", index)
		}
	}

	return nil
}

// nearestManagedAncestor walks a folder path upward ("a/b/c/" -> "a/b/" ->
// "a/") and returns the first ancestor this module also manages, so the
// child can depend on it. Ancestors NOT in the spec are a manifest error
// the API reports at create time (parents are never auto-created).
func nearestManagedAncestor(created map[string]*storage.Folder, name string) *storage.Folder {
	current := strings.TrimSuffix(name, "/")
	for {
		idx := strings.LastIndex(current, "/")
		if idx < 0 {
			return nil
		}
		current = current[:idx]
		if folder, ok := created[current+"/"]; ok {
			return folder
		}
	}
}
