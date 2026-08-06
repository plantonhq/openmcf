package module

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

const (
	OpNamespace         = "namespace"
	OpCollectorName     = "collector_name"
	OpService           = "service"
	OpOtlpGrpcEndpoint  = "otlp_grpc_endpoint"
	OpOtlpHttpEndpoint  = "otlp_http_endpoint"
	OpHeadlessService   = "headless_service"
	OpMonitoringService = "monitoring_service"
)

// exportOutputs publishes the collector's composition handles. Sidecar
// mode creates NO standalone workload — every Service-derived handle is
// empty there (the operator injects the collector into annotated pods
// instead). Twin: the Terraform module's outputs.tf.
func exportOutputs(ctx *pulumi.Context, locals *Locals) error {
	service := ""
	otlpGrpc := ""
	otlpHttp := ""
	headless := ""
	monitoring := ""
	if locals.EffectiveMode != "sidecar" {
		service = locals.ResourceName + "-collector"
		otlpGrpc = service + "." + locals.Namespace + ".svc.cluster.local:4317"
		otlpHttp = "http://" + service + "." + locals.Namespace + ".svc.cluster.local:4318"
		headless = service + "-headless"
		monitoring = service + "-monitoring"
	}

	ctx.Export(OpNamespace, pulumi.String(locals.Namespace))
	ctx.Export(OpCollectorName, pulumi.String(locals.ResourceName))
	ctx.Export(OpService, pulumi.String(service))
	ctx.Export(OpOtlpGrpcEndpoint, pulumi.String(otlpGrpc))
	ctx.Export(OpOtlpHttpEndpoint, pulumi.String(otlpHttp))
	ctx.Export(OpHeadlessService, pulumi.String(headless))
	ctx.Export(OpMonitoringService, pulumi.String(monitoring))
	return nil
}
