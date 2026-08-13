package awsbedrockagentcoretoolsv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/planton/shared/foreignkey/v1"
)

func TestAwsBedrockAgentCoreToolsSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsBedrockAgentCoreToolsSpec Validation Suite")
}

func svr(val string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: val},
	}
}

// minimalTools is the smallest valid manifest: region plus one sandboxed
// code interpreter.
func minimalTools() *AwsBedrockAgentCoreToolsSpec {
	return &AwsBedrockAgentCoreToolsSpec{
		Region: "us-west-2",
		CodeInterpreters: []*AwsBedrockAgentCoreCodeInterpreter{
			{
				Name:    "python_sandbox",
				Network: &AwsBedrockAgentCoreCodeInterpreterNetwork{Mode: "SANDBOX"},
			},
		},
	}
}

var _ = ginkgo.Describe("AwsBedrockAgentCoreToolsSpec validations", func() {

	ginkgo.Describe("When valid input is passed", func() {

		ginkgo.Context("with minimal required fields", func() {
			ginkgo.It("should not return a validation error", func() {
				err := protovalidate.Validate(minimalTools())
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with the full surface configured", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := minimalTools()
				spec.Browsers = []*AwsBedrockAgentCoreBrowser{
					{
						Name:             "research-browser",
						Description:      "web research sessions",
						ExecutionRoleArn: svr("arn:aws:iam::123456789012:role/agentcore-browser"),
						Network:          &AwsBedrockAgentCoreBrowserNetwork{Mode: "PUBLIC"},
						SigningEnabled:   boolPtr(true),
						Recording: &AwsBedrockAgentCoreBrowserRecording{
							Enabled: boolPtr(true),
							S3Location: &AwsBedrockAgentCoreToolS3Location{
								Bucket: svr("recordings-bucket"),
								Prefix: "browser-sessions/",
							},
						},
						EnterprisePolicies: []*AwsBedrockAgentCoreBrowserEnterprisePolicy{
							{
								Type: "MANAGED",
								S3: &AwsBedrockAgentCoreToolS3Object{
									Bucket: svr("policies-bucket"),
									Prefix: "chrome/policy.json",
								},
							},
						},
						Certificates: []*AwsBedrockAgentCoreToolCertificate{
							{SecretArn: svr("arn:aws:secretsmanager:us-west-2:123456789012:secret:mtls-cert")},
						},
					},
					{
						Name: "vpc-browser",
						Network: &AwsBedrockAgentCoreBrowserNetwork{
							Mode: "VPC",
							VpcConfig: &AwsBedrockAgentCoreToolVpcConfig{
								Subnets:        []*foreignkeyv1.StringValueOrRef{svr("subnet-abc")},
								SecurityGroups: []*foreignkeyv1.StringValueOrRef{svr("sg-abc")},
							},
						},
					},
				}
				spec.BrowserProfiles = []*AwsBedrockAgentCoreBrowserProfile{
					{Name: "logged_in_docs", Description: "session state with docs site login"},
				}
				spec.CodeInterpreters[0].Description = "runs model-written python"
				spec.CodeInterpreters[0].ExecutionRoleArn = svr("arn:aws:iam::123456789012:role/agentcore-ci")
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})
	})

	ginkgo.Describe("When invalid input is passed", func() {

		ginkgo.Context("with an empty bundle", func() {
			ginkgo.It("should return a validation error", func() {
				spec := &AwsBedrockAgentCoreToolsSpec{Region: "us-west-2"}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("with a browser missing its network", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalTools()
				spec.Browsers = []*AwsBedrockAgentCoreBrowser{{Name: "b"}}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("with a browser in VPC mode and no vpc_config", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalTools()
				spec.Browsers = []*AwsBedrockAgentCoreBrowser{
					{Name: "b", Network: &AwsBedrockAgentCoreBrowserNetwork{Mode: "VPC"}},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("with a browser in SANDBOX mode", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalTools()
				spec.Browsers = []*AwsBedrockAgentCoreBrowser{
					{Name: "b", Network: &AwsBedrockAgentCoreBrowserNetwork{Mode: "SANDBOX"}},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("with duplicate browser names", func() {
			ginkgo.It("should return a validation error", func() {
				b := &AwsBedrockAgentCoreBrowser{
					Name:    "b",
					Network: &AwsBedrockAgentCoreBrowserNetwork{Mode: "PUBLIC"},
				}
				spec := minimalTools()
				spec.Browsers = []*AwsBedrockAgentCoreBrowser{b, b}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("with a hyphenated browser profile name", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalTools()
				spec.BrowserProfiles = []*AwsBedrockAgentCoreBrowserProfile{{Name: "logged-in"}}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("with a code interpreter in PUBLIC mode carrying vpc_config", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalTools()
				spec.CodeInterpreters[0].Network = &AwsBedrockAgentCoreCodeInterpreterNetwork{
					Mode: "PUBLIC",
					VpcConfig: &AwsBedrockAgentCoreToolVpcConfig{
						Subnets:        []*foreignkeyv1.StringValueOrRef{svr("subnet-abc")},
						SecurityGroups: []*foreignkeyv1.StringValueOrRef{svr("sg-abc")},
					},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("with an enterprise policy missing its S3 object", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalTools()
				spec.Browsers = []*AwsBedrockAgentCoreBrowser{
					{
						Name:    "b",
						Network: &AwsBedrockAgentCoreBrowserNetwork{Mode: "PUBLIC"},
						EnterprisePolicies: []*AwsBedrockAgentCoreBrowserEnterprisePolicy{
							{Type: "MANAGED"},
						},
					},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("with a certificate missing its secret", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalTools()
				spec.CodeInterpreters[0].Certificates = []*AwsBedrockAgentCoreToolCertificate{{}}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})
	})
})

func boolPtr(v bool) *bool {
	return &v
}
