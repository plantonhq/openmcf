package module

import (
	"strings"

	alicloudvpngatewayv1alpha1 "github.com/plantonhq/planton/apis/dev/planton/provider/alicloud/alicloudvpngateway/v1alpha1"
	"github.com/plantonhq/planton/apis/dev/planton/shared/cloudresourcekind"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AliCloudVpnGateway *alicloudvpngatewayv1alpha1.AliCloudVpnGateway
	Tags               map[string]string
}

func initializeLocals(ctx *pulumi.Context, stackInput *alicloudvpngatewayv1alpha1.AliCloudVpnGatewayStackInput) *Locals {
	locals := &Locals{}
	locals.AliCloudVpnGateway = stackInput.Target
	target := stackInput.Target

	locals.Tags = map[string]string{
		"resource":      "true",
		"resource_name": target.Metadata.Name,
		"resource_kind": strings.ToLower(cloudresourcekind.CloudResourceKind_AliCloudVpnGateway.String()),
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

func paymentType(spec *alicloudvpngatewayv1alpha1.AliCloudVpnGatewaySpec) string {
	if spec.PaymentType != nil {
		return *spec.PaymentType
	}
	return "PayAsYouGo"
}
