package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/compute"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/projects"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// backendBucket provisions the Compute Engine backend bucket — the node that
// serves a Cloud Storage bucket's objects through an external HTTP(S) load
// balancer, optionally cached at Google's edge by Cloud CDN — plus its
// signed-URL keys.
//
// name and project are immutable (ForceNew in the provider): changing either
// destroys and recreates the backend bucket, briefly breaking every URL map
// referencing the old self_link. bucket_name is deliberately mutable —
// origin swaps (blue/green static releases) are in-place updates.
//
// CDN is a policy ON this resource, not a separate GCP object: enable_cdn
// turns edge caching on and cdn_policy tunes it. TTL/cache-mode coherence
// (e.g. USE_ORIGIN_HEADERS forbids explicit TTLs) is enforced by the spec
// before deploy, so no TTL-stripping logic lives here.
func backendBucket(ctx *pulumi.Context, locals *Locals, gcpProvider *gcp.Provider) error {
	spec := locals.GcpBackendBucket.Spec

	// Enable the Compute Engine API so a fresh project can host the backend
	// bucket. disable_on_destroy stays false (the provider default): tearing
	// down one backend bucket must never disable the API for everything else
	// in the project. Matches the Terraform module.
	serviceArgs := &projects.ServiceArgs{
		Service:                  pulumi.String("compute.googleapis.com"),
		DisableDependentServices: pulumi.BoolPtr(true),
	}
	if spec.ProjectId.GetValue() != "" {
		serviceArgs.Project = pulumi.String(spec.ProjectId.GetValue())
	}
	createdProjectService, err := projects.NewService(ctx,
		"backendbucket-compute.googleapis.com", serviceArgs, pulumi.Provider(gcpProvider))
	if err != nil {
		return errors.Wrap(err, "failed to enable compute.googleapis.com api")
	}

	args := &compute.BackendBucketArgs{
		Name:       pulumi.String(locals.BackendBucketName),
		BucketName: pulumi.String(spec.BucketName.GetValue()),
		EnableCdn:  pulumi.Bool(spec.EnableCdn),
	}

	// Omitted optionals stay unset (matching the Terraform module's null)
	// rather than being sent as empty strings the API would reject or
	// misread.
	if spec.Description != "" {
		args.Description = pulumi.String(spec.Description)
	}
	if spec.CompressionMode != "" {
		args.CompressionMode = pulumi.String(spec.CompressionMode)
	}
	if spec.LoadBalancingScheme != "" {
		args.LoadBalancingScheme = pulumi.String(spec.LoadBalancingScheme)
	}
	// Attaches by reference — the Cloud Armor EDGE policy itself is its own
	// composable node (GcpCloudArmorPolicy), never embedded here.
	if spec.EdgeSecurityPolicy.GetValue() != "" {
		args.EdgeSecurityPolicy = pulumi.String(spec.EdgeSecurityPolicy.GetValue())
	}

	// Honor the spec contract: an empty project_id falls back to the provider's
	// default project. Leaving Project unset lets the gcp provider resolve its
	// own project (configuration or the GOOGLE_PROJECT / GOOGLE_CLOUD_PROJECT
	// environment chain); an empty string would be sent verbatim and rejected.
	if spec.ProjectId.GetValue() != "" {
		args.Project = pulumi.String(spec.ProjectId.GetValue())
	}

	if len(spec.CustomResponseHeaders) > 0 {
		customResponseHeaders := pulumi.StringArray{}
		for _, customResponseHeader := range spec.CustomResponseHeaders {
			customResponseHeaders = append(customResponseHeaders, pulumi.String(customResponseHeader))
		}
		args.CustomResponseHeaders = customResponseHeaders
	}

	if spec.CdnPolicy != nil {
		cdnPolicy := &compute.BackendBucketCdnPolicyArgs{}

		if spec.CdnPolicy.CacheMode != "" {
			cdnPolicy.CacheMode = pulumi.String(spec.CdnPolicy.CacheMode)
		}
		// The tfvars/proto contract treats 0 as unset for TTLs, letting the
		// GCP API apply its own defaults — identical to the Terraform module.
		if spec.CdnPolicy.ClientTtl != 0 {
			cdnPolicy.ClientTtl = pulumi.Int(int(spec.CdnPolicy.ClientTtl))
		}
		if spec.CdnPolicy.DefaultTtl != 0 {
			cdnPolicy.DefaultTtl = pulumi.Int(int(spec.CdnPolicy.DefaultTtl))
		}
		if spec.CdnPolicy.MaxTtl != 0 {
			cdnPolicy.MaxTtl = pulumi.Int(int(spec.CdnPolicy.MaxTtl))
		}
		if spec.CdnPolicy.NegativeCaching {
			cdnPolicy.NegativeCaching = pulumi.Bool(true)
		}
		if spec.CdnPolicy.ServeWhileStale != 0 {
			cdnPolicy.ServeWhileStale = pulumi.Int(int(spec.CdnPolicy.ServeWhileStale))
		}
		if spec.CdnPolicy.RequestCoalescing {
			cdnPolicy.RequestCoalescing = pulumi.Bool(true)
		}
		if spec.CdnPolicy.SignedUrlCacheMaxAgeSec != 0 {
			cdnPolicy.SignedUrlCacheMaxAgeSec = pulumi.Int(int(spec.CdnPolicy.SignedUrlCacheMaxAgeSec))
		}

		if len(spec.CdnPolicy.NegativeCachingPolicy) > 0 {
			negativeCachingPolicies := compute.BackendBucketCdnPolicyNegativeCachingPolicyArray{}
			for _, negativeCachingPolicy := range spec.CdnPolicy.NegativeCachingPolicy {
				negativeCachingPolicies = append(negativeCachingPolicies,
					&compute.BackendBucketCdnPolicyNegativeCachingPolicyArgs{
						Code: pulumi.Int(int(negativeCachingPolicy.Code)),
						// A 0 TTL means don't-cache-this-code to GCP, so it is
						// passed as-is (unlike the top-level TTLs).
						Ttl: pulumi.Int(int(negativeCachingPolicy.Ttl)),
					})
			}
			cdnPolicy.NegativeCachingPolicies = negativeCachingPolicies
		}

		if spec.CdnPolicy.CacheKeyPolicy != nil {
			cacheKeyPolicy := &compute.BackendBucketCdnPolicyCacheKeyPolicyArgs{}
			if len(spec.CdnPolicy.CacheKeyPolicy.QueryStringWhitelist) > 0 {
				queryStringWhitelist := pulumi.StringArray{}
				for _, queryString := range spec.CdnPolicy.CacheKeyPolicy.QueryStringWhitelist {
					queryStringWhitelist = append(queryStringWhitelist, pulumi.String(queryString))
				}
				// Pulumi pluralizes this field (queryStringWhitelists) vs the
				// provider's singular query_string_whitelist — same wire field.
				cacheKeyPolicy.QueryStringWhitelists = queryStringWhitelist
			}
			if len(spec.CdnPolicy.CacheKeyPolicy.IncludeHttpHeaders) > 0 {
				includeHttpHeaders := pulumi.StringArray{}
				for _, includeHttpHeader := range spec.CdnPolicy.CacheKeyPolicy.IncludeHttpHeaders {
					includeHttpHeaders = append(includeHttpHeaders, pulumi.String(includeHttpHeader))
				}
				cacheKeyPolicy.IncludeHttpHeaders = includeHttpHeaders
			}
			cdnPolicy.CacheKeyPolicy = cacheKeyPolicy
		}

		if len(spec.CdnPolicy.BypassCacheOnRequestHeaders) > 0 {
			bypassCacheOnRequestHeaders := compute.BackendBucketCdnPolicyBypassCacheOnRequestHeaderArray{}
			for _, bypassCacheOnRequestHeader := range spec.CdnPolicy.BypassCacheOnRequestHeaders {
				bypassCacheOnRequestHeaders = append(bypassCacheOnRequestHeaders,
					&compute.BackendBucketCdnPolicyBypassCacheOnRequestHeaderArgs{
						HeaderName: pulumi.String(bypassCacheOnRequestHeader.HeaderName),
					})
			}
			cdnPolicy.BypassCacheOnRequestHeaders = bypassCacheOnRequestHeaders
		}

		args.CdnPolicy = cdnPolicy
	}

	createdBackendBucket, err := compute.NewBackendBucket(ctx, "backend-bucket", args,
		pulumi.Provider(gcpProvider), pulumi.DependsOn([]pulumi.Resource{createdProjectService}))
	if err != nil {
		return errors.Wrap(err, "failed to create backend bucket")
	}

	// Signed-URL keys — folded into this kind rather than modeled as a
	// separate node: keys are never referenced by other resources, GCP caps
	// them at 3 per bucket, and their lifecycle is the bucket's. Each key is
	// immutable in GCP (add/delete only), which is exactly the rotation
	// semantics signed URLs need (add new key -> re-sign -> remove old).
	for _, signedUrlKey := range spec.SignedUrlKeys {
		signedUrlKeyArgs := &compute.BackendBucketSignedUrlKeyArgs{
			Name: pulumi.String(signedUrlKey.Name),
			// Secret material; never surfaced in outputs. ToSecret marks it
			// in the Pulumi state as well.
			KeyValue:      pulumi.ToSecret(pulumi.String(signedUrlKey.KeyValue)).(pulumi.StringOutput),
			BackendBucket: createdBackendBucket.Name,
		}
		if spec.ProjectId.GetValue() != "" {
			signedUrlKeyArgs.Project = pulumi.String(spec.ProjectId.GetValue())
		}
		if _, err := compute.NewBackendBucketSignedUrlKey(ctx, "signed-url-key-"+signedUrlKey.Name,
			signedUrlKeyArgs, pulumi.Provider(gcpProvider)); err != nil {
			return errors.Wrap(err, "failed to create signed url key "+signedUrlKey.Name)
		}
	}

	ctx.Export(OpSelfLink, createdBackendBucket.SelfLink)
	ctx.Export(OpBackendBucketName, createdBackendBucket.Name)
	ctx.Export(OpBucketName, createdBackendBucket.BucketName)

	return nil
}
