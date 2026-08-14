package module

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// runtimeOutputs is what every variant builder hands back to the
// shared export step. The managed flavors fill the key fields with
// empty strings (Azure issues authorization keys only for a primary
// self-hosted registration).
type runtimeOutputs struct {
	id                        pulumi.StringInput
	name                      pulumi.StringInput
	primaryAuthorizationKey   pulumi.StringInput
	secondaryAuthorizationKey pulumi.StringInput
}

// stringPtrWhenSet leaves an optional string argument unsent when the
// spec field is empty, so the provider's own defaults apply (and the
// arguments that reject an explicit empty string never receive one).
func stringPtrWhenSet(value string) pulumi.StringPtrInput {
	if value == "" {
		return nil
	}
	return pulumi.StringPtr(value)
}

// intPtrWhenSet leaves an optional numeric argument unsent when the
// spec field is 0 (unset in proto3), so the provider's own default
// applies.
func intPtrWhenSet(value int32) pulumi.IntPtrInput {
	if value == 0 {
		return nil
	}
	return pulumi.IntPtr(int(value))
}

// boolPtrWhenTrue sends a flag only when it is on -- unset means
// false (the provider's default for the plain boolean arguments).
func boolPtrWhenTrue(value bool) pulumi.BoolPtrInput {
	if !value {
		return nil
	}
	return pulumi.BoolPtr(true)
}
