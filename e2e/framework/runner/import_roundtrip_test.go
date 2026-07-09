package runner

import (
	"reflect"
	"testing"
)

// The round-trip oracle's tolerance hinges on this differ naming EXACTLY the
// attributes an in-place update touches: name too few and a real drift slips
// through as "tolerated"; degrade on unexpected shapes and the check must
// fail closed, never tolerate blindly.
func TestChangedTopLevelAttributes(t *testing.T) {
	cases := []struct {
		name   string
		before interface{}
		after  interface{}
		want   []string
	}{
		{
			name:   "single config-only flip (the post-import force_destroy case)",
			before: map[string]interface{}{"force_destroy": false, "bucket": "b", "tags": map[string]interface{}{"a": "1"}},
			after:  map[string]interface{}{"force_destroy": true, "bucket": "b", "tags": map[string]interface{}{"a": "1"}},
			want:   []string{"force_destroy"},
		},
		{
			name:   "nested change surfaces as its top-level attribute",
			before: map[string]interface{}{"tags": map[string]interface{}{"a": "1"}},
			after:  map[string]interface{}{"tags": map[string]interface{}{"a": "2"}},
			want:   []string{"tags"},
		},
		{
			name:   "attribute present on one side only",
			before: map[string]interface{}{"bucket": "b"},
			after:  map[string]interface{}{"bucket": "b", "acl": "private"},
			want:   []string{"acl"},
		},
		{
			name:   "no change",
			before: map[string]interface{}{"bucket": "b"},
			after:  map[string]interface{}{"bucket": "b"},
			want:   nil,
		},
		{
			name:   "non-object shapes fail closed",
			before: nil,
			after:  map[string]interface{}{"bucket": "b"},
			want:   []string{"<non-object change>"},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := changedTopLevelAttributes(tc.before, tc.after)
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("changedTopLevelAttributes = %v, want %v", got, tc.want)
			}
		})
	}
}
