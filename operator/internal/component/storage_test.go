package component

import (
	"testing"

	"k8s.io/apimachinery/pkg/api/resource"

	v1 "github.com/plantonhq/planton/operator/api/v1"
)

func platformWithStorage(spec *v1.StorageSpec) *v1.PlantonPlatform {
	return &v1.PlantonPlatform{Spec: v1.PlantonPlatformSpec{Storage: spec}}
}

// The resolution contract every volume rides: component explicit -> global
// spec.storage -> the component's built-in default.
func TestEffectiveStorageSize(t *testing.T) {
	tests := []struct {
		name          string
		planton       *v1.PlantonPlatform
		componentSize string // "" = unset
		want          string
	}{
		{"nothing set anywhere -> component default",
			platformWithStorage(nil), "", "10Gi"},
		{"global only -> global wins over the default",
			platformWithStorage(&v1.StorageSpec{Size: resource.MustParse("800Gi")}), "", "800Gi"},
		{"component beats global",
			platformWithStorage(&v1.StorageSpec{Size: resource.MustParse("800Gi")}), "2Ti", "2Ti"},
		{"component beats default with no global",
			platformWithStorage(nil), "50Gi", "50Gi"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var componentSize resource.Quantity
			if tt.componentSize != "" {
				componentSize = resource.MustParse(tt.componentSize)
			}
			if got := effectiveStorageSize(tt.planton, componentSize, "10Gi"); got != tt.want {
				t.Errorf("effectiveStorageSize = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestEffectiveStorageClass(t *testing.T) {
	tests := []struct {
		name           string
		planton        *v1.PlantonPlatform
		componentClass string
		want           string
	}{
		{"nothing set -> empty (cluster default; callers omit the key)",
			platformWithStorage(nil), "", ""},
		{"global only",
			platformWithStorage(&v1.StorageSpec{StorageClassName: "trident"}), "", "trident"},
		{"component beats global",
			platformWithStorage(&v1.StorageSpec{StorageClassName: "trident"}), "fast-ssd", "fast-ssd"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := effectiveStorageClass(tt.planton, tt.componentClass); got != tt.want {
				t.Errorf("effectiveStorageClass = %q, want %q", got, tt.want)
			}
		})
	}
}
