package module

import (
	"strings"

	alicloudramrolev1alpha1 "github.com/plantonhq/planton/catalog/alicloud/alicloudramrole/v1alpha1"
	"github.com/plantonhq/planton/shared/cloudresourcekind"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AliCloudRamRole *alicloudramrolev1alpha1.AliCloudRamRole
	Tags            map[string]string
}

func initializeLocals(ctx *pulumi.Context, stackInput *alicloudramrolev1alpha1.AliCloudRamRoleStackInput) *Locals {
	locals := &Locals{}
	locals.AliCloudRamRole = stackInput.Target
	target := stackInput.Target

	locals.Tags = map[string]string{
		"resource":      "true",
		"resource_name": target.Metadata.Name,
		"resource_kind": strings.ToLower(cloudresourcekind.CloudResourceKind_AliCloudRamRole.String()),
	}

	if target.Metadata.Id != "" {
		locals.Tags["resource_id"] = target.Metadata.Id
	}

	if target.Metadata.Org != "" {
		locals.Tags["organization"] = target.Metadata.Org
	}

	if target.Metadata.Env != "" {
		locals.Tags["environment"] = target.Metadata.Env
	}

	for k, v := range target.Spec.Tags {
		locals.Tags[k] = v
	}

	return locals
}

func maxSessionDuration(spec *alicloudramrolev1alpha1.AliCloudRamRoleSpec) int {
	if spec.MaxSessionDuration != nil {
		return int(*spec.MaxSessionDuration)
	}
	return 3600
}

func forceDelete(spec *alicloudramrolev1alpha1.AliCloudRamRoleSpec) bool {
	if spec.Force != nil {
		return *spec.Force
	}
	return false
}

func policyType(pa *alicloudramrolev1alpha1.AliCloudRamRolePolicyAttachment) string {
	if pa.PolicyType != nil {
		return *pa.PolicyType
	}
	return "System"
}
