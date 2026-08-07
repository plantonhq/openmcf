package module

import (
	gcpprovider "github.com/plantonhq/planton/catalog/gcp"
	gcpcloudtasksqueuev1alpha1 "github.com/plantonhq/planton/catalog/gcp/gcpcloudtasksqueue/v1alpha1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	GcpProviderConfig  *gcpprovider.GcpProviderConfig
	GcpCloudTasksQueue *gcpcloudtasksqueuev1alpha1.GcpCloudTasksQueue
}

func initializeLocals(_ *pulumi.Context, stackInput *gcpcloudtasksqueuev1alpha1.GcpCloudTasksQueueStackInput) *Locals {
	locals := &Locals{}
	locals.GcpCloudTasksQueue = stackInput.Target
	locals.GcpProviderConfig = stackInput.ProviderConfig
	// Note: Cloud Tasks queues do NOT support GCP labels.
	// No label computation needed (unlike most GCP components).
	return locals
}
