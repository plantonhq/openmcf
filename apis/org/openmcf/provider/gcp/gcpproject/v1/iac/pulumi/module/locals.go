package module

import (
	"strconv"
	"strings"

	gcpprojectv1 "github.com/plantonhq/openmcf/apis/org/openmcf/provider/gcp/gcpproject/v1"
	"github.com/plantonhq/openmcf/apis/org/openmcf/shared/cloudresourcekind"
	"github.com/plantonhq/openmcf/pkg/iac/pulumi/pulumimodule/provider/gcp/gcplabelkeys"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	GcpProject *gcpprojectv1.GcpProject
	GcpLabels  map[string]string
}

func initializeLocals(ctx *pulumi.Context, stackInput *gcpprojectv1.GcpProjectStackInput) *Locals {
	locals := &Locals{}
	locals.GcpProject = stackInput.Target

	target := stackInput.Target

	// Start with user-provided labels from spec
	locals.GcpLabels = make(map[string]string)
	for k, v := range target.Spec.Labels {
		locals.GcpLabels[k] = v
	}

	// Add OpenMCF internal labels (these override user labels if there's a conflict)
	locals.GcpLabels[gcplabelkeys.Resource] = strconv.FormatBool(true)
	locals.GcpLabels[gcplabelkeys.ResourceName] = target.Metadata.Name
	locals.GcpLabels[gcplabelkeys.ResourceKind] = strings.ToLower(cloudresourcekind.CloudResourceKind_GcpProject.String())

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
