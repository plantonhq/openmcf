package root

import (
	"github.com/plantonhq/planton/cmd/planton/root/tofu"
	"github.com/plantonhq/planton/internal/cli/flag"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
)

var Tofu = &cobra.Command{
	Use:   "tofu",
	Short: "run open-tofu commands",
}

func init() {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal("failed to get current working directory")
	}

	Tofu.PersistentFlags().String(string(flag.Manifest), "", "path of the deployment-component manifest file")

	Tofu.PersistentFlags().String(string(flag.InputDir), "", "directory containing target.yaml and credential yaml files")
	Tofu.PersistentFlags().String(string(flag.KustomizeDir), "", "directory containing kustomize configuration")
	Tofu.PersistentFlags().String(string(flag.Overlay), "", "kustomize overlay to use (e.g., prod, dev, staging)")
	Tofu.PersistentFlags().String(string(flag.ModuleDir), pwd, "directory containing the terraform module")
	Tofu.PersistentFlags().StringToString(string(flag.Set), map[string]string{}, "override resource manifest values using key=value pairs")

	// Stack input file flag: the shared manifest resolver reads this on every
	// command, so it must be registered here (as on the pulumi command group)
	// or resolution fails before --manifest is even considered.
	Tofu.PersistentFlags().StringP(string(flag.StackInput), "i", "", "path to a YAML file containing the stack input (bypasses building stack input from manifest)")

	// Provider config flag (unified)
	Tofu.PersistentFlags().StringP(string(flag.ProviderConfig), "p", "", "path to provider credentials file")

	Tofu.AddCommand(
		tofu.Apply,
		tofu.Destroy,
		tofu.GenerateModule,
		tofu.GenerateVariables,
		tofu.Init,
		tofu.LoadTfVars,
		tofu.Plan,
		tofu.Refresh,
	)
}
