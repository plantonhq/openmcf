package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/kms"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/projects"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func keyRing(ctx *pulumi.Context, locals *Locals, gcpProvider *gcp.Provider) error {
	spec := locals.GcpKmsKeyRing.Spec

	// Enable the Cloud KMS API — the control plane that owns key rings and
	// keys. DisableOnDestroy stays false: tearing down one key ring must
	// never disable the API for everything else in the project (other
	// rings' keys may be actively encrypting production data).
	cloudkmsApiArgs := &projects.ServiceArgs{
		Service:                  pulumi.String("cloudkms.googleapis.com"),
		DisableDependentServices: pulumi.BoolPtr(true),
		DisableOnDestroy:         pulumi.BoolPtr(false),
	}
	if spec.ProjectId.GetValue() != "" {
		cloudkmsApiArgs.Project = pulumi.String(spec.ProjectId.GetValue())
	}
	createdCloudkmsApi, err := projects.NewService(ctx,
		"gcpkr-cloudkms.googleapis.com", cloudkmsApiArgs, pulumi.Provider(gcpProvider))
	if err != nil {
		return errors.Wrap(err, "failed to enable cloudkms.googleapis.com api")
	}

	// The KMS key ring — the permanent organizational container for crypto
	// keys. Every field is immutable, and GCP has no delete API for key
	// rings: on destroy the engine removes the ring from state only,
	// leaving the (free, inert) ring in the project forever. IAM granted
	// on the ring flows down to every key inside it, which makes the ring
	// the blast-radius boundary — group keys by environment or data
	// domain, not one ring per key.
	args := &kms.KeyRingArgs{
		Name:     pulumi.String(spec.KeyRingName),
		Location: pulumi.String(spec.Location),
	}

	// Honor the spec contract: an empty project_id falls back to the
	// provider's default project.
	if spec.ProjectId.GetValue() != "" {
		args.Project = pulumi.StringPtr(spec.ProjectId.GetValue())
	}

	createdKeyRing, err := kms.NewKeyRing(ctx, "key-ring", args,
		pulumi.Provider(gcpProvider),
		pulumi.DependsOn([]pulumi.Resource{createdCloudkmsApi}))
	if err != nil {
		return errors.Wrap(err, "failed to create kms key ring")
	}

	// The resource ID is the fully qualified key ring path
	// (projects/{project}/locations/{location}/keyRings/{name}) — the
	// exact string a GcpKmsKey's key_ring_id reference consumes.
	ctx.Export(OpKeyRingId, createdKeyRing.ID())
	ctx.Export(OpKeyRingName, createdKeyRing.Name)
	// Consumers that take a bare ring name plus a separately supplied
	// location compose from key_ring_name + location.
	ctx.Export(OpLocation, createdKeyRing.Location)

	return nil
}
