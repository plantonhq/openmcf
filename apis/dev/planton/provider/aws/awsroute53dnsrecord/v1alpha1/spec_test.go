package awsroute53dnsrecordv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/planton/apis/dev/planton/shared/foreignkey/v1"
)

func TestAwsRoute53DnsRecordSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsRoute53DnsRecordSpec Validation Suite")
}

// helper to create a StringValueOrRef with a literal value.
func strRef(val string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: val},
	}
}

// helper to build a valid alias target (an ALB pair).
func albAlias() *AwsRoute53AliasTarget {
	return &AwsRoute53AliasTarget{
		DnsName: strRef("my-alb-1234567890.us-east-1.elb.amazonaws.com"),
		ZoneId:  strRef("Z35SXDOTRQ7X7K"),
	}
}

var _ = ginkgo.Describe("AwsRoute53DnsRecordSpec validations", func() {
	var spec *AwsRoute53DnsRecordSpec

	ginkgo.BeforeEach(func() {
		// Minimal valid spec: a simple A record.
		spec = &AwsRoute53DnsRecordSpec{
			Region: "us-west-2",
			ZoneId: strRef("Z1234567890ABC"),
			Name:   "www.example.com",
			Type:   "A",
			Ttl:    300,
			Values: []string{"192.0.2.1"},
		}
	})

	// -------------------------------------------------------------------------
	// Happy path: standard records
	// -------------------------------------------------------------------------

	ginkgo.It("accepts a simple A record", func() {
		gomega.Expect(protovalidate.Validate(spec)).To(gomega.Succeed())
	})

	ginkgo.It("accepts every supported record type", func() {
		for _, t := range []string{"A", "AAAA", "CAA", "CNAME", "DS", "HTTPS", "MX", "NAPTR", "NS", "PTR", "SOA", "SPF", "SRV", "SSHFP", "SVCB", "TLSA", "TXT"} {
			spec.Type = t
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.Succeed(), "type %s should be valid", t)
		}
	})

	ginkgo.It("accepts a multi-value A record and a wildcard name", func() {
		spec.Name = "*.example.com"
		spec.Values = []string{"192.0.2.1", "192.0.2.2"}
		gomega.Expect(protovalidate.Validate(spec)).To(gomega.Succeed())
	})

	ginkgo.It("accepts underscore-prefixed service-record names", func() {
		for _, n := range []string{"_dmarc.example.com", "token._domainkey.example.com", "_sip._tcp.example.com"} {
			spec.Name = n
			spec.Type = "TXT"
			spec.Values = []string{"v=DMARC1; p=none"}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.Succeed(), "name %s should be valid", n)
		}
	})

	ginkgo.It("accepts allow_overwrite", func() {
		spec.AllowOverwrite = true
		gomega.Expect(protovalidate.Validate(spec)).To(gomega.Succeed())
	})

	// -------------------------------------------------------------------------
	// Happy path: alias records
	// -------------------------------------------------------------------------

	ginkgo.It("accepts an alias record to an ALB", func() {
		spec.Ttl = 0
		spec.Values = nil
		spec.AliasTarget = albAlias()
		gomega.Expect(protovalidate.Validate(spec)).To(gomega.Succeed())
	})

	ginkgo.It("accepts an alias record with target-health evaluation", func() {
		spec.Ttl = 0
		spec.Values = nil
		spec.AliasTarget = albAlias()
		spec.AliasTarget.EvaluateTargetHealth = true
		gomega.Expect(protovalidate.Validate(spec)).To(gomega.Succeed())
	})

	// -------------------------------------------------------------------------
	// Happy path: routing policies
	// -------------------------------------------------------------------------

	ginkgo.It("accepts weighted routing with a set identifier", func() {
		spec.SetIdentifier = "weight-70"
		spec.RoutingPolicy = &AwsRoute53RoutingPolicy{
			Policy: &AwsRoute53RoutingPolicy_Weighted{Weighted: &AwsRoute53WeightedPolicy{Weight: 70}},
		}
		gomega.Expect(protovalidate.Validate(spec)).To(gomega.Succeed())
	})

	ginkgo.It("accepts latency routing", func() {
		spec.SetIdentifier = "us-east-1"
		spec.RoutingPolicy = &AwsRoute53RoutingPolicy{
			Policy: &AwsRoute53RoutingPolicy_Latency{Latency: &AwsRoute53LatencyPolicy{Region: "us-east-1"}},
		}
		gomega.Expect(protovalidate.Validate(spec)).To(gomega.Succeed())
	})

	ginkgo.It("accepts failover routing with a health check ref", func() {
		spec.SetIdentifier = "primary"
		spec.HealthCheckId = strRef("abcdef11-2222-3333-4444-555555fedcba")
		spec.RoutingPolicy = &AwsRoute53RoutingPolicy{
			Policy: &AwsRoute53RoutingPolicy_Failover{Failover: &AwsRoute53FailoverPolicy{FailoverType: "PRIMARY"}},
		}
		gomega.Expect(protovalidate.Validate(spec)).To(gomega.Succeed())
	})

	ginkgo.It("accepts geolocation routing by country", func() {
		spec.SetIdentifier = "eu-users"
		spec.RoutingPolicy = &AwsRoute53RoutingPolicy{
			Policy: &AwsRoute53RoutingPolicy_Geolocation{Geolocation: &AwsRoute53GeolocationPolicy{Country: "DE"}},
		}
		gomega.Expect(protovalidate.Validate(spec)).To(gomega.Succeed())
	})

	ginkgo.It("accepts geoproximity routing by region with bias", func() {
		spec.SetIdentifier = "us-east"
		spec.RoutingPolicy = &AwsRoute53RoutingPolicy{
			Policy: &AwsRoute53RoutingPolicy_Geoproximity{Geoproximity: &AwsRoute53GeoproximityPolicy{
				AwsRegion: "us-east-1",
				Bias:      25,
			}},
		}
		gomega.Expect(protovalidate.Validate(spec)).To(gomega.Succeed())
	})

	ginkgo.It("accepts geoproximity routing by coordinates", func() {
		spec.SetIdentifier = "nyc-dc"
		spec.RoutingPolicy = &AwsRoute53RoutingPolicy{
			Policy: &AwsRoute53RoutingPolicy_Geoproximity{Geoproximity: &AwsRoute53GeoproximityPolicy{
				Coordinates: &AwsRoute53Coordinates{Latitude: "40.71", Longitude: "-74.01"},
			}},
		}
		gomega.Expect(protovalidate.Validate(spec)).To(gomega.Succeed())
	})

	ginkgo.It("accepts CIDR routing", func() {
		spec.SetIdentifier = "office-network"
		spec.RoutingPolicy = &AwsRoute53RoutingPolicy{
			Policy: &AwsRoute53RoutingPolicy_Cidr{Cidr: &AwsRoute53CidrPolicy{
				CollectionId: "cf1234ab-cdef-5678-90ab-cdef12345678",
				LocationName: "office",
			}},
		}
		gomega.Expect(protovalidate.Validate(spec)).To(gomega.Succeed())
	})

	ginkgo.It("accepts multivalue answer routing on a standard record", func() {
		spec.SetIdentifier = "server-1"
		spec.RoutingPolicy = &AwsRoute53RoutingPolicy{
			Policy: &AwsRoute53RoutingPolicy_MultivalueAnswer{MultivalueAnswer: &AwsRoute53MultivalueAnswerPolicy{}},
		}
		gomega.Expect(protovalidate.Validate(spec)).To(gomega.Succeed())
	})

	// -------------------------------------------------------------------------
	// Required fields
	// -------------------------------------------------------------------------

	ginkgo.It("rejects a missing zone_id", func() {
		spec.ZoneId = nil
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.Succeed())
	})

	ginkgo.It("rejects a missing name", func() {
		spec.Name = ""
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.Succeed())
	})

	ginkgo.It("rejects a missing type", func() {
		spec.Type = ""
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.Succeed())
	})

	ginkgo.It("rejects an unknown record type", func() {
		spec.Type = "ALIAS"
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.Succeed())
	})

	ginkgo.It("rejects an invalid record name", func() {
		spec.Name = "not a domain!"
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.Succeed())
	})

	// -------------------------------------------------------------------------
	// Standard XOR alias
	// -------------------------------------------------------------------------

	ginkgo.It("rejects a record with neither values nor alias", func() {
		spec.Values = nil
		spec.Ttl = 0
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.Succeed())
	})

	ginkgo.It("rejects a record with both values and alias", func() {
		spec.AliasTarget = albAlias()
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.Succeed())
	})

	ginkgo.It("rejects a ttl on an alias record", func() {
		spec.Values = nil
		spec.AliasTarget = albAlias()
		spec.Ttl = 300
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.Succeed())
	})

	ginkgo.It("rejects a standard record without a ttl", func() {
		spec.Ttl = 0
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.Succeed())
	})

	ginkgo.It("rejects a ttl above one week", func() {
		spec.Ttl = 604801
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.Succeed())
	})

	ginkgo.It("rejects an alias missing its target zone id", func() {
		spec.Values = nil
		spec.Ttl = 0
		spec.AliasTarget = &AwsRoute53AliasTarget{
			DnsName: strRef("d1234abcd.cloudfront.net"),
		}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.Succeed())
	})

	// -------------------------------------------------------------------------
	// Routing policy contracts
	// -------------------------------------------------------------------------

	ginkgo.It("rejects a routing policy without a set identifier", func() {
		spec.RoutingPolicy = &AwsRoute53RoutingPolicy{
			Policy: &AwsRoute53RoutingPolicy_Weighted{Weighted: &AwsRoute53WeightedPolicy{Weight: 50}},
		}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.Succeed())
	})

	ginkgo.It("rejects a weight above 255", func() {
		spec.SetIdentifier = "w"
		spec.RoutingPolicy = &AwsRoute53RoutingPolicy{
			Policy: &AwsRoute53RoutingPolicy_Weighted{Weighted: &AwsRoute53WeightedPolicy{Weight: 256}},
		}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.Succeed())
	})

	ginkgo.It("rejects latency routing without a region", func() {
		spec.SetIdentifier = "l"
		spec.RoutingPolicy = &AwsRoute53RoutingPolicy{
			Policy: &AwsRoute53RoutingPolicy_Latency{Latency: &AwsRoute53LatencyPolicy{}},
		}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.Succeed())
	})

	ginkgo.It("rejects an invalid failover type", func() {
		spec.SetIdentifier = "p"
		spec.RoutingPolicy = &AwsRoute53RoutingPolicy{
			Policy: &AwsRoute53RoutingPolicy_Failover{Failover: &AwsRoute53FailoverPolicy{FailoverType: "ACTIVE"}},
		}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.Succeed())
	})

	ginkgo.It("rejects geolocation routing without any location", func() {
		spec.SetIdentifier = "g"
		spec.RoutingPolicy = &AwsRoute53RoutingPolicy{
			Policy: &AwsRoute53RoutingPolicy_Geolocation{Geolocation: &AwsRoute53GeolocationPolicy{}},
		}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.Succeed())
	})

	ginkgo.It("rejects geoproximity with two location determinants", func() {
		spec.SetIdentifier = "gp"
		spec.RoutingPolicy = &AwsRoute53RoutingPolicy{
			Policy: &AwsRoute53RoutingPolicy_Geoproximity{Geoproximity: &AwsRoute53GeoproximityPolicy{
				AwsRegion:      "us-east-1",
				LocalZoneGroup: "us-east-1-bue-1",
			}},
		}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.Succeed())
	})

	ginkgo.It("rejects geoproximity with no location determinant", func() {
		spec.SetIdentifier = "gp"
		spec.RoutingPolicy = &AwsRoute53RoutingPolicy{
			Policy: &AwsRoute53RoutingPolicy_Geoproximity{Geoproximity: &AwsRoute53GeoproximityPolicy{Bias: 10}},
		}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.Succeed())
	})

	ginkgo.It("rejects a geoproximity bias outside -99..99", func() {
		spec.SetIdentifier = "gp"
		spec.RoutingPolicy = &AwsRoute53RoutingPolicy{
			Policy: &AwsRoute53RoutingPolicy_Geoproximity{Geoproximity: &AwsRoute53GeoproximityPolicy{
				AwsRegion: "us-east-1",
				Bias:      100,
			}},
		}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.Succeed())
	})

	ginkgo.It("rejects CIDR routing without a location name", func() {
		spec.SetIdentifier = "c"
		spec.RoutingPolicy = &AwsRoute53RoutingPolicy{
			Policy: &AwsRoute53RoutingPolicy_Cidr{Cidr: &AwsRoute53CidrPolicy{CollectionId: "cf1234"}},
		}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.Succeed())
	})

	ginkgo.It("rejects multivalue answer routing on an alias record", func() {
		spec.Values = nil
		spec.Ttl = 0
		spec.AliasTarget = albAlias()
		spec.SetIdentifier = "server-1"
		spec.RoutingPolicy = &AwsRoute53RoutingPolicy{
			Policy: &AwsRoute53RoutingPolicy_MultivalueAnswer{MultivalueAnswer: &AwsRoute53MultivalueAnswerPolicy{}},
		}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.Succeed())
	})
})
