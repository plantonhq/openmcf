package module

import (
	azurecosmosdbsqlroleassignmentv1alpha1 "github.com/plantonhq/planton/apis/dev/planton/provider/azure/azurecosmosdbsqlroleassignment/v1alpha1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AzureCosmosdbSqlRoleAssignment *azurecosmosdbsqlroleassignmentv1alpha1.AzureCosmosdbSqlRoleAssignment
	CosmosdbAccountId              string
	RoleDefinitionId               string
	PrincipalId                    string
	Scope                          string
}

func initializeLocals(ctx *pulumi.Context, stackInput *azurecosmosdbsqlroleassignmentv1alpha1.AzureCosmosdbSqlRoleAssignmentStackInput) *Locals {
	locals := &Locals{}

	locals.AzureCosmosdbSqlRoleAssignment = stackInput.Target
	locals.CosmosdbAccountId = stackInput.Target.Spec.CosmosdbAccountId.GetValue()
	locals.RoleDefinitionId = stackInput.Target.Spec.RoleDefinitionId.GetValue()
	locals.PrincipalId = stackInput.Target.Spec.PrincipalId.GetValue()
	locals.Scope = stackInput.Target.Spec.Scope.GetValue()

	// No Azure tags: ARM does not support tags on Cosmos child
	// resources, so the platform's identity tags live on the account.

	return locals
}
