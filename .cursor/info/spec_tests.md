# Spec Tests Authoring Guide

Purpose: create Go tests that exercise `spec.proto` validations (field options and CEL) for a resource kind.

## Scope
- File: `apis/project/planton/provider/<provider>/<kindfolder>/v1/spec_test.go`
- Package: `package <kindfolder>v1`
- Use Ginkgo v2 and Gomega with protovalidate-go.

## Imports
- `buf.build/go/protovalidate`
- `. "github.com/onsi/ginkgo/v2"`
- `. "github.com/onsi/gomega"`
- Optionally: `github.com/plantonhq/openmcf/apis/org/openmcf/shared/validateutil`

## Test Coverage
- Positive: minimal valid `<Kind>Spec` passes.
- Negative: key validations fail as expected:
  - Required strings (min_len)
  - Repeated min_items and uniqueness
  - Enum defined_only and non-zero enforcement
  - CEL conditionals (required-when, ordering constraints, mutually exclusive fields)

## Pattern
- Suite: `Test<Kind>Spec(t *testing.T)` with `RegisterFailHandler(Fail)` and `RunSpecs`.
- `Describe`/`Context`/`It` blocks for scenarios.
- Validate via `protovalidate.Validate(spec)`; assert nil/non-nil.
- For detailed assertions, consider `protovalidate.New().Validate()` and matching violations.

## Example Skeleton
```go
package awscloudfrontv1

import (
  "testing"
  "buf.build/go/protovalidate"
  . "github.com/onsi/ginkgo/v2"
  . "github.com/onsi/gomega"
)

func TestAwsCloudFrontSpec(t *testing.T) {
  RegisterFailHandler(Fail)
  RunSpecs(t, "AwsCloudFrontSpec Validation Suite")
}

var _ = Describe("AwsCloudFrontSpec validations", func() {
  ginkgo.It("accepts a valid spec", func() {
    // construct minimal valid spec then
    Expect(protovalidate.Validate(&AwsCloudFrontSpec{})).To(BeNil())
  })
})
```

## Testing Optional Fields with Default Options

When proto fields are marked as `optional` (required for fields with default options), the generated Go code uses pointer types (`*string`, `*int32`, etc.). Tests must assign values via pointers.

### Problem

```go
// WRONG: Cannot assign string literal to *string field
spec.RunnerGroup = "production"  // Compile error!
```

### Solution

```go
// CORRECT: Use pointer variable
runnerGroup := "production"
spec.RunnerGroup = &runnerGroup

// Or for inline assignment
spec.Image = &KubernetesGhaRunnerScaleSetRunnerImage{
    Repository: stringPtr("my-registry.com/runner"),
    Tag:        stringPtr("v1.0.0"),
}
```

### Helper Pattern

Consider defining a helper in your test file:

```go
func stringPtr(s string) *string {
    return &s
}

func int32Ptr(i int32) *int32 {
    return &i
}
```

### Example Test with Optional Fields

```go
ginkgo.Context("with runner group", func() {
    ginkgo.It("should not return a validation error", func() {
        runnerGroup := "production"
        spec.RunnerGroup = &runnerGroup
        err := protovalidate.Validate(spec)
        gomega.Expect(err).To(gomega.BeNil())
    })
})

ginkgo.Context("with custom runner image", func() {
    ginkgo.It("should not return a validation error", func() {
        repository := "my-registry.com/custom-runner"
        tag := "v1.0.0"
        pullPolicy := "IfNotPresent"
        spec.Runner = &KubernetesGhaRunnerScaleSetRunner{
            Image: &KubernetesGhaRunnerScaleSetRunnerImage{
                Repository: &repository,
                Tag:        &tag,
                PullPolicy: &pullPolicy,
            },
        }
        err := protovalidate.Validate(spec)
        gomega.Expect(err).To(gomega.BeNil())
    })
})
```

### Why This Matters

Fields with `(org.openmcf.shared.options.default)` must be `optional`. This changes the generated Go type from `string` to `*string`. Tests that directly assign string literals will fail to compile after proto changes.

## Notes
- Keep tests focused and robust; avoid brittle provider-format regexes unless necessary.
