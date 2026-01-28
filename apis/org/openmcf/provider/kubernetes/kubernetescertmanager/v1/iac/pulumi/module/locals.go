package module

import (
	"fmt"

	kubernetescertmanagerv1 "github.com/plantonhq/openmcf/apis/org/openmcf/provider/kubernetes/kubernetescertmanager/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals holds computed values used throughout the module
type Locals struct {
	KubernetesCertManager *kubernetescertmanagerv1.KubernetesCertManager
	Namespace             string
	// Computed resource names to avoid conflicts when multiple instances share a namespace
	// Format: {metadata.name}-{purpose}
	ServiceAccountName string
}

// CloudflareSecretName returns the computed secret name for a Cloudflare provider.
// Format: {metadata.name}-{provider-name}-credentials
func (l *Locals) CloudflareSecretName(providerName string) string {
	return fmt.Sprintf("%s-%s-credentials", l.KubernetesCertManager.Metadata.Name, providerName)
}

// initializeLocals creates and initializes the Locals struct
func initializeLocals(ctx *pulumi.Context, stackInput *kubernetescertmanagerv1.KubernetesCertManagerStackInput) *Locals {
	locals := &Locals{}

	locals.KubernetesCertManager = stackInput.Target

	target := stackInput.Target

	// Get namespace from spec (required field)
	locals.Namespace = target.Spec.Namespace.GetValue()

	// Computed resource names to avoid conflicts when multiple instances share a namespace
	// Users can prefix metadata.name with component type if needed (e.g., "cert-manager-prod")
	locals.ServiceAccountName = target.Metadata.Name

	return locals
}
