//go:build !codegen
// +build !codegen

package refcheck

import (
	"testing"

	"github.com/plantonhq/planton/shared/cloudresourcekind"
)

// Manifest-authored valueFrom fieldPaths must resolve the way the platform's
// output flattening does -- including MAP outputs addressed by entry key
// (e.g. a load balancer's name-keyed pool ids referenced as
// "status.outputs.backend_pool_ids.web"). A resolver stricter than the
// deploy-time semantics rejects compositions that work; looser accepts
// references that dangle.
func TestResolveValueFromPath(t *testing.T) {
	cases := []struct {
		name      string
		kind      cloudresourcekind.CloudResourceKind
		fieldPath string
		wantOk    bool
	}{
		{
			name:      "plain string output resolves",
			kind:      cloudresourcekind.CloudResourceKind_AzureLoadBalancer,
			fieldPath: "status.outputs.load_balancer_id",
			wantOk:    true,
		},
		{
			name:      "map output addressed by key resolves",
			kind:      cloudresourcekind.CloudResourceKind_AzureLoadBalancer,
			fieldPath: "status.outputs.backend_pool_ids.web",
			wantOk:    true,
		},
		{
			name:      "camelCase map field with key resolves",
			kind:      cloudresourcekind.CloudResourceKind_AzureLoadBalancer,
			fieldPath: "status.outputs.backendPoolIds.web",
			wantOk:    true,
		},
		{
			name:      "map output without a key is rejected",
			kind:      cloudresourcekind.CloudResourceKind_AzureLoadBalancer,
			fieldPath: "status.outputs.backend_pool_ids",
			wantOk:    false,
		},
		{
			name:      "descending past a string map value is rejected",
			kind:      cloudresourcekind.CloudResourceKind_AzureLoadBalancer,
			fieldPath: "status.outputs.backend_pool_ids.web.extra",
			wantOk:    false,
		},
		{
			name:      "unknown field is rejected",
			kind:      cloudresourcekind.CloudResourceKind_AzureLoadBalancer,
			fieldPath: "status.outputs.no_such_output",
			wantOk:    false,
		},
		{
			name:      "non-string terminal is rejected",
			kind:      cloudresourcekind.CloudResourceKind_AzureLoadBalancer,
			fieldPath: "status.outputs",
			wantOk:    false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			reason := ResolveValueFromPath(tc.kind, tc.fieldPath)
			if tc.wantOk && reason != "" {
				t.Errorf("expected %q to resolve, got: %s", tc.fieldPath, reason)
			}
			if !tc.wantOk && reason == "" {
				t.Errorf("expected %q to be rejected, but it resolved", tc.fieldPath)
			}
		})
	}
}
