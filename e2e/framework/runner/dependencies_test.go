package runner

import (
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

// writeScenario creates a loadable KRM scenario manifest (real kind, optional
// extra-prerequisites annotation) under a fake repo root and returns its path.
func writeScenario(t *testing.T, repoRoot, relPath, kind, extraPrereqs string) string {
	t.Helper()
	content := "apiVersion: azure.planton.dev/v1\nkind: " + kind + "\nmetadata:\n  name: scenario-under-test\n"
	if extraPrereqs != "" {
		content += "  annotations:\n    planton.dev/e2e-extra-prerequisites: \"" + extraPrereqs + "\"\n"
	}
	full := filepath.Join(repoRoot, relPath)
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		t.Fatalf("mkdir for %s: %v", relPath, err)
	}
	if err := os.WriteFile(full, []byte(content), 0o600); err != nil {
		t.Fatalf("write %s: %v", relPath, err)
	}
	return full
}

// A kind-name entry deploys through its standard install profile, preceded by
// its own transitive prerequisites -- with prerequisites the registry chain
// already schedules deduplicated (here: the resource group is shared).
func TestResolveAllDependencies_ExtraKindEntry(t *testing.T) {
	repoRoot := t.TempDir()
	writeManifest(t, repoRoot, "apis/dev/planton/provider/azure/azureresourcegroup/v1/e2e/prerequisite.yaml")
	pdnsProfile := writeManifest(t, repoRoot, "apis/dev/planton/provider/azure/azureprivatednszone/v1/e2e/prerequisite.yaml")
	scenario := writeScenario(t, repoRoot, "scenario.yaml", "AzureRouteTable", "AzurePrivateDnsZone")

	deps, err := ResolveAllDependencies(repoRoot, "azure", "azureroutetable", scenario)
	if err != nil {
		t.Fatalf("ResolveAllDependencies: %v", err)
	}
	got := make([]string, len(deps))
	for i, d := range deps {
		got[i] = d.KindSlug
	}
	want := []string{"azureresourcegroup", "azureprivatednszone"}
	if len(got) != len(want) {
		t.Fatalf("dependency slugs = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("dependency order = %v, want %v", got, want)
		}
	}
	if deps[1].ManifestPath != pdnsProfile {
		t.Errorf("extra fixture manifest = %q, want the kind's install profile %q", deps[1].ManifestPath, pdnsProfile)
	}
}

// A kind-name entry the registry chain already deploys is skipped entirely.
func TestResolveAllDependencies_ExtraKindAlreadyScheduled(t *testing.T) {
	repoRoot := t.TempDir()
	writeManifest(t, repoRoot, "apis/dev/planton/provider/azure/azureresourcegroup/v1/e2e/prerequisite.yaml")
	scenario := writeScenario(t, repoRoot, "scenario.yaml", "AzureRouteTable", "AzureResourceGroup")

	deps, err := ResolveAllDependencies(repoRoot, "azure", "azureroutetable", scenario)
	if err != nil {
		t.Fatalf("ResolveAllDependencies: %v", err)
	}
	if len(deps) != 1 || deps[0].KindSlug != "azureresourcegroup" {
		t.Fatalf("expected only the registry chain, got %+v", deps)
	}
}

// A path entry deploys as an EXTRA INSTANCE of its declared kind even when the
// registry chain already deploys that kind, and its own transitive
// prerequisites are deduplicated against the chain.
func TestResolveAllDependencies_ExtraPathInstance(t *testing.T) {
	repoRoot := t.TempDir()
	writeManifest(t, repoRoot, "apis/dev/planton/provider/azure/azureresourcegroup/v1/e2e/prerequisite.yaml")
	writeManifest(t, repoRoot, "apis/dev/planton/provider/azure/azurevirtualnetwork/v1/e2e/prerequisite.yaml")
	remoteRel := "apis/dev/planton/provider/azure/azurevirtualnetworkpeering/v1/e2e/fixtures/remote-network.yaml"
	remoteAbs := writeScenario(t, repoRoot, remoteRel, "AzureVirtualNetwork", "")
	scenario := writeScenario(t, repoRoot, "scenario.yaml", "AzureVirtualNetwork", remoteRel)

	// Component under test declares AzureVirtualNetwork in its registry chain
	// (transitively through the DNS-zone link kind, which needs zone + network).
	writeManifest(t, repoRoot, "apis/dev/planton/provider/azure/azureprivatednszone/v1/e2e/prerequisite.yaml")
	deps, err := ResolveAllDependencies(repoRoot, "azure", "azureprivatednszonevirtualnetworklink", scenario)
	if err != nil {
		t.Fatalf("ResolveAllDependencies: %v", err)
	}
	got := make([]string, len(deps))
	for i, d := range deps {
		got[i] = d.KindSlug
	}
	want := []string{"azureresourcegroup", "azureprivatednszone", "azurevirtualnetwork", "azurevirtualnetwork"}
	if len(got) != len(want) {
		t.Fatalf("dependency slugs = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("dependency order = %v, want %v", got, want)
		}
	}
	if deps[3].ManifestPath != remoteAbs {
		t.Errorf("extra instance manifest = %q, want %q", deps[3].ManifestPath, remoteAbs)
	}
}

// An entry that is neither a registered kind nor an existing manifest path
// fails resolution loudly.
func TestResolveAllDependencies_UnknownEntryErrors(t *testing.T) {
	repoRoot := t.TempDir()
	writeManifest(t, repoRoot, "apis/dev/planton/provider/azure/azureresourcegroup/v1/e2e/prerequisite.yaml")
	scenario := writeScenario(t, repoRoot, "scenario.yaml", "AzureRouteTable", "NotARealKind")

	if _, err := ResolveAllDependencies(repoRoot, "azure", "azureroutetable", scenario); err == nil {
		t.Fatal("expected an error for an unknown annotation entry, got nil")
	}
}

// A scenario without the annotation resolves identically to the registry chain.
func TestResolveAllDependencies_NoAnnotation(t *testing.T) {
	repoRoot := t.TempDir()
	writeManifest(t, repoRoot, "apis/dev/planton/provider/azure/azureresourcegroup/v1/e2e/prerequisite.yaml")
	scenario := writeScenario(t, repoRoot, "scenario.yaml", "AzureRouteTable", "")

	deps, err := ResolveAllDependencies(repoRoot, "azure", "azureroutetable", scenario)
	if err != nil {
		t.Fatalf("ResolveAllDependencies: %v", err)
	}
	if len(deps) != 1 || deps[0].KindSlug != "azureresourcegroup" {
		t.Fatalf("expected only the registry chain, got %+v", deps)
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
