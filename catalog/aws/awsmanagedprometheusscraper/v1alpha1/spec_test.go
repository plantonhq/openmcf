package awsmanagedprometheusscraperv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/planton/shared/foreignkey/v1"
)

func TestAwsManagedPrometheusScraperSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsManagedPrometheusScraperSpec Validation Suite")
}

func svr(val string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: val},
	}
}

func eksScraper() *AwsManagedPrometheusScraperSpec {
	return &AwsManagedPrometheusScraperSpec{
		Region: "us-east-1",
		SourceEks: &AwsManagedPrometheusScraperEksSource{
			ClusterArn: svr("arn:aws:eks:us-east-1:123456789012:cluster/platform"),
			SubnetIds:  []*foreignkeyv1.StringValueOrRef{svr("subnet-0abc"), svr("subnet-0def")},
		},
		AmpWorkspaceArn: svr("arn:aws:aps:us-east-1:123456789012:workspace/ws-abc"),
	}
}

func vpcScraper() *AwsManagedPrometheusScraperSpec {
	return &AwsManagedPrometheusScraperSpec{
		Region: "us-east-1",
		SourceVpc: &AwsManagedPrometheusScraperVpcSource{
			SubnetIds:        []*foreignkeyv1.StringValueOrRef{svr("subnet-0abc"), svr("subnet-0def")},
			SecurityGroupIds: []*foreignkeyv1.StringValueOrRef{svr("sg-0abc")},
		},
		CloudwatchDatasetArn: svr("arn:aws:cloudwatch:us-east-1:123456789012:dataset/prom"),
		ScrapeConfiguration:  "scrape_configs:\n  - job_name: app\n    static_configs:\n      - targets: ['10.0.0.5:9090']\n",
	}
}

var _ = ginkgo.Describe("AwsManagedPrometheusScraperSpec validations", func() {

	ginkgo.Describe("When valid input is passed", func() {

		ginkgo.It("accepts an EKS-sourced scraper without a scrape configuration (AWS default)", func() {
			gomega.Expect(protovalidate.Validate(eksScraper())).To(gomega.BeNil())
		})

		ginkgo.It("accepts a VPC-sourced scraper with its own scrape configuration", func() {
			gomega.Expect(protovalidate.Validate(vpcScraper())).To(gomega.BeNil())
		})

		ginkgo.It("accepts the cross-account role pair together", func() {
			spec := eksScraper()
			spec.RoleConfiguration = &AwsManagedPrometheusScraperRoleConfiguration{
				SourceRoleArn: svr("arn:aws:iam::123456789012:role/scrape-source"),
				TargetRoleArn: svr("arn:aws:iam::210987654321:role/scrape-target"),
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts scraper logging with components", func() {
			spec := eksScraper()
			spec.Logging = &AwsManagedPrometheusScraperLogging{
				LogGroupArn: svr("arn:aws:logs:us-east-1:123456789012:log-group:/amp/scraper"),
				Components:  []string{"COLLECTOR", "EXPORTER"},
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})
	})

	ginkgo.Describe("When invalid input is passed", func() {

		ginkgo.It("rejects both sources at once", func() {
			spec := eksScraper()
			spec.SourceVpc = vpcScraper().SourceVpc
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a single subnet (CreateScraper requires at least two)", func() {
			spec := vpcScraper()
			spec.SourceVpc.SubnetIds = spec.SourceVpc.SubnetIds[:1]
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())

			eks := eksScraper()
			eks.SourceEks.SubnetIds = eks.SourceEks.SubnetIds[:1]
			gomega.Expect(protovalidate.Validate(eks)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a spec with no source", func() {
			spec := eksScraper()
			spec.SourceEks = nil
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects both destinations at once", func() {
			spec := eksScraper()
			spec.CloudwatchDatasetArn = svr("arn:aws:cloudwatch:us-east-1:123456789012:dataset/prom")
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a spec with no destination", func() {
			spec := eksScraper()
			spec.AmpWorkspaceArn = nil
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a VPC source without a scrape configuration", func() {
			spec := vpcScraper()
			spec.ScrapeConfiguration = ""
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a VPC source without security groups", func() {
			spec := vpcScraper()
			spec.SourceVpc.SecurityGroupIds = nil
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a lone cross-account role", func() {
			spec := eksScraper()
			spec.RoleConfiguration = &AwsManagedPrometheusScraperRoleConfiguration{
				SourceRoleArn: svr("arn:aws:iam::123456789012:role/scrape-source"),
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an unknown logging component", func() {
			spec := eksScraper()
			spec.Logging = &AwsManagedPrometheusScraperLogging{
				LogGroupArn: svr("arn:aws:logs:us-east-1:123456789012:log-group:/amp/scraper"),
				Components:  []string{"INGESTER"},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})
	})
})
