package verify

import "testing"

// The operator's phase Error is transient (any component error for one
// reconcile cycle, requeued and retried); only VersionSupported=False is
// terminal. The wait must keep going through the former and stop at once
// on the latter.
func TestPlatformBootVerdict(t *testing.T) {
	cases := []struct {
		name             string
		phase            string
		versionSupported string
		want             bootVerdict
	}{
		{"ready", "Ready", "True", bootReady},
		{"deploying", "Deploying", "True", bootWaiting},
		{"pending before the first reconcile", "Pending", "", bootWaiting},
		{"a component error is a retry, not a stop", "Error", "True", bootWaiting},
		{"a component error before the condition exists", "Error", "", bootWaiting},
		{"the version floor refuses the declaration", "Error", "False", bootRefused},
		{"a refusal wins over any phase", "Deploying", "False", bootRefused},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := platformBootVerdict(tc.phase, tc.versionSupported); got != tc.want {
				t.Fatalf("platformBootVerdict(%q, %q) = %v, want %v", tc.phase, tc.versionSupported, got, tc.want)
			}
		})
	}
}

// The component trace must read the same whatever order the API server's
// map arrives in, and must not swallow output it cannot parse.
func TestComponentPhaseSummary(t *testing.T) {
	raw := `{"identity":{"message":"Waiting for the identity server","phase":"Deploying"},` +
		`"controlPlane":{"message":"Waiting for dependency: identity","phase":"Pending"},` +
		`"gateway":{"message":"Front door ready","phase":"Ready"}}`
	if got, want := componentPhaseSummary(raw), "controlPlane=Pending gateway=Ready identity=Deploying"; got != want {
		t.Fatalf("componentPhaseSummary = %q, want %q", got, want)
	}
	if got := componentPhaseSummary("  not json  "); got != "not json" {
		t.Fatalf("unparseable input must pass through trimmed, got %q", got)
	}
	if got := componentPhaseSummary(""); got != "" {
		t.Fatalf("empty input must stay empty, got %q", got)
	}
}
