package crkreflect

import (
	"sort"

	"github.com/plantonhq/planton/shared/cloudresourcekind"
)

// KindsList returns every registered kind (excluding the unspecified zero
// value), sorted by enum number so every consumer sees the same order on
// every run. The source enum table is a Go map whose iteration order is
// randomized per process — never iterate it directly where order can leak
// into behavior or output.
func KindsList() []cloudresourcekind.CloudResourceKind {
	resp := make([]cloudresourcekind.CloudResourceKind, 0, len(cloudresourcekind.CloudResourceKind_value))
	for _, enumValue := range cloudresourcekind.CloudResourceKind_value {
		if enumValue == 0 {
			continue
		}
		resp = append(resp, cloudresourcekind.CloudResourceKind(enumValue))
	}
	sort.Slice(resp, func(i, j int) bool { return resp[i] < resp[j] })
	return resp
}
