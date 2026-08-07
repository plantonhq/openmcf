package module

import (
	civoprovider "github.com/plantonhq/planton/catalog/civo"
	civocomputeinstancev1alpha1 "github.com/plantonhq/planton/catalog/civo/civocomputeinstance/v1alpha1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals bundles pointers frequently reused across the module.
type Locals struct {
	CivoProviderConfig  *civoprovider.CivoProviderConfig
	CivoComputeInstance *civocomputeinstancev1alpha1.CivoComputeInstance
}

// initializeLocals mirrors the simple pattern used in other Planton modules.
func initializeLocals(
	_ *pulumi.Context,
	stackInput *civocomputeinstancev1alpha1.CivoComputeInstanceStackInput,
) *Locals {
	return &Locals{
		CivoComputeInstance: stackInput.Target,
		CivoProviderConfig:  stackInput.ProviderConfig,
	}
}
