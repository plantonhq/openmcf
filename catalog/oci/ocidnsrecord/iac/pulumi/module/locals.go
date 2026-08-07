package module

import (
	ocidnsrecordv1alpha1 "github.com/plantonhq/planton/catalog/oci/ocidnsrecord/v1alpha1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	OciDnsRecord *ocidnsrecordv1alpha1.OciDnsRecord
	ResourceName string
}

func initializeLocals(_ *pulumi.Context, stackInput *ocidnsrecordv1alpha1.OciDnsRecordStackInput) *Locals {
	return &Locals{
		OciDnsRecord: stackInput.Target,
		ResourceName: stackInput.Target.Metadata.Name,
	}
}
