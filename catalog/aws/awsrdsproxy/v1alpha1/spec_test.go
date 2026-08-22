package awsrdsproxyv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/planton/shared/foreignkey/v1"
	"google.golang.org/protobuf/proto"
)

func TestAwsRdsProxySpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsRdsProxySpec Validation Suite")
}

func literal(value string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: value},
	}
}

func minimalConfig() *AwsRdsProxySpec {
	return &AwsRdsProxySpec{
		Region:       "us-west-2",
		EngineFamily: "POSTGRESQL",
		RoleArn:      literal("arn:aws:iam::111111111111:role/rds-proxy-secrets"),
		VpcSubnetIds: []*foreignkeyv1.StringValueOrRef{
			literal("subnet-0123456789abcdef0"),
			literal("subnet-0123456789abcdef1"),
		},
		Auth: []*AwsRdsProxyAuth{
			{SecretArn: literal("arn:aws:secretsmanager:us-west-2:111111111111:secret:db-creds-AbCdEf")},
		},
	}
}

var _ = ginkgo.Describe("AwsRdsProxySpec validations", func() {

	ginkgo.Describe("When valid input is passed", func() {

		ginkgo.It("accepts a minimal configuration", func() {
			gomega.Expect(protovalidate.Validate(minimalConfig())).To(gomega.BeNil())
		})

		ginkgo.It("accepts full auth, pool tuning, and network dials", func() {
			spec := minimalConfig()
			spec.Auth[0].IamAuth = "REQUIRED"
			spec.Auth[0].ClientPasswordAuthType = "POSTGRES_SCRAM_SHA_256"
			spec.Auth[0].Username = "app_rw"
			spec.RequireTls = true
			spec.IdleClientTimeout = 900
			spec.DefaultAuthScheme = "IAM_AUTH"
			spec.EndpointNetworkType = "DUAL"
			spec.TargetConnectionNetworkType = "IPV4"
			spec.ConnectionPool = &AwsRdsProxyConnectionPool{
				ConnectionBorrowTimeout:   proto.Int64(0),
				MaxConnectionsPercent:     proto.Int64(90),
				MaxIdleConnectionsPercent: proto.Int64(10),
				SessionPinningFilters:     []string{"EXCLUDE_VARIABLE_SETS"},
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts a read-only endpoint and an instance target", func() {
			spec := minimalConfig()
			spec.Endpoints = []*AwsRdsProxyEndpoint{
				{
					Name:       "readers",
					TargetRole: "READ_ONLY",
					VpcSubnetIds: []*foreignkeyv1.StringValueOrRef{
						literal("subnet-0123456789abcdef0"),
						literal("subnet-0123456789abcdef1"),
					},
				},
			}
			spec.Target = &AwsRdsProxyTarget{
				DbInstanceIdentifier: literal("orders-db"),
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts a cluster target", func() {
			spec := minimalConfig()
			spec.Target = &AwsRdsProxyTarget{
				DbClusterIdentifier: literal("orders-aurora"),
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})
	})

	ginkgo.Describe("When invalid input is passed", func() {

		ginkgo.It("rejects an unknown engine family", func() {
			spec := minimalConfig()
			spec.EngineFamily = "MARIADB"
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a missing role", func() {
			spec := minimalConfig()
			spec.RoleArn = nil
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects fewer than two subnets", func() {
			spec := minimalConfig()
			spec.VpcSubnetIds = spec.VpcSubnetIds[:1]
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects empty auth", func() {
			spec := minimalConfig()
			spec.Auth = nil
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects auth without a secret", func() {
			spec := minimalConfig()
			spec.Auth[0].SecretArn = nil
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an unknown iam_auth mode", func() {
			spec := minimalConfig()
			spec.Auth[0].IamAuth = "MAYBE"
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an out-of-range idle timeout", func() {
			spec := minimalConfig()
			spec.IdleClientTimeout = 30000
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an out-of-range borrow timeout", func() {
			spec := minimalConfig()
			spec.ConnectionPool = &AwsRdsProxyConnectionPool{
				ConnectionBorrowTimeout: proto.Int64(4000),
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an unknown session pinning filter", func() {
			spec := minimalConfig()
			spec.ConnectionPool = &AwsRdsProxyConnectionPool{
				SessionPinningFilters: []string{"EXCLUDE_ALL"},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an endpoint name with consecutive hyphens", func() {
			spec := minimalConfig()
			spec.Endpoints = []*AwsRdsProxyEndpoint{
				{
					Name: "read--only",
					VpcSubnetIds: []*foreignkeyv1.StringValueOrRef{
						literal("subnet-0123456789abcdef0"),
						literal("subnet-0123456789abcdef1"),
					},
				},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an endpoint name ending with a hyphen", func() {
			spec := minimalConfig()
			spec.Endpoints = []*AwsRdsProxyEndpoint{
				{
					Name: "readers-",
					VpcSubnetIds: []*foreignkeyv1.StringValueOrRef{
						literal("subnet-0123456789abcdef0"),
						literal("subnet-0123456789abcdef1"),
					},
				},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects duplicate endpoint names", func() {
			spec := minimalConfig()
			subnets := []*foreignkeyv1.StringValueOrRef{
				literal("subnet-0123456789abcdef0"),
				literal("subnet-0123456789abcdef1"),
			}
			spec.Endpoints = []*AwsRdsProxyEndpoint{
				{Name: "readers", VpcSubnetIds: subnets},
				{Name: "readers", VpcSubnetIds: subnets},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an endpoint with fewer than two subnets", func() {
			spec := minimalConfig()
			spec.Endpoints = []*AwsRdsProxyEndpoint{
				{
					Name:         "readers",
					VpcSubnetIds: []*foreignkeyv1.StringValueOrRef{literal("subnet-0123456789abcdef0")},
				},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a target naming both an instance and a cluster", func() {
			spec := minimalConfig()
			spec.Target = &AwsRdsProxyTarget{
				DbInstanceIdentifier: literal("orders-db"),
				DbClusterIdentifier:  literal("orders-aurora"),
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a target naming neither an instance nor a cluster", func() {
			spec := minimalConfig()
			spec.Target = &AwsRdsProxyTarget{}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})
	})
})
