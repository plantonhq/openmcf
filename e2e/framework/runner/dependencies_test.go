package runner

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// writeManifest creates a placeholder manifest file (with parent dirs) under a
// fake repo root and returns its absolute path.
func writeManifest(t *testing.T, repoRoot, relPath string) string {
	t.Helper()
	full := filepath.Join(repoRoot, relPath)
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		t.Fatalf("mkdir for %s: %v", relPath, err)
	}
	if err := os.WriteFile(full, []byte("apiVersion: kubernetes.planton.dev/v1\nkind: Placeholder\n"), 0o600); err != nil {
		t.Fatalf("write %s: %v", relPath, err)
	}
	return full
}

const (
	gwCrdsPrereqRel  = "apis/dev/planton/provider/kubernetes/kubernetesgatewayapicrds/v1/e2e/prerequisite.yaml"
	gwCrdsMinimalRel = "apis/dev/planton/provider/kubernetes/kubernetesgatewayapicrds/v1/e2e/scenarios/minimal.yaml"
)

func TestResolveDependencies_RegistryPrerequisite(t *testing.T) {
	repoRoot := t.TempDir()
	want := writeManifest(t, repoRoot, gwCrdsPrereqRel)

	deps, err := ResolveDependencies(repoRoot, "kubernetes", "kuberneteshttproute")
	if err != nil {
		t.Fatalf("ResolveDependencies: %v", err)
	}
	if len(deps) != 1 {
		t.Fatalf("expected 1 dependency, got %d: %+v", len(deps), deps)
	}
	got := deps[0]
	if got.KindSlug != "kubernetesgatewayapicrds" {
		t.Errorf("kind slug = %q, want kubernetesgatewayapicrds", got.KindSlug)
	}
	if got.ManifestPath != want {
		t.Errorf("manifest path = %q, want %q", got.ManifestPath, want)
	}
}

func TestResolveDependencies_FallbackToMinimalScenario(t *testing.T) {
	repoRoot := t.TempDir()
	want := writeManifest(t, repoRoot, gwCrdsMinimalRel)

	deps, err := ResolveDependencies(repoRoot, "kubernetes", "kuberneteshttproute")
	if err != nil {
		t.Fatalf("ResolveDependencies: %v", err)
	}
	if len(deps) != 1 {
		t.Fatalf("expected 1 dependency, got %d", len(deps))
	}
	if deps[0].ManifestPath != want {
		t.Errorf("manifest path = %q, want fallback %q", deps[0].ManifestPath, want)
	}
}

func TestResolveDependencies_PrerequisiteYamlWinsOverMinimal(t *testing.T) {
	repoRoot := t.TempDir()
	prereq := writeManifest(t, repoRoot, gwCrdsPrereqRel)
	writeManifest(t, repoRoot, gwCrdsMinimalRel)

	deps, err := ResolveDependencies(repoRoot, "kubernetes", "kuberneteshttproute")
	if err != nil {
		t.Fatalf("ResolveDependencies: %v", err)
	}
	if len(deps) != 1 || deps[0].ManifestPath != prereq {
		t.Fatalf("expected prerequisite.yaml to win, got %+v", deps)
	}
}

func TestResolveDependencies_NoPrerequisites(t *testing.T) {
	repoRoot := t.TempDir()
	deps, err := ResolveDependencies(repoRoot, "kubernetes", "kubernetesnamespace")
	if err != nil {
		t.Fatalf("ResolveDependencies: %v", err)
	}
	if len(deps) != 0 {
		t.Fatalf("expected no dependencies, got %d: %+v", len(deps), deps)
	}
}

func TestResolveDependencies_MissingInstallManifestErrors(t *testing.T) {
	repoRoot := t.TempDir()
	// httproute has a registry prereq but we create no install manifest for it.
	if _, err := ResolveDependencies(repoRoot, "kubernetes", "kuberneteshttproute"); err == nil {
		t.Fatal("expected an error when the prerequisite install manifest is missing, got nil")
	}
}

// A single-document install profile passes through untouched -- the common
// case pays no temp-file cost and keeps error messages pointing at the real file.
func TestSplitManifestDocuments_SingleDocumentPassesThrough(t *testing.T) {
	repoRoot := t.TempDir()
	path := writeManifest(t, repoRoot, "single.yaml")

	docs, err := splitManifestDocuments(path)
	if err != nil {
		t.Fatalf("splitManifestDocuments: %v", err)
	}
	if len(docs) != 1 || docs[0] != path {
		t.Fatalf("expected the original path back, got %v", docs)
	}
}

// A multi-document install profile (e.g. the two different-AZ subnets a load
// balancer requires) splits into one deployable file per document, in order,
// with empty documents from leading separators skipped.
func TestSplitManifestDocuments_MultiDocumentSplits(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "prerequisite.yaml")
	content := "---\n" + // leading separator produces an empty doc that must be skipped
		"apiVersion: aws.planton.dev/v1\nkind: AwsSubnet\nmetadata:\n  name: subnet-a\n" +
		"---\n" +
		"apiVersion: aws.planton.dev/v1\nkind: AwsSubnet\nmetadata:\n  name: subnet-b\n"
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write profile: %v", err)
	}

	docs, err := splitManifestDocuments(path)
	if err != nil {
		t.Fatalf("splitManifestDocuments: %v", err)
	}
	if len(docs) != 2 {
		t.Fatalf("expected 2 documents, got %d: %v", len(docs), docs)
	}
	for i, want := range []string{"name: subnet-a", "name: subnet-b"} {
		raw, err := os.ReadFile(docs[i])
		if err != nil {
			t.Fatalf("read split doc %d: %v", i, err)
		}
		if !strings.Contains(string(raw), want) {
			t.Errorf("doc %d missing %q:\n%s", i, want, raw)
		}
	}
}

// TestTeardownDependencies_AggregatesFailures guards the teardown contract:
// one dependency's destroy failure must not stop the remaining teardowns
// (stopping early would leak everything deployed before it), yet every
// failure must surface in the returned error so the run FAILS instead of
// silently leaking cloud resources -- the exact failure mode when an
// ephemeral backend's state disappears before teardown ("no stack named").
func TestTeardownDependencies_AggregatesFailures(t *testing.T) {
	origDestroy, origRemove := pulumiDestroyFn, pulumiRemoveStackFn
	t.Cleanup(func() { pulumiDestroyFn, pulumiRemoveStackFn = origDestroy, origRemove })

	var destroyed []string
	pulumiDestroyFn = func(moduleDir, stackName, backendURL, stackInputFilePath string) (*PulumiResult, error) {
		destroyed = append(destroyed, stackName)
		if stackName == "stack-b" {
			return nil, errors.New("no stack named 'stack-b' found")
		}
		return &PulumiResult{}, nil
	}
	var removed []string
	pulumiRemoveStackFn = func(moduleDir, stackName, backendURL string) error {
		removed = append(removed, stackName)
		return nil
	}

	deployed := []DependencyState{
		{Dependency: Dependency{KindSlug: "awsvpc"}, StackName: "stack-a"},
		{Dependency: Dependency{KindSlug: "awssubnet"}, StackName: "stack-b"},
		{Dependency: Dependency{KindSlug: "awsiamrole"}, StackName: "stack-c"},
	}

	err := TeardownDependencies(deployed)
	if err == nil {
		t.Fatal("expected an aggregated error when a destroy fails, got nil")
	}
	if !strings.Contains(err.Error(), "stack-b") || !strings.Contains(err.Error(), "awssubnet") {
		t.Errorf("aggregated error should identify the failed dependency and stack, got: %v", err)
	}
	// Reverse order, and the failure in the middle must not stop stack-a.
	wantDestroyed := []string{"stack-c", "stack-b", "stack-a"}
	if len(destroyed) != len(wantDestroyed) {
		t.Fatalf("destroyed = %v, want %v", destroyed, wantDestroyed)
	}
	for i := range wantDestroyed {
		if destroyed[i] != wantDestroyed[i] {
			t.Fatalf("destroy order = %v, want %v", destroyed, wantDestroyed)
		}
	}
	// Stack removal runs only for the successful destroys.
	if len(removed) != 2 || removed[0] != "stack-c" || removed[1] != "stack-a" {
		t.Errorf("removed = %v, want [stack-c stack-a]", removed)
	}
}

// A fully clean teardown returns nil so healthy runs keep passing.
func TestTeardownDependencies_AllCleanReturnsNil(t *testing.T) {
	origDestroy, origRemove := pulumiDestroyFn, pulumiRemoveStackFn
	t.Cleanup(func() { pulumiDestroyFn, pulumiRemoveStackFn = origDestroy, origRemove })

	pulumiDestroyFn = func(moduleDir, stackName, backendURL, stackInputFilePath string) (*PulumiResult, error) {
		return &PulumiResult{}, nil
	}
	pulumiRemoveStackFn = func(moduleDir, stackName, backendURL string) error { return nil }

	deployed := []DependencyState{
		{Dependency: Dependency{KindSlug: "awsvpc"}, StackName: "stack-a"},
	}
	if err := TeardownDependencies(deployed); err != nil {
		t.Fatalf("expected nil for a clean teardown, got %v", err)
	}
}

// TestResolveDependencies_TransitiveDeployOrder guards the deep-composition
// ordering DeployDependencies relies on: AwsNatGateway ->
// [AwsSubnet, AwsElasticIp, AwsInternetGateway] with AwsSubnet -> [AwsVpc] and
// AwsInternetGateway -> [AwsVpc] must resolve to
// [AwsVpc, AwsSubnet, AwsElasticIp, AwsInternetGateway], so the VPC's outputs are
// accumulated before the Subnet/IGW whose vpc_id references them (and the IGW is
// attached before the public NAT that requires it).
func TestResolveDependencies_TransitiveDeployOrder(t *testing.T) {
	repoRoot := t.TempDir()
	writeManifest(t, repoRoot, "apis/dev/planton/provider/aws/awsvpc/v1/e2e/prerequisite.yaml")
	writeManifest(t, repoRoot, "apis/dev/planton/provider/aws/awssubnet/v1/e2e/scenarios/minimal.yaml")
	writeManifest(t, repoRoot, "apis/dev/planton/provider/aws/awselasticip/v1/e2e/prerequisite.yaml")
	writeManifest(t, repoRoot, "apis/dev/planton/provider/aws/awsinternetgateway/v1/e2e/prerequisite.yaml")

	deps, err := ResolveDependencies(repoRoot, "aws", "awsnatgateway")
	if err != nil {
		t.Fatalf("ResolveDependencies: %v", err)
	}

	got := make([]string, len(deps))
	for i, d := range deps {
		got[i] = d.KindSlug
	}
	want := []string{"awsvpc", "awssubnet", "awselasticip", "awsinternetgateway"}
	if len(got) != len(want) {
		t.Fatalf("dependency count = %d (%v), want %d (%v)", len(got), got, len(want), want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("dependency order = %v, want %v (VPC must precede Subnet/IGW)", got, want)
		}
	}
}
