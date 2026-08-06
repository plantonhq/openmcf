package manifest

import (
	"strings"
	"testing"
)

// The validation contract: a manifest is checked as a whole document — the
// envelope (apiVersion and kind constants, metadata presence) and the spec
// rules together. The envelope constants are part of each kind's schema; a
// wrong or missing apiVersion is invalid input and must fail with an error
// that names the exact line to fix.

func validateBytes(t *testing.T, manifestYaml string) error {
	t.Helper()
	msg, err := LoadManifestBytes([]byte(manifestYaml), "probe.yaml")
	if err != nil {
		t.Fatalf("manifest failed to load before validation: %v", err)
	}
	return ValidateLoaded(msg)
}

func requireValidationError(t *testing.T, err error, wants ...string) {
	t.Helper()
	if err == nil {
		t.Fatal("expected validation failure")
	}
	for _, want := range wants {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("error missing %q:\n%s", want, err.Error())
		}
	}
}

func TestValidateLoadedValidManifest(t *testing.T) {
	err := validateBytes(t, `apiVersion: aws.planton.dev/v1alpha1
kind: AwsVpc
metadata:
  name: probe
spec:
  region: us-east-1
  cidrBlock: 10.0.0.0/16
`)
	if err != nil {
		t.Fatalf("valid manifest failed validation: %v", err)
	}
}

func TestValidateLoadedRejectsWrongApiVersion(t *testing.T) {
	err := validateBytes(t, `apiVersion: aws.planton.dev/v99
kind: AwsVpc
metadata:
  name: probe
spec:
  region: us-east-1
  cidrBlock: 10.0.0.0/16
`)
	requireValidationError(t, err,
		"apiVersion 'aws.planton.dev/v99' does not match kind AwsVpc",
		"'apiVersion: aws.planton.dev/v1alpha1'")
}

func TestValidateLoadedRejectsMissingApiVersion(t *testing.T) {
	err := validateBytes(t, `kind: AwsVpc
metadata:
  name: probe
spec:
  region: us-east-1
  cidrBlock: 10.0.0.0/16
`)
	requireValidationError(t, err,
		"manifest is missing apiVersion",
		"'apiVersion: aws.planton.dev/v1alpha1'")
}

// The manifest loader resolves kind names tolerantly (case and separators),
// but the document contract requires the canonical name — the same rule the
// server and MCP surfaces enforce. The error names the exact spelling.
func TestValidateLoadedRejectsNonCanonicalKind(t *testing.T) {
	err := validateBytes(t, `apiVersion: aws.planton.dev/v1alpha1
kind: awsvpc
metadata:
  name: probe
spec:
  region: us-east-1
  cidrBlock: 10.0.0.0/16
`)
	requireValidationError(t, err,
		"kind 'awsvpc' is not the canonical name",
		"'kind: AwsVpc'")
}

func TestValidateLoadedRejectsMissingMetadata(t *testing.T) {
	err := validateBytes(t, `apiVersion: aws.planton.dev/v1alpha1
kind: AwsVpc
spec:
  region: us-east-1
  cidrBlock: 10.0.0.0/16
`)
	requireValidationError(t, err, "metadata")
}

// Spec-rule violations must keep flowing through the same report, attributed
// with their full path from the document root.
func TestValidateLoadedStillReportsSpecViolations(t *testing.T) {
	err := validateBytes(t, `apiVersion: aws.planton.dev/v1alpha1
kind: AwsVpc
metadata:
  name: probe
spec:
  cidrBlock: 10.0.0.0/16
`)
	requireValidationError(t, err, "spec.region")
}
