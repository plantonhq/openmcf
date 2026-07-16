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

func TestResolveDependencies_ConsumerScopedPrerequisiteWins(t *testing.T) {
	repoRoot := t.TempDir()
	consumerPrereq := writeManifest(t, repoRoot,
		"apis/dev/planton/provider/gcp/gcpservicenetworkingconnection/v1/e2e/prerequisites/gcpglobaladdress.yaml")
	writeManifest(t, repoRoot, "apis/dev/planton/provider/gcp/gcpglobaladdress/v1/e2e/prerequisite.yaml")
	writeManifest(t, repoRoot, "apis/dev/planton/provider/gcp/gcpvpcnetwork/v1/e2e/prerequisite.yaml")

	deps, err := ResolveDependencies(repoRoot, "gcp", "gcpservicenetworkingconnection", "")
	if err != nil {
		t.Fatalf("ResolveDependencies: %v", err)
	}
	if len(deps) != 2 {
		t.Fatalf("expected 2 dependencies, got %d: %+v", len(deps), deps)
	}
	got := make([]string, len(deps))
	for i, d := range deps {
		got[i] = d.KindSlug
	}
	want := []string{"gcpvpcnetwork", "gcpglobaladdress"}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("dependency order = %v, want %v", got, want)
		}
	}
	if deps[1].ManifestPath != consumerPrereq {
		t.Errorf("gcpglobaladdress manifest = %q, want consumer override %q", deps[1].ManifestPath, consumerPrereq)
	}
}

func TestResolveDependencies_RegistryPrerequisite(t *testing.T) {
	repoRoot := t.TempDir()
	want := writeManifest(t, repoRoot, gwCrdsPrereqRel)

	deps, err := ResolveDependencies(repoRoot, "kubernetes", "kuberneteshttproute", "")
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

	deps, err := ResolveDependencies(repoRoot, "kubernetes", "kuberneteshttproute", "")
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

	deps, err := ResolveDependencies(repoRoot, "kubernetes", "kuberneteshttproute", "")
	if err != nil {
		t.Fatalf("ResolveDependencies: %v", err)
	}
	if len(deps) != 1 || deps[0].ManifestPath != prereq {
		t.Fatalf("expected prerequisite.yaml to win, got %+v", deps)
	}
}

func TestResolveDependencies_NoPrerequisites(t *testing.T) {
	repoRoot := t.TempDir()
	deps, err := ResolveDependencies(repoRoot, "kubernetes", "kubernetesnamespace", "")
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
	if _, err := ResolveDependencies(repoRoot, "kubernetes", "kuberneteshttproute", ""); err == nil {
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
	origBackoff := dependencyDestroyBackoff
	dependencyDestroyBackoff = 0 // retries must not sleep in unit tests
	t.Cleanup(func() {
		pulumiDestroyFn, pulumiRemoveStackFn = origDestroy, origRemove
		dependencyDestroyBackoff = origBackoff
	})

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
	// Reverse order, the failing destroy retried to its full budget, and the
	// failure in the middle must not stop stack-a.
	wantDestroyed := []string{"stack-c"}
	for i := 0; i < dependencyDestroyAttempts; i++ {
		wantDestroyed = append(wantDestroyed, "stack-b")
	}
	wantDestroyed = append(wantDestroyed, "stack-a")
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

// writeScenarioManifest creates a real, loadable AwsEcsCluster scenario
// manifest carrying the given annotations block (pass "" for none).
func writeScenarioManifest(t *testing.T, dir, annotationsYaml string) string {
	t.Helper()
	content := "apiVersion: aws.planton.dev/v1\n" +
		"kind: AwsEcsCluster\n" +
		"metadata:\n" +
		"  name: scenario-under-test\n" +
		annotationsYaml +
		"spec:\n" +
		"  region: us-west-2\n"
	path := filepath.Join(dir, "scenario.yaml")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write scenario manifest: %v", err)
	}
	return path
}

// TestResolveDependencies_ScenarioDeclaredPrerequisites guards the
// scenario-annotation source: a scenario that composes an optional reference
// (here an auto-scaling group) declares it via the e2e-prerequisites
// annotation, and the harness expands the declared kind through its OWN
// registry prerequisites -- so naming AwsAutoScalingGroup alone yields the
// full VPC -> Subnet -> LaunchTemplate -> ASG chain in deploy order, without
// the component kind carrying a false registry prerequisite.
func TestResolveDependencies_ScenarioDeclaredPrerequisites(t *testing.T) {
	repoRoot := t.TempDir()
	writeManifest(t, repoRoot, "apis/dev/planton/provider/aws/awsvpc/v1/e2e/prerequisite.yaml")
	writeManifest(t, repoRoot, "apis/dev/planton/provider/aws/awssubnet/v1/e2e/prerequisite.yaml")
	writeManifest(t, repoRoot, "apis/dev/planton/provider/aws/awslaunchtemplate/v1/e2e/prerequisite.yaml")
	writeManifest(t, repoRoot, "apis/dev/planton/provider/aws/awsautoscalinggroup/v1/e2e/scenarios/minimal.yaml")

	scenario := writeScenarioManifest(t, t.TempDir(),
		"  annotations:\n"+
			"    planton.dev/e2e-prerequisites: \"AwsAutoScalingGroup\"\n")

	deps, err := ResolveDependencies(repoRoot, "aws", "awsecscluster", scenario)
	if err != nil {
		t.Fatalf("ResolveDependencies: %v", err)
	}

	got := make([]string, len(deps))
	for i, d := range deps {
		got[i] = d.KindSlug
	}
	want := []string{"awsvpc", "awssubnet", "awslaunchtemplate", "awsautoscalinggroup"}
	if len(got) != len(want) {
		t.Fatalf("dependency count = %d (%v), want %d (%v)", len(got), got, len(want), want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("dependency order = %v, want %v", got, want)
		}
	}
}

// A scenario without the annotation resolves exactly as before -- the
// registry graph alone (empty for a Fargate-capable ECS cluster, which is
// honestly a leaf).
func TestResolveDependencies_AnnotationAbsentUsesRegistryOnly(t *testing.T) {
	repoRoot := t.TempDir()
	scenario := writeScenarioManifest(t, t.TempDir(), "")

	deps, err := ResolveDependencies(repoRoot, "aws", "awsecscluster", scenario)
	if err != nil {
		t.Fatalf("ResolveDependencies: %v", err)
	}
	if len(deps) != 0 {
		t.Fatalf("expected no dependencies, got %d: %+v", len(deps), deps)
	}
}

// An unknown kind in the annotation must fail loudly -- silently skipping it
// would deploy the scenario without a dependency it relies on.
func TestResolveDependencies_UnknownAnnotationKindErrors(t *testing.T) {
	repoRoot := t.TempDir()
	scenario := writeScenarioManifest(t, t.TempDir(),
		"  annotations:\n"+
			"    planton.dev/e2e-prerequisites: \"AwsNoSuchKind\"\n")

	if _, err := ResolveDependencies(repoRoot, "aws", "awsecscluster", scenario); err == nil {
		t.Fatal("expected an error for an unknown annotation kind, got nil")
	} else if !strings.Contains(err.Error(), "AwsNoSuchKind") {
		t.Errorf("error should name the unknown kind, got: %v", err)
	}
}

// Declaring the component's own kind as its prerequisite is a modeling
// mistake and must be rejected rather than deploying the kind twice.
func TestResolveDependencies_SelfAnnotationErrors(t *testing.T) {
	repoRoot := t.TempDir()
	scenario := writeScenarioManifest(t, t.TempDir(),
		"  annotations:\n"+
			"    planton.dev/e2e-prerequisites: \"AwsEcsCluster\"\n")

	if _, err := ResolveDependencies(repoRoot, "aws", "awsecscluster", scenario); err == nil {
		t.Fatal("expected an error when a scenario declares its own kind, got nil")
	}
}

// A kind already in the registry graph and also named by the annotation must
// deploy once -- the closure dedupes across the two sources.
func TestResolveDependencies_AnnotationMergesWithRegistry(t *testing.T) {
	repoRoot := t.TempDir()
	writeManifest(t, repoRoot, "apis/dev/planton/provider/aws/awsvpc/v1/e2e/prerequisite.yaml")
	writeManifest(t, repoRoot, "apis/dev/planton/provider/aws/awssubnet/v1/e2e/prerequisite.yaml")
	writeManifest(t, repoRoot, "apis/dev/planton/provider/aws/awselasticip/v1/e2e/prerequisite.yaml")
	writeManifest(t, repoRoot, "apis/dev/planton/provider/aws/awsinternetgateway/v1/e2e/prerequisite.yaml")

	// NAT gateway's registry graph already includes AwsSubnet; the annotation
	// naming it again must not produce a duplicate deployment.
	scenario := writeScenarioManifest(t, t.TempDir(),
		"  annotations:\n"+
			"    planton.dev/e2e-prerequisites: \"AwsSubnet\"\n")

	deps, err := ResolveDependencies(repoRoot, "aws", "awsnatgateway", scenario)
	if err != nil {
		t.Fatalf("ResolveDependencies: %v", err)
	}
	seen := map[string]int{}
	for _, d := range deps {
		seen[d.KindSlug]++
	}
	if seen["awssubnet"] != 1 {
		t.Fatalf("awssubnet deployed %d times, want exactly once: %+v", seen["awssubnet"], deps)
	}
}

// TestResolveDependencies_InstallManifestPrerequisites guards edge source 3:
// a prerequisite whose OWN install manifest composes fixtures of other kinds
// (the zip-backed Lambda function referencing the S3 object-set fixture)
// declares them via its manifest's e2e-prerequisites annotation, and those
// fixtures must deploy BEFORE the declaring kind so its value_from references
// resolve. Here the event-source-mapping's registry prerequisite (AwsLambda)
// installs via a manifest declaring AwsS3ObjectSet, whose registry graph pulls
// in AwsS3Bucket -- so the resolved order must place the bucket and object set
// ahead of the function, with the scenario-declared SQS queue after.
func TestResolveDependencies_InstallManifestPrerequisites(t *testing.T) {
	repoRoot := t.TempDir()
	writeManifest(t, repoRoot, "apis/dev/planton/provider/aws/awsiamrole/v1/e2e/prerequisite.yaml")
	writeManifest(t, repoRoot, "apis/dev/planton/provider/aws/awss3bucket/v1/e2e/prerequisite.yaml")
	writeManifest(t, repoRoot, "apis/dev/planton/provider/aws/awss3objectset/v1/e2e/prerequisite.yaml")
	writeManifest(t, repoRoot, "apis/dev/planton/provider/aws/awssqsqueue/v1/e2e/prerequisite.yaml")

	// The Lambda install manifest is a REAL loadable manifest carrying the
	// annotation -- placeholder manifests are unparseable and contribute no
	// edges by design.
	lambdaManifest := "apiVersion: aws.planton.dev/v1\n" +
		"kind: AwsLambda\n" +
		"metadata:\n" +
		"  name: lambda-fixture\n" +
		"  annotations:\n" +
		"    planton.dev/e2e-prerequisites: \"AwsS3ObjectSet\"\n" +
		"spec:\n" +
		"  region: us-west-2\n"
	lambdaPath := filepath.Join(repoRoot, "apis/dev/planton/provider/aws/awslambda/v1/e2e/scenarios/minimal.yaml")
	if err := os.MkdirAll(filepath.Dir(lambdaPath), 0o755); err != nil {
		t.Fatalf("mkdir lambda scenario dir: %v", err)
	}
	if err := os.WriteFile(lambdaPath, []byte(lambdaManifest), 0o600); err != nil {
		t.Fatalf("write lambda install manifest: %v", err)
	}

	scenario := "apiVersion: aws.planton.dev/v1\n" +
		"kind: AwsLambdaEventSourceMapping\n" +
		"metadata:\n" +
		"  name: esm-under-test\n" +
		"  annotations:\n" +
		"    planton.dev/e2e-prerequisites: \"AwsSqsQueue\"\n" +
		"spec:\n" +
		"  region: us-west-2\n"
	scenarioPath := filepath.Join(t.TempDir(), "scenario.yaml")
	if err := os.WriteFile(scenarioPath, []byte(scenario), 0o600); err != nil {
		t.Fatalf("write scenario manifest: %v", err)
	}

	deps, err := ResolveDependencies(repoRoot, "aws", "awslambdaeventsourcemapping", scenarioPath)
	if err != nil {
		t.Fatalf("ResolveDependencies: %v", err)
	}

	got := make([]string, len(deps))
	for i, d := range deps {
		got[i] = d.KindSlug
	}
	want := []string{"awsiamrole", "awss3bucket", "awss3objectset", "awslambda", "awssqsqueue"}
	if len(got) != len(want) {
		t.Fatalf("dependency count = %d (%v), want %d (%v)", len(got), got, len(want), want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("dependency order = %v, want %v (the S3 chain must precede the function whose manifest references it)", got, want)
		}
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

	deps, err := ResolveDependencies(repoRoot, "aws", "awsnatgateway", "")
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

// writeAzureScenario creates a real, loadable KRM scenario manifest of the
// given Azure kind (optionally carrying the e2e-prerequisites annotation)
// under a fake repo root and returns its absolute path.
func writeAzureScenario(t *testing.T, repoRoot, relPath, kind, prereqs string) string {
	t.Helper()
	content := "apiVersion: azure.planton.dev/v1\nkind: " + kind + "\nmetadata:\n  name: scenario-under-test\n"
	if prereqs != "" {
		content += "  annotations:\n    planton.dev/e2e-prerequisites: \"" + prereqs + "\"\n"
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

// A manifest-path entry deploys as an EXTRA INSTANCE of its declared kind even
// when the registry chain already deploys that kind, and its own transitive
// prerequisites are deduplicated against the chain -- the seam that lets a
// scenario compose a second instance (e.g. a virtual-network peering's remote
// network) without polluting any kind's install profile.
func TestResolveDependencies_PathEntryExtraInstance(t *testing.T) {
	repoRoot := t.TempDir()
	writeManifest(t, repoRoot, "apis/dev/planton/provider/azure/azureresourcegroup/v1/e2e/prerequisite.yaml")
	writeManifest(t, repoRoot, "apis/dev/planton/provider/azure/azurevirtualnetwork/v1/e2e/prerequisite.yaml")
	writeManifest(t, repoRoot, "apis/dev/planton/provider/azure/azureprivatednszone/v1/e2e/prerequisite.yaml")
	remoteRel := "apis/dev/planton/provider/azure/azurevirtualnetworkpeering/v1/e2e/fixtures/remote-network.yaml"
	remoteAbs := writeAzureScenario(t, repoRoot, remoteRel, "AzureVirtualNetwork", "")
	scenario := writeAzureScenario(t, repoRoot, "scenario.yaml", "AzurePrivateDnsZoneVirtualNetworkLink", remoteRel)

	deps, err := ResolveDependencies(repoRoot, "azure", "azureprivatednszonevirtualnetworklink", scenario)
	if err != nil {
		t.Fatalf("ResolveDependencies: %v", err)
	}
	got := make([]string, len(deps))
	for i, d := range deps {
		got[i] = d.KindSlug
	}
	want := []string{"azureresourcegroup", "azureprivatednszone", "azurevirtualnetwork", "azurevirtualnetwork"}
	if strings.Join(got, ",") != strings.Join(want, ",") {
		t.Fatalf("dependency order = %v, want %v", got, want)
	}
	if deps[3].ManifestPath != remoteAbs {
		t.Errorf("extra instance manifest = %q, want %q", deps[3].ManifestPath, remoteAbs)
	}
}

// A path entry pointing at a file that does not exist fails resolution loudly
// instead of deploying a partial fixture chain.
func TestResolveDependencies_PathEntryMissingFileErrors(t *testing.T) {
	repoRoot := t.TempDir()
	writeManifest(t, repoRoot, "apis/dev/planton/provider/azure/azureresourcegroup/v1/e2e/prerequisite.yaml")
	scenario := writeAzureScenario(t, repoRoot, "scenario.yaml", "AzureRouteTable", "apis/does/not/exist.yaml")

	if _, err := ResolveDependencies(repoRoot, "azure", "azureroutetable", scenario); err == nil {
		t.Fatal("expected an error for a missing path entry, got nil")
	}
}

// A manifest-path entry on an INSTALL manifest (not a scenario) is rejected:
// install-manifest edges must be kind names so the graph can dedup them.
func TestResolveDependencies_InstallManifestPathEntryErrors(t *testing.T) {
	repoRoot := t.TempDir()
	writeManifest(t, repoRoot, "apis/dev/planton/provider/azure/azureresourcegroup/v1/e2e/prerequisite.yaml")
	// The virtual network's install profile illegally declares a path entry.
	writeAzureScenario(t, repoRoot,
		"apis/dev/planton/provider/azure/azurevirtualnetwork/v1/e2e/prerequisite.yaml",
		"AzureVirtualNetwork", "apis/some/fixture.yaml")
	scenario := writeAzureScenario(t, repoRoot, "scenario.yaml", "AzureSubnet", "")

	_, err := ResolveDependencies(repoRoot, "azure", "azuresubnet", scenario)
	if err == nil || !strings.Contains(err.Error(), "path entries are only valid on scenario manifests") {
		t.Fatalf("expected the install-manifest path-entry rejection, got %v", err)
	}
}
