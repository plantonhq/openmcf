package module

import (
	alicloudcontainerregistryv1alpha1 "github.com/plantonhq/planton/apis/dev/planton/provider/alicloud/alicloudcontainerregistry/v1alpha1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AliCloudContainerRegistry *alicloudcontainerregistryv1alpha1.AliCloudContainerRegistry
}

func initializeLocals(_ *pulumi.Context, stackInput *alicloudcontainerregistryv1alpha1.AliCloudContainerRegistryStackInput) *Locals {
	return &Locals{
		AliCloudContainerRegistry: stackInput.Target,
	}
}

func paymentType(spec *alicloudcontainerregistryv1alpha1.AliCloudContainerRegistrySpec) string {
	if spec.PaymentType != nil {
		return *spec.PaymentType
	}
	return "Subscription"
}

func namespaceAutoCreate(ns *alicloudcontainerregistryv1alpha1.AliCloudContainerRegistryNamespace) bool {
	if ns.AutoCreate != nil {
		return *ns.AutoCreate
	}
	return false
}

func namespaceDefaultVisibility(ns *alicloudcontainerregistryv1alpha1.AliCloudContainerRegistryNamespace) string {
	if ns.DefaultVisibility != nil {
		return *ns.DefaultVisibility
	}
	return "PRIVATE"
}

func optionalString(s string) pulumi.StringPtrInput {
	if s == "" {
		return nil
	}
	return pulumi.String(s)
}
