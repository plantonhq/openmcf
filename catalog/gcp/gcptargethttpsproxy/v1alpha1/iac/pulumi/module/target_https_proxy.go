package module

import (
	"strconv"

	"github.com/pkg/errors"
	foreignkeyv1 "github.com/plantonhq/planton/shared/foreignkey/v1"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/compute"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/projects"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// targetHttpsProxy provisions the global Compute Engine target HTTPS proxy —
// the TLS-termination node that binds a global forwarding rule (the VIP) to
// a URL map (the routing brain) and owns the client-facing TLS handshake:
// certificates, SSL policy, QUIC negotiation, and TLS 1.3 early data.
//
// Certificates attach through exactly one of three mechanisms (enforced
// pre-deploy by the spec's CEL): the classic ssl_certificates list, the
// cross-region-internal-ALB certificate_manager_certificates list, or an
// SNI-scale certificate_map. Traffic Director proxies skip certificates and
// drive TLS through server_tls_policy instead.
//
// url_map, the certificate wiring, ssl_policy, server_tls_policy, and
// quic_override all update in place via dedicated API calls — certificate
// rotation is attach-new-then-detach-old with zero VIP churn. name,
// description, keep-alive, tls_early_data, and proxy_bind are immutable.
func targetHttpsProxy(ctx *pulumi.Context, locals *Locals, gcpProvider *gcp.Provider) error {
	spec := locals.GcpTargetHttpsProxy.Spec

	// Enable the Compute Engine API first so a fresh project works on the
	// first deploy. disable_on_destroy stays false: tearing down one proxy
	// must never disable the API for everything else in the project.
	serviceArgs := &projects.ServiceArgs{
		Service:                  pulumi.String("compute.googleapis.com"),
		DisableDependentServices: pulumi.BoolPtr(true),
	}
	if spec.ProjectId.GetValue() != "" {
		serviceArgs.Project = pulumi.String(spec.ProjectId.GetValue())
	}
	createdProjectService, err := projects.NewService(ctx,
		"targethttpsproxy-compute.googleapis.com", serviceArgs, pulumi.Provider(gcpProvider))
	if err != nil {
		return errors.Wrap(err, "failed to enable compute.googleapis.com api")
	}

	args := &compute.TargetHttpsProxyArgs{
		Name: pulumi.String(locals.ProxyName),
		// The URL map ref arrives resolved to a literal self-link (or a
		// plain name, which the provider expands against the project).
		UrlMap: pulumi.String(spec.UrlMap.GetValue()),
	}

	// An empty project falls back to the provider's default project — the
	// ambient-project contract every GCP kind honors.
	if spec.ProjectId.GetValue() != "" {
		args.Project = pulumi.String(spec.ProjectId.GetValue())
	}
	if spec.Description != "" {
		args.Description = pulumi.String(spec.Description)
	}
	if certs := refListValues(spec.SslCertificates); len(certs) > 0 {
		args.SslCertificates = certs
	}
	if certs := refListValues(spec.CertificateManagerCertificates); len(certs) > 0 {
		args.CertificateManagerCertificates = certs
	}
	if spec.CertificateMap != "" {
		args.CertificateMap = pulumi.String(spec.CertificateMap)
	}
	if spec.SslPolicy.GetValue() != "" {
		args.SslPolicy = pulumi.String(spec.SslPolicy.GetValue())
	}
	if spec.ServerTlsPolicy.GetValue() != "" {
		args.ServerTlsPolicy = pulumi.String(spec.ServerTlsPolicy.GetValue())
	}
	// The middleware default (NONE) matches GCP's own default, so an unset
	// value can simply be omitted — the API computes NONE either way.
	if spec.QuicOverride != nil && spec.GetQuicOverride() != "" {
		args.QuicOverride = pulumi.String(spec.GetQuicOverride())
	}
	// Empty lets GCP apply its default (DISABLED); the field is immutable,
	// so an explicit value is only sent when the user chose a mode.
	if spec.TlsEarlyData != "" {
		args.TlsEarlyData = pulumi.String(spec.TlsEarlyData)
	}
	// 0 means "let GCP apply its default" (610s on EXTERNAL_MANAGED); the
	// field is only honored by the envoy-based external ALB.
	if spec.HttpKeepAliveTimeoutSec != 0 {
		args.HttpKeepAliveTimeoutSec = pulumi.Int(int(spec.HttpKeepAliveTimeoutSec))
	}
	// proxy_bind is a Traffic Director lever; the API default is false, so
	// only an explicit true is worth sending.
	if spec.ProxyBind {
		args.ProxyBind = pulumi.Bool(true)
	}

	createdProxy, err := compute.NewTargetHttpsProxy(ctx, "target-https-proxy", args,
		pulumi.Provider(gcpProvider), pulumi.DependsOn([]pulumi.Resource{createdProjectService}))
	if err != nil {
		return errors.Wrap(err, "failed to create target https proxy")
	}

	ctx.Export(OpSelfLink, createdProxy.SelfLink)
	ctx.Export(OpProxyName, createdProxy.Name)
	ctx.Export(OpProxyId, createdProxy.ProxyId.ApplyT(func(id int) string {
		return strconv.Itoa(id)
	}).(pulumi.StringOutput))
	ctx.Export(OpFingerprint, createdProxy.Fingerprint)

	return nil
}

// refListValues flattens a repeated StringValueOrRef into the literal values
// the refs resolved to, skipping empties so an unresolved placeholder never
// reaches the API.
func refListValues(refs []*foreignkeyv1.StringValueOrRef) pulumi.StringArray {
	values := pulumi.StringArray{}
	for _, ref := range refs {
		if ref.GetValue() != "" {
			values = append(values, pulumi.String(ref.GetValue()))
		}
	}
	return values
}
