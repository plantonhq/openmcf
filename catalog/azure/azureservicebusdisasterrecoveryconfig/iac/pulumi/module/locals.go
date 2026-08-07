package module

import (
	azureservicebusdisasterrecoveryconfigv1alpha1 "github.com/plantonhq/planton/catalog/azure/azureservicebusdisasterrecoveryconfig/v1alpha1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AzureServiceBusDisasterRecoveryConfig *azureservicebusdisasterrecoveryconfigv1alpha1.AzureServiceBusDisasterRecoveryConfig
	PrimaryNamespaceId                    string
	PartnerNamespaceId                    string
	AliasAuthorizationRuleId              string
}

func initializeLocals(ctx *pulumi.Context, stackInput *azureservicebusdisasterrecoveryconfigv1alpha1.AzureServiceBusDisasterRecoveryConfigStackInput) *Locals {
	locals := &Locals{}

	locals.AzureServiceBusDisasterRecoveryConfig = stackInput.Target
	spec := stackInput.Target.Spec

	locals.PrimaryNamespaceId = spec.PrimaryNamespaceId.GetValue()
	locals.PartnerNamespaceId = spec.PartnerNamespaceId.GetValue()
	locals.AliasAuthorizationRuleId = spec.AliasAuthorizationRuleId.GetValue()

	// Geo-DR configs carry no Azure tags: ARM does not support tags on
	// disasterRecoveryConfigs, so the platform's identity tags live on the
	// paired namespaces.

	return locals
}
