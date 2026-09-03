package component

import (
	"k8s.io/apimachinery/pkg/api/resource"

	v1 "github.com/plantonhq/planton/operator/api/v1"
)

// Storage settings resolve identically for every volume the operator creates:
// the component's own setting wins, then the platform-wide spec.storage
// block, then the component's documented default. Defaults live HERE rather
// than as CRD defaults: server-side defaulting would bake sizes into stored
// specs, making "the user chose this" indistinguishable from "auto-filled" --
// which would silently defeat the global spec.storage.size override.

// effectiveStorageSize returns one volume's size as a quantity string.
// componentSize is the component's own spec value (zero when unset);
// defaultSize is the component's built-in default.
func effectiveStorageSize(planton *v1.PlantonPlatform, componentSize resource.Quantity, defaultSize string) string {
	if !componentSize.IsZero() {
		return componentSize.String()
	}
	if planton.Spec.Storage != nil && !planton.Spec.Storage.Size.IsZero() {
		return planton.Spec.Storage.Size.String()
	}
	return defaultSize
}

// effectiveStorageClass returns the StorageClass one volume is pinned to, or
// "" for the cluster default. Callers must OMIT the class entirely when this
// returns "" -- an explicit empty string means "disable dynamic provisioning"
// to the Bitnami chart family, which is never what an unset field intends.
func effectiveStorageClass(planton *v1.PlantonPlatform, componentClass string) string {
	if componentClass != "" {
		return componentClass
	}
	if planton.Spec.Storage != nil {
		return planton.Spec.Storage.StorageClassName
	}
	return ""
}
