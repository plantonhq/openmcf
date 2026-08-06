package runner

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// writeTempManifest is shared with refresolve_test.go.

func TestExpandManifestTokens_ReplacesAllOccurrences(t *testing.T) {
	src := writeTempManifest(t, strings.Join([]string{
		"apiVersion: gcp.planton.dev/v1alpha1",
		"kind: GcpWorkloadIdentityPool",
		"metadata:",
		"  name: fixed-name",
		"spec:",
		"  workloadIdentityPoolId: e2e-pool-${E2E_RUN_ID}",
		"  displayName: pool ${E2E_RUN_ID}",
	}, "\n"))

	out, err := ExpandManifestTokens(src, "ab12cd34")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out == src {
		t.Fatal("expected a new expanded file, got the original path")
	}
	t.Cleanup(func() { os.Remove(out) })

	expanded, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("failed to read expanded manifest: %v", err)
	}
	got := string(expanded)
	if strings.Contains(got, RunIDToken) {
		t.Errorf("expanded manifest still contains %s:\n%s", RunIDToken, got)
	}
	if !strings.Contains(got, "workloadIdentityPoolId: e2e-pool-ab12cd34") {
		t.Errorf("identifier was not expanded:\n%s", got)
	}
	if !strings.Contains(got, "displayName: pool ab12cd34") {
		t.Errorf("second occurrence was not expanded:\n%s", got)
	}

	// The source manifest must be untouched.
	original, err := os.ReadFile(src)
	if err != nil {
		t.Fatalf("failed to re-read source manifest: %v", err)
	}
	if !strings.Contains(string(original), RunIDToken) {
		t.Error("source manifest was modified by expansion")
	}
}

func TestExpandManifestTokens_PassthroughWithoutToken(t *testing.T) {
	src := writeTempManifest(t, "apiVersion: gcp.planton.dev/v1alpha1\nkind: GcpServiceAccount\n")

	out, err := ExpandManifestTokens(src, "ab12cd34")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != src {
		t.Errorf("expected the original path for a token-free manifest, got %s", out)
	}
}

func TestExpandManifestTokens_ErrorsOnTokenWithoutRunID(t *testing.T) {
	src := writeTempManifest(t, "spec:\n  id: e2e-${E2E_RUN_ID}\n")

	if _, err := ExpandManifestTokens(src, ""); err == nil {
		t.Fatal("expected an error when the manifest uses the token but no run id is provided")
	}
}

func TestExpandManifestTokens_KeepsScenarioBasename(t *testing.T) {
	src := writeTempManifest(t, "spec:\n  id: e2e-${E2E_RUN_ID}\n")

	out, err := ExpandManifestTokens(src, "ab12cd34")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	t.Cleanup(func() { os.Remove(out) })

	// Verifier dispatch keys behavioral variants off the scenario name in
	// the manifest path — the expanded copy must preserve the basename.
	if filepath.Base(out) != filepath.Base(src) {
		t.Errorf("expanded copy basename = %q, want %q", filepath.Base(out), filepath.Base(src))
	}
}

func TestExpandManifestTokens_ExpandsEnvTokens(t *testing.T) {
	t.Setenv("PLANTON_E2E_TEST_ROLE_ARN", "arn:aws:iam::123456789012:role/e2e-test")
	src := writeTempManifest(t, strings.Join([]string{
		"spec:",
		"  irsaRoleArn: ${E2E_ENV:PLANTON_E2E_TEST_ROLE_ARN}",
		"  bucket: e2e-${E2E_RUN_ID}",
	}, "\n"))

	out, err := ExpandManifestTokens(src, "ab12cd34")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	t.Cleanup(func() { os.Remove(out) })

	expanded, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("failed to read expanded manifest: %v", err)
	}
	got := string(expanded)
	if !strings.Contains(got, "irsaRoleArn: arn:aws:iam::123456789012:role/e2e-test") {
		t.Errorf("env token was not expanded:\n%s", got)
	}
	if !strings.Contains(got, "bucket: e2e-ab12cd34") {
		t.Errorf("run-id token must still expand alongside env tokens:\n%s", got)
	}
}

func TestExpandManifestTokens_ErrorsOnUnsetEnvToken(t *testing.T) {
	src := writeTempManifest(t, "spec:\n  arn: ${E2E_ENV:PLANTON_E2E_DEFINITELY_UNSET_VAR}\n")

	if _, err := ExpandManifestTokens(src, "ab12cd34"); err == nil {
		t.Fatal("expected an error when an env token's variable is unset")
	}
}

func TestExpandManifestTokens_RejectsNonPrefixedEnvToken(t *testing.T) {
	t.Setenv("AWS_SECRET_ACCESS_KEY", "leak-me")
	src := writeTempManifest(t, "spec:\n  arn: ${E2E_ENV:AWS_SECRET_ACCESS_KEY}\n")

	if _, err := ExpandManifestTokens(src, "ab12cd34"); err == nil {
		t.Fatal("expected an error for env tokens outside the allowed prefix")
	}
}

func TestExpandManifestTokens_EnvTokensWithoutRunID(t *testing.T) {
	t.Setenv("PLANTON_E2E_TEST_VALUE", "some-value")
	src := writeTempManifest(t, "spec:\n  v: ${E2E_ENV:PLANTON_E2E_TEST_VALUE}\n")

	out, err := ExpandManifestTokens(src, "")
	if err != nil {
		t.Fatalf("env-token-only manifests must not require a run id: %v", err)
	}
	t.Cleanup(func() { os.Remove(out) })
}

func TestEngineScopedRunID_DistinctPerEngine(t *testing.T) {
	pulumiID := EngineScopedRunID("ab12cd34", "pulumi")
	terraformID := EngineScopedRunID("ab12cd34", "terraform")
	if pulumiID == terraformID {
		t.Errorf("engine-scoped ids must differ per engine, both were %q", pulumiID)
	}
	if got, want := pulumiID, "ab12cd34-p"; got != want {
		t.Errorf("pulumi id = %q, want %q", got, want)
	}
	if got, want := terraformID, "ab12cd34-t"; got != want {
		t.Errorf("terraform id = %q, want %q", got, want)
	}
	if got := EngineScopedRunID("ab12cd34", ""); got != "ab12cd34" {
		t.Errorf("empty engine should pass the run id through, got %q", got)
	}
}
