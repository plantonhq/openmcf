package awscognitoresourceserverv1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/planton/apis/dev/planton/shared/foreignkey/v1"
)

func TestAwsCognitoResourceServerSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsCognitoResourceServerSpec Validation Suite")
}

// helper to create a StringValueOrRef with a literal value.
func strRef(val string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: val},
	}
}

var _ = ginkgo.Describe("AwsCognitoResourceServerSpec validations", func() {
	var spec *AwsCognitoResourceServerSpec

	ginkgo.BeforeEach(func() {
		// Minimal valid spec: region + pool + identifier + name.
		spec = &AwsCognitoResourceServerSpec{
			Region:     "us-west-2",
			UserPoolId: strRef("us-east-1_Ab1Cd2EfG"),
			Identifier: "https://api.example.com",
			Name:       "orders-api",
		}
	})

	// -------------------------------------------------------------------------
	// Happy path
	// -------------------------------------------------------------------------

	ginkgo.It("accepts a minimal spec without scopes", func() {
		gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
	})

	ginkgo.It("accepts custom scopes", func() {
		spec.Scopes = []*AwsCognitoResourceServerScope{
			{ScopeName: "read", ScopeDescription: "Read access to orders"},
			{ScopeName: "orders:write", ScopeDescription: "Write access to orders"},
		}
		gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
	})

	// -------------------------------------------------------------------------
	// Required fields
	// -------------------------------------------------------------------------

	ginkgo.It("fails when region is empty", func() {
		spec.Region = ""
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when user_pool_id is missing", func() {
		spec.UserPoolId = nil
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when identifier is empty", func() {
		spec.Identifier = ""
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when name is empty", func() {
		spec.Name = ""
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	// -------------------------------------------------------------------------
	// CEL: scopes
	// -------------------------------------------------------------------------

	ginkgo.It("fails on duplicate scope names", func() {
		spec.Scopes = []*AwsCognitoResourceServerScope{
			{ScopeName: "read", ScopeDescription: "Read access"},
			{ScopeName: "read", ScopeDescription: "Read access again"},
		}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when a scope name contains the reserved separator", func() {
		spec.Scopes = []*AwsCognitoResourceServerScope{
			{ScopeName: "orders/read", ScopeDescription: "Read access"},
		}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when a scope name contains a space", func() {
		spec.Scopes = []*AwsCognitoResourceServerScope{
			{ScopeName: "read orders", ScopeDescription: "Read access"},
		}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when a scope lacks a description", func() {
		spec.Scopes = []*AwsCognitoResourceServerScope{
			{ScopeName: "read"},
		}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})
})
