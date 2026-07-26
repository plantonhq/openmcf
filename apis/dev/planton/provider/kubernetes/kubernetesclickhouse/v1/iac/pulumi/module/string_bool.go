package module

import (
	"context"
	"reflect"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// crdStringBool renders a SCALAR into a CHI field that crd2pulumi typed as
// an object. The operator's polymorphic StringBool fields (spec.stop,
// clusters[].secret.auto, ...) carry x-kubernetes-preserve-unknown-fields
// in the CRD — structural schemas cannot express the bool|int|string union
// — so the generated SDK exposes them as pulumi.MapInput, which cannot hold
// the "yes"/"true" scalar the operator actually expects on the wire.
//
// This type claims the map element type (so the SDK's input/destination
// type check passes) while holding a plain string; the SDK's marshaler
// serializes non-Output inputs by their dynamic Go kind, so the property
// reaches the provider as the scalar string. Verified against the pinned
// pulumi SDK's rpc marshaling and proven by the module's preview flow.
type crdStringBool string

func (crdStringBool) ElementType() reflect.Type {
	return reflect.TypeOf((*map[string]interface{})(nil)).Elem()
}

func (b crdStringBool) ToMapOutput() pulumi.MapOutput {
	return b.ToMapOutputWithContext(context.Background())
}

// ToMapOutputWithContext exists only to satisfy pulumi.MapInput. It is
// unreachable through resource registration (the marshaler never converts
// an input whose element type already matches the destination) and a map
// output cannot faithfully carry the scalar — fail loudly rather than
// render a config the operator would reject.
func (b crdStringBool) ToMapOutputWithContext(context.Context) pulumi.MapOutput {
	panic("crdStringBool carries a scalar for the CRD's polymorphic StringBool fields and cannot be converted to a MapOutput")
}
