package module

import (
	awsroute53dnsrecordv1 "github.com/plantonhq/project-planton/apis/org/project_planton/provider/aws/awsroute53dnsrecord/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AwsRoute53DnsRecord *awsroute53dnsrecordv1.AwsRoute53DnsRecord
}

func initializeLocals(ctx *pulumi.Context, stackInput *awsroute53dnsrecordv1.AwsRoute53DnsRecordStackInput) *Locals {
	locals := &Locals{}
	locals.AwsRoute53DnsRecord = stackInput.Target
	return locals
}
