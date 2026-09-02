package protobufyaml

import (
	"testing"
)

// TestYAMLToJSON_Yaml12Semantics locks the deliberate semantics of the
// canonical conversion: YAML 1.1 boolean tokens stay strings, exotic 1.1
// numeric forms stay strings, and only true/false are booleans. Each case
// here is a defect class the YAML 1.1 lineage produced in real manifests.
func TestYAMLToJSON_Yaml12Semantics(t *testing.T) {
	cases := []struct {
		name string
		yaml string
		want string
	}{
		{
			// The key-position mangle: a CloudWatch dashboard widget's `y`
			// coordinate reached AWS as the key "true" under YAML 1.1.
			name: "y stays a string key",
			yaml: "y: 2",
			want: `{"y":2}`,
		},
		{
			// The Norway problem: 1.1 resolved unquoted NO to false.
			name: "NO stays a string value",
			yaml: "country: NO",
			want: `{"country":"NO"}`,
		},
		{
			// 1.1 resolved off/on/yes/no to booleans; a real preset carried
			// `steeringPolicy: off` for an enum whose value name is "off".
			name: "off/on/yes/no stay string values",
			yaml: "a: off\nb: on\nc: yes\nd: no",
			want: `{"a":"off","b":"on","c":"yes","d":"no"}`,
		},
		{
			name: "true and false are the only booleans",
			yaml: "a: true\nb: false\nc: True\nd: FALSE",
			want: `{"a":true,"b":false,"c":true,"d":false}`,
		},
		{
			name: "quoted booleans stay strings",
			yaml: `a: "true"`,
			want: `{"a":"true"}`,
		},
		{
			// 1.1 sexagesimal: `3:00` meant the integer 180. Time-of-day
			// window strings depend on this staying a string.
			name: "sexagesimal-looking values stay strings",
			yaml: "window: 3:00",
			want: `{"window":"3:00"}`,
		},
		{
			// yaml.v3 keeps underscore numerals as integers — identical to
			// the previous converter's behavior, locked here so an upstream
			// resolver change surfaces as a test failure, not a silent shift.
			name: "underscore numerals parse as integers (unchanged behavior)",
			yaml: "a: 1_000",
			want: `{"a":1000}`,
		},
		{
			name: "hex and 1.2 octal parse as integers",
			yaml: "a: 0x1A\nb: 0o17",
			want: `{"a":26,"b":15}`,
		},
		{
			// Dates must never become timestamp objects; protojson expects
			// the string form, and yaml.v3's PLAIN decode would rewrite
			// this to "2026-01-01T00:00:00Z" — the node walk exists to
			// prevent exactly that.
			name: "dates stay the author's exact string",
			yaml: "d: 2026-01-01\nt: 2026-01-01T10:30:00Z",
			want: `{"d":"2026-01-01","t":"2026-01-01T10:30:00Z"}`,
		},
		{
			name: "anchors and aliases resolve",
			yaml: "base: &b\n  region: us-east-1\ncopy: *b",
			want: `{"base":{"region":"us-east-1"},"copy":{"region":"us-east-1"}}`,
		},
		{
			// JSON objects require string keys; YAML 1.2 mapping keys may
			// be numbers or booleans and are stringified.
			name: "non-string keys are stringified",
			yaml: "1: a\ntrue: b",
			want: `{"1":"a","true":"b"}`,
		},
		{
			name: "nested structures convert recursively",
			yaml: "spec:\n  items:\n    - y: 1\n      country: NO\n  enabled: true",
			want: `{"spec":{"enabled":true,"items":[{"country":"NO","y":1}]}}`,
		},
		{
			name: "empty document is null",
			yaml: "",
			want: `null`,
		},
		{
			// The previous converter read only the first document of a
			// multi-document stream; chart-template loaders rely on it.
			name: "multi-document streams read the first document",
			yaml: "a: 1\n---\nb: 2\n",
			want: `{"a":1}`,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := YAMLToJSON([]byte(tc.yaml))
			if err != nil {
				t.Fatalf("YAMLToJSON(%q) errored: %v", tc.yaml, err)
			}
			if string(got) != tc.want {
				t.Fatalf("YAMLToJSON(%q) = %s, want %s", tc.yaml, got, tc.want)
			}
		})
	}
}

// TestYAMLToJSON_DuplicateKeysAreLoud locks the second deliberate semantics
// change: the YAML 1.1 lineage silently kept the LAST duplicate key; the
// canonical conversion refuses the document.
func TestYAMLToJSON_DuplicateKeysAreLoud(t *testing.T) {
	_, err := YAMLToJSON([]byte("a: 1\na: 2"))
	if err == nil {
		t.Fatal("expected duplicate mapping keys to error, got nil")
	}
}

// TestJSONToYAML_RoundTripsThroughYAMLToJSON proves the write direction
// agrees with the read direction: whatever JSONToYAML emits, YAMLToJSON
// reads back byte-identically. This is what the manifest override path
// (load -> override -> write temp file -> re-load) depends on.
func TestJSONToYAML_RoundTripsThroughYAMLToJSON(t *testing.T) {
	// Keys are alphabetical because json.Marshal sorts object keys.
	cases := []string{
		`{"country":"NO","enabled":true,"y":2}`,
		`{"a":"off","b":"true","c":"3:00","d":"2026-01-01"}`,
		`{"spec":{"items":[{"count":3,"name":"x"}],"labels":{"y":"n"}}}`,
	}
	for _, jsonIn := range cases {
		yamlBytes, err := JSONToYAML([]byte(jsonIn))
		if err != nil {
			t.Fatalf("JSONToYAML(%s) errored: %v", jsonIn, err)
		}
		jsonOut, err := YAMLToJSON(yamlBytes)
		if err != nil {
			t.Fatalf("YAMLToJSON of emitted YAML %q errored: %v", yamlBytes, err)
		}
		if string(jsonOut) != jsonIn {
			t.Fatalf("round trip changed the document:\n in: %s\nout: %s\nvia: %s", jsonIn, jsonOut, yamlBytes)
		}
	}
}
