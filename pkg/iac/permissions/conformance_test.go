package permissions

import (
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"testing"

	permissionsv1 "github.com/plantonhq/planton/iac/componentpermissions/v1"
)

var (
	// awsActionPattern is IAM's "service:Action" spelling. Wildcards are
	// syntactically legal; the gate separately demands a defense in notes.
	awsActionPattern = regexp.MustCompile(`^[a-z0-9-]+:[A-Za-z0-9*]+$`)
	// gcpPermissionPattern is GCP IAM's dotted permission form,
	// e.g. "container.clusters.create" (3+ segments).
	gcpPermissionPattern = regexp.MustCompile(`^[a-z][a-zA-Z0-9]*(\.[a-zA-Z0-9]+){2,}$`)
	// azureActionPattern is Azure RBAC's "Provider.Namespace/type/action"
	// form, e.g. "Microsoft.ContainerService/managedClusters/write".
	azureActionPattern = regexp.MustCompile(`^[A-Za-z]+\.[A-Za-z]+(/[A-Za-z0-9*]+)+$`)
	// kubernetesVerbs is the closed RBAC verb vocabulary. "*" is absent
	// deliberately -- a wildcard verb is never least privilege.
	kubernetesVerbs = map[string]bool{
		"get": true, "list": true, "watch": true, "create": true,
		"update": true, "patch": true, "delete": true, "deletecollection": true,
	}
	// tokenScopedProviders are catalog providers whose modules authenticate
	// with a bearer API token (scoped by the provider's own permission-group
	// or scope vocabulary) and touch NONE of the schema's modeled provider
	// APIs. Per the schema's own absence semantics, each absent section
	// claims "the modules touch no such API" -- literally true across the
	// board for these providers -- so their manifests legally carry no
	// provider section, with the least-privilege token-scope inventory
	// carried as a structured comment header until the schema grows an arm
	// for them. When an arm lands, remove the provider here so the strict
	// at-least-one-section rule resumes holding its manifests.
	tokenScopedProviders = map[string]bool{
		"auth0": true, "cloudflare": true, "digitalocean": true, "openfga": true,
	}
)

// TestPermissionsConformance holds every authored permissions manifest to
// its contract, offline:
//
//  1. The manifest parses strictly against its proto schema and names its
//     component (metadata.name equals the component directory).
//  2. It declares at least one provider section -- a permissions file that
//     grants nothing describes no module.
//  3. Every entry is structurally sound for its provider (action spelling,
//     RBAC verb vocabulary), carries provenance, and defends its trust
//     posture: derived entries cite the module resources they cover, and
//     any wildcard resource scope carries a defense in notes.
//
// What this gate deliberately CANNOT prove: that the actions are SUFFICIENT
// or MINIMAL against a live cloud. That proof is the capture harness's job
// (e2e running under a role built from this very manifest); until an entry
// graduates to proven, its provenance says so honestly.
func TestPermissionsConformance(t *testing.T) {
	root := repoRoot(t)
	if _, err := os.Stat(filepath.Join(root, "catalog")); err != nil {
		t.Skip("catalog source tree not present (bazel sandbox); runs under go test and the lint.catalog-data lane")
	}

	discovered, err := Discover(root)
	if err != nil {
		t.Fatalf("discovering permissions manifests: %v", err)
	}
	if len(discovered) == 0 {
		t.Skip("no permissions manifests authored yet")
	}

	for provider, components := range discovered {
		for _, component := range components {
			component := component
			t.Run(provider+"/"+component, func(t *testing.T) {
				manifest, err := Load(root, provider, component)
				if err != nil {
					t.Fatalf("permissions manifest: %v", err)
				}
				if manifest.GetKind() != "ComponentPermissions" {
					t.Fatalf("kind is %q, want ComponentPermissions", manifest.GetKind())
				}
				if manifest.GetMetadata().GetName() != component {
					t.Errorf("metadata.name is %q, want %q", manifest.GetMetadata().GetName(), component)
				}

				spec := manifest.GetSpec()
				if spec.GetAws() == nil && spec.GetGcp() == nil && spec.GetAzure() == nil && spec.GetKubernetes() == nil {
					if !tokenScopedProviders[provider] {
						t.Fatal("manifest declares no provider section -- a permissions file that grants nothing describes no module")
					}
				}

				checkAws(t, spec.GetAws())
				checkGcp(t, spec.GetGcp())
				checkAzure(t, spec.GetAzure())
				checkKubernetes(t, spec.GetKubernetes())
			})
		}
	}
}

func checkProvenance(t *testing.T, where string, provenance permissionsv1.Provenance, notes string) {
	t.Helper()
	switch provenance {
	case permissionsv1.Provenance_provenance_unspecified:
		t.Errorf("%s: provenance is unspecified -- derived and proven are never blurred", where)
	case permissionsv1.Provenance_derived:
		if strings.TrimSpace(notes) == "" {
			t.Errorf("%s: derived entry cites no module resources in notes -- a derivation without its source is unverifiable", where)
		}
	}
}

func checkAws(t *testing.T, aws *permissionsv1.AwsPermissions) {
	t.Helper()
	if aws == nil {
		return
	}
	if len(aws.GetStatements()) == 0 {
		t.Error("aws section is present but declares no statements")
	}
	sids := map[string]bool{}
	for _, statement := range aws.GetStatements() {
		sid := statement.GetSid()
		if sid == "" {
			t.Error("aws statement with empty sid")
			continue
		}
		if sids[sid] {
			t.Errorf("aws: duplicate sid %q", sid)
		}
		sids[sid] = true
		if len(statement.GetActions()) == 0 {
			t.Errorf("aws %s: no actions", sid)
		}
		for _, action := range statement.GetActions() {
			if !awsActionPattern.MatchString(action) {
				t.Errorf("aws %s: action %q is not service:Action spelling", sid, action)
			}
		}
		if len(statement.GetResources()) == 0 {
			t.Errorf("aws %s: no resources -- scope the statement or defend '*'", sid)
		}
		for _, resource := range statement.GetResources() {
			if resource == "*" && strings.TrimSpace(statement.GetNotes()) == "" {
				t.Errorf("aws %s: resource '*' without a defense in notes -- least privilege demands the reason", sid)
			}
		}
		checkProvenance(t, "aws "+sid, statement.GetProvenance(), statement.GetNotes())
	}
}

func checkGcp(t *testing.T, gcp *permissionsv1.GcpPermissions) {
	t.Helper()
	if gcp == nil {
		return
	}
	if len(gcp.GetGroups()) == 0 {
		t.Error("gcp section is present but declares no groups")
	}
	for _, group := range gcp.GetGroups() {
		if strings.TrimSpace(group.GetPurpose()) == "" {
			t.Error("gcp group with empty purpose")
		}
		if len(group.GetPermissions()) == 0 {
			t.Errorf("gcp %s: no permissions", group.GetPurpose())
		}
		for _, permission := range group.GetPermissions() {
			if !gcpPermissionPattern.MatchString(permission) {
				t.Errorf("gcp %s: permission %q is not dotted IAM form", group.GetPurpose(), permission)
			}
		}
		checkProvenance(t, "gcp "+group.GetPurpose(), group.GetProvenance(), group.GetNotes())
	}
}

func checkAzure(t *testing.T, azure *permissionsv1.AzurePermissions) {
	t.Helper()
	if azure == nil {
		return
	}
	if len(azure.GetGroups()) == 0 {
		t.Error("azure section is present but declares no groups")
	}
	for _, group := range azure.GetGroups() {
		if strings.TrimSpace(group.GetPurpose()) == "" {
			t.Error("azure group with empty purpose")
		}
		if len(group.GetActions()) == 0 && len(group.GetDataActions()) == 0 {
			t.Errorf("azure %s: no actions or data_actions", group.GetPurpose())
		}
		for _, action := range append(append([]string{}, group.GetActions()...), group.GetDataActions()...) {
			if !azureActionPattern.MatchString(action) {
				t.Errorf("azure %s: action %q is not Provider.Namespace/type/action form", group.GetPurpose(), action)
			}
		}
		checkProvenance(t, "azure "+group.GetPurpose(), group.GetProvenance(), group.GetNotes())
	}
}

func checkKubernetes(t *testing.T, kubernetes *permissionsv1.KubernetesPermissions) {
	t.Helper()
	if kubernetes == nil {
		return
	}
	if len(kubernetes.GetRules()) == 0 {
		t.Error("kubernetes section is present but declares no rules")
	}
	for i, rule := range kubernetes.GetRules() {
		if len(rule.GetResources()) == 0 {
			t.Errorf("kubernetes rule %d: no resources", i)
		}
		if len(rule.GetVerbs()) == 0 {
			t.Errorf("kubernetes rule %d: no verbs", i)
		}
		for _, verb := range rule.GetVerbs() {
			if !kubernetesVerbs[verb] {
				t.Errorf("kubernetes rule %d: verb %q is not in the RBAC verb vocabulary (wildcards are never least privilege)", i, verb)
			}
		}
		checkProvenance(t, "kubernetes rule "+ruleLabel(rule), rule.GetProvenance(), rule.GetNotes())
	}
}

func ruleLabel(rule *permissionsv1.KubernetesRule) string {
	return strings.Join(rule.GetResources(), ",")
}

// repoRoot walks up from this test file to the directory containing go.mod.
func repoRoot(t *testing.T) string {
	t.Helper()
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	dir := filepath.Dir(thisFile)
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("go.mod not found walking up from test file")
		}
		dir = parent
	}
}
