package awsmanagedprefixlistv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

func TestAwsManagedPrefixListSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsManagedPrefixListSpec Validation Suite")
}

func minimalList() *AwsManagedPrefixListSpec {
	return &AwsManagedPrefixListSpec{
		Region:        "us-east-1",
		AddressFamily: "IPv4",
		MaxEntries:    10,
		Entries: []*AwsManagedPrefixListEntry{
			{Cidr: "10.20.0.0/16", Description: "office network"},
		},
	}
}

var _ = ginkgo.Describe("AwsManagedPrefixListSpec validations", func() {

	ginkgo.Describe("When valid input is passed", func() {

		ginkgo.It("accepts the minimal IPv4 list", func() {
			gomega.Expect(protovalidate.Validate(minimalList())).To(gomega.BeNil())
		})

		ginkgo.It("accepts an IPv6 list with IPv6 entries", func() {
			spec := minimalList()
			spec.AddressFamily = "IPv6"
			spec.Entries = []*AwsManagedPrefixListEntry{
				{Cidr: "2001:db8::/32"},
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts an empty entry set (capacity reserved, entries later)", func() {
			spec := minimalList()
			spec.Entries = nil
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts entries exactly at max_entries", func() {
			spec := minimalList()
			spec.MaxEntries = 2
			spec.Entries = []*AwsManagedPrefixListEntry{
				{Cidr: "10.20.0.0/16"},
				{Cidr: "10.30.0.0/16"},
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})
	})

	ginkgo.Describe("When invalid input is passed", func() {

		ginkgo.It("rejects an unknown address family", func() {
			spec := minimalList()
			spec.AddressFamily = "DualStack"
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects more entries than max_entries", func() {
			spec := minimalList()
			spec.MaxEntries = 1
			spec.Entries = []*AwsManagedPrefixListEntry{
				{Cidr: "10.20.0.0/16"},
				{Cidr: "10.30.0.0/16"},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an IPv6 entry in an IPv4 list", func() {
			spec := minimalList()
			spec.Entries = append(spec.Entries, &AwsManagedPrefixListEntry{Cidr: "2001:db8::/32"})
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an IPv4 entry in an IPv6 list", func() {
			spec := minimalList()
			spec.AddressFamily = "IPv6"
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects duplicate cidr entries", func() {
			spec := minimalList()
			spec.Entries = append(spec.Entries, &AwsManagedPrefixListEntry{Cidr: "10.20.0.0/16"})
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a cidr without a prefix length", func() {
			spec := minimalList()
			spec.Entries[0].Cidr = "10.20.0.0"
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects zero max_entries", func() {
			spec := minimalList()
			spec.MaxEntries = 0
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})
	})
})
