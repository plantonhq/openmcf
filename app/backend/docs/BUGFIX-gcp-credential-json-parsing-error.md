# Bugfix: GCP Credential JSON Parsing Error

## Issue

Deployments of GCP cloud resources were failing with the following error:

```
error: an unhandled error occurred: program failed:
1 error occurred:
	* failed to load stack-input: failed to load json into proto message:
	  proto: syntax error (line 1:3754): unexpected token "openmcf-testing"
```

The error occurred at character position 3754 in the JSON, where an unquoted string "openmcf-testing" (the GCP project ID from the service account key) was encountered.

## Root Cause

The issue was caused by YAML line folding of long base64-encoded strings during the credential marshalling process:

1. **GCP Service Account Keys are Long**: A base64-encoded GCP service account key JSON is typically 3000-5000+ characters
2. **YAML Line Folding**: When `yaml.v3.Marshal()` encountered such a long string, it automatically folded it across multiple lines according to YAML spec
3. **JSON Conversion Failure**: When this folded YAML was converted to JSON using `sigs.k8s.io/yaml.YAMLToJSON()`, the line folding wasn't properly converted, resulting in malformed JSON with unquoted strings

### Flow Where Error Occurred

```
1. Backend receives GCP credential (base64-encoded service account key)
2. user_provider.go creates temp credential file:
   - Uses yaml.v3.Marshal() to marshal the credentials
   - Long base64 string gets folded across multiple lines
3. gcp_provider.go reads the file and adds it to stack input map
4. stack_input.go builds final stack input YAML
5. Pulumi module reads and converts YAML → JSON
6. JSON parsing fails due to malformed structure from line folding
```

### Example of the Problem

**Before Fix** - YAML output with line folding (simplified):
```yaml
serviceAccountKeyBase64: eyJhdXRoX3Byb3ZpZGVyX3g1MDlfY2VydF91cmwiOiJodHRwczov
  L3d3dy5nb29nbGVhcGlzLmNvbS9vYXV0aDIvdjEvY2VydHMiLCJhdXRoX3VyaSI6Imh0
  dHBzOi8vYWNjb3VudHMuZ29vZ2xlLmNvbS9vL29hdXRoMi9hdXRoIiwiY2xpZW50X2Vt
  ...
```

When converted to JSON, this becomes:
```json
{
  "serviceAccountKeyBase64": "eyJhdXRoX3Byb3ZpZGVyX3g1MDlfY2VydF91cmwiOiJodHRwczov
  L3d3dy5nb29nbGVhcGlzLmNvbS9vYXV0aDIvdjEvY2VydHMiLCJhdXRoX3VyaSI6Imh0..."
}
```

The newlines in the string value cause JSON parsing to fail.

**After Fix** - YAML output with double-quoted style (no folding):
```yaml
serviceAccountKeyBase64: "eyJhdXRoX3Byb3ZpZGVyX3g1MDlfY2VydF91cmwiOiJodHRwczovL3d3dy5nb29nbGVhcGlzLmNvbS9vYXV0aDIvdjEvY2VydHMiLCJhdXRoX3VyaSI6Imh0dHBzOi8vYWNjb3VudHMuZ29vZ2xlLmNvbS9vL29hdXRoMi9hdXRoIiwiY2xpZW50X2VtYWlsIjoidGVzdEBwcm9qZWN0LXBsYW50b24tdGVzdGluZy5pYW0uZ3NlcnZpY2VhY2NvdW50LmNvbSIsImNsaWVudF9pZCI6IjEyMzQ1Njc4OTAxMjM0NTY3ODkwMSIsImNsaWVudF94NTA5X2NlcnRfdXJsIjoiaHR0cHM6Ly93d3cuZ29vZ2xlYXBpcy5jb20vcm9ib3QvdjEvbWV0YWRhdGEveDUwOS90ZXN0JTQwcHJvamVjdC1wbGFudG9uLXRlc3RpbmcuaWFtLmdzZXJ2aWNlYWNjb3VudC5jb20iLCJwcml2YXRlX2tleSI6Ii0tLS0tQkVHSU4gUFJJVkFURSBLRVktLS0tLVxuTUlJRXZRSUJBREFOQmdrcWhraUc5dzBCQVFFRkFBU0NCS2N3Z2dTakFnRUFBb0lCQVFDN1ZKVFV0OVVzOGNLalxuLS0tLS1FTkQgUFJJVkFURSBLRVktLS0tLVxuIiwicHJpdmF0ZV9rZXlfaWQiOiJhYmMxMjNkZWY0NTZnaGk3ODlqa2wwMTJtbm8zNDVwcXI2NzhzdHU5MDF2d3gyMzR5eiIsInByb2plY3RfaWQiOiJwcm9qZWN0LXBsYW50b24tdGVzdGluZyIsInRva2VuX3VyaSI6Imh0dHBzOi8vb2F1dGgyLmdvb2dsZWFwaXMuY29tL3Rva2VuIiwidHlwZSI6InNlcnZpY2VfYWNjb3VudCJ9"
```

This converts cleanly to valid JSON:
```json
{
  "serviceAccountKeyBase64": "eyJhdXRoX3Byb3ZpZGVyX3g1MDlfY2VydF91cmwiOiJodHRwczov..."
}
```

## Solution

Modified `createGcpProviderConfigFileFromProto()` in `pkg/iac/stackinput/stackinputproviderconfig/user_provider.go` to use YAML's `DoubleQuotedStyle`, which prevents line folding for long strings.

### Code Changes

**Before:**
```go
func createGcpProviderConfigFileFromProto(gcpConfig *gcpv1.GcpProviderConfig) (string, func(), error) {
    // ... setup code ...

    gcpCredMap := map[string]interface{}{
        "serviceAccountKeyBase64": gcpConfig.ServiceAccountKeyBase64,
    }

    yamlBytes, err := yaml.Marshal(gcpCredMap)
    // ... write to file ...
}
```

**After:**
```go
func createGcpProviderConfigFileFromProto(gcpConfig *gcpv1.GcpProviderConfig) (string, func(), error) {
    // ... setup code ...

    // Use custom encoder to prevent line folding for long base64 strings
    encoder := yaml.NewEncoder(tmpFile)
    encoder.SetIndent(2)

    // Create a yaml.Node with DoubleQuotedStyle to prevent line folding
    node := &yaml.Node{
        Kind: yaml.MappingNode,
        Content: []*yaml.Node{
            {
                Kind:  yaml.ScalarNode,
                Value: "serviceAccountKeyBase64",
            },
            {
                Kind:  yaml.ScalarNode,
                Style: yaml.DoubleQuotedStyle, // Force double-quoted style
                Value: gcpConfig.ServiceAccountKeyBase64,
            },
        },
    }

    if err := encoder.Encode(node); err != nil {
        // ... error handling ...
    }
    // ... cleanup ...
}
```

### Key Changes

1. **Replaced `yaml.Marshal()`** with `yaml.NewEncoder()` for fine-grained control
2. **Used `yaml.Node` structure** instead of `map[string]interface{}`
3. **Set `Style: yaml.DoubleQuotedStyle`** on the value node to force double-quoting
4. **Double-quoted strings don't get folded** by YAML encoder, ensuring single-line output

## Testing

Added comprehensive tests in `user_provider_test.go`:

1. **`TestCreateGcpProviderConfigFileFromProto`**: Tests with realistic GCP service account key (~944 chars)
2. **`TestCreateGcpProviderConfigWithVeryLongKey`**: Tests with very long key (5232 chars) to simulate real-world scenarios

Both tests verify:
- YAML can be marshalled without errors
- YAML can be unmarshalled back to Go map
- YAML → JSON conversion works correctly
- JSON can be parsed without "unexpected token" errors
- Base64 string integrity is preserved

### Test Results

```
=== RUN   TestCreateGcpProviderConfigFileFromProto
    user_provider_test.go:97: Test passed! Base64 key length: 944 characters
--- PASS: TestCreateGcpProviderConfigFileFromProto (0.00s)

=== RUN   TestCreateGcpProviderConfigWithVeryLongKey
    user_provider_test.go:129: Testing with base64 key of length: 5232 characters
    user_provider_test.go:173: Test passed! Successfully handled base64 key with 5232 characters
--- PASS: TestCreateGcpProviderConfigWithVeryLongKey (0.00s)
```

## Impact

- **Fixes**: GCP cloud resource deployments that were failing with "unexpected token" errors
- **Improves**: Reliability of credential handling for any provider with long credential strings
- **No Breaking Changes**: Output format remains valid YAML, just with different quoting style
- **Performance**: Minimal impact - encoding is slightly more explicit but negligible overhead

## Related Files

- **Fixed**: `pkg/iac/stackinput/stackinputproviderconfig/user_provider.go`
- **Tests**: `pkg/iac/stackinput/stackinputproviderconfig/user_provider_test.go`
- **Affected Flow**: All GCP cloud resource deployments using the backend API

## Future Considerations

1. **Other Providers**: AWS, Azure, and other providers may have similar issues with long credential strings, but less likely since:
   - AWS credentials are typically shorter (access keys ~128 chars)
   - Azure credentials are moderate length
   - GCP's base64-encoded JSON service account keys are uniquely long (3000-5000+ chars)

2. **Consistent YAML Package Usage**: Currently using both `gopkg.in/yaml.v3` and `sigs.k8s.io/yaml` in different files. Consider standardizing on one package.

3. **Alternative Approaches Considered**:
   - Using literal style (`|` or `>`) - rejected because harder to parse back
   - Base64 string splitting - rejected as it changes the data structure
   - Different YAML encoder options - rejected as DoubleQuotedStyle is cleanest

## Verification

To verify the fix works:

1. Create a GCP credential in the web console or via CLI
2. Deploy a GCP cloud resource (e.g., GcpCloudSql)
3. Confirm deployment succeeds without JSON parsing errors
4. Check stack update logs for successful credential resolution

The fix ensures that long base64-encoded credentials are properly handled throughout the deployment pipeline, from storage through YAML marshalling to JSON conversion and finally protobuf unmarshalling in the Pulumi Go module.

