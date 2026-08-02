package importmap

import (
	"reflect"
	"testing"
)

func TestSplitAttributePath(t *testing.T) {
	cases := []struct {
		name string
		path string
		want []string
	}{
		{
			name: "plain dotted path",
			path: "spec.update_strategy",
			want: []string{"spec", "update_strategy"},
		},
		{
			name: "three plain segments",
			path: "data.REGISTRY_HTPASSWD",
			want: []string{"data", "REGISTRY_HTPASSWD"},
		},
		{
			name: "bracket-quoted key containing dots",
			path: `data["password.db"]`,
			want: []string{"data", "password.db"},
		},
		{
			name: "bracket-quoted key followed by a plain segment",
			path: `data["a.b"].nested`,
			want: []string{"data", "a.b", "nested"},
		},
		{
			name: "single attribute (no sub-path)",
			path: "data",
			want: []string{"data"},
		},
		{
			name: "unterminated bracket kept verbatim (fails loud downstream)",
			path: `data["password.db`,
			want: []string{"data", `["password.db`},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := SplitAttributePath(tc.path)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("SplitAttributePath(%q) = %v, want %v", tc.path, got, tc.want)
			}
		})
	}
}
