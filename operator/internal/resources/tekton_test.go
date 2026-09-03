package resources

import "testing"

// The sink URL is a two-process contract: Tekton posts here, the runner's
// webhook mux serves exactly this path. A bare Service root would 404 every
// event and builds would silently degrade to the reconciliation safety net --
// the path is as load-bearing as the host.
func TestTektonEventsSinkURL(t *testing.T) {
	got := TektonEventsSinkURL("planton", "default")
	want := "http://planton-runner.default.svc.cluster.local/service-hub/tekton/cloud-event"
	if got != want {
		t.Errorf("sink URL = %q, want %q", got, want)
	}
}

// The fragment carries ONLY the sink key: applied with its own field manager,
// the operator owns exactly this key on Tekton's config-defaults and nothing
// else -- a Tekton version bump re-applying the full ConfigMap under the
// install's manager can never claw it back.
func TestTektonConfigDefaultsSinkFragment(t *testing.T) {
	fragment := TektonConfigDefaultsSinkFragment("http://sink")
	if fragment.Name != "config-defaults" || fragment.Namespace != TektonPipelinesNamespace {
		t.Errorf("target = %s/%s, want tekton-pipelines/config-defaults", fragment.Namespace, fragment.Name)
	}
	if len(fragment.Data) != 1 || fragment.Data["default-cloud-events-sink"] != "http://sink" {
		t.Errorf("data = %v, want exactly the default-cloud-events-sink key", fragment.Data)
	}
	if len(fragment.OwnerReferences) != 0 {
		t.Error("the fragment must carry no owner reference -- config-defaults is Tekton's object")
	}
}

// The two field managers must never collide: shared ownership of
// config-defaults between the install apply and the sink write is the exact
// mechanism by which a version bump would erase the sink.
func TestTektonEventsSinkFieldManagerIsDistinct(t *testing.T) {
	if TektonEventsSinkFieldManager == "planton-operator" {
		t.Fatal("the sink write must not share the install apply's field manager")
	}
}

// Without this grant the runner's readiness probe can only report "could not
// determine" for the events sink -- scoped to the single ConfigMap the probe
// reads.
func TestTektonEventsSinkReadRBAC(t *testing.T) {
	role := TektonEventsSinkReadRole("planton", "default")
	if role.Namespace != TektonPipelinesNamespace {
		t.Errorf("role namespace = %s, want tekton-pipelines (where config-defaults lives)", role.Namespace)
	}
	// The pair lives in Tekton's SHARED namespace, so its name must carry
	// the platform's namespace -- two same-named platforms must never share
	// a binding whose subject each would force-apply to its own runner.
	if role.Name != "default-planton-runner-events-sink-read" {
		t.Errorf("role name = %s, want the namespace-qualified name", role.Name)
	}
	if len(role.Rules) != 1 {
		t.Fatalf("rules = %+v, want exactly one", role.Rules)
	}
	rule := role.Rules[0]
	if len(rule.ResourceNames) != 1 || rule.ResourceNames[0] != "config-defaults" ||
		len(rule.Verbs) != 1 || rule.Verbs[0] != "get" {
		t.Errorf("rule = %+v, want get on exactly config-defaults", rule)
	}
	if len(role.OwnerReferences) != 0 {
		t.Error("no owner reference -- the Role lives outside the CR's namespace")
	}

	binding := TektonEventsSinkReadRoleBinding("planton", "default")
	if binding.Namespace != TektonPipelinesNamespace {
		t.Errorf("binding namespace = %s, want tekton-pipelines", binding.Namespace)
	}
	if len(binding.Subjects) != 1 || binding.Subjects[0].Name != "planton-runner" ||
		binding.Subjects[0].Namespace != "default" {
		t.Errorf("subjects = %+v, want the runner ServiceAccount in the CR namespace", binding.Subjects)
	}
}
