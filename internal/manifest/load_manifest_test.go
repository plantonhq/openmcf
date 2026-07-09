package manifest

import (
	"strings"
	"testing"
)

// TestLoadManifestBytesDiagnosesFailures locks the load-error contract: when
// the manifest does not fit its schema, the error itself -- not just a styled
// display path -- names the real YAML line, the field path, and the fix.
func TestLoadManifestBytesDiagnosesFailures(t *testing.T) {
	manifest := `apiVersion: aws.planton.dev/v1
kind: AwsSecurityGroup
metadata:
  name: probe
spec:
  region: us-east-1
  vpcId: vpc-12345
  description: probe group
`
	_, err := LoadManifestBytes([]byte(manifest), "probe.yaml")
	if err == nil {
		t.Fatal("expected load failure")
	}
	mle, ok := err.(*ManifestLoadError)
	if !ok {
		t.Fatalf("expected ManifestLoadError, got %T: %v", err, err)
	}
	if mle.Kind != "AwsSecurityGroup" || len(mle.Mismatches) == 0 {
		t.Fatalf("diagnosis missing: kind=%q mismatches=%d", mle.Kind, len(mle.Mismatches))
	}
	text := mle.Error()
	for _, want := range []string{"spec.vpcId", "(line 7)", "write as {value:", "planton explain AwsSecurityGroup.spec.vpcId"} {
		if !strings.Contains(text, want) {
			t.Errorf("Error() missing %q:\n%s", want, text)
		}
	}
	// The parser's original error stays reachable for callers that need it.
	if mle.Err == nil {
		t.Error("original parser error must be preserved")
	}
}

// TestLoadManifestBytesKeepsParserErrorWhenUndiagnosed: value-format
// failures inside legal structure belong to protojson; the error text must
// remain the parser's, never a guess.
func TestLoadManifestBytesKeepsParserErrorWhenUndiagnosed(t *testing.T) {
	manifest := `apiVersion: aws.planton.dev/v1
kind: AwsVpc
metadata:
  name: probe
spec:
  region: us-east-1
  ipv4NetmaskLength: abc
`
	_, err := LoadManifestBytes([]byte(manifest), "probe.yaml")
	if err == nil {
		t.Fatal("expected load failure")
	}
	mle, ok := err.(*ManifestLoadError)
	if !ok {
		t.Fatalf("expected ManifestLoadError, got %T", err)
	}
	if len(mle.Mismatches) != 0 {
		t.Fatalf("coercion failure must not be diagnosed: %+v", mle.Mismatches)
	}
	if mle.Error() != mle.Err.Error() {
		t.Error("undiagnosed error text must be the parser's own")
	}
}

func TestLoadManifestBytesValid(t *testing.T) {
	manifest := `apiVersion: aws.planton.dev/v1
kind: AwsVpc
metadata:
  name: probe
spec:
  region: us-east-1
  cidrBlock: 10.0.0.0/16
`
	msg, err := LoadManifestBytes([]byte(manifest), "probe.yaml")
	if err != nil {
		t.Fatalf("valid manifest failed: %v", err)
	}
	if msg == nil {
		t.Fatal("nil message for valid manifest")
	}
}
