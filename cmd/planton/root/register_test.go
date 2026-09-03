package root

import (
	"testing"

	"github.com/plantonhq/planton/internal/cli/version"
	"github.com/spf13/cobra"
)

// The engine set is the embedding contract: every user-facing engine command
// must be present, and binary self-management (version/upgrade/downgrade) and
// developer tools (e2e, provider-parity) must not be.
func TestRegisterCommands_EngineSet(t *testing.T) {
	parent := &cobra.Command{Use: "host"}
	RegisterCommands(parent, Options{})

	want := []string{
		"apply", "checkout", "destroy", "init", "kustomize", "load-manifest",
		"module", "modules-version", "plan", "pull", "pulumi", "refresh",
		"secret-coverage", "terraform", "tofu", "validate-manifest",
		"validate-outputs", "validate-refs",
	}
	got := map[string]bool{}
	for _, c := range parent.Commands() {
		got[c.Name()] = true
	}
	for _, name := range want {
		if !got[name] {
			t.Errorf("engine command %q not registered", name)
		}
	}
	for _, excluded := range []string{"version", "upgrade", "downgrade", "e2e", "provider-parity"} {
		if got[excluded] {
			t.Errorf("command %q must not be part of the engine set", excluded)
		}
	}
}

func TestRegisterCommands_PersistentFlags(t *testing.T) {
	parent := &cobra.Command{Use: "host"}
	RegisterCommands(parent, Options{})

	if f := parent.PersistentFlags().Lookup("local-module"); f == nil {
		t.Error("persistent flag --local-module not registered")
	}
	f := parent.PersistentFlags().Lookup("planton-git-repo")
	if f == nil {
		t.Fatal("persistent flag --planton-git-repo not registered")
	}
	if f.DefValue != DefaultPlantonGitRepo {
		t.Errorf("--planton-git-repo default = %q, want %q", f.DefValue, DefaultPlantonGitRepo)
	}
}

func TestSetModulesVersion(t *testing.T) {
	original := version.Version
	defer func() { version.Version = original }()

	version.Version = ""
	SetModulesVersion("v0.3.0")
	if version.Version != "v0.3.0" {
		t.Errorf("version = %q, want v0.3.0", version.Version)
	}

	// Empty input must never erase a stamped version.
	SetModulesVersion("")
	if version.Version != "v0.3.0" {
		t.Errorf("empty SetModulesVersion overwrote stamped version: %q", version.Version)
	}
}

// Every engine command that RUNS an IaC module reads --module-dir with an
// EMPTY default: empty means "no explicit choice", and the module runtime
// then probes the current directory before falling through to the published
// module for this release. A group that registered a required or defaulted
// --module-dir would break that contract for every command under it, and a
// handler that refused an empty value would contradict the flag's own help.
// (The module-authoring tools, `module verify` and `validate-outputs`,
// inspect a directory the author names and are rightly required.)
func TestEngineCommands_ModuleDirIsOptionalEverywhere(t *testing.T) {
	parent := &cobra.Command{Use: "host"}
	RegisterCommands(parent, Options{})

	runsAModule := map[string]bool{
		"apply": true, "destroy": true, "plan": true, "refresh": true, "init": true,
		"pulumi": true, "tofu": true, "terraform": true,
	}

	var walk func(c *cobra.Command)
	walk = func(c *cobra.Command) {
		if f := c.InheritedFlags().Lookup("module-dir"); f != nil || c.Flags().Lookup("module-dir") != nil {
			if f == nil {
				f = c.Flags().Lookup("module-dir")
			}
			if f.DefValue != "" {
				t.Errorf("%s: --module-dir default = %q, want empty (the resolver owns the default)", c.CommandPath(), f.DefValue)
			}
			if ann := f.Annotations[cobra.BashCompOneRequiredFlag]; len(ann) > 0 {
				t.Errorf("%s: --module-dir is marked required; the module runtime resolves an empty value", c.CommandPath())
			}
		}
		for _, sub := range c.Commands() {
			walk(sub)
		}
	}
	for _, c := range parent.Commands() {
		if runsAModule[c.Name()] {
			walk(c)
		}
	}
}

// Every command that deploys to Kubernetes reads --kube-context; each engine
// group must therefore register it, or the handler silently sees "" and the
// deploy lands on whatever context the kubeconfig currently selects.
func TestEngineGroups_RegisterKubeContext(t *testing.T) {
	parent := &cobra.Command{Use: "host"}
	RegisterCommands(parent, Options{})

	for _, group := range []string{"pulumi", "tofu", "terraform"} {
		var found *cobra.Command
		for _, c := range parent.Commands() {
			if c.Name() == group {
				found = c
			}
		}
		if found == nil {
			t.Fatalf("engine group %q not registered", group)
		}
		if found.PersistentFlags().Lookup("kube-context") == nil {
			t.Errorf("%s: --kube-context not registered on the group; its handlers read it", group)
		}
	}
	for _, lifecycle := range []string{"apply", "destroy", "plan", "refresh", "init"} {
		for _, c := range parent.Commands() {
			if c.Name() == lifecycle && c.PersistentFlags().Lookup("kube-context") == nil {
				t.Errorf("%s: --kube-context not registered", lifecycle)
			}
		}
	}
}
