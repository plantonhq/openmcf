package module

import (
	aliclouddnsrecordv1alpha1 "github.com/plantonhq/planton/catalog/alicloud/aliclouddnsrecord/v1alpha1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AliCloudDnsRecord *aliclouddnsrecordv1alpha1.AliCloudDnsRecord
}

func initializeLocals(ctx *pulumi.Context, stackInput *aliclouddnsrecordv1alpha1.AliCloudDnsRecordStackInput) *Locals {
	return &Locals{
		AliCloudDnsRecord: stackInput.Target,
	}
}
