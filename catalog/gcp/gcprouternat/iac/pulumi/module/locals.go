package module

import (
	gcpprovider "github.com/plantonhq/planton/catalog/gcp"
	gcprouternatv1alpha1 "github.com/plantonhq/planton/catalog/gcp/gcprouternat/v1alpha1"
)

// Locals collects frequently used input values (mirrors the Terraform
// "locals" pattern). Routers and NATs accept no labels in the GCP API, so
// there is no label set to derive for this kind.
type Locals struct {
	GcpProviderConfig *gcpprovider.GcpProviderConfig
	GcpRouterNat      *gcprouternatv1alpha1.GcpRouterNat
}

// initializeLocals converts the stack-input into a struct that is easy to
// reference across the module.
func initializeLocals(stackInput *gcprouternatv1alpha1.GcpRouterNatStackInput) *Locals {
	return &Locals{
		GcpProviderConfig: stackInput.ProviderConfig,
		GcpRouterNat:      stackInput.Target,
	}
}
