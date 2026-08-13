package module

import (
	"strconv"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/compute"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/projects"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// targetHttpProxy provisions the global Compute Engine target HTTP proxy —
// the plaintext-HTTP frontend adapter that binds a global forwarding rule
// (the VIP) to a URL map (the routing brain). The proxy is deliberately
// thin: TLS lives on the target HTTPS proxy sibling, routing on the URL
// map, traffic policy on the backend service.
//
// url_map is the only mutable field — GCP repoints it in place via a
// dedicated setUrlMap call, so a live frontend can move to a new routing
// table with no downtime. Everything else (name, description, keep-alive,
// proxy_bind, project) is immutable and forces destroy-and-recreate,
// briefly breaking any forwarding rule that references the old self_link.
func targetHttpProxy(ctx *pulumi.Context, locals *Locals, gcpProvider *gcp.Provider) error {
	spec := locals.GcpTargetHttpProxy.Spec

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
		"targethttpproxy-compute.googleapis.com", serviceArgs, pulumi.Provider(gcpProvider))
	if err != nil {
		return errors.Wrap(err, "failed to enable compute.googleapis.com api")
	}

	args := &compute.TargetHttpProxyArgs{
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
	// Empty defers to the provider default (DELETE).
	if spec.DeletionPolicy != "" {
		args.DeletionPolicy = pulumi.StringPtr(spec.DeletionPolicy)
	}

	createdProxy, err := compute.NewTargetHttpProxy(ctx, "target-http-proxy", args,
		pulumi.Provider(gcpProvider), pulumi.DependsOn([]pulumi.Resource{createdProjectService}))
	if err != nil {
		return errors.Wrap(err, "failed to create target http proxy")
	}

	ctx.Export(OpSelfLink, createdProxy.SelfLink)
	ctx.Export(OpProxyName, createdProxy.Name)
	ctx.Export(OpProxyId, createdProxy.ProxyId.ApplyT(func(id int) string {
		return strconv.Itoa(id)
	}).(pulumi.StringOutput))
	ctx.Export(OpFingerprint, createdProxy.Fingerprint)

	return nil
}
