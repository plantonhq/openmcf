package module

import (
	"strconv"
	"strings"

	gcpcertmanagercertv1 "github.com/plantonhq/openmcf/apis/org/openmcf/provider/gcp/gcpcertmanagercert/v1"
	"github.com/plantonhq/openmcf/apis/org/openmcf/shared/cloudresourcekind"
	"github.com/plantonhq/openmcf/pkg/iac/pulumi/pulumimodule/provider/gcp/gcplabelkeys"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	GcpCertManagerCert *gcpcertmanagercertv1.GcpCertManagerCert
	GcpLabels          map[string]string
}

func initializeLocals(ctx *pulumi.Context, stackInput *gcpcertmanagercertv1.GcpCertManagerCertStackInput) *Locals {
	locals := &Locals{}

	locals.GcpCertManagerCert = stackInput.Target

	target := stackInput.Target

	locals.GcpLabels = map[string]string{
		gcplabelkeys.Resource:     strconv.FormatBool(true),
		gcplabelkeys.ResourceName: target.Metadata.Name,
		gcplabelkeys.ResourceKind: strings.ToLower(cloudresourcekind.CloudResourceKind_GcpCertManagerCert.String()),
	}

	if target.Metadata.Id != "" {
		locals.GcpLabels[gcplabelkeys.ResourceId] = target.Metadata.Id
	}

	if target.Metadata.Org != "" {
		locals.GcpLabels[gcplabelkeys.Organization] = target.Metadata.Org
	}

	if target.Metadata.Env != "" {
		locals.GcpLabels[gcplabelkeys.Environment] = target.Metadata.Env
	}

	return locals
}
