package awsecrregistrysettingsv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/planton/shared/foreignkey/v1"
	"google.golang.org/protobuf/proto"
)

func TestAwsEcrRegistrySettingsSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsEcrRegistrySettingsSpec Validation Suite")
}

func literal(value string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: value},
	}
}

func minimalConfig() *AwsEcrRegistrySettingsSpec {
	return &AwsEcrRegistrySettingsSpec{
		Region: "us-west-2",
		Scanning: &AwsEcrRegistryScanning{
			ScanType: "BASIC",
			Rules: []*AwsEcrRegistryScanningRule{
				{ScanFrequency: "SCAN_ON_PUSH", Filters: []string{"*"}},
			},
		},
	}
}

var _ = ginkgo.Describe("AwsEcrRegistrySettingsSpec validations", func() {

	ginkgo.Describe("When valid input is passed", func() {

		ginkgo.It("accepts a minimal scanning-only configuration", func() {
			gomega.Expect(protovalidate.Validate(minimalConfig())).To(gomega.BeNil())
		})

		ginkgo.It("accepts enhanced continuous scanning", func() {
			spec := minimalConfig()
			spec.Scanning = &AwsEcrRegistryScanning{
				ScanType: "ENHANCED",
				Rules: []*AwsEcrRegistryScanningRule{
					{ScanFrequency: "CONTINUOUS_SCAN", Filters: []string{"prod-*"}},
					{ScanFrequency: "SCAN_ON_PUSH", Filters: []string{"*"}},
				},
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts replication, cache rules, templates, settings, and exclusions", func() {
			spec := minimalConfig()
			spec.ReplicationRules = []*AwsEcrReplicationRule{
				{
					Destinations: []*AwsEcrReplicationDestination{
						{Region: "us-east-1", RegistryId: "111111111111"},
					},
					RepositoryFilters: []string{"prod-"},
				},
			}
			spec.PullThroughCacheRules = []*AwsEcrPullThroughCacheRule{
				{
					EcrRepositoryPrefix: "docker-hub",
					UpstreamRegistryUrl: "registry-1.docker.io",
					CredentialArn:       literal("arn:aws:secretsmanager:us-west-2:111111111111:secret:ecr-pullthroughcache/dh-AbCdEf"),
				},
				{
					EcrRepositoryPrefix: "k8s",
					UpstreamRegistryUrl: "registry.k8s.io",
				},
			}
			spec.RepositoryCreationTemplates = []*AwsEcrRepositoryCreationTemplate{
				{
					Prefix:                             "docker-hub",
					AppliedFor:                         []string{"PULL_THROUGH_CACHE"},
					ImageTagMutability:                 "IMMUTABLE_WITH_EXCLUSION",
					ImageTagMutabilityExclusionFilters: []string{"latest", "dev-*"},
					Encryption: &AwsEcrTemplateEncryption{
						Type:   "KMS",
						KmsKey: literal("arn:aws:kms:us-west-2:111111111111:key/1234abcd-12ab-34cd-56ef-1234567890ab"),
					},
				},
			}
			spec.AccountSettings = &AwsEcrAccountSettings{
				BasicScanTypeVersion: "AWS_NATIVE",
				BlobMounting:         proto.Bool(true),
				RegistryPolicyScope:  "V2",
			}
			spec.PullTimeUpdateExclusions = []*foreignkeyv1.StringValueOrRef{
				literal("arn:aws:iam::111111111111:role/ci-push-role"),
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts a ROOT cache prefix", func() {
			spec := minimalConfig()
			spec.PullThroughCacheRules = []*AwsEcrPullThroughCacheRule{
				{EcrRepositoryPrefix: "ROOT", UpstreamRegistryUrl: "public.ecr.aws"},
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})
	})

	ginkgo.Describe("When invalid input is passed", func() {

		ginkgo.It("rejects an instance managing nothing", func() {
			spec := &AwsEcrRegistrySettingsSpec{Region: "us-west-2"}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects continuous scanning on the basic scanner", func() {
			spec := minimalConfig()
			spec.Scanning.Rules[0].ScanFrequency = "CONTINUOUS_SCAN"
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a scanning rule without filters", func() {
			spec := minimalConfig()
			spec.Scanning.Rules[0].Filters = nil
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an uppercase scanning filter", func() {
			spec := minimalConfig()
			spec.Scanning.Rules[0].Filters = []string{"Prod-*"}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a replication destination with a malformed region", func() {
			spec := minimalConfig()
			spec.ReplicationRules = []*AwsEcrReplicationRule{
				{Destinations: []*AwsEcrReplicationDestination{{Region: "useast1", RegistryId: "111111111111"}}},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a replication destination with a malformed account id", func() {
			spec := minimalConfig()
			spec.ReplicationRules = []*AwsEcrReplicationRule{
				{Destinations: []*AwsEcrReplicationDestination{{Region: "us-east-1", RegistryId: "12345"}}},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a replication rule without destinations", func() {
			spec := minimalConfig()
			spec.ReplicationRules = []*AwsEcrReplicationRule{{}}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects duplicate cache rule prefixes", func() {
			spec := minimalConfig()
			spec.PullThroughCacheRules = []*AwsEcrPullThroughCacheRule{
				{EcrRepositoryPrefix: "docker-hub", UpstreamRegistryUrl: "registry-1.docker.io"},
				{EcrRepositoryPrefix: "docker-hub", UpstreamRegistryUrl: "public.ecr.aws"},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a cache rule with both credential and custom role", func() {
			spec := minimalConfig()
			spec.PullThroughCacheRules = []*AwsEcrPullThroughCacheRule{
				{
					EcrRepositoryPrefix: "docker-hub",
					UpstreamRegistryUrl: "registry-1.docker.io",
					CredentialArn:       literal("arn:aws:secretsmanager:us-west-2:111111111111:secret:x-AbCdEf"),
					CustomRoleArn:       literal("arn:aws:iam::111111111111:role/cache-role"),
				},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a cache prefix with uppercase characters", func() {
			spec := minimalConfig()
			spec.PullThroughCacheRules = []*AwsEcrPullThroughCacheRule{
				{EcrRepositoryPrefix: "DockerHub", UpstreamRegistryUrl: "registry-1.docker.io"},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects duplicate template prefixes", func() {
			spec := minimalConfig()
			spec.RepositoryCreationTemplates = []*AwsEcrRepositoryCreationTemplate{
				{Prefix: "docker-hub", AppliedFor: []string{"PULL_THROUGH_CACHE"}},
				{Prefix: "docker-hub", AppliedFor: []string{"REPLICATION"}},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a template without applied_for", func() {
			spec := minimalConfig()
			spec.RepositoryCreationTemplates = []*AwsEcrRepositoryCreationTemplate{
				{Prefix: "docker-hub"},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects exclusion filters without an exclusion mutability mode", func() {
			spec := minimalConfig()
			spec.RepositoryCreationTemplates = []*AwsEcrRepositoryCreationTemplate{
				{
					Prefix:                             "docker-hub",
					AppliedFor:                         []string{"PULL_THROUGH_CACHE"},
					ImageTagMutability:                 "IMMUTABLE",
					ImageTagMutabilityExclusionFilters: []string{"latest"},
				},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an exclusion filter with three wildcards", func() {
			spec := minimalConfig()
			spec.RepositoryCreationTemplates = []*AwsEcrRepositoryCreationTemplate{
				{
					Prefix:                             "docker-hub",
					AppliedFor:                         []string{"PULL_THROUGH_CACHE"},
					ImageTagMutability:                 "MUTABLE_WITH_EXCLUSION",
					ImageTagMutabilityExclusionFilters: []string{"*a*b*"},
				},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a template KMS key on AES256 encryption", func() {
			spec := minimalConfig()
			spec.RepositoryCreationTemplates = []*AwsEcrRepositoryCreationTemplate{
				{
					Prefix:     "docker-hub",
					AppliedFor: []string{"PULL_THROUGH_CACHE"},
					Encryption: &AwsEcrTemplateEncryption{
						Type:   "AES256",
						KmsKey: literal("arn:aws:kms:us-west-2:111111111111:key/1234abcd-12ab-34cd-56ef-1234567890ab"),
					},
				},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an unknown basic scanner version", func() {
			spec := minimalConfig()
			spec.AccountSettings = &AwsEcrAccountSettings{BasicScanTypeVersion: "TRIVY"}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an unknown registry policy scope", func() {
			spec := minimalConfig()
			spec.AccountSettings = &AwsEcrAccountSettings{RegistryPolicyScope: "V3"}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})
	})
})
