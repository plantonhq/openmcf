package module

import (
	gcpprovider "github.com/plantonhq/planton/apis/dev/planton/provider/gcp"
	gcppubsubschemav1alpha1 "github.com/plantonhq/planton/apis/dev/planton/provider/gcp/gcppubsubschema/v1alpha1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	GcpProviderConfig *gcpprovider.GcpProviderConfig
	GcpPubSubSchema   *gcppubsubschemav1alpha1.GcpPubSubSchema
}

func initializeLocals(_ *pulumi.Context, stackInput *gcppubsubschemav1alpha1.GcpPubSubSchemaStackInput) *Locals {
	locals := &Locals{}
	locals.GcpPubSubSchema = stackInput.Target

	// The schema resource has no labels surface in the Pub/Sub API — no
	// platform attribution labels are stamped, identically on both engines.

	locals.GcpProviderConfig = stackInput.ProviderConfig
	return locals
}
