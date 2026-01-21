# Revert Auto-Initialization of Unset Nested Messages

**Date**: January 21, 2026
**Type**: Bug Fix
**Components**: Manifest Processing, Proto Defaults

## Summary

Reverted the auto-initialization behavior for unset nested messages that was introduced in the previous change. Unset nested messages now remain unset to preserve user intent, fixing regressions in components like AWS ECS Service where optional feature messages were being inadvertently enabled.

## Problem Statement

The previous implementation auto-initialized unset nested messages when they contained fields with default values. This caused semantic issues:

### The Regression

Using AWS ECS Service as an example:

```protobuf
message AwsEcsServiceSpec {
  // autoscaling is OPTIONAL - user may not want it
  AwsEcsServiceAutoscaling autoscaling = 7;
}

message AwsEcsServiceAutoscaling {
  bool enabled = 1;
  optional int32 target_cpu_percent = 4 [(default) = "75"];  // Has a default!
}
```

**Before this fix:**
1. User provides manifest WITHOUT `autoscaling` (they don't want autoscaling)
2. System sees `target_cpu_percent` has a default
3. System auto-creates `autoscaling` message and sets `target_cpu_percent = 75`
4. IaC module sees `autoscaling` is set, may attempt to configure autoscaling
5. Unexpected infrastructure changes occur

**The core issue:** Auto-initialization conflated two different user intents:
- "I want this feature with defaults" → User sets empty message `autoscaling: {}`
- "I don't want this feature" → User omits the message entirely

## Solution

Reverted to the original behavior: **only apply defaults to messages that are explicitly set by the user**.

```go
// In applyDefaultsToMessage
if field.Kind() == protoreflect.MessageKind {
    if msgReflect.Has(field) {
        // If the field is set, recurse into it to apply defaults
        nestedMsg := msgReflect.Get(field).Message()
        applyDefaultsToMessage(nestedMsg)
    }
    // If NOT set, leave it unset - respect user's choice
    continue
}
```

## Semantic Rules (Final Design)

| User Action | System Behavior |
|-------------|-----------------|
| Omits optional message | Message remains `nil` - feature disabled |
| Sets empty message `{}` | Message created, defaults applied - feature enabled with defaults |
| Sets partial message | User values preserved, defaults applied to unset fields |

## Files Changed

| File | Change |
|------|--------|
| `internal/manifest/protodefaults/applier.go` | Removed `hasFieldsWithDefaults` function and auto-init logic |
| `internal/manifest/protodefaults/applier_test.go` | Updated tests to verify unset messages remain unset |

## Impact

### Fixed Regressions

- AWS ECS Service `autoscaling` no longer auto-enables
- AWS ECS Service `alb` no longer auto-enables
- Any optional feature message with nested defaults now respects user intent

### User Experience

Users who want defaults on optional features simply set the message to empty:

```yaml
# I want autoscaling with defaults
spec:
  autoscaling: {}  # Triggers defaults

# After loading:
spec:
  autoscaling:
    target_cpu_percent: 75  # Default applied
```

```yaml
# I don't want autoscaling
spec:
  # autoscaling omitted - remains nil

# After loading:
spec:
  # autoscaling still omitted - no side effects
```

## Related Work

- `2026-01-14-154952-protodefaults-unset-nested-messages-and-load-alias.md` - Previous change that introduced the regression
- `2025-10-18-01.proto-field-defaults-support.md` - Original proto defaults implementation

---

**Status**: ✅ Production Ready
**Timeline**: ~15 minutes
