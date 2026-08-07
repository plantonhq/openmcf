package gcpurlmapv1alpha1

import (
	"strings"
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"github.com/plantonhq/planton/shared"
	foreignkeyv1 "github.com/plantonhq/planton/shared/foreignkey/v1"
)

func TestSuite(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "GcpUrlMapSpec Suite")
}

func backendSelfLink() *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{
			Value: "https://www.googleapis.com/compute/v1/projects/p/global/backendServices/web",
		},
	}
}

var _ = ginkgo.Describe("GcpUrlMapSpec", func() {
	var validator protovalidate.Validator

	ginkgo.BeforeEach(func() {
		var err error
		validator, err = protovalidate.New()
		gomega.Expect(err).ToNot(gomega.HaveOccurred())
	})

	minimal := func() *GcpUrlMap {
		return &GcpUrlMap{
			ApiVersion: "gcp.planton.dev/v1alpha1",
			Kind:       "GcpUrlMap",
			Metadata: &shared.CloudResourceMetadata{
				Name: "test-url-map",
			},
			Spec: &GcpUrlMapSpec{
				DefaultService: backendSelfLink(),
			},
		}
	}

	// ──────────────── Positive Cases ────────────────

	ginkgo.It("should accept a minimal valid spec with default_service", func() {
		gomega.Expect(validator.Validate(minimal())).To(gomega.Succeed())
	})

	ginkgo.It("should accept default_url_redirect for apex-to-www", func() {
		target := minimal()
		target.Spec.DefaultService = nil
		target.Spec.DefaultUrlRedirect = &GcpUrlMapUrlRedirect{
			HostRedirect:  "www.example.com",
			HttpsRedirect: true,
			StripQuery:    false,
		}
		gomega.Expect(validator.Validate(target)).To(gomega.Succeed())
	})

	ginkgo.It("should accept default_route_action with weighted backends", func() {
		target := minimal()
		target.Spec.DefaultService = nil
		target.Spec.DefaultRouteAction = &GcpUrlMapRouteAction{
			WeightedBackendServices: []*GcpUrlMapWeightedBackendService{
				{BackendService: backendSelfLink(), Weight: 900},
				{BackendService: backendSelfLink(), Weight: 100},
			},
		}
		gomega.Expect(validator.Validate(target)).To(gomega.Succeed())
	})

	ginkgo.It("should accept host rules and path_rules fan-out", func() {
		target := minimal()
		target.Spec.HostRules = []*GcpUrlMapHostRule{
			{Hosts: []string{"www.example.com"}, PathMatcher: "main"},
		}
		target.Spec.PathMatchers = []*GcpUrlMapPathMatcher{
			{
				Name:           "main",
				DefaultService: backendSelfLink(),
				PathRules: []*GcpUrlMapPathRule{
					{Paths: []string{"/api/*"}, Service: backendSelfLink()},
				},
			},
		}
		gomega.Expect(validator.Validate(target)).To(gomega.Succeed())
	})

	ginkgo.It("should accept route_rules with rich match rules", func() {
		target := minimal()
		target.Spec.PathMatchers = []*GcpUrlMapPathMatcher{
			{
				Name: "api",
				RouteRules: []*GcpUrlMapRouteRule{
					{
						Priority: 1,
						Service:  backendSelfLink(),
						MatchRules: []*GcpUrlMapRouteRuleMatchRule{
							{
								PrefixMatch: "/v1",
								HeaderMatches: []*GcpUrlMapHeaderMatch{
									{HeaderName: "X-Api-Version", ExactMatch: "1"},
								},
							},
						},
					},
				},
			},
		}
		gomega.Expect(validator.Validate(target)).To(gomega.Succeed())
	})

	ginkgo.It("should accept path_template_rewrite inside a route rule route_action", func() {
		target := minimal()
		target.Spec.PathMatchers = []*GcpUrlMapPathMatcher{
			{
				Name: "tmpl",
				RouteRules: []*GcpUrlMapRouteRule{
					{
						Priority: 1,
						RouteAction: &GcpUrlMapRouteAction{
							UrlRewrite: &GcpUrlMapUrlRewrite{
								PathTemplateRewrite: "/v2/{country}",
							},
						},
						MatchRules: []*GcpUrlMapRouteRuleMatchRule{
							{PathTemplateMatch: "/v1/{country}/**"},
						},
					},
				},
			},
		}
		gomega.Expect(validator.Validate(target)).To(gomega.Succeed())
	})

	ginkgo.It("should accept header_action at the URL map level", func() {
		target := minimal()
		target.Spec.HeaderAction = &GcpUrlMapHeaderAction{
			RequestHeadersToAdd: []*GcpUrlMapHeaderValue{
				{HeaderName: "X-Edge", HeaderValue: "global", Replace: true},
			},
			RequestHeadersToRemove: []string{"X-Internal"},
		}
		gomega.Expect(validator.Validate(target)).To(gomega.Succeed())
	})

	ginkgo.It("should accept custom error response policy with override code", func() {
		target := minimal()
		target.Spec.DefaultCustomErrorResponsePolicy = &GcpUrlMapCustomErrorResponsePolicy{
			ErrorService: &foreignkeyv1.StringValueOrRef{
				LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{
					Value: "https://www.googleapis.com/compute/v1/projects/p/global/backendBuckets/errors",
				},
			},
			ErrorResponseRules: []*GcpUrlMapCustomErrorResponseRule{
				{MatchResponseCodes: []string{"404", "5xx"}, Path: "/errors/404.html", OverrideResponseCode: 200},
			},
		}
		gomega.Expect(validator.Validate(target)).To(gomega.Succeed())
	})

	ginkgo.It("should accept routing self-tests", func() {
		target := minimal()
		target.Spec.Tests = []*GcpUrlMapTest{
			{
				Host:        "www.example.com",
				Path:        "/",
				Service:     backendSelfLink(),
				Description: "home resolves to web backend",
				Headers:     []*GcpUrlMapTestHeader{{Name: "Accept", Value: "text/html"}},
			},
		}
		gomega.Expect(validator.Validate(target)).To(gomega.Succeed())
	})

	ginkgo.It("should accept metadata filters on route rule match rules", func() {
		target := minimal()
		target.Spec.PathMatchers = []*GcpUrlMapPathMatcher{
			{
				Name: "mesh",
				RouteRules: []*GcpUrlMapRouteRule{
					{
						Priority: 10,
						Service:  backendSelfLink(),
						MatchRules: []*GcpUrlMapRouteRuleMatchRule{
							{
								PrefixMatch: "/",
								MetadataFilters: []*GcpUrlMapMetadataFilter{
									{
										FilterMatchCriteria: "MATCH_ALL",
										FilterLabels: []*GcpUrlMapMetadataFilterLabel{
											{Name: "region", Value: "us"},
										},
									},
								},
							},
						},
					},
				},
			},
		}
		gomega.Expect(validator.Validate(target)).To(gomega.Succeed())
	})

	ginkgo.It("should accept a path matcher with default_url_redirect only", func() {
		target := minimal()
		target.Spec.PathMatchers = []*GcpUrlMapPathMatcher{
			{
				Name: "redirect-matcher",
				DefaultUrlRedirect: &GcpUrlMapUrlRedirect{
					PrefixRedirect: "/landing",
					StripQuery:     true,
				},
			},
		}
		gomega.Expect(validator.Validate(target)).To(gomega.Succeed())
	})

	ginkgo.It("should accept path_rule route_action with weighted backends", func() {
		target := minimal()
		target.Spec.PathMatchers = []*GcpUrlMapPathMatcher{
			{
				Name: "canary",
				PathRules: []*GcpUrlMapPathRule{
					{
						Paths: []string{"/release/*"},
						RouteAction: &GcpUrlMapRouteAction{
							WeightedBackendServices: []*GcpUrlMapWeightedBackendService{
								{BackendService: backendSelfLink(), Weight: 500},
								{BackendService: backendSelfLink(), Weight: 500},
							},
						},
					},
				},
			},
		}
		gomega.Expect(validator.Validate(target)).To(gomega.Succeed())
	})

	// ──────────────── Negative Cases ────────────────

	ginkgo.It("should reject a spec with no default target", func() {
		target := minimal()
		target.Spec.DefaultService = nil
		err := validator.Validate(target)
		gomega.Expect(err).To(gomega.HaveOccurred())
		gomega.Expect(strings.Contains(err.Error(), "exactly one default target")).To(gomega.BeTrue())
	})

	ginkgo.It("should reject both default_service and default_url_redirect", func() {
		target := minimal()
		target.Spec.DefaultUrlRedirect = &GcpUrlMapUrlRedirect{HttpsRedirect: true, StripQuery: false}
		err := validator.Validate(target)
		gomega.Expect(err).To(gomega.HaveOccurred())
	})

	ginkgo.It("should reject both default_route_action and default_url_redirect", func() {
		target := minimal()
		target.Spec.DefaultService = nil
		target.Spec.DefaultUrlRedirect = &GcpUrlMapUrlRedirect{StripQuery: false}
		target.Spec.DefaultRouteAction = &GcpUrlMapRouteAction{
			WeightedBackendServices: []*GcpUrlMapWeightedBackendService{
				{BackendService: backendSelfLink(), Weight: 100},
			},
		}
		err := validator.Validate(target)
		gomega.Expect(err).To(gomega.HaveOccurred())
		gomega.Expect(strings.Contains(err.Error(), "mutually exclusive")).To(gomega.BeTrue())
	})

	ginkgo.It("should reject default_route_action without weighted backends at top level", func() {
		target := minimal()
		target.Spec.DefaultService = nil
		target.Spec.DefaultRouteAction = &GcpUrlMapRouteAction{
			UrlRewrite: &GcpUrlMapUrlRewrite{HostRewrite: "internal.example.com"},
		}
		err := validator.Validate(target)
		gomega.Expect(err).To(gomega.HaveOccurred())
	})

	ginkgo.It("should reject an invalid url_map_name", func() {
		target := minimal()
		target.Spec.UrlMapName = "Invalid_Name"
		err := validator.Validate(target)
		gomega.Expect(err).To(gomega.HaveOccurred())
		gomega.Expect(strings.Contains(err.Error(), "RFC1035")).To(gomega.BeTrue())
	})

	ginkgo.It("should reject path_redirect and prefix_redirect together", func() {
		target := minimal()
		target.Spec.DefaultService = nil
		target.Spec.DefaultUrlRedirect = &GcpUrlMapUrlRedirect{
			PathRedirect:   "/new",
			PrefixRedirect: "/old",
			StripQuery:     false,
		}
		err := validator.Validate(target)
		gomega.Expect(err).To(gomega.HaveOccurred())
		gomega.Expect(strings.Contains(err.Error(), "mutually exclusive")).To(gomega.BeTrue())
	})

	ginkgo.It("should reject an invalid redirect_response_code", func() {
		target := minimal()
		target.Spec.DefaultService = nil
		target.Spec.DefaultUrlRedirect = &GcpUrlMapUrlRedirect{
			RedirectResponseCode: "INVALID",
			StripQuery:           false,
		}
		err := validator.Validate(target)
		gomega.Expect(err).To(gomega.HaveOccurred())
	})

	ginkgo.It("should reject a path matcher with both path_rules and route_rules", func() {
		target := minimal()
		target.Spec.PathMatchers = []*GcpUrlMapPathMatcher{
			{
				Name: "mixed",
				PathRules: []*GcpUrlMapPathRule{
					{Paths: []string{"/a/*"}, Service: backendSelfLink()},
				},
				RouteRules: []*GcpUrlMapRouteRule{
					{Priority: 1, Service: backendSelfLink(), MatchRules: []*GcpUrlMapRouteRuleMatchRule{{PrefixMatch: "/b"}}},
				},
			},
		}
		err := validator.Validate(target)
		gomega.Expect(err).To(gomega.HaveOccurred())
		gomega.Expect(strings.Contains(err.Error(), "either path_rules or route_rules")).To(gomega.BeTrue())
	})

	ginkgo.It("should reject a path rule with no target", func() {
		target := minimal()
		target.Spec.PathMatchers = []*GcpUrlMapPathMatcher{
			{
				Name: "bad-rule",
				PathRules: []*GcpUrlMapPathRule{
					{Paths: []string{"/orphan/*"}},
				},
			},
		}
		err := validator.Validate(target)
		gomega.Expect(err).To(gomega.HaveOccurred())
	})

	ginkgo.It("should reject a path rule with multiple targets", func() {
		target := minimal()
		target.Spec.PathMatchers = []*GcpUrlMapPathMatcher{
			{
				Name: "multi",
				PathRules: []*GcpUrlMapPathRule{
					{
						Paths:       []string{"/both/*"},
						Service:     backendSelfLink(),
						UrlRedirect: &GcpUrlMapUrlRedirect{StripQuery: false},
					},
				},
			},
		}
		err := validator.Validate(target)
		gomega.Expect(err).To(gomega.HaveOccurred())
	})

	ginkgo.It("should reject a route rule with no target", func() {
		target := minimal()
		target.Spec.PathMatchers = []*GcpUrlMapPathMatcher{
			{
				Name: "empty-route",
				RouteRules: []*GcpUrlMapRouteRule{
					{Priority: 1, MatchRules: []*GcpUrlMapRouteRuleMatchRule{{PrefixMatch: "/x"}}},
				},
			},
		}
		err := validator.Validate(target)
		gomega.Expect(err).To(gomega.HaveOccurred())
	})

	ginkgo.It("should reject multiple path matchers on one match rule", func() {
		target := minimal()
		target.Spec.PathMatchers = []*GcpUrlMapPathMatcher{
			{
				Name: "bad-match",
				RouteRules: []*GcpUrlMapRouteRule{
					{
						Priority: 1,
						Service:  backendSelfLink(),
						MatchRules: []*GcpUrlMapRouteRuleMatchRule{
							{PrefixMatch: "/a", FullPathMatch: "/b"},
						},
					},
				},
			},
		}
		err := validator.Validate(target)
		gomega.Expect(err).To(gomega.HaveOccurred())
	})

	ginkgo.It("should reject path_prefix_rewrite and path_template_rewrite together", func() {
		target := minimal()
		target.Spec.PathMatchers = []*GcpUrlMapPathMatcher{
			{
				Name: "bad-rewrite",
				RouteRules: []*GcpUrlMapRouteRule{
					{
						Priority: 1,
						RouteAction: &GcpUrlMapRouteAction{
							UrlRewrite: &GcpUrlMapUrlRewrite{
								PathPrefixRewrite:   "/v2",
								PathTemplateRewrite: "/v2/{id}",
							},
						},
						MatchRules: []*GcpUrlMapRouteRuleMatchRule{{PrefixMatch: "/v1"}},
					},
				},
			},
		}
		err := validator.Validate(target)
		gomega.Expect(err).To(gomega.HaveOccurred())
	})

	ginkgo.It("should reject a path matcher with multiple default targets", func() {
		target := minimal()
		target.Spec.PathMatchers = []*GcpUrlMapPathMatcher{
			{
				Name:           "two-defaults",
				DefaultService: backendSelfLink(),
				DefaultUrlRedirect: &GcpUrlMapUrlRedirect{
					StripQuery: false,
				},
			},
		}
		err := validator.Validate(target)
		gomega.Expect(err).To(gomega.HaveOccurred())
	})

	ginkgo.It("should reject an invalid filter_match_criteria", func() {
		target := minimal()
		target.Spec.PathMatchers = []*GcpUrlMapPathMatcher{
			{
				Name: "bad-filter",
				RouteRules: []*GcpUrlMapRouteRule{
					{
						Priority: 1,
						Service:  backendSelfLink(),
						MatchRules: []*GcpUrlMapRouteRuleMatchRule{
							{
								PrefixMatch: "/",
								MetadataFilters: []*GcpUrlMapMetadataFilter{
									{
										FilterMatchCriteria: "MATCH_SOME",
										FilterLabels:        []*GcpUrlMapMetadataFilterLabel{{Name: "k", Value: "v"}},
									},
								},
							},
						},
					},
				},
			},
		}
		err := validator.Validate(target)
		gomega.Expect(err).To(gomega.HaveOccurred())
	})

	ginkgo.It("should reject override_response_code outside 200-599", func() {
		target := minimal()
		target.Spec.DefaultCustomErrorResponsePolicy = &GcpUrlMapCustomErrorResponsePolicy{
			ErrorResponseRules: []*GcpUrlMapCustomErrorResponseRule{
				{MatchResponseCodes: []string{"404"}, Path: "/e.html", OverrideResponseCode: 199},
			},
		}
		err := validator.Validate(target)
		gomega.Expect(err).To(gomega.HaveOccurred())
	})

	ginkgo.It("should reject weighted backend weight above 1000", func() {
		target := minimal()
		target.Spec.DefaultService = nil
		target.Spec.DefaultRouteAction = &GcpUrlMapRouteAction{
			WeightedBackendServices: []*GcpUrlMapWeightedBackendService{
				{BackendService: backendSelfLink(), Weight: 1001},
			},
		}
		err := validator.Validate(target)
		gomega.Expect(err).To(gomega.HaveOccurred())
	})

	ginkgo.It("should reject a route rule without match_rules", func() {
		target := minimal()
		target.Spec.PathMatchers = []*GcpUrlMapPathMatcher{
			{
				Name: "no-match",
				RouteRules: []*GcpUrlMapRouteRule{
					{Priority: 1, Service: backendSelfLink()},
				},
			},
		}
		err := validator.Validate(target)
		gomega.Expect(err).To(gomega.HaveOccurred())
	})
})
