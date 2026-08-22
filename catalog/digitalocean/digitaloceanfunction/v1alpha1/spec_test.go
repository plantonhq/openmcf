package digitaloceanfunctionv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"github.com/plantonhq/planton/catalog/digitalocean"
)

func TestDigitalOceanFunctionSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "DigitalOceanFunctionSpec Validation Suite")
}

func validGitFunction() *DigitalOceanFunctionSpec {
	return &DigitalOceanFunctionSpec{
		FunctionName: "hello",
		Region:       digitalocean.DigitalOceanRegion_nyc3,
		Git: &digitalocean.DigitalOceanAppGitSource{
			RepoCloneUrl: "https://github.com/digitalocean/sample-functions-nodejs-helloworld.git",
			Branch:       "master",
		},
		SourceDirectory: "packages",
	}
}

var _ = ginkgo.Describe("DigitalOceanFunctionSpec", func() {
	ginkgo.It("accepts a public git source", func() {
		gomega.Expect(protovalidate.Validate(validGitFunction())).To(gomega.BeNil())
	})

	ginkgo.It("accepts a linked github source", func() {
		spec := &DigitalOceanFunctionSpec{
			FunctionName: "hello",
			Region:       digitalocean.DigitalOceanRegion_nyc3,
			Github: &digitalocean.DigitalOceanAppGithubSource{
				Repo:   "myorg/my-functions",
				Branch: "main",
			},
			SourceDirectory: "packages",
		}
		gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
	})

	ginkgo.It("rejects an empty function_name", func() {
		spec := validGitFunction()
		spec.FunctionName = ""
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("rejects a missing region", func() {
		spec := validGitFunction()
		spec.Region = digitalocean.DigitalOceanRegion_digital_ocean_region_unspecified
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("rejects a missing source_directory", func() {
		spec := validGitFunction()
		spec.SourceDirectory = ""
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("rejects two sources", func() {
		spec := validGitFunction()
		spec.Github = &digitalocean.DigitalOceanAppGithubSource{
			Repo:   "myorg/my-functions",
			Branch: "main",
		}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("rejects a spec with no source", func() {
		spec := validGitFunction()
		spec.Git = nil
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("rejects an env var with no value", func() {
		spec := validGitFunction()
		spec.Envs = []*digitalocean.DigitalOceanAppEnvVar{{Key: "EMPTY"}}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("accepts a plaintext env var", func() {
		spec := validGitFunction()
		spec.Envs = []*digitalocean.DigitalOceanAppEnvVar{
			{Key: "NODE_ENV", Value: &digitalocean.DigitalOceanAppEnvVar_Plaintext{Plaintext: "production"}},
		}
		gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
	})
})
