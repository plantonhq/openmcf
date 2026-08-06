package module

import (
	"github.com/pkg/errors"
	gcphealthcheckv1alpha1 "github.com/plantonhq/planton/apis/dev/planton/provider/gcp/gcphealthcheck/v1alpha1"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/compute"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/projects"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// healthCheck provisions the Compute Engine health check — the probe backend
// services consult before sending traffic to a backend, and managed instance
// groups consult before auto-healing an instance.
//
// One kind, two provider resources: GCP models global and regional health
// checks as separate API collections with an otherwise identical surface, so
// an empty spec.region creates compute.HealthCheck and a set region creates
// compute.RegionHealthCheck — mirroring the Terraform module's count guards.
//
// name and project are immutable (ForceNew in the provider): changing either
// destroys and recreates the check, briefly breaking every backend service
// referencing the old self_link. All probe knobs update in place.
//
// Ports deliberately fall through to the API defaults (http/tcp 80,
// https/http2/ssl 443) when unset — hardcoding them here would silently pin
// behavior the provider may evolve.
func healthCheck(ctx *pulumi.Context, locals *Locals, gcpProvider *gcp.Provider) error {
	spec := locals.GcpHealthCheck.Spec

	// Enable the Compute Engine API so a fresh project can host the health
	// check. disable_on_destroy stays false (the provider default): tearing
	// down one health check must never disable the API for everything else in
	// the project. Matches the Terraform module.
	serviceArgs := &projects.ServiceArgs{
		Service:                  pulumi.String("compute.googleapis.com"),
		DisableDependentServices: pulumi.BoolPtr(true),
	}
	if spec.ProjectId.GetValue() != "" {
		serviceArgs.Project = pulumi.String(spec.ProjectId.GetValue())
	}
	createdProjectService, err := projects.NewService(ctx,
		"healthcheck-compute.googleapis.com", serviceArgs, pulumi.Provider(gcpProvider))
	if err != nil {
		return errors.Wrap(err, "failed to enable compute.googleapis.com api")
	}

	if spec.Region == "" {
		return globalHealthCheck(ctx, locals, gcpProvider, createdProjectService)
	}
	return regionalHealthCheck(ctx, locals, gcpProvider, createdProjectService)
}

func globalHealthCheck(ctx *pulumi.Context, locals *Locals, gcpProvider *gcp.Provider, computeApiService pulumi.Resource) error {
	spec := locals.GcpHealthCheck.Spec

	args := &compute.HealthCheckArgs{
		Name:               pulumi.String(locals.HealthCheckName),
		CheckIntervalSec:   pulumi.Int(spec.GetCheckIntervalSec()),
		TimeoutSec:         pulumi.Int(spec.GetTimeoutSec()),
		HealthyThreshold:   pulumi.Int(spec.GetHealthyThreshold()),
		UnhealthyThreshold: pulumi.Int(spec.GetUnhealthyThreshold()),
		// The block is always emitted so disabling logging is an explicit
		// false, not an absent block the API back-fills (which would show as
		// drift). Matches the Terraform module.
		LogConfig: &compute.HealthCheckLogConfigArgs{
			Enable: pulumi.Bool(spec.EnableLogging),
		},
	}

	// Omitted description stays unset (matching the Terraform module's null)
	// rather than being sent as an empty string.
	if spec.Description != "" {
		args.Description = pulumi.String(spec.Description)
	}

	// Honor the spec contract: an empty project_id falls back to the provider's
	// default project. Leaving Project unset lets the gcp provider resolve its
	// own project (configuration or the GOOGLE_PROJECT / GOOGLE_CLOUD_PROJECT
	// environment chain); an empty string would be sent verbatim and rejected.
	if spec.ProjectId.GetValue() != "" {
		args.Project = pulumi.String(spec.ProjectId.GetValue())
	}

	// Global-only: pin probing to exactly 3 named regions so a regional outage
	// cannot flip global health verdicts.
	if len(spec.SourceRegions) > 0 {
		sourceRegions := pulumi.StringArray{}
		for _, sourceRegion := range spec.SourceRegions {
			sourceRegions = append(sourceRegions, pulumi.String(sourceRegion))
		}
		args.SourceRegions = sourceRegions
	}

	// Exactly one protocol arm is set (enforced by the proto oneof).
	switch protocol := spec.Protocol.(type) {
	case *gcphealthcheckv1alpha1.GcpHealthCheckSpec_Http:
		args.HttpHealthCheck = &compute.HealthCheckHttpHealthCheckArgs{
			Host:              emptyAsNilString(protocol.Http.Host),
			Port:              zeroAsNilInt(protocol.Http.Port),
			PortName:          emptyAsNilString(protocol.Http.PortName),
			PortSpecification: emptyAsNilString(protocol.Http.PortSpecification),
			ProxyHeader:       emptyAsNilString(protocol.Http.ProxyHeader),
			RequestPath:       emptyAsNilString(protocol.Http.RequestPath),
			Response:          emptyAsNilString(protocol.Http.Response),
		}
	case *gcphealthcheckv1alpha1.GcpHealthCheckSpec_Https:
		args.HttpsHealthCheck = &compute.HealthCheckHttpsHealthCheckArgs{
			Host:              emptyAsNilString(protocol.Https.Host),
			Port:              zeroAsNilInt(protocol.Https.Port),
			PortName:          emptyAsNilString(protocol.Https.PortName),
			PortSpecification: emptyAsNilString(protocol.Https.PortSpecification),
			ProxyHeader:       emptyAsNilString(protocol.Https.ProxyHeader),
			RequestPath:       emptyAsNilString(protocol.Https.RequestPath),
			Response:          emptyAsNilString(protocol.Https.Response),
		}
	case *gcphealthcheckv1alpha1.GcpHealthCheckSpec_Http2:
		args.Http2HealthCheck = &compute.HealthCheckHttp2HealthCheckArgs{
			Host:              emptyAsNilString(protocol.Http2.Host),
			Port:              zeroAsNilInt(protocol.Http2.Port),
			PortName:          emptyAsNilString(protocol.Http2.PortName),
			PortSpecification: emptyAsNilString(protocol.Http2.PortSpecification),
			ProxyHeader:       emptyAsNilString(protocol.Http2.ProxyHeader),
			RequestPath:       emptyAsNilString(protocol.Http2.RequestPath),
			Response:          emptyAsNilString(protocol.Http2.Response),
		}
	case *gcphealthcheckv1alpha1.GcpHealthCheckSpec_Tcp:
		args.TcpHealthCheck = &compute.HealthCheckTcpHealthCheckArgs{
			Port:              zeroAsNilInt(protocol.Tcp.Port),
			PortName:          emptyAsNilString(protocol.Tcp.PortName),
			PortSpecification: emptyAsNilString(protocol.Tcp.PortSpecification),
			ProxyHeader:       emptyAsNilString(protocol.Tcp.ProxyHeader),
			Request:           emptyAsNilString(protocol.Tcp.Request),
			Response:          emptyAsNilString(protocol.Tcp.Response),
		}
	case *gcphealthcheckv1alpha1.GcpHealthCheckSpec_Ssl:
		args.SslHealthCheck = &compute.HealthCheckSslHealthCheckArgs{
			Port:              zeroAsNilInt(protocol.Ssl.Port),
			PortName:          emptyAsNilString(protocol.Ssl.PortName),
			PortSpecification: emptyAsNilString(protocol.Ssl.PortSpecification),
			ProxyHeader:       emptyAsNilString(protocol.Ssl.ProxyHeader),
			Request:           emptyAsNilString(protocol.Ssl.Request),
			Response:          emptyAsNilString(protocol.Ssl.Response),
		}
	case *gcphealthcheckv1alpha1.GcpHealthCheckSpec_Grpc:
		args.GrpcHealthCheck = &compute.HealthCheckGrpcHealthCheckArgs{
			GrpcServiceName:   emptyAsNilString(protocol.Grpc.GrpcServiceName),
			Port:              zeroAsNilInt(protocol.Grpc.Port),
			PortName:          emptyAsNilString(protocol.Grpc.PortName),
			PortSpecification: emptyAsNilString(protocol.Grpc.PortSpecification),
		}
	case *gcphealthcheckv1alpha1.GcpHealthCheckSpec_GrpcTls:
		args.GrpcTlsHealthCheck = &compute.HealthCheckGrpcTlsHealthCheckArgs{
			GrpcServiceName:   emptyAsNilString(protocol.GrpcTls.GrpcServiceName),
			Port:              zeroAsNilInt(protocol.GrpcTls.Port),
			PortSpecification: emptyAsNilString(protocol.GrpcTls.PortSpecification),
		}
	}

	createdHealthCheck, err := compute.NewHealthCheck(ctx, "health-check", args,
		pulumi.Provider(gcpProvider), pulumi.DependsOn([]pulumi.Resource{computeApiService}))
	if err != nil {
		return errors.Wrap(err, "failed to create global health check")
	}

	ctx.Export(OpSelfLink, createdHealthCheck.SelfLink)
	ctx.Export(OpHealthCheckName, createdHealthCheck.Name)
	ctx.Export(OpType, createdHealthCheck.Type)
	// Empty region marks the global scope for downstream composition checks.
	ctx.Export(OpRegion, pulumi.String(""))

	return nil
}

func regionalHealthCheck(ctx *pulumi.Context, locals *Locals, gcpProvider *gcp.Provider, computeApiService pulumi.Resource) error {
	spec := locals.GcpHealthCheck.Spec

	args := &compute.RegionHealthCheckArgs{
		Name:               pulumi.String(locals.HealthCheckName),
		Region:             pulumi.String(spec.Region),
		CheckIntervalSec:   pulumi.Int(spec.GetCheckIntervalSec()),
		TimeoutSec:         pulumi.Int(spec.GetTimeoutSec()),
		HealthyThreshold:   pulumi.Int(spec.GetHealthyThreshold()),
		UnhealthyThreshold: pulumi.Int(spec.GetUnhealthyThreshold()),
		LogConfig: &compute.RegionHealthCheckLogConfigArgs{
			Enable: pulumi.Bool(spec.EnableLogging),
		},
	}

	if spec.Description != "" {
		args.Description = pulumi.String(spec.Description)
	}

	if spec.ProjectId.GetValue() != "" {
		args.Project = pulumi.String(spec.ProjectId.GetValue())
	}

	switch protocol := spec.Protocol.(type) {
	case *gcphealthcheckv1alpha1.GcpHealthCheckSpec_Http:
		args.HttpHealthCheck = &compute.RegionHealthCheckHttpHealthCheckArgs{
			Host:              emptyAsNilString(protocol.Http.Host),
			Port:              zeroAsNilInt(protocol.Http.Port),
			PortName:          emptyAsNilString(protocol.Http.PortName),
			PortSpecification: emptyAsNilString(protocol.Http.PortSpecification),
			ProxyHeader:       emptyAsNilString(protocol.Http.ProxyHeader),
			RequestPath:       emptyAsNilString(protocol.Http.RequestPath),
			Response:          emptyAsNilString(protocol.Http.Response),
		}
	case *gcphealthcheckv1alpha1.GcpHealthCheckSpec_Https:
		args.HttpsHealthCheck = &compute.RegionHealthCheckHttpsHealthCheckArgs{
			Host:              emptyAsNilString(protocol.Https.Host),
			Port:              zeroAsNilInt(protocol.Https.Port),
			PortName:          emptyAsNilString(protocol.Https.PortName),
			PortSpecification: emptyAsNilString(protocol.Https.PortSpecification),
			ProxyHeader:       emptyAsNilString(protocol.Https.ProxyHeader),
			RequestPath:       emptyAsNilString(protocol.Https.RequestPath),
			Response:          emptyAsNilString(protocol.Https.Response),
		}
	case *gcphealthcheckv1alpha1.GcpHealthCheckSpec_Http2:
		args.Http2HealthCheck = &compute.RegionHealthCheckHttp2HealthCheckArgs{
			Host:              emptyAsNilString(protocol.Http2.Host),
			Port:              zeroAsNilInt(protocol.Http2.Port),
			PortName:          emptyAsNilString(protocol.Http2.PortName),
			PortSpecification: emptyAsNilString(protocol.Http2.PortSpecification),
			ProxyHeader:       emptyAsNilString(protocol.Http2.ProxyHeader),
			RequestPath:       emptyAsNilString(protocol.Http2.RequestPath),
			Response:          emptyAsNilString(protocol.Http2.Response),
		}
	case *gcphealthcheckv1alpha1.GcpHealthCheckSpec_Tcp:
		args.TcpHealthCheck = &compute.RegionHealthCheckTcpHealthCheckArgs{
			Port:              zeroAsNilInt(protocol.Tcp.Port),
			PortName:          emptyAsNilString(protocol.Tcp.PortName),
			PortSpecification: emptyAsNilString(protocol.Tcp.PortSpecification),
			ProxyHeader:       emptyAsNilString(protocol.Tcp.ProxyHeader),
			Request:           emptyAsNilString(protocol.Tcp.Request),
			Response:          emptyAsNilString(protocol.Tcp.Response),
		}
	case *gcphealthcheckv1alpha1.GcpHealthCheckSpec_Ssl:
		args.SslHealthCheck = &compute.RegionHealthCheckSslHealthCheckArgs{
			Port:              zeroAsNilInt(protocol.Ssl.Port),
			PortName:          emptyAsNilString(protocol.Ssl.PortName),
			PortSpecification: emptyAsNilString(protocol.Ssl.PortSpecification),
			ProxyHeader:       emptyAsNilString(protocol.Ssl.ProxyHeader),
			Request:           emptyAsNilString(protocol.Ssl.Request),
			Response:          emptyAsNilString(protocol.Ssl.Response),
		}
	case *gcphealthcheckv1alpha1.GcpHealthCheckSpec_Grpc:
		args.GrpcHealthCheck = &compute.RegionHealthCheckGrpcHealthCheckArgs{
			GrpcServiceName:   emptyAsNilString(protocol.Grpc.GrpcServiceName),
			Port:              zeroAsNilInt(protocol.Grpc.Port),
			PortName:          emptyAsNilString(protocol.Grpc.PortName),
			PortSpecification: emptyAsNilString(protocol.Grpc.PortSpecification),
		}
	case *gcphealthcheckv1alpha1.GcpHealthCheckSpec_GrpcTls:
		args.GrpcTlsHealthCheck = &compute.RegionHealthCheckGrpcTlsHealthCheckArgs{
			GrpcServiceName:   emptyAsNilString(protocol.GrpcTls.GrpcServiceName),
			Port:              zeroAsNilInt(protocol.GrpcTls.Port),
			PortSpecification: emptyAsNilString(protocol.GrpcTls.PortSpecification),
		}
	}

	createdHealthCheck, err := compute.NewRegionHealthCheck(ctx, "health-check", args,
		pulumi.Provider(gcpProvider), pulumi.DependsOn([]pulumi.Resource{computeApiService}))
	if err != nil {
		return errors.Wrap(err, "failed to create regional health check")
	}

	ctx.Export(OpSelfLink, createdHealthCheck.SelfLink)
	ctx.Export(OpHealthCheckName, createdHealthCheck.Name)
	ctx.Export(OpType, createdHealthCheck.Type)
	ctx.Export(OpRegion, createdHealthCheck.Region)

	return nil
}

// emptyAsNilString leaves unset proto strings out of the API payload —
// matching the Terraform module's "" -> null normalization so both engines
// send identical requests.
func emptyAsNilString(value string) pulumi.StringPtrInput {
	if value == "" {
		return nil
	}
	return pulumi.String(value)
}

// zeroAsNilInt leaves unset proto ints out of the API payload so GCP applies
// its own protocol-default ports — matching the Terraform module.
func zeroAsNilInt(value int32) pulumi.IntPtrInput {
	if value == 0 {
		return nil
	}
	return pulumi.Int(int(value))
}
