package kuberneteshelmreleasev1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/planton/shared/foreignkey/v1"
)

func TestKubernetesHelmReleaseSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "KubernetesHelmReleaseSpec Validation Suite")
}

func literal(v string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: v},
	}
}

func i32(v int32) *int32 { return &v }

// minimalSpec returns a valid baseline spec; tests mutate one aspect each.
func minimalSpec() *KubernetesHelmReleaseSpec {
	return &KubernetesHelmReleaseSpec{
		Namespace: literal("podinfo"),
		Repo:      "https://stefanprodan.github.io/podinfo",
		Chart:     "podinfo",
		Version:   "6.9.2",
	}
}

var _ = ginkgo.Describe("KubernetesHelmReleaseSpec validations", func() {

	ginkgo.Context("When valid specs are provided", func() {

		ginkgo.It("accepts a minimal chart install", func() {
			gomega.Expect(protovalidate.Validate(minimalSpec())).To(gomega.BeNil())
		})

		ginkgo.It("accepts every CRD lifecycle position, including accepting Helm-managed CRDs", func() {
			spec := minimalSpec()
			install, keep, allow := true, false, true
			spec.Crds = &KubernetesHelmReleaseCrds{Install: &install, KeepOnUninstall: &keep, AllowHelmManaged: &allow}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts an OCI repository", func() {
			spec := minimalSpec()
			spec.Repo = "oci://ghcr.io/stefanprodan/charts"
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts the full values model", func() {
			spec := minimalSpec()
			spec.ValuesYaml = "replicaCount: 2\nresources:\n  requests:\n    cpu: 100m\n"
			spec.Set = map[string]string{"ui.color": "#34577c", "replicaCount": "3"}
			spec.SetString = map[string]string{"image.tag": "6.9.2"}
			spec.SetSensitive = map[string]string{"secret.apiKey": "s3cr3t"}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts private-repo auth set as a pair", func() {
			spec := minimalSpec()
			spec.RepositoryUsername = "robot"
			spec.RepositoryPassword = "token"
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts a release-name override and lifecycle knobs", func() {
			spec := minimalSpec()
			spec.ReleaseName = "podinfo-canary"
			spec.Atomic = true
			spec.CleanupOnFail = true
			spec.WaitForJobs = true
			spec.TimeoutSeconds = i32(600)
			spec.DependencyUpdate = true
			spec.MaxHistory = i32(0)
			spec.TakeOwnership = true
			spec.Description = "canary lane"
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts skip_await without atomic", func() {
			spec := minimalSpec()
			spec.SkipAwait = true
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts a namespace expressed as a foreign-key reference", func() {
			spec := minimalSpec()
			spec.Namespace = &foreignkeyv1.StringValueOrRef{
				LiteralOrRef: &foreignkeyv1.StringValueOrRef_ValueFrom{
					ValueFrom: &foreignkeyv1.ValueFromRef{Name: "apps-namespace"},
				},
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})
	})

	ginkgo.Context("When invalid specs are provided", func() {

		ginkgo.It("rejects a missing namespace", func() {
			spec := minimalSpec()
			spec.Namespace = nil
			gomega.Expect(protovalidate.Validate(spec)).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects a missing repo", func() {
			spec := minimalSpec()
			spec.Repo = ""
			gomega.Expect(protovalidate.Validate(spec)).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects a repo without an http(s):// or oci:// scheme", func() {
			spec := minimalSpec()
			spec.Repo = "stefanprodan.github.io/podinfo"
			gomega.Expect(protovalidate.Validate(spec)).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects a missing chart name", func() {
			spec := minimalSpec()
			spec.Chart = ""
			gomega.Expect(protovalidate.Validate(spec)).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects a missing version (unpinned installs are not reproducible)", func() {
			spec := minimalSpec()
			spec.Version = ""
			gomega.Expect(protovalidate.Validate(spec)).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects an uppercase release name", func() {
			spec := minimalSpec()
			spec.ReleaseName = "Podinfo"
			gomega.Expect(protovalidate.Validate(spec)).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects a release name longer than Helm's 53-character limit", func() {
			spec := minimalSpec()
			spec.ReleaseName = "a-very-long-release-name-that-exceeds-helm-fifty-three"
			gomega.Expect(protovalidate.Validate(spec)).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects atomic combined with skip_await", func() {
			spec := minimalSpec()
			spec.Atomic = true
			spec.SkipAwait = true
			gomega.Expect(protovalidate.Validate(spec)).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects reuse_values combined with reset_values", func() {
			spec := minimalSpec()
			spec.ReuseValues = true
			spec.ResetValues = true
			gomega.Expect(protovalidate.Validate(spec)).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects a repository password without a username", func() {
			spec := minimalSpec()
			spec.RepositoryPassword = "token"
			gomega.Expect(protovalidate.Validate(spec)).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects a repository username without a password", func() {
			spec := minimalSpec()
			spec.RepositoryUsername = "robot"
			gomega.Expect(protovalidate.Validate(spec)).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects a zero timeout", func() {
			spec := minimalSpec()
			spec.TimeoutSeconds = i32(0)
			gomega.Expect(protovalidate.Validate(spec)).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects a negative max_history", func() {
			spec := minimalSpec()
			spec.MaxHistory = i32(-1)
			gomega.Expect(protovalidate.Validate(spec)).ToNot(gomega.BeNil())
		})
	})
})
